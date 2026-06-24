package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"beo-echo/backend/src/auth/pat"
	"beo-echo/backend/src/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MCPOAuthHandler implements a minimal OAuth 2.0 Authorization Server scoped to
// the MCP integration. It supports the pieces an OAuth-aware MCP client (e.g.
// Claude.ai) needs to connect with just the endpoint URL:
//
//   - Authorization Server / Protected Resource discovery metadata
//   - Dynamic Client Registration (RFC 7591) — accepted permissively
//   - Authorization Code grant with PKCE (RFC 7636)
//
// On a successful token exchange the access token returned is a Beo Echo PAT,
// so the rest of the API authenticates it through the existing middleware with
// no special casing.
type MCPOAuthHandler struct {
	db  *gorm.DB
	pat *pat.Service
}

// NewMCPOAuthHandler builds the handler. The server's public URL is derived per
// request from the incoming Host / X-Forwarded-* headers, so no base-URL
// configuration is needed — it always matches the URL the client used.
func NewMCPOAuthHandler(db *gorm.DB) *MCPOAuthHandler {
	return &MCPOAuthHandler{
		db:  db,
		pat: pat.NewService(db),
	}
}

// baseURL reconstructs the externally reachable base URL of this server from the
// request, honoring reverse-proxy headers (X-Forwarded-Proto / X-Forwarded-Host).
func baseURL(c *gin.Context) string {
	scheme := c.GetHeader("X-Forwarded-Proto")
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}
	return scheme + "://" + host
}

const (
	authRequestTTL = 10 * time.Minute
	// oauthTokenTTL is how long an OAuth-issued access token lives. 0 would mean
	// never; we cap it so revoking access naturally happens on expiry too.
	oauthTokenTTL = 90 * 24 * time.Hour
)

// ----- Discovery metadata -----

// AuthorizationServerMetadata serves /.well-known/oauth-authorization-server
// (RFC 8414). MCP clients fetch this to learn the authorize/token endpoints.
func (h *MCPOAuthHandler) AuthorizationServerMetadata(c *gin.Context) {
	base := baseURL(c)
	c.JSON(http.StatusOK, gin.H{
		"issuer":                                base,
		"authorization_endpoint":                base + "/api/oauth/mcp/authorize",
		"token_endpoint":                        base + "/api/oauth/mcp/token",
		"registration_endpoint":                 base + "/api/oauth/mcp/register",
		"scopes_supported":                      []string{"mcp"},
		"response_types_supported":              []string{"code"},
		"grant_types_supported":                 []string{"authorization_code"},
		"code_challenge_methods_supported":      []string{"S256", "plain"},
		"token_endpoint_auth_methods_supported": []string{"none"},
	})
}

// ProtectedResourceMetadata serves /.well-known/oauth-protected-resource
// (MCP authorization spec) pointing clients at the authorization server.
func (h *MCPOAuthHandler) ProtectedResourceMetadata(c *gin.Context) {
	base := baseURL(c)
	c.JSON(http.StatusOK, gin.H{
		"resource":              base + "/mcp",
		"authorization_servers": []string{base},
		"scopes_supported":      []string{"mcp"},
	})
}

// Register implements permissive Dynamic Client Registration (RFC 7591).
// We do not gate on a pre-shared client; any MCP client may register and we
// echo back a generated client_id. Real authorization happens at consent time.
func (h *MCPOAuthHandler) Register(c *gin.Context) {
	var body map[string]interface{}
	_ = c.ShouldBindJSON(&body)

	clientID := "mcp_" + uuid.New().String()

	resp := gin.H{
		"client_id":                  clientID,
		"token_endpoint_auth_method": "none",
		"grant_types":                []string{"authorization_code"},
		"response_types":             []string{"code"},
	}
	// Echo back redirect_uris / client_name if provided so the client is happy.
	if v, ok := body["redirect_uris"]; ok {
		resp["redirect_uris"] = v
	}
	if v, ok := body["client_name"]; ok {
		resp["client_name"] = v
	}
	c.JSON(http.StatusCreated, resp)
}

// ----- Authorization Code grant -----

// Authorize is the front door of the auth-code flow. It validates the request,
// persists a pending OAuthAuthRequest, then redirects the user's browser to the
// SPA consent page (which holds the user's session) to approve.
func (h *MCPOAuthHandler) Authorize(c *gin.Context) {
	q := c.Request.URL.Query()
	respType := q.Get("response_type")
	clientID := q.Get("client_id")
	redirectURI := q.Get("redirect_uri")
	state := q.Get("state")
	scope := q.Get("scope")
	challenge := q.Get("code_challenge")
	challengeMethod := q.Get("code_challenge_method")

	if respType != "code" {
		h.redirectError(c, redirectURI, state, "unsupported_response_type", "only response_type=code is supported")
		return
	}
	if clientID == "" || redirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "client_id and redirect_uri are required"})
		return
	}
	if challengeMethod == "" {
		challengeMethod = "plain"
	}

	req := database.OAuthAuthRequest{
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		State:               state,
		Scope:               scope,
		CodeChallenge:       challenge,
		CodeChallengeMethod: challengeMethod,
		Status:              "pending",
		ExpiresAt:           time.Now().Add(authRequestTTL),
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": err.Error()})
		return
	}

	// Send the browser to the SPA consent page on the same origin the request
	// came in on. The SPA reads its stored JWT, shows what is being authorized,
	// and calls /approve.
	consent := baseURL(c) + "/oauth/consent?request_id=" + url.QueryEscape(req.ID)
	c.Redirect(http.StatusFound, consent)
}

