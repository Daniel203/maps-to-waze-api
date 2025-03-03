package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Geometry struct {
	Location struct {
		Lat  float64 `json:"lat"`
		Lng  float64 `json:"lng"`
	} `json:"location"`
}

type PlaceResponse struct {
	Status string    `json:"status"`
	Result struct {
		Geometry Geometry `json:"geometry"`
	} `json:"result"`
}

func ConvertLink(URL string) (string, error) {
	var wazeLink string
	var err error

	// Try converting the link from URL
	wazeLink, err = convertLinkFromURL(URL)
	if err == nil && len(wazeLink) != 0 {
		log.Printf("Successfully converted link from URL: %s", URL)
		return wazeLink, nil
	}

	// Fallback to converting the link via the API
	wazeLink, err = convertLinkFromApi(URL)
	if err == nil && len(wazeLink) != 0 {
		log.Printf("Successfully converted link from API: %s", URL)
		return wazeLink, nil
	}

	// If no conversion method worked, log the error
	log.Printf("Failed to convert the link: %s", URL)
	return "", errors.New("Cannot convert the link")
}

func convertLinkFromURL(URL string) (string, error) {
	// Step 1: Get the redirect URL from the original URL
	redirectURL, err := getRedirectURL(URL)
	if err != nil {
		log.Printf("Error getting redirect URL: %v", err)
		return "", err
	}

	// Step 2: Convert the maps link to a Waze link
	wazeLink, err := convertMapsToWazeLink(redirectURL)
	if err != nil {
		log.Printf("Error converting maps URL to Waze URL: %v", err)
		return "", err
	}

	log.Printf("Redirect URL: %s, Waze link: %s", redirectURL, wazeLink)
	return wazeLink, nil
}

func convertLinkFromApi(URL string) (string, error) {
	redirectURL, err := getRedirectURL(URL)
	if err != nil {
		log.Printf("Error getting redirect URL: %v", err)
		return "", err
	}

	// Load environment variables
	err = godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		return "", fmt.Errorf("Error loading .env file: %v", err)
	}

	// Retrieve the API key from the environment
	mapsApiKey := os.Getenv("MAPS_API_KEY")
	if mapsApiKey == "" {
		log.Printf("MAPS_API_KEY is not set in the environment")
		return "", fmt.Errorf("MAPS_API_KEY is not set in the environment")
	}

	// Get the place ID from the URL
	placeID, err := getPlaceIdFromURL(redirectURL)
	if err != nil {
		log.Printf("Error extracting place ID from URL: %v", err)
		return "", err
	}

    baseUrl := "https://maps.googleapis.com/maps/api/place/details/json"
	params := url.Values{}
	params.Add("cid", placeID)
	params.Add("key", mapsApiKey)

    apiURL := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Error making the request to API: %v", err)
		return "", fmt.Errorf("Error making the request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK status from API: %s", resp.Status)
		return "", fmt.Errorf("Received non-OK HTTP status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading the response body: %v", err)
		return "", fmt.Errorf("Error reading the response body: %v", err)
	}

	// Unmarshal the response into the PlaceResponse struct
	var placeResp PlaceResponse
	if err := json.Unmarshal(body, &placeResp); err != nil {
		log.Printf("Error unmarshaling the response body: %v", err)
		return "", fmt.Errorf("Error unmarshaling the response: %v", err)
	}

    if placeResp.Status != "OK" {
        log.Printf("Received non-OK status from API: %s", placeResp.Status)
        return "", fmt.Errorf("Received non-OK status from API: %s", placeResp.Status)
    }

    // Extract the latitude and longitude from the response
    latitude := fmt.Sprintf("%f", placeResp.Result.Geometry.Location.Lat)
    longitude := fmt.Sprintf("%f", placeResp.Result.Geometry.Location.Lng)

	wazeLink := fmt.Sprintf("https://www.waze.com/ul?ll=%s,%s&navigate=yes", latitude, longitude)
	return wazeLink, nil
}

func getRedirectURL(URL string) (string, error) {
	client := &http.Client{}

	// Perform the HTTP GET request
	resp, err := client.Get(URL)
	if err != nil {
		log.Printf("Failed to make the request to URL: %v", err)
		return "", fmt.Errorf("Failed to make the request: %v", err)
	}
	defer resp.Body.Close()

	// Return the redirect URL
	return resp.Request.URL.String(), nil
}

func convertMapsToWazeLink(URL string) (string, error) {
	var pattern *regexp.Regexp

	// Decide which regex pattern to use based on the URL
	if strings.Contains(URL, "search/") || strings.Contains(URL, "place") {
		pattern = regexp.MustCompile(`([-+]?\d{1,2}\.\d+),\s*([-+]?\d{1,3}\.\d+)`)
	} else {
		pattern = regexp.MustCompile(`@(-?\d+\.\d+),(-?\d+\.\d+)`)
	}

	// Find matching coordinates in the URL
	match := pattern.FindStringSubmatch(URL)
	if match == nil {
		log.Printf("Coordinates not found in URL: %s", URL)
		return "", fmt.Errorf("coordinates not found in URL")
	}

	// Extract latitude and longitude
	latitude, longitude := match[1], match[2]
	wazeLink := fmt.Sprintf("https://www.waze.com/ul?ll=%s,%s&navigate=yes", latitude, longitude)

	log.Printf("Extracted coordinates: %s, %s -> Waze link: %s", latitude, longitude, wazeLink)
	return wazeLink, nil
}

func getPlaceIdFromURL(URL string) (string, error) {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`ftid.*:(\w+)`),
		regexp.MustCompile(`/data=.*0x(\w+)`),
	}

	for _, pattern := range patterns {
		match := pattern.FindStringSubmatch(URL)
		if match != nil {
			placeIdHex := match[1]
			if placeIdHex != "" {
                if strings.HasPrefix(placeIdHex, "0x") {
                    placeIdHex = placeIdHex[2:]
                }

				placeIdInt, err := strconv.ParseInt(placeIdHex, 16, 64)

				if err == nil {
					return fmt.Sprintf("%d", placeIdInt), nil
				}
			}
		}
	}

	return "", errors.New("Place ID not found in URL")
}
