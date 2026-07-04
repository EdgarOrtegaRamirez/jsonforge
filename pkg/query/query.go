// Package query provides dot-notation JSON path queries.
package query

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// PathSegment represents one segment of a JSON path.
type PathSegment struct {
	Key       string
	Index     int
	Last      bool
	Wildcard  bool
	Recursive bool
}

// ParsePath parses a dot-notation path into segments.
func ParsePath(path string) ([]PathSegment, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	var segments []PathSegment
	i := 0
	n := len(path)

	for i < n {
		// Check for recursive descent ".."
		if i+1 < n && path[i] == '.' && path[i+1] == '.' {
			segments = append(segments, PathSegment{Recursive: true})
			i += 2
			// Skip optional trailing dot after ".."
			if i < n && path[i] == '.' {
				i++
			}
			continue
		}

		// Skip single dot separator
		if path[i] == '.' {
			i++
			continue
		}

		// Find end of current segment
		start := i
		inBracket := false
		for i < n {
			ch := path[i]
			if ch == '[' {
				// If we're not in a bracket and haven't collected anything, this starts a bracket segment
				if !inBracket && i > start {
					break
				}
				inBracket = true
			} else if ch == ']' {
				inBracket = false
				i++
				break
			} else if ch == '.' && !inBracket {
				break
			}
			i++
		}

		segment := path[start:i]
		if segment != "" {
			seg, err := parseSegment(segment)
			if err != nil {
				return nil, err
			}
			segments = append(segments, seg)
		}
	}

	return segments, nil
}

func parseSegment(s string) (PathSegment, error) {
	if s == "" {
		return PathSegment{}, fmt.Errorf("empty path segment")
	}

	// Wildcard
	if s == "*" {
		return PathSegment{Wildcard: true}, nil
	}

	// Bracket notation: ["key"] or [key] or [0] or [last] or [-1]
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		inner := s[1 : len(s)-1]
		inner = strings.TrimSpace(inner)

		// Remove quotes
		if len(inner) >= 2 && (inner[0] == '"' || inner[0] == '\'') {
			inner = inner[1 : len(inner)-1]
			return PathSegment{Key: inner}, nil
		}

		if inner == "last" || inner == "-1" {
			return PathSegment{Last: true, Index: -1}, nil
		}

		if i, err := strconv.Atoi(inner); err == nil {
			return PathSegment{Index: i}, nil
		}

		return PathSegment{Key: inner}, nil
	}

	// Mixed notation: key[number] or key[last]
	if idx := strings.Index(s, "["); idx > 0 && strings.HasSuffix(s, "]") {
		key := s[:idx]
		inner := s[idx+1 : len(s)-1]
		inner = strings.TrimSpace(inner)

		if inner == "last" || inner == "-1" {
			return PathSegment{Key: key, Last: true, Index: -1}, nil
		}
		if i, err := strconv.Atoi(inner); err == nil {
			return PathSegment{Key: key, Index: i}, nil
		}
		return PathSegment{Key: key}, nil
	}

	// "last" keyword
	if s == "last" {
		return PathSegment{Last: true, Index: -1}, nil
	}

	// Numeric index
	if i, err := strconv.Atoi(s); err == nil {
		return PathSegment{Index: i}, nil
	}

	return PathSegment{Key: s}, nil
}

// Get retrieves a value from JSON data using a dot-notation path.
func Get(data interface{}, path string) (interface{}, error) {
	segments, err := ParsePath(path)
	if err != nil {
		return nil, err
	}
	return resolveSegments(data, segments)
}

// Set sets a value in JSON data at the given path.
func Set(data interface{}, path string, value interface{}) (interface{}, error) {
	segments, err := ParsePath(path)
	if err != nil {
		return nil, err
	}
	return setSegments(data, segments, value)
}

// Delete removes a value from JSON data at the given path.
func Delete(data interface{}, path string) (interface{}, error) {
	segments, err := ParsePath(path)
	if err != nil {
		return nil, err
	}
	return deleteSegments(data, segments)
}

// Exists checks if a path exists in JSON data.
func Exists(data interface{}, path string) bool {
	_, err := Get(data, path)
	return err == nil
}

func resolveSegments(data interface{}, segments []PathSegment) (interface{}, error) {
	if len(segments) == 0 {
		return data, nil
	}

	seg := segments[0]
	rest := segments[1:]

	if seg.Recursive {
		return resolveRecursive(data, rest)
	}

	switch current := data.(type) {
	case map[string]interface{}:
		if seg.Wildcard {
			var results []interface{}
			for _, v := range current {
				if len(rest) == 0 {
					results = append(results, v)
				} else {
					val, err := resolveSegments(v, rest)
					if err == nil {
						results = append(results, val)
					}
				}
			}
			return results, nil
		}
		val, ok := current[seg.Key]
		if !ok {
			return nil, fmt.Errorf("key %q not found", seg.Key)
		}
		return resolveSegments(val, rest)

	case []interface{}:
		if seg.Wildcard {
			var results []interface{}
			for _, v := range current {
				if len(rest) == 0 {
					results = append(results, v)
				} else {
					val, err := resolveSegments(v, rest)
					if err == nil {
						results = append(results, val)
					}
				}
			}
			return results, nil
		}
		idx := seg.Index
		if seg.Last || idx == -1 {
			idx = len(current) - 1
		}
		if idx < 0 || idx >= len(current) {
			return nil, fmt.Errorf("index %d out of range (len=%d)", seg.Index, len(current))
		}
		return resolveSegments(current[idx], rest)

	default:
		return nil, fmt.Errorf("cannot navigate into %T", data)
	}
}

func resolveRecursive(data interface{}, rest []PathSegment) (interface{}, error) {
	var results []interface{}

	var walk func(node interface{})
	walk = func(node interface{}) {
		switch v := node.(type) {
		case map[string]interface{}:
			if len(rest) > 0 {
				val, err := resolveSegments(v, rest)
				if err == nil {
					results = append(results, val)
				}
			} else {
				results = append(results, v)
			}
			for _, child := range v {
				walk(child)
			}
		case []interface{}:
			if len(rest) > 0 {
				val, err := resolveSegments(v, rest)
				if err == nil {
					results = append(results, val)
				}
			} else {
				results = append(results, v)
			}
			for _, child := range v {
				walk(child)
			}
		}
	}

	walk(data)

	if len(results) == 0 {
		return nil, fmt.Errorf("no matches found for recursive path")
	}
	if len(results) == 1 {
		return results[0], nil
	}
	return results, nil
}

func setSegments(data interface{}, segments []PathSegment, value interface{}) (interface{}, error) {
	if len(segments) == 0 {
		return value, nil
	}

	seg := segments[0]
	rest := segments[1:]

	if seg.Wildcard || seg.Recursive {
		return nil, fmt.Errorf("cannot use wildcards or recursive paths in set operations")
	}

	switch current := data.(type) {
	case map[string]interface{}:
		if len(rest) == 0 {
			current[seg.Key] = value
			return current, nil
		}
		child, ok := current[seg.Key]
		if !ok {
			child = make(map[string]interface{})
			current[seg.Key] = child
		}
		newChild, err := setSegments(child, rest, value)
		if err != nil {
			return nil, err
		}
		current[seg.Key] = newChild
		return current, nil

	case []interface{}:
		idx := seg.Index
		if seg.Last || idx == -1 {
			idx = len(current) - 1
		}
		if idx < 0 || idx >= len(current) {
			return nil, fmt.Errorf("index %d out of range (len=%d)", seg.Index, len(current))
		}
		if len(rest) == 0 {
			current[idx] = value
			return current, nil
		}
		newChild, err := setSegments(current[idx], rest, value)
		if err != nil {
			return nil, err
		}
		current[idx] = newChild
		return current, nil

	default:
		return nil, fmt.Errorf("cannot set on %T", data)
	}
}

func deleteSegments(data interface{}, segments []PathSegment) (interface{}, error) {
	if len(segments) == 0 {
		return data, nil
	}

	seg := segments[0]
	rest := segments[1:]

	if seg.Wildcard || seg.Recursive {
		return nil, fmt.Errorf("cannot use wildcards or recursive paths in delete operations")
	}

	switch current := data.(type) {
	case map[string]interface{}:
		if len(rest) == 0 {
			delete(current, seg.Key)
			return current, nil
		}
		child, ok := current[seg.Key]
		if !ok {
			return current, nil
		}
		newChild, err := deleteSegments(child, rest)
		if err != nil {
			return nil, err
		}
		current[seg.Key] = newChild
		return current, nil

	case []interface{}:
		idx := seg.Index
		if seg.Last || idx == -1 {
			idx = len(current) - 1
		}
		if idx < 0 || idx >= len(current) {
			return current, nil
		}
		if len(rest) == 0 {
			return append(current[:idx], current[idx+1:]...), nil
		}
		newChild, err := deleteSegments(current[idx], rest)
		if err != nil {
			return nil, err
		}
		current[idx] = newChild
		return current, nil

	default:
		return data, nil
	}
}

// ParseJSON parses a JSON string into an interface{}.
func ParseJSON(input string) (interface{}, error) {
	var result interface{}
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.UseNumber()
	err := decoder.Decode(&result)
	return result, err
}
