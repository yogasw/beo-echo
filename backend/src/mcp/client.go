// Package mcp implements a Model Context Protocol server for Beo Echo.
//
// The server is a thin client over the existing Beo Echo REST API: every tool
// forwards the caller's bearer token (a PAT or OAuth-issued token) to the API,
// so all authorization and permission checks happen in one place — the API
// itself. No business logic is duplicated here.
package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// Client invokes the Beo Echo REST API in-process by serving requests directly
// through the API's http.Handler (the same Gin router), so calls pass through
// all the API middleware (auth, permission checks, logging) without a network
// round-trip and without needing to know the server's host/port. It is
// stateless with respect to auth: the bearer token is passed per call, so a
// single Client can serve many concurrent MCP sessions for different users.
type Client struct {
	handler http.Handler
}

// NewClient builds a REST client that dispatches to the given API handler.
func NewClient(handler http.Handler) *Client {
	return &Client{handler: handler}
}

// APIError represents a non-2xx response from the API, carrying the status code
// and the server's message so tools can surface it usefully to the model.
type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Body)
}

// do performs a request to the API, attaching the bearer token, and decodes the
// JSON response into out (which may be nil to ignore the body). query is
// optional and appended to the path.
func (c *Client) do(ctx context.Context, token, method, path string, query url.Values, body any, out any) error {
	if token == "" {
		return &APIError{StatusCode: http.StatusUnauthorized, Message: "missing bearer token; configure your MCP client with an Authorization header or connect via OAuth"}
	}

	u := path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var reader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reader = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Serve the request in-process through the API handler.
	rec := httptest.NewRecorder()
	c.handler.ServeHTTP(rec, req)

	raw := rec.Body.Bytes()

	if rec.Code < 200 || rec.Code >= 300 {
		return &APIError{
			StatusCode: rec.Code,
			Message:    extractMessage(raw),
			Body:       string(raw),
		}
	}

	if out != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("decode response: %w (body: %s)", err, truncate(string(raw), 500))
		}
	}
	return nil
}

// Get issues a GET and decodes the result.
func (c *Client) Get(ctx context.Context, token, path string, query url.Values, out any) error {
	return c.do(ctx, token, http.MethodGet, path, query, nil, out)
}

// Post issues a POST with a JSON body.
func (c *Client) Post(ctx context.Context, token, path string, body, out any) error {
	return c.do(ctx, token, http.MethodPost, path, nil, body, out)
}

// Put issues a PUT with a JSON body.
func (c *Client) Put(ctx context.Context, token, path string, body, out any) error {
	return c.do(ctx, token, http.MethodPut, path, nil, body, out)
}

// Patch issues a PATCH with a JSON body.
func (c *Client) Patch(ctx context.Context, token, path string, body, out any) error {
	return c.do(ctx, token, http.MethodPatch, path, nil, body, out)
}

// Delete issues a DELETE.
func (c *Client) Delete(ctx context.Context, token, path string, out any) error {
	return c.do(ctx, token, http.MethodDelete, path, nil, nil, out)
}

// extractMessage pulls a "message" or "error" field out of a JSON error body.
func extractMessage(raw []byte) string {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return ""
	}
	if v, ok := m["message"].(string); ok && v != "" {
		return v
	}
	if v, ok := m["error"].(string); ok && v != "" {
		return v
	}
	if v, ok := m["error_description"].(string); ok && v != "" {
		return v
	}
	return ""
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
