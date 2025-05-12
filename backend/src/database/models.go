package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
	ModeMock      ProjectMode = "mock"      // Serve mock responses
	ModeProxy     ProjectMode = "proxy"     // Forward to target, no modification
	ModeForwarder ProjectMode = "forwarder" // Forward with recording
	ModeDisabled  ProjectMode = "disabled"  // Disabled, no responses
)

// Project represents one group of endpoints, accessible via subdomain or alias
type Project struct {
	ID            string         `gorm:"type:string;primaryKey" json:"id"`
	Name          string         `gorm:"type:string" json:"name"`
	Mode          ProjectMode    `gorm:"type:string" json:"mode"`
	Status        string         `gorm:"type:string;default:'running'" json:"status"` // running, stopped, error
	ActiveProxyID *string        `gorm:"type:string" json:"active_proxy_id"`
	ActiveProxy   *ProxyTarget   `gorm:"foreignKey:ActiveProxyID" json:"active_proxy"`
	Endpoints     []MockEndpoint `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"endpoints"`
	ProxyTargets  []ProxyTarget  `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"proxy_targets"`
	Alias         string         `gorm:"type:string;uniqueIndex;not null" json:"alias"` // Subdomain or alias for the project
	URL           string         `json:"url"`                                           // URL for the project, e.g. "https://example.com" this is used for FE only
	Documentation string         `gorm:"type:string" json:"documentation"`              // Documentation URL or text
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
	Responses     []MockResponse `gorm:"foreignKey:EndpointID;constraint:OnDelete:CASCADE;" json:"responses"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
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
	ID            string     `gorm:"type:string;primaryKey" json:"id"`
	EndpointID    string     `gorm:"type:string" json:"endpoint_id"`
	StatusCode    int        `json:"status_code"`                    // HTTP status code
	Body          string     `gorm:"type:text" json:"body"`          // Response body, stored as JSON
	Headers       string     `gorm:"type:text" json:"headers"`       // Headers stored as JSON
	Priority      int        `json:"priority"`                       // Priority if ResponseMode = static
	DelayMS       int        `json:"delay_ms"`                       // Delay before response (milliseconds)
	Stream        bool       `json:"stream"`                         // True if response is stream (e.g. SSE, chunked)
	Documentation string     `gorm:"type:text" json:"documentation"` // Documentation URL or text
	Enabled       bool       `json:"enabled" gorm:"default:true"`    // Whether enabled or not
	Rules         []MockRule `gorm:"foreignKey:ResponseID;constraint:OnDelete:CASCADE" json:"rules"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
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

// TODO multi user, multi workspace support and multi sso
// User represents an individual who can log in to the system via SSO or password.
// A user can belong to multiple workspaces.
type User struct {
	ID         string          `gorm:"type:string;primaryKey" json:"id"`    // Unique user ID
	Email      string          `gorm:"uniqueIndex" json:"email"`            // Unique email (used for login/identity)
	Name       string          `json:"name"`                                // Display name
	Password   string          `json:"-"`                                   // Optional (if login via password)
	IsOwner    bool            `gorm:"default:false" json:"is_owner"`       // System-wide owner (can manage SSO configs, manage all workspaces and etc)
	Identities []UserIdentity  `gorm:"foreignKey:UserID" json:"identities"` // Linked SSO accounts
	Workspaces []UserWorkspace `gorm:"foreignKey:UserID" json:"workspaces"` // Memberships in workspaces
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
	ID        string          `gorm:"type:string;primaryKey" json:"id"`                                   // Unique workspace ID
	Name      string          `gorm:"uniqueIndex" json:"name"`                                            // Unique workspace name
	Projects  []Project       `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"projects"` // Projects under this workspace
	Members   []UserWorkspace `gorm:"foreignKey:WorkspaceID" json:"members"`                              // User membership records
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
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
