package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

/*
UpdateResponseHandler updates an existing response

Sample curl:

	curl -X PUT "http://localhost:3600/api/api/projects/my-new-project/endpoints/2/responses/1" \
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
		StatusCode *int    `json:"status_code"`
		Body       *string `json:"body"`
		Headers    *string `json:"headers"` // Allow headers to be null
		Priority   *int    `json:"priority"`
		DelayMS    *int    `json:"delay_ms"`
		Stream     *bool   `json:"stream"`
		Enabled    *bool   `json:"enabled"`
		Note       *string `json:"note"`
		IsFallback *bool   `json:"is_fallback"`
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
		// Check if headers are empty
		var headers map[string]string
		if *updateData.Headers != "" {
			if err := json.Unmarshal([]byte(*updateData.Headers), &headers); err != nil {
				fmt.Println("Error unmarshalling headers:", err)
			} else {
				// Check if headers are empty
				if len(headers) == 0 {
					existingResponse.Headers = ""
				} else {
					headersJSON, _ := json.Marshal(headers)
					existingResponse.Headers = string(headersJSON)
				}
			}
		}
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

	if updateData.Note != nil {
		existingResponse.Note = *updateData.Note
	}

	if updateData.IsFallback != nil {
		existingResponse.IsFallback = *updateData.IsFallback
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
	// when updating a fallback response, ensure no other fallback exists
	if updateData.IsFallback != nil {
		if *updateData.IsFallback {
			// Find and disable any other fallback responses for this endpoint
			if err := database.GetDB().Model(&database.MockResponse{}).
				Where("endpoint_id = ? AND is_fallback = ? AND id != ?", existingResponse.EndpointID, true, existingResponse.ID).
				Update("is_fallback", false).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   true,
					"message": "Failed to update other fallback responses: " + err.Error(),
				})
				return
			}
		}

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
