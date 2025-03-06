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

    var wazeLink, err = services.ConvertUrl(ctx, requestData.URL)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    slog.DebugContext(ctx, fmt.Sprintf("Waze link: %s", wazeLink))

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(wazeLink)
}
