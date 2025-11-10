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

	"github.com/dop251/goja"
)

// RunJavascriptConfig represents the configuration for run_javascript action type
type RunJavascriptConfig struct {
	Script string `json:"script"` // JavaScript code to execute
}

// JavascriptContext represents the context passed to JavaScript execution
type JavascriptContext struct {
	Request  *JavascriptRequest  `json:"request,omitempty"`
	Response *JavascriptResponse `json:"response,omitempty"`
}

// JavascriptRequest represents the HTTP request in JavaScript context
type JavascriptRequest struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Query   map[string][]string `json:"query"`
	Headers map[string]string   `json:"headers"`
	Body    string              `json:"body"`
}

// JavascriptResponse represents the HTTP response in JavaScript context
type JavascriptResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// JavascriptResult represents the result from JavaScript execution
type JavascriptResult struct {
	Request  *JavascriptRequest  `json:"request,omitempty"`
	Response *JavascriptResponse `json:"response,omitempty"`
	Logs     []string            `json:"logs,omitempty"`
}

// ExecuteRunJavascriptAction executes a run_javascript action
func (m *ModulesAction) ExecuteRunJavascriptAction(action *database.Action, req *http.Request, resp *http.Response) error {
	var config RunJavascriptConfig
	if err := json.Unmarshal([]byte(action.Config), &config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Prepare context based on execution point
	ctx := &JavascriptContext{}

	// Always include request if available
	if req != nil {
		jsReq, err := convertHTTPRequestToJS(req)
		if err != nil {
			return fmt.Errorf("failed to convert request: %w", err)
		}
		ctx.Request = jsReq
	}

	// Include response only for after_request execution point
	if resp != nil && action.ExecutionPoint == database.ExecutionPointAfterRequest {
		jsResp, err := convertHTTPResponseToJS(resp)
		if err != nil {
			return fmt.Errorf("failed to convert response: %w", err)
		}
		ctx.Response = jsResp
	}

	// Execute JavaScript
	result, err := m.executeJavaScript(config.Script, ctx)
	if err != nil {
		return fmt.Errorf("script execution failed: %w", err)
	}

	// Apply modifications back to request/response
	if result.Request != nil && req != nil {
		if err := applyJSRequestToHTTP(result.Request, req); err != nil {
			return fmt.Errorf("failed to apply request modifications: %w", err)
		}
	}

	if result.Response != nil && resp != nil {
		if err := applyJSResponseToHTTP(result.Response, resp); err != nil {
			return fmt.Errorf("failed to apply response modifications: %w", err)
		}
	}

	return nil
}

// executeJavaScript executes JavaScript code in a sandboxed environment
func (m *ModulesAction) executeJavaScript(script string, ctx *JavascriptContext) (*JavascriptResult, error) {
	vm := goja.New()

	// Track console.log output
	logs := []string{}
	consoleObj := vm.NewObject()
	consoleObj.Set("log", func(call goja.FunctionCall) goja.Value {
		args := make([]string, len(call.Arguments))
		for i, arg := range call.Arguments {
			args[i] = arg.String()
		}
		logMsg := strings.Join(args, " ")
		logs = append(logs, logMsg)
		return goja.Undefined()
	})
	vm.Set("console", consoleObj)

	// Clone context to avoid modifying the original
	ctxJSON, err := json.Marshal(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context: %w", err)
	}

	// Set context in VM
	vm.Set("context", string(ctxJSON))

	// Create wrapper script that parses context and returns result
	wrapperScript := `
		(function() {
			// Parse the context
			var ctx = JSON.parse(context);

			// Make request and response available globally
			var request = ctx.request || null;
			var response = ctx.response || null;

			// Execute user script
			(function() {
				` + script + `
			})();

			// Return modified context
			return JSON.stringify({
				request: request,
				response: response
			});
		})();
	`

	// Execute script with timeout
	done := make(chan struct {
		result goja.Value
		err    error
	}, 1)

	go func() {
		val, err := vm.RunString(wrapperScript)
		done <- struct {
			result goja.Value
			err    error
		}{val, err}
	}()

	select {
	case res := <-done:
		if res.err != nil {
			return &JavascriptResult{Logs: logs}, fmt.Errorf("script error: %w", res.err)
		}

		// Parse result
		resultJSON := res.result.String()
		var result JavascriptResult
		if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
			return &JavascriptResult{Logs: logs}, fmt.Errorf("failed to parse script result: %w", err)
		}

		result.Logs = logs
		return &result, nil

	case <-time.After(5 * time.Second):
		vm.Interrupt("script timeout")
		return &JavascriptResult{Logs: logs}, errors.New("script execution timeout (5 seconds)")
	}
}

// ValidateRunJavascriptConfig validates the configuration for run_javascript actions
func (m *ModulesAction) ValidateRunJavascriptConfig(configJSON string) error {
	// TODO : Implement actual validation logic if needed
	return nil
}

// convertHTTPRequestToJS converts an http.Request to JavascriptRequest
func convertHTTPRequestToJS(req *http.Request) (*JavascriptRequest, error) {
	path := ""
	query := make(map[string][]string)
	if req.URL != nil {
		path = req.URL.Path
		query = req.URL.Query()
	}

	jsReq := &JavascriptRequest{
		Method:  req.Method,
		Path:    path,
		Query:   query,
		Headers: make(map[string]string),
	}

	// Convert headers to map
	for key, values := range req.Header {
		if len(values) > 0 {
			jsReq.Headers[key] = values[0]
		}
	}

	// Read body if present
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
		jsReq.Body = string(bodyBytes)

		// Restore body for further use
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return jsReq, nil
}

// convertHTTPResponseToJS converts an http.Response to JavascriptResponse
func convertHTTPResponseToJS(resp *http.Response) (*JavascriptResponse, error) {
	jsResp := &JavascriptResponse{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string]string),
	}

	// Convert headers to map
	for key, values := range resp.Header {
		if len(values) > 0 {
			jsResp.Headers[key] = values[0]
		}
	}

	// Read body if present
	if resp.Body != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()
		jsResp.Body = string(bodyBytes)

		// Restore body for further use
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return jsResp, nil
}

// applyJSRequestToHTTP applies JavaScript modifications back to http.Request
func applyJSRequestToHTTP(jsReq *JavascriptRequest, req *http.Request) error {
	// Update method
	if jsReq.Method != "" {
		req.Method = jsReq.Method
	}

	// Update path
	if jsReq.Path != "" {
		req.URL.Path = jsReq.Path
	}

	// Update headers
	if jsReq.Headers != nil {
		for key, value := range jsReq.Headers {
			req.Header.Set(key, value)
		}
	}

	// Update body
	if jsReq.Body != "" {
		bodyBytes := []byte(jsReq.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		req.ContentLength = int64(len(bodyBytes))
		req.Header.Set("Content-Length", strconv.Itoa(len(bodyBytes)))
	}

	return nil
}

// applyJSResponseToHTTP applies JavaScript modifications back to http.Response
func applyJSResponseToHTTP(jsResp *JavascriptResponse, resp *http.Response) error {
	// Update status code
	if jsResp.StatusCode > 0 {
		resp.StatusCode = jsResp.StatusCode
	}

	// Update headers
	if jsResp.Headers != nil {
		for key, value := range jsResp.Headers {
			resp.Header.Set(key, value)
		}
	}

	// Update body
	if jsResp.Body != "" {
		bodyBytes := []byte(jsResp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		resp.ContentLength = int64(len(bodyBytes))
		resp.Header.Set("Content-Length", strconv.Itoa(len(bodyBytes)))
	}

	return nil
}
