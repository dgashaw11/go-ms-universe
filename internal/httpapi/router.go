package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/dgashaw11/go-ms-universe/internal/metrics"
	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(
	svc ProductService,
	m *metrics.Metrics,
	swaggerSpec []byte,
	logger *slog.Logger,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(ContextLogger(logger))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/swagger/*", http.StripPrefix("/swagger", swaggerui.Handler(swaggerSpec)))

	r.Route("/api/v1/products", func(r chi.Router) {
		r.Post("/", CreateProduct(svc, m, logger))
		r.Delete("/{id}", DeleteProduct(svc, m, logger))
		r.Get("/", ListProducts(svc, logger))
	})

	return r
}
