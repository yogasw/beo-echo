package handlers

import (
	"net/http"

	"beo-echo/backend/src/replay/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// CreateReplayHandler handles POST /projects/{projectId}/replays
func CreateReplayHandler(service *services.ReplayService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := zerolog.Ctx(c.Request.Context())
		projectID := c.Param("projectId")

		if projectID == "" {
			log.Error().Msg("missing project ID")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
			return
		}

		var req services.CreateReplayRequest
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
			Str("alias", req.Alias).
			Str("protocol", req.Protocol).
			Str("method", req.Method).
			Msg("handling create replay request")

		replay, err := service.CreateReplay(c.Request.Context(), projectID, req)
		if err != nil {
			log.Error().
				Err(err).
				Str("project_id", projectID).
				Str("alias", req.Alias).
				Msg("failed to create replay")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Info().
			Str("replay_id", replay.ID).
			Str("project_id", projectID).
			Str("alias", req.Alias).
			Msg("successfully created replay")

		c.JSON(http.StatusCreated, gin.H{
			"replay":  replay,
			"message": "Replay created successfully",
		})
	}
}
