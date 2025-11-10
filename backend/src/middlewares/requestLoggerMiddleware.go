package middlewares

import (
	"beo-echo/backend/src/auth"
	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/handler"
	handlerLogs "beo-echo/backend/src/logs/handlers"
	systemConfig "beo-echo/backend/src/systemConfigs"
	"beo-echo/backend/src/utils"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

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
		l := log.With().Str("func", "middlewares.RequestLoggerMiddleware").Logger()
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
		path, _ := c.Get(handler.KeyPath)

		// If no project ID, skip logging
		if projectID == nil || projectID == "" || path == nil {
			return
		}

		var id string
		requestId := c.GetString("request_id")
		if requestId != "" {
			id = requestId
		} else {
			id = uuid.New().String()
		}

		// Get response body - prefer raw response from handler context
		var responseBody string
		// Check if response is compressed and decompress if needed
		responseBody = respBodyBuf.String()
		contentEncoding := c.Writer.Header().Get("Content-Encoding")

		if contentEncoding != "" {
			compressedData := respBodyBuf.Bytes()

			switch strings.ToLower(contentEncoding) {
			case "gzip":
				if reader, err := gzip.NewReader(bytes.NewReader(compressedData)); err == nil {
					if decompressed, err := io.ReadAll(reader); err == nil {
						responseBody = string(decompressed)
					}
					reader.Close()
				}
			case "br":
				decompressed := brotli.NewReader(bytes.NewReader(compressedData))
				if data, err := io.ReadAll(decompressed); err == nil {
					responseBody = string(data)
				}
			}
		}

		// Save log with hashed JWTs
		logEntry := &database.RequestLog{
			ID:              id,
			ProjectID:       toString(projectID),
			Method:          c.Request.Method,
			Path:            toString(path),
			QueryParams:     c.Request.URL.RawQuery,
			RequestHeaders:  MapSliceToJSONJoined(c.Request.Header),
			RequestBody:     requestBody,
			ResponseStatus:  c.Writer.Status(),
			ResponseHeaders: MapSliceToJSONJoined(c.Writer.Header()),
			ResponseBody:    responseBody, // Use the responseBody variable
			LatencyMS:       int(latency),
			ExecutionMode:   database.ProjectMode(toString(executionMode)),
			Matched:         toBool(matched),
			CreatedAt:       time.Now(),
		}

		entry, err := json.Marshal(logEntry)
		if err != nil {
			l.Error().Err(err).Msg("Failed to marshal log entry to JSON")
		} else {
			md5Hash := utils.HashMD5(string(entry))
			logHash, errJwt := auth.GenerateJWTFromString(md5Hash)
			if errJwt != nil {
				l.Error().Err(errJwt).Msg("Failed to generate JWT from MD5 hash")
			}
			logEntry.LogsHash = logHash
		}

		handlerLogs.EnsureLogService()
		if ls := handlerLogs.LogService(); ls != nil {
			ls.NotifySubscribers(*logEntry)
		}

		// Check if auto-save is enabled
		autoSaveEnabled, err := systemConfig.GetSystemConfigWithType[bool](systemConfig.AUTO_SAVE_LOGS_IN_DB_ENABLED)
		if err != nil {
			l.Error().Err(err).Msg("Failed to get AUTO_SAVE_LOGS_IN_DB_ENABLED config")
			return
		}

		if autoSaveEnabled {
			// Save to database
			if err := db.Create(logEntry).Error; err != nil {
				// Log error if saving to DB fails
				l.Error().Err(err).
					Str("project_id", toString(projectID)).
					Msg("Failed to save request log to database")
			}
		}
	}
}

func MapSliceToJSONJoined(m map[string][]string) string {
	flat := make(map[string]string, len(m))
	for key, values := range m {
		flat[key] = strings.Join(values, "; ")
	}
	jsonBytes, err := json.Marshal(flat)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return "{}"
	}
	return string(jsonBytes)
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
