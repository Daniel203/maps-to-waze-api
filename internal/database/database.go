package database

import (
	"database/sql"
	"fmt"
	"os"
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

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	return db, nil
}

