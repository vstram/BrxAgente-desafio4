package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Encrypt encrypts plaintext using the provided key
func Encrypt(key []byte, plaintext string) (string, error) {
	// Create a new AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encode plaintext to bytes
	plaintextBytes := []byte(plaintext)

	// Pad plaintext to be multiple of block size
	plaintextBytes = pad(plaintextBytes, aes.BlockSize)

	// Encrypt the plaintext
	ciphertext := make([]byte, len(plaintextBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintextBytes)

	// Prepend IV to ciphertext
	result := append(iv, ciphertext...)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts ciphertext using the provided key
func Decrypt(key []byte, ciphertext string) (string, error) {
	// Decode from base64
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Extract IV
	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	// Decrypt the ciphertext
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertextBytes, ciphertextBytes)

	// Remove padding
	ciphertextBytes, err = unpad(ciphertextBytes)
	if err != nil {
		return "", err
	}

	return string(ciphertextBytes), nil
}

// pad adds PKCS#7 padding to the plaintext
func pad(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(plaintext, padtext...)
}

// unpad removes PKCS#7 padding from the plaintext
func unpad(plaintext []byte) ([]byte, error) {
	length := len(plaintext)
	if length == 0 {
		return nil, fmt.Errorf("plaintext is empty")
	}
	unpadding := int(plaintext[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("invalid padding")
	}
	return plaintext[:(length - unpadding)], nil
}