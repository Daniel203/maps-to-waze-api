package main

import (
	"fmt"
	"log/slog"
	"maps-to-waze-api/internal/database"
	"maps-to-waze-api/internal/utils"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvironmentVariables()
	utils.InitLogging()

    slog.Info("starting the server")

	database.MigrateDb()
	utils.InitRouter()

    slog.Info("server started successfully")
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
