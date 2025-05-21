package repositories

import (
	"beo-echo/backend/src/database"
	"beo-echo/backend/src/users"
	"context"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// userRepository implements users.UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) users.UserRepository {
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
func (r *userRepository) DeleteUser(ctx context.Context, userID string) error {
	return r.db.Delete(&database.User{}, "id = ?", userID).Error
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
