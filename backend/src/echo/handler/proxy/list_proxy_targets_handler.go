package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

// ListProxyTargetsHandler lists all proxy targets for a project
//
// Sample curl:
// curl -X GET "http://localhost:8000/mock/api/workspaces/{workspaceID}/projects/{projectId}/proxies" -H "Content-Type: application/json" -H "Authorization: Bearer {token}"
func ListProxyTargetsHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectId := c.Param("projectId")
	if projectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	// Find project first
	var project database.Project
	result := database.GetDB().Where("id = ?", projectId).First(&project)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not found",
		})
		return
	}

	// Get proxy targets for this project
	var proxyTargets []database.ProxyTarget
	result = database.GetDB().
		Where("project_id = ?", project.ID).
		Order("label").
		Find(&proxyTargets)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve proxy targets: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    proxyTargets,
	})
}
