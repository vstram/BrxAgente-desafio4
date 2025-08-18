// Package config provides functionality for handling application configuration
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	// OpenAI API key
	OpenAIKey string `json:"openai_key,omitempty"`
	
	// Ollama configuration
	OllamaConfig OllamaConfig `json:"ollama_config,omitempty"`
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
		return &Config{}, nil
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
	
	return &config, nil
}

// SaveConfig saves the configuration to file
func SaveConfig(config *Config) error {
	// Get config file path
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	
	// Convert config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
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