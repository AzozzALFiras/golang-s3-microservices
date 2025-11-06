package main

import (
	"log"
	"os"

	"app-service/controllers"
	"app-service/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Public route to get a test JWT token (for testing only)
	r.POST("/auth/token", controllers.GenerateTestToken)

	// JWT middleware for protected routes
	auth := r.Group("/")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("/upload-url", controllers.RequestUploadURL)
		auth.POST("/products", controllers.CreateProduct)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("App Service starting on port %s", port)
	r.Run(":" + port)
}
