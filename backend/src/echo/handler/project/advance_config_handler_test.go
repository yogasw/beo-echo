package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
)

// Helper function to generate unique alias for tests
func generateUniqueAliasAdvance(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano()/int64(time.Millisecond))
}

func TestGetProjectAdvanceConfigHandler(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Initialize database
	err := database.CheckAndHandle()
	require.NoError(t, err, "Failed to setup test database")

	t.Run("Get Project Advance Config - Empty Config", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test1@example.com", "Test User 1", "Test Workspace 1")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project without advance config
		project, err := database.CreateTestProject(workspace.ID, "Test Project 1", generateUniqueAliasAdvance("test-project-advance-1"))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.GET("/api/projects/:projectId/advance-config", GetProjectAdvanceConfigHandler)

		// Create HTTP request
		req, err := http.NewRequest("GET", "/api/projects/"+project.ID+"/advance-config", nil)
		require.NoError(t, err)

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, project.ID, data["project_id"])
		assert.Equal(t, "Test Project 1", data["project_name"])
		assert.Nil(t, data["advance_config"]) // Should be null for empty config
	})

	t.Run("Get Project Advance Config - With Config", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test2@example.com", "Test User 2", "Test Workspace 2")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create advance config
		advanceConfig := map[string]interface{}{
			"global_timeout": 5000,
			"rate_limit": map[string]interface{}{
				"enabled":          true,
				"requests_per_min": 100,
				"burst_size":       10,
			},
			"cors": map[string]interface{}{
				"enabled":         true,
				"allowed_origins": []string{"https://example.com"},
			},
		}

		advanceConfigJSON, err := json.Marshal(advanceConfig)
		require.NoError(t, err)

		// Create test project with advance config
		project, err := database.CreateTestProjectWithConfig(workspace.ID, "Test Project 2", generateUniqueAliasAdvance("test-project-advance-2"), string(advanceConfigJSON))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.GET("/api/projects/:projectId/advance-config", GetProjectAdvanceConfigHandler)

		// Create HTTP request
		req, err := http.NewRequest("GET", "/api/projects/"+project.ID+"/advance-config", nil)
		require.NoError(t, err)

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, project.ID, data["project_id"])
		assert.Equal(t, "Test Project 2", data["project_name"])

		// Verify advance config structure
		config := data["advance_config"].(map[string]interface{})
		assert.Equal(t, float64(5000), config["global_timeout"])

		rateLimitConfig := config["rate_limit"].(map[string]interface{})
		assert.True(t, rateLimitConfig["enabled"].(bool))
		assert.Equal(t, float64(100), rateLimitConfig["requests_per_min"])
	})

	t.Run("Get Project Advance Config - Project Not Found", func(t *testing.T) {
		// Setup Gin router
		router := gin.New()
		router.GET("/api/projects/:projectId/advance-config", GetProjectAdvanceConfigHandler)

		// Create HTTP request with non-existent project ID
		req, err := http.NewRequest("GET", "/api/projects/non-existent-id/advance-config", nil)
		require.NoError(t, err)

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Equal(t, "Project not found", response["message"])
	})
}

