package services

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
)

// ExecuteReplay executes a replay request with the provided configuration
func (s *ReplayService) ExecuteReplay(ctx context.Context, projectID string, req ExecuteReplayRequest) (*ExecuteReplayResponse, error) {
	log := zerolog.Ctx(ctx)

	log.Info().
		Str("project_id", projectID).
		Str("protocol", req.Protocol).
		Str("method", req.Method).
		Str("url", req.URL).
		Msg("executing replay request")

	// Validate project exists
	_, err := s.repo.FindProjectByID(ctx, projectID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("project not found")
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Validate protocol
	protocol := strings.ToLower(req.Protocol)
	if protocol != "http" && protocol != "https" {
		return nil, fmt.Errorf("unsupported protocol: %s (supported: http, https)", req.Protocol)
	}

	startTime := time.Now()
	replayID := uuid.New().String()

	// Build URL with query parameters
	targetURL := req.URL
	if len(req.Query) > 0 {
		u, err := url.Parse(req.URL)
		if err != nil {
			log.Error().
				Err(err).
				Str("url", req.URL).
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
	resp, err := s.client.Do(httpReq)
	latencyMS := int(time.Since(startTime).Milliseconds())

	response := &ExecuteReplayResponse{
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

	// Convert response headers
	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = values[0] // Take first value
		}
	}
	response.ResponseHeaders = respHeaders

	log.Info().
		Str("replay_id", replayID).
		Str("project_id", projectID).
		Int("status_code", resp.StatusCode).
		Int("latency_ms", latencyMS).
		Msg("successfully executed replay request")

	return response, nil
}
