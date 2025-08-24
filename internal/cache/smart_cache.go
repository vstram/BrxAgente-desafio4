package cache

import (
	"sync"
	"time"
)

// SmartCache combina LLM cache e data cache em uma interface unificada
type SmartCache struct {
	llmCache  *LLMCache
	dataCache *DataCache
	mutex     sync.RWMutex
}

// SmartCacheConfig configuração para o smart cache
type SmartCacheConfig struct {
	LLMMaxSize       int
	LLMTTLHours      int
	DataMaxSizeBytes int64
	DataTTLHours     int
}

// NewSmartCache cria uma nova instância de smart cache
func NewSmartCache(config SmartCacheConfig) *SmartCache {
	llmTTL := time.Duration(config.LLMTTLHours) * time.Hour
	dataTTL := time.Duration(config.DataTTLHours) * time.Hour

	return &SmartCache{
		llmCache:  NewLLMCache(config.LLMMaxSize, llmTTL),
		dataCache: NewDataCache(config.DataMaxSizeBytes, dataTTL),
	}
}

// GetLLMResponse recupera resposta do LLM do cache
func (sc *SmartCache) GetLLMResponse(prompt string) (string, bool) {
	return sc.llmCache.Get(prompt)
}

// SetLLMResponse armazena resposta do LLM no cache
func (sc *SmartCache) SetLLMResponse(prompt, response string, metadata CacheMetadata) {
	sc.llmCache.Set(prompt, response, metadata)
}

// FindSimilarLLMResponse procura resposta similar no cache LLM
func (sc *SmartCache) FindSimilarLLMResponse(prompt string, threshold float64) (string, bool) {
	return sc.llmCache.FindSimilar(prompt, threshold)
}

// GetProcessedData recupera dados processados do cache
func (sc *SmartCache) GetProcessedData(key string) (interface{}, bool) {
	return sc.dataCache.GetProcessedData(key)
}

// SetProcessedData armazena dados processados no cache
func (sc *SmartCache) SetProcessedData(key string, data interface{}, estimatedSize int64) {
	sc.dataCache.SetProcessedData(key, data, estimatedSize)
}

// InvalidateData invalida dados do cache baseado em padrão
func (sc *SmartCache) InvalidateData(pattern string) int {
	return sc.dataCache.InvalidatePattern(pattern)
}

// GetCombinedStats retorna estatísticas combinadas dos caches
func (sc *SmartCache) GetCombinedStats() CombinedCacheStats {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	llmStats := sc.llmCache.GetStats()
	dataStats := sc.dataCache.GetStats()

	return CombinedCacheStats{
		LLMStats:  llmStats,
		DataStats: dataStats,
		Timestamp: time.Now(),
	}
}

// CombinedCacheStats estatísticas combinadas dos caches
type CombinedCacheStats struct {
	LLMStats  CacheStats
	DataStats DataCacheStats
	Timestamp time.Time
}

// ClearAll limpa todos os caches
func (sc *SmartCache) ClearAll() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.llmCache.Clear()
	sc.dataCache.Clear()
}

// OptimizeCaches executa otimizações nos caches
func (sc *SmartCache) OptimizeCaches() CacheOptimizationResult {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	result := CacheOptimizationResult{
		StartTime: time.Now(),
	}

	// Analisa hit ratio do LLM cache
	llmStats := sc.llmCache.GetStats()
	if llmStats.HitRatio < 0.5 && llmStats.Size > llmStats.MaxSize/2 {
		// Se hit ratio é baixo mas cache está meio cheio, limpa entradas antigas
		oldSize := llmStats.Size
		sc.llmCache.evictOldestEntries(llmStats.MaxSize / 4)
		newStats := sc.llmCache.GetStats()
		result.LLMEntriesRemoved = oldSize - newStats.Size
	}

	// Analisa uso do data cache
	dataStats := sc.dataCache.GetStats()
	if dataStats.UsagePercent > 80 {
		// Se uso está alto, força limpeza de entradas antigas
		oldEntryCount := dataStats.EntryCount
		sc.dataCache.evictLRU()
		newStats := sc.dataCache.GetStats()
		result.DataEntriesRemoved = oldEntryCount - newStats.EntryCount
	}

	result.EndTime = time.Now()
	result.Success = true

	return result
}

