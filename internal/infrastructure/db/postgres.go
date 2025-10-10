package db

import (
	"context"
	"time"

	"book-store-api/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildPoolConn(ctx context.Context, cfg *config.DBConfig) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, err
	}
	// Настройка соединения пула
	poolCfg.MaxConns = cfg.MaxOpenConns
	poolCfg.MinConns = cfg.MaxIdleConns
	poolCfg.MaxConnIdleTime = cfg.ConnMaxIdleTime * time.Minute
	poolCfg.MaxConnLifetime = cfg.ConnMaxLifetime * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	// Проверка соединения
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
