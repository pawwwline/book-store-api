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
	ctx := context.Background()

	validBook := &models.Book{
		ID:     uuid.New(),
		Title:  "Test Book",
		Author: "Author",
	}

	tests := []struct {
		name     string
		cacheGet func(ctx context.Context, key string) (interface{}, error)
		cacheSet func(ctx context.Context, key string, value interface{}) error
		repoGet  func(ctx context.Context, id string) (models.Book, error)
		wantBook *models.Book
		wantErr  bool
	}{
		{
			name: "found in cache",
			cacheGet: func(ctx context.Context, key string) (interface{}, error) {
				return validBook, nil
			},
			cacheSet: func(ctx context.Context, key string, value interface{}) error { return nil },
			repoGet:  func(ctx context.Context, id string) (models.Book, error) { return models.Book{}, nil },
			wantBook: validBook,
			wantErr:  false,
		},
		{
			name: "not in cache, found in repo",
			cacheGet: func(ctx context.Context, key string) (interface{}, error) {
				return nil, nil
			},
			cacheSet: func(ctx context.Context, key string, value interface{}) error { return nil },
			repoGet: func(ctx context.Context, id string) (models.Book, error) {
				return *validBook, nil
			},
			wantBook: validBook,
			wantErr:  false,
		},
		{
			name: "cache error, fallback to repo",
			cacheGet: func(ctx context.Context, key string) (interface{}, error) {
				return nil, fmt.Errorf("cache failure")
			},
			cacheSet: func(ctx context.Context, key string, value interface{}) error { return nil },
			repoGet: func(ctx context.Context, id string) (models.Book, error) {
				return *validBook, nil
			},
			wantBook: validBook,
			wantErr:  false,
		},
		{
			name: "invalid type in cache, fallback to repo",
			cacheGet: func(ctx context.Context, key string) (interface{}, error) {
				return "string_instead_of_book", nil
			},
			cacheSet: func(ctx context.Context, key string, value interface{}) error { return nil },
			repoGet: func(ctx context.Context, id string) (models.Book, error) {
				return *validBook, nil
			},
			wantBook: validBook,
			wantErr:  false,
		},
		{
			name:     "repo error",
			cacheGet: func(ctx context.Context, key string) (interface{}, error) { return nil, nil },
			cacheSet: func(ctx context.Context, key string, value interface{}) error { return nil },
			repoGet: func(ctx context.Context, id string) (models.Book, error) {
				return models.Book{}, fmt.Errorf("db failure")
			},
			wantBook: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				cache: &CacheMock{
					GetFunc: tt.cacheGet,
					SetFunc: tt.cacheSet,
				},
				repository: &RepositoryMock{
					GetByIdFunc: tt.repoGet,
				},
				logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
			}

			got, err := svc.GetByID(ctx, "uuid-123")
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.ID != tt.wantBook.ID {
				t.Errorf("expected book ID %v, got %v", tt.wantBook.ID, got.ID)
			}
		})
	}
}
