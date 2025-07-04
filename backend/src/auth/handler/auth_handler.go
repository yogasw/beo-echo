package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"beo-echo/backend/src/auth"
	"beo-echo/backend/src/auth/services"
	"beo-echo/backend/src/users"
)

// Global auth service instance
var authService *services.AuthService
var userService *users.UserService

// InitAuthService initializes the auth service
func InitAuthService(db *gorm.DB, userRepo users.UserRepository) {
	authService = services.NewAuthService(userRepo, db)
	userService = users.NewUserService(userRepo)
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

	// Generate JWT access token (15 minutes)
	token, err := auth.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate token: " + err.Error(),
		})
		return
	}

	// Generate refresh token (30 days) and save to database
	refreshToken, err := userService.SaveRefreshToken(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate refresh token: " + err.Error(),
		})
		return
	}

	// Set refresh token in HTTP-only cookie (30 days)
	c.SetCookie(
		"refresh_token", // name
		refreshToken,    // value
		30*24*60*60,     // maxAge (30 days in seconds)
		"/",             // path
		"",              // domain (empty = same domain)
		true,            // secure (true for HTTPS)
		true,            // httpOnly
	)

	// Return the access token and user info
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

// RefreshTokenHandler handles refresh token requests and returns a new access token
func RefreshTokenHandler(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "No refresh token found",
		})
		return
	}

	// Validate refresh token and get user
	user, err := userService.ValidateRefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		// Clear invalid refresh token cookie
		c.SetCookie("refresh_token", "", -1, "/", "", true, true)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid or expired refresh token",
		})
		return
	}

	// Generate new access token (15 minutes)
	newToken, err := auth.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate new token: " + err.Error(),
		})
		return
	}

	// Check if refresh token should be rotated (every 7 days)
	shouldRotate := auth.ShouldRotateRefreshToken(refreshToken)

	if shouldRotate {
		// Generate new refresh token for rotation (only after 7 days)
		newRefreshToken, err := userService.SaveRefreshToken(c.Request.Context(), user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to generate new refresh token: " + err.Error(),
			})
			return
		}

		// Update refresh token cookie with new token
		c.SetCookie(
			"refresh_token",
			newRefreshToken,
			30*24*60*60, // 30 days
			"/",
			"",
			true, // secure
			true, // httpOnly
		)
	}

	// Return new access token
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token refreshed successfully",
		"token":   newToken,
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"is_active": user.IsActive,
		},
	})
}

// LogoutHandler handles user logout and clears refresh token
func LogoutHandler(c *gin.Context) {
	// Get user ID from JWT (middleware should have already validated this)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Clear refresh token from database
	err := userService.ClearRefreshToken(c.Request.Context(), userID.(string))
	if err != nil {
		// Log error but don't fail the logout
		c.Header("X-Warning", "Failed to clear refresh token from database")
	}

	// Clear refresh token cookie
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}
