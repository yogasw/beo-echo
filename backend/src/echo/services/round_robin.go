package services

import (
	"beo-echo/backend/src/database"
	"sync"
	"time"
)

// endpointState holds the counter shared by the round_robin and sequence modes of an endpoint.
type endpointState struct {
	mu        sync.Mutex
	lastIndex int   // index of the last response served
	lastUsed  int64 // Unix timestamp of last usage
}

// Global state map for round-robin/sequence selection per endpoint
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

	return sortedResponses[advanceEndpointCounter(endpointID, len(sortedResponses), database.SequenceLoop)]
}

// getNextSequenceResponse serves responses in priority order. When the last one was served,
// onEnd decides whether to start over (loop) or keep serving the last one (stay_last, default).
func getNextSequenceResponse(endpointID string, responses []database.MockResponse, onEnd database.SequenceOnEnd) database.MockResponse {
	if len(responses) == 0 {
		return database.MockResponse{}
	}

	if onEnd == "" {
		onEnd = database.SequenceStayLast
	}

	sortedResponses := make([]database.MockResponse, len(responses))
	copy(sortedResponses, responses)
	sortByPriority(sortedResponses)

	return sortedResponses[advanceEndpointCounter(endpointID, len(sortedResponses), onEnd)]
}

// advanceEndpointCounter moves the endpoint counter one step and returns the index to serve.
func advanceEndpointCounter(endpointID string, size int, onEnd database.SequenceOnEnd) int {
	now := time.Now().Unix()

	// Load or initialize endpoint state
	val, _ := endpointStates.LoadOrStore(endpointID, &endpointState{
		lastIndex: -1, // Start with -1 so first request gets index 0
		lastUsed:  now,
	})
	state := val.(*endpointState)

	state.mu.Lock()
	state.lastUsed = now
	nextIndex := state.lastIndex + 1
	if nextIndex >= size {
		if onEnd == database.SequenceStayLast {
			nextIndex = size - 1
		} else {
			nextIndex = 0 // Wrap around
		}
	}
	state.lastIndex = nextIndex
	state.mu.Unlock()

	// cleanup stale endpoints
	cleanupStaleEndpoints()

	return nextIndex
}

// ResetEndpointCounter forgets the round_robin/sequence position of an endpoint so the
// next request is served from the first response again.
func ResetEndpointCounter(endpointID string) {
	endpointStates.Delete(endpointID)
}

// cleanupStaleEndpoints removes entries from endpointStates that haven't been used within stateTimeout
func cleanupStaleEndpoints() {
	now := time.Now().Unix()
	endpointStates.Range(func(key, value any) bool {
		state := value.(*endpointState)
		state.mu.Lock()
		stale := now-state.lastUsed > stateTimeout
		state.mu.Unlock()
		if stale {
			endpointStates.Delete(key)
		}
		return true
	})
}
