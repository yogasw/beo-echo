package middlewares

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitConstants(t *testing.T) {
	t.Run("verifies rate limit constants are unchanged", func(t *testing.T) {
		// Test that the production rate limits are set correctly
		// These values should not change without explicit approval

		// Verify API rate limit
		assert.Equal(t, 100, apiRateLimit, "API rate limit should be 100 requests per minute")

		// Verify general rate limit
		assert.Equal(t, 200, generalRateLimit, "General rate limit should be 200 requests per minute")

		// Verify rate window
		assert.Equal(t, int64(60), rateWindow, "Rate window should be 60 seconds")

		// Verify state timeout
		assert.Equal(t, int64(600), stateTimeout, "State timeout should be 600 seconds (10 minutes)")
	})

	t.Run("verifies rate limit behavior with production values", func(t *testing.T) {
		// Setup test router without modifying global constants
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.Use(RateLimitByIP())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "api success"})
		})

		// Clear state for clean test
		ipRateStates = sync.Map{}

		// Test that we can make requests up to the production limits
		// For API endpoints (100 requests)
		testIP := "192.168.100.1:12345"

		// Make 99 requests to API endpoint (should all succeed)
		for i := 0; i < 99; i++ {
			req := httptest.NewRequest("GET", "/api/test", nil)
			req.RemoteAddr = testIP

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("Request %d failed with status %d, expected %d", i+1, w.Code, http.StatusOK)
			}
		}

		// 100th request should still succeed
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = testIP
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "100th API request should succeed")

		// 101st request should be blocked
		req = httptest.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = testIP
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code, "101st API request should be blocked")
	})

	t.Run("verifies general endpoint rate limit", func(t *testing.T) {
		// Setup test router
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.Use(RateLimitByIP())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Clear state for clean test
		ipRateStates = sync.Map{}

		// Test general endpoints (200 requests)
		testIP := "192.168.101.1:12345"

		// Make 199 requests to general endpoint (should all succeed)
		for i := 0; i < 199; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = testIP

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("Request %d failed with status %d, expected %d", i+1, w.Code, http.StatusOK)
			}
		}

		// 200th request should still succeed
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = testIP
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "200th general request should succeed")

		// 201st request should be blocked
		req = httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = testIP
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code, "201st general request should be blocked")
	})
}

func TestRateLimitByIP(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Store original values to restore later
	originalRateWindow := rateWindow
	originalStateTimeout := stateTimeout
	originalApiRateLimit := apiRateLimit
	originalGeneralRateLimit := generalRateLimit

	// Set shorter intervals for testing
	rateWindow = 2        // 2 seconds instead of 60
	stateTimeout = 5      // 5 seconds instead of 600
	apiRateLimit = 5      // 5 requests instead of 100
	generalRateLimit = 10 // 10 requests instead of 200

	// Restore original values after test
	defer func() {
		rateWindow = originalRateWindow
		stateTimeout = originalStateTimeout
		apiRateLimit = originalApiRateLimit
		generalRateLimit = originalGeneralRateLimit
	}()

	// Clear any existing state
	ipRateStates = sync.Map{}

	// Create router with rate limit middleware
	router := gin.New()
	router.Use(RateLimitByIP())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "api success"})
	})

	t.Run("allows requests under limit", func(t *testing.T) {
		// Clear state
		ipRateStates = sync.Map{}

		// Make 8 requests (under limit of 10 for general endpoints)
		for i := 0; i < 8; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:12345"

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("blocks requests over limit", func(t *testing.T) {
		// Clear state
		ipRateStates = sync.Map{}

		// Make 10 requests (at limit for general endpoints)
		for i := 0; i < 10; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.2:12345"

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}

		// 11th request should be blocked
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.2:12345"

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("resets counter after window expires", func(t *testing.T) {
		// Clear state
		ipRateStates = sync.Map{}

		// Override rate window for testing (1 second)
		originalWindow := rateWindow
		rateWindow = 1
		defer func() { rateWindow = originalWindow }()

		// Make 10 requests (at limit for general endpoints)
		for i := 0; i < 10; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.3:12345"

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}

		// Next request should be blocked
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.3:12345"

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)

		// Wait for window to expire
		time.Sleep(1100 * time.Millisecond)

		// Should be allowed again
		req = httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.3:12345"

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("handles different IPs separately", func(t *testing.T) {
		// Clear state
		ipRateStates = sync.Map{}

		// IP1 makes 10 requests (at limit for general endpoints)
		for i := 0; i < 10; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.4:12345"

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}

		// IP1 should be blocked
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.4:12345"

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)

		// IP2 should still be allowed
		req = httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.5:12345"

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("applies different limits for API endpoints", func(t *testing.T) {
		// Clear state
		ipRateStates = sync.Map{}

		// Make 5 requests to API endpoint (at limit for API endpoints)
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest("GET", "/api/test", nil)
			req.RemoteAddr = "192.168.1.6:12345"

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}

		// 6th request should be blocked
		req := httptest.NewRequest("GET", "/api/test", nil)
		req.RemoteAddr = "192.168.1.6:12345"

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("cleanup removes stale states", func(t *testing.T) {
		// Clear state
		ipRateStates = sync.Map{}

		// Override stateTimeout for testing (1 second)
		originalTimeout := stateTimeout
		stateTimeout = 1
		defer func() { stateTimeout = originalTimeout }()

		// Manually create a state entry (simulating what the middleware would do)
		now := time.Now().Unix()
		testIP := "192.168.1.7"

		ipRateStates.Store(testIP, &rateLimitState{
			count:     1,
			lastReset: now,
			lastUsed:  now - 2, // Make it 2 seconds old (stale)
		})

		// Verify state exists
		_, exists := ipRateStates.Load(testIP)
		assert.True(t, exists)

		// Call cleanup function
		CleanupStaleIPStates()

		// State should be cleaned up since it's older than stateTimeout
		_, exists = ipRateStates.Load(testIP)
		assert.False(t, exists)
	})
}

func TestGetClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("extracts IP from RemoteAddr", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.RemoteAddr = "192.168.1.1:12345"

		ip := getClientIP(c)
		assert.Equal(t, "192.168.1.1", ip)
	})

	t.Run("handles IP without port", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.RemoteAddr = "192.168.1.1"

		ip := getClientIP(c)
		assert.Equal(t, "192.168.1.1", ip)
	})
}
