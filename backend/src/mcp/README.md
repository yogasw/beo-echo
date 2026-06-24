# Beo Echo MCP Server

Exposes Beo Echo over the [Model Context Protocol](https://modelcontextprotocol.io)
so AI clients (Claude.ai, Claude Desktop, Cursor, VS Code, Claude Code) can drive
your mocks, routes, logs, replays, actions, config, and workspaces with tools.

## TODO / status

- [x] PAT auth (`beo_pat_…`) + profile endpoints
- [x] OAuth authorization-code flow (for "paste URL" clients like Claude.ai)
- [x] `/mcp` streamable HTTP endpoint (mounted on the main server)
- [x] Tools for all 7 areas (51 tools)
- [ ] Frontend: token management UI on the profile page + `/oauth/consent` route
- [ ] Optional: streaming logs tool over SSE

## The endpoint

The MCP endpoint is served at **`/mcp`** on the same host/port as the API:

```
https://<your-beo-echo-host>/mcp
```

It is the URL every MCP client needs. Both auth modes below use it.

## Auth — two modes

### 1. Static bearer token (Claude Desktop, Cursor, VS Code, curl)

1. Open the web app → **Profile → Access Tokens → Generate**.
2. Copy the token (`beo_pat_…`) — it is shown once.
3. Put it in your client config as an `Authorization: Bearer` header.

Tokens are managed via the API:

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/users/me/tokens` | Create a token (`{"name": "...", "expires_days": 0}`; `0` = never expires). Returns the plaintext **once**. |
| `GET`  | `/api/users/me/tokens` | List your tokens (metadata only). |
| `DELETE` | `/api/users/me/tokens/:tokenId` | Revoke a token. |

### 2. OAuth (Claude.ai and other OAuth-aware clients)

Paste just the endpoint URL — the client discovers the OAuth server via
`/.well-known/oauth-authorization-server`, sends you to Beo Echo to log in and
**Approve**, then receives a token automatically. No manual token paste.

OAuth endpoints (all under the API host):

| Path | Purpose |
|------|---------|
| `GET /.well-known/oauth-authorization-server` | Authorization server metadata (RFC 8414) |
| `GET /.well-known/oauth-protected-resource` | Protected resource metadata (MCP spec) |
| `POST /api/oauth/mcp/register` | Dynamic client registration (RFC 7591) |
| `GET /api/oauth/mcp/authorize` | Start the flow; redirects to the consent page |
| `POST /api/oauth/mcp/approve` | Consent approval (called by the SPA, requires login) |
| `POST /api/oauth/mcp/token` | Exchange the auth code for an access token (PKCE) |

The access token returned by OAuth is a PAT under the hood, so it authenticates
through the same middleware as a static token.

## Client config examples

### Claude Desktop / Cursor / VS Code (static token)

```json
{
  "mcpServers": {
    "beo-echo": {
      "url": "https://<your-beo-echo-host>/mcp",
      "headers": { "Authorization": "Bearer beo_pat_YOUR_TOKEN" }
    }
  }
}
```

### Claude.ai (OAuth)

Settings → Integrations → Add custom MCP server → paste `https://<host>/mcp`.

### curl smoke test

```bash
# Initialize a session (note the Mcp-Session-Id response header)
curl -X POST "https://<host>/mcp" \
  -H "Authorization: Bearer beo_pat_YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"curl","version":"1"}}}'
```

## Running

The MCP endpoint is mounted on the main server and starts with it
(`beo-echo server`) — there is nothing separate to run, and no MCP-specific
configuration.

The OAuth discovery metadata and consent redirect URLs are derived per request
from the incoming `Host` / `X-Forwarded-*` headers, so they always match the URL
clients use — no base-URL env to set. (If you run behind a reverse proxy, make
sure it forwards `X-Forwarded-Proto` and `X-Forwarded-Host`.)

Internally, MCP tools dispatch to the API in-process via the shared router, so
calls pass through all API auth/permission middleware without a network
round-trip.

## Tools

51 tools across 7 areas. All project-scoped tools take `workspace_id` +
`project_id`.

- **workspace** — `workspace_list`, `workspace_create`, `workspace_check_role`, `workspace_add_member`, `workspace_list_users`
- **project** — `project_list`, `project_get`, `project_create`, `project_update`, `project_delete`, `project_get_advance_config`, `project_update_advance_config`
- **routes** — endpoints (`route_*_endpoint`), responses (`route_*_response`, `route_duplicate_response`, `route_reorder_responses`), rules (`route_*_rule`), proxies (`route_*_proxy`)
- **logs** — `logs_list`, `logs_clear`, `logs_list_bookmarks`, `logs_add_bookmark`, `logs_delete_bookmark`
- **replay** — `replay_list`, `replay_get`, `replay_create`, `replay_update`, `replay_delete`, `replay_execute`, `replay_get_logs`
- **action** — `action_list_types`, `action_list`, `action_get`, `action_create`, `action_update`, `action_delete`, `action_toggle`, `action_set_priority`
- **config** — `config_whoami`, `config_public`, `config_list_system`, `config_get_system`, `config_update_system`, `config_get_auto_invite`, `config_update_auto_invite`

Every tool forwards the caller's bearer token to the REST API, so all permission
checks (workspace access, owner-only routes) happen in the API — the MCP layer
adds no privileges of its own.
