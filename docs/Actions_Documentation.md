# Actions - User Guide (v2.6.0)

## What are Actions?

Actions allow you to automatically modify requests and responses. You can transform data, add headers, or run custom logic on every request coming to your project.

---

## Architecture

```
┌─────────────────────────────────────────────────┐
│              Incoming Request                    │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│   BEFORE_REQUEST Actions (by priority)          │
│   - Filter evaluation (OR logic)                │
│   - Execute enabled actions                     │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│         Process Request (Mock/Proxy)            │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│   AFTER_REQUEST Actions (by priority)           │
│   - Filter evaluation (OR logic)                │
│   - Execute enabled actions                     │
└─────────────────┬───────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────┐
│              Return Response                     │
└─────────────────────────────────────────────────┘
```

---

## When are Actions Executed?

Actions can run at 2 different times:

### 1. **Before Request**
Executed **before** the request is processed by mock or proxy.

**Use for:**
- Modify incoming request body
- Add/modify request headers
- Validate input data
- Log requests

### 2. **After Request** - *Default*
Executed **after** the response is received from mock or proxy.

**Use for:**
- Modify response body
- Add/modify response headers
- Add debug information
- Hide sensitive data

---

## Action Types

### 1. Replace Text

Find and replace text in requests or responses using plain string or regex patterns.

#### Available Targets:
- **request_body** - Request body content
- **response_body** - Response body content
- **request_header** - Specific request header
- **response_header** - Specific response header

#### Usage Examples:

**Simple text replacement:**
```json
{
  "target": "response_body",
  "pattern": "Hello",
  "replacement": "Hi",
  "use_regex": false
}
```
Input: `{"message": "Hello World"}`
Output: `{"message": "Hi World"}`

**Using Regex:**
```json
{
  "target": "response_body",
  "pattern": "user_\\d+",
  "replacement": "customer_$0",
  "use_regex": true
}
```
Input: `{"id": "user_123"}`
Output: `{"id": "customer_user_123"}`

**Modify Headers:**
```json
{
  "target": "response_header",
  "pattern": "application/json",
  "replacement": "application/vnd.api+json",
  "use_regex": false,
  "header_key": "Content-Type"
}
```

---

### 2. Run JavaScript

Execute custom JavaScript code to modify requests or responses.

#### Available Objects:

```javascript
// Request (available in both before_request and after_request)
request.method       // "GET", "POST", etc
request.path         // "/api/users"
request.query        // {"page": ["1"]}
request.headers      // {"Content-Type": "application/json"}
request.body         // String body

// Response (only available in after_request)
response.status_code // 200, 404, etc
response.headers     // {"Content-Type": "application/json"}
response.body        // String body

// Console for debugging
console.log("Debug message");
```

#### Script Examples:

**1. Add timestamp to response:**
```javascript
var data = JSON.parse(response.body);
data.timestamp = new Date().toISOString();
response.body = JSON.stringify(data);
```

**2. Add custom headers:**
```javascript
response.headers["X-API-Version"] = "2.6.0";
response.headers["X-Server"] = "Beo-Echo";
```

**3. Transform request format:**
```javascript
var data = JSON.parse(request.body);

// Convert snake_case to camelCase
if (data.user_name) {
  data.userName = data.user_name;
  delete data.user_name;
}

request.body = JSON.stringify(data);
```

**4. Hide sensitive data:**
```javascript
var data = JSON.parse(response.body);

// Hide email
if (data.email) {
  data.email = data.email.replace(/(.{2})(.*)(@.*)/, "$1***$3");
}

// Hide credit card number
if (data.card_number) {
  data.card_number = "****-****-****-" + data.card_number.slice(-4);
}

response.body = JSON.stringify(data);
```

**5. Add debug info (development only):**
```javascript
var data = JSON.parse(response.body);

// Add debug info if X-Dev-Mode header exists
if (request.headers["X-Dev-Mode"] === "true") {
  data._debug = {
    timestamp: new Date().toISOString(),
    path: request.path,
    method: request.method
  };
  console.log("Debug mode enabled");
}

response.body = JSON.stringify(data);
```

#### JavaScript Limitations:
- **Timeout:** 5 seconds maximum
- **No async/await** - Synchronous code only
- **No external libraries** - Standard JavaScript only (ES5)
- **Sandboxed** - No file system or network access

---

## Filters (When Actions Execute)

Filters control when an action should execute. Without filters, the action runs for **all requests**.

### Filter Logic: **OR** (Any Match)
If **any** filter matches, the action will execute.

### Filter Types:

#### 1. Method
Match by HTTP method.
```json
{
  "type": "method",
  "operator": "equals",
  "value": "POST"
}
```

#### 2. Path
Match by URL path.
```json
{
  "type": "path",
  "operator": "starts_with",
  "value": "/api/users"
}
```

#### 3. Header
Match by header value.
```json
{
  "type": "header",
  "key": "Content-Type",
  "operator": "contains",
  "value": "json"
}
```

#### 4. Status Code
Match by response status code (only for after_request).
```json
{
  "type": "status_code",
  "operator": "equals",
  "value": "200"
}
```

### Available Operators:

| Operator | Description | Example |
|----------|-------------|---------|
| `equals` | Exact match | `"POST"` = `"POST"` |
| `contains` | Contains substring | `"application/json"` contains `"json"` |
| `starts_with` | Starts with prefix | `"/api/users"` starts with `"/api"` |
| `ends_with` | Ends with suffix | `"image.png"` ends with `".png"` |
| `regex` | Regular expression | Match pattern `"^user_\\d+$"` |

