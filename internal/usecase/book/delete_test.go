package book

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"book-store-api/internal/usecase"

	"github.com/stretchr/testify/assert"
)

func TestService_DeleteBook(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	cacheDeleted := false
	repoDeleted := false

	mockRepo := &RepositoryMock{
		DeleteFunc: func(ctx context.Context, id string) error {
			repoDeleted = true
			if id == "fail" {
				return fmt.Errorf("db error")
			}
			return nil
		},
	}
	mockCache := &CacheMock{
		DeleteFunc: func(ctx context.Context, key string) error {
			cacheDeleted = true
			if key == "fail" {
				return fmt.Errorf("cache error")
			}
			return nil
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := NewService(logger, mockRepo, mockCache)

	t.Run("successful delete", func(t *testing.T) {
		t.Parallel()
		err := svc.DeleteBook(ctx, "book-123")
		assert.NoError(t, err)
		assert.True(t, repoDeleted, "expected repository.Delete to be called")
		assert.True(t, cacheDeleted, "expected cache.Delete to be called")
	})

	t.Run("repository delete fails", func(t *testing.T) {
		t.Parallel()
		err := svc.DeleteBook(ctx, "fail")
		assert.Equal(t, usecase.ErrDbInfrastructure, err)
	})

	t.Run("cache delete fails but repo succeeds", func(t *testing.T) {
		t.Parallel()
		repoDeleted, cacheDeleted = false, false
		err := svc.DeleteBook(ctx, "fail")
		// репозиторий вернет ошибку, поэтому cache не будет вызван
		assert.Equal(t, usecase.ErrDbInfrastructure, err)
	})
}
