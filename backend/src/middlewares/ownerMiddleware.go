package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OwnerOnlyMiddleware restricts access to routes that should only be accessible by system owners
func OwnerOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Check if the user is an owner - this flag is set by JWTAuthMiddleware
		isOwner, exists := c.Get("isOwner")
		if !exists || isOwner != true {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Access denied: Owner privileges required",
			})
			c.Abort()
			return
		}

		// User is an owner, proceed
		c.Next()
	}
}
