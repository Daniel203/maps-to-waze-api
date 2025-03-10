package middleware

import (
	"context"
	"log/slog"
	"maps-to-waze-api/internal/database"
	"net/http"
)

func Database(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		db, err := database.InitDb()
		if err != nil {
			slog.ErrorContext(ctx, "Failed to open database connection", "error", err)
			http.Error(w, "Failed to open database connection", http.StatusInternalServerError)
		}

		ctx = context.WithValue(r.Context(), "db", db)

		// Add context to request
		r.WithContext(ctx)

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r.WithContext(ctx))
	})
}
