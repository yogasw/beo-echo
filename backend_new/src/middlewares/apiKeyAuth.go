package middlewares

import (
	"encoding/base64"
	"net/http"
	"strings"

	"mockoon-control-panel/backend_new/src/lib"

	"github.com/gin-gonic/gin"
)

// ApiKeyAuth middleware for authenticating API requests using Basic Auth
func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")

		// Check if header exists and starts with "Basic "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid authorization header",
			})
			c.Abort()
			return
		}

		// Extract the base64 encoded credentials
		base64Credentials := strings.TrimPrefix(authHeader, "Basic ")
		credentials, err := base64.StdEncoding.DecodeString(base64Credentials)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials format",
			})
			c.Abort()
			return
		}

		// Split the credentials into username and password
		credentialsStr := string(credentials)
		parts := strings.SplitN(credentialsStr, ":", 2)
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials format",
			})
			c.Abort()
			return
		}

		username, password := parts[0], parts[1]

		// Get expected credentials from environment
		expectedCreds := strings.SplitN(lib.API_KEY, ":", 2)
		if len(expectedCreds) != 2 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Server configuration error",
			})
			c.Abort()
			return
		}

		expectedUsername, expectedPassword := expectedCreds[0], expectedCreds[1]

		// Validate credentials
		if username != expectedUsername || password != expectedPassword {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			c.Abort()
			return
		}

		// Authentication successful, continue to the next middleware/handler
		c.Next()
	}
}
