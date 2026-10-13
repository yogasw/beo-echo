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

// mergeConfig merges src into dst. Nested objects are merged key by key (so {"cors":{"maxAge":600}}
// keeps cors.mode), a null removes the key, and any other value (including arrays) replaces it.
func mergeConfig(dst, src map[string]interface{}) {
	for key, value := range src {
		if value == nil {
			delete(dst, key)
			continue
		}
		srcMap, srcIsMap := value.(map[string]interface{})
		dstMap, dstIsMap := dst[key].(map[string]interface{})
		if srcIsMap && dstIsMap {
			mergeConfig(dstMap, srcMap)
			continue
		}
		dst[key] = value
	}
}

/*
UpdateProjectAdvanceConfigHandler updates the advance configuration for a project.
The config is deep-merged into the stored one: nested objects are merged key by key, a key sent
as null is removed, arrays and scalars replace the stored value.

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
	      "mode": "reflect"
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

	// Parse request body as a JSON object
	var configData map[string]interface{}
	if err := c.ShouldBindJSON(&configData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid JSON data: " + err.Error(),
		})
		return
	}

	// Merge into the stored config: keys that are absent stay untouched, so setting delayMs
	// does not wipe cors (and the other way around). A key set to null removes it.
	merged := map[string]interface{}{}
	if project.AdvanceConfig != "" {
		_ = json.Unmarshal([]byte(project.AdvanceConfig), &merged)
	}
	mergeConfig(merged, configData)
	configData = merged

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
