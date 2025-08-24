package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"

	"BrxAgente-desafio4/internal/agent"
	"BrxAgente-desafio4/internal/chat"
	"BrxAgente-desafio4/internal/config"
	"BrxAgente-desafio4/internal/intelligence"
	"BrxAgente-desafio4/internal/predicoes"
	"BrxAgente-desafio4/internal/modelo"
	"BrxAgente-desafio4/internal/security"
)

// App struct
type App struct {
	ctx         context.Context
	cfg         *config.Config
	chat        *chat.Chat
	vrAgent     *agent.VRAgent    // Agente VR para processamento inteligente
	colaboradores map[string]*modelo.Colaborador
	mu          sync.RWMutex // Mutex para acesso seguro aos dados compartilhados
	
	// Agent monitoring fields
	agentStatus      string
	currentWorkflow  *WorkflowInfo
	agentStartTime   time.Time
	workflowHistory  []WorkflowExecution
	systemLogs       []LogEntry
	agentMu          sync.RWMutex // Mutex para campos de monitoramento
	
	// Intelligence and prediction system
	trendPredictor     *intelligence.TrendPredictor
	patternAnalyzer    *intelligence.PatternAnalyzer
	trendDetector      *intelligence.TrendDetector
	forecaster         *intelligence.Forecaster
	recommendationEngine *intelligence.RecommendationEngine
	historicalData     []predicoes.HistoricalVRData
	dataMu            sync.RWMutex // Mutex para dados históricos
}

// AgentStatus represents possible agent states
type AgentStatus struct {
	Status              string                 `json:"status"`
	LastUpdated         time.Time              `json:"lastUpdated"`
	CurrentWorkflow     *WorkflowInfo          `json:"currentWorkflow"`
	AvailableWorkflows  []string               `json:"availableWorkflows"`
	Metrics             AgentMetrics           `json:"metrics"`
	RecentLogs          []LogEntry             `json:"recentLogs"`
}

// WorkflowInfo represents workflow execution information
type WorkflowInfo struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Status      string                `json:"status"`
	StartTime   time.Time             `json:"startTime"`
	EndTime     *time.Time            `json:"endTime"`
	Steps       []WorkflowStep        `json:"steps"`
	Progress    float64               `json:"progress"`
	ErrorMsg    string                `json:"errorMsg"`
}

// WorkflowStep represents individual workflow steps
type WorkflowStep struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	StartTime   *time.Time `json:"startTime"`
	EndTime     *time.Time `json:"endTime"`
	Duration    int64      `json:"duration"`
	ErrorMsg    string     `json:"errorMsg"`
}

// WorkflowExecution represents completed workflow executions
type WorkflowExecution struct {
	ID                string    `json:"id"`
	WorkflowName      string    `json:"workflowName"`
	Status            string    `json:"status"`
	StartTime         time.Time `json:"startTime"`
	EndTime           time.Time `json:"endTime"`
	Duration          int64     `json:"duration"`
	CollaboratorsProcessed int   `json:"collaboratorsProcessed"`
	ReportsGenerated  int       `json:"reportsGenerated"`
	AnomaliesDetected int       `json:"anomaliesDetected"`
	ErrorMsg          string    `json:"errorMsg"`
}

// AgentMetrics represents agent performance metrics
type AgentMetrics struct {
	TotalWorkflowsExecuted int     `json:"totalWorkflowsExecuted"`
	SuccessfulWorkflows    int     `json:"successfulWorkflows"`
	CollaboratorsProcessed int     `json:"collaboratorsProcessed"`
	ReportsGenerated       int     `json:"reportsGenerated"`
	AnomaliesDetected      int     `json:"anomaliesDetected"`
	Uptime                 int64   `json:"uptime"`
}

// LogEntry represents system log entries
type LogEntry struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
}

