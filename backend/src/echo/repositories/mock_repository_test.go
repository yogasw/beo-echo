package repositories

import (
	"testing"

	"beo-echo/backend/src/database"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePathMatchScore(t *testing.T) {
	tests := []struct {
		name         string
		endpointPath string
		requestPath  string
		expectedScore int
		description  string
	}{
		// Exact matches
		{
			name:          "exact match simple",
			endpointPath:  "/users",
			requestPath:   "/users",
			expectedScore: 100,
			description:   "Exact match should get highest score",
		},
		{
			name:          "exact match complex",
			endpointPath:  "/api/v2/customer_rooms/broadcast_history",
			requestPath:   "/api/v2/customer_rooms/broadcast_history",
			expectedScore: 100,
			description:   "Complex exact match should get highest score",
		},
		
		// Wildcard patterns
		{
			name:          "wildcard single segment",
			endpointPath:  "/api/v2/customer_rooms/*/broadcast_history",
			requestPath:   "/api/v2/customer_rooms/331931307/broadcast_history",
			expectedScore: 106, // 60 (wildcard base) + 40 (4 exact matches) + 6 (1 wildcard match)
			description:   "Wildcard should match numeric values",
		},
		{
			name:          "wildcard string value",
			endpointPath:  "/api/v2/customer_rooms/*/broadcast_history",
			requestPath:   "/api/v2/customer_rooms/salam/broadcast_history",
			expectedScore: 106, // 60 (wildcard base) + 40 (4 exact matches) + 6 (1 wildcard match)
			description:   "Wildcard should match string values",
		},
		{
			name:          "multiple wildcards",
			endpointPath:  "/api/*/users/*/profile",
			requestPath:   "/api/v2/users/123/profile",
			expectedScore: 102, // 60 (wildcard base) + 30 (3 exact matches) + 12 (2 wildcard matches)
			description:   "Multiple wildcards should work",
		},
		
		// Path parameters
		{
			name:          "path parameter",
			endpointPath:  "/users/:id",
			requestPath:   "/users/123",
			expectedScore: 98, // 80 (param base) + 10 (1 exact match) + 8 (1 param match)
			description:   "Path parameters should work",
		},
		{
			name:          "mixed exact and parameters",
			endpointPath:  "/api/v2/users/:id/settings",
			requestPath:   "/api/v2/users/123/settings",
			expectedScore: 128, // 80 (param base) + 40 (4 exact matches) + 8 (1 param match)
			description:   "Mixed exact and parameters should work",
		},
		
		// Regex patterns
		{
			name:          "regex numeric pattern",
			endpointPath:  "/api/v\\d+/users/\\d+",
			requestPath:   "/api/v2/users/123",
			expectedScore: 40,
			description:   "Regex patterns should work for numeric matches",
		},
		{
			name:          "regex complex pattern",
			endpointPath:  "/api/v[12]/customer_rooms/\\w+/broadcast_history",
			requestPath:   "/api/v2/customer_rooms/salam/broadcast_history",
			expectedScore: 40,
			description:   "Complex regex patterns should work",
		},
		
		// No matches
		{
			name:          "no match different paths",
			endpointPath:  "/users",
			requestPath:   "/posts",
			expectedScore: -1,
			description:   "Different paths should not match",
		},
		{
			name:          "no match different lengths",
			endpointPath:  "/users/123",
			requestPath:   "/users",
			expectedScore: -1,
			description:   "Different path lengths should not match",
		},
		{
			name:          "regex no match",
			endpointPath:  "/api/v\\d+/users",
			requestPath:   "/api/vX/users",
			expectedScore: -1,
			description:   "Regex should not match invalid patterns",
		},
		
		// Edge cases
		{
			name:          "empty paths",
			endpointPath:  "",
			requestPath:   "",
			expectedScore: 100,
			description:   "Empty paths should match exactly",
		},
		{
			name:          "root path",
			endpointPath:  "/",
			requestPath:   "/",
			expectedScore: 100,
			description:   "Root paths should match exactly",
		},
		{
			name:          "trailing slashes ignored",
			endpointPath:  "/users/",
			requestPath:   "/users",
			expectedScore: 100,
			description:   "Trailing slashes should be ignored",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calculatePathMatchScore(tt.endpointPath, tt.requestPath)
			assert.Equal(t, tt.expectedScore, score, tt.description)
		})
	}
}

