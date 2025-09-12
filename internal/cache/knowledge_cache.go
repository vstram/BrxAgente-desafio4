package cache

import (
	"crypto/sha256"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// KnowledgeCache gerencia o cache de consultas frequentes para otimizar performance
type KnowledgeCache struct {
	cache      map[string]*KnowledgeCacheEntry
	lruOrder   []string // Ordem LRU para eviction
	mutex      sync.RWMutex
	config     CacheConfig
	metrics    KnowledgeCacheMetrics
	logger     *log.Logger
	enabled    bool
}

// KnowledgeCacheEntry representa uma entrada no cache de conhecimento
type KnowledgeCacheEntry struct {
	Response     string    `json:"response"`
	Timestamp    time.Time `json:"timestamp"`
	LastAccessed time.Time `json:"last_accessed"`
	HitCount     int64     `json:"hit_count"`
	Confidence   float64   `json:"confidence"`
	QueryHash    string    `json:"query_hash"`
	OriginalQuery string   `json:"original_query"`
}

// CacheConfig define a configuração do cache
type CacheConfig struct {
	MaxEntries int           `json:"max_entries"`
	TTL        time.Duration `json:"ttl"`
	Enabled    bool          `json:"enabled"`
	CleanupInterval time.Duration `json:"cleanup_interval"`
}

// KnowledgeCacheMetrics coleta métricas de performance do cache
type KnowledgeCacheMetrics struct {
	HitCount         int64         `json:"hit_count"`
	MissCount        int64         `json:"miss_count"`
	TotalRequests    int64         `json:"total_requests"`
	EvictionCount    int64         `json:"eviction_count"`
	ExpiredCount     int64         `json:"expired_count"`
	AverageHitTime   time.Duration `json:"average_hit_time"`
	AverageMissTime  time.Duration `json:"average_miss_time"`
	CacheSize        int           `json:"cache_size"`
	LastCleanup      time.Time     `json:"last_cleanup"`
	mutex            sync.RWMutex
}

// DefaultCacheConfig retorna configuração padrão do cache
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		MaxEntries:      1000,                // Máximo 1000 entradas
		TTL:             24 * time.Hour,      // 24 horas de TTL
		Enabled:         true,                // Cache habilitado por padrão
		CleanupInterval: 1 * time.Hour,       // Limpeza a cada hora
	}
}

// NewKnowledgeCache cria uma nova instância do cache
func NewKnowledgeCache(config CacheConfig) *KnowledgeCache {
	if config.MaxEntries == 0 {
		config = DefaultCacheConfig()
	}

	cache := &KnowledgeCache{
		cache:    make(map[string]*KnowledgeCacheEntry),
		lruOrder: make([]string, 0),
		config:   config,
		metrics:  KnowledgeCacheMetrics{LastCleanup: time.Now()},
		logger:   log.Default(),
		enabled:  config.Enabled,
	}

	// Iniciar limpeza automática
	if config.CleanupInterval > 0 {
		go cache.startCleanupRoutine()
	}

	cache.logger.Printf("KnowledgeCache inicializado: MaxEntries=%d, TTL=%s, Enabled=%t", 
		config.MaxEntries, config.TTL, config.Enabled)

	return cache
}

// normalizeQuestion normaliza uma pergunta para melhor matching no cache
func (kc *KnowledgeCache) normalizeQuestion(question string) string {
	if question == "" {
		return ""
	}

	// Converter para minúsculas
	normalized := strings.ToLower(question)

	// Remover acentos usando unicode normalization
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalized, _, _ = transform.String(t, normalized)

	// Remover pontuação desnecessária (incluindo interrogações múltiplas)
	re := regexp.MustCompile(`[^\w\s]`)
	normalized = re.ReplaceAllString(normalized, " ")

	// Padronizar espaços múltiplos
	spaceRe := regexp.MustCompile(`\s+`)
	normalized = spaceRe.ReplaceAllString(normalized, " ")

	// Trim espaços
	normalized = strings.TrimSpace(normalized)

	// Ordenar palavras para capturar perguntas similares com ordem diferente
	// Exceto palavras interrogativas que devem ficar no início
	words := strings.Fields(normalized)
	if len(words) > 1 {
		interrogatives := []string{"como", "qual", "quando", "onde", "por", "que", "quem", "quantos", "porque"}
		
		var interrogativeWords []string
		var otherWords []string
		
		for _, word := range words {
			isInterrogative := false
			for _, inter := range interrogatives {
				if strings.HasPrefix(word, inter) {
					isInterrogative = true
					break
				}
			}
			
			if isInterrogative {
				interrogativeWords = append(interrogativeWords, word)
			} else {
				otherWords = append(otherWords, word)
			}
		}
		
		// Ordenar apenas as palavras não-interrogativas
		sort.Strings(otherWords)
		
		// Reagrupar: interrogativas primeiro, depois outras ordenadas
		result := append(interrogativeWords, otherWords...)
		normalized = strings.Join(result, " ")
	}

	return normalized
}

