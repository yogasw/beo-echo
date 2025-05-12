package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mockoon-control-panel/backend_new/src/database"
)

// GetCurrentUserHandler returns the authenticated user's information
func GetCurrentUserHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Fetch user details
	var user database.User
	if err := database.DB.Select("id, email, name, is_owner, created_at, updated_at").
		Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve user information: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"is_owner":   user.IsOwner,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

// GetWorkspaceHandler returns details of a specific workspace
func GetWorkspaceHandler(c *gin.Context) {
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
	// Query the database directly to check if the user is an owner
	var user database.User
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID format",
		})
		return
	}

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

	// Get workspace details
	var workspace database.Workspace
	if err := database.DB.First(&workspace, "id = ?", workspaceID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Workspace not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve workspace: " + err.Error(),
		})
		return
	}

	// Count members and projects
	var memberCount int64
	database.DB.Model(&database.UserWorkspace{}).Where("workspace_id = ?", workspaceID).Count(&memberCount)

	var projectCount int64
	database.DB.Model(&database.Project{}).Where("workspace_id = ?", workspaceID).Count(&projectCount)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":            workspace.ID,
			"name":          workspace.Name,
			"created_at":    workspace.CreatedAt,
			"updated_at":    workspace.UpdatedAt,
			"member_count":  memberCount,
			"project_count": projectCount,
		},
	})
}
