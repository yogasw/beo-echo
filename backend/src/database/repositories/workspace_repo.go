package repositories

import (
	"beo-echo/backend/src/database"
	"beo-echo/backend/src/workspaces"
	"context"

	"gorm.io/gorm"
)

// workspaceRepository implements the WorkspaceRepository interface
type workspaceRepository struct {
	db *gorm.DB
}

// GetAllWorkspaces implements workspaces.WorkspaceRepository.
func (r *workspaceRepository) GetAllWorkspaces(ctx context.Context) ([]database.Workspace, error) {
	var workspaces []database.Workspace
	err := r.db.Find(&workspaces).Error
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

// NewWorkspaceRepository creates a new workspace repository
func NewWorkspaceRepository(db *gorm.DB) workspaces.WorkspaceRepository {
	return &workspaceRepository{db: db}
}

// GetUserWorkspaces retrieves all workspaces accessible to a user
func (r *workspaceRepository) GetUserWorkspaces(ctx context.Context, userID string) ([]database.Workspace, error) {
	var workspaces []database.Workspace
	err := r.db.Joins("JOIN user_workspaces ON user_workspaces.workspace_id = workspaces.id").
		Where("user_workspaces.user_id = ?", userID).
		Find(&workspaces).Error

	return workspaces, err
}

// GetUserWorkspacesWithRoles retrieves all workspaces accessible to a user along with their role in each workspace
func (r *workspaceRepository) GetUserWorkspacesWithRoles(ctx context.Context, userID string) ([]workspaces.WorkspaceWithRole, error) {
	var result []struct {
		database.Workspace
		Role string
	}

	err := r.db.Table("workspaces").
		Select("workspaces.*, user_workspaces.role").
		Joins("JOIN user_workspaces ON user_workspaces.workspace_id = workspaces.id").
		Where("user_workspaces.user_id = ?", userID).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// Convert to WorkspaceWithRole
	workspacesWithRoles := make([]workspaces.WorkspaceWithRole, len(result))
	for i, item := range result {
		workspacesWithRoles[i] = workspaces.WorkspaceWithRole{
			ID:   item.ID,
			Name: item.Name,
			Role: item.Role,
		}
	}

	return workspacesWithRoles, nil
}

// CreateWorkspace creates a new workspace and adds the user as an admin
func (r *workspaceRepository) CreateWorkspace(ctx context.Context, workspace *database.Workspace, userID string) error {
	// Create workspace in a transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create the workspace
		if err := tx.Create(workspace).Error; err != nil {
			return err
		}

		// Add the current user as an admin of this workspace
		userWorkspace := database.UserWorkspace{
			UserID:      userID,
			WorkspaceID: workspace.ID,
			Role:        "admin",
		}

		if err := tx.Create(&userWorkspace).Error; err != nil {
			return err
		}

		return nil
	})
}

// CheckWorkspaceRole returns a user's role in a specific workspace
func (r *workspaceRepository) CheckWorkspaceRole(ctx context.Context, userID string, workspaceID string) (*database.UserWorkspace, error) {
	var userWorkspace database.UserWorkspace
	err := r.db.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&userWorkspace).Error
	if err != nil {
		return nil, err
	}
	return &userWorkspace, nil
}

// IsUserWorkspaceAdmin checks if a user is an admin in a specific workspace
func (r *workspaceRepository) IsUserWorkspaceAdmin(ctx context.Context, userID string, workspaceID string) (bool, error) {
	var userWorkspace database.UserWorkspace
	result := r.db.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&userWorkspace)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}

	return userWorkspace.Role == "admin", nil
}

// GetUserByEmail retrieves a user by their email address
func (r *workspaceRepository) GetUserByEmail(ctx context.Context, email string) (*database.User, error) {
	var user database.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// AddUserToWorkspace adds a user to a workspace with the specified role
func (r *workspaceRepository) AddUserToWorkspace(ctx context.Context, workspaceID string, userID string, role string) error {
	userWorkspace := database.UserWorkspace{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Role:        role,
	}

	// First check if the user is already in the workspace
	var existing database.UserWorkspace
	result := r.db.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&existing)

	if result.Error == nil {
		// User already exists in workspace, update their role
		return r.db.Model(&existing).Update("role", role).Error
	} else if result.Error == gorm.ErrRecordNotFound {
		// User is not in the workspace, add them
		return r.db.Create(&userWorkspace).Error
	}

	// Some other error occurred
	return result.Error
}

// GetWorkspaceMembers retrieves all members of a workspace with their user details
func (r *workspaceRepository) GetWorkspaceMembers(ctx context.Context, workspaceID string) ([]workspaces.WorkspaceMember, error) {
	// Join users and user_workspaces tables to get user details with their roles
	var results []struct {
		UserID string
		Name   string
		Email  string
		Role   string
	}

	err := r.db.Table("users").
		Select("users.id as user_id, users.name, users.email, user_workspaces.role").
		Joins("JOIN user_workspaces ON user_workspaces.user_id = users.id").
		Where("user_workspaces.workspace_id = ?", workspaceID).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert to WorkspaceMember slice
	members := make([]workspaces.WorkspaceMember, len(results))
	for i, result := range results {
		members[i] = workspaces.WorkspaceMember{
			ID:    result.UserID,
			Name:  result.Name,
			Email: result.Email,
			Role:  result.Role,
		}
	}

	return members, nil
}

// CountUserWorkspaces counts the number of workspaces a user belongs to
func (r *workspaceRepository) CountUserWorkspaces(ctx context.Context, userID string) (int, error) {
	var count int64
	err := r.db.Table("user_workspaces").
		Where("user_id = ?", userID).
		Count(&count).Error

	return int(count), err
}

// GetUserByID retrieves a user by their ID
func (r *workspaceRepository) GetUserByID(ctx context.Context, userID string) (*database.User, error) {
	var user database.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetWorkspacesWithAutoInvite retrieves all workspaces that have auto-invite enabled
func (r *workspaceRepository) GetWorkspacesWithAutoInvite(ctx context.Context) ([]database.Workspace, error) {
	var workspaces []database.Workspace
	err := r.db.Where("auto_invite_enabled = ?", true).Find(&workspaces).Error
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

// CheckUserWorkspaceMembership checks if a user is already a member of a workspace
func (r *workspaceRepository) CheckUserWorkspaceMembership(ctx context.Context, userID string, workspaceID string) (*database.UserWorkspace, error) {
	var membership database.UserWorkspace
	err := r.db.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&membership).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User is not a member
		}
		return nil, err
	}
	return &membership, nil
}

// CreateUserWorkspaceMembership creates a new user workspace membership record
func (r *workspaceRepository) CreateUserWorkspaceMembership(ctx context.Context, membership *database.UserWorkspace) error {
	return r.db.Create(membership).Error
}

// CreateProject creates a new project
func (r *workspaceRepository) CreateProject(ctx context.Context, project *database.Project) error {
	return r.db.Create(project).Error
}

// CreateEndpoint creates a new mock endpoint
func (r *workspaceRepository) CreateEndpoint(ctx context.Context, endpoint *database.MockEndpoint) error {
	return r.db.Create(endpoint).Error
}

// CreateResponse creates a new mock response
func (r *workspaceRepository) CreateResponse(ctx context.Context, response *database.MockResponse) error {
	return r.db.Create(response).Error
}

// CheckProjectAliasExists checks if a project alias already exists
func (r *workspaceRepository) CheckProjectAliasExists(ctx context.Context, alias string) (bool, error) {
	var count int64
	err := r.db.Model(&database.Project{}).Where("alias = ?", alias).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
