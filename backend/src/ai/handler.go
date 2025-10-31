package ai

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// AIHandler handles HTTP requests for AI operations
type AIHandler struct {
	service *AIService
}

// NewAIHandler creates a new AI handler
func NewAIHandler(service *AIService) *AIHandler {
	return &AIHandler{service: service}
}

// GenerateHandler handles POST /ai/generate
func (h *AIHandler) GenerateHandler(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Msg("invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	log.Info().
		Str("template", req.Template).
		Msg("handling AI generate request")

	resp, err := h.service.Generate(c.Request.Context(), req)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to generate content")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Str("model", resp.Model).
		Int("token_used", resp.TokenUsed).
		Msg("successfully generated content")

	c.JSON(http.StatusOK, resp)
}
