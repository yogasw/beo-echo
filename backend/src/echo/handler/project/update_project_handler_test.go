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
func generateUniqueAliasUpdate(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano()/int64(time.Millisecond))
}

func TestUpdateProjectHandler(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Initialize database
	err := database.CheckAndHandle()
	require.NoError(t, err, "Failed to setup test database")

	t.Run("Update Project Basic Fields", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test1@example.com", "Test User 1", "Test Workspace 1")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project
		originalAlias := generateUniqueAliasUpdate("original-project-update-1")
		project, err := database.CreateTestProject(workspace.ID, "Original Project", originalAlias)
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId", UpdateProjectHandler)

		// Create update data
		updatedAlias := generateUniqueAliasUpdate("updated-project-1")
		updateData := map[string]interface{}{
			"name":   "Updated Project Name",
			"alias":  updatedAlias,
			"mode":   "proxy",
			"status": "stopped",
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID, bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Debug: Print response if not success
		if w.Code != http.StatusOK {
			t.Logf("Response Code: %d", w.Code)
			t.Logf("Response Body: %s", w.Body.String())
		}

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Check if response indicates success
		if success, ok := response["success"].(bool); ok && success {
			assert.True(t, response["success"].(bool))
			assert.Equal(t, "Project updated successfully", response["message"])

			// Verify updated project data
			projectData := response["data"].(map[string]interface{})
			assert.Equal(t, "Updated Project Name", projectData["name"])
			assert.Equal(t, updatedAlias, projectData["alias"])
			assert.Equal(t, "proxy", projectData["mode"])
			assert.Equal(t, "stopped", projectData["status"])
		} else {
			t.Fatalf("Expected success response, got error: %v", response["message"])
		}
	})

	t.Run("Update Project Advance Config", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test2@example.com", "Test User 2", "Test Workspace 2")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project with initial advance config
		initialConfig := map[string]interface{}{
			"global_timeout": 3000,
		}
		initialConfigJSON, err := json.Marshal(initialConfig)
		require.NoError(t, err)

		project, err := database.CreateTestProjectWithConfig(workspace.ID, "Test Project 2", generateUniqueAliasUpdate("test-project-update-2"), string(initialConfigJSON))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId", UpdateProjectHandler)

		// Create advance config update
		newAdvanceConfig := map[string]interface{}{
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

		newAdvanceConfigJSON, err := json.Marshal(newAdvanceConfig)
		require.NoError(t, err)

		updateData := map[string]interface{}{
			"advance_config": string(newAdvanceConfigJSON),
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID, bytes.NewBuffer(jsonData))
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
		assert.Equal(t, "Project updated successfully", response["message"])

		// Verify updated advance config
		projectData := response["data"].(map[string]interface{})
		assert.NotEmpty(t, projectData["advance_config"])

		var storedAdvanceConfig map[string]interface{}
		err = json.Unmarshal([]byte(projectData["advance_config"].(string)), &storedAdvanceConfig)
		require.NoError(t, err)

		// Verify the new config values
		assert.Equal(t, float64(10000), storedAdvanceConfig["global_timeout"])

		rateLimitConfig := storedAdvanceConfig["rate_limit"].(map[string]interface{})
		assert.True(t, rateLimitConfig["enabled"].(bool))
		assert.Equal(t, float64(200), rateLimitConfig["requests_per_min"])
		assert.Equal(t, float64(20), rateLimitConfig["burst_size"])

		corsConfig := storedAdvanceConfig["cors"].(map[string]interface{})
		assert.True(t, corsConfig["enabled"].(bool))
		allowedOrigins := corsConfig["allowed_origins"].([]interface{})
		assert.Len(t, allowedOrigins, 2)
		assert.Equal(t, "https://example.com", allowedOrigins[0])
		assert.Equal(t, "https://app.example.com", allowedOrigins[1])
	})

	t.Run("Update Project With Invalid Advance Config JSON", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test3@example.com", "Test User 3", "Test Workspace 3")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project
		project, err := database.CreateTestProject(workspace.ID, "Test Project 3", generateUniqueAliasUpdate("test-project-update-3"))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId", UpdateProjectHandler)

		// Create update data with invalid JSON
		updateData := map[string]interface{}{
			"advance_config": `{"global_timeout": 5000, "invalid":}`, // Invalid JSON
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID, bytes.NewBuffer(jsonData))
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
		assert.Contains(t, response["message"].(string), "invalid JSON format in advance_config")
	})

	t.Run("Update Project Clear Advance Config", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test4@example.com", "Test User 4", "Test Workspace 4")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project with initial advance config
		initialConfig := map[string]interface{}{
			"global_timeout": 5000,
			"rate_limit": map[string]interface{}{
				"enabled":          true,
				"requests_per_min": 100,
			},
		}
		initialConfigJSON, err := json.Marshal(initialConfig)
		require.NoError(t, err)

		project, err := database.CreateTestProjectWithConfig(workspace.ID, "Test Project 4", generateUniqueAliasUpdate("test-project-update-4"), string(initialConfigJSON))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId", UpdateProjectHandler)

		// Clear advance config by setting it to empty string
		updateData := map[string]interface{}{
			"advance_config": "",
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID, bytes.NewBuffer(jsonData))
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

		// Verify advance config is cleared
		projectData := response["data"].(map[string]interface{})
		assert.Equal(t, "", projectData["advance_config"])
	})

	t.Run("Update Project With Invalid Status", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test5@example.com", "Test User 5", "Test Workspace 5")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project
		project, err := database.CreateTestProject(workspace.ID, "Test Project 5", generateUniqueAliasUpdate("test-project-update-5"))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId", UpdateProjectHandler)

		// Create update data with invalid status
		updateData := map[string]interface{}{
			"status": "invalid_status",
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID, bytes.NewBuffer(jsonData))
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
		assert.Contains(t, response["message"].(string), "Invalid status value")
	})

	t.Run("Update Project With Invalid Alias", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test6@example.com", "Test User 6", "Test Workspace 6")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project
		project, err := database.CreateTestProject(workspace.ID, "Test Project 6", generateUniqueAliasUpdate("test-project-update-6"))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId", UpdateProjectHandler)

		// Create update data with invalid alias
		updateData := map[string]interface{}{
			"alias": "invalid alias with spaces@#$",
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID, bytes.NewBuffer(jsonData))
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
		assert.Contains(t, response["message"].(string), "Invalid alias format")
	})

	t.Run("Update Non-Existent Project", func(t *testing.T) {
		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId", UpdateProjectHandler)

		// Create update data
		updateData := map[string]interface{}{
			"name": "Updated Name",
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// Create HTTP request with non-existent project ID
		req, err := http.NewRequest("PUT", "/api/projects/non-existent-id", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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

	t.Run("Update Project Partial Update", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test7@example.com", "Test User 7", "Test Workspace 7")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project with initial data
		initialConfig := map[string]interface{}{
			"global_timeout": 3000,
		}
		initialConfigJSON, err := json.Marshal(initialConfig)
		require.NoError(t, err)

		originalAlias := generateUniqueAliasUpdate("original-project-update-7")
		project, err := database.CreateTestProjectWithConfig(workspace.ID, "Original Project", originalAlias, string(initialConfigJSON))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId", UpdateProjectHandler)

		// Only update the name, leave other fields unchanged
		updateData := map[string]interface{}{
			"name": "Updated Project Name Only",
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID, bytes.NewBuffer(jsonData))
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

		// Verify only name was updated, other fields remain unchanged
		projectData := response["data"].(map[string]interface{})
		assert.Equal(t, "Updated Project Name Only", projectData["name"])
		assert.Equal(t, originalAlias, projectData["alias"]) // Should remain unchanged
		assert.Equal(t, "mock", projectData["mode"])         // Should remain unchanged
		assert.NotEmpty(t, projectData["advance_config"])    // Should remain unchanged

		// Verify advance config is still the same
		var storedAdvanceConfig map[string]interface{}
		err = json.Unmarshal([]byte(projectData["advance_config"].(string)), &storedAdvanceConfig)
		require.NoError(t, err)
		assert.Equal(t, float64(3000), storedAdvanceConfig["global_timeout"])
	})

	t.Run("Update Project With Invalid DelayMs Too Low", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test_update_timeout_low@example.com", "Test User Update DelayMs Low", "Test Workspace Update DelayMs Low")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project
		project, err := database.CreateTestProject(workspace.ID, "Test Project Update DelayMs Low", generateUniqueAliasUpdate("test-project-update-delayMs-low"))
		require.NoError(t, err)

		// Debug: Ensure project ID is not empty
		require.NotEmpty(t, project.ID, "Project ID should not be empty")
		t.Logf("Created project with ID: %s", project.ID)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId", UpdateProjectHandler)

		// Create update data with invalid delayMs (negative)
		updateData := map[string]interface{}{
			"advance_config": `{"delayMs": -1500}`, // Negative delayMs
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// Create HTTP request
		url := "/api/projects/" + project.ID
		t.Logf("Request URL: %s", url)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Debug: Print response if not expected
		if w.Code != http.StatusBadRequest {
			t.Logf("Response Code: %d", w.Code)
			t.Logf("Response Body: %s", w.Body.String())
		}

		// Verify response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Debug: Print actual message if assertion fails
		if !assert.True(t, response["error"].(bool)) {
			t.Logf("Actual response: %+v", response)
		}

		actualMessage := response["message"].(string)
		if !assert.Contains(t, actualMessage, "delayMs cannot be negative") {
			t.Logf("Actual message: %s", actualMessage)
		}
	})

	t.Run("Update Project With Invalid DelayMs Too High", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test_update_timeout_high@example.com", "Test User Update DelayMs High", "Test Workspace Update DelayMs High")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project
		project, err := database.CreateTestProject(workspace.ID, "Test Project Update DelayMs High", generateUniqueAliasUpdate("test-project-update-delayMs-high"))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId", UpdateProjectHandler)

		// Create update data with invalid delayMs (too high)
		updateData := map[string]interface{}{
			"advance_config": `{"delayMs": 130000}`, // Above maximum 120000ms (2 minutes)
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID, bytes.NewBuffer(jsonData))
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
		assert.Contains(t, response["message"].(string), "delayMs cannot exceed 120000ms")
	})

	t.Run("Update Project With Valid DelayMs", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test_update_timeout_valid@example.com", "Test User Update DelayMs Valid", "Test Workspace Update DelayMs Valid")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create test project
		project, err := database.CreateTestProject(workspace.ID, "Test Project Update DelayMs Valid", generateUniqueAliasUpdate("test-project-update-delayMs-valid"))
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.PUT("/api/projects/:projectId", UpdateProjectHandler)

		// Create update data with valid delayMs
		updateData := map[string]interface{}{
			"advance_config": `{"delayMs": 10000}`, // Valid 10 seconds delayMs
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("PUT", "/api/projects/"+project.ID, bytes.NewBuffer(jsonData))
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
		assert.Equal(t, "Project updated successfully", response["message"])

		// Verify updated project data
		projectData := response["data"].(map[string]interface{})
		assert.Equal(t, `{"delayMs": 10000}`, projectData["advance_config"])
	})
}
