package modules

import (
	"beo-echo/backend/src/database"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ReplaceTextConfig represents the configuration for replace_text action type
type ReplaceTextConfig struct {
	Target      string `json:"target"`      // request_body, response_body, request_header, response_header
	Pattern     string `json:"pattern"`     // String or regex pattern to find
	Replacement string `json:"replacement"` // Replacement text
	UseRegex    bool   `json:"use_regex"`   // Whether to use regex matching
	HeaderKey   string `json:"header_key"`  // Header key (only for header targets)
}

// executeReplaceTextAction executes a replace_text action
func (m *ModulesAction) ExecuteReplaceTextAction(action *database.Action, req *http.Request, resp *http.Response) error {
	var config ReplaceTextConfig
	if err := json.Unmarshal([]byte(action.Config), &config); err != nil {
		return err
	}

	// Perform replacement based on target
	switch config.Target {
	case "request_body":
		if req != nil && req.Body != nil {
			// Read body
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return err
			}
			req.Body.Close()

			// Replace text
			bodyStr := string(bodyBytes)
			bodyStr = m.replaceText(bodyStr, config.Pattern, config.Replacement, config.UseRegex)

			// Set new body
			req.Body = io.NopCloser(bytes.NewBufferString(bodyStr))
			req.ContentLength = int64(len(bodyStr))
			req.Header.Set("Content-Length", strconv.Itoa(len(bodyStr)))
		}
	case "response_body":
		if resp != nil && resp.Body != nil {
			// Read body
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			resp.Body.Close()

			// Replace text
			bodyStr := string(bodyBytes)
			bodyStr = m.replaceText(bodyStr, config.Pattern, config.Replacement, config.UseRegex)

			// Set new body
			resp.Body = io.NopCloser(bytes.NewBufferString(bodyStr))
			resp.ContentLength = int64(len(bodyStr))
			resp.Header.Set("Content-Length", strconv.Itoa(len(bodyStr)))
		}
	case "request_header":
		if req != nil && req.Header != nil {
			value := req.Header.Get(config.HeaderKey)
			if value != "" {
				req.Header.Set(config.HeaderKey, m.replaceText(value, config.Pattern, config.Replacement, config.UseRegex))
			}
		}
	case "response_header":
		if resp != nil && resp.Header != nil {
			value := resp.Header.Get(config.HeaderKey)
			if value != "" {
				resp.Header.Set(config.HeaderKey, m.replaceText(value, config.Pattern, config.Replacement, config.UseRegex))
			}
		}
	}

	return nil
}

// replaceText performs the actual text replacement
func (m *ModulesAction) replaceText(text, pattern, replacement string, useRegex bool) string {
	if useRegex {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return text // Return original on error
		}
		return re.ReplaceAllString(text, replacement)
	}
	return strings.ReplaceAll(text, pattern, replacement)
}

// validateReplaceTextConfig validates the configuration for replace_text actions
func (m *ModulesAction) ValidateReplaceTextConfig(configJSON string) error {
	if configJSON == "" {
		return errors.New("config is required for replace_text action")
	}

	var config ReplaceTextConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return errors.New("invalid config JSON: " + err.Error())
	}

	// Validate target
	validTargets := map[string]bool{
		"request_body":    true,
		"response_body":   true,
		"request_header":  true,
		"response_header": true,
	}
	if !validTargets[config.Target] {
		return errors.New("invalid target: must be request_body, response_body, request_header, or response_header")
	}

	// Validate header key for header targets
	if (config.Target == "request_header" || config.Target == "response_header") && config.HeaderKey == "" {
		return errors.New("header_key is required for header targets")
	}

	// Validate pattern
	if config.Pattern == "" {
		return errors.New("pattern is required")
	}

	// Validate regex if used
	if config.UseRegex {
		if _, err := regexp.Compile(config.Pattern); err != nil {
			return errors.New("invalid regex pattern: " + err.Error())
		}
	}

	return nil
}
