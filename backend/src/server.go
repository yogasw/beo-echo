package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/git-sync/handler"
	"mockoon-control-panel/backend_new/src/health"
	"mockoon-control-panel/backend_new/src/lib"
	"mockoon-control-panel/backend_new/src/middlewares"
	mockHandler "mockoon-control-panel/backend_new/src/mocks/handler"
	"mockoon-control-panel/backend_new/src/traefik"
	"mockoon-control-panel/backend_new/src/utils"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(filepath.Join("..", ".env")); err != nil {
		log.Println("Warning: .env file not found or could not be loaded")
	}

	// Setup required directories
	if err := utils.EnsureRequiredFoldersAndEnv(); err != nil {
		log.Fatalf("Failed to setup required folders: %v", err)
	}

	// Check if mockoon CLI is available
	mockoonAvailable, err := utils.CheckMockoonCli()
	if err != nil || !mockoonAvailable {
		log.Fatalf("Error: mockoon-cli is not available. Please install it first: %v", err)
	}

	// Sync configs to git
	if err := handler.SyncConfigsToGit(); err != nil {
		log.Printf("Error syncing to Git: %v", err)
	} else {
		log.Println("Sync to Git completed successfully")
	}

	// Setup database connection
	if err := database.CheckAndHandlePrisma(); err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	// Generate Traefik config
	if err := traefik.GenerateDynamicTraefikConfig(); err != nil {
		log.Fatalf("Error generating Traefik config: %v", err)
	}

	if err := traefik.GenerateStaticTraefikConfig(); err != nil {
		log.Fatalf("Error generating static Traefik config: %v", err)
	}

	// Create Gin router with default middleware
	router := gin.Default()

	// Add request logging middleware
	router.Use(func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Log request details
		log.Printf(
			"[%s] %s %s %d %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			time.Since(startTime),
		)
	})

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{lib.CORS_ORIGIN},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "Accept"},
		ExposeHeaders:    []string{"Content-Range", "X-Content-Range"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup file upload directory
	if err := os.MkdirAll(lib.UPLOAD_DIR, os.ModePerm); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Basic route for checking if server is running
	router.GET("/mock", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is running!")
	})

	// Health check route
	router.GET("/mock/api/health", health.HealthCheckHandler)

	// Authentication route
	router.POST("/mock/api/auth", func(c *gin.Context) {
		var loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request",
			})
			return
		}

		if loginRequest.Username != "" && loginRequest.Password != "" {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Login successful",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid credentials",
			})
		}
	})

	// Protected API routes group
	apiGroup := router.Group("/mock/api")
	apiGroup.Use(middlewares.ApiKeyAuth())
	{
		apiGroup.POST("/start", mockHandler.StartMockHandler)
		apiGroup.POST("/stop", mockHandler.StopMockHandler)
		apiGroup.GET("/status", mockHandler.StatusMockHandler)
		apiGroup.POST("/upload", mockHandler.UploadMockHandler)
		apiGroup.GET("/configs", mockHandler.ListConfigsHandler)
		apiGroup.DELETE("/configs/:filename", mockHandler.DeleteConfigHandler)
		apiGroup.GET("/configs/:filename/download", mockHandler.DownloadConfigHandler)
		apiGroup.POST("/sync", handler.SyncToGitHttpHandler)

		apiGroup.POST("/git/save-config", handler.SaveGitConfigHandler)
		apiGroup.POST("/git/save-and-test-sync", handler.SaveAndTestSyncGitHandler)
		apiGroup.GET("/git/config", handler.GetGitConfigHandler)
	}

	// Start the server
	serverAddr := lib.SERVER_HOSTNAME + ":" + lib.SERVER_PORT
	log.Printf("Server is running on http://%s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
