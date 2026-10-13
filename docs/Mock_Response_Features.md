# CORS Modes, Response Templating and Sequence Mode

Three features for mocks that are called from a browser or that must behave in a scripted way.

## 1. CORS preflight per project

A browser sends an `OPTIONS` preflight before a request with a custom header such as `X-Api-Key`.
The preflight is answered by the project's own CORS config (`advance_config.cors`), so each project
decides which headers are allowed. The admin API (`/api`, `/mcp`) keeps the server-wide CORS config.

| `cors.mode` | `Access-Control-Allow-Headers` |
|---|---|
| `static` (default when `cors` is not set) | `Origin, Content-Length, Content-Type, Authorization, X-Requested-With, Accept` |
| `reflect` | the value of `Access-Control-Request-Headers`, echoed back |
| `custom` | the list in `cors.allowHeaders` |

Optional keys: `allowOrigins` (default: the server-wide `CORS_ORIGIN`, normally `*`) and `maxAge` (seconds, 0-86400).
New projects are created with `{"cors":{"mode":"reflect"}}`; existing projects keep `static`.

```bash
curl -X PUT "$BASE/api/workspaces/$WS/projects/$PROJECT_ID/advance-config" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"cors":{"mode":"reflect"}}'
```

`PUT .../advance-config` **merges** top-level keys: sending `{"delayMs":300}` keeps `cors`, and sending
`{"cors":...}` keeps `delayMs`. Send a key as `null` to remove it. The MCP tool `project_update_advance_config`
takes `delay_ms`, `cors` and `clear_cors` and behaves the same way.

A response header set on a mock response still wins on the actual request, so a response with
`"Access-Control-Allow-Origin": ""` removes the CORS header (useful to test a "no CORS" scenario).
The header override does not apply to the preflight itself.

## 2. Response templating

Turn on **Templated body** on a response (API/MCP field `templated`, default `false`) and the body is rendered
as a [Handlebars](https://handlebarsjs.com) template for every request. Responses that are not templated are
sent untouched, even if they contain `{{`.

| Helper | Example | Result |
|---|---|---|
| `queryParam` | `{{queryParam 'room_id'}}` / `{{queryParam 'n' default='10'}}` | query value, or the default (empty if missing) |
| `header` | `{{header 'X-Api-Key'}}` / `{{header 'X-Name' default='none'}}` | request header value |
| `repeat` | `{{#repeat (queryParam 'n') comma=true}}{{@index}}{{/repeat}}` | block repeated `n` times (`@index`, `@first`, `@last`); `comma=true` joins with `,` |

```text
GET /log?room_id=123&app_id=abc
[{"label":"room_id={{queryParam 'room_id'}}","value":"{{queryParam 'app_id'}}"}]
=> [{"label":"room_id=123","value":"abc"}]

GET /many?n=500
[{{#repeat (queryParam 'n') comma=true}}{"label":"Opsi {{@index}}","value":"OPT-{{@index}}"}{{/repeat}}]
=> 500 items
```

- `repeat` is capped at 10000 iterations. A missing, negative or non-numeric count renders nothing.
- When the response `Content-Type` is JSON, `queryParam` / `header` values are JSON-escaped, so a value such as `a"b` cannot break the document.
- A syntax error or an unknown helper (for example a typo like `queryParms`) returns `500` with a message that starts with `template error`/`unknown template helper`, instead of an empty value.
- Rendering happens before `Content-Encoding` compression (gzip/br).

## 3. Sequence mode and counter reset

Endpoint `response_mode` can now be `sequence`: responses are served in priority order (highest first).
What happens after the last one is set in the endpoint `advance_config`:

```json
{"sequence": {"onEnd": "stay_last"}}
```

- `stay_last` (default): keep serving the last response.
- `loop`: start over (this is what `round_robin` does).

"Fail N times, then succeed" = N error responses followed by one success response in `sequence` mode.
For a fail-once scenario: a `500` response with the higher priority and a `200` response below it.

The counter of `round_robin` and `sequence` endpoints is kept in memory. Reset it with:

- UI: the undo button next to the response mode buttons
- API: `POST /api/workspaces/{ws}/projects/{projectId}/endpoints/{endpointId}/reset-counter`
- MCP: `route_reset_counter`

The counter also resets by itself when a response of the endpoint is created, updated, duplicated, deleted or
reordered, or when the endpoint's `response_mode` / `advance_config` changes. It is lost on restart and is not
shared between instances.
