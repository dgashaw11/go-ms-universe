// Notifications service — consumes product events from Kafka and logs them.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dgashaw11/go-ms-universe/internal/config"
	"github.com/dgashaw11/go-ms-universe/internal/kafka"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(config.LogLevel()),
	}))

	cfg, err := config.LoadNotifications()
	if err != nil {
		logger.Error("loading config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	consumer := kafka.NewConsumer(cfg.KafkaBroker, cfg.KafkaTopic, cfg.ConsumerGroup, logger)
	defer consumer.Close()

	logger.Info("notifications service started",
		"broker", cfg.KafkaBroker,
		"topic", cfg.KafkaTopic,
		"group", cfg.ConsumerGroup,
	)

	consumer.Listen(ctx)

	logger.Info("notifications service stopped")
}
