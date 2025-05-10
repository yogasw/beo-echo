package handler

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/repositories"
	"mockoon-control-panel/backend_new/src/mocks/services"
)

var mockService *services.MockService

// InitMockService initializes the mock service
func InitMockService() {
	db := database.GetDB() // Get the database connection
	if db == nil {
		log.Println("Warning: Database connection not available for mock service")
		return
	}

	repo := repositories.NewMockRepository(db)
	mockService = services.NewMockService(repo)
}

// MockRequestHandler is a catch-all handler for mock API endpoints
func MockRequestHandler(c *gin.Context) {
	// Check if mock service is initialized
	if mockService == nil {
		InitMockService()
		if mockService == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Mock service is not available",
			})
			return
		}
	}

	// Get project name from subdomain or path
	projectName := extractProjectName(c.Request)
	if projectName == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not specified",
		})
		return
	}

	// Get path (remove /mock prefix if present)
	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/mock") {
		path = path[5:]
	}
	if path == "" {
		path = "/"
	}

	// Process the request
	resp, err := mockService.HandleRequest(projectName, c.Request.Method, path, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Error processing mock request: " + err.Error(),
		})
		return
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Send response with proper status code
	c.Status(resp.StatusCode)

	// Copy response body
	if resp.Body != nil {
		defer resp.Body.Close()
		// Copy body to response writer
		if body, err := io.ReadAll(resp.Body); err == nil {
			c.Writer.Write(body)
		}
	}
}

// RegisterMockHandlers registers all mock API handlers
func RegisterMockHandlers(router *gin.Engine) {
	// Initialize the mock service
	InitMockService()

	// Register API endpoints for managing mocks
	api := router.Group("/mock/api")
	{
		api.GET("/projects", ListProjectsHandler)
		api.POST("/projects", CreateProjectHandler)
		api.GET("/projects/:name", GetProjectHandler)
		api.PUT("/projects/:name", UpdateProjectHandler)
		api.DELETE("/projects/:name", DeleteProjectHandler)

		// Endpoint management
		api.GET("/projects/:name/endpoints", ListEndpointsHandler)
		api.POST("/projects/:name/endpoints", CreateEndpointHandler)
		api.GET("/projects/:name/endpoints/:id", GetEndpointHandler)
		api.PUT("/projects/:name/endpoints/:id", UpdateEndpointHandler)
		api.DELETE("/projects/:name/endpoints/:id", DeleteEndpointHandler)

		// Response management
		api.GET("/endpoints/:id/responses", ListResponsesHandler)
		api.POST("/endpoints/:id/responses", CreateResponseHandler)
		api.GET("/responses/:id", GetResponseHandler)
		api.PUT("/responses/:id", UpdateResponseHandler)
		api.DELETE("/responses/:id", DeleteResponseHandler)
	}

	// Catch-all route for all unhandled requests (low priority)
	// This should be registered last so it doesn't interfere with other routes
	router.NoRoute(MockRequestHandler)
}

// extractProjectName extracts project name from request (subdomain or path)
func extractProjectName(req *http.Request) string {
	// Try to extract from Host header (subdomain)
	host := req.Host

	// Check for subdomain
	parts := strings.Split(host, ".")
	if len(parts) > 2 {
		return parts[0]
	}

	// Try to extract from path
	path := req.URL.Path

	// If path starts with /mock, remove it
	if strings.HasPrefix(path, "/mock") {
		path = path[5:]
	}

	// Extract first part of path
	parts = strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}

	// Default project
	return "default"
}

// Project APIs
func ListProjectsHandler(c *gin.Context) {
	if mockService == nil {
		InitMockService()
	}

	var projects []database.Project
	result := database.GetDB().Find(&projects)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve projects: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projects,
	})
}

func CreateProjectHandler(c *gin.Context) {
	if mockService == nil {
		InitMockService()
	}

	var project database.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Validate project data
	if project.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project name is required",
		})
		return
	}

	// Default to mock mode if not specified
	if project.Mode == "" {
		project.Mode = database.ModeMock
	}

	// Create the project
	result := database.GetDB().Create(&project)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create project: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Project created successfully",
		"data":    project,
	})
}

func GetProjectHandler(c *gin.Context) {
	if mockService == nil {
		InitMockService()
	}

	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project name is required",
		})
		return
	}

	var project database.Project
	result := database.GetDB().
		Preload("Endpoints").
		Preload("ProxyTargets").
		Preload("ActiveProxy").
		Where("name = ?", name).
		First(&project)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    project,
	})
}

func UpdateProjectHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func DeleteProjectHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func ListEndpointsHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func CreateEndpointHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func GetEndpointHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func UpdateEndpointHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func DeleteEndpointHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func ListResponsesHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func CreateResponseHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func GetResponseHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func UpdateResponseHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

func DeleteResponseHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}
