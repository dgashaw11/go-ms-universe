package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/dgashaw11/go-ms-universe/internal/metrics"
	"github.com/dgashaw11/go-ms-universe/internal/product"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

//go:generate go tool mockgen -destination=mock/service_mock.go -package=mock . ProductService

type ProductService interface {
	Create(ctx context.Context, params product.CreateParams) (product.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params product.ListParams) (product.Page, error)
}

type createProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type productResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type listResponse struct {
	Items    []productResponse `json:"items"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

func toProductResponse(p product.Product) productResponse {
	return productResponse{
		ID:          p.ID(),
		Name:        p.Name(),
		Description: p.Description(),
		Price:       p.Price(),
		CreatedAt:   p.CreatedAt().UTC().Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt().UTC().Format(time.RFC3339),
	}
}

func CreateProduct(svc ProductService, m *metrics.Metrics, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createProductRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		p, err := svc.Create(r.Context(), product.CreateParams{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
		})
		if err != nil {
			handleServiceError(r.Context(), w, err, logger)
			return
		}

		m.ProductsCreated.Inc()

		writeJSON(w, http.StatusCreated, toProductResponse(p))
	}
}

func DeleteProduct(svc ProductService, m *metrics.Metrics, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid product id")
			return
		}

		if err := svc.Delete(r.Context(), id); err != nil {
			handleServiceError(r.Context(), w, err, logger)
			return
		}

		m.ProductsDeleted.Inc()

		w.WriteHeader(http.StatusNoContent)
	}
}

func ListProducts(svc ProductService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := queryIntOrDefault(r, "page", 1)
		pageSize := min(queryIntOrDefault(r, "page_size", 20), 100)

		result, err := svc.List(r.Context(), product.ListParams{
			Page:     page,
			PageSize: pageSize,
		})
		if err != nil {
			handleServiceError(r.Context(), w, err, logger)
			return
		}

		items := make([]productResponse, 0, len(result.Items))
		for _, p := range result.Items {
			items = append(items, toProductResponse(p))
		}

		writeJSON(w, http.StatusOK, listResponse{
			Items:    items,
			Total:    result.Total,
			Page:     page,
			PageSize: pageSize,
		})
	}
}

func handleServiceError(ctx context.Context, w http.ResponseWriter, err error, fallback *slog.Logger) {
	var validationErr product.ValidationError
	if errors.As(err, &validationErr) {
		writeError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	var notFoundErr product.NotFoundError
	if errors.As(err, &notFoundErr) {
		writeError(w, http.StatusNotFound, notFoundErr.Error())
		return
	}

	LoggerFromContext(ctx, fallback).Error("unexpected error", "error", err)
	writeError(w, http.StatusInternalServerError, "internal server error")
}

func queryIntOrDefault(r *http.Request, key string, fallback int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return fallback
	}

	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return fallback
	}

	return n
}
