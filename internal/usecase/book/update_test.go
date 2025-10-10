package book

import (
	"book-store-api/internal/models"
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_Update(t *testing.T) {
	ctx := context.Background()

	validBookParams := models.BookParams{
		ID:          uuid.New(),
		Title:       "Book",
		Description: "desc",
		Author:      "auth",
		ISBN:        "isbn",
		Price:       10,
	}

	t.Run("successful update triggers cache set", func(t *testing.T) {
		cacheSetCalled := false

		mockRepo := &RepositoryMock{
			UpdateFunc: func(ctx context.Context, b models.Book) error {
				return nil
			},
		}
		mockCache := &CacheMock{
			SetFunc: func(ctx context.Context, key string, val interface{}) error {
				cacheSetCalled = true
				return nil
			},
		}

		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		svc := NewService(logger, mockRepo, mockCache)

		err := svc.Update(ctx, validBookParams)
		assert.NoError(t, err)
		assert.True(t, cacheSetCalled, "expected cache.Set to be called")
	})

	t.Run("update fails with invalid book", func(t *testing.T) {
		mockRepo := &RepositoryMock{}
		mockCache := &CacheMock{}
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		svc := NewService(logger, mockRepo, mockCache)

		invalidBook := models.BookParams{Title: ""}
		err := svc.Update(ctx, invalidBook)
		assert.Error(t, err)
	})
}
