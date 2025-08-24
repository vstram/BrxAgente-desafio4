package tools

import (
	"strings"
	"testing"
)

// MockTool implementa VRTool para testes
type MockTool struct {
	name        string
	description string
}

func (m *MockTool) Name() string {
	return m.name
}

func (m *MockTool) Description() string {
	return m.description
}

func (m *MockTool) Execute(input string) (string, error) {
	return `{"result": "mock executed with input: ` + input + `"}`, nil
}

func (m *MockTool) Validate(input string) error {
	if input == "invalid" {
		return NewToolError(m.name, "input inválido", "INVALID_INPUT")
	}
	return nil
}

func (m *MockTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"test": map[string]interface{}{
				"type": "string",
			},
		},
	}
}

func TestNewToolRegistry(t *testing.T) {
	registry := NewToolRegistry()

	if registry == nil {
		t.Fatal("NewToolRegistry should not return nil")
	}

	if registry.Count() != 0 {
		t.Error("New registry should be empty")
	}
}

func TestToolRegistry_Register(t *testing.T) {
	registry := NewToolRegistry()

	// Testar registro bem-sucedido
	tool := &MockTool{name: "test_tool", description: "Test tool"}
	err := registry.Register(tool)
	if err != nil {
		t.Errorf("Register should succeed, got error: %v", err)
	}

	if registry.Count() != 1 {
		t.Errorf("Registry should have 1 tool, got %d", registry.Count())
	}

	// Testar registro duplicado
	duplicateTool := &MockTool{name: "test_tool", description: "Duplicate tool"}
	err = registry.Register(duplicateTool)
	if err == nil {
		t.Error("Register should fail for duplicate tool name")
	}

	// Testar registro com nome vazio
	emptyNameTool := &MockTool{name: "", description: "Empty name tool"}
	err = registry.Register(emptyNameTool)
	if err == nil {
		t.Error("Register should fail for empty tool name")
	}
}

func TestToolRegistry_Get(t *testing.T) {
	registry := NewToolRegistry()
	tool := &MockTool{name: "test_tool", description: "Test tool"}

	// Registrar ferramenta
	err := registry.Register(tool)
	if err != nil {
		t.Fatalf("Failed to register tool: %v", err)
	}

	// Testar Get bem-sucedido
	retrievedTool, err := registry.Get("test_tool")
	if err != nil {
		t.Errorf("Get should succeed, got error: %v", err)
	}

	if retrievedTool.Name() != "test_tool" {
		t.Errorf("Expected tool name 'test_tool', got '%s'", retrievedTool.Name())
	}

	// Testar Get para ferramenta inexistente
	_, err = registry.Get("nonexistent_tool")
	if err == nil {
		t.Error("Get should fail for nonexistent tool")
	}
}

func TestToolRegistry_List(t *testing.T) {
	registry := NewToolRegistry()

	// Registrar algumas ferramentas
	tool1 := &MockTool{name: "tool1", description: "First tool"}
	tool2 := &MockTool{name: "tool2", description: "Second tool"}

	registry.Register(tool1)
	registry.Register(tool2)

	tools := registry.List()
	if len(tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(tools))
	}

	// Verificar se ambas as ferramentas estão na lista
	foundTool1, foundTool2 := false, false
	for _, tool := range tools {
		if tool.Name() == "tool1" {
			foundTool1 = true
		}
		if tool.Name() == "tool2" {
			foundTool2 = true
		}
	}

	if !foundTool1 || !foundTool2 {
		t.Error("List should contain both registered tools")
	}
}

func TestToolRegistry_ListNames(t *testing.T) {
	registry := NewToolRegistry()

	// Registrar algumas ferramentas
	tool1 := &MockTool{name: "tool1", description: "First tool"}
	tool2 := &MockTool{name: "tool2", description: "Second tool"}

	registry.Register(tool1)
	registry.Register(tool2)

	names := registry.ListNames()
	if len(names) != 2 {
		t.Errorf("Expected 2 names, got %d", len(names))
	}

	// Verificar se ambos os nomes estão na lista
	foundTool1, foundTool2 := false, false
	for _, name := range names {
		if name == "tool1" {
			foundTool1 = true
		}
		if name == "tool2" {
			foundTool2 = true
		}
	}

	if !foundTool1 || !foundTool2 {
		t.Error("ListNames should contain both tool names")
	}
}

