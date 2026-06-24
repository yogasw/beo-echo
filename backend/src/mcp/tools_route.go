package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// "Routes" in Beo Echo terms covers endpoints (method + path), their responses,
// the matching rules on each response, and the proxy targets a project forwards
// to. This file exposes full CRUD over all four.

func endpointsBase(ws, proj string) string {
	return projectPath(ws, proj) + "/endpoints"
}
func endpointPath(ws, proj, id string) string {
	return endpointsBase(ws, proj) + "/" + id
}
func responsesBase(ws, proj, endpointID string) string {
	return endpointPath(ws, proj, endpointID) + "/responses"
}
func responsePath(ws, proj, endpointID, respID string) string {
	return responsesBase(ws, proj, endpointID) + "/" + respID
}

func (s *Server) registerRouteTools() {
	s.registerEndpointTools()
	s.registerResponseTools()
	s.registerRuleTools()
	s.registerProxyTools()
}

// ---- Endpoints ----

func (s *Server) registerEndpointTools() {
	type projIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
	}
	addTool(s, "route_list_endpoints",
		"List all endpoints (routes) in a project, including their responses.",
		func(ctx context.Context, req *mcp.CallToolRequest, in projIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, endpointsBase(in.WorkspaceID, in.ProjectID), nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type epIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		EndpointID  string `json:"endpoint_id" jsonschema:"the endpoint id"`
	}
	addTool(s, "route_get_endpoint",
		"Get a single endpoint with its responses and proxy target.",
		func(ctx context.Context, req *mcp.CallToolRequest, in epIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, endpointPath(in.WorkspaceID, in.ProjectID, in.EndpointID), nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type createEpIn struct {
		WorkspaceID   string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID     string `json:"project_id" jsonschema:"the project id"`
		Method        string `json:"method" jsonschema:"HTTP method: GET, POST, PUT, DELETE, PATCH, etc."`
		Path          string `json:"path" jsonschema:"endpoint path, e.g. /users/:id"`
		Enabled       *bool  `json:"enabled,omitempty" jsonschema:"whether the endpoint is enabled (default true)"`
		ResponseMode  string `json:"response_mode,omitempty" jsonschema:"how responses are picked: static, random, or round_robin"`
		Documentation string `json:"documentation,omitempty" jsonschema:"optional documentation"`
	}
	addTool(s, "route_create_endpoint",
		"Create a new endpoint (route) in a project.",
		func(ctx context.Context, req *mcp.CallToolRequest, in createEpIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"method": in.Method, "path": in.Path}
			if in.Enabled != nil {
				body["enabled"] = *in.Enabled
			}
			if in.ResponseMode != "" {
				body["response_mode"] = in.ResponseMode
			}
			if in.Documentation != "" {
				body["documentation"] = in.Documentation
			}
			var out raw
			if err := s.client.Post(ctx, token, endpointsBase(in.WorkspaceID, in.ProjectID), body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type updateEpIn struct {
		WorkspaceID   string  `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID     string  `json:"project_id" jsonschema:"the project id"`
		EndpointID    string  `json:"endpoint_id" jsonschema:"the endpoint id"`
		Method        *string `json:"method,omitempty" jsonschema:"new HTTP method"`
		Path          *string `json:"path,omitempty" jsonschema:"new path"`
		Enabled       *bool   `json:"enabled,omitempty" jsonschema:"enable/disable the endpoint"`
		ResponseMode  *string `json:"response_mode,omitempty" jsonschema:"static, random, or round_robin"`
		Documentation *string `json:"documentation,omitempty" jsonschema:"new documentation for the endpoint"`
		AdvanceConfig *string `json:"advance_config,omitempty" jsonschema:"endpoint advanced config as a JSON string"`
		UseProxy      *bool   `json:"use_proxy,omitempty" jsonschema:"forward this endpoint to a proxy target"`
		ProxyTargetID *string `json:"proxy_target_id,omitempty" jsonschema:"proxy target id when use_proxy is true"`
	}
	addTool(s, "route_update_endpoint",
		"Update an endpoint. Only provided fields are changed.",
		func(ctx context.Context, req *mcp.CallToolRequest, in updateEpIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{}
			if in.Method != nil {
				body["method"] = *in.Method
			}
			if in.Path != nil {
				body["path"] = *in.Path
			}
			if in.Enabled != nil {
				body["enabled"] = *in.Enabled
			}
			if in.ResponseMode != nil {
				body["response_mode"] = *in.ResponseMode
			}
			if in.Documentation != nil {
				body["documentation"] = *in.Documentation
			}
			if in.AdvanceConfig != nil {
				body["advance_config"] = *in.AdvanceConfig
			}
			if in.UseProxy != nil {
				body["use_proxy"] = *in.UseProxy
			}
			if in.ProxyTargetID != nil {
				body["proxy_target_id"] = *in.ProxyTargetID
			}
			var out raw
			if err := s.client.Put(ctx, token, endpointPath(in.WorkspaceID, in.ProjectID, in.EndpointID), body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "route_delete_endpoint",
		"Delete an endpoint and its responses and rules.",
		func(ctx context.Context, req *mcp.CallToolRequest, in epIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Delete(ctx, token, endpointPath(in.WorkspaceID, in.ProjectID, in.EndpointID), &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})
}

// ---- Responses ----

func (s *Server) registerResponseTools() {
	type epIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		EndpointID  string `json:"endpoint_id" jsonschema:"the endpoint id"`
	}
	addTool(s, "route_list_responses",
		"List all responses configured for an endpoint, with their rules.",
		func(ctx context.Context, req *mcp.CallToolRequest, in epIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, responsesBase(in.WorkspaceID, in.ProjectID, in.EndpointID), nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type createRespIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		EndpointID  string `json:"endpoint_id" jsonschema:"the endpoint id"`
		StatusCode  int    `json:"status_code" jsonschema:"HTTP status code to return, e.g. 200"`
		Body        string `json:"body,omitempty" jsonschema:"response body (often JSON as a string)"`
		Headers     string `json:"headers,omitempty" jsonschema:"response headers as a JSON string"`
		Priority    int    `json:"priority,omitempty" jsonschema:"match priority (higher wins)"`
		DelayMs     int    `json:"delay_ms,omitempty" jsonschema:"per-response delay in milliseconds"`
		Note        string `json:"note,omitempty" jsonschema:"human note describing this response"`
		Enabled     *bool  `json:"enabled,omitempty" jsonschema:"whether this response is active"`
		IsFallback  *bool  `json:"is_fallback,omitempty" jsonschema:"use this when no other response matches"`
		RulesLogic  string `json:"rules_logic,omitempty" jsonschema:"how rules combine: and / or"`
		Stream      *bool  `json:"stream,omitempty" jsonschema:"stream the response body instead of sending it at once"`
	}
	addTool(s, "route_create_response",
		"Create a new response for an endpoint.",
		func(ctx context.Context, req *mcp.CallToolRequest, in createRespIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"status_code": in.StatusCode}
			if in.Body != "" {
				body["body"] = in.Body
			}
			if in.Headers != "" {
				body["headers"] = in.Headers
			}
			if in.Priority != 0 {
				body["priority"] = in.Priority
			}
			if in.DelayMs != 0 {
				body["delay_ms"] = in.DelayMs
			}
			if in.Note != "" {
				body["note"] = in.Note
			}
			if in.Enabled != nil {
				body["enabled"] = *in.Enabled
			}
			if in.IsFallback != nil {
				body["is_fallback"] = *in.IsFallback
			}
			if in.RulesLogic != "" {
				body["rules_logic"] = in.RulesLogic
			}
			if in.Stream != nil {
				body["stream"] = *in.Stream
			}
			var out raw
			if err := s.client.Post(ctx, token, responsesBase(in.WorkspaceID, in.ProjectID, in.EndpointID), body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type respIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		EndpointID  string `json:"endpoint_id" jsonschema:"the endpoint id"`
		ResponseID  string `json:"response_id" jsonschema:"the response id"`
	}
	addTool(s, "route_get_response",
		"Get a single response with its rules.",
		func(ctx context.Context, req *mcp.CallToolRequest, in respIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, responsePath(in.WorkspaceID, in.ProjectID, in.EndpointID, in.ResponseID), nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type updateRespIn struct {
		WorkspaceID string  `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string  `json:"project_id" jsonschema:"the project id"`
		EndpointID  string  `json:"endpoint_id" jsonschema:"the endpoint id"`
		ResponseID  string  `json:"response_id" jsonschema:"the response id"`
		StatusCode  *int    `json:"status_code,omitempty" jsonschema:"new status code"`
		Body        *string `json:"body,omitempty" jsonschema:"new body"`
		Headers     *string `json:"headers,omitempty" jsonschema:"new headers (JSON string)"`
		Priority    *int    `json:"priority,omitempty" jsonschema:"new match priority"`
		DelayMs     *int    `json:"delay_ms,omitempty" jsonschema:"new per-response delay (ms)"`
		Note        *string `json:"note,omitempty" jsonschema:"new note"`
		Enabled     *bool   `json:"enabled,omitempty" jsonschema:"enable/disable"`
		IsFallback  *bool   `json:"is_fallback,omitempty" jsonschema:"mark as fallback"`
		RulesLogic  *string `json:"rules_logic,omitempty" jsonschema:"and / or"`
		Stream      *bool   `json:"stream,omitempty" jsonschema:"stream the response body instead of sending it at once"`
	}
	addTool(s, "route_update_response",
		"Update a response. Only provided fields are changed.",
		func(ctx context.Context, req *mcp.CallToolRequest, in updateRespIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{}
			if in.StatusCode != nil {
				body["status_code"] = *in.StatusCode
			}
			if in.Body != nil {
				body["body"] = *in.Body
			}
			if in.Headers != nil {
				body["headers"] = *in.Headers
			}
			if in.Priority != nil {
				body["priority"] = *in.Priority
			}
			if in.DelayMs != nil {
				body["delay_ms"] = *in.DelayMs
			}
			if in.Note != nil {
				body["note"] = *in.Note
			}
			if in.Enabled != nil {
				body["enabled"] = *in.Enabled
			}
			if in.IsFallback != nil {
				body["is_fallback"] = *in.IsFallback
			}
			if in.RulesLogic != nil {
				body["rules_logic"] = *in.RulesLogic
			}
			if in.Stream != nil {
				body["stream"] = *in.Stream
			}
			var out raw
			if err := s.client.Put(ctx, token, responsePath(in.WorkspaceID, in.ProjectID, in.EndpointID, in.ResponseID), body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "route_delete_response",
		"Delete a response and its rules.",
		func(ctx context.Context, req *mcp.CallToolRequest, in respIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Delete(ctx, token, responsePath(in.WorkspaceID, in.ProjectID, in.EndpointID, in.ResponseID), &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "route_duplicate_response",
		"Duplicate a response, copying its rules.",
		func(ctx context.Context, req *mcp.CallToolRequest, in respIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Post(ctx, token, responsePath(in.WorkspaceID, in.ProjectID, in.EndpointID, in.ResponseID)+"/duplicate", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type reorderIn struct {
		WorkspaceID string   `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string   `json:"project_id" jsonschema:"the project id"`
		EndpointID  string   `json:"endpoint_id" jsonschema:"the endpoint id"`
		Order       []string `json:"order" jsonschema:"response ids in the desired priority order"`
	}
	addTool(s, "route_reorder_responses",
		"Reorder an endpoint's responses by priority.",
		func(ctx context.Context, req *mcp.CallToolRequest, in reorderIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"order": in.Order}
			var out raw
			if err := s.client.Put(ctx, token, responsesBase(in.WorkspaceID, in.ProjectID, in.EndpointID)+"/reorder", body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})
}

// ---- Rules ----

func (s *Server) registerRuleTools() {
	rulesBase := func(ws, proj, ep, resp string) string {
		return responsePath(ws, proj, ep, resp) + "/rules"
	}
	type respIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		EndpointID  string `json:"endpoint_id" jsonschema:"the endpoint id"`
		ResponseID  string `json:"response_id" jsonschema:"the response id"`
	}
	addTool(s, "route_list_rules",
		"List the matching rules on a response.",
		func(ctx context.Context, req *mcp.CallToolRequest, in respIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, rulesBase(in.WorkspaceID, in.ProjectID, in.EndpointID, in.ResponseID), nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type createRuleIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		EndpointID  string `json:"endpoint_id" jsonschema:"the endpoint id"`
		ResponseID  string `json:"response_id" jsonschema:"the response id"`
		Type        string `json:"type" jsonschema:"what to match: header, body, query, or path"`
		Key         string `json:"key,omitempty" jsonschema:"the key to match (e.g. header name or query param)"`
		Operator    string `json:"operator" jsonschema:"comparison: equals, contains, or regex"`
		Value       string `json:"value" jsonschema:"value to compare against"`
	}
	addTool(s, "route_create_rule",
		"Add a matching rule to a response (used to pick the response based on request content).",
		func(ctx context.Context, req *mcp.CallToolRequest, in createRuleIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"type": in.Type, "operator": in.Operator, "value": in.Value}
			if in.Key != "" {
				body["key"] = in.Key
			}
			var out raw
			if err := s.client.Post(ctx, token, rulesBase(in.WorkspaceID, in.ProjectID, in.EndpointID, in.ResponseID), body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type ruleIn struct {
		WorkspaceID string  `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string  `json:"project_id" jsonschema:"the project id"`
		EndpointID  string  `json:"endpoint_id" jsonschema:"the endpoint id"`
		ResponseID  string  `json:"response_id" jsonschema:"the response id"`
		RuleID      string  `json:"rule_id" jsonschema:"the rule id"`
		Type        *string `json:"type,omitempty" jsonschema:"header, body, query, or path"`
		Key         *string `json:"key,omitempty" jsonschema:"the key to match"`
		Operator    *string `json:"operator,omitempty" jsonschema:"equals, contains, or regex"`
		Value       *string `json:"value,omitempty" jsonschema:"value to compare against"`
	}
	addTool(s, "route_update_rule",
		"Update a matching rule on a response.",
		func(ctx context.Context, req *mcp.CallToolRequest, in ruleIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{}
			if in.Type != nil {
				body["type"] = *in.Type
			}
			if in.Key != nil {
				body["key"] = *in.Key
			}
			if in.Operator != nil {
				body["operator"] = *in.Operator
			}
			if in.Value != nil {
				body["value"] = *in.Value
			}
			path := rulesBase(in.WorkspaceID, in.ProjectID, in.EndpointID, in.ResponseID) + "/" + in.RuleID
			var out raw
			if err := s.client.Put(ctx, token, path, body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "route_delete_rule",
		"Delete a single matching rule from a response.",
		func(ctx context.Context, req *mcp.CallToolRequest, in ruleIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			path := rulesBase(in.WorkspaceID, in.ProjectID, in.EndpointID, in.ResponseID) + "/" + in.RuleID
			var out raw
			if err := s.client.Delete(ctx, token, path, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})
}

// ---- Proxy targets ----

func (s *Server) registerProxyTools() {
	proxiesBase := func(ws, proj string) string { return projectPath(ws, proj) + "/proxies" }
	type projIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
	}
	addTool(s, "route_list_proxies",
		"List proxy targets configured for a project.",
		func(ctx context.Context, req *mcp.CallToolRequest, in projIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, proxiesBase(in.WorkspaceID, in.ProjectID), nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type createProxyIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		Label       string `json:"label" jsonschema:"human label for the proxy target"`
		URL         string `json:"url" jsonschema:"base URL to forward requests to"`
	}
	addTool(s, "route_create_proxy",
		"Create a proxy target (an upstream the project can forward to).",
		func(ctx context.Context, req *mcp.CallToolRequest, in createProxyIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"label": in.Label, "url": in.URL}
			var out raw
			if err := s.client.Post(ctx, token, proxiesBase(in.WorkspaceID, in.ProjectID), body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type proxyIn struct {
		WorkspaceID string  `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string  `json:"project_id" jsonschema:"the project id"`
		ProxyID     string  `json:"proxy_id" jsonschema:"the proxy target id"`
		Label       *string `json:"label,omitempty" jsonschema:"new label"`
		URL         *string `json:"url,omitempty" jsonschema:"new base URL"`
	}
	addTool(s, "route_update_proxy",
		"Update a proxy target.",
		func(ctx context.Context, req *mcp.CallToolRequest, in proxyIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{}
			if in.Label != nil {
				body["label"] = *in.Label
			}
			if in.URL != nil {
				body["url"] = *in.URL
			}
			var out raw
			if err := s.client.Put(ctx, token, proxiesBase(in.WorkspaceID, in.ProjectID)+"/"+in.ProxyID, body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "route_delete_proxy",
		"Delete a proxy target.",
		func(ctx context.Context, req *mcp.CallToolRequest, in proxyIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Delete(ctx, token, proxiesBase(in.WorkspaceID, in.ProjectID)+"/"+in.ProxyID, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})
}
