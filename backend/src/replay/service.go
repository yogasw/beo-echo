package replay

import "beo-echo/backend/src/replay/services"

// NewReplayService creates a new replay service
func NewReplayService(repo services.ReplayRepository) *services.ReplayService {
	return services.NewReplayService(repo)
}
