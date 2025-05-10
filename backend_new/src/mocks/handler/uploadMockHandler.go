package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/lib"
	"mockoon-control-panel/backend_new/src/mocks/repositories"
	"mockoon-control-panel/backend_new/src/utils"
)

// UploadMockHandler handles uploading a new mock configuration
func UploadMockHandler(c *gin.Context) {
	// Make sure upload directory exists
	if err := utils.EnsureDirectoryExists(lib.UPLOAD_DIR); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create upload directory: " + err.Error(),
		})
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("config")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No config file provided: " + err.Error(),
		})
		return
	}

	// Validate file is a JSON file
	if !strings.HasSuffix(file.Filename, ".json") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only JSON files are allowed",
		})
		return
	}

	// Create a safe filename
	safeName := strings.Replace(file.Filename, "/", "_", -1)
	safeName = strings.Replace(safeName, "\\", "_", -1)
	safeName = strings.Replace(safeName, "..", "", -1)

	// Create temporary path
	uploadPath := filepath.Join(lib.UPLOAD_DIR, safeName)

	// Save the file
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save uploaded file: " + err.Error(),
		})
		return
	}

	// Copy the file to the configs directory
	if err := repositories.FileRepo.SaveConfigFromUpload(uploadPath, safeName); err != nil {
		os.Remove(uploadPath) // Clean up upload
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save configuration: " + err.Error(),
		})
		return
	}

	// Clean up the temporary file
	os.Remove(uploadPath)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"filename": safeName,
		},
	})
}
