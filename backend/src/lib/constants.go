package lib

import (
	"os"
	"path/filepath"
)

// Path constants
var (
	IS_TEST        = false
	IS_DESKTOP_APP = false // Flag to indicate if running as desktop app
	// Derived paths
	CONFIGS_DIR = filepath.Join(CURRENT_DIR(), "..", "configs")
	UPLOAD_DIR  = filepath.Join(CURRENT_DIR(), "uploads")
	CANDY_DIR   = filepath.Join(CONFIGS_DIR, "caddy")
	JWT_SECRET  = "" // ini from db or from env
)

// ResetPaths re-initializes all path constants based on current working directory
// This is useful when the working directory changes after package initialization
func ResetPaths() {
	if IS_DESKTOP_APP {
		// For desktop apps, use the app data directory directly
		appDataDir := CURRENT_DIR()
		CONFIGS_DIR = filepath.Join(appDataDir, "configs")
		UPLOAD_DIR = filepath.Join(appDataDir, "uploads")
		CANDY_DIR = filepath.Join(CONFIGS_DIR, "caddy")
	} else {
		// For regular server mode, use relative paths
		CONFIGS_DIR = filepath.Join(CURRENT_DIR(), "..", "configs")
		UPLOAD_DIR = filepath.Join(CURRENT_DIR(), "uploads")
		CANDY_DIR = filepath.Join(CONFIGS_DIR, "caddy")
	}
}

// SetDesktopMode sets the application to desktop mode and resets paths accordingly
func SetDesktopMode(isDesktop bool) {
	IS_DESKTOP_APP = isDesktop
	ResetPaths()
}

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
	} else if IS_DESKTOP_APP {
		// For desktop applications, always use user home directory
		// regardless of where the app is launched from (Applications folder, etc.)
		return getDesktopAppDataDir()
	} else {
		dir, _ := os.Getwd()
		return dir
	}
}

// getDesktopAppDataDir returns the appropriate app data directory for desktop apps
// based on the operating system
func getDesktopAppDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current working directory if home dir can't be determined
		dir, _ := os.Getwd()
		return dir
	}

	// Use .beoecho directory in user's home for all platforms
	appDataDir := filepath.Join(homeDir, ".beoecho")

	// Ensure the directory exists
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		// Fallback to current working directory if can't create app data dir
		dir, _ := os.Getwd()
		return dir
	}

	return appDataDir
}

// when env JWT_SECRET is not set, it will be generated and saved in the database
func GetJWTSecret() []byte {
	fromEnv := getEnvOrDefault("JWT_SECRET", "")
	if fromEnv == "" {
		return []byte(JWT_SECRET)
	} else {
		return []byte(fromEnv)
	}
}

func SetJWTSecret(secret string) {
	JWT_SECRET = secret
}
