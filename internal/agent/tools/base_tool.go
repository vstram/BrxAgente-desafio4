package tools

import (
	"encoding/json"
	"fmt"
)

// VRTool define a interface base para todas as ferramentas do agente
type VRTool interface {
	// Name retorna o nome único da ferramenta
	Name() string
	
	// Description retorna uma descrição clara do que a ferramenta faz
	Description() string
	
	// Execute executa a ferramenta com o input fornecido
	Execute(input string) (string, error)
	
	// Validate valida se o input está correto antes da execução
	Validate(input string) error
	
	// Schema retorna o schema JSON do input esperado (opcional)
	Schema() map[string]interface{}
}

// BaseTool implementa funcionalidades comuns para todas as ferramentas
type BaseTool struct {
	name        string
	description string
	schema      map[string]interface{}
}

// NewBaseTool cria uma nova BaseTool
func NewBaseTool(name, description string, schema map[string]interface{}) *BaseTool {
	return &BaseTool{
		name:        name,
		description: description,
		schema:      schema,
	}
}

// Name retorna o nome da ferramenta
func (bt *BaseTool) Name() string {
	return bt.name
}

// Description retorna a descrição da ferramenta
func (bt *BaseTool) Description() string {
	return bt.description
}

// Schema retorna o schema da ferramenta
func (bt *BaseTool) Schema() map[string]interface{} {
	return bt.schema
}

// ValidateJSON valida se o input é um JSON válido
func (bt *BaseTool) ValidateJSON(input string) error {
	var temp interface{}
	if err := json.Unmarshal([]byte(input), &temp); err != nil {
		return fmt.Errorf("input não é um JSON válido: %w", err)
	}
	return nil
}

// ParseJSONInput faz parse do input JSON para um map
func (bt *BaseTool) ParseJSONInput(input string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do JSON: %w", err)
	}
	return data, nil
}

// FormatJSONOutput formata o output como JSON
func (bt *BaseTool) FormatJSONOutput(data interface{}) (string, error) {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao formatar output JSON: %w", err)
	}
	return string(output), nil
}

// ToolError representa um erro específico de ferramenta
type ToolError struct {
	Tool    string `json:"tool"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

func (e *ToolError) Error() string {
	return fmt.Sprintf("erro na ferramenta '%s': %s", e.Tool, e.Message)
}

// NewToolError cria um novo erro de ferramenta
func NewToolError(tool, message, code string) *ToolError {
	return &ToolError{
		Tool:    tool,
		Message: message,
		Code:    code,
	}
}