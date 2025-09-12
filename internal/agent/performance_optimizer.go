package agent

import (
	"crypto/md5"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"BrxAgente-desafio4/internal/knowledge"
	"BrxAgente-desafio4/internal/modelo"
)

// PerformanceOptimizer gerencia otimizações de performance do sistema
type PerformanceOptimizer struct {
	// Context formatting cache
	contextCache     map[string]*CachedContext
	contextCacheMu   sync.RWMutex
	
	// Knowledge base lazy loader
	lazyKnowledgeBase *LazyKnowledgeBase
	
	// Cleanup configuration
	cleanupInterval  time.Duration
	cacheMaxAge      time.Duration
	cacheMaxEntries  int
	
	// Statistics
	stats            *PerformanceStats
	statsMu          sync.RWMutex
	
	// Lifecycle
	stopCleanup      chan bool
	logger           *log.Logger
}

// CachedContext representa um contexto formatado em cache
type CachedContext struct {
	FormattedData string
	Hash          string
	CreatedAt     time.Time
	LastAccessed  time.Time
	AccessCount   int64
}

// LazyKnowledgeBase implementa carregamento lazy da base de conhecimento
type LazyKnowledgeBase struct {
	dataDir       string
	loaded        bool
	loadMutex     sync.Mutex
	knowledge     *knowledge.KnowledgeBaseManager
	indexCache    map[string][]knowledge.SearchResult
	indexCacheMu  sync.RWMutex
	lastLoadTime  time.Time
}

// PerformanceStats coleta estatísticas de performance
type PerformanceStats struct {
	ContextCacheHits      int64
	ContextCacheMisses    int64
	KnowledgeLoadCount    int64
	IndexCacheHits        int64
	IndexCacheMisses      int64
	AverageResponseTime   time.Duration
	TotalQueries          int64
	MemoryUsage           int64
	CleanupCycles         int64
	LastOptimizationTime  time.Time
}

// NewPerformanceOptimizer cria uma nova instância do otimizador
func NewPerformanceOptimizer() *PerformanceOptimizer {
	po := &PerformanceOptimizer{
		contextCache:      make(map[string]*CachedContext),
		cleanupInterval:   15 * time.Minute,
		cacheMaxAge:       30 * time.Minute,
		cacheMaxEntries:   1000,
		stats:             &PerformanceStats{},
		stopCleanup:       make(chan bool, 1),
		logger:            log.Default(),
	}
	
	po.lazyKnowledgeBase = &LazyKnowledgeBase{
		indexCache: make(map[string][]knowledge.SearchResult),
	}
	
	// Iniciar limpeza automática
	go po.startCacheCleanup()
	
	return po
}

// SetKnowledgeDataDir configura o diretório de dados da base de conhecimento
func (po *PerformanceOptimizer) SetKnowledgeDataDir(dataDir string) {
	po.lazyKnowledgeBase.dataDir = dataDir
}

// FormatContextWithCache formata dados de contexto usando cache
func (po *PerformanceOptimizer) FormatContextWithCache(data map[string]*modelo.Colaborador, maxDetailed int) string {
	startTime := time.Now()
	defer func() {
		po.updateStats(time.Since(startTime))
	}()
	
	// Gerar hash dos dados
	hash := po.generateContextHash(data, maxDetailed)
	
	// Verificar cache primeiro
	po.contextCacheMu.RLock()
	if cached, exists := po.contextCache[hash]; exists {
		cached.LastAccessed = time.Now()
		cached.AccessCount++
		po.contextCacheMu.RUnlock()
		
		po.incrementCacheHit(true)
		return cached.FormattedData
	}
	po.contextCacheMu.RUnlock()
	
	po.incrementCacheHit(false)
	
	// Formatar dados
	formatted := po.formatContextData(data, maxDetailed)
	
	// Armazenar no cache
	po.setCachedContext(hash, formatted)
	
	return formatted
}

