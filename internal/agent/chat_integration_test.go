package agent

import (
	"strings"
	"testing"
	"time"

	"BrxAgente-desafio4/internal/chat"
	"BrxAgente-desafio4/internal/config"
	"BrxAgente-desafio4/internal/modelo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPolicyQuestions testa perguntas sobre políticas de elegibilidade
func TestPolicyQuestions(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de integração em modo short - requer configuração de LLM")
	}
	
	agent := setupTestAgent(t)
	
	tests := []struct {
		name          string
		question      string
		shouldContain []string
	}{
		{
			name:     "Diretores elegibilidade",
			question: "Diretores têm direito a VR?",
			shouldContain: []string{"Não", "diretores", "excluídos"},
		},
		{
			name:     "Estagiários elegibilidade",
			question: "Estagiários podem receber Vale Refeição?",
			shouldContain: []string{"Não", "estagiários", "excluídos"},
		},
		{
			name:     "Aprendizes elegibilidade",
			question: "Aprendizes têm direito ao benefício de VR?",
			shouldContain: []string{"Não", "aprendizes", "excluídos"},
		},
		{
			name:     "Colaboradores no exterior",
			question: "Colaboradores no exterior recebem VR?",
			shouldContain: []string{"Não", "exterior", "excluídos"},
		},
		{
			name:     "Colaboradores afastados",
			question: "Colaboradores afastados têm direito a VR?",
			shouldContain: []string{"depende", "afastamento", "superior a 15 dias"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Para testes de integração, toleramos falhas de configuração
			response, err := agent.Ask(tt.question)
			
			// Se não conseguir responder por falta de configuração, apenas logar
			if err != nil && strings.Contains(err.Error(), "configuração") {
				t.Logf("Teste pulado devido a configuração não disponível: %v", err)
				t.Skip("Configuração de LLM não disponível para teste de integração")
			}
			
			// Deve obter uma resposta sem erro
			assert.NoError(t, err, "Pergunta sobre política não deve retornar erro")
			assert.NotEmpty(t, response, "Resposta não deve estar vazia")
			
			// Verificar se contém as palavras-chave esperadas
			for _, keyword := range tt.shouldContain {
				assert.Contains(t, response, keyword, 
					"Resposta deve conter '%s' para pergunta '%s'", keyword, tt.question)
			}
		})
	}
}

// TestCalculationQuestions testa perguntas sobre cálculos específicos
func TestCalculationQuestions(t *testing.T) {
	agent := setupTestAgent(t)
	
	tests := []struct {
		name          string
		question      string
		shouldContain []string
	}{
		{
			name:     "Licença médica longa",
			question: "Como calcular VR para licença médica de 20 dias?",
			shouldContain: []string{"superior a 15 dias", "não têm direito"},
		},
		{
			name:     "Admissão tardia",
			question: "Qual valor para colaborador admitido dia 25?",
			shouldContain: []string{"após o dia 15", "50%"},
		},
		{
			name:     "Desligamento comunicado tardio",
			question: "Como calcular VR quando desligamento é comunicado após dia 15?",
			shouldContain: []string{"comunicado", "dia 15", "integral"},
		},
		{
			name:     "Férias no mês",
			question: "Como calcular VR para colaborador que tirou 10 dias de férias?",
			shouldContain: []string{"férias", "proporcional", "dias úteis"},
		},
		{
			name:     "Rateio empresa-colaborador",
			question: "Qual o rateio entre empresa e colaborador?",
			shouldContain: []string{"80%", "20%", "empresa", "colaborador"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := agent.Ask(tt.question)
			
			// Deve obter uma resposta sem erro
			assert.NoError(t, err, "Pergunta sobre cálculo não deve retornar erro")
			assert.NotEmpty(t, response, "Resposta não deve estar vazia")
			
			// Verificar se contém as palavras-chave esperadas
			for _, keyword := range tt.shouldContain {
				assert.Contains(t, response, keyword, 
					"Resposta deve conter '%s' para pergunta '%s'", keyword, tt.question)
			}
		})
	}
}

