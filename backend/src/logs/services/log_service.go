package services

import (
	"encoding/json"
	"sync"
	"time"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/repositories"
)

// LogService handles log data retrieval and streaming
type LogService struct {
	Repo             *repositories.LogRepository
	subscribers      map[string][]chan database.RequestLog
	subscribersMutex sync.RWMutex
}

// NewLogService creates a new log service
func NewLogService(repo *repositories.LogRepository) *LogService {
	return &LogService{
		Repo:        repo,
		subscribers: make(map[string][]chan database.RequestLog),
	}
}

// GetPaginatedLogs retrieves logs with pagination
func (s *LogService) GetPaginatedLogs(page, pageSize int, projectID string) ([]database.RequestLog, int64, error) {
	// Default to page 1 if invalid
	if page < 1 {
		page = 1
	}

	// Limit maximum page size
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 100
	}

	return s.Repo.GetLogs(page, pageSize, projectID)
}

// GetLatestLogs retrieves the most recent logs
func (s *LogService) GetLatestLogs(limit int, projectID string) ([]database.RequestLog, error) {
	// Default and limit
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	return s.Repo.GetLatestLogs(limit, projectID)
}

// SubscribeToLogs registers a channel to receive new logs
func (s *LogService) SubscribeToLogs(projectID string) chan database.RequestLog {
	channel := make(chan database.RequestLog, 100) // Buffer to prevent blocking

	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()

	if _, exists := s.subscribers[projectID]; !exists {
		s.subscribers[projectID] = make([]chan database.RequestLog, 0)
	}

	s.subscribers[projectID] = append(s.subscribers[projectID], channel)

	return channel
}

// UnsubscribeFromLogs removes a channel from receiving logs
func (s *LogService) UnsubscribeFromLogs(projectID string, channel chan database.RequestLog) {
	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()

	if channels, exists := s.subscribers[projectID]; exists {
		for i, ch := range channels {
			if ch == channel {
				// Remove channel from slice
				s.subscribers[projectID] = append(channels[:i], channels[i+1:]...)
				break
			}
		}
	}

	close(channel)
}

// NotifySubscribers sends a log to all subscribers for a project
// This should be called whenever a new log is created
func (s *LogService) NotifySubscribers(log database.RequestLog) {
	s.subscribersMutex.RLock()
	defer s.subscribersMutex.RUnlock()

	// Notify subscribers for this project
	if channels, exists := s.subscribers[log.ProjectID]; exists {
		for _, channel := range channels {
			// Use non-blocking send to avoid deadlocks
			select {
			case channel <- log:
				// Log sent successfully
			default:
				// Channel is full, skip this one
			}
		}
	}

	// Also notify subscribers for all projects (empty project ID)
	if channels, exists := s.subscribers[""]; exists {
		for _, channel := range channels {
			select {
			case channel <- log:
				// Log sent successfully
			default:
				// Channel is full, skip this one
			}
		}
	}
}

// FormatSSEEvent formats a log as a Server-Sent Event message
func FormatSSEEvent(log database.RequestLog, eventType string) string {
	// Serialize the log to JSON
	data, err := json.Marshal(log)
	if err != nil {
		return ""
	}

	// Format as SSE message
	result := "event: " + eventType + "\n"
	result += "data: " + string(data) + "\n\n"

	return result
}

// FormatSSEPingEvent creates a ping event for keeping the connection alive
func FormatSSEPingEvent() string {
	data := map[string]string{
		"time": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return "event: ping\ndata: " + string(jsonData) + "\n\n"
}
