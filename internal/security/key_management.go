package security

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// GenerateKey generates a new 32-byte key for AES-256
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// GetKeyPath returns the path to the encryption key file
func GetKeyPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	
	// Create config directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".brxagente")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Return path to key file
	return filepath.Join(configDir, "key"), nil
}

// SaveKey saves the encryption key to a file
func SaveKey(key []byte) error {
	keyPath, err := GetKeyPath()
	if err != nil {
		return err
	}
	
	// Encode key to hex for storage
	encodedKey := hex.EncodeToString(key)
	
	// Write to file with restricted permissions
	if err := os.WriteFile(keyPath, []byte(encodedKey), 0600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}
	
	return nil
}

// LoadKey loads the encryption key from a file
func LoadKey() ([]byte, error) {
	keyPath, err := GetKeyPath()
	if err != nil {
		return nil, err
	}
	
	// Check if key file exists
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("key file does not exist")
	}
	
	// Read key file
	encodedKey, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}
	
	// Decode key from hex
	key, err := hex.DecodeString(string(encodedKey))
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}
	
	// Verify key length
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: expected 32 bytes, got %d", len(key))
	}
	
	return key, nil
}