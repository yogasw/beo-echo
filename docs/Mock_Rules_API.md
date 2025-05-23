# Mock Rules API

This document provides examples for using the Mock Rules API in Beo Echo.

## Overview

Rules allow you to define conditions that determine when a mock response should be served. Each rule is associated with a specific response and can check HTTP headers, query parameters, or the request body.

## API Endpoints

All endpoints require authentication and proper permissions.

### List Rules for a Response

```
GET /mock/api/projects/:projectId/endpoints/:id/responses/:responseId/rules
```

#### Example Response

```json
{
  "success": true,
  "data": [
    {
      "id": "6c45f12b-3a0c-4f7a-8e5b-7d1d1234abcd",
      "response_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
      "type": "header",
      "key": "X-API-Key",
      "operator": "equals",
      "value": "secret-key-123"
    },
    {
      "id": "8f21a5c7-9b6d-4e3f-2a1b-3c4d5e6f7g8h",
      "response_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
      "type": "query",
      "key": "version",
      "operator": "contains",
      "value": "v2"
    }
  ]
}
```

### Get a Single Rule

```
GET /mock/api/projects/:projectId/endpoints/:id/responses/:responseId/rules/:ruleId
```

#### Example Response

```json
{
  "success": true,
  "data": {
    "id": "6c45f12b-3a0c-4f7a-8e5b-7d1d1234abcd",
    "response_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
    "type": "header",
    "key": "X-API-Key",
    "operator": "equals",
    "value": "secret-key-123"
  }
}
```

### Create a Rule

```
POST /mock/api/projects/:projectId/endpoints/:id/responses/:responseId/rules
```

#### Request Body

```json
{
  "type": "header",
  "key": "Authorization",
  "operator": "contains",
  "value": "Bearer"
}
```

#### Example Response

```json
{
  "success": true,
  "message": "Rule created successfully",
  "data": {
    "id": "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
    "response_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
    "type": "header",
    "key": "Authorization",
    "operator": "contains",
    "value": "Bearer"
  }
}
```

### Update a Rule

```
PUT /mock/api/projects/:projectId/endpoints/:id/responses/:responseId/rules/:ruleId
```

#### Request Body

```json
{
  "type": "query",
  "key": "auth",
  "operator": "equals",
  "value": "true"
}
```

#### Example Response

```json
{
  "success": true,
  "message": "Rule updated successfully",
  "data": {
    "id": "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
    "response_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
    "type": "query",
    "key": "auth",
    "operator": "equals",
    "value": "true"
  }
}
```

### Delete a Rule

```
DELETE /mock/api/projects/:projectId/endpoints/:id/responses/:responseId/rules/:ruleId
```

#### Example Response

```json
{
  "success": true,
  "message": "Rule deleted successfully"
}
```

### Delete All Rules for a Response

```
DELETE /mock/api/projects/:projectId/endpoints/:id/responses/:responseId/rules
```

#### Example Response

```json
{
  "success": true,
  "message": "All rules deleted successfully"
}
```

## Rule Types and Operators

### Rule Types

- `header`: Match against an HTTP header
- `query`: Match against a query parameter
- `body`: Match against a value in the request body (supports nested JSON paths)

### Operators

- `equals`: Exact match
- `contains`: String contains
- `regex`: Regular expression match (not fully implemented yet)

## Notes

- Rules are evaluated in the order they are defined
- A response is selected only if all rules match
- If multiple responses have matching rules, they are sorted by priority and selected according to the endpoint's response mode
