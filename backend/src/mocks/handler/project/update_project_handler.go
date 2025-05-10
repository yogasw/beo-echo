package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

/*
UpdateProjectHandler updates an existing project

Sample curl:

	curl -X PUT "http://localhost:3600/mock/api/projects/my-new-project" \
	  -H "Content-Type: application/json" \
	  -d '{
	    "mode": "proxy",
	    "activeProxyID": 1
	  }'
*/
func UpdateProjectHandler(c *gin.Context) {
	handler.EnsureMockService()

	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project name is required",
		})
		return
	}

	// Check if project exists
	var existingProject database.Project
	result := database.GetDB().Where("name = ?", name).First(&existingProject)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not found",
		})
		return
	}

	// Parse update data
	var updateData struct {
		Mode          database.ProjectMode `json:"mode"`
		ActiveProxyID *uint                `json:"activeProxyID"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Apply updates
	if updateData.Mode != "" {
		existingProject.Mode = updateData.Mode
	}

	if updateData.ActiveProxyID != nil {
		// Validate that the proxy target exists and belongs to this project
		if *updateData.ActiveProxyID > 0 {
			var proxyTarget database.ProxyTarget
			result = database.GetDB().Where("id = ? AND project_id = ?",
				*updateData.ActiveProxyID, existingProject.ID).First(&proxyTarget)

			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   true,
					"message": "Proxy target not found or doesn't belong to this project",
				})
				return
			}
		}

		existingProject.ActiveProxyID = updateData.ActiveProxyID
	}

	// Save updates
	result = database.GetDB().Save(&existingProject)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to update project: " + result.Error.Error(),
		})
		return
	}

	// Reload project with relationships
	database.GetDB().
		Preload("Endpoints").
		Preload("ProxyTargets").
		Preload("ActiveProxy").
		Where("id = ?", existingProject.ID).
		First(&existingProject)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project updated successfully",
		"data":    existingProject,
	})
}
