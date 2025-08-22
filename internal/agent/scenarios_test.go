package agent

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"BrxAgente-desafio4/internal/chat"
	"BrxAgente-desafio4/internal/config"
)

// TestAgentRealScenarios testa cenários reais de uso do agente
func TestAgentRealScenarios(t *testing.T) {
	// Setup do ambiente de teste
	cfg := &config.Config{
		OpenAIKey: "",
		OllamaConfig: config.OllamaConfig{
			BaseURL: "",
			Model:   "",
		},
	}
	
	chatSvc := chat.NewChat(cfg)
	
	// Configuração do agente para testes
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
	
	// Criar agente
	agent, err := NewVRAgent(agentConfig, chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Execute subtestes
	t.Run("Cenario1_ConsultaBasica", func(t *testing.T) {
		testConsultaBasica(t, agent)
	})
	
	t.Run("Cenario2_ValidacaoDados", func(t *testing.T) {
		testValidacaoDados(t, agent)
	})
	
	t.Run("Cenario3_CalculoVR", func(t *testing.T) {
		testCalculoVR(t, agent)
	})
	
	t.Run("Cenario4_WorkflowCompleto", func(t *testing.T) {
		testWorkflowCompleto(t, agent)
	})
	
	t.Run("Cenario5_MemoryPersistence", func(t *testing.T) {
		testMemoryPersistence(t, agent)
	})
}

// testConsultaBasica testa consulta básica sobre dados
func testConsultaBasica(t *testing.T, agent *VRAgent) {
	// Teste de execução da ferramenta ReadExcel
	ativosPath := filepath.Join("../../files", "ATIVOS.xlsx")
	result, err := agent.ExecuteTool("read_excel", ativosPath)
	
	if err != nil {
		// É esperado que falhe sem arquivo real, mas deve retornar erro específico
		if !strings.Contains(err.Error(), "arquivo não encontrado") &&
		   !strings.Contains(err.Error(), "no such file") {
			t.Errorf("Erro inesperado ao executar ReadExcel: %v", err)
		}
		t.Logf("Resultado esperado - arquivo não encontrado: %v", err)
	} else {
		t.Logf("ReadExcel executado com sucesso: %s", result)
	}
	
	// Verificar se a ferramenta está disponível
	tools := agent.GetAvailableTools()
	found := false
	for _, tool := range tools {
		if tool == "read_excel" {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Ferramenta read_excel não está disponível")
	}
}

// testValidacaoDados testa validação de dados
func testValidacaoDados(t *testing.T, agent *VRAgent) {
	// Teste de execução da ferramenta ValidateData
	testData := `{"planilha": "ATIVOS.xlsx", "dados": {"teste": "valor"}}`
	
	result, err := agent.ExecuteTool("validate_data", testData)
	if err != nil {
		t.Logf("Erro esperado na validação (dados de teste): %v", err)
	} else {
		t.Logf("Validação executada: %s", result)
	}
	
	// Verificar se a ferramenta retorna estrutura esperada
	if result != "" && !strings.Contains(result, "status") {
		t.Error("Resultado da validação deveria conter 'status'")
	}
}

// testCalculoVR testa cálculo de VR para colaborador
func testCalculoVR(t *testing.T, agent *VRAgent) {
	// Teste de cálculo de VR com dados de exemplo
	testInput := `{
		"matricula": "12345",
		"sindicato": "SINDPD",
		"diasUteis": 22,
		"afastamentos": []
	}`
	
	result, err := agent.ExecuteTool("calculate_vr", testInput)
	if err != nil {
		t.Logf("Erro no cálculo VR (dados de teste): %v", err)
	} else {
		t.Logf("Cálculo VR executado: %s", result)
		
		// Verificar se o resultado contém campos esperados
		expectedFields := []string{"valorTotal", "valorEmpresa", "valorColaborador"}
		for _, field := range expectedFields {
			if !strings.Contains(result, field) {
				t.Errorf("Resultado deveria conter campo '%s'", field)
			}
		}
	}
}

// testWorkflowCompleto testa execução de workflow completo
func testWorkflowCompleto(t *testing.T, agent *VRAgent) {
	// Testar workflow de processamento VR
	err := agent.ExecuteWorkflow("processar-vr-mensal")
	if err != nil {
		t.Errorf("Erro ao executar workflow processar-vr-mensal: %v", err)
	} else {
		t.Log("Workflow processar-vr-mensal executado com sucesso")
	}
	
	// Testar workflow de validação
	err = agent.ExecuteWorkflow("validar-dados")
	if err != nil {
		t.Errorf("Erro ao executar workflow validar-dados: %v", err)
	} else {
		t.Log("Workflow validar-dados executado com sucesso")
	}
	
	// Verificar se as estatísticas foram atualizadas
	status := agent.GetStatus()
	if status.TotalRequests < 2 {
		t.Error("Número de requests deveria ter aumentado após workflows")
	}
}

// testMemoryPersistence testa persistência de memória
func testMemoryPersistence(t *testing.T, agent *VRAgent) {
	// Limpar memória inicial
	err := agent.ClearMemory()
	if err != nil {
		t.Fatalf("Erro ao limpar memória: %v", err)
	}
	
	// Verificar memória vazia
	memory, err := agent.GetMemory()
	if err != nil {
		t.Fatalf("Erro ao obter memória: %v", err)
	}
	
	if len(memory) != 0 {
		t.Error("Memória deveria estar vazia após ClearMemory()")
	}
	
	// Testar reset completo
	err = agent.Reset()
	if err != nil {
		t.Errorf("Erro ao resetar agente: %v", err)
	}
	
	// Verificar se status foi resetado
	status := agent.GetStatus()
	if status.TotalRequests != 0 {
		t.Error("TotalRequests deveria ser 0 após reset")
	}
}

// TestAgentPerformance testa performance básica do agente
func TestAgentPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de performance em modo short")
	}
	
	// Setup
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Teste de performance para múltiplas operações
	start := time.Now()
	
	for i := 0; i < 10; i++ {
		err := agent.ExecuteWorkflow("validar-dados")
		if err != nil {
			t.Logf("Erro esperado no workflow %d: %v", i, err)
		}
	}
	
	duration := time.Since(start)
	t.Logf("Tempo para 10 workflows: %v", duration)
	
	// Verificar se está dentro do limite aceitável (< 30s para 10 operações)
	if duration > 30*time.Second {
		t.Errorf("Performance abaixo do esperado: %v > 30s", duration)
	}
}

