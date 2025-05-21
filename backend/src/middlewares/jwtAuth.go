package middlewares

import (
	"net/http"
	"strings"

	"beo-echo/backend/src/auth"
	"beo-echo/backend/src/database"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware authenticates requests using JWT
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			authHeader = c.Query("auth")
		}

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if the header format is valid
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization format, should be 'Bearer {token}'",
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := parts[1]

		// Validate the token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user info in the context for handlers to use
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)

		// Get user from database to check if they're an owner
		var user database.User
		if err := database.DB.Where("id = ?", claims.UserID).First(&user).Error; err == nil {
			c.Set("isOwner", user.IsOwner)
		} else {
			c.Set("isOwner", false)
		}

		c.Next()
	}
}
