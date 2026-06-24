package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerConfigTools wires configuration access: the authenticated user's
// profile, public instance config, and (owner-only) system configuration
// read/update plus workspace auto-invite settings.
func (s *Server) registerConfigTools() {
	addTool(s, "config_whoami",
		"Get the authenticated user's profile (id, email, name, owner flag, feature flags).",
		func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, "/api/auth/me", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "config_public",
		"Get the public instance configuration (auth state, landing page, mock URL format).",
		func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, "/api/config/public", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "config_list_system",
		"List all system configuration entries. Requires instance owner.",
		func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, "/api/system-configs", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type keyIn struct {
		Key string `json:"key" jsonschema:"the system config key"`
	}
	addTool(s, "config_get_system",
		"Get a single system configuration entry by key. Requires instance owner.",
		func(ctx context.Context, req *mcp.CallToolRequest, in keyIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, "/api/system-config/"+in.Key, nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type setConfigIn struct {
		Key   string `json:"key" jsonschema:"the system config key"`
		Value string `json:"value" jsonschema:"the new value"`
	}
	addTool(s, "config_update_system",
		"Update a system configuration entry. Requires instance owner.",
		func(ctx context.Context, req *mcp.CallToolRequest, in setConfigIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"value": in.Value}
			var out raw
			if err := s.client.Put(ctx, token, "/api/system-config/"+in.Key, body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type wsIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
	}
	addTool(s, "config_get_auto_invite",
		"Get a workspace's auto-invite configuration. Requires instance owner.",
		func(ctx context.Context, req *mcp.CallToolRequest, in wsIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, "/api/workspaces/"+in.WorkspaceID+"/auto-invite", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type setAutoInviteIn struct {
		WorkspaceID string   `json:"workspace_id" jsonschema:"the workspace id"`
		Enabled     bool     `json:"enabled" jsonschema:"whether auto-invite is enabled"`
		Domains     []string `json:"domains,omitempty" jsonschema:"email domains to auto-invite"`
		Role        string   `json:"role,omitempty" jsonschema:"role for auto-invited users: admin or member"`
	}
	addTool(s, "config_update_auto_invite",
		"Update a workspace's auto-invite configuration. Requires instance owner.",
		func(ctx context.Context, req *mcp.CallToolRequest, in setAutoInviteIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"enabled": in.Enabled}
			if in.Domains != nil {
				body["domains"] = in.Domains
			}
			if in.Role != "" {
				body["role"] = in.Role
			}
			var out raw
			if err := s.client.Put(ctx, token, "/api/workspaces/"+in.WorkspaceID+"/auto-invite", body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})
}
