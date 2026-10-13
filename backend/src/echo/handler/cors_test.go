package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"beo-echo/backend/src/database"
)

func preflightRequest(path, requestHeaders string) *http.Request {
	req := httptest.NewRequest(http.MethodOptions, path, nil)
	req.Header.Set("Origin", "https://omnichannel.qiscus.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	if requestHeaders != "" {
		req.Header.Set("Access-Control-Request-Headers", requestHeaders)
	}
	return req
}

func TestIsMockPreflight(t *testing.T) {
	assert.True(t, IsMockPreflight(preflightRequest("/omniform-mock-stag/ok", "x-api-key")))
	assert.False(t, IsMockPreflight(preflightRequest("/api/projects", "x-api-key")), "admin API keeps the global middleware")
	assert.False(t, IsMockPreflight(preflightRequest("/mcp", "")), "mcp keeps the global middleware")

	get := preflightRequest("/omniform-mock-stag/ok", "")
	get.Method = http.MethodGet
	assert.False(t, IsMockPreflight(get))

	noOrigin := preflightRequest("/omniform-mock-stag/ok", "")
	noOrigin.Header.Del("Origin")
	assert.False(t, IsMockPreflight(noOrigin))
}

func TestPreflightAllowHeaders(t *testing.T) {
	static := "Origin, Content-Length, Content-Type, Authorization, X-Requested-With, Accept"
	tests := []struct {
		name      string
		cfg       database.CorsConfig
		requested string
		want      string
	}{
		{"empty config is static", database.CorsConfig{}, "x-api-key", static},
		{"static ignores request", database.CorsConfig{Mode: database.CorsModeStatic}, "x-api-key", static},
		{"reflect echoes request", database.CorsConfig{Mode: database.CorsModeReflect}, "x-api-key", "x-api-key"},
		{"reflect echoes several", database.CorsConfig{Mode: database.CorsModeReflect}, "x-api-key, content-type", "x-api-key, content-type"},
		{"reflect without request falls back to static", database.CorsConfig{Mode: database.CorsModeReflect}, "", static},
		{"custom uses list", database.CorsConfig{Mode: database.CorsModeCustom, AllowHeaders: []string{"X-Api-Key", "X-Trace"}}, "other", "X-Api-Key, X-Trace"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, preflightAllowHeaders(tc.cfg, tc.requested))
		})
	}
}

func TestAllowedOrigin(t *testing.T) {
	v, ok := allowedOrigin(database.CorsConfig{AllowOrigins: []string{"https://a.example"}}, "https://a.example")
	assert.True(t, ok)
	assert.Equal(t, "https://a.example", v)

	_, ok = allowedOrigin(database.CorsConfig{AllowOrigins: []string{"https://a.example"}}, "https://b.example")
	assert.False(t, ok)

	v, ok = allowedOrigin(database.CorsConfig{AllowOrigins: []string{"*"}}, "https://b.example")
	assert.True(t, ok)
	assert.Equal(t, "*", v)
}

// newCORSTestRouter wires the server-wide CORS middleware (with the mock preflight bypass) in
// front of a catch-all mock route, the same way server.go does.
func newCORSTestRouter(project *database.Project) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	global := cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: StaticAllowMethods,
		AllowHeaders: StaticAllowHeaders,
	})
	router.Use(SkipMockPreflight(global))
	router.Any("/:project/*path", func(c *gin.Context) {
		if IsMockPreflight(c.Request) {
			HandleMockPreflight(c, project)
			return
		}
		// Per-response header override is applied after the middleware, like MockRequestHandler does.
		if c.Param("path") == "/nocors" {
			c.Header("Access-Control-Allow-Origin", "")
		}
		c.String(http.StatusOK, "ok")
	})
	router.GET("/api/health", func(c *gin.Context) { c.String(http.StatusOK, "up") })
	router.OPTIONS("/api/health", func(c *gin.Context) { c.String(http.StatusOK, "up") })
	return router
}

func TestMockPreflight_ReflectProject(t *testing.T) {
	project := &database.Project{AdvanceConfig: `{"cors":{"mode":"reflect","maxAge":600}}`}
	router := newCORSTestRouter(project)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, preflightRequest("/omniform-mock-stag/ok", "x-api-key"))

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "x-api-key", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "600", w.Header().Get("Access-Control-Max-Age"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
}

func TestMockPreflight_ProjectWithoutConfigKeepsStaticList(t *testing.T) {
	router := newCORSTestRouter(&database.Project{})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, preflightRequest("/legacy/ok", "x-api-key"))

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "Origin, Content-Length, Content-Type, Authorization, X-Requested-With, Accept", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestMockPreflight_UnknownProjectAnsweredAsStatic(t *testing.T) {
	router := newCORSTestRouter(nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, preflightRequest("/nope/ok", "x-api-key"))

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	assert.NotContains(t, w.Header().Get("Access-Control-Allow-Headers"), "x-api-key")
}

func TestMockPreflight_DisallowedOrigin(t *testing.T) {
	router := newCORSTestRouter(&database.Project{AdvanceConfig: `{"cors":{"mode":"reflect","allowOrigins":["https://only.example"]}}`})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, preflightRequest("/p/ok", "x-api-key"))

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAdminPreflightStillUsesGlobalMiddleware(t *testing.T) {
	router := newCORSTestRouter(&database.Project{AdvanceConfig: `{"cors":{"mode":"reflect"}}`})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, preflightRequest("/api/health", "x-api-key"))

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.NotContains(t, w.Header().Get("Access-Control-Allow-Headers"), "x-api-key", "admin API must not reflect headers")
}

func TestActualRequest_PerResponseAllowOriginOverrideStillWins(t *testing.T) {
	router := newCORSTestRouter(&database.Project{AdvanceConfig: `{"cors":{"mode":"reflect"}}`})

	normal := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/p/ok", nil)
	req.Header.Set("Origin", "https://omnichannel.qiscus.com")
	router.ServeHTTP(normal, req)
	assert.Equal(t, "*", normal.Header().Get("Access-Control-Allow-Origin"))

	nocors := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/p/nocors", nil)
	req.Header.Set("Origin", "https://omnichannel.qiscus.com")
	router.ServeHTTP(nocors, req)
	assert.Equal(t, "", nocors.Header().Get("Access-Control-Allow-Origin"))
}
