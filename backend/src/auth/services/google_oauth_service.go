package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	auth "beo-echo/backend/src/auth"
	"beo-echo/backend/src/database"
	systemConfig "beo-echo/backend/src/systemConfigs"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

// GoogleUserInfo represents user information from Google
type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

// HandleOAuthCallback processes the OAuth callback flow
func (s *GoogleOAuthService) HandleOAuthCallback(code string, baseURL string) (*database.User, string, error) {
	// 1. Exchange code for tokens
	tokens, err := s.exchangeCodeForTokens(code, baseURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange code for tokens: %w", err)
	}

	// 2. Get user info
	userInfo, err := s.fetchGoogleUserInfo(tokens.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user info: %w", err)
	}

	// 3. Validate domain
	if err := s.validateUserDomain(userInfo.Email); err != nil {
		return nil, "", err
	}

	// 4. Create/update user and identity with auto-register check
	user, err := s.handleUserCreation(userInfo, tokens.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to handle user creation: %w", err)
	}

	// 5. Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return user, token, nil
}

// Internal helper functions

func (s *GoogleOAuthService) exchangeCodeForTokens(code string, baseURL string) (*oauth2.Token, error) {
	config, err := s.GetConfig()
	if err != nil {
		return nil, err
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/mock/api/auth/google/callback", baseURL),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return token, nil
}

func (s *GoogleOAuthService) fetchGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %d", resp.StatusCode)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (s *GoogleOAuthService) validateUserDomain(email string) error {
	isValid, err := s.ValidateDomain(email)
	if err != nil {
		return fmt.Errorf("failed to validate domain: %w", err)
	}

	if !isValid {
		return fmt.Errorf("email domain not allowed")
	}

	return nil
}

func (s *GoogleOAuthService) handleUserCreation(userInfo *GoogleUserInfo, accessToken string) (*database.User, error) {
	// Check if user exists by identity
	var identity database.UserIdentity
	err := s.db.Where("provider = ? AND provider_id = ?", "google", userInfo.Sub).
		Preload("User").First(&identity).Error

	if err == nil {
		// Update existing identity
		identity.AccessToken = accessToken
		identity.Email = userInfo.Email
		identity.Name = userInfo.Name
		identity.AvatarURL = userInfo.Picture
		if err := s.db.Save(&identity).Error; err != nil {
			return nil, err
		}
		return &identity.User, nil
	}

	autoRegisterEnabled, err := systemConfig.GetSystemConfigWithType[bool](systemConfig.FEATURE_OAUTH_AUTO_REGISTER)
	if err != nil {
		return nil, fmt.Errorf("failed to get auto-register config: %w", err)
	}

	if !autoRegisterEnabled {
		return nil, fmt.Errorf("auto-registration is disabled and user does not exist")
	}

	// Create new user and identity
	newUser := &database.User{
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		IsEnabled: true,
	}

	// Start transaction
	tx := s.db.Begin()

	if err := tx.Create(newUser).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	newIdentity := &database.UserIdentity{
		UserID:      newUser.ID,
		Provider:    "google",
		ProviderID:  userInfo.Sub,
		Email:       userInfo.Email,
		Name:        userInfo.Name,
		AvatarURL:   userInfo.Picture,
		AccessToken: accessToken,
	}

	if err := tx.Create(newIdentity).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return newUser, nil
}
