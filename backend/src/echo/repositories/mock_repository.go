package repositories

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"beo-echo/backend/src/database"

	"gorm.io/gorm"
)

// MockRepository handles database operations for mock endpoints
type MockRepository struct {
	DB *gorm.DB
}

// NewMockRepository creates a new repository instance
func NewMockRepository(db *gorm.DB) *MockRepository {
	return &MockRepository{
		DB: db,
	}
}

// FindProjectByAlias finds a project by its alias (slug/subdomain)
func (r *MockRepository) FindProjectByAlias(alias string) (*database.Project, error) {
	var project database.Project
	result := r.DB.Preload("ActiveProxy").Where("alias = ?", alias).First(&project)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

// FindMatchingEndpoint finds an endpoint that matches the given method and path
func (r *MockRepository) FindMatchingEndpoint(projectID string, method, path string) (*database.MockEndpoint, error) {
	var endpoints []database.MockEndpoint

	// Preload ProxyTarget for endpoints that use proxy
	result := r.DB.Preload("ProxyTarget").Where("project_id = ? AND method = ? AND enabled = ?", projectID, strings.ToUpper(method), true).Find(&endpoints)
	if result.Error != nil {
		return nil, result.Error
	}

	// Find best matching path (handle path params like /users/:id)
	bestMatch := findBestPathMatch(endpoints, path)
	if bestMatch == nil {
		return nil, fmt.Errorf("no matching endpoint found")
	}

	return bestMatch, nil
}

// FindResponsesByEndpointID gets all responses for an endpoint
func (r *MockRepository) FindResponsesByEndpointID(endpointID string) ([]database.MockResponse, error) {
	var responses []database.MockResponse
	result := r.DB.Preload("Rules").Where("endpoint_id = ? AND enabled = ?", endpointID, true).Find(&responses)
	if result.Error != nil {
		return nil, result.Error
	}
	return responses, nil
}

// FindProxyTarget gets a proxy target by ID
func (r *MockRepository) GetProxyTarget(proxyTargetID string) (*database.ProxyTarget, error) {
	var proxyTarget database.ProxyTarget
	result := r.DB.Where("id = ?", proxyTargetID).First(&proxyTarget)
	if result.Error != nil {
		return nil, result.Error
	}
	return &proxyTarget, nil
}

// Helper functions

// findBestPathMatch finds the best matching endpoint from a list of endpoints
// Supporting various path patterns:
// 1. Exact match: /users/123
// 2. Path parameters: /users/:id
// 3. Wildcard: /api/v2/customer_rooms/*/broadcast_history
// 4. Regex: /api/v\d+/users/\d+
func findBestPathMatch(endpoints []database.MockEndpoint, requestPath string) *database.MockEndpoint {
	var bestMatch *database.MockEndpoint
	bestScore := -1

	for i := range endpoints {
		endpoint := &endpoints[i]
		
		if score := calculatePathMatchScore(endpoint.Path, requestPath); score > bestScore {
			bestScore = score
			bestMatch = endpoint
		}
	}

	// If we have a match
	if bestScore >= 0 {
		return bestMatch
	}

	return nil
}

// calculatePathMatchScore calculates how well an endpoint path matches a request path
// Higher score means better match. Supports different pattern types:
// - Exact match (highest score: 100)
// - Path parameters with : (score: 80)
// - Wildcard patterns with * (score: 60)
// - Regex patterns (score: 40)
func calculatePathMatchScore(endpointPath, requestPath string) int {
	// Clean paths
	endpointPath = strings.Trim(endpointPath, "/")
	requestPath = strings.Trim(requestPath, "/")
	
	// Exact match gets highest priority
	if endpointPath == requestPath {
		return 100
	}
	
	// Check for regex pattern (contains regex metacharacters)
	if isRegexPattern(endpointPath) {
		if matchesRegex(endpointPath, requestPath) {
			return 40
		}
		return -1
	}
	
	// Check for wildcard or path parameter patterns
	return calculateSegmentMatchScore(endpointPath, requestPath)
}

// isRegexPattern checks if a path contains regex metacharacters
func isRegexPattern(path string) bool {
	// Common regex metacharacters that indicate it's a regex pattern
	regexChars := []string{"\\d", "\\w", "\\s", "[", "]", "(", ")", "+", "?", "^", "$", "\\", "|"}
	for _, char := range regexChars {
		if strings.Contains(path, char) {
			return true
		}
	}
	return false
}

// matchesRegex tests if request path matches the regex pattern
func matchesRegex(pattern, requestPath string) bool {
	// Compile and test regex
	regex, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		return false
	}
	return regex.MatchString(requestPath)
}

// calculateSegmentMatchScore handles path parameters and wildcard matching
func calculateSegmentMatchScore(endpointPath, requestPath string) int {
	endpointParts := strings.Split(endpointPath, "/")
	requestParts := strings.Split(requestPath, "/")
	
	// If lengths don't match, this can't be a match
	if len(endpointParts) != len(requestParts) {
		return -1
	}
	
	score := 0
	wildcardCount := 0
	paramCount := 0
	
	for i := 0; i < len(endpointParts); i++ {
		endpointPart := endpointParts[i]
		requestPart := requestParts[i]
		
		// Exact match of path part (highest score)
		if endpointPart == requestPart {
			score += 10
			continue
		}
		
		// Path parameter (starts with :)
		if strings.HasPrefix(endpointPart, ":") {
			score += 8
			paramCount++
			continue
		}
		
		// Wildcard match (*)
		if endpointPart == "*" {
			score += 6
			wildcardCount++
			continue
		}
		
		// Not a match
		return -1
	}
	
	// Calculate final score based on match type
	if wildcardCount > 0 {
		return 60 + score // Wildcard patterns
	} else if paramCount > 0 {
		return 80 + score // Path parameter patterns
	}
	
	return score // Should not reach here for valid matches
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
