package tools

import (
	"testing"
)

func TestNewBaseTool(t *testing.T) {
	name := "test_tool"
	description := "Test tool description"
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"test": map[string]interface{}{
				"type": "string",
			},
		},
	}

	tool := NewBaseTool(name, description, schema)

	if tool == nil {
		t.Fatal("NewBaseTool should not return nil")
	}

	if tool.Name() != name {
		t.Errorf("Expected name '%s', got '%s'", name, tool.Name())
	}

	if tool.Description() != description {
		t.Errorf("Expected description '%s', got '%s'", description, tool.Description())
	}

	if tool.Schema() == nil {
		t.Error("Schema should not be nil")
	}
}

func TestBaseTool_ValidateJSON(t *testing.T) {
	tool := NewBaseTool("test", "test", nil)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "JSON válido",
			input:   `{"key": "value"}`,
			wantErr: false,
		},
		{
			name:    "JSON array válido",
			input:   `["item1", "item2"]`,
			wantErr: false,
		},
		{
			name:    "JSON string válido",
			input:   `"simple string"`,
			wantErr: false,
		},
		{
			name:    "JSON number válido",
			input:   `123`,
			wantErr: false,
		},
		{
			name:    "JSON boolean válido",
			input:   `true`,
			wantErr: false,
		},
		{
			name:    "JSON null válido",
			input:   `null`,
			wantErr: false,
		},
		{
			name:    "JSON inválido - chave sem aspas",
			input:   `{key: "value"}`,
			wantErr: true,
		},
		{
			name:    "JSON inválido - vírgula extra",
			input:   `{"key": "value",}`,
			wantErr: true,
		},
		{
			name:    "JSON inválido - aspas não fechadas",
			input:   `{"key": "value}`,
			wantErr: true,
		},
		{
			name:    "String vazia",
			input:   ``,
			wantErr: true,
		},
		{
			name:    "Texto plano",
			input:   `plain text`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tool.ValidateJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBaseTool_ParseJSONInput(t *testing.T) {
	tool := NewBaseTool("test", "test", nil)

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    map[string]interface{}
	}{
		{
			name:    "JSON object válido",
			input:   `{"key": "value", "number": 123}`,
			wantErr: false,
			want: map[string]interface{}{
				"key":    "value",
				"number": float64(123), // JSON numbers are parsed as float64
			},
		},
		{
			name:    "JSON object vazio",
			input:   `{}`,
			wantErr: false,
			want:    map[string]interface{}{},
		},
		{
			name:    "JSON array - deve falhar",
			input:   `["item1", "item2"]`,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "JSON string - deve falhar",
			input:   `"simple string"`,
			wantErr: true,
			want:    nil,
		},
		{
			name:    "JSON inválido",
			input:   `{invalid json}`,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tool.ParseJSONInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSONInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("ParseJSONInput() got %d keys, want %d", len(got), len(tt.want))
					return
				}

				for key, expectedValue := range tt.want {
					if got[key] != expectedValue {
						t.Errorf("ParseJSONInput() got[%s] = %v, want %v", key, got[key], expectedValue)
					}
				}
			}
		})
	}
}

func TestBaseTool_FormatJSONOutput(t *testing.T) {
	tool := NewBaseTool("test", "test", nil)

	tests := []struct {
		name    string
		data    interface{}
		wantErr bool
	}{
		{
			name: "map válido",
			data: map[string]interface{}{
				"key":    "value",
				"number": 123,
			},
			wantErr: false,
		},
		{
			name:    "slice válido",
			data:    []string{"item1", "item2"},
			wantErr: false,
		},
		{
			name:    "string válido",
			data:    "simple string",
			wantErr: false,
		},
		{
			name:    "number válido",
			data:    123,
			wantErr: false,
		},
		{
			name:    "boolean válido",
			data:    true,
			wantErr: false,
		},
		{
			name:    "nil válido",
			data:    nil,
			wantErr: false,
		},
		{
			name:    "função - deve falhar",
			data:    func() {},
			wantErr: true,
		},
		{
			name:    "channel - deve falhar",
			data:    make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tool.FormatJSONOutput(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatJSONOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar se o output é um JSON válido
				if err := tool.ValidateJSON(got); err != nil {
					t.Errorf("FormatJSONOutput() produced invalid JSON: %v", err)
				}

				// Verificar se tem formatação (indentação) - apenas para objetos/arrays complexos
				if (tt.name == "map válido" || tt.name == "slice válido") && len(got) > 10 && !contains(got, "\n") {
					t.Error("FormatJSONOutput() should produce indented JSON for complex objects")
				}
			}
		})
	}
}

func TestToolError(t *testing.T) {
	toolName := "test_tool"
	message := "Test error message"
	code := "TEST_ERROR"

	err := NewToolError(toolName, message, code)

	if err.Tool != toolName {
		t.Errorf("Expected tool '%s', got '%s'", toolName, err.Tool)
	}

	if err.Message != message {
		t.Errorf("Expected message '%s', got '%s'", message, err.Message)
	}

	if err.Code != code {
		t.Errorf("Expected code '%s', got '%s'", code, err.Code)
	}

	// Testar método Error()
	errorString := err.Error()
	if !contains(errorString, toolName) {
		t.Error("Error string should contain tool name")
	}

	if !contains(errorString, message) {
		t.Error("Error string should contain message")
	}
}

func TestToolError_EmptyCode(t *testing.T) {
	err := NewToolError("test_tool", "test message", "")

	if err.Code != "" {
		t.Errorf("Expected empty code, got '%s'", err.Code)
	}

	// Verificar se ainda funciona sem código
	errorString := err.Error()
	if errorString == "" {
		t.Error("Error string should not be empty")
	}
}

// Helper function para verificar se uma string contém outra
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				func() bool {
					for i := 1; i <= len(s)-len(substr); i++ {
						if s[i:i+len(substr)] == substr {
							return true
						}
					}
					return false
				}())))
}
