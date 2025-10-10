package book

import (
	"book-store-api/internal/repository"
	"context"
	"errors"

	"book-store-api/internal/models"
	"book-store-api/internal/usecase"
)

func (s *Service) Update(ctx context.Context, bookInfo models.BookParams) error {
	book, err := models.NewBook(bookInfo)
	if err != nil {
		return err
	}
	err = s.repository.Update(ctx, book)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return err
		}
		s.logger.Error("db error", "update error", err)
		return usecase.ErrDbInfrastructure
	}

	if err := s.cache.Set(ctx, book.ID.String(), book); err != nil {
		s.logger.Error("cache async set error", "err", err)
	}

	return nil
}
