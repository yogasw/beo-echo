package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

/*
CreateProjectHandler creates a new project

Sample curl:

	curl -X POST "http://localhost:3600/api/api/projects" \
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
	if project.Name == "" || project.Alias == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project name and alias are required",
		})
		return
	}

	// project alias only allow alphanumeric characters, dashes and underscores
	if !handler.IsValidAlias(project.Alias) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project alias can only contain alphanumeric characters, dashes and underscores",
		})
		return
	}

	// Check if project alias already exists
	existingProject := &database.Project{}
	result := database.GetDB().Where("alias = ?", project.Alias).First(existingProject)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "Project alias already exists",
		})
		return
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to check project alias: " + result.Error.Error(),
		})
		return
	}

	// Default to mock mode if not specified
	if project.Mode == "" {
		project.Mode = database.ModeMock
	}

	// Validate advance config JSON format if provided
	if project.AdvanceConfig != "" {
		_, err := database.ParseProjectAdvanceConfig(project.AdvanceConfig)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
	}

	// Create the project
	result = database.GetDB().Create(&project)
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
