package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// projectBase builds the workspace-scoped project collection path.
func projectBase(workspaceID string) string {
	return "/api/workspaces/" + workspaceID + "/projects"
}

// projectPath builds the path to a single project.
func projectPath(workspaceID, projectID string) string {
	return projectBase(workspaceID) + "/" + projectID
}

// registerProjectTools wires the project area: list, get, create, update,
// delete, and advance-config read/write.
func (s *Server) registerProjectTools() {
	type wsOnly struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
	}
	addTool(s, "project_list",
		"List all projects in a workspace.",
		func(ctx context.Context, req *mcp.CallToolRequest, in wsOnly) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, projectBase(in.WorkspaceID), nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type projIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
	}
	addTool(s, "project_get",
		"Get a single project's full details. Returns its general info — name, alias, "+
			"base URL (the 'url' field), mode (mock/proxy/forwarder/disabled), status, and "+
			"created_at — plus its endpoints (with responses), its proxy_targets, and which "+
			"proxy is active (active_proxy_id / active_proxy). When reporting to the user, "+
			"summarize these key fields: mode, base URL, status, when it was created, and the "+
			"proxy targets with the active one highlighted.",
		func(ctx context.Context, req *mcp.CallToolRequest, in projIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, projectPath(in.WorkspaceID, in.ProjectID), nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type createProjectIn struct {
		WorkspaceID   string `json:"workspace_id" jsonschema:"the workspace id"`
		Name          string `json:"name" jsonschema:"project name"`
		Alias         string `json:"alias,omitempty" jsonschema:"unique alias used in the mock URL (subdomain or path segment)"`
		Mode          string `json:"mode,omitempty" jsonschema:"operation mode: mock, proxy, forwarder, or disabled"`
		Documentation string `json:"documentation,omitempty" jsonschema:"optional project documentation"`
	}
	addTool(s, "project_create",
		"Create a new project in a workspace.",
		func(ctx context.Context, req *mcp.CallToolRequest, in createProjectIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"name": in.Name}
			if in.Alias != "" {
				body["alias"] = in.Alias
			}
			if in.Mode != "" {
				body["mode"] = in.Mode
			}
			if in.Documentation != "" {
				body["documentation"] = in.Documentation
			}
			var out raw
			if err := s.client.Post(ctx, token, projectBase(in.WorkspaceID), body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type updateProjectIn struct {
		WorkspaceID   string  `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID     string  `json:"project_id" jsonschema:"the project id"`
		Name          *string `json:"name,omitempty" jsonschema:"new project name"`
		Alias         *string `json:"alias,omitempty" jsonschema:"new alias"`
		Mode          *string `json:"mode,omitempty" jsonschema:"new mode: mock, proxy, forwarder, or disabled"`
		Status        *string `json:"status,omitempty" jsonschema:"new status: running, stopped, or error"`
		ActiveProxyID *string `json:"active_proxy_id,omitempty" jsonschema:"id of the proxy target to make active"`
	}
	addTool(s, "project_update",
		"Update a project. Only provided fields are changed.",
		func(ctx context.Context, req *mcp.CallToolRequest, in updateProjectIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{}
			if in.Name != nil {
				body["name"] = *in.Name
			}
			if in.Alias != nil {
				body["alias"] = *in.Alias
			}
			if in.Mode != nil {
				body["mode"] = *in.Mode
			}
			if in.Status != nil {
				body["status"] = *in.Status
			}
			if in.ActiveProxyID != nil {
				body["active_proxy_id"] = *in.ActiveProxyID
			}
			var out raw
			if err := s.client.Put(ctx, token, projectPath(in.WorkspaceID, in.ProjectID), body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "project_delete",
		"Delete a project and all its endpoints, responses, and rules.",
		func(ctx context.Context, req *mcp.CallToolRequest, in projIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Delete(ctx, token, projectPath(in.WorkspaceID, in.ProjectID), &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "project_get_advance_config",
		"Get a project's advanced config (e.g. global response delay).",
		func(ctx context.Context, req *mcp.CallToolRequest, in projIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, projectPath(in.WorkspaceID, in.ProjectID)+"/advance-config", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type advConfigIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		DelayMs     int    `json:"delay_ms" jsonschema:"global response delay in milliseconds (0-120000)"`
	}
	addTool(s, "project_update_advance_config",
		"Update a project's advanced config (sets the global response delay in ms).",
		func(ctx context.Context, req *mcp.CallToolRequest, in advConfigIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"delayMs": in.DelayMs}
			var out raw
			if err := s.client.Put(ctx, token, projectPath(in.WorkspaceID, in.ProjectID)+"/advance-config", body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})
}
