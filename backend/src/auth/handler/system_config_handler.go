package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
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
	if !strings.HasPrefix(key, "feature_") {
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

	// Prepare query
	query := database.DB.Model(&database.SystemConfig{})

	// If user is not an owner, only return feature flags and non-hidden values
	if !isOwner {
		query = query.Where("key LIKE ? OR hide_value = ?", "feature_%", false)
	}

	var configs []database.SystemConfig
	result := query.Find(&configs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve configurations: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    configs,
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

	// Find the config
	var config database.SystemConfig
	result := database.DB.Where("key = ?", key).First(&config)
	if result.Error != nil {
		// Config doesn't exist, create a new one
		config = database.SystemConfig{
			Key:         key,
			Value:       req.Value,
			Type:        "string", // Default to string
			Description: "",
			HideValue:   false,
		}
		result = database.DB.Create(&config)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create configuration: " + result.Error.Error(),
			})
			return
		}
	} else {
		// Update existing config
		config.Value = req.Value
		result = database.DB.Save(&config)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to update configuration: " + result.Error.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Configuration updated successfully",
		"data":    config,
	})
}
