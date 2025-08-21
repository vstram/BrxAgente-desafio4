package agent

import (
	"testing"
	"time"
	
	"BrxAgente-desafio4/internal/chat"
	"BrxAgente-desafio4/internal/config"
)

func TestNewVRAgent(t *testing.T) {
	// Criar configuração de teste
	cfg := &config.Config{
		AgentConfig: config.AgentConfig{
			Enabled:        true,
			Model:          "gpt-3.5-turbo",
			Temperature:    0.7,
			MaxTokens:      1000,
			WorkerPoolSize: 2,
			CacheEnabled:   true,
			CacheSize:      100,
			ToolsEnabled:   []string{"excel"},
		},
	}
	
	// Criar chat service mock
	chatSvc := chat.NewChat(cfg)
	
	// Converter para AgentConfig
	agentConfig := &AgentConfig{
		Enabled:        cfg.AgentConfig.Enabled,
		Model:          cfg.AgentConfig.Model,
		Temperature:    cfg.AgentConfig.Temperature,
		MaxTokens:      cfg.AgentConfig.MaxTokens,
		Timeout:        30 * time.Second,
		MemorySize:     100,
		MemoryTTL:      24 * time.Hour,
		ContextWindow:  4000,
		WorkerPoolSize: cfg.AgentConfig.WorkerPoolSize,
		CacheEnabled:   cfg.AgentConfig.CacheEnabled,
		CacheSize:      cfg.AgentConfig.CacheSize,
		CacheTTL:       1 * time.Hour,
		LogLevel:       "info",
		DebugMode:      false,
		ToolsEnabled:   cfg.AgentConfig.ToolsEnabled,
	}
	
	// Criar agente
	agent, err := NewVRAgent(agentConfig, chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Verificar se agente foi criado corretamente
	if agent == nil {
		t.Fatal("Agente não foi criado")
	}
	
	if !agent.IsEnabled() {
		t.Error("Agente deveria estar habilitado")
	}
	
	status := agent.GetStatus()
	if status.State != "idle" {
		t.Errorf("Estado inicial deveria ser 'idle', mas é '%s'", status.State)
	}
}

func TestAgentConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *AgentConfig
		wantErr bool
	}{
		{
			name:    "configuração válida",
			config:  DefaultAgentConfig(),
			wantErr: false,
		},
		{
			name: "temperatura inválida (muito alta)",
			config: &AgentConfig{
				Temperature:    3.0, // Inválida
				MaxTokens:      1000,
				MemorySize:     100,
				WorkerPoolSize: 4,
			},
			wantErr: true,
		},
		{
			name: "max_tokens inválido",
			config: &AgentConfig{
				Temperature:    0.7,
				MaxTokens:      0, // Inválido
				MemorySize:     100,
				WorkerPoolSize: 4,
			},
			wantErr: true,
		},
		{
			name: "memory_size inválido",
			config: &AgentConfig{
				Temperature:    0.7,
				MaxTokens:      1000,
				MemorySize:     0, // Inválido
				WorkerPoolSize: 4,
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVRAgent_EnableDisable(t *testing.T) {
	// Criar agente de teste
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Testar disable
	agent.Disable()
	if agent.IsEnabled() {
		t.Error("Agente deveria estar desabilitado")
	}
	
	// Testar enable
	agent.Enable()
	if !agent.IsEnabled() {
		t.Error("Agente deveria estar habilitado")
	}
}

func TestVRAgent_Reset(t *testing.T) {
	// Criar agente de teste
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	// Simular algumas operações
	agent.status.TotalRequests = 10
	agent.status.ErrorCount = 2
	
	// Reset
	err = agent.Reset()
	if err != nil {
		t.Errorf("Erro ao resetar agente: %v", err)
	}
	
	// Verificar se foi resetado
	status := agent.GetStatus()
	if status.TotalRequests != 0 {
		t.Errorf("TotalRequests deveria ser 0, mas é %d", status.TotalRequests)
	}
	if status.ErrorCount != 0 {
		t.Errorf("ErrorCount deveria ser 0, mas é %d", status.ErrorCount)
	}
}

func TestDefaultAgentConfig(t *testing.T) {
	config := DefaultAgentConfig()
	
	// Verificar valores padrão
	if !config.Enabled {
		t.Error("Configuração padrão deveria ter agente habilitado")
	}
	
	if config.Model != "gpt-3.5-turbo" {
		t.Errorf("Modelo padrão deveria ser 'gpt-3.5-turbo', mas é '%s'", config.Model)
	}
	
	if config.Temperature != 0.7 {
		t.Errorf("Temperatura padrão deveria ser 0.7, mas é %f", config.Temperature)
	}
	
	// Validar configuração padrão
	if err := config.Validate(); err != nil {
		t.Errorf("Configuração padrão deveria ser válida: %v", err)
	}
}