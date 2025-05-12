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
	"mockoon-control-panel/backend_new/src/middlewares"
	"mockoon-control-panel/backend_new/src/mocks/handler"
	"mockoon-control-panel/backend_new/src/mocks/handler/endpoint"
	"mockoon-control-panel/backend_new/src/mocks/handler/project"
	"mockoon-control-panel/backend_new/src/mocks/handler/response"
	"mockoon-control-panel/backend_new/src/traefik"
	"mockoon-control-panel/backend_new/src/utils"

	// New imports for auth and workspace management
	authHandler "mockoon-control-panel/backend_new/src/auth/handler"
)

// SetupRouter creates and configures a new Gin router
func SetupRouter() *gin.Engine {
	// Create Gin router with default middleware
	router := gin.Default()

	// middleware to log requests to the database
	router.Use(middlewares.RequestLoggerMiddleware(database.DB))

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

	// Authentication routes
	router.POST("/mock/api/auth/login", authHandler.LoginHandler)
	router.POST("/mock/api/auth/register", authHandler.RegisterHandler)

	// Protected API routes group
	apiGroup := router.Group("/mock/api")
	apiGroup.Use(middlewares.JWTAuthMiddleware())
	{
		// User-related routes
		apiGroup.GET("/auth/me", authHandler.GetCurrentUserHandler)

		// General workspace-related routes
		apiGroup.GET("/workspaces", authHandler.GetUserWorkspacesHandler)
		apiGroup.POST("/workspaces", authHandler.CreateWorkspaceHandler)
		apiGroup.GET("/workspaces/:workspaceID/role", authHandler.CheckWorkspaceRoleHandler)

		// Workspace-project hierarchy routes (nested)
		workspaceRoutes := apiGroup.Group("/workspaces/:workspaceID")
		{
			// Projects list and creation for a workspace
			workspaceRoutes.GET("/projects", project.ListProjectsHandler)
			workspaceRoutes.POST("/projects", project.CreateProjectWithWorkspaceHandler)

			// Project-specific routes with workspace context
			projectRoutes := workspaceRoutes.Group("/projects/:projectId")
			projectRoutes.Use(middlewares.WorkspaceProjectAccessMiddleware())
			{
				// Project CRUD
				projectRoutes.GET("", project.GetProjectHandler)
				projectRoutes.PUT("", project.UpdateProjectHandler)
				projectRoutes.DELETE("", project.DeleteProjectHandler)

				// Endpoint management
				projectRoutes.GET("/endpoints", endpoint.ListEndpointsHandler)
				projectRoutes.POST("/endpoints", endpoint.CreateEndpointHandler)
				projectRoutes.GET("/endpoints/:id", endpoint.GetEndpointHandler)
				projectRoutes.PUT("/endpoints/:id", endpoint.UpdateEndpointHandler)
				projectRoutes.DELETE("/endpoints/:id", endpoint.DeleteEndpointHandler)

				// Response management
				projectRoutes.GET("/endpoints/:id/responses", response.ListResponsesHandler)
				projectRoutes.POST("/endpoints/:id/responses", response.CreateResponseHandler)
				projectRoutes.GET("/endpoints/:id/responses/:responseId", response.GetResponseHandler)
				projectRoutes.PUT("/endpoints/:id/responses/:responseId", response.UpdateResponseHandler)
				projectRoutes.DELETE("/endpoints/:id/responses/:responseId", response.DeleteResponseHandler)

				// Request Logs management
				projectRoutes.GET("/logs", handler.GetLogsHandler)
				projectRoutes.GET("/logs/stream", handler.StreamLogsHandler)
			}
		}
	}

	// Register the catch-all handler for mock API endpoints
	// We need to avoid conflict with the /mock path, so we'll create a separate group
	// for the mock project endpoints
	mockProjectGroup := router.Group("")
	{
		// This handler will catch any request that doesn't match the above routes
		// particularly targeting project-specific mock endpoints
		mockProjectGroup.Any("/:project/*path", handler.MockRequestHandler)
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

	// Initialize services
	handler.InitLogService()

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
