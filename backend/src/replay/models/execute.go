package models

// ExecuteReplayRequest represents the request payload for executing a replay test
type ExecuteReplayRequest struct {
	Protocol string            `json:"protocol" binding:"required"` // http, https, ws, grpc
	Method   string            `json:"method" binding:"required"`   // HTTP method or operation type
	URL      string            `json:"url" binding:"required"`      // Target URL
	Headers  map[string]string `json:"headers"`                     // Request headers
	Payload  string            `json:"payload"`                     // Request body/payload/content
	Query    map[string]string `json:"query"`                       // Query parameters
	Metadata map[string]string `json:"metadata"`                    // Additional protocol-specific metadata
}

// ExecuteReplayResponse represents the response from executing a replay
type ExecuteReplayResponse struct {
	ReplayID        string            `json:"replay_id"`
	StatusCode      int               `json:"status_code"`
	StatusText      string            `json:"status_text"`
	ResponseBody    string            `json:"response_body"`
	ResponseHeaders map[string]string `json:"response_headers"`
	LatencyMS       int               `json:"latency_ms"`
	Size            int64             `json:"size"`
	Error           string            `json:"error,omitempty"`
	LogID           string            `json:"log_id"`
}
