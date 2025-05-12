package database

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"mockoon-control-panel/backend_new/src/lib"
)

var DB *gorm.DB

// CheckAndHandlePrisma initializes the database connection
func CheckAndHandlePrisma() error {
	// Get the database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Set default database URL if not provided
		dbDir := filepath.Join(lib.CONFIGS_DIR, "db")
		dbPath := filepath.Join(dbDir, "db.sqlite")
		dbURL = "file:" + dbPath
		os.Setenv("DATABASE_URL", dbURL)
	}

	// Remove "file:" prefix for GORM SQLite
	sqlitePath := strings.TrimPrefix(dbURL, "file:")

	// Open the database connection
	var err error
	DB, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return errors.New("Failed to connect to the database: " + err.Error())
	}

	// Auto migrate the models
	if err := DB.AutoMigrate(
		&SystemConfig{},
		&Project{},
		&ProxyTarget{},
		&MockEndpoint{},
		&MockResponse{},
		&MockRule{},
		&RequestLog{},
		&User{},
		&UserIdentity{},
		&Workspace{},
		&UserWorkspace{},
		&SSOConfig{},
	); err != nil {
		return errors.New("Failed to migrate database schema: " + err.Error())
	}

	log.Println("Database connection established and migrations completed")
	return nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}
