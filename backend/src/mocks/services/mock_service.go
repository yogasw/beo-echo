package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/repositories"
)

// MockService handles mock response logic
type MockService struct {
	Repo *repositories.MockRepository
}

// NewMockService creates a new mock service
func NewMockService(repo *repositories.MockRepository) *MockService {
	return &MockService{
		Repo: repo,
	}
}

// HandleRequest processes an incoming request and returns a mock response or proxies it
func (s *MockService) HandleRequest(alias, method, path string, req *http.Request) (*http.Response, error) {
	// Find project by alias
	project, err := s.Repo.FindProjectByAlias(alias)
	if err != nil {
		return createErrorResponse(http.StatusNotFound, "Project not found"), nil
	}

	// Extract the actual API endpoint path
	// Path comes in like "/api/users" or "/users" - we need just the endpoint part
	// First trim any project alias prefix if it exists
	cleanPath := strings.TrimPrefix(path, "/"+project.Alias)

	// Check project mode
	switch project.Mode {
	case database.ModeMock:
		return s.handleMockMode(project.ID, method, cleanPath, req)
	case database.ModeProxy, database.ModeForwarder:
		return s.handleProxyMode(project, req)
	case database.ModeDisabled:
		return createErrorResponse(http.StatusServiceUnavailable, "Service is disabled"), nil
	default:
		return createErrorResponse(http.StatusInternalServerError, "Invalid project mode"), nil
	}
}

// handleMockMode generates mock response
func (s *MockService) handleMockMode(projectID string, method, path string, req *http.Request) (*http.Response, error) {
	endpoint, err := s.Repo.FindMatchingEndpoint(projectID, method, path)
	if err != nil {
		return createErrorResponse(http.StatusNotFound, "Endpoint not found"), nil
	}

	// Get all responses for this endpoint
	responses, err := s.Repo.FindResponsesByEndpointID(endpoint.ID)
	if err != nil || len(responses) == 0 {
		return createErrorResponse(http.StatusInternalServerError, "No responses configured"), nil
	}

	// Select response based on ResponseMode
	response := selectResponse(responses, endpoint.ResponseMode, req)

	// Apply delay if configured
	if response.DelayMS > 0 {
		time.Sleep(time.Duration(response.DelayMS) * time.Millisecond)
	}

	// Create and return HTTP response
	return createMockResponse(response)
}

// handleProxyMode forwards the request to target
func (s *MockService) handleProxyMode(project *database.Project, req *http.Request) (*http.Response, error) {
	if project.ActiveProxy == nil {
		return createErrorResponse(http.StatusInternalServerError, "No proxy target configured"), nil
	}

	targetURL, err := url.Parse(project.ActiveProxy.URL)
	if err != nil {
		return createErrorResponse(http.StatusInternalServerError, "Invalid proxy URL"), nil
	}

	// Create proxy director
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		originalDirector(r)
		r.Host = targetURL.Host
	}

	// Execute the request
	resp, err := executeProxyRequest(req, proxy)
	if err != nil {
		return createErrorResponse(http.StatusBadGateway, "Proxy error: "+err.Error()), nil
	}

	// If in forwarder mode, record the response (implement if needed)
	if project.Mode == database.ModeForwarder {
		// TODO: Record response for later use
	}

	return resp, nil
}

// Helper functions

// selectResponse selects a response based on mode and rules
func selectResponse(responses []database.MockResponse, mode string, req *http.Request) database.MockResponse {
	// Filter responses by rules first
	validResponses := filterResponsesByRules(responses, req)
	if len(validResponses) == 0 {
		// If no responses match rules, fallback to all responses
		validResponses = responses
	}

	// Sort by priority (higher is more important)
	sortByPriority(validResponses)

	// Select based on mode
	switch strings.ToLower(mode) {
	case "static":
		// Return highest priority
		return validResponses[0]
	case "random":
		// Return random response
		return validResponses[rand.Intn(len(validResponses))]
	case "round_robin":
		// TODO: Implement round-robin selection
		// For now, return first
		return validResponses[0]
	default:
		// Default to random
		return validResponses[rand.Intn(len(validResponses))]
	}
}

// filterResponsesByRules filters responses that match request rules
func filterResponsesByRules(responses []database.MockResponse, req *http.Request) []database.MockResponse {
	if req == nil {
		return responses
	}

	var validResponses []database.MockResponse
	for _, resp := range responses {
		if matchesRules(resp, req) {
			validResponses = append(validResponses, resp)
		}
	}

	return validResponses
}

// sortByPriority sorts responses by priority (higher first)
func sortByPriority(responses []database.MockResponse) {
	// Simple bubble sort
	for i := 0; i < len(responses)-1; i++ {
		for j := 0; j < len(responses)-i-1; j++ {
			if responses[j].Priority < responses[j+1].Priority {
				responses[j], responses[j+1] = responses[j+1], responses[j]
			}
		}
	}
}

// matchesRules checks if a response matches all rules against the request
func matchesRules(response database.MockResponse, req *http.Request) bool {
	if len(response.Rules) == 0 {
		return true // No rules means always match
	}

	for _, rule := range response.Rules {
		switch rule.Type {
		case "header":
			if !matchHeaderRule(rule, req) {
				return false
			}
		case "query":
			if !matchQueryRule(rule, req) {
				return false
			}
		case "body":
			if !matchBodyRule(rule, req) {
				return false
			}
		}
		// Path rules are handled earlier during endpoint matching
	}

	return true
}

