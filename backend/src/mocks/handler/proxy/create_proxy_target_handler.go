package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// CreateProxyTargetHandler creates a new proxy target for a project
//
// Sample curl:
//
//	curl -X POST "http://localhost:8000/api/projects/my-project/proxies" \
//	  -H "Content-Type: application/json" \
//	  -d '{
//	    "label": "Production",
//	    "url": "https://api.example.com"
//	  }'
func CreateProxyTargetHandler(c *gin.Context) {
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

	// Parse proxy target data
	var proxyTarget database.ProxyTarget
	if err := c.ShouldBindJSON(&proxyTarget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Validate proxy target data
	if proxyTarget.Label == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Proxy target label is required",
		})
		return
	}

	if proxyTarget.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Proxy target URL is required",
		})
		return
	}

	// Assign to project
	proxyTarget.ProjectID = project.ID

	// Create proxy target
	result = database.GetDB().Create(&proxyTarget)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create proxy target: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Proxy target created successfully",
		"data":    proxyTarget,
	})
}
