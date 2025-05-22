package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"beo-echo/backend/src/auth"
	"beo-echo/backend/src/auth/services"
)

// Global auth service instance
var authService *services.AuthService

// InitAuthService initializes the auth service
func InitAuthService(db *gorm.DB, userRepo services.UserRepository) {
	authService = services.NewAuthService(userRepo, db)
}

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

	// Find the user by email - using direct DB method for now
	// TODO: Update to use context-based method in future
	user, err := authService.GetUserByEmailDirect(request.Email)
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
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.Name,
			"is_owner":  user.IsOwner,
			"is_active": user.IsActive,
		},
	})
}