func TestUpdateProjectAdvanceConfigHandler(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Initialize database
	err := database.CheckAndHandle()
	require.NoError(t, err, "Failed to setup test database")

	t.Run("Update Project Advance Config - Success", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test3@example.com", "Test User 3", "Test Workspace 3")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project
		project, err := database.CreateTestProject(workspace.ID, "Test Project 3", generateUniqueAliasAdvance("test-project-advance-3"))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId/advance-config", UpdateProjectAdvanceConfigHandler)

		// Create new advance config
		newConfig := map[string]interface{}{
			"global_timeout": 10000,
			"rate_limit": map[string]interface{}{
				"enabled":          true,
				"requests_per_min": 200,
				"burst_size":       20,
			},
			"cors": map[string]interface{}{
				"enabled":         true,
				"allowed_origins": []string{"https://example.com", "https://app.example.com"},
				"allowed_methods": []string{"GET", "POST", "PUT", "DELETE"},
			},
			"security": map[string]interface{}{
				"enable_https": true,
				"enable_hsts":  true,
			},
		}

		jsonData, err := json.Marshal(newConfig)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID+"/advance-config", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Project advance config updated successfully", response["message"])

		// Verify response data
		data := response["data"].(map[string]interface{})
		assert.Equal(t, project.ID, data["project_id"])
		assert.Equal(t, "Test Project 3", data["project_name"])

		// Verify advance config structure
		config := data["advance_config"].(map[string]interface{})
		assert.Equal(t, float64(10000), config["global_timeout"])

		rateLimitConfig := config["rate_limit"].(map[string]interface{})
		assert.True(t, rateLimitConfig["enabled"].(bool))
		assert.Equal(t, float64(200), rateLimitConfig["requests_per_min"])
		assert.Equal(t, float64(20), rateLimitConfig["burst_size"])

		corsConfig := config["cors"].(map[string]interface{})
		assert.True(t, corsConfig["enabled"].(bool))
		allowedOrigins := corsConfig["allowed_origins"].([]interface{})
		assert.Len(t, allowedOrigins, 2)
		assert.Equal(t, "https://example.com", allowedOrigins[0])
		assert.Equal(t, "https://app.example.com", allowedOrigins[1])

		securityConfig := config["security"].(map[string]interface{})
		assert.True(t, securityConfig["enable_https"].(bool))
		assert.True(t, securityConfig["enable_hsts"].(bool))
	})

	t.Run("Update Project Advance Config - Invalid JSON", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test4@example.com", "Test User 4", "Test Workspace 4")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project
		project, err := database.CreateTestProject(workspace.ID, "Test Project 4", generateUniqueAliasAdvance("test-project-advance-4"))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId/advance-config", UpdateProjectAdvanceConfigHandler)

		// Create invalid JSON data
		invalidJSON := `{"global_timeout": 5000, "invalid":}`

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID+"/advance-config", bytes.NewBufferString(invalidJSON))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"].(string), "Invalid JSON data")
	})

	t.Run("Update Project Advance Config - Invalid DelayMs Too High", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test_advance_delayms_high@example.com", "Test User Advance DelayMs High", "Test Workspace Advance DelayMs High")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project
		project, err := database.CreateTestProject(workspace.ID, "Test Project Advance DelayMs High", generateUniqueAliasAdvance("test-project-advance-delayms-high"))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId/advance-config", UpdateProjectAdvanceConfigHandler)

		// Invalid config with delayMs too high
		invalidConfig := map[string]interface{}{
			"delayMs": 400000, // Above maximum 120000ms
		}

		jsonData, err := json.Marshal(invalidConfig)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID+"/advance-config", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"].(string), "delayMs cannot exceed 120000ms (2 minutes)")
	})

	t.Run("Update Project Advance Config - Valid DelayMs", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test_advance_delayms_valid@example.com", "Test User Advance DelayMs Valid", "Test Workspace Advance DelayMs Valid")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project
		project, err := database.CreateTestProject(workspace.ID, "Test Project Advance DelayMs Valid", generateUniqueAliasAdvance("test-project-advance-delayms-valid"))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId/advance-config", UpdateProjectAdvanceConfigHandler)

		// Valid config with proper delayMs
		validConfig := map[string]interface{}{
			"delayMs": 15000, // Valid 15 seconds delay
		}

		jsonData, err := json.Marshal(validConfig)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID+"/advance-config", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Project advance config updated successfully", response["message"])

		// Verify config in database
		var updatedProject database.Project
		result := database.GetDB().Where("id = ?", project.ID).First(&updatedProject)
		require.NoError(t, result.Error)
		assert.Equal(t, `{"delayMs":15000}`, updatedProject.AdvanceConfig)
	})

}
