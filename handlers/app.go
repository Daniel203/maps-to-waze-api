package handlers

import (
	"database/sql"
	"net/http"
)

type App struct {
	DB *sql.DB
	HTTPClient *http.Client
}