func TestToolRegistry_Exists(t *testing.T) {
	registry := NewToolRegistry()
	tool := &MockTool{name: "test_tool", description: "Test tool"}

	// Verificar que não existe antes do registro
	if registry.Exists("test_tool") {
		t.Error("Tool should not exist before registration")
	}

	// Registrar e verificar que existe
	registry.Register(tool)
	if !registry.Exists("test_tool") {
		t.Error("Tool should exist after registration")
	}

	// Verificar ferramenta inexistente
	if registry.Exists("nonexistent_tool") {
		t.Error("Nonexistent tool should not exist")
	}
}

func TestToolRegistry_Unregister(t *testing.T) {
	registry := NewToolRegistry()
	tool := &MockTool{name: "test_tool", description: "Test tool"}

	// Registrar ferramenta
	registry.Register(tool)

	// Verificar que existe
	if !registry.Exists("test_tool") {
		t.Error("Tool should exist after registration")
	}

	// Desregistrar
	err := registry.Unregister("test_tool")
	if err != nil {
		t.Errorf("Unregister should succeed, got error: %v", err)
	}

	// Verificar que não existe mais
	if registry.Exists("test_tool") {
		t.Error("Tool should not exist after unregistration")
	}

	// Tentar desregistrar novamente
	err = registry.Unregister("test_tool")
	if err == nil {
		t.Error("Unregister should fail for already unregistered tool")
	}
}

func TestToolRegistry_Execute(t *testing.T) {
	registry := NewToolRegistry()
	tool := &MockTool{name: "test_tool", description: "Test tool"}

	// Registrar ferramenta
	registry.Register(tool)

	// Testar execução bem-sucedida
	result, err := registry.Execute("test_tool", "valid_input")
	if err != nil {
		t.Errorf("Execute should succeed, got error: %v", err)
	}

	if !strings.Contains(result, "valid_input") {
		t.Error("Result should contain input")
	}

	// Testar execução com input inválido
	_, err = registry.Execute("test_tool", "invalid")
	if err == nil {
		t.Error("Execute should fail for invalid input")
	}

	// Testar execução de ferramenta inexistente
	_, err = registry.Execute("nonexistent_tool", "input")
	if err == nil {
		t.Error("Execute should fail for nonexistent tool")
	}
}

func TestToolRegistry_GetToolInfo(t *testing.T) {
	registry := NewToolRegistry()
	tool := &MockTool{name: "test_tool", description: "Test tool description"}

	// Registrar ferramenta
	registry.Register(tool)

	// Obter informações
	info, err := registry.GetToolInfo("test_tool")
	if err != nil {
		t.Errorf("GetToolInfo should succeed, got error: %v", err)
	}

	// Verificar informações
	if info["name"] != "test_tool" {
		t.Error("Info should contain correct name")
	}

	if info["description"] != "Test tool description" {
		t.Error("Info should contain correct description")
	}

	if info["schema"] == nil {
		t.Error("Info should contain schema")
	}

	// Testar para ferramenta inexistente
	_, err = registry.GetToolInfo("nonexistent_tool")
	if err == nil {
		t.Error("GetToolInfo should fail for nonexistent tool")
	}
}

func TestToolRegistry_GetAllToolsInfo(t *testing.T) {
	registry := NewToolRegistry()

	// Registrar algumas ferramentas
	tool1 := &MockTool{name: "tool1", description: "First tool"}
	tool2 := &MockTool{name: "tool2", description: "Second tool"}

	registry.Register(tool1)
	registry.Register(tool2)

	// Obter todas as informações
	allInfo := registry.GetAllToolsInfo()

	if len(allInfo) != 2 {
		t.Errorf("Expected info for 2 tools, got %d", len(allInfo))
	}

	// Verificar se ambas as ferramentas estão presentes
	if _, exists := allInfo["tool1"]; !exists {
		t.Error("AllInfo should contain tool1")
	}

	if _, exists := allInfo["tool2"]; !exists {
		t.Error("AllInfo should contain tool2")
	}
}

func TestToolRegistry_ConcurrentAccess(t *testing.T) {
	registry := NewToolRegistry()

	// Testar acesso concorrente básico
	done := make(chan bool)

	// Goroutine para registrar ferramentas
	go func() {
		for i := 0; i < 10; i++ {
			tool := &MockTool{name: strings.ReplaceAll("tool_PLACEHOLDER", "PLACEHOLDER", string(rune(i+'0'))), description: "Concurrent tool"}
			registry.Register(tool)
		}
		done <- true
	}()

	// Goroutine para listar ferramentas
	go func() {
		for i := 0; i < 10; i++ {
			registry.ListNames()
		}
		done <- true
	}()

	// Esperar conclusão
	<-done
	<-done

	// Verificar estado final
	if registry.Count() != 10 {
		t.Errorf("Expected 10 tools after concurrent operations, got %d", registry.Count())
	}
}
