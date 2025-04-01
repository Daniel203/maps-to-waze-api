package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps-to-waze-api/internal/database"
	services_models "maps-to-waze-api/services/models"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const GeoapifyStaticMapRequestTypeId = 2

func GetStaticMap(ctx context.Context, latitude float64, longitude float64) ([]byte, error) {
	slog.InfoContext(ctx, fmt.Sprintf("getting static map for coordinates: %f, %f", latitude, longitude))

	if !CheckNumberRequests(ctx) {
		slog.ErrorContext(ctx, "Number of requests exceeded")
		return nil, fmt.Errorf("number of requests exceeded")
	}

	apiKey := os.Getenv("GEOAPIFY_API_KEY")
	baseUrl := "https://maps.geoapify.com/v1/staticmap"
	params := url.Values{
		"apiKey": {apiKey},
	}
	apiUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	body := services_models.GeoapifyStaticMapRequest{
		Style:       "osm-liberty",
		ScaleFactor: 2,
		Width:       400,
		Height:      200,
		Zoom:        12,
		Center: services_models.Center{
			Lat: latitude,
			Lon: longitude,
		},
		Markers: []services_models.Marker{
			{
				Lat:   latitude,
				Lon:   longitude,
				Color: "#ff3421",
				Size:  "small",
			},
		},
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to marshal the request body: %w", err)
	}

	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		return []byte{}, fmt.Errorf("failed to make the request to API: %w", err)
	}
	defer resp.Body.Close()

	// Track the request in the database
	requestId := ctx.Value("request_id").(string)
	err = database.InsertRequest(ctx, requestId, GeoapifyStaticMapRequestTypeId)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to insert the request in the database: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	// Read the response body
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read the response body: %w", err)
	}

	return resp_body, nil
}

func CheckNumberRequests(ctx context.Context) bool {
	slog.InfoContext(ctx, "checking number of requests")

	creditsPerRequestStr, isPresent := os.LookupEnv("GEOAPIFY_CREDIT_PER_REQUEST_STATIC_MAP")
	if !isPresent && creditsPerRequestStr == "" {
		slog.ErrorContext(ctx, "GEOAPIFY_CREDIT_PER_REQUEST_STATIC_MAP environment variable is not set")
		return false
	}

	creditsPerRequest, err := strconv.ParseFloat(creditsPerRequestStr, 32)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to convert GEOAPIFY_CREDIT_PER_REQUEST_STATIC_MAP to float: %s", err))
		return false
	}

	// Check that the number of requests this month is below the limit
	requestLimit, isPresent := os.LookupEnv("GEOAPIFY_MAX_CREDITS_PER_MONTH")
	if !isPresent && requestLimit == "" {
		slog.ErrorContext(ctx, "GEOAPIFY_MAX_CREDITS_PER_MONTH environment variable is not set")
		return false
	}

	requestLimitInt, err := strconv.Atoi(requestLimit)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to convert the request limit to int: %s", err))
		return false
	}

	canProcede, err := checkNumberOfRequestsThisMonth(ctx, GeoapifyStaticMapRequestTypeId, &creditsPerRequest, requestLimitInt)

	if err != nil || !canProcede {
		return false
	}

	// Check that the number of requests today is below the limit
	requestLimit, isPresent = os.LookupEnv("GEOAPIFY_MAX_CREDITS_PER_DAY")
	if !isPresent && requestLimit == "" {
		slog.ErrorContext(ctx, "GEOAPIFY_MAX_CREDITS_PER_DAY environment variable is not set")
		return false
	}

	requestLimitInt, err = strconv.Atoi(requestLimit)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to convert the request limit to int: %s", err))
		return false
	}

	canProcede, err = checkNumberOfRequestsToday(ctx, GeoapifyStaticMapRequestTypeId, &creditsPerRequest, requestLimitInt)
	if err != nil || !canProcede {
		return false
	}

	return true
}
