package product_test

import (
	"testing"

	"github.com/dgashaw11/go-ms-universe/internal/product"
	"github.com/dgashaw11/go-ms-universe/internal/product/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_Create(t *testing.T) {
	t.Run("saves product and publishes event", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mock.NewMockRepository(ctrl)
		publisher := mock.NewMockEventPublisher(ctrl)
		svc := product.NewService(repo, publisher)

		repo.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(nil)

		publisher.EXPECT().
			ProductCreated(gomock.Any(), gomock.Any()).
			Return(nil)

		created, err := svc.Create(t.Context(), product.CreateParams{
			Name:        "Widget",
			Description: "A test widget",
			Price:       9.99,
		})

		require.NoError(t, err)
		assert.Equal(t, "Widget", created.Name())
		assert.Equal(t, "A test widget", created.Description())
		assert.InDelta(t, 9.99, created.Price(), 0.001)
		assert.NotEqual(t, uuid.Nil, created.ID())
	})

	t.Run("rejects empty name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mock.NewMockRepository(ctrl)
		publisher := mock.NewMockEventPublisher(ctrl)
		svc := product.NewService(repo, publisher)

		_, err := svc.Create(t.Context(), product.CreateParams{
			Name:  "",
			Price: 5.00,
		})

		require.Error(t, err)
		assert.ErrorContains(t, err, "product name must not be empty")
	})

	t.Run("rejects negative price", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mock.NewMockRepository(ctrl)
		publisher := mock.NewMockEventPublisher(ctrl)
		svc := product.NewService(repo, publisher)

		_, err := svc.Create(t.Context(), product.CreateParams{
			Name:  "Widget",
			Price: -1.00,
		})

		require.Error(t, err)
		assert.ErrorContains(t, err, "product price must not be negative")
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("deletes product and publishes event", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mock.NewMockRepository(ctrl)
		publisher := mock.NewMockEventPublisher(ctrl)
		svc := product.NewService(repo, publisher)

		id := uuid.New()

		repo.EXPECT().
			Delete(gomock.Any(), id).
			Return(nil)

		publisher.EXPECT().
			ProductDeleted(gomock.Any(), id).
			Return(nil)

		err := svc.Delete(t.Context(), id)
		require.NoError(t, err)
	})

	t.Run("returns not found error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mock.NewMockRepository(ctrl)
		publisher := mock.NewMockEventPublisher(ctrl)
		svc := product.NewService(repo, publisher)

		id := uuid.New()

		repo.EXPECT().
			Delete(gomock.Any(), id).
			Return(product.NotFoundError{ID: id})

		err := svc.Delete(t.Context(), id)
		require.Error(t, err)
		assert.ErrorContains(t, err, "not found")
	})
}

func TestService_List(t *testing.T) {
	t.Run("returns paginated products", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := mock.NewMockRepository(ctrl)
		publisher := mock.NewMockEventPublisher(ctrl)
		svc := product.NewService(repo, publisher)

		params := product.ListParams{Page: 1, PageSize: 10}
		expected := product.Page{Total: 0}

		repo.EXPECT().
			List(gomock.Any(), params).
			Return(expected, nil)

		page, err := svc.List(t.Context(), params)

		require.NoError(t, err)
		assert.Equal(t, 0, page.Total)
	})
}
