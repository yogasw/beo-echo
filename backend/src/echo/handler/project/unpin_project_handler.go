package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
)

// UnpinProjectHandler removes a project from the authenticated user's pinned list.
//
// POST /api/workspaces/:workspaceID/projects/:projectID/unpin
func UnpinProjectHandler(c *gin.Context) {
	workspaceID := c.Param("workspaceID")
	projectID := c.Param("projectId")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not authenticated"})
		return
	}
	userIDStr, _ := userID.(string)

	// workspaceID is used for validation (the project must belong to this workspace)
	var proj database.Project
	if err := database.DB.Where("id = ? AND workspace_id = ?", projectID, workspaceID).First(&proj).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Project not found"})
		return
	}

	result := database.DB.Where("user_id = ? AND project_id = ?", userIDStr, projectID).
		Delete(&database.UserPinnedProject{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to unpin: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Project unpinned"})
}
