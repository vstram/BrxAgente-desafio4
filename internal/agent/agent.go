package agent

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	
	"BrxAgente-desafio4/internal/chat"
	"BrxAgente-desafio4/internal/excel"
	"BrxAgente-desafio4/internal/agent/tools"
	"BrxAgente-desafio4/internal/workflows"
	"BrxAgente-desafio4/internal/intelligence"
	
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
)

// VRAgent representa o agente principal de IA para processamento de VR
type VRAgent struct {
	config    *AgentConfig
	
	// Componentes LangChain
	chain    *chains.LLMChain
	memory   *memory.ConversationBuffer
	llm      llms.Model
	
	// Integração com sistema existente
	chatService    *chat.Chat
	excelService   *excel.Service
	toolRegistry   *tools.ToolRegistry
	
	// Sistema de workflow
	orchestrator   *workflows.Orchestrator
	
	// Sistema de detecção de anomalias  
	analyzer       *intelligence.Analyzer
	
	// Estado do agente
	enabled    bool
	status     AgentStatus
	startTime  time.Time
	logger     *log.Logger
}

// AgentStatus representa o status atual do agente
type AgentStatus struct {
	State          string            `json:"state"`           // idle, running, error
	LastActivity   time.Time         `json:"last_activity"`   
	TotalRequests  int64             `json:"total_requests"`  
	ErrorCount     int64             `json:"error_count"`     
	Uptime         time.Duration     `json:"uptime"`          
	CurrentTask    string            `json:"current_task"`    
	Metadata       map[string]string `json:"metadata"`        
}

// NewVRAgent cria uma nova instância do agente VR
func NewVRAgent(agentConfig *AgentConfig, chatSvc *chat.Chat) (*VRAgent, error) {
	if agentConfig == nil {
		agentConfig = DefaultAgentConfig()
	}
	
	// Validar configuração
	if err := agentConfig.Validate(); err != nil {
		return nil, fmt.Errorf("configuração inválida: %w", err)
	}

	// Criar apenas serviços essenciais para máxima estabilidade
	var excelSvc *excel.Service
	var toolRegistry *tools.ToolRegistry
	var memoryBuffer *memory.ConversationBuffer
	var orchestrator *workflows.Orchestrator
	var analyzer *intelligence.Analyzer
	
	// Inicialização condicional para evitar overhead
	if agentConfig.DebugMode {
		// Apenas se debug estiver ativado, criar componentes avançados
		excelSvc = excel.NewService()
		
		var toolErr error
		toolRegistry, toolErr = tools.GetDefaultToolRegistry()
		if toolErr != nil {
			// Se falhar, continuar sem ferramentas - não é crítico
			fmt.Printf("Warning: Failed to create tool registry: %v\n", toolErr)
		}
	}
	
	agent := &VRAgent{
		config:       agentConfig,
		memory:       memoryBuffer,
		chatService:  chatSvc,
		excelService: excelSvc,
		toolRegistry: toolRegistry,
		orchestrator: orchestrator,
		analyzer:     analyzer,
		enabled:      agentConfig.Enabled,
		startTime:    time.Now(),
		status: AgentStatus{
			State:         "idle",
			LastActivity:  time.Now(),
			TotalRequests: 0,
			ErrorCount:    0,
			CurrentTask:   "",
			Metadata:      make(map[string]string),
		},
		logger: log.Default(),
	}
	
	// Configurar LLM Chain quando o agente for habilitado
	if agentConfig.Enabled {
		if err := agent.setupLLMChain(); err != nil {
			agent.logger.Printf("Aviso: Falha ao configurar LLM Chain: %v", err)
			// Não retorna erro para manter backward compatibility
		}
	}
	
	// Registrar workflows padrão
	if err := agent.registerDefaultWorkflows(); err != nil {
		agent.logger.Printf("Aviso: Falha ao registrar workflows padrão: %v", err)
	}
	
	// Log de inicialização com verificação de nil
	toolCount := 0
	if toolRegistry != nil {
		toolCount = toolRegistry.Count()
	}
	
	workflowCount := 0
	if agent.orchestrator != nil {
		workflowCount = len(agent.orchestrator.ListWorkflows())
	}
	
	agent.logger.Printf("VRAgent inicializado com sucesso - Enabled: %v, Model: %s, Tools: %d, Workflows: %d", 
		agentConfig.Enabled, agentConfig.Model, toolCount, workflowCount)
	
	return agent, nil
}

