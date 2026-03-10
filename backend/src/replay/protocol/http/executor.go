package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"beo-echo/backend/src/replay/models"
	"beo-echo/backend/src/replay/protocol"
)

// ensure HttpExecutor implements protocol.Executor
var _ protocol.Executor = (*HttpExecutor)(nil)

// HttpExecutor implements the Executor interface for HTTP/HTTPS protocols
type HttpExecutor struct {
	client *http.Client
}

// NewExecutor creates a new HttpExecutor
func NewExecutor() *HttpExecutor {
	return &HttpExecutor{
		client: &http.Client{
			// Note: timeout could also be configurable per request later
			Timeout: 30 * time.Second,
		},
	}
}

// Execute performs an HTTP request based on the replay configuration
func (e *HttpExecutor) Execute(ctx context.Context, projectID string, req models.ExecuteReplayRequest) (*models.ExecuteReplayResponse, error) {
	log := zerolog.Ctx(ctx)
	startTime := time.Now()
	replayID := uuid.New().String()

	// Build URL with query parameters
	targetURL := req.URL
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "http://" + targetURL
	}

	if len(req.Query) > 0 {
		u, err := url.Parse(targetURL)
		if err != nil {
			log.Error().
				Err(err).
				Str("url", targetURL).
				Msg("invalid URL format")
			return nil, fmt.Errorf("invalid URL format: %w", err)
		}

		q := u.Query()
		for key, value := range req.Query {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
		targetURL = u.String()
	}

	// Create HTTP request
	var reqBody io.Reader
	if req.Body != "" {
		reqBody = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, strings.ToUpper(req.Method), targetURL, reqBody)
	if err != nil {
		log.Error().
			Err(err).
			Str("target_url", targetURL).
			Msg("failed to create HTTP request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Execute request
	resp, err := e.client.Do(httpReq)
	latencyMS := int(time.Since(startTime).Milliseconds())

	response := &models.ExecuteReplayResponse{
		ReplayID:  replayID,
		LatencyMS: latencyMS,
	}

	if err != nil {
		log.Error().
			Err(err).
			Str("replay_id", replayID).
			Str("target_url", targetURL).
			Int("latency_ms", latencyMS).
			Msg("request execution failed")

		response.Error = err.Error()
		response.StatusCode = 0

		return response, nil
	}
	defer resp.Body.Close()

	// Read response body
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().
			Err(err).
			Str("replay_id", replayID).
			Msg("failed to read response body")
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	respBody := string(respBodyBytes)
	response.StatusCode = resp.StatusCode
	response.ResponseBody = respBody
	response.StatusText = resp.Status
	response.Size = int64(len(respBodyBytes))

	// Convert response headers
	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = values[0] // Take first value
		}
	}
	response.ResponseHeaders = respHeaders

	return response, nil
}
