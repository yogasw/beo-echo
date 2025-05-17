package systemConfig

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"beo-echo/backend/src/database"
)

// GetSystemConfig retrieves a system configuration value from the database with type conversion
func GetSystemConfig(key string) (interface{}, error) {

	// Try to get from database first
	var config database.SystemConfig
	result := database.DB.Where("key = ?", key).First(&config)
	if result.Error == nil {
		// Return with proper type conversion
		return convertValue(config.Value, config.Type)
	}

	// If not found in database, check defaults
	defaultValue, exists := DefaultConfigSettings[SystemConfigKey(key)]
	if exists {
		// Convert default value to the requested type
		convertedValue, err := convertValue(defaultValue.Value, string(defaultValue.Type))
		if err != nil {
			return nil, fmt.Errorf("failed to convert default value: %w", err)
		}
		return convertedValue, nil
	}

	return nil, fmt.Errorf("configuration key %s not found", key)
}

// GetSystemConfigWithType retrieves a system configuration value with automatic type conversion to T
// T can be string, bool, float64, or []string
func GetSystemConfigWithType[T any](key string) (T, error) {
	var empty T

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
func SetSystemConfig(key, value string) error {
	// Check if the key exists
	var config database.SystemConfig
	result := database.DB.Where("key = ?", key).First(&config)

	if result.Error == nil {
		// Update existing config
		config.Value = value
		if err := database.DB.Save(&config).Error; err != nil {
			return fmt.Errorf("failed to update system config: %w", err)
		}
	} else {
		var config database.SystemConfig
		// Create new value
		defaultValue, exists := DefaultConfigSettings[SystemConfigKey(key)]
		if !exists {
			config = database.SystemConfig{
				Key:   key,
				Value: value,
			}
		} else {
			config = database.SystemConfig{
				Key:         key,
				Value:       value,
				Type:        string(defaultValue.Type),
				Description: defaultValue.Description,
				HideValue:   defaultValue.HideValue,
			}
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

// GetAllConfigSettings retrieves all configuration settings with metadata
// This combines database values with default settings metadata
func GetAllConfigSettings() ([]ConfigSetting, error) {
	settingsMap := DefaultConfigSettings

	// Get all configs from database
	dbConfigs, err := GetAllSystemConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch configs: %w", err)
	}

	// Update with database values
	for _, dbConfig := range dbConfigs {
		keyName := SystemConfigKey(dbConfig.Key)

		if setting, exists := settingsMap[keyName]; exists {
			// Update existing setting with DB value
			setting.Value = dbConfig.Value
			setting.Description = dbConfig.Description
			setting.HideValue = dbConfig.HideValue
			settingsMap[keyName] = setting
		} else {
			// Add new setting from database
			configType := ConfigType(dbConfig.Type)
			if configType == "" {
				configType = TypeString
			}

			// Create new setting and add to result
			newSetting := ConfigSetting{
				Key:         keyName,
				Type:        configType,
				Value:       dbConfig.Value,
				Description: dbConfig.Description,
				HideValue:   dbConfig.HideValue,
				Category:    "Custom", // Default category for DB-only settings
			}
			settingsMap[keyName] = newSetting
		}
	}
	result := []ConfigSetting{}
	// Convert map to slice
	for _, setting := range settingsMap {
		result = append(result, setting)
	}

	return result, nil
}

// GetConfigSetting retrieves a specific configuration setting with metadata
func GetConfigSetting(key string) (*ConfigSetting, error) {
	keyName := SystemConfigKey(key)
	for _, setting := range DefaultConfigSettings {
		if setting.Key == keyName {
			// Create a copy to avoid modifying the default
			result := setting

			// Check if there's a database override
			var dbConfig database.SystemConfig
			if err := database.DB.Where("key = ?", keyName).First(&dbConfig).Error; err == nil {
				// Update with database values
				result.Value = dbConfig.Value
				result.Description = dbConfig.Description
				result.HideValue = dbConfig.HideValue
			}

			return &result, nil
		}
	}

	// If not found in defaults, check database
	var dbConfig database.SystemConfig
	if err := database.DB.Where("key = ?", keyName).First(&dbConfig).Error; err == nil {
		configType := ConfigType(dbConfig.Type)
		if configType == "" {
			configType = TypeString
		}

		return &ConfigSetting{
			Key:         keyName,
			Type:        configType,
			Value:       dbConfig.Value,
			Description: dbConfig.Description,
			HideValue:   dbConfig.HideValue,
			Category:    "Custom", // Default category for DB-only settings
		}, nil
	}

	return nil, fmt.Errorf("configuration key %s not found", key)
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
		if strings.HasPrefix(config.Key, "FEATURE_") {
			// Convert the value to boolean
			enabled, _ := strconv.ParseBool(config.Value)
			featureFlags[config.Key] = enabled
		}
	}

	// Add default feature flags that don't exist in the database
	for key, v := range DefaultConfigSettings {
		// skip if already exists
		if strings.HasPrefix(string(key), "FEATURE_") && !featureFlags[string(key)] {
			enabled, _ := strconv.ParseBool(v.Value)
			featureFlags[string(key)] = enabled
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
func AddConfig(key, value, description string, valueType ConfigType, hideValue bool) (*database.SystemConfig, error) {

	config := database.SystemConfig{
		Key:         key,
		Value:       value,
		Description: description,
		Type:        valueType.String(),
		HideValue:   hideValue,
	}

	if err := database.DB.Create(&config).Error; err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	return &config, nil
}
