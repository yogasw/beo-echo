package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"beo-echo/backend/src/database"
)

// System configuration keys
const (
	// Base config
	BaseURL = "BASE_URL:string"

	// Git sync
	GitURL    = "GIT_URL:string"
	GitBranch = "GIT_BRANCH:string"
	SSHKey    = "SSH_KEY:string"
	GitName   = "GIT_NAME:string"
	GitEmail  = "GIT_EMAIL:string"

	FeatureShowPasswordRequirements = "FEATURE_SHOW_PASSWORD_REQUIREMENTS:boolean"
	FeatureEmailUpdatesEnabled      = "FEATURE_EMAIL_UPDATES_ENABLED:boolean"
	FEATURE_REGISTER_EMAIL_ENABLED  = "FEATURE_REGISTER_EMAIL_ENABLED:boolean"
)

// DefaultVariables contains default values for system configuration
var DefaultVariables = map[string]string{
	BaseURL:   "",
	GitURL:    "",
	GitBranch: "main",
	SSHKey:    "",
	GitName:   "BeoEcho",
	GitEmail:  "noreply@example.com",

	// Default values for feature flags
	FeatureShowPasswordRequirements: "true",
	FeatureEmailUpdatesEnabled:      "true",
	FEATURE_REGISTER_EMAIL_ENABLED:  "false",
}

// GetSystemConfig retrieves a system configuration value from the database with type conversion
func GetSystemConfig(key string) (interface{}, error) {
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid key format, expected key:type")
	}

	keyName := parts[0]
	keyType := parts[1]

	// Try to get from database first
	var config database.SystemConfig
	result := database.DB.Where("key = ?", keyName).First(&config)
	if result.Error == nil {
		// Return with proper type conversion
		return convertValue(config.Value, config.Type)
	}

	// If not found in database, check defaults
	if defaultVal, exists := DefaultVariables[key]; exists {
		return convertValue(defaultVal, keyType)
	}

	return nil, fmt.Errorf("configuration key %s not found", key)
}

// GetConfig retrieves a system configuration value with automatic type conversion to T
// T can be string, bool, float64, or []string
func GetConfig[T any](key string) (T, error) {
	var empty T

	// Add the type suffix if not already present
	if !strings.Contains(key, ":") {
		// Determine type suffix based on T
		switch any(empty).(type) {
		case string:
			key += ":string"
		case bool:
			key += ":boolean"
		case float64:
			key += ":number"
		case []string:
			key += ":array"
		default:
			return empty, fmt.Errorf("unsupported type for key %s", key)
		}
	}

	// Get the config using the original function
	value, err := GetSystemConfig(key)
	if err != nil {
		return empty, err
	}

	// Type assert to the requested type
	result, ok := value.(T)
	if !ok {
		return empty, fmt.Errorf("unable to convert value to requested type for key %s", key)
	}

	return result, nil
}

// SetSystemConfig sets a system configuration value in the database with type validation
func SetSystemConfig(key string, value interface{}) error {
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		return errors.New("invalid key format, expected key:type")
	}

	keyName := parts[0]
	keyType := parts[1]

	// Convert value to string based on type
	var stringValue string
	switch keyType {
	case "array":
		bytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal array value: %w", err)
		}
		stringValue = string(bytes)
	default:
		stringValue = fmt.Sprintf("%v", value)
	}

	// Check if the key exists
	var config database.SystemConfig
	result := database.DB.Where("key = ?", keyName).First(&config)

	if result.Error == nil {
		// Update existing config
		config.Value = stringValue
		config.Type = keyType
		if err := database.DB.Save(&config).Error; err != nil {
			return fmt.Errorf("failed to update system config: %w", err)
		}
	} else {
		// Create new config
		config = database.SystemConfig{
			Key:   keyName,
			Value: stringValue,
			Type:  keyType,
		}
		if err := database.DB.Create(&config).Error; err != nil {
			return fmt.Errorf("failed to create system config: %w", err)
		}
	}

	return nil
}

// convertValue converts a string value to the appropriate type based on type string
func convertValue(value string, valueType string) (interface{}, error) {
	switch valueType {
	case "number":
		return strconv.ParseFloat(value, 64)
	case "boolean":
		return strings.ToLower(value) == "true", nil
	case "array":
		var result []string
		err := json.Unmarshal([]byte(value), &result)
		if err != nil {
			log.Printf("Error parsing array value: %v", err)
			return []string{}, nil
		}
		return result, nil
	default:
		return value, nil
	}
}

// GetAllSystemConfigs retrieves all system configurations from the database
func GetAllSystemConfigs() ([]database.SystemConfig, error) {
	var configs []database.SystemConfig
	if err := database.DB.Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch system configs: %w", err)
	}
	return configs, nil
}

// GetFeatureFlags retrieves all feature flags from the system configuration
func GetFeatureFlags() (map[string]bool, error) {
	featureFlags := make(map[string]bool)

	// Get all configs
	configs, err := GetAllSystemConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feature flags: %w", err)
	}

	// Filter out feature flags
	for _, config := range configs {
		if strings.HasPrefix(strings.ToLower(config.Key), "feature_") ||
			strings.HasPrefix(config.Key, "FEATURE_") {
			// Convert the value to boolean
			enabled, _ := strconv.ParseBool(config.Value)
			featureFlags[config.Key] = enabled
		}
	}

	// Add default feature flags that don't exist in the database
	for key, value := range DefaultVariables {
		// skip if already exists
		if strings.HasPrefix(key, "FEATURE_") && !featureFlags[key] {
			enabled, _ := strconv.ParseBool(value)
			featureFlags[key] = enabled
		}
	}

	// remove type suffix from keys
	for key := range featureFlags {
		if strings.Contains(key, ":") {
			parts := strings.Split(key, ":")
			featureFlags[parts[0]] = featureFlags[key]
			delete(featureFlags, key)
		}
	}

	return featureFlags, nil
}

// SetConfigByID updates a system configuration by its ID
func SetConfigByID(id, key, value, description string) (*database.SystemConfig, error) {
	var config database.SystemConfig
	if err := database.DB.Where("id = ?", id).First(&config).Error; err != nil {
		return nil, fmt.Errorf("config with ID %v not found: %w", id, err)
	}

	config.Key = key
	config.Value = value
	config.Description = description

	if err := database.DB.Save(&config).Error; err != nil {
		return nil, fmt.Errorf("failed to update config: %w", err)
	}

	return &config, nil
}

// AddConfig creates a new system configuration
func AddConfig(key, value, description, valueType string) (*database.SystemConfig, error) {
	config := database.SystemConfig{
		Key:         key,
		Value:       value,
		Description: description,
		Type:        valueType,
	}

	if err := database.DB.Create(&config).Error; err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	return &config, nil
}
