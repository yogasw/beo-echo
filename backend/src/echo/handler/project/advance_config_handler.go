package project

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

/*
GetProjectAdvanceConfigHandler retrieves the advance configuration for a project

Sample curl:

	curl -X GET "http://localhost:3600/api/api/projects/my-project/advance-config" \
	  -H "Content-Type: application/json"
*/
func GetProjectAdvanceConfigHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	// Check if project exists
	var project database.Project
	result := database.GetDB().Where("id = ?", projectID).First(&project)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not found",
		})
		return
	}

	// Parse advance config if it exists
	var advanceConfig interface{}
	if project.AdvanceConfig != "" {
		if err := json.Unmarshal([]byte(project.AdvanceConfig), &advanceConfig); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to parse advance config: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"project_id":     project.ID,
			"project_name":   project.Name,
			"advance_config": advanceConfig,
		},
	})
}

/*
UpdateProjectAdvanceConfigHandler updates the advance configuration for a project

Sample curl:

	curl -X PUT "http://localhost:3600/api/api/projects/my-project/advance-config" \
	  -H "Content-Type: application/json" \
	  -d '{
	    "global_timeout": 5000,
	    "rate_limit": {
	      "enabled": true,
	      "requests_per_min": 100,
	      "burst_size": 10
	    },
	    "cors": {
	      "enabled": true,
	      "allowed_origins": ["https://example.com"],
	      "allowed_methods": ["GET", "POST"]
	    }
	  }'
*/
func UpdateProjectAdvanceConfigHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	// Check if project exists
	var project database.Project
	result := database.GetDB().Where("id = ?", projectID).First(&project)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not found",
		})
		return
	}

	// Parse request body as raw JSON
	var configData interface{}
	if err := c.ShouldBindJSON(&configData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid JSON data: " + err.Error(),
		})
		return
	}

	// Convert back to JSON string for storage and validation
	configJSON, err := json.Marshal(configData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to process config data: " + err.Error(),
		})
		return
	}

	// Validate the advance config using our validation function
	_, err = database.ParseProjectAdvanceConfig(string(configJSON))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Update project's advance config
	project.AdvanceConfig = string(configJSON)
	result = database.GetDB().Save(&project)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to update project config: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project advance config updated successfully",
		"data": gin.H{
			"project_id":     project.ID,
			"project_name":   project.Name,
			"advance_config": configData,
		},
	})
}