func TestFindBestPathMatch(t *testing.T) {
	endpoints := []database.MockEndpoint{
		{
			ID:     "1",
			Method: "GET",
			Path:   "/users",
		},
		{
			ID:     "2", 
			Method: "GET",
			Path:   "/users/:id",
		},
		{
			ID:     "3",
			Method: "GET", 
			Path:   "/api/v2/customer_rooms/*/broadcast_history",
		},
		{
			ID:     "4",
			Method: "GET",
			Path:   "/api/v\\d+/users/\\d+",
		},
		{
			ID:     "5",
			Method: "GET",
			Path:   "/api/*/users/*",
		},
		{
			ID:     "6",
			Method: "GET",
			Path:   "/api/v\\d+/specific/\\d+",
		},
	}

	tests := []struct {
		name            string
		requestPath     string
		expectedEndpointID string
		description     string
	}{
		{
			name:               "exact match wins over parameters",
			requestPath:        "/users",
			expectedEndpointID: "1", // Exact match should win
			description:        "Exact match should have higher priority than parameters",
		},
		{
			name:               "parameter match when no exact",
			requestPath:        "/users/123",
			expectedEndpointID: "2", // Parameter match
			description:        "Parameter match should work when no exact match",
		},
		{
			name:               "wildcard customer rooms",
			requestPath:        "/api/v2/customer_rooms/331931307/broadcast_history",
			expectedEndpointID: "3", // Wildcard match
			description:        "Wildcard should match the customer rooms pattern",
		},
		{
			name:               "wildcard customer rooms string",
			requestPath:        "/api/v2/customer_rooms/salam/broadcast_history",
			expectedEndpointID: "3", // Wildcard match
			description:        "Wildcard should match string values too",
		},
		{
			name:               "regex pattern match",
			requestPath:        "/api/v3/users/456",
			expectedEndpointID: "5", // Wildcard match wins over regex (higher score: 102 vs 40)
			description:        "Wildcard should have higher priority than regex for this path",
		},
		{
			name:               "wildcard vs regex priority",
			requestPath:        "/api/v2/users/123",
			expectedEndpointID: "5", // Wildcard should win over regex (higher score)
			description:        "Wildcard should have higher priority than regex",
		},
		{
			name:               "regex only match",
			requestPath:        "/api/v2/specific/456",
			expectedEndpointID: "6", // Only regex pattern matches this path
			description:        "Regex pattern should match when no wildcard available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := findBestPathMatch(endpoints, tt.requestPath)
			if tt.expectedEndpointID == "" {
				assert.Nil(t, match, tt.description)
			} else {
				assert.NotNil(t, match, tt.description)
				assert.Equal(t, tt.expectedEndpointID, match.ID, tt.description)
			}
		})
	}
}

func TestIsRegexPattern(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "simple path not regex",
			path:     "/users/123",
			expected: false,
		},
		{
			name:     "wildcard not regex",
			path:     "/users/*/profile",
			expected: false,
		},
		{
			name:     "parameter not regex",
			path:     "/users/:id",
			expected: false,
		},
		{
			name:     "digit regex pattern",
			path:     "/api/v\\d+/users",
			expected: true,
		},
		{
			name:     "word regex pattern",
			path:     "/users/\\w+/profile",
			expected: true,
		},
		{
			name:     "character class regex",
			path:     "/api/v[12]/users",
			expected: true,
		},
		{
			name:     "group regex pattern",
			path:     "/api/(v1|v2)/users",
			expected: true,
		},
		{
			name:     "anchor regex pattern",
			path:     "^/users$",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRegexPattern(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMatchesRegex(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		requestPath string
		expected    bool
	}{
		{
			name:        "digit pattern matches",
			pattern:     "/api/v\\d+/users/\\d+",
			requestPath: "/api/v2/users/123",
			expected:    true,
		},
		{
			name:        "digit pattern no match",
			pattern:     "/api/v\\d+/users/\\d+",
			requestPath: "/api/vX/users/123",
			expected:    false,
		},
		{
			name:        "word pattern matches",
			pattern:     "/users/\\w+/profile",
			requestPath: "/users/john123/profile",
			expected:    true,
		},
		{
			name:        "character class matches",
			pattern:     "/api/v[123]/users",
			requestPath: "/api/v2/users",
			expected:    true,
		},
		{
			name:        "character class no match",
			pattern:     "/api/v[123]/users",
			requestPath: "/api/v4/users",
			expected:    false,
		},
		{
			name:        "invalid regex returns false",
			pattern:     "/api/v[/users", // Invalid regex
			requestPath: "/api/v2/users",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesRegex(tt.pattern, tt.requestPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}
