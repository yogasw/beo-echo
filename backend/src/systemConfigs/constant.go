package systemConfig

// System configuration keys
const (
	// Custom Subdomain configuration
	CUSTOM_SUBDOMAIN_ENABLED = "CUSTOM_SUBDOMAIN_ENABLED" // Controls whether custom subdomains are enabled
	CUSTOM_SUBDOMAIN_DOMAIN  = "CUSTOM_SUBDOMAIN_DOMAIN"  // Base domain for custom subdomains (e.g., "beo-echo.com")

	// Feature flags
	FEATURE_SHOW_PASSWORD_REQUIREMENTS = "FEATURE_SHOW_PASSWORD_REQUIREMENTS"
	FEATURE_EMAIL_UPDATES_ENABLED      = "FEATURE_EMAIL_UPDATES_ENABLED"
	FEATURE_REGISTER_EMAIL_ENABLED     = "FEATURE_REGISTER_EMAIL_ENABLED"
	FEATURE_OAUTH_AUTO_REGISTER        = "FEATURE_OAUTH_AUTO_REGISTER" // Controls whether new users can register through OAuth

	AUTO_SAVE_LOGS_IN_DB_ENABLED = "AUTO_SAVE_LOGS_IN_DB_ENABLED" // Enable auto-saving of logs

	// Workspace and Project Limits
	AUTO_CREATE_WORKSPACE_ON_REGISTER = "AUTO_CREATE_WORKSPACE_ON_REGISTER" // Automatically create a workspace for new users
	MAX_USER_WORKSPACES               = "MAX_USER_WORKSPACES"               // Maximum number of workspaces a user can create
	MAX_WORKSPACE_PROJECTS            = "MAX_WORKSPACE_PROJECTS"            // Maximum number of projects allowed in a workspace

	JWT_SECRET = "JWT_SECRET" // JWT secret for signing tokens

	// AI Configuration
	AI_PROVIDER     = "AI_PROVIDER"     // AI provider (gemini, openai, custom)
	AI_API_KEY      = "AI_API_KEY"      // API key for AI service
	AI_API_ENDPOINT = "AI_API_ENDPOINT" // API endpoint for AI service
	AI_MODEL        = "AI_MODEL"        // AI model to use for generation

	// Default Response Configuration
	DEFAULT_RESPONSE_PROJECT_NOT_FOUND      = "This is a default response from Beo Echo mock service." // Default response when project is not found
	DEFAULT_RESPONSE_ENDPOINT_NOT_FOUND     = "This is a default response from Beo Echo mock service." // Default response when endpoint is not found
	DEFAULT_RESPONSE_NO_RESPONSE_CONFIGURED = "This is a default response from Beo Echo mock service." // Default response when no response is configured

	// Landing Page Configuration
	LANDING_PAGE_ENABLED = "LANDING_PAGE_ENABLED" // Enable/disable landing page
	MOCK_URL_FORMAT      = "MOCK_URL_FORMAT"      // URL format: "subdomain" or "path"
)

