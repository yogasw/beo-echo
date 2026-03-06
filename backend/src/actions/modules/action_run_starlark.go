package modules

import (
	"beo-echo/backend/src/database"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.starlark.net/starlark"
	starlarkjson "go.starlark.net/lib/json"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
)

// RunStarlarkConfig represents the configuration for run_starlark action type
type RunStarlarkConfig struct {
	Script string `json:"script"` // Starlark code to execute
}

// StarlarkContext represents the context passed to Starlark execution
type StarlarkContext struct {
	Request  *StarlarkRequest  `json:"request,omitempty"`
	Response *StarlarkResponse `json:"response,omitempty"`
}

// StarlarkRequest represents the HTTP request in Starlark context
type StarlarkRequest struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Query   map[string][]string `json:"query"`
	Headers map[string]string   `json:"headers"`
	Body    string              `json:"body"`
}

// StarlarkResponse represents the HTTP response in Starlark context
type StarlarkResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// StarlarkResult represents the result from Starlark execution
type StarlarkResult struct {
	Request  *StarlarkRequest  `json:"request,omitempty"`
	Response *StarlarkResponse `json:"response,omitempty"`
	Logs     []string          `json:"logs,omitempty"`
}

// ExecuteRunStarlarkAction executes a run_starlark action
func (m *ModulesAction) ExecuteRunStarlarkAction(action *database.Action, req *http.Request, resp *http.Response) error {
	var config RunStarlarkConfig
	if err := json.Unmarshal([]byte(action.Config), &config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Prepare context based on execution point
	ctx := &StarlarkContext{}

	// Always include request if available
	if req != nil {
		starlarkReq, err := convertHTTPRequestToStarlark(req)
		if err != nil {
			return fmt.Errorf("failed to convert request: %w", err)
		}
		ctx.Request = starlarkReq
	}

	// Include response only for after_request execution point
	if resp != nil && action.ExecutionPoint == database.ExecutionPointAfterRequest {
		starlarkResp, err := convertHTTPResponseToStarlark(resp)
		if err != nil {
			return fmt.Errorf("failed to convert response: %w", err)
		}
		ctx.Response = starlarkResp
	}

	// Execute Starlark
	result, err := m.executeStarlark(config.Script, ctx)
	if err != nil {
		return fmt.Errorf("script execution failed: %w", err)
	}

	// Apply modifications back to request/response
	if result.Request != nil && req != nil {
		if err := applyStarlarkRequestToHTTP(result.Request, req); err != nil {
			return fmt.Errorf("failed to apply request modifications: %w", err)
		}
	}

	if result.Response != nil && resp != nil {
		if err := applyStarlarkResponseToHTTP(result.Response, resp); err != nil {
			return fmt.Errorf("failed to apply response modifications: %w", err)
		}
	}

	return nil
}

// executeStarlark executes Starlark code in a sandboxed environment
func (m *ModulesAction) executeStarlark(script string, ctx *StarlarkContext) (*StarlarkResult, error) {
	// Track console output
	logs := []string{}

	// Create print function that captures output
	printFunc := func(thread *starlark.Thread, msg string) {
		logs = append(logs, msg)
	}

	// Convert context to JSON for Starlark
	ctxJSON, err := json.Marshal(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context: %w", err)
	}

	// Create wrapper script that provides context and captures modifications
	wrapperScript := `
# Parse context
import json
ctx = json.decode(context_json)

# Make request and response available globally
request = ctx.get("request")
response = ctx.get("response")

# User script execution
def run_user_script():
	` + indentScript(script) + `

run_user_script()

# Return modified context
result = json.encode({
	"request": request,
	"response": response
})
`

	// Execute with timeout
	done := make(chan struct {
		result string
		err    error
	}, 1)

	go func() {
		thread := &starlark.Thread{
			Print: printFunc,
		}

		// Create predeclared environment with built-ins
		predeclared := starlark.StringDict{
			"json":         starlarkjson.Module,
			"context_json": starlark.String(string(ctxJSON)),
			"struct":       starlark.NewBuiltin("struct", starlarkstruct.Make),
		}

		// Execute the script with options
		opts := &syntax.FileOptions{
			Set:               true,
			While:             true,
			TopLevelControl:   true,
			GlobalReassign:    true,
			LoadBindsGlobally: true,
			Recursion:         true,
		}
		globals, err := starlark.ExecFileOptions(opts, thread, "script.star", wrapperScript, predeclared)
		if err != nil {
			done <- struct {
				result string
				err    error
			}{"", fmt.Errorf("starlark execution error: %w", err)}
			return
		}

		// Get result
		resultVal, ok := globals["result"]
		if !ok {
			done <- struct {
				result string
				err    error
			}{"", errors.New("script did not produce result")}
			return
		}

		resultStr, ok := resultVal.(starlark.String)
		if !ok {
			done <- struct {
				result string
				err    error
			}{"", errors.New("result is not a string")}
			return
		}

		done <- struct {
			result string
			err    error
		}{string(resultStr), nil}
	}()

	select {
	case res := <-done:
		if res.err != nil {
			return &StarlarkResult{Logs: logs}, res.err
		}

		// Parse result
		var result StarlarkResult
		if err := json.Unmarshal([]byte(res.result), &result); err != nil {
			return &StarlarkResult{Logs: logs}, fmt.Errorf("failed to parse script result: %w", err)
		}

		result.Logs = logs
		return &result, nil

	case <-time.After(5 * time.Second):
		return &StarlarkResult{Logs: logs}, errors.New("script execution timeout (5 seconds)")
	}
}

// indentScript indents user script for proper Python-like syntax
func indentScript(script string) string {
	lines := strings.Split(script, "\n")
	indented := make([]string, len(lines))
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			indented[i] = "\t" + line
		} else {
			indented[i] = line
		}
	}
	return strings.Join(indented, "\n")
}

