package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

/*
CreateProjectHandler creates a new project

Sample curl:

	curl -X POST "http://localhost:3600/mock/api/projects" \
		-H "Content-Type: application/json" \
		-d '{
		 "name": "my-new-project",
		 "mode": "mock"
		}'
*/
func CreateProjectHandler(c *gin.Context) {
	handler.EnsureMockService()

	var project database.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Validate project data
	if project.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project name is required",
		})
		return
	}

	// Default to mock mode if not specified
	if project.Mode == "" {
		project.Mode = database.ModeMock
	}

	// Create the project
	result := database.GetDB().Create(&project)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create project: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Project created successfully",
		"data":    project,
	})
}
