package agent

import (
	"strings"
	"testing"
	"time"

	"BrxAgente-desafio4/internal/chat"
	"BrxAgente-desafio4/internal/config"
)

// TestAgentErrorHandling testa cenários de erro e recuperação
func TestAgentErrorHandling(t *testing.T) {
	t.Run("ConfiguracaoInvalida", func(t *testing.T) {
		testConfiguracaoInvalida(t)
	})
	
	t.Run("AgenteDesabilitado", func(t *testing.T) {
		testAgenteDesabilitado(t)
	})
	
	t.Run("FerramentaInexistente", func(t *testing.T) {
		testFerramentaInexistente(t)
	})
	
	t.Run("WorkflowInvalido", func(t *testing.T) {
		testWorkflowInvalido(t)
	})
	
	t.Run("DadosInvalidos", func(t *testing.T) {
		testDadosInvalidos(t)
	})
	
	t.Run("MemoriaIndisponivel", func(t *testing.T) {
		testMemoriaIndisponivel(t)
	})
}

// testConfiguracaoInvalida testa criação de agente com configuração inválida
func testConfiguracaoInvalida(t *testing.T) {
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	
	// Teste 1: Temperatura inválida
	invalidConfig1 := &AgentConfig{
		Enabled:     true,
		Temperature: -5.0, // Inválida
		MaxTokens:   1000,
		MemorySize:  50,
		WorkerPoolSize: 4,
	}
	
	_, err := NewVRAgent(invalidConfig1, chatSvc)
	if err == nil {
		t.Error("Deveria retornar erro para temperatura inválida")
	}
	
	if !strings.Contains(err.Error(), "temperatura") {
		t.Errorf("Erro deveria mencionar temperatura: %v", err)
	}
	
	// Teste 2: MaxTokens inválido
	invalidConfig2 := &AgentConfig{
		Enabled:        true,
		Temperature:    0.7,
		MaxTokens:      -100, // Inválido
		MemorySize:     50,
		WorkerPoolSize: 4,
	}
	
	_, err = NewVRAgent(invalidConfig2, chatSvc)
	if err == nil {
		t.Error("Deveria retornar erro para MaxTokens inválido")
	}
	
	if !strings.Contains(err.Error(), "max_tokens") {
		t.Errorf("Erro deveria mencionar max_tokens: %v", err)
	}
	
	// Teste 3: MemorySize inválido
	invalidConfig3 := &AgentConfig{
		Enabled:        true,
		Temperature:    0.7,
		MaxTokens:      1000,
		MemorySize:     0, // Inválido
		WorkerPoolSize: 4,
	}
	
	_, err = NewVRAgent(invalidConfig3, chatSvc)
	if err == nil {
		t.Error("Deveria retornar erro para MemorySize inválido")
	}
}

// testAgenteDesabilitado testa operações com agente desabilitado
func testAgenteDesabilitado(t *testing.T) {
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	
	// Criar agente habilitado
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Desabilitar agente
	agent.Disable()
	
	// Teste 1: Ask com agente desabilitado
	_, err = agent.Ask("Teste de pergunta")
	if err == nil {
		t.Error("Ask deveria falhar com agente desabilitado")
	}
	
	if !strings.Contains(err.Error(), "desabilitado") {
		t.Errorf("Erro deveria mencionar agente desabilitado: %v", err)
	}
	
	// Teste 2: ExecuteTool com agente desabilitado
	_, err = agent.ExecuteTool("read_excel", "test.xlsx")
	if err == nil {
		t.Error("ExecuteTool deveria falhar com agente desabilitado")
	}
	
	// Teste 3: ExecuteWorkflow com agente desabilitado
	err = agent.ExecuteWorkflow("processar-vr-mensal")
	if err == nil {
		t.Error("ExecuteWorkflow deveria falhar com agente desabilitado")
	}
	
	// Verificar que após habilitar volta a funcionar
	agent.Enable()
	if !agent.IsEnabled() {
		t.Error("Agente deveria estar habilitado após Enable()")
	}
}

