package security

import (
	"fmt"
)

// SecureString represents a string that is encrypted in memory
type SecureString struct {
	encryptedData string
}

// NewSecureString creates a new SecureString from a plain text value
func NewSecureString(plaintext string) (*SecureString, error) {
	// Load or generate encryption key
	key, err := LoadKey()
	if err != nil {
		// If key doesn't exist, generate a new one
		if err.Error() == "key file does not exist" {
			key, err = GenerateKey()
			if err != nil {
				return nil, fmt.Errorf("failed to generate key: %w", err)
			}
			
			// Save the new key
			if err := SaveKey(key); err != nil {
				return nil, fmt.Errorf("failed to save key: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to load key: %w", err)
		}
	}
	
	// Encrypt the plaintext
	encryptedData, err := Encrypt(key, plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}
	
	return &SecureString{
		encryptedData: encryptedData,
	}, nil
}

// GetValue decrypts and returns the plain text value
func (s *SecureString) GetValue() (string, error) {
	// Load encryption key
	key, err := LoadKey()
	if err != nil {
		return "", fmt.Errorf("failed to load key: %w", err)
	}
	
	// Decrypt the data
	plaintext, err := Decrypt(key, s.encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %w", err)
	}
	
	return plaintext, nil
}

// Destroy clears the encrypted data
func (s *SecureString) Destroy() {
	s.encryptedData = ""
}