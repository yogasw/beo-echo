package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExecuteProxyRequest_DropsForwardedHeaders makes sure the headers that only
// describe the hop between the client and beo-echo never reach the target.
func TestExecuteProxyRequest_DropsForwardedHeaders(t *testing.T) {
	// Given - a target that records what it received
	var got http.Header
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	req := httptest.NewRequest(http.MethodGet, "http://beo-echo.local/api/v1/calendars", nil)
	req.Header.Set("X-Forwarded-Proto", "http")
	req.Header.Set("X-Forwarded-Host", "beo-echo.local")
	req.Header.Set("X-Forwarded-Port", "443")
	req.Header.Set("X-Real-Ip", "10.0.0.1")
	req.Header.Set("Cdn-Loop", "cloudflare; loops=1")
	req.Header.Set("Referer", "http://beo-echo.local/api/v1/calendars")
	req.Header.Set("X-Api-Key", "secret")
	req.Header.Set("X-Forwarded-For", "203.0.113.7")

	// When - the request is forwarded
	resp, err := executeProxyRequest(context.Background(), target.URL, http.MethodGet, "/api/v1/calendars", "", req)

	// Then - hop headers are gone, the caller's own headers survive
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	for _, h := range []string{"X-Forwarded-Proto", "X-Forwarded-Host", "X-Forwarded-Port", "X-Real-Ip", "Cdn-Loop", "Referer"} {
		assert.Empty(t, got.Get(h), "%s must not be forwarded to the target", h)
	}
	assert.Equal(t, "secret", got.Get("X-Api-Key"))
	assert.Equal(t, "203.0.113.7", got.Get("X-Forwarded-For"))
}

// TestExecuteProxyRequest_DoesNotFollowRedirects covers the loop a TLS-enforcing
// target used to cause: it answers 301 to the same URL, the client followed it
// ten times and the whole request failed with "stopped after 10 redirects".
func TestExecuteProxyRequest_DoesNotFollowRedirects(t *testing.T) {
	// Given - a target that always redirects to itself
	hits := 0
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Location", "https://example.test"+r.URL.Path)
		w.WriteHeader(http.StatusMovedPermanently)
	}))
	defer target.Close()

	req := httptest.NewRequest(http.MethodGet, "http://beo-echo.local/api/v1/calendars", nil)

	// When - the request is forwarded
	resp, err := executeProxyRequest(context.Background(), target.URL, http.MethodGet, "/api/v1/calendars", "", req)

	// Then - the 301 is handed back untouched instead of being chased
	require.NoError(t, err)
	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode)
	assert.Equal(t, "https://example.test/api/v1/calendars", resp.Header.Get("Location"))
	assert.Equal(t, 1, hits, "target must be hit exactly once")
}