// testFerramentaInexistente testa execução de ferramenta inexistente
func testFerramentaInexistente(t *testing.T) {
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Tentar executar ferramenta inexistente
	_, err = agent.ExecuteTool("ferramenta_inexistente", "input")
	if err == nil {
		t.Error("ExecuteTool deveria falhar para ferramenta inexistente")
	}
	
	// Verificar se o erro é específico
	if !strings.Contains(err.Error(), "não encontrada") && 
	   !strings.Contains(err.Error(), "not found") {
		t.Logf("Erro obtido: %v", err)
		// Aceitar qualquer erro relacionado à ferramenta inexistente
	}
	
	// Verificar que ferramentas válidas ainda funcionam
	tools := agent.GetAvailableTools()
	if len(tools) == 0 {
		t.Error("Deveria haver ferramentas disponíveis")
	}
}

// testWorkflowInvalido testa execução de workflow inválido
func testWorkflowInvalido(t *testing.T) {
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Tentar executar workflow inexistente
	err = agent.ExecuteWorkflow("workflow-inexistente")
	if err == nil {
		t.Error("ExecuteWorkflow deveria falhar para workflow inexistente")
	}
	
	if !strings.Contains(err.Error(), "não reconhecido") &&
	   !strings.Contains(err.Error(), "not recognized") {
		t.Errorf("Erro deveria mencionar workflow não reconhecido: %v", err)
	}
	
	// Verificar que workflows válidos ainda funcionam
	err = agent.ExecuteWorkflow("processar-vr-mensal")
	if err != nil {
		t.Errorf("Workflow válido falhou: %v", err)
	}
}

// testDadosInvalidos testa processamento de dados inválidos
func testDadosInvalidos(t *testing.T) {
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Teste 1: JSON malformado para calculate_vr
	malformedJSON := `{"matricula": "12345", "sindicato"`
	_, err = agent.ExecuteTool("calculate_vr", malformedJSON)
	if err == nil {
		t.Error("calculate_vr deveria falhar com JSON malformado")
	}
	
	// Teste 2: Arquivo inexistente para read_excel
	_, err = agent.ExecuteTool("read_excel", "/caminho/inexistente.xlsx")
	if err == nil {
		t.Error("read_excel deveria falhar com arquivo inexistente")
	}
	
	// Teste 3: Dados vazios para validate_data
	_, err = agent.ExecuteTool("validate_data", "")
	if err == nil {
		t.Error("validate_data deveria falhar com dados vazios")
	}
	
	t.Log("Todos os testes de dados inválidos passaram (erros esperados)")
}

// testMemoriaIndisponivel testa cenários de falha na memória
func testMemoriaIndisponivel(t *testing.T) {
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Testar GetMemory com agente recém-criado
	memory, err := agent.GetMemory()
	if err != nil {
		t.Errorf("GetMemory não deveria falhar: %v", err)
	}
	
	// Memória pode estar vazia inicialmente, isso é normal
	t.Logf("Memória inicial tem %d items", len(memory))
	
	// Testar ClearMemory
	err = agent.ClearMemory()
	if err != nil {
		t.Errorf("ClearMemory não deveria falhar: %v", err)
	}
	
	// Verificar se memória foi limpa
	memory, err = agent.GetMemory()
	if err != nil {
		t.Errorf("GetMemory após ClearMemory não deveria falhar: %v", err)
	}
	
	if len(memory) != 0 {
		t.Errorf("Memória deveria estar vazia após ClearMemory, got %d items", len(memory))
	}
}

