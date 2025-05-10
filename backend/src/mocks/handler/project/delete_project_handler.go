package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// DeleteProjectHandler removes a project
//
// Sample curl:
// curl -X DELETE "http://localhost:3600/mock/api/projects/my-new-project" -H "Content-Type: application/json"
func DeleteProjectHandler(c *gin.Context) {
	handler.EnsureMockService()

	id := c.Param("projectId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project id is required",
		})
		return
	}

	// Check if project exists
	var project database.Project
	result := database.GetDB().Where("id = ?", id).First(&project)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not found",
		})
		return
	}

	// Delete the project (GORM will cascade delete related records due to constraints)
	result = database.GetDB().Delete(&project)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to delete project: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project deleted successfully",
	})
}
