package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Alias model for mapping filenames to aliases
type Alias struct {
	ID       string `gorm:"type:string;primaryKey"`
	FileName string `gorm:"uniqueIndex"`
	Alias    string `gorm:"uniqueIndex"`
	Port     int
	IsActive bool `gorm:"default:false"`
}

// BeforeCreate hook to generate UUID string
func (a *Alias) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// SystemConfig model for storing system configuration
type SystemConfig struct {
	ID          string `gorm:"type:string;primaryKey"`
	Key         string `gorm:"uniqueIndex"`
	Value       string
	Type        string    `gorm:"default:string"` // string, number, boolean, json
	Description string    `gorm:"default:''"`     // optional description
	HideValue   bool      `gorm:"default:false"`  // hide value in the UI
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
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
	ID            string         `gorm:"type:string;primaryKey"`
	Name          string         `gorm:"uniqueIndex"` // Used as subdomain or slug
	Mode          ProjectMode    `gorm:"type:text"`
	ActiveProxyID *string        `gorm:"type:string"`
	ActiveProxy   *ProxyTarget   `gorm:"foreignKey:ActiveProxyID"`
	Endpoints     []MockEndpoint `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	ProxyTargets  []ProxyTarget  `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
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
	ID        string `gorm:"type:string;primaryKey"`
	ProjectID string `gorm:"type:string"`
	Label     string // Example: "Staging", "Production"
	URL       string // Example: "https://staging.example.com"
	CreatedAt time.Time
	UpdatedAt time.Time
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
	ID           string         `gorm:"type:string;primaryKey"`
	ProjectID    string         `gorm:"type:string"`
	Method       string         // GET, POST, PUT, DELETE, etc
	Path         string         // Example: "/users/:id"
	Enabled      bool           // Whether endpoint is active or not
	ResponseMode string         // "static", "random", "round_robin"
	Responses    []MockResponse `gorm:"foreignKey:EndpointID;constraint:OnDelete:CASCADE"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	ID         string `gorm:"type:string;primaryKey"`
	EndpointID string `gorm:"type:string"`
	StatusCode int
	Body       string     // Response body, can be raw text or JSON
	Headers    string     // JSON string: {"Content-Type":"application/json"}
	Priority   int        // Priority if ResponseMode = static
	DelayMS    int        // Delay before response (milliseconds)
	Stream     bool       // True if response is stream (e.g. SSE, chunked)
	Active     bool       // Whether active or not
	Rules      []MockRule `gorm:"foreignKey:ResponseID;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
	ID         string `gorm:"type:string;primaryKey"`
	ResponseID string `gorm:"type:string"`
	Type       string // "header", "body", "query", "path"
	Key        string // Example: "X-Auth", "q", "user.id"
	Operator   string // "equals", "contains", "regex"
	Value      string
}

// BeforeCreate hook to generate UUID string
func (mr *MockRule) BeforeCreate(tx *gorm.DB) error {
	if mr.ID == "" {
		mr.ID = uuid.New().String()
	}
	return nil
}
