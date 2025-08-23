package cache

import (
	"fmt"
	"sync"
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

// DataCacheEntry representa uma entrada no cache de dados
type DataCacheEntry struct {
	Data      interface{}
	CreatedAt time.Time
	LastUsed  time.Time
	Size      int64 // Tamanho estimado em bytes
}

// DataCache implementa cache para dados processados
type DataCache struct {
	cache    map[string]*DataCacheEntry
	mutex    sync.RWMutex
	ttl      time.Duration
	maxSize  int64 // Tamanho máximo em bytes
	currSize int64 // Tamanho atual em bytes
}

// NewDataCache cria uma nova instância de cache de dados
func NewDataCache(maxSizeBytes int64, ttl time.Duration) *DataCache {
	cache := &DataCache{
		cache:   make(map[string]*DataCacheEntry),
		ttl:     ttl,
		maxSize: maxSizeBytes,
	}

	// Inicia limpeza periódica
	go cache.cleanupExpired()

	return cache
}

// GetColaboradores recupera colaboradores do cache
func (c *DataCache) GetColaboradores(key string) (map[string]*modelo.Colaborador, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	if time.Since(entry.CreatedAt) > c.ttl {
		delete(c.cache, key)
		c.currSize -= entry.Size
		return nil, false
	}

	entry.LastUsed = time.Now()

	if colaboradores, ok := entry.Data.(map[string]*modelo.Colaborador); ok {
		return colaboradores, true
	}

	return nil, false
}

// SetColaboradores armazena colaboradores no cache
func (c *DataCache) SetColaboradores(key string, colaboradores map[string]*modelo.Colaborador) {
	size := c.estimateColaboradoresSize(colaboradores)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Remove entrada existente se houver
	if existing, exists := c.cache[key]; exists {
		c.currSize -= existing.Size
	}

	// Limpa espaço se necessário
	for c.currSize+size > c.maxSize && len(c.cache) > 0 {
		c.evictLRU()
	}

	c.cache[key] = &DataCacheEntry{
		Data:      colaboradores,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
		Size:      size,
	}

	c.currSize += size
}

// GetProcessedData recupera dados processados genéricos
func (c *DataCache) GetProcessedData(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	if time.Since(entry.CreatedAt) > c.ttl {
		delete(c.cache, key)
		c.currSize -= entry.Size
		return nil, false
	}

	entry.LastUsed = time.Now()
	return entry.Data, true
}

// SetProcessedData armazena dados processados genéricos
func (c *DataCache) SetProcessedData(key string, data interface{}, estimatedSize int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Remove entrada existente se houver
	if existing, exists := c.cache[key]; exists {
		c.currSize -= existing.Size
	}

	// Limpa espaço se necessário
	for c.currSize+estimatedSize > c.maxSize && len(c.cache) > 0 {
		c.evictLRU()
	}

	c.cache[key] = &DataCacheEntry{
		Data:      data,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
		Size:      estimatedSize,
	}

	c.currSize += estimatedSize
}

// InvalidatePattern invalida todas as chaves que correspondem ao padrão
func (c *DataCache) InvalidatePattern(pattern string) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	removed := 0
	for key, entry := range c.cache {
		if c.matchPattern(key, pattern) {
			c.currSize -= entry.Size
			delete(c.cache, key)
			removed++
		}
	}

	return removed
}

// GetStats retorna estatísticas do cache de dados
func (c *DataCache) GetStats() DataCacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return DataCacheStats{
		EntryCount:    len(c.cache),
		CurrentSize:   c.currSize,
		MaxSize:       c.maxSize,
		UsagePercent:  float64(c.currSize) / float64(c.maxSize) * 100,
		OldestEntry:   c.getOldestEntryAge(),
		NewestEntry:   c.getNewestEntryAge(),
	}
}

// DataCacheStats contém estatísticas do cache de dados
type DataCacheStats struct {
	EntryCount   int
	CurrentSize  int64
	MaxSize      int64
	UsagePercent float64
	OldestEntry  time.Duration
	NewestEntry  time.Duration
}

// Clear limpa todo o cache de dados
func (c *DataCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]*DataCacheEntry)
	c.currSize = 0
}

// estimateColaboradoresSize estima o tamanho de um map de colaboradores
func (c *DataCache) estimateColaboradoresSize(colaboradores map[string]*modelo.Colaborador) int64 {
	// Estimativa grosseira: ~200 bytes por colaborador
	const bytesPerColaborador = 200
	return int64(len(colaboradores) * bytesPerColaborador)
}

// evictLRU remove a entrada menos recentemente usada
func (c *DataCache) evictLRU() {
	if len(c.cache) == 0 {
		return
	}

	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, entry := range c.cache {
		if first || entry.LastUsed.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.LastUsed
			first = false
		}
	}

	if entry, exists := c.cache[oldestKey]; exists {
		c.currSize -= entry.Size
		delete(c.cache, oldestKey)
	}
}

// matchPattern verifica se uma chave corresponde ao padrão (implementação simples)
func (c *DataCache) matchPattern(key, pattern string) bool {
	// Implementação simples de matching - pode ser expandida para regex
	return key == pattern || (pattern[len(pattern)-1] == '*' && 
		len(key) >= len(pattern)-1 &&
		key[:len(pattern)-1] == pattern[:len(pattern)-1])
}

// cleanupExpired limpa periodicamente entradas expiradas
func (c *DataCache) cleanupExpired() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()

		for key, entry := range c.cache {
			if now.Sub(entry.CreatedAt) > c.ttl {
				c.currSize -= entry.Size
				delete(c.cache, key)
			}
		}
		c.mutex.Unlock()
	}
}

// getOldestEntryAge retorna a idade da entrada mais antiga
func (c *DataCache) getOldestEntryAge() time.Duration {
	if len(c.cache) == 0 {
		return 0
	}

	var oldest time.Time
	first := true

	for _, entry := range c.cache {
		if first || entry.CreatedAt.Before(oldest) {
			oldest = entry.CreatedAt
			first = false
		}
	}

	return time.Since(oldest)
}

// getNewestEntryAge retorna a idade da entrada mais nova
func (c *DataCache) getNewestEntryAge() time.Duration {
	if len(c.cache) == 0 {
		return 0
	}

	var newest time.Time
	first := true

	for _, entry := range c.cache {
		if first || entry.CreatedAt.After(newest) {
			newest = entry.CreatedAt
			first = false
		}
	}

	return time.Since(newest)
}

// GenerateKey gera uma chave de cache baseada em múltiplos parâmetros
func GenerateKey(prefix string, params ...interface{}) string {
	key := prefix
	for _, param := range params {
		key += fmt.Sprintf(":%v", param)
	}
	return key
}