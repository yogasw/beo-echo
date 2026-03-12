package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ListReplaysHandler handles GET /projects/{projectId}/replays
func (s *replayHandler) ListReplaysHandler(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())
	projectID := c.Param("projectId")

	if projectID == "" {
		log.Error().Msg("missing project ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	result, err := s.service.ListReplays(c.Request.Context(), projectID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("failed to list replays")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"replays":      result.Replays,
		"folders":      result.Folders,
		"replay_count": len(result.Replays),
		"folder_count": len(result.Folders),
	})
}
