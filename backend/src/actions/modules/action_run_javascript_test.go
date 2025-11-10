package modules

import (
	"beo-echo/backend/src/database"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteRunJavascriptAction_ModifyRequestHeaders(t *testing.T) {
	m := NewActionModules()

	script := `if (request) { request.headers['X-Custom-Header'] = 'TestValue'; }`

	configJSON := fmt.Sprintf(`{"script":"%s"}`, script)

	action := &database.Action{
		Type:           database.ActionTypeRunJavascript,
		ExecutionPoint: database.ExecutionPointBeforeRequest,
		Config:         configJSON,
	}

	req := &http.Request{
		Method: "GET",
		Header: make(http.Header),
		URL:    &url.URL{Path: "/test"},
	}
	req.Header.Set("Content-Type", "application/json")

	err := m.ExecuteRunJavascriptAction(action, req, nil)
	assert.NoError(t, err)
	assert.Equal(t, "TestValue", req.Header.Get("X-Custom-Header"))
}

func TestExecuteRunJavascriptAction_ModifyRequestBody(t *testing.T) {
	m := NewActionModules()

	script := `if (request && request.body) { var data = JSON.parse(request.body); data.modified = true; request.body = JSON.stringify(data); }`

	configJSON := fmt.Sprintf(`{"script":"%s"}`, script)

	action := &database.Action{
		Type:           database.ActionTypeRunJavascript,
		ExecutionPoint: database.ExecutionPointBeforeRequest,
		Config:         configJSON,
	}

	originalBody := `{"name":"test"}`
	req := &http.Request{
		Method: "POST",
		Header: make(http.Header),
		URL:    &url.URL{Path: "/test"},
		Body:   io.NopCloser(bytes.NewBufferString(originalBody)),
	}

	err := m.ExecuteRunJavascriptAction(action, req, nil)
	assert.NoError(t, err)

	// Read modified body
	bodyBytes, err := io.ReadAll(req.Body)
	require.NoError(t, err)
	bodyStr := string(bodyBytes)

	assert.Contains(t, bodyStr, `"modified":true`)
	assert.Contains(t, bodyStr, `"name":"test"`)
}

func TestExecuteRunJavascriptAction_ModifyResponseStatus(t *testing.T) {
	m := NewActionModules()

	script := `if (response) { response.status_code = 418; }`

	configJSON := fmt.Sprintf(`{"script":"%s"}`, script)

	action := &database.Action{
		Type:           database.ActionTypeRunJavascript,
		ExecutionPoint: database.ExecutionPointAfterRequest,
		Config:         configJSON,
	}

	req := &http.Request{
		Method: "GET",
		Header: make(http.Header),
		URL:    &url.URL{Path: "/test"},
	}

	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString("")),
	}

	err := m.ExecuteRunJavascriptAction(action, req, resp)
	assert.NoError(t, err)
	assert.Equal(t, 418, resp.StatusCode)
}

func TestExecuteRunJavascriptAction_BeforeRequest_NoResponse(t *testing.T) {
	m := NewActionModules()

	script := `if (request) { console.log('Processing request'); }`

	configJSON := fmt.Sprintf(`{"script":"%s"}`, script)

	action := &database.Action{
		Type:           database.ActionTypeRunJavascript,
		ExecutionPoint: database.ExecutionPointBeforeRequest,
		Config:         configJSON,
	}

	req := &http.Request{
		Method: "GET",
		Header: make(http.Header),
		URL:    &url.URL{Path: "/test"},
	}

	err := m.ExecuteRunJavascriptAction(action, req, nil)
	assert.NoError(t, err)
}

func TestExecuteRunJavascriptAction_ScriptError(t *testing.T) {
	m := NewActionModules()

	script := `throw new Error('Script error');`

	configJSON := fmt.Sprintf(`{"script":"%s"}`, script)

	action := &database.Action{
		Type:           database.ActionTypeRunJavascript,
		ExecutionPoint: database.ExecutionPointBeforeRequest,
		Config:         configJSON,
	}

	req := &http.Request{
		Method: "GET",
		Header: make(http.Header),
	}

	err := m.ExecuteRunJavascriptAction(action, req, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "script execution failed")
}
