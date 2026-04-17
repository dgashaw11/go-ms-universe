package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgashaw11/go-ms-universe/internal/product"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) Save(ctx context.Context, p product.Product) error {
	const query = `
		INSERT INTO products (id, name, description, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.pool.Exec(ctx, query,
		p.ID(), p.Name(), p.Description(), p.Price(), p.CreatedAt(), p.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("inserting product: %w", err)
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM products WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return product.NotFoundError{ID: id}
	}

	return nil
}

// List returns products ordered by creation time (newest first).
func (r *ProductRepository) List(ctx context.Context, params product.ListParams) (product.Page, error) {
	const query = `
		SELECT id, name, description, price, created_at, updated_at,
		       count(*) OVER () AS total
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	offset := (params.Page - 1) * params.PageSize

	rows, err := r.pool.Query(ctx, query, params.PageSize, offset)
	if err != nil {
		return product.Page{}, fmt.Errorf("querying products: %w", err)
	}

	var total int

	products, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (product.Product, error) {
		var p product.StorageParams

		err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt, &total)
		if err != nil {
			return product.Product{}, fmt.Errorf("scanning product row: %w", err)
		}

		return product.FromStorage(p), nil
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return product.Page{}, fmt.Errorf("collecting products: %w", err)
	}

	return product.Page{Items: products, Total: total}, nil
}
