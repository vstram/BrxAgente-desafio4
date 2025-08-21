package tools

import (
	"fmt"
	"sync"
)

// ToolRegistry gerencia o registro e descoberta de ferramentas
type ToolRegistry struct {
	tools map[string]VRTool
	mutex sync.RWMutex
}

// NewToolRegistry cria um novo registro de ferramentas
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]VRTool),
	}
}

// Register registra uma nova ferramenta
func (tr *ToolRegistry) Register(tool VRTool) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	
	name := tool.Name()
	if name == "" {
		return fmt.Errorf("nome da ferramenta não pode ser vazio")
	}
	
	if _, exists := tr.tools[name]; exists {
		return fmt.Errorf("ferramenta '%s' já está registrada", name)
	}
	
	tr.tools[name] = tool
	return nil
}

// Get retorna uma ferramenta pelo nome
func (tr *ToolRegistry) Get(name string) (VRTool, error) {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()
	
	tool, exists := tr.tools[name]
	if !exists {
		return nil, fmt.Errorf("ferramenta '%s' não encontrada", name)
	}
	
	return tool, nil
}

// List retorna a lista de todas as ferramentas registradas
func (tr *ToolRegistry) List() []VRTool {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()
	
	tools := make([]VRTool, 0, len(tr.tools))
	for _, tool := range tr.tools {
		tools = append(tools, tool)
	}
	
	return tools
}

// ListNames retorna a lista de nomes das ferramentas registradas
func (tr *ToolRegistry) ListNames() []string {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()
	
	names := make([]string, 0, len(tr.tools))
	for name := range tr.tools {
		names = append(names, name)
	}
	
	return names
}

// Exists verifica se uma ferramenta está registrada
func (tr *ToolRegistry) Exists(name string) bool {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()
	
	_, exists := tr.tools[name]
	return exists
}

// Unregister remove uma ferramenta do registro
func (tr *ToolRegistry) Unregister(name string) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	
	if _, exists := tr.tools[name]; !exists {
		return fmt.Errorf("ferramenta '%s' não está registrada", name)
	}
	
	delete(tr.tools, name)
	return nil
}

// Count retorna o número de ferramentas registradas
func (tr *ToolRegistry) Count() int {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()
	
	return len(tr.tools)
}

// Execute executa uma ferramenta pelo nome
func (tr *ToolRegistry) Execute(toolName, input string) (string, error) {
	tool, err := tr.Get(toolName)
	if err != nil {
		return "", err
	}
	
	// Validar input antes da execução
	if err := tool.Validate(input); err != nil {
		return "", fmt.Errorf("input inválido para ferramenta '%s': %w", toolName, err)
	}
	
	// Executar ferramenta
	return tool.Execute(input)
}

// GetToolInfo retorna informações sobre uma ferramenta
func (tr *ToolRegistry) GetToolInfo(name string) (map[string]interface{}, error) {
	tool, err := tr.Get(name)
	if err != nil {
		return nil, err
	}
	
	info := map[string]interface{}{
		"name":        tool.Name(),
		"description": tool.Description(),
		"schema":      tool.Schema(),
	}
	
	return info, nil
}

// GetAllToolsInfo retorna informações sobre todas as ferramentas
func (tr *ToolRegistry) GetAllToolsInfo() map[string]interface{} {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()
	
	info := make(map[string]interface{})
	for name, tool := range tr.tools {
		info[name] = map[string]interface{}{
			"name":        tool.Name(),
			"description": tool.Description(),
			"schema":      tool.Schema(),
		}
	}
	
	return info
}