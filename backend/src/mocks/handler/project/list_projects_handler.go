package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// ListProjectsHandler lists all projects
//
// Sample curl:
// curl -X GET "http://localhost:3600/mock/api/projects" -H "Content-Type: application/json"
func ListProjectsHandler(c *gin.Context) {
	handler.EnsureMockService()

	var projects []database.Project
	result := database.GetDB().Find(&projects)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve projects: " + result.Error.Error(),
		})
		return
	}
	//loop prject and insert url
	for i := range projects {
		projects[i].URL = handler.GetProjectURL(c.Request.Host, projects[i])
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projects,
	})
}
