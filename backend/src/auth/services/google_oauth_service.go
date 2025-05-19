package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"beo-echo/backend/src/database"

	"gorm.io/gorm"
)

// GoogleOAuthConfig represents the configuration for Google OAuth
type GoogleOAuthConfig struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	AllowDomains []string `json:"allow_domains"` // List of allowed email domains
	Instructions string   `json:"instructions"`  // Setup instructions/notes
}

// GoogleOAuthService handles business logic for Google OAuth operations
type GoogleOAuthService struct {
	db *gorm.DB
}

// NewGoogleOAuthService creates a new GoogleOAuthService instance
func NewGoogleOAuthService(db *gorm.DB) *GoogleOAuthService {
	return &GoogleOAuthService{db: db}
}

// SaveGoogleConfig saves Google OAuth configuration
func (s *GoogleOAuthService) SaveConfig(config GoogleOAuthConfig) error {
	// Convert config to JSON string
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Update or create SSO config
	result := s.db.Model(&database.SSOConfig{}).
		Where("provider = ?", "google").
		Updates(map[string]interface{}{
			"config": string(configJSON),
		})

	if result.RowsAffected == 0 {
		// Create new config if it doesn't exist
		ssoConfig := database.SSOConfig{
			Provider: "google",
			Config:   string(configJSON),
			Enabled:  true,
		}
		if err := s.db.Create(&ssoConfig).Error; err != nil {
			return fmt.Errorf("failed to create SSO config: %w", err)
		}
	}

	return result.Error
}

// GetGoogleConfig retrieves Google OAuth configuration
func (s *GoogleOAuthService) GetConfig() (*GoogleOAuthConfig, error) {
	var ssoConfig database.SSOConfig
	if err := s.db.Where("provider = ?", "google").First(&ssoConfig).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch SSO config: %w", err)
	}

	var config GoogleOAuthConfig
	if err := json.Unmarshal([]byte(ssoConfig.Config), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// UpdateGoogleState enables/disables Google OAuth
func (s *GoogleOAuthService) UpdateState(enabled bool) error {
	result := s.db.Model(&database.SSOConfig{}).
		Where("provider = ?", "google").
		Update("enabled", enabled)

	if result.Error != nil {
		return fmt.Errorf("failed to update SSO config state: %w", result.Error)
	}

	return nil
}

// GetState checks if Google OAuth is enabled
func (s *GoogleOAuthService) GetState() (bool, error) {
	var ssoConfig database.SSOConfig
	if err := s.db.Where("provider = ?", "google").First(&ssoConfig).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to fetch SSO config: %w", err)
	}

	return ssoConfig.Enabled, nil
}

// ValidateDomain checks if the email domain is allowed
func (s *GoogleOAuthService) ValidateDomain(email string) (bool, error) {
	config, err := s.GetConfig()
	if err != nil {
		return false, err
	}

	if config == nil || len(config.AllowDomains) == 0 {
		return true, nil // If no domains configured, allow all
	}

	// Extract domain from email
	domain := email[len(email)-strings.Index(email, "@"):]

	for _, allowedDomain := range config.AllowDomains {
		if allowedDomain == domain {
			return true, nil
		}
	}

	return false, nil
}