// generateCacheKey gera uma chave única para o cache baseada na pergunta normalizada
func (kc *KnowledgeCache) generateCacheKey(question string) string {
	normalized := kc.normalizeQuestion(question)
	hash := sha256.Sum256([]byte(normalized))
	return fmt.Sprintf("%x", hash)
}

// Get recupera uma entrada do cache se existir e não tiver expirado
func (kc *KnowledgeCache) Get(question string) *KnowledgeCacheEntry {
	if !kc.enabled {
		return nil
	}

	start := time.Now()
	kc.mutex.Lock()
	defer kc.mutex.Unlock()

	// Atualizar métricas
	kc.metrics.mutex.Lock()
	kc.metrics.TotalRequests++
	kc.metrics.mutex.Unlock()

	key := kc.generateCacheKey(question)
	entry, exists := kc.cache[key]

	if !exists {
		kc.recordMiss(time.Since(start))
		return nil
	}

	// Verificar TTL
	if time.Since(entry.Timestamp) > kc.config.TTL {
		delete(kc.cache, key)
		kc.removeLRUEntry(key)
		kc.recordExpired()
		kc.recordMiss(time.Since(start))
		return nil
	}

	// Atualizar estatísticas da entrada
	entry.HitCount++
	entry.LastAccessed = time.Now()

	// Atualizar posição LRU
	kc.updateLRU(key)

	kc.recordHit(time.Since(start))
	
	kc.logger.Printf("Cache HIT: %s (hits: %d)", 
		kc.truncateString(question, 50), entry.HitCount)

	return entry
}

// Set armazena uma nova entrada no cache
func (kc *KnowledgeCache) Set(question, response string, confidence float64) {
	if !kc.enabled || question == "" || response == "" {
		return
	}

	kc.mutex.Lock()
	defer kc.mutex.Unlock()

	key := kc.generateCacheKey(question)

	// Verificar se precisa fazer eviction
	if len(kc.cache) >= kc.config.MaxEntries && kc.cache[key] == nil {
		kc.evictLRU()
	}

	entry := &KnowledgeCacheEntry{
		Response:      response,
		Timestamp:     time.Now(),
		LastAccessed:  time.Now(),
		HitCount:      0,
		Confidence:    confidence,
		QueryHash:     key,
		OriginalQuery: question,
	}

	kc.cache[key] = entry
	kc.updateLRU(key)

	kc.logger.Printf("Cache SET: %s", kc.truncateString(question, 50))
}

// updateLRU atualiza a ordem LRU movendo a chave para o final
func (kc *KnowledgeCache) updateLRU(key string) {
	// Remover da posição atual se existir
	for i, k := range kc.lruOrder {
		if k == key {
			kc.lruOrder = append(kc.lruOrder[:i], kc.lruOrder[i+1:]...)
			break
		}
	}
	
	// Adicionar no final (mais recentemente usado)
	kc.lruOrder = append(kc.lruOrder, key)
}

// removeLRUEntry remove uma chave da ordem LRU
func (kc *KnowledgeCache) removeLRUEntry(key string) {
	for i, k := range kc.lruOrder {
		if k == key {
			kc.lruOrder = append(kc.lruOrder[:i], kc.lruOrder[i+1:]...)
			break
		}
	}
}

// evictLRU remove a entrada menos recentemente usada
func (kc *KnowledgeCache) evictLRU() {
	if len(kc.lruOrder) == 0 {
		return
	}

	// Primeira entrada é a menos recentemente usada
	oldestKey := kc.lruOrder[0]
	delete(kc.cache, oldestKey)
	kc.lruOrder = kc.lruOrder[1:]

	kc.metrics.mutex.Lock()
	kc.metrics.EvictionCount++
	kc.metrics.mutex.Unlock()

	kc.logger.Printf("Cache EVICT: %s (LRU)", oldestKey[:16]+"...")
}

