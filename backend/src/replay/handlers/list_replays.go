package handlers

import (
	"net/http"

	"beo-echo/backend/src/replay/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ListReplaysHandler handles GET /projects/{projectId}/replays
func ListReplaysHandler(service *services.ReplayService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := zerolog.Ctx(c.Request.Context())
		projectID := c.Param("projectId")

		if projectID == "" {
			log.Error().Msg("missing project ID")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
			return
		}

		log.Info().
			Str("project_id", projectID).
			Msg("handling list replays request")

		replays, err := service.ListReplays(c.Request.Context(), projectID)
		if err != nil {
			log.Error().
				Err(err).
				Str("project_id", projectID).
				Msg("failed to list replays")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Info().
			Str("project_id", projectID).
			Int("count", len(replays)).
			Msg("successfully listed replays")

		c.JSON(http.StatusOK, gin.H{
			"replays": replays,
			"count":   len(replays),
		})
	}
}
