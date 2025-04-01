package utils

import (
	"fmt"
	"log/slog"
	"maps-to-waze-api/handlers"
	"maps-to-waze-api/middleware"
	"net/http"
	"os"
)

func InitRouter() {
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
	router.HandleFunc("POST /convertUrl", handlers.PostConvertUrl)
	router.HandleFunc("GET /staticMap", handlers.GetStaticMap)

	stack := middleware.CreateStack(middleware.Logging, middleware.Database)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: stack(router),
	}

    slog.Info("Router initialized successfully")
	server.ListenAndServe()
}
