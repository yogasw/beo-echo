package src

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
	"mockoon-control-panel/backend_new/src/health"
	"mockoon-control-panel/backend_new/src/lib"
	"mockoon-control-panel/backend_new/src/mocks/handler/endpoint"
	"mockoon-control-panel/backend_new/src/mocks/handler/project"
	"mockoon-control-panel/backend_new/src/mocks/handler/response"
	"mockoon-control-panel/backend_new/src/traefik"
	"mockoon-control-panel/backend_new/src/utils"
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
	// apiGroup.Use(middlewares.ApiKeyAuth())
	{
		apiGroup.GET("/projects", project.ListProjectsHandler)
		apiGroup.POST("/projects", project.CreateProjectHandler)
		apiGroup.GET("/projects/:name", project.GetProjectHandler)
		apiGroup.PUT("/projects/:name", project.UpdateProjectHandler)
		apiGroup.DELETE("/projects/:name", project.DeleteProjectHandler)

		// Endpoint management
		apiGroup.GET("/projects/:name/endpoints", endpoint.ListEndpointsHandler)
		apiGroup.POST("/projects/:name/endpoints", endpoint.CreateEndpointHandler)
		apiGroup.GET("/projects/:name/endpoints/:id", endpoint.GetEndpointHandler)
		apiGroup.PUT("/projects/:name/endpoints/:id", endpoint.UpdateEndpointHandler)
		apiGroup.DELETE("/projects/:name/endpoints/:id", endpoint.DeleteEndpointHandler)

		// Response management
		apiGroup.GET("/projects/:name/endpoints/:id/responses", response.ListResponsesHandler)
		apiGroup.POST("/projects/:name/endpoints/:id/responses", response.CreateResponseHandler)
		apiGroup.GET("/projects/:name/endpoints/:id/responses/:responseId", response.GetResponseHandler)
		apiGroup.PUT("/projects/:name/endpoints/:id/responses/:responseId", response.UpdateResponseHandler)
		apiGroup.DELETE("/projects/:name/endpoints/:id/responses/:responseId", response.DeleteResponseHandler)
	}

	return router
}

// StartServer initializes and starts the HTTP server
func StartServer() error {
	// Load environment variables from .env file
	if err := godotenv.Load(filepath.Join("..", ".env")); err != nil {
		log.Println("Warning: .env file not found or could not be loaded")
	}

	// Setup required directories
	if err := utils.EnsureRequiredFoldersAndEnv(); err != nil {
		log.Fatalf("Failed to setup required folders: %v", err)
	}

	// Setup database connection
	if err := database.CheckAndHandlePrisma(); err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	if err := traefik.GenerateStaticTraefikConfig(); err != nil {
		log.Fatalf("Error generating static Traefik config: %v", err)
	}

	router := SetupRouter()

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
