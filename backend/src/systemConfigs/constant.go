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
	AUTO_SAVE_LOGS_IN_DB_ENABLED       = "AUTO_SAVE_LOGS_IN_DB_ENABLED" // Enable auto-saving of logs
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
		Value:       "false",
		Description: "Enable email registration confirmation",
		Category:    "Features",
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
