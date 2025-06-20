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

func TestReorderResponsesTestSuite(t *testing.T) {
	suite.Run(t, new(ReorderResponsesTestSuite))
}
