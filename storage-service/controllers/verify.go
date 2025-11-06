package controllers

import (
	"net/http"

	"storage-service/models"

	"github.com/gin-gonic/gin"
)

func VerifyUpload(c *gin.Context) {
	imageID := c.Param("id")
	image, exists := models.Images[imageID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// For demo, assume upload is successful if metadata exists
	// In real scenario, check S3 for object existence
	image.Uploaded = true

	c.JSON(http.StatusOK, gin.H{"valid": true})
}
