package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	// Test SQLite connection
	t.Run("SQLite Connection", func(t *testing.T) {
		// Clear any existing DATABASE_URL
		os.Unsetenv("DATABASE_URL")

		// Create necessary directory structure for testing
		currentDir, _ := os.Getwd()
		configsDir := filepath.Join(currentDir, "..", "configs")
		dbDir := filepath.Join(configsDir, "db")

		// Create the directories if they don't exist
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			t.Fatalf("Failed to create database directory: %v", err)
		}

		t.Cleanup(func() {
			// Clean up the test database file after test
			dbPath := filepath.Join(dbDir, "db.sqlite")
			os.Remove(dbPath)
		})

		// Initialize database
		err := CheckAndHandle()
		if err != nil {
			t.Fatalf("Failed to connect to SQLite database: %v", err)
		}

		// Verify connection works
		db := GetDB()
		if db == nil {
			t.Fatal("Database connection is nil")
		}

		sqlDB, err := db.DB()
		if err != nil {
			t.Fatalf("Failed to get SQL DB: %v", err)
		}

		err = sqlDB.Ping()
		if err != nil {
			t.Fatalf("Failed to ping SQLite database: %v", err)
		}

		t.Log("Successfully connected to SQLite database")
	})

	// Skip PostgreSQL test if no environment variable is set
	// To run this test, set the DATABASE_URL environment variable to a valid PostgreSQL connection string
	t.Run("PostgreSQL Connection", func(t *testing.T) {
		postgresURL := os.Getenv("TEST_POSTGRES_URL")
		if postgresURL == "" {
			t.Skip("Skipping PostgreSQL test - TEST_POSTGRES_URL environment variable not set")
		}

		// Set DATABASE_URL for PostgreSQL
		os.Setenv("DATABASE_URL", postgresURL)

		// Initialize database
		err := CheckAndHandle()
		if err != nil {
			t.Fatalf("Failed to connect to PostgreSQL database: %v", err)
		}

		// Verify connection works
		db := GetDB()
		if db == nil {
			t.Fatal("Database connection is nil")
		}

		sqlDB, err := db.DB()
		if err != nil {
			t.Fatalf("Failed to get SQL DB: %v", err)
		}

		err = sqlDB.Ping()
		if err != nil {
			t.Fatalf("Failed to ping PostgreSQL database: %v", err)
		}

		t.Log("Successfully connected to PostgreSQL database")

		// Reset DATABASE_URL
		os.Unsetenv("DATABASE_URL")
	})
}