// CacheOptimizationResult resultado da otimização de cache
type CacheOptimizationResult struct {
	Success            bool
	StartTime          time.Time
	EndTime            time.Time
	LLMEntriesRemoved  int
	DataEntriesRemoved int
}

// WarmupCache pré-aquece o cache com dados frequentemente usados
func (sc *SmartCache) WarmupCache(warmupData []WarmupEntry) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	for _, entry := range warmupData {
		switch entry.Type {
		case "llm":
			if prompt, ok := entry.Key.(string); ok {
				if response, ok := entry.Value.(string); ok {
					metadata := CacheMetadata{
						TokensUsed: 100, // Valor padrão para warmup
						Duration:   time.Millisecond * 100,
						Model:      "warmup",
						Quality:    1.0,
					}
					sc.llmCache.Set(prompt, response, metadata)
				}
			}
		case "data":
			if key, ok := entry.Key.(string); ok {
				sc.dataCache.SetProcessedData(key, entry.Value, entry.EstimatedSize)
			}
		}
	}
}

// WarmupEntry entrada para pré-aquecimento do cache
type WarmupEntry struct {
	Type          string      // "llm" ou "data"
	Key           interface{} // Chave do cache
	Value         interface{} // Valor a ser armazenado
	EstimatedSize int64       // Tamanho estimado (para data cache)
}

// GetCacheHealth verifica a saúde geral dos caches
func (sc *SmartCache) GetCacheHealth() CacheHealthReport {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	llmStats := sc.llmCache.GetStats()
	dataStats := sc.dataCache.GetStats()

	health := CacheHealthReport{
		Timestamp: time.Now(),
		Overall:   "healthy",
	}

	// Avalia saúde do LLM cache
	if llmStats.HitRatio < 0.3 {
		health.LLMHealth = "poor"
		health.Issues = append(health.Issues, "LLM cache hit ratio muito baixo")
	} else if llmStats.HitRatio < 0.7 {
		health.LLMHealth = "fair"
	} else {
		health.LLMHealth = "good"
	}

	// Avalia saúde do data cache
	if dataStats.UsagePercent > 90 {
		health.DataHealth = "poor"
		health.Issues = append(health.Issues, "Data cache quase cheio")
	} else if dataStats.UsagePercent > 75 {
		health.DataHealth = "fair"
	} else {
		health.DataHealth = "good"
	}

	// Determina saúde geral
	if health.LLMHealth == "poor" || health.DataHealth == "poor" {
		health.Overall = "poor"
	} else if health.LLMHealth == "fair" || health.DataHealth == "fair" {
		health.Overall = "fair"
	}

	health.Recommendations = sc.generateRecommendations(llmStats, dataStats)

	return health
}

// CacheHealthReport relatório de saúde dos caches
type CacheHealthReport struct {
	Timestamp       time.Time
	Overall         string // "healthy", "fair", "poor"
	LLMHealth       string
	DataHealth      string
	Issues          []string
	Recommendations []string
}

// generateRecommendations gera recomendações baseadas nas estatísticas
func (sc *SmartCache) generateRecommendations(llmStats CacheStats, dataStats DataCacheStats) []string {
	var recommendations []string

	if llmStats.HitRatio < 0.5 {
		recommendations = append(recommendations, "Considere aumentar o TTL do LLM cache")
		recommendations = append(recommendations, "Revise os prompts para maior consistência")
	}

	if dataStats.UsagePercent > 80 {
		recommendations = append(recommendations, "Aumente o tamanho máximo do data cache")
		recommendations = append(recommendations, "Implemente invalidação mais agressiva")
	}

	if llmStats.Size < llmStats.MaxSize/4 && dataStats.EntryCount < 10 {
		recommendations = append(recommendations, "Considere pré-aquecer os caches com dados comuns")
	}

	return recommendations
}
