package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type convertLinkRequest struct {
	URL string `json:"url"`
}

func main() {
	router := gin.Default()
	router.POST("convertLink", postConvertLink)

	router.Run("localhost:8080")
}

func postConvertLink(c *gin.Context) {
	var req convertLinkRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redirectURL, err := getRedirectURL(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
	}

	wazeLink, err := convertMapsToWazeLink(redirectURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
	}

	c.String(http.StatusOK, wazeLink)
}

func getRedirectURL(URL string) (string, error) {
	client := &http.Client{}

	resp, err := client.Get(URL)
	if err != nil {
		return "", fmt.Errorf("failed to make the request: %v", err)
	}
	defer resp.Body.Close()

	return resp.Request.URL.String(), nil
}

func convertMapsToWazeLink(URL string) (string, error) {
    var pattern *regexp.Regexp

	if strings.Contains(URL, "search/") || strings.Contains(URL, "place") {
        pattern = regexp.MustCompile(`([-+]?\d{1,2}\.\d+),\s*([-+]?\d{1,3}\.\d+)`)
	} else {
		pattern = regexp.MustCompile(`@(-?\d+\.\d+),(-?\d+\.\d+)`)
	}

	match := pattern.FindStringSubmatch(URL)

	if match == nil {
		return "", fmt.Errorf("coordinates not found in URL")
	}

    latitude, longitude := match[1], match[2]
	wazeLink := fmt.Sprintf("https://www.waze.com/ul?ll=%s,%s&navigate=yes", latitude, longitude)

	return wazeLink, nil
}
