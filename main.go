package main

import (
	"github.com/gin-gonic/gin"
	"maps-to-waze-api/handlers"
)

func main() {
	router := gin.Default()
	router.POST("convertUrl", handlers.PostConvertUrl)

	router.Run("localhost:8080")
}
