// Integration tests: real gin router (server-wide CORS middleware + mock catch-all route +
// admin handlers) on top of the real mock service and the test database. They replay the
// QA Omniform scenarios: custom-header preflight, templated bodies and a flaky endpoint.
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
	"beo-echo/backend/src/echo/handler/endpoint"
	"beo-echo/backend/src/echo/handler/project"
	"beo-echo/backend/src/echo/handler/response"
)

func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(handler.SkipMockPreflight(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: handler.StaticAllowMethods,
		AllowHeaders: handler.StaticAllowHeaders,
	})))

	admin := router.Group("/api/projects/:projectId")
	admin.PUT("/advance-config", project.UpdateProjectAdvanceConfigHandler)
	admin.PUT("/endpoints/:id", endpoint.UpdateEndpointHandler)
	admin.POST("/endpoints/:id/reset-counter", endpoint.ResetCounterHandler)
	admin.PUT("/endpoints/:id/responses/:responseId", response.UpdateResponseHandler)

	router.Any("/:project/*path", handler.MockRequestHandler)
	return router
}

func do(router *gin.Engine, method, target string, headers map[string]string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

type fixture struct {
	router   *gin.Engine
	project  *database.Project
	alias    string
	cleanup  func()
	endpoint func(method, path, mode, advanceConfig string, responses ...database.MockResponse) *database.MockEndpoint
}

func newFixture(t *testing.T, alias, projectAdvanceConfig string) *fixture {
	t.Helper()
	setup, err := database.InitTestWorkspaceWithProject(alias+"@example.com", alias, alias+" ws", alias, alias)
	require.NoError(t, err)

	require.NoError(t, database.DB.Model(&database.Project{}).Where("id = ?", setup.Project.ID).
		Update("advance_config", projectAdvanceConfig).Error)

	f := &fixture{router: newRouter(), project: setup.Project, alias: alias, cleanup: setup.Cleanup}
	f.endpoint = func(method, path, mode, advanceConfig string, responses ...database.MockResponse) *database.MockEndpoint {
		ep, err := database.CreateTestEndpointWithConfig(setup.Project.ID, method, path, advanceConfig)
		require.NoError(t, err)
		require.NoError(t, database.DB.Model(ep).Update("response_mode", mode).Error)
		for i := range responses {
			responses[i].EndpointID = ep.ID
			responses[i].Enabled = true
			if responses[i].Headers == "" {
				responses[i].Headers = `{"Content-Type":"application/json"}`
			}
			require.NoError(t, database.DB.Create(&responses[i]).Error)
		}
		return ep
	}
	return f
}

func TestIntegration_PreflightReflectsCustomHeaders(t *testing.T) {
	f := newFixture(t, "it-cors-reflect", `{"cors":{"mode":"reflect"}}`)
	defer f.cleanup()
	f.endpoint("GET", "/ok", "static", "", database.MockResponse{StatusCode: 200, Body: `{"ok":true}`})

	w := do(f.router, http.MethodOptions, "/it-cors-reflect/ok", map[string]string{
		"Origin":                         "https://omnichannel.qiscus.com",
		"Access-Control-Request-Method":  "GET",
		"Access-Control-Request-Headers": "x-api-key",
	}, "")

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "x-api-key", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))

	actual := do(f.router, http.MethodGet, "/it-cors-reflect/ok", map[string]string{"Origin": "https://omnichannel.qiscus.com", "X-Api-Key": "k"}, "")
	assert.Equal(t, http.StatusOK, actual.Code)
	assert.Equal(t, "*", actual.Header().Get("Access-Control-Allow-Origin"))
	assert.JSONEq(t, `{"ok":true}`, actual.Body.String())
}

