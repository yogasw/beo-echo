package proxy

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// GetProxyTargetHandler gets a proxy target by ID
//
// Sample curl:
// curl -X GET "http://localhost:8000/api/projects/my-project/proxies/1" -H "Content-Type: application/json"
func GetProxyTargetHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectName := c.Param("name")
	if projectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project name is required",
		})
		return
	}

	// Parse proxy target ID
	proxyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid proxy target ID",
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
