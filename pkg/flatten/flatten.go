// Package flatten provides JSON flattening and unflattening.
package flatten

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Flatten converts nested JSON to dot-notation keys.
func Flatten(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	flattenValue(data, "", result)
	return result
}

func flattenValue(data interface{}, prefix string, result map[string]interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			childPath := key
			if prefix != "" {
				childPath = prefix + "." + key
			}
			flattenValue(val, childPath, result)
		}
	case []interface{}:
		for i, val := range v {
			childPath := fmt.Sprintf("%s[%d]", prefix, i)
			flattenValue(val, childPath, result)
		}
	default:
		result[prefix] = data
	}
}

// Unflatten converts dot-notation keys back to nested JSON.
func Unflatten(data map[string]interface{}) interface{} {
	if len(data) == 0 {
		return nil
	}

	result := make(map[string]interface{})

	for key, value := range data {
		parts := parseFlattenedKey(key)
		setValue(result, parts, value)
	}

	return result
}

func parseFlattenedKey(key string) []string {
	var parts []string
	remaining := key

	for remaining != "" {
		// Find next dot or bracket
		dotIdx := strings.Index(remaining, ".")
		bracketIdx := strings.Index(remaining, "[")

		if dotIdx == -1 && bracketIdx == -1 {
			parts = append(parts, remaining)
			break
		}

		if dotIdx == -1 || (bracketIdx != -1 && bracketIdx < dotIdx) {
			// Bracket comes first
			end := strings.Index(remaining[bracketIdx:], "]")
			if end == -1 {
				parts = append(parts, remaining)
				break
			}
			end += bracketIdx + 1
			parts = append(parts, remaining[:end])
			if end < len(remaining) && remaining[end] == '.' {
				end++
			}
			remaining = remaining[end:]
		} else {
			// Dot comes first
			parts = append(parts, remaining[:dotIdx])
			remaining = remaining[dotIdx+1:]
		}
	}

	return parts
}

func setValue(data interface{}, parts []string, value interface{}) interface{} {
	if len(parts) == 0 {
		return value
	}

	part := parts[0]
	rest := parts[1:]

	// Check if this is an array index like [0]
	if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
		idxStr := part[1 : len(part)-1]
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			return data
		}

		arr, ok := data.([]interface{})
		if !ok {
			arr = make([]interface{}, 0)
		}

		// Extend array if needed
		for len(arr) <= idx {
			arr = append(arr, nil)
		}

		if len(rest) == 0 {
			arr[idx] = value
		} else {
			child := arr[idx]
			if child == nil {
				child = make(map[string]interface{})
			}
			arr[idx] = setValue(child, rest, value)
		}
		return arr
	}

	// Check if this is a mixed key like items[0]
	if bracketIdx := strings.Index(part, "["); bracketIdx > 0 {
		key := part[:bracketIdx]
		bracketContent := part[bracketIdx:]

		obj, ok := data.(map[string]interface{})
		if !ok {
			obj = make(map[string]interface{})
		}

		child, exists := obj[key]
		if !exists {
			child = make([]interface{}, 0)
		}

		// Parse the bracket content
		allParts := append([]string{bracketContent}, rest...)
		obj[key] = setValue(child, allParts, value)
		return obj
	}

	// Object key
	obj, ok := data.(map[string]interface{})
	if !ok {
		obj = make(map[string]interface{})
	}

	if len(rest) == 0 {
		obj[part] = value
	} else {
		child, exists := obj[part]
		if !exists {
			// Check if next part is array index
			if len(rest) > 0 && strings.HasPrefix(rest[0], "[") {
				child = make([]interface{}, 0)
			} else {
				child = make(map[string]interface{})
			}
		}
		obj[part] = setValue(child, rest, value)
	}

	return obj
}

// Keys returns sorted list of flattened keys.
func Keys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
