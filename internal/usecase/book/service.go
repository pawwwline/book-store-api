package book

import (
	"log/slog"

	"book-store-api/internal/usecase/book/interfaces"
)

type Service struct {
	logger     *slog.Logger
	repository interfaces.Repository
	cache      interfaces.Cache
}

func NewService(logger slog.Logger, repo interfaces.Repository, cache interfaces.Cache) *Service {
	return &Service{
		logger:     &logger,
		repository: repo,
		cache:      cache,
	}
}
