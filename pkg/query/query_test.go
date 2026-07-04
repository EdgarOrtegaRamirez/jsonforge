package query

import (
	"encoding/json"
	"testing"
)

func TestParsePath(t *testing.T) {
	tests := []struct {
		path     string
		expected []PathSegment
		wantErr  bool
	}{
		{"users.0.name", []PathSegment{{Key: "users"}, {Index: 0}, {Key: "name"}}, false},
		{"data.*.id", []PathSegment{{Key: "data"}, {Wildcard: true}, {Key: "id"}}, false},
		{"data..name", []PathSegment{{Key: "data"}, {Recursive: true}, {Key: "name"}}, false},
		{"items[-1]", []PathSegment{{Key: "items"}, {Index: -1}}, false},
		{"items[last]", []PathSegment{{Key: "items"}, {Index: -1}}, false},
		{"", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			segments, err := ParsePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(segments) != len(tt.expected) {
				t.Errorf("ParsePath(%q) got %d segments, want %d", tt.path, len(segments), len(tt.expected))
			}
		})
	}
}

func TestGet(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"address": map[string]interface{}{
			"city":   "NYC",
			"zip":    "10001",
			"coords": []interface{}{40.7128, -74.0060},
		},
		"tags":   []interface{}{"admin", "user", "moderator"},
		"nested": map[string]interface{}{"deep": map[string]interface{}{"value": "found"}},
	}

	tests := []struct {
		path    string
		want    interface{}
		wantErr bool
	}{
		{"name", "John", false},
		{"age", float64(30), false},
		{"address.city", "NYC", false},
		{"address.coords.0", 40.7128, false},
		{"tags.1", "user", false},
		{"tags.-1", "moderator", false},
		{"tags.last", "moderator", false},
		{"nested.deep.value", "found", false},
		{"missing", nil, true},
		{"address.state", nil, true},
		{"tags.10", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, err := Get(data, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !jsonEqual(t, got, tt.want) {
				t.Errorf("Get(%q) = %v (%T), want %v (%T)", tt.path, got, got, tt.want, tt.want)
			}
		})
	}
}

func TestGetWildcard(t *testing.T) {
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice"},
			map[string]interface{}{"name": "Bob"},
		},
	}

	got, err := Get(data, "users.*.name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names, ok := got.([]interface{})
	if !ok {
		t.Fatalf("expected array, got %T", got)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestGetRecursive(t *testing.T) {
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"name": "first",
			"b": map[string]interface{}{
				"name": "second",
				"c": map[string]interface{}{
					"name": "third",
				},
			},
		},
	}

	got, err := Get(data, "a..name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names, ok := got.([]interface{})
	if !ok {
		t.Fatalf("expected array, got %T", got)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 names, got %d", len(names))
	}
}

func TestSet(t *testing.T) {
	data := map[string]interface{}{
		"name":    "John",
		"address": map[string]interface{}{"city": "NYC"},
	}

	result, err := Set(data, "name", "Jane")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := result.(map[string]interface{})
	if m["name"] != "Jane" {
		t.Errorf("expected name=Jane, got %v", m["name"])
	}

	result, err = Set(data, "address.zip", "10001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m = result.(map[string]interface{})
	addr := m["address"].(map[string]interface{})
	if addr["zip"] != "10001" {
		t.Errorf("expected zip=10001, got %v", addr["zip"])
	}
}

func TestSetCreateIntermediate(t *testing.T) {
	data := map[string]interface{}{}

	result, err := Set(data, "a.b.c", "deep")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := result.(map[string]interface{})
	a := m["a"].(map[string]interface{})
	b := a["b"].(map[string]interface{})
	if b["c"] != "deep" {
		t.Errorf("expected c=deep, got %v", b["c"])
	}
}

func TestDelete(t *testing.T) {
	data := map[string]interface{}{
		"name":  "John",
		"age":   float64(30),
		"email": "john@example.com",
	}

	result, err := Delete(data, "email")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := result.(map[string]interface{})
	if _, ok := m["email"]; ok {
		t.Error("expected email to be deleted")
	}
	if m["name"] != "John" {
		t.Error("name should not be affected")
	}
}

func TestExists(t *testing.T) {
	data := map[string]interface{}{
		"name":   "John",
		"nested": map[string]interface{}{"value": "found"},
	}

	if !Exists(data, "name") {
		t.Error("expected name to exist")
	}
	if !Exists(data, "nested.value") {
		t.Error("expected nested.value to exist")
	}
	if Exists(data, "missing") {
		t.Error("expected missing to not exist")
	}
}

func jsonEqual(t *testing.T, a, b interface{}) bool {
	t.Helper()
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	return string(aj) == string(bj)
}
