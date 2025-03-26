package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"maps-to-waze-api/models"
	"maps-to-waze-api/services"
	"net/http"
)

func PostConvertUrl(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context();
    var requestData models.ConvertUrlRequest

    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return 
    }

    var data, err = services.ConvertUrl(ctx, requestData.URL)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    jsonData, err := json.Marshal(data)

    if err != nil {
		slog.ErrorContext(ctx, "Error marshaling JSON:", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

    slog.DebugContext(ctx, fmt.Sprintf("Waze link: %s, Coordinates: %+v", data.URL, data.Coordinates))

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, err = w.Write(jsonData);

    if err != nil {
        slog.ErrorContext(ctx, "Error writing response:", "error", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}
