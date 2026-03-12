package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
)

// ListProjectsHandler lists projects accessible to the authenticated user
//
// Sample curl:
// curl -X GET "http://localhost:3600/api/api/projects" -H "Content-Type: application/json" -H "Authorization: Bearer TOKEN"

// ProjectListItem is the response shape for a single project in the list.
// It mirrors database.Project but adds the per-user computed `IsPinned` flag.
type ProjectListItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	WorkspaceID   string `json:"workspace_id"`
	Mode          string `json:"mode"`
	Status        string `json:"status"`
	ActiveProxyID *string `json:"active_proxy_id"`
	Alias         string `json:"alias"`
	URL           string `json:"url"`
	Documentation string `json:"documentation"`
	AdvanceConfig string `json:"advance_config"`
	IsPinned      bool   `json:"is_pinned"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// GetWorkspaceProjectsHandler returns all projects in a workspace
func ListProjectsHandler(c *gin.Context) {
	workspaceID := c.Param("workspaceID")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Workspace ID is required",
		})
		return
	}

	// Get user ID from context (set by JWTAuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Check if user is a system admin (can access all workspaces)
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID format",
		})
		return
	}

	// Directly query database to check if user is an owner
	var user database.User
	err := database.DB.Where("id = ?", userIDStr).First(&user).Error
	isSystemOwner := err == nil && user.IsOwner

	if !isSystemOwner {
		// Check if the user is a member of this workspace
		var userWorkspace database.UserWorkspace
		err := database.DB.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&userWorkspace).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "You do not have access to this workspace",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to verify workspace access: " + err.Error(),
			})
			return
		}
	}

	var projects []database.Project
	if err := database.DB.Where("workspace_id = ?", workspaceID).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch projects: " + err.Error(),
		})
		return
	}

	scheme := c.Request.Header.Get("X-Forwarded-Scheme")
	if scheme == "" {
		scheme = "http"
	}

	// Build a set of pinned project IDs for the current user in this workspace
	var pinnedRows []database.UserPinnedProject
	database.DB.Where("user_id = ? AND workspace_id = ?", userIDStr, workspaceID).Find(&pinnedRows)
	pinnedSet := make(map[string]bool, len(pinnedRows))
	for _, row := range pinnedRows {
		pinnedSet[row.ProjectID] = true
	}

	// Build response items with URL and IsPinned enrichment
	items := make([]ProjectListItem, 0, len(projects))
	for i := range projects {
		p := &projects[i]
		p.URL = handler.GetProjectURL(scheme, c.Request.Host, *p)
		items = append(items, ProjectListItem{
			ID:            p.ID,
			Name:          p.Name,
			WorkspaceID:   p.WorkspaceID,
			Mode:          string(p.Mode),
			Status:        p.Status,
			ActiveProxyID: p.ActiveProxyID,
			Alias:         p.Alias,
			URL:           p.URL,
			Documentation: p.Documentation,
			AdvanceConfig: p.AdvanceConfig,
			IsPinned:      pinnedSet[p.ID],
			CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
	})
}

