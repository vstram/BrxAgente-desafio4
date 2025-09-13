package tools

import (
	"context"
)

// VRPolicyConsultantTool é um wrapper que adapta PolicyConsultantTool para a interface VRTool
type VRPolicyConsultantTool struct {
	tool *PolicyConsultantTool
}

// NewVRPolicyConsultantTool cria uma nova instância do wrapper
func NewVRPolicyConsultantTool(dataDir string) *VRPolicyConsultantTool {
	return &VRPolicyConsultantTool{
		tool: NewPolicyConsultantTool(dataDir),
	}
}

// Name implementa VRTool.Name
func (vrpct *VRPolicyConsultantTool) Name() string {
	return vrpct.tool.Name()
}

// Description implementa VRTool.Description
func (vrpct *VRPolicyConsultantTool) Description() string {
	return vrpct.tool.Description()
}

// Execute implementa VRTool.Execute (adapta para remover Context)
func (vrpct *VRPolicyConsultantTool) Execute(input string) (string, error) {
	return vrpct.tool.Execute(context.Background(), input)
}

// Validate implementa VRTool.Validate
func (vrpct *VRPolicyConsultantTool) Validate(input string) error {
	return vrpct.tool.ValidateInput(input)
}

// Schema implementa VRTool.Schema
func (vrpct *VRPolicyConsultantTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type":        "object",
		"description": "Input para consulta de políticas de VR",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "Pergunta sobre políticas de VR",
				"required":    true,
			},
			"type": map[string]interface{}{
				"type":        "string",
				"description": "Tipo de consulta: simple, complex, whatif, compliance, conflict",
				"enum":        []string{"simple", "complex", "whatif", "compliance", "conflict"},
			},
		},
	}
}

// GetInternalTool retorna a instância interna do PolicyConsultantTool para controle direto do cache
func (vrpct *VRPolicyConsultantTool) GetInternalTool() *PolicyConsultantTool {
	return vrpct.tool
}

// RegisterDefaultTools registra todas as ferramentas padrão no registry
func RegisterDefaultTools(registry *ToolRegistry) error {
	// Registrar ReadExcelTool
	excelTool := NewReadExcelTool()
	if err := registry.Register(excelTool); err != nil {
		return err
	}

	// Registrar CalculateVRTool
	calculateTool := NewCalculateVRTool()
	if err := registry.Register(calculateTool); err != nil {
		return err
	}

	// Registrar ValidateDataTool
	validateTool := NewValidateDataTool()
	if err := registry.Register(validateTool); err != nil {
		return err
	}

	// Registrar PolicyConsultantTool
	policyTool := NewVRPolicyConsultantTool("internal/data/policies")
	if err := registry.Register(policyTool); err != nil {
		return err
	}

	return nil
}

// GetDefaultToolRegistry cria um registry com todas as ferramentas padrão registradas
func GetDefaultToolRegistry() (*ToolRegistry, error) {
	registry := NewToolRegistry()

	if err := RegisterDefaultTools(registry); err != nil {
		return nil, err
	}

	return registry, nil
}
