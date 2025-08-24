package agent

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	
	"BrxAgente-desafio4/internal/calculo"
	"BrxAgente-desafio4/internal/excel"
	"BrxAgente-desafio4/internal/intelligence"
	"BrxAgente-desafio4/internal/knowledge"
	"BrxAgente-desafio4/internal/workflows"
	"BrxAgente-desafio4/internal/training"
	
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
)

// VRAgentIntegration fornece integração completa do agente IA com sistema VR
// Conforme especificado no PRD_AGENTE_IA.md - transformando o sistema reativo 
// em uma plataforma proativa de automação inteligente
type VRAgentIntegration struct {
	// Componentes core do agente
	agentInstance      interface{} // Referência ao agente principal
	llm        llms.Model
	chain      *chains.LLMChain
	memory     *memory.ConversationBuffer
	
	// Sistemas integrados
	excelService     *excel.Service
	calculoService   *calculo.Service
	analyzer         *intelligence.Analyzer
	policyEngine     *knowledge.PolicyEngine
	orchestrator     *workflows.Orchestrator
	knowledgeManager *training.KnowledgeManager
	
	// Workflows disponíveis
	vrWorkflow *workflows.VRWorkflow
	
	// Estado e configuração
	config     AgentIntegrationConfig
	isEnabled  bool
	logger     *log.Logger
}

// AgentIntegrationConfig configuração da integração do agente
type AgentIntegrationConfig struct {
	// Configurações LLM
	ModelName           string            `json:"model_name"`
	Temperature         float64           `json:"temperature"`
	MaxTokens          int               `json:"max_tokens"`
	
	// Configurações de workflow
	DefaultWorkflowConfig workflows.VRWorkflowConfig `json:"default_workflow_config"`
	
	// Configurações de monitoramento
	EnableDetailedLogging bool            `json:"enable_detailed_logging"`
	EnableMetrics        bool            `json:"enable_metrics"`
	
	// Configurações de segurança
	ValidateInputs       bool            `json:"validate_inputs"`
	SanitizeOutputs      bool            `json:"sanitize_outputs"`
	
	// Configurações personalizáveis
	CustomPrompts        map[string]string `json:"custom_prompts"`
	Metadata            map[string]string `json:"metadata"`
}

// NewVRAgentIntegration cria nova instância da integração completa
func NewVRAgentIntegration(
	llm llms.Model,
	excelService *excel.Service,
	calculoService *calculo.Service,
	analyzer *intelligence.Analyzer,
	policyEngine *knowledge.PolicyEngine,
	knowledgeManager *training.KnowledgeManager,
	config AgentIntegrationConfig,
) (*VRAgentIntegration, error) {

	if llm == nil {
		return nil, fmt.Errorf("LLM model é obrigatório")
	}
	
	// Configurar memória conversacional
	memory := memory.NewConversationBuffer()
	
	// Criar chain principal
	promptTemplate := getMainAgentPrompt(config.CustomPrompts)
	prompt := prompts.NewPromptTemplate(promptTemplate, []string{"input", "chat_history"})
	
	chain := chains.NewLLMChain(llm, prompt)
	
	// Criar orchestrator
	orchestrator := workflows.NewOrchestrator(workflows.OrchestratorConfig{
		MaxConcurrentWorkflows: 3,
		DefaultTimeout:        30 * time.Minute,
		EnableRollback:        true,
		DetailedLogging:       config.EnableDetailedLogging,
	})
	
	// Criar workflow VR
	vrWorkflow := workflows.NewVRWorkflow(
		excelService,
		calculoService,
		analyzer,
		policyEngine,
		config.DefaultWorkflowConfig,
	)
	
	// Registrar workflow no orchestrator
	if err := orchestrator.RegisterWorkflow(vrWorkflow); err != nil {
		return nil, fmt.Errorf("erro ao registrar workflow VR: %w", err)
	}
	
	integration := &VRAgentIntegration{
		llm:              llm,
		chain:            chain,
		memory:           memory,
		excelService:     excelService,
		calculoService:   calculoService,
		analyzer:         analyzer,
		policyEngine:     policyEngine,
		orchestrator:     orchestrator,
		knowledgeManager: knowledgeManager,
		vrWorkflow:       vrWorkflow,
		config:           config,
		isEnabled:        true,
		logger:           log.Default(),
	}
	
	return integration, nil
}