// EnsureKnowledgeBaseLoaded garante que a base de conhecimento está carregada
func (po *PerformanceOptimizer) EnsureKnowledgeBaseLoaded() error {
	lkb := po.lazyKnowledgeBase
	
	if lkb.loaded {
		return nil
	}
	
	lkb.loadMutex.Lock()
	defer lkb.loadMutex.Unlock()
	
	// Double-check pattern
	if lkb.loaded {
		return nil
	}
	
	startTime := time.Now()
	
	// Carregar base de conhecimento
	kb := knowledge.NewKnowledgeBaseManager()
	err := kb.LoadFromFiles(lkb.dataDir)
	if err != nil {
		return fmt.Errorf("erro ao carregar base de conhecimento: %w", err)
	}
	
	lkb.knowledge = kb
	lkb.loaded = true
	lkb.lastLoadTime = time.Now()
	
	po.statsMu.Lock()
	po.stats.KnowledgeLoadCount++
	po.statsMu.Unlock()
	
	po.logger.Printf("Base de conhecimento carregada em %v", time.Since(startTime))
	
	return nil
}

// SearchKnowledgeWithCache busca na base de conhecimento usando cache
func (po *PerformanceOptimizer) SearchKnowledgeWithCache(query string, limit int) ([]knowledge.SearchResult, error) {
	// Garantir que a base está carregada
	if err := po.EnsureKnowledgeBaseLoaded(); err != nil {
		return nil, err
	}
	
	// Verificar cache de índice
	cacheKey := fmt.Sprintf("%s:%d", query, limit)
	
	po.lazyKnowledgeBase.indexCacheMu.RLock()
	if cached, exists := po.lazyKnowledgeBase.indexCache[cacheKey]; exists {
		po.lazyKnowledgeBase.indexCacheMu.RUnlock()
		
		po.statsMu.Lock()
		po.stats.IndexCacheHits++
		po.statsMu.Unlock()
		
		return cached, nil
	}
	po.lazyKnowledgeBase.indexCacheMu.RUnlock()
	
	po.statsMu.Lock()
	po.stats.IndexCacheMisses++
	po.statsMu.Unlock()
	
	// Buscar na base de conhecimento
	results, err := po.lazyKnowledgeBase.knowledge.Search(query, limit)
	if err != nil {
		return nil, err
	}
	
	// Cachear resultados
	po.lazyKnowledgeBase.indexCacheMu.Lock()
	po.lazyKnowledgeBase.indexCache[cacheKey] = results
	po.lazyKnowledgeBase.indexCacheMu.Unlock()
	
	return results, nil
}

// GetKnowledgeBase retorna a base de conhecimento (carregando se necessário)
func (po *PerformanceOptimizer) GetKnowledgeBase() (*knowledge.KnowledgeBaseManager, error) {
	if err := po.EnsureKnowledgeBaseLoaded(); err != nil {
		return nil, err
	}
	return po.lazyKnowledgeBase.knowledge, nil
}

// generateContextHash gera hash para dados de contexto
func (po *PerformanceOptimizer) generateContextHash(data map[string]*modelo.Colaborador, maxDetailed int) string {
	h := md5.New()
	
	// Adicionar parâmetros ao hash
	h.Write([]byte(fmt.Sprintf("maxDetailed:%d", maxDetailed)))
	
	// Adicionar dados de forma determinística
	for matricula, colaborador := range data {
		h.Write([]byte(fmt.Sprintf("m:%s:e:%s:s:%s:v:%.2f:d:%d",
			matricula,
			colaborador.Empresa,
			colaborador.Sindicato,
			colaborador.ValorTotalVR,
			colaborador.DiasUteisEfetivos)))
	}
	
	return fmt.Sprintf("%x", h.Sum(nil))
}

