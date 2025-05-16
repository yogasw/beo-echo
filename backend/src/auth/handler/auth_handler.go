package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mockoon-control-panel/backend_new/src/auth"
	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/utils"
)

// LoginRequest represents the login form data
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the registration form data
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginHandler handles user authentication and returns a JWT token
func LoginHandler(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// Find the user by email
	user, err := database.GetUserByEmail(request.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid email or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to authenticate: " + err.Error(),
		})
		return
	}

	// Verify the password
	if !user.VerifyPassword(request.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate token: " + err.Error(),
		})
		return
	}

	// Return the token and user info
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"is_owner":   user.IsOwner,
			"is_enabled": user.IsEnabled,
		},
	})
}

// RegisterHandler handles user registration
func RegisterHandler(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}
	//check feature flag
	enable, err := utils.GetConfig[bool](utils.FEATURE_REGISTER_EMAIL_ENABLED)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to check feature flag: " + err.Error(),
		})
		return
	}
	if enable == false {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Registration is disabled",
		})
		return
	}

	// Check if a user with this email already exists
	existingUser, _ := database.GetUserByEmail(request.Email)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email is already registered",
		})
		return
	}

	// Hash the password
	hashedPassword, err := database.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to process registration: " + err.Error(),
		})
		return
	}

	// Create the user within a transaction
	tx := database.DB.Begin()

	user := database.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
		IsOwner:  false, // Default non-admin account
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create user: " + err.Error(),
		})
		return
	}

	// Create a workspace with the user's name
	workspace := database.Workspace{
		Name: request.Name + "'s Workspace",
	}

	if err := tx.Create(&workspace).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create workspace: " + err.Error(),
		})
		return
	}

	// Add user to workspace as admin
	userWorkspace := database.UserWorkspace{
		UserID:      user.ID,
		WorkspaceID: workspace.ID,
		Role:        "admin", // User is admin of their own workspace
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
			"message": "Failed to complete registration: " + err.Error(),
		})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Registration successful but failed to generate token: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registration successful",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}
