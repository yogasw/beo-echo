package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/mocks/handler"
)

/*
CreateResponseHandler creates a new response for an endpoint

Sample curl:

	curl -X POST "http://localhost:3600/mock/api/projects/my-new-project/endpoints/9585df96-32e6-4d30-8f63-7bf5cd783b05/responses" \
	  -H "Content-Type: application/json" \
	  -d '{
	    "statusCode": 200,
	    "body": "{\"message\":\"Hello World\"}",
	    "headers": "{\"Content-Type\":\"application/json\"}",
	    "priority": 1,
	    "delayMS": 0,
	    "stream": false,
	    "active": true
	  }'
*/
func CreateResponseHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectId := c.Param("projectId")
	if projectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	// Check if project exists
	if err := database.GetDB().
		Where("id = ?", projectId).
		First(&database.Project{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Project not found: " + err.Error(),
		})
		return
	}

	// Parse endpoint ID
	endpointIDStr := c.Param("id")
	// Check if endpoint exists and belongs to this project
	var endpoint database.MockEndpoint
	result := database.GetDB().Where("id = ? AND project_id = ?", endpointIDStr, projectId).First(&endpoint)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Endpoint not found or doesn't belong to this project",
		})
		return
	}

	// Parse response data
	var response database.MockResponse
	if err := c.ShouldBindJSON(&response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Basic validation
	if response.StatusCode == 0 {
		response.StatusCode = 200 // Default to 200 OK
	}

	// Assign to endpoint
	response.EndpointID = endpointIDStr

	// Create response
	result = database.GetDB().Create(&response)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create response: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Response created successfully",
		"data":    response,
	})
}
