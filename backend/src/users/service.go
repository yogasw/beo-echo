package users

import (
	"beo-echo/backend/src/database"
	systemConfig "beo-echo/backend/src/systemConfigs"
	"context"
	"errors"
)

// UserRepository defines data access requirements for user operations
type UserRepository interface {
	// User Management
	GetUserByID(ctx context.Context, id string) (*database.User, error)
	GetUserByEmail(ctx context.Context, email string) (*database.User, error)
	GetAllUsers(ctx context.Context) ([]database.User, error)
	UpdatePassword(ctx context.Context, userID string, newPassword string) error
	UpdateUserFields(ctx context.Context, userID string, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, userID string) error
	VerifyPassword(ctx context.Context, userID string, password string) (bool, error)

	// Workspace-User Relationship
	GetWorkspaceUsers(ctx context.Context, workspaceID string) ([]database.User, error)
	GetWorkspaceUser(ctx context.Context, workspaceID string, userID string) (*database.UserWorkspace, error)
	RemoveUserFromWorkspace(ctx context.Context, workspaceID string, userID string) error
	UpdateUserWorkspaceRole(ctx context.Context, workspaceID string, userID string, role string) error
}

// UserService implements user business operations
type UserService struct {
	repo UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetCurrentUser retrieves the authenticated user's information
func (s *UserService) GetCurrentUser(ctx context.Context, userID string) (*database.User, map[string]bool, error) {
	// Fetch user details
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	// Get feature flags from system config
	featureFlags, err := systemConfig.GetFeatureFlags()
	if err != nil {
		// Log the error but don't fail the request
		featureFlags = make(map[string]bool)
	}

	return user, featureFlags, nil
}

// GetAllUsers retrieves all users in the system (admin/owner only)
func (s *UserService) GetAllUsers(ctx context.Context) ([]database.User, error) {
	return s.repo.GetAllUsers(ctx)
}

// UpdatePassword updates the user's password
func (s *UserService) UpdatePassword(ctx context.Context, userID string, currentPassword, newPassword string) error {
	// Verify current password
	isValid, err := s.repo.VerifyPassword(ctx, userID, currentPassword)
	if err != nil {
		return err
	}

	if !isValid {
		return errors.New("current password is incorrect")
	}

	// Update the password
	return s.repo.UpdatePassword(ctx, userID, newPassword)
}

// UpdateUserProfile updates user profile information
func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, name, email string, isOwner bool) error {
	// Get current user
	currentUser, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Prepare update fields
	updates := make(map[string]interface{})

	// Always allow name update
	if name != "" {
		updates["name"] = name
	}

	// Email updates are restricted based on system config
	// Check if email updates are enabled for non-owner users
	if email != "" && email != currentUser.Email {
		emailUpdatesEnabled := false
		if currentUser.IsOwner || isOwner {
			emailUpdatesEnabled = true // Owners can always update email
		} else {
			var err error
			emailUpdatesEnabled, err = systemConfig.GetSystemConfigWithType[bool](string(systemConfig.FEATURE_EMAIL_UPDATES_ENABLED))
			if err != nil {
				return err
			}
		}

		if !emailUpdatesEnabled {
			return errors.New("email updates are disabled by system administrator")
		}

		// Check if email is already in use by another user
		existingUser, _ := s.repo.GetUserByEmail(ctx, email)
		if existingUser != nil && existingUser.ID != userID {
			return errors.New("email address is already in use")
		}

		updates["email"] = email
	}

	// If there's nothing to update
	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	// Update the user
	return s.repo.UpdateUserFields(ctx, userID, updates)
}

// UpdateUserFields directly updates the specified fields for a user
func (s *UserService) UpdateUserFields(ctx context.Context, userID string, updates map[string]interface{}) error {
	return s.repo.UpdateUserFields(ctx, userID, updates)
}

// GetWorkspaceUsers retrieves all users in a specific workspace
func (s *UserService) GetWorkspaceUsers(ctx context.Context, workspaceID string) ([]database.User, error) {
	return s.repo.GetWorkspaceUsers(ctx, workspaceID)
}

// RemoveUserFromWorkspace removes a user from a workspace
func (s *UserService) RemoveUserFromWorkspace(ctx context.Context, workspaceID string, userID string) error {
	return s.repo.RemoveUserFromWorkspace(ctx, workspaceID, userID)
}

// DeleteUser completely removes a user from the system
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	return s.repo.DeleteUser(ctx, userID)
}

// UpdateUserWorkspaceRole updates a user's role in a workspace
func (s *UserService) UpdateUserWorkspaceRole(ctx context.Context, workspaceID string, userID string, role string) error {
	// Check how many admins are in the workspace
	if role == "member" {
		// Get the current user workspace
		userWorkspace, err := s.repo.GetWorkspaceUser(ctx, workspaceID, userID)
		if err != nil {
			return err
		}

		// If user is already a member, no change needed
		if userWorkspace.Role == "member" {
			return nil
		}

		// If current role is admin, ensure there will still be at least one admin
		if userWorkspace.Role == "admin" {
			// We'd need to implement a check in the repository for this
			// This would be better implemented in the repository, but for demonstration:
			users, err := s.repo.GetWorkspaceUsers(ctx, workspaceID)
			if err != nil {
				return err
			}

			// Count admins
			adminCount := 0
			for _, user := range users {
				for _, ws := range user.Workspaces {
					if ws.WorkspaceID == workspaceID && ws.Role == "admin" {
						adminCount++
					}
				}
			}

			if adminCount <= 1 {
				return errors.New("cannot demote the last admin in a workspace")
			}
		}
	}

	return s.repo.UpdateUserWorkspaceRole(ctx, workspaceID, userID, role)
}
