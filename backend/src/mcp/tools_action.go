package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerActionTools wires the actions area: list available action types,
// and full CRUD + toggle + priority over a project's actions. Actions
// transform requests/responses (e.g. replace text, run JavaScript).
func (s *Server) registerActionTools() {
	actionsBase := func(ws, proj string) string { return projectPath(ws, proj) + "/actions" }
	actionPath := func(ws, proj, id string) string { return actionsBase(ws, proj) + "/" + id }

	addTool(s, "action_list_types",
		"List the available action types (e.g. replace_text, run_javascript) and their config fields.",
		func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, "/api/action-types", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type projIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
	}
	addTool(s, "action_list",
		"List actions configured for a project, with their filters.",
		func(ctx context.Context, req *mcp.CallToolRequest, in projIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, actionsBase(in.WorkspaceID, in.ProjectID), nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type actionIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		ActionID    string `json:"action_id" jsonschema:"the action id"`
	}
	addTool(s, "action_get",
		"Get a single action with its filters.",
		func(ctx context.Context, req *mcp.CallToolRequest, in actionIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, actionPath(in.WorkspaceID, in.ProjectID, in.ActionID), nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type filterIn struct {
		Type     string `json:"type" jsonschema:"what to match: method, path, header, query, or status_code"`
		Key      string `json:"key,omitempty" jsonschema:"key to match (e.g. header name)"`
		Operator string `json:"operator" jsonschema:"equals, contains, regex, starts_with, or ends_with"`
		Value    string `json:"value" jsonschema:"value to compare against"`
	}
	type createActionIn struct {
		WorkspaceID    string     `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID      string     `json:"project_id" jsonschema:"the project id"`
		Name           string     `json:"name" jsonschema:"action name"`
		Type           string     `json:"type" jsonschema:"action type: replace_text or run_javascript"`
		ExecutionPoint string     `json:"execution_point,omitempty" jsonschema:"when it runs: before_request or after_request"`
		Config         string     `json:"config" jsonschema:"action config as a JSON string (shape depends on type; see action_list_types)"`
		Enabled        *bool      `json:"enabled,omitempty" jsonschema:"whether the action is enabled"`
		Priority       int        `json:"priority,omitempty" jsonschema:"execution priority"`
		Filters        []filterIn `json:"filters,omitempty" jsonschema:"filters that decide when the action applies"`
	}
	addTool(s, "action_create",
		"Create an action that transforms matching requests or responses.",
		func(ctx context.Context, req *mcp.CallToolRequest, in createActionIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"name": in.Name, "type": in.Type, "config": in.Config}
			if in.ExecutionPoint != "" {
				body["execution_point"] = in.ExecutionPoint
			}
			if in.Enabled != nil {
				body["enabled"] = *in.Enabled
			}
			if in.Priority != 0 {
				body["priority"] = in.Priority
			}
			if len(in.Filters) > 0 {
				body["filters"] = in.Filters
			}
			var out raw
			if err := s.client.Post(ctx, token, actionsBase(in.WorkspaceID, in.ProjectID), body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type updateActionIn struct {
		WorkspaceID    string     `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID      string     `json:"project_id" jsonschema:"the project id"`
		ActionID       string     `json:"action_id" jsonschema:"the action id"`
		Name           *string    `json:"name,omitempty" jsonschema:"new name"`
		ExecutionPoint *string    `json:"execution_point,omitempty" jsonschema:"before_request or after_request"`
		Config         *string    `json:"config,omitempty" jsonschema:"new config (JSON string)"`
		Enabled        *bool      `json:"enabled,omitempty" jsonschema:"enable/disable"`
		Priority       *int       `json:"priority,omitempty" jsonschema:"new priority"`
		Filters        []filterIn `json:"filters,omitempty" jsonschema:"replace the action's filters"`
	}
	addTool(s, "action_update",
		"Update an action. Only provided fields are changed.",
		func(ctx context.Context, req *mcp.CallToolRequest, in updateActionIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{}
			if in.Name != nil {
				body["name"] = *in.Name
			}
			if in.ExecutionPoint != nil {
				body["execution_point"] = *in.ExecutionPoint
			}
			if in.Config != nil {
				body["config"] = *in.Config
			}
			if in.Enabled != nil {
				body["enabled"] = *in.Enabled
			}
			if in.Priority != nil {
				body["priority"] = *in.Priority
			}
			if in.Filters != nil {
				body["filters"] = in.Filters
			}
			var out raw
			if err := s.client.Put(ctx, token, actionPath(in.WorkspaceID, in.ProjectID, in.ActionID), body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "action_delete",
		"Delete an action.",
		func(ctx context.Context, req *mcp.CallToolRequest, in actionIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Delete(ctx, token, actionPath(in.WorkspaceID, in.ProjectID, in.ActionID), &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "action_toggle",
		"Toggle an action's enabled state.",
		func(ctx context.Context, req *mcp.CallToolRequest, in actionIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Post(ctx, token, actionPath(in.WorkspaceID, in.ProjectID, in.ActionID)+"/toggle", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type priorityIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		ActionID    string `json:"action_id" jsonschema:"the action id"`
		Priority    int    `json:"priority" jsonschema:"new priority (>= 1)"`
	}
	addTool(s, "action_set_priority",
		"Set an action's execution priority.",
		func(ctx context.Context, req *mcp.CallToolRequest, in priorityIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"priority": in.Priority}
			var out raw
			if err := s.client.Patch(ctx, token, actionPath(in.WorkspaceID, in.ProjectID, in.ActionID)+"/priority", body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})
}
