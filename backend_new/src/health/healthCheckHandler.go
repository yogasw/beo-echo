package health

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler returns information about the system health
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
		"system": gin.H{
			"os":   runtime.GOOS,
			"arch": runtime.GOARCH,
			"go":   runtime.Version(),
		},
	})
}
