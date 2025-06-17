package database

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ProjectAdvanceConfig represents the structure of advance config for projects
type ProjectAdvanceConfig struct {
	GlobalTimeout    int                    `json:"global_timeout,omitempty"`    // Global timeout in milliseconds
	RateLimit        *RateLimitConfig       `json:"rate_limit,omitempty"`        // Rate limiting configuration
	CORS             *CORSConfig            `json:"cors,omitempty"`              // CORS configuration
	Security         *SecurityConfig        `json:"security,omitempty"`          // Security configuration
	CustomHeaders    map[string]string      `json:"custom_headers,omitempty"`    // Global custom headers
	Middleware       []string               `json:"middleware,omitempty"`        // List of middleware to apply
	ProxyTimeout     int                    `json:"proxy_timeout,omitempty"`     // Proxy timeout in milliseconds
	RetryPolicy      *RetryPolicyConfig     `json:"retry_policy,omitempty"`      // Retry policy for proxy/forward mode
	Logging          *LoggingConfig         `json:"logging,omitempty"`           // Logging configuration
	Cache            *CacheConfig           `json:"cache,omitempty"`             // Caching configuration
	Compression      *CompressionConfig     `json:"compression,omitempty"`       // Response compression settings
	CustomProperties map[string]interface{} `json:"custom_properties,omitempty"` // Additional custom properties
}

type RateLimitConfig struct {
	Enabled     bool `json:"enabled"`          // Whether rate limiting is enabled
	RequestsRPM int  `json:"requests_per_min"` // Requests per minute
	BurstSize   int  `json:"burst_size"`       // Burst size for rate limiting
}

type CORSConfig struct {
	Enabled          bool     `json:"enabled"`                     // Whether CORS is enabled
	AllowedOrigins   []string `json:"allowed_origins,omitempty"`   // Allowed origins
	AllowedMethods   []string `json:"allowed_methods,omitempty"`   // Allowed HTTP methods
	AllowedHeaders   []string `json:"allowed_headers,omitempty"`   // Allowed headers
	ExposedHeaders   []string `json:"exposed_headers,omitempty"`   // Headers exposed to the browser
	AllowCredentials bool     `json:"allow_credentials,omitempty"` // Whether to allow credentials
	MaxAge           int      `json:"max_age,omitempty"`           // Max age for preflight requests
}

type SecurityConfig struct {
	EnableHTTPS           bool                   `json:"enable_https,omitempty"`           // Force HTTPS
	EnableHSTS            bool                   `json:"enable_hsts,omitempty"`            // HTTP Strict Transport Security
	CSP                   map[string]string      `json:"csp,omitempty"`                    // Content Security Policy
	AllowedIPRanges       []string               `json:"allowed_ip_ranges,omitempty"`      // Allowed IP ranges
	BlockedIPRanges       []string               `json:"blocked_ip_ranges,omitempty"`      // Blocked IP ranges
	RequireAuthentication bool                   `json:"require_authentication,omitempty"` // Require authentication for all endpoints
	CustomSecurityRules   map[string]interface{} `json:"custom_security_rules,omitempty"`  // Custom security rules
}

type RetryPolicyConfig struct {
	Enabled     bool   `json:"enabled"`      // Whether retry is enabled
	MaxRetries  int    `json:"max_retries"`  // Maximum number of retries
	RetryDelay  int    `json:"retry_delay"`  // Delay between retries in milliseconds
	BackoffType string `json:"backoff_type"` // "fixed", "exponential", "linear"
}

type LoggingConfig struct {
	Enabled         bool     `json:"enabled"`                    // Whether custom logging is enabled
	Level           string   `json:"level,omitempty"`            // Log level: "debug", "info", "warn", "error"
	IncludeHeaders  bool     `json:"include_headers,omitempty"`  // Include request/response headers in logs
	IncludeBody     bool     `json:"include_body,omitempty"`     // Include request/response body in logs
	ExcludePaths    []string `json:"exclude_paths,omitempty"`    // Paths to exclude from logging
	SensitiveFields []string `json:"sensitive_fields,omitempty"` // Fields to mask in logs
}

