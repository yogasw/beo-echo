package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// GetEndpointHandler retrieves an endpoint by ID
//
// Sample curl:
// curl -X GET "http://localhost:8000/api/projects/my-project/endpoints/1" -H "Content-Type: application/json"
func GetEndpointHandler(c *gin.Context) {
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

	// Get endpoint
	var endpoint database.MockEndpoint
	result := database.GetDB().
		Preload("Responses").
		Where("id = ? AND project_id = ?", endpointID, projectId).
		First(&endpoint)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Endpoint not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    endpoint,
	})
}
