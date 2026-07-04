package flatten

import (
	"testing"
)

func TestFlattenSimple(t *testing.T) {
	data := map[string]interface{}{
		"a": float64(1),
		"b": "hello",
	}
	result := Flatten(data)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if result["a"] != float64(1) {
		t.Errorf("expected a=1, got %v", result["a"])
	}
}

func TestFlattenNested(t *testing.T) {
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": float64(1),
			"c": map[string]interface{}{
				"d": "deep",
			},
		},
	}
	result := Flatten(data)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if result["a.b"] != float64(1) {
		t.Errorf("expected a.b=1, got %v", result["a.b"])
	}
	if result["a.c.d"] != "deep" {
		t.Errorf("expected a.c.d=deep, got %v", result["a.c.d"])
	}
}

func TestFlattenArray(t *testing.T) {
	data := map[string]interface{}{
		"items": []interface{}{float64(1), float64(2), float64(3)},
	}
	result := Flatten(data)
	if len(result) != 3 {
		t.Errorf("expected 3 keys, got %d", len(result))
	}
	if result["items[0]"] != float64(1) {
		t.Errorf("expected items[0]=1, got %v", result["items[0]"])
	}
}

func TestUnflatten(t *testing.T) {
	data := map[string]interface{}{
		"a.b":   float64(1),
		"a.c.d": "deep",
	}
	result := Unflatten(data)
	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}
	a, ok := obj["a"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected a to be map, got %T", obj["a"])
	}
	if a["b"] != float64(1) {
		t.Errorf("expected a.b=1, got %v", a["b"])
	}
}

func TestUnflattenArray(t *testing.T) {
	data := map[string]interface{}{
		"items[0]": float64(1),
		"items[1]": float64(2),
	}
	result := Unflatten(data)
	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}
	arr, ok := obj["items"].([]interface{})
	if !ok {
		t.Fatalf("expected items to be array, got %T", obj["items"])
	}
	if len(arr) != 2 {
		t.Errorf("expected 2 items, got %d", len(arr))
	}
}

func TestKeys(t *testing.T) {
	data := map[string]interface{}{
		"c": float64(1),
		"a": float64(2),
		"b": float64(3),
	}
	keys := Keys(data)
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "a" || keys[1] != "b" || keys[2] != "c" {
		t.Errorf("expected [a, b, c], got %v", keys)
	}
}
