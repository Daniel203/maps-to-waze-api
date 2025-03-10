package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

const DB_CONTEXT_KEY string = "db"

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

func ContextWithDB(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, DB_CONTEXT_KEY, db)
}

// Retrieve the DB connection from the context
func DBFromContext(ctx context.Context) (*sql.DB, error) {
	db, ok := ctx.Value(DB_CONTEXT_KEY).(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("could not get DB from context")
	}
	return db, nil
}
