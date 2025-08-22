package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"BrxAgente-desafio4/internal/config"
	"BrxAgente-desafio4/internal/modelo"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIRequest represents a request to the OpenAI API
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// OpenAIResponse represents a response from the OpenAI API
type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int     `json:"index"`
		Message Message `json:"message"`
	} `json:"choices"`
}

// OllamaRequest represents a request to the Ollama API
type OllamaRequest struct {
	Model   string        `json:"model"`
	Prompt  string        `json:"prompt"`
	Stream  bool          `json:"stream"`
	System  string        `json:"system,omitempty"`
	Format  string        `json:"format,omitempty"`
	Options OllamaOptions `json:"options,omitempty"`
}

// OllamaOptions represents options for Ollama requests
type OllamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
}

// OllamaResponse represents a response from the Ollama API
type OllamaResponse struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	Context            []int     `json:"context,omitempty"`
	TotalDuration      int64     `json:"total_duration,omitempty"`
	LoadDuration       int64     `json:"load_duration,omitempty"`
	PromptEvalCount    int       `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64     `json:"prompt_eval_duration,omitempty"`
	EvalCount          int       `json:"eval_count,omitempty"`
	EvalDuration       int64     `json:"eval_duration,omitempty"`
}

// Chat handles chat interactions with AI services
type Chat struct {
	cfg         *config.Config
	contextData map[string]*modelo.Colaborador
	agent       AgentInterface // Interface para integração com o agente
}

// AgentInterface define a interface para integração com o agente de IA
type AgentInterface interface {
	Ask(question string) (string, error)
	IsEnabled() bool
	GetStatusInterface() interface{}
}

// NewChat creates a new Chat instance
func NewChat(cfg *config.Config) *Chat {
	return &Chat{
		cfg:         cfg,
		contextData: make(map[string]*modelo.Colaborador),
	}
}

// AskOpenAI sends a question to OpenAI and returns the response
func (c *Chat) AskOpenAI(question string, context []Message) (string, error) {
	// Check if OpenAI key is configured
	if c.cfg.OpenAIKey == "" {
		return "", fmt.Errorf("chave da API do OpenAI não configurada")
	}

	// Prepare the context data as a string
	contextDataStr := c.formatContextData()

	// Add context data as a system message if available
	if contextDataStr != "" {
		context = append([]Message{
			{
				Role:    "system",
				Content: fmt.Sprintf("Contexto dos dados:\n%s", contextDataStr),
			},
		}, context...)
	}

	// Prepare messages
	messages := append(context, Message{
		Role:    "user",
		Content: question,
	})

	// Create request
	request := OpenAIRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("falha ao serializar requisição: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("falha ao criar requisição HTTP: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.OpenAIKey)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("falha ao enviar requisição para OpenAI: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API retornou status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("falha ao parsear resposta da OpenAI: %w", err)
	}

	// Return the response content
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("resposta da OpenAI vazia")
}

// AskOllama sends a question to Ollama and returns the response
func (c *Chat) AskOllama(question string, systemPrompt string) (string, error) {
	// Check if Ollama is configured
	if c.cfg.OllamaConfig.BaseURL == "" {
		return "", fmt.Errorf("configuração do Ollama não encontrada")
	}

	// Prepare the context data as a string
	contextDataStr := c.formatContextData()

	// Print debug information
	fmt.Printf("AskOllama: Tamanho dos dados de contexto formatados: %d caracteres\n", len(contextDataStr))

	// Prepend context data to the system prompt
	fullSystemPrompt := systemPrompt
	if contextDataStr != "" {
		fullSystemPrompt = fmt.Sprintf("Contexto dos dados:\n%s\n\n%s", contextDataStr, systemPrompt)
	}

	// Print debug information
	fmt.Printf("AskOllama: Tamanho do prompt completo: %d caracteres\n", len(fullSystemPrompt))

	// Create request
	request := OllamaRequest{
		Model:  c.cfg.OllamaConfig.Model,
		Prompt: question,
		Stream: false,
		System: fullSystemPrompt,
		Options: OllamaOptions{
			Temperature: 0.7,
		},
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("falha ao serializar requisição: %w", err)
	}

	// Create full URL
	url := c.cfg.OllamaConfig.BaseURL
	if url[len(url)-1] != '/' {
		url += "/"
	}
	url += "api/generate"

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("falha ao criar requisição HTTP: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("falha ao enviar requisição para Ollama: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API retornou status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("falha ao parsear resposta do Ollama: %w", err)
	}

	return response.Response, nil
}

// SetContextData updates the context data for the chat
func (c *Chat) SetContextData(data map[string]*modelo.Colaborador) error {
	// For now, just store a reference to the data
	// In a real implementation, we might want to process or copy the data
	c.contextData = data

	// Print debug information
	fmt.Printf("SetContextData: Recebido %d colaboradores\n", len(data))

	return nil
}

// GetContextDataAsString returns the context data formatted as a string
func (c *Chat) GetContextDataAsString() string {
	return c.formatContextData()
}

// formatContextData formats the context data as a string for inclusion in the prompt
func (c *Chat) formatContextData() string {
	if len(c.contextData) == 0 {
		return ""
	}

	// Create a more detailed summary of the data
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Dados de %d colaboradores disponíveis:\n", len(c.contextData)))

	// Add details for each collaborator (limiting to 10 for brevity)
	count := 0
	for _, colaborador := range c.contextData {
		// if count >= 10 {
		// 	summary.WriteString(fmt.Sprintf("... e mais %d colaboradores\n", len(c.contextData)-10))
		// 	break
		// }

		// Format the collaborator data in a more structured way
		summary.WriteString(fmt.Sprintf(
			"Colaborador %d:\n  Matrícula: %s\n  Empresa: %s\n  Sindicato: %s\n  Valor Total VR: R$ %.2f\n  Valor Empresa: R$ %.2f\n  Valor Colaborador: R$ %.2f\n  Dias Úteis: %d\n\n",
			count+1,
			colaborador.Matricula,
			colaborador.Empresa,
			colaborador.Sindicato,
			colaborador.ValorTotalVR,
			colaborador.ValorEmpresa,
			colaborador.ValorColaborador,
			colaborador.DiasUteisEfetivos,
		))
		count++
	}

	return summary.String()
}

// SetAgent configura o agente de IA para o chat
func (c *Chat) SetAgent(agent AgentInterface) {
	c.agent = agent
}

// Ask sends a question to the configured AI service and returns the response
func (c *Chat) Ask(question string, systemPrompt string, context []Message) (string, error) {
	// Print debug information
	fmt.Printf("Ask: Recebida pergunta: %s\n", question)
	fmt.Printf("Ask: Tamanho do contexto: %d\n", len(context))
	fmt.Printf("Ask: Tamanho dos dados de contexto: %d\n", len(c.contextData))

	// Try agent first if configured and enabled
	if c.agent != nil && c.agent.IsEnabled() {
		response, err := c.agent.Ask(question)
		if err == nil {
			fmt.Printf("Ask: Resposta obtida via agente\n")
			return response, nil
		}
		// If agent fails, log the error and fallback to other services
		fmt.Printf("Warning: Agent request failed, fallback to other services: %v\n", err)
	}

	// Try OpenAI first if configured
	if c.cfg.OpenAIKey != "" {
		response, err := c.AskOpenAI(question, context)
		if err == nil {
			return response, nil
		}
		// If OpenAI fails, log the error and try Ollama
		fmt.Printf("Warning: OpenAI request failed: %v\n", err)
	}

	// Try Ollama if configured
	if c.cfg.OllamaConfig.BaseURL != "" && c.cfg.OllamaConfig.Model != "" {
		response, err := c.AskOllama(question, systemPrompt)
		if err == nil {
			return response, nil
		}
		// If Ollama fails, return the error
		return "", fmt.Errorf("falha ao obter resposta de ambos os serviços de IA: %w", err)
	}

	// If neither is configured, return an error
	return "", fmt.Errorf("nenhum serviço de IA configurado (OpenAI ou Ollama)")
}
