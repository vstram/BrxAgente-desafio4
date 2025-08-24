package agent

import (
	"testing"
	"time"

	"BrxAgente-desafio4/internal/chat"
	"BrxAgente-desafio4/internal/config"
)

// TestAgentMemoryManagement demonstra o funcionamento do memory management
func TestAgentMemoryManagement(t *testing.T) {
	// Configurar ambiente de teste
	cfg := &config.Config{
		OpenAIKey: "",
		OllamaConfig: config.OllamaConfig{
			BaseURL: "",
			Model:   "",
		},
	}

	chatSvc := chat.NewChat(cfg)

	// Criar configuração do agente para teste
	agentConfig := &AgentConfig{
		Enabled:         true,
		Model:           "test-model",
		Temperature:     0.7,
		MaxTokens:       1000,
		Timeout:         10 * time.Second,
		MemorySize:      50,
		MemoryTTL:       1 * time.Hour,
		ContextWindow:   2000,
		MaxMemoryTokens: 1000,
		WorkerPoolSize:  2,
		CacheEnabled:    true,
		CacheSize:       100,
		CacheTTL:        30 * time.Minute,
		LogLevel:        "info",
		DebugMode:       true,
		ToolsEnabled:    []string{"excel", "calculation"},
	}

	// Criar agente
	agent, err := NewVRAgent(agentConfig, chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}

	// Verificar se o agente foi criado corretamente
	if !agent.IsEnabled() {
		t.Error("Agente deveria estar habilitado")
	}

	// Verificar status inicial
	status := agent.GetStatus()
	if status.State != "idle" {
		t.Errorf("Estado inicial deveria ser 'idle', got %s", status.State)
	}

	if status.TotalRequests != 0 {
		t.Errorf("Total de requests inicial deveria ser 0, got %d", status.TotalRequests)
	}

	// Verificar memória inicial (deve estar vazia)
	memory, err := agent.GetMemory()
	if err != nil {
		t.Fatalf("Erro ao obter memória: %v", err)
	}

	if len(memory) != 0 {
		t.Errorf("Memória inicial deveria estar vazia, got %d items", len(memory))
	}

	// Testar reset do agente
	err = agent.Reset()
	if err != nil {
		t.Errorf("Erro ao resetar agente: %v", err)
	}

	// Verificar se status foi resetado
	statusAfterReset := agent.GetStatus()
	if statusAfterReset.TotalRequests != 0 {
		t.Errorf("Total de requests após reset deveria ser 0, got %d", statusAfterReset.TotalRequests)
	}

	// Testar clear memory
	err = agent.ClearMemory()
	if err != nil {
		t.Errorf("Erro ao limpar memória: %v", err)
	}

	// Testar workflow execution
	err = agent.ExecuteWorkflow("processar-vr-mensal")
	if err != nil {
		t.Errorf("Erro ao executar workflow: %v", err)
	}

	// Verificar se request foi contabilizado
	statusAfterWorkflow := agent.GetStatus()
	if statusAfterWorkflow.TotalRequests == 0 {
		t.Error("Total de requests deveria ter aumentado após workflow")
	}

	// Testar ferramentas disponíveis
	tools := agent.GetAvailableTools()
	if len(tools) == 0 {
		t.Error("Deveria haver ferramentas disponíveis")
	}

	// Testar informações de ferramentas
	allToolsInfo := agent.GetAllToolsInfo()
	if len(allToolsInfo) == 0 {
		t.Error("Deveria haver informações sobre as ferramentas")
	}

	// Testar disable/enable
	agent.Disable()
	if agent.IsEnabled() {
		t.Error("Agente deveria estar desabilitado")
	}

	agent.Enable()
	if !agent.IsEnabled() {
		t.Error("Agente deveria estar habilitado novamente")
	}

	t.Log("✅ Todos os testes de memory management passaram!")
}

// TestAgentConfiguration testa diferentes configurações do agente
func TestAgentConfiguration(t *testing.T) {
	chatSvc := chat.NewChat(&config.Config{})

	// Testar configuração padrão
	defaultAgent, err := NewVRAgent(nil, chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente com configuração padrão: %v", err)
	}

	defaultConfig := defaultAgent.GetConfig()
	if defaultConfig.Model != "gpt-3.5-turbo" {
		t.Errorf("Modelo padrão deveria ser gpt-3.5-turbo, got %s", defaultConfig.Model)
	}

	if defaultConfig.Temperature != 0.7 {
		t.Errorf("Temperatura padrão deveria ser 0.7, got %f", defaultConfig.Temperature)
	}

	// Testar configuração inválida
	invalidConfig := &AgentConfig{
		Temperature: -1.0, // Temperatura inválida
	}

	_, err = NewVRAgent(invalidConfig, chatSvc)
	if err == nil {
		t.Error("Deveria ter retornado erro para configuração inválida")
	}

	t.Log("✅ Todos os testes de configuração passaram!")
}
