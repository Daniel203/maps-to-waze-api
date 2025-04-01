package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps-to-waze-api/internal/database"
	"maps-to-waze-api/models"
	"maps-to-waze-api/services/models"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const MapsPlacesRequestTypeId = 1

func ConvertUrl(ctx context.Context, Url string) (models.ConvertUrlResponse, error) {
	// Step 1: Follow the redirect to get the decompressed google maps Url
	slog.InfoContext(ctx, fmt.Sprintf("Obtaining redirect URL from %s", Url))
	redirectUrl, err := getRedirectUrl(Url)
	if err != nil {
		slog.ErrorContext(ctx, "ConvertUrl failed to get redirect URL", "error", err)
		return models.ConvertUrlResponse{}, fmt.Errorf("ConvertUrl failed to get redirect URL: %w ", err)
	}
	slog.DebugContext(ctx, fmt.Sprintf("Redirect URL: %s", redirectUrl))

	// Step 2: Try to get the coordinates parsing the Url
	slog.InfoContext(ctx, "Trying to get the coordinates from the URL")
	coordinates, err := getCoordinatesFromUrl(redirectUrl)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("ConvertUrl failed to get coordinates from URL: %v ", err))
	}
	if coordinates.Latitude != "" && coordinates.Longitude != "" {
		slog.InfoContext(ctx, fmt.Sprintf("Coordinates found: %v", coordinates))
		url := getWazeLinkFromCoordinates(coordinates)
		return models.ConvertUrlResponse{URL: url, Coordinates: coordinates}, nil
	}

	// Step 3: Try to get the coordinates from the Google Maps API
	slog.InfoContext(ctx, "Trying to get the coordinates from the Google Maps API")
	coordinates, err = getCoordinatesFromApi(ctx, redirectUrl)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("ConvertUrl failed to get coordinates from API: %v ", err))
	}
	if coordinates.Latitude != "" && coordinates.Longitude != "" {
		slog.InfoContext(ctx, fmt.Sprintf("Coordinates found: %v", coordinates))
		url := getWazeLinkFromCoordinates(coordinates)
		return models.ConvertUrlResponse{URL: url, Coordinates: coordinates}, nil
	}

	// Step 4: If no coordinates were found, return an error
	slog.WarnContext(ctx, "No coordinates found")
	return models.ConvertUrlResponse{}, fmt.Errorf("ConvertUrl failed")
}

func getWazeLinkFromCoordinates(coordinates models.Coordinates) string {
	return fmt.Sprintf("https://www.waze.com/ul?ll=%s,%s&navigate=yes", coordinates.Latitude, coordinates.Longitude)
}

func getRedirectUrl(Url string) (string, error) {
	client := &http.Client{}

	resp, err := client.Get(Url)
	if err != nil {
		return "", fmt.Errorf("failed to make the request to %s: %w", Url, err)
	}
	defer resp.Body.Close()

	if resp.Request == nil {
		return "", fmt.Errorf("failed to obtain the redirect URL from response")
	}

	redirectUrl := resp.Request.URL
	decodedUrl, err := url.QueryUnescape(redirectUrl.String())

	if err != nil {
		return "", fmt.Errorf("failed to decode the redirect URL: %w", err)
	}

	return decodedUrl, nil
}

func getCoordinatesFromUrl(Url string) (models.Coordinates, error) {
	var pattern *regexp.Regexp

	// Decide which regex pattern to use based on the Url
	if strings.Contains(Url, "search/") || strings.Contains(Url, "place") {
		pattern = regexp.MustCompile(`([-+]?\d{1,2}\.\d+),\s*([-+]?\d{1,3}\.\d+)`)
	} else {
		pattern = regexp.MustCompile(`@(-?\d+\.\d+),(-?\d+\.\d+)`)
	}

	// Find matching coordinates in the Url
	match := pattern.FindStringSubmatch(Url)
	if match == nil {
		return models.Coordinates{}, fmt.Errorf("failed to match the coordinates in the URL")
	}

	// Extract latitude and longitude
	latitude, longitude := match[1], match[2]

	if latitude == "" || longitude == "" {
		return models.Coordinates{}, fmt.Errorf("failed to get the latitude and longitude from the URL")
	}

	return models.Coordinates{Latitude: latitude, Longitude: longitude}, nil
}

