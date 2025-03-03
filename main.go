package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"maps-to-waze-api/lib"
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

    var wazeLink, err = lib.ConvertLink(req.URL);

    if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
    }

    if len(wazeLink) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible to convert the link"})
        return 
    }

	c.String(http.StatusOK, wazeLink)
}
