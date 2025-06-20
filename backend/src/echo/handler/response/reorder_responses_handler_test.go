package response

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"beo-echo/backend/src/database"
)

type ReorderResponsesTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *ReorderResponsesTestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize test database
	database.SetupTestEnvironment(suite.T())

	// Create router and register routes
	suite.router = gin.New()
	suite.router.PUT("/api/workspaces/:workspaceId/projects/:projectId/endpoints/:id/responses/reorder", ReorderResponsesHandler)
}

func (suite *ReorderResponsesTestSuite) TestReorderResponsesHandler_MissingProjectId() {
	// Test with missing project ID
	requestBody := ReorderResponsesRequest{
		Order: []string{"response1", "response2"},
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/workspaces/ws1/projects//endpoints/endpoint1/responses/reorder", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, response["error"])
	assert.Contains(suite.T(), response["message"], "Project ID is required")
}

func (suite *ReorderResponsesTestSuite) TestReorderResponsesHandler_MissingEndpointId() {
	// Test with missing endpoint ID
	requestBody := ReorderResponsesRequest{
		Order: []string{"response1", "response2"},
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/workspaces/ws1/projects/project1/endpoints//responses/reorder", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, response["error"])
	assert.Contains(suite.T(), response["message"], "Endpoint ID is required")
}

func (suite *ReorderResponsesTestSuite) TestReorderResponsesHandler_EmptyOrder() {
	// Test with empty order array
	requestBody := ReorderResponsesRequest{
		Order: []string{},
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/workspaces/ws1/projects/project1/endpoints/endpoint1/responses/reorder", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, response["error"])
	assert.Contains(suite.T(), response["message"], "Order array cannot be empty")
}

func (suite *ReorderResponsesTestSuite) TestReorderResponsesHandler_InvalidJSON() {
	// Test with invalid JSON
	req, _ := http.NewRequest("PUT", "/api/workspaces/ws1/projects/project1/endpoints/endpoint1/responses/reorder", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, response["error"])
	assert.Contains(suite.T(), response["message"], "Invalid request format")
}

func (suite *ReorderResponsesTestSuite) TestReorderResponsesHandler_PriorityOrdering() {
	// Initialize test database for this specific test
	database.SetupTestEnvironment(suite.T())

	// Create test project
	project := &database.Project{
		ID:          "test-project-priority",
		Name:        "Test Project Priority",
		WorkspaceID: "test-workspace",
		Alias:       "test-priority-alias",
		Mode:        database.ModeMock,
		Status:      "running",
	}
	result := database.DB.Create(project)
	assert.NoError(suite.T(), result.Error)

	// Create test endpoint
	endpoint := &database.MockEndpoint{
		ID:           "test-endpoint-priority",
		ProjectID:    project.ID,
		Path:         "/test-priority",
		Method:       "GET",
		ResponseMode: "static",
		Enabled:      true,
	}
	result = database.DB.Create(endpoint)
	assert.NoError(suite.T(), result.Error)

	// Create test responses with scrambled priorities to test sorting
	responses := []database.MockResponse{
		{
			ID:         "response-low",
			EndpointID: endpoint.ID,
			StatusCode: 200,
			Body:       `{"message": "Low Priority"}`,
			Headers:    `{"Content-Type": "application/json"}`,
			Priority:   2, // Low priority
			Enabled:    true,
			Note:       "Low Priority Response",
		},
		{
			ID:         "response-highest", 
			EndpointID: endpoint.ID,
			StatusCode: 201,
			Body:       `{"message": "Highest Priority"}`,
			Headers:    `{"Content-Type": "application/json"}`,
			Priority:   10, // Highest priority
			Enabled:    true,
			Note:       "Highest Priority Response",
		},
		{
			ID:         "response-medium",
			EndpointID: endpoint.ID,
			StatusCode: 202,
			Body:       `{"message": "Medium Priority"}`,
			Headers:    `{"Content-Type": "application/json"}`,
			Priority:   5, // Medium priority
			Enabled:    true,
			Note:       "Medium Priority Response",
		},
		{
			ID:         "response-high",
			EndpointID: endpoint.ID,
			StatusCode: 203,
			Body:       `{"message": "High Priority"}`,
			Headers:    `{"Content-Type": "application/json"}`,
			Priority:   8, // High priority
			Enabled:    true,
			Note:       "High Priority Response",
		},
	}

	for _, response := range responses {
		result = database.DB.Create(&response)
		assert.NoError(suite.T(), result.Error)
	}

	// Test reordering - frontend sends desired visual order
	requestBody := ReorderResponsesRequest{
		Order: []string{"response-highest", "response-high", "response-medium", "response-low"},
	}
	jsonBody, err := json.Marshal(requestBody)
	assert.NoError(suite.T(), err)

	req, err := http.NewRequest("PUT", "/api/workspaces/test-workspace/projects/test-project-priority/endpoints/test-endpoint-priority/responses/reorder", bytes.NewBuffer(jsonBody))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert successful reordering
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), false, response["error"])
	assert.Equal(suite.T(), "Responses reordered successfully", response["message"])

	// Verify that responses are returned ordered by priority DESC (highest first)
	data := response["data"].([]interface{})
	assert.Len(suite.T(), data, 4)

	// Expected order after reordering: responses should have priorities assigned based on order array
	// Order: ["response-highest", "response-high", "response-medium", "response-low"]
	// This means: response-highest gets priority 4, response-high gets priority 3, etc.
	expectedOrder := []struct {
		id       string
		priority float64
		note     string
	}{
		{"response-highest", 4, "Highest Priority Response"}, // Index 0 → Priority 4 (highest)
		{"response-high", 3, "High Priority Response"},       // Index 1 → Priority 3  
		{"response-medium", 2, "Medium Priority Response"},   // Index 2 → Priority 2
		{"response-low", 1, "Low Priority Response"},         // Index 3 → Priority 1 (lowest)
	}

	for i, expected := range expectedOrder {
		responseData := data[i].(map[string]interface{})
		assert.Equal(suite.T(), expected.id, responseData["id"], "Response at index %d should have ID %s", i, expected.id)
		assert.Equal(suite.T(), expected.priority, responseData["priority"], "Response at index %d should have priority %f", i, expected.priority)
		assert.Equal(suite.T(), expected.note, responseData["note"], "Response at index %d should have note %s", i, expected.note)
	}

	// Verify database consistency - check that responses are stored with correct priorities
	var dbResponses []database.MockResponse
	result = database.DB.Where("endpoint_id = ?", endpoint.ID).Order("priority DESC").Find(&dbResponses)
	assert.NoError(suite.T(), result.Error)
	assert.Len(suite.T(), dbResponses, 4)

	// Verify database order matches expected priority order after reordering
	expectedDBOrder := []struct {
		id       string
		priority int
	}{
		{"response-highest", 4}, // Highest priority (first in order array)
		{"response-high", 3},    // Second highest priority
		{"response-medium", 2},  // Third highest priority
		{"response-low", 1},     // Lowest priority (last in order array)
	}

	for i, expected := range expectedDBOrder {
		assert.Equal(suite.T(), expected.id, dbResponses[i].ID, "DB Response at index %d should have ID %s", i, expected.id)
		assert.Equal(suite.T(), expected.priority, dbResponses[i].Priority, "DB Response at index %d should have priority %d", i, expected.priority)
	}
}

func (suite *ReorderResponsesTestSuite) TestReorderResponsesHandler_SingleResponse() {
	// Test reordering with only one response
	database.SetupTestEnvironment(suite.T())

	// Create test project
	project := &database.Project{
		ID:          "test-project-single",
		Name:        "Test Project Single Response",
		WorkspaceID: "test-workspace",
		Alias:       "test-single-alias",
		Mode:        database.ModeMock,
		Status:      "running",
	}
	result := database.DB.Create(project)
	assert.NoError(suite.T(), result.Error)

	// Create test endpoint
	endpoint := &database.MockEndpoint{
		ID:           "test-endpoint-single",
		ProjectID:    project.ID,
		Path:         "/test-single",
		Method:       "GET",
		ResponseMode: "static",
		Enabled:      true,
	}
	result = database.DB.Create(endpoint)
	assert.NoError(suite.T(), result.Error)

	// Create single response
	response := &database.MockResponse{
		ID:         "response-only",
		EndpointID: endpoint.ID,
		StatusCode: 200,
		Body:       `{"message": "Only Response"}`,
		Headers:    `{"Content-Type": "application/json"}`,
		Priority:   5,
		Enabled:    true,
		Note:       "Only Response",
	}
	result = database.DB.Create(response)
	assert.NoError(suite.T(), result.Error)

	// Test reordering with single response
	requestBody := ReorderResponsesRequest{
		Order: []string{"response-only"},
	}
	jsonBody, err := json.Marshal(requestBody)
	assert.NoError(suite.T(), err)

	req, err := http.NewRequest("PUT", "/api/workspaces/test-workspace/projects/test-project-single/endpoints/test-endpoint-single/responses/reorder", bytes.NewBuffer(jsonBody))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert successful operation
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), false, resp["error"])
	
	// Verify single response is returned
	data := resp["data"].([]interface{})
	assert.Len(suite.T(), data, 1)
	
	responseData := data[0].(map[string]interface{})
	assert.Equal(suite.T(), "response-only", responseData["id"])
	assert.Equal(suite.T(), float64(1), responseData["priority"]) // Single response gets priority 1
}

func TestReorderResponsesTestSuite(t *testing.T) {
	suite.Run(t, new(ReorderResponsesTestSuite))
}
