package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
)

var cloudflareRequestHeaders = map[string]string{
	"CF-Connecting-IP": "203.0.113.7",
	"CF-Ray":           "8f1a2b3c4d5e6f70-SIN",
	"CF-IPCountry":     "ID",
	"CF-Visitor":       `{"scheme":"https"}`,
	"True-Client-IP":   "203.0.113.7",
	"CDN-Loop":         "cloudflare; loops=1",
}

func withHeaders(extra map[string]string) map[string]string {
	h := map[string]string{"X-Api-Key": "secret"}
	for k, v := range cloudflareRequestHeaders {
		h[k] = v
	}
	for k, v := range extra {
		h[k] = v
	}
	return h
}

func assertNoCloudflareHeaders(t *testing.T, got http.Header) {
	t.Helper()
	require.NotNil(t, got, "upstream was not called")
	for name := range cloudflareRequestHeaders {
		assert.Empty(t, got.Values(name), "%s must not reach the upstream", name)
	}
	assert.Equal(t, "secret", got.Get("X-Api-Key"), "the caller's own headers still reach the upstream")
}

// The built-in proxy must not forward Cloudflare headers: an upstream that sits behind Cloudflare
// itself rejects CF-Connecting-IP with error 1000 (double proxy).
func TestIntegration_ProxyDoesNotForwardCloudflareHeaders(t *testing.T) {
	u := newUpstream(t)
	f := proxyProject(t, "it-proxy-cf", u)
	defer f.cleanup()

	w := do(f.router, http.MethodGet, "/it-proxy-cf/ping", withHeaders(nil), "")

	assert.Equal(t, http.StatusOK, w.Code)
	assertNoCloudflareHeaders(t, u.received())
	assert.Equal(t, upstreamHTML, decodeBrotli(t, w.Body.Bytes()), "an untouched response passes through intact")
}

// A before_request script that blanks the CF headers works on top of the built-in strip.
func TestIntegration_BeforeRequestScriptDoesNotBreakProxy(t *testing.T) {
	u := newUpstream(t)
	f := proxyProject(t, "it-proxy-before", u)
	defer f.cleanup()
	addScriptAction(t, f.project.ID, database.ExecutionPointBeforeRequest, `
		for (var k in request.headers) { if (k.toLowerCase().indexOf("cf-") === 0) { request.headers[k] = ""; } }
	`)

	w := do(f.router, http.MethodGet, "/it-proxy-before/ping", withHeaders(nil), "")

	assert.Equal(t, http.StatusOK, w.Code)
	assertNoCloudflareHeaders(t, u.received())
	assert.Equal(t, upstreamHTML, decodeBrotli(t, w.Body.Bytes()))
}
