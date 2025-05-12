package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mockoon-control-panel/backend_new/src/database"
)

// GetUserWorkspacesHandler returns all workspaces accessible to the authenticated user
func GetUserWorkspacesHandler(c *gin.Context) {
	// Get user ID from context (set by JWTAuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Fetch user's workspaces
	workspaces, err := database.GetUserWorkspaces(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch workspaces: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    workspaces,
	})
}

// CreateWorkspaceHandler creates a new workspace
func CreateWorkspaceHandler(c *gin.Context) {
	// Only system admins or authorized users can create workspaces
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Parse request
	var request struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// Create workspace in a transaction
	tx := database.DB.Begin()

	// Create the workspace
	workspace := database.Workspace{
		Name: request.Name,
	}

	if err := tx.Create(&workspace).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create workspace: " + err.Error(),
		})
		return
	}

	// Add the current user as an admin of this workspace
	userWorkspace := database.UserWorkspace{
		UserID:      userID.(string),
		WorkspaceID: workspace.ID,
		Role:        "admin",
	}

	if err := tx.Create(&userWorkspace).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to assign user to workspace: " + err.Error(),
		})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to complete workspace creation: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Workspace created successfully",
		"data":    workspace,
	})
}

// CheckWorkspaceRoleHandler returns the user's role in a specific workspace
func CheckWorkspaceRoleHandler(c *gin.Context) {
	workspaceID := c.Param("workspaceID")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Workspace ID is required",
		})
		return
	}

	// Get user ID from request (this allows admins to check other users' roles)
	requestedUserID := c.Query("user_id")
	if requestedUserID == "" {
		// If no specific user ID is provided, use the authenticated user's ID
		userIDValue, exists := c.Get("userID")
		if !exists || userIDValue == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User not authenticated",
			})
			return
		}
		requestedUserID = userIDValue.(string)
	} else {
		// If checking another user's role, ensure the requesting user is an admin
		userIDValue, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User not authenticated",
			})
			return
		}

		// Check if the authenticated user is a system admin or workspace admin
		authenticatedUserID := userIDValue.(string)

		// Directly query database to check if user is an owner
		var user database.User
		err := database.DB.Where("id = ?", authenticatedUserID).First(&user).Error
		isSystemOwner := err == nil && user.IsOwner

		if !isSystemOwner {
			// Not a system admin, check if workspace admin
			isAdmin, err := database.IsUserWorkspaceAdmin(userIDValue.(string), workspaceID)
			if err != nil || !isAdmin {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "You do not have admin privileges to check other users' roles",
				})
				return
			}
		}
	}

	// Now check the requested user's role
	var userWorkspace database.UserWorkspace
	err := database.DB.Where("user_id = ? AND workspace_id = ?", requestedUserID, workspaceID).First(&userWorkspace).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User is not a member of this workspace",
				"data": gin.H{
					"role": "none",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to verify workspace role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"user_id":      userWorkspace.UserID,
			"workspace_id": userWorkspace.WorkspaceID,
			"role":         userWorkspace.Role,
		},
	})
}
