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

	// Default Response Configuration
	DEFAULT_RESPONSE_PROJECT_NOT_FOUND      = "This is a default response from Beo Echo mock service."       // Default response when project is not found
	DEFAULT_RESPONSE_ENDPOINT_NOT_FOUND     = "This is a default response from Beo Echo mock service."       // Default response when endpoint is not found
	DEFAULT_RESPONSE_NO_RESPONSE_CONFIGURED = "This endpoint exists but no specific response is configured." // Default response when no response is configured
)

// DefaultConfigSettings contains all system configuration settings with metadata
var DefaultConfigSettings = map[SystemConfigKey]ConfigSetting{
	// Custom Subdomain Configuration
	CUSTOM_SUBDOMAIN_ENABLED: {
		Type:        TypeBoolean,
		Value:       "false",
		Description: "Enable custom subdomains for projects (security implications)",
		Category:    "Domains",
	},
	CUSTOM_SUBDOMAIN_DOMAIN: {
		Type:        TypeString,
		Value:       "*.beo-echo.com",
		Description: "Base domain for custom subdomains",
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
		HideValue:   true,
	},
	MAX_WORKSPACE_PROJECTS: {
		Type:        TypeNumber,
		Value:       "100",
		Description: "Maximum number of projects allowed in a workspace",
		Category:    "Limits",
		HideValue:   true,
	},
	JWT_SECRET: {
		Type:        TypeString,
		Value:       "",
		Description: "JWT secret for signing tokens",
		Category:    "Security",
		HideValue:   true,
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
	Key         SystemConfigKey // Unique key for the configuration
	Type        ConfigType      // Data type of the configuration
	Value       string          // Default value
	Description string          // Human-readable description
	HideValue   bool            // Whether to hide value in UI (for sensitive data)
	Category    string          // Category for grouping related settings
}
