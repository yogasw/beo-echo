package workspaces

import (
	"beo-echo/backend/src/database"
	systemConfig "beo-echo/backend/src/systemConfigs"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// WorkspaceWithRole extends Workspace with user role information
type WorkspaceWithRole struct {
	ID   string `json:"id"` // Unique identifier for the workspace
	Name string `json:"name"`
	Role string `json:"role"` // Role of the current user in this workspace
}

// WorkspaceMember represents a member of a workspace with user details
type WorkspaceMember struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// WorkspaceRepository defines the data access requirements for workspace operations
type WorkspaceRepository interface {
	GetUserWorkspaces(ctx context.Context, userID string) ([]database.Workspace, error)
	GetUserWorkspacesWithRoles(ctx context.Context, userID string) ([]WorkspaceWithRole, error)
	CreateWorkspace(ctx context.Context, workspace *database.Workspace, userID string) error
	CountUserWorkspaces(ctx context.Context, userID string) (int, error)
	GetUserByID(ctx context.Context, userID string) (*database.User, error)
	CheckWorkspaceRole(ctx context.Context, userID string, workspaceID string) (*database.UserWorkspace, error)
	IsUserWorkspaceAdmin(ctx context.Context, userID string, workspaceID string) (bool, error)
	GetAllWorkspaces(ctx context.Context) ([]database.Workspace, error)
	// New methods for invitation
	GetUserByEmail(ctx context.Context, email string) (*database.User, error)
	AddUserToWorkspace(ctx context.Context, workspaceID string, userID string, role string) error
	GetWorkspaceMembers(ctx context.Context, workspaceID string) ([]WorkspaceMember, error)
	// Methods for auto-invite functionality
	GetWorkspacesWithAutoInvite(ctx context.Context) ([]database.Workspace, error)
	CheckUserWorkspaceMembership(ctx context.Context, userID string, workspaceID string) (*database.UserWorkspace, error)
	CreateUserWorkspaceMembership(ctx context.Context, membership *database.UserWorkspace) error
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

// GetUserWorkspacesWithRoles retrieves all workspaces accessible to a user along with their role in each workspace
func (s *WorkspaceService) GetUserWorkspacesWithRoles(ctx context.Context, userID string) ([]WorkspaceWithRole, error) {
	return s.repo.GetUserWorkspacesWithRoles(ctx, userID)
}

// CreateWorkspace creates a new workspace and adds the user as an admin
func (s *WorkspaceService) CreateWorkspace(ctx context.Context, workspace *database.Workspace, userID string) error {
	// Get effective workspace limit for the user (user-specific or system default)
	maxUserWorkspaces, err := s.GetUserWorkspaceLimit(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get workspace limit: %w", err)
	}

	// Count current user workspaces
	currentWorkspaceCount, err := s.repo.CountUserWorkspaces(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to count user workspaces: %w", err)
	}

	if currentWorkspaceCount >= maxUserWorkspaces {
		return fmt.Errorf("workspace limit exceeded: maximum %d workspaces allowed. Please contact admin for more information", maxUserWorkspaces)
	}

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

func (s *WorkspaceService) GetAllWorkspaces(ctx context.Context) ([]database.Workspace, error) {
	return s.repo.GetAllWorkspaces(ctx)
}

// GetWorkspaceMembers retrieves all members of a workspace with their details
func (s *WorkspaceService) GetWorkspaceMembers(ctx context.Context, workspaceID string) ([]WorkspaceMember, error) {
	return s.repo.GetWorkspaceMembers(ctx, workspaceID)
}

// AddMember adds an existing user to a workspace
// Returns an error if the user doesn't exist
func (s *WorkspaceService) AddMember(ctx context.Context, workspaceID string, email string, role string) (map[string]interface{}, error) {
	// Check if the user already exists
	user, err := s.repo.GetUserByEmail(ctx, email)

	// If user not found, return an error
	if err != nil || user == nil {
		return nil, err
	}

	// User found - check if already in workspace
	existingRole, err := s.repo.CheckWorkspaceRole(ctx, user.ID, workspaceID)
	if err == nil && existingRole != nil {
		// User is already in the workspace
		return map[string]interface{}{
			"user_id": user.ID,
			"email":   email,
			"status":  "already_member",
			"role":    existingRole.Role,
		}, nil
	}

	// Add the user to the workspace
	err = s.repo.AddUserToWorkspace(ctx, workspaceID, user.ID, role)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user_id": user.ID,
		"email":   email,
		"status":  "added",
		"role":    role,
	}, nil
}

// AutoCreateWorkspaceOnRegister creates a default workspace for a new user if enabled
func (s *WorkspaceService) AutoCreateWorkspaceOnRegister(ctx context.Context, userID string, userName string) error {
	// Check if auto-create workspace is enabled
	autoCreateEnabled, err := systemConfig.GetSystemConfigWithType[bool](systemConfig.AUTO_CREATE_WORKSPACE_ON_REGISTER)
	if err != nil {
		return fmt.Errorf("failed to get auto-create workspace configuration: %w", err)
	}

	if !autoCreateEnabled {
		return nil // Auto-create is disabled, nothing to do
	}

	// Create a default workspace for the user
	defaultWorkspace := &database.Workspace{
		Name: fmt.Sprintf("%s's Workspace", userName),
	}

	return s.CreateWorkspace(ctx, defaultWorkspace, userID)
}

// GetUserWorkspaceLimit returns the effective workspace limit for a user
// Uses user-specific limit if set, otherwise falls back to system default
func (s *WorkspaceService) GetUserWorkspaceLimit(ctx context.Context, userID string) (int, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	// Use user-specific limit if set
	if user.MaxWorkspaces != nil {
		return *user.MaxWorkspaces, nil
	}

	// Fall back to system default
	return systemConfig.GetSystemConfigWithType[int](systemConfig.MAX_USER_WORKSPACES)
}

// GetUserProjectLimit returns the effective project limit for a user in a workspace
// Uses user-specific limit if set, otherwise falls back to system default
func (s *WorkspaceService) GetUserProjectLimit(ctx context.Context, userID string) (int, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	// Use user-specific limit if set
	if user.MaxProjectsWorkspace != nil {
		return *user.MaxProjectsWorkspace, nil
	}

	// Fall back to system default
	return systemConfig.GetSystemConfigWithType[int](systemConfig.MAX_WORKSPACE_PROJECTS)
}

// ProcessUserAutoInvite checks if a user should be automatically invited to any workspaces
// based on their email domain, and creates the necessary UserWorkspace records if so
func (s *WorkspaceService) ProcessUserAutoInvite(ctx context.Context, user *database.User) error {

	if user == nil || user.Email == "" {
		return nil // No user or email, nothing to do
	}

	// Extract the domain from the user's email
	parts := strings.Split(user.Email, "@")
	if len(parts) != 2 {
		return nil // Invalid email format, nothing to do
	}
	userDomain := strings.ToLower(parts[1]) // Convert to lowercase for case-insensitive comparison

	// Find all workspaces with auto-invite enabled
	workspaces, err := s.repo.GetWorkspacesWithAutoInvite(ctx)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID).Msg("failed to get workspaces for auto-invite")
		return err
	}

	for _, workspace := range workspaces {
		// Check if user is already a member of this workspace
		existingMembership, err := s.repo.CheckUserWorkspaceMembership(ctx, user.ID, workspace.ID)
		if err == nil && existingMembership != nil {
			// User already has a membership record for this workspace, skip
			continue
		}

		// Check if the user's domain matches any in the workspace's auto-invite domains
		if workspace.AutoInviteDomains == "" {
			continue // No domains configured
		}

		domains := strings.Split(workspace.AutoInviteDomains, ",")
		for _, domain := range domains {
			// Trim whitespace and convert to lowercase for comparison
			domainToCheck := strings.ToLower(strings.TrimSpace(domain))

			if domainToCheck == userDomain {
				// Create a new UserWorkspace record
				newMembership := &database.UserWorkspace{
					UserID:      user.ID,
					WorkspaceID: workspace.ID,
					Role:        workspace.AutoInviteRole,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				if err := s.repo.CreateUserWorkspaceMembership(ctx, newMembership); err != nil {
					log.Error().
						Err(err).
						Str("user_id", user.ID).
						Str("workspace_id", workspace.ID).
						Str("domain", userDomain).
						Msg("failed to create auto-invite membership")
					return err
				}

				log.Info().
					Str("user_id", user.ID).
					Str("email", user.Email).
					Str("workspace_id", workspace.ID).
					Str("workspace_name", workspace.Name).
					Str("role", workspace.AutoInviteRole).
					Msg("user auto-invited to workspace based on email domain")

				break // No need to check other domains for this workspace
			}
		}
	}

	return nil
}
