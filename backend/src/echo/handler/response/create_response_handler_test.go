package response

import (
	"bytes"
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

func TestCreateResponseHandler_Enabled(t *testing.T) {
	database.SetupTestEnvironment(t)
	handler.EnsureMockService()

	project := &database.Project{ID: uuid.New().String(), Name: "Test Project"}
	require.NoError(t, database.GetDB().Create(project).Error)
	endpoint := &database.MockEndpoint{ID: uuid.New().String(), ProjectID: project.ID, Path: "/test", Method: "GET"}
	require.NoError(t, database.GetDB().Create(endpoint).Error)

	router := gin.New()
	router.POST("/projects/:projectId/endpoints/:id/responses", CreateResponseHandler)

	create := func(t *testing.T, body map[string]interface{}) (map[string]interface{}, database.MockResponse) {
		t.Helper()
		raw, err := json.Marshal(body)
		require.NoError(t, err)
		req, _ := http.NewRequest("POST", "/projects/"+project.ID+"/endpoints/"+endpoint.ID+"/responses", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code, w.Body.String())

		var out map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
		data := out["data"].(map[string]interface{})

		var stored database.MockResponse
		require.NoError(t, database.GetDB().First(&stored, "id = ?", data["id"]).Error)
		return data, stored
	}

	t.Run("explicit enabled false is stored disabled", func(t *testing.T) {
		data, stored := create(t, map[string]interface{}{"status_code": 200, "priority": 40, "enabled": false})
		assert.Equal(t, false, data["enabled"])
		assert.False(t, stored.Enabled)
	})

	t.Run("omitted enabled defaults to enabled", func(t *testing.T) {
		data, stored := create(t, map[string]interface{}{"status_code": 200})
		assert.Equal(t, true, data["enabled"])
		assert.True(t, stored.Enabled)
	})

	t.Run("explicit enabled true is stored enabled", func(t *testing.T) {
		data, stored := create(t, map[string]interface{}{"status_code": 200, "enabled": true})
		assert.Equal(t, true, data["enabled"])
		assert.True(t, stored.Enabled)
	})
}

func TestDuplicateResponseHandler_KeepsDisabled(t *testing.T) {
	database.SetupTestEnvironment(t)
	handler.EnsureMockService()

	project := &database.Project{ID: uuid.New().String(), Name: "Test Project"}
	require.NoError(t, database.GetDB().Create(project).Error)
	endpoint := &database.MockEndpoint{ID: uuid.New().String(), ProjectID: project.ID, Path: "/test", Method: "GET"}
	require.NoError(t, database.GetDB().Create(endpoint).Error)

	original := &database.MockResponse{ID: uuid.New().String(), EndpointID: endpoint.ID, StatusCode: 200, Note: "off"}
	require.NoError(t, database.GetDB().Create(original).Error)
	// Create cannot store false for a default:true column, so disable it explicitly.
	require.NoError(t, database.GetDB().Model(original).Update("enabled", false).Error)

	router := gin.New()
	router.POST("/projects/:projectId/endpoints/:id/responses/:responseId/duplicate", DuplicateResponseHandler)
	req, _ := http.NewRequest("POST", "/projects/"+project.ID+"/endpoints/"+endpoint.ID+"/responses/"+original.ID+"/duplicate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())

	var out map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	data := out["data"].(map[string]interface{})
	assert.Equal(t, false, data["enabled"])

	var stored database.MockResponse
	require.NoError(t, database.GetDB().First(&stored, "id = ?", data["id"]).Error)
	assert.False(t, stored.Enabled)
}
