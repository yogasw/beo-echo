package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerWorkspaceTools wires the workspace area:
// list workspaces, create a workspace, check role, add a member.
func (s *Server) registerWorkspaceTools() {
	type empty struct{}

	addTool(s, "workspace_list",
		"List all workspaces the authenticated user belongs to, with their role in each.",
		func(ctx context.Context, req *mcp.CallToolRequest, _ empty) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, "/api/workspaces", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type createWorkspaceIn struct {
		Name string `json:"name" jsonschema:"name of the new workspace"`
	}
	addTool(s, "workspace_create",
		"Create a new workspace owned by the authenticated user.",
		func(ctx context.Context, req *mcp.CallToolRequest, in createWorkspaceIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Post(ctx, token, "/api/workspaces", map[string]any{"name": in.Name}, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type roleIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id to check"`
	}
	addTool(s, "workspace_check_role",
		"Get the authenticated user's role (admin/member/readonly) in a workspace.",
		func(ctx context.Context, req *mcp.CallToolRequest, in roleIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, "/api/workspaces/"+in.WorkspaceID+"/role", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type addMemberIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"target workspace id"`
		Email       string `json:"email" jsonschema:"email of the user to add"`
		Role        string `json:"role" jsonschema:"role to grant: admin, member, or readonly"`
	}
	addTool(s, "workspace_add_member",
		"Add a user to a workspace by email with a role. Requires workspace admin or owner.",
		func(ctx context.Context, req *mcp.CallToolRequest, in addMemberIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"email": in.Email, "role": in.Role}
			var out raw
			if err := s.client.Post(ctx, token, "/api/workspaces/"+in.WorkspaceID+"/members", body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "workspace_list_users",
		"List all users in a workspace with their roles. Requires workspace admin or owner.",
		func(ctx context.Context, req *mcp.CallToolRequest, in roleIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, "/api/workspaces/"+in.WorkspaceID+"/users", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})
}