// formatContextData formata dados de contexto (versão otimizada)
func (po *PerformanceOptimizer) formatContextData(data map[string]*modelo.Colaborador, maxDetailed int) string {
	if len(data) == 0 {
		return ""
	}

	// Pre-allocate builder com tamanho estimado
	estimatedSize := len(data)*100 + 1000 // Estimativa conservadora
	var summary strings.Builder
	summary.Grow(estimatedSize)
	
	total := len(data)
	summary.WriteString(fmt.Sprintf("Dados de %d colaboradores disponíveis:\n", total))

	// Calculate totals in single pass
	var totalVR, totalEmpresa, totalColaborador float64
	var totalDiasUteis int
	sindicatos := make(map[string]int, 10)
	empresas := make(map[string]int, 10)
	
	// Use slice for deterministic iteration order
	matriculas := make([]string, 0, len(data))
	for matricula := range data {
		matriculas = append(matriculas, matricula)
	}
	
	// Sort for deterministic output
	for i := 0; i < len(matriculas)-1; i++ {
		for j := i + 1; j < len(matriculas); j++ {
			if matriculas[i] > matriculas[j] {
				matriculas[i], matriculas[j] = matriculas[j], matriculas[i]
			}
		}
	}
	
	for _, matricula := range matriculas {
		colaborador := data[matricula]
		totalVR += colaborador.ValorTotalVR
		totalEmpresa += colaborador.ValorEmpresa
		totalColaborador += colaborador.ValorColaborador
		totalDiasUteis += colaborador.DiasUteisEfetivos
		sindicatos[colaborador.Sindicato]++
		empresas[colaborador.Empresa]++
	}

	// Add summary statistics
	summary.WriteString("\n=== RESUMO GERAL ===\n")
	summary.WriteString(fmt.Sprintf("Total de colaboradores: %d\n", total))
	summary.WriteString(fmt.Sprintf("Valor total VR: R$ %.2f\n", totalVR))
	summary.WriteString(fmt.Sprintf("Valor total empresa: R$ %.2f\n", totalEmpresa))
	summary.WriteString(fmt.Sprintf("Valor total colaborador: R$ %.2f\n", totalColaborador))
	summary.WriteString(fmt.Sprintf("Total dias úteis: %d\n", totalDiasUteis))

	summary.WriteString("\nDistribuição por sindicato:\n")
	// Sort sindicatos for deterministic output
	sindicatoKeys := make([]string, 0, len(sindicatos))
	for sindicato := range sindicatos {
		sindicatoKeys = append(sindicatoKeys, sindicato)
	}
	for i := 0; i < len(sindicatoKeys)-1; i++ {
		for j := i + 1; j < len(sindicatoKeys); j++ {
			if sindicatoKeys[i] > sindicatoKeys[j] {
				sindicatoKeys[i], sindicatoKeys[j] = sindicatoKeys[j], sindicatoKeys[i]
			}
		}
	}
	for _, sindicato := range sindicatoKeys {
		summary.WriteString(fmt.Sprintf("  %s: %d colaboradores\n", sindicato, sindicatos[sindicato]))
	}

	summary.WriteString("\nDistribuição por empresa:\n")
	// Sort empresas for deterministic output
	empresaKeys := make([]string, 0, len(empresas))
	for empresa := range empresas {
		empresaKeys = append(empresaKeys, empresa)
	}
	for i := 0; i < len(empresaKeys)-1; i++ {
		for j := i + 1; j < len(empresaKeys); j++ {
			if empresaKeys[i] > empresaKeys[j] {
				empresaKeys[i], empresaKeys[j] = empresaKeys[j], empresaKeys[i]
			}
		}
	}
	for _, empresa := range empresaKeys {
		summary.WriteString(fmt.Sprintf("  %s: %d colaboradores\n", empresa, empresas[empresa]))
	}

	// Add detailed examples (optimized loop)
	if total > 0 && maxDetailed > 0 {
		detailedCount := maxDetailed
		if total < maxDetailed {
			detailedCount = total
		}
		summary.WriteString(fmt.Sprintf("\n=== EXEMPLOS DETALHADOS (primeiros %d) ===\n", detailedCount))
		
		for i, matricula := range matriculas {
			if i >= maxDetailed {
				summary.WriteString(fmt.Sprintf("... e mais %d colaboradores com dados detalhados disponíveis\n", total-maxDetailed))
				break
			}
			
			colaborador := data[matricula]
			summary.WriteString(fmt.Sprintf(
				"Colaborador %d: Matrícula=%s, Empresa=%s, Sindicato=%s, VR=R$%.2f, Dias=%d\n",
				i+1,
				colaborador.Matricula,
				colaborador.Empresa,
				colaborador.Sindicato,
				colaborador.ValorTotalVR,
				colaborador.DiasUteisEfetivos,
			))
		}
	}

	return summary.String()
}

