# Beo Echo Replay API Documentation

The Replay API allows you to execute HTTP requests to external endpoints and record their responses. This feature is useful for testing, debugging, and verifying your API endpoints before creating mocks.

## Table of Contents

- [API Endpoints](#api-endpoints)
- [Execute Replay](#execute-replay)
- [Request Examples](#request-examples)
  - [GET Requests](#get-requests)
  - [POST Requests](#post-requests)
  - [Authentication](#authentication)
  - [Special Cases](#special-cases)
- [Response Format](#response-format)
- [Error Handling](#error-handling)
- [Curl Examples](#curl-examples)

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/projects/{projectId}/replays/execute` | Execute a replay request |
| GET | `/api/projects/{projectId}/replays` | List all replays for a project |
| GET | `/api/projects/{projectId}/replays/{replayId}` | Get details of a specific replay |
| GET | `/api/projects/{projectId}/replays/{replayId}/logs` | Get logs for a specific replay |

## Execute Replay

### Request Format

To execute a replay request, send a POST request to `/api/projects/{projectId}/replays/execute` with the following JSON payload:

```json
{
  "protocol": "http",          // Required: http or https
  "method": "GET",             // Required: HTTP method (GET, POST, PUT, DELETE, etc.)
  "url": "https://example.com/api/resource", // Required: Target URL
  "headers": {                 // Optional: Request headers
    "Content-Type": "application/json",
    "Authorization": "Bearer token"
  },
  "query": {                   // Optional: Query parameters
    "param1": "value1",
    "param2": "value2"
  },
  "body": "request body data", // Optional: Request body
  "metadata": {                // Optional: Additional metadata
    "key1": "value1",
    "key2": "value2"
  }
}
```

### Response Format

The API responds with a JSON object containing details about the executed request:

```json
{
  "replay_id": "62e9cbf9-b970-4fff-add2-3a8517ee76d1", // Unique ID for this replay
  "status_code": 200,                                  // HTTP status code
  "response_body": "response content",                 // Response body
  "response_headers": {                                // Response headers
    "Content-Type": "application/json",
    "Server": "nginx"
  },
  "latency_ms": 120,                                   // Response time in milliseconds
  "error": "",                                         // Error message (if any)
  "log_id": "log-uuid"                                 // ID of the request log entry
}
```

## Request Examples

### GET Requests

#### Simple GET Request

```json
{
  "protocol": "http",
  "method": "GET",
  "url": "https://api.example.com/users"
}
```

#### GET Request with Query Parameters

```json
{
  "protocol": "http",
  "method": "GET",
  "url": "https://api.example.com/users",
  "query": {
    "page": "1",
    "limit": "10",
    "sort": "name:asc",
    "filter": "status:active"
  }
}
```

#### GET Request with Headers

```json
{
  "protocol": "http",
  "method": "GET",
  "url": "https://api.example.com/users",
  "headers": {
    "Accept": "application/json",
    "User-Agent": "beo-echo-client"
  }
}
```

### POST Requests

#### POST with JSON Body

```json
{
  "protocol": "http",
  "method": "POST",
  "url": "https://api.example.com/users",
  "headers": {
    "Content-Type": "application/json",
    "Accept": "application/json"
  },
  "body": "{\"name\":\"John Doe\",\"email\":\"john@example.com\",\"age\":30}"
}
```

#### POST with Form Data

```json
{
  "protocol": "http",
  "method": "POST",
  "url": "https://api.example.com/form",
  "headers": {
    "Content-Type": "application/x-www-form-urlencoded"
  },
  "body": "username=testuser&password=secret&remember=true"
}
```

#### POST with Complex JSON

```json
{
  "protocol": "http",
  "method": "POST",
  "url": "https://api.example.com/complex",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": "{\"user\":{\"name\":\"Test User\",\"email\":\"test@example.com\",\"preferences\":{\"theme\":\"dark\",\"notifications\":true},\"tags\":[\"test\",\"api\",\"json\"]},\"metadata\":{\"client\":\"beo-echo-test\",\"version\":\"1.0.0\"}}"
}
```

### Authentication

#### Bearer Token Authentication

```json
{
  "protocol": "http",
  "method": "GET",
  "url": "https://api.example.com/secure",
  "headers": {
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
  }
}
```

#### Basic Authentication

```json
{
  "protocol": "http",
  "method": "GET",
  "url": "https://api.example.com/basic-auth",
  "headers": {
    "Authorization": "Basic dXNlcjpwYXNz"
  }
}
```

#### API Key Authentication

```json
{
  "protocol": "http",
  "method": "GET",
  "url": "https://api.example.com/secure",
  "headers": {
    "X-API-Key": "your-api-key-here"
  }
}
```

#### Multiple Authentication Methods

```json
{
  "protocol": "http",
  "method": "GET",
  "url": "https://api.example.com/secure",
  "headers": {
    "Authorization": "Bearer token-here",
    "X-API-Key": "api-key-here",
    "Cookie": "session=abc123; user=testuser"
  }
}
```

### Special Cases

#### Handling Network Errors

If the request encounters a network error (e.g., non-existent domain), the response will have:
- `status_code`: 0
- `error`: Error message describing the network issue
- A valid `replay_id` and `latency_ms` will still be provided

Example response:

```json
{
  "replay_id": "be31b4ed-64fd-4c1e-8d42-dd2928922083",
  "status_code": 0,
  "response_body": "",
  "response_headers": null,
  "latency_ms": 120,
  "error": "Get \"https://non-existent-domain-12345.com\": dial tcp: lookup non-existent-domain-12345.com: no such host",
  "log_id": "log-uuid"
}
```

#### HTTP Errors

HTTP errors (e.g., 404, 500) are not considered service errors and will return:
- The actual `status_code` (e.g., 404)
- Any response body or headers returned by the server
- No `error` field will be populated

## Error Handling

The Replay API differentiates between several types of errors:

1. **Validation Errors**: Occur before the request is sent, such as:
   - Missing required fields (protocol, method, URL)
   - Invalid project ID
   - Unsupported protocol

2. **Network Errors**: Occur during request execution, such as:
   - Invalid domain
   - Connection timeout
   - DNS resolution failure

3. **HTTP Errors**: These are valid HTTP responses with error status codes (4xx, 5xx)
   - These are NOT considered service errors
   - The full response is captured and returned

## Curl Examples

### Basic GET Request

```bash
curl -X POST https://api.beo-echo.example.com/api/projects/project-uuid/replays/execute \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "protocol": "http",
    "method": "GET",
    "url": "https://api.example.com/users"
  }'
```

### POST with JSON Body

```bash
curl -X POST https://api.beo-echo.example.com/api/projects/project-uuid/replays/execute \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "protocol": "http",
    "method": "POST",
    "url": "https://api.example.com/users",
    "headers": {
      "Content-Type": "application/json",
      "Accept": "application/json"
    },
    "body": "{\"name\":\"John Doe\",\"email\":\"john@example.com\"}"
  }'
```

### GET with Query Parameters

```bash
curl -X POST https://api.beo-echo.example.com/api/projects/project-uuid/replays/execute \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "protocol": "http",
    "method": "GET",
    "url": "https://api.example.com/search",
    "query": {
      "q": "test query",
      "page": "1",
      "limit": "25",
      "sort": "relevance"
    }
  }'
```

### POST with Form Data

```bash
curl -X POST https://api.beo-echo.example.com/api/projects/project-uuid/replays/execute \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "protocol": "http",
    "method": "POST",
    "url": "https://api.example.com/login",
    "headers": {
      "Content-Type": "application/x-www-form-urlencoded"
    },
    "body": "username=testuser&password=secret&remember=true"
  }'
```

### Authentication Request

```bash
curl -X POST https://api.beo-echo.example.com/api/projects/project-uuid/replays/execute \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "protocol": "http",
    "method": "GET",
    "url": "https://api.example.com/secured-endpoint",
    "headers": {
      "Authorization": "Bearer api-auth-token",
      "X-API-Key": "api-key-value"
    }
  }'
```
