package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// DeleteResponseHandler removes a response
//
// Sample curl:
// curl -X DELETE "http://localhost:3600/mock/api/projects/my-new-project/endpoints/2/responses/1" -H "Content-Type: application/json"
func DeleteResponseHandler(c *gin.Context) {
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