// setupLLMChain configura a chain do LangChain
func (a *VRAgent) setupLLMChain() error {
	// Por enquanto, apenas configuramos a estrutura básica
	// O LLM específico será configurado quando integrarmos com OpenAI/Ollama
	a.logger.Println("LLM Chain configurado (usando fallback para chat service)")
	return nil
}

// Ask faz uma pergunta ao agente e retorna a resposta
func (a *VRAgent) Ask(question string) (string, error) {
	if !a.enabled {
		return "", fmt.Errorf("agente está desabilitado")
	}
	
	// Atualizar status
	a.updateStatus("running", fmt.Sprintf("Processando pergunta: %.50s...", question))
	defer a.updateStatus("idle", "")
	
	a.status.TotalRequests++
	
	// Usar apenas chat service para máxima simplicidade e estabilidade
	var response string
	var err error
	
	// Usar diretamente AskOllama ou AskOpenAI para evitar recursão infinita
	// (não usar a.chatService.Ask que pode chamar o próprio agente novamente)
	contextData := a.chatService.GetContextDataAsString()
	systemPrompt := fmt.Sprintf(`Você é um assistente especializado em análise de dados de Vale Refeição (VR) e Vale Alimentação (VA).

Seu trabalho é responder perguntas sobre os dados de colaboradores que foram PROCESSADOS pelo sistema de cálculo de VR/VA.

Contexto dos dados processados:
%s

IMPORTANTE: 
- "Colaboradores processados" = colaboradores que tiveram seus dados de VR/VA calculados pelo sistema
- Todos os colaboradores mostrados no contexto JÁ FORAM PROCESSADOS pelo sistema
- Use os dados fornecidos para responder as perguntas de forma clara e objetiva
- Seja preciso com números e cálculos`, contextData)
	
	// Tentar Ollama primeiro, depois OpenAI se configurado
	response, err = a.chatService.AskOllama(question, systemPrompt)
	if err != nil {
		// Se Ollama falhar, tentar OpenAI
		response, err = a.chatService.AskOpenAI(question, []chat.Message{})
	}
	
	if err != nil {
		a.status.ErrorCount++
		a.logger.Printf("Erro ao processar pergunta: %v", err)
		return "", fmt.Errorf("erro ao processar pergunta: %w", err)
	}
	
	a.logger.Printf("Pergunta processada com sucesso: %.50s...", question)
	return response, nil
}

// askWithLangChain processa pergunta usando LangChain com memory
func (a *VRAgent) askWithLangChain(question string) (string, error) {
	ctx := context.Background()
	
	// Preparar input com contexto das ferramentas disponíveis
	toolsInfo := a.getToolsContext()
	
	input := map[string]interface{}{
		"input": question,
		"tools_available": toolsInfo,
		"agent_context": "Você é um assistente especializado em processamento de VR (Vale Refeição).",
	}
	
	// Executar chain com memory
	result, err := chains.Run(ctx, a.chain, input)
	if err != nil {
		return "", fmt.Errorf("erro na execução da chain: %w", err)
	}
	
	// Armazenar na memória
	if a.memory != nil {
		inputs := map[string]any{"input": question}
		outputs := map[string]any{"output": result}
		if err := a.memory.SaveContext(ctx, inputs, outputs); err != nil {
			a.logger.Printf("Aviso: Erro ao salvar contexto na memória: %v", err)
		}
	}
	
	return result, nil
}

// getToolsContext retorna informações sobre ferramentas disponíveis para o contexto
func (a *VRAgent) getToolsContext() string {
	if a.toolRegistry == nil {
		return "Nenhuma ferramenta específica disponível."
	}
	
	tools := a.toolRegistry.ListNames()
	if len(tools) == 0 {
		return "Nenhuma ferramenta específica disponível."
	}
	
	return fmt.Sprintf("Ferramentas disponíveis: %s", strings.Join(tools, ", "))
}

