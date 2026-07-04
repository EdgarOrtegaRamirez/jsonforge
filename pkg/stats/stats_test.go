package stats

import (
	"testing"
)

func TestAnalyzeObject(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"tags": []interface{}{"admin", "user"},
	}
	result := Analyze(data, false)
	if result.Type != "object" {
		t.Errorf("expected object, got %s", result.Type)
	}
	if result.KeyCount != 3 {
		t.Errorf("expected 3 keys, got %d", result.KeyCount)
	}
	if result.Depth != 2 {
		t.Errorf("expected depth 2, got %d", result.Depth)
	}
}

func TestAnalyzeArray(t *testing.T) {
	data := []interface{}{float64(1), float64(2), float64(3)}
	result := Analyze(data, false)
	if result.Type != "array" {
		t.Errorf("expected array, got %s", result.Type)
	}
	if result.ArrayLen != 3 {
		t.Errorf("expected array length 3, got %d", result.ArrayLen)
	}
}

func TestAnalyzeWithPaths(t *testing.T) {
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": "value",
		},
	}
	result := Analyze(data, true)
	if len(result.Paths) != 1 {
		t.Errorf("expected 1 path, got %d", len(result.Paths))
	}
	if result.Paths[0] != "a.b" {
		t.Errorf("expected path a.b, got %s", result.Paths[0])
	}
}

func TestAnalyzeNested(t *testing.T) {
	data := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"level3": "deep",
			},
		},
	}
	result := Analyze(data, false)
	if result.Depth != 3 {
		t.Errorf("expected depth 3, got %d", result.Depth)
	}
}

func TestTypeCounts(t *testing.T) {
	data := map[string]interface{}{
		"str":   "hello",
		"num":   float64(42),
		"bool":  true,
		"null":  nil,
		"array": []interface{}{float64(1)},
	}
	result := Analyze(data, false)
	if result.TypeCounts["string"] != 1 {
		t.Errorf("expected 1 string, got %d", result.TypeCounts["string"])
	}
	if result.TypeCounts["number"] != 2 {
		t.Errorf("expected 2 numbers, got %d", result.TypeCounts["number"])
	}
	if result.TypeCounts["boolean"] != 1 {
		t.Errorf("expected 1 boolean, got %d", result.TypeCounts["boolean"])
	}
	if result.TypeCounts["null"] != 1 {
		t.Errorf("expected 1 null, got %d", result.TypeCounts["null"])
	}
}

func TestFormatText(t *testing.T) {
	data := map[string]interface{}{"a": float64(1)}
	result := Analyze(data, false)
	text := FormatText(result)
	if text == "" {
		t.Error("expected non-empty text")
	}
}

func TestFormatCompact(t *testing.T) {
	data := map[string]interface{}{"a": float64(1)}
	result := Analyze(data, false)
	compact := FormatCompact(result)
	if compact == "" {
		t.Error("expected non-empty compact")
	}
}
