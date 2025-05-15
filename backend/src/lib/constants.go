package lib

import (
	"os"
	"path/filepath"
)

// Path constants
var (
	// Get the current working directory
	CURRENT_DIR, _ = os.Getwd()

	// Derived paths
	CONFIGS_DIR = filepath.Join(CURRENT_DIR, "..", "configs")
	UPLOAD_DIR  = filepath.Join(CURRENT_DIR, "uploads")
)

// Server configuration
var (
	IS_DEBUG        = getEnvOrDefault("IS_DEBUG", "false")
	SERVER_PORT     = getEnvOrDefault("SERVER_PORT", "3600")
	SERVER_HOSTNAME = getEnvOrDefault("SERVER_HOSTNAME", "0.0.0.0")
	CORS_ORIGIN     = getEnvOrDefault("CORS_ORIGIN", "*")
	PROXY_MODE      = getEnvOrDefault("PROXY_MODE", "true") != "false"
	PROXY_BASE_URL  = getEnvOrDefault("PROXY_BASE_URL", "")
)

// Authentication
var (
	API_KEY = getEnvOrDefault("API_KEY", "admin:admin")
)

// Helper function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
