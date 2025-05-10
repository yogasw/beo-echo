package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// GetProjectHandler retrieves a project by name
//
// Sample curl:
// curl -X GET "http://localhost:8000/api/projects/my-project" -H "Content-Type: application/json"
func GetProjectHandler(c *gin.Context) {
	handler.EnsureMockService()

	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project name is required",
		})
		return
	}

	var project database.Project
	result := database.GetDB().
		Preload("Endpoints").
		Preload("ProxyTargets").
		Preload("ActiveProxy").
		Where("name = ?", name).
		First(&project)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    project,
	})
}
