package response

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// UpdateResponseHandler updates an existing response
//
// Sample curl:
//
//	curl -X PUT "http://localhost:8000/api/projects/my-project/endpoints/1/responses/1" \
//	  -H "Content-Type: application/json" \
//	  -d '{
//	    "statusCode": 201,
//	    "body": "{\"message\":\"Resource created\"}",
//	    "headers": "{\"Content-Type\":\"application/json\",\"Location\":\"/api/resources/123\"}",
//	    "priority": 2,
//	    "delayMS": 50,
//	    "active": true
//	  }'
func UpdateResponseHandler(c *gin.Context) {
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
		Active     *bool   `json:"active"`
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

	if updateData.Active != nil {
		existingResponse.Active = *updateData.Active
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
