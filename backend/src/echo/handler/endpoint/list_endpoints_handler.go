package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

// ListEndpointsHandler lists all endpoints for a project
//
// Sample curl:
// curl -X GET "http://localhost:3600/api/api/projects/my-new-project/endpoints" -H "Content-Type: application/json"
func ListEndpointsHandler(c *gin.Context) {
	handler.EnsureMockService()

	projectId := c.Param("projectId")
	if projectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project id is required",
		})
		return
	}

	// Get endpoints for this project
	var endpoints []database.MockEndpoint
	result := database.GetDB().
		Where("project_id = ?", projectId).
		Order("path").
		Order("method").
		Find(&endpoints)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve endpoints: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    endpoints,
	})
}
