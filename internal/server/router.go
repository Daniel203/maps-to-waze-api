package server

import (
	"context"
	"fmt"
	"log/slog"
	"maps-to-waze-api/handlers"
	"maps-to-waze-api/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	router.HandleFunc("GET /health", app.GetHealth)
	router.HandleFunc("POST /convertUrl", app.PostConvertUrl)
	router.HandleFunc("GET /staticMap", app.GetStaticMap)
	router.HandleFunc("GET /placeDetails", app.GetPlaceDetails)

	stack := middleware.CreateStack(middleware.Logging)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: stack(router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second, 
		IdleTimeout:  120 * time.Second,
	}

	// Run the server in a goroutine
	go func() {
		slog.Info("router initialized successfully, listening on " + server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Set up a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block execution until a signal is received
	<-quit
	slog.Info("shutting down server gracefully...")

	// Create a timeout context for the shutdown process
	// This gives active requests 10 seconds to finish before forcing termination
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	slog.Info("server exited cleanly")

	return nil
}