// DefaultConfigSettings contains all system configuration settings with metadata
var DefaultConfigSettings = map[SystemConfigKey]ConfigSetting{
	// Custom Subdomain Configuration
	CUSTOM_SUBDOMAIN_ENABLED: {
		Type:        TypeBoolean,
		Value:       "false",
		Description: "Enable custom subdomains for projects (security implications) - required restart service to apply",
		Category:    "Domains",
	},
	CUSTOM_SUBDOMAIN_DOMAIN: {
		Type:        TypeString,
		Value:       "*.beo-echo.com",
		Description: "Base domain for custom subdomains - required restart service to apply",
		Category:    "Domains",
	},

	// Feature Flags
	FEATURE_SHOW_PASSWORD_REQUIREMENTS: {
		Type:        TypeBoolean,
		Value:       "true",
		Description: "Show password requirements during account creation",
		Category:    "Features",
	},
	FEATURE_EMAIL_UPDATES_ENABLED: {
		Type:        TypeBoolean,
		Value:       "true",
		Description: "Enable email notification updates",
		Category:    "Features",
	},
	FEATURE_REGISTER_EMAIL_ENABLED: {
		Type:        TypeBoolean,
		Value:       "true",
		Description: "Enable email registration confirmation",
		Category:    "Features",
	},
	FEATURE_OAUTH_AUTO_REGISTER: {
		Type:        TypeBoolean,
		Value:       "true",
		Description: "Allow new users to register through OAuth providers. If disabled, only existing users can login",
		Category:    "Features",
	},
	AUTO_SAVE_LOGS_IN_DB_ENABLED: {
		Type:        TypeBoolean,
		Value:       "false",
		Description: "Automatically persist request logs to database (may affect performance)",
		Category:    "Logging",
	},

	// Workspace and Project Limits
	AUTO_CREATE_WORKSPACE_ON_REGISTER: {
		Type:        TypeBoolean,
		Value:       "false",
		Description: "Automatically create a workspace for new users upon registration",
		Category:    "Limits",
	},
	MAX_USER_WORKSPACES: {
		Type:        TypeNumber,
		Value:       "2",
		Description: "Maximum number of workspaces a user can create",
		Category:    "Limits",
	},
	MAX_WORKSPACE_PROJECTS: {
		Type:        TypeNumber,
		Value:       "100",
		Description: "Maximum number of projects allowed in a workspace",
		Category:    "Limits",
	},
	JWT_SECRET: {
		Type:        TypeString,
		Value:       "",
		Description: "JWT secret for signing tokens",
		Category:    "Security",
		HideValue:   true,
	},

	// AI Configuration
	AI_PROVIDER: {
		Type:        TypeString,
		Value:       "gemini",
		Description: "AI provider type: gemini, openai, or custom",
		Category:    "AI",
	},
	AI_API_KEY: {
		Type:        TypeString,
		Value:       "",
		Description: "API key for AI service (Gemini, OpenAI, Claude, etc.)",
		Category:    "AI",
		HideValue:   true,
	},
	AI_API_ENDPOINT: {
		Type:        TypeString,
		Value:       "https://generativelanguage.googleapis.com/v1beta",
		Description: "API endpoint for AI service (Gemini free tier default)",
		Category:    "AI",
	},
	AI_MODEL: {
		Type:        TypeString,
		Value:       "gemini-pro",
		Description: "AI model to use for generation (gemini-pro is free)",
		Category:    "AI",
	},

	// Landing Page Configuration
	LANDING_PAGE_ENABLED: {
		Type:        TypeBoolean,
		Value:       "false",
		Description: "Enable or disable the landing page display",
		Category:    "Landing Page",
	},
	MOCK_URL_FORMAT: {
		Type:        TypeString,
		Value:       "subdomain",
		Description: "Mock URL format: 'subdomain' (alias.domain.com) or 'path' (domain.com/alias)",
		Category:    "Landing Page",
	},
}

// ConfigType represents the data type of a configuration value
type ConfigType string
type SystemConfigKey string

func (k SystemConfigKey) String() string {
	return string(k)
}
func (k ConfigType) String() string {
	return string(k)
}

const (
	TypeString  ConfigType = "string"
	TypeBoolean ConfigType = "boolean"
	TypeNumber  ConfigType = "number"
	TypeArray   ConfigType = "array"
)

// ConfigSetting represents a single configuration setting with metadata
type ConfigSetting struct {
	Key         SystemConfigKey `json:"key"`         // Unique key for the configuration
	Type        ConfigType      `json:"type"`        // Data type of the configuration
	Value       string          `json:"value"`       // Default value
	Description string          `json:"description"` // Human-readable description
	HideValue   bool            `json:"hide_value"`  // Whether to hide value in UI (for sensitive data)
	Category    string          `json:"category"`    // Category for grouping related settings
}
