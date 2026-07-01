package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func InitDb() (*sql.DB, error) {
	dbUrl, isPresent := os.LookupEnv("DATABASE_URL")
	if !isPresent || dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