func TestIntegration_LegacyProjectKeepsStaticPreflight(t *testing.T) {
	f := newFixture(t, "it-cors-legacy", "")
	defer f.cleanup()
	f.endpoint("GET", "/ok", "static", "", database.MockResponse{StatusCode: 200, Body: `{}`})

	w := do(f.router, http.MethodOptions, "/it-cors-legacy/ok", map[string]string{
		"Origin":                         "https://omnichannel.qiscus.com",
		"Access-Control-Request-Method":  "GET",
		"Access-Control-Request-Headers": "x-api-key",
	}, "")

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "Origin, Content-Length, Content-Type, Authorization, X-Requested-With, Accept", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestIntegration_NoCorsOverrideStillWorks(t *testing.T) {
	f := newFixture(t, "it-nocors", `{"cors":{"mode":"reflect"}}`)
	defer f.cleanup()
	f.endpoint("GET", "/nocors", "static", "", database.MockResponse{
		StatusCode: 200, Body: `{}`,
		Headers: `{"Content-Type":"application/json","Access-Control-Allow-Origin":""}`,
	})

	w := do(f.router, http.MethodGet, "/it-nocors/nocors", map[string]string{"Origin": "https://omnichannel.qiscus.com"}, "")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"), "per-response override must win over the CORS middleware")
}

func TestIntegration_SettingDelayDoesNotWipeCors(t *testing.T) {
	f := newFixture(t, "it-merge", "")
	defer f.cleanup()
	f.endpoint("GET", "/ok", "static", "", database.MockResponse{StatusCode: 200, Body: `{}`})

	put := func(body string) {
		w := do(f.router, http.MethodPut, "/api/projects/"+f.project.ID+"/advance-config", map[string]string{"Content-Type": "application/json"}, body)
		require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	}
	put(`{"cors":{"mode":"reflect"}}`)
	put(`{"delayMs":0}`)

	w := do(f.router, http.MethodOptions, "/it-merge/ok", map[string]string{
		"Origin":                         "https://omnichannel.qiscus.com",
		"Access-Control-Request-Method":  "GET",
		"Access-Control-Request-Headers": "x-api-key",
	}, "")
	assert.Equal(t, "x-api-key", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestIntegration_TemplatedBodies(t *testing.T) {
	f := newFixture(t, "it-template", "")
	defer f.cleanup()
	f.endpoint("GET", "/log", "static", "", database.MockResponse{
		StatusCode: 200, Templated: true,
		Body: `[{"label":"room_id={{queryParam 'room_id'}}","value":"{{queryParam 'app_id'}}"}]`,
	})
	f.endpoint("GET", "/many", "static", "", database.MockResponse{
		StatusCode: 200, Templated: true,
		Body: `[{{#repeat (queryParam 'n') comma=true}}{"label":"Opsi {{@index}}","value":"OPT-{{@index}}"}{{/repeat}}]`,
	})
	f.endpoint("GET", "/raw", "static", "", database.MockResponse{
		StatusCode: 200, Body: `{"x":"{{queryParam 'room_id'}}"}`,
	})
	f.endpoint("GET", "/broken", "static", "", database.MockResponse{
		StatusCode: 200, Templated: true, Body: `{{queryParms 'room_id'}}`,
	})

	log := do(f.router, http.MethodGet, "/it-template/log?room_id=123&app_id=abc", nil, "")
	assert.Equal(t, http.StatusOK, log.Code)
	assert.JSONEq(t, `[{"label":"room_id=123","value":"abc"}]`, log.Body.String())

	many := do(f.router, http.MethodGet, "/it-template/many?n=500", nil, "")
	assert.Equal(t, http.StatusOK, many.Code)
	var items []map[string]string
	require.NoError(t, json.Unmarshal(many.Body.Bytes(), &items))
	require.Len(t, items, 500)
	assert.Equal(t, "OPT-499", items[499]["value"])

	raw := do(f.router, http.MethodGet, "/it-template/raw?room_id=123", nil, "")
	assert.Equal(t, `{"x":"{{queryParam 'room_id'}}"}`, raw.Body.String(), "not templated = untouched")

	broken := do(f.router, http.MethodGet, "/it-template/broken", nil, "")
	assert.Equal(t, http.StatusInternalServerError, broken.Code)
	assert.Contains(t, broken.Body.String(), "unknown template helper")
}

func TestIntegration_FlakyEndpointSequenceAndReset(t *testing.T) {
	f := newFixture(t, "it-flaky", "")
	defer f.cleanup()
	ep := f.endpoint("GET", "/flaky", "sequence", "",
		database.MockResponse{StatusCode: 500, Priority: 2, Body: `{"error":"boom"}`},
		database.MockResponse{StatusCode: 200, Priority: 1, Body: `{"ok":true}`},
	)

	codes := func(n int) []int {
		var out []int
		for i := 0; i < n; i++ {
			out = append(out, do(f.router, http.MethodGet, "/it-flaky/flaky", nil, "").Code)
		}
		return out
	}

	assert.Equal(t, []int{500, 200, 200, 200}, codes(4), "first call fails, then it keeps succeeding")

	reset := do(f.router, http.MethodPost, "/api/projects/"+f.project.ID+"/endpoints/"+ep.ID+"/reset-counter", nil, "")
	require.Equal(t, http.StatusOK, reset.Code)
	assert.Equal(t, []int{500, 200}, codes(2), "after reset it starts from the first response again")

	// loop mode comes from the endpoint advance config; changing it also resets the counter
	upd := do(f.router, http.MethodPut, "/api/projects/"+f.project.ID+"/endpoints/"+ep.ID,
		map[string]string{"Content-Type": "application/json"}, `{"advance_config":"{\"sequence\":{\"onEnd\":\"loop\"}}"}`)
	require.Equal(t, http.StatusOK, upd.Code, upd.Body.String())
	assert.Equal(t, []int{500, 200, 500, 200}, codes(4))

	// editing a response of the endpoint resets the counter too
	var first database.MockResponse
	require.NoError(t, database.DB.Where("endpoint_id = ? AND status_code = ?", ep.ID, 200).First(&first).Error)
	assert.Equal(t, 500, do(f.router, http.MethodGet, "/it-flaky/flaky", nil, "").Code)
	edit := do(f.router, http.MethodPut, "/api/projects/"+f.project.ID+"/endpoints/"+ep.ID+"/responses/"+first.ID,
		map[string]string{"Content-Type": "application/json"}, `{"note":"edited"}`)
	require.Equal(t, http.StatusOK, edit.Code, edit.Body.String())
	assert.Equal(t, 500, do(f.router, http.MethodGet, "/it-flaky/flaky", nil, "").Code, "counter restarted after the edit")
}
