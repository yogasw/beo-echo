package endpoint

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

/*
CreateEndpointHandler creates a new endpoint for a project

Sample curl:

	curl -X POST "http://localhost:3600/mock/api/projects/my-new-project/endpoints" \
	  -H "Content-Type: application/json" \
	  -d '{
	    "method": "GET",
	    "path": "/api/users",
	    "enabled": true,
	    "responseMode": "static"
	  }'
*/
func CreateEndpointHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectId := c.Param("projectId")
	if projectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project id is required",
		})
		return
	}

	// Check if project exists
	if err := database.GetDB().
		Where("id = ?", projectId).
		First(&database.Project{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not found: " + err.Error(),
		})
		return
	}

	// Parse endpoint data
	var endpoint database.MockEndpoint
	if err := c.ShouldBindJSON(&endpoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Validate endpoint data
	if endpoint.Method == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "HTTP method is required",
		})
		return
	}

	// Normalize method to uppercase
	endpoint.Method = strings.ToUpper(endpoint.Method)

	if endpoint.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Path is required",
		})
		return
	}

	// Make sure path starts with /
	if !strings.HasPrefix(endpoint.Path, "/") {
		endpoint.Path = "/" + endpoint.Path
	}

	// Assign to project
	endpoint.ProjectID = projectId

	// Default values
	if endpoint.ResponseMode == "" {
		endpoint.ResponseMode = "static"
	}

	// Create endpoint
	result := database.GetDB().Create(&endpoint)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create endpoint: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Endpoint created successfully",
		"data":    endpoint,
	})
}
