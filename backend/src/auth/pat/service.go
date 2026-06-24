package pat

import (
	"context"
	"errors"
	"time"

	"beo-echo/backend/src/database"

	"gorm.io/gorm"
)

var (
	// ErrInvalidToken is returned when a presented token does not resolve to an
	// active, unexpired record.
	ErrInvalidToken = errors.New("invalid api token")
	// ErrNotFound is returned when a token record does not exist for the user.
	ErrNotFound = errors.New("api token not found")
)

// Service manages personal access tokens.
type Service struct {
	db *gorm.DB
}

// NewService creates a PAT service backed by the given DB handle.
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateResult carries the one-time plaintext token alongside its stored record.
type CreateResult struct {
	PlainToken string                 // Shown to the user exactly once
	Token      database.UserApiToken  // Persisted metadata (no secret)
}

// Create mints a new PAT for the user. ttl of 0 means the token never expires.
func (s *Service) Create(ctx context.Context, userID, name string, ttl time.Duration) (*CreateResult, error) {
	plain, hash, err := GenerateToken()
	if err != nil {
		return nil, err
	}

	token := database.UserApiToken{
		UserID:    userID,
		Name:      name,
		TokenHash: hash,
		Prefix:    DisplayPrefix(plain),
		Source:    "pat",
	}
	if ttl > 0 {
		exp := time.Now().Add(ttl)
		token.ExpiresAt = &exp
	}

	if err := s.db.WithContext(ctx).Create(&token).Error; err != nil {
		return nil, err
	}

	return &CreateResult{PlainToken: plain, Token: token}, nil
}

// CreateOAuth mints a token issued through the OAuth flow, tagged with the
// originating client id so it can be listed/revoked separately from manual PATs.
func (s *Service) CreateOAuth(ctx context.Context, userID, clientID, name string, ttl time.Duration) (*CreateResult, error) {
	plain, hash, err := GenerateToken()
	if err != nil {
		return nil, err
	}

	token := database.UserApiToken{
		UserID:    userID,
		Name:      name,
		TokenHash: hash,
		Prefix:    DisplayPrefix(plain),
		Source:    "oauth",
		ClientID:  clientID,
	}
	if ttl > 0 {
		exp := time.Now().Add(ttl)
		token.ExpiresAt = &exp
	}

	if err := s.db.WithContext(ctx).Create(&token).Error; err != nil {
		return nil, err
	}

	return &CreateResult{PlainToken: plain, Token: token}, nil
}

// List returns all tokens belonging to a user, newest first.
func (s *Service) List(ctx context.Context, userID string) ([]database.UserApiToken, error) {
	var tokens []database.UserApiToken
	err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tokens).Error
	return tokens, err
}

// Revoke deletes a token by id, scoped to the owning user.
func (s *Service) Revoke(ctx context.Context, userID, tokenID string) error {
	res := s.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", tokenID, userID).
		Delete(&database.UserApiToken{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Authenticate resolves a plaintext token to its owning user. It rejects
// expired tokens and updates last_used_at on success (best effort).
func (s *Service) Authenticate(ctx context.Context, plaintext string) (*database.User, error) {
	if !IsPAT(plaintext) {
		return nil, ErrInvalidToken
	}

	hash := HashToken(plaintext)

	var token database.UserApiToken
	err := s.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		First(&token).Error
	if err != nil {
		return nil, ErrInvalidToken
	}

	now := time.Now()
	if token.IsExpired(now) {
		return nil, ErrInvalidToken
	}

	var user database.User
	if err := s.db.WithContext(ctx).Where("id = ?", token.UserID).First(&user).Error; err != nil {
		return nil, ErrInvalidToken
	}
	if !user.IsActive {
		return nil, ErrInvalidToken
	}

	// Best-effort touch; never block auth on this.
	s.db.WithContext(ctx).Model(&database.UserApiToken{}).
		Where("id = ?", token.ID).
		Update("last_used_at", now)

	return &user, nil
}
