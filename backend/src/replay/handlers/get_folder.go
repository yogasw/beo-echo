package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// GetFolderHandler handles GET /projects/{projectId}/replays/folder/{folderId}
func (s *replayHandler) GetFolderHandler(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())
	projectID := c.Param("projectId")
	folderID := c.Param("folderId")

	if projectID == "" || folderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID and folder ID are required"})
		return
	}

	folder, err := s.service.GetFolder(c.Request.Context(), projectID, folderID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get folder")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"folder": folder})
}
