package handlers

import "beo-echo/backend/src/replay/services"

// ReplayHandler handles HTTP requests for replay operations
type replayHandler struct {
	service *services.ReplayService
}

// NewReplayHandler creates a new replay handler
func NewReplayHandler(service *services.ReplayService) *replayHandler {
	return &replayHandler{service: service}
}
