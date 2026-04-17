package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dgashaw11/go-ms-universe/internal/product"
	kafkago "github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafkago.Reader
	logger *slog.Logger
}

func NewConsumer(broker, topic, group string, logger *slog.Logger) *Consumer {
	return &Consumer{
		reader: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:  []string{broker},
			Topic:    topic,
			GroupID:  group,
			MinBytes: 1,
			MaxBytes: 10e6, // 10 MB
		}),
		logger: logger,
	}
}

// Listen blocks until the context is canceled.
func (c *Consumer) Listen(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			c.logger.Error("reading kafka message", "error", err)

			continue
		}

		var event product.Event

		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.logger.Error("unmarshaling event", "error", err, "raw", string(msg.Value))

			continue
		}

		c.logger.Info("received product event",
			"type", event.Type,
			"product_id", event.ProductID,
			"name", event.Name,
			"price", event.Price,
			"timestamp", event.Timestamp,
		)
	}
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("closing kafka consumer: %w", err)
	}

	return nil
}
