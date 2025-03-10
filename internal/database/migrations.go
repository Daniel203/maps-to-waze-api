package database

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDb() {
    slog.Info("Migrating the database")

	dbUrl, isPresent := os.LookupEnv("DATABASE_URL")
	if !isPresent || dbUrl == "" {
		panic("DATABASE_URL is not set")
	}

	migrationsPath := "file://db/migrations"

	// Create a new migration instance
	m, err := migrate.New(migrationsPath, dbUrl)
	if err != nil {
		panic(fmt.Sprintf("Failed to create migration instance: %v", err))
	}

	// Apply migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
        panic(fmt.Sprintf("Failed to apply migrations: %v", err))
	}

    slog.Info("Database migrated successfully")
}
