package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// DeleteFolderEndpoint creates a gin handler for deleting a replay folder
func (s *replayHandler) DeleteFolderEndpoint(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	projectID := c.Param("projectId")
	folderID := c.Param("folderId")

	if projectID == "" || folderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id and folder_id are required"})
		return
	}

	err := s.service.DeleteFolder(c.Request.Context(), projectID, folderID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Str("folder_id", folderID).
			Msg("failed to delete folder")

		if err.Error() == "project not found" || err.Error() == "folder not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "folder deleted successfully",
	})
}
