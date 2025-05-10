package endpoint

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

/*
UpdateEndpointHandler updates an existing endpoint

Sample curl:

	curl -X PUT "http://localhost:3600/mock/api/projects/my-new-project/endpoints/1" \
	  -H "Content-Type: application/json" \
	  -d '{
	    "method": "POST",
	    "path": "/api/users/new",
	    "enabled": true,
	    "responseMode": "random"
	  }'
*/
func UpdateEndpointHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectId := c.Param("projectId")
	if projectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	// Parse endpoint ID
	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Endpoint ID is required",
		})
		return
	}

	// Check if endpoint exists
	var existingEndpoint database.MockEndpoint
	result := database.GetDB().Where("id = ? AND project_id = ?", endpointID, projectId).First(&existingEndpoint)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Endpoint not found",
		})
		return
	}

	// Parse update data
	var updateData struct {
		Method       string `json:"method"`
		Path         string `json:"path"`
		Enabled      *bool  `json:"enabled"`
		ResponseMode string `json:"responseMode"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Apply updates
	if updateData.Method != "" {
		existingEndpoint.Method = strings.ToUpper(updateData.Method)
	}

	if updateData.Path != "" {
		// Make sure path starts with /
		if !strings.HasPrefix(updateData.Path, "/") {
			existingEndpoint.Path = "/" + updateData.Path
		} else {
			existingEndpoint.Path = updateData.Path
		}
	}

	if updateData.Enabled != nil {
		existingEndpoint.Enabled = *updateData.Enabled
	}

	if updateData.ResponseMode != "" {
		existingEndpoint.ResponseMode = updateData.ResponseMode
	}

	// Save updates
	result = database.GetDB().Save(&existingEndpoint)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to update endpoint: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Endpoint updated successfully",
		"data":    existingEndpoint,
	})
}
