package middlewares

import (
	"bytes"
	"io"
	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/mocks/handler"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// bodyWriter wraps gin.ResponseWriter to capture the response body
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)                  // Copy to buffer
	return w.ResponseWriter.Write(b) // Write as usual
}

// RequestLoggerMiddleware logs each HTTP request and response
func RequestLoggerMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Clone request body
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Wrap response writer
		respBodyBuf := new(bytes.Buffer)
		bw := &bodyWriter{body: respBodyBuf, ResponseWriter: c.Writer}
		c.Writer = bw

		// Continue
		c.Next()

		// Calculate latency
		latency := time.Since(start).Milliseconds()

		// Extract context vars (must be set by handler)
		projectID, _ := c.Get(handler.KeyProjectID)
		executionMode, _ := c.Get(handler.KeyExecutionMode)
		matched, _ := c.Get(handler.KeyMatched)

		// If no project ID, skip logging
		if projectID == nil || projectID == "" {
			return
		}

		// Save log
		logEntry := &database.RequestLog{
			ProjectID:       toString(projectID),
			Method:          c.Request.Method,
			Path:            c.Request.URL.Path,
			QueryParams:     c.Request.URL.RawQuery,
			RequestHeaders:  headersToJSON(c.Request.Header),
			RequestBody:     requestBody,
			ResponseStatus:  c.Writer.Status(),
			ResponseHeaders: headersToJSON(c.Writer.Header()),
			ResponseBody:    respBodyBuf.String(),
			LatencyMS:       int(latency),
			ExecutionMode:   database.ProjectMode(toString(executionMode)),
			Matched:         toBool(matched),
		}

		// Save to database
		if err := db.Create(logEntry).Error; err == nil {
			// Notify log subscribers if log service is available
			handler.EnsureLogService()
			if ls := handler.LogService(); ls != nil {
				ls.NotifySubscribers(*logEntry)
			}
		}
	}
}

// Helper: Convert headers to JSON-like string (simple)
func headersToJSON(h map[string][]string) string {
	s := "{"
	for k, v := range h {
		s += `"` + k + `":"` + v[0] + `",`
	}
	if len(h) > 0 {
		s = s[:len(s)-1]
	}
	return s + "}"
}

func toString(v interface{}) string {
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

func toBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
