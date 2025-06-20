package services

import (
	"beo-echo/backend/src/database"
	"sync"
	"time"
)

// endpointState holds round-robin state for each endpoint
// NOTE: This version is intentionally not thread-safe — acceptable for mock/testing use where race conditions are tolerable.
type endpointState struct {
	lastIndex int   // index of the last response served
	lastUsed  int64 // Unix timestamp of last usage
}

// Global state map for round-robin selection per endpoint
var endpointStates sync.Map

// State timeout for lazy cleanup (1 hour)
const stateTimeout = int64(3600)

// getNextRoundRobinResponse implements simple round-robin selection per endpoint
func getNextRoundRobinResponse(endpointID string, responses []database.MockResponse) database.MockResponse {
	if len(responses) == 0 {
		// This should not happen as we check for empty responses elsewhere,
		// but return empty response as fallback
		return database.MockResponse{}
	}

	// Create a copy of responses to avoid modifying the original slice
	sortedResponses := make([]database.MockResponse, len(responses))
	copy(sortedResponses, responses)

	// Sort by priority (higher priority first)
	sortByPriority(sortedResponses)

	now := time.Now().Unix()

	// Load or initialize endpoint state
	val, _ := endpointStates.LoadOrStore(endpointID, &endpointState{
		lastIndex: -1, // Start with -1 so first request gets index 0
		lastUsed:  now,
	})
	state := val.(*endpointState)

	// Update last used time
	state.lastUsed = now

	// Get next index (round-robin)
	nextIndex := state.lastIndex + 1
	if nextIndex >= len(sortedResponses) {
		nextIndex = 0 // Wrap around
	}

	// Update state with the new index
	state.lastIndex = nextIndex

	// cleanup stale endpoints
	cleanupStaleEndpoints()

	// Return the selected response from sorted array
	return sortedResponses[nextIndex]
}

// cleanupStaleEndpoints removes entries from endpointStates that haven't been used within stateTimeout
func cleanupStaleEndpoints() {
	now := time.Now().Unix()
	endpointStates.Range(func(key, value any) bool {
		state := value.(*endpointState)
		if now-state.lastUsed > stateTimeout {
			endpointStates.Delete(key)
		}
		return true
	})
}
