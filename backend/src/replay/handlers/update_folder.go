package handlers

import (
	"beo-echo/backend/src/replay/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// UpdateFolderHandler handles PATCH /projects/{projectId}/replays/folder/{folderId}
func (s *replayHandler) UpdateFolderHandler(c *gin.Context) {
	projectID := c.Param("projectId")
	folderID := c.Param("folderId")

	if projectID == "" || folderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id and folder_id are required"})
		return
	}

	var req services.UpdateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Msg("invalid request payload for update folder")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder, err := s.service.UpdateFolder(c.Request.Context(), projectID, folderID, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Str("folder_id", folderID).
			Msg("failed to update replay folder")

		if err.Error() == "project not found" || err.Error() == "folder not found" || err.Error() == "invalid parent folder" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update replay folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"folder":  folder,
		"message": "Folder updated successfully",
	})
}
