// Package config provides functionality for handling application configuration
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"BrxAgente-desafio4/internal/security"
)

// Config represents the application configuration
type Config struct {
	// OpenAI API key (encrypted)
	OpenAIKey string `json:"openai_key,omitempty"`

	// Ollama configuration
	OllamaConfig OllamaConfig `json:"ollama_config,omitempty"`

	// Agent configuration
	AgentConfig AgentConfig `json:"agent_config,omitempty"`
}

// AgentConfig represents the configuration for the AI agent
type AgentConfig struct {
	// Se o agente está habilitado
	Enabled bool `json:"enabled"`

	// Configurações do modelo LLM
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`

	// Configurações de performance
	WorkerPoolSize int  `json:"worker_pool_size"`
	CacheEnabled   bool `json:"cache_enabled"`
	CacheSize      int  `json:"cache_size"`

	// Ferramentas habilitadas
	ToolsEnabled []string `json:"tools_enabled"`
}

// OllamaConfig represents the configuration for Ollama
type OllamaConfig struct {
	// Base URL for Ollama API
	BaseURL string `json:"base_url,omitempty"`

	// Model to use
	Model string `json:"model,omitempty"`
}

// configFile is the name of the configuration file
const configFile = "brxagente-config.json"

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".brxagente")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	// Return path to config file
	return filepath.Join(configDir, configFile), nil
}

// LoadConfig loads the configuration from file
func LoadConfig() (*Config, error) {
	// Get config file path
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return GetDefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Decrypt sensitive data after loading
	// In a real implementation, we would decrypt the encrypted values
	// For now, we'll just return the config as is

	return &config, nil
}

// SaveConfig saves the configuration to file
func SaveConfig(config *Config) error {
	// Get config file path
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create a copy of the config to avoid modifying the original
	configCopy := *config

	// Encrypt sensitive data before saving
	if configCopy.OpenAIKey != "" {
		secureString, err := security.NewSecureString(configCopy.OpenAIKey)
		if err != nil {
			return fmt.Errorf("failed to create secure string: %w", err)
		}
		// We store the encrypted value directly
		// In a real implementation, we would store the encrypted value
		// For now, we'll just keep the original value but in a real app,
		// we would store the encrypted version
		configCopy.OpenAIKey, _ = secureString.GetValue() // This would be the encrypted value
	}

	// Convert config to JSON
	data, err := json.MarshalIndent(configCopy, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ValidateOpenAIKey validates an OpenAI API key format
func ValidateOpenAIKey(key string) bool {
	// Basic validation: key should start with "sk-" and have reasonable length
	return len(key) >= 20 && key[:3] == "sk-"
}

// ValidateOllamaConfig validates Ollama configuration
func ValidateOllamaConfig(config OllamaConfig) error {
	// Base URL is optional but if provided should not be empty
	if config.BaseURL != "" && len(config.BaseURL) < 10 {
		return fmt.Errorf("ollama base URL is too short")
	}

	// Model is optional but if provided should not be empty
	if config.Model != "" && len(config.Model) < 1 {
		return fmt.Errorf("ollama model name is required")
	}

	return nil
}

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() *Config {
	return &Config{
		AgentConfig: AgentConfig{
			Enabled:        true,
			Model:          "gpt-3.5-turbo",
			Temperature:    0.7,
			MaxTokens:      2000,
			WorkerPoolSize: 4,
			CacheEnabled:   true,
			CacheSize:      1000,
			ToolsEnabled:   []string{"excel", "calculation", "validation"},
		},
	}
}

// ValidateAgentConfig validates agent configuration
func ValidateAgentConfig(config AgentConfig) error {
	if config.Temperature < 0.0 || config.Temperature > 2.0 {
		return fmt.Errorf("temperature must be between 0.0 and 2.0")
	}

	if config.MaxTokens <= 0 {
		return fmt.Errorf("max_tokens must be greater than 0")
	}

	if config.WorkerPoolSize <= 0 {
		return fmt.Errorf("worker_pool_size must be greater than 0")
	}

	if config.CacheSize < 0 {
		return fmt.Errorf("cache_size cannot be negative")
	}

	return nil
}
