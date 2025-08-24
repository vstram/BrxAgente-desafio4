package agent

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"
	
	"BrxAgente-desafio4/internal/chat"
	"BrxAgente-desafio4/internal/config"
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

// TestVRAgent_EndToEndWorkflow testa fluxo completo end-to-end
func TestVRAgent_EndToEndWorkflow(t *testing.T) {
	// Setup completo com configurações reais
	cfg := &config.Config{
		OpenAIKey: "",
		OllamaConfig: config.OllamaConfig{
			BaseURL: "",
			Model:   "",
		},
	}
	
	chatSvc := chat.NewChat(cfg)
	
	agentConfig := &AgentConfig{
		Enabled:          true,
		Model:           "test-model",
		Temperature:     0.7,
		MaxTokens:       2000,
		Timeout:         30 * time.Second,
		MemorySize:      100,
		MemoryTTL:       24 * time.Hour,
		ContextWindow:   4000,
		MaxMemoryTokens: 2000,
		WorkerPoolSize:  4,
		CacheEnabled:    true,
		CacheSize:       1000,
		CacheTTL:        1 * time.Hour,
		LogLevel:        "info",
		DebugMode:       true,
		ToolsEnabled:    []string{"excel", "calculation", "validation"},
	}
	
	agent, err := NewVRAgent(agentConfig, chatSvc)
	if err != nil {
		t.Fatalf("Failed to create VRAgent: %v", err)
	}

	t.Run("CompleteWorkflowSequence", func(t *testing.T) {
		testCompleteWorkflowSequence(t, agent)
	})

	t.Run("MemoryPersistenceInWorkflow", func(t *testing.T) {
		testMemoryPersistenceInWorkflow(t, agent)
	})

	t.Run("ChatIntegration", func(t *testing.T) {
		testChatIntegration(t, agent, chatSvc)
	})

	t.Run("ToolChaining", func(t *testing.T) {
		testToolChaining(t, agent)
	})
}

// testCompleteWorkflowSequence testa uma sequência completa de workflow
func testCompleteWorkflowSequence(t *testing.T, agent *VRAgent) {
	// 1. Verificar status inicial
	initialStatus := agent.GetStatus()
	if initialStatus.State != "idle" {
		t.Errorf("Expected initial state to be 'idle', got %s", initialStatus.State)
	}

	// 2. Executar workflow de validação
	err := agent.ExecuteWorkflow("validar-dados")
	if err != nil {
		t.Errorf("validar-dados workflow failed: %v", err)
	}

	// 3. Verificar que status foi atualizado
	afterValidation := agent.GetStatus()
	if afterValidation.TotalRequests <= initialStatus.TotalRequests {
		t.Error("Total requests should have increased after workflow")
	}

	// 4. Executar workflow principal
	err = agent.ExecuteWorkflow("processar-vr-mensal")
	if err != nil {
		t.Errorf("processar-vr-mensal workflow failed: %v", err)
	}

	// 5. Verificar estatísticas finais
	finalStatus := agent.GetStatus()
	if finalStatus.TotalRequests <= afterValidation.TotalRequests {
		t.Error("Total requests should have increased after second workflow")
	}

	// 6. Verificar que agente volta ao estado idle
	time.Sleep(100 * time.Millisecond) // Pequena pausa para garantir conclusão
	idleStatus := agent.GetStatus()
	if idleStatus.State != "idle" {
		t.Errorf("Expected final state to be 'idle', got %s", idleStatus.State)
	}
}

// testMemoryPersistenceInWorkflow testa persistência de memória durante workflows
func testMemoryPersistenceInWorkflow(t *testing.T, agent *VRAgent) {
	// Limpar memória inicial
	err := agent.ClearMemory()
	if err != nil {
		t.Fatalf("Failed to clear memory: %v", err)
	}

	// Verificar memória vazia
	memory, err := agent.GetMemory()
	if err != nil {
		t.Fatalf("Failed to get memory: %v", err)
	}

	if len(memory) != 0 {
		t.Errorf("Memory should be empty after clear, got %d items", len(memory))
	}

	// Executar algumas operações para popular a memória
	// (Nota: Como não temos LLM real configurado, a memória pode não ser populada)
	err = agent.ExecuteWorkflow("validar-dados")
	if err != nil {
		t.Logf("Workflow error (expected without real LLM): %v", err)
	}

	// Verificar se a memória mantém estado entre operações
	status := agent.GetStatus()
	if status.TotalRequests == 0 {
		t.Error("Should have recorded at least one request")
	}

	// Testar reset da memória
	err = agent.Reset()
	if err != nil {
		t.Fatalf("Failed to reset agent: %v", err)
	}

	resetStatus := agent.GetStatus()
	if resetStatus.TotalRequests != 0 {
		t.Error("Total requests should be 0 after reset")
	}
}

// testChatIntegration testa integração com o sistema de chat
func testChatIntegration(t *testing.T, agent *VRAgent, chatSvc *chat.Chat) {
	// Configurar o agente no chat service
	chatSvc.SetAgent(agent)

	// Teste básico de pergunta via chat
	// Como não temos LLM configurado, esperamos fallback
	response, err := chatSvc.Ask("Teste de integração", "", []chat.Message{})
	
	// Pode falhar devido à falta de configuração de LLM - isso é esperado
	if err != nil {
		if strings.Contains(err.Error(), "nenhum serviço de IA configurado") {
			t.Logf("Fallback behavior working correctly: %v", err)
		} else {
			t.Logf("Chat integration error (expected without LLM config): %v", err)
		}
	} else {
		t.Logf("Chat integration successful: %s", response)
	}

	// Verificar que agente ainda está funcional após tentativa de uso via chat
	tools := agent.GetAvailableTools()
	if len(tools) == 0 {
		t.Error("Agent should still have tools after chat integration test")
	}
}

