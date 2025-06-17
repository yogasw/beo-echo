package endpoint

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

/*
CreateEndpointHandler creates a new endpoint for a project

Sample curl:

	curl -X POST "http://localhost:3600/api/api/projects/my-new-project/endpoints" \
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
		endpoint.ResponseMode = "random"
	}

	// Validate JSON format in AdvanceConfig if provided
	if endpoint.AdvanceConfig != "" {
		var temp interface{}
		if err := json.Unmarshal([]byte(endpoint.AdvanceConfig), &temp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid JSON format in advance_config: " + err.Error(),
			})
			return
		}
	}

	// Validate proxy target if proxy is enabled
	if endpoint.UseProxy && endpoint.ProxyTargetID != nil {
		var proxyTarget database.ProxyTarget
		result := database.GetDB().Where("id = ? AND project_id = ?", *endpoint.ProxyTargetID, projectId).First(&proxyTarget)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid proxy target: " + result.Error.Error(),
			})
			return
		}
	} else if endpoint.UseProxy && endpoint.ProxyTargetID == nil {
		// Can't enable proxy without a target
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "ProxyTargetID is required when UseProxy is enabled",
		})
		return
	}

	// validation alredy added or not when creating the project response error

	check := database.GetDB().
		Where("project_id = ? AND path = ?", projectId, endpoint.Path).
		Where("method = ?", endpoint.Method).
		First(&database.MockEndpoint{})
	if check.Error != nil {
		if !strings.Contains(check.Error.Error(), "record not found") {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to check existing endpoint: " + check.Error.Error(),
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Endpoint already exists for this project",
		})
		return
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
