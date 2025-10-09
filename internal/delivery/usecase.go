package delivery

import (
	"context"

	"book-store-api/internal/models"
)

type Usecase interface {
	Create(ctx context.Context, bookInfo models.BookParams) (string, error)
	DeleteBook(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]models.Book, error)
	Update(ctx context.Context, bookInfo models.BookParams) error
	GetByID(ctx context.Context, id string) (*models.Book, error)
}
