package agent

import (
	"encoding/json"
	"testing"
	
	"BrxAgente-desafio4/internal/chat"
)

// TestVRAgent_ToolIntegration testa a integração completa das ferramentas com o agente
func TestVRAgent_ToolIntegration(t *testing.T) {
	// Criar agente
	chatSvc := &chat.Chat{} // Mock simples
	agent, err := NewVRAgent(nil, chatSvc)
	if err != nil {
		t.Fatalf("Failed to create VRAgent: %v", err)
	}

	// Verificar se as ferramentas foram carregadas
	tools := agent.GetAvailableTools()
	expectedTools := []string{"read_excel", "calculate_vr", "validate_data"}
	
	if len(tools) != len(expectedTools) {
		t.Errorf("Expected %d tools, got %d", len(expectedTools), len(tools))
	}
	
	for _, expected := range expectedTools {
		found := false
		for _, tool := range tools {
			if tool == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Tool '%s' not found in available tools", expected)
		}
	}
}

// TestVRAgent_ExecuteReadExcelTool testa execução da ferramenta ReadExcel
func TestVRAgent_ExecuteReadExcelTool(t *testing.T) {
	chatSvc := &chat.Chat{}
	agent, err := NewVRAgent(nil, chatSvc)
	if err != nil {
		t.Fatalf("Failed to create VRAgent: %v", err)
	}

	// Testar execução de ferramenta conhecida (ATIVOS.xlsx) - usar simulação
	input := `{"file_path": "ATIVOS.xlsx"}`
	result, err := agent.ExecuteTool("read_excel", input)
	if err != nil {
		t.Fatalf("ExecuteTool() error = %v", err)
	}

	// Parse result
	var output map[string]interface{}
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}
	
	// Para arquivo ATIVOS.xlsx, mesmo que não exista, deve usar simulação
	// Verificar se retornou dados simulados (ou erro real)
	if success, ok := output["success"].(bool); ok && success {
		// Sucesso - dados simulados
		if rowCount, ok := output["row_count"].(float64); !ok || rowCount <= 0 {
			t.Errorf("Expected row_count to be greater than 0, got: %v", output["row_count"])
		}
	} else {
		// Se falhou, pelo menos deve ter indicado o erro claramente
		if errorMsg, exists := output["error"]; !exists || errorMsg == "" {
			t.Error("Expected error message when file not found")
		}
		t.Logf("ReadExcel failed as expected (file not found): %v", output["error"])
	}
}

// TestVRAgent_ExecuteCalculateVRTool testa execução da ferramenta CalculateVR
func TestVRAgent_ExecuteCalculateVRTool(t *testing.T) {
	chatSvc := &chat.Chat{}
	agent, err := NewVRAgent(nil, chatSvc)
	if err != nil {
		t.Fatalf("Failed to create VRAgent: %v", err)
	}

	// Input válido para cálculo
	// ⚠️ CONFIDENCIALIDADE: Remove campo nome conforme PRD.md
	input := `{
		"colaborador": {
			"matricula": "12345",
			"sindicato": "SINDPD"
		},
		"valor_por_sindicato": {"SINDPD": 21.25},
		"dias_uteis_por_sindicato": {"SINDPD": 22},
		"mes_referencia": "2025-09"
	}`

	result, err := agent.ExecuteTool("calculate_vr", input)
	if err != nil {
		t.Fatalf("ExecuteTool() error = %v", err)
	}

	// Parse result
	var output map[string]interface{}
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verificar campos obrigatórios
	if success, ok := output["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true, got: %v", output["error"])
	}

	if valorTotal, ok := output["valor_total"].(float64); !ok || valorTotal <= 0 {
		t.Error("Expected valor_total to be greater than 0")
	}
}

// TestVRAgent_ExecuteValidateDataTool testa execução da ferramenta ValidateData
func TestVRAgent_ExecuteValidateDataTool(t *testing.T) {
	chatSvc := &chat.Chat{}
	agent, err := NewVRAgent(nil, chatSvc)
	if err != nil {
		t.Fatalf("Failed to create VRAgent: %v", err)
	}

	// Input para validação de colaborador
	// ⚠️ CONFIDENCIALIDADE: Remove campo nome conforme PRD.md
	input := `{
		"tipo_validacao": "colaborador",
		"colaborador": {
			"matricula": "12345",
			"sindicato": "SINDPD",
			"empresa": "Empresa ABC"
		},
		"validar_campos_obrigatorios": true
	}`

	result, err := agent.ExecuteTool("validate_data", input)
	if err != nil {
		t.Fatalf("ExecuteTool() error = %v", err)
	}

	// Parse result
	var output map[string]interface{}
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verificar campos obrigatórios
	if success, ok := output["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true, got: %v", output["error"])
	}

	if totalRegistros, ok := output["total_registros"].(float64); !ok || totalRegistros != 1 {
		t.Error("Expected total_registros to be 1")
	}
}

// TestVRAgent_GetToolInfo testa obtenção de informações das ferramentas
func TestVRAgent_GetToolInfo(t *testing.T) {
	chatSvc := &chat.Chat{}
	agent, err := NewVRAgent(nil, chatSvc)
	if err != nil {
		t.Fatalf("Failed to create VRAgent: %v", err)
	}

	// Testar informações de uma ferramenta específica
	info, err := agent.GetToolInfo("read_excel")
	if err != nil {
		t.Fatalf("GetToolInfo() error = %v", err)
	}

	// Verificar campos obrigatórios
	if name, ok := info["name"].(string); !ok || name != "read_excel" {
		t.Error("Expected name to be 'read_excel'")
	}

	if description, ok := info["description"].(string); !ok || description == "" {
		t.Error("Expected description to be non-empty")
	}

	if schema, ok := info["schema"].(map[string]interface{}); !ok || schema == nil {
		t.Error("Expected schema to be present")
	}
}

// TestVRAgent_GetAllToolsInfo testa obtenção de informações de todas as ferramentas
func TestVRAgent_GetAllToolsInfo(t *testing.T) {
	chatSvc := &chat.Chat{}
	agent, err := NewVRAgent(nil, chatSvc)
	if err != nil {
		t.Fatalf("Failed to create VRAgent: %v", err)
	}

	// Obter informações de todas as ferramentas
	allInfo := agent.GetAllToolsInfo()

	expectedTools := []string{"read_excel", "calculate_vr", "validate_data"}
	for _, toolName := range expectedTools {
		if info, ok := allInfo[toolName]; !ok || info == nil {
			t.Errorf("Expected tool info for '%s' to be present", toolName)
		}
	}
}

// TestVRAgent_DisabledToolExecution testa que ferramentas não podem ser executadas quando agente está desabilitado
func TestVRAgent_DisabledToolExecution(t *testing.T) {
	chatSvc := &chat.Chat{}
	agent, err := NewVRAgent(nil, chatSvc)
	if err != nil {
		t.Fatalf("Failed to create VRAgent: %v", err)
	}

	// Desabilitar agente
	agent.Disable()

	// Tentar executar ferramenta
	input := `{"file_path": "test.xlsx"}`
	_, err = agent.ExecuteTool("read_excel", input)
	if err == nil {
		t.Error("Expected error when executing tool with disabled agent")
	}

	if err.Error() != "agente está desabilitado" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}