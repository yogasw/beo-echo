package users

import (
	"beo-echo/backend/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	service *UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetCurrentUser handles retrieving the authenticated user's information
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	user, featureFlags, err := h.service.GetCurrentUser(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve user information: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":            user.ID,
			"email":         user.Email,
			"name":          user.Name,
			"is_owner":      user.IsOwner,
			"is_active":     user.IsActive,
			"feature_flags": featureFlags,
		},
	})
}

// UpdatePasswordRequest represents the update password request
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// UpdateProfileRequest represents the update user profile request
type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdatePassword handles password update requests
func (h *UserHandler) UpdatePassword(c *gin.Context) {
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

	// Update the password
	err := h.service.UpdatePassword(c.Request.Context(), userID.(string), req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password updated successfully",
	})
}

// UpdateUser handles user profile updates
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// Current user data to determine if they are an owner
	userData, _, err := h.service.GetCurrentUser(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve user: " + err.Error(),
		})
		return
	}

	// Update the user
	err = h.service.UpdateUserProfile(c.Request.Context(), userID.(string), req.Name, req.Email, userData.IsOwner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Get updated user
	updatedUser, _, err := h.service.GetCurrentUser(c.Request.Context(), userID.(string))
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
			"id":        updatedUser.ID,
			"email":     updatedUser.Email,
			"name":      updatedUser.Name,
			"is_owner":  updatedUser.IsOwner,
			"is_active": updatedUser.IsActive,
		},
	})
}

// GetAllUsers handles retrieving all users (admin/owner only)
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Check if user is an owner
	currentUser, _, err := h.service.GetCurrentUser(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve user: " + err.Error(),
		})
		return
	}

	// Only owners can see all users
	if !currentUser.IsOwner {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Permission denied: Only instance owners can view all users",
		})
		return
	}

	users, err := h.service.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve users: " + err.Error(),
		})
		return
	}

	// Transform users for response
	var usersList []gin.H
	for _, user := range users {
		usersList = append(usersList, gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"is_owner":   user.IsOwner,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    usersList,
	})
}

// GetWorkspaceUsersRequest represents the request to get users in a workspace
type WorkspaceUserRequest struct {
	WorkspaceID string `uri:"workspace_id" binding:"required"`
}

// GetWorkspaceUsers handles retrieving all users in a workspace
func (h *UserHandler) GetWorkspaceUsers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	var req WorkspaceUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid workspace ID",
		})
		return
	}

	// Check if user has access to this workspace
	var userWorkspace database.UserWorkspace
	if err := database.DB.Where("user_id = ? AND workspace_id = ?", userID, req.WorkspaceID).First(&userWorkspace).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "You do not have access to this workspace",
		})
		return
	}

	users, err := h.service.GetWorkspaceUsers(c.Request.Context(), req.WorkspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve workspace users: " + err.Error(),
		})
		return
	}

	// Transform users for response
	var usersList []gin.H
	for _, user := range users {
		// Find user's role in this workspace
		var role string = "member"
		for _, workspace := range user.Workspaces {
			if workspace.WorkspaceID == req.WorkspaceID {
				role = workspace.Role
				break
			}
		}

		usersList = append(usersList, gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"name":     user.Name,
			"role":     role,
			"is_owner": user.IsOwner,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    usersList,
	})
}

// RemoveWorkspaceUserRequest represents the request to remove a user from a workspace
type RemoveWorkspaceUserRequest struct {
	WorkspaceID string `uri:"workspaceID" binding:"required"`
	UserID      string `uri:"user_id" binding:"required"`
}

// RemoveWorkspaceUser handles removing a user from a workspace
func (h *UserHandler) RemoveWorkspaceUser(c *gin.Context) {
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	var req RemoveWorkspaceUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request parameters",
		})
		return
	}

	// TODO: Check if current user is admin in this workspace
	// This would typically be handled by middleware

	// Prevent users from removing themselves
	if currentUserID.(string) == req.UserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Cannot remove yourself from workspace",
		})
		return
	}

	err := h.service.RemoveUserFromWorkspace(c.Request.Context(), req.WorkspaceID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to remove user from workspace: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User removed from workspace successfully",
	})
}

// DeleteUserRequest represents the request to delete a user
type DeleteUserRequest struct {
	UserID string `uri:"user_id" binding:"required"`
}

// DeleteUser handles completely removing a user from the system
func (h *UserHandler) DeleteUser(c *gin.Context) {
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	var req DeleteUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Prevent self-deletion
	if currentUserID.(string) == req.UserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Cannot delete your own account",
		})
		return
	}

	err := h.service.DeleteUser(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

// UpdateWorkspaceUserRoleRequest represents the request to update a user's role in a workspace
type UpdateWorkspaceUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member"` // Only admin or member roles are allowed
}

// UpdateWorkspaceUserRole handles updating a user's role in a workspace
func (h *UserHandler) UpdateWorkspaceUserRole(c *gin.Context) {
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	var pathParams RemoveWorkspaceUserRequest
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request parameters",
		})
		return
	}

	var req UpdateWorkspaceUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// TODO: Check if current user is admin in this workspace
	// This would typically be handled by middleware

	// If trying to update own role, prevent demoting self if last admin
	if currentUserID.(string) == pathParams.UserID && req.Role == "member" {
		// This check is also performed in the service, but we provide a clearer message here
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Cannot demote yourself if you are the last admin",
		})
		return
	}

	err := h.service.UpdateUserWorkspaceRole(c.Request.Context(), pathParams.WorkspaceID, pathParams.UserID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User role updated successfully",
	})
}

// UpdateUserOwnerRequest represents the request to update a user's owner status
type UpdateUserRequest struct {
	IsOwner  bool `json:"is_owner" omitempty:"true"`
	IsActive bool `json:"is_active" omitempty:"true"`
}

// UpdateUser handles updating a user's owner status
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("user_id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}
	updates := make(map[string]interface{})
	if req.IsOwner {
		updates["is_owner"] = req.IsOwner
	}
	if req.IsActive {
		updates["is_active"] = req.IsActive
	}

	if err := h.service.UpdateUserFields(c.Request.Context(), userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User owner status updated successfully",
	})
}
