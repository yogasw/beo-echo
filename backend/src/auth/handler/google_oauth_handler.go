package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/auth/services"
)

type GoogleOAuthHandler struct {
	service *services.GoogleOAuthService
}

func NewGoogleOAuthHandler(service *services.GoogleOAuthService) *GoogleOAuthHandler {
	return &GoogleOAuthHandler{service: service}
}

// GetConfig handles GET request for Google OAuth config
func (h *GoogleOAuthHandler) GetConfig(c *gin.Context) {
	config, err := h.service.GetConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to fetch Google OAuth config",
			"error":   services.ErrOAuthConfigRetrieval,
		})
		return
	}

	enabled, err := h.service.GetState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to fetch Google OAuth state",
			"error":   services.ErrOAuthConfigRetrieval,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"config":  config,
			"enabled": enabled,
		},
	})
}

// UpdateConfig handles PUT request for Google OAuth config
func (h *GoogleOAuthHandler) UpdateConfig(c *gin.Context) {
	var config services.GoogleOAuthConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		oauthErr := services.NewInvalidRequestBodyError()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": oauthErr.Error(),
			"error":   services.ErrInvalidRequestBody,
		})
		return
	}

	if err := h.service.SaveConfig(config); err != nil {
		oauthErr := services.NewOAuthConfigRetrievalError("failed to save configuration")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": oauthErr.Error(),
			"error":   services.ErrOAuthConfigRetrieval,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Google OAuth config updated successfully",
	})
}

// UpdateState handles PUT request to enable/disable Google OAuth
func (h *GoogleOAuthHandler) UpdateState(c *gin.Context) {
	var req struct {
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		oauthErr := services.NewInvalidRequestBodyError()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": oauthErr.Error(),
			"error":   services.ErrInvalidRequestBody,
		})
		return
	}

	if err := h.service.UpdateState(req.Enabled); err != nil {
		oauthErr := services.NewOAuthConfigRetrievalError("failed to update state")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": oauthErr.Error(),
			"error":   services.ErrOAuthConfigRetrieval,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Google OAuth state updated successfully",
	})
}

// InitiateLogin starts the OAuth flow by redirecting to Google
func (h *GoogleOAuthHandler) InitiateLogin(c *gin.Context) {
	// Get frontend redirect URL from query
	frontendRedirectURI := c.Query("redirect_uri")
	if frontendRedirectURI == "" {
		oauthErr := services.NewMissingRedirectURIError()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": oauthErr.Error(),
			"error":   services.ErrMissingRedirectURI,
		})
		return
	}

	// Get scheme and host for backend callback URL
	scheme := "http"
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Scheme") == "https" || strings.Contains(c.Request.Header.Get("Referer"), "https") {
		scheme = "https"
	}
	backendBaseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	backendCallbackURI := fmt.Sprintf("%s/api/oauth/google/callback", backendBaseURL)

	// Check if OAuth is configured
	config, err := h.service.GetConfig()
	if err != nil || config == nil {
		// Create a proper error object and cast it to OAuthError to access GetType()
		oauthErr := services.NewGoogleOAuthNotConfiguredError()
		typedErr, ok := oauthErr.(*services.OAuthError)
		errorType := services.ErrGoogleOAuthNotConfigured
		if ok {
			errorType = typedErr.GetType()
		}

		// Use our error constants for consistent error handling
		errorURL := fmt.Sprintf("%s?error=%s&message=%s",
			frontendRedirectURI,
			errorType,
			url.QueryEscape(oauthErr.Error()))
		c.Redirect(http.StatusTemporaryRedirect, errorURL)
		return
	}

	// Check if OAuth is enabled
	enabled, err := h.service.GetState()
	if err != nil || !enabled {
		// Create a proper error object and cast it to OAuthError to access GetType()
		oauthErr := services.NewGoogleOAuthDisabledError()
		typedErr, ok := oauthErr.(*services.OAuthError)
		errorType := services.ErrGoogleOAuthDisabled
		if ok {
			errorType = typedErr.GetType()
		}

		// Use our error constants for consistent error handling
		errorURL := fmt.Sprintf("%s?error=%s&message=%s",
			frontendRedirectURI,
			errorType,
			url.QueryEscape(oauthErr.Error()))
		c.Redirect(http.StatusTemporaryRedirect, errorURL)
		return
	}

	// Get the OAuth URL from service with backend callback
	loginURL, err := h.service.GetLoginURL(backendCallbackURI, frontendRedirectURI)
	if err != nil {
		var errorURL string
		var oauthError *services.OAuthError

		if errors.As(err, &oauthError) {
			switch oauthError.GetType() {
			case services.ErrGoogleOAuthDisabled:
				errorURL = fmt.Sprintf("%s?error=google_oauth_disabled&message=%s",
					frontendRedirectURI, url.QueryEscape(oauthError.Error()))
			case services.ErrGoogleOAuthNotConfigured:
				errorURL = fmt.Sprintf("%s?error=google_oauth_not_configured&message=%s",
					frontendRedirectURI, url.QueryEscape(oauthError.Error()))
			case services.ErrOAuthConfigRetrieval:
				errorURL = fmt.Sprintf("%s?error=google_oauth_config_error&message=%s",
					frontendRedirectURI, url.QueryEscape(oauthError.Error()))
			default:
				errorURL = fmt.Sprintf("%s?error=google_oauth_error&message=%s",
					frontendRedirectURI, url.QueryEscape(err.Error()))
			}
		} else {
			// Fallback for non-typed errors
			errorURL = fmt.Sprintf("%s?error=google_oauth_error&message=%s",
				frontendRedirectURI, url.QueryEscape(err.Error()))
		}

		c.Redirect(http.StatusTemporaryRedirect, errorURL)
		return
	}

	// Redirect to Google's OAuth page
	c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

