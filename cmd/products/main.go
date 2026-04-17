// Products service — REST API for product management with Kafka event publishing.
package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dgashaw11/go-ms-universe/internal/config"
	"github.com/dgashaw11/go-ms-universe/internal/httpapi"
	"github.com/dgashaw11/go-ms-universe/internal/kafka"
	"github.com/dgashaw11/go-ms-universe/internal/metrics"
	"github.com/dgashaw11/go-ms-universe/internal/postgres"
	"github.com/dgashaw11/go-ms-universe/internal/product"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

//go:embed swagger.yml
var swaggerSpec []byte

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(config.LogLevel()),
	}))

	cfg, err := config.LoadProducts()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer pool.Close()

	if err := runMigrations(cfg.DatabaseURL); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	logger.Info("migrations applied")

	producer := kafka.NewProducer(cfg.KafkaBroker, cfg.KafkaTopic)
	defer producer.Close()

	repo := postgres.NewProductRepository(pool)
	svc := product.NewService(repo, producer)
	m := metrics.New(prometheus.DefaultRegisterer)

	router := httpapi.NewRouter(svc, m, swaggerSpec, logger)
	server := httpapi.NewServer(cfg.HTTPAddr, router, logger)

	go server.Start()

	logger.Info("products service ready", "addr", cfg.HTTPAddr)
	<-ctx.Done()

	logger.Info("products service shutting down")

	if err := server.Stop(); err != nil {
		return fmt.Errorf("stopping server: %w", err)
	}

	return nil
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New(config.MigrationsPath(), "pgx5://"+databaseURL[len("postgres://"):])
	if err != nil {
		return fmt.Errorf("creating migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("applying migrations: %w", err)
	}

	return nil
}
