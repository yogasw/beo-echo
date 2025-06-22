package systemConfig

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	systemConfig "beo-echo/backend/src/systemConfigs"
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
	config, err := systemConfig.GetConfigSetting(key)
	if err != nil || config == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Configuration not found",
		})
		return
	}

	// Only return HideValue: true configs to owners
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

	// Get configs from services
	configs, err := systemConfig.GetAllSystemConfigs()
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
		// skip when hide value
		if config.HideValue {
			continue
		}

		// push data to visibleConfigs
		visibleConfigs = append(visibleConfigs, config)
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

	// when key contain : get first part to support legacy configs
	if strings.Contains(key, ":") {
		parts := strings.SplitN(key, ":", 2)
		if len(parts) > 0 {
			key = parts[0]
		}
	}

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

	// validate key and value
	// Get default config to validate key and type
	defaultConfig, exists := systemConfig.DefaultConfigSettings[systemConfig.SystemConfigKey(key)]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid configuration key",
		})
		return
	}

	// disable update hide value from api
	if defaultConfig.HideValue {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "You do not have permission to update this configuration",
		})
		return
	}

	// Ensure value is not empty
	if req.Value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Configuration value cannot be empty",
		})
		return
	}
	// validate type value
	switch defaultConfig.Type {
	case systemConfig.TypeString:
		// String type can accept any value, no additional validation needed
	case systemConfig.TypeBoolean:
		// Boolean type must be "true" or "false"
		if req.Value != "true" && req.Value != "false" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Boolean configuration value must be 'true' or 'false'",
			})
			return
		}
	case systemConfig.TypeNumber:
		// Number type must be a valid integer
		if _, err := strconv.Atoi(req.Value); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Number configuration value must be a valid integer",
			})
			return
		}
	}
	err := systemConfig.SetSystemConfig(key, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update configuration: " + err.Error(),
		})
		return
	}

	config, err := systemConfig.GetConfigSetting(key)
	if err != nil || config == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Configuration not found",
		})
		return
	}

	// validate if the config is updated
	if config.Value != req.Value {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Configuration value was not updated successfully",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Configuration updated successfully",
		"data":    config,
	})
}