type CacheConfig struct {
	Enabled    bool   `json:"enabled"`               // Whether caching is enabled
	TTL        int    `json:"ttl,omitempty"`         // Time to live in seconds
	MaxSize    int    `json:"max_size,omitempty"`    // Maximum cache size in MB
	CacheKey   string `json:"cache_key,omitempty"`   // Custom cache key pattern
	VaryHeader string `json:"vary_header,omitempty"` // Header to vary cache by
}

type CompressionConfig struct {
	Enabled      bool     `json:"enabled"`                 // Whether compression is enabled
	MinSize      int      `json:"min_size,omitempty"`      // Minimum response size to compress (bytes)
	ContentTypes []string `json:"content_types,omitempty"` // Content types to compress
	Level        int      `json:"level,omitempty"`         // Compression level (1-9)
}

func TestProjectCreation(t *testing.T) {
	// Setup test database
	err := CheckAndHandle()
	require.NoError(t, err, "Failed to setup test database")

	t.Run("Create Project Without Advance Config", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := CreateTestWorkspace("test@example.com", "Test User", "Test Workspace")
		require.NoError(t, err)
		defer CleanupTestData(user.ID, workspace.ID, "", "")

		// Create project
		project, err := CreateTestProject(workspace.ID, "Test Project", "test-project")
		require.NoError(t, err)
		assert.NotEmpty(t, project.ID)
		assert.Equal(t, "Test Project", project.Name)
		assert.Equal(t, workspace.ID, project.WorkspaceID)
		assert.Equal(t, ModeMock, project.Mode)
		assert.Equal(t, "running", project.Status)
		assert.Equal(t, "test-project", project.Alias)
		assert.Empty(t, project.AdvanceConfig)
	})

	t.Run("Create Project With Empty Advance Config", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := CreateTestWorkspace("test2@example.com", "Test User 2", "Test Workspace 2")
		require.NoError(t, err)
		defer CleanupTestData(user.ID, workspace.ID, "", "")

		// Create project with empty advance config
		project, err := CreateTestProjectWithConfig(workspace.ID, "Test Project 2", "test-project-2", "")
		require.NoError(t, err)
		assert.NotEmpty(t, project.ID)
		assert.Empty(t, project.AdvanceConfig)
	})

	t.Run("Create Project With Basic Advance Config", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := CreateTestWorkspace("test3@example.com", "Test User 3", "Test Workspace 3")
		require.NoError(t, err)
		defer CleanupTestData(user.ID, workspace.ID, "", "")

		// Create advance config
		advanceConfig := ProjectAdvanceConfig{
			GlobalTimeout: 5000,
			RateLimit: &RateLimitConfig{
				Enabled:     true,
				RequestsRPM: 100,
				BurstSize:   10,
			},
		}

		configJSON, err := json.Marshal(advanceConfig)
		require.NoError(t, err)

		// Create project with advance config
		project, err := CreateTestProjectWithConfig(workspace.ID, "Test Project 3", "test-project-3", string(configJSON))
		require.NoError(t, err)
		assert.NotEmpty(t, project.ID)
		assert.NotEmpty(t, project.AdvanceConfig)

		// Verify the stored config can be unmarshaled
		var storedConfig ProjectAdvanceConfig
		err = json.Unmarshal([]byte(project.AdvanceConfig), &storedConfig)
		require.NoError(t, err)
		assert.Equal(t, 5000, storedConfig.GlobalTimeout)
		assert.True(t, storedConfig.RateLimit.Enabled)
		assert.Equal(t, 100, storedConfig.RateLimit.RequestsRPM)
		assert.Equal(t, 10, storedConfig.RateLimit.BurstSize)
	})

	t.Run("Create Project With Full Advance Config", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := CreateTestWorkspace("test4@example.com", "Test User 4", "Test Workspace 4")
		require.NoError(t, err)
		defer CleanupTestData(user.ID, workspace.ID, "", "")

		// Create comprehensive advance config
		advanceConfig := ProjectAdvanceConfig{
			GlobalTimeout: 10000,
			ProxyTimeout:  15000,
			RateLimit: &RateLimitConfig{
				Enabled:     true,
				RequestsRPM: 200,
				BurstSize:   20,
			},
			CORS: &CORSConfig{
				Enabled:          true,
				AllowedOrigins:   []string{"https://example.com", "https://app.example.com"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
				AllowedHeaders:   []string{"Content-Type", "Authorization"},
				ExposedHeaders:   []string{"X-Total-Count"},
				AllowCredentials: true,
				MaxAge:           3600,
			},
			Security: &SecurityConfig{
				EnableHTTPS:           true,
				EnableHSTS:            true,
				AllowedIPRanges:       []string{"192.168.1.0/24", "10.0.0.0/8"},
				BlockedIPRanges:       []string{"192.168.100.0/24"},
				RequireAuthentication: false,
				CSP: map[string]string{
					"default-src": "'self'",
					"script-src":  "'self' 'unsafe-inline'",
				},
			},
			CustomHeaders: map[string]string{
				"X-API-Version": "v1",
				"X-Service":     "beo-echo",
			},
			Middleware: []string{"cors", "rate-limit", "security"},
			RetryPolicy: &RetryPolicyConfig{
				Enabled:     true,
				MaxRetries:  3,
				RetryDelay:  1000,
				BackoffType: "exponential",
			},
			Logging: &LoggingConfig{
				Enabled:         true,
				Level:           "info",
				IncludeHeaders:  true,
				IncludeBody:     false,
				ExcludePaths:    []string{"/health", "/metrics"},
				SensitiveFields: []string{"password", "token", "authorization"},
			},
			Cache: &CacheConfig{
				Enabled:    true,
				TTL:        300,
				MaxSize:    100,
				CacheKey:   "method:path:query",
				VaryHeader: "Accept-Encoding",
			},
			Compression: &CompressionConfig{
				Enabled:      true,
				MinSize:      1024,
				ContentTypes: []string{"application/json", "text/html", "text/plain"},
				Level:        6,
			},
			CustomProperties: map[string]interface{}{
				"environment": "production",
				"version":     "1.0.0",
				"features":    []string{"analytics", "monitoring"},
			},
		}

		configJSON, err := json.Marshal(advanceConfig)
		require.NoError(t, err)

		// Create project with full advance config
		project, err := CreateTestProjectWithConfig(workspace.ID, "Test Project 4", "test-project-4", string(configJSON))
		require.NoError(t, err)
		assert.NotEmpty(t, project.ID)
		assert.NotEmpty(t, project.AdvanceConfig)

		// Verify the stored config can be unmarshaled and contains all fields
		var storedConfig ProjectAdvanceConfig
		err = json.Unmarshal([]byte(project.AdvanceConfig), &storedConfig)
		require.NoError(t, err)

		// Verify basic settings
		assert.Equal(t, 10000, storedConfig.GlobalTimeout)
		assert.Equal(t, 15000, storedConfig.ProxyTimeout)

		// Verify rate limit config
		require.NotNil(t, storedConfig.RateLimit)
		assert.True(t, storedConfig.RateLimit.Enabled)
		assert.Equal(t, 200, storedConfig.RateLimit.RequestsRPM)
		assert.Equal(t, 20, storedConfig.RateLimit.BurstSize)

		// Verify CORS config
		require.NotNil(t, storedConfig.CORS)
		assert.True(t, storedConfig.CORS.Enabled)
		assert.Equal(t, []string{"https://example.com", "https://app.example.com"}, storedConfig.CORS.AllowedOrigins)
		assert.Equal(t, []string{"GET", "POST", "PUT", "DELETE"}, storedConfig.CORS.AllowedMethods)
		assert.True(t, storedConfig.CORS.AllowCredentials)
		assert.Equal(t, 3600, storedConfig.CORS.MaxAge)

		// Verify security config
		require.NotNil(t, storedConfig.Security)
		assert.True(t, storedConfig.Security.EnableHTTPS)
		assert.True(t, storedConfig.Security.EnableHSTS)
		assert.Equal(t, []string{"192.168.1.0/24", "10.0.0.0/8"}, storedConfig.Security.AllowedIPRanges)
		assert.Equal(t, []string{"192.168.100.0/24"}, storedConfig.Security.BlockedIPRanges)
		assert.False(t, storedConfig.Security.RequireAuthentication)

		// Verify custom headers
		assert.Equal(t, "v1", storedConfig.CustomHeaders["X-API-Version"])
		assert.Equal(t, "beo-echo", storedConfig.CustomHeaders["X-Service"])

		// Verify middleware
		assert.Equal(t, []string{"cors", "rate-limit", "security"}, storedConfig.Middleware)

		// Verify retry policy
		require.NotNil(t, storedConfig.RetryPolicy)
		assert.True(t, storedConfig.RetryPolicy.Enabled)
		assert.Equal(t, 3, storedConfig.RetryPolicy.MaxRetries)
		assert.Equal(t, 1000, storedConfig.RetryPolicy.RetryDelay)
		assert.Equal(t, "exponential", storedConfig.RetryPolicy.BackoffType)

		// Verify logging config
		require.NotNil(t, storedConfig.Logging)
		assert.True(t, storedConfig.Logging.Enabled)
		assert.Equal(t, "info", storedConfig.Logging.Level)
		assert.True(t, storedConfig.Logging.IncludeHeaders)
		assert.False(t, storedConfig.Logging.IncludeBody)
		assert.Equal(t, []string{"/health", "/metrics"}, storedConfig.Logging.ExcludePaths)
		assert.Equal(t, []string{"password", "token", "authorization"}, storedConfig.Logging.SensitiveFields)

		// Verify cache config
		require.NotNil(t, storedConfig.Cache)
		assert.True(t, storedConfig.Cache.Enabled)
		assert.Equal(t, 300, storedConfig.Cache.TTL)
		assert.Equal(t, 100, storedConfig.Cache.MaxSize)
		assert.Equal(t, "method:path:query", storedConfig.Cache.CacheKey)
		assert.Equal(t, "Accept-Encoding", storedConfig.Cache.VaryHeader)

		// Verify compression config
		require.NotNil(t, storedConfig.Compression)
		assert.True(t, storedConfig.Compression.Enabled)
		assert.Equal(t, 1024, storedConfig.Compression.MinSize)
		assert.Equal(t, []string{"application/json", "text/html", "text/plain"}, storedConfig.Compression.ContentTypes)
		assert.Equal(t, 6, storedConfig.Compression.Level)

		// Verify custom properties
		assert.Equal(t, "production", storedConfig.CustomProperties["environment"])
		assert.Equal(t, "1.0.0", storedConfig.CustomProperties["version"])
		features, ok := storedConfig.CustomProperties["features"].([]interface{})
		require.True(t, ok)
		assert.Len(t, features, 2)
		assert.Equal(t, "analytics", features[0])
		assert.Equal(t, "monitoring", features[1])
	})
}

