package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
)

// ListProjectsHandler lists projects accessible to the authenticated user
//
// Sample curl:
// curl -X GET "http://localhost:3600/mock/api/projects" -H "Content-Type: application/json" -H "Authorization: Bearer TOKEN"
func ListProjectsHandler(c *gin.Context) {
	handler.EnsureMockService()

	// Get the authenticated user ID from context (set by JWTAuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "User not authenticated",
		})
		return
	}

	// Check if user is a system owner (can see all projects)
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Invalid user ID format",
		})
		return
	}

	// Directly query database to check if user is an owner
	var user database.User
	err := database.GetDB().Where("id = ?", userIDStr).First(&user).Error
	isSystemOwner := err == nil && user.IsOwner

	var projects []database.Project
	var result error

	if isSystemOwner {
		// System owners can see all projects
		result = database.GetDB().Find(&projects).Error
	} else {
		// Regular users can only see projects in their workspaces
		// Get workspace IDs the user has access to
		var workspaceIDs []string
		if err := database.GetDB().Table("user_workspaces").
			Where("user_id = ?", userID).
			Pluck("workspace_id", &workspaceIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to retrieve user workspaces: " + err.Error(),
			})
			return
		}

		// Check if the user has any workspaces
		if len(workspaceIDs) == 0 {
			// User doesn't belong to any workspace, return empty list
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    []database.Project{},
			})
			return
		}

		// Get all projects in the user's workspaces
		result = database.GetDB().Where("workspace_id IN ?", workspaceIDs).Find(&projects).Error
	}

	// Check for database errors
	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve projects: " + result.Error(),
		})
		return
	}

	// Add project URLs
	for i := range projects {
		projects[i].URL = handler.GetProjectURL(c.Request.Host, projects[i])
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projects,
	})
}
