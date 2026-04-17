package product

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeCreated EventType = "product.created"
	EventTypeDeleted EventType = "product.deleted"
)

type Event struct {
	Type      EventType `json:"type"`
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name,omitempty"`
	Price     float64   `json:"price,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