func TestProjectAdvanceConfigValidation(t *testing.T) {
	t.Run("Invalid JSON in Advance Config", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := CreateTestWorkspace("test5@example.com", "Test User 5", "Test Workspace 5")
		require.NoError(t, err)
		defer CleanupTestData(user.ID, workspace.ID, "", "")

		// Create project with invalid JSON advance config
		invalidJSON := `{"global_timeout": 5000, "rate_limit": {`
		project, err := CreateTestProjectWithConfig(workspace.ID, "Test Project 5", "test-project-5", invalidJSON)

		// Project creation should succeed (database doesn't validate JSON)
		require.NoError(t, err)
		assert.Equal(t, invalidJSON, project.AdvanceConfig)

		// But JSON parsing should fail
		var config ProjectAdvanceConfig
		err = json.Unmarshal([]byte(project.AdvanceConfig), &config)
		assert.Error(t, err)
	})

	t.Run("Partial Config - Only Rate Limit", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := CreateTestWorkspace("test6@example.com", "Test User 6", "Test Workspace 6")
		require.NoError(t, err)
		defer CleanupTestData(user.ID, workspace.ID, "", "")

		// Create advance config with only rate limit
		advanceConfig := ProjectAdvanceConfig{
			RateLimit: &RateLimitConfig{
				Enabled:     true,
				RequestsRPM: 50,
				BurstSize:   5,
			},
		}

		configJSON, err := json.Marshal(advanceConfig)
		require.NoError(t, err)

		project, err := CreateTestProjectWithConfig(workspace.ID, "Test Project 6", "test-project-6", string(configJSON))
		require.NoError(t, err)

		// Verify only rate limit is set
		var storedConfig ProjectAdvanceConfig
		err = json.Unmarshal([]byte(project.AdvanceConfig), &storedConfig)
		require.NoError(t, err)

		assert.Equal(t, 0, storedConfig.GlobalTimeout) // Should be zero value
		require.NotNil(t, storedConfig.RateLimit)
		assert.True(t, storedConfig.RateLimit.Enabled)
		assert.Equal(t, 50, storedConfig.RateLimit.RequestsRPM)
		assert.Nil(t, storedConfig.CORS)     // Should be nil
		assert.Nil(t, storedConfig.Security) // Should be nil
	})
}

