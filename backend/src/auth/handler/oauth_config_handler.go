package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"beo-echo/backend/src/database"
)

// OAuthConfigHandler handles generic OAuth configuration management
type OAuthConfigHandler struct {
	db *gorm.DB
}

// NewOAuthConfigHandler creates a new OAuthConfigHandler instance
func NewOAuthConfigHandler(db *gorm.DB) *OAuthConfigHandler {
	return &OAuthConfigHandler{db: db}
}

// ListConfigs handles GET request to list all available OAuth configurations
func (h *OAuthConfigHandler) ListConfigs(c *gin.Context) {
	var configs []database.SSOConfig
	if err := h.db.Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch OAuth configurations",
		})
		return
	}

	var configResponse []gin.H
	for _, config := range configs {
		configResponse = append(configResponse, gin.H{
			"provider":   config.Provider,
			"config":     config.Config,
			"enabled":    config.Enabled,
			"created_at": config.CreatedAt,
			"updated_at": config.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    configResponse,
	})
}
