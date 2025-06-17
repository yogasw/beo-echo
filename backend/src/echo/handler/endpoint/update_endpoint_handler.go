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
UpdateEndpointHandler updates an existing endpoint

Sample curl:

	curl -X PUT "http://localhost:3600/api/api/projects/my-new-project/endpoints/1" \
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
		Method        string  `json:"method"`
		Path          string  `json:"path"`
		Enabled       *bool   `json:"enabled"`
		ResponseMode  string  `json:"response_mode"`
		Documentation string  `json:"documentation"`
		AdvanceConfig *string `json:"advance_config"`  // Changed to pointer to detect if field is provided
		UseProxy      *bool   `json:"use_proxy"`       // Whether to use proxy for this endpoint
		ProxyTargetID *string `json:"proxy_target_id"` // ID of the proxy target to use
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

	if updateData.Documentation != "" {
		existingEndpoint.Documentation = updateData.Documentation
	}

	// Handle advance_config: allow empty string to clear the config
	// Check if advance_config field is provided (not nil)
	if updateData.AdvanceConfig != nil {
		advanceConfig := *updateData.AdvanceConfig

		// Only validate JSON format if the string is not empty
		if advanceConfig != "" {
			// Validate JSON format before saving
			var temp interface{}
			if err := json.Unmarshal([]byte(advanceConfig), &temp); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   true,
					"message": "Invalid JSON format in advance_config: " + err.Error(),
				})
				return
			}
		}
		// Update with the new value (could be empty string)
		existingEndpoint.AdvanceConfig = advanceConfig
	}

	// Update proxy settings
	if updateData.UseProxy != nil {
		existingEndpoint.UseProxy = *updateData.UseProxy

		// If proxy is disabled, clear proxy target
		if !*updateData.UseProxy {
			existingEndpoint.ProxyTargetID = nil
		}
	}

	// Update proxy target if provided and proxy is enabled
	if updateData.ProxyTargetID != nil && (updateData.UseProxy == nil || *updateData.UseProxy) {
		// Verify if the proxy target exists and belongs to this project
		var proxyTarget database.ProxyTarget
		result := database.GetDB().Where("id = ? AND project_id = ?", *updateData.ProxyTargetID, projectId).First(&proxyTarget)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid proxy target: " + result.Error.Error(),
			})
			return
		}
		existingEndpoint.ProxyTargetID = updateData.ProxyTargetID
		existingEndpoint.UseProxy = true
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