func TestProjectModelValidation(t *testing.T) {
	// Setup test database
	err := CheckAndHandle()
	require.NoError(t, err, "Failed to setup test database")

	t.Run("UUID Generation", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := CreateTestWorkspace("test7@example.com", "Test User 7", "Test Workspace 7")
		require.NoError(t, err)
		defer CleanupTestData(user.ID, workspace.ID, "", "")

		// Create project and verify UUID is generated
		project, err := CreateTestProject(workspace.ID, "Test Project 7", "test-project-7")
		require.NoError(t, err)

		// Verify UUID format
		_, err = uuid.Parse(project.ID)
		assert.NoError(t, err, "Project ID should be a valid UUID")
		assert.NotEmpty(t, project.CreatedAt)
		assert.NotEmpty(t, project.UpdatedAt)
	})

	t.Run("Unique Alias Constraint", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := CreateTestWorkspace("test8@example.com", "Test User 8", "Test Workspace 8")
		require.NoError(t, err)
		defer CleanupTestData(user.ID, workspace.ID, "", "")

		// Create first project
		_, err = CreateTestProject(workspace.ID, "Test Project 8A", "duplicate-alias")
		require.NoError(t, err)

		// Try to create second project with same alias - should fail
		_, err = CreateTestProject(workspace.ID, "Test Project 8B", "duplicate-alias")
		assert.Error(t, err, "Creating project with duplicate alias should fail")
	})

	t.Run("Default Values", func(t *testing.T) {
		// Create test workspace and user
		user, workspace, err := CreateTestWorkspace("test9@example.com", "Test User 9", "Test Workspace 9")
		require.NoError(t, err)
		defer CleanupTestData(user.ID, workspace.ID, "", "")

		// Create project and verify default values
		project, err := CreateTestProject(workspace.ID, "Test Project 9", "test-project-9")
		require.NoError(t, err)

		assert.Equal(t, ModeMock, project.Mode)
		assert.Equal(t, "running", project.Status)
		assert.Empty(t, project.AdvanceConfig)
		assert.Nil(t, project.ActiveProxyID)
	})
}
