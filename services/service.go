package services

import (
	"database/sql"
	"maps-to-waze-api/models"
	"net/http"
)

type Service struct {
	DB         *sql.DB
	HTTPClient *http.Client
	Config     models.Config
}

func NewService(db *sql.DB, client *http.Client, config models.Config) *Service {
	return &Service{
		DB:         db,
		HTTPClient: client,
		Config:     config,
	}
}
