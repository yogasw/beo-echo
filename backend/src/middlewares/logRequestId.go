package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func LogRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()

		// Create logger with request context
		logger := log.With().
			Str("request_id", requestID).
			Logger()

		// Add logger and request ID to gin context
		c.Set("request_id", requestID)

		// Add logger to request context for deeper integration
		ctx := logger.WithContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		// Continue to next handler
		c.Next()

	}
}
