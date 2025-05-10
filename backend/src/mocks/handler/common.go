package handler

import (
	"log"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/repositories"
	"mockoon-control-panel/backend_new/src/mocks/services"
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
func GetProjectURL(host string, project database.Project) string {
	// Direct access mode
	return "http://" + host + "/" + project.Alias
}
