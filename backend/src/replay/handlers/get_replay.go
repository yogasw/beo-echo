package handlers

import (
	"net/http"

	"beo-echo/backend/src/replay/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// GetReplayHandler handles GET /projects/{projectId}/replays/{replayId}
func GetReplayHandler(service *services.ReplayService) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		log.Info().
			Str("project_id", projectID).
			Str("replay_id", replayID).
			Msg("handling get replay request")

		replay, err := service.GetReplay(c.Request.Context(), replayID)
		if err != nil {
			log.Error().
				Err(err).
				Str("project_id", projectID).
				Str("replay_id", replayID).
				Msg("failed to get replay")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// Verify replay belongs to the specified project
		if replay.ProjectID != projectID {
			log.Error().
				Str("project_id", projectID).
				Str("replay_id", replayID).
				Str("replay_project_id", replay.ProjectID).
				Msg("replay does not belong to specified project")
			c.JSON(http.StatusNotFound, gin.H{"error": "Replay not found in the specified project"})
			return
		}

		log.Info().
			Str("project_id", projectID).
			Str("replay_id", replayID).
			Str("alias", replay.Alias).
			Msg("successfully retrieved replay")

		c.JSON(http.StatusOK, gin.H{
			"replay": replay,
		})
	}
}
