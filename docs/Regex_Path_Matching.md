# Regex Support Documentation

## Overview
The mock repository now supports multiple path matching patterns with a priority-based scoring system for efficient and flexible API endpoint matching.

## Quick Start

### Basic Wildcard Example (Most Common Use Case)
```
✅ RECOMMENDED: Use wildcard (*) for simple variable segments
Endpoint: /api/v2/customer_rooms/*/broadcast_history
✅ Matches: /api/v2/customer_rooms/331931307/broadcast_history
✅ Matches: /api/v2/customer_rooms/salam/broadcast_history
```

### When to Use Each Pattern Type
- **Wildcard (`*`)**: Simple variable segments (IDs, names, slugs)
- **Path Parameters (`:id`)**: When you need named parameters in your code
- **Regex (`\d+`, `\w+`)**: Complex validation (numbers only, specific formats)
- **Exact Match**: Static paths that never change

## Supported Pattern Types

### 1. Exact Match (Highest Priority - Score: 100)
```
Endpoint: /api/v2/users
Request:  /api/v2/users
✅ Match - Score: 100
```

### 2. Path Parameters (Score: 80 + bonus)
```
Endpoint: /users/:id/profile
Request:  /users/123/profile
✅ Match - Score: 98 (80 base + 10 exact + 8 param)
```

### 3. Wildcard Patterns (Score: 60 + bonus)
```
Endpoint: /api/v2/customer_rooms/*/broadcast_history
Request:  /api/v2/customer_rooms/331931307/broadcast_history
Request:  /api/v2/customer_rooms/salam/broadcast_history
✅ Both Match - Score: 106 (60 base + 40 exact + 6 wildcard)
```

### 4. Regex Patterns (Score: 40)
```
Endpoint: /api/v\d+/users/\d+
Request:  /api/v2/users/123
✅ Match - Score: 40

Endpoint: /api/v[12]/customer_rooms/\w+/broadcast_history
Request:  /api/v2/customer_rooms/salam/broadcast_history
✅ Match - Score: 40
```

## How It Works

### Priority System
1. **Exact Match** always wins (score: 100)
2. **Path Parameters** have high priority (score: 80+)
3. **Wildcard** patterns come next (score: 60+)
4. **Regex** patterns have lowest priority (score: 40)

### Scoring Calculation
- **Exact segment match**: +10 points each
- **Path parameter match** (`:id`): +8 points each
- **Wildcard match** (`*`): +6 points each
- **Base scores**: Exact (100), Parameters (80), Wildcards (60), Regex (40)

### Performance Optimization
- **Fast rejection**: Invalid patterns return -1 immediately
- **Early winner**: Exact matches (score 100) stop further evaluation
- **Simple operations**: Uses string operations before regex compilation
- **Regex detection**: Only compiles regex when pattern contains metacharacters

## Usage Examples

### Customer Rooms Use Case (Your Original Request)
```go
// Setup in mock endpoint configuration
Endpoint: "/api/v2/customer_rooms/*/broadcast_history"

// ✅ These requests will match:
GET /api/v2/customer_rooms/331931307/broadcast_history
GET /api/v2/customer_rooms/salam/broadcast_history
GET /api/v2/customer_rooms/room-abc-123/broadcast_history

// ❌ These will NOT match:
GET /api/v2/customer_rooms/broadcast_history        (missing segment)
GET /api/v2/customer_rooms/123/456/broadcast_history (extra segment)
```

### Multiple Pattern Competition
```go
// If you have multiple endpoints:
endpoints := []string{
    "/api/v2/customer_rooms/*/broadcast_history",  // Wildcard (Score: 106)
    "/api/v\\d+/customer_rooms/\\d+/history",      // Regex (Score: 40)
}

// Request: "/api/v2/customer_rooms/123/broadcast_history"
// → Wildcard wins (106 > 40)

// Request: "/api/v3/customer_rooms/123/history" 
// → Regex wins (only option that matches)
```

### Real-World API Patterns
```go
// E-commerce API examples
"/api/v1/products/*/reviews"          // Any product ID
"/api/v1/users/*/orders"              // Any user ID  
"/api/v1/categories/*/products/*/info" // Category + Product

// Validation with Regex  
"/api/v\\d+/users/\\d+/orders"        // Only numeric IDs
"/api/v1/products/[a-z0-9\\-]+/info"  // Slugified product names
```

