package services

import (
	"net/http"
	"net/url"
	"testing"

	"beo-echo/backend/src/database"

	"github.com/stretchr/testify/assert"
)

func TestMatchesRules_ANDLogic(t *testing.T) {
	tests := []struct {
		name     string
		response database.MockResponse
		req      *http.Request
		expected bool
	}{
		{
			name: "no rules should always match",
			response: database.MockResponse{
				RulesLogic: "and",
				Rules:      []database.MockRule{},
			},
			req:      createTestRequest("GET", "/test", map[string]string{}, map[string]string{}),
			expected: true,
		},
		{
			name: "single header rule matches",
			response: database.MockResponse{
				RulesLogic: "and",
				Rules: []database.MockRule{
					{Type: "header", Key: "Authorization", Operator: "equals", Value: "Bearer token123"},
				},
			},
			req:      createTestRequest("GET", "/test", map[string]string{"Authorization": "Bearer token123"}, map[string]string{}),
			expected: true,
		},
		{
			name: "single header rule does not match",
			response: database.MockResponse{
				RulesLogic: "and",
				Rules: []database.MockRule{
					{Type: "header", Key: "Authorization", Operator: "equals", Value: "Bearer token123"},
				},
			},
			req:      createTestRequest("GET", "/test", map[string]string{"Authorization": "Bearer wrong"}, map[string]string{}),
			expected: false,
		},
		{
			name: "multiple rules all match (AND logic)",
			response: database.MockResponse{
				RulesLogic: "and",
				Rules: []database.MockRule{
					{Type: "header", Key: "Authorization", Operator: "equals", Value: "Bearer token123"},
					{Type: "query", Key: "user_id", Operator: "equals", Value: "123"},
				},
			},
			req:      createTestRequest("GET", "/test?user_id=123", map[string]string{"Authorization": "Bearer token123"}, map[string]string{"user_id": "123"}),
			expected: true,
		},
		{
			name: "multiple rules one fails (AND logic)",
			response: database.MockResponse{
				RulesLogic: "and",
				Rules: []database.MockRule{
					{Type: "header", Key: "Authorization", Operator: "equals", Value: "Bearer token123"},
					{Type: "query", Key: "user_id", Operator: "equals", Value: "123"},
				},
			},
			req:      createTestRequest("GET", "/test?user_id=456", map[string]string{"Authorization": "Bearer token123"}, map[string]string{"user_id": "456"}),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesRules(tt.response, tt.req)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMatchesRules_ORLogic(t *testing.T) {
	tests := []struct {
		name     string
		response database.MockResponse
		req      *http.Request
		expected bool
	}{
		{
			name: "no rules should always match (OR logic)",
			response: database.MockResponse{
				RulesLogic: "or",
				Rules:      []database.MockRule{},
			},
			req:      createTestRequest("GET", "/test", map[string]string{}, map[string]string{}),
			expected: true,
		},
		{
			name: "single rule matches (OR logic)",
			response: database.MockResponse{
				RulesLogic: "or",
				Rules: []database.MockRule{
					{Type: "header", Key: "Authorization", Operator: "equals", Value: "Bearer token123"},
				},
			},
			req:      createTestRequest("GET", "/test", map[string]string{"Authorization": "Bearer token123"}, map[string]string{}),
			expected: true,
		},
		{
			name: "multiple rules one matches (OR logic)",
			response: database.MockResponse{
				RulesLogic: "or",
				Rules: []database.MockRule{
					{Type: "header", Key: "Authorization", Operator: "equals", Value: "Bearer token123"},
					{Type: "query", Key: "user_id", Operator: "equals", Value: "123"},
				},
			},
			req:      createTestRequest("GET", "/test?user_id=456", map[string]string{"Authorization": "Bearer token123"}, map[string]string{"user_id": "456"}),
			expected: true,
		},
		{
			name: "multiple rules all fail (OR logic)",
			response: database.MockResponse{
				RulesLogic: "or",
				Rules: []database.MockRule{
					{Type: "header", Key: "Authorization", Operator: "equals", Value: "Bearer token123"},
					{Type: "query", Key: "user_id", Operator: "equals", Value: "123"},
				},
			},
			req:      createTestRequest("GET", "/test?user_id=456", map[string]string{"Authorization": "Bearer wrong"}, map[string]string{"user_id": "456"}),
			expected: false,
		},
		{
			name: "multiple rules all match (OR logic)",
			response: database.MockResponse{
				RulesLogic: "or",
				Rules: []database.MockRule{
					{Type: "header", Key: "Authorization", Operator: "equals", Value: "Bearer token123"},
					{Type: "query", Key: "user_id", Operator: "equals", Value: "123"},
				},
			},
			req:      createTestRequest("GET", "/test?user_id=123", map[string]string{"Authorization": "Bearer token123"}, map[string]string{"user_id": "123"}),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesRules(tt.response, tt.req)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMatchesRules_DefaultLogic(t *testing.T) {
	tests := []struct {
		name     string
		response database.MockResponse
		req      *http.Request
		expected bool
	}{
		{
			name: "empty rule_logic defaults to OR - one rule matches",
			response: database.MockResponse{
				RulesLogic: "",
				Rules: []database.MockRule{
					{Type: "header", Key: "Authorization", Operator: "equals", Value: "Bearer token123"},
					{Type: "query", Key: "user_id", Operator: "equals", Value: "123"},
				},
			},
			req:      createTestRequest("GET", "/test?user_id=456", map[string]string{"Authorization": "Bearer token123"}, map[string]string{"user_id": "456"}),
			expected: true, // Should pass because OR logic is default and header matches
		},
		{
			name: "invalid rule_logic defaults to OR - one rule matches",
			response: database.MockResponse{
				RulesLogic: "invalid",
				Rules: []database.MockRule{
					{Type: "header", Key: "Authorization", Operator: "equals", Value: "Bearer token123"},
					{Type: "query", Key: "user_id", Operator: "equals", Value: "123"},
				},
			},
			req:      createTestRequest("GET", "/test?user_id=456", map[string]string{"Authorization": "Bearer token123"}, map[string]string{"user_id": "456"}),
			expected: true, // Should pass because OR logic is default and header matches
		},
		{
			name: "empty rule_logic defaults to OR - no rules match",
			response: database.MockResponse{
				RulesLogic: "",
				Rules: []database.MockRule{
					{Type: "header", Key: "Authorization", Operator: "equals", Value: "Bearer token123"},
					{Type: "query", Key: "user_id", Operator: "equals", Value: "123"},
				},
			},
			req:      createTestRequest("GET", "/test?user_id=456", map[string]string{"Authorization": "Bearer wrong"}, map[string]string{"user_id": "456"}),
			expected: false, // Should fail because no rules match in OR logic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesRules(tt.response, tt.req)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to create test HTTP requests
func createTestRequest(method, urlStr string, headers map[string]string, queryParams map[string]string) *http.Request {
	parsedURL, _ := url.Parse(urlStr)

	// Add query parameters
	q := parsedURL.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	parsedURL.RawQuery = q.Encode()

	req := &http.Request{
		Method: method,
		URL:    parsedURL,
		Header: make(http.Header),
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req
}
