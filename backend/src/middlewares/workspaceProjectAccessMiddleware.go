package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mockoon-control-panel/backend_new/src/database"
)

// WorkspaceProjectAccessMiddleware verifies that:
// 1. The user has access to the specified workspace
// 2. The project belongs to the specified workspace
// 3. The user has access to the project through workspace membership
func WorkspaceProjectAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the workspace and project IDs from the URL parameters
		workspaceID := c.Param("workspaceID")
		projectID := c.Param("projectId")

		if workspaceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Workspace ID is required",
			})
			c.Abort()
			return
		}

		if projectID == "" {
			// This middleware should only be used on routes with both workspace and project IDs
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Project ID is required",
			})
			c.Abort()
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
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Invalid user ID format",
			})
			c.Abort()
			return
		}

		// Directly query database to check if user is an owner
		var user database.User
		if err := database.GetDB().Where("id = ?", userIDStr).First(&user).Error; err == nil && user.IsOwner {
			// System owners can access all projects but we still need to check if project exists
			// and belongs to the specified workspace
			var project database.Project
			if err := database.GetDB().First(&project, "id = ? AND workspace_id = ?", projectID, workspaceID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusNotFound, gin.H{
						"error":   true,
						"message": "Project not found in the specified workspace",
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

			// Project exists and belongs to the workspace
			c.Next()
			return
		}

		// For regular users, check workspace access first
		var userWorkspace database.UserWorkspace
		if err := database.GetDB().
			Where("user_id = ? AND workspace_id = ?", userID, workspaceID).
			First(&userWorkspace).Error; err != nil {

			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   true,
					"message": "You do not have access to this workspace",
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to verify workspace access: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Now check if the project belongs to this workspace
		var project database.Project
		if err := database.GetDB().First(&project, "id = ? AND workspace_id = ?", projectID, workspaceID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"error":   true,
					"message": "Project not found in the specified workspace",
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

		// Store the workspace ID and project ID in context for handlers
		c.Set("workspaceID", workspaceID)
		c.Set("projectID", projectID)

		// User has access to the workspace and project, continue
		c.Next()
	}
}
