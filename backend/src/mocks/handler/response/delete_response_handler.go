package response

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// DeleteResponseHandler removes a response
//
// Sample curl:
// curl -X DELETE "http://localhost:8000/api/projects/my-project/endpoints/1/responses/1" -H "Content-Type: application/json"
func DeleteResponseHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectName := c.Param("projectName")
	if projectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project name is required",
		})
		return
	}

	// Parse endpoint ID and response ID
	endpointID, err := strconv.ParseUint(c.Param("endpointID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid endpoint ID",
		})
		return
	}

	responseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid response ID",
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

	// Check if response exists
	var response database.MockResponse
	result = database.GetDB().Where("id = ? AND endpoint_id = ?", responseID, endpoint.ID).First(&response)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Response not found",
		})
		return
	}

	// Delete the response (GORM will cascade delete related rules)
	result = database.GetDB().Delete(&response)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to delete response: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Response deleted successfully",
	})
}
