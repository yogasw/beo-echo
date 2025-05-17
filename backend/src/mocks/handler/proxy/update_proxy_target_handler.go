package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/mocks/handler"
)

// UpdateProxyTargetHandler updates a proxy target
//
// Sample curl:
//
//	curl -X PUT "http://localhost:8000/mock/api/workspaces/{workspaceID}/projects/{projectId}/proxies/{proxyId}" \
//	  -H "Content-Type: application/json" \
//	  -H "Authorization: Bearer {token}" \
//	  -d '{
//	    "label": "Staging",
//	    "url": "https://staging.example.com"
//	  }'
func UpdateProxyTargetHandler(c *gin.Context) {
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
	var existingProxy database.ProxyTarget
	result = database.GetDB().Where("id = ? AND project_id = ?", proxyID, project.ID).First(&existingProxy)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Proxy target not found",
		})
		return
	}

	// Parse update data
	var updateData struct {
		Label string `json:"label"`
		URL   string `json:"url"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Apply updates
	if updateData.Label != "" {
		existingProxy.Label = updateData.Label
	}

	if updateData.URL != "" {
		existingProxy.URL = updateData.URL
	}

	// Save updates
	result = database.GetDB().Save(&existingProxy)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to update proxy target: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Proxy target updated successfully",
		"data":    existingProxy,
	})
}
