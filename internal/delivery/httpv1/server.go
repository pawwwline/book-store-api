package http

import (
	"book-store-api/internal/config"

	"log/slog"
	"net/http"
)

func InitServer(cfg *config.HTTPConfig, logger *slog.Logger, bookHandler *Handler) *http.Server {
	router := NewRouter(bookHandler, logger)

	return &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: router,
	}

}
