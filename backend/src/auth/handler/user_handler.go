package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	systemConfig "beo-echo/backend/src/systemConfigs"
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
	featureFlags, err := systemConfig.GetFeatureFlags()
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
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to retrieve user: " + err.Error(),
		})
		return
	}

	// Verify current password
	if !user.VerifyPassword(req.CurrentPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
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
		emailUpdatesEnabled, err = systemConfig.GetSystemConfigWithType[bool](string(systemConfig.FEATURE_EMAIL_UPDATES_ENABLED))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to retrieve system config: " + err.Error(),
			})
			return
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