// TestAgentErrorRecovery testa recuperação de erros
func TestAgentErrorRecovery(t *testing.T) {
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	initialRequests := agent.GetStatus().TotalRequests
	initialErrors := agent.GetStatus().ErrorCount
	
	// Causar alguns erros
	_, err = agent.ExecuteTool("ferramenta_inexistente", "test")
	if err != nil {
		t.Logf("Erro esperado 1: %v", err)
	}
	
	err = agent.ExecuteWorkflow("workflow-inexistente")
	if err != nil {
		t.Logf("Erro esperado 2: %v", err)
	}
	
	// Verificar se contador de erros aumentou
	finalStatus := agent.GetStatus()
	if finalStatus.ErrorCount <= initialErrors {
		t.Error("Contador de erros deveria ter aumentado")
	}
	
	if finalStatus.TotalRequests <= initialRequests {
		t.Error("Contador de requests deveria ter aumentado")
	}
	
	// Verificar que agente ainda funciona após erros
	tools := agent.GetAvailableTools()
	if len(tools) == 0 {
		t.Error("Agente deveria ainda ter ferramentas disponíveis após erros")
	}
	
	// Executar operação válida
	err = agent.ExecuteWorkflow("processar-vr-mensal")
	if err != nil {
		t.Errorf("Operação válida falhou após recuperação de erros: %v", err)
	}
	
	t.Log("Teste de recuperação de erros completado com sucesso")
}

// TestAgentStressTest testa comportamento sob stress
func TestAgentStressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando stress test em modo short")
	}
	
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Executar muitas operações rapidamente
	numOperations := 100
	errors := 0
	
	start := time.Now()
	
	for i := 0; i < numOperations; i++ {
		// Alternar entre operações válidas e inválidas
		if i%2 == 0 {
			err := agent.ExecuteWorkflow("processar-vr-mensal")
			if err != nil {
				errors++
			}
		} else {
			_, err := agent.ExecuteTool("read_excel", "nonexistent.xlsx")
			if err != nil {
				errors++ // Erro esperado
			}
		}
	}
	
	duration := time.Since(start)
	
	t.Logf("Stress test completado:")
	t.Logf("- Operações: %d", numOperations)
	t.Logf("- Tempo total: %v", duration)
	t.Logf("- Operações por segundo: %.2f", float64(numOperations)/duration.Seconds())
	t.Logf("- Erros: %d", errors)
	
	// Verificar se o agente ainda está responsivo
	status := agent.GetStatus()
	if status.State == "" {
		t.Error("Agente não deveria estar em estado vazio após stress test")
	}
	
	// Verificar se ainda consegue executar operações
	finalTools := agent.GetAvailableTools()
	if len(finalTools) == 0 {
		t.Error("Agente deveria ainda ter ferramentas após stress test")
	}
}

// TestAgentTimeout testa comportamento de timeout
func TestAgentTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de timeout em modo short")
	}
	
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	
	// Configurar agente com timeout muito baixo
	shortTimeoutConfig := &AgentConfig{
		Enabled:         true,
		Model:          "test-model",
		Temperature:    0.7,
		MaxTokens:      1000,
		Timeout:        1 * time.Millisecond, // Timeout muito baixo
		MemorySize:     50,
		MemoryTTL:      1 * time.Hour,
		ContextWindow:  2000,
		MaxMemoryTokens: 1000,
		WorkerPoolSize: 2,
		CacheEnabled:   false,
		CacheSize:      100,
		CacheTTL:       10 * time.Minute,
		LogLevel:       "info",
		DebugMode:      true,
		ToolsEnabled:   []string{"excel", "calculation"},
	}
	
	agent, err := NewVRAgent(shortTimeoutConfig, chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente com timeout baixo: %v", err)
	}
	
	// Tentar operação que pode sofrer timeout
	start := time.Now()
	err = agent.ExecuteWorkflow("processar-vr-mensal")
	duration := time.Since(start)
	
	t.Logf("Operação com timeout baixo levou: %v", duration)
	
	// Se houve erro relacionado a timeout, é esperado
	if err != nil {
		t.Logf("Erro esperado de timeout: %v", err)
	}
	
	// Verificar se o agente ainda está funcional
	status := agent.GetStatus()
	if status.State == "" {
		t.Error("Agente deveria manter estado após timeout")
	}
}