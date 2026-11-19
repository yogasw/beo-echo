package endpoint

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

func TestCreateEndpointHandler_Enabled(t *testing.T) {
	database.SetupTestEnvironment(t)
	handler.EnsureMockService()

	project := &database.Project{ID: uuid.New().String(), Name: "Test Project"}
	require.NoError(t, database.GetDB().Create(project).Error)

	router := gin.New()
	router.POST("/projects/:projectId/endpoints", CreateEndpointHandler)

	create := func(t *testing.T, body map[string]interface{}) (map[string]interface{}, database.MockEndpoint) {
		t.Helper()
		raw, err := json.Marshal(body)
		require.NoError(t, err)
		req, _ := http.NewRequest("POST", "/projects/"+project.ID+"/endpoints", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code, w.Body.String())

		var out map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
		data := out["data"].(map[string]interface{})

		var stored database.MockEndpoint
		require.NoError(t, database.GetDB().First(&stored, "id = ?", data["id"]).Error)
		return data, stored
	}

	t.Run("explicit enabled false is stored disabled", func(t *testing.T) {
		data, stored := create(t, map[string]interface{}{"method": "GET", "path": "/off", "enabled": false})
		assert.Equal(t, false, data["enabled"])
		assert.False(t, stored.Enabled)
	})

	t.Run("omitted enabled defaults to enabled", func(t *testing.T) {
		data, stored := create(t, map[string]interface{}{"method": "GET", "path": "/on"})
		assert.Equal(t, true, data["enabled"])
		assert.True(t, stored.Enabled)
	})
}
