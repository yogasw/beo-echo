package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"mockoon-control-panel/backend_new/src/prisma"
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
)

// DefaultVariables contains default values for system configuration
var DefaultVariables = map[string]string{
	BaseURL:   "",
	GitURL:    "",
	GitBranch: "main",
	SSHKey:    "",
	GitName:   "Mockoon Control Panel",
	GitEmail:  "noreply@example.com",
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
	var config prisma.SystemConfig
	result := prisma.DB.Where("key = ?", keyName).First(&config)
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
	var config prisma.SystemConfig
	result := prisma.DB.Where("key = ?", keyName).First(&config)

	if result.Error == nil {
		// Update existing config
		config.Value = stringValue
		config.Type = keyType
		if err := prisma.DB.Save(&config).Error; err != nil {
			return fmt.Errorf("failed to update system config: %w", err)
		}
	} else {
		// Create new config
		config = prisma.SystemConfig{
			Key:   keyName,
			Value: stringValue,
			Type:  keyType,
		}
		if err := prisma.DB.Create(&config).Error; err != nil {
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
func GetAllSystemConfigs() ([]prisma.SystemConfig, error) {
	var configs []prisma.SystemConfig
	if err := prisma.DB.Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch system configs: %w", err)
	}
	return configs, nil
}

// SetConfigByID updates a system configuration by its ID
func SetConfigByID(id int, key, value, description string) (*prisma.SystemConfig, error) {
	var config prisma.SystemConfig
	if err := prisma.DB.First(&config, id).Error; err != nil {
		return nil, fmt.Errorf("config with ID %d not found: %w", id, err)
	}

	config.Key = key
	config.Value = value
	config.Description = description

	if err := prisma.DB.Save(&config).Error; err != nil {
		return nil, fmt.Errorf("failed to update config: %w", err)
	}

	return &config, nil
}

// AddConfig creates a new system configuration
func AddConfig(key, value, description, valueType string) (*prisma.SystemConfig, error) {
	config := prisma.SystemConfig{
		Key:         key,
		Value:       value,
		Description: description,
		Type:        valueType,
	}

	if err := prisma.DB.Create(&config).Error; err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	return &config, nil
}
