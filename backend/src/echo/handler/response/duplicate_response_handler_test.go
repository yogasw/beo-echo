package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

func TestDuplicateResponseHandler(t *testing.T) {
	// Setup test database
	database.SetupTestEnvironment(t)

	// Initialize mock service
	handler.EnsureMockService()

	// Create test data
	project := &database.Project{
		ID:   uuid.New().String(),
		Name: "Test Project",
	}
	require.NoError(t, database.GetDB().Create(project).Error)

	endpoint := &database.MockEndpoint{
		ID:        uuid.New().String(),
		ProjectID: project.ID,
		Path:      "/test",
		Method:    "GET",
	}
	require.NoError(t, database.GetDB().Create(endpoint).Error)

	// Create original response with rules
	originalResponse := &database.MockResponse{
		ID:         uuid.New().String(),
		EndpointID: endpoint.ID,
		StatusCode: 200,
		Body:       `{"message": "Hello World"}`,
		Headers:    `{"Content-Type": "application/json"}`,
		Priority:   1,
		DelayMS:    100,
		Stream:     false,
		Note:       "Original response",
		Enabled:    true,
	}
	require.NoError(t, database.GetDB().Create(originalResponse).Error)

	// Create rules for the original response
	rule1 := &database.MockRule{
		ID:         uuid.New().String(),
		ResponseID: originalResponse.ID,
		Type:       "header",
		Key:        "X-Auth",
		Operator:   "equals",
		Value:      "bearer-token",
	}
	rule2 := &database.MockRule{
		ID:         uuid.New().String(),
		ResponseID: originalResponse.ID,
		Type:       "query",
		Key:        "version",
		Operator:   "contains",
		Value:      "v2",
	}
	require.NoError(t, database.GetDB().Create(rule1).Error)
	require.NoError(t, database.GetDB().Create(rule2).Error)

	t.Run("Successfully duplicate response with rules", func(t *testing.T) {
		// Setup router
		router := gin.New()
		router.POST("/api/workspaces/:workspaceID/projects/:projectId/endpoints/:id/responses/:responseId/duplicate", DuplicateResponseHandler)

		// Create request
		req, _ := http.NewRequest("POST", "/api/workspaces/test-workspace/projects/"+project.ID+"/endpoints/"+endpoint.ID+"/responses/"+originalResponse.ID+"/duplicate", nil)
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Response duplicated successfully", response["message"])

		// Verify response data
		responseData := response["data"].(map[string]interface{})
		assert.NotEqual(t, originalResponse.ID, responseData["id"]) // Should have new ID
		assert.Equal(t, originalResponse.EndpointID, responseData["endpoint_id"])
		assert.Equal(t, float64(originalResponse.StatusCode), responseData["status_code"])
		assert.Equal(t, originalResponse.Body, responseData["body"])
		assert.Equal(t, originalResponse.Headers, responseData["headers"])
		assert.Equal(t, float64(originalResponse.Priority), responseData["priority"])
		assert.Equal(t, float64(originalResponse.DelayMS), responseData["delay_ms"])
		assert.Equal(t, originalResponse.Stream, responseData["stream"])
		assert.Equal(t, "Original response (Copy)", responseData["note"]) // Should have "(Copy)" appended
		assert.Equal(t, originalResponse.Enabled, responseData["enabled"])

		// Verify rules were duplicated
		rules := responseData["rules"].([]interface{})
		assert.Len(t, rules, 2)

		// Check that rules have different IDs but same content
		duplicatedResponse := responseData["id"].(string)
		var duplicatedRules []database.MockRule
		database.GetDB().Where("response_id = ?", duplicatedResponse).Find(&duplicatedRules)

		assert.Len(t, duplicatedRules, 2)
		for _, rule := range duplicatedRules {
			assert.Equal(t, duplicatedResponse, rule.ResponseID)
			assert.NotEqual(t, rule1.ID, rule.ID) // Should have new ID
			assert.NotEqual(t, rule2.ID, rule.ID) // Should have new ID
		}
	})

	t.Run("Error when project not found", func(t *testing.T) {
		router := gin.New()
		router.POST("/api/workspaces/:workspaceID/projects/:projectId/endpoints/:id/responses/:responseId/duplicate", DuplicateResponseHandler)

		req, _ := http.NewRequest("POST", "/api/workspaces/test-workspace/projects/non-existent/endpoints/"+endpoint.ID+"/responses/"+originalResponse.ID+"/duplicate", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"].(string), "Project not found")
	})

	t.Run("Error when response not found", func(t *testing.T) {
		router := gin.New()
		router.POST("/api/workspaces/:workspaceID/projects/:projectId/endpoints/:id/responses/:responseId/duplicate", DuplicateResponseHandler)

		req, _ := http.NewRequest("POST", "/api/workspaces/test-workspace/projects/"+project.ID+"/endpoints/"+endpoint.ID+"/responses/non-existent/duplicate", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"].(string), "Response not found")
	})

	t.Run("Error when endpoint not found", func(t *testing.T) {
		router := gin.New()
		router.POST("/api/workspaces/:workspaceID/projects/:projectId/endpoints/:id/responses/:responseId/duplicate", DuplicateResponseHandler)

		req, _ := http.NewRequest("POST", "/api/workspaces/test-workspace/projects/"+project.ID+"/endpoints/non-existent/responses/"+originalResponse.ID+"/duplicate", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"].(string), "Endpoint not found")
	})

	t.Run("Error when missing required parameters", func(t *testing.T) {
		router := gin.New()
		router.POST("/api/workspaces/:workspaceID/projects/:projectId/endpoints/:id/responses/:responseId/duplicate", DuplicateResponseHandler)

		// Test missing project ID
		req, _ := http.NewRequest("POST", "/api/workspaces/test-workspace/projects//endpoints/"+endpoint.ID+"/responses/"+originalResponse.ID+"/duplicate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test missing endpoint ID
		req, _ = http.NewRequest("POST", "/api/workspaces/test-workspace/projects/"+project.ID+"/endpoints//responses/"+originalResponse.ID+"/duplicate", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test missing response ID
		req, _ = http.NewRequest("POST", "/api/workspaces/test-workspace/projects/"+project.ID+"/endpoints/"+endpoint.ID+"/responses//duplicate", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
