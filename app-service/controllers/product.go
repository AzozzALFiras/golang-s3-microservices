package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"app-service/models"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		ImageID     string  `json:"image_id"`
		Price       float64 `json:"price"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate image with storage-service
	storageURL := os.Getenv("STORAGE_SERVICE_URL")
	resp, err := http.Get(fmt.Sprintf("%s/verify/%s", storageURL, req.ImageID))
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}
	defer resp.Body.Close()

	var verifyResp struct {
		Valid bool `json:"valid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil || !verifyResp.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	// Create product
	models.ProductIDCounter++
	id := strconv.Itoa(models.ProductIDCounter)
	product := models.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		ImageID:     req.ImageID,
		Price:       req.Price,
	}
	models.Products[id] = product

	c.JSON(http.StatusCreated, product)
}
