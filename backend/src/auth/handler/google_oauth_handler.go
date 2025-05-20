package handler

import (
	"fmt"
	"net/http"

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
			"message": "Failed to fetch Google OAuth config",
		})
		return
	}

	enabled, err := h.service.GetState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch Google OAuth state",
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
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if err := h.service.SaveConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save Google OAuth config",
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
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if err := h.service.UpdateState(req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update Google OAuth state",
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
	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "redirect_uri is required",
		})
		return
	}

	// First check if OAuth is configured
	config, err := h.service.GetConfig()
	if err != nil || config == nil {
		// Redirect back to login with configuration error
		errorURL := fmt.Sprintf("%s/login?error=google_oauth_not_configured&message=%s",
			redirectURI, "Google OAuth credentials are not configured. Please contact your administrator.")
		c.Redirect(http.StatusTemporaryRedirect, errorURL)
		return
	}

	// Check if OAuth is enabled
	enabled, err := h.service.GetState()
	if err != nil || !enabled {
		// Redirect back to login with disabled error
		errorURL := fmt.Sprintf("%s/login?error=google_oauth_not_configured&message=%s",
			redirectURI, "Google OAuth service is disabled. Please contact your administrator.")
		c.Redirect(http.StatusTemporaryRedirect, errorURL)
		return
	}

	// Get the OAuth URL from service
	loginURL, err := h.service.GetLoginURL(redirectURI)
	if err != nil {
		var errorURL string

		// Handle specific error cases
		switch err.Error() {
		case "google OAuth service is disabled":
			errorURL = fmt.Sprintf("%s/login?error=google_oauth_not_configured&message=%s",
				redirectURI, "Google OAuth service is disabled. Please contact your administrator.")
		case "google OAuth credentials are not configured":
			errorURL = fmt.Sprintf("%s/login?error=google_oauth_not_configured&message=%s",
				redirectURI, "Google OAuth credentials are not configured. Please contact your administrator.")
		case "failed to get OAuth config":
			errorURL = fmt.Sprintf("%s/login?error=google_oauth_not_configured&message=%s",
				redirectURI, "Failed to retrieve Google OAuth configuration. Please contact your administrator.")
		default:
			errorURL = fmt.Sprintf("%s/login?error=google_oauth_error&message=%s",
				redirectURI, err.Error())
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
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No authorization code provided",
		})
		return
	}

	// Get scheme and host from request
	scheme := "http"
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

	_, token, err := h.service.HandleOAuthCallback(code, baseURL)
	if err != nil {
		var errorURL string
		switch err.Error() {
		case "auto-registration is disabled and user does not exist":
			errorURL = fmt.Sprintf("/login?error=registration_disabled&message=%s",
				"Auto-registration is disabled. Please contact your administrator.")
		case "email domain not allowed":
			errorURL = fmt.Sprintf("/login?error=domain_not_allowed&message=%s",
				"Your email domain is not allowed. Please contact your administrator.")
		case "google OAuth service is disabled":
			errorURL = fmt.Sprintf("/login?error=google_oauth_not_configured&message=%s",
				"Google OAuth service is disabled. Please contact your administrator.")
		default:
			errorURL = fmt.Sprintf("/login?error=google_oauth_error&message=%s", err.Error())
		}
		c.Redirect(http.StatusTemporaryRedirect, errorURL)
		return
	}

	// Success! Set token in cookie and redirect to root path
	c.SetCookie(
		"jwt_token",
		token,
		86400, // 24 hours
		"/",
		"",   // domain
		true, // secure
		true, // httpOnly
	)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
