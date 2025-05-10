package response

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

/*
CreateResponseHandler creates a new response for an endpoint

Sample curl:

	curl -X POST "http://localhost:3600/mock/api/projects/my-new-project/endpoints/2/responses" \
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

	projectName := c.Param("name")
	if projectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project name is required",
		})
		return
	}

	// Parse endpoint ID
	endpointID, err := strconv.ParseUint(c.Param("id"), 10, 32)
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

	// Validate headers are valid JSON if present
	if response.Headers != "" {
		var headersMap map[string]string
		if err := json.Unmarshal([]byte(response.Headers), &headersMap); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Headers must be a valid JSON object",
			})
			return
		}
	}

	// Assign to endpoint
	response.EndpointID = uint(endpointID)

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
