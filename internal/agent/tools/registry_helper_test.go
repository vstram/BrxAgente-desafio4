package tools

import (
	"testing"
)

func TestRegisterDefaultTools(t *testing.T) {
	registry := NewToolRegistry()

	// Verificar registry vazio inicialmente
	if registry.Count() != 0 {
		t.Error("Registry should be empty initially")
	}

	// Registrar ferramentas padrão
	err := RegisterDefaultTools(registry)
	if err != nil {
		t.Fatalf("RegisterDefaultTools should succeed, got error: %v", err)
	}

	// Verificar se todas as ferramentas foram registradas
	expectedTools := []string{"read_excel", "calculate_vr", "validate_data", "policy_consultant"}
	expectedCount := len(expectedTools)

	if registry.Count() != expectedCount {
		t.Errorf("Expected %d tools, got %d", expectedCount, registry.Count())
	}

	// Verificar se cada ferramenta específica foi registrada
	for _, toolName := range expectedTools {
		if !registry.Exists(toolName) {
			t.Errorf("Tool '%s' should be registered", toolName)
		}

		// Verificar se pode obter a ferramenta
		tool, err := registry.Get(toolName)
		if err != nil {
			t.Errorf("Should be able to get tool '%s': %v", toolName, err)
		}

		// Verificar se o nome está correto
		if tool.Name() != toolName {
			t.Errorf("Tool name mismatch: expected '%s', got '%s'", toolName, tool.Name())
		}
	}
}

func TestGetDefaultToolRegistry(t *testing.T) {
	registry, err := GetDefaultToolRegistry()
	if err != nil {
		t.Fatalf("GetDefaultToolRegistry should succeed, got error: %v", err)
	}

	if registry == nil {
		t.Fatal("Registry should not be nil")
	}

	// Verificar se tem as ferramentas esperadas
	expectedTools := []string{"read_excel", "calculate_vr", "validate_data", "policy_consultant"}
	expectedCount := len(expectedTools)

	if registry.Count() != expectedCount {
		t.Errorf("Expected %d tools, got %d", expectedCount, registry.Count())
	}

	// Verificar funcionalidade básica de cada ferramenta
	for _, toolName := range expectedTools {
		tool, err := registry.Get(toolName)
		if err != nil {
			t.Errorf("Should be able to get tool '%s': %v", toolName, err)
			continue
		}

		// Verificar se tem descrição
		description := tool.Description()
		if description == "" {
			t.Errorf("Tool '%s' should have description", toolName)
		}

		// Verificar se tem schema
		schema := tool.Schema()
		if schema == nil {
			t.Errorf("Tool '%s' should have schema", toolName)
		}
	}
}

func TestGetDefaultToolRegistry_ToolExecution(t *testing.T) {
	registry, err := GetDefaultToolRegistry()
	if err != nil {
		t.Fatalf("GetDefaultToolRegistry should succeed, got error: %v", err)
	}

	// Teste básico de execução para cada ferramenta
	tests := []struct {
		toolName    string
		input       string
		expectError bool
	}{
		{
			toolName:    "read_excel",
			input:       `{"file_path": "nonexistent.xlsx"}`,
			expectError: false, // Should return controlled error, not execution error
		},
		{
			toolName:    "calculate_vr",
			input:       `{"invalid": "input"}`,
			expectError: true, // Will fail validation and return execution error
		},
		{
			toolName:    "validate_data",
			input:       `{"tipo_validacao": "invalid"}`,
			expectError: true, // Will fail validation and return execution error
		},
		{
			toolName:    "policy_consultant",
			input:       `{"query": "Diretores têm direito a VR?", "type": "simple"}`,
			expectError: true, // May fail if knowledge base files are not available in test environment
		},
	}

	for _, tt := range tests {
		t.Run(tt.toolName, func(t *testing.T) {
			result, err := registry.Execute(tt.toolName, tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("Execute('%s') error = %v, expectError %v", tt.toolName, err, tt.expectError)
			}

			if !tt.expectError && result == "" {
				t.Errorf("Execute('%s') should return result", tt.toolName)
			}
		})
	}
}

func TestGetDefaultToolRegistry_ToolInfo(t *testing.T) {
	registry, err := GetDefaultToolRegistry()
	if err != nil {
		t.Fatalf("GetDefaultToolRegistry should succeed, got error: %v", err)
	}

	// Obter informações de todas as ferramentas
	allInfo := registry.GetAllToolsInfo()

	expectedTools := []string{"read_excel", "calculate_vr", "validate_data", "policy_consultant"}
	for _, toolName := range expectedTools {
		info, exists := allInfo[toolName]
		if !exists {
			t.Errorf("Tool info for '%s' should exist", toolName)
			continue
		}

		infoMap, ok := info.(map[string]interface{})
		if !ok {
			t.Errorf("Tool info for '%s' should be a map", toolName)
			continue
		}

		// Verificar se tem campos obrigatórios
		requiredFields := []string{"name", "description", "schema"}
		for _, field := range requiredFields {
			if _, exists := infoMap[field]; !exists {
				t.Errorf("Tool info for '%s' should have field '%s'", toolName, field)
			}
		}
	}
}

func TestRegisterDefaultTools_Twice(t *testing.T) {
	registry := NewToolRegistry()

	// Registrar primeira vez
	err := RegisterDefaultTools(registry)
	if err != nil {
		t.Fatalf("First RegisterDefaultTools should succeed, got error: %v", err)
	}

	// Tentar registrar segunda vez (deve falhar por duplicata)
	err = RegisterDefaultTools(registry)
	if err == nil {
		t.Error("Second RegisterDefaultTools should fail due to duplicates")
	}

	// Registry deve manter apenas um conjunto de ferramentas
	expectedCount := 4 // read_excel, calculate_vr, validate_data, policy_consultant
	if registry.Count() != expectedCount {
		t.Errorf("Expected %d tools after duplicate registration attempt, got %d", expectedCount, registry.Count())
	}
}
