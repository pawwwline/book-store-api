package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type contextKey string

const requestIDKey contextKey = "requestID"

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

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
			lrw := &loggingResponseWriter{ResponseWriter: w, status: 200}
			start := time.Now()

			reqID, _ := r.Context().Value(requestIDKey).(string)
			if reqID == "" {
				reqID = "unknown"
			}

			logger.Info("request started",
				"method", r.Method,
				"path", r.URL.Path,
				"request_id", reqID,
			)

			next.ServeHTTP(lrw, r)

			logger.Info("request finished",
				"method", r.Method,
				"path", r.URL.Path,
				"status", lrw.status,
				"duration", time.Since(start),
				"request_id", reqID,
			)
		})
	}
}