### Mixed Pattern Priority
```go
endpoints := []string{
    "/users",              // Exact match
    "/users/:id",          // Path parameter  
    "/users/*",            // Wildcard
    "/users/\\d+",         // Regex
}

// Request: "/users" → Exact match wins (score: 100)
// Request: "/users/123" → Parameter wins (score: 98 vs 76 vs 40)
// Request: "/users/abc" → Wildcard wins (score: 76 vs 40)
```

## Implementation Benefits

1. **Simple to Use**: Just define patterns as strings
2. **Efficient**: Fast path matching with minimal regex compilation
3. **Predictable**: Clear priority system for pattern resolution
4. **Flexible**: Supports simple wildcards to complex regex patterns
5. **Backward Compatible**: Existing path parameters continue to work

## Regex Pattern Guidelines

### Supported Regex Metacharacters
- `\d+` - One or more digits (123, 456789)
- `\d*` - Zero or more digits (123, "", empty)
- `\w+` - One or more word characters (abc, user123, test_name)
- `[abc]` - Character classes (a, b, or c)
- `[0-9]+` - Numeric range (equivalent to \d+)
- `[a-z]+` - Lowercase letters only
- `(a|b)` - Groups and alternation (a or b)
- `^$` - Anchors (automatically added by system)
- `+?*` - Quantifiers (one or more, zero or one, zero or more)

### Pattern Examples with Explanations
```
/api/v\d+/users           → /api/v1/users, /api/v22/users
/products/[a-z]+          → /products/laptop, /products/phone
/users/\d+/orders/\d+     → /users/123/orders/456
/files/[\w\-\.]+          → /files/doc.pdf, /files/my-file_v2.txt
/(en|id|fr)/products      → /en/products, /id/products, /fr/products
```

### Best Practices
1. **Use wildcards first**: `*` for simple cases, regex for validation
2. **Keep it simple**: Complex regex can slow down matching
3. **Test your patterns**: Use the test endpoints to verify behavior
4. **Escape special chars**: Use `\\` for literal backslashes in JSON
5. **Consider performance**: Wildcard (score 60+) beats regex (score 40)

### ⚠️ Common Pitfalls
```
❌ BAD:  /api/v\d/users     (matches /api/v1/users but not /api/v10/users)
✅ GOOD: /api/v\d+/users    (matches both /api/v1/users and /api/v10/users)

❌ BAD:  /users/\w/profile  (only single character)
✅ GOOD: /users/\w+/profile (one or more characters)

❌ BAD:  /api/v[1-9]/users  (doesn't match v10, v11, etc.)
✅ GOOD: /api/v\d+/users    (matches any version number)
```

## Testing
Comprehensive test suite covers:
- All pattern types and combinations
- Priority resolution between different patterns
- Edge cases and error conditions
- Performance with large endpoint lists

### Testing Your Patterns
```bash
# Run the regex support tests
cd backend && go test ./src/echo/repositories -v

# Test specific pattern matching
go test ./src/echo/repositories -run "TestCalculatePathMatchScore"

# Test endpoint resolution
go test ./src/echo/repositories -run "TestFindBestPathMatch"
```

## Troubleshooting

### Pattern Not Matching?
1. **Check the syntax**: Ensure regex metacharacters are properly escaped
2. **Verify priority**: Higher scoring patterns win (exact > param > wildcard > regex)
3. **Test segments**: Each path segment must match (same number of `/` separators)
4. **Case sensitivity**: Patterns are case-sensitive

### Performance Issues?
1. **Prefer wildcards**: Use `*` instead of `.*` when possible
2. **Optimize regex**: Simple patterns like `\d+` are faster than complex ones
3. **Limit endpoints**: More endpoints = longer search time
4. **Use exact matches**: They get score 100 and stop further evaluation

### Debug Pattern Matching
```go
// Add debug logging to see which patterns are being evaluated
func debugPathMatch(endpointPath, requestPath string) {
    score := calculatePathMatchScore(endpointPath, requestPath)
    fmt.Printf("Pattern: %s, Request: %s, Score: %d\n", 
               endpointPath, requestPath, score)
}
```

## Migration Guide

### From Old System
If you're upgrading from the previous path matching system:

1. **Existing patterns still work**: `:id` style parameters unchanged
2. **New wildcard support**: Replace complex patterns with simple `*`
3. **Regex for validation**: Use regex only when you need strict validation
4. **Test thoroughly**: Verify your endpoints resolve as expected

### Recommended Upgrade Path
```go
// Old approach - complex regex for simple cases
"/api/v\\d+/users/.*"              

// New approach - use wildcards for simplicity  
"/api/v*/users/*"                   // Simpler and higher priority

// Keep regex only for strict validation
"/api/v\\d+/users/\\d+"            // When you need numeric validation
```
