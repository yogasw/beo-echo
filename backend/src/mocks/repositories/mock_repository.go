package repositories

import (
	"encoding/json"
	"fmt"
	"strings"

	"mockoon-control-panel/backend_new/src/database"

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

// FindProjectByName finds a project by its name (slug/subdomain)
func (r *MockRepository) FindProjectByName(name string) (*database.Project, error) {
	var project database.Project
	result := r.DB.Preload("ActiveProxy").Where("name = ?", name).First(&project)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

// FindMatchingEndpoint finds an endpoint that matches the given method and path
func (r *MockRepository) FindMatchingEndpoint(projectID string, method, path string) (*database.MockEndpoint, error) {
	var endpoints []database.MockEndpoint

	result := r.DB.Where("project_id = ? AND method = ? AND enabled = ?", projectID, strings.ToUpper(method), true).Find(&endpoints)
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
	result := r.DB.Preload("Rules").Where("endpoint_id = ? AND active = ?", endpointID, true).Find(&responses)
	if result.Error != nil {
		return nil, result.Error
	}
	return responses, nil
}

// Helper functions

// findBestPathMatch finds the best matching endpoint from a list of endpoints
// considering path parameters (e.g., /users/:id)
func findBestPathMatch(endpoints []database.MockEndpoint, requestPath string) *database.MockEndpoint {
	requestParts := strings.Split(strings.Trim(requestPath, "/"), "/")

	var bestMatch *database.MockEndpoint
	bestScore := -1

	for i := range endpoints {
		endpoint := &endpoints[i]
		endpointParts := strings.Split(strings.Trim(endpoint.Path, "/"), "/")

		if score := calculatePathMatchScore(endpointParts, requestParts); score > bestScore {
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
// Higher score means better match
func calculatePathMatchScore(endpointParts, requestParts []string) int {
	// If lengths don't match, this can't be a match
	if len(endpointParts) != len(requestParts) {
		return -1
	}

	score := 0
	for i := 0; i < len(endpointParts); i++ {
		// Exact match of path part
		if endpointParts[i] == requestParts[i] {
			score += 10
			continue
		}

		// Path parameter (starts with :)
		if strings.HasPrefix(endpointParts[i], ":") {
			score += 5
			continue
		}

		// Not a match
		return -1
	}

	return score
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