// recordHit registra um cache hit nas métricas
func (kc *KnowledgeCache) recordHit(duration time.Duration) {
	kc.metrics.mutex.Lock()
	kc.metrics.HitCount++
	
	// Calcular média móvel do tempo de hit
	if kc.metrics.HitCount == 1 {
		kc.metrics.AverageHitTime = duration
	} else {
		kc.metrics.AverageHitTime = time.Duration(
			(int64(kc.metrics.AverageHitTime)*(kc.metrics.HitCount-1) + int64(duration)) / kc.metrics.HitCount,
		)
	}
	kc.metrics.mutex.Unlock()
}

// recordMiss registra um cache miss nas métricas
func (kc *KnowledgeCache) recordMiss(duration time.Duration) {
	kc.metrics.mutex.Lock()
	kc.metrics.MissCount++
	
	// Calcular média móvel do tempo de miss
	if kc.metrics.MissCount == 1 {
		kc.metrics.AverageMissTime = duration
	} else {
		kc.metrics.AverageMissTime = time.Duration(
			(int64(kc.metrics.AverageMissTime)*(kc.metrics.MissCount-1) + int64(duration)) / kc.metrics.MissCount,
		)
	}
	kc.metrics.mutex.Unlock()
}

// recordExpired registra uma entrada expirada nas métricas
func (kc *KnowledgeCache) recordExpired() {
	kc.metrics.mutex.Lock()
	kc.metrics.ExpiredCount++
	kc.metrics.mutex.Unlock()
}

// GetMetrics retorna as métricas atuais do cache
func (kc *KnowledgeCache) GetMetrics() KnowledgeCacheMetrics {
	kc.metrics.mutex.RLock()
	kc.mutex.RLock()
	
	metrics := kc.metrics
	metrics.CacheSize = len(kc.cache)
	
	kc.mutex.RUnlock()
	kc.metrics.mutex.RUnlock()
	
	return metrics
}

// GetHitRate calcula a taxa de acertos do cache
func (kc *KnowledgeCache) GetHitRate() float64 {
	metrics := kc.GetMetrics()
	if metrics.TotalRequests == 0 {
		return 0.0
	}
	return float64(metrics.HitCount) / float64(metrics.TotalRequests)
}

// Clear limpa todo o cache
func (kc *KnowledgeCache) Clear() {
	kc.mutex.Lock()
	defer kc.mutex.Unlock()

	kc.cache = make(map[string]*KnowledgeCacheEntry)
	kc.lruOrder = make([]string, 0)
	
	kc.logger.Println("Cache limpo completamente")
}

// startCleanupRoutine inicia rotina de limpeza automática
func (kc *KnowledgeCache) startCleanupRoutine() {
	ticker := time.NewTicker(kc.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		kc.cleanup()
	}
}

// cleanup remove entradas expiradas do cache
func (kc *KnowledgeCache) cleanup() {
	kc.mutex.Lock()
	defer kc.mutex.Unlock()

	now := time.Now()
	expiredKeys := make([]string, 0)

	for key, entry := range kc.cache {
		if now.Sub(entry.Timestamp) > kc.config.TTL {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		delete(kc.cache, key)
		kc.removeLRUEntry(key)
	}

	if len(expiredKeys) > 0 {
		kc.metrics.mutex.Lock()
		kc.metrics.ExpiredCount += int64(len(expiredKeys))
		kc.metrics.LastCleanup = now
		kc.metrics.mutex.Unlock()
		
		kc.logger.Printf("Cache cleanup: %d entradas expiradas removidas", len(expiredKeys))
	}
}

// Enable habilita o cache
func (kc *KnowledgeCache) Enable() {
	kc.enabled = true
	kc.logger.Println("Cache habilitado")
}

// Disable desabilita o cache
func (kc *KnowledgeCache) Disable() {
	kc.enabled = false
	kc.logger.Println("Cache desabilitado")
}

// IsEnabled retorna se o cache está habilitado
func (kc *KnowledgeCache) IsEnabled() bool {
	return kc.enabled
}

// GetConfig retorna a configuração atual do cache
func (kc *KnowledgeCache) GetConfig() CacheConfig {
	return kc.config
}

// truncateString trunca uma string para exibição em logs
func (kc *KnowledgeCache) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// GetCacheEntries retorna todas as entradas do cache (para debugging)
func (kc *KnowledgeCache) GetCacheEntries() map[string]*KnowledgeCacheEntry {
	kc.mutex.RLock()
	defer kc.mutex.RUnlock()

	entries := make(map[string]*KnowledgeCacheEntry)
	for k, v := range kc.cache {
		entryCopy := *v
		entries[k] = &entryCopy
	}

	return entries
}