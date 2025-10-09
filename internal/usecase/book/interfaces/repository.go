package interfaces

import (
	"book-store-api/internal/models"
	"context"
)

type Repository interface {
	Create(ctx context.Context, book models.Book) error
	GetAll(ctx context.Context) ([]models.Book, error)
	GetById(ctx context.Context, id string) (models.Book, error)
	Update(ctx context.Context, book models.Book) error
	Delete(ctx context.Context, id string) error
	GetAllWithLimit(ctx context.Context, limit int) ([]models.Book, error)
}