// ValidateRunStarlarkConfig validates the configuration for run_starlark actions
func (m *ModulesAction) ValidateRunStarlarkConfig(configJSON string) error {
	var config RunStarlarkConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("invalid JSON config: %w", err)
	}

	if config.Script == "" {
		return errors.New("script is required")
	}

	return nil
}

// convertHTTPRequestToStarlark converts an http.Request to StarlarkRequest
func convertHTTPRequestToStarlark(req *http.Request) (*StarlarkRequest, error) {
	path := ""
	query := make(map[string][]string)
	if req.URL != nil {
		path = req.URL.Path
		query = req.URL.Query()
	}

	starlarkReq := &StarlarkRequest{
		Method:  req.Method,
		Path:    path,
		Query:   query,
		Headers: make(map[string]string),
	}

	// Convert headers to map
	for key, values := range req.Header {
		if len(values) > 0 {
			starlarkReq.Headers[key] = values[0]
		}
	}

	// Read body if present
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
		starlarkReq.Body = string(bodyBytes)

		// Restore body for further use
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return starlarkReq, nil
}

// convertHTTPResponseToStarlark converts an http.Response to StarlarkResponse
func convertHTTPResponseToStarlark(resp *http.Response) (*StarlarkResponse, error) {
	starlarkResp := &StarlarkResponse{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string]string),
	}

	// Convert headers to map
	for key, values := range resp.Header {
		if len(values) > 0 {
			starlarkResp.Headers[key] = values[0]
		}
	}

	// Read body if present
	if resp.Body != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()
		starlarkResp.Body = string(bodyBytes)

		// Restore body for further use
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return starlarkResp, nil
}

// applyStarlarkRequestToHTTP applies Starlark modifications back to http.Request
func applyStarlarkRequestToHTTP(starlarkReq *StarlarkRequest, req *http.Request) error {
	// Update method
	if starlarkReq.Method != "" {
		req.Method = starlarkReq.Method
	}

	// Update path
	if starlarkReq.Path != "" {
		req.URL.Path = starlarkReq.Path
	}

	// Update headers
	if starlarkReq.Headers != nil {
		for key, value := range starlarkReq.Headers {
			req.Header.Set(key, value)
		}
	}

	// Update body
	if starlarkReq.Body != "" {
		bodyBytes := []byte(starlarkReq.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		req.ContentLength = int64(len(bodyBytes))
		req.Header.Set("Content-Length", strconv.Itoa(len(bodyBytes)))
	}

	return nil
}

// applyStarlarkResponseToHTTP applies Starlark modifications back to http.Response
func applyStarlarkResponseToHTTP(starlarkResp *StarlarkResponse, resp *http.Response) error {
	// Update status code
	if starlarkResp.StatusCode > 0 {
		resp.StatusCode = starlarkResp.StatusCode
	}

	// Update headers
	if starlarkResp.Headers != nil {
		for key, value := range starlarkResp.Headers {
			resp.Header.Set(key, value)
		}
	}

	// Update body
	if starlarkResp.Body != "" {
		bodyBytes := []byte(starlarkResp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		resp.ContentLength = int64(len(bodyBytes))
		resp.Header.Set("Content-Length", strconv.Itoa(len(bodyBytes)))
	}

	return nil
}
