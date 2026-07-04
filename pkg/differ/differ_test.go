package differ

import (
	"testing"
)

func TestDiffIdentical(t *testing.T) {
	old := map[string]interface{}{"a": float64(1), "b": "hello"}
	new := map[string]interface{}{"a": float64(1), "b": "hello"}
	result := Diff(old, new)
	if !result.Same {
		t.Error("expected same=true")
	}
	if result.Count != 0 {
		t.Errorf("expected 0 changes, got %d", result.Count)
	}
}

func TestDiffAdded(t *testing.T) {
	old := map[string]interface{}{"a": float64(1)}
	new := map[string]interface{}{"a": float64(1), "b": "new"}
	result := Diff(old, new)
	if result.Same {
		t.Error("expected same=false")
	}
	if result.Count != 1 {
		t.Errorf("expected 1 change, got %d", result.Count)
	}
	if result.Changes[0].Type != Added {
		t.Errorf("expected Added, got %s", result.Changes[0].Type)
	}
}

func TestDiffRemoved(t *testing.T) {
	old := map[string]interface{}{"a": float64(1), "b": "old"}
	new := map[string]interface{}{"a": float64(1)}
	result := Diff(old, new)
	if result.Count != 1 {
		t.Errorf("expected 1 change, got %d", result.Count)
	}
	if result.Changes[0].Type != Removed {
		t.Errorf("expected Removed, got %s", result.Changes[0].Type)
	}
}

func TestDiffModified(t *testing.T) {
	old := map[string]interface{}{"a": float64(1)}
	new := map[string]interface{}{"a": float64(2)}
	result := Diff(old, new)
	if result.Count != 1 {
		t.Errorf("expected 1 change, got %d", result.Count)
	}
	if result.Changes[0].Type != Modified {
		t.Errorf("expected Modified, got %s", result.Changes[0].Type)
	}
}

func TestDiffNested(t *testing.T) {
	old := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  float64(30),
		},
	}
	new := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Jane",
			"age":  float64(30),
		},
	}
	result := Diff(old, new)
	if result.Count != 1 {
		t.Errorf("expected 1 change, got %d", result.Count)
	}
	if result.Changes[0].Path != "user.name" {
		t.Errorf("expected path user.name, got %s", result.Changes[0].Path)
	}
}

func TestDiffArrays(t *testing.T) {
	old := []interface{}{float64(1), float64(2), float64(3)}
	new := []interface{}{float64(1), float64(99), float64(3), float64(4)}
	result := Diff(old, new)
	// Should have: modified [1], added [3]
	if result.Count != 2 {
		t.Errorf("expected 2 changes, got %d", result.Count)
	}
}

func TestSummary(t *testing.T) {
	old := map[string]interface{}{"a": float64(1), "b": "old"}
	new := map[string]interface{}{"a": float64(2), "c": "new"}
	result := Diff(old, new)
	summary := Summary(result)
	if summary["modified"] != 1 {
		t.Errorf("expected 1 modified, got %d", summary["modified"])
	}
	if summary["added"] != 1 {
		t.Errorf("expected 1 added, got %d", summary["added"])
	}
	if summary["removed"] != 1 {
		t.Errorf("expected 1 removed, got %d", summary["removed"])
	}
}

func TestFormatText(t *testing.T) {
	old := map[string]interface{}{"a": float64(1)}
	new := map[string]interface{}{"a": float64(2)}
	result := Diff(old, new)
	text := FormatText(result)
	if text == "" {
		t.Error("expected non-empty text")
	}
}

func TestPaths(t *testing.T) {
	old := map[string]interface{}{"b": float64(1), "a": float64(1)}
	new := map[string]interface{}{"b": float64(2), "a": float64(2)}
	result := Diff(old, new)
	paths := Paths(result)
	if len(paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(paths))
	}
	if paths[0] != "a" || paths[1] != "b" {
		t.Errorf("expected [a, b], got %v", paths)
	}
}
