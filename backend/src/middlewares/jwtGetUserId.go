package middlewares

import (
	"strings"

	"beo-echo/backend/src/auth"

	"github.com/gin-gonic/gin"
)

// JWTGetUserIdMiddleware extracts the user ID from the JWT token
func JWTGetUserIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			authHeader = c.Query("auth")
		}

		if authHeader == "" {
			c.Set("isAuthenticated", false)
			c.Next()
			return
		}

		// Check if the header format is valid
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Set("isAuthenticated", false)
			c.Next()
			return
		}

		// Extract the token
		tokenString := parts[1]

		// Validate the token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.Set("isAuthenticated", false)
			c.Next()
			return
		}

		// Set user info in the context for handlers to use
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("isAuthenticated", true)

		c.Next()
	}
}
