package response

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/database/repositories"
	"beo-echo/backend/src/echo/handler"
)

// ReorderResponsesRequest represents the request body for reordering responses
type ReorderResponsesRequest struct {
	Order []string `json:"order" binding:"required"`
}

/*
ReorderResponsesHandler reorders responses by updating their priority values

Sample curl:

	curl -X PUT "http://localhost:3600/api/workspaces/ws1/projects/my-project/endpoints/endpoint1/responses/reorder" \
	  -H "Content-Type: application/json" \
	  -d '{
	    "order": ["response-id-1", "response-id-2", "response-id-3"]
	  }'
*/
func ReorderResponsesHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectId := c.Param("projectId")
	if projectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID is required",
		})
		return
	}

	endpointID := c.Param("id")
	if endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Endpoint ID is required",
		})
		return
	}

	// Parse request body
	var req ReorderResponsesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request format: " + err.Error(),
		})
		return
	}

	if len(req.Order) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Order array cannot be empty",
		})
		return
	}

	// Create response repository
	db := database.DB
	responseRepo := repositories.NewResponseRepository(db)

	// Validate that all responses belong to the endpoint
	for _, responseID := range req.Order {
		valid, err := responseRepo.ValidateResponseHierarchy(projectId, endpointID, responseID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to validate response: " + err.Error(),
			})
			return
		}
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Response " + responseID + " does not belong to this endpoint",
			})
			return
		}
	}

	// Reorder responses
	err := responseRepo.ReorderResponses(endpointID, req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to reorder responses: " + err.Error(),
		})
		return
	}

	// Fetch updated responses to return
	var responses []database.MockResponse
	result := db.Preload("Rules").Where("endpoint_id = ?", endpointID).Order("priority ASC").Find(&responses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to fetch updated responses: " + result.Error.Error(),
		})
		return
	}

	// Convert to JSON response format
	var jsonResponses []map[string]interface{}
	for _, response := range responses {
		responseMap := map[string]interface{}{
			"id":          response.ID,
			"endpoint_id": response.EndpointID,
			"status_code": response.StatusCode,
			"body":        response.Body,
			"headers":     response.Headers,
			"priority":    response.Priority,
			"delay_ms":    response.DelayMS,
			"stream":      response.Stream,
			"enabled":     response.Enabled,
			"note":        response.Note,
			"created_at":  response.CreatedAt,
			"updated_at":  response.UpdatedAt,
		}

		// Add rules if they exist
		if response.Rules != nil {
			var rules []map[string]interface{}
			rulesData, _ := json.Marshal(response.Rules)
			json.Unmarshal(rulesData, &rules)
			responseMap["rules"] = rules
		} else {
			responseMap["rules"] = nil
		}

		jsonResponses = append(jsonResponses, responseMap)
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Responses reordered successfully",
		"data":    jsonResponses,
	})
}