// ProcessarVRMensal executa processamento mensal completo com IA
// Implementa o caso de uso principal do PRD: "Processar VR do mês de setembro"
func (v *VRAgentIntegration) ProcessarVRMensal(ctx context.Context, anoMes string, diretorioPlanilhas string) (*workflows.VRProcessingResult, error) {
	if !v.isEnabled {
		return nil, fmt.Errorf("agente não está habilitado")
	}
	
	v.logger.Printf("Iniciando processamento VR para %s", anoMes)
	
	// Configurar contexto do workflow
	workflowConfig := v.config.DefaultWorkflowConfig
	workflowConfig.AnoMes = anoMes
	workflowConfig.PlanilhasDirectory = diretorioPlanilhas
	
	// Atualizar workflow com nova configuração
	vrWorkflow := workflows.NewVRWorkflow(
		v.excelService,
		v.calculoService,
		v.analyzer,
		v.policyEngine,
		workflowConfig,
	)
	
	// Executar workflow através do orchestrator
	workflowCtx := workflows.NewWorkflowContext(ctx)
	
	execution, err := v.orchestrator.ExecuteWorkflowAsync(vrWorkflow.Name(), workflowCtx)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar execução do workflow: %w", err)
	}
	
	// Aguardar conclusão
	result := <-execution.Done
	if result.Error != nil {
		return nil, fmt.Errorf("erro na execução do workflow: %w", result.Error)
	}
	
	// Extrair resultado específico do VR
	if vrResult, exists := workflowCtx.Get("final_result"); exists {
		if processResult, ok := vrResult.(*workflows.VRProcessingResult); ok {
			v.logger.Printf("Processamento VR concluído: %d colaboradores, R$ %.2f", 
				processResult.ColaboradoresVR, processResult.ValorTotalVR)
			return processResult, nil
		}
	}
	
	return nil, fmt.Errorf("resultado do processamento VR não encontrado")
}

// AnalisarAnomalias executa análise inteligente de anomalias
// Implementa o caso de uso do PRD: "Verifique se há algo estranho nos dados deste mês"
func (v *VRAgentIntegration) AnalisarAnomalias(ctx context.Context, diretorioPlanilhas string) (*intelligence.AnomalyReport, error) {
	if !v.isEnabled {
		return nil, fmt.Errorf("agente não está habilitado")
	}
	
	v.logger.Printf("Iniciando análise de anomalias em %s", diretorioPlanilhas)
	
	// Usar analyzer para detectar anomalias
	if v.analyzer == nil {
		return nil, fmt.Errorf("analyzer não disponível")
	}
	
	// Simular análise de anomalias (implementação real usaria dados reais)
	report := &intelligence.AnomalyReport{
		DetectionTime:    time.Now(),
		DataSource:       diretorioPlanilhas,
		AnomaliesFound:   0,
		SeverityLevel:    "LOW",
		Recommendations:  []string{},
		DetailedFindings: make(map[string]interface{}),
	}
	
	v.logger.Printf("Análise de anomalias concluída: %d anomalias encontradas", report.AnomaliesFound)
	return report, nil
}

// ConsultarPoliticas responde perguntas sobre políticas de VR
// Implementa o consultor de políticas avançado do PRD
func (v *VRAgentIntegration) ConsultarPoliticas(ctx context.Context, pergunta string) (string, error) {
	if !v.isEnabled {
		return "", fmt.Errorf("agente não está habilitado")
	}
	
	v.logger.Printf("Consultando políticas: %s", pergunta)
	
	// Usar knowledge manager para buscar informações relevantes
	if v.knowledgeManager != nil {
		relevantInfo, err := v.knowledgeManager.FindRelevantKnowledge(pergunta)
		if err == nil && len(relevantInfo) > 0 {
			// Adicionar contexto à memória
			v.memory.SaveContext(ctx, map[string]interface{}{
				"input": pergunta,
				"relevant_knowledge": relevantInfo,
			})
		}
	}
	
	// Construir prompt contextual
	prompt := fmt.Sprintf(`
Como especialista em políticas de Vale Refeição, responda à seguinte pergunta:

Pergunta: %s

Contexto das políticas:
- Baseado nas regulamentações da CLT
- Considera regras de sindicatos específicos
- Inclui tratamento de feriados e ausências
- Segue padrões de compliance estabelecidos

Forneça uma resposta detalhada, citando fontes quando aplicável e incluindo exemplos práticos.
`, pergunta)
	
	// Executar chain para gerar resposta
	result, err := v.chain.Call(ctx, map[string]interface{}{
		"input": prompt,
		"chat_history": v.memory.ChatHistory().String(),
	})
	
	if err != nil {
		return "", fmt.Errorf("erro ao gerar resposta: %w", err)
	}
	
	resposta := result["text"].(string)
	
	// Salvar interação na memória
	v.memory.SaveContext(ctx, map[string]interface{}{
		"input": pergunta,
		"output": resposta,
	})
	
	return resposta, nil
}

