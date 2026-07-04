// Package stats provides JSON statistics analysis.
package stats

import (
	"fmt"
	"sort"
)

// Stats contains statistics about a JSON value.
type Stats struct {
	Type       string         `json:"type"`
	Size       int            `json:"size"`
	Depth      int            `json:"depth"`
	KeyCount   int            `json:"key_count,omitempty"`
	ArrayLen   int            `json:"array_length,omitempty"`
	LeafCount  int            `json:"leaf_count"`
	PathCount  int            `json:"path_count"`
	TypeCounts map[string]int `json:"type_counts"`
	Paths      []string       `json:"paths,omitempty"`
}

// Analyze computes statistics for a JSON value.
func Analyze(data interface{}, includePaths bool) *Stats {
	stats := &Stats{
		TypeCounts: make(map[string]int),
		Type:       getTypeName(data),
	}
	analyzeValue(data, "", 0, stats, includePaths)
	return stats
}

func analyzeValue(data interface{}, path string, depth int, stats *Stats, includePaths bool) {
	stats.TypeCounts[getTypeName(data)]++

	if depth > stats.Depth {
		stats.Depth = depth
	}

	switch v := data.(type) {
	case map[string]interface{}:
		stats.KeyCount += len(v)
		stats.Size++
		for key, val := range v {
			childPath := key
			if path != "" {
				childPath = path + "." + key
			}
			analyzeValue(val, childPath, depth+1, stats, includePaths)
		}
	case []interface{}:
		stats.ArrayLen = len(v)
		stats.Size++
		for i, val := range v {
			childPath := fmt.Sprintf("%s[%d]", path, i)
			analyzeValue(val, childPath, depth+1, stats, includePaths)
		}
	default:
		stats.LeafCount++
		stats.PathCount++
		if includePaths && path != "" {
			stats.Paths = append(stats.Paths, path)
		}
	}
}

func getTypeName(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch v.(type) {
	case bool:
		return "boolean"
	case float64:
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

// FormatText formats stats as human-readable text.
func FormatText(s *Stats) string {
	result := fmt.Sprintf("Type:        %s\n", s.Type)
	result += fmt.Sprintf("Depth:       %d\n", s.Depth)
	result += fmt.Sprintf("Size:        %d\n", s.Size)
	if s.KeyCount > 0 {
		result += fmt.Sprintf("Keys:        %d\n", s.KeyCount)
	}
	if s.ArrayLen > 0 {
		result += fmt.Sprintf("Array len:   %d\n", s.ArrayLen)
	}
	result += fmt.Sprintf("Leaf values: %d\n", s.LeafCount)
	result += fmt.Sprintf("Paths:       %d\n", s.PathCount)
	result += "\nType distribution:\n"

	types := make([]string, 0, len(s.TypeCounts))
	for t := range s.TypeCounts {
		types = append(types, t)
	}
	sort.Strings(types)

	for _, t := range types {
		result += fmt.Sprintf("  %-10s %d\n", t, s.TypeCounts[t])
	}

	return result
}

// FormatCompact returns a single-line summary.
func FormatCompact(s *Stats) string {
	return fmt.Sprintf("%s | depth=%d | keys=%d | leaves=%d | paths=%d",
		s.Type, s.Depth, s.KeyCount, s.LeafCount, s.PathCount)
}
