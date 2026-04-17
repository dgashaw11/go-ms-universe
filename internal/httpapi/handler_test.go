package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgashaw11/go-ms-universe/internal/httpapi"
	"github.com/dgashaw11/go-ms-universe/internal/httpapi/mock"
	"github.com/dgashaw11/go-ms-universe/internal/metrics"
	"github.com/dgashaw11/go-ms-universe/internal/product"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	svc := mock.NewMockProductService(ctrl)
	m := metrics.NewNoop()

	body := map[string]any{"name": "Widget", "description": "A widget", "price": 9.99}

	svc.EXPECT().
		Create(gomock.Any(), product.CreateParams{
			Name:        "Widget",
			Description: "A widget",
			Price:       9.99,
		}).
		Return(product.FromStorage(product.StorageParams{
			ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			Name:        "Widget",
			Description: "A widget",
			Price:       9.99,
			CreatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		}), nil)

	rec := httptest.NewRecorder()
	req := jsonRequest(t, http.MethodPost, "/api/v1/products", body)

	handler := httpapi.CreateProduct(svc, m, discardLogger())
	handler(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "Widget", resp["name"])
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", resp["id"])
}

func TestCreateProduct_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	svc := mock.NewMockProductService(ctrl)
	m := metrics.NewNoop()

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/v1/products", bytes.NewBufferString("not json"))

	handler := httpapi.CreateProduct(svc, m, discardLogger())
	handler(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	svc := mock.NewMockProductService(ctrl)
	m := metrics.NewNoop()

	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	svc.EXPECT().
		Delete(gomock.Any(), id).
		Return(nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(
		withChiURLParam(t.Context(), "id", id.String()),
		http.MethodDelete,
		"/api/v1/products/"+id.String(),
		nil,
	)

	handler := httpapi.DeleteProduct(svc, m, discardLogger())
	handler(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestListProducts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	svc := mock.NewMockProductService(ctrl)

	svc.EXPECT().
		List(gomock.Any(), product.ListParams{Page: 1, PageSize: 20}).
		Return(product.Page{
			Items: []product.Product{
				product.FromStorage(product.StorageParams{
					ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					Name:      "Widget",
					Price:     9.99,
					CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				}),
			},
			Total: 1,
		}, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/products", nil)

	handler := httpapi.ListProducts(svc, discardLogger())
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.InDelta(t, float64(1), resp["total"], 0)
	assert.Len(t, resp["items"], 1)
}

func jsonRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()

	var buf bytes.Buffer

	require.NoError(t, json.NewEncoder(&buf).Encode(body))

	return httptest.NewRequestWithContext(t.Context(), method, path, &buf)
}

func withChiURLParam(parent context.Context, key, value string) context.Context {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)

	return context.WithValue(parent, chi.RouteCtxKey, rctx)
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
