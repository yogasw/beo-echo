package response

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// ListResponsesHandler lists all responses for an endpoint
//
// Sample curl:
// curl -X GET "http://localhost:8000/api/projects/my-project/endpoints/1/responses" -H "Content-Type: application/json"
func ListResponsesHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectName := c.Param("projectName")
	if projectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project name is required",
		})
		return
	}

	// Parse endpoint ID
	endpointID, err := strconv.ParseUint(c.Param("endpointID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid endpoint ID",
		})
		return
	}

	// Find project first
	var project database.Project
	result := database.GetDB().Where("name = ?", projectName).First(&project)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not found",
		})
		return
	}

	// Check if endpoint exists and belongs to this project
	var endpoint database.MockEndpoint
	result = database.GetDB().Where("id = ? AND project_id = ?", endpointID, project.ID).First(&endpoint)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Endpoint not found",
		})
		return
	}

	// Get responses for this endpoint
	var responses []database.MockResponse
	result = database.GetDB().
		Preload("Rules").
		Where("endpoint_id = ?", endpoint.ID).
		Order("priority").
		Find(&responses)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve responses: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responses,
	})
}
