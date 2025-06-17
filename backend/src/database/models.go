package database

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AdvanceConfigProject defines advance configuration structure for projects
type AdvanceConfigProject struct {
	DelayMs int `json:"delayMs,omitempty"` // Response delay in milliseconds (0-120000)
}

// AdvanceConfigEndpoint defines advance configuration structure for endpoints
type AdvanceConfigEndpoint struct {
	DelayMs int `json:"delayMs,omitempty"` // Response delay in milliseconds (0-120000)
}

// Validate validates the project advance configuration
func (a *AdvanceConfigProject) Validate() error {
	if a.DelayMs < 0 {
		return errors.New("delayMs cannot be negative")
	}
	if a.DelayMs > 120000 {
		return errors.New("delayMs cannot exceed 120000ms (2 minutes)")
	}
	return nil
}

// Validate validates the endpoint advance configuration
func (a *AdvanceConfigEndpoint) Validate() error {
	if a.DelayMs < 0 {
		return errors.New("delayMs cannot be negative")
	}
	if a.DelayMs > 120000 {
		return errors.New("delayMs cannot exceed 120000ms (2 minutes)")
	}
	return nil
}

// ParseProjectAdvanceConfig parses JSON string to AdvanceConfigProject struct
func ParseProjectAdvanceConfig(configJSON string) (*AdvanceConfigProject, error) {
	if configJSON == "" {
		return &AdvanceConfigProject{}, nil
	}

	var config AdvanceConfigProject
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, errors.New("invalid JSON format in advance_config")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// ParseEndpointAdvanceConfig parses JSON string to AdvanceConfigEndpoint struct
func ParseEndpointAdvanceConfig(configJSON string) (*AdvanceConfigEndpoint, error) {
	if configJSON == "" {
		return &AdvanceConfigEndpoint{}, nil
	}

	var config AdvanceConfigEndpoint
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, errors.New("invalid JSON format in advance_config")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// ToJSON converts AdvanceConfigProject to JSON string
func (a *AdvanceConfigProject) ToJSON() (string, error) {
	if a.DelayMs == 0 {
		return "", nil
	}

	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJSON converts AdvanceConfigEndpoint to JSON string
func (a *AdvanceConfigEndpoint) ToJSON() (string, error) {
	if a.DelayMs == 0 {
		return "", nil
	}

	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SystemConfig model for storing system configuration
type SystemConfig struct {
	ID          string    `gorm:"type:string;primaryKey" json:"id"`
	Key         string    `gorm:"uniqueIndex" json:"key"`          // Unique key for the config
	Value       string    `json:"value"`                           // Value of the config
	Type        string    `gorm:"default:string" json:"type"`      // string, number, boolean, json
	Description string    `gorm:"default:''" json:"description"`   // optional description
	HideValue   bool      `gorm:"default:false" json:"hide_value"` // hide value in the UI
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate hook to generate UUID string
func (s *SystemConfig) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// ProjectMode defines operation mode of mock system per project
type ProjectMode string

const (
	ModeMock      ProjectMode = "mock"      // Serves predefined mock responses only
	ModeProxy     ProjectMode = "proxy"     // Uses mocks when available, otherwise forwards requests
	ModeForwarder ProjectMode = "forwarder" // Always forwards all requests to target endpoint
	ModeDisabled  ProjectMode = "disabled"  // Endpoint inactive - no responses served
)

// Project represents one group of endpoints, accessible via subdomain or alias
type Project struct {
	ID            string         `gorm:"type:string;primaryKey" json:"id"`
	Name          string         `gorm:"type:string" json:"name"`
	WorkspaceID   string         `gorm:"type:string;index" json:"workspace_id"`       // Foreign key to the associated workspace
	Mode          ProjectMode    `gorm:"type:string;default:'mock'" json:"mode"`      // default: mock
	Status        string         `gorm:"type:string;default:'running'" json:"status"` // running, stopped, error
	ActiveProxyID *string        `gorm:"type:string" json:"active_proxy_id"`
	ActiveProxy   *ProxyTarget   `gorm:"foreignKey:ActiveProxyID" json:"active_proxy"`
	Endpoints     []MockEndpoint `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"endpoints"`
	ProxyTargets  []ProxyTarget  `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"proxy_targets"`
	Alias         string         `gorm:"type:string;uniqueIndex;not null" json:"alias"` // Subdomain or alias for the project
	URL           string         `json:"url"`                                           // URL for the project, e.g. "https://example.com" this is used for FE only
	Documentation string         `gorm:"type:string" json:"documentation"`              // Documentation URL or text
	AdvanceConfig string         `gorm:"type:text" json:"advance_config"`               // Advanced configuration (e.g. global timeout, rate limiting) as JSON string
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate hook to generate UUID string
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// project table name
func (p *Project) TableName() string {
	return "projects"
}

// ProxyTarget defines forward request destination if project mode is proxy or forwarder
type ProxyTarget struct {
	ID        string    `gorm:"type:string;primaryKey" json:"id"`
	ProjectID string    `gorm:"type:string" json:"project_id"`
	Label     string    `json:"label"` // Example: "Staging", "Production"
	URL       string    `json:"url"`   // Example: "https://staging.example.com"
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate hook to generate UUID string
func (pt *ProxyTarget) BeforeCreate(tx *gorm.DB) error {
	if pt.ID == "" {
		pt.ID = uuid.New().String()
	}
	return nil
}

// MockEndpoint represents an HTTP route that is mocked
type MockEndpoint struct {
	ID            string         `gorm:"type:string;primaryKey" json:"id"`
	ProjectID     string         `gorm:"type:string" json:"project_id"`
	Method        string         `json:"method"`                                // GET, POST, PUT, DELETE, etc
	Path          string         `json:"path"`                                  // Example: "/users/:id"
	Enabled       bool           `json:"enabled" gorm:"default:true"`           // Whether endpoint is active or not
	ResponseMode  string         `json:"response_mode" gorm:"default:'random'"` // "static", "random", "round_robin"
	Documentation string         `gorm:"type:text" json:"documentation"`        // Documentation URL or text
	AdvanceConfig string         `gorm:"type:text" json:"advance_config"`       // Advanced configuration (e.g. timeout) as JSON string
	Responses     []MockResponse `gorm:"foreignKey:EndpointID;constraint:OnDelete:CASCADE;" json:"responses"`
	// Proxy configuration for endpoint-level proxying
	UseProxy      bool         `json:"use_proxy" gorm:"default:false"`               // Whether to use proxy for this endpoint
	ProxyTargetID *string      `gorm:"type:string" json:"proxy_target_id"`           // ID of the proxy target to use
	ProxyTarget   *ProxyTarget `gorm:"foreignKey:ProxyTargetID" json:"proxy_target"` // The associated proxy target
	CreatedAt     time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate hook to generate UUID string
func (me *MockEndpoint) BeforeCreate(tx *gorm.DB) error {
	if me.ID == "" {
		me.ID = uuid.New().String()
	}
	return nil
}

// MockResponse represents possible responses from an endpoint
type MockResponse struct {
	ID         string     `gorm:"type:string;primaryKey" json:"id"`
	EndpointID string     `gorm:"type:string" json:"endpoint_id"`
	StatusCode int        `json:"status_code"`                 // HTTP status code
	Body       string     `gorm:"type:text" json:"body"`       // Response body, stored as JSON
	Headers    string     `gorm:"type:text" json:"headers"`    // Headers stored as JSON
	Priority   int        `json:"priority"`                    // Priority if ResponseMode = static
	DelayMS    int        `json:"delay_ms"`                    // Delay before response (milliseconds)
	Stream     bool       `json:"stream"`                      // True if response is stream (e.g. SSE, chunked)
	Note       string     `gorm:"type:text" json:"note"`       // Optional note for the response
	Enabled    bool       `json:"enabled" gorm:"default:true"` // Whether enabled or not
	Rules      []MockRule `gorm:"foreignKey:ResponseID;constraint:OnDelete:CASCADE" json:"rules"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate hook to generate UUID string
func (mr *MockResponse) BeforeCreate(tx *gorm.DB) error {
	if mr.ID == "" {
		mr.ID = uuid.New().String()
	}
	return nil
}

// MockRule represents filter rules for selecting responses
type MockRule struct {
	ID         string `gorm:"type:string;primaryKey" json:"id"`
	ResponseID string `gorm:"type:string" json:"response_id"`
	Type       string `json:"type"`     // "header", "body", "query", "path"
	Key        string `json:"key"`      // Example: "X-Auth", "q", "user.id"
	Operator   string `json:"operator"` // "equals", "contains", "regex"
	Value      string `json:"value"`
}

// BeforeCreate hook to generate UUID string
func (mr *MockRule) BeforeCreate(tx *gorm.DB) error {
	if mr.ID == "" {
		mr.ID = uuid.New().String()
	}
	return nil
}

// SourceRequest defines the source of the request log.
type SourceRequest string

const (
	RequestSourceUnknown SourceRequest = ""
	RequestSourceEcho    SourceRequest = "echo"
	RequestSourceReplay  SourceRequest = "replay"
)

// RequestLog stores detailed information about each incoming HTTP request.
// It captures how the request was handled (mock, proxy, forwarder), whether it matched a mock endpoint,
// and includes raw request/response data for auditing or debugging.
type RequestLog struct {
	ID              string `gorm:"type:string;primaryKey" json:"id"`    // Unique identifier (UUID)
	ProjectID       string `gorm:"type:string;index" json:"project_id"` // Foreign key to the associated project
	Method          string `json:"method"`                              // HTTP method (GET, POST, etc.)
	Path            string `json:"path"`                                // Request path (e.g. "/api/users")
	QueryParams     string `gorm:"type:text" json:"query_params"`       // Query parameters (stored as JSON string)
	RequestHeaders  string `gorm:"type:text" json:"request_headers"`    // Request headers as array of key-value pairs
	RequestBody     string `gorm:"type:text" json:"request_body"`       // Raw request body
	ResponseStatus  int    `json:"response_status"`                     // HTTP status code returned
	ResponseBody    string `gorm:"type:text" json:"response_body"`      // Raw response body
	ResponseHeaders string `gorm:"type:text" json:"response_headers"`   // Response headers as array of key-value pairs
	LatencyMS       int    `json:"latency_ms"`                          // Time taken to respond or delay applied (in milliseconds)
	Bookmark        bool   `gorm:"type:bool" json:"bookmark"`           // Optional bookmark for easy reference
	LogsHash        string `gorm:"type:string" json:"logs_hash"`        // Hash of the response body for integrity checks + jwt signature

	Source SourceRequest `gorm:"size:50;not null default:''" json:"source"` // Source of the request: "replay", "echo", etc.

	// ExecutionMode indicates the handling logic used for this request.
	// Values follow ProjectMode: "mock", "proxy", "forwarder", etc.
	ExecutionMode ProjectMode `gorm:"type:string" json:"execution_mode"`

	// Matched is true if the request matched an existing mock endpoint.
	Matched   bool      `gorm:"default:false" json:"matched"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // Timestamp of the request

	// Association to the Project
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate hook generates UUID before inserting into database
func (rl *RequestLog) BeforeCreate(tx *gorm.DB) error {
	if rl.ID == "" {
		rl.ID = uuid.New().String()
	}
	return nil
}

// TODO: multi sso
// User represents an individual who can log in to the system via SSO or password.
// A user can belong to multiple workspaces.
// password is Argon2id hashed password (when using password login).
// Salt is generated using unixtime + random number (8 bytes total).
// NOT recommended for cryptographic use — better to use crypto/rand if possible.
type User struct {
	ID         string          `gorm:"type:string;primaryKey" json:"id"`                                // Unique user ID
	Email      string          `gorm:"uniqueIndex" json:"email"`                                        // Unique email (used for login/identity)
	Name       string          `json:"name"`                                                            // Display name
	Password   string          `json:"-"`                                                               // Argon2id hashed password (when using password login)
	IsOwner    bool            `gorm:"default:false" json:"is_owner"`                                   // System-wide owner (can manage SSO configs, manage all workspaces and etc)
	IsActive   bool            `gorm:"default:true" json:"is_active"`                                   // Whether this user account is active
	Identities []UserIdentity  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"identities"` // Linked SSO accounts
	Workspaces []UserWorkspace `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"workspaces"` // Memberships in workspaces
	CreatedAt  time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// UserIdentity links a user to an external SSO provider (Google, GitHub, etc).
// Used to authenticate users without passwords.
// Each combination of Provider + ProviderID must be unique.
type UserIdentity struct {
	ID          string    `gorm:"type:string;primaryKey" json:"id"`              // Unique identity ID
	UserID      string    `gorm:"index" json:"user_id"`                          // Linked internal user ID
	Provider    string    `gorm:"index" json:"provider"`                         // e.g. "google", "github"
	ProviderID  string    `gorm:"index" json:"provider_id"`                      // Unique user ID from provider (e.g. Google "sub" claim)
	Email       string    `json:"email"`                                         // Email from provider (for display/debug)
	Name        string    `json:"name"`                                          // Display name from provider
	AvatarURL   string    `json:"avatar_url"`                                    // Profile image
	AccessToken string    `json:"-"`                                             // Optional: OAuth access token (not returned in JSON)
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // Associated user
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (ui *UserIdentity) BeforeCreate(tx *gorm.DB) error {
	if ui.ID == "" {
		ui.ID = uuid.New().String()
	}
	return nil
}

// Workspace represents a shared space (team or organization) that can contain projects.
// A user can belong to multiple workspaces, and each workspace can have multiple users.
type Workspace struct {
	ID                string          `gorm:"type:string;primaryKey" json:"id"`                                   // Unique workspace ID
	Name              string          `json:"name"`                                                               // Unique workspace name
	Projects          []Project       `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"projects"` // Projects under this workspace
	Members           []UserWorkspace `gorm:"foreignKey:WorkspaceID" json:"members"`                              // User membership records
	AutoInviteDomains string          `gorm:"type:text" json:"auto_invite_domains"`                               // Comma-separated list of email domains for auto-invitation
	AutoInviteEnabled bool            `gorm:"default:false" json:"auto_invite_enabled"`                           // Whether auto-invitation is enabled for this workspace
	AutoInviteRole    string          `gorm:"default:'member'" json:"auto_invite_role"`                           // Role assigned to auto-invited users ("admin" or "member")
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

func (w *Workspace) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}

// UserWorkspace defines the relationship and role of a user in a workspace.
// Used to implement multi-workspace and per-workspace roles (e.g., owner/admin/member).
type UserWorkspace struct {
	ID          string    `gorm:"type:string;primaryKey" json:"id"` // Unique record ID
	UserID      string    `gorm:"index" json:"user_id"`             // Linked user
	WorkspaceID string    `gorm:"index" json:"workspace_id"`        // Linked workspace
	Role        string    `gorm:"default:'member'" json:"role"`     // "admin", or "member" and only admin can manage user in workspace
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	// User and Workspace are backrefs to the User and Workspace models.
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`      // Backref to User
	Workspace Workspace `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE"` // Backref to Workspace
}

func (uw *UserWorkspace) BeforeCreate(tx *gorm.DB) error {
	if uw.ID == "" {
		uw.ID = uuid.New().String()
	}
	return nil
}

// SSOConfig stores global SSO (OAuth) settings for each provider.
// Used to configure OAuth for Google, GitHub, etc.
// Only one record per provider is allowed.
type SSOConfig struct {
	ID        string    `gorm:"type:string;primaryKey" json:"id"` // Unique config ID
	Provider  string    `gorm:"uniqueIndex" json:"provider"`      // e.g. "google", "github"
	Config    string    `gorm:"type:text" json:"config"`          // JSON config (client_id, secret, redirect_url, etc.)
	Enabled   bool      `gorm:"default:true" json:"enabled"`      // Whether this provider is active
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (s *SSOConfig) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

type ReplayProtocol string

const (
	ReplayProtocolHTTP ReplayProtocol = "http" // HTTP/HTTPS protocol
)

// Replay stores a preset request configuration to be executed for testing or mocking purposes.
// Simplified model for easier management and API interaction.
type Replay struct {
	ID        string  `gorm:"primaryKey;type:TEXT" json:"id"`   // Unique identifier (UUID)
	Name      string  `json:"name"`                             // User-defined name for this replay
	ProjectID string  `gorm:"index;not null" json:"project_id"` // Project scoping
	FolderID  *string `gorm:"index" json:"folder_id"`           // Optional folder location

	Protocol ReplayProtocol `gorm:"not null;default:'http'" json:"protocol"` // Protocol: http, https (simplified to HTTP only for now)
	Method   string         `gorm:"size:20;not null" json:"method"`          // HTTP method (GET, POST, PUT, DELETE, etc.)
	Url      string         `gorm:"not null" json:"url"`                     // Target URL or endpoint

	Config   string `gorm:"type:text" json:"config"`   // Additional configuration (e.g. timeout, retries) as JSON string
	Metadata string `gorm:"type:text" json:"metadata"` // Additional metadata (e.g. tags, notes) as JSON string

	Headers string `gorm:"type:text" json:"headers"` // Headers as JSON string (key-value pairs)
	Payload string `gorm:"type:text" json:"payload"` // Request payload/body

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // Timestamp of creation
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // Timestamp of last update
}

func (uw *Replay) BeforeCreate(tx *gorm.DB) error {
	if uw.ID == "" {
		uw.ID = uuid.New().String()
	}
	return nil
}

// ReplayFolder represents a hierarchical folder to group replays.
// Folders can be nested via ParentID and scoped per project.
type ReplayFolder struct {
	ID        string  `gorm:"primaryKey;type:TEXT" json:"id"`   // Unique identifier (UUID)
	Name      string  `gorm:"not null" json:"name"`             // Folder name
	ParentID  *string `gorm:"type:TEXT;index" json:"parent_id"` // Optional parent folder (null = root)
	ProjectID string  `gorm:"index;not null" json:"project_id"` // Project scoping

	Children []ReplayFolder `gorm:"foreignKey:ParentID" json:"children,omitempty"` // Subfolders
	Replays  []Replay       `gorm:"foreignKey:FolderID" json:"replays,omitempty"`  // Replays inside this folder

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // Timestamp of creation
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // Timestamp of last update
}

func (uw *ReplayFolder) BeforeCreate(tx *gorm.DB) error {
	if uw.ID == "" {
		uw.ID = uuid.New().String()
	}
	return nil
}
