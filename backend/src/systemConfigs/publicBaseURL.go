package systemConfig

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// PublicBaseURL returns the externally reachable base URL of this server
// (scheme://host, no trailing slash).
//
// PUBLIC_BASE_URL wins when it is set. Otherwise the scheme is whatever the
// fronting Caddy reports in X-Forwarded-Proto: it trusts the proxies in front
// of it (see trusted_proxies in the Caddy config), so an https visitor stays
// https and a plain http one stays http.
func PublicBaseURL(c *gin.Context) string {
	if v, err := GetSystemConfigWithType[string](PUBLIC_BASE_URL); err == nil {
		if v = strings.TrimRight(strings.TrimSpace(v), "/"); v != "" {
			return v
		}
	}

	scheme := strings.ToLower(strings.TrimSpace(strings.Split(c.GetHeader("X-Forwarded-Proto"), ",")[0]))
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}
	return scheme + "://" + host
}
