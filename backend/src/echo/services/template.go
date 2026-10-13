package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/aymerick/raymond/ast"
	"github.com/aymerick/raymond/parser"

	"beo-echo/backend/src/database"
)

// maxTemplateRepeat caps the iteration count of the repeat helper so a request such as
// ?n=99999999 cannot exhaust memory.
const maxTemplateRepeat = 10000

// renderResponseBody renders the body of a templated mock response (Handlebars) for the
// incoming request. Helpers: queryParam, header, repeat.
func renderResponseBody(mockResp database.MockResponse, req *http.Request) (string, error) {
	tpl, err := raymond.Parse(mockResp.Body)
	if err != nil {
		return "", fmt.Errorf("invalid template: %w", err)
	}
	if err := rejectUnknownHelpers(mockResp.Body); err != nil {
		return "", err
	}

	jsonOutput := isJSONContentType(mockResp.Headers)
	tpl.RegisterHelpers(map[string]interface{}{
		// {{queryParam 'name'}} or {{queryParam 'name' default='x'}}
		"queryParam": func(name string, options *raymond.Options) raymond.SafeString {
			values, present := req.URL.Query()[name]
			value := ""
			if present && len(values) > 0 {
				value = values[0]
			} else if def := options.HashStr("default"); def != "" {
				value = def
			}
			return raymond.SafeString(escapeForOutput(value, jsonOutput))
		},
		// {{header 'X-Name'}} or {{header 'X-Name' default='x'}}
		"header": func(name string, options *raymond.Options) raymond.SafeString {
			value := req.Header.Get(name)
			if value == "" {
				value = options.HashStr("default")
			}
			return raymond.SafeString(escapeForOutput(value, jsonOutput))
		},
		// {{#repeat n comma=true}}...{{@index}}...{{/repeat}}
		"repeat": func(n interface{}, options *raymond.Options) raymond.SafeString {
			count := toRepeatCount(n)
			if count > maxTemplateRepeat {
				count = maxTemplateRepeat
			}
			separator := ""
			if comma, _ := options.HashProp("comma").(bool); comma {
				separator = ","
			}

			var out strings.Builder
			for i := 0; i < count; i++ {
				frame := options.NewDataFrame()
				frame.Set("index", i)
				frame.Set("key", i)
				frame.Set("first", i == 0)
				frame.Set("last", i == count-1)
				if i > 0 {
					out.WriteString(separator)
				}
				out.WriteString(options.FnData(frame))
			}
			return raymond.SafeString(out.String())
		},
	})

	result, err := tpl.Exec(nil)
	if err != nil {
		return "", fmt.Errorf("template error: %w", err)
	}
	return result, nil
}

// allowedHelpers are the helpers a template may call: ours plus raymond's safe built-ins.
var allowedHelpers = map[string]bool{
	"queryParam": true, "header": true, "repeat": true,
	"if": true, "unless": true, "each": true, "with": true, "lookup": true, "equal": true,
}

// rejectUnknownHelpers fails on a helper call such as {{queryParms 'x'}} (typo) or {{log 'x'}}.
// raymond would silently render an unknown helper as an empty string, which hides mistakes.
func rejectUnknownHelpers(source string) error {
	program, err := parser.Parse(source)
	if err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}
	return checkHelpers(program)
}

func checkHelpers(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
		if n == nil {
			return nil
		}
		for _, stmt := range n.Body {
			if err := checkHelpers(stmt); err != nil {
				return err
			}
		}
	case *ast.MustacheStatement:
		return checkExpression(n.Expression)
	case *ast.BlockStatement:
		if err := checkExpression(n.Expression); err != nil {
			return err
		}
		if err := checkHelpers(n.Program); err != nil {
			return err
		}
		if n.Inverse != nil {
			return checkHelpers(n.Inverse)
		}
	case *ast.SubExpression:
		return checkExpression(n.Expression)
	}
	return nil
}

func checkExpression(expr *ast.Expression) error {
	if expr == nil {
		return nil
	}
	if len(expr.Params) > 0 || expr.Hash != nil {
		if path, ok := expr.Path.(*ast.PathExpression); ok && !allowedHelpers[path.Original] {
			return fmt.Errorf("unknown template helper %q (available: queryParam, header, repeat)", path.Original)
		}
	}
	for _, param := range expr.Params {
		if err := checkHelpers(param); err != nil {
			return err
		}
	}
	return nil
}

// toRepeatCount converts a repeat helper argument ("500", 500, ...) into a non-negative count.
func toRepeatCount(v interface{}) int {
	n, err := strconv.Atoi(strings.TrimSpace(fmt.Sprint(v)))
	if err != nil || n < 0 {
		return 0
	}
	return n
}

// isJSONContentType reports whether the response headers (JSON string) declare a JSON body.
func isJSONContentType(headersJSON string) bool {
	var headers map[string]string
	if err := json.Unmarshal([]byte(headersJSON), &headers); err != nil {
		return false
	}
	for key, value := range headers {
		if strings.EqualFold(key, "content-type") {
			return strings.Contains(strings.ToLower(value), "json")
		}
	}
	return false
}

// escapeForOutput makes a request value safe to embed inside a JSON string literal when the
// response is JSON, so a query value such as a"b cannot break the document. Other content
// types get the value as-is.
func escapeForOutput(value string, jsonOutput bool) string {
	if !jsonOutput {
		return value
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(value); err != nil {
		return value
	}
	encoded := strings.TrimSuffix(buf.String(), "\n")
	return encoded[1 : len(encoded)-1] // drop the surrounding quotes
}

// buildMockResponse renders the template when the response is templated, then builds the HTTP response.
func buildMockResponse(mockResp database.MockResponse, req *http.Request) (*http.Response, error) {
	if mockResp.Templated {
		rendered, err := renderResponseBody(mockResp, req)
		if err != nil {
			return createErrorResponse(http.StatusInternalServerError, err.Error()), nil
		}
		mockResp.Body = rendered
	}
	return createMockResponse(mockResp)
}
