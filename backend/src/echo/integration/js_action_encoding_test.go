package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"beo-echo/backend/src/database"
)

// An after_request script that only reads the response must not change a compressed body.
func TestIntegration_AfterRequestScriptKeepsCompressedBodyIntact(t *testing.T) {
	u := newUpstream(t)
	f := proxyProject(t, "it-proxy-after", u)
	defer f.cleanup()
	addScriptAction(t, f.project.ID, database.ExecutionPointAfterRequest, `console.log("status " + response.status_code);`)

	w := do(f.router, http.MethodGet, "/it-proxy-after/ping", map[string]string{"X-Api-Key": "secret"}, "")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, upstreamHTML, decodeBrotli(t, w.Body.Bytes()), "a read-only script must not corrupt the Brotli body")
}