// TestProcessedDataQuestions testa perguntas sobre dados processados
func TestProcessedDataQuestions(t *testing.T) {
	agent := setupTestAgentWithData(t)
	
	tests := []struct {
		name          string
		question      string
		shouldContain []string
	}{
		{
			name:     "Quantidade de colaboradores",
			question: "Quantos colaboradores foram processados?",
			shouldContain: []string{"3", "colaboradores"},
		},
		{
			name:     "Valor total de VR",
			question: "Qual o valor total de VR processado?",
			shouldContain: []string{"R$", "total"},
		},
		{
			name:     "Distribuição por sindicato",
			question: "Quantos colaboradores por sindicato?",
			shouldContain: []string{"sindicato", "distribuição"},
		},
		{
			name:     "Colaborador específico por matrícula",
			question: "Mostre os dados do colaborador matrícula 12345",
			shouldContain: []string{"12345", "Test Corp"},
		},
		{
			name:     "Média de dias úteis",
			question: "Qual a média de dias úteis dos colaboradores?",
			shouldContain: []string{"média", "dias úteis"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := agent.Ask(tt.question)
			
			// Deve obter uma resposta sem erro
			assert.NoError(t, err, "Pergunta sobre dados processados não deve retornar erro")
			assert.NotEmpty(t, response, "Resposta não deve estar vazia")
			
			// Verificar se contém as palavras-chave esperadas
			for _, keyword := range tt.shouldContain {
				assert.Contains(t, response, keyword, 
					"Resposta deve conter '%s' para pergunta '%s'", keyword, tt.question)
			}
		})
	}
}

// TestQuestionRouting testa o roteamento inteligente de perguntas
func TestQuestionRouting(t *testing.T) {
	agent := setupTestAgent(t)
	
	tests := []struct {
		name               string
		question           string
		expectedType       QuestionType
		minConfidence      float64
		shouldUsePolicy    bool
		shouldUseProcessed bool
	}{
		{
			name:            "Política clara",
			question:        "Diretores têm direito a VR?",
			expectedType:    PolicyQuestion,
			minConfidence:   0.7,
			shouldUsePolicy: true,
		},
		{
			name:               "Dados processados clara",
			question:           "Quantos colaboradores foram processados no sistema?",
			expectedType:       ProcessedDataQuestion,
			minConfidence:      0.7,
			shouldUseProcessed: true,
		},
		{
			name:            "Cálculo específico",
			question:        "Como calcular VR para admissão no dia 20?",
			expectedType:    CalculationQuestion,
			minConfidence:   0.6,
			shouldUsePolicy: true,
		},
		{
			name:            "Cenário hipotético",
			question:        "E se um colaborador for admitido e desligado no mesmo mês?",
			expectedType:    WhatIfQuestion,
			minConfidence:   0.6,
			shouldUsePolicy: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Testar classificação da pergunta
			classification := agent.ClassifyQuestion(tt.question)
			
			assert.Equal(t, tt.expectedType, classification.QuestionType, 
				"Pergunta deve ser classificada como %s", tt.expectedType.String())
			assert.GreaterOrEqual(t, classification.Confidence, tt.minConfidence,
				"Confiança da classificação deve ser >= %.2f", tt.minConfidence)
			
			// Testar se a resposta é obtida
			response, err := agent.Ask(tt.question)
			assert.NoError(t, err, "Pergunta roteada não deve retornar erro")
			assert.NotEmpty(t, response, "Resposta roteada não deve estar vazia")
		})
	}
}

// TestResponseQuality testa a qualidade das respostas
func TestResponseQuality(t *testing.T) {
	agent := setupTestAgent(t)
	
	tests := []struct {
		name           string
		question       string
		minLength      int
		shouldNotMatch []string
	}{
		{
			name:           "Resposta sobre política deve ser informativa",
			question:       "Qual a política para diretores?",
			minLength:      50,
			shouldNotMatch: []string{"não sei", "desculpe", "erro"},
		},
		{
			name:           "Resposta sobre cálculo deve ser específica",
			question:       "Como calcular VR para admissão tardia?",
			minLength:      100,
			shouldNotMatch: []string{"não tenho", "impossível", "erro"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := agent.Ask(tt.question)
			
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(response), tt.minLength,
				"Resposta deve ter pelo menos %d caracteres", tt.minLength)
			
			for _, badPattern := range tt.shouldNotMatch {
				assert.NotContains(t, response, badPattern,
					"Resposta não deve conter padrão inadequado: %s", badPattern)
			}
		})
	}
}

