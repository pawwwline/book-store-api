package httpv1

import (
	"log/slog"
	"net/http"
	"time"

	"book-store-api/internal/config"

	"github.com/gorilla/handlers"
)

func InitServer(cfg config.HTTPConfig, logger *slog.Logger, bookHandler *Handler) *http.Server {
	router := NewRouter(bookHandler, logger)
	corsAllowed := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	return &http.Server{
		Addr:              cfg.Host + ":" + cfg.Port,
		Handler:           corsAllowed(router),
		ReadHeaderTimeout: time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.IdleTimeout) * time.Second,
	}

}
