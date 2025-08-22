package agent

import (
	"time"
)

// AgentConfig contém as configurações específicas do agente de IA
type AgentConfig struct {
	// Configurações gerais do agente
	Enabled     bool          `json:"enabled"`      // Se o agente está habilitado
	Model       string        `json:"model"`        // Modelo LLM a ser usado
	Temperature float64       `json:"temperature"`  // Temperatura para geração de texto
	MaxTokens   int           `json:"max_tokens"`   // Máximo de tokens por resposta
	Timeout     time.Duration `json:"timeout"`      // Timeout para chamadas LLM

	// Configurações de memória
	MemorySize       int           `json:"memory_size"`        // Tamanho do buffer de memória
	MemoryTTL        time.Duration `json:"memory_ttl"`         // TTL da memória de conversação
	ContextWindow    int           `json:"context_window"`     // Janela de contexto
	MaxMemoryTokens  int           `json:"max_memory_tokens"`  // Máximo de tokens na memória
	
	// Configurações de performance
	WorkerPoolSize  int           `json:"worker_pool_size"`  // Tamanho do pool de workers
	CacheEnabled    bool          `json:"cache_enabled"`     // Se o cache está habilitado
	CacheSize       int           `json:"cache_size"`        // Tamanho do cache LLM
	CacheTTL        time.Duration `json:"cache_ttl"`         // TTL do cache

	// Configurações de logging
	LogLevel        string        `json:"log_level"`         // Nível de log
	DebugMode       bool          `json:"debug_mode"`        // Modo debug
	
	// Configurações de ferramentas
	ToolsEnabled    []string      `json:"tools_enabled"`     // Lista de ferramentas habilitadas
}

// DefaultAgentConfig retorna a configuração padrão do agente
func DefaultAgentConfig() *AgentConfig {
	return &AgentConfig{
		Enabled:         true,
		Model:          "gpt-3.5-turbo",
		Temperature:    0.7,
		MaxTokens:      2000,
		Timeout:        30 * time.Second,
		
		MemorySize:      100,
		MemoryTTL:       24 * time.Hour,
		ContextWindow:   4000,
		MaxMemoryTokens: 2000,
		
		WorkerPoolSize: 4,
		CacheEnabled:   true,
		CacheSize:      1000,
		CacheTTL:       1 * time.Hour,
		
		LogLevel:       "info",
		DebugMode:      false,
		
		ToolsEnabled:   []string{"excel", "calculation", "validation"},
	}
}

// Validate valida se a configuração do agente é válida
func (c *AgentConfig) Validate() error {
	if c.Temperature < 0.0 || c.Temperature > 2.0 {
		return &AgentConfigError{
			Field:   "temperature",
			Message: "temperatura deve estar entre 0.0 e 2.0",
		}
	}
	
	if c.MaxTokens <= 0 {
		return &AgentConfigError{
			Field:   "max_tokens",
			Message: "max_tokens deve ser maior que 0",
		}
	}
	
	if c.MemorySize <= 0 {
		return &AgentConfigError{
			Field:   "memory_size",
			Message: "memory_size deve ser maior que 0",
		}
	}
	
	if c.WorkerPoolSize <= 0 {
		return &AgentConfigError{
			Field:   "worker_pool_size",
			Message: "worker_pool_size deve ser maior que 0",
		}
	}
	
	return nil
}

// AgentConfigError representa um erro de configuração do agente
type AgentConfigError struct {
	Field   string
	Message string
}

func (e *AgentConfigError) Error() string {
	return "erro de configuração do agente no campo '" + e.Field + "': " + e.Message
}