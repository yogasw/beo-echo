package handler

import (
	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/repositories"
	"mockoon-control-panel/backend_new/src/mocks/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var logService *services.LogService

// InitLogService initializes the log service
func InitLogService() {
	db := database.GetDB()
	if db == nil {
		return
	}

	repo := repositories.NewLogRepository(db)
	logService = services.NewLogService(repo)
}

// EnsureLogService ensures log service is initialized
func EnsureLogService() {
	if logService == nil {
		InitLogService()
	}
}

// LogService returns the log service instance
func LogService() *services.LogService {
	return logService
}

// GetLogsHandler handles retrieving logs with pagination
func GetLogsHandler(c *gin.Context) {
	EnsureLogService()
	if logService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Log service is not available",
		})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))
	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project id is required",
		})
		return
	}

	// Get logs
	logs, total, err := logService.GetPaginatedLogs(page, pageSize, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Error retrieving logs: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// StreamLogsHandler handles streaming logs using Server-Sent Events
func StreamLogsHandler(c *gin.Context) {
	EnsureLogService()
	if logService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Log service is not available",
		})
		return
	}

	// Get the project ID from query parameters
	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Project id is required",
		})
		return
	}

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Create channel for this client
	logChannel := logService.SubscribeToLogs(projectID)

	// Send initial batch of logs (most recent 100 first)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	initialLogs, err := logService.GetLatestLogs(limit, projectID)
	if err == nil {
		// Send initial logs from oldest to newest
		for i := len(initialLogs) - 1; i >= 0; i-- {
			sseData := services.FormatSSEEvent(initialLogs[i], "log")
			c.Writer.Write([]byte(sseData))
			c.Writer.Flush()
		}
	}

	// Create a client connection close notifier
	clientGone := c.Writer.CloseNotify()

	// Send ping to keep the connection alive
	go func() {
		pingTicker := time.NewTicker(30 * time.Second)
		defer pingTicker.Stop()

		for {
			select {
			case <-pingTicker.C:
				// Send ping
				pingEvent := services.FormatSSEPingEvent()
				c.Writer.Write([]byte(pingEvent))
				c.Writer.Flush()
			case <-clientGone:
				// Client disconnected
				logService.UnsubscribeFromLogs(projectID, logChannel)
				return
			}
		}
	}()

	// Main loop to listen for new logs
	for {
		select {
		case log, ok := <-logChannel:
			// Check if channel is closed
			if !ok {
				return
			}

			// Send log as SSE event
			sseData := services.FormatSSEEvent(log, "log")
			c.Writer.Write([]byte(sseData))
			c.Writer.Flush()

		case <-clientGone:
			// Client disconnected
			logService.UnsubscribeFromLogs(projectID, logChannel)
			return
		}
	}
}
