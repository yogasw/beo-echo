package handler

import (
	"log"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/lib"
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
func GetProjectURL(project database.Project) string {
	// Import configuration values from lib package
	if lib.PROXY_MODE {
		if lib.PROXY_BASE_URL != "" {
			return lib.PROXY_BASE_URL + "/" + project.Name
		} else {
			// When in PROXY_MODE but no PROXY_BASE_URL specified, we can't determine
			// the full URL here since we don't have access to HTTP request info
			// We'll return a placeholder that will be updated by the frontend
			return project.Name
		}
	} else {
		// Direct access mode
		return "http://" + lib.SERVER_HOSTNAME + ":" + lib.SERVER_PORT + "/" + project.Name
	}
}
