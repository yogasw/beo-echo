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
func generateUniqueAlias(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano()/int64(time.Millisecond))
}

func TestCreateProjectHandler(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Initialize database
	err := database.CheckAndHandle()
	require.NoError(t, err, "Failed to setup test database")

	t.Run("Create Project Without Advance Config", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test1@example.com", "Test User 1", "Test Workspace 1")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Setup Gin router
		router := gin.New()
		router.POST("/api/projects", CreateProjectHandler)

		// Create request data
		projectAlias := generateUniqueAlias("test-project-create-1")
		projectData := map[string]interface{}{
			"name":         "Test Project",
			"alias":        projectAlias,
			"mode":         "mock",
			"workspace_id": workspace.ID,
		}

		jsonData, err := json.Marshal(projectData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Debug: Print response body if test fails
		if w.Code != http.StatusCreated {
			t.Logf("Response body: %s", w.Body.String())
		}

		// Verify response
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Check if response is success or error
		if success, ok := response["success"].(bool); ok && success {
			assert.True(t, response["success"].(bool))
			assert.Equal(t, "Project created successfully", response["message"])

			// Verify project data
			projectData = response["data"].(map[string]interface{})
			assert.Equal(t, "Test Project", projectData["name"])
			assert.Equal(t, projectAlias, projectData["alias"])
			assert.Equal(t, "mock", projectData["mode"])
			assert.Equal(t, workspace.ID, projectData["workspace_id"])
			assert.Equal(t, "", projectData["advance_config"]) // Should be empty by default
		} else {
			t.Fatalf("Expected success response, got error: %s", response["message"])
		}
	})

	t.Run("Create Project With Valid Advance Config", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test2@example.com", "Test User 2", "Test Workspace 2")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Setup Gin router
		router := gin.New()
		router.POST("/api/projects", CreateProjectHandler)

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

		// Create request data
		projectAlias := generateUniqueAlias("test-project-create-2")
		projectData := map[string]interface{}{
			"name":           "Test Project 2",
			"alias":          projectAlias,
			"mode":           "proxy",
			"workspace_id":   workspace.ID,
			"advance_config": string(advanceConfigJSON),
		}

		jsonData, err := json.Marshal(projectData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Project created successfully", response["message"])

		// Verify project data
		responseProjectData := response["data"].(map[string]interface{})
		assert.Equal(t, "Test Project 2", responseProjectData["name"])
		assert.Equal(t, projectAlias, responseProjectData["alias"])
		assert.Equal(t, "proxy", responseProjectData["mode"])
		assert.Equal(t, workspace.ID, responseProjectData["workspace_id"])
		assert.NotEmpty(t, responseProjectData["advance_config"])

		// Verify advance config can be parsed
		var storedAdvanceConfig map[string]interface{}
		err = json.Unmarshal([]byte(responseProjectData["advance_config"].(string)), &storedAdvanceConfig)
		require.NoError(t, err)
		assert.Equal(t, float64(5000), storedAdvanceConfig["global_timeout"])
	})

	t.Run("Create Project With Invalid JSON Advance Config", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test3@example.com", "Test User 3", "Test Workspace 3")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Setup Gin router
		router := gin.New()
		router.POST("/api/projects", CreateProjectHandler)

		// Create request data with invalid JSON in advance_config
		projectData := map[string]interface{}{
			"name":           "Test Project 3",
			"alias":          "test-project-create-3",
			"mode":           "mock",
			"workspace_id":   workspace.ID,
			"advance_config": `{"global_timeout": 5000, "invalid":}`, // Invalid JSON
		}

		jsonData, err := json.Marshal(projectData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
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

	t.Run("Create Project With Missing Required Fields", func(t *testing.T) {
		// Setup Gin router
		router := gin.New()
		router.POST("/api/projects", CreateProjectHandler)

		// Create request data without name
		projectData := map[string]interface{}{
			"alias": "test-project-create-4",
			"mode":  "mock",
		}

		jsonData, err := json.Marshal(projectData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
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
		assert.Contains(t, response["message"].(string), "Project name and alias are required")
	})

	t.Run("Create Project With Invalid Alias", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test4@example.com", "Test User 4", "Test Workspace 4")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Setup Gin router
		router := gin.New()
		router.POST("/api/projects", CreateProjectHandler)

		// Create request data with invalid alias (contains spaces and special chars)
		projectData := map[string]interface{}{
			"name":         "Test Project 4",
			"alias":        "test project@#$",
			"mode":         "mock",
			"workspace_id": workspace.ID,
		}

		jsonData, err := json.Marshal(projectData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
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
		assert.Contains(t, response["message"].(string), "alphanumeric characters, dashes and underscores")
	})

	t.Run("Create Project With Duplicate Alias", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test5@example.com", "Test User 5", "Test Workspace 5")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Create first project with unique alias first
		duplicateAlias := generateUniqueAlias("duplicate-alias-test-create")
		_, err = database.CreateTestProject(workspace.ID, "Test Project 5A", duplicateAlias)
		require.NoError(t, err)

		// Setup Gin router
		router := gin.New()
		router.POST("/api/projects", CreateProjectHandler)

		// Try to create second project with same alias
		projectData := map[string]interface{}{
			"name":         "Test Project 5B",
			"alias":        duplicateAlias, // Use the same alias
			"mode":         "mock",
			"workspace_id": workspace.ID,
		}

		jsonData, err := json.Marshal(projectData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusConflict, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"].(string), "Project alias already exists")
	})

	t.Run("Create Project With Invalid delayMs Too Low", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test_delayms_low@example.com", "Test User DelayMs Low", "Test Workspace DelayMs Low")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Setup Gin router
		router := gin.New()
		router.POST("/api/projects", CreateProjectHandler)

		// Create request data with negative delayMs
		projectAlias := generateUniqueAlias("test-project-delayms-low")
		projectData := map[string]interface{}{
			"name":           "Test Project DelayMs Low",
			"alias":          projectAlias,
			"mode":           "mock",
			"workspace_id":   workspace.ID,
			"advance_config": `{"delayMs": -1000}`, // Negative delayMs
		}

		jsonData, err := json.Marshal(projectData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
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
		assert.Contains(t, response["message"].(string), "delayMs cannot be negative")
	})

	t.Run("Create Project With Invalid DelayMs Too High", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test_delayms_high@example.com", "Test User DelayMs High", "Test Workspace DelayMs High")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Setup Gin router
		router := gin.New()
		router.POST("/api/projects", CreateProjectHandler)

		// Create request data with delayMs too high
		projectAlias := generateUniqueAlias("test-project-delayms-high")
		projectData := map[string]interface{}{
			"name":           "Test Project DelayMs High",
			"alias":          projectAlias,
			"mode":           "mock",
			"workspace_id":   workspace.ID,
			"advance_config": `{"delayMs": 130000}`, // Above maximum 120000ms (2 minutes)
		}

		jsonData, err := json.Marshal(projectData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
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

	t.Run("Create Project With Valid DelayMs Range", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := database.CreateTestWorkspace("test_delayms_valid@example.com", "Test User DelayMs Valid", "Test Workspace DelayMs Valid")
		require.NoError(t, err)
		defer database.CleanupTestData(user.ID, workspace.ID, "", "")

		// Setup Gin router
		router := gin.New()
		router.POST("/api/projects", CreateProjectHandler)

		// Create request data with valid delayMs
		projectAlias := generateUniqueAlias("test-project-delayms-valid")
		projectData := map[string]interface{}{
			"name":           "Test Project DelayMs Valid",
			"alias":          projectAlias,
			"mode":           "mock",
			"workspace_id":   workspace.ID,
			"advance_config": `{"delayMs": 5000}`, // Valid delayMs 5 seconds
		}

		jsonData, err := json.Marshal(projectData)
		require.NoError(t, err)

		// Create HTTP request
		req, err := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Project created successfully", response["message"])

		// Verify project data
		responseData := response["data"].(map[string]interface{})
		assert.Equal(t, "Test Project DelayMs Valid", responseData["name"])
		assert.Equal(t, projectAlias, responseData["alias"])
		assert.Contains(t, responseData["advance_config"].(string), "5000") // Check delayMs value exists
	})
}
