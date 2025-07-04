package repositories

import (
	"beo-echo/backend/src/database"
	"context"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository implements users.UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

// GetUserByID retrieves a user by ID
func (r *userRepository) GetUserByID(ctx context.Context, id string) (*database.User, error) {
	var user database.User
	err := r.db.Preload("Workspaces").Where("id = ?", id).First(&user).Error
	return &user, err
}

// GetUserByEmail retrieves a user by their email address
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*database.User, error) {
	var user database.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// GetAllUsers retrieves all users in the system
func (r *userRepository) GetAllUsers(ctx context.Context) ([]database.User, error) {
	var users []database.User
	err := r.db.Find(&users).Error
	return users, err
}

// UpdatePassword updates a user's password with a bcrypt hash
func (r *userRepository) UpdatePassword(ctx context.Context, userID string, newPassword string) error {
	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the user's password
	result := r.db.Model(&database.User{}).
		Where("id = ?", userID).
		Update("password", string(hashedPassword))

	return result.Error
}

// UpdateUserFields updates specified user fields
func (r *userRepository) UpdateUserFields(ctx context.Context, userID string, updates map[string]interface{}) error {
	result := r.db.Model(&database.User{}).
		Where("id = ?", userID).
		Updates(updates)

	return result.Error
}

// DeleteUser completely removes a user from the system
// This will also cascade delete related UserIdentity and UserWorkspace records
// due to the constraint:OnDelete:CASCADE in the model definitions
func (r *userRepository) DeleteUser(ctx context.Context, userID string) error {
	tx := r.db.Begin()

	// Check if user exists
	var user database.User
	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// First delete related records manually to ensure they are removed
	if err := tx.Unscoped().Where("user_id = ?", userID).Delete(&database.UserIdentity{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Where("user_id = ?", userID).Delete(&database.UserWorkspace{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Now delete the user
	if err := tx.Unscoped().Delete(&database.User{}, "id = ?", userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// VerifyPassword checks if the provided password matches the hashed password in the database
func (r *userRepository) VerifyPassword(ctx context.Context, userID string, password string) (bool, error) {
	var user database.User
	err := r.db.Select("password").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return false, err
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil, nil
}

// GetWorkspaceUsers retrieves all users in a specific workspace
func (r *userRepository) GetWorkspaceUsers(ctx context.Context, workspaceID string) ([]database.User, error) {
	var users []database.User
	err := r.db.Joins("JOIN user_workspaces ON user_workspaces.user_id = users.id").
		Preload("Workspaces", "workspace_id = ?", workspaceID).
		Where("user_workspaces.workspace_id = ?", workspaceID).
		Find(&users).Error

	return users, err
}

// GetWorkspaceUser retrieves a specific user-workspace relationship
func (r *userRepository) GetWorkspaceUser(ctx context.Context, workspaceID string, userID string) (*database.UserWorkspace, error) {
	var userWorkspace database.UserWorkspace
	err := r.db.Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&userWorkspace).Error
	return &userWorkspace, err
}

// RemoveUserFromWorkspace removes a user from a workspace
func (r *userRepository) RemoveUserFromWorkspace(ctx context.Context, workspaceID string, userID string) error {
	return r.db.Delete(&database.UserWorkspace{}, "workspace_id = ? AND user_id = ?", workspaceID, userID).Error
}

// UpdateUserWorkspaceRole updates a user's role in a workspace
func (r *userRepository) UpdateUserWorkspaceRole(ctx context.Context, workspaceID string, userID string, role string) error {
	return r.db.Model(&database.UserWorkspace{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Update("role", role).Error
}

// UpdateRefreshToken updates the user's refresh token
func (r *userRepository) UpdateRefreshToken(ctx context.Context, userID string, hashedToken string) error {
	return r.db.Model(&database.User{}).
		Where("id = ?", userID).
		Update("refresh_token", hashedToken).Error
}

// GetUserByRefreshToken retrieves a user by their refresh token
// JWT expiry validation is handled by the JWT validation process
func (r *userRepository) GetUserByRefreshToken(ctx context.Context, hashedToken string) (*database.User, error) {
	var user database.User
	err := r.db.Where("refresh_token = ?", hashedToken).First(&user).Error
	return &user, err
}

// ClearRefreshToken removes the refresh token from a user
func (r *userRepository) ClearRefreshToken(ctx context.Context, userID string) error {
	return r.db.Model(&database.User{}).
		Where("id = ?", userID).
		Update("refresh_token", nil).Error
}
