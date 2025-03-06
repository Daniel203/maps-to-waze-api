package main

import (
	"fmt"
	"maps-to-waze-api/handlers"
	"maps-to-waze-api/middleware"
	"net/http"
	"os"
)

func initRouter() {
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

	stack := middleware.CreateStack(middleware.Logging)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: stack(router),
	}

	server.ListenAndServe()
}