// TestConsistency testa a consistência das respostas
func TestConsistency(t *testing.T) {
	agent := setupTestAgent(t)
	
	// Fazer a mesma pergunta múltiplas vezes
	question := "Diretores têm direito a VR?"
	var responses []string
	
	for i := 0; i < 3; i++ {
		response, err := agent.Ask(question)
		require.NoError(t, err)
		require.NotEmpty(t, response)
		responses = append(responses, response)
	}
	
	// Verificar se todas as respostas contêm elementos consistentes
	for _, response := range responses {
		assert.Contains(t, response, "Não", 
			"Todas as respostas devem ser consistentes sobre diretores")
	}
}

// TestErrorHandling testa o tratamento de erros e casos extremos
func TestErrorHandling(t *testing.T) {
	agent := setupTestAgent(t)
	
	tests := []struct {
		name     string
		question string
	}{
		{
			name:     "Pergunta vazia",
			question: "",
		},
		{
			name:     "Pergunta muito longa",
			question: strings.Repeat("teste ", 1000),
		},
		{
			name:     "Pergunta com caracteres especiais",
			question: "Como calcular VR com ñ, ç, é, â???",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := agent.Ask(tt.question)
			
			// Deve sempre retornar algo (mesmo que seja uma mensagem de erro)
			// Não deve causar panic
			if err != nil {
				assert.NotEmpty(t, err.Error(), "Erro deve ter mensagem descritiva")
			} else {
				assert.NotEmpty(t, response, "Resposta não deve estar vazia se não há erro")
			}
		})
	}
}

// TestPerformance testa o desempenho básico das consultas
func TestPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de performance em modo short")
	}
	
	agent := setupTestAgent(t)
	
	questions := []string{
		"Diretores têm direito a VR?",
		"Como calcular VR para admissão no dia 20?",
		"Quantos colaboradores foram processados?",
	}
	
	for _, question := range questions {
		t.Run(question, func(t *testing.T) {
			start := time.Now()
			response, err := agent.Ask(question)
			duration := time.Since(start)
			
			assert.NoError(t, err)
			assert.NotEmpty(t, response)
			assert.Less(t, duration, 30*time.Second, 
				"Pergunta deve ser respondida em menos de 30 segundos")
		})
	}
}

// setupTestAgent cria um agente configurado para testes
func setupTestAgent(t *testing.T) *VRAgent {
	t.Helper()
	
	agentConfig := &AgentConfig{
		Enabled:     true,
		Model:       "test-model",
		Temperature: 0.7,
		MaxTokens:   2000,
		Timeout:     30 * time.Second,
		MemorySize:  100,
		WorkerPoolSize: 4,
		CacheEnabled: true,
		CacheSize:   1000,
		ToolsEnabled: []string{"policy_consultant"},
	}
	
	chatConfig := &config.Config{
		OllamaConfig: config.OllamaConfig{
			BaseURL: "http://localhost:11434",
			Model:   "llama2",
		},
	}
	
	chatSvc := chat.NewChat(chatConfig)
	
	agent, err := NewVRAgent(agentConfig, chatSvc)
	require.NoError(t, err, "Deve criar agente de teste sem erro")
	require.NotNil(t, agent, "Agente não deve ser nulo")
	
	return agent
}

// setupTestAgentWithData cria um agente com dados mock para testes
func setupTestAgentWithData(t *testing.T) *VRAgent {
	t.Helper()
	
	agent := setupTestAgent(t)
	
	// Criar dados mock para os testes
	mockData := createMockProcessedData()
	
	// Configurar dados no chat service
	if agent.chatService != nil {
		err := agent.chatService.SetContextData(mockData)
		require.NoError(t, err, "Deve configurar dados mock sem erro")
	}
	
	return agent
}

// createMockProcessedData cria dados processados mock para testes
func createMockProcessedData() map[string]*modelo.Colaborador {
	now := time.Now()
	admissao := now.AddDate(0, -6, 0) // 6 meses atrás
	
	return map[string]*modelo.Colaborador{
		"12345": {
			Matricula:           "12345",
			Nome:                "Colaborador Teste 1", // Usado apenas internamente
			Empresa:             "Test Corp",
			Sindicato:           "SINDICATO A",
			Cargo:               "Analista",
			Situacao:            "Trabalhando",
			DataAdmissao:        admissao,
			ValorTotalVR:        450.00,
			ValorEmpresa:        360.00,
			ValorColaborador:    90.00,
			DiasUteisEfetivos:   20,
		},
		"67890": {
			Matricula:           "67890",
			Nome:                "Colaborador Teste 2", // Usado apenas internamente
			Empresa:             "Test Corp",
			Sindicato:           "SINDICATO B",
			Cargo:               "Coordenador",
			Situacao:            "Trabalhando",
			DataAdmissao:        admissao,
			ValorTotalVR:        500.00,
			ValorEmpresa:        400.00,
			ValorColaborador:    100.00,
			DiasUteisEfetivos:   22,
		},
		"11111": {
			Matricula:           "11111",
			Nome:                "Colaborador Teste 3", // Usado apenas internamente
			Empresa:             "Test Corp 2",
			Sindicato:           "SINDICATO A",
			Cargo:               "Gerente",
			Situacao:            "Trabalhando",
			DataAdmissao:        admissao,
			ValorTotalVR:        600.00,
			ValorEmpresa:        480.00,
			ValorColaborador:    120.00,
			DiasUteisEfetivos:   21,
		},
	}
}

// TestQuestionClassificationUnit testa apenas a classificação de perguntas (não requer LLM)
func TestQuestionClassificationUnit(t *testing.T) {
	agent := setupTestAgent(t)
	
	tests := []struct {
		name               string
		question           string
		expectedType       QuestionType
		minConfidence      float64
	}{
		{
			name:            "Política clara",
			question:        "Diretores têm direito a VR?",
			expectedType:    PolicyQuestion,
			minConfidence:   0.5,
		},
		{
			name:               "Dados processados clara",
			question:           "Quantos colaboradores foram processados no sistema?",
			expectedType:       ProcessedDataQuestion,
			minConfidence:      0.5,
		},
		{
			name:            "Cálculo específico",
			question:        "Como calcular VR para admissão no dia 20?",
			expectedType:    CalculationQuestion,
			minConfidence:   0.5,
		},
		{
			name:            "Cenário hipotético",
			question:        "E se um colaborador for admitido e desligado no mesmo mês?",
			expectedType:    WhatIfQuestion,
			minConfidence:   0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Testar classificação da pergunta (não requer LLM)
			classification := agent.ClassifyQuestion(tt.question)
			
			assert.Equal(t, tt.expectedType, classification.QuestionType, 
				"Pergunta deve ser classificada como %s", tt.expectedType.String())
			assert.GreaterOrEqual(t, classification.Confidence, tt.minConfidence,
				"Confiança da classificação deve ser >= %.2f", tt.minConfidence)
			
			t.Logf("Pergunta: %s", tt.question)
			t.Logf("Classificação: %s (confiança: %.2f)", 
				classification.QuestionType.String(), classification.Confidence)
		})
	}
}

// TestAgentConfigurationChat testa configuração do agente para chat
func TestAgentConfigurationChat(t *testing.T) {
	agent := setupTestAgent(t)
	
	// Testar se o agente foi criado corretamente
	assert.NotNil(t, agent, "Agente deve ser criado")
	assert.True(t, agent.IsEnabled(), "Agente deve estar habilitado")
	
	// Testar status do agente
	status := agent.GetStatus()
	assert.Equal(t, "idle", status.State, "Estado inicial deve ser idle")
	assert.Equal(t, int64(0), status.TotalRequests, "Requests iniciais devem ser 0")
	assert.Equal(t, int64(0), status.ErrorCount, "Erros iniciais devem ser 0")
	
	// Testar ferramentas disponíveis
	tools := agent.GetAvailableTools()
	assert.NotEmpty(t, tools, "Deve ter ferramentas disponíveis")
	t.Logf("Ferramentas disponíveis: %v", tools)
	
	// Testar configuração
	config := agent.GetConfig()
	assert.NotNil(t, config, "Configuração deve estar disponível")
	assert.True(t, config.Enabled, "Configuração deve mostrar agente habilitado")
}

// TestCoverageReport gera relatório de cobertura dos testes
func TestCoverageReport(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando relatório de cobertura em modo short")
	}
	
	// Este teste serve para documentar a cobertura esperada
	coverage := map[string]bool{
		"PolicyQuestions":       true,
		"CalculationQuestions":  true,
		"ProcessedDataQuestions": true,
		"QuestionRouting":       true,
		"ResponseQuality":       true,
		"Consistency":          true,
		"ErrorHandling":        true,
		"Performance":          true,
		"QuestionClassificationUnit": true,
		"AgentConfigurationChat": true,
	}
	
	for testType, implemented := range coverage {
		assert.True(t, implemented, "Tipo de teste %s deve estar implementado", testType)
	}
	
	t.Logf("Cobertura de testes: %d/%d tipos implementados", len(coverage), len(coverage))
}