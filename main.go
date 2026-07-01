package main

import (
	"fmt"
	"log/slog"
	"maps-to-waze-api/handlers"
	"maps-to-waze-api/internal/database"
	"maps-to-waze-api/internal/logger"
	"maps-to-waze-api/internal/server"
	"maps-to-waze-api/models"
	"maps-to-waze-api/services"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvironmentVariables()

	logger.InitLogging()
	slog.Info("starting the server")

	config, err := loadConfig()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

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

	service := services.NewService(db, httpClient, config)
	app := &handlers.App{
		Service: service,
	}

	if err := server.InitRouter(app); err != nil {
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

func loadConfig() (models.Config, error) {
	mapsMonthLimitStr := os.Getenv("MAPS_MAX_REQUESTS_PER_MONTH")
	if mapsMonthLimitStr == "" {
		return models.Config{}, fmt.Errorf("MAPS_MAX_REQUESTS_PER_MONTH is required")
	}
	mapsMonthLimit, err := strconv.Atoi(mapsMonthLimitStr)
	if err != nil {
		return models.Config{}, fmt.Errorf("MAPS_MAX_REQUESTS_PER_MONTH must be an integer")
	}

	mapsDayLimitStr := os.Getenv("MAPS_MAX_REQUESTS_PER_DAY")
	if mapsDayLimitStr == "" {
		return models.Config{}, fmt.Errorf("MAPS_MAX_REQUESTS_PER_DAY is required")
	}
	mapsDayLimit, err := strconv.Atoi(mapsDayLimitStr)
	if err != nil {
		return models.Config{}, fmt.Errorf("MAPS_MAX_REQUESTS_PER_DAY must be an integer")
	}

	mapsAPIKey := os.Getenv("MAPS_API_KEY")
	if mapsAPIKey == "" {
		return models.Config{}, fmt.Errorf("MAPS_API_KEY is required")
	}

	return models.Config{
		MapsMaxRequestsPerMonth: mapsMonthLimit,
		MapsMaxRequestsPerDay:   mapsDayLimit,
		MapsAPIKey:              mapsAPIKey,
	}, nil
}