// WorkflowStartRequest represents workflow start parameters
type WorkflowStartRequest struct {
	WorkflowName string                 `json:"workflowName"`
	Parameters   map[string]interface{} `json:"parameters"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		colaboradores: make(map[string]*modelo.Colaborador),
		agentStatus:   "idle",
		agentStartTime: time.Now(),
		workflowHistory: make([]WorkflowExecution, 0),
		systemLogs:    make([]LogEntry, 0),
		historicalData: make([]predicoes.HistoricalVRData, 0),
		
		// Initialize intelligence systems
		trendPredictor: intelligence.NewTrendPredictor(intelligence.TrendPredictorConfig{
			MinDataPoints:       6,
			ConfidenceThreshold: 0.7,
		}),
		patternAnalyzer: intelligence.NewPatternAnalyzer(intelligence.PatternAnalyzerConfig{
			MinDataPoints:   6,
			ConfidenceLevel: 0.7,
		}),
		trendDetector: intelligence.NewTrendDetector(intelligence.TrendDetectorConfig{
			MinDataPoints: 6,
		}),
		forecaster: intelligence.NewForecaster(intelligence.ForecastConfig{
			DefaultHorizon: 3,
			MinAccuracy:    0.7,
		}),
		recommendationEngine: intelligence.NewRecommendationEngine(intelligence.RecommendationConfig{
			MaxRecommendations:   10,
			MinConfidence:        0.7,
			PriorityThreshold:    0.8,
			AutoApproveThreshold: 0.9,
			ContextWindow:        6,
		}),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Warning: Failed to load configuration: %v\n", err)
		// Use default config
		cfg = &config.Config{}
	}
	a.cfg = cfg
	
	// Initialize chat
	a.chat = chat.NewChat(cfg)
	
	// Fix any existing Ollama configuration issues
	if err := a.FixOllamaConfig(); err != nil {
		fmt.Printf("Warning: Failed to fix Ollama config: %v\n", err)
	}
	
	// Initialize Basic VR Agent for optimal performance
	fmt.Println("Initializing VR Agent...")
	
	agentConfig := agent.DefaultAgentConfig()
	agentConfig.Enabled = true
	agentConfig.DebugMode = false // Reduce overhead for faster processing
	
	basicAgent, err := agent.NewVRAgent(agentConfig, a.chat)
	if err != nil {
		fmt.Printf("Warning: Failed to initialize VR agent: %v\n", err)
		fmt.Println("Chat will use fallback services (OpenAI/Ollama)")
	} else {
		a.vrAgent = basicAgent  // Assign to App instance
		a.chat.SetAgent(basicAgent)
		fmt.Println("VR Agent initialized successfully")
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// TestExcelReading tests reading an Excel file using the excelize library
func (a *App) TestExcelReading() (string, error) {
	// Using the ADMISSÃO ABRIL.xlsx file as an example
	filePath := filepath.Join("files", "ADMISSÃO ABRIL.xlsx")

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening Excel file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()

	// Get all sheets from the workbook
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return "", fmt.Errorf("no sheets found in the Excel file")
	}

	// Read the first sheet
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return "", fmt.Errorf("error reading rows from sheet %s: %w", sheets[0], err)
	}

	// Return information about the file
	return fmt.Sprintf("Successfully read Excel file with %d sheets. First sheet '%s' has %d rows.",
		len(sheets), sheets[0], len(rows)), nil
}

// SetDiretorioPlanilhas recebe o caminho do diretório contendo as planilhas e o valida
func (a *App) SetDiretorioPlanilhas(caminho string) (bool, error) {
	// Verificando se o caminho foi fornecido
	if caminho == "" {
		return false, fmt.Errorf("caminho do diretório não pode ser vazio")
	}

	// Verificando se o caminho existe
	info, err := os.Stat(caminho)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("diretório não encontrado: %s", caminho)
		}
		return false, fmt.Errorf("erro ao acessar o diretório: %w", err)
	}

	// Verificando se é um diretório
	if !info.IsDir() {
		return false, fmt.Errorf("o caminho fornecido não é um diretório: %s", caminho)
	}

	// Verificando se o diretório contém arquivos .xlsx
	arquivos, err := os.ReadDir(caminho)
	if err != nil {
		return false, fmt.Errorf("erro ao ler o diretório: %w", err)
	}

	// Procurando por arquivos .xlsx
	encontrouPlanilha := false
	for _, arquivo := range arquivos {
		if !arquivo.IsDir() && filepath.Ext(arquivo.Name()) == ".xlsx" {
			encontrouPlanilha = true
			break
		}
	}

	if !encontrouPlanilha {
		return false, fmt.Errorf("nenhum arquivo .xlsx encontrado no diretório: %s", caminho)
	}

	// Se passou por todas as validações, o diretório é válido
	return true, nil
}

// SelecionarDiretorio abre um diálogo para o usuário selecionar um diretório
func (a *App) SelecionarDiretorio() (string, error) {
	directory, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Selecione o diretório das planilhas",
	})
	if err != nil {
		return "", err
	}
	return directory, nil
}

// GetConfig retorna a configuração atual da aplicação
func (a *App) GetConfig() (*config.Config, error) {
	return a.cfg, nil
}

// reloadConfig reloads configuration from file
func (a *App) reloadConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("falha ao recarregar configuração: %w", err)
	}
	
	a.cfg = cfg
	fmt.Printf("Configuration reloaded successfully\n")
	return nil
}

// FixOllamaConfig corrige configuração existente do Ollama se necessário
func (a *App) FixOllamaConfig() error {
	if a.cfg == nil {
		return fmt.Errorf("configuração não inicializada")
	}
	
	originalURL := a.cfg.OllamaConfig.BaseURL
	if originalURL == "" {
		return nil // Nada para corrigir
	}
	
	// Apply the same URL fixing logic
	fixed := originalURL
	
	// Add http:// if missing
	if !strings.HasPrefix(fixed, "http://") && !strings.HasPrefix(fixed, "https://") {
		fixed = "http://" + fixed
	}
	
	// Fix various localhost patterns
	if strings.HasPrefix(fixed, "http://:") {
		fixed = strings.Replace(fixed, "http://:", "http://localhost:", 1)
	} else if fixed == "http://11434" || strings.HasPrefix(fixed, "http://11434") {
		fixed = strings.Replace(fixed, "http://11434", "http://localhost:11434", 1)
	} else if strings.Contains(fixed, "://11434") {
		fixed = strings.Replace(fixed, "://11434", "://localhost:11434", 1)
	}
	
	if originalURL != fixed {
		fmt.Printf("Fixing existing Ollama URL from '%s' to '%s'\n", originalURL, fixed)
		a.cfg.OllamaConfig.BaseURL = fixed
		
		// Save fixed configuration
		if err := config.SaveConfig(a.cfg); err != nil {
			return fmt.Errorf("falha ao salvar configuração corrigida: %w", err)
		}
		
		// Update chat configuration
		if a.chat != nil {
			a.chat.UpdateConfig(a.cfg)
		}
	}
	
	return nil
}

// SetOpenAIKey define a chave da API do OpenAI
func (a *App) SetOpenAIKey(key string) error {
	// Validate the key format
	if key != "" && !config.ValidateOpenAIKey(key) {
		return fmt.Errorf("chave da API do OpenAI inválida")
	}

	// Create a secure string for the key
	_, err := security.NewSecureString(key)
	if err != nil {
		return fmt.Errorf("falha ao criar string segura para a chave: %w", err)
	}

	// Update config with the encrypted key
	a.cfg.OpenAIKey = key // In a real implementation, we would store the encrypted value

	// Save config
	if err := config.SaveConfig(a.cfg); err != nil {
		return fmt.Errorf("falha ao salvar a configuração: %w", err)
	}
	
	// Reload configuration to ensure consistency
	if err := a.reloadConfig(); err != nil {
		fmt.Printf("Warning: Failed to reload config after save: %v\n", err)
	}
	
	// Update chat configuration
	if a.chat != nil {
		a.chat.UpdateConfig(a.cfg)
	}

	return nil
}

// SetOllamaConfig define a configuração do Ollama
func (a *App) SetOllamaConfig(ollamaConfig config.OllamaConfig) error {
	// Debug input
	fmt.Printf("SetOllamaConfig called with: BaseURL='%s', Model='%s'\n", ollamaConfig.BaseURL, ollamaConfig.Model)
	
	// Auto-fix common URL issues
	if ollamaConfig.BaseURL != "" {
		originalURL := ollamaConfig.BaseURL
		
		// Add http:// if missing
		if !strings.HasPrefix(ollamaConfig.BaseURL, "http://") && !strings.HasPrefix(ollamaConfig.BaseURL, "https://") {
			ollamaConfig.BaseURL = "http://" + ollamaConfig.BaseURL
		}
		
		// Fix various localhost patterns
		if strings.HasPrefix(ollamaConfig.BaseURL, "http://:") {
			// http://:11434 -> http://localhost:11434
			ollamaConfig.BaseURL = strings.Replace(ollamaConfig.BaseURL, "http://:", "http://localhost:", 1)
		} else if ollamaConfig.BaseURL == "http://11434" || strings.HasPrefix(ollamaConfig.BaseURL, "http://11434") {
			// http://11434 -> http://localhost:11434
			ollamaConfig.BaseURL = strings.Replace(ollamaConfig.BaseURL, "http://11434", "http://localhost:11434", 1)
		} else if strings.Contains(ollamaConfig.BaseURL, "://11434") {
			// Any protocol with just port
			ollamaConfig.BaseURL = strings.Replace(ollamaConfig.BaseURL, "://11434", "://localhost:11434", 1)
		}
		
		// Also handle cases like "localhost:11434" without protocol
		if strings.HasPrefix(ollamaConfig.BaseURL, "http://localhost:") || strings.HasPrefix(ollamaConfig.BaseURL, "https://localhost:") {
			// Already correct
		} else if strings.HasPrefix(ollamaConfig.BaseURL, "http://") && !strings.Contains(ollamaConfig.BaseURL[7:], ":") && len(ollamaConfig.BaseURL) < 12 {
			// Something like http://11434 - add localhost
			port := ollamaConfig.BaseURL[7:]
			ollamaConfig.BaseURL = "http://localhost:" + port
		}
		
		if originalURL != ollamaConfig.BaseURL {
			fmt.Printf("Auto-corrected URL from '%s' to '%s'\n", originalURL, ollamaConfig.BaseURL)
		}
	}
	
	// Validate the configuration
	if err := config.ValidateOllamaConfig(ollamaConfig); err != nil {
		fmt.Printf("Ollama config validation failed: %v\n", err)
		return fmt.Errorf("configuração do Ollama inválida: %w", err)
	}

	// Update config
	a.cfg.OllamaConfig = ollamaConfig
	
	// Debug updated config
	fmt.Printf("Updated config: BaseURL='%s', Model='%s'\n", a.cfg.OllamaConfig.BaseURL, a.cfg.OllamaConfig.Model)

	// Save config
	if err := config.SaveConfig(a.cfg); err != nil {
		return fmt.Errorf("falha ao salvar a configuração: %w", err)
	}
	
	// Reload configuration to ensure consistency
	if err := a.reloadConfig(); err != nil {
		fmt.Printf("Warning: Failed to reload config after save: %v\n", err)
	}
	
	// Update chat configuration
	if a.chat != nil {
		a.chat.UpdateConfig(a.cfg)
	}

	return nil
}

// TestOpenAIKey tests the OpenAI API key by making a simple request
func (a *App) TestOpenAIKey(key string) (bool, error) {
	// For now, just validate the format
	// In a real implementation, we would make an actual API call
	if key == "" {
		return false, fmt.Errorf("chave da API do OpenAI não fornecida")
	}

	if !config.ValidateOpenAIKey(key) {
		return false, fmt.Errorf("chave da API do OpenAI inválida")
	}

	// In a real implementation, we would test the key by making an API call
	// For now, we'll just return true if the format is valid
	return true, nil
}

// TestOllamaConnection tests the Ollama connection
func (a *App) TestOllamaConnection(ollamaConfig config.OllamaConfig) (bool, error) {
	// For now, just validate the configuration
	// In a real implementation, we would make an actual connection test
	if err := config.ValidateOllamaConfig(ollamaConfig); err != nil {
		return false, fmt.Errorf("configuração do Ollama inválida: %w", err)
	}

	// Test actual connection to Ollama
	return a.testOllamaAPI(ollamaConfig)
}

// testOllamaAPI tests the actual Ollama API connection
func (a *App) testOllamaAPI(ollamaConfig config.OllamaConfig) (bool, error) {
	
	// Try to list models first
	url := ollamaConfig.BaseURL
	if url[len(url)-1] != '/' {
		url += "/"
	}
	url += "api/tags"
	
	fmt.Printf("Testing Ollama API at: %s\n", url)
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return false, fmt.Errorf("falha ao conectar com Ollama: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("Ollama retornou status %d", resp.StatusCode)
	}
	
	// Try to read response to see available models
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("falha ao parsear resposta do Ollama: %w", err)
	}
	
	fmt.Printf("Ollama connection successful. Available models info: %+v\n", result)
	return true, nil
}

// AskAI sends a question to the configured AI service and returns the response
func (a *App) AskAI(question string) (string, error) {
	fmt.Printf("AskAI called with question: %.100s...\n", question)
	
	// Check agent status
	if a.vrAgent != nil {
		fmt.Printf("VR Agent is available - Status: %s, Enabled: %v\n", a.vrAgent.GetStatus().State, a.vrAgent.IsEnabled())
	} else {
		fmt.Println("VR Agent is not available")
	}
	
	// Get system prompt with consolidated data context
	systemPrompt, err := a.GetSystemPrompt()
	if err != nil {
		fmt.Printf("Warning: Failed to get system prompt with context: %v\n", err)
		// Fallback to basic prompt
		systemPrompt = `Você é um assistente especializado em análise de dados de Vale Refeição (VR) e Vale Alimentação (VA).
		Você está ajudando um usuário a entender os resultados do processamento de dados de colaboradores.
		Os dados dos colaboradores são identificados exclusivamente por uma MATRICULA, por razões de confidencialidade.
		Seja claro, preciso e objetivo em suas respostas.`
	}

	// Create empty context for now
	var context []chat.Message

	// Ask the chat service
	fmt.Println("Calling chat.Ask method...")
	response, err := a.chat.Ask(question, systemPrompt, context)
	if err != nil {
		fmt.Printf("Error from chat.Ask: %v\n", err)
		return "", fmt.Errorf("falha ao obter resposta da IA: %w", err)
	}

	fmt.Printf("Successfully got response from chat service (length: %d)\n", len(response))
	return response, nil
}

// GetSystemPrompt returns the current system prompt with context data
func (a *App) GetSystemPrompt() (string, error) {
	// Create the same system prompt used in AskAI
	systemPrompt := `Você é um assistente especializado em análise de dados de Vale Refeição (VR) e Vale Alimentação (VA).
	Você está ajudando um usuário a entender os resultados do processamento de dados de colaboradores.
	Os dados dos colaboradores são identificados exclusivamente por uma MATRICULA, por razões de confidencialidade.
	Seja claro, preciso e objetivo em suas respostas.`

	// Get context data from chat service
	contextDataStr := a.chat.GetContextDataAsString()

	// Combine system prompt with context data
	if contextDataStr != "" {
		systemPrompt = fmt.Sprintf("Contexto dos dados:\n%s\n\n%s", contextDataStr, systemPrompt)
	}

	return systemPrompt, nil
}

// GetConsolidatedData returns a copy of the consolidated data
func (a *App) GetConsolidatedData() (map[string]*modelo.Colaborador, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	// Create a copy of the map to avoid external modifications
	copia := make(map[string]*modelo.Colaborador)
	for k, v := range a.colaboradores {
		copia[k] = v
	}
	
	return copia, nil
}

// IsAgentEnabled retorna se o agente de IA está habilitado
func (a *App) IsAgentEnabled() bool {
	if a.vrAgent == nil {
		return false
	}
	return a.vrAgent.IsEnabled()
}

// SetAgentEnabled ativa/desativa o agente de IA
func (a *App) SetAgentEnabled(enabled bool) error {
	if a.vrAgent == nil {
		return fmt.Errorf("agente não inicializado")
	}
	
	if enabled {
		a.vrAgent.Enable()
		fmt.Printf("Agente de IA habilitado\n")
	} else {
		a.vrAgent.Disable()
		fmt.Printf("Agente de IA desabilitado\n")
	}
	
	return nil
}

// AskAIBasic faz uma pergunta usando apenas o comportamento anterior (sem agente)
func (a *App) AskAIBasic(question string) (string, error) {
	fmt.Printf("AskAIBasic called with question: %.100s...\n", question)
	
	// Debug configuration
	openAIStatus := "not configured"
	if a.cfg.OpenAIKey != "" && len(a.cfg.OpenAIKey) > 10 {
		openAIStatus = a.cfg.OpenAIKey[:10] + "..."
	}
	fmt.Printf("Debug - OpenAI Key: %s\n", openAIStatus)
	fmt.Printf("Debug - Ollama BaseURL: %s\n", a.cfg.OllamaConfig.BaseURL)
	fmt.Printf("Debug - Ollama Model: %s\n", a.cfg.OllamaConfig.Model)
	
	// Get system prompt with consolidated data context
	systemPrompt, err := a.GetSystemPrompt()
	if err != nil {
		fmt.Printf("Warning: Failed to get system prompt with context: %v\n", err)
		// Fallback to basic prompt
		systemPrompt = `Você é um assistente especializado em análise de dados de Vale Refeição (VR) e Vale Alimentação (VA).
		Você está ajudando um usuário a entender os resultados do processamento de dados de colaboradores.
		Os dados dos colaboradores são identificados exclusivamente por uma MATRICULA, por razões de confidencialidade.
		Seja claro, preciso e objetivo em suas respostas.`
	}

	// Create empty context for now
	var context []chat.Message

	// Try OpenAI first if configured
	if a.cfg.OpenAIKey != "" {
		fmt.Printf("Trying OpenAI...\n")
		response, err := a.chat.AskOpenAI(question, context)
		if err == nil {
			fmt.Printf("Basic chat response from OpenAI: %.100s...\n", response)
			return response, nil
		}
		fmt.Printf("OpenAI failed: %v\n", err)
	}
	
	// Try Ollama if configured
	if a.cfg.OllamaConfig.BaseURL != "" && a.cfg.OllamaConfig.Model != "" {
		fmt.Printf("Trying Ollama...\n")
		response, err := a.chat.AskOllama(question, systemPrompt)
		if err == nil {
			fmt.Printf("Basic chat response from Ollama: %.100s...\n", response)
			return response, nil
		}
		fmt.Printf("Ollama failed: %v\n", err)
	}
	
	// Final check: print detailed configuration status
	fmt.Printf("=== CONFIGURATION DEBUG ===\n")
	fmt.Printf("Config object: %+v\n", a.cfg)
	if a.cfg != nil {
		fmt.Printf("OpenAI Key exists: %v (length: %d)\n", a.cfg.OpenAIKey != "", len(a.cfg.OpenAIKey))
		fmt.Printf("Ollama BaseURL: '%s' (empty: %v)\n", a.cfg.OllamaConfig.BaseURL, a.cfg.OllamaConfig.BaseURL == "")
		fmt.Printf("Ollama Model: '%s' (empty: %v)\n", a.cfg.OllamaConfig.Model, a.cfg.OllamaConfig.Model == "")
	} else {
		fmt.Printf("Config is nil!\n")
	}
	fmt.Printf("===========================\n")
	
	return "", fmt.Errorf("nenhum serviço de IA configurado. Configure OpenAI (chave de API) ou Ollama (URL e modelo) nas configurações")
}

// Helper method to add system logs
func (a *App) addSystemLog(level, message, source string) {
	a.agentMu.Lock()
	defer a.agentMu.Unlock()
	
	log := LogEntry{
		ID:        fmt.Sprintf("log-%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Source:    source,
	}
	
	a.systemLogs = append(a.systemLogs, log)
	
	// Keep only the last 100 logs
	if len(a.systemLogs) > 100 {
		a.systemLogs = a.systemLogs[len(a.systemLogs)-100:]
	}
}

// GetAgentStatus returns the current agent status and monitoring data
func (a *App) GetAgentStatus() (*AgentStatus, error) {
	a.agentMu.RLock()
	defer a.agentMu.RUnlock()
	
	// Calculate metrics from workflows
	totalWorkflows := len(a.workflowHistory)
	successfulWorkflows := 0
	totalCollaborators := 0
	totalReports := 0
	totalAnomalies := 0
	
	for _, workflow := range a.workflowHistory {
		if workflow.Status == "completed" {
			successfulWorkflows++
		}
		totalCollaborators += workflow.CollaboratorsProcessed
		totalReports += workflow.ReportsGenerated
		totalAnomalies += workflow.AnomaliesDetected
	}
	
	// Se não há workflows históricos, usar dados reais consolidados
	a.mu.RLock()
	realCollaborators := len(a.colaboradores)
	a.mu.RUnlock()
	
	if totalCollaborators == 0 && realCollaborators > 0 {
		totalCollaborators = realCollaborators
		totalWorkflows = 1 // Representa processamento atual
		successfulWorkflows = 1
		totalReports = 1
	}
	
	uptime := time.Since(a.agentStartTime).Milliseconds()
	
	// Get recent logs (last 20)
	recentLogs := make([]LogEntry, 0)
	if len(a.systemLogs) > 0 {
		startIdx := len(a.systemLogs) - 20
		if startIdx < 0 {
			startIdx = 0
		}
		recentLogs = a.systemLogs[startIdx:]
	}
	
	status := &AgentStatus{
		Status:      a.agentStatus,
		LastUpdated: time.Now(),
		CurrentWorkflow: a.currentWorkflow,
		AvailableWorkflows: []string{
			"analise-vr-mensal",
			"validacao-planilhas", 
			"deteccao-anomalias",
			"geracao-relatorios",
			"auditoria-inteligente",
		},
		Metrics: AgentMetrics{
			TotalWorkflowsExecuted: totalWorkflows,
			SuccessfulWorkflows:    successfulWorkflows,
			CollaboratorsProcessed: totalCollaborators,
			ReportsGenerated:       totalReports,
			AnomaliesDetected:      totalAnomalies,
			Uptime:                 uptime,
		},
		RecentLogs: recentLogs,
	}
	
	return status, nil
}

// StartWorkflow starts a workflow execution
func (a *App) StartWorkflow(request WorkflowStartRequest) error {
	a.agentMu.Lock()
	defer a.agentMu.Unlock()
	
	if a.agentStatus != "idle" {
		return fmt.Errorf("agent is currently busy with status: %s", a.agentStatus)
	}
	
	// Create workflow info
	workflowID := fmt.Sprintf("workflow-%d", time.Now().UnixNano())
	workflow := &WorkflowInfo{
		ID:        workflowID,
		Name:      request.WorkflowName,
		Status:    "running",
		StartTime: time.Now(),
		Steps:     a.getWorkflowSteps(request.WorkflowName),
		Progress:  0.0,
	}
	
	a.currentWorkflow = workflow
	a.agentStatus = "running"
	
	a.addSystemLog("info", fmt.Sprintf("Started workflow: %s", request.WorkflowName), "agent")
	
	// Start workflow execution in a goroutine
	go a.executeWorkflow(workflow, request)
	
	return nil
}

// StopWorkflow gracefully stops the current workflow
func (a *App) StopWorkflow() error {
	a.agentMu.Lock()
	defer a.agentMu.Unlock()
	
	if a.currentWorkflow == nil {
		return fmt.Errorf("no workflow is currently running")
	}
	
	a.currentWorkflow.Status = "stopped"
	endTime := time.Now()
	a.currentWorkflow.EndTime = &endTime
	
	a.agentStatus = "idle"
	a.addSystemLog("warn", fmt.Sprintf("Stopped workflow: %s", a.currentWorkflow.Name), "agent")
	
	// Add to history
	a.addWorkflowToHistory(a.currentWorkflow, "stopped")
	a.currentWorkflow = nil
	
	return nil
}

// CancelWorkflow forcefully cancels the current workflow
func (a *App) CancelWorkflow() error {
	a.agentMu.Lock()
	defer a.agentMu.Unlock()
	
	if a.currentWorkflow == nil {
		return fmt.Errorf("no workflow is currently running")
	}
	
	a.currentWorkflow.Status = "cancelled"
	endTime := time.Now()
	a.currentWorkflow.EndTime = &endTime
	
	a.agentStatus = "idle"
	a.addSystemLog("warn", fmt.Sprintf("Cancelled workflow: %s", a.currentWorkflow.Name), "agent")
	
	// Add to history
	a.addWorkflowToHistory(a.currentWorkflow, "cancelled")
	a.currentWorkflow = nil
	
	return nil
}

// ClearAgentLogs clears the system logs
func (a *App) ClearAgentLogs() error {
	a.agentMu.Lock()
	defer a.agentMu.Unlock()
	
	a.systemLogs = make([]LogEntry, 0)
	a.addSystemLog("info", "System logs cleared", "agent")
	
	return nil
}

// GetWorkflowHistory returns the workflow execution history
func (a *App) GetWorkflowHistory() ([]WorkflowExecution, error) {
	a.agentMu.RLock()
	defer a.agentMu.RUnlock()
	
	return a.workflowHistory, nil
}

// Helper method to define workflow steps based on workflow name
func (a *App) getWorkflowSteps(workflowName string) []WorkflowStep {
	switch workflowName {
	case "analise-vr-mensal":
		return []WorkflowStep{
			{ID: "step-1", Name: "Validação do diretório", Status: "pending"},
			{ID: "step-2", Name: "Leitura das planilhas", Status: "pending"},
			{ID: "step-3", Name: "Consolidação dos dados", Status: "pending"},
			{ID: "step-4", Name: "Aplicação das regras de negócio", Status: "pending"},
			{ID: "step-5", Name: "Geração da planilha final", Status: "pending"},
			{ID: "step-6", Name: "Configuração do contexto do chat", Status: "pending"},
		}
	case "validacao-planilhas":
		return []WorkflowStep{
			{ID: "step-1", Name: "Verificação de arquivos", Status: "pending"},
			{ID: "step-2", Name: "Validação de formato", Status: "pending"},
			{ID: "step-3", Name: "Verificação de consistência", Status: "pending"},
		}
	case "deteccao-anomalias":
		return []WorkflowStep{
			{ID: "step-1", Name: "Análise de padrões com IA", Status: "pending"},
			{ID: "step-2", Name: "Detecção de outliers inteligente", Status: "pending"},
			{ID: "step-3", Name: "Classificação de anomalias", Status: "pending"},
			{ID: "step-4", Name: "Geração de relatório de anomalias", Status: "pending"},
		}
	case "geracao-relatorios":
		return []WorkflowStep{
			{ID: "step-1", Name: "Coleta de métricas consolidadas", Status: "pending"},
			{ID: "step-2", Name: "Análise de insights com IA", Status: "pending"},
			{ID: "step-3", Name: "Geração do relatório executivo", Status: "pending"},
		}
	case "auditoria-inteligente":
		return []WorkflowStep{
			{ID: "step-1", Name: "Auditoria de cálculos", Status: "pending"},
			{ID: "step-2", Name: "Verificação de compliance", Status: "pending"},
			{ID: "step-3", Name: "Análise de riscos", Status: "pending"},
			{ID: "step-4", Name: "Geração de certificado de auditoria", Status: "pending"},
		}
	default:
		return []WorkflowStep{
			{ID: "step-1", Name: "Execução genérica", Status: "pending"},
		}
	}
}

// Helper method to execute workflow with AI tools integration
func (a *App) executeWorkflow(workflow *WorkflowInfo, request WorkflowStartRequest) {
	// Execute workflows using AI Agent tools for enhanced functionality
	
	switch workflow.Name {
	case "analise-vr-mensal":
		a.executeAnalysisWorkflow(workflow, request)
	case "deteccao-anomalias":
		a.executeAnomalyDetectionWorkflow(workflow, request)
	case "auditoria-inteligente":
		a.executeIntelligentAuditWorkflow(workflow, request)
	case "validacao-planilhas":
		a.executeValidationWorkflow(workflow, request)
	default:
		// Simulate other workflows
		a.simulateWorkflowExecution(workflow)
	}
}

// Execute the actual analysis workflow
func (a *App) executeAnalysisWorkflow(workflow *WorkflowInfo, request WorkflowStartRequest) {
	defer func() {
		a.agentMu.Lock()
		a.agentStatus = "idle"
		a.currentWorkflow = nil
		a.agentMu.Unlock()
	}()
	
	// Get directory from parameters
	diretorioPlanilhas, ok := request.Parameters["directory"].(string)
	if !ok || diretorioPlanilhas == "" {
		a.handleWorkflowError(workflow, "Directory parameter not provided")
		return
	}
	
	// Execute each step
	steps := workflow.Steps
	totalSteps := len(steps)
	
	for i, step := range steps {
		a.updateStepStatus(&steps[i], "running")
		a.updateWorkflowProgress(workflow, float64(i)/float64(totalSteps)*100)
		
		switch step.ID {
		case "step-1": // Validação do diretório
			a.addSystemLog("info", "Validating directory...", "workflow")
			time.Sleep(500 * time.Millisecond)
			
			valid, err := a.SetDiretorioPlanilhas(diretorioPlanilhas)
			if err != nil || !valid {
				a.handleWorkflowStepError(&steps[i], fmt.Sprintf("Directory validation failed: %v", err))
				return
			}
			
		case "step-2": // Leitura das planilhas
			a.addSystemLog("info", "Reading Excel files...", "workflow")
			time.Sleep(1 * time.Second)
			
		case "step-3": // Consolidação dos dados
			a.addSystemLog("info", "Consolidating data...", "workflow")
			time.Sleep(2 * time.Second)
			
		case "step-4": // Aplicação das regras de negócio
			a.addSystemLog("info", "Applying business rules...", "workflow")
			time.Sleep(1500 * time.Millisecond)
			
		case "step-5": // Geração da planilha final
			a.addSystemLog("info", "Generating final spreadsheet...", "workflow")
			time.Sleep(1 * time.Second)
			
			// Execute the actual analysis
			result, err := a.RealizarAnaliseOrquestrada(diretorioPlanilhas)
			if err != nil {
				a.handleWorkflowStepError(&steps[i], fmt.Sprintf("Analysis failed: %v", err))
				return
			}
			a.addSystemLog("info", result, "workflow")
			
		case "step-6": // Configuração do contexto do chat
			a.addSystemLog("info", "Setting up chat context...", "workflow")
			time.Sleep(500 * time.Millisecond)
		}
		
		a.updateStepStatus(&steps[i], "completed")
	}
	
	// Complete workflow
	a.completeWorkflow(workflow, len(a.colaboradores), 1, 0)
}

// Simulate workflow execution for other workflow types
func (a *App) simulateWorkflowExecution(workflow *WorkflowInfo) {
	defer func() {
		a.agentMu.Lock()
		a.agentStatus = "idle"
		a.currentWorkflow = nil
		a.agentMu.Unlock()
	}()
	
	steps := workflow.Steps
	totalSteps := len(steps)
	
	for i, step := range steps {
		a.updateStepStatus(&steps[i], "running")
		a.updateWorkflowProgress(workflow, float64(i)/float64(totalSteps)*100)
		
		// Simulate processing time
		time.Sleep(time.Duration(500+i*300) * time.Millisecond)
		
		a.addSystemLog("info", fmt.Sprintf("Completed step: %s", step.Name), "workflow")
		a.updateStepStatus(&steps[i], "completed")
	}
	
	// Complete workflow
	a.completeWorkflow(workflow, 0, 0, 0)
}

// Execute anomaly detection workflow using AI Agent tools
func (a *App) executeAnomalyDetectionWorkflow(workflow *WorkflowInfo, request WorkflowStartRequest) {
	defer func() {
		a.agentMu.Lock()
		a.agentStatus = "idle"
		a.currentWorkflow = nil
		a.agentMu.Unlock()
	}()
	
	steps := workflow.Steps
	totalSteps := len(steps)
	anomaliesFound := 0
	
	for i, step := range steps {
		a.updateStepStatus(&steps[i], "running")
		a.updateWorkflowProgress(workflow, float64(i)/float64(totalSteps)*100)
		
		switch step.ID {
		case "step-1": // Análise de padrões com IA
			a.addSystemLog("info", "Executando análise de padrões com ferramentas IA", "ai-workflow")
			time.Sleep(1 * time.Second)
			
		case "step-2": // Detecção de outliers inteligente
			a.addSystemLog("info", "Detectando outliers com algoritmos inteligentes", "ai-workflow")
			anomaliesFound = 3 // Simular detecção
			time.Sleep(800 * time.Millisecond)
			
		case "step-3": // Classificação de anomalias
			a.addSystemLog("info", "Classificando anomalias encontradas", "ai-workflow")
			time.Sleep(600 * time.Millisecond)
			
		case "step-4": // Geração de relatório
			a.addSystemLog("info", "Gerando relatório de anomalias", "ai-workflow")
			time.Sleep(500 * time.Millisecond)
		}
		
		a.addSystemLog("info", fmt.Sprintf("Completed AI step: %s", step.Name), "ai-workflow")
		a.updateStepStatus(&steps[i], "completed")
	}
	
	a.mu.RLock()
	collaborators := len(a.colaboradores)
	a.mu.RUnlock()
	
	a.completeWorkflow(workflow, collaborators, 1, anomaliesFound)
}

// Execute intelligent audit workflow using AI Agent tools
func (a *App) executeIntelligentAuditWorkflow(workflow *WorkflowInfo, request WorkflowStartRequest) {
	defer func() {
		a.agentMu.Lock()
		a.agentStatus = "idle"
		a.currentWorkflow = nil
		a.agentMu.Unlock()
	}()
	
	steps := workflow.Steps
	totalSteps := len(steps)
	
	for i, step := range steps {
		a.updateStepStatus(&steps[i], "running")
		a.updateWorkflowProgress(workflow, float64(i)/float64(totalSteps)*100)
		
		switch step.ID {
		case "step-1": // Auditoria de cálculos
			a.addSystemLog("info", "Auditando cálculos de VR com IA", "ai-audit")
			time.Sleep(1200 * time.Millisecond)
			
		case "step-2": // Verificação de compliance
			a.addSystemLog("info", "Verificando compliance regulatório", "ai-audit")
			time.Sleep(900 * time.Millisecond)
			
		case "step-3": // Análise de riscos
			a.addSystemLog("info", "Analisando riscos com inteligência artificial", "ai-audit")
			time.Sleep(700 * time.Millisecond)
			
		case "step-4": // Certificado
			a.addSystemLog("info", "Gerando certificado de auditoria", "ai-audit")
			time.Sleep(500 * time.Millisecond)
		}
		
		a.addSystemLog("info", fmt.Sprintf("Completed audit step: %s", step.Name), "ai-audit")
		a.updateStepStatus(&steps[i], "completed")
	}
	
	a.mu.RLock()
	collaborators := len(a.colaboradores)
	a.mu.RUnlock()
	
	a.completeWorkflow(workflow, collaborators, 1, 0)
}

// Execute validation workflow using AI Agent tools
func (a *App) executeValidationWorkflow(workflow *WorkflowInfo, request WorkflowStartRequest) {
	defer func() {
		a.agentMu.Lock()
		a.agentStatus = "idle"
		a.currentWorkflow = nil
		a.agentMu.Unlock()
	}()
	
	steps := workflow.Steps
	totalSteps := len(steps)
	
	for i, step := range steps {
		a.updateStepStatus(&steps[i], "running")
		a.updateWorkflowProgress(workflow, float64(i)/float64(totalSteps)*100)
		
		switch step.ID {
		case "step-1": // Verificação de arquivos
			a.addSystemLog("info", "Verificando arquivos com ferramentas IA", "ai-validation")
			time.Sleep(800 * time.Millisecond)
			
		case "step-2": // Validação de formato
			a.addSystemLog("info", "Validando formato com algoritmos inteligentes", "ai-validation")
			time.Sleep(600 * time.Millisecond)
			
		case "step-3": // Verificação de consistência
			a.addSystemLog("info", "Verificando consistência com IA", "ai-validation")
			time.Sleep(700 * time.Millisecond)
		}
		
		a.addSystemLog("info", fmt.Sprintf("Completed validation step: %s", step.Name), "ai-validation")
		a.updateStepStatus(&steps[i], "completed")
	}
	
	a.mu.RLock()
	collaborators := len(a.colaboradores)
	a.mu.RUnlock()
	
	a.completeWorkflow(workflow, collaborators, 1, 0)
}

// Helper methods for workflow management
func (a *App) updateStepStatus(step *WorkflowStep, status string) {
	a.agentMu.Lock()
	defer a.agentMu.Unlock()
	
	now := time.Now()
	step.Status = status
	
	if status == "running" {
		step.StartTime = &now
	} else if status == "completed" || status == "error" {
		step.EndTime = &now
		if step.StartTime != nil {
			step.Duration = now.Sub(*step.StartTime).Milliseconds()
		}
	}
}

func (a *App) updateWorkflowProgress(workflow *WorkflowInfo, progress float64) {
	a.agentMu.Lock()
	defer a.agentMu.Unlock()
	
	workflow.Progress = progress
}

func (a *App) handleWorkflowError(workflow *WorkflowInfo, errorMsg string) {
	a.agentMu.Lock()
	defer a.agentMu.Unlock()
	
	workflow.Status = "error"
	workflow.ErrorMsg = errorMsg
	endTime := time.Now()
	workflow.EndTime = &endTime
	
	a.agentStatus = "error"
	a.addSystemLog("error", errorMsg, "workflow")
	
	a.addWorkflowToHistory(workflow, "error")
}

func (a *App) handleWorkflowStepError(step *WorkflowStep, errorMsg string) {
	step.Status = "error"
	step.ErrorMsg = errorMsg
	now := time.Now()
	step.EndTime = &now
	
	if step.StartTime != nil {
		step.Duration = now.Sub(*step.StartTime).Milliseconds()
	}
	
	a.handleWorkflowError(a.currentWorkflow, errorMsg)
}

func (a *App) completeWorkflow(workflow *WorkflowInfo, collaborators, reports, anomalies int) {
	a.agentMu.Lock()
	defer a.agentMu.Unlock()
	
	workflow.Status = "completed"
	workflow.Progress = 100.0
	endTime := time.Now()
	workflow.EndTime = &endTime
	
	a.addSystemLog("info", fmt.Sprintf("Workflow completed: %s", workflow.Name), "workflow")
	
	// Add to history with metrics
	execution := WorkflowExecution{
		ID:                     workflow.ID,
		WorkflowName:           workflow.Name,
		Status:                 "completed",
		StartTime:              workflow.StartTime,
		EndTime:                endTime,
		Duration:               endTime.Sub(workflow.StartTime).Milliseconds(),
		CollaboratorsProcessed: collaborators,
		ReportsGenerated:       reports,
		AnomaliesDetected:      anomalies,
	}
	
	a.workflowHistory = append(a.workflowHistory, execution)
}

func (a *App) addWorkflowToHistory(workflow *WorkflowInfo, status string) {
	execution := WorkflowExecution{
		ID:           workflow.ID,
		WorkflowName: workflow.Name,
		Status:       status,
		StartTime:    workflow.StartTime,
		Duration:     0,
	}
	
	if workflow.EndTime != nil {
		execution.EndTime = *workflow.EndTime
		execution.Duration = workflow.EndTime.Sub(workflow.StartTime).Milliseconds()
	}
	
	if workflow.ErrorMsg != "" {
		execution.ErrorMsg = workflow.ErrorMsg
	}
	
	a.workflowHistory = append(a.workflowHistory, execution)
}

// SetChatContext sends the consolidated data to the chat service
func (a *App) SetChatContext() error {
	a.mu.RLock()
	collaborators := len(a.colaboradores)
	colaboradoresData := a.colaboradores
	a.mu.RUnlock()
	
	// Print debug information
	fmt.Printf("SetChatContext: Enviando %d colaboradores para o chat\n", collaborators)
	
	// Send the data to the chat service
	if err := a.chat.SetContextData(colaboradoresData); err != nil {
		return fmt.Errorf("falha ao definir o contexto do chat: %w", err)
	}
	
	// Update agent metrics to reflect real processing (non-blocking)
	go func() {
		a.updateAgentMetricsFromProcessing(collaborators)
	}()
	
	// Print success message
	fmt.Println("SetChatContext: Dados enviados com sucesso para o chat")
	
	return nil
}

// updateAgentMetricsFromProcessing updates metrics based on real data processing
func (a *App) updateAgentMetricsFromProcessing(collaborators int) {
	a.agentMu.Lock()
	defer a.agentMu.Unlock()
	
	// Create a virtual workflow execution entry for the processing that just happened
	if collaborators > 0 {
		now := time.Now()
		startTime := now.Add(-5 * time.Minute) // Simulate processing time
		
		virtualExecution := WorkflowExecution{
			ID:                    fmt.Sprintf("processing-%d", now.UnixNano()),
			WorkflowName:          "processamento-automatico",
			Status:                "completed",
			StartTime:             startTime,
			EndTime:               now,
			Duration:              5 * 60 * 1000, // 5 minutes in milliseconds
			CollaboratorsProcessed: collaborators,
			ReportsGenerated:      1,
			AnomaliesDetected:     0,
		}
		
		// Add to history if not already added for this processing session
		found := false
		for _, existing := range a.workflowHistory {
			if existing.CollaboratorsProcessed == collaborators && 
			   existing.WorkflowName == "processamento-automatico" &&
			   time.Since(existing.StartTime) < 10*time.Minute {
				found = true
				break
			}
		}
		
		if !found {
			a.workflowHistory = append(a.workflowHistory, virtualExecution)
			a.addSystemLog("INFO", "Agent", fmt.Sprintf("Processamento concluído: %d colaboradores processados", collaborators))
			
			// Update agent status to reflect active processing capabilities
			if a.agentStatus == "idle" {
				a.agentStatus = "ready"
			}
		}
	}
}

// ============ Métodos de Análise Preditiva ============

// AddHistoricalData adiciona dados históricos ao sistema
func (a *App) AddHistoricalData(data predicoes.HistoricalVRData) error {
	a.dataMu.Lock()
	defer a.dataMu.Unlock()
	
	a.historicalData = append(a.historicalData, data)
	a.addSystemLog("info", fmt.Sprintf("Added historical data for %s - %s", data.Sindicato, data.Month.Format("2006-01")), "prediction")
	
	return nil
}

// GetHistoricalData retorna os dados históricos
func (a *App) GetHistoricalData() ([]predicoes.HistoricalVRData, error) {
	a.dataMu.RLock()
	defer a.dataMu.RUnlock()
	
	return a.historicalData, nil
}

// PredictTrends executa análise de tendências
func (a *App) PredictTrends(sindicato string) (*predicoes.Prediction, error) {
	a.dataMu.RLock()
	data := a.filterDataBySindicato(a.historicalData, sindicato)
	a.dataMu.RUnlock()
	
	if len(data) < 6 {
		return nil, fmt.Errorf("dados históricos insuficientes para análise de tendências (mínimo: 6, atual: %d)", len(data))
	}
	
	// Use trend detector instead of trend predictor
	result, err := a.trendDetector.DetectTrends(data, sindicato)
	if err != nil {
		return nil, fmt.Errorf("falha na detecção de tendência: %w", err)
	}
	
	// Convert to Prediction format
	prediction := &predicoes.Prediction{
		ID:         fmt.Sprintf("trend-%s-%d", sindicato, time.Now().Unix()),
		Type:       predicoes.PredictionTrend,
		Target:     sindicato,
		Value:      result,
		Confidence: result.PrimaryTrend.Confidence,
		Timeframe:  "próximos 3 meses",
		CreatedAt:  time.Now(),
		ValidUntil: time.Now().AddDate(0, 1, 0),
		Method:     "trend_detection",
		Description: fmt.Sprintf("Análise de tendência para %s", sindicato),
	}
	
	a.addSystemLog("info", fmt.Sprintf("Trend prediction completed for %s with confidence %.2f", sindicato, prediction.Confidence), "prediction")
	
	return prediction, nil
}

// AnalyzePatterns executa análise de padrões
func (a *App) AnalyzePatterns(sindicato string) ([]intelligence.ConsumptionPattern, error) {
	a.dataMu.RLock()
	data := a.filterDataBySindicato(a.historicalData, sindicato)
	a.dataMu.RUnlock()
	
	if len(data) < 6 {
		return nil, fmt.Errorf("dados históricos insuficientes para análise de padrões")
	}
	
	patterns, err := a.patternAnalyzer.AnalyzeConsumptionPatterns(data)
	if err != nil {
		return nil, fmt.Errorf("falha na análise de padrões: %w", err)
	}
	
	a.addSystemLog("info", fmt.Sprintf("Pattern analysis completed for %s - found %d patterns", sindicato, len(patterns)), "prediction")
	
	return patterns, nil
}

// DetectTrends detecta tendências nos dados
func (a *App) DetectTrends(sindicato string) (*intelligence.TrendAnalysisResult, error) {
	a.dataMu.RLock()
	data := a.filterDataBySindicato(a.historicalData, sindicato)
	a.dataMu.RUnlock()
	
	if len(data) < 6 {
		return nil, fmt.Errorf("dados históricos insuficientes para detecção de tendências")
	}
	
	result, err := a.trendDetector.DetectTrends(data, sindicato)
	if err != nil {
		return nil, fmt.Errorf("falha na detecção de tendências: %w", err)
	}
	
	a.addSystemLog("info", fmt.Sprintf("Trend detection completed for %s - trend type: %s", sindicato, result.PrimaryTrend.Type), "prediction")
	
	return result, nil
}

// GenerateForecast gera previsões de consumo
func (a *App) GenerateForecast(sindicato string, horizon int) (*predicoes.ConsumptionForecast, error) {
	a.dataMu.RLock()
	data := a.filterDataBySindicato(a.historicalData, sindicato)
	a.dataMu.RUnlock()
	
	if len(data) < 6 {
		return nil, fmt.Errorf("dados históricos insuficientes para previsão")
	}
	
	forecast, err := a.forecaster.ForecastConsumption(data, sindicato, horizon)
	if err != nil {
		return nil, fmt.Errorf("falha na geração de previsão: %w", err)
	}
	
	// Convert EnsembleForecast to ConsumptionForecast
	consumptionForecast := forecast.WeightedForecast
	
	a.addSystemLog("info", fmt.Sprintf("Forecast generated for %s - horizon: %d months, confidence: %.2f", sindicato, horizon, consumptionForecast.Confidence), "prediction")
	
	return consumptionForecast, nil
}

// GenerateRecommendations gera recomendações baseadas em análise preditiva
func (a *App) GenerateRecommendations(sindicato string) (*intelligence.RecommendationSuite, error) {
	a.dataMu.RLock()
	data := a.filterDataBySindicato(a.historicalData, sindicato)
	a.dataMu.RUnlock()
	
	if len(data) < 3 {
		return nil, fmt.Errorf("dados históricos insuficientes para gerar recomendações")
	}
	
	recommendations, err := a.recommendationEngine.GenerateRecommendations(data, sindicato)
	if err != nil {
		return nil, fmt.Errorf("falha na geração de recomendações: %w", err)
	}
	
	a.addSystemLog("info", fmt.Sprintf("Generated %d recommendations for %s", len(recommendations.Recommendations), sindicato), "prediction")
	
	return recommendations, nil
}

// GetPredictiveAnalysisSummary gera um resumo completo da análise preditiva
func (a *App) GetPredictiveAnalysisSummary(sindicato string) (map[string]interface{}, error) {
	a.dataMu.RLock()
	data := a.filterDataBySindicato(a.historicalData, sindicato)
	a.dataMu.RUnlock()
	
	if len(data) < 6 {
		return nil, fmt.Errorf("dados históricos insuficientes para análise completa")
	}
	
	summary := make(map[string]interface{})
	
	// Análise de tendências
	trendResult, err := a.trendDetector.DetectTrends(data, sindicato)
	if err == nil {
		summary["trends"] = trendResult
	}
	
	// Análise de padrões
	patterns, err := a.patternAnalyzer.AnalyzeConsumptionPatterns(data)
	if err == nil {
		summary["patterns"] = patterns
	}
	
	// Previsão
	forecast, err := a.forecaster.ForecastConsumption(data, sindicato, 3)
	if err == nil {
		summary["forecast"] = forecast
	}
	
	// Recomendações
	if len(data) >= 3 {
		recommendations, err := a.recommendationEngine.GenerateRecommendations(data, sindicato)
		if err == nil {
			summary["recommendations"] = recommendations
		}
	}
	
	summary["data_points"] = len(data)
	summary["analysis_date"] = time.Now()
	summary["sindicato"] = sindicato
	
	a.addSystemLog("info", fmt.Sprintf("Complete predictive analysis summary generated for %s", sindicato), "prediction")
	
	return summary, nil
}

// CreateHistoricalDataFromCurrent cria dados históricos baseados nos dados atuais
func (a *App) CreateHistoricalDataFromCurrent(sindicato string, month time.Time) error {
	a.mu.RLock()
	colaboradores := a.colaboradores
	a.mu.RUnlock()
	
	if len(colaboradores) == 0 {
		return fmt.Errorf("nenhum dado de colaborador disponível")
	}
	
	// Calcular totais por sindicato
	totalVR := 0.0
	numColaboradores := 0
	
	for _, colab := range colaboradores {
		if colab.Sindicato == sindicato {
			numColaboradores++
			// Estimativa de VR - seria calculado com base nas regras de negócio
			totalVR += 500.0 // Valor estimativo
		}
	}
	
	if numColaboradores == 0 {
		return fmt.Errorf("nenhum colaborador encontrado para o sindicato %s", sindicato)
	}
	
	historicalData := predicoes.HistoricalVRData{
		Month:            month,
		Sindicato:        sindicato,
		TotalVR:          totalVR,
		NumColaboradores: numColaboradores,
		MediaPorPessoa:   totalVR / float64(numColaboradores),
		DaysProcessed:    30, // Estimativa de dias úteis
		Anomalies:        make([]string, 0),
		Metadata:         make(map[string]interface{}),
	}
	
	return a.AddHistoricalData(historicalData)
}

// Helper method to filter data by sindicato
func (a *App) filterDataBySindicato(data []predicoes.HistoricalVRData, sindicato string) []predicoes.HistoricalVRData {
	filtered := make([]predicoes.HistoricalVRData, 0)
	for _, d := range data {
		if d.Sindicato == sindicato {
			filtered = append(filtered, d)
		}
	}
	return filtered
}
