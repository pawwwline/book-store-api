package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"book-store-api/internal/cache"
	"book-store-api/internal/config"
	"book-store-api/internal/delivery/httpv1"
	"book-store-api/internal/infrastructure/db"
	"book-store-api/internal/repository"
	"book-store-api/internal/usecase/book"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	httpServer *http.Server
	db         *pgxpool.Pool
	usecase    *book.Service
	logger     *slog.Logger
}

func BuildApp(cfg *config.Config) (*App, error) {
	logger, err := buildLogger(cfg.Env)
	if err != nil {
		return nil, err
	}

	pool, err := buildDB(&cfg.DB)
	if err != nil {
		return nil, err
	}

	redisCache := buildCache(cfg.Redis)
	repo := buildRepo(pool)

	usecase := buildUseCase(logger, repo, redisCache)
	httpServer := buildHTTP(cfg.HTTP, logger, usecase)

	return &App{
		httpServer: httpServer,
		db:         pool,
		usecase:    usecase,
		logger:     logger,
	}, nil
}

func buildLogger(env string) (*slog.Logger, error) {
	return config.InitLogger(env)
}

func buildDB(cfg *config.DBConfig) (*pgxpool.Pool, error) {
	ctx := context.Background()
	pool, err := db.BuildPoolConn(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func buildRepo(pool *pgxpool.Pool) *repository.BookRepository {

	return repository.NewBookRepository(pool)
}

func buildCache(cfg config.RedisConfig) *cache.Cache {
	return cache.NewCache(cfg)
}

func buildUseCase(logger *slog.Logger, db *repository.BookRepository, cache *cache.Cache) *book.Service {
	return book.NewService(logger, db, cache)
}

func buildHTTP(cfg config.HTTPConfig, logger *slog.Logger, service *book.Service) *http.Server {
	handler := httpv1.NewBookHandler(service, logger)
	return httpv1.InitServer(cfg, logger, handler)
}

func (a *App) Run(ctx context.Context, cacheConfig config.CacheConfig) error {

	go func() {
		if err := a.usecase.LoadCache(ctx, cacheConfig.Limit); err != nil {

			a.logger.Error("cache loading error", "error", err)
		}

		a.logger.Info("orders cache loaded")

	}()

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {

			a.logger.Error("server error", "err", err)
		}
	}()

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	var errList []error

	if err := a.httpServer.Shutdown(ctx); err != nil {
		errList = append(errList, err)
	}
	a.logger.Info("httpv1 server shutdown")

	a.db.Close()
	a.logger.Info("db shutdown")

	if len(errList) > 0 {
		return fmt.Errorf("shutdown errors: %v", errList)
	}

	return nil
}
