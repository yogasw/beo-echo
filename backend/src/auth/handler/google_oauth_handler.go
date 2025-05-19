package handler

import (
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Google OAuth config"})
		return
	}

	enabled, err := h.service.GetState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Google OAuth state"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"config":  config,
		"enabled": enabled,
	})
}

// UpdateConfig handles PUT request for Google OAuth config
func (h *GoogleOAuthHandler) UpdateConfig(c *gin.Context) {
	var config services.GoogleOAuthConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.SaveConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Google OAuth config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Google OAuth config updated successfully"})
}

// UpdateState handles PUT request to enable/disable Google OAuth
func (h *GoogleOAuthHandler) UpdateState(c *gin.Context) {
	var req struct {
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.UpdateState(req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Google OAuth state"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Google OAuth state updated successfully"})
}

// HandleGoogleCallback handles OAuth callback from Google
func (h *GoogleOAuthHandler) HandleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No authorization code provided"})
		return
	}

	// TODO: Implement OAuth flow completion
	// 1. Exchange code for tokens
	// 2. Get user info from Google
	// 3. Validate email domain
	// 4. Create/update user and identity
	// 5. Generate JWT token
	// 6. Return token or redirect with token

	c.JSON(http.StatusOK, gin.H{"message": "OAuth callback received"})
}
