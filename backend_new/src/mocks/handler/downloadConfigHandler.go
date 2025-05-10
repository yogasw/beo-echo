package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/mocks/repositories"
)

// DownloadConfigHandler handles requests to download a configuration file
func DownloadConfigHandler(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No filename provided",
		})
		return
	}

	// Check if the file exists
	if !repositories.FileRepo.ConfigExists(fileName) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Configuration file not found",
		})
		return
	}

	// Get the file path
	configPath := repositories.FileRepo.GetConfigPath(fileName)

	// Set response headers for download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/json")

	// Serve the file
	c.File(configPath)
}
