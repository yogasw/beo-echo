package services

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
)

const jsonHeaders = `{"Content-Type":"application/json"}`

func templatedResponse(body string) database.MockResponse {
	return database.MockResponse{StatusCode: 200, Body: body, Headers: jsonHeaders, Templated: true}
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return string(b)
}

func TestTemplate_QueryParamEcho(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/omniform-mock-stag/log?room_id=123&app_id=abc", nil)
	body := `[{"label":"room_id={{queryParam 'room_id'}}","value":"{{queryParam 'app_id'}}"}]`

	resp, err := buildMockResponse(templatedResponse(body), req)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, `[{"label":"room_id=123","value":"abc"}]`, readBody(t, resp))
}

func TestTemplate_RepeatManyItems(t *testing.T) {
	body := `[{{#repeat (queryParam 'n') comma=true}}{"label":"Opsi {{@index}}","value":"OPT-{{@index}}"}{{/repeat}}]`
	req := httptest.NewRequest(http.MethodGet, "/omniform-mock-stag/many?n=500", nil)

	resp, err := buildMockResponse(templatedResponse(body), req)
	require.NoError(t, err)

	var items []map[string]string
	require.NoError(t, json.Unmarshal([]byte(readBody(t, resp)), &items), "rendered body must be valid JSON")
	require.Len(t, items, 500)
	assert.Equal(t, map[string]string{"label": "Opsi 0", "value": "OPT-0"}, items[0])
	assert.Equal(t, map[string]string{"label": "Opsi 499", "value": "OPT-499"}, items[499])
}

func TestTemplate_RepeatEdgeCases(t *testing.T) {
	body := `[{{#repeat (queryParam 'n') comma=true}}{{@index}}{{/repeat}}]`
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{"zero", "n=0", "[]"},
		{"missing param", "", "[]"},
		{"negative", "n=-3", "[]"},
		{"not a number", "n=abc", "[]"},
		{"one", "n=1", "[0]"},
		{"three", "n=3", "[0,1,2]"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/p/x?"+tc.query, nil)
			got, err := renderResponseBody(templatedResponse(body), req)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTemplate_RepeatIsCapped(t *testing.T) {
	body := `[{{#repeat (queryParam 'n') comma=true}}1{{/repeat}}]`
	req := httptest.NewRequest(http.MethodGet, "/p/x?n=99999999", nil)

	got, err := renderResponseBody(templatedResponse(body), req)
	require.NoError(t, err)

	var items []int
	require.NoError(t, json.Unmarshal([]byte(got), &items))
	assert.Len(t, items, maxTemplateRepeat)
}

func TestTemplate_RepeatFirstLast(t *testing.T) {
	body := `{{#repeat 3}}{{#if @first}}F{{/if}}{{@index}}{{#if @last}}L{{/if}}{{/repeat}}`
	got, err := renderResponseBody(templatedResponse(body), httptest.NewRequest(http.MethodGet, "/p/x", nil))
	require.NoError(t, err)
	assert.Equal(t, "F012L", got)
}

func TestTemplate_QueryParamDefaultAndHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/p/x?present=1", nil)
	req.Header.Set("X-Api-Key", "secret-key")

	got, err := renderResponseBody(templatedResponse(`{{queryParam 'present' default='d'}}|{{queryParam 'absent' default='d'}}|{{queryParam 'absent'}}|{{header 'X-Api-Key'}}|{{header 'X-None' default='none'}}`), req)
	require.NoError(t, err)
	assert.Equal(t, "1|d||secret-key|none", got)
}

func TestTemplate_JSONOutputEscapesRequestValues(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, `/p/x?name=a%22b%5Cc`, nil)

	got, err := renderResponseBody(templatedResponse(`{"name":"{{queryParam 'name'}}"}`), req)
	require.NoError(t, err)

	var parsed map[string]string
	require.NoError(t, json.Unmarshal([]byte(got), &parsed), "a quote in the query must not break the JSON: %s", got)
	assert.Equal(t, `a"b\c`, parsed["name"])
}

