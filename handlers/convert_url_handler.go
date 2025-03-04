package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
    "maps-to-waze-api/models"
    "maps-to-waze-api/services"
)

func PostConvertUrl(c *gin.Context) {
	var requestData models.ConvertUrlRequest

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wazeLink, err = services.ConvertUrl(requestData.URL)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, wazeLink)
}
