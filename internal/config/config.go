package config

import (
	"fmt"
	"os"
	"strconv"
)

type Products struct {
	HTTPAddr    string
	MetricsAddr string
	DatabaseURL string
	KafkaBroker string
	KafkaTopic  string
}

func LoadProducts() (Products, error) {
	cfg := Products{
		HTTPAddr:    envOrDefault("HTTP_ADDR", ":8080"),
		MetricsAddr: envOrDefault("METRICS_ADDR", ":9091"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		KafkaBroker: envOrDefault("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:  envOrDefault("KAFKA_TOPIC", "product-events"),
	}

	if cfg.DatabaseURL == "" {
		return Products{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

type Notifications struct {
	KafkaBroker   string
	KafkaTopic    string
	ConsumerGroup string
}

func LoadNotifications() (Notifications, error) {
	return Notifications{
		KafkaBroker:   envOrDefault("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:    envOrDefault("KAFKA_TOPIC", "product-events"),
		ConsumerGroup: envOrDefault("KAFKA_CONSUMER_GROUP", "notifications"),
	}, nil
}

func MigrationsPath() string {
	return envOrDefault("MIGRATIONS_PATH", "file://migrations")
}

func LogLevel() int {
	level, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return 0 // slog.LevelInfo
	}

	return level
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
