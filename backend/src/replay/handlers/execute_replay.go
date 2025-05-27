package handlers

import (
	"net/http"

	"beo-echo/backend/src/replay/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ExecuteReplayHandler handles POST /projects/{projectId}/replays/execute
func ExecuteReplayHandler(service *services.ReplayService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := zerolog.Ctx(c.Request.Context())
		projectID := c.Param("projectId")

		if projectID == "" {
			log.Error().Msg("missing project ID")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
			return
		}

		var req services.ExecuteReplayRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error().
				Err(err).
				Str("project_id", projectID).
				Msg("invalid request payload")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
			return
		}

		log.Info().
			Str("project_id", projectID).
			Str("protocol", req.Protocol).
			Str("method", req.Method).
			Str("url", req.URL).
			Msg("handling execute replay request")

		result, err := service.ExecuteReplay(c.Request.Context(), projectID, req)
		if err != nil {
			log.Error().
				Err(err).
				Str("project_id", projectID).
				Str("protocol", req.Protocol).
				Str("url", req.URL).
				Msg("failed to execute replay request")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Info().
			Str("project_id", projectID).
			Str("protocol", req.Protocol).
			Int("status_code", result.StatusCode).
			Int("latency_ms", result.LatencyMS).
			Msg("successfully executed replay")

		c.JSON(http.StatusOK, gin.H{
			"result":  result,
			"message": "Replay executed successfully",
		})
	}
}