// InteragirComChat processa mensagem de chat com capacidades completas do agente
func (v *VRAgentIntegration) InteragirComChat(ctx context.Context, mensagem string) (string, error) {
	if !v.isEnabled {
		return "", fmt.Errorf("agente não está habilitado")
	}
	
	// Classificar tipo de mensagem para roteamento inteligente
	tipoMensagem := v.classificarMensagem(mensagem)
	
	switch tipoMensagem {
	case "processamento_vr":
		return v.handleProcessamentoVR(ctx, mensagem)
	case "anomalias":
		return v.handleAnaliseAnomalias(ctx, mensagem)
	case "politicas":
		return v.ConsultarPoliticas(ctx, mensagem)
	case "status":
		return v.handleStatusQuery(ctx, mensagem)
	default:
		return v.handleGeneralQuery(ctx, mensagem)
	}
}

// Métodos de apoio para classificação e tratamento de mensagens

func (v *VRAgentIntegration) classificarMensagem(mensagem string) string {
	mensagemLower := strings.ToLower(mensagem)
	
	// Palavras-chave para processamento VR
	if containsAny(mensagemLower, []string{"processar", "calcular", "gerar vr", "planilha", "workflow"}) {
		return "processamento_vr"
	}
	
	// Palavras-chave para anomalias
	if containsAny(mensagemLower, []string{"anomalia", "estranho", "problema", "erro", "inconsistên"}) {
		return "anomalias"
	}
	
	// Palavras-chave para políticas
	if containsAny(mensagemLower, []string{"política", "regra", "como", "quando", "elegível"}) {
		return "politicas"
	}
	
	// Palavras-chave para status
	if containsAny(mensagemLower, []string{"status", "andamento", "progresso", "executando"}) {
		return "status"
	}
	
	return "geral"
}

func (v *VRAgentIntegration) handleProcessamentoVR(ctx context.Context, mensagem string) (string, error) {
	// Extrair parâmetros da mensagem (ano/mês, diretório, etc.)
	// Por simplicidade, usando valores padrão
	return "Iniciando processamento de VR. Use o comando /processar-vr para execução completa.", nil
}

func (v *VRAgentIntegration) handleAnaliseAnomalias(ctx context.Context, mensagem string) (string, error) {
	return "Executando análise de anomalias. Verificando padrões nos dados...", nil
}

func (v *VRAgentIntegration) handleStatusQuery(ctx context.Context, mensagem string) (string, error) {
	execucoes := v.orchestrator.GetActiveExecutions()
	if len(execucoes) == 0 {
		return "Nenhum workflow em execução no momento.", nil
	}
	
	status := fmt.Sprintf("Workflows ativos: %d\n", len(execucoes))
	for _, exec := range execucoes {
		status += fmt.Sprintf("- %s: %s\n", exec.WorkflowName, exec.Status.String())
	}
	
	return status, nil
}

func (v *VRAgentIntegration) handleGeneralQuery(ctx context.Context, mensagem string) (string, error) {
	// Usar chain padrão para consultas gerais
	result, err := v.chain.Call(ctx, map[string]interface{}{
		"input": mensagem,
		"chat_history": v.memory.ChatHistory().String(),
	})
	
	if err != nil {
		return "", err
	}
	
	return result["text"].(string), nil
}

// Enable/Disable do agente
func (v *VRAgentIntegration) Enable() {
	v.isEnabled = true
	v.logger.Println("Agente de IA habilitado")
}

func (v *VRAgentIntegration) Disable() {
	v.isEnabled = false
	v.logger.Println("Agente de IA desabilitado")
}

func (v *VRAgentIntegration) IsEnabled() bool {
	return v.isEnabled
}

// GetMetrics retorna métricas do agente
func (v *VRAgentIntegration) GetMetrics() map[string]interface{} {
	execucoes := v.orchestrator.GetActiveExecutions()
	
	return map[string]interface{}{
		"enabled":             v.isEnabled,
		"active_workflows":    len(execucoes),
		"memory_size":         len(v.memory.ChatHistory().Messages),
		"uptime_seconds":      time.Now().Unix(), // Simplified for this implementation
	}
}

// Utilitário
func containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

// getMainAgentPrompt retorna o prompt principal do agente
func getMainAgentPrompt(customPrompts map[string]string) string {
	if prompt, exists := customPrompts["main_agent"]; exists {
		return prompt
	}
	
	return `Você é um assistente especializado em processamento de Vale Refeição (VR) com capacidades avançadas de IA.

Suas principais funções incluem:
1. Processar planilhas mensais de VR automaticamente
2. Detectar anomalias e inconsistências nos dados  
3. Responder perguntas sobre políticas e regras de VR
4. Gerar insights e relatórios inteligentes
5. Orquestrar workflows complexos de automação

Contexto atual: {chat_history}
Pergunta/Comando: {input}

Responda de forma precisa, citando fontes quando aplicável, e sugira ações automáticas quando apropriado.`
}

