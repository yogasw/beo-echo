package handlers

import (
	"net/http"

	"beo-echo/backend/src/replay/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// DeleteReplayHandler handles DELETE /projects/{projectId}/replays/{replayId}
func DeleteReplayHandler(service *services.ReplayService) gin.HandlerFunc {
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
			Msg("handling delete replay request")

		// First verify the replay belongs to the project
		replay, err := service.GetReplay(c.Request.Context(), replayID)
		if err != nil {
			log.Error().
				Err(err).
				Str("project_id", projectID).
				Str("replay_id", replayID).
				Msg("replay not found")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if replay.ProjectID != projectID {
			log.Error().
				Str("project_id", projectID).
				Str("replay_id", replayID).
				Str("replay_project_id", replay.ProjectID).
				Msg("replay does not belong to specified project")
			c.JSON(http.StatusNotFound, gin.H{"error": "Replay not found in the specified project"})
			return
		}

		err = service.DeleteReplay(c.Request.Context(), replayID)
		if err != nil {
			log.Error().
				Err(err).
				Str("project_id", projectID).
				Str("replay_id", replayID).
				Msg("failed to delete replay")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Info().
			Str("project_id", projectID).
			Str("replay_id", replayID).
			Msg("successfully deleted replay")

		c.JSON(http.StatusOK, gin.H{
			"message": "Replay deleted successfully",
		})
	}
}