// setCachedContext armazena contexto no cache
func (po *PerformanceOptimizer) setCachedContext(hash, formatted string) {
	po.contextCacheMu.Lock()
	defer po.contextCacheMu.Unlock()
	
	// Verificar limite de entradas
	if len(po.contextCache) >= po.cacheMaxEntries {
		po.evictOldestEntries(po.cacheMaxEntries / 4) // Remove 25%
	}
	
	po.contextCache[hash] = &CachedContext{
		FormattedData: formatted,
		Hash:          hash,
		CreatedAt:     time.Now(),
		LastAccessed:  time.Now(),
		AccessCount:   1,
	}
}

// evictOldestEntries remove as entradas mais antigas do cache
func (po *PerformanceOptimizer) evictOldestEntries(count int) {
	if len(po.contextCache) == 0 {
		return
	}
	
	// Encontrar entradas mais antigas
	type cacheEntry struct {
		hash       string
		lastAccess time.Time
	}
	
	entries := make([]cacheEntry, 0, len(po.contextCache))
	for hash, cached := range po.contextCache {
		entries = append(entries, cacheEntry{
			hash:       hash,
			lastAccess: cached.LastAccessed,
		})
	}
	
	// Ordenar por último acesso
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].lastAccess.After(entries[j].lastAccess) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
	
	// Remover as mais antigas
	removed := 0
	for _, entry := range entries {
		if removed >= count {
			break
		}
		delete(po.contextCache, entry.hash)
		removed++
	}
}

// incrementCacheHit atualiza estatísticas de cache
func (po *PerformanceOptimizer) incrementCacheHit(hit bool) {
	po.statsMu.Lock()
	defer po.statsMu.Unlock()
	
	if hit {
		po.stats.ContextCacheHits++
	} else {
		po.stats.ContextCacheMisses++
	}
}

// updateStats atualiza estatísticas gerais
func (po *PerformanceOptimizer) updateStats(responseTime time.Duration) {
	po.statsMu.Lock()
	defer po.statsMu.Unlock()
	
	po.stats.TotalQueries++
	
	// Calcular média móvel do tempo de resposta
	if po.stats.TotalQueries == 1 {
		po.stats.AverageResponseTime = responseTime
	} else {
		// Média ponderada com peso maior para valores recentes
		weight := float64(0.1)
		po.stats.AverageResponseTime = time.Duration(
			float64(po.stats.AverageResponseTime)*(1-weight) +
			float64(responseTime)*weight)
	}
}

// startCacheCleanup inicia rotina de limpeza de cache
func (po *PerformanceOptimizer) startCacheCleanup() {
	ticker := time.NewTicker(po.cleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			po.cleanupExpiredCache()
		case <-po.stopCleanup:
			return
		}
	}
}

// cleanupExpiredCache remove entradas expiradas do cache
func (po *PerformanceOptimizer) cleanupExpiredCache() {
	startTime := time.Now()
	
	po.contextCacheMu.Lock()
	defer po.contextCacheMu.Unlock()
	
	expiredCount := 0
	cutoff := time.Now().Add(-po.cacheMaxAge)
	
	// Remover entradas expiradas
	for hash, cached := range po.contextCache {
		if cached.LastAccessed.Before(cutoff) {
			delete(po.contextCache, hash)
			expiredCount++
		}
	}
	
	// Cleanup do cache de índice também
	po.lazyKnowledgeBase.indexCacheMu.Lock()
	// Cache de índice expira mais rapidamente (5 minutos)
	indexCutoff := time.Now().Add(-5 * time.Minute)
	indexExpired := 0
	
	// Como não temos timestamp no cache de índice, limpar tudo se passou muito tempo
	if po.lazyKnowledgeBase.lastLoadTime.Before(indexCutoff) {
		indexExpired = len(po.lazyKnowledgeBase.indexCache)
		po.lazyKnowledgeBase.indexCache = make(map[string][]knowledge.SearchResult)
	}
	po.lazyKnowledgeBase.indexCacheMu.Unlock()
	
	po.statsMu.Lock()
	po.stats.CleanupCycles++
	po.stats.LastOptimizationTime = time.Now()
	po.statsMu.Unlock()
	
	if expiredCount > 0 || indexExpired > 0 {
		po.logger.Printf("Cache cleanup: removed %d context entries and %d index entries in %v", 
			expiredCount, indexExpired, time.Since(startTime))
	}
}

