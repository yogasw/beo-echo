package project

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
)

func putAdvanceConfig(t *testing.T, projectID, body string) *httptest.ResponseRecorder {
	t.Helper()
	router := gin.New()
	router.PUT("/api/projects/:projectId/advance-config", UpdateProjectAdvanceConfigHandler)

	req, err := http.NewRequest(http.MethodPut, "/api/projects/"+projectID+"/advance-config", bytes.NewBufferString(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func storedAdvanceConfig(t *testing.T, projectID string) string {
	t.Helper()
	var p database.Project
	require.NoError(t, database.DB.First(&p, "id = ?", projectID).Error)
	return p.AdvanceConfig
}

func TestUpdateProjectAdvanceConfig_MergesKeys(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setup, err := database.InitTestWorkspaceWithProject("adv-merge@example.com", "Adv Merge", "Adv Merge WS", "Adv Merge Project", generateUniqueAlias("adv-merge"))
	require.NoError(t, err)
	defer setup.Cleanup()
	id := setup.Project.ID

	// Set cors first, then delay: cors must survive the delay update.
	require.Equal(t, http.StatusOK, putAdvanceConfig(t, id, `{"cors":{"mode":"reflect"}}`).Code)
	require.Equal(t, http.StatusOK, putAdvanceConfig(t, id, `{"delayMs":300}`).Code)
	assert.JSONEq(t, `{"cors":{"mode":"reflect"},"delayMs":300}`, storedAdvanceConfig(t, id))

	// ...and the other way around: changing cors keeps the delay.
	require.Equal(t, http.StatusOK, putAdvanceConfig(t, id, `{"cors":{"mode":"custom","allowHeaders":["X-Api-Key"]}}`).Code)
	assert.JSONEq(t, `{"cors":{"mode":"custom","allowHeaders":["X-Api-Key"]},"delayMs":300}`, storedAdvanceConfig(t, id))

	// null removes a key and keeps the rest.
	require.Equal(t, http.StatusOK, putAdvanceConfig(t, id, `{"cors":null}`).Code)
	assert.JSONEq(t, `{"delayMs":300}`, storedAdvanceConfig(t, id))
}

func TestUpdateProjectAdvanceConfig_InvalidCorsRejectedAndKeepsStored(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setup, err := database.InitTestWorkspaceWithProject("adv-bad@example.com", "Adv Bad", "Adv Bad WS", "Adv Bad Project", generateUniqueAlias("adv-bad"))
	require.NoError(t, err)
	defer setup.Cleanup()
	id := setup.Project.ID

	require.Equal(t, http.StatusOK, putAdvanceConfig(t, id, `{"cors":{"mode":"reflect"}}`).Code)

	w := putAdvanceConfig(t, id, `{"cors":{"mode":"nonsense"}}`)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"cors":{"mode":"reflect"}}`, storedAdvanceConfig(t, id))
}

func TestUpdateProjectAdvanceConfig_DeepMergesCors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setup, err := database.InitTestWorkspaceWithProject("adv-deep@example.com", "Adv Deep", "Adv Deep WS", "Adv Deep Project", generateUniqueAlias("adv-deep"))
	require.NoError(t, err)
	defer setup.Cleanup()
	id := setup.Project.ID

	require.Equal(t, http.StatusOK, putAdvanceConfig(t, id, `{"cors":{"mode":"custom","allowHeaders":["X-Api-Key"]}}`).Code)

	// Only maxAge is sent: mode and allowHeaders must survive.
	require.Equal(t, http.StatusOK, putAdvanceConfig(t, id, `{"cors":{"maxAge":600}}`).Code)
	assert.JSONEq(t, `{"cors":{"mode":"custom","allowHeaders":["X-Api-Key"],"maxAge":600}}`, storedAdvanceConfig(t, id))

	// Arrays are replaced, not appended.
	require.Equal(t, http.StatusOK, putAdvanceConfig(t, id, `{"cors":{"allowHeaders":["X-Trace"]}}`).Code)
	assert.JSONEq(t, `{"cors":{"mode":"custom","allowHeaders":["X-Trace"],"maxAge":600}}`, storedAdvanceConfig(t, id))

	// A nested null removes just that key.
	require.Equal(t, http.StatusOK, putAdvanceConfig(t, id, `{"cors":{"maxAge":null}}`).Code)
	assert.JSONEq(t, `{"cors":{"mode":"custom","allowHeaders":["X-Trace"]}}`, storedAdvanceConfig(t, id))
}
