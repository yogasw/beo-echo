package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
)

const upstreamHTML = `<html><body>hello from upstream — ünïcödé ✓</body></html>`

func brotliBytes(t *testing.T, s string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := brotli.NewWriter(&buf)
	_, err := w.Write([]byte(s))
	require.NoError(t, err)
	require.NoError(t, w.Close())
	return buf.Bytes()
}

// upstream records the request headers it receives and answers with a Brotli-compressed HTML page,
// like a site behind Cloudflare does.
type upstream struct {
	server *httptest.Server
	mu     sync.Mutex
	got    http.Header
}

func newUpstream(t *testing.T) *upstream {
	u := &upstream{}
	body := brotliBytes(t, upstreamHTML)
	u.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u.mu.Lock()
		u.got = r.Header.Clone()
		u.mu.Unlock()
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.Header().Set("Content-Encoding", "br")
		w.Header().Set("Server", "cloudflare")
		w.Header().Set("Cf-Ray", "abc-SIN")
		_, _ = w.Write(body)
	}))
	t.Cleanup(u.server.Close)
	return u
}

func (u *upstream) received() http.Header {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.got
}

// proxyProject creates a proxy-mode project that forwards to the upstream.
func proxyProject(t *testing.T, alias string, u *upstream) *fixture {
	t.Helper()
	f := newFixture(t, alias, "")
	target := &database.ProxyTarget{ProjectID: f.project.ID, Label: "upstream", URL: u.server.URL}
	require.NoError(t, database.DB.Create(target).Error)
	require.NoError(t, database.DB.Model(&database.Project{}).Where("id = ?", f.project.ID).
		Updates(map[string]interface{}{"mode": string(database.ModeProxy), "active_proxy_id": target.ID}).Error)
	return f
}

func addScriptAction(t *testing.T, projectID string, point database.ExecutionPoint, script string) {
	t.Helper()
	cfg, err := json.Marshal(map[string]string{"script": script})
	require.NoError(t, err)
	action := &database.Action{
		ProjectID: projectID, Name: "it-script", Type: database.ActionTypeRunJavascript,
		ExecutionPoint: point, Enabled: true, Config: string(cfg),
	}
	require.NoError(t, database.DB.Create(action).Error)
}

func decodeBrotli(t *testing.T, raw []byte) string {
	t.Helper()
	out, err := io.ReadAll(brotli.NewReader(bytes.NewReader(raw)))
	require.NoError(t, err, "body must still be valid Brotli")
	return string(out)
}
