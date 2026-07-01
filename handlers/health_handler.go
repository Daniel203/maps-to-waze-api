package handlers

import (
	"net/http"
)

func (app *App) GetHealth(w http.ResponseWriter, r *http.Request) {
	// Verify the database connection is alive
	if err := app.Service.DB.PingContext(r.Context()); err != nil {
		http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
