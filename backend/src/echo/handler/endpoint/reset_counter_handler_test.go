package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
)

func TestResetCounterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setup, err := database.InitTestWorkspaceWithProject("reset@example.com", "Reset", "Reset WS", "Reset Project", "reset-counter-project")
	require.NoError(t, err)
	defer setup.Cleanup()

	ep, err := database.CreateTestEndpoint(setup.Project.ID, "GET", "/flaky")
	require.NoError(t, err)

	router := gin.New()
	router.POST("/api/projects/:projectId/endpoints/:id/reset-counter", ResetCounterHandler)

	t.Run("existing endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/projects/"+setup.Project.ID+"/endpoints/"+ep.ID+"/reset-counter", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("endpoint of another project is not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/projects/other-project/endpoints/"+ep.ID+"/reset-counter", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
