package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// DeleteProxyTargetHandler deletes a proxy target
//
// Sample curl:
// curl -X DELETE "http://localhost:8000/mock/api/workspaces/{workspaceID}/projects/{projectId}/proxies/{proxyId}" -H "Content-Type: application/json" -H "Authorization: Bearer {token}"
func DeleteProxyTargetHandler(c *gin.Context) {
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

	// Check if proxy target exists
	var proxyTarget database.ProxyTarget
	result = database.GetDB().Where("id = ? AND project_id = ?", proxyID, project.ID).First(&proxyTarget)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Proxy target not found",
		})
		return
	}

	// Check if this is the active proxy target for the project
	if project.ActiveProxyID != nil && *project.ActiveProxyID == proxyID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Cannot delete active proxy target. Set another proxy target as active first.",
		})
		return
	}

	// Delete the proxy target
	result = database.GetDB().Delete(&proxyTarget)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to delete proxy target: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Proxy target deleted successfully",
	})
}
