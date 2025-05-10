package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/git-sync/services"
)

// SaveGitConfigHandler handles requests to save Git configuration
func SaveGitConfigHandler(c *gin.Context) {
	var config services.GitConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	if err := services.SaveGitConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save Git configuration: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"success": true,
			"message": "Git configuration saved successfully",
		},
	})
}

// GetGitConfigHandler handles requests to get the Git configuration
func GetGitConfigHandler(c *gin.Context) {
	config, err := services.GetGitConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get Git configuration: " + err.Error(),
		})
		return
	}

	// Mask sensitive information
	if config.AuthType == "basic" {
		config.Password = "********"
	} else if config.AuthType == "ssh" {
		config.SSHKey = "********"
	}

	c.JSON(http.StatusOK, gin.H{
		"data": config,
	})
}

// SaveAndTestSyncGitHandler saves Git configuration and tests the sync
func SaveAndTestSyncGitHandler(c *gin.Context) {
	var config services.GitConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	// Save the configuration
	if err := services.SaveGitConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save Git configuration: " + err.Error(),
		})
		return
	}

	// Test the sync
	if err := SyncConfigsToGit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Git sync test failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"success": true,
			"message": "Git configuration saved and sync test completed successfully",
		},
	})
}

// SyncToGitHttpHandler handles requests to manually trigger a Git sync
func SyncToGitHttpHandler(c *gin.Context) {
	if err := SyncConfigsToGit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Git sync failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"success": true,
			"message": "Git sync completed successfully",
		},
	})
}

// SyncConfigsToGit synchronizes configs with Git repository
func SyncConfigsToGit() error {
	return services.SyncConfigsToGit()
}
