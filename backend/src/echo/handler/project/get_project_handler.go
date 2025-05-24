package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

// GetProjectHandler retrieves a project by name
//
// Sample curl:
// curl -X GET "http://localhost:3600/api/api/projects/my-new-project" -H "Content-Type: application/json"
func GetProjectHandler(c *gin.Context) {
	handler.EnsureMockService()

	id := c.Param("projectId")
	if id == "" {
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
		Preload("Endpoints.Responses").
		Where("id = ?", id).
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
