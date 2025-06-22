package services

import (
	"context"
	"crypto/rand"
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
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid email format: %s", email)
	}

	domain := parts[1]

	for _, allowedDomain := range config.AllowDomains {
		if strings.EqualFold(allowedDomain, domain) { // Case-insensitive comparison
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

	// 5. Process auto-invite based on email domain
	// We need the AutoInviteService, but to avoid circular dependencies,
	// We'll handle the auto-invite directly here
	if err := s.processAutoInvite(user); err != nil {
		// Log but don't fail the auth flow
		fmt.Printf("Warning: Failed to process auto-invite for user %s: %v\n", user.ID, err)
	}

	// 6. Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return user, token, nil
}

// GetLoginURL generates the Google OAuth login URL
func (s *GoogleOAuthService) GetLoginURL(backendCallbackURI string, frontendRedirectURI string) (string, error) {
	config, err := s.GetConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get OAuth config: %w", err)
	}

	if config == nil || config.ClientID == "" || config.ClientSecret == "" {
		return "", NewGoogleOAuthNotConfiguredError()
	}

	enabled, err := s.GetState()
	if err != nil {
		return "", fmt.Errorf("failed to check OAuth state: %w", err)
	}

	if !enabled {
		return "", NewGoogleOAuthDisabledError()
	}

	// Generate a random state to prevent CSRF and include frontend redirect URL
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	stateStr := fmt.Sprintf("%x&%s", stateBytes, frontendRedirectURI)

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  backendCallbackURI,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return oauth2Config.AuthCodeURL(stateStr), nil
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
		RedirectURL:  fmt.Sprintf("%s/api/oauth/google/callback", baseURL),
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
		return nil, NewOAuthUserInfoRetrievalError()
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
		// Extract the domain part from the email
		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			return fmt.Errorf("invalid email format")
		}
		domain := parts[1]
		return NewDomainNotAllowedError(domain)
	}

	return nil
}

// processAutoInvite checks if a user should be automatically invited to any workspaces
// based on their email domain, and creates the necessary UserWorkspace records if needed
func (s *GoogleOAuthService) processAutoInvite(user *database.User) error {
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
	var workspaces []database.Workspace
	if err := s.db.Where("auto_invite_enabled = ?", true).Find(&workspaces).Error; err != nil {
		return fmt.Errorf("failed to get workspaces for auto-invite: %w", err)
	}

	for _, workspace := range workspaces {
		// Check if user is already a member of this workspace
		var existingMembership database.UserWorkspace
		existingMembershipCount := s.db.Where("user_id = ? AND workspace_id = ?", user.ID, workspace.ID).
			First(&existingMembership).RowsAffected

		if existingMembershipCount > 0 {
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
				newMembership := database.UserWorkspace{
					UserID:      user.ID,
					WorkspaceID: workspace.ID,
					Role:        workspace.AutoInviteRole,
				}

				if err := s.db.Create(&newMembership).Error; err != nil {
					return fmt.Errorf("failed to create auto-invite membership for workspace %s: %w", workspace.ID, err)
				}

				break // No need to check other domains for this workspace
			}
		}
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
		return nil, NewAutoRegistrationDisabledError()
	}

	// Create new user and identity
	newUser := &database.User{
		Email:    userInfo.Email,
		Name:     userInfo.Name,
		IsActive: true,
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

	// Check if auto-create workspace is enabled and create workspace for new users
	autoCreateWorkspace, err := systemConfig.GetSystemConfigWithType[bool](systemConfig.AUTO_CREATE_WORKSPACE_ON_REGISTER)
	if err == nil && autoCreateWorkspace {
		// Create a default workspace for the new user
		if err := s.createDefaultWorkspaceForUser(newUser); err != nil {
			// Log but don't fail the auth flow
			fmt.Printf("Warning: Failed to create default workspace for user %s: %v\n", newUser.ID, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return newUser, nil
}

// createDefaultWorkspaceForUser creates a default workspace for a new user
func (s *GoogleOAuthService) createDefaultWorkspaceForUser(user *database.User) error {
	workspace := &database.Workspace{
		Name: fmt.Sprintf("%s's Workspace", user.Name),
	}

	// Create workspace in a transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Create the workspace
		if err := tx.Create(workspace).Error; err != nil {
			return fmt.Errorf("failed to create workspace: %w", err)
		}

		// Add the user as an admin
		userWorkspace := database.UserWorkspace{
			UserID:      user.ID,
			WorkspaceID: workspace.ID,
			Role:        "admin",
		}

		if err := tx.Create(&userWorkspace).Error; err != nil {
			return fmt.Errorf("failed to add user to workspace: %w", err)
		}

		return nil
	})
}