// matchHeaderRule checks if a header rule matches
func matchHeaderRule(rule database.MockRule, req *http.Request) bool {
	headerValue := req.Header.Get(rule.Key)
	return matchRuleValue(rule.Operator, headerValue, rule.Value)
}

// matchQueryRule checks if a query parameter rule matches
func matchQueryRule(rule database.MockRule, req *http.Request) bool {
	queryValue := req.URL.Query().Get(rule.Key)
	return matchRuleValue(rule.Operator, queryValue, rule.Value)
}

// matchBodyRule checks if a body rule matches
func matchBodyRule(rule database.MockRule, req *http.Request) bool {
	// Get body content (this is a simplistic approach; in real-world you'd want to cache)
	if req.Body == nil {
		return false
	}

	// Read body
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return false
	}

	// Restore body for subsequent reads
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// For JSON bodies, try to extract nested values
	var bodyData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &bodyData); err == nil {
		// Try to find nested value
		if value := getNestedValue(bodyData, rule.Key); value != "" {
			return matchRuleValue(rule.Operator, value, rule.Value)
		}
	}

	// Fallback to treating body as string
	return matchRuleValue(rule.Operator, string(bodyBytes), rule.Value)
}

// matchRuleValue compares values based on operator
func matchRuleValue(operator, actual, expected string) bool {
	switch strings.ToLower(operator) {
	case "equals":
		return actual == expected
	case "contains":
		return strings.Contains(actual, expected)
	case "regex":
		// TODO: Implement regex matching
		// For simplicity, fallback to contains for now
		return strings.Contains(actual, expected)
	default:
		return actual == expected // Default to equals
	}
}

// getNestedValue extracts a value from nested JSON using dot notation
// e.g., "user.address.city" => data["user"]["address"]["city"]
func getNestedValue(data map[string]interface{}, key string) string {
	parts := strings.Split(key, ".")
	var current interface{} = data

	for _, part := range parts {
		// Try to navigate the object
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[part]
		default:
			return "" // Can't navigate further
		}

		if current == nil {
			return ""
		}
	}

	// Convert result to string
	switch v := current.(type) {
	case string:
		return v
	case bool, int, float64:
		return fmt.Sprintf("%v", v)
	default:
		bytes, err := json.Marshal(current)
		if err != nil {
			return ""
		}
		return string(bytes)
	}
}

// createMockResponse builds an HTTP response from a mock response
func createMockResponse(mockResp database.MockResponse) (*http.Response, error) {
	// Create response body
	body := io.NopCloser(strings.NewReader(mockResp.Body))

	// Create response
	resp := &http.Response{
		StatusCode: mockResp.StatusCode,
		Body:       body,
		Header:     make(http.Header),
	}

	// Add headers
	headers, err := ParseHeaders(mockResp.Headers)
	if err == nil {
		for key, value := range headers {
			resp.Header.Set(key, value)
		}
	}

	// Set content length
	resp.ContentLength = int64(len(mockResp.Body))

	return resp, nil
}

// createErrorResponse creates a standard error response
func createErrorResponse(statusCode int, message string) *http.Response {
	respBody := map[string]interface{}{
		"error":   true,
		"message": message,
	}

	jsonBody, _ := json.Marshal(respBody)
	body := io.NopCloser(bytes.NewBuffer(jsonBody))

	resp := &http.Response{
		StatusCode:    statusCode,
		Body:          body,
		Header:        make(http.Header),
		ContentLength: int64(len(jsonBody)),
	}

	resp.Header.Set("Content-Type", "application/json")
	return resp
}

// executeProxyRequest executes a request through a proxy
func executeProxyRequest(originalReq *http.Request, proxy *httputil.ReverseProxy) (*http.Response, error) {
	// Create in-memory pipe
	pr, pw := io.Pipe()

	// Create a response recorder
	recorderResp := &responseRecorder{
		headers: make(http.Header),
		pipe:    pw,
	}

	// Use goroutine to execute proxy request
	go func() {
		proxy.ServeHTTP(recorderResp, originalReq)
		pw.Close()
	}()

	// Create response
	resp := &http.Response{
		StatusCode:    recorderResp.statusCode,
		Header:        recorderResp.headers,
		Body:          pr,
		ContentLength: -1, // Unknown content length
	}

	return resp, nil
}

// responseRecorder implements a ResponseWriter to capture response details
type responseRecorder struct {
	headers    http.Header
	statusCode int
	pipe       *io.PipeWriter
}

func (r *responseRecorder) Header() http.Header {
	return r.headers
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.statusCode == 0 {
		r.statusCode = http.StatusOK // Default status
	}
	return r.pipe.Write(b)
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

// ParseHeaders converts a JSON string to a map of headers
func ParseHeaders(headersJSON string) (map[string]string, error) {
	headers := make(map[string]string)
	if headersJSON == "" {
		return headers, nil
	}

	err := json.Unmarshal([]byte(headersJSON), &headers)
	if err != nil {
		return nil, err
	}

	return headers, nil
}
