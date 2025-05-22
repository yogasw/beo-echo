package handler

import (
	"log"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/mocks/repositories"
	"beo-echo/backend/src/mocks/services"
)

// Gin context keys
const (
	KeyProjectID     = "projectID"
	KeyExecutionMode = "executionMode"
	KeyMatched       = "matched"
	KeyPath          = "path"
)

var mockService *services.MockService

// InitMockService initializes the mock service
func InitMockService() {
	db := database.GetDB() // Get the database connection
	if db == nil {
		log.Println("Warning: Database connection not available for mock service")
		return
	}

	repo := repositories.NewMockRepository(db)
	mockService = services.NewMockService(repo)
}

// EnsureMockService ensures that the mock service is initialized
func EnsureMockService() {
	if mockService == nil {
		InitMockService()
	}
}

// GetProjectURL returns the URL for accessing a project's API
// It handles different URL formats based on PROXY_MODE configuration
func GetProjectURL(scheme, host string, project database.Project) string {
	// Direct access mode
	return scheme + "://" + host + "/" + project.Alias
}

// add validation for project alias
func IsValidAlias(alias string) bool {
	// Check if the alias contains only alphanumeric characters, dashes, and underscores
	for _, char := range alias {
		if !(char >= 'a' && char <= 'z') && !(char >= 'A' && char <= 'Z') && !(char >= '0' && char <= '9') && char != '-' && char != '_' {
			return false
		}
	}
	return true
}
