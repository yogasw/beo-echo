package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/repositories"
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
// Also returns project ID, execution mode, and whether the request matched an endpoint
func (s *MockService) HandleRequest(alias, method, reqPath string, req *http.Request) (*http.Response, error, string, database.ProjectMode, bool) {
	// Find project by alias
	project, err := s.Repo.FindProjectByAlias(alias)
	if err != nil {
		return createErrorResponse(http.StatusNotFound, "Project not found"), nil, "", "", false
	}

	// Extract the actual API endpoint path
	// Path comes in like "/api/users" or "/users" - we need just the endpoint part
	// First trim any project alias prefix if it exists
	cleanPath := strings.TrimPrefix(reqPath, "/"+project.Alias)

	// Check project mode
	switch project.Mode {
	case database.ModeMock:
		resp, err, mode, matched := s.handleMockMode(project, method, cleanPath, req)
		return resp, err, project.ID, mode, matched
	case database.ModeProxy:
		resp, err, matched := s.handleProxyMode(project, method, cleanPath, req)
		return resp, err, project.ID, project.Mode, matched // Matched is true only if handled by a mock endpoint
	case database.ModeForwarder:
		resp, err := s.handleForwarderMode(project, method, cleanPath, req)
		return resp, err, project.ID, project.Mode, false // Forwarder requests are always considered "not matched"
	case database.ModeDisabled:
		return createErrorResponse(http.StatusServiceUnavailable, "Service is disabled"), nil, project.ID, project.Mode, false
	default:
		return createErrorResponse(http.StatusInternalServerError, "Invalid project mode"), nil, project.ID, project.Mode, false
	}
}

// handleMockMode generates mock response and returns if the request matched an endpoint
func (s *MockService) handleMockMode(project *database.Project, method, path string, req *http.Request) (*http.Response, error, database.ProjectMode, bool) {
	endpoint, err := s.Repo.FindMatchingEndpoint(project.ID, method, path)
	if err != nil {
		// No matching endpoint found - apply project-level delay before returning error
		s.applyDelay(project, nil, nil)
		return createErrorResponse(http.StatusNotFound, "Endpoint not found"), nil, database.ModeMock, false
	}

	// Check if endpoint is configured for proxying
	if endpoint.UseProxy && endpoint.ProxyTarget != nil {
		// Apply delays before proxying
		s.applyDelay(project, endpoint, nil)
		// Forward the request to the proxy target
		resp, err := executeProxyRequest(endpoint.ProxyTarget.URL, method, path, req.URL.RawQuery, req)
		return resp, err, database.ModeProxy, true
	}

	// Get all responses for this endpoint
	responses, err := s.Repo.FindResponsesByEndpointID(endpoint.ID)
	if err != nil || len(responses) == 0 {
		// Apply delays before returning error
		s.applyDelay(project, endpoint, nil)
		return createErrorResponse(http.StatusInternalServerError, "No responses configured"), nil, database.ModeMock, true
	}

	// Select response based on ResponseMode
	response := selectResponse(responses, endpoint.ResponseMode, req)

	// Apply delays (response-level delay overrides endpoint-level delay, which overrides project-level delay)
	s.applyDelay(project, endpoint, &response)

	// Create and return HTTP response with match indicator
	resp, err := createMockResponse(response)
	return resp, err, database.ModeMock, true
}

// handleProxyMode checks for mock endpoint first, if not found forwards the request to target
func (s *MockService) handleProxyMode(project *database.Project, method, path string, req *http.Request) (*http.Response, error, bool) {
	if project.ActiveProxy == nil {
		return createErrorResponse(http.StatusInternalServerError, "No proxy target configured"), nil, false
	}

	// Check for recursive proxy loops by checking for any header with beo-echo prefix
	for name := range req.Header {
		if strings.HasPrefix(strings.ToLower(name), "beo-echo") {
			return createErrorResponse(http.StatusLoopDetected, "Proxy loop detected: request contains beo-echo header"), nil, false
		}
	}

	// First check if a mock endpoint exists for this request
	endpoint, err := s.Repo.FindMatchingEndpoint(project.ID, method, path)
	if err == nil {
		// Found a matching endpoint, use the mock response
		responses, err := s.Repo.FindResponsesByEndpointID(endpoint.ID)
		if err == nil && len(responses) > 0 {
			// Select response based on ResponseMode
			response := selectResponse(responses, endpoint.ResponseMode, req)

			// Apply delay with proper priority: Response > Endpoint > Project
			s.applyDelay(project, endpoint, &response)

			// Create and return HTTP response from mock
			resp, err := createMockResponse(response)
			if err == nil {
				// Add header to indicate response was mocked
				resp.Header.Set("beo-echo-response-type", "mock")
				return resp, nil, true // True because it was handled by a mock endpoint
			}
		}
	}

	// No matching mock endpoint found or error occurred, forward to target
	// Apply project-level delay before forwarding
	s.applyDelay(project, nil, nil)
	resp, err := executeProxyRequest(project.ActiveProxy.URL, method, path, req.URL.RawQuery, req)
	if err == nil && resp != nil && resp.Header != nil {
		// Add header to indicate response was proxied
		resp.Header.Set("beo-echo-response-type", "proxy")
	}
	return resp, err, false // False because it was forwarded to target, not handled by a mock
}