func TestTemplate_NonJSONOutputIsRaw(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, `/p/x?name=a%22b`, nil)
	mock := database.MockResponse{Body: `{{queryParam 'name'}}`, Headers: `{"Content-Type":"text/plain"}`, Templated: true}

	got, err := renderResponseBody(mock, req)
	require.NoError(t, err)
	assert.Equal(t, `a"b`, got)
}

func TestTemplate_NotTemplatedStaysRaw(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/p/x?room_id=123", nil)
	mock := database.MockResponse{StatusCode: 200, Body: `{"x":"{{queryParam 'room_id'}}"}`, Headers: jsonHeaders}

	resp, err := buildMockResponse(mock, req)
	require.NoError(t, err)
	assert.Equal(t, `{"x":"{{queryParam 'room_id'}}"}`, readBody(t, resp))
}

func TestTemplate_SyntaxErrorReturns500(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/p/x", nil)

	resp, err := buildMockResponse(templatedResponse(`{{#repeat 3}}no close`), req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Contains(t, readBody(t, resp), "template")
}

func TestTemplate_UnknownHelperReturns500(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/p/x", nil)

	resp, err := buildMockResponse(templatedResponse(`{{noSuchHelper 'x'}}`), req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestTemplate_GzipAfterRender(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/p/x?room_id=123", nil)
	mock := templatedResponse(`{"room":"{{queryParam 'room_id'}}"}`)
	mock.Headers = `{"Content-Type":"application/json","Content-Encoding":"gzip"}`

	resp, err := buildMockResponse(mock, req)
	require.NoError(t, err)

	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))
	// Compressed bytes must decode to the rendered body, not the template source.
	zr, err := gzipReader(resp.Body)
	require.NoError(t, err)
	decoded, err := io.ReadAll(zr)
	require.NoError(t, err)
	assert.Equal(t, `{"room":"123"}`, string(decoded))
}

func gzipReader(r io.Reader) (io.Reader, error) {
	return gzip.NewReader(r)
}

func TestTemplate_JSONEscapeWithoutContentTypeHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, `/p/x?q=a%22b`, nil)
	mock := database.MockResponse{Body: `{"v":"{{queryParam 'q'}}"}`, Headers: `{}`, Templated: true}

	got, err := renderResponseBody(mock, req)
	require.NoError(t, err)

	var parsed map[string]string
	require.NoError(t, json.Unmarshal([]byte(got), &parsed), "body that looks like JSON is escaped even without a Content-Type: %s", got)
	assert.Equal(t, `a"b`, parsed["v"])
}

func TestTemplate_NoContentTypeAndPlainBodyIsRaw(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, `/p/x?q=a%22b`, nil)
	mock := database.MockResponse{Body: `hello {{queryParam 'q'}}`, Headers: `{}`, Templated: true}

	got, err := renderResponseBody(mock, req)
	require.NoError(t, err)
	assert.Equal(t, `hello a"b`, got)
}

func TestTemplate_NestedRepeatSharesOneBudget(t *testing.T) {
	body := `[{{#repeat 10000 comma=true}}{{#repeat 10000 comma=true}}1{{/repeat}}{{/repeat}}]`
	req := httptest.NewRequest(http.MethodGet, "/p/x", nil)

	got, err := renderResponseBody(templatedResponse(body), req)
	require.NoError(t, err)

	// The outer repeat spends the whole budget, so nothing is left for the inner ones.
	assert.LessOrEqual(t, len(got), 2*maxTemplateRepeat+2, "output stays bounded instead of 10000 x 10000")
}

func TestTemplate_RepeatCountParsing(t *testing.T) {
	assert.Equal(t, maxTemplateRepeat, toRepeatCount("99999999999999999999999"), "overflowing numbers count as very many")
	assert.Equal(t, 0, toRepeatCount("-99999999999999999999999"))
	assert.Equal(t, maxTemplateRepeat, toRepeatCount("500000"))
	assert.Equal(t, 500, toRepeatCount(" 500 "))
	assert.Equal(t, 7, toRepeatCount(7))
	assert.Equal(t, 0, toRepeatCount("abc"))
	assert.Equal(t, 0, toRepeatCount(""))
	assert.Equal(t, 0, toRepeatCount("1.5"))
}
