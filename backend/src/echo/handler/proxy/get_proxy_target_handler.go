package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

// GetProxyTargetHandler gets a proxy target by ID
//
// Sample curl:
// curl -X GET "http://localhost:8000/mock/api/workspaces/{workspaceID}/projects/{projectId}/proxies/{proxyId}" -H "Content-Type: application/json" -H "Authorization: Bearer {token}"
func GetProxyTargetHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectId := c.Param("projectId")
	if projectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	// Parse proxy target ID
	proxyID := c.Param("proxyId")
	if proxyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Proxy target ID is required",
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

	// Get proxy target
	var proxyTarget database.ProxyTarget
	result = database.GetDB().
		Where("id = ? AND project_id = ?", proxyID, project.ID).
		First(&proxyTarget)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Proxy target not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    proxyTarget,
	})
}
