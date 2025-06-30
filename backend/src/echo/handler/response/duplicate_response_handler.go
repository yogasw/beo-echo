package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

// DuplicateResponseHandler duplicates an existing response with all its rules
//
// Sample curl:
//
//	curl -X POST "http://localhost:3600/api/workspaces/{workspaceID}/projects/{projectId}/endpoints/{id}/responses/{responseId}/duplicate" \
//	  -H "Content-Type: application/json" \
//	  -H "Authorization: Bearer {token}"
func DuplicateResponseHandler(c *gin.Context) {
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

	// Get the original response with all its rules
	var originalResponse database.MockResponse
	result = database.GetDB().
		Preload("Rules").
		Where("id = ? AND endpoint_id = ?", responseID, endpoint.ID).
		First(&originalResponse)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Response not found",
		})
		return
	}

	// Create a new response by copying the original
	duplicatedResponse := database.MockResponse{
		ID:         uuid.New().String(), // Generate new ID
		EndpointID: originalResponse.EndpointID,
		StatusCode: originalResponse.StatusCode,
		Body:       originalResponse.Body,
		Headers:    originalResponse.Headers,
		Priority:   originalResponse.Priority,
		DelayMS:    originalResponse.DelayMS,
		Stream:     originalResponse.Stream,
		Note:       originalResponse.Note + " (Copy)", // Add "(Copy)" to distinguish
		Enabled:    originalResponse.Enabled,
		// Don't copy Rules here - we'll handle them separately
	}

	// Start a database transaction to ensure all operations succeed or fail together
	tx := database.GetDB().Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to start transaction: " + tx.Error.Error(),
		})
		return
	}

	// Create the duplicated response
	if err := tx.Create(&duplicatedResponse).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create duplicated response: " + err.Error(),
		})
		return
	}

	// Duplicate all rules from the original response
	for _, originalRule := range originalResponse.Rules {
		duplicatedRule := database.MockRule{
			ID:         uuid.New().String(),   // Generate new ID
			ResponseID: duplicatedResponse.ID, // Link to the new response
			Type:       originalRule.Type,
			Key:        originalRule.Key,
			Operator:   originalRule.Operator,
			Value:      originalRule.Value,
		}

		if err := tx.Create(&duplicatedRule).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to create duplicated rule: " + err.Error(),
			})
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to commit transaction: " + err.Error(),
		})
		return
	}

	// Reload the duplicated response with its rules for the response
	var responseWithRules database.MockResponse
	database.GetDB().
		Preload("Rules").
		Where("id = ?", duplicatedResponse.ID).
		First(&responseWithRules)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Response duplicated successfully",
		"data":    responseWithRules,
	})
}
