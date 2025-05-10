package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// ListResponsesHandler lists all responses for an endpoint
//
// Sample curl:
// curl -X GET "http://localhost:3600/mock/api/projects/my-new-project/endpoints/2/responses" -H "Content-Type: application/json"
func ListResponsesHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectId := c.Param("projectId")
	if projectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	// Parse endpoint ID
	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Endpoint ID is required",
		})
		return
	}

	// Check if endpoint exists and belongs to this project
	var endpoint database.MockEndpoint
	result := database.GetDB().Where("id = ? AND project_id = ?", endpointID, projectId).First(&endpoint)
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