// handleForwarderMode always forwards requests to the target without checking for mock endpoints
func (s *MockService) handleForwarderMode(project *database.Project, method, path string, req *http.Request) (*http.Response, error) {
	if project.ActiveProxy == nil {
		return createErrorResponse(http.StatusInternalServerError, "No proxy target configured"), nil
	}

	// Check for recursive proxy loops by checking for any header with beo-echo prefix
	for name := range req.Header {
		if strings.HasPrefix(strings.ToLower(name), "beo-echo") {
			return createErrorResponse(http.StatusLoopDetected, "Proxy loop detected: request contains beo-echo header"), nil
		}
	}

	// Use the common executeProxyRequest helper function, but with the path parameter
	// which might differ from req.URL.Path in this context
	// Note: handleForwarderMode always returns false for match status in HandleRequest
	// Apply project-level delay before forwarding
	s.applyDelay(project, nil, nil)

	return executeProxyRequest(project.ActiveProxy.URL, method, path, req.URL.RawQuery, req)
}

// executeProxyRequest is a common helper function to forward requests to a target URL
// with proper header and body copying. This centralizes the forwarding logic for both
// proxy and forwarder modes.
func executeProxyRequest(targetURLString, method, pathStr, queryString string, req *http.Request) (*http.Response, error) {
	// Check for recursive proxy loops by checking for any header with beo-echo prefix
	for name := range req.Header {
		if strings.HasPrefix(strings.ToLower(name), "beo-echo") {
			return createErrorResponse(http.StatusLoopDetected, "Proxy loop detected: request contains beo-echo header"), nil
		}
	}

	targetURL, err := url.Parse(targetURLString)
	if err != nil {
		return createErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Invalid proxy URL: %s", err.Error())), nil
	}

	// Create a new client with desired configuration
	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Disable SSL verification
			},
		},
	}

	// Create new URL for the target
	forwardURL := *targetURL
	// Join the target base path with the requested path
	forwardURL.Path = path.Join(forwardURL.Path, pathStr)
	forwardURL.RawQuery = queryString

	// Read the original request body if present
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return createErrorResponse(http.StatusBadGateway, fmt.Sprintf("Failed to read request body: %s", err.Error())), nil
		}
		// Restore the original body
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	// Create a new request with all original attributes
	newReq, err := http.NewRequestWithContext(
		req.Context(),
		method,
		forwardURL.String(),
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return createErrorResponse(http.StatusBadGateway, fmt.Sprintf("Failed to create request: %s", err.Error())), nil
	}

	// Copy all headers
	for key, values := range req.Header {
		for _, value := range values {
			if key != "Referer" {
				newReq.Header.Add(key, value)
			}
		}
	}

	// Set host header to target host
	newReq.Host = targetURL.Host

	// Add loop detection header to prevent recursive proxying
	newReq.Header.Set("beo-echo-loop-detect", "true")

	// Track request time for latency measurement
	startTime := time.Now()

	// Execute the request
	resp, err := client.Do(newReq)
	if err != nil {
		return createErrorResponse(http.StatusBadGateway, fmt.Sprintf("Request error: %s", err.Error())), nil
	}

	latencyMS := time.Since(startTime).Milliseconds()

	// Log the latency in the header for debugging purposes
	// Using a simpler header name without X- prefix
	if resp != nil && resp.Header != nil {
		resp.Header.Set("beo-echo-latency-ms", fmt.Sprintf("%d", latencyMS))
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

	var headers map[string]string
	if err := json.Unmarshal([]byte(mockResp.Headers), &headers); err != nil {
		fmt.Println("Error unmarshalling headers:", err)
	} else {
		// Set headers
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

// applyDelay applies delay based on priority: Response DelayMS > Endpoint DelayMs > Project DelayMs
// Response parameter is optional - pass nil when response delay is not applicable
func (s *MockService) applyDelay(project *database.Project, endpoint *database.MockEndpoint, response *database.MockResponse) {
	var delayMs int

	// Response delay has highest priority
	if response != nil {
		if response.DelayMS > 0 {
			delayMs = response.DelayMS
		}
	}

	// Endpoint delay overrides project delay
	if endpoint != nil && delayMs == 0 {
		if endpoint.AdvanceConfig != "" {
			if endpointConfig, err := database.ParseEndpointAdvanceConfig(endpoint.AdvanceConfig); err == nil && endpointConfig.DelayMs > 0 {
				delayMs = endpointConfig.DelayMs
			}
		}
	}

	// Get project-level delay from advance config
	if project != nil && delayMs == 0 {
		if project.AdvanceConfig != "" {
			if projectConfig, err := database.ParseProjectAdvanceConfig(project.AdvanceConfig); err == nil {
				delayMs = projectConfig.DelayMs
			}
		}
	}

	// Apply delay if configured
	if delayMs > 0 {
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
}
