package services

import (
	"database/sql"
	"net/http"
)

type Service struct {
	DB         *sql.DB
	HTTPClient *http.Client
}

func NewService(db *sql.DB, client *http.Client) *Service {
	return &Service{
		DB:         db,
		HTTPClient: client,
	}
}
