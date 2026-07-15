// Package output provides multiple output formats for JSON data.
package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// Format represents an output format.
type Format int

const (
	// JSON is the standard JSON format with pretty printing.
	JSON Format = iota
	// JSONL is JSON Lines format (one JSON object per line).
	JSONL
	// TEXT is a human-readable text format.
	TEXT
)

// ParseFormat parses a format string into a Format enum.
func ParseFormat(s string) Format {
	switch s {
	case "jsonl":
		return JSONL
	case "text":
		return TEXT
	case "json", "":
		return JSON
	default:
		return JSON
	}
}

// Writer is the interface for writing JSON data to an output.
type Writer interface {
	Write(data interface{}, w io.Writer) error
}

// JSONWriter writes data in standard JSON format.
type JSONWriter struct {
	Compact bool
}

// NewJSONWriter creates a new JSON writer.
func NewJSONWriter(compact bool) *JSONWriter {
	return &JSONWriter{Compact: compact}
}

// Write writes the data as JSON.
func (w *JSONWriter) Write(data interface{}, out io.Writer) error {
	var err error
	var output []byte
	if w.Compact {
		output, err = json.Marshal(data)
	} else {
		output, err = json.MarshalIndent(data, "", "  ")
	}
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if _, err := out.Write(output); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if !w.Compact {
		if _, err := out.Write([]byte("\n")); err != nil {
			return fmt.Errorf("newline: %w", err)
		}
	}
	return nil
}

// JSONLWriter writes data in JSON Lines format.
type JSONLWriter struct{}

// NewJSONLWriter creates a new JSONL writer.
func NewJSONLWriter() *JSONLWriter {
	return &JSONLWriter{}
}

// Write writes each array element as a separate JSON line.
func (w *JSONLWriter) Write(data interface{}, out io.Writer) error {
	arr, ok := data.([]interface{})
	if !ok {
		// Non-array data: write as a single line
		line, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("marshal: %w", err)
		}
		if _, err := out.Write(line); err != nil {
			return fmt.Errorf("write: %w", err)
		}
		if _, err := out.Write([]byte("\n")); err != nil {
			return fmt.Errorf("newline: %w", err)
		}
		return nil
	}

	for _, item := range arr {
		line, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("marshal: %w", err)
		}
		if _, err := out.Write(line); err != nil {
			return fmt.Errorf("write: %w", err)
		}
		if _, err := out.Write([]byte("\n")); err != nil {
			return fmt.Errorf("newline: %w", err)
		}
	}
	return nil
}

// TextWriter writes data in a human-readable text format.
type TextWriter struct{}

// NewTextWriter creates a new text writer.
func NewTextWriter() *TextWriter {
	return &TextWriter{}
}

// Write writes the data in text format.
func (w *TextWriter) Write(data interface{}, out io.Writer) error {
	arr, ok := data.([]interface{})
	if !ok {
		if _, err := fmt.Fprintf(out, "%v\n", data); err != nil {
			return err
		}
		return nil
	}

	for i, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			for k, v := range m {
				if _, err := fmt.Fprintf(out, "%s[%d].%s: %v\n", "result", i, k, v); err != nil {
					return err
				}
			}
		} else {
			if _, err := fmt.Fprintf(out, "%s[%d]: %v\n", "result", i, item); err != nil {
				return err
			}
		}
	}
	return nil
}

// InfoWriter writes a summary of JSON structure.
type InfoWriter struct{}

// NewInfoWriter creates a new info writer.
func NewInfoWriter() *InfoWriter {
	return &InfoWriter{}
}

// Write writes structure info.
func (w *InfoWriter) Write(data interface{}, out io.Writer) error {
	info := analyzeStructure(data, 0)
	fmt.Fprintln(out, "Structure Summary:")
	fmt.Fprintf(out, "  Type: %s\n", info.Type)
	fmt.Fprintf(out, "  Depth: %d\n", info.Depth)
	if info.Keys > 0 {
		fmt.Fprintf(out, "  Keys: %d\n", info.Keys)
	}
	if info.ArrayLength > 0 {
		fmt.Fprintf(out, "  Array Length: %d\n", info.ArrayLength)
	}
	return nil
}

type structureInfo struct {
	Type         string
	Depth        int
	Keys         int
	ArrayLength  int
}

func analyzeStructure(data interface{}, depth int) structureInfo {
	info := structureInfo{Depth: depth}

	switch v := data.(type) {
	case map[string]interface{}:
		info.Type = "object"
		info.Keys = len(v)
		maxDepth := depth
		for _, val := range v {
			sub := analyzeStructure(val, depth+1)
			if sub.Depth > maxDepth {
				maxDepth = sub.Depth
			}
		}
		info.Depth = maxDepth
	case []interface{}:
		info.Type = "array"
		info.ArrayLength = len(v)
		maxDepth := depth
		for _, val := range v {
			sub := analyzeStructure(val, depth+1)
			if sub.Depth > maxDepth {
				maxDepth = sub.Depth
			}
		}
		info.Depth = maxDepth
	case string:
		info.Type = "string"
	case float64:
		info.Type = "number"
	case bool:
		info.Type = "boolean"
	case nil:
		info.Type = "null"
	default:
		info.Type = fmt.Sprintf("%T", v)
	}

	return info
}
