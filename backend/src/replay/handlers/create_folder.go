package handlers

import (
	"net/http"

	"beo-echo/backend/src/replay/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// CreateFolderHandler handles POST /projects/{projectId}/replays/folder
func (s *replayHandler) CreateFolderHandler(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())
	projectID := c.Param("projectId")

	if projectID == "" {
		log.Error().Msg("missing project ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	var req services.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("invalid request payload for folder")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	folder, err := s.service.CreateFolder(c.Request.Context(), projectID, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Str("name", req.Name).
			Msg("failed to create replay folder")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"folder":  folder,
		"message": "Folder created successfully",
	})
}
