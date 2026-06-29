package utils

import (
	"fmt"
	"log/slog"
	"maps-to-waze-api/handlers"
	"maps-to-waze-api/middleware"
	"net/http"
	"os"
	"time"
)

func InitRouter(app *handlers.App) error {
    slog.Info("Initializing the router")

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := http.NewServeMux()
	router.HandleFunc("POST /convertUrl", app.PostConvertUrl)
	router.HandleFunc("GET /staticMap", app.GetStaticMap)
	router.HandleFunc("GET /placeDetails", app.GetPlaceDetails)

	stack := middleware.CreateStack(middleware.Logging)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: stack(router),
		// Time from connection accepted to request body fully read
		ReadTimeout:  5 * time.Second,
		// Time from end of request header read to end of response write
		WriteTimeout: 15 * time.Second, 
		// Time to keep idle connections alive (keep-alive)
		IdleTimeout:  120 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

    slog.Info("router initialized successfully")
	return nil
}
