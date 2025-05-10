package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// ListEndpointsHandler lists all endpoints for a project
//
// Sample curl:
// curl -X GET "http://localhost:3600/mock/api/projects/my-new-project/endpoints" -H "Content-Type: application/json"
func ListEndpointsHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectName := c.Param("name")
	if projectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project name is required",
		})
		return
	}

	// Find project first
	var project database.Project
	result := database.GetDB().Where("name = ?", projectName).First(&project)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not found",
		})
		return
	}

	// Get endpoints for this project
	var endpoints []database.MockEndpoint
	result = database.GetDB().
		Where("project_id = ?", project.ID).
		Order("path").
		Order("method").
		Find(&endpoints)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve endpoints: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    endpoints,
	})
}
