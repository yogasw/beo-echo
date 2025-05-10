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
	Method        string         `json:"method"`                         // GET, POST, PUT, DELETE, etc
	Path          string         `json:"path"`                           // Example: "/users/:id"
	Enabled       bool           `json:"enabled" gorm:"default:true"`    // Whether endpoint is active or not
	ResponseMode  string         `json:"response_mode"`                  // "static", "random", "round_robin"
	Documentation string         `gorm:"type:text" json:"documentation"` // Documentation URL or text
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
	Body          string     `json:"body"`                           // Response body, can be raw text or JSON
	Headers       string     `json:"headers"`                        // JSON string: "Content-Type":"application/json"}
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
