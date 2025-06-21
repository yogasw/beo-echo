package project

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
	systemConfig "beo-echo/backend/src/systemConfigs"
)

// CreateProjectWithWorkspaceHandler creates a new project within a specified workspace
// This is a variation of the regular project creation that enforces workspace permission checks
func CreateProjectWithWorkspaceHandler(c *gin.Context) {
	handler.EnsureMockService()

	// Get workspace ID from path
	workspaceID := c.Param("workspaceID")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Workspace ID is required",
		})
		return
	}

	// Get authenticated user
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "User not authenticated",
		})
		return
	}

	// Check if user is system admin or workspace admin
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
	isAllowed := err == nil && user.IsOwner

	if !isAllowed {
		// Check if user is workspace admin
		var userWorkspace database.UserWorkspace
		if err := database.GetDB().
			Where("user_id = ? AND workspace_id = ? AND role = ?", userID, workspaceID, "admin").
			First(&userWorkspace).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   true,
					"message": "You don't have permission to create projects in this workspace",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to verify workspace permissions: " + err.Error(),
			})
			return
		}
	}

	// Check if workspace exists
	var workspace database.Workspace
	if err := database.GetDB().First(&workspace, "id = ?", workspaceID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "Workspace not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to verify workspace: " + err.Error(),
		})
		return
	}

	// Check project limit for the workspace - get user-specific limit if available
	var maxProjectsWorkspace int

	// Use user-specific limit if set, otherwise system default
	if user.MaxProjectsWorkspace != nil {
		maxProjectsWorkspace = *user.MaxProjectsWorkspace
	} else {
		maxProjectsWorkspace, err = systemConfig.GetSystemConfigWithType[int](systemConfig.MAX_WORKSPACE_PROJECTS)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to get project limit configuration: " + err.Error(),
			})
			return
		}
	}

	// Count current projects in the workspace
	var currentProjectCount int64
	if err := database.GetDB().Model(&database.Project{}).Where("workspace_id = ?", workspaceID).Count(&currentProjectCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to count workspace projects: " + err.Error(),
		})
		return
	}

	if int(currentProjectCount) >= maxProjectsWorkspace {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": fmt.Sprintf("Project limit exceeded: maximum %d projects allowed per workspace", maxProjectsWorkspace),
		})
		return
	}

	// Parse request
	var project database.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Set the workspace ID
	project.WorkspaceID = workspaceID

	// Check if alias is already used
	var existingProject database.Project
	result := database.GetDB().Where("alias = ? AND id != ?", project.Alias, project.ID).First(&existingProject)
	if result.Error == nil {
		// Found a project with the same alias
		c.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "Project alias already exists",
		})
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		// Database error other than "not found"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to check alias uniqueness: " + result.Error.Error(),
		})
		return
	}

	// Create the project
	if err := database.GetDB().Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create project: " + err.Error(),
		})
		return
	}

	// combine X-Forwarded-Scheme and host

	scheme := c.Request.Header.Get("X-Forwarded-Scheme")
	if scheme == "" {
		scheme = "http"
	}

	// Add project URL
	project.URL = handler.GetProjectURL(scheme, c.Request.Host, project)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Project created successfully",
		"data":    project,
	})
}
