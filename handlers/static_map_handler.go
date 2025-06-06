package handlers

import (
	"maps-to-waze-api/services"
	"net/http"
	"strconv"
)

func GetStaticMap(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context();
	latitudeStr := r.URL.Query().Get("lat");
	longitudeStr := r.URL.Query().Get("lon");

	if latitudeStr == "" || longitudeStr == "" {
		http.Error(w, "Missing latitude or longitude", http.StatusBadRequest)
		return
	}

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude format", http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude format", http.StatusBadRequest)
		return
	}

	data, err := services.GetStaticMap(ctx, latitude, longitude);

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "image/jpg")
    w.WriteHeader(http.StatusOK)
	_, err = w.Write(data);

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
