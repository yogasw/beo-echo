package response

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

/*
UpdateResponseHandler updates an existing response

Sample curl:

	curl -X PUT "http://localhost:3600/mock/api/projects/my-new-project/endpoints/2/responses/1" \
	  -H "Content-Type: application/json" \
	  -d '{
	    "statusCode": 201,
	    "body": "{\"message\":\"Resource created\"}",
	    "headers": "{\"Content-Type\":\"application/json\",\"Location\":\"/api/resources/123\"}",
	    "priority": 2,
	    "delayMS": 50,
	    "active": true
	  }'
*/
func UpdateResponseHandler(c *gin.Context) {
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
			"message": "Endpoint not found or doesn't belong to this project",
		})
		return
	}

	// Check if response exists
	var existingResponse database.MockResponse
	result = database.GetDB().Where("id = ? AND endpoint_id = ?", responseID, endpoint.ID).First(&existingResponse)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Response not found",
		})
		return
	}

	// Parse update data
	var updateData struct {
		StatusCode *int    `json:"statusCode"`
		Body       *string `json:"body"`
		Headers    *string `json:"headers"`
		Priority   *int    `json:"priority"`
		DelayMS    *int    `json:"delayMS"`
		Stream     *bool   `json:"stream"`
		Enabled    *bool   `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Apply updates
	if updateData.StatusCode != nil {
		existingResponse.StatusCode = *updateData.StatusCode
	}

	if updateData.Body != nil {
		existingResponse.Body = *updateData.Body
	}

	if updateData.Headers != nil {
		// Validate headers are valid JSON
		var headersMap map[string]string
		if err := json.Unmarshal([]byte(*updateData.Headers), &headersMap); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Headers must be a valid JSON object",
			})
			return
		}
		existingResponse.Headers = *updateData.Headers
	}

	if updateData.Priority != nil {
		existingResponse.Priority = *updateData.Priority
	}

	if updateData.DelayMS != nil {
		existingResponse.DelayMS = *updateData.DelayMS
	}

	if updateData.Stream != nil {
		existingResponse.Stream = *updateData.Stream
	}

	if updateData.Enabled != nil {
		existingResponse.Enabled = *updateData.Enabled
	}

	// Save updates
	result = database.GetDB().Save(&existingResponse)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to update response: " + result.Error.Error(),
		})
		return
	}

	// Reload response with rules
	database.GetDB().
		Preload("Rules").
		Where("id = ?", existingResponse.ID).
		First(&existingResponse)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Response updated successfully",
		"data":    existingResponse,
	})
}
