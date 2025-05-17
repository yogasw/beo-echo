package utils

import (
	"os"
	"path/filepath"
	"testing"

	"beo-echo/backend/src/lib"

	"github.com/stretchr/testify/assert"
)

func TestEnsureRequiredFoldersAndEnv(t *testing.T) {
	// Setup temporary directories for testing
	tempDir := t.TempDir()
	lib.CONFIGS_DIR = filepath.Join(tempDir, "configs")
	lib.UPLOAD_DIR = filepath.Join(tempDir, "uploads")
	lib.CANDY_DIR = filepath.Join(tempDir, "candy")

	// Call the function
	err := EnsureRequiredFoldersAndEnv()
	assert.NoError(t, err, "EnsureRequiredFoldersAndEnv should not return an error")

	// Verify directories are created
	assert.DirExists(t, lib.CONFIGS_DIR, "CONFIGS_DIR should exist")
	assert.DirExists(t, lib.UPLOAD_DIR, "UPLOAD_DIR should exist")
	assert.DirExists(t, lib.CANDY_DIR, "CANDY_DIR should exist")

	// Verify database directory and environment variable
	dbDir := filepath.Join(lib.CONFIGS_DIR, "db")
	assert.DirExists(t, dbDir, "Database directory should exist")
	dbPath := "file:" + filepath.Join(dbDir, "db.sqlite")
	assert.Equal(t, dbPath, os.Getenv("DATABASE_URL"), "DATABASE_URL should be set correctly")
}

func TestEnsureDirectoryExists(t *testing.T) {
	tempDir := t.TempDir()
	dirPath := filepath.Join(tempDir, "testDir")

	err := EnsureDirectoryExists(dirPath)
	assert.NoError(t, err, "EnsureDirectoryExists should not return an error")
	assert.DirExists(t, dirPath, "Directory should exist")
}