// TestAgentConcurrency testa uso concorrente do agente
func TestAgentConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de concorrência em modo short")
	}
	
	// Setup
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Teste de concorrência
	numGoroutines := 5
	done := make(chan bool, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			// Executar operações em paralelo
			tools := agent.GetAvailableTools()
			if len(tools) == 0 {
				t.Errorf("Goroutine %d: Nenhuma ferramenta disponível", id)
				return
			}
			
			status := agent.GetStatus()
			if status.State == "" {
				t.Errorf("Goroutine %d: Status vazio", id)
				return
			}
			
			t.Logf("Goroutine %d executada com sucesso", id)
		}(i)
	}
	
	// Aguardar conclusão
	timeout := time.After(10 * time.Second)
	completed := 0
	
	for completed < numGoroutines {
		select {
		case <-done:
			completed++
		case <-timeout:
			t.Fatalf("Timeout: apenas %d/%d goroutines completaram", completed, numGoroutines)
		}
	}
	
	t.Logf("Teste de concorrência completado: %d goroutines", completed)
}

// TestAgentIntegrationWithRealFiles testa integração com arquivos reais
func TestAgentIntegrationWithRealFiles(t *testing.T) {
	// Setup
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Lista de arquivos para testar
	testFiles := []string{
		"ATIVOS.xlsx",
		"AFASTAMENTOS.xlsx",
		"Base sindicato x valor.xlsx",
		"Base dias uteis.xlsx",
	}
	
	for _, filename := range testFiles {
		t.Run("Arquivo_"+filename, func(t *testing.T) {
			filePath := filepath.Join("../../files", filename)
			
			// Testar leitura do arquivo
			result, err := agent.ExecuteTool("read_excel", filePath)
			if err != nil {
				// Log do erro (esperado se arquivo não existir no ambiente de teste)
				t.Logf("Erro ao ler %s: %v", filename, err)
			} else {
				t.Logf("Arquivo %s lido com sucesso: %.100s...", filename, result)
			}
		})
	}
}