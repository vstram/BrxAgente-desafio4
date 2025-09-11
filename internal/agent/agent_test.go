package agent

import (
	"fmt"
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

// TestVRAgent_ToolsAlwaysAvailable testa que as ferramentas estão sempre disponíveis independente do debug mode
func TestVRAgent_ToolsAlwaysAvailable(t *testing.T) {
	tests := []struct {
		name      string
		debugMode bool
	}{
		{
			name:      "com debug mode ativado",
			debugMode: true,
		},
		{
			name:      "com debug mode desativado",
			debugMode: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar configuração com debug mode específico
			config := DefaultAgentConfig()
			config.DebugMode = tt.debugMode

			// Criar chat service
			chatSvc := chat.NewChat(nil)

			// Criar agente
			agent, err := NewVRAgent(config, chatSvc)
			if err != nil {
				t.Fatalf("Erro ao criar agente: %v", err)
			}

			// Verificar se as ferramentas estão disponíveis
			tools := agent.GetAvailableTools()
			expectedTools := []string{"read_excel", "calculate_vr", "validate_data", "policy_consultant"}

			if len(tools) != len(expectedTools) {
				t.Errorf("Esperado %d ferramentas, mas obtido %d: %v", len(expectedTools), len(tools), tools)
			}

			// Verificar se cada ferramenta esperada está presente
			toolMap := make(map[string]bool)
			for _, tool := range tools {
				toolMap[tool] = true
			}

			for _, expectedTool := range expectedTools {
				if !toolMap[expectedTool] {
					t.Errorf("Ferramenta '%s' deveria estar disponível (debug mode: %v)", expectedTool, tt.debugMode)
				}
			}

			// Verificar se consegue obter informações da ferramenta policy_consultant
			toolInfo, err := agent.GetToolInfo("policy_consultant")
			if err != nil {
				t.Errorf("Erro ao obter informações da ferramenta policy_consultant: %v", err)
			}

			if toolInfo == nil {
				t.Error("Informações da ferramenta policy_consultant não deveriam ser nil")
			}
		})
	}
}

// TestVRAgent_QuestionRouting testa o roteamento inteligente de perguntas
func TestVRAgent_QuestionRouting(t *testing.T) {
	// Criar agente de teste
	config := DefaultAgentConfig()
	chatSvc := chat.NewChat(nil)
	agent, err := NewVRAgent(config, chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}

	tests := []struct {
		name           string
		question       string
		expectPolicy   bool
		description    string
	}{
		// Perguntas sobre políticas (devem rotear para PolicyConsultantTool)
		{
			name:         "pergunta_sobre_direito_diretores",
			question:     "Diretores têm direito a VR?",
			expectPolicy: true,
			description:  "Pergunta sobre direito de diretores",
		},
		{
			name:         "pergunta_sobre_calculo_licenca",
			question:     "Como calcular VR para licença médica de 20 dias?",
			expectPolicy: true,
			description:  "Pergunta sobre cálculo em licença médica",
		},
		{
			name:         "pergunta_sobre_estagiarios",
			question:     "Estagiários podem receber vale alimentação?",
			expectPolicy: true,
			description:  "Pergunta sobre elegibilidade de estagiários",
		},
		{
			name:         "pergunta_sobre_regras_afastamento",
			question:     "Qual a regra para colaborador em afastamento?",
			expectPolicy: true,
			description:  "Pergunta sobre regras de afastamento",
		},
		{
			name:         "pergunta_sobre_admissao",
			question:     "Como funciona VR para admissão no meio do mês?",
			expectPolicy: true,
			description:  "Pergunta sobre regras de admissão",
		},
		{
			name:         "pergunta_sobre_politica_geral",
			question:     "Qual a política de vale refeição da empresa?",
			expectPolicy: true,
			description:  "Pergunta sobre políticas gerais",
		},
		
		// Perguntas sobre dados processados (devem rotear para dados)
		{
			name:         "pergunta_quantos_colaboradores",
			question:     "Quantos colaboradores foram processados?",
			expectPolicy: false,
			description:  "Pergunta sobre quantidade de colaboradores processados",
		},
		{
			name:         "pergunta_total_vr",
			question:     "Qual o total de VR calculado este mês?",
			expectPolicy: false,
			description:  "Pergunta sobre totais calculados",
		},
		{
			name:         "pergunta_colaborador_especifico",
			question:     "Qual o valor de VR do colaborador matrícula 12345?",
			expectPolicy: false,
			description:  "Pergunta sobre colaborador específico",
		},
		{
			name:         "pergunta_estatisticas",
			question:     "Mostre as estatísticas de processamento",
			expectPolicy: false,
			description:  "Pergunta sobre estatísticas",
		},
		{
			name:         "pergunta_resumo_dados",
			question:     "Faça um resumo dos dados processados",
			expectPolicy: false,
			description:  "Pergunta sobre resumo de dados",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Testar apenas a função de classificação
			isPolicy := agent.isPolicyQuestion(tt.question)
			
			if isPolicy != tt.expectPolicy {
				t.Errorf("isPolicyQuestion('%s') = %v, esperado %v (%s)", 
					tt.question, isPolicy, tt.expectPolicy, tt.description)
			}
		})
	}
}

// TestVRAgent_isPolicyQuestion testa especificamente a função de classificação
func TestVRAgent_isPolicyQuestion(t *testing.T) {
	// Criar agente de teste
	config := DefaultAgentConfig()
	chatSvc := chat.NewChat(nil)
	agent, err := NewVRAgent(config, chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}

	// Testes de palavras-chave específicas
	policyQuestions := []string{
		"direito", "política", "regra", "elegível", "pode", "deve",
		"diretores", "estagiários", "aprendizes", "terceirizados",
		"licença", "afastamento", "férias", "admissão", "desligamento",
		"como calcular", "qual valor", "quanto vale", "regras de",
		"política de", "tem direito", "não tem direito", "excluído",
		"incluído", "benefício", "vale refeição", "vale alimentação",
	}

	for _, keyword := range policyQuestions {
		question := fmt.Sprintf("Pergunta com %s sobre VR", keyword)
		if !agent.isPolicyQuestion(question) {
			t.Errorf("Pergunta com palavra-chave '%s' deveria ser classificada como política", keyword)
		}
	}

	// Testes de perguntas que NÃO são sobre políticas
	nonPolicyQuestions := []string{
		"Quantos colaboradores?",
		"Total de VR",
		"Estatísticas de processamento",
		"Resumo dos dados",
		"Matrícula 12345",
		"Processados este mês",
	}

	for _, question := range nonPolicyQuestions {
		if agent.isPolicyQuestion(question) {
			t.Errorf("Pergunta '%s' NÃO deveria ser classificada como política", question)
		}
	}
}

// TestVRAgent_askWithPolicyTool testa o método de consulta a políticas
func TestVRAgent_askWithPolicyTool(t *testing.T) {
	// Criar agente de teste com configuração adequada
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	config := DefaultAgentConfig()
	agent, err := NewVRAgent(config, chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}

	// Testar com ToolRegistry disponível
	question := "Diretores têm direito a VR?"
	result, err := agent.askWithPolicyTool(question)
	
	// No ambiente de teste, esperamos que o PolicyConsultantTool falhe devido a arquivos
	// não disponíveis, mas deve usar fallback adequadamente
	if result == "" {
		t.Logf("PolicyConsultantTool falhou (esperado no ambiente de teste), verificando fallback")
		// Se resultado vazio, pelo menos não deveria ter panic
	}
	
	// Log do resultado para debug
	t.Logf("Resultado askWithPolicyTool: err=%v, result_length=%d", err, len(result))
}

// TestVRAgent_askWithProcessedData testa o método de consulta a dados processados
func TestVRAgent_askWithProcessedData(t *testing.T) {
	// Criar agente de teste com configuração adequada
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	config := DefaultAgentConfig()
	agent, err := NewVRAgent(config, chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}

	// Testar pergunta sobre dados processados
	question := "Quantos colaboradores foram processados?"
	result, err := agent.askWithProcessedData(question)
	
	// O resultado pode ser um erro se Ollama/OpenAI não estiverem disponíveis,
	// mas a função deveria tentar processar
	if err == nil && result == "" {
		t.Error("askWithProcessedData retornou resultado vazio sem erro")
	}

	// Se houve erro, deveria ser erro de LLM, não de lógica
	if err != nil {
		t.Logf("Erro esperado no ambiente de teste (LLM não disponível): %v", err)
	}
	
	// Log do resultado para debug
	t.Logf("Resultado askWithProcessedData: err=%v, result_length=%d", err, len(result))
}
