package src

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	authServices "beo-echo/backend/src/auth/services"

	"beo-echo/backend/src/caddy/scripts"
	"beo-echo/backend/src/database"
	"beo-echo/backend/src/database/repositories"
	"beo-echo/backend/src/echo/handler"
	"beo-echo/backend/src/echo/handler/endpoint"
	"beo-echo/backend/src/echo/handler/project"
	"beo-echo/backend/src/echo/handler/proxy"
	"beo-echo/backend/src/echo/handler/response"
	"beo-echo/backend/src/echo/services"
	"beo-echo/backend/src/health"
	"beo-echo/backend/src/lib"
	"beo-echo/backend/src/middlewares"
	"beo-echo/backend/src/users"
	"beo-echo/backend/src/utils"
	"beo-echo/backend/src/workspaces"

	// New imports for auth and workspace management
	authHandler "beo-echo/backend/src/auth/handler"
	handlerLogs "beo-echo/backend/src/logs/handlers"
	replayHandler "beo-echo/backend/src/replay/handlers"
	replayServices "beo-echo/backend/src/replay/services"
	systemConfigHandler "beo-echo/backend/src/systemConfigs/handler"
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
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is running!")
	})

	// Health check route
	router.GET("/api/health", health.HealthCheckHandler)

	// Initialize user repository for auth service
	userRepo := repositories.NewUserRepository(database.DB)
	userService := users.NewUserService(userRepo)
	userHandler := users.NewUserHandler(userService)

	// Initialize OAuth and Auth services
	googleOAuthService := authServices.NewGoogleOAuthService(database.DB)
	googleOAuthHandler := authHandler.NewGoogleOAuthHandler(googleOAuthService)
	oauthConfigHandler := authHandler.NewOAuthConfigHandler(database.DB)

	// Initialize workspace module
	workspaceRepo := repositories.NewWorkspaceRepository(database.DB)
	workspaceService := workspaces.NewWorkspaceService(workspaceRepo)
	workspaceHandler := workspaces.NewWorkspaceHandler(workspaceService)

	// Initialize auto-invite handler
	autoInviteHandler := workspaces.NewAutoInviteHandler(database.DB)

	// Initialize Auth service with user repository
	authHandler.InitAuthService(database.DB, userRepo)

	// Initialize system configuration handle
	ruleRepo := repositories.NewRuleRepository(database.DB)
	responseRepo := repositories.NewResponseRepository(database.DB)
	ruleService := services.NewRuleService(ruleRepo, responseRepo)
	ruleHandler := handler.NewRuleHandler(ruleService)

	// Initialize replay service and handler
	replayRepo := repositories.NewReplayRepository(database.DB)
	replayService := replayServices.NewReplayService(replayRepo)
	replayHandler := replayHandler.NewReplayHandler(replayService)

	// Authentication routes
	router.POST("/api/auth/login", authHandler.LoginHandler)

	// Public OAuth routes
	router.GET("/api/oauth/google/login", googleOAuthHandler.InitiateLogin)
	router.GET("/api/oauth/google/callback", googleOAuthHandler.HandleCallback)

	// Protected API routes group
	apiGroup := router.Group("/api")
	apiGroup.Use(middlewares.JWTAuthMiddleware())
	{
		// Owner-only system configuration routes
		ownerGroup := apiGroup.Group("")
		ownerGroup.Use(middlewares.OwnerOnlyMiddleware())
		{
			apiGroup.GET("/system-config/:key", systemConfigHandler.GetSystemConfigHandler)
			apiGroup.GET("/system-configs", systemConfigHandler.GetAllSystemConfigsHandler)
			ownerGroup.PUT("/system-config/:key", systemConfigHandler.UpdateSystemConfigHandler)

			// OAuth Configuration Routes
			ownerGroup.GET("/oauth/config", oauthConfigHandler.ListConfigs)

			// Provider-specific OAuth Configuration Routes
			ownerGroup.GET("/oauth/google/config", googleOAuthHandler.GetConfig)
			ownerGroup.PUT("/oauth/google/config", googleOAuthHandler.UpdateConfig)
			ownerGroup.PUT("/oauth/google/state", googleOAuthHandler.UpdateState)
		}

		// User-related routes
		apiGroup.GET("/auth/me", userHandler.GetCurrentUser)
		apiGroup.PATCH("/users/profile", userHandler.UpdateProfile)
		apiGroup.POST("/users/change-password", userHandler.UpdatePassword)

		// Admin/Owner only user management
		ownerGroup.GET("/users", userHandler.GetAllUsers)
		ownerGroup.DELETE("/users/:user_id", userHandler.DeleteUser)
		ownerGroup.PATCH("/users/:user_id", userHandler.UpdateUser)

		// Member invitation and management (accessible by workspace admins and system owners)
		workspaceAdminGroup := apiGroup.Group("/workspaces/:workspaceID")
		workspaceAdminGroup.Use(middlewares.OwnerOrWorkspaceAdminMiddleware())
		{
			workspaceAdminGroup.DELETE("/users/:user_id", userHandler.RemoveWorkspaceUser)
			workspaceAdminGroup.PUT("/users/:user_id/role", userHandler.UpdateWorkspaceUserRole)
			// Workspace-User management
			workspaceAdminGroup.GET("/users", userHandler.GetWorkspaceUsers)
			workspaceAdminGroup.POST("/members", workspaceHandler.AddMember)

		}

		// Register workspace routes directly
		workspacesGroup := apiGroup.Group("/workspaces")
		{
			workspacesGroup.GET("", workspaceHandler.GetUserWorkspacesWithRoles)
			workspacesGroup.POST("", workspaceHandler.CreateWorkspace)
			workspacesGroup.GET("/:workspaceID/role", workspaceHandler.CheckWorkspaceRole)
			workspacesGroup.GET("/all", middlewares.OwnerOnlyMiddleware(), workspaceHandler.GetAllWorkspaces)

			// Auto-invite configuration (only accessible by system owners)
			workspacesGroup.GET("/:workspaceID/auto-invite", middlewares.OwnerOnlyMiddleware(), autoInviteHandler.GetAutoInviteConfig)
			workspacesGroup.PUT("/:workspaceID/auto-invite", middlewares.OwnerOnlyMiddleware(), autoInviteHandler.UpdateAutoInviteConfig)
		}

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

				// Project Advance Config management
				projectRoutes.GET("/advance-config", project.GetProjectAdvanceConfigHandler)
				projectRoutes.PUT("/advance-config", project.UpdateProjectAdvanceConfigHandler)

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

				// Rule management
				projectRoutes.GET("/endpoints/:id/responses/:responseId/rules", ruleHandler.ListRulesHandler)
				projectRoutes.POST("/endpoints/:id/responses/:responseId/rules", ruleHandler.CreateRuleHandler)
				projectRoutes.GET("/endpoints/:id/responses/:responseId/rules/:ruleId", ruleHandler.GetRuleHandler)
				projectRoutes.PUT("/endpoints/:id/responses/:responseId/rules/:ruleId", ruleHandler.UpdateRuleHandler)
				projectRoutes.DELETE("/endpoints/:id/responses/:responseId/rules/:ruleId", ruleHandler.DeleteRuleHandler)
				projectRoutes.DELETE("/endpoints/:id/responses/:responseId/rules", ruleHandler.DeleteAllRulesHandler)

				// Proxy management
				projectRoutes.GET("/proxies", proxy.ListProxyTargetsHandler)
				projectRoutes.POST("/proxies", proxy.CreateProxyTargetHandler)
				projectRoutes.GET("/proxies/:proxyId", proxy.GetProxyTargetHandler)
				projectRoutes.PUT("/proxies/:proxyId", proxy.UpdateProxyTargetHandler)
				projectRoutes.DELETE("/proxies/:proxyId", proxy.DeleteProxyTargetHandler)

				// Request Logs management
				projectRoutes.GET("/logs", handlerLogs.GetLogsHandler)
				projectRoutes.GET("/logs/stream", handlerLogs.StreamLogsHandler)

				// Bookmark Logs management
				projectRoutes.GET("/logs/bookmark", handlerLogs.GetBookmarksHandler)
				projectRoutes.POST("/logs/bookmark", handlerLogs.AddBookmarkHandler)
				projectRoutes.DELETE("/logs/bookmark/:bookmarkId", handlerLogs.DeleteBookmarkHandler)

				// Replay management
				projectRoutes.GET("/replays", replayHandler.ListReplaysHandler)
				projectRoutes.POST("/replays", replayHandler.CreateReplayHandler)
				projectRoutes.GET("/replays/:replayId", replayHandler.GetReplayHandler)
				projectRoutes.PUT("/replays/:replayId", replayHandler.UpdateReplayHandler)
				projectRoutes.POST("/replays/execute", replayHandler.ExecuteReplayHandler)
				projectRoutes.DELETE("/replays/:replayId", replayHandler.DeleteReplayHandler)
				projectRoutes.GET("/replays/:replayId/logs", replayHandler.GetReplayLogsHandler)

			}
		}
	}

	// Register the catch-all handler for mock API endpoints
	// We need to avoid conflict with the /api path, so we'll create a separate group
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
		log.Info().Msgf("Warning: .env file not found or could not be loaded")
	}

	// Setup required directories
	if err := utils.EnsureRequiredFoldersAndEnv(); err != nil {
		log.Fatal().Msgf("Failed to setup required folders: %v", err)
	}

	// Setup database connection
	if err := database.CheckAndHandle(); err != nil {
		log.Fatal().Msgf("Failed to setup database: %v", err)
	}

	// Initialize services
	handlerLogs.InitLogService()

	router := SetupRouter()
	// zero log context
	ctxLog := log.With().
		Str("script", "Candy Setup").
		Logger().
		WithContext(context.Background())

	if err := scripts.InitCaddyConfig(ctxLog); err != nil {
		log.Error().Msgf("Failed to initialize Caddy config: %v", err)
	}

	// Add request logging middleware
	router.Use(func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()
		log.Info().
			// Log request details
			Msgf(
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
	log.Printf("🚀 BeoEcho server is starting up!")
	log.Printf("🔗 Server URL: http://%s", serverAddr)
	log.Printf("📄 API endpoint: http://%s/api", serverAddr)
	log.Printf("🔍 Health check: http://%s/api/health", serverAddr)
	log.Printf("=================================================")

	// This will block until the server is stopped
	return router.Run(serverAddr)
}