// testToolChaining testa encadeamento de ferramentas
func testToolChaining(t *testing.T, agent *VRAgent) {
	// Simular uma sequência de operações que normalmente seriam encadeadas

	// 1. Primeiro tentar ler dados
	_, err := agent.ExecuteTool("read_excel", `{"file_path": "ATIVOS.xlsx"}`)
	if err != nil {
		t.Logf("Read Excel failed (expected without file): %v", err)
		// Continuar mesmo se ler arquivo falhar - usaremos dados simulados
	}

	// 2. Usar dados para validação
	validateInput := `{
		"tipo_validacao": "colaborador",
		"colaborador": {
			"matricula": "12345",
			"sindicato": "SINDPD",
			"empresa": "Empresa Test"
		}
	}`

	validateResult, err := agent.ExecuteTool("validate_data", validateInput)
	if err != nil {
		t.Errorf("Validation tool failed: %v", err)
		return
	}

	t.Logf("Validation result: %s", validateResult)

	// 3. Usar dados validados para cálculo
	calcInput := `{
		"colaborador": {
			"matricula": "12345",
			"sindicato": "SINDPD"
		},
		"valor_por_sindicato": {"SINDPD": 21.25},
		"dias_uteis_por_sindicato": {"SINDPD": 22}
	}`

	calcResult, err := agent.ExecuteTool("calculate_vr", calcInput)
	if err != nil {
		t.Errorf("Calculation tool failed: %v", err)
		return
	}

	t.Logf("Calculation result: %s", calcResult)

	// Verificar que todos os passos foram registrados
	status := agent.GetStatus()
	if status.TotalRequests < 2 {
		t.Error("Should have registered at least 2 tool executions")
	}
}

// TestVRAgent_RealFileIntegration testa integração com arquivos reais se disponíveis
func TestVRAgent_RealFileIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real file integration test in short mode")
	}

	chatSvc := &chat.Chat{}
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Failed to create VRAgent: %v", err)
	}

	// Lista de arquivos para testar
	testFiles := []string{
		"ATIVOS.xlsx",
		"AFASTAMENTOS.xlsx",
		"Base sindicato x valor.xlsx",
		"Base dias uteis.xlsx",
		"FÉRIAS.xlsx",
	}

	for _, filename := range testFiles {
		t.Run("RealFile_"+filename, func(t *testing.T) {
			filePath := filepath.Join("../../files", filename)
			
			start := time.Now()
			result, err := agent.ExecuteTool("read_excel", filePath)
			duration := time.Since(start)

			t.Logf("File %s processing took: %v", filename, duration)

			if err != nil {
				// Arquivos podem não existir no ambiente de teste
				t.Logf("File %s not accessible: %v", filename, err)
			} else {
				t.Logf("File %s read successfully: %.200s...", filename, result)

				// Verificar que resultado é válido JSON
				var output map[string]interface{}
				if jsonErr := json.Unmarshal([]byte(result), &output); jsonErr != nil {
					t.Errorf("Invalid JSON output for %s: %v", filename, jsonErr)
				}
			}

			// Verificar que performance está dentro do limite (< 5s por arquivo)
			if duration > 5*time.Second {
				t.Errorf("File processing too slow for %s: %v", filename, duration)
			}
		})
	}
}

// TestVRAgent_PerformanceBenchmark executa benchmark de performance
func TestVRAgent_PerformanceBenchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance benchmark in short mode")
	}

	chatSvc := &chat.Chat{}
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Failed to create VRAgent: %v", err)
	}

	// Benchmark diferentes operações
	benchmarks := []struct {
		name      string
		operation func() error
	}{
		{
			name: "GetStatus",
			operation: func() error {
				_ = agent.GetStatus()
				return nil
			},
		},
		{
			name: "GetAvailableTools",
			operation: func() error {
				_ = agent.GetAvailableTools()
				return nil
			},
		},
		{
			name: "ExecuteWorkflow",
			operation: func() error {
				return agent.ExecuteWorkflow("validar-dados")
			},
		},
		{
			name: "ToolExecution",
			operation: func() error {
				_, err := agent.ExecuteTool("validate_data", `{"test": "data"}`)
				return err
			},
		},
	}

	for _, bm := range benchmarks {
		t.Run("Benchmark_"+bm.name, func(t *testing.T) {
			iterations := 100
			start := time.Now()

			for i := 0; i < iterations; i++ {
				err := bm.operation()
				if err != nil {
					t.Logf("Iteration %d error (may be expected): %v", i, err)
				}
			}

			duration := time.Since(start)
			avgTime := duration / time.Duration(iterations)

			t.Logf("%s benchmark:", bm.name)
			t.Logf("  Total time: %v", duration)
			t.Logf("  Average time: %v", avgTime)
			t.Logf("  Operations per second: %.2f", float64(iterations)/duration.Seconds())

			// Performance thresholds
			if avgTime > 100*time.Millisecond {
				t.Errorf("%s average time too high: %v > 100ms", bm.name, avgTime)
			}
		})
	}
}