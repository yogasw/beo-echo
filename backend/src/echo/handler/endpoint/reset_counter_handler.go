package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
	"beo-echo/backend/src/echo/services"
)

/*
ResetCounterHandler resets the round_robin / sequence counter of an endpoint so the next
request is served from the first response again.

Sample curl:

	curl -X POST "http://localhost:3600/api/workspaces/my-workspace/projects/my-project-id/endpoints/9585df96-32e6-4d30-8f63-7bf5cd783b05/reset-counter" \
	  -H "Authorization: Bearer $TOKEN"
*/
func ResetCounterHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectId := c.Param("projectId")
	endpointID := c.Param("id")
	if projectId == "" || endpointID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project ID and endpoint ID are required",
		})
		return
	}

	var endpoint database.MockEndpoint
	if err := database.GetDB().Where("id = ? AND project_id = ?", endpointID, projectId).First(&endpoint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Endpoint not found",
		})
		return
	}

	services.ResetEndpointCounter(endpoint.ID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Endpoint counter reset successfully",
	})
}
