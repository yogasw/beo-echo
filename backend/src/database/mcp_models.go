package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserApiToken is a long-lived personal access token (PAT) used by external
// clients such as the MCP server / CLI to authenticate as a user without going
// through the short-lived JWT login flow.
//
// It also backs OAuth-issued access tokens: when a user approves an MCP OAuth
// client in the browser, the access token handed back is one of these records
// (Source == "oauth"). Static tokens minted from the profile page are
// Source == "pat".
//
// Only the SHA-256 hash of the token is stored; the plaintext is shown exactly
// once at creation time. The plaintext format is "beo_pat_<random>".
type UserApiToken struct {
	ID         string     `gorm:"type:string;primaryKey" json:"id"`
	UserID     string     `gorm:"index;not null" json:"user_id"` // Owner of the token
	Name       string     `json:"name"`                          // Human label, e.g. "MCP on laptop"
	TokenHash  string     `gorm:"uniqueIndex;not null" json:"-"` // SHA-256 hex of plaintext token
	Prefix     string     `gorm:"index" json:"prefix"`           // First chars of token for display ("beo_pat_ab12")
	Source     string     `gorm:"default:'pat'" json:"source"`   // "pat" (manual) or "oauth" (issued via OAuth flow)
	ClientID   string     `json:"client_id"`                     // OAuth client id when Source == "oauth"
	LastUsedAt *time.Time `json:"last_used_at"`                  // Updated on each successful auth
	ExpiresAt  *time.Time `json:"expires_at"`                    // Optional expiry; nil = never expires
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (t *UserApiToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// IsExpired reports whether the token has passed its expiry time.
func (t *UserApiToken) IsExpired(now time.Time) bool {
	return t.ExpiresAt != nil && now.After(*t.ExpiresAt)
}

// OAuthAuthRequest backs the OAuth 2.0 authorization-code grant (with PKCE)
// used by OAuth-aware MCP clients (e.g. Claude.ai). The client only needs the
// MCP endpoint URL; it discovers the authorization server via
// /.well-known/oauth-authorization-server, sends the user to /authorize, the
// user logs in and approves on a consent page, and the client exchanges the
// returned code for an access token at /token.
//
// Lifecycle: a row is created when /authorize is hit (Status == "pending"),
// flipped to "approved" (with UserID + Code set) once the user consents, and
// the Code is single-use — cleared on a successful /token exchange.
type OAuthAuthRequest struct {
	ID                  string    `gorm:"type:string;primaryKey" json:"id"`
	ClientID            string    `gorm:"index" json:"client_id"`             // OAuth client identifier (the MCP client)
	RedirectURI         string    `json:"redirect_uri"`                       // Where to send the user back with the code
	State               string    `json:"state"`                              // Opaque value echoed back to the client
	Scope               string    `json:"scope"`                              // Requested scope (e.g. "mcp")
	CodeChallenge       string    `json:"code_challenge"`                     // PKCE challenge
	CodeChallengeMethod string    `json:"code_challenge_method"`              // "S256" or "plain"
	Status              string    `gorm:"index;default:'pending'" json:"status"` // "pending" | "approved" | "denied"
	UserID              *string   `gorm:"index" json:"user_id"`               // Set once the user approves
	Code                *string   `gorm:"uniqueIndex" json:"-"`               // Authorization code, single-use
	ExpiresAt           time.Time `json:"expires_at"`                         // Request/code expiry (~10 min)
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (o *OAuthAuthRequest) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}

// IsExpired reports whether the authorization request/code has expired.
func (o *OAuthAuthRequest) IsExpired(now time.Time) bool {
	return now.After(o.ExpiresAt)
}
