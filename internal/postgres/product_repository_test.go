package postgres_test

import (
	"os"
	"testing"
	"time"

	"github.com/dgashaw11/go-ms-universe/internal/postgres"
	"github.com/dgashaw11/go-ms-universe/internal/product"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductRepository_SaveListDelete_Integration(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set")
	}

	ctx := t.Context()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	t.Cleanup(func() { pool.Close() })

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id          UUID PRIMARY KEY,
			name        TEXT        NOT NULL,
			description TEXT        NOT NULL DEFAULT '',
			price       NUMERIC     NOT NULL CHECK (price >= 0),
			created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
		)`)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM products`) //nolint:errcheck
	})

	repo := postgres.NewProductRepository(pool)
	now := time.Now().Truncate(time.Microsecond)

	p, err := product.New(product.CreateParams{
		Name:        "Integration Widget",
		Description: "Created by integration test",
		Price:       42.50,
	}, now)
	require.NoError(t, err)

	err = repo.Save(ctx, p)
	require.NoError(t, err)

	page, err := repo.List(ctx, product.ListParams{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	assert.Equal(t, 1, page.Total)
	assert.Equal(t, "Integration Widget", page.Items[0].Name())
	assert.InDelta(t, 42.50, page.Items[0].Price(), 0.001)

	err = repo.Delete(ctx, p.ID())
	require.NoError(t, err)

	page, err = repo.List(ctx, product.ListParams{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Empty(t, page.Items)
	assert.Equal(t, 0, page.Total)

	err = repo.Delete(ctx, uuid.New())
	require.Error(t, err)
	assert.ErrorAs(t, err, &product.NotFoundError{})
}
