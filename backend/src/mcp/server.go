package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Version reported to MCP clients.
const Version = "0.1.0"

// serverInstructions is returned in the initialize response. It gives the model
// the mental model it needs to use the tools well: the resource hierarchy, the
// project modes, and the usual workflow (resolve ids before acting).
const serverInstructions = `Beo Echo is an API mocking, proxy, and replay service. Use these tools to manage mocks and inspect traffic.

Resource hierarchy:
  workspace -> project -> endpoint (a route: method + path) -> response -> rule
A project also has proxy targets, request logs, replays, and actions.

Project modes (set via project_update "mode"):
  - mock: serve configured responses only
  - proxy: use a mock if one matches, otherwise forward to a proxy target
  - forwarder: always forward to a proxy target, logging traffic
  - disabled: inactive

Typical workflow:
  1. Call workspace_list to get a workspace_id.
  2. Call project_list (with that workspace_id) to get a project_id, or project_create.
  3. Most tools are project-scoped and need BOTH workspace_id and project_id.
  4. To mock an API: route_create_endpoint, then route_create_response, optionally route_create_rule for conditional matching.
  5. Inspect traffic with logs_list; to watch for new traffic, call logs_wait (long-poll). Test live requests with replay_execute.

Notes:
  - Always resolve ids with the list tools before acting; do not guess ids.
  - Response "headers" and action "config" are JSON encoded as a string.
  - System config and auto-invite tools require an instance owner.`

// Server bundles the MCP server with the REST client it drives.
type Server struct {
	mcp    *mcp.Server
	client *Client
}

// NewServer builds the MCP server, registering every tool. apiHandler is the
// Beo Echo API's http.Handler (the Gin router); tools dispatch to it in-process,
// so every call passes through the API's auth and permission middleware.
func NewServer(apiHandler http.Handler) *Server {
	impl := &mcp.Implementation{
		Name:    "beo-echo",
		Title:   "Beo Echo",
		Version: Version,
	}
	s := &Server{
		mcp:    mcp.NewServer(impl, &mcp.ServerOptions{Instructions: serverInstructions}),
		client: NewClient(apiHandler),
	}

	// Register tools grouped by area.
	s.registerWorkspaceTools()
	s.registerProjectTools()
	s.registerRouteTools()
	s.registerLogTools()
	s.registerReplayTools()
	s.registerActionTools()
	s.registerConfigTools()

	// Attach the category catalogue to tools/list responses.
	s.mcp.AddReceivingMiddleware(s.withCategoryMeta())

	return s
}

// HTTPHandler returns an http.Handler serving the streamable MCP transport.
// The same server instance is reused across sessions; per-request auth is read
// from the Authorization header inside each tool.
func (s *Server) HTTPHandler() http.Handler {
	return mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return s.mcp },
		&mcp.StreamableHTTPOptions{
			// The endpoint is always reached through a reverse proxy / tunnel, so
			// the connection arrives via localhost while carrying a public Host
			// header (e.g. "local2.yogasw.my.id"). The SDK's DNS-rebinding guard
			// rejects that with 403, so disable it — request authorization is
			// enforced by the bearer token (PAT/OAuth) inside each tool, not by
			// this Host check.
			DisableLocalhostProtection: true,
		},
	)
}

// tokenFromRequest extracts the bearer token from the MCP request's transport
// headers. Returns "" when absent.
func tokenFromRequest(req *mcp.CallToolRequest) string {
	if req == nil || req.Extra == nil || req.Extra.Header == nil {
		return ""
	}
	authz := req.Extra.Header.Get("Authorization")
	if authz == "" {
		return ""
	}
	parts := strings.SplitN(authz, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	// Allow a bare token too, for lenient clients.
	return strings.TrimSpace(authz)
}

// jsonResult renders any value as a pretty-printed JSON text result. This is the
// standard success shape for read/write tools.
func jsonResult(v any) (*mcp.CallToolResult, any, error) {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return errorResult("failed to encode result: " + err.Error()), nil, nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(buf)}},
	}, nil, nil
}

// textResult returns a plain text result.
func textResult(msg string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}, nil, nil
}

// errorResult marks the tool call as an error so the model sees it failed.
func errorResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}
}

// handleErr converts a client/API error into a tool error result. Returns
// (result, true) when err is non-nil so callers can early-return.
func handleErr(err error) (*mcp.CallToolResult, any, error, bool) {
	if err == nil {
		return nil, nil, nil, false
	}
	return errorResult(err.Error()), nil, nil, true
}

// raw is a convenience alias for decoding arbitrary JSON responses.
type raw = map[string]any

// addTool is a thin wrapper around the generic mcp.AddTool to keep registration
// sites compact. The In type's json/jsonschema struct tags drive the schema.
// Each tool is tagged with a category (derived from its name prefix) under
// _meta.category so clients can group the many tools by area.
func addTool[In any](s *Server, name, desc string, fn func(context.Context, *mcp.CallToolRequest, In) (*mcp.CallToolResult, any, error)) {
	tool := &mcp.Tool{Name: name, Description: desc}
	if cat := categoryForTool(name); cat != "" {
		tool.Meta = mcp.Meta{"category": cat}
	}
	mcp.AddTool(s.mcp, tool, fn)
}

// toolCategory describes one tool group, surfaced in the tools/list result
// under _meta.categories.
type toolCategory struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	// prefixes that map a tool name to this category.
	prefixes []string
}

// toolCategories is the ordered catalogue of tool groups. Order also defines
// display order for clients that respect it.
var toolCategories = []toolCategory{
	{Title: "Workspaces", Description: "Workspaces, members, and roles", prefixes: []string{"workspace_"}},
	{Title: "Projects", Description: "Projects and their advanced config", prefixes: []string{"project_"}},
	{Title: "Routes", Description: "Endpoints, responses, rules, and proxy targets", prefixes: []string{"route_"}},
	{Title: "Logs", Description: "Request logs and bookmarks", prefixes: []string{"logs_"}},
	{Title: "Replay", Description: "Saved replays and live request execution", prefixes: []string{"replay_"}},
	{Title: "Actions", Description: "Request/response transform actions", prefixes: []string{"action_"}},
	{Title: "Config", Description: "User, instance, and workspace configuration", prefixes: []string{"config_"}},
}

// categoryForTool returns the category title for a tool name, or "" if none match.
func categoryForTool(name string) string {
	for _, c := range toolCategories {
		for _, p := range c.prefixes {
			if strings.HasPrefix(name, p) {
				return c.Title
			}
		}
	}
	return ""
}

// categoriesMeta returns the value placed under tools/list _meta.categories.
func categoriesMeta() []map[string]any {
	out := make([]map[string]any, 0, len(toolCategories))
	for _, c := range toolCategories {
		out = append(out, map[string]any{
			"title":       c.Title,
			"description": c.Description,
		})
	}
	return out
}

// withCategoryMeta wraps the server's method handling so the tools/list result
// carries the category catalogue under _meta.categories. Per-tool categories
// are set at registration time via addTool.
func (s *Server) withCategoryMeta() mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			res, err := next(ctx, method, req)
			if err != nil || method != "tools/list" {
				return res, err
			}
			if lt, ok := res.(*mcp.ListToolsResult); ok && lt != nil {
				if lt.Meta == nil {
					lt.Meta = mcp.Meta{}
				}
				lt.Meta["categories"] = categoriesMeta()
			}
			return res, err
		}
	}
}

var _ = context.Background // keep context imported for tool signatures
