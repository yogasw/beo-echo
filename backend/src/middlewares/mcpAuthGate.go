package middlewares

import (
	"net/http"
	"strings"

	"beo-echo/backend/src/auth"
	"beo-echo/backend/src/auth/pat"
	"beo-echo/backend/src/database"

	"github.com/gin-gonic/gin"
)

// MCPAuthGate requires a valid bearer token (a Personal Access Token or a JWT)
// on every request to the MCP endpoint — including the MCP discovery handshake
// (initialize / tools-list). Requests without a valid token are rejected with
// 401 before they reach the MCP handler, so the tool catalogue is never exposed
// to unauthenticated callers.
//
// This is stricter than the per-tool token forwarding the MCP tools already do:
// it blocks the protocol itself, not just tool execution.
func MCPAuthGate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abortMCPUnauthorized(c, "Authorization bearer token is required")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			abortMCPUnauthorized(c, "Invalid authorization format, expected 'Bearer <token>'")
			return
		}
		token := strings.TrimSpace(parts[1])

		// PATs are resolved against the database; anything else is treated as a JWT.
		if pat.IsPAT(token) {
			if _, err := pat.NewService(database.DB).Authenticate(c.Request.Context(), token); err != nil {
				abortMCPUnauthorized(c, "Invalid or expired token")
				return
			}
			c.Next()
			return
		}

		if _, err := auth.ValidateToken(token); err != nil {
			abortMCPUnauthorized(c, "Invalid or expired token")
			return
		}
		c.Next()
	}
}

// abortMCPUnauthorized returns a 401 in a shape MCP clients can read.
func abortMCPUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"jsonrpc": "2.0",
		"error": gin.H{
			"code":    -32001,
			"message": message,
		},
	})
	c.Abort()
}