// approveRequest is the body for POST /api/oauth/mcp/approve, called by the
// authenticated SPA on the consent page.
type approveRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Approve   bool   `json:"approve"`
}

// Approve records the user's decision. On approval it mints a single-use
// authorization code and returns the redirect URL the SPA should navigate to.
// Requires JWTAuthMiddleware (the user must be logged in).
func (h *MCPOAuthHandler) Approve(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not authenticated"})
		return
	}

	var body approveRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request: " + err.Error()})
		return
	}

	var req database.OAuthAuthRequest
	if err := h.db.WithContext(c.Request.Context()).Where("id = ?", body.RequestID).First(&req).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Authorization request not found"})
		return
	}
	if req.IsExpired(time.Now()) || req.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Authorization request is no longer valid"})
		return
	}

	if !body.Approve {
		h.db.WithContext(c.Request.Context()).Model(&req).Update("status", "denied")
		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"redirect_url": appendQuery(req.RedirectURI, map[string]string{"error": "access_denied", "state": req.State}),
		})
		return
	}

	uid := userID.(string)
	code := uuid.New().String() + uuid.New().String()
	if err := h.db.WithContext(c.Request.Context()).Model(&req).Updates(map[string]interface{}{
		"status":  "approved",
		"user_id": uid,
		"code":    code,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"redirect_url": appendQuery(req.RedirectURI, map[string]string{"code": code, "state": req.State}),
	})
}

// tokenRequest is the body for POST /api/oauth/mcp/token. We accept both form
// and JSON encodings since clients vary.
type tokenRequest struct {
	GrantType    string `form:"grant_type" json:"grant_type"`
	Code         string `form:"code" json:"code"`
	RedirectURI  string `form:"redirect_uri" json:"redirect_uri"`
	ClientID     string `form:"client_id" json:"client_id"`
	CodeVerifier string `form:"code_verifier" json:"code_verifier"`
}

// Token exchanges an authorization code for an access token. It enforces PKCE
// when a challenge was supplied at /authorize, and returns a Beo Echo PAT as
// the bearer access_token.
func (h *MCPOAuthHandler) Token(c *gin.Context) {
	var body tokenRequest
	// ShouldBind picks form or JSON based on Content-Type.
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": err.Error()})
		return
	}

	if body.GrantType != "authorization_code" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_grant_type"})
		return
	}
	if body.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "code is required"})
		return
	}

	var req database.OAuthAuthRequest
	if err := h.db.WithContext(c.Request.Context()).Where("code = ?", body.Code).First(&req).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant", "error_description": "unknown code"})
		return
	}
	if req.IsExpired(time.Now()) || req.Status != "approved" || req.UserID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant", "error_description": "code is no longer valid"})
		return
	}

	// PKCE verification.
	if req.CodeChallenge != "" {
		if !verifyPKCE(req.CodeChallenge, req.CodeChallengeMethod, body.CodeVerifier) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant", "error_description": "PKCE verification failed"})
			return
		}
	}

	// Mint the access token (a PAT) for the approving user.
	name := "MCP OAuth (" + req.ClientID + ")"
	result, err := h.pat.CreateOAuth(c.Request.Context(), *req.UserID, req.ClientID, name, oauthTokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error", "error_description": err.Error()})
		return
	}

	// Burn the code so it cannot be replayed.
	h.db.WithContext(c.Request.Context()).Model(&req).Updates(map[string]interface{}{
		"status": "used",
		"code":   nil,
	})

	c.JSON(http.StatusOK, gin.H{
		"access_token": result.PlainToken,
		"token_type":   "Bearer",
		"expires_in":   int(oauthTokenTTL.Seconds()),
		"scope":        req.Scope,
	})
}

// ----- helpers -----

func verifyPKCE(challenge, method, verifier string) bool {
	if verifier == "" {
		return false
	}
	switch method {
	case "S256":
		sum := sha256.Sum256([]byte(verifier))
		computed := base64.RawURLEncoding.EncodeToString(sum[:])
		return pat.ConstantTimeEqual(computed, challenge)
	default: // "plain"
		return pat.ConstantTimeEqual(verifier, challenge)
	}
}

func appendQuery(rawURL string, params map[string]string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		// Fall back to naive concatenation; better to send something than nothing.
		sep := "?"
		if strings.Contains(rawURL, "?") {
			sep = "&"
		}
		q := url.Values{}
		for k, v := range params {
			if v != "" {
				q.Set(k, v)
			}
		}
		return rawURL + sep + q.Encode()
	}
	q := u.Query()
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (h *MCPOAuthHandler) redirectError(c *gin.Context, redirectURI, state, code, desc string) {
	if redirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": code, "error_description": desc})
		return
	}
	c.Redirect(http.StatusFound, appendQuery(redirectURI, map[string]string{"error": code, "error_description": desc, "state": state}))
}
