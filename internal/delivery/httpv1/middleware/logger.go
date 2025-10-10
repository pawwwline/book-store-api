package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"time"
)

type contextKey string

const requestIDKey contextKey = "requestID"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggerMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			reqID := r.Context().Value(requestIDKey).(string)

			logger.Info("request started",
				"method", r.Method,
				"path", r.URL.Path,
				"request_id", reqID,
			)

			next.ServeHTTP(w, r)

			logger.Info("request finished",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
				"request_id", reqID,
			)
		})
	}
}
