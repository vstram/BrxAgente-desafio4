// Package cache implementa cache inteligente para otimização de performance
package cache

import (
	"crypto/md5"
	"fmt"
	"strings"
	"sync"
	"time"
)

// CacheEntry representa uma entrada no cache
type CacheEntry struct {
	Response  string
	Metadata  CacheMetadata
	CreatedAt time.Time
	LastUsed  time.Time
	UseCount  int
}

// CacheMetadata contém metadados da entrada do cache
type CacheMetadata struct {
	TokensUsed int
	Duration   time.Duration
	Model      string
	Quality    float64
}

// LLMCache implementa cache inteligente para respostas de LLM
type LLMCache struct {
	cache     map[string]*CacheEntry
	mutex     sync.RWMutex
	ttl       time.Duration
	maxSize   int
	hitCount  int64
	missCount int64
}

// NewLLMCache cria uma nova instância de cache LLM
func NewLLMCache(maxSize int, ttl time.Duration) *LLMCache {
	cache := &LLMCache{
		cache:   make(map[string]*CacheEntry),
		ttl:     ttl,
		maxSize: maxSize,
	}

	// Inicia goroutine para limpeza periódica
	go cache.cleanupExpiredEntries()

	return cache
}

// generateKey gera uma chave única para o prompt
func (c *LLMCache) generateKey(prompt string) string {
	// Normaliza o prompt removendo espaços extras
	normalized := strings.TrimSpace(strings.ReplaceAll(prompt, "\n", " "))
	normalized = strings.ReplaceAll(normalized, "  ", " ")

	// Gera hash MD5 para a chave
	hash := md5.Sum([]byte(normalized))
	return fmt.Sprintf("%x", hash)
}

// Get recupera uma resposta do cache
func (c *LLMCache) Get(prompt string) (string, bool) {
	key := c.generateKey(prompt)

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		c.missCount++
		return "", false
	}

	// Verifica se a entrada expirou
	if time.Since(entry.CreatedAt) > c.ttl {
		delete(c.cache, key)
		c.missCount++
		return "", false
	}

	// Atualiza estatísticas de uso
	entry.LastUsed = time.Now()
	entry.UseCount++
	c.hitCount++

	return entry.Response, true
}

// Set armazena uma resposta no cache
func (c *LLMCache) Set(prompt, response string, metadata CacheMetadata) {
	key := c.generateKey(prompt)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Remove entradas antigas se o cache estiver cheio
	if len(c.cache) >= c.maxSize {
		c.evictOldestEntries(len(c.cache) - c.maxSize + 1)
	}

	c.cache[key] = &CacheEntry{
		Response:  response,
		Metadata:  metadata,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
		UseCount:  1,
	}
}

// FindSimilar procura por prompts similares no cache
func (c *LLMCache) FindSimilar(prompt string, threshold float64) (string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	normalizedPrompt := strings.ToLower(strings.TrimSpace(prompt))
	bestMatch := ""
	bestScore := 0.0

	for cachedPrompt, entry := range c.cache {
		// Reconstrói o prompt original a partir da chave (simplificado)
		similarity := c.calculateSimilarity(normalizedPrompt, cachedPrompt)

		if similarity > threshold && similarity > bestScore {
			bestScore = similarity
			bestMatch = entry.Response
		}
	}

	if bestMatch != "" {
		c.hitCount++
		return bestMatch, true
	}

	c.missCount++
	return "", false
}

// calculateSimilarity calcula similaridade entre dois prompts
func (c *LLMCache) calculateSimilarity(prompt1, prompt2 string) float64 {
	// Implementação simples de similaridade baseada em palavras comuns
	words1 := strings.Fields(prompt1)
	words2 := strings.Fields(prompt2)

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	commonWords := 0
	wordMap := make(map[string]bool)

	for _, word := range words1 {
		wordMap[word] = true
	}

	for _, word := range words2 {
		if wordMap[word] {
			commonWords++
		}
	}

	// Similaridade de Jaccard simplificada
	union := len(words1) + len(words2) - commonWords
	if union == 0 {
		return 1.0
	}

	return float64(commonWords) / float64(union)
}

// GetHitRatio retorna a taxa de acertos do cache
func (c *LLMCache) GetHitRatio() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	total := c.hitCount + c.missCount
	if total == 0 {
		return 0.0
	}

	return float64(c.hitCount) / float64(total)
}

// GetStats retorna estatísticas do cache
func (c *LLMCache) GetStats() CacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return CacheStats{
		Size:      len(c.cache),
		MaxSize:   c.maxSize,
		HitCount:  c.hitCount,
		MissCount: c.missCount,
		HitRatio:  c.GetHitRatio(),
	}
}

// CacheStats contém estatísticas do cache
type CacheStats struct {
	Size      int
	MaxSize   int
	HitCount  int64
	MissCount int64
	HitRatio  float64
}

// Clear limpa todo o cache
func (c *LLMCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]*CacheEntry)
	c.hitCount = 0
	c.missCount = 0
}

// evictOldestEntries remove as entradas mais antigas
func (c *LLMCache) evictOldestEntries(count int) {
	if count <= 0 {
		return
	}

	// Coleta todas as entradas com seus timestamps
	type entryWithKey struct {
		key   string
		entry *CacheEntry
	}

	entries := make([]entryWithKey, 0, len(c.cache))
	for key, entry := range c.cache {
		entries = append(entries, entryWithKey{key, entry})
	}

	// Ordena por último uso (mais antigo primeiro)
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].entry.LastUsed.After(entries[j].entry.LastUsed) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	// Remove as entradas mais antigas
	toRemove := count
	if toRemove > len(entries) {
		toRemove = len(entries)
	}

	for i := 0; i < toRemove; i++ {
		delete(c.cache, entries[i].key)
	}
}

// cleanupExpiredEntries limpa periodicamente entradas expiradas
func (c *LLMCache) cleanupExpiredEntries() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()

		for key, entry := range c.cache {
			if now.Sub(entry.CreatedAt) > c.ttl {
				delete(c.cache, key)
			}
		}
		c.mutex.Unlock()
	}
}
