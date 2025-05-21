package workspaces

import (
	"beo-echo/backend/src/database"
	"context"
)

// WorkspaceRepository defines the data access requirements for workspace operations
type WorkspaceRepository interface {
	GetUserWorkspaces(ctx context.Context, userID string) ([]database.Workspace, error)
	CreateWorkspace(ctx context.Context, workspace *database.Workspace, userID string) error
	CheckWorkspaceRole(ctx context.Context, userID string, workspaceID string) (*database.UserWorkspace, error)
	IsUserWorkspaceAdmin(ctx context.Context, userID string, workspaceID string) (bool, error)
}

// WorkspaceService implements the workspace business operations
type WorkspaceService struct {
	repo WorkspaceRepository
}

// NewWorkspaceService creates a new workspace service
func NewWorkspaceService(repo WorkspaceRepository) *WorkspaceService {
	return &WorkspaceService{repo: repo}
}

// GetUserWorkspaces retrieves all workspaces accessible to a user
func (s *WorkspaceService) GetUserWorkspaces(ctx context.Context, userID string) ([]database.Workspace, error) {
	return s.repo.GetUserWorkspaces(ctx, userID)
}

// CreateWorkspace creates a new workspace and adds the user as an admin
func (s *WorkspaceService) CreateWorkspace(ctx context.Context, workspace *database.Workspace, userID string) error {
	return s.repo.CreateWorkspace(ctx, workspace, userID)
}

// CheckWorkspaceRole returns a user's role in a specific workspace
func (s *WorkspaceService) CheckWorkspaceRole(ctx context.Context, userID string, workspaceID string) (*database.UserWorkspace, error) {
	return s.repo.CheckWorkspaceRole(ctx, userID, workspaceID)
}

// IsUserWorkspaceAdmin checks if a user is an admin in a specific workspace
func (s *WorkspaceService) IsUserWorkspaceAdmin(ctx context.Context, userID string, workspaceID string) (bool, error) {
	return s.repo.IsUserWorkspaceAdmin(ctx, userID, workspaceID)
}