func getCoordinatesFromApi(ctx context.Context, Url string) (models.Coordinates, error) {
	// Check that the number of requests this month is below the limit
	requestLimit, isPresent := os.LookupEnv("MAPS_MAX_REQUESTS_PER_MONTH")
	if !isPresent && requestLimit == "" {
		return models.Coordinates{}, fmt.Errorf("MAPS_MAX_REQUESTS_PER_MONTH environment variable is not set")
	}

	requestLimitInt, err := strconv.Atoi(requestLimit)
	if err != nil {
		return models.Coordinates{}, fmt.Errorf("failed to convert the request limit to int: %w", err)
	}

	canProcede, err := checkNumberOfRequestsThisMonth(ctx, MapsPlacesRequestTypeId, nil, requestLimitInt)
	if err != nil {
		return models.Coordinates{}, fmt.Errorf("failed to check the number of requests this month: %w", err)
	}

	if !canProcede {
		return models.Coordinates{}, fmt.Errorf("exceeded the number of requests this month")
	}

	// Check that the number of requests today is below the limit

	requestLimit, isPresent = os.LookupEnv("MAPS_MAX_REQUESTS_PER_DAY")
	if !isPresent && requestLimit == "" {
		return models.Coordinates{}, fmt.Errorf("MAPS_MAX_REQUESTS_PER_DAY environment variable is not set")
	}

	requestLimitInt, err = strconv.Atoi(requestLimit)
	if err != nil {
		return models.Coordinates{}, fmt.Errorf("failed to convert the request limit to int: %w", err)
	}

	canProcede, err = checkNumberOfRequestsToday(ctx, MapsPlacesRequestTypeId, nil, requestLimitInt)
	if err != nil {
		return models.Coordinates{}, fmt.Errorf("failed to check the number of requests today: %w", err)
	}

	if !canProcede {
		return models.Coordinates{}, fmt.Errorf("exceeded the number of requests today")
	}

	// Get the api key from environment variables
	apiKey := os.Getenv("MAPS_API_KEY")

	// Get the place ID from the Url
	placeID, err := getPlaceIdFromUrl(Url)
	if err != nil || placeID == "" {
		return models.Coordinates{}, fmt.Errorf("failed to extract place ID from URL")
	}

	// Call google maps api and get the coordinates
	baseUrl := "https://maps.googleapis.com/maps/api/place/details/json"
	params := url.Values{
		"cid":    {placeID},
		"key":    {apiKey},
		"fields": {"geometry"},
	}
	apiUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	resp, err := http.Get(apiUrl)
	if err != nil {
		return models.Coordinates{}, fmt.Errorf("failed to make the request to API: %w", err)
	}
	defer resp.Body.Close()

	// Track the request in the database
	requestId := ctx.Value("request_id").(string)
	err = database.InsertRequest(ctx, requestId, MapsPlacesRequestTypeId)
	if err != nil {
		return models.Coordinates{}, fmt.Errorf("failed to insert the request in the database: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return models.Coordinates{}, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Coordinates{}, fmt.Errorf("failed to read the response body: %w", err)
	}

	// Unmarshal the response into the PlaceResponse struct
	var placeResp services_models.GooglePlacesResponse
	if err := json.Unmarshal(body, &placeResp); err != nil {
		return models.Coordinates{}, fmt.Errorf("failed to unmarshal the response: %w", err)
	}

	if placeResp.Status != "OK" {
		return models.Coordinates{}, fmt.Errorf("received non-OK status from API: %s", placeResp.Status)
	}

	// Extract the latitude and longitude from the response
	latitude := fmt.Sprintf("%f", placeResp.Result.Geometry.Location.Lat)
	longitude := fmt.Sprintf("%f", placeResp.Result.Geometry.Location.Lng)

	return models.Coordinates{Latitude: latitude, Longitude: longitude}, nil
}

func getPlaceIdFromUrl(Url string) (string, error) {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`ftid.*:(\w+)`),
		regexp.MustCompile(`data=.*0x(\w+)`),
		regexp.MustCompile(`:0x(\w+)`),
	}

	for _, pattern := range patterns {
		match := pattern.FindStringSubmatch(Url)
		if match != nil && len(match) > 1 {
			placeIdHex := match[1]

			if placeIdHex != "" {
				if strings.HasPrefix(placeIdHex, "0x") {
					placeIdHex = placeIdHex[2:]
				}

				placeIdInt := new(big.Int)
				placeIdInt, success := placeIdInt.SetString(placeIdHex, 16)

				if success {
					return fmt.Sprintf("%d", placeIdInt), nil
				}
			}
		}
	}

	return "", nil
}