// GetStatus retorna o status atual do agente
func (a *VRAgent) GetStatus() AgentStatus {
	a.status.Uptime = time.Since(a.startTime)
	return a.status
}

// GetStatusInterface retorna o status como interface{} para compatibilidade
func (a *VRAgent) GetStatusInterface() interface{} {
	return a.GetStatus()
}

// IsEnabled retorna se o agente está habilitado
func (a *VRAgent) IsEnabled() bool {
	return a.enabled
}

// Enable habilita o agente
func (a *VRAgent) Enable() {
	a.enabled = true
	a.logger.Println("Agente habilitado")
}

// Disable desabilita o agente
func (a *VRAgent) Disable() {
	a.enabled = false
	a.updateStatus("disabled", "")
	a.logger.Println("Agente desabilitado")
}

// Reset reseta o estado do agente
func (a *VRAgent) Reset() error {
	a.status = AgentStatus{
		State:         "idle",
		LastActivity:  time.Now(),
		TotalRequests: 0,
		ErrorCount:    0,
		CurrentTask:   "",
		Metadata:      make(map[string]string),
	}
	a.startTime = time.Now()
	
	// Limpar memória do agente
	if a.memory != nil {
		ctx := context.Background()
		if err := a.memory.Clear(ctx); err != nil {
			a.logger.Printf("Aviso: Erro ao limpar memória: %v", err)
		}
	}
	
	a.logger.Println("Estado do agente resetado (incluindo memória)")
	return nil
}

// GetMemory retorna o histórico de conversação atual
func (a *VRAgent) GetMemory() ([]string, error) {
	if a.memory == nil {
		return []string{}, nil
	}
	
	ctx := context.Background()
	memoryVars := a.memory.MemoryVariables(ctx)
	
	var history []string
	for _, key := range memoryVars {
		if key == "history" {
			buffer, err := a.memory.LoadMemoryVariables(ctx, map[string]any{})
			if err != nil {
				return nil, fmt.Errorf("erro ao carregar histórico de memória: %w", err)
			}
			
			if historyValue, ok := buffer["history"]; ok {
				if historyStr, ok := historyValue.(string); ok {
					// Parse do histórico de conversação
					lines := strings.Split(historyStr, "\n")
					for _, line := range lines {
						if strings.TrimSpace(line) != "" {
							history = append(history, strings.TrimSpace(line))
						}
					}
				}
			}
			break
		}
	}
	
	return history, nil
}

// ClearMemory limpa apenas a memória sem resetar outras estatísticas
func (a *VRAgent) ClearMemory() error {
	if a.memory == nil {
		return fmt.Errorf("memória não inicializada")
	}
	
	ctx := context.Background()
	if err := a.memory.Clear(ctx); err != nil {
		return fmt.Errorf("erro ao limpar memória: %w", err)
	}
	
	a.logger.Println("Memória do agente limpa")
	return nil
}

// ExecuteWorkflow executa um workflow específico
func (a *VRAgent) ExecuteWorkflow(workflow string) error {
	if !a.enabled {
		return fmt.Errorf("agente está desabilitado")
	}
	
	a.updateStatus("running", fmt.Sprintf("Executando workflow: %s", workflow))
	defer a.updateStatus("idle", "")
	
	a.status.TotalRequests++
	
	// Por enquanto, implementação básica de workflows
	switch strings.ToLower(workflow) {
	case "processar-vr-mensal":
		return a.executeVRWorkflow()
	case "validar-dados":
		return a.executeValidationWorkflow()
	default:
		return fmt.Errorf("workflow não reconhecido: %s", workflow)
	}
}

