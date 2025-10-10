package book

import (
	"book-store-api/internal/repository"
	"book-store-api/internal/usecase"
	"context"
	"errors"

	"book-store-api/internal/models"
)

func (s *Service) GetAll(ctx context.Context) ([]models.Book, error) {
	books, err := s.repository.GetAll(ctx)
	if err != nil {
		s.logger.Error("db error", "GetAll err", err)
		return nil, usecase.ErrDbInfrastructure
	}
	return books, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*models.Book, error) {
	cached, err := s.cache.Get(ctx, id)
	if err != nil {
		s.logger.Error("cache error", "err", err)
	}

	if cached != nil {
		book, ok := cached.(*models.Book)
		if !ok {
			s.logger.Error("cache: invalid type")
		} else {
			return book, nil
		}
	}

	bookRepo, err := s.repository.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
		s.logger.Error("db error", "getById err", err)
		return nil, usecase.ErrDbInfrastructure
	}

	//Асинхронно добавляем в кэш
	go func() {
		if err := s.cache.Set(ctx, id, bookRepo); err != nil {
			s.logger.Error("cache set error", "err", err)
		}
	}()

	return &bookRepo, nil
}
