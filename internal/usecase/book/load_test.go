package book

import (
	"book-store-api/internal/models"
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_LoadCache(t *testing.T) {
	ctx := context.Background()

	mockRepo := &RepositoryMock{}
	mockCache := &CacheMock{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	svc := NewService(*logger, mockRepo, mockCache)

	t.Run("successfully loads cache", func(t *testing.T) {
		mockRepo.GetAllWithLimitFunc = func(ctx context.Context, limit int) ([]models.Book, error) {
			return []models.Book{
				{ID: uuid.New(), Title: "Book1"},
				{ID: uuid.New(), Title: "Book2"},
			}, nil
		}

		calledSet := []string{}
		mockCache.SetFunc = func(ctx context.Context, key string, value interface{}) error {
			calledSet = append(calledSet, key)
			return nil
		}

		err := svc.LoadCache(ctx, 10)
		assert.NoError(t, err)
		assert.Len(t, calledSet, 2)
	})

	t.Run("repo error returns error", func(t *testing.T) {
		mockRepo.GetAllWithLimitFunc = func(ctx context.Context, limit int) ([]models.Book, error) {
			return nil, fmt.Errorf("db failure")
		}

		err := svc.LoadCache(ctx, 5)
		assert.Error(t, err)
	})

	t.Run("cache errors are logged but not returned", func(t *testing.T) {
		mockRepo.GetAllWithLimitFunc = func(ctx context.Context, limit int) ([]models.Book, error) {
			return []models.Book{{ID: uuid.New(), Title: "Book1"}}, nil
		}

		cacheErr := fmt.Errorf("cache failure")
		cacheCalled := false
		mockCache.SetFunc = func(ctx context.Context, key string, value interface{}) error {
			cacheCalled = true
			return cacheErr
		}

		err := svc.LoadCache(ctx, 1)
		assert.NoError(t, err)
		assert.True(t, cacheCalled)
	})
}
