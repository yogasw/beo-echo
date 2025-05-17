package lib

import (
	"os"
	"path/filepath"
)

// Path constants
var (
	IS_TEST = false
	// Derived paths
	CONFIGS_DIR = filepath.Join(CURRENT_DIR(), "..", "configs")
	UPLOAD_DIR  = filepath.Join(CURRENT_DIR(), "uploads")
	CANDY_DIR   = filepath.Join(CONFIGS_DIR, "caddy")
)

// Server configuration
var (
	IS_DEBUG        = getEnvOrDefault("IS_DEBUG", "false")
	SERVER_PORT     = getEnvOrDefault("SERVER_PORT", "3600")
	SERVER_HOSTNAME = getEnvOrDefault("SERVER_HOSTNAME", "0.0.0.0")
	CORS_ORIGIN     = getEnvOrDefault("CORS_ORIGIN", "*")
)

// Helper function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func CURRENT_DIR() string {
	if IS_TEST {
		return "/tmp/beo"
	} else {
		dir, _ := os.Getwd()
		return dir
	}
}
