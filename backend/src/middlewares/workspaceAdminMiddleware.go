package middlewares

import (
	"beo-echo/backend/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WorkspaceAdminMiddleware checks if the current user is an admin in the specified workspace.
// This middleware should be applied to routes that require workspace admin privileges.
func WorkspaceAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the authenticated user ID from the context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Get the workspace ID from the request parameters
		workspaceID := c.Param("workspaceID")
		if workspaceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Workspace ID is required",
			})
			c.Abort()
			return
		}

		// Check if the user is an admin in the workspace
		var userWorkspace database.UserWorkspace
		result := database.DB.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&userWorkspace)

		if result.Error != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Access denied: User does not belong to this workspace",
			})
			c.Abort()
			return
		}

		if userWorkspace.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Access denied: Only workspace admins can perform this action",
			})
			c.Abort()
			return
		}

		// User is an admin, continue
		c.Next()
	}
}
