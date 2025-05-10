package database

import (
	"time"

	"github.com/google/uuid"
)

// Alias model for mapping filenames to aliases
type Alias struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FileName string    `gorm:"uniqueIndex"`
	Alias    string    `gorm:"uniqueIndex"`
	Port     int
	IsActive bool `gorm:"default:false"`
}

// SystemConfig model for storing system configuration
type SystemConfig struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Key         string    `gorm:"uniqueIndex"`
	Value       string
	Type        string    `gorm:"default:string"` // string, number, boolean, json
	Description string    `gorm:"default:''"`     // optional description
	HideValue   bool      `gorm:"default:false"`  // hide value in the UI
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
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
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string         `gorm:"uniqueIndex"` // Used as subdomain or slug
	Mode          ProjectMode    `gorm:"type:text"`
	ActiveProxyID *uuid.UUID     `gorm:"type:uuid"`
	ActiveProxy   *ProxyTarget   `gorm:"foreignKey:ActiveProxyID"`
	Endpoints     []MockEndpoint `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	ProxyTargets  []ProxyTarget  `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ProxyTarget defines forward request destination if project mode is proxy or forwarder
type ProxyTarget struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID uuid.UUID `gorm:"type:uuid"`
	Label     string    // Example: "Staging", "Production"
	URL       string    // Example: "https://staging.example.com"
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MockEndpoint represents an HTTP route that is mocked
type MockEndpoint struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID    uuid.UUID      `gorm:"type:uuid"`
	Method       string         // GET, POST, PUT, DELETE, etc
	Path         string         // Example: "/users/:id"
	Enabled      bool           // Whether endpoint is active or not
	ResponseMode string         // "static", "random", "round_robin"
	Responses    []MockResponse `gorm:"foreignKey:EndpointID;constraint:OnDelete:CASCADE"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// MockResponse represents possible responses from an endpoint
type MockResponse struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EndpointID uuid.UUID `gorm:"type:uuid"`
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

// MockRule represents filter rules for selecting responses
type MockRule struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ResponseID uuid.UUID `gorm:"type:uuid"`
	Type       string    // "header", "body", "query", "path"
	Key        string    // Example: "X-Auth", "q", "user.id"
	Operator   string    // "equals", "contains", "regex"
	Value      string
}
