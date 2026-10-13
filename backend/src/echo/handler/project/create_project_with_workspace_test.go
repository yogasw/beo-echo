package project

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
)

// createViaWorkspaceRoute calls the handler that POST /api/workspaces/:workspaceID/projects uses,
// which is the route the UI, the API and the MCP server go through.
func createViaWorkspaceRoute(t *testing.T, userID, workspaceID string, body map[string]interface{}) map[string]interface{} {
	t.Helper()
	router := gin.New()
	router.POST("/api/workspaces/:workspaceID/projects", func(c *gin.Context) {
		c.Set("userID", userID)
		CreateProjectWithWorkspaceHandler(c)
	})

	payload, err := json.Marshal(body)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, "/api/workspaces/"+workspaceID+"/projects", bytes.NewBuffer(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	return response["data"].(map[string]interface{})
}

func TestCreateProjectWithWorkspace_DefaultsToReflectCors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setup, err := database.InitTestWorkspaceOnly("ws-default-cors@example.com", "WS Default", "WS Default CORS")
	require.NoError(t, err)
	defer setup.Cleanup()

	data := createViaWorkspaceRoute(t, setup.User.ID, setup.Workspace.ID, map[string]interface{}{
		"name":  "Default CORS",
		"alias": generateUniqueAlias("ws-default-cors"),
		"mode":  "mock",
	})
	defer database.DB.Unscoped().Delete(&database.Project{}, "id = ?", data["id"])

	assert.JSONEq(t, `{"cors":{"mode":"reflect"}}`, data["advance_config"].(string))
}

func TestCreateProjectWithWorkspace_KeepsProvidedAdvanceConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setup, err := database.InitTestWorkspaceOnly("ws-keep-cors@example.com", "WS Keep", "WS Keep CORS")
	require.NoError(t, err)
	defer setup.Cleanup()

	data := createViaWorkspaceRoute(t, setup.User.ID, setup.Workspace.ID, map[string]interface{}{
		"name":           "Keep CORS",
		"alias":          generateUniqueAlias("ws-keep-cors"),
		"mode":           "mock",
		"advance_config": `{"delayMs":250}`,
	})
	defer database.DB.Unscoped().Delete(&database.Project{}, "id = ?", data["id"])

	assert.JSONEq(t, `{"delayMs":250}`, data["advance_config"].(string))
}
