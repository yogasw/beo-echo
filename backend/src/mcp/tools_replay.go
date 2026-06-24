package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerReplayTools wires replay management: list, get, create, update,
// delete, execute (fire the request live), and fetch a replay's logs.
func (s *Server) registerReplayTools() {
	replaysBase := func(ws, proj string) string { return projectPath(ws, proj) + "/replays" }
	replayPath := func(ws, proj, id string) string { return replaysBase(ws, proj) + "/" + id }

	type projIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
	}
	addTool(s, "replay_list",
		"List saved replays (preset requests) and folders for a project.",
		func(ctx context.Context, req *mcp.CallToolRequest, in projIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, replaysBase(in.WorkspaceID, in.ProjectID), nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type replayIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		ReplayID    string `json:"replay_id" jsonschema:"the replay id"`
	}
	addTool(s, "replay_get",
		"Get a single saved replay.",
		func(ctx context.Context, req *mcp.CallToolRequest, in replayIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, replayPath(in.WorkspaceID, in.ProjectID, in.ReplayID), nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type headerKV struct {
		Key   string `json:"key" jsonschema:"header name"`
		Value string `json:"value" jsonschema:"header value"`
	}
	type createReplayIn struct {
		WorkspaceID string     `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string     `json:"project_id" jsonschema:"the project id"`
		Name        string     `json:"name,omitempty" jsonschema:"name for the replay"`
		Method      string     `json:"method" jsonschema:"HTTP method"`
		URL         string     `json:"url" jsonschema:"full request URL"`
		Headers     []headerKV `json:"headers,omitempty" jsonschema:"request headers"`
		Payload     string     `json:"payload,omitempty" jsonschema:"request body"`
		FolderID    *string    `json:"folder_id,omitempty" jsonschema:"optional folder id to place the replay in"`
	}
	addTool(s, "replay_create",
		"Save a new replay (a preset HTTP request) for later execution.",
		func(ctx context.Context, req *mcp.CallToolRequest, in createReplayIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{
				"protocol": "http",
				"method":   in.Method,
				"url":      in.URL,
			}
			if in.Name != "" {
				body["name"] = in.Name
			}
			if len(in.Headers) > 0 {
				body["headers"] = in.Headers
			}
			if in.Payload != "" {
				body["payload"] = in.Payload
			}
			if in.FolderID != nil {
				body["folder_id"] = *in.FolderID
			}
			var out raw
			if err := s.client.Post(ctx, token, replaysBase(in.WorkspaceID, in.ProjectID), body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type updateReplayIn struct {
		WorkspaceID string     `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string     `json:"project_id" jsonschema:"the project id"`
		ReplayID    string     `json:"replay_id" jsonschema:"the replay id"`
		Name        *string    `json:"name,omitempty" jsonschema:"new name"`
		Method      *string    `json:"method,omitempty" jsonschema:"new HTTP method"`
		URL         *string    `json:"url,omitempty" jsonschema:"new URL"`
		Headers     []headerKV `json:"headers,omitempty" jsonschema:"replace request headers"`
		Payload     *string    `json:"payload,omitempty" jsonschema:"new request body"`
	}
	addTool(s, "replay_update",
		"Update a saved replay. Only provided fields are changed.",
		func(ctx context.Context, req *mcp.CallToolRequest, in updateReplayIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{}
			if in.Name != nil {
				body["name"] = *in.Name
			}
			if in.Method != nil {
				body["method"] = *in.Method
			}
			if in.URL != nil {
				body["url"] = *in.URL
			}
			if in.Headers != nil {
				body["headers"] = in.Headers
			}
			if in.Payload != nil {
				body["payload"] = *in.Payload
			}
			var out raw
			if err := s.client.Put(ctx, token, replayPath(in.WorkspaceID, in.ProjectID, in.ReplayID), body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "replay_delete",
		"Delete a saved replay.",
		func(ctx context.Context, req *mcp.CallToolRequest, in replayIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Delete(ctx, token, replayPath(in.WorkspaceID, in.ProjectID, in.ReplayID), &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type executeIn struct {
		WorkspaceID string            `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string            `json:"project_id" jsonschema:"the project id"`
		Method      string            `json:"method" jsonschema:"HTTP method"`
		URL         string            `json:"url" jsonschema:"full request URL to send"`
		Headers     map[string]string `json:"headers,omitempty" jsonschema:"request headers as key/value"`
		Query       map[string]string `json:"query,omitempty" jsonschema:"query parameters as key/value"`
		Payload     string            `json:"payload,omitempty" jsonschema:"request body"`
	}
	addTool(s, "replay_execute",
		"Execute an HTTP request live and return the response (status, headers, body, latency). Useful for ad-hoc testing.",
		func(ctx context.Context, req *mcp.CallToolRequest, in executeIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{
				"protocol": "http",
				"method":   in.Method,
				"url":      in.URL,
			}
			if len(in.Headers) > 0 {
				body["headers"] = in.Headers
			}
			if len(in.Query) > 0 {
				body["query"] = in.Query
			}
			if in.Payload != "" {
				body["payload"] = in.Payload
			}
			var out raw
			if err := s.client.Post(ctx, token, replaysBase(in.WorkspaceID, in.ProjectID)+"/execute", body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "replay_get_logs",
		"Get the execution logs recorded for a replay.",
		func(ctx context.Context, req *mcp.CallToolRequest, in replayIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, replayPath(in.WorkspaceID, in.ProjectID, in.ReplayID)+"/logs", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})
}
