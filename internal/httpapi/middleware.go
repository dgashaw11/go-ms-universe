package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type ctxKey struct{}

func ContextLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())

			reqLogger := logger.With(
				"request_id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
			)

			ctx := context.WithValue(r.Context(), ctxKey{}, reqLogger)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r.WithContext(ctx))

			reqLogger.Info("request completed",
				"status", ww.Status(),
				"duration", time.Since(start),
			)
		})
	}
}

func LoggerFromContext(ctx context.Context, fallback *slog.Logger) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}

	return fallback
}
