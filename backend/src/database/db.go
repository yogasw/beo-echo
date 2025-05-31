package database

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"beo-echo/backend/src/lib"
)

var DB *gorm.DB

// CheckAndHandle initializes the database connection
func CheckAndHandle() error {
	// Get the database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	var err error

	// Set default log level to Silent, but allow overriding via DB_LOG_LEVEL env var
	// Valid values: "silent", "error", "warn", "info"
	logLevelStr := strings.ToLower(os.Getenv("DB_LOG_LEVEL"))
	logLevel := logger.Silent // Default to silent

	// Configure log level based on environment variable if set
	switch logLevelStr {
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	}

	// Force silent logging during tests to suppress SQL query logs
	if os.Getenv("GO_TEST") == "1" {
		logLevel = logger.Silent
	}

	// Check if the database URL is empty or contains "sqlite"
	if dbURL == "" || strings.Contains(strings.ToLower(dbURL), "sqlite") {
		// If no URL is provided or if it contains "sqlite", use SQLite
		if dbURL == "" {
			// Set default database URL if not provided
			dbDir := filepath.Join(lib.CONFIGS_DIR, "db")
			dbPath := filepath.Join(dbDir, "db.sqlite")

			// Create database directory if it doesn't exist
			if _, err := os.Stat(dbDir); os.IsNotExist(err) {
				log.Println("Database directory doesn't exist, creating:", dbDir)
				if err := os.MkdirAll(dbDir, 0755); err != nil {
					return errors.New("Failed to create database directory: " + err.Error())
				}
			}

			dbURL = "file:" + dbPath
			os.Setenv("DATABASE_URL", dbURL)
		}

		log.Println("Using SQLite database:", dbURL)

		// Process SQLite connection string based on format
		var sqlitePath string

		// Handle different SQLite URL formats
		if strings.HasPrefix(dbURL, "sqlite:") {
			sqlitePath = strings.TrimPrefix(dbURL, "sqlite:")
		} else if strings.HasPrefix(dbURL, "file:") {
			sqlitePath = strings.TrimPrefix(dbURL, "file:")
		} else {
			sqlitePath = dbURL
		}

		// Ensure the directory for the SQLite file exists
		dbDir := filepath.Dir(sqlitePath)
		if _, err := os.Stat(dbDir); os.IsNotExist(err) {
			log.Println("SQLite directory doesn't exist, creating:", dbDir)
			if err := os.MkdirAll(dbDir, 0755); err != nil {
				return errors.New("Failed to create SQLite database directory: " + err.Error())
			}
		}

		// Open SQLite database connection
		DB, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
			// Disable SQL logs in tests
			DisableForeignKeyConstraintWhenMigrating: os.Getenv("GO_TEST") == "1",
		})
	} else {
		// Use PostgreSQL for all other database URLs
		log.Println("Using PostgreSQL database")
		DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
			// Disable SQL logs in tests
			DisableForeignKeyConstraintWhenMigrating: os.Getenv("GO_TEST") == "1",
		})
	}
	if err != nil {
		return errors.New("Failed to connect to the database: " + err.Error())
	}

	DB.AutoMigrate(&ProxyTarget{})
	DB.AutoMigrate(&Project{})
	DB.AutoMigrate(&Replay{})

	// Auto migrate the models
	if err := DB.AutoMigrate(
		&SystemConfig{},
		&MockEndpoint{},
		&MockResponse{},
		&MockRule{},
		&RequestLog{},
		&User{},
		&UserIdentity{},
		&Workspace{},
		&UserWorkspace{},
		&SSOConfig{},
		&ReplayFolder{},
	); err != nil {
		return errors.New("Failed to migrate database schema: " + err.Error())
	}

	log.Println("Database connection established and migrations completed")

	// Initialize default user and workspace if no users exist
	if err := InitializeDefaultUserAndWorkspace(DB); err != nil {
		log.Printf("Warning: Failed to initialize default user: %v", err)
	}

	return nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}
