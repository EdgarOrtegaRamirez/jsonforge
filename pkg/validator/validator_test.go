package validator

import (
	"testing"
)

func TestValidateValid(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
		"required": []interface{}{"name", "age"},
		"properties": map[string]interface{}{
			"name": map[string]interface{}{"type": "string"},
			"age":  map[string]interface{}{"type": "number"},
		},
	}
	data := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
	}
	result := Validate(data, schema)
	if !result.Valid {
		t.Errorf("expected valid, got %d errors", len(result.Errors))
	}
}

func TestValidateMissingRequired(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
		"required": []interface{}{"name", "age"},
	}
	data := map[string]interface{}{
		"name": "John",
	}
	result := Validate(data, schema)
	if result.Valid {
		t.Error("expected invalid")
	}
	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(result.Errors))
	}
}

func TestValidateWrongType(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{"type": "string"},
		},
	}
	data := map[string]interface{}{
		"name": float64(123),
	}
	result := Validate(data, schema)
	if result.Valid {
		t.Error("expected invalid")
	}
}

func TestValidateMinimum(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"age": map[string]interface{}{
				"type":    "number",
				"minimum": float64(0),
				"maximum": float64(150),
			},
		},
	}
	data := map[string]interface{}{
		"age": float64(-5),
	}
	result := Validate(data, schema)
	if result.Valid {
		t.Error("expected invalid for negative age")
	}
}

func TestValidateEnum(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"status": map[string]interface{}{
				"type": "string",
				"enum": []interface{}{"active", "inactive"},
			},
		},
	}
	data := map[string]interface{}{
		"status": "pending",
	}
	result := Validate(data, schema)
	if result.Valid {
		t.Error("expected invalid for invalid enum value")
	}
}

func TestValidateMinLength(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":      "string",
				"minLength": float64(3),
			},
		},
	}
	data := map[string]interface{}{
		"name": "Jo",
	}
	result := Validate(data, schema)
	if result.Valid {
		t.Error("expected invalid for short name")
	}
}

func TestValidateArrayItems(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"tags": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}
	data := map[string]interface{}{
		"tags": []interface{}{"valid", float64(123)},
	}
	result := Validate(data, schema)
	if result.Valid {
		t.Error("expected invalid for non-string array item")
	}
}

func TestFormatValidation(t *testing.T) {
	result := &ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{Path: "name", Message: "required field missing"},
		},
	}
	text := FormatValidation(result)
	if text == "" {
		t.Error("expected non-empty text")
	}
}
