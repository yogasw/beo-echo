package services

import (
	"context"
)

// ProjectSearchRepository defines data access requirements for project search operations
type ProjectSearchRepository interface {
	CheckAliasAvailability(ctx context.Context, alias string) (bool, error)
	SearchProjectsByNameInUserWorkspaces(ctx context.Context, userID, searchQuery string) ([]ProjectSearchResult, error)
}

// ProjectSearchResult represents a project search result with workspace information
type ProjectSearchResult struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Alias         string `json:"alias"`
	WorkspaceID   string `json:"workspace_id"`
	WorkspaceName string `json:"workspace_name"`
}

// AliasAvailabilityResponse represents the response for alias availability check
type AliasAvailabilityResponse struct {
	Available bool                  `json:"available"`
	Projects  []ProjectSearchResult `json:"projects"`
}

// ProjectSearchService implements project search business operations
type ProjectSearchService struct {
	repo ProjectSearchRepository
}

// NewProjectSearchService creates a new project search service
func NewProjectSearchService(repo ProjectSearchRepository) *ProjectSearchService {
	return &ProjectSearchService{repo: repo}
}

// CheckAliasAndSearchProjects checks if an alias is available and searches for projects with similar names
func (s *ProjectSearchService) CheckAliasAndSearchProjects(ctx context.Context, userID, query string) (*AliasAvailabilityResponse, error) {
	// Check if alias is available globally
	aliasAvailable, err := s.repo.CheckAliasAvailability(ctx, query)
	if err != nil {
		return nil, err
	}

	// Search for projects containing the query in user's workspaces
	projects, err := s.repo.SearchProjectsByNameInUserWorkspaces(ctx, userID, query)
	if err != nil {
		return nil, err
	}

	return &AliasAvailabilityResponse{
		Available: aliasAvailable,
		Projects:  projects,
	}, nil
}
