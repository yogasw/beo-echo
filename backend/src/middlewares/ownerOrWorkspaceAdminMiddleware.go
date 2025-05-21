package middlewares

import (
	"beo-echo/backend/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OwnerOrWorkspaceAdminMiddleware checks if the current user is either a system owner
// or an admin in the specified workspace. This middleware combines privileges for routes
// that can be accessed by either system owners or workspace admins.
func OwnerOrWorkspaceAdminMiddleware() gin.HandlerFunc {
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

		// Check if the user is a system owner
		var user database.User
		if err := database.DB.Select("is_owner").Where("id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to verify user permissions",
			})
			c.Abort()
			return
		}

		// If the user is a system owner, allow access regardless of workspace role
		if user.IsOwner {
			c.Next()
			return
		}

		// If not a system owner, check if they are a workspace admin
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
				"message": "Access denied: You need to be a workspace admin to perform this action",
			})
			c.Abort()
			return
		}

		// User is a workspace admin, continue
		c.Next()
	}
}
