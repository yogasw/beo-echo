package models

import (
	"time"
)

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
	ID            uint        `gorm:"primaryKey"`
	Name          string      `gorm:"uniqueIndex"` // Used as subdomain or slug
	Mode          ProjectMode `gorm:"type:text"`
	ActiveProxyID *uint
	ActiveProxy   *ProxyTarget   `gorm:"foreignKey:ActiveProxyID"`
	Endpoints     []MockEndpoint `gorm:"constraint:OnDelete:CASCADE"`
	ProxyTargets  []ProxyTarget  `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ProxyTarget defines forward request destination if project mode is proxy or forwarder
type ProxyTarget struct {
	ID        uint `gorm:"primaryKey"`
	ProjectID uint
	Label     string // Example: "Staging", "Production"
	URL       string // Example: "https://staging.example.com"
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MockEndpoint represents an HTTP route that is mocked
type MockEndpoint struct {
	ID           uint `gorm:"primaryKey"`
	ProjectID    uint
	Method       string         // GET, POST, PUT, DELETE, etc
	Path         string         // Example: "/users/:id"
	Enabled      bool           // Whether endpoint is active or not
	ResponseMode string         // "static", "random", "round_robin"
	Responses    []MockResponse `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// MockResponse represents possible responses from an endpoint
type MockResponse struct {
	ID         uint `gorm:"primaryKey"`
	EndpointID uint
	StatusCode int
	Body       string     // Response body, can be raw text or JSON
	Headers    string     // JSON string: {"Content-Type":"application/json"}
	Priority   int        // Priority if ResponseMode = static
	DelayMS    int        // Delay before response (milliseconds)
	Stream     bool       // True if response is stream (e.g. SSE, chunked)
	Active     bool       // Whether active or not
	Rules      []MockRule `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// MockRule represents filter rules for selecting responses
type MockRule struct {
	ID         uint `gorm:"primaryKey"`
	ResponseID uint
	Type       string // "header", "body", "query", "path"
	Key        string // Example: "X-Auth", "q", "user.id"
	Operator   string // "equals", "contains", "regex"
	Value      string
}