// GetStats retorna estatísticas de performance
func (po *PerformanceOptimizer) GetStats() PerformanceStats {
	po.statsMu.RLock()
	defer po.statsMu.RUnlock()
	
	stats := *po.stats // Copy
	
	// Adicionar informações em tempo real
	po.contextCacheMu.RLock()
	stats.MemoryUsage = int64(len(po.contextCache))
	po.contextCacheMu.RUnlock()
	
	return stats
}

// GetCacheInfo retorna informações sobre o cache
func (po *PerformanceOptimizer) GetCacheInfo() map[string]interface{} {
	po.contextCacheMu.RLock()
	contextEntries := len(po.contextCache)
	po.contextCacheMu.RUnlock()
	
	po.lazyKnowledgeBase.indexCacheMu.RLock()
	indexEntries := len(po.lazyKnowledgeBase.indexCache)
	po.lazyKnowledgeBase.indexCacheMu.RUnlock()
	
	stats := po.GetStats()
	
	hitRate := float64(0)
	if stats.ContextCacheHits+stats.ContextCacheMisses > 0 {
		hitRate = float64(stats.ContextCacheHits) / float64(stats.ContextCacheHits+stats.ContextCacheMisses)
	}
	
	return map[string]interface{}{
		"context_cache_entries": contextEntries,
		"index_cache_entries":   indexEntries,
		"cache_hit_rate":        hitRate,
		"knowledge_loaded":      po.lazyKnowledgeBase.loaded,
		"total_queries":         stats.TotalQueries,
		"avg_response_time_ms":  float64(stats.AverageResponseTime.Nanoseconds()) / 1000000,
		"cleanup_cycles":        stats.CleanupCycles,
		"last_optimization":     stats.LastOptimizationTime.Format(time.RFC3339),
	}
}

// Stop para o otimizador e limpa recursos
func (po *PerformanceOptimizer) Stop() {
	close(po.stopCleanup)
	
	po.contextCacheMu.Lock()
	po.contextCache = make(map[string]*CachedContext)
	po.contextCacheMu.Unlock()
	
	po.lazyKnowledgeBase.indexCacheMu.Lock()
	po.lazyKnowledgeBase.indexCache = make(map[string][]knowledge.SearchResult)
	po.lazyKnowledgeBase.indexCacheMu.Unlock()
}

// ClearAllCaches limpa todos os caches
func (po *PerformanceOptimizer) ClearAllCaches() {
	po.contextCacheMu.Lock()
	po.contextCache = make(map[string]*CachedContext)
	po.contextCacheMu.Unlock()
	
	po.lazyKnowledgeBase.indexCacheMu.Lock()
	po.lazyKnowledgeBase.indexCache = make(map[string][]knowledge.SearchResult)
	po.lazyKnowledgeBase.indexCacheMu.Unlock()
	
	po.logger.Println("All performance caches cleared")
}

// ReloadKnowledgeBase recarrega a base de conhecimento
func (po *PerformanceOptimizer) ReloadKnowledgeBase() error {
	lkb := po.lazyKnowledgeBase
	
	lkb.loadMutex.Lock()
	defer lkb.loadMutex.Unlock()
	
	// Limpar cache de índice
	lkb.indexCacheMu.Lock()
	lkb.indexCache = make(map[string][]knowledge.SearchResult)
	lkb.indexCacheMu.Unlock()
	
	// Recarregar
	lkb.loaded = false
	return po.EnsureKnowledgeBaseLoaded()
}