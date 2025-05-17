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

// AdminRoleMiddleware verifies that the authenticated user is a system owner
func AdminRoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Directly query database to check if user is an owner
		var user database.User
		if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil || !user.IsOwner {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Access denied: Admin role required",
			})
			c.Abort()
			return
		}

		// Set the isOwner flag in the context for downstream handlers
		c.Set("isOwner", true)
		c.Next()
	}
}

// WorkspaceAdminCheck verifies that the user is an admin in a specific workspace
func WorkspaceAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This middleware expects the workspaceID to be in the URL parameters
		workspaceID := c.Param("workspaceID")
		if workspaceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Workspace ID is required",
			})
			c.Abort()
			return
		}

		// Get authenticated user ID from context (set by JWTAuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Check if user is a system-wide admin (can access all workspaces)
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Invalid user ID format",
			})
			c.Abort()
			return
		}

		// Directly query database to check if user is an owner
		var user database.User
		if err := database.DB.Where("id = ?", userIDStr).First(&user).Error; err == nil && user.IsOwner {
			// System owners can access all workspaces
			c.Set("isOwner", true)
			c.Next()
			return
		}

		// Get the database instance from context or import it directly
		// For now, we'll use a direct import to check workspace admin status
		isAdmin, err := database.IsUserWorkspaceAdmin(userID.(string), workspaceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to verify workspace access",
			})
			c.Abort()
			return
		}

		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "You do not have admin privileges in this workspace",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
