package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/mocks/handler"
)

// DeleteEndpointHandler removes an endpoint
//
// Sample curl:
// curl -X DELETE "http://localhost:3600/mock/api/projects/my-new-project/endpoints/1" -H "Content-Type: application/json"
func DeleteEndpointHandler(c *gin.Context) {
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

	// Check if endpoint exists
	var endpoint database.MockEndpoint
	result := database.GetDB().Where("id = ? AND project_id = ?", endpointID, projectId).First(&endpoint)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Endpoint not found",
		})
		return
	}

	// Delete the endpoint (GORM will cascade delete related responses and rules)
	result = database.GetDB().Delete(&endpoint)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to delete endpoint: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Endpoint deleted successfully",
	})
}
