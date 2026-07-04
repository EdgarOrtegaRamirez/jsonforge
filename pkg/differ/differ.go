// Package differ provides semantic JSON diffing.
package differ

import (
	"fmt"
	"sort"
	"strings"
)

// ChangeType represents the type of change.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

// Change represents a single diff change.
type Change struct {
	Path     string      `json:"path"`
	Type     ChangeType  `json:"type"`
	OldValue interface{} `json:"old_value,omitempty"`
	NewValue interface{} `json:"new_value,omitempty"`
}

// DiffResult contains the result of diffing two JSON values.
type DiffResult struct {
	Changes []Change `json:"changes"`
	Count   int      `json:"count"`
	Same    bool     `json:"same"`
}

// Diff computes the semantic diff between two JSON values.
func Diff(old, new interface{}) *DiffResult {
	changes := []Change{}
	diffValue("", old, new, &changes)
	return &DiffResult{
		Changes: changes,
		Count:   len(changes),
		Same:    len(changes) == 0,
	}
}

func diffValue(path string, old, new interface{}, changes *[]Change) {
	if old == nil && new == nil {
		return
	}

	if old == nil {
		*changes = append(*changes, Change{Path: path, Type: Added, NewValue: new})
		return
	}

	if new == nil {
		*changes = append(*changes, Change{Path: path, Type: Removed, OldValue: old})
		return
	}

	oldMap, oldIsMap := old.(map[string]interface{})
	newMap, newIsMap := new.(map[string]interface{})

	if oldIsMap && newIsMap {
		diffMaps(path, oldMap, newMap, changes)
		return
	}

	oldArr, oldIsArr := old.([]interface{})
	newArr, newIsArr := new.([]interface{})

	if oldIsArr && newIsArr {
		diffArrays(path, oldArr, newArr, changes)
		return
	}

	// Primitive comparison
	if fmt.Sprintf("%v", old) != fmt.Sprintf("%v", new) {
		*changes = append(*changes, Change{
			Path:     path,
			Type:     Modified,
			OldValue: old,
			NewValue: new,
		})
	}
}

func diffMaps(path string, old, new map[string]interface{}, changes *[]Change) {
	// Find added and modified keys
	for key, newVal := range new {
		keyPath := joinPath(path, key)
		oldVal, exists := old[key]
		if !exists {
			*changes = append(*changes, Change{Path: keyPath, Type: Added, NewValue: newVal})
		} else {
			diffValue(keyPath, oldVal, newVal, changes)
		}
	}

	// Find removed keys
	for key := range old {
		if _, exists := new[key]; !exists {
			*changes = append(*changes, Change{
				Path:     joinPath(path, key),
				Type:     Removed,
				OldValue: old[key],
			})
		}
	}
}

func diffArrays(path string, old, new []interface{}, changes *[]Change) {
	maxLen := len(old)
	if len(new) > maxLen {
		maxLen = len(new)
	}

	for i := 0; i < maxLen; i++ {
		itemPath := fmt.Sprintf("%s[%d]", path, i)
		if i >= len(old) {
			*changes = append(*changes, Change{Path: itemPath, Type: Added, NewValue: new[i]})
		} else if i >= len(new) {
			*changes = append(*changes, Change{Path: itemPath, Type: Removed, OldValue: old[i]})
		} else {
			diffValue(itemPath, old[i], new[i], changes)
		}
	}
}

func joinPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}

// FormatText formats a diff result as human-readable text.
func FormatText(result *DiffResult) string {
	if result.Same {
		return "No differences found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d difference(s):\n\n", result.Count))

	for _, c := range result.Changes {
		switch c.Type {
		case Added:
			sb.WriteString(fmt.Sprintf("  + %s: %v\n", c.Path, c.NewValue))
		case Removed:
			sb.WriteString(fmt.Sprintf("  - %s: %v\n", c.Path, c.OldValue))
		case Modified:
			sb.WriteString(fmt.Sprintf("  ~ %s: %v → %v\n", c.Path, c.OldValue, c.NewValue))
		}
	}

	return sb.String()
}

// FormatJSON formats a diff result as JSON-like output.
func FormatJSON(result *DiffResult) string {
	if result.Same {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[\n")
	for i, c := range result.Changes {
		sb.WriteString(fmt.Sprintf(`  {"path": %q, "type": %q`, c.Path, c.Type))
		if c.OldValue != nil {
			sb.WriteString(fmt.Sprintf(`, "old": %v`, c.OldValue))
		}
		if c.NewValue != nil {
			sb.WriteString(fmt.Sprintf(`, "new": %v`, c.NewValue))
		}
		sb.WriteString("}")
		if i < len(result.Changes)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]")
	return sb.String()
}

// Summary returns a summary of changes by type.
func Summary(result *DiffResult) map[string]int {
	summary := map[string]int{
		"added":    0,
		"removed":  0,
		"modified": 0,
	}
	for _, c := range result.Changes {
		summary[string(c.Type)]++
	}
	return summary
}

// Paths returns sorted list of changed paths.
func Paths(result *DiffResult) []string {
	paths := make([]string, len(result.Changes))
	for i, c := range result.Changes {
		paths[i] = c.Path
	}
	sort.Strings(paths)
	return paths
}
