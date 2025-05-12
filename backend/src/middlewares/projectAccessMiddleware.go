package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mockoon-control-panel/backend_new/src/database"
)

// ProjectAccessMiddleware verifies that the user has access to the specified project
// This middleware should be used on all routes containing a projectId parameter
func ProjectAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the project ID from the URL parameters
		projectID := c.Param("projectId")
		if projectID == "" {
			// No project ID in the URL, skip this middleware
			c.Next()
			return
		}

		// Get authenticated user ID from context (set by JWTAuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Check if user is a system owner (can access all projects)
		isOwner, ownerExists := c.Get("isOwner")
		if ownerExists && isOwner == true {
			// System owners can access all projects
			c.Next()
			return
		}

		// For regular users, check if the project belongs to one of their workspaces
		var project database.Project
		if err := database.GetDB().First(&project, "id = ?", projectID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"error":   true,
					"message": "Project not found",
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to retrieve project information: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Get the workspaceID of the project
		workspaceID := project.WorkspaceID

		// Check if user is a member of this workspace
		var userWorkspace database.UserWorkspace
		if err := database.GetDB().
			Where("user_id = ? AND workspace_id = ?", userID, workspaceID).
			First(&userWorkspace).Error; err != nil {

			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   true,
					"message": "You do not have access to this project",
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to verify project access: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Store the workspace ID in the context for potential use in handlers
		c.Set("workspaceID", workspaceID)

		// User has access to the project, continue
		c.Next()
	}
}
