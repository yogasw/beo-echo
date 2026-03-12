package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"beo-echo/backend/src/database"
)

// PinProjectHandler adds a project to the authenticated user's pinned list.
//
// POST /api/workspaces/:workspaceID/projects/:projectID/pin
func PinProjectHandler(c *gin.Context) {
	workspaceID := c.Param("workspaceID")
	projectID := c.Param("projectId")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not authenticated"})
		return
	}
	userIDStr, _ := userID.(string)

	// Verify the project exists and belongs to this workspace
	var proj database.Project
	if err := database.DB.Where("id = ? AND workspace_id = ?", projectID, workspaceID).First(&proj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		}
		return
	}

	// Idempotent: skip if already pinned
	var existing database.UserPinnedProject
	err := database.DB.Where("user_id = ? AND project_id = ?", userIDStr, projectID).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Already pinned"})
		return
	}

	pin := database.UserPinnedProject{
		UserID:      userIDStr,
		ProjectID:   projectID,
		WorkspaceID: workspaceID,
	}
	if err := database.DB.Create(&pin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to pin project: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Project pinned"})
}
