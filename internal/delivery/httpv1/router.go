package httpv1

import (
	"book-store-api/internal/delivery/httpv1/middleware"
	"github.com/gorilla/mux"
	"log/slog"
)

func NewRouter(bookHandler *Handler, logger *slog.Logger) *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.RequestIDMiddleware)
	api.Use(middleware.LoggerMiddleware(logger))
	bookHandler.RegisterRoutes(api)

	return router
}