---

## Priority (Execution Order)

Actions execute based on **priority** (order). Priority starts from **1** (highest).

### Priority Rules:
- **Lower number = Higher priority** (Priority 1 executes first)
- Separate priorities for `before_request` and `after_request`
- New actions automatically get the last priority
- You can change priority to reorder execution

### Example:
```
Priority 1: Validate Input      → Executed first
Priority 2: Transform Data       → Executed second
Priority 3: Add Debug Headers    → Executed last
```

### Ordering Tips:
1. **Validation first** - Place validation at low priority numbers (1-2)
2. **Transform middle** - Data transformation in middle priorities
3. **Logging/Debug last** - Logging at high priority numbers (last)

---

## Use Case Examples

### 1. Add CORS Headers
**Action Type:** Run JavaScript
**Execution Point:** after_request
**Filter:** All requests

```javascript
response.headers["Access-Control-Allow-Origin"] = "*";
response.headers["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE";
response.headers["Access-Control-Allow-Headers"] = "Content-Type, Authorization";
```

---

### 2. Hide Credit Card Numbers
**Action Type:** Replace Text
**Execution Point:** after_request

```json
{
  "target": "response_body",
  "pattern": "\\d{4}-\\d{4}-\\d{4}-\\d{4}",
  "replacement": "****-****-****-****",
  "use_regex": true
}
```

---

### 3. Add API Version Header
**Action Type:** Run JavaScript
**Execution Point:** after_request
**Filter:** Path starts with `/api/`

```javascript
response.headers["X-API-Version"] = "2.6.0";
response.headers["X-Powered-By"] = "Beo-Echo";
```

---

### 4. Standardize Error Responses
**Action Type:** Run JavaScript
**Execution Point:** after_request
**Filter:** Status code >= 400

```javascript
if (response.status_code >= 400) {
  var errorResponse = {
    error: true,
    code: response.status_code,
    message: "Request failed",
    timestamp: new Date().toISOString(),
    path: request.path
  };

  response.body = JSON.stringify(errorResponse);
  response.headers["Content-Type"] = "application/json";
}
```

---

### 5. Transform Snake Case to Camel Case
**Action Type:** Run JavaScript
**Execution Point:** before_request

```javascript
var data = JSON.parse(request.body);

// Transform all keys
var transformed = {};
for (var key in data) {
  var camelKey = key.replace(/_([a-z])/g, function(g) {
    return g[1].toUpperCase();
  });
  transformed[camelKey] = data[key];
}

request.body = JSON.stringify(transformed);
```

---

### 6. Add Request ID Tracking
**Action Type:** Run JavaScript
**Execution Point:** before_request

```javascript
// Add unique request ID
request.headers["X-Request-ID"] = "req-" + Date.now() + "-" + Math.random().toString(36).substr(2, 9);

console.log("Request ID:", request.headers["X-Request-ID"]);
```

---

## Tips & Best Practices

### ✅ DO (Recommended)
- **Use clear names** - Descriptive action names (e.g., "Add CORS Headers")
- **Use filters** - Make actions run only when needed
- **Test first** - Test actions with sample data before enabling
- **Use Replace Text** - For simple replacements (faster than JavaScript)
- **Add console.log** - For debugging JavaScript code
- **Disable when not needed** - Toggle off unused actions

### ❌ DON'T (Avoid)
- **Don't create complex scripts** - Keep it simple, max 5 second timeout
- **Don't log sensitive data** - Don't console.log passwords/tokens
- **Don't use too many actions** - Many actions = slower responses
- **Don't forget to parse JSON** - Always JSON.parse before modifying data
- **Don't hardcode values** - Use variables for frequently changing values

---

## Troubleshooting

### Action Not Running?
1. ✓ Ensure action is enabled (green toggle)
2. ✓ Check filters - do they match the request?
3. ✓ Check execution point (before/after request)
4. ✓ Verify priority - is the order correct?

### JavaScript Error?
1. ✓ Check console.log output in debug panel
2. ✓ Ensure JSON.parse and JSON.stringify are correct
3. ✓ Check timeout - script must finish within 5 seconds
4. ✓ Test script in browser console first

### Response Not as Expected?
1. ✓ Check priority order - actions can override each other
2. ✓ Disable actions one by one to test
3. ✓ Ensure target (body/header) is correct
4. ✓ Check regex pattern if using Replace Text

---

## How to Use in UI

1. **Open Project** → Select the project to add actions
2. **Click "Actions" Tab** → View all existing actions
3. **Click "Add Action"** → Wizard will open
4. **Choose Action Type** → Replace Text or Run JavaScript
5. **Fill Configuration:**
   - Action name
   - Execution point (before/after)
   - Configuration (based on action type)
   - Filters (optional)
6. **Save** → Action will be automatically enabled
7. **Test** → Send requests to test the action
8. **Adjust Priority** → Drag & drop to reorder

### UI Features:
- 🔄 **Drag & Drop** - Change priority by dragging
- 🎛️ **Toggle Switch** - Quick enable/disable
- 🗑️ **Delete** - Remove action
- ✏️ **Edit** - Modify configuration
- 🐛 **Debug Panel** - View console.log output (JavaScript)

---

**Happy using Actions! 🚀**

For more questions, contact the development team or refer to the main project documentation.
