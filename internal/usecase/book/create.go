package book

import (
	"context"

	"book-store-api/internal/models"
	"book-store-api/internal/usecase"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, bookInfo models.BookParams) (string, error) {
	uid := uuid.New()
	bookInfo.ID = uid
	book, err := models.NewBook(bookInfo)
	if err != nil {
		return "", err
	}
	err = s.repository.Create(ctx, book)
	if err != nil {
		s.logger.Error("db error", "Create err", err)
		return "", usecase.ErrDbInfrastructure
	}

	err = s.cache.Set(ctx, uid.String(), book)
	if err != nil {
		s.logger.Error("cache error", "err", err)
	}

	return uid.String(), nil

}
