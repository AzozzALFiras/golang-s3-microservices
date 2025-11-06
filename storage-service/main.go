package main

import (
	"log"
	"os"

	"storage-service/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/presigned-url", controllers.GeneratePresignedURL)
	r.GET("/verify/:id", controllers.VerifyUpload)

	port := os.Getenv("STORAGE_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Storage Service starting on port %s", port)
	r.Run(":" + port)
}
