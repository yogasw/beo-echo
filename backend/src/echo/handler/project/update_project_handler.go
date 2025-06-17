package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

/*
UpdateProjectHandler updates an existing project

Sample curl:

	curl -X PUT "http://localhost:3600/api/api/projects/my-new-project" \
	  -H "Content-Type: application/json" \
	  -d '{
	    "mode": "proxy",
	    "activeProxyID": 1
	  }'
*/
func UpdateProjectHandler(c *gin.Context) {
	handler.EnsureMockService()

	id := c.Param("projectId")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project id is required",
		})
		return
	}

	// Check if project exists
	var existingProject database.Project
	result := database.GetDB().Where("id = ?", id).First(&existingProject)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not found",
		})
		return
	}

	// Parse update data - supporting partial updates
	var updateData struct {
		Name          *string               `json:"name"`
		Alias         *string               `json:"alias"`
		Mode          *database.ProjectMode `json:"mode"`
		ActiveProxyID *string               `json:"active_proxy_id"`
		Status        *string               `json:"status"`
		AdvanceConfig *string               `json:"advance_config"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Apply updates only for fields that were provided
	if updateData.Name != nil {
		existingProject.Name = *updateData.Name
	}

	if updateData.Alias != nil {
		if !handler.IsValidAlias(*updateData.Alias) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid alias format for project. Only alphanumeric characters, dashes and underscores are allowed",
			})
			return
		}

		existingProject.Alias = *updateData.Alias
	}

	if updateData.Mode != nil {
		existingProject.Mode = *updateData.Mode
	}

	if updateData.Status != nil {
		// Validate status value
		if *updateData.Status == "running" || *updateData.Status == "stopped" || *updateData.Status == "error" {
			existingProject.Status = *updateData.Status
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid status value. Must be 'running', 'stopped', or 'error'",
			})
			return
		}
	}

	if updateData.ActiveProxyID != nil {
		// Validate that the proxy target exists and belongs to this project
		if *updateData.ActiveProxyID != "" {
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

	if updateData.AdvanceConfig != nil {
		// Validate advance config if provided and not empty
		if *updateData.AdvanceConfig != "" {
			_, err := database.ParseProjectAdvanceConfig(*updateData.AdvanceConfig)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   true,
					"message": err.Error(),
				})
				return
			}
		}

		existingProject.AdvanceConfig = *updateData.AdvanceConfig
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
