package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/database/repositories"
	"beo-echo/backend/src/echo/services"
)

// CheckAliasAvailabilityRequest represents the request body for alias availability check
type CheckAliasAvailabilityRequest struct {
	Query string `json:"query" binding:"required"`
}

// CheckAliasAvailabilityHandler checks if a project alias is available and searches for similar projects
//
// @Summary Check alias availability and search projects
// @Description Checks if a project alias is available globally and returns projects with similar names in user's workspaces
// @Tags projects
// @Accept json
// @Produce json
// @Param request body CheckAliasAvailabilityRequest true "Search query"
// @Success 200 {object} services.AliasAvailabilityResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/projects/check-alias [post]
//
// Sample curl:
// curl -X POST "http://localhost:3600/api/projects/check-alias" \
//   -H "Content-Type: application/json" \
//   -H "Authorization: Bearer TOKEN" \
//   -d '{"query": "raya"}'
func CheckAliasAvailabilityHandler(c *gin.Context) {
	// Get user ID from context (set by JWTAuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID format",
		})
		return
	}

	// Parse request body
	var request CheckAliasAvailabilityRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("Failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Validate query length
	if len(request.Query) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Query must be at least 1 character long",
		})
		return
	}

	// Initialize project search service
	projectSearchRepo := repositories.NewProjectSearchRepository(database.DB)
	projectSearchService := services.NewProjectSearchService(projectSearchRepo)

	// Check alias availability and search for similar projects
	result, err := projectSearchService.CheckAliasAndSearchProjects(c.Request.Context(), userIDStr, request.Query)
	if err != nil {
		log.Error().Err(err).Str("user_id", userIDStr).Str("query", request.Query).Msg("Failed to check alias availability")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to check alias availability",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
