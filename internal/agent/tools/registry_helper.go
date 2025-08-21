package tools

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