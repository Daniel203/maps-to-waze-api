package main

import (
	"log/slog"
	"maps-to-waze-api/handlers"
	"maps-to-waze-api/internal/database"
	"maps-to-waze-api/internal/utils"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvironmentVariables()
	utils.InitLogging()

	slog.Info("starting the server")

	db, err := database.InitDb()
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.MigrateDb(db); err != nil {
		slog.Error("database migration failed", "error", err)
		os.Exit(1)
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	app := &handlers.App{
		DB: db,
		HTTPClient: httpClient,
	}

	if err := utils.InitRouter(app); err != nil {
		slog.Error("router failed", "error", err)
		os.Exit(1)
	}
}

func loadEnvironmentVariables() {
	// Try to load a .env file for local development.
	// If it's not found (Docker, Cloud Run, etc.), env vars are
	// expected to be injected by the runtime — that's fine.
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, using variables from runtime environment")
	}
}
