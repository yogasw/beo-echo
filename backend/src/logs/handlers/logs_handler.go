package handler

import (
	"beo-echo/backend/src/database"
	"beo-echo/backend/src/logs/services"
	"beo-echo/backend/src/mocks/repositories"
	"net/http"
	"strconv"
	"sync"
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

	// Send initial ping event to establish SSE connection
	pingEvent := services.FormatSSEPingEvent()
	if c.Writer != nil {
		_, err := c.Writer.Write([]byte(pingEvent))
		if err == nil {
			if flusher, ok := c.Writer.(http.Flusher); ok && flusher != nil {
				flusher.Flush()
			}
		}
	}

	// Create a client connection close notifier
	clientGone := c.Writer.CloseNotify()

	// Create a mutex to protect access to the writer
	var writerMutex sync.Mutex

	// Create a channel to signal ping goroutine to stop
	stopPing := make(chan struct{})

	// Create a once guard to ensure unsubscription happens only once
	var unsubscribeOnce sync.Once

	// Send ping to keep the connection alive
	go func() {
		pingTicker := time.NewTicker(30 * time.Second)
		defer pingTicker.Stop()

		for {
			select {
			case <-pingTicker.C:
				// Send ping with mutex protection
				writerMutex.Lock()
				// Check if connection is still alive before writing
				select {
				case <-clientGone:
					writerMutex.Unlock()
					return
				default:
					// Connection is still alive, proceed with writing
					pingEvent := services.FormatSSEPingEvent()
					if c.Writer != nil {
						_, err := c.Writer.Write([]byte(pingEvent))
						if err == nil {
							// Only flush if Write was successful
							if flusher, ok := c.Writer.(http.Flusher); ok && flusher != nil {
								flusher.Flush()
							}
						}
					}
					writerMutex.Unlock()
				}
			case <-clientGone:
				// Client disconnected
				unsubscribeOnce.Do(func() {
					logService.UnsubscribeFromLogs(projectID, logChannel)
				})
				return
			case <-stopPing:
				// Explicit signal to stop
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
				unsubscribeOnce.Do(func() {
					// Channel already closed, just cleanup
					close(stopPing)
				})
				return
			}

			// Send log as SSE event with mutex protection
			writerMutex.Lock()
			sseData := services.FormatSSEEvent(log, "log")
			if c.Writer != nil {
				_, err := c.Writer.Write([]byte(sseData))
				if err == nil {
					// Only flush if Write was successful
					if flusher, ok := c.Writer.(http.Flusher); ok && flusher != nil {
						flusher.Flush()
					}
				}
			}
			writerMutex.Unlock()

		case <-clientGone:
			// Client disconnected
			unsubscribeOnce.Do(func() {
				logService.UnsubscribeFromLogs(projectID, logChannel)
				close(stopPing) // Signal ping goroutine to stop
			})
			return
		}
	}
}
