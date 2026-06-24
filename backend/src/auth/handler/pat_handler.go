package handler

import (
	"net/http"
	"time"

	"beo-echo/backend/src/auth/pat"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PATHandler exposes personal-access-token management on the profile.
type PATHandler struct {
	service *pat.Service
}

// NewPATHandler builds a PAT handler backed by the given DB.
func NewPATHandler(db *gorm.DB) *PATHandler {
	return &PATHandler{service: pat.NewService(db)}
}

// createTokenRequest is the body for POST /api/users/me/tokens.
type createTokenRequest struct {
	Name       string `json:"name" binding:"required"` // Human label
	ExpiresDays int   `json:"expires_days"`            // 0 = never expires
}

// CreateToken mints a new PAT for the authenticated user and returns the
// plaintext exactly once.
func (h *PATHandler) CreateToken(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not authenticated"})
		return
	}

	var req createTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	var ttl time.Duration
	if req.ExpiresDays > 0 {
		ttl = time.Duration(req.ExpiresDays) * 24 * time.Hour
	}

	result, err := h.service.Create(c.Request.Context(), userID.(string), req.Name, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create token: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Token created. Copy it now — it will not be shown again.",
		"data": gin.H{
			"token":      result.PlainToken, // shown once
			"id":         result.Token.ID,
			"name":       result.Token.Name,
			"prefix":     result.Token.Prefix,
			"expires_at": result.Token.ExpiresAt,
			"created_at": result.Token.CreatedAt,
		},
	})
}

// ListTokens returns metadata for all the user's tokens (never the secret).
func (h *PATHandler) ListTokens(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not authenticated"})
		return
	}

	tokens, err := h.service.List(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to list tokens: " + err.Error()})
		return
	}

	list := make([]gin.H, 0, len(tokens))
	for _, t := range tokens {
		list = append(list, gin.H{
			"id":           t.ID,
			"name":         t.Name,
			"prefix":       t.Prefix,
			"source":       t.Source,
			"client_id":    t.ClientID,
			"last_used_at": t.LastUsedAt,
			"expires_at":   t.ExpiresAt,
			"created_at":   t.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}

// RevokeToken deletes one of the user's tokens.
func (h *PATHandler) RevokeToken(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not authenticated"})
		return
	}

	tokenID := c.Param("tokenId")
	if err := h.service.Revoke(c.Request.Context(), userID.(string), tokenID); err != nil {
		if err == pat.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Token not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to revoke token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Token revoked successfully"})
}
