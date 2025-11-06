package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func signMetadata(filename string, size int64, contentType string) string {
	secret := os.Getenv("JWT_SECRET")
	data := fmt.Sprintf("%s:%d:%s", filename, size, contentType)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func RequestUploadURL(c *gin.Context) {
	var req struct {
		Filename    string `json:"filename"`
		Size        int64  `json:"size"`
		ContentType string `json:"content_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sign the metadata
	signature := signMetadata(req.Filename, req.Size, req.ContentType)

	// Call storage-service
	storageURL := os.Getenv("STORAGE_SERVICE_URL")
	payload := fmt.Sprintf(`{"filename":"%s","size":%d,"content_type":"%s","signature":"%s"}`, req.Filename, req.Size, req.ContentType, signature)
	resp, err := http.Post(storageURL+"/presigned-url", "application/json", strings.NewReader(payload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pre-signed URL"})
		return
	}
	defer resp.Body.Close()

	var presignedResp struct {
		URL     string `json:"upload_url"`
		ImageID string `json:"image_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&presignedResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"upload_url": presignedResp.URL, "image_id": presignedResp.ImageID})
}
