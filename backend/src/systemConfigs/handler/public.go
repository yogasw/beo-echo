package systemConfig

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	systemConfig "beo-echo/backend/src/systemConfigs"
)

// PublicConfigResponse represents the structure of public configuration data
type PublicConfigResponse struct {
	IsAuthenticated bool   `json:"is_authenticated"` // Whether the user is authenticated
	LandingEnabled  bool   `json:"landing_enabled"`  // Whether landing page is enabled
	MockURLFormat   string `json:"mock_url_format"`  // Final URL format: "alias.domain.com" or "domain.com/alias"
}

// GetPublicConfigHandler returns public configuration for landing page
// This endpoint is accessible without authentication to allow landing page to configure itself
func GetPublicConfigHandler(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.Ctx(ctx)

	// Check if user is authenticated (optional for this endpoint)
	isAuthenticated := false
	if authValue, exists := c.Get("isAuthenticated"); exists {
		if authBool, ok := authValue.(bool); ok {
			isAuthenticated = authBool
		}
	}

	// Get required configurations
	landingEnabled, err := systemConfig.GetSystemConfigWithType[bool](systemConfig.LANDING_PAGE_ENABLED)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get landing page enabled config")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve configuration",
		})
		return
	}

	customSubdomainDomain, err := systemConfig.GetSystemConfigWithType[string](systemConfig.CUSTOM_SUBDOMAIN_DOMAIN)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get custom subdomain domain config")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve configuration",
		})
		return
	}

	customSubdomainEnabled, err := systemConfig.GetSystemConfigWithType[bool](systemConfig.CUSTOM_SUBDOMAIN_ENABLED)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get custom subdomain enabled config")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve configuration",
		})
		return
	}

	// Compute the final URL format based on custom subdomain setting
	var mockURLFormat string
	if customSubdomainEnabled {
		// Format: alias.CUSTOM_SUBDOMAIN_DOMAIN
		mockURLFormat = "alias." + customSubdomainDomain
	} else {
		// Format: c.Request.Host/alias
		mockURLFormat = c.Request.Host + "/alias"
	}

	// Parse boolean values - no longer needed since GetSystemConfigWithType returns proper types
	response := PublicConfigResponse{
		IsAuthenticated: isAuthenticated,
		LandingEnabled:  landingEnabled,
		MockURLFormat:   mockURLFormat,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
