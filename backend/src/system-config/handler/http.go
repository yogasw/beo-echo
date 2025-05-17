package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/system-config/services"
)

// GetSystemConfigHandler returns a specific system configuration by key
func GetSystemConfigHandler(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Config key is required",
		})
		return
	}

	// Check if user is authenticated
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Check if the user is an owner for non-feature configs
	if !strings.HasPrefix(strings.ToLower(key), "feature_") && !strings.HasPrefix(key, "FEATURE_") {
		isOwner, exists := c.Get("isOwner")
		if !exists || isOwner != true {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Only system owners can access this configuration",
			})
			return
		}
	}

	// Find the config
	var config database.SystemConfig
	result := database.DB.Where("key = ?", key).First(&config)
	if result.Error != nil {
		// Try to get from default configs
		for defaultKey := range services.DefaultVariables {
			parts := strings.Split(defaultKey, ":")
			if parts[0] == key {
				configType := "string"
				if len(parts) > 1 {
					configType = parts[1]
				}

				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"data": gin.H{
						"key":         key,
						"value":       services.DefaultVariables[defaultKey],
						"type":        configType,
						"description": "Default configuration",
						"hide_value":  false,
					},
				})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Configuration not found",
		})
		return
	}

	// Only return hideValue: true configs to owners
	if config.HideValue {
		isOwner, exists := c.Get("isOwner")
		if !exists || isOwner != true {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "You do not have permission to view this configuration",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    config,
	})
}

// GetAllSystemConfigsHandler returns all visible system configurations
func GetAllSystemConfigsHandler(c *gin.Context) {
	// Check if user is authenticated
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Get owner status from context (set by JWTAuthMiddleware)
	isOwnerValue, exists := c.Get("isOwner")
	isOwner := exists && isOwnerValue == true

	// Get configs from services
	configs, err := services.GetAllSystemConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve configurations: " + err.Error(),
		})
		return
	}

	// Filter configs based on user permissions
	var visibleConfigs []database.SystemConfig
	for _, config := range configs {
		// If user is not an owner, only show feature flags and non-hidden configs
		if isOwner || strings.HasPrefix(strings.ToLower(config.Key), "feature_") ||
			strings.HasPrefix(config.Key, "FEATURE_") || !config.HideValue {
			visibleConfigs = append(visibleConfigs, config)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    visibleConfigs,
	})
}

// UpdateSystemConfigRequest represents the update system config request
type UpdateSystemConfigRequest struct {
	Value string `json:"value" binding:"required"`
}

// UpdateSystemConfigHandler updates a system configuration
func UpdateSystemConfigHandler(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Config key is required",
		})
		return
	}

	// Authentication already verified by middleware

	var req UpdateSystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// Check if this is a feature flag
	isFeatureFlag := strings.HasPrefix(strings.ToLower(key), "feature_") || strings.HasPrefix(key, "FEATURE_")

	// If it's a feature flag, ensure it's set as a boolean type
	configType := "string"
	if isFeatureFlag {
		configType = "boolean"

		// Validate that the value is a valid boolean for feature flags
		if req.Value != "true" && req.Value != "false" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Feature flag value must be 'true' or 'false'",
			})
			return
		}
	}

	// Find the config
	var config database.SystemConfig
	result := database.DB.Where("key = ?", key).First(&config)
	if result.Error != nil {
		// Config doesn't exist, create a new one
		description := ""
		if isFeatureFlag {
			description = "Feature flag created via API"
		}

		// Use the utility function to create new config
		newConfig, err := services.AddConfig(key, req.Value, description, configType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create configuration: " + err.Error(),
			})
			return
		}
		config = *newConfig
	} else {
		// Use the utility function to update existing config
		keyWithType := key
		if config.Type != "" {
			keyWithType = key + ":" + config.Type
		}

		err := services.SetSystemConfig(keyWithType, req.Value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to update configuration: " + err.Error(),
			})
			return
		}

		// Refresh config object to return in response
		database.DB.Where("key = ?", key).First(&config)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Configuration updated successfully",
		"data":    config,
	})
}
