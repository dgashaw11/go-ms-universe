package httpapi

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const shutdownTimeout = 5 * time.Second

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func NewServer(addr string, handler http.Handler, logger *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		logger: logger,
	}
}

// Start blocks until the server is shut down.
func (s *Server) Start() {
	s.logger.Info("starting HTTP server", "addr", s.httpServer.Addr)

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("HTTP server failed", "error", err)
	}
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	s.logger.Info("shutting down HTTP server")

	err := s.httpServer.Shutdown(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("shutting down HTTP server: %w", err)
	}

	return nil
}