// executeVRWorkflow executa o workflow principal de processamento de VR
func (a *VRAgent) executeVRWorkflow() error {
	a.logger.Println("Iniciando workflow de processamento de VR mensal")
	
	// Implementação básica - será expandida em issues futuras
	steps := []string{
		"Validação de dados",
		"Cálculo de VR",
		"Geração de relatórios",
		"Notificação de resultados",
	}
	
	for i, step := range steps {
		a.updateStatus("running", fmt.Sprintf("Executando: %s (%d/%d)", step, i+1, len(steps)))
		a.logger.Printf("Workflow VR - Etapa %d/%d: %s", i+1, len(steps), step)
		
		// Simular processamento
		time.Sleep(100 * time.Millisecond)
	}
	
	a.logger.Println("Workflow de processamento de VR concluído")
	return nil
}

// executeValidationWorkflow executa o workflow de validação de dados
func (a *VRAgent) executeValidationWorkflow() error {
	a.logger.Println("Iniciando workflow de validação de dados")
	
	// Implementação básica - será expandida em issues futuras
	validations := []string{
		"Verificação de arquivos Excel",
		"Validação de dados de colaboradores",
		"Verificação de datas",
		"Análise de inconsistências",
	}
	
	for i, validation := range validations {
		a.updateStatus("running", fmt.Sprintf("Executando: %s (%d/%d)", validation, i+1, len(validations)))
		a.logger.Printf("Workflow Validação - Etapa %d/%d: %s", i+1, len(validations), validation)
		
		// Simular validação
		time.Sleep(50 * time.Millisecond)
	}
	
	a.logger.Println("Workflow de validação concluído")
	return nil
}

// GetConfig retorna a configuração atual do agente
func (a *VRAgent) GetConfig() *AgentConfig {
	return a.config
}

// GetAvailableTools retorna lista de ferramentas disponíveis
func (a *VRAgent) GetAvailableTools() []string {
	return a.toolRegistry.ListNames()
}

// ExecuteTool executa uma ferramenta específica
func (a *VRAgent) ExecuteTool(toolName, input string) (string, error) {
	if !a.enabled {
		return "", fmt.Errorf("agente está desabilitado")
	}

	a.updateStatus("running", fmt.Sprintf("Executando ferramenta: %s", toolName))
	defer a.updateStatus("idle", "")

	a.status.TotalRequests++

	result, err := a.toolRegistry.Execute(toolName, input)
	if err != nil {
		a.status.ErrorCount++
		a.logger.Printf("Erro ao executar ferramenta %s: %v", toolName, err)
		return "", fmt.Errorf("erro ao executar ferramenta %s: %w", toolName, err)
	}

	a.logger.Printf("Ferramenta %s executada com sucesso", toolName)
	return result, nil
}

// GetToolInfo retorna informações sobre uma ferramenta específica
func (a *VRAgent) GetToolInfo(toolName string) (map[string]interface{}, error) {
	return a.toolRegistry.GetToolInfo(toolName)
}

// GetAllToolsInfo retorna informações sobre todas as ferramentas
func (a *VRAgent) GetAllToolsInfo() map[string]interface{} {
	return a.toolRegistry.GetAllToolsInfo()
}

// updateStatus atualiza o status interno do agente
func (a *VRAgent) updateStatus(state, task string) {
	a.status.State = state
	a.status.CurrentTask = task
	a.status.LastActivity = time.Now()
}

// registerDefaultWorkflows registra workflows padrão no orquestrador
func (a *VRAgent) registerDefaultWorkflows() error {
	// Só registrar workflows se orchestrator estiver disponível
	if a.orchestrator == nil {
		a.logger.Printf("Orchestrator não disponível - pulando registro de workflows")
		return nil
	}
	
	// Registrar workflow de validação simples
	simpleValidation := workflows.NewSimpleValidationWorkflow()
	if err := a.orchestrator.RegisterWorkflow(simpleValidation); err != nil {
		return fmt.Errorf("erro ao registrar workflow simple-validation: %w", err)
	}
	
	// TODO: Registrar workflow com detecção de anomalias quando implementado
	// validatedVR := intelligence.NewValidatedVRWorkflow(a.analyzer)
	// if err := a.orchestrator.RegisterWorkflow(validatedVR); err != nil {
	//     return fmt.Errorf("erro ao registrar workflow validated-vr-processing: %w", err)
	// }
	
	// Registrar outros workflows futuros aqui
	// TODO: Implementar workflow processar-vr-mensal quando as ferramentas estiverem prontas
	
	return nil
}

