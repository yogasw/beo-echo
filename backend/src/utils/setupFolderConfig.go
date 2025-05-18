package utils

import (
	"errors"
	"os"
	"path/filepath"

	"beo-echo/backend/src/lib"
)

// EnsureRequiredFoldersAndEnv ensures that all required folders and environment variables exist
func EnsureRequiredFoldersAndEnv() error {
	// Create required directories
	directories := []string{
		lib.CONFIGS_DIR,
		lib.UPLOAD_DIR,
		lib.CANDY_DIR,
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

// EnsureDirectoryExists creates a directory if it doesn't exist
func EnsureDirectoryExists(dirPath string) error {
	return os.MkdirAll(dirPath, os.ModePerm)
}

func SetupFolderConfigForTest() {
	lib.IS_TEST = true
	EnsureRequiredFoldersAndEnv()
}

func CleanupTestFolders() {
	os.RemoveAll(lib.CONFIGS_DIR)
	os.RemoveAll(lib.UPLOAD_DIR)
}
