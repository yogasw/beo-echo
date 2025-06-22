package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// rateLimitState holds request count and window start time
type rateLimitState struct {
	count     int
	lastReset int64
	lastUsed  int64 // Unix timestamp of last usage for cleanup
}

var (
	ipRateStates     sync.Map
	apiRateLimit     = 100        // 100 requests for /api endpoints
	generalRateLimit = 200        // 200 requests for other endpoints
	rateWindow       = int64(60)  // per 60 seconds
	stateTimeout     = int64(600) // cleanup IPs idle for 10 minutes
)

// RateLimitByIP applies rate limiting per IP address
// 100 requests per minute for /api endpoints, 200 for others
func RateLimitByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		now := time.Now().Unix()

		// Determine rate limit based on path
		limit := generalRateLimit
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			limit = apiRateLimit
		}

		// Load or initialize rate state
		val, _ := ipRateStates.LoadOrStore(ip, &rateLimitState{
			count:     0,
			lastReset: now,
			lastUsed:  now,
		})
		state := val.(*rateLimitState)

		// Update last used time
		state.lastUsed = now

		// Reset count if window expired
		if now-state.lastReset >= rateWindow {
			state.count = 0
			state.lastReset = now
		}

		// Block if over limit
		if state.count >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Max " + fmt.Sprintf("%d", limit) + " requests per minute.",
			})
			c.Abort()
			return
		}

		// Allow and increment
		state.count++

		// Cleanup stale IPs periodically
		cleanupStaleIPStates()

		c.Next()
	}
}

// getClientIP extracts the client IP address
func getClientIP(c *gin.Context) string {
	// Try to get IP from various sources
	ip := c.ClientIP()
	if ip == "" {
		// Fallback to RemoteAddr if ClientIP() returns empty
		ip = c.Request.RemoteAddr
	}

	// Remove port if present
	if strings.Contains(ip, ":") {
		parts := strings.Split(ip, ":")
		return parts[0]
	}
	return ip
}

// CleanupStaleIPStates removes entries from ipRateStates that haven't been used within stateTimeout
// This function is exported for testing purposes
func CleanupStaleIPStates() {
	cleanupStaleIPStates()
}

// cleanupStaleIPStates removes entries from ipRateStates that haven't been used within stateTimeout
func cleanupStaleIPStates() {
	now := time.Now().Unix()
	ipRateStates.Range(func(key, value any) bool {
		state := value.(*rateLimitState)
		if now-state.lastUsed >= stateTimeout {
			ipRateStates.Delete(key)
		}
		return true
	})
}
