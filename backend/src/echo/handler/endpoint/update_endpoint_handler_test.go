package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"beo-echo/backend/src/database"
)

func TestUpdateEndpointHandler_ValidAdvanceConfig(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Initialize test workspace with project and defer cleanup
	setup, err := database.InitTestWorkspaceWithProject(
		"test@example.com",
		"Test User",
		"Test Workspace",
		"Test Project",
		"test-project-123",
	)
	if err != nil {
		t.Fatalf("Failed to initialize test workspace: %v", err)
	}
	defer setup.Cleanup()

	testProject := setup.Project

	// Create test endpoint using helper function
	testEndpoint, err := database.CreateTestEndpointWithConfig(
		testProject.ID,
		"GET",
		"/test",
		`{"timeout": 30000}`,
	)
	if err != nil {
		t.Fatalf("Failed to create test endpoint: %v", err)
	}

	updateData := map[string]interface{}{
		"method":         "POST",
		"path":           "/updated",
		"advance_config": `{"timeout": 60000, "retries": 3}`,
	}

	// Create request
	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/api/projects/"+testProject.ID+"/endpoints/"+testEndpoint.ID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{
		{Key: "projectId", Value: testProject.ID},
		{Key: "id", Value: testEndpoint.ID},
	}

	// Execute
	UpdateEndpointHandler(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the endpoint was updated in database
	var updatedEndpoint database.MockEndpoint
	database.DB.First(&updatedEndpoint, "id = ?", testEndpoint.ID)
	assert.Equal(t, "POST", updatedEndpoint.Method)
	assert.Equal(t, "/updated", updatedEndpoint.Path)
	assert.Equal(t, `{"timeout": 60000, "retries": 3}`, updatedEndpoint.AdvanceConfig)
}

func TestUpdateEndpointHandler_InvalidAdvanceConfig(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Initialize test workspace with project and defer cleanup
	setup, err := database.InitTestWorkspaceWithProject(
		"test2@example.com",
		"Test User 2",
		"Test Workspace 2",
		"Test Project 2",
		"test-project-456",
	)
	if err != nil {
		t.Fatalf("Failed to initialize test workspace: %v", err)
	}
	defer setup.Cleanup()

	testProject := setup.Project

	// Create test endpoint using helper function
	testEndpoint, err := database.CreateTestEndpointWithConfig(
		testProject.ID,
		"GET",
		"/test",
		`{"timeout": 30000}`,
	)
	if err != nil {
		t.Fatalf("Failed to create test endpoint: %v", err)
	}

	updateData := map[string]interface{}{
		"method":         "POST",
		"path":           "/updated",
		"advance_config": `{"timeout": 60000, "retries": 3`, // Invalid JSON - missing closing brace
	}

	// Create request
	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/api/projects/"+testProject.ID+"/endpoints/"+testEndpoint.ID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{
		{Key: "projectId", Value: testProject.ID},
		{Key: "id", Value: testEndpoint.ID},
	}

	// Execute
	UpdateEndpointHandler(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Verify response contains error message
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"], "Invalid JSON format in advance_config")

	// Verify original endpoint was not modified
	var originalEndpoint database.MockEndpoint
	database.DB.First(&originalEndpoint, "id = ?", testEndpoint.ID)
	assert.Equal(t, "GET", originalEndpoint.Method)                       // Should remain unchanged
	assert.Equal(t, "/test", originalEndpoint.Path)                       // Should remain unchanged
	assert.Equal(t, `{"timeout": 30000}`, originalEndpoint.AdvanceConfig) // Should remain unchanged
}

func TestUpdateEndpointHandler_EmptyAdvanceConfig(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Initialize test workspace with project and defer cleanup
	setup, err := database.InitTestWorkspaceWithProject(
		"test3@example.com",
		"Test User 3",
		"Test Workspace 3",
		"Test Project 3",
		"test-project-789",
	)
	if err != nil {
		t.Fatalf("Failed to initialize test workspace: %v", err)
	}
	defer setup.Cleanup()

	testProject := setup.Project

	// Create test endpoint using helper function
	testEndpoint, err := database.CreateTestEndpointWithConfig(
		testProject.ID,
		"GET",
		"/test",
		`{"timeout": 30000}`,
	)
	if err != nil {
		t.Fatalf("Failed to create test endpoint: %v", err)
	}

	updateData := map[string]interface{}{
		"method":         "POST",
		"path":           "/updated",
		"advance_config": "", // Empty string should be valid
	}

	// Create request
	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/api/projects/"+testProject.ID+"/endpoints/"+testEndpoint.ID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{
		{Key: "projectId", Value: testProject.ID},
		{Key: "id", Value: testEndpoint.ID},
	}

	// Execute
	UpdateEndpointHandler(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the endpoint was updated with empty config
	var updatedEndpoint database.MockEndpoint
	database.DB.First(&updatedEndpoint, "id = ?", testEndpoint.ID)
	assert.Equal(t, "POST", updatedEndpoint.Method)
	assert.Equal(t, "/updated", updatedEndpoint.Path)
	assert.Equal(t, "", updatedEndpoint.AdvanceConfig)
}

func TestUpdateEndpointHandler_EndpointNotFound(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Initialize test workspace with project and defer cleanup
	setup, err := database.InitTestWorkspaceWithProject(
		"test4@example.com",
		"Test User 4",
		"Test Workspace 4",
		"Test Project 4",
		"test-project-404",
	)
	if err != nil {
		t.Fatalf("Failed to initialize test workspace: %v", err)
	}
	defer setup.Cleanup()

	testProject := setup.Project

	updateData := map[string]interface{}{
		"method":         "POST",
		"advance_config": `{"timeout": 60000}`,
	}

	// Create request
	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/api/projects/"+testProject.ID+"/endpoints/endpoint-404", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{
		{Key: "projectId", Value: testProject.ID},
		{Key: "id", Value: "endpoint-404"},
	}

	// Execute
	UpdateEndpointHandler(c)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Verify response contains error message
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["error"].(bool))
	assert.Contains(t, response["message"], "Endpoint not found")
}