// HandleCallback handles OAuth callback from Google
func (h *GoogleOAuthHandler) HandleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	if code == "" {
		oauthErr := services.NewMissingCodeError()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": oauthErr.Error(),
			"error":   services.ErrMissingCode,
		})
		return
	}

	// Extract frontend redirect URL from state
	stateParts := strings.Split(state, "&")
	if len(stateParts) != 2 {
		oauthErr := services.NewInvalidStateError()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": oauthErr.Error(),
			"error":   services.ErrInvalidState,
		})
		return
	}
	frontendRedirectURI := stateParts[1]

	// Get scheme and host from request
	scheme := "http"
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Scheme") == "https" || strings.Contains(c.Request.Header.Get("Referer"), "https") {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

	user, token, err := h.service.HandleOAuthCallback(code, baseURL)
	if err != nil {
		var errorURL string

		// Try to cast to OAuthError to get more specific information
		var oauthError *services.OAuthError
		if errors.As(err, &oauthError) {
			switch oauthError.GetType() {
			case services.ErrAutoRegistrationDisabled:
				errorURL = fmt.Sprintf("%s?error=registration_disabled&message=%s",
					frontendRedirectURI, url.QueryEscape(oauthError.Error()))
			case services.ErrDomainNotAllowed:
				errorURL = fmt.Sprintf("%s?error=domain_not_allowed&message=%s",
					frontendRedirectURI, url.QueryEscape(oauthError.Error()))
			case services.ErrGoogleOAuthDisabled:
				errorURL = fmt.Sprintf("%s?error=google_oauth_disabled&message=%s",
					frontendRedirectURI, url.QueryEscape(oauthError.Error()))
			default:
				errorURL = fmt.Sprintf("%s?error=google_oauth_error&message=%s",
					frontendRedirectURI, url.QueryEscape(err.Error()))
			}
		} else {
			// Fallback for non-typed errors
			errorURL = fmt.Sprintf("%s?error=google_oauth_error&message=%s",
				frontendRedirectURI, url.QueryEscape(err.Error()))
		}

		c.Redirect(http.StatusTemporaryRedirect, errorURL)
		return
	}

	// Send token in URL for frontend to handle with SSO flag
	successURL := fmt.Sprintf("%s?success=true&token=%s&user=%s&sso=google",
		frontendRedirectURI,
		url.QueryEscape(token),
		url.QueryEscape(user.Email))
	c.Redirect(http.StatusTemporaryRedirect, successURL)
}
