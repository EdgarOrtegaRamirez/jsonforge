// Package validator provides JSON Schema validation.
package validator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ValidationError represents a single validation error.
type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// ValidationResult contains the result of validation.
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// Validate validates a JSON value against a JSON Schema.
func Validate(data interface{}, schema map[string]interface{}) *ValidationResult {
	errors := []ValidationError{}
	validateValue(data, schema, "", &errors)
	return &ValidationResult{
		Valid:  len(errors) == 0,
		Errors: errors,
	}
}

func validateValue(data interface{}, schema map[string]interface{}, path string, errors *[]ValidationError) {
	// Check type
	if typeName, ok := schema["type"]; ok {
		if !checkType(data, typeName.(string)) {
			*errors = append(*errors, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("expected type %s, got %s", typeName, getTypeName(data)),
			})
			return
		}
	}

	// Check required fields
	if required, ok := schema["required"]; ok {
		if reqList, ok := required.([]interface{}); ok {
			if obj, ok := data.(map[string]interface{}); ok {
				for _, field := range reqList {
					fieldName := field.(string)
					if _, exists := obj[fieldName]; !exists {
						*errors = append(*errors, ValidationError{
							Path:    path,
							Message: fmt.Sprintf("required field %q is missing", fieldName),
						})
					}
				}
			}
		}
	}

	// Check properties
	if properties, ok := schema["properties"]; ok {
		if propMap, ok := properties.(map[string]interface{}); ok {
			if obj, ok := data.(map[string]interface{}); ok {
				for key, propSchema := range propMap {
					if val, exists := obj[key]; exists {
						childPath := key
						if path != "" {
							childPath = path + "." + key
						}
						if propSchemaMap, ok := propSchema.(map[string]interface{}); ok {
							validateValue(val, propSchemaMap, childPath, errors)
						}
					}
				}
			}
		}
	}

	// Check minimum/maximum for numbers
	if num, ok := data.(float64); ok {
		if min, ok := schema["minimum"]; ok {
			if minVal, ok := min.(float64); ok && num < minVal {
				*errors = append(*errors, ValidationError{
					Path:    path,
					Message: fmt.Sprintf("value %v is less than minimum %v", num, minVal),
				})
			}
		}
		if max, ok := schema["maximum"]; ok {
			if maxVal, ok := max.(float64); ok && num > maxVal {
				*errors = append(*errors, ValidationError{
					Path:    path,
					Message: fmt.Sprintf("value %v is greater than maximum %v", num, maxVal),
				})
			}
		}
	}

	// Check minLength/maxLength for strings
	if str, ok := data.(string); ok {
		if minLen, ok := schema["minLength"]; ok {
			if minL, ok := minLen.(float64); ok && float64(len(str)) < minL {
				*errors = append(*errors, ValidationError{
					Path:    path,
					Message: fmt.Sprintf("string length %d is less than minLength %d", len(str), int(minL)),
				})
			}
		}
		if maxLen, ok := schema["maxLength"]; ok {
			if maxL, ok := maxLen.(float64); ok && float64(len(str)) > maxL {
				*errors = append(*errors, ValidationError{
					Path:    path,
					Message: fmt.Sprintf("string length %d is greater than maxLength %d", len(str), int(maxL)),
				})
			}
		}
		if pattern, ok := schema["pattern"]; ok {
			if patStr, ok := pattern.(string); ok {
				matched, _ := matchSimplePattern(str, patStr)
				if !matched {
					*errors = append(*errors, ValidationError{
						Path:    path,
						Message: fmt.Sprintf("string does not match pattern %q", patStr),
					})
				}
			}
		}
	}

	// Check enum
	if enum, ok := schema["enum"]; ok {
		if enumList, ok := enum.([]interface{}); ok {
			matched := false
			for _, v := range enumList {
				if fmt.Sprintf("%v", data) == fmt.Sprintf("%v", v) {
					matched = true
					break
				}
			}
			if !matched {
				*errors = append(*errors, ValidationError{
					Path:    path,
					Message: fmt.Sprintf("value not in enum %v", enumList),
				})
			}
		}
	}

	// Check array items
	if arr, ok := data.([]interface{}); ok {
		if items, ok := schema["items"]; ok {
			if itemSchema, ok := items.(map[string]interface{}); ok {
				for i, item := range arr {
					childPath := fmt.Sprintf("%s[%d]", path, i)
					validateValue(item, itemSchema, childPath, errors)
				}
			}
		}
		if minItems, ok := schema["minItems"]; ok {
			if minI, ok := minItems.(float64); ok && float64(len(arr)) < minI {
				*errors = append(*errors, ValidationError{
					Path:    path,
					Message: fmt.Sprintf("array has %d items, minimum is %d", len(arr), int(minI)),
				})
			}
		}
		if maxItems, ok := schema["maxItems"]; ok {
			if maxI, ok := maxItems.(float64); ok && float64(len(arr)) > maxI {
				*errors = append(*errors, ValidationError{
					Path:    path,
					Message: fmt.Sprintf("array has %d items, maximum is %d", len(arr), int(maxI)),
				})
			}
		}
	}
}

func checkType(data interface{}, expected string) bool {
	actual := getTypeName(data)
	return actual == expected
}

func getTypeName(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch v.(type) {
	case bool:
		return "boolean"
	case float64, int, int64:
		return "number"
	case string:
		return "string"
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	default:
		return "unknown"
	}
}

// Simple pattern matching (supports basic * and ? wildcards)
func matchSimplePattern(s, pattern string) (bool, error) {
	// Convert simple pattern to regex-like matching
	// * matches any sequence, ? matches single character
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		// No wildcards, exact match
		return s == pattern, nil
	}

	remaining := s
	for i, part := range parts {
		if part == "" {
			continue
		}
		idx := strings.Index(remaining, part)
		if idx == -1 {
			return false, nil
		}
		remaining = remaining[idx+len(part):]
		_ = i
	}
	return true, nil
}

// FormatValidation formats validation errors as text.
func FormatValidation(result *ValidationResult) string {
	if result.Valid {
		return "Valid."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Invalid: %d error(s)\n\n", len(result.Errors)))
	for _, e := range result.Errors {
		path := e.Path
		if path == "" {
			path = "(root)"
		}
		sb.WriteString(fmt.Sprintf("  %s: %s\n", path, e.Message))
	}
	return sb.String()
}

// ParseSchema parses a JSON string into a schema map.
func ParseSchema(input string) (map[string]interface{}, error) {
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(input), &schema); err != nil {
		return nil, err
	}
	return schema, nil
}
