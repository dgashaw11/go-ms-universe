package product

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	id          uuid.UUID
	name        string
	description string
	price       float64
	createdAt   time.Time
	updatedAt   time.Time
}

type CreateParams struct {
	Name        string
	Description string
	Price       float64
}

func New(params CreateParams, now time.Time) (Product, error) {
	p := Product{
		id:          uuid.New(),
		name:        params.Name,
		description: params.Description,
		price:       params.Price,
		createdAt:   now,
		updatedAt:   now,
	}

	if err := p.validate(); err != nil {
		return Product{}, err
	}

	return p, nil
}

// StorageParams reconstructs without validation.
type StorageParams struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func FromStorage(p StorageParams) Product {
	return Product{
		id:          p.ID,
		name:        p.Name,
		description: p.Description,
		price:       p.Price,
		createdAt:   p.CreatedAt,
		updatedAt:   p.UpdatedAt,
	}
}

func (p Product) ID() uuid.UUID        { return p.id }
func (p Product) Name() string         { return p.name }
func (p Product) Description() string  { return p.description }
func (p Product) Price() float64       { return p.price }
func (p Product) CreatedAt() time.Time { return p.createdAt }
func (p Product) UpdatedAt() time.Time { return p.updatedAt }

func (p Product) validate() error {
	var msgs []string

	if p.name == "" {
		msgs = append(msgs, "product name must not be empty")
	}

	if p.price < 0 {
		msgs = append(msgs, "product price must not be negative")
	}

	if len(msgs) > 0 {
		return ValidationError{Messages: msgs}
	}

	return nil
}

type ValidationError struct {
	Messages []string
}

func (e ValidationError) Error() string {
	return strings.Join(e.Messages, "; ")
}

type NotFoundError struct {
	ID uuid.UUID
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("product %s not found", e.ID)
}

type Page struct {
	Items []Product
	Total int
}

type ListParams struct {
	Page     int
	PageSize int
}
