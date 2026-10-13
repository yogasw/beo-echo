package modules

import (
	"beo-echo/backend/src/database"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const encodedPage = `<html><body>hello — ünïcödé ✓ {"n":1}</body></html>`

func compress(t *testing.T, encoding, s string) []byte {
	t.Helper()
	var buf bytes.Buffer
	var w io.WriteCloser
	switch encoding {
	case "gzip":
		w = gzip.NewWriter(&buf)
	case "br":
		w = brotli.NewWriter(&buf)
	case "deflate":
		w = zlib.NewWriter(&buf)
	case "raw-deflate":
		fw, err := flate.NewWriter(&buf, flate.DefaultCompression)
		require.NoError(t, err)
		w = fw
	default:
		t.Fatalf("unknown encoding %s", encoding)
	}
	_, err := w.Write([]byte(s))
	require.NoError(t, err)
	require.NoError(t, w.Close())
	return buf.Bytes()
}

func encodedResponse(t *testing.T, encoding, body string) (*http.Response, []byte) {
	t.Helper()
	headerEncoding := encoding
	if encoding == "raw-deflate" {
		headerEncoding = "deflate"
	}
	raw := compress(t, encoding, body)
	h := make(http.Header)
	h.Set("Content-Type", "text/html; charset=UTF-8")
	h.Set("Content-Encoding", headerEncoding)
	h.Set("Content-Length", strconv.Itoa(len(raw)))
	return &http.Response{StatusCode: 200, Header: h, Body: io.NopCloser(bytes.NewReader(raw)), ContentLength: int64(len(raw))}, raw
}

func runAfterRequest(t *testing.T, script string, resp *http.Response) {
	t.Helper()
	cfg := fmt.Sprintf(`{"script":%q}`, script)
	action := &database.Action{Type: database.ActionTypeRunJavascript, ExecutionPoint: database.ExecutionPointAfterRequest, Config: cfg}
	req := &http.Request{Method: "GET", Header: make(http.Header), URL: &url.URL{Path: "/p"}}
	require.NoError(t, NewActionModules().ExecuteRunJavascriptAction(action, req, resp))
}

func TestRunJavascript_ReadOnlyScriptKeepsCompressedBodyIntact(t *testing.T) {
	for _, encoding := range []string{"br", "gzip", "deflate", "raw-deflate"} {
		t.Run(encoding, func(t *testing.T) {
			resp, raw := encodedResponse(t, encoding, encodedPage)

			runAfterRequest(t, `console.log("length " + response.body.length);`, resp)

			got, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, raw, got, "the compressed bytes must be exactly what the upstream sent")
			assert.NotEmpty(t, resp.Header.Get("Content-Encoding"), "Content-Encoding stays: the body is still encoded")
		})
	}
}

func TestRunJavascript_ScriptSeesDecodedBodyAndCanModifyIt(t *testing.T) {
	for _, encoding := range []string{"br", "gzip", "deflate"} {
		t.Run(encoding, func(t *testing.T) {
			resp, _ := encodedResponse(t, encoding, `{"name":"upstream"}`)

			runAfterRequest(t, `var b = JSON.parse(response.body); b.modified = true; response.body = JSON.stringify(b);`, resp)

			got, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.JSONEq(t, `{"name":"upstream","modified":true}`, string(got), "the script worked on the decoded JSON")
			assert.Empty(t, resp.Header.Get("Content-Encoding"), "the new body is plain, so Content-Encoding is dropped")
			assert.Equal(t, strconv.Itoa(len(got)), resp.Header.Get("Content-Length"))
		})
	}
}

func TestRunJavascript_PlainBodyBehaviourUnchanged(t *testing.T) {
	h := make(http.Header)
	h.Set("Content-Type", "application/json")
	resp := &http.Response{StatusCode: 200, Header: h, Body: io.NopCloser(bytes.NewBufferString(`{"a":1}`))}

	runAfterRequest(t, `var b = JSON.parse(response.body); b.b = 2; response.body = JSON.stringify(b);`, resp)

	got, _ := io.ReadAll(resp.Body)
	assert.JSONEq(t, `{"a":1,"b":2}`, string(got))
}

func TestDecodeContentEncoding(t *testing.T) {
	_, ok := decodeContentEncoding("", []byte("x"))
	assert.False(t, ok, "no encoding")
	_, ok = decodeContentEncoding("zstd", []byte("x"))
	assert.False(t, ok, "unsupported encoding keeps raw bytes")
	_, ok = decodeContentEncoding("gzip", []byte("not gzip"))
	assert.False(t, ok, "undecodable bytes keep raw bytes")

	out, ok := decodeContentEncoding(" GZIP ", compress(t, "gzip", "hi"))
	assert.True(t, ok)
	assert.Equal(t, "hi", string(out))
}
