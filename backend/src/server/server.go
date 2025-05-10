package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/git-sync/handler"
	"mockoon-control-panel/backend_new/src/health"
	"mockoon-control-panel/backend_new/src/lib"
	"mockoon-control-panel/backend_new/src/middlewares"
	mockHandler "mockoon-control-panel/backend_new/src/mocks/handler"
)

// SetupRouter creates and configures a new Gin router
func SetupRouter() *gin.Engine {
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
		log.Printf("Warning: Failed to create upload directory: %v", err)
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

	// Register mock handlers to handle all routes
	mockHandler.RegisterMockHandlers(router)

	return router
}

// StartServer initializes and starts the HTTP server
func StartServer() error {
	router := SetupRouter()

	// Start the server
	serverAddr := lib.SERVER_HOSTNAME + ":" + lib.SERVER_PORT

	log.Printf("=================================================")
	log.Printf("🚀 Mockoon Control Panel server is starting up!")
	log.Printf("🔗 Server URL: http://%s", serverAddr)
	log.Printf("📄 API endpoint: http://%s/mock/api", serverAddr)
	log.Printf("🔍 Health check: http://%s/mock/api/health", serverAddr)
	log.Printf("=================================================")

	// This will block until the server is stopped
	return router.Run(serverAddr)
}
