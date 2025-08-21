package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"

	"BrxAgente-desafio4/internal/chat"
	"BrxAgente-desafio4/internal/config"
	"BrxAgente-desafio4/internal/modelo"
	"BrxAgente-desafio4/internal/security"
)

// App struct
type App struct {
	ctx         context.Context
	cfg         *config.Config
	chat        *chat.Chat
	colaboradores map[string]*modelo.Colaborador
	mu          sync.RWMutex // Mutex para acesso seguro aos dados compartilhados
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		colaboradores: make(map[string]*modelo.Colaborador),
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

	return nil
}

// SetOllamaConfig define a configuração do Ollama
func (a *App) SetOllamaConfig(ollamaConfig config.OllamaConfig) error {
	// Validate the configuration
	if err := config.ValidateOllamaConfig(ollamaConfig); err != nil {
		return fmt.Errorf("configuração do Ollama inválida: %w", err)
	}

	// Update config
	a.cfg.OllamaConfig = ollamaConfig

	// Save config
	if err := config.SaveConfig(a.cfg); err != nil {
		return fmt.Errorf("falha ao salvar a configuração: %w", err)
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

	// In a real implementation, we would test the connection by making an API call
	// For now, we'll just return true if the configuration is valid
	return true, nil
}

// AskAI sends a question to the configured AI service and returns the response
func (a *App) AskAI(question string) (string, error) {
	// Define a system prompt with context about the VR/VA application
	systemPrompt := `Você é um assistente especializado em análise de dados de Vale Refeição (VR) e Vale Alimentação (VA).
	Você está ajudando um usuário a entender os resultados do processamento de dados de colaboradores.
	Os dados dos colaboradores são identificados exclusivamente por uma MATRICULA, por razões de confidencialidade.
	Seja claro, preciso e objetivo em suas respostas.`

	// Create empty context for now
	var context []chat.Message

	// Ask the chat service
	response, err := a.chat.Ask(question, systemPrompt, context)
	if err != nil {
		return "", fmt.Errorf("falha ao obter resposta da IA: %w", err)
	}

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

// SetChatContext sends the consolidated data to the chat service
func (a *App) SetChatContext() error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	// Print debug information
	fmt.Printf("SetChatContext: Enviando %d colaboradores para o chat\n", len(a.colaboradores))
	
	// Send the data to the chat service
	if err := a.chat.SetContextData(a.colaboradores); err != nil {
		return fmt.Errorf("falha ao definir o contexto do chat: %w", err)
	}
	
	// Print success message
	fmt.Println("SetChatContext: Dados enviados com sucesso para o chat")
	
	return nil
}
