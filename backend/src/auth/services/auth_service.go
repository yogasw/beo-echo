package services

import (
	"beo-echo/backend/src/database"
	"context"

	"gorm.io/gorm"
)

// UserRepository defines the data access interface needed for authentication
type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*database.User, error)
	UpdateRefreshToken(ctx context.Context, userID string, hashedToken string) error
	GetUserByRefreshToken(ctx context.Context, hashedToken string) (*database.User, error)
	ClearRefreshToken(ctx context.Context, userID string) error
}

// AuthService provides authentication related services
type AuthService struct {
	repo UserRepository
	db   *gorm.DB
}

// NewAuthService creates a new auth service
func NewAuthService(repo UserRepository, db *gorm.DB) *AuthService {
	return &AuthService{
		repo: repo,
		db:   db,
	}
}

// GetUserByEmail retrieves a user by their email address
func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*database.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

// Direct DB access for backwards compatibility - to be removed in future versions
// as we refactor more code to use the repository pattern
func (s *AuthService) GetUserByEmailDirect(email string) (*database.User, error) {
	var user database.User
	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserRepository returns the user repository interface
func (s *AuthService) GetUserRepository() UserRepository {
	return s.repo
}
