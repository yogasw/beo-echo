package workspaces

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"beo-echo/backend/src/database"
)

// AutoInviteHandler provides endpoints to manage workspace auto-invitation settings
type AutoInviteHandler struct {
	db *gorm.DB
}

// NewAutoInviteHandler creates a new AutoInviteHandler instance
func NewAutoInviteHandler(db *gorm.DB) *AutoInviteHandler {
	return &AutoInviteHandler{db: db}
}

// GetAutoInviteConfig retrieves the auto-invite configuration for a workspace
func (h *AutoInviteHandler) GetAutoInviteConfig(c *gin.Context) {
	workspaceID := c.Param("workspaceID")

	// Verify that user is an owner (this endpoint should only be accessible by system owners)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "User not authenticated",
		})
		return
	}

	var user database.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve user information",
		})
		return
	}

	if !user.IsOwner {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "Only system owners can manage auto-invite configurations",
		})
		return
	}

	// Retrieve workspace with auto-invite settings
	var workspace database.Workspace
	if err := h.db.Where("id = ?", workspaceID).First(&workspace).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Workspace not found",
		})
		return
	}

	// Return the auto-invite configuration
	domainList := []string{}
	if workspace.AutoInviteDomains != "" {
		domains := strings.Split(workspace.AutoInviteDomains, ",")
		for _, domain := range domains {
			domainList = append(domainList, strings.TrimSpace(domain))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"enabled":        workspace.AutoInviteEnabled,
			"domains":        domainList,
			"role":           workspace.AutoInviteRole,
			"workspace_id":   workspace.ID,
			"workspace_name": workspace.Name,
		},
	})
}

// UpdateAutoInviteConfig updates the auto-invite configuration for a workspace
func (h *AutoInviteHandler) UpdateAutoInviteConfig(c *gin.Context) {
	workspaceID := c.Param("workspaceID")

	// Verify that user is an owner
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "User not authenticated",
		})
		return
	}

	var user database.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve user information",
		})
		return
	}

	if !user.IsOwner {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "Only system owners can manage auto-invite configurations",
		})
		return
	}

	// Parse request body
	var request struct {
		Enabled bool     `json:"enabled"`
		Domains []string `json:"domains"`
		Role    string   `json:"role"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate role
	if request.Role != "admin" && request.Role != "member" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid role: must be 'admin' or 'member'",
		})
		return
	}

	// Validate domains format
	for _, domain := range request.Domains {
		// Basic validation for domain format
		trimmedDomain := strings.TrimSpace(domain)
		if trimmedDomain == "" {
			continue
		}

		// Check for valid domain format
		if !strings.Contains(trimmedDomain, ".") || strings.Contains(trimmedDomain, "@") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid domain format: " + trimmedDomain,
			})
			return
		}
	}

	// Convert domains to comma-separated string, removing any empty entries
	var cleanedDomains []string
	for _, domain := range request.Domains {
		domain = strings.TrimSpace(domain)
		if domain != "" {
			cleanedDomains = append(cleanedDomains, domain)
		}
	}
	domainsStr := strings.Join(cleanedDomains, ",")

	// Retrieve workspace
	var workspace database.Workspace
	if err := h.db.Where("id = ?", workspaceID).First(&workspace).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Workspace not found",
		})
		return
	}

	// Update workspace with new settings
	workspace.AutoInviteEnabled = request.Enabled
	workspace.AutoInviteDomains = domainsStr
	workspace.AutoInviteRole = request.Role

	if err := h.db.Save(&workspace).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to update auto-invite configuration: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Auto-invite configuration updated successfully",
		"data": gin.H{
			"enabled":        workspace.AutoInviteEnabled,
			"domains":        cleanedDomains,
			"role":           workspace.AutoInviteRole,
			"workspace_id":   workspace.ID,
			"workspace_name": workspace.Name,
		},
	})
}
