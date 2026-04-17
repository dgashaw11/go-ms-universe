package product

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

//go:generate go tool mockgen -destination=mock/repository_mock.go -package=mock . Repository

type Repository interface {
	Save(ctx context.Context, p Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params ListParams) (Page, error)
}

//go:generate go tool mockgen -destination=mock/event_publisher_mock.go -package=mock . EventPublisher

type EventPublisher interface {
	ProductCreated(ctx context.Context, p Product) error
	ProductDeleted(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo      Repository
	publisher EventPublisher
}

func NewService(repo Repository, publisher EventPublisher) *Service {
	return &Service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *Service) Create(ctx context.Context, params CreateParams) (Product, error) {
	p, err := New(params, time.Now())
	if err != nil {
		return Product{}, fmt.Errorf("creating product: %w", err)
	}

	if err := s.repo.Save(ctx, p); err != nil {
		return Product{}, fmt.Errorf("saving product: %w", err)
	}

	if err := s.publisher.ProductCreated(ctx, p); err != nil {
		return Product{}, fmt.Errorf("publishing product created event: %w", err)
	}

	return p, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting product: %w", err)
	}

	if err := s.publisher.ProductDeleted(ctx, id); err != nil {
		return fmt.Errorf("publishing product deleted event: %w", err)
	}

	return nil
}

func (s *Service) List(ctx context.Context, params ListParams) (Page, error) {
	page, err := s.repo.List(ctx, params)
	if err != nil {
		return Page{}, fmt.Errorf("listing products: %w", err)
	}

	return page, nil
}
