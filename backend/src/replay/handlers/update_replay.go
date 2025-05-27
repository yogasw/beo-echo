package handlers

import (
	"net/http"

	"beo-echo/backend/src/replay/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// UpdateReplayHandler handles PUT /projects/{projectId}/replays/{replayId}
func (s *replayHandler) UpdateReplayHandler(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())
	projectID := c.Param("projectId")
	replayID := c.Param("replayId")

	if projectID == "" || replayID == "" {
		log.Error().
			Str("project_id", projectID).
			Str("replay_id", replayID).
			Msg("missing required parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID and Replay ID are required"})
		return
	}

	var req services.UpdateReplayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Str("replay_id", replayID).
			Msg("invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	log.Info().
		Str("project_id", projectID).
		Str("replay_id", replayID).
		Msg("handling update replay request")

	// First verify the replay belongs to the project
	existingReplay, err := s.service.GetReplay(c.Request.Context(), replayID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Str("replay_id", replayID).
			Msg("replay not found")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if existingReplay.ProjectID != projectID {
		log.Error().
			Str("project_id", projectID).
			Str("replay_id", replayID).
			Str("replay_project_id", existingReplay.ProjectID).
			Msg("replay does not belong to specified project")
		c.JSON(http.StatusNotFound, gin.H{"error": "Replay not found in the specified project"})
		return
	}

	replay, err := s.service.UpdateReplay(c.Request.Context(), replayID, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Str("replay_id", replayID).
			Msg("failed to update replay")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Str("project_id", projectID).
		Str("replay_id", replayID).
		Str("name", replay.Name).
		Msg("successfully updated replay")

	c.JSON(http.StatusOK, gin.H{
		"replay":  replay,
		"message": "Replay updated successfully",
	})
}
