package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// GetReplayLogsHandler handles GET /projects/{projectId}/replays/logs
func (s *replayHandler) GetReplayLogsHandler(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())
	projectID := c.Param("projectId")

	if projectID == "" {
		log.Error().Msg("missing project ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

		// Optional query parameter to filter by specific replay
		replayIDParam := c.Query("replay_id")
		var replayID *string
		if replayIDParam != "" {
			replayID = &replayIDParam
		}

		log.Info().
			Str("project_id", projectID).
			Str("replay_id", func() string {
				if replayID != nil {
					return *replayID
				}
				return "all"
			}()).
			Msg("handling get replay logs request")

		logs, err := s.service.GetReplayLogs(c.Request.Context(), projectID, replayID)
		if err != nil {
			log.Error().
				Err(err).
				Str("project_id", projectID).
				Msg("failed to get replay logs")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Info().
			Str("project_id", projectID).
			Int("count", len(logs)).
			Msg("successfully retrieved replay logs")

		c.JSON(http.StatusOK, gin.H{
			"logs":  logs,
			"count": len(logs),
		})
	}
}
