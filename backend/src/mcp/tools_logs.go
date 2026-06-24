package mcp

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerLogTools wires request-log access: list (paginated), clear, and
// bookmark management. The streaming endpoint is intentionally omitted — SSE
// doesn't map cleanly onto a single tool call.
func (s *Server) registerLogTools() {
	logsBase := func(ws, proj string) string { return projectPath(ws, proj) + "/logs" }

	type listLogsIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		Page        int    `json:"page,omitempty" jsonschema:"page number (default 1)"`
		PageSize    int    `json:"page_size,omitempty" jsonschema:"logs per page (default 100)"`
	}
	addTool(s, "logs_list",
		"List request logs for a project, paginated. Each log includes method, path, status, latency, and bodies.",
		func(ctx context.Context, req *mcp.CallToolRequest, in listLogsIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			q := url.Values{}
			if in.Page > 0 {
				q.Set("page", strconv.Itoa(in.Page))
			}
			if in.PageSize > 0 {
				q.Set("pageSize", strconv.Itoa(in.PageSize))
			}
			var out raw
			if err := s.client.Get(ctx, token, logsBase(in.WorkspaceID, in.ProjectID), q, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type projIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
	}
	addTool(s, "logs_clear",
		"Delete all request logs for a project.",
		func(ctx context.Context, req *mcp.CallToolRequest, in projIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Delete(ctx, token, logsBase(in.WorkspaceID, in.ProjectID)+"/clear", &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	addTool(s, "logs_list_bookmarks",
		"List bookmarked request logs for a project.",
		func(ctx context.Context, req *mcp.CallToolRequest, in projIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Get(ctx, token, logsBase(in.WorkspaceID, in.ProjectID)+"/bookmark", nil, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type addBookmarkIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		Log         string `json:"log" jsonschema:"the full log object as a JSON string (including its logs_hash field), as returned by logs_list"`
	}
	addTool(s, "logs_add_bookmark",
		"Bookmark a request log. Pass the log object (JSON string) from logs_list, including its logs_hash.",
		func(ctx context.Context, req *mcp.CallToolRequest, in addBookmarkIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			body := map[string]any{"logs": in.Log}
			var out raw
			if err := s.client.Post(ctx, token, logsBase(in.WorkspaceID, in.ProjectID)+"/bookmark", body, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	type delBookmarkIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		BookmarkID  string `json:"bookmark_id" jsonschema:"the bookmarked log id to remove"`
	}
	addTool(s, "logs_delete_bookmark",
		"Remove a bookmark from a request log.",
		func(ctx context.Context, req *mcp.CallToolRequest, in delBookmarkIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)
			var out raw
			if err := s.client.Delete(ctx, token, logsBase(in.WorkspaceID, in.ProjectID)+"/bookmark/"+in.BookmarkID, &out); err != nil {
				r, _, e, _ := handleErr(err)
				return r, nil, e
			}
			return jsonResult(out)
		})

	s.registerLogWaitTool(logsBase)
}

// logsPage mirrors the GET /logs response shape so we can read log ids.
type logsPage struct {
	Logs []struct {
		ID string `json:"id"`
	} `json:"logs"`
}

// registerLogWaitTool adds logs_wait: a long-poll that blocks until a request
// log newer than a given baseline appears, then returns the new logs. Because
// MCP tool calls are request/response (no streaming), this approximates the
// SSE log stream by polling in-process until something arrives or it times out.
func (s *Server) registerLogWaitTool(logsBase func(ws, proj string) string) {
	type waitIn struct {
		WorkspaceID string `json:"workspace_id" jsonschema:"the workspace id"`
		ProjectID   string `json:"project_id" jsonschema:"the project id"`
		AfterLogID  string `json:"after_log_id,omitempty" jsonschema:"return logs newer than this log id (from a prior logs_list/logs_wait). If empty, waits for the next log after the current newest."`
		TimeoutSec  int    `json:"timeout_sec,omitempty" jsonschema:"how long to wait before returning empty, in seconds (default 25, max 55)"`
	}

	// fetchNewest returns the ids on page 1 (newest first).
	fetchNewest := func(ctx context.Context, token, ws, proj string) ([]string, error) {
		var page logsPage
		if err := s.client.Get(ctx, token, logsBase(ws, proj), nil, &page); err != nil {
			return nil, err
		}
		ids := make([]string, 0, len(page.Logs))
		for _, l := range page.Logs {
			ids = append(ids, l.ID)
		}
		return ids, nil
	}

	// newerThan returns the ids that appear before baseline in the DESC list
	// (i.e. logs created after it). If baseline is "" or not found, returns nil.
	newerThan := func(ids []string, baseline string) []string {
		if baseline == "" {
			return nil
		}
		var newer []string
		for _, id := range ids {
			if id == baseline {
				return newer
			}
			newer = append(newer, id)
		}
		// Baseline not on this page: everything is newer (logs rotated/cleared).
		return ids
	}

	addTool(s, "logs_wait",
		"Wait (long-poll) for new request logs in a project, then return them. Pass after_log_id from a prior logs_list to only get newer logs. Returns an empty list if nothing arrives before the timeout. MCP has no streaming, so this polls in-process.",
		func(ctx context.Context, req *mcp.CallToolRequest, in waitIn) (*mcp.CallToolResult, any, error) {
			token := tokenFromRequest(req)

			timeout := in.TimeoutSec
			if timeout <= 0 {
				timeout = 25
			}
			if timeout > 55 {
				timeout = 55
			}
			deadline := time.Now().Add(time.Duration(timeout) * time.Second)

			// Establish a baseline: caller's after_log_id, or the current newest.
			baseline := in.AfterLogID
			if baseline == "" {
				ids, err := fetchNewest(ctx, token, in.WorkspaceID, in.ProjectID)
				if err != nil {
					r, _, e, _ := handleErr(err)
					return r, nil, e
				}
				if len(ids) > 0 {
					baseline = ids[0]
				}
			}

			const pollInterval = 1500 * time.Millisecond
			for {
				ids, err := fetchNewest(ctx, token, in.WorkspaceID, in.ProjectID)
				if err != nil {
					r, _, e, _ := handleErr(err)
					return r, nil, e
				}
				if newer := newerThan(ids, baseline); len(newer) > 0 {
					return jsonResult(map[string]any{
						"new_log_ids": newer,
						"newest_log_id": ids[0],
						"count":       len(newer),
					})
				}
				// Special case: started with no logs at all, then some appeared.
				if baseline == "" && len(ids) > 0 {
					return jsonResult(map[string]any{
						"new_log_ids":   ids,
						"newest_log_id": ids[0],
						"count":         len(ids),
					})
				}

				if time.Now().After(deadline) {
					return jsonResult(map[string]any{
						"new_log_ids":   []string{},
						"newest_log_id": baseline,
						"count":         0,
						"timed_out":     true,
					})
				}

				select {
				case <-ctx.Done():
					return textResult("logs_wait cancelled")
				case <-time.After(pollInterval):
				}
			}
		})
}
