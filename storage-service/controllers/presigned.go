package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"storage-service/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func GeneratePresignedURL(c *gin.Context) {
	var req struct {
		Filename    string `json:"filename"`
		Size        int64  `json:"size"`
		ContentType string `json:"content_type"`
		Signature   string `json:"signature"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify signature
	if !models.VerifySignature(req.Filename, req.Size, req.ContentType, req.Signature) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// Generate unique ID
	models.ImageIDCounter++
	id := strconv.Itoa(models.ImageIDCounter)

	// Store metadata
	models.Images[id] = &models.ImageMetadata{
		ID:          id,
		Filename:    req.Filename,
		Size:        req.Size,
		ContentType: req.ContentType,
		Signature:   req.Signature,
		Uploaded:    false,
	}

	// Generate pre-signed URL
	cfg, err := config.LoadDefaultConfig(c.Request.Context(),
		config.WithRegion(os.Getenv("AWS_DEFAULT_REGION")),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load AWS config"})
		return
	}

	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)

	key := fmt.Sprintf("uploads/%s/%s", id, req.Filename)
	reqParams := &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("AWS_BUCKET")),
		Key:         aws.String(key),
		ContentType: aws.String(req.ContentType),
	}

	presignedReq, err := presignClient.PresignPutObject(c.Request.Context(), reqParams,
		s3.WithPresignExpires(time.Minute*15),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate pre-signed URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_url": presignedReq.URL,
		"image_id":   id,
	})
}
