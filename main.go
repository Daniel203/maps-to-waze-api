package main

import (
	"fmt"
	"maps-to-waze-api/handlers"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
	router.POST("convertUrl", handlers.PostConvertUrl)

    loadEnvironmentVariables();

    host := os.Getenv("HOST")
    if host == "" {
        host = "localhost"
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    router.Run(host + ":" + port)
}

func loadEnvironmentVariables() {
    env, isPresent := os.LookupEnv("ENV") 

    // If production do nothing
    if isPresent && env == "prod" {
        return
    }

	if err := godotenv.Load(); err != nil {
        panic(fmt.Errorf("Failed to load the environment variables: %w", err))
	}
}
