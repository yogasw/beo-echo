package utils

import (
	"errors"
	"os"
	"path/filepath"

	"mockoon-control-panel/backend_new/src/lib"
)

// EnsureRequiredFoldersAndEnv ensures that all required folders and environment variables exist
func EnsureRequiredFoldersAndEnv() error {
	// Create required directories
	directories := []string{
		lib.CONFIGS_DIR,
		lib.UPLOAD_DIR,
		lib.LOGS_DIR,
		filepath.Join(lib.CONFIGS_DIR, "traefik"),
	}

	for _, dir := range directories {
		if err := EnsureDirectoryExists(dir); err != nil {
			return errors.New("Failed to create directory " + dir + ": " + err.Error())
		}
	}

	// Check or create SQLite database directory
	dbDir := filepath.Join(lib.CONFIGS_DIR, "db")
	if err := EnsureDirectoryExists(dbDir); err != nil {
		return errors.New("Failed to create database directory: " + err.Error())
	}

	// Ensure database file path is set in environment
	if os.Getenv("DATABASE_URL") == "" {
		dbPath := "file:" + filepath.Join(dbDir, "db.sqlite")
		os.Setenv("DATABASE_URL", dbPath)
	}

	return nil
}
