// Package converter provides JSON format conversion.
package converter

import (
	"encoding/csv"
	"fmt"
	"sort"
	"strings"
)

// ToYAML converts JSON data to YAML format.
func ToYAML(data interface{}) string {
	var sb strings.Builder
	writeYAML(&sb, data, 0)
	return sb.String()
}

func writeYAML(sb *strings.Builder, data interface{}, indent int) {
	prefix := strings.Repeat("  ", indent)

	switch v := data.(type) {
	case nil:
		sb.WriteString("null\n")
	case bool:
		if v {
			sb.WriteString("true\n")
		} else {
			sb.WriteString("false\n")
		}
	case float64:
		if v == float64(int(v)) {
			sb.WriteString(fmt.Sprintf("%d\n", int(v)))
		} else {
			sb.WriteString(fmt.Sprintf("%g\n", v))
		}
	case string:
		if strings.ContainsAny(v, "\n:\"#{}[]|>&*!%@`") || strings.TrimSpace(v) != v {
			sb.WriteString(fmt.Sprintf("%q\n", v))
		} else {
			sb.WriteString(v + "\n")
		}
	case map[string]interface{}:
		if len(v) == 0 {
			sb.WriteString("{}\n")
			return
		}
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			val := v[k]
			sb.WriteString(prefix + k + ": ")
			switch val.(type) {
			case map[string]interface{}, []interface{}:
				sb.WriteString("\n")
				writeYAML(sb, val, indent+1)
			default:
				writeYAML(sb, val, 0)
			}
		}
	case []interface{}:
		if len(v) == 0 {
			sb.WriteString("[]\n")
			return
		}
		for _, item := range v {
			sb.WriteString(prefix + "- ")
			switch item.(type) {
			case map[string]interface{}:
				sb.WriteString("\n")
				writeYAML(sb, item, indent+1)
			default:
				line := fmt.Sprintf("%v", item)
				sb.WriteString(line + "\n")
			}
		}
	default:
		sb.WriteString(fmt.Sprintf("%v\n", v))
	}
}

// ToTOML converts JSON data to TOML format.
func ToTOML(data interface{}) string {
	var sb strings.Builder
	if obj, ok := data.(map[string]interface{}); ok {
		writeTOMLTable(&sb, obj, "")
	} else {
		sb.WriteString(fmt.Sprintf("value = %v\n", data))
	}
	return sb.String()
}

func writeTOMLTable(sb *strings.Builder, data map[string]interface{}, prefix string) {
	// Simple values first
	for _, key := range sortedKeys(data) {
		val := data[key]
		switch v := val.(type) {
		case nil:
			sb.WriteString(fmt.Sprintf("%s = null\n", key))
		case bool:
			sb.WriteString(fmt.Sprintf("%s = %v\n", key, v))
		case float64:
			if v == float64(int(v)) {
				sb.WriteString(fmt.Sprintf("%s = %d\n", key, int(v)))
			} else {
				sb.WriteString(fmt.Sprintf("%s = %g\n", key, v))
			}
		case string:
			sb.WriteString(fmt.Sprintf("%s = %q\n", key, v))
		}
	}

	// Tables
	for _, key := range sortedKeys(data) {
		val := data[key]
		if subObj, ok := val.(map[string]interface{}); ok {
			tableName := key
			if prefix != "" {
				tableName = prefix + "." + key
			}
			sb.WriteString(fmt.Sprintf("\n[%s]\n", tableName))
			writeTOMLTable(sb, subObj, tableName)
		}
	}
}

// ToCSV converts a JSON array of objects to CSV.
func ToCSV(data interface{}) (string, error) {
	arr, ok := data.([]interface{})
	if !ok {
		return "", fmt.Errorf("data must be a JSON array for CSV conversion")
	}

	if len(arr) == 0 {
		return "", nil
	}

	// Collect all keys
	keySet := make(map[string]bool)
	for _, item := range arr {
		if obj, ok := item.(map[string]interface{}); ok {
			for k := range obj {
				keySet[k] = true
			}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	// Header
	writer.Write(keys)

	// Rows
	for _, item := range arr {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		row := make([]string, len(keys))
		for i, k := range keys {
			if val, exists := obj[k]; exists {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		writer.Write(row)
	}

	writer.Flush()
	return sb.String(), nil
}

// ToHTMLTable converts a JSON array of objects to an HTML table.
func ToHTMLTable(data interface{}) (string, error) {
	arr, ok := data.([]interface{})
	if !ok {
		return "", fmt.Errorf("data must be a JSON array for HTML table conversion")
	}

	if len(arr) == 0 {
		return "<table></table>", nil
	}

	// Collect all keys
	keySet := make(map[string]bool)
	for _, item := range arr {
		if obj, ok := item.(map[string]interface{}); ok {
			for k := range obj {
				keySet[k] = true
			}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString("<table>\n<thead>\n<tr>\n")
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("  <th>%s</th>\n", k))
	}
	sb.WriteString("</tr>\n</thead>\n<tbody>\n")

	for _, item := range arr {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		sb.WriteString("<tr>\n")
		for _, k := range keys {
			val := ""
			if v, exists := obj[k]; exists {
				val = fmt.Sprintf("%v", v)
			}
			sb.WriteString(fmt.Sprintf("  <td>%s</td>\n", val))
		}
		sb.WriteString("</tr>\n")
	}

	sb.WriteString("</tbody>\n</table>")
	return sb.String(), nil
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