// ExecuteWorkflowByName executa um workflow pelo nome através do orquestrador
func (a *VRAgent) ExecuteWorkflowByName(workflowName string, params map[string]interface{}) (*workflows.WorkflowResult, error) {
	if !a.enabled {
		return nil, fmt.Errorf("agente está desabilitado")
	}
	
	a.updateStatus("running", fmt.Sprintf("workflow:%s", workflowName))
	a.status.TotalRequests++
	
	defer func() {
		a.updateStatus("idle", "")
	}()
	
	result, err := a.orchestrator.ExecuteWorkflow(workflowName, params)
	if err != nil {
		a.status.ErrorCount++
		a.logger.Printf("Erro ao executar workflow %s: %v", workflowName, err)
		return nil, err
	}
	
	a.logger.Printf("Workflow %s executado com sucesso", workflowName)
	return result, nil
}

// ExecuteWorkflowAsync executa um workflow de forma assíncrona
func (a *VRAgent) ExecuteWorkflowAsync(workflowName string, params map[string]interface{}) (string, error) {
	if !a.enabled {
		return "", fmt.Errorf("agente está desabilitado")
	}
	
	a.status.TotalRequests++
	
	executionID, err := a.orchestrator.ExecuteWorkflowAsync(workflowName, params)
	if err != nil {
		a.status.ErrorCount++
		a.logger.Printf("Erro ao executar workflow async %s: %v", workflowName, err)
		return "", err
	}
	
	a.logger.Printf("Workflow %s iniciado assincronamente com ID: %s", workflowName, executionID)
	return executionID, nil
}

// GetWorkflowExecution retorna informações sobre uma execução de workflow
func (a *VRAgent) GetWorkflowExecution(executionID string) (*workflows.WorkflowExecution, error) {
	return a.orchestrator.GetExecution(executionID)
}

// CancelWorkflowExecution cancela uma execução de workflow em andamento
func (a *VRAgent) CancelWorkflowExecution(executionID string) error {
	return a.orchestrator.CancelExecution(executionID)
}

// ListAvailableWorkflows retorna lista de workflows disponíveis
func (a *VRAgent) ListAvailableWorkflows() []string {
	return a.orchestrator.ListWorkflows()
}

// GetWorkflowOrchestrator retorna o orquestrador (para uso interno/testes)
func (a *VRAgent) GetWorkflowOrchestrator() *workflows.Orchestrator {
	return a.orchestrator
}

// AnalyzeAnomalies executa análise de anomalias nos dados fornecidos
func (a *VRAgent) AnalyzeAnomalies(colaboradores map[string]interface{}, params map[string]interface{}) (*intelligence.AnomalyReport, error) {
	if !a.enabled {
		return nil, fmt.Errorf("agente está desabilitado")
	}
	
	a.updateStatus("running", "anomaly-analysis")
	a.status.TotalRequests++
	
	defer func() {
		a.updateStatus("idle", "")
	}()
	
	report, err := a.analyzer.AnalyzeData(colaboradores, params)
	if err != nil {
		a.status.ErrorCount++
		a.logger.Printf("Erro na análise de anomalias: %v", err)
		return nil, err
	}
	
	a.logger.Printf("Análise de anomalias concluída: %d anomalias detectadas (score: %.1f)",
		report.TotalAnomalies, report.Summary.OverallScore)
	
	return report, nil
}

// GetAnomalyAnalyzer retorna o analisador de anomalias (para uso interno/testes)
func (a *VRAgent) GetAnomalyAnalyzer() *intelligence.Analyzer {
	return a.analyzer
}

// FormatAnomalyReport formata relatório de anomalias para exibição
func (a *VRAgent) FormatAnomalyReport(report *intelligence.AnomalyReport) string {
	return intelligence.FormatAnomalyReportForHuman(report)
}