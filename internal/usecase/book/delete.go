package book

import (
	"context"
	"errors"

	"book-store-api/internal/repository"
	"book-store-api/internal/usecase"
)

func (s *Service) DeleteBook(ctx context.Context, id string) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return err
		}
		s.logger.Error("db error", "delete book err", err)
		return usecase.ErrDbInfrastructure
	}

	err = s.cache.Delete(ctx, id)
	if err != nil {

		s.logger.Error("cache delete error", "delete book err", err)
	}

	return nil
}
