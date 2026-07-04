package converter

import (
	"strings"
	"testing"
)

func TestToYAML(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"address": map[string]interface{}{
			"city": "NYC",
		},
	}
	result := ToYAML(data)
	if !strings.Contains(result, "name: John") {
		t.Error("YAML should contain name: John")
	}
	if !strings.Contains(result, "age: 30") {
		t.Error("YAML should contain age: 30")
	}
}

func TestToYAMLArray(t *testing.T) {
	data := []interface{}{float64(1), float64(2), float64(3)}
	result := ToYAML(data)
	if !strings.Contains(result, "- 1") {
		t.Error("YAML array should contain - 1")
	}
}

func TestToYAMLEmpty(t *testing.T) {
	data := map[string]interface{}{}
	result := ToYAML(data)
	if result != "{}\n" {
		t.Errorf("expected empty object YAML, got %q", result)
	}
}

func TestToTOML(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
	}
	result := ToTOML(data)
	if !strings.Contains(result, "name = \"John\"") {
		t.Error("TOML should contain name = \"John\"")
	}
}

func TestToCSV(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"name": "John", "age": float64(30)},
		map[string]interface{}{"name": "Jane", "age": float64(25)},
	}
	result, err := ToCSV(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "name") || !strings.Contains(result, "age") {
		t.Error("CSV should contain name and age headers")
	}
	if !strings.Contains(result, "John") {
		t.Error("CSV should contain John")
	}
}

func TestToCSVNotArray(t *testing.T) {
	data := map[string]interface{}{"a": float64(1)}
	_, err := ToCSV(data)
	if err == nil {
		t.Error("expected error for non-array input")
	}
}

func TestToHTMLTable(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"name": "John", "age": float64(30)},
	}
	result, err := ToHTMLTable(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "<table>") {
		t.Error("HTML should contain <table>")
	}
	if !strings.Contains(result, "John") {
		t.Error("HTML should contain John")
	}
}

func TestToHTMLTableEmpty(t *testing.T) {
	data := []interface{}{}
	result, err := ToHTMLTable(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "<table></table>") {
		t.Error("HTML should contain empty table")
	}
}
