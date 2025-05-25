package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"beo-echo/backend/src/database"

	"github.com/rs/zerolog"
)

// ExecuteReplay executes a replay and logs the request/response
func (s *ReplayService) ExecuteReplay(ctx context.Context, replayID string) (*ExecuteReplayResponse, error) {
	log := zerolog.Ctx(ctx)

	log.Info().
		Str("replay_id", replayID).
		Msg("executing replay")

	// Get replay configuration
	replay, err := s.repo.FindByID(ctx, replayID)
	if err != nil {
		log.Error().
			Err(err).
			Str("replay_id", replayID).
			Msg("replay not found")
		return nil, fmt.Errorf("replay not found: %w", err)
	}

	// Currently only support HTTP protocol
	if strings.ToLower(string(replay.Protocol)) != "http" {
		return nil, fmt.Errorf("unsupported protocol: %s", replay.Protocol)
	}

	startTime := time.Now()

	// Parse headers
	var headers map[string]string
	if replay.Headers != "" {
		err = json.Unmarshal([]byte(replay.Headers), &headers)
		if err != nil {
			log.Error().
				Err(err).
				Str("replay_id", replayID).
				Msg("failed to parse headers")
			return nil, fmt.Errorf("invalid headers in replay: %w", err)
		}
	}

	// Create HTTP request
	var reqBody io.Reader
	if replay.Payload != "" {
		reqBody = strings.NewReader(replay.Payload)
	}

	httpReq, err := http.NewRequestWithContext(ctx, replay.Method, replay.Url, reqBody)
	if err != nil {
		log.Error().
			Err(err).
			Str("replay_id", replayID).
			Str("target_url", replay.Url).
			Msg("failed to create HTTP request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
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
			Str("target_url", replay.Url).
			Int("latency_ms", latencyMS).
			Msg("request execution failed")

		response.Error = err.Error()
		response.StatusCode = 0

		// Log the failed execution
		s.logReplayExecution(ctx, replay, "", "", map[string]string{}, 0, "", map[string]string{}, latencyMS, err.Error())

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

	// Log the execution
	logID, err := s.logReplayExecution(ctx, replay, replay.Payload, "", headers, resp.StatusCode, respBody, respHeaders, latencyMS, "")
	if err != nil {
		log.Warn().
			Err(err).
			Str("replay_id", replayID).
			Msg("failed to log replay execution")
	} else {
		response.LogID = logID
	}

	log.Info().
		Str("replay_id", replayID).
		Int("status_code", resp.StatusCode).
		Int("latency_ms", latencyMS).
		Msg("successfully executed replay")

	return response, nil
}

// logReplayExecution creates a request log entry for replay execution
func (s *ReplayService) logReplayExecution(ctx context.Context, replay *database.Replay, requestBody, queryParams string, requestHeaders map[string]string, statusCode int, responseBody string, responseHeaders map[string]string, latencyMS int, errorMsg string) (string, error) {
	log := zerolog.Ctx(ctx)

	// Convert headers to JSON
	reqHeadersJSON, _ := json.Marshal(requestHeaders)
	respHeadersJSON, _ := json.Marshal(responseHeaders)

	// Create request log
	requestLog := &database.RequestLog{
		ProjectID:       replay.ProjectID,
		Method:          replay.Method,
		Path:            replay.Url,
		QueryParams:     queryParams,
		RequestHeaders:  string(reqHeadersJSON),
		RequestBody:     requestBody,
		ResponseStatus:  statusCode,
		ResponseBody:    responseBody,
		ResponseHeaders: string(respHeadersJSON),
		LatencyMS:       latencyMS,
		Source:          database.RequestSourceReplay,
		ExecutionMode:   database.ModeProxy, // Replay acts like proxy mode
		Matched:         true,               // Replay always "matches" its configuration
	}

	// Generate hash for integrity
	requestLog.LogsHash = fmt.Sprintf("%d_%s_%s", time.Now().Unix(), replay.ID, requestLog.ID)

	err := s.repo.CreateRequestLog(ctx, requestLog)
	if err != nil {
		log.Error().
			Err(err).
			Str("replay_id", replay.ID).
			Msg("failed to create request log")
		return "", err
	}

	return requestLog.ID, nil
}
