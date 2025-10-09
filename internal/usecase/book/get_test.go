package book

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"book-store-api/internal/models"
	"book-store-api/internal/usecase"
)

func TestGetAll(t *testing.T) {
	ctx := context.Background()
	mockRepo := &RepositoryMock{}
	cacheMock := &CacheMock{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	svc := NewService(*logger, mockRepo, cacheMock)

	t.Run("success", func(t *testing.T) {
		expectedBooks := []models.Book{
			{ID: uuid.New(), Title: "Book1"},
			{ID: uuid.New(), Title: "Book2"},
		}

		mockRepo.GetAllFunc = func(ctx context.Context) ([]models.Book, error) {
			return expectedBooks, nil
		}

		books, err := svc.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedBooks, books)

		calls := mockRepo.GetAllCalls()
		assert.Len(t, calls, 1)
		assert.Equal(t, ctx, calls[0].Ctx)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.GetAllFunc = func(ctx context.Context) ([]models.Book, error) {
			return nil, errors.New("db error")
		}

		books, err := svc.GetAll(ctx)
		assert.Nil(t, books)
		assert.Equal(t, usecase.ErrDbInfrastructure, err)

		calls := mockRepo.GetAllCalls()
		assert.Len(t, calls, 2)
		assert.Equal(t, ctx, calls[1].Ctx)
	})
}

func TestService_GetByID(t *testing.T) {
	t.Run("found in cache", testGetByIDFoundInCache)
	t.Run("not in cache, found in repo", testGetByIDNotInCacheFoundInRepo)
	t.Run("cache error, fallback to repo", testGetByIDCacheErrorFallback)
	t.Run("invalid type in cache, fallback to repo", testGetByIDInvalidTypeFallback)
	t.Run("repo error", testGetByIDRepoError)
}

func newService(cacheGet func(ctx context.Context, key string) (interface{}, error), cacheSet func(ctx context.Context, key string, value interface{}) error, repoGet func(ctx context.Context, id string) (models.Book, error)) *Service {
	return &Service{
		cache: &CacheMock{
			GetFunc: cacheGet,
			SetFunc: cacheSet,
		},
		repository: &RepositoryMock{
			GetByIdFunc: repoGet,
		},
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func runGetByIDTest(t *testing.T,
	cacheGet func(ctx context.Context, key string) (interface{}, error),
	cacheSet func(ctx context.Context, key string, value interface{}) error,
	repoGet func(ctx context.Context, id string) (models.Book, error),
	wantBook *models.Book,
	wantErr bool,
) {
	svc := newService(cacheGet, cacheSet, repoGet)
	got, err := svc.GetByID(context.Background(), "uuid-123")

	if (err != nil) != wantErr {
		t.Fatalf("GetByID() error = %v, wantErr %v", err, wantErr)
	}

	if !wantErr && wantBook != nil {
		if got == nil {
			t.Fatalf("expected book, got nil")
		}
		if got.ID != wantBook.ID {
			t.Errorf("expected book ID %v, got %v", wantBook.ID, got.ID)
		}
	}

}

func testGetByIDFoundInCache(t *testing.T) {
	validBook := &models.Book{ID: uuid.New(), Title: "Test Book", Author: "Author"}
	runGetByIDTest(t,
		func(ctx context.Context, key string) (interface{}, error) { return validBook, nil },
		func(ctx context.Context, key string, value interface{}) error { return nil },
		func(ctx context.Context, id string) (models.Book, error) { return models.Book{}, nil },
		validBook,
		false,
	)
}

func testGetByIDNotInCacheFoundInRepo(t *testing.T) {
	validBook := &models.Book{ID: uuid.New(), Title: "Test Book", Author: "Author"}
	runGetByIDTest(t,
		func(ctx context.Context, key string) (interface{}, error) { return nil, nil },
		func(ctx context.Context, key string, value interface{}) error { return nil },
		func(ctx context.Context, id string) (models.Book, error) { return *validBook, nil },
		validBook,
		false,
	)
}

func testGetByIDCacheErrorFallback(t *testing.T) {
	validBook := &models.Book{ID: uuid.New(), Title: "Test Book", Author: "Author"}
	runGetByIDTest(t,
		func(ctx context.Context, key string) (interface{}, error) { return nil, fmt.Errorf("cache failure") },
		func(ctx context.Context, key string, value interface{}) error { return nil },
		func(ctx context.Context, id string) (models.Book, error) { return *validBook, nil },
		validBook,
		false,
	)
}

func testGetByIDInvalidTypeFallback(t *testing.T) {
	validBook := &models.Book{ID: uuid.New(), Title: "Test Book", Author: "Author"}
	runGetByIDTest(t,
		func(ctx context.Context, key string) (interface{}, error) { return "string_instead_of_book", nil },
		func(ctx context.Context, key string, value interface{}) error { return nil },
		func(ctx context.Context, id string) (models.Book, error) { return *validBook, nil },
		validBook,
		false,
	)
}

func testGetByIDRepoError(t *testing.T) {
	runGetByIDTest(t,
		func(ctx context.Context, key string) (interface{}, error) { return nil, nil },
		func(ctx context.Context, key string, value interface{}) error { return nil },
		func(ctx context.Context, id string) (models.Book, error) {
			return models.Book{}, fmt.Errorf("db failure")
		},
		nil,
		true,
	)
}
