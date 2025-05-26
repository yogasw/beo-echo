package replay

import (
	"beo-echo/backend/src/replay/handlers"
	"beo-echo/backend/src/replay/services"

	"github.com/gin-gonic/gin"
)

// ReplayHandler handles HTTP requests for replay operations
type ReplayHandler struct {
	service *services.ReplayService
}

// NewReplayHandler creates a new replay handler
func NewReplayHandler(service *services.ReplayService) *ReplayHandler {
	return &ReplayHandler{service: service}
}

// ListReplays handles GET /projects/{projectId}/replays
func (h *ReplayHandler) ListReplays(c *gin.Context) {
	handlers.ListReplaysHandler(h.service)(c)
}

// CreateReplay handles POST /projects/{projectId}/replays
func (h *ReplayHandler) CreateReplay(c *gin.Context) {
	handlers.CreateReplayHandler(h.service)(c)
}

// GetReplay handles GET /projects/{projectId}/replays/{replayId}
func (h *ReplayHandler) GetReplay(c *gin.Context) {
	handlers.GetReplayHandler(h.service)(c)
}

// ExecuteReplay handles POST /projects/{projectId}/replays/{replayId}/execute
func (h *ReplayHandler) ExecuteReplay(c *gin.Context) {
	handlers.ExecuteReplayHandler(h.service)(c)
}

// DeleteReplay handles DELETE /projects/{projectId}/replays/{replayId}
func (h *ReplayHandler) DeleteReplay(c *gin.Context) {
	handlers.DeleteReplayHandler(h.service)(c)
}

// GetReplayLogs handles GET /projects/{projectId}/replays/logs
func (h *ReplayHandler) GetReplayLogs(c *gin.Context) {
	handlers.GetReplayLogsHandler(h.service)(c)
}
