package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/mocks/repositories"
)

// DeleteConfigHandler handles requests to delete a configuration file
func DeleteConfigHandler(c *gin.Context) {
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

	// Check if the file is in use
	running, port := repositories.MockInstanceRepo.IsRunningWithConfig(fileName)
	if running {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Cannot delete configuration file that is currently in use (running on port " + strconv.Itoa(port) + ")",
		})
		return
	}

	// Delete the file
	if err := repositories.FileRepo.DeleteConfig(fileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete configuration file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"deleted": fileName,
		},
	})
}
