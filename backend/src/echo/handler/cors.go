package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/lib"
)

// StaticAllowHeaders is the fixed Access-Control-Allow-Headers list used by the
// server-wide CORS middleware and by projects in "static" CORS mode.
var StaticAllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "Accept"}

// StaticAllowMethods is the Access-Control-Allow-Methods list of the server.
var StaticAllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// IsMockPreflight reports whether the request is a CORS preflight aimed at a mock
// route. Admin routes (/api, /mcp and the root) stay with the server-wide CORS middleware.
func IsMockPreflight(r *http.Request) bool {
	if r.Method != http.MethodOptions || r.Header.Get("Origin") == "" || r.Header.Get("Access-Control-Request-Method") == "" {
		return false
	}
	// On a mirror instance every path belongs to the project
	if lib.MIRROR_PROJECT_ALIAS != "" {
		return true
	}
	p := r.URL.Path
	if p == "/" || p == "/api" || strings.HasPrefix(p, "/api/") || p == "/mcp" || strings.HasPrefix(p, "/mcp/") {
		return false
	}
	return true
}

// SkipMockPreflight wraps the server-wide CORS middleware so it ignores preflights aimed at
// mock routes; MockRequestHandler answers those with the project's own CORS config.
func SkipMockPreflight(global gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsMockPreflight(c.Request) {
			c.Next()
			return
		}
		global(c)
	}
}

// corsConfigOf returns the CORS config of a project; nil project or empty config means static.
func corsConfigOf(project *database.Project) database.CorsConfig {
	if project == nil || project.AdvanceConfig == "" {
		return database.CorsConfig{}
	}
	cfg, err := database.ParseProjectAdvanceConfig(project.AdvanceConfig)
	if err != nil || cfg.Cors == nil {
		return database.CorsConfig{}
	}
	return *cfg.Cors
}

// allowedOrigin resolves the Access-Control-Allow-Origin value for the request origin.
// ok=false means the origin is not allowed.
func allowedOrigin(cfg database.CorsConfig, origin string) (value string, ok bool) {
	origins := cfg.AllowOrigins
	if len(origins) == 0 {
		origins = []string{lib.CORS_ORIGIN}
	}
	for _, o := range origins {
		if o == "*" {
			return "*", true
		}
		if strings.EqualFold(o, origin) {
			return origin, true
		}
	}
	return "", false
}

// preflightAllowHeaders picks Access-Control-Allow-Headers for the project's CORS mode.
func preflightAllowHeaders(cfg database.CorsConfig, requested string) string {
	switch cfg.Mode {
	case database.CorsModeReflect:
		if strings.TrimSpace(requested) != "" {
			return requested
		}
	case database.CorsModeCustom:
		return strings.Join(cfg.AllowHeaders, ", ")
	}
	return strings.Join(StaticAllowHeaders, ", ")
}

// HandleMockPreflight answers a CORS preflight for a mock route using the project's CORS
// config. An unknown project (nil) is answered like the server-wide static config.
func HandleMockPreflight(c *gin.Context, project *database.Project) {
	cfg := corsConfigOf(project)
	origin := c.Request.Header.Get("Origin")

	allowOrigin, ok := allowedOrigin(cfg, origin)
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	h := c.Writer.Header()
	h.Set("Access-Control-Allow-Origin", allowOrigin)
	h.Set("Access-Control-Allow-Methods", strings.Join(StaticAllowMethods, ", "))
	if allowHeaders := preflightAllowHeaders(cfg, c.Request.Header.Get("Access-Control-Request-Headers")); allowHeaders != "" {
		h.Set("Access-Control-Allow-Headers", allowHeaders)
	}
	if cfg.MaxAge > 0 {
		h.Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
	}
	h.Add("Vary", "Origin")
	h.Add("Vary", "Access-Control-Request-Headers")
	c.AbortWithStatus(http.StatusNoContent)
}
