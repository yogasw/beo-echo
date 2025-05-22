package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/mocks/handler"
)

// ListProjectsHandler lists projects accessible to the authenticated user
//
// Sample curl:
// curl -X GET "http://localhost:3600/mock/api/projects" -H "Content-Type: application/json" -H "Authorization: Bearer TOKEN"

// GetWorkspaceProjectsHandler returns all projects in a workspace
func ListProjectsHandler(c *gin.Context) {
	workspaceID := c.Param("workspaceID")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Workspace ID is required",
		})
		return
	}

	// Get user ID from context (set by JWTAuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Check if user is a system admin (can access all workspaces)
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID format",
		})
		return
	}

	// Directly query database to check if user is an owner
	var user database.User
	err := database.DB.Where("id = ?", userIDStr).First(&user).Error
	isSystemOwner := err == nil && user.IsOwner

	if !isSystemOwner {
		// Check if the user is a member of this workspace
		var userWorkspace database.UserWorkspace
		err := database.DB.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&userWorkspace).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "You do not have access to this workspace",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to verify workspace access: " + err.Error(),
			})
			return
		}
	}

	var projects []database.Project
	if err := database.DB.Where("workspace_id = ?", workspaceID).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch projects: " + err.Error(),
		})
		return
	}

	// Add project URLs
	for i := range projects {
		projects[i].URL = handler.GetProjectURL(c.Request.Host, projects[i])
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projects,
	})
}
