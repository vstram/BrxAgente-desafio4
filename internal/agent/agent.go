package agent

import (
	"fmt"
	"log"
	"time"
	
	"BrxAgente-desafio4/internal/chat"
)

// VRAgent representa o agente principal de IA para processamento de VR
type VRAgent struct {
	config    *AgentConfig
	
	// Integração com sistema existente
	chatService    *chat.Chat
	// excelService   será implementado na próxima issue
	// calculoService será implementado na próxima issue
	
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
	
	agent := &VRAgent{
		config:      agentConfig,
		chatService: chatSvc,
		enabled:     agentConfig.Enabled,
		startTime:   time.Now(),
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
	
	// Log de inicialização
	agent.logger.Printf("VRAgent inicializado com sucesso - Enabled: %v, Model: %s", 
		agentConfig.Enabled, agentConfig.Model)
	
	return agent, nil
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
	
	// Por enquanto, delegar para o chat service existente
	// Nas próximas issues, implementaremos LangChain aqui
	response, err := a.chatService.Ask(question, "", []chat.Message{})
	if err != nil {
		a.status.ErrorCount++
		a.logger.Printf("Erro ao processar pergunta: %v", err)
		return "", fmt.Errorf("erro ao processar pergunta: %w", err)
	}
	
	a.logger.Printf("Pergunta processada com sucesso: %.50s...", question)
	return response, nil
}

// GetStatus retorna o status atual do agente
func (a *VRAgent) GetStatus() AgentStatus {
	a.status.Uptime = time.Since(a.startTime)
	return a.status
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
	a.logger.Println("Estado do agente resetado")
	return nil
}

// GetConfig retorna a configuração atual do agente
func (a *VRAgent) GetConfig() *AgentConfig {
	return a.config
}

// updateStatus atualiza o status interno do agente
func (a *VRAgent) updateStatus(state, task string) {
	a.status.State = state
	a.status.CurrentTask = task
	a.status.LastActivity = time.Now()
}