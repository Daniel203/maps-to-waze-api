package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestId := uuid.New().String()

        ctx := r.Context();
        if ctx == nil {
            ctx = context.Background()
        }

		ctx = context.WithValue(r.Context(), "request_id", requestId)

		// Add context to request
		r.WithContext(ctx)

		slog.InfoContext(
			ctx,
			"Request received",
			slog.Group("request",
				"method", r.Method,
				"url", r.URL.Path,
			),
		)

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r.WithContext(ctx))

		slog.InfoContext(
			ctx,
			"Request processed",
			slog.Group("response",
				"status", wrapped.statusCode,
				"duration", time.Since(start).String(),
			),
		)
	})
}
