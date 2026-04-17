package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgashaw11/go-ms-universe/internal/product"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(broker, topic string) *Producer {
	return &Producer{
		writer: &kafkago.Writer{
			Addr:         kafkago.TCP(broker),
			Topic:        topic,
			Balancer:     &kafkago.LeastBytes{},
			RequiredAcks: kafkago.RequireAll,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (p *Producer) ProductCreated(ctx context.Context, prod product.Product) error {
	return p.publish(ctx, product.Event{
		Type:      product.EventTypeCreated,
		ProductID: prod.ID(),
		Name:      prod.Name(),
		Price:     prod.Price(),
		Timestamp: prod.CreatedAt(),
	})
}

func (p *Producer) ProductDeleted(ctx context.Context, id uuid.UUID) error {
	return p.publish(ctx, product.Event{
		Type:      product.EventTypeDeleted,
		ProductID: id,
		Timestamp: time.Now(),
	})
}

func (p *Producer) publish(ctx context.Context, event product.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshaling event: %w", err)
	}

	msg := kafkago.Message{
		Key:   []byte(event.ProductID.String()),
		Value: payload,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("writing message to kafka: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("closing kafka producer: %w", err)
	}

	return nil
}
