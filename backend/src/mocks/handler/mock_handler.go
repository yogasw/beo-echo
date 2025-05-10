package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// MockRequestHandler is a catch-all handler for mock API endpoints
//
// Sample curl:
// curl -X GET "http://localhost:8000/mock/myproject/api/users"
// curl -X POST "http://myproject.localhost:8000/api/users" -H "Content-Type: application/json" -d '{"name":"John Doe"}'
func MockRequestHandler(c *gin.Context) {
	// Check if mock service is initialized
	EnsureMockService()
	if mockService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Mock service is not available",
		})
		return
	}

	// Get project name from subdomain or path
	projectName := extractProjectName(c.Request)
	if projectName == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not specified",
		})
		return
	}

	// Get path (remove /mock prefix if present)
	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/mock") {
		path = path[5:]
	}
	if path == "" {
		path = "/"
	}

	// Process the request
	resp, err := mockService.HandleRequest(projectName, c.Request.Method, path, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Error processing mock request: " + err.Error(),
		})
		return
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Send response with proper status code
	c.Status(resp.StatusCode)

	// Copy response body
	if resp.Body != nil {
		defer resp.Body.Close()
		// Copy body to response writer
		if body, err := io.ReadAll(resp.Body); err == nil {
			c.Writer.Write(body)
		}
	}
}

// extractProjectName extracts project name from request (subdomain or path)
func extractProjectName(req *http.Request) string {
	// Try to extract from Host header (subdomain)
	host := req.Host

	// Check for subdomain
	parts := strings.Split(host, ".")
	if len(parts) > 2 {
		return parts[0]
	}

	// Try to extract from path
	path := req.URL.Path

	// If path starts with /mock, remove it
	if strings.HasPrefix(path, "/mock") {
		path = path[5:]
	}

	// Extract first part of path
	parts = strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}

	// Default project
	return "default"
}
