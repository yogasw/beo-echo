package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/utils"
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

	// Get feature flags from system config
	featureFlags, err := utils.GetFeatureFlags()
	if err != nil {
		// Log the error but don't fail the request
		featureFlags = make(map[string]bool)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":            user.ID,
			"email":         user.Email,
			"name":          user.Name,
			"is_owner":      user.IsOwner,
			"is_enabled":    user.IsEnabled,
			"created_at":    user.CreatedAt,
			"updated_at":    user.UpdatedAt,
			"feature_flags": featureFlags,
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

// UpdatePasswordRequest represents the update password request
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// UpdateUserRequest represents the update user profile request
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdatePasswordHandler handles password update requests
func UpdatePasswordHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// Get the user from database
	user, err := database.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve user: " + err.Error(),
		})
		return
	}

	// Verify current password
	if !user.VerifyPassword(req.CurrentPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Current password is incorrect",
		})
		return
	}

	// Update the password
	err = database.UpdatePassword(userID.(string), req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update password: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password updated successfully",
	})
}

// UpdateUserHandler handles user profile updates
func UpdateUserHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// Get current user to check permissions
	currentUser, err := database.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve user: " + err.Error(),
		})
		return
	}

	// Prepare update fields
	updates := make(map[string]interface{})

	// Always allow name update
	if req.Name != "" {
		updates["name"] = req.Name
	}

	// Email updates are restricted based on system config
	// Check if email updates are enabled for non-owner users
	emailUpdatesEnabled := false
	if currentUser.IsOwner {
		emailUpdatesEnabled = true // Owners can always update email
	} else {
		// For non-owners, check system config
		var config database.SystemConfig
		err := database.DB.Where("key = ?", "feature_EMAIL_UPDATES_ENABLED").First(&config).Error
		if err == nil && config.Value == "true" {
			emailUpdatesEnabled = true
		}
	}

	if req.Email != "" && req.Email != currentUser.Email {
		if !emailUpdatesEnabled {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Email updates are disabled by system administrator",
			})
			return
		}

		// Check if email is already in use by another user
		existingUser, _ := database.GetUserByEmail(req.Email)
		if existingUser != nil && existingUser.ID != userID.(string) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Email address is already in use",
			})
			return
		}

		updates["email"] = req.Email
	}

	// If there's nothing to update
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No fields to update",
		})
		return
	}

	// Update the user
	err = database.UpdateUserFields(userID.(string), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user: " + err.Error(),
		})
		return
	}

	// Get updated user
	updatedUser, err := database.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve updated user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated successfully",
		"data": gin.H{
			"id":         updatedUser.ID,
			"email":      updatedUser.Email,
			"name":       updatedUser.Name,
			"is_owner":   updatedUser.IsOwner,
			"is_enabled": updatedUser.IsEnabled,
		},
	})
}
