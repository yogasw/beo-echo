package services

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
)

func TestCreateMockResponse_NoCompression(t *testing.T) {
	// Given - Mock response with no Content-Encoding header
	mockResp := database.MockResponse{
		StatusCode: 200,
		Body:       `{"message": "Hello World"}`,
		Headers:    `{"Content-Type": "application/json"}`,
	}

	// When - Create HTTP response
	resp, err := createMockResponse(mockResp)

	// Then - Response should be created without compression
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, int64(len(mockResp.Body)), resp.ContentLength)

	// Read and verify body
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, mockResp.Body, string(bodyBytes))
}

func TestCreateMockResponse_GzipCompression(t *testing.T) {
	// Given - Mock response with gzip Content-Encoding
	mockResp := database.MockResponse{
		StatusCode: 200,
		Body:       `{"message": "Hello World"}`,
		Headers:    `{"Content-Type": "application/json", "Content-Encoding": "gzip"}`,
	}

	// When - Create HTTP response
	resp, err := createMockResponse(mockResp)

	// Then - Response should be gzip compressed
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))

	// Verify compressed body can be decompressed
	gzipReader, err := gzip.NewReader(resp.Body)
	require.NoError(t, err)
	defer gzipReader.Close()

	decompressedBytes, err := io.ReadAll(gzipReader)
	require.NoError(t, err)
	assert.Equal(t, mockResp.Body, string(decompressedBytes))

	// Content length should be different (compressed size)
	assert.NotEqual(t, int64(len(mockResp.Body)), resp.ContentLength)
}

func TestCreateMockResponse_BrotliCompression(t *testing.T) {
	// Given - Mock response with brotli Content-Encoding
	mockResp := database.MockResponse{
		StatusCode: 200,
		Body:       `{"message": "Hello World"}`,
		Headers:    `{"Content-Type": "application/json", "Content-Encoding": "br"}`,
	}

	// When - Create HTTP response
	resp, err := createMockResponse(mockResp)

	// Then - Response should be brotli compressed
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, "br", resp.Header.Get("Content-Encoding"))

	// Verify compressed body can be decompressed
	brotliReader := brotli.NewReader(resp.Body)
	decompressedBytes, err := io.ReadAll(brotliReader)
	require.NoError(t, err)
	assert.Equal(t, mockResp.Body, string(decompressedBytes))

	// Content length should be different (compressed size)
	assert.NotEqual(t, int64(len(mockResp.Body)), resp.ContentLength)
}

func TestCreateMockResponse_CaseInsensitiveHeaders(t *testing.T) {
	// Given - Mock response with case-insensitive Content-Encoding header
	testCases := []struct {
		name           string
		headerKey      string
		headerValue    string
		expectedLength bool // true if length should be different (compressed)
	}{
		{"lowercase", "content-encoding", "gzip", true},
		{"uppercase", "CONTENT-ENCODING", "gzip", true},
		{"mixed case", "Content-Encoding", "gzip", true},
		{"camel case", "Content-encoding", "gzip", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create headers with different case variations
			headers := map[string]string{
				"Content-Type": "application/json",
				tc.headerKey:   tc.headerValue,
			}
			headersJSON, _ := json.Marshal(headers)

			mockResp := database.MockResponse{
				StatusCode: 200,
				Body:       `{"message": "Hello World"}`,
				Headers:    string(headersJSON),
			}

			// When - Create HTTP response
			resp, err := createMockResponse(mockResp)

			// Then - Should handle case-insensitive headers correctly
			require.NoError(t, err)
			assert.Equal(t, tc.headerValue, resp.Header.Get("Content-Encoding"))

			if tc.expectedLength {
				assert.NotEqual(t, int64(len(mockResp.Body)), resp.ContentLength)
			} else {
				assert.Equal(t, int64(len(mockResp.Body)), resp.ContentLength)
			}
		})
	}
}

func TestCreateMockResponse_InvalidHeaders(t *testing.T) {
	// Given - Mock response with invalid JSON headers
	mockResp := database.MockResponse{
		StatusCode: 200,
		Body:       `{"message": "Hello World"}`,
		Headers:    `invalid json`,
	}

	// When - Create HTTP response
	resp, err := createMockResponse(mockResp)

	// Then - Should not fail and use raw body (no compression)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, int64(len(mockResp.Body)), resp.ContentLength)

	// Read and verify body
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, mockResp.Body, string(bodyBytes))
}

func TestCreateMockResponse_UnsupportedCompression(t *testing.T) {
	// Given - Mock response with unsupported compression type
	mockResp := database.MockResponse{
		StatusCode: 200,
		Body:       `{"message": "Hello World"}`,
		Headers:    `{"Content-Type": "application/json", "Content-Encoding": "deflate"}`,
	}

	// When - Create HTTP response
	resp, err := createMockResponse(mockResp)

	// Then - Should use raw body (no compression)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "deflate", resp.Header.Get("Content-Encoding"))
	assert.Equal(t, int64(len(mockResp.Body)), resp.ContentLength)

	// Read and verify body is not compressed
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, mockResp.Body, string(bodyBytes))
}

func TestCreateMockResponse_EmptyBody(t *testing.T) {
	// Given - Mock response with empty body and compression
	mockResp := database.MockResponse{
		StatusCode: 204,
		Body:       "",
		Headers:    `{"Content-Encoding": "gzip"}`,
	}

	// When - Create HTTP response
	resp, err := createMockResponse(mockResp)

	// Then - Should handle empty body compression correctly
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))

	// Verify compressed empty body can be decompressed
	gzipReader, err := gzip.NewReader(resp.Body)
	require.NoError(t, err)
	defer gzipReader.Close()

	decompressedBytes, err := io.ReadAll(gzipReader)
	require.NoError(t, err)
	assert.Equal(t, "", string(decompressedBytes))
}
