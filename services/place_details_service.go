package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps-to-waze-api/internal/database"
	"maps-to-waze-api/models"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func (s *Service) GetPlaceDetails(ctx context.Context, latitude float64, longitude float64) (models.PlaceDetailsResponse, error) {
	slog.InfoContext(ctx, fmt.Sprintf("getting place details for coordinates: %f, %f", latitude, longitude))

	if !s.checkNumberRequestsReverseGeocoding(ctx) {
		slog.ErrorContext(ctx, "number of requests exceeded")
		return models.PlaceDetailsResponse{}, fmt.Errorf("number of requests exceeded")
	}

	apiKey, isPresent := os.LookupEnv("GEOAPIFY_API_KEY")
	if !isPresent {
		slog.ErrorContext(ctx, "GEOAPIFY_API_KEY environment variable is not set")
		return models.PlaceDetailsResponse{}, fmt.Errorf("GEOAPIFY_API_KEY environment variable is not set")
	}

	baseUrl := "https://api.geoapify.com/v1/geocode/reverse?"
	params := url.Values{
		"apiKey": {apiKey},
		"lat":    {fmt.Sprintf("%f", latitude)},
		"lon":    {fmt.Sprintf("%f", longitude)},
		"limit":  {"1"},
		"lang":   {"en"},
		"format": {"json"},
	}
	apiUrl := fmt.Sprintf("%s&%s", baseUrl, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		return models.PlaceDetailsResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to make the request: %s", err))
		return models.PlaceDetailsResponse{}, fmt.Errorf("failed to make the request: %w", err)
	}
	defer resp.Body.Close()

	// Track the request in the database
	requestId := ctx.Value("request_id").(string)
	err = database.InsertRequest(ctx, s.DB, requestId, GeoapifyStaticMapRequestTypeId)
	if err != nil {
		return models.PlaceDetailsResponse{}, fmt.Errorf("failed to insert the request in the database: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return models.PlaceDetailsResponse{}, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.PlaceDetailsResponse{}, fmt.Errorf("failed to read the response body: %w", err)
	}

	// Unmarshal the response into the GeoapifyReverseGeocodingResponse struct
	var geoapifyResp models.GeoapifyReverseGeocodingResponse
	if err := json.Unmarshal(body, &geoapifyResp); err != nil {
		return models.PlaceDetailsResponse{}, fmt.Errorf("failed to unmarshal the response: %w", err)
	}

	// Check if the response contains any results
	if len(geoapifyResp.Results) == 0 {
		return models.PlaceDetailsResponse{}, fmt.Errorf("no results found")
	}

	// Extract the relevant information from the response
	placeDetailsRaw := geoapifyResp.Results[0]
	placeDetails := models.PlaceDetailsResponse{
		Formatted: placeDetailsRaw.Formatted,
		AddressLine1: placeDetailsRaw.AddressLine1,
		AddressLine2: placeDetailsRaw.AddressLine2,
	}

	return placeDetails, nil
}

func (s *Service) checkNumberRequestsReverseGeocoding(ctx context.Context) bool {
	slog.InfoContext(ctx, "checking number of requests")

	creditsPerRequestStr, isPresent := os.LookupEnv("GEOAPIFY_CREDIT_PER_REQUEST_REVERSE_GEOCODING")
	if !isPresent && creditsPerRequestStr == "" {
		slog.ErrorContext(ctx, "GEOAPIFY_CREDIT_PER_REQUEST_REVERSE_GEOCODING environment variable is not set")
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

	canProcede, err := s.checkNumberOfRequestsThisMonth(ctx, GeoapifyStaticMapRequestTypeId, &creditsPerRequest, requestLimitInt)

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

	canProcede, err = s.checkNumberOfRequestsToday(ctx, GeoapifyStaticMapRequestTypeId, &creditsPerRequest, requestLimitInt)
	if err != nil || !canProcede {
		return false
	}

	return true
}
