package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

// GetResponseHandler retrieves a response by ID
//
// Sample curl:
// curl -X GET "http://localhost:3600/api/api/projects/my-new-project/endpoints/2/responses/1" -H "Content-Type: application/json"
func GetResponseHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectId := c.Param("projectId")
	if projectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	// Parse endpoint ID and response ID
	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Endpoint ID is required",
		})
		return
	}

	responseID := c.Param("responseId")
	if responseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Response ID is required",
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

	// Get response
	var response database.MockResponse
	result = database.GetDB().
		Preload("Rules").
		Where("id = ? AND endpoint_id = ?", responseID, endpoint.ID).
		First(&response)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Response not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
