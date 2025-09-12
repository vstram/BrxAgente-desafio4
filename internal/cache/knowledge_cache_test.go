package cache

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNewKnowledgeCache(t *testing.T) {
	config := DefaultCacheConfig()
	cache := NewKnowledgeCache(config)

	if cache == nil {
		t.Fatal("NewKnowledgeCache retornou nil")
	}

	if !cache.IsEnabled() {
		t.Error("Cache deveria estar habilitado por padrão")
	}

	if cache.GetConfig().MaxEntries != config.MaxEntries {
		t.Errorf("MaxEntries = %d, esperado %d", cache.GetConfig().MaxEntries, config.MaxEntries)
	}
}

func TestKnowledgeCache_SetAndGet(t *testing.T) {
	cache := NewKnowledgeCache(DefaultCacheConfig())

	// Teste básico de Set/Get
	question := "Teste de cache"
	response := "Resposta de teste"
	confidence := 0.9

	cache.Set(question, response, confidence)

	entry := cache.Get(question)
	if entry == nil {
		t.Fatal("Cache miss inesperado")
	}

	if entry.Response != response {
		t.Errorf("Response = %s, esperado %s", entry.Response, response)
	}

	if entry.Confidence != confidence {
		t.Errorf("Confidence = %f, esperado %f", entry.Confidence, confidence)
	}

	if entry.HitCount != 1 {
		t.Errorf("HitCount = %d, esperado 1", entry.HitCount)
	}
}

func TestKnowledgeCache_QuestionNormalization(t *testing.T) {
	cache := NewKnowledgeCache(DefaultCacheConfig())

	// Perguntas similares que devem ser normalizadas para o mesmo cache key
	variations := []string{
		"Diretores têm direito a VR?",
		"DIRETORES TEM DIREITO A VR?",
		"Diretores  têm   direito   a   VR ?",
		"diretores têm direito a vr?",
		"Diretores têm direito a VR ???",
	}

	response := "Resposta sobre diretores"

	// Set usando a primeira variação
	cache.Set(variations[0], response, 0.9)

	// Testar se todas as variações resultam em cache hit
	for i, variation := range variations {
		entry := cache.Get(variation)
		if entry == nil {
			t.Errorf("Cache miss para variação %d: %s", i, variation)
			continue
		}

		if entry.Response != response {
			t.Errorf("Resposta incorreta para variação %d: %s", i, variation)
		}
	}

	// Verificar métricas
	metrics := cache.GetMetrics()
	if metrics.HitCount != int64(len(variations)) {
		t.Errorf("HitCount = %d, esperado %d", metrics.HitCount, len(variations))
	}
}

func TestKnowledgeCache_TTLExpiration(t *testing.T) {
	config := DefaultCacheConfig()
	config.TTL = 50 * time.Millisecond // TTL muito curto para teste
	config.CleanupInterval = 0         // Desabilitar limpeza automática
	cache := NewKnowledgeCache(config)

	question := "Teste TTL"
	response := "Resposta que vai expirar"

	// Set entrada no cache
	cache.Set(question, response, 0.9)

	// Verificar que está no cache
	entry := cache.Get(question)
	if entry == nil {
		t.Fatal("Cache miss inesperado antes da expiração")
	}

	// Aguardar expiração
	time.Sleep(100 * time.Millisecond)

	// Verificar que expirou
	entry = cache.Get(question)
	if entry != nil {
		t.Error("Cache hit inesperado após expiração")
	}

	// Verificar métricas
	metrics := cache.GetMetrics()
	if metrics.ExpiredCount == 0 {
		t.Error("ExpiredCount deveria ser > 0")
	}
}

func TestKnowledgeCache_LRUEviction(t *testing.T) {
	config := DefaultCacheConfig()
	config.MaxEntries = 3 // Limite pequeno para testar eviction
	cache := NewKnowledgeCache(config)

	// Adicionar entradas até o limite
	for i := 0; i < config.MaxEntries; i++ {
		question := fmt.Sprintf("Pergunta %d", i)
		response := fmt.Sprintf("Resposta %d", i)
		cache.Set(question, response, 0.9)
	}

	// Verificar que todas estão no cache
	for i := 0; i < config.MaxEntries; i++ {
		question := fmt.Sprintf("Pergunta %d", i)
		entry := cache.Get(question)
		if entry == nil {
			t.Errorf("Cache miss inesperado para pergunta %d", i)
		}
	}

	// Adicionar uma entrada extra para forçar eviction
	cache.Set("Pergunta extra", "Resposta extra", 0.9)

	// A primeira entrada (menos recentemente usada) deveria ter sido removida
	entry := cache.Get("Pergunta 0")
	if entry != nil {
		t.Error("Pergunta 0 deveria ter sido removida por LRU eviction")
	}

	// As outras deveriam ainda estar no cache
	for i := 1; i < config.MaxEntries; i++ {
		question := fmt.Sprintf("Pergunta %d", i)
		entry := cache.Get(question)
		if entry == nil {
			t.Errorf("Cache miss inesperado para pergunta %d após eviction", i)
		}
	}

	// A nova entrada deveria estar no cache
	entry = cache.Get("Pergunta extra")
	if entry == nil {
		t.Error("Cache miss inesperado para pergunta extra")
	}

	// Verificar métricas de eviction
	metrics := cache.GetMetrics()
	if metrics.EvictionCount == 0 {
		t.Error("EvictionCount deveria ser > 0")
	}
}

func TestKnowledgeCache_HitRate(t *testing.T) {
	cache := NewKnowledgeCache(DefaultCacheConfig())

	// Adicionar algumas entradas
	cache.Set("Q1", "R1", 0.9)
	cache.Set("Q2", "R2", 0.9)

	// Fazer algumas consultas (mix de hits e misses)
	cache.Get("Q1") // Hit
	cache.Get("Q2") // Hit
	cache.Get("Q1") // Hit
	cache.Get("Q3") // Miss
	cache.Get("Q4") // Miss

	hitRate := cache.GetHitRate()
	expectedHitRate := 3.0 / 5.0 // 3 hits de 5 requests

	if hitRate != expectedHitRate {
		t.Errorf("HitRate = %f, esperado %f", hitRate, expectedHitRate)
	}

	metrics := cache.GetMetrics()
	if metrics.HitCount != 3 {
		t.Errorf("HitCount = %d, esperado 3", metrics.HitCount)
	}

	if metrics.MissCount != 2 {
		t.Errorf("MissCount = %d, esperado 2", metrics.MissCount)
	}

	if metrics.TotalRequests != 5 {
		t.Errorf("TotalRequests = %d, esperado 5", metrics.TotalRequests)
	}
}

func TestKnowledgeCache_DisabledCache(t *testing.T) {
	cache := NewKnowledgeCache(DefaultCacheConfig())
	cache.Disable()

	// Tentar definir entrada com cache desabilitado
	cache.Set("Q1", "R1", 0.9)

	// Verificar que não foi armazenado
	entry := cache.Get("Q1")
	if entry != nil {
		t.Error("Cache hit inesperado com cache desabilitado")
	}

	// Habilitar novamente
	cache.Enable()

	// Agora deveria funcionar
	cache.Set("Q1", "R1", 0.9)
	entry = cache.Get("Q1")
	if entry == nil {
		t.Error("Cache miss inesperado com cache habilitado")
	}
}

func TestKnowledgeCache_Clear(t *testing.T) {
	cache := NewKnowledgeCache(DefaultCacheConfig())

	// Adicionar algumas entradas
	cache.Set("Q1", "R1", 0.9)
	cache.Set("Q2", "R2", 0.9)

	// Verificar que estão no cache
	if cache.Get("Q1") == nil || cache.Get("Q2") == nil {
		t.Fatal("Falha ao adicionar entradas no cache")
	}

	// Limpar cache
	cache.Clear()

	// Verificar que foram removidas
	if cache.Get("Q1") != nil || cache.Get("Q2") != nil {
		t.Error("Entradas ainda presentes após Clear()")
	}

	metrics := cache.GetMetrics()
	if metrics.CacheSize != 0 {
		t.Errorf("CacheSize = %d após Clear(), esperado 0", metrics.CacheSize)
	}
}

func TestKnowledgeCache_EdgeCases(t *testing.T) {
	cache := NewKnowledgeCache(DefaultCacheConfig())

	// Teste com strings vazias
	cache.Set("", "response", 0.9)
	entry := cache.Get("")
	if entry != nil {
		t.Error("Cache hit inesperado para pergunta vazia")
	}

	cache.Set("question", "", 0.9)
	entry = cache.Get("question")
	if entry != nil {
		t.Error("Cache hit inesperado para resposta vazia")
	}

	// Teste com pergunta muito longa
	longQuestion := strings.Repeat("palavra ", 1000)
	cache.Set(longQuestion, "resposta", 0.9)
	entry = cache.Get(longQuestion)
	if entry == nil {
		t.Error("Cache miss inesperado para pergunta longa")
	}
}

func TestKnowledgeCache_MultipleHits(t *testing.T) {
	cache := NewKnowledgeCache(DefaultCacheConfig())

	question := "Pergunta frequente"
	response := "Resposta frequente"

	// Set inicial
	cache.Set(question, response, 0.9)

	// Fazer múltiplos Gets para testar contador de hits
	initialEntry := cache.Get(question)
	if initialEntry.HitCount != 1 {
		t.Errorf("HitCount inicial = %d, esperado 1", initialEntry.HitCount)
	}

	// Mais hits
	for i := 0; i < 5; i++ {
		entry := cache.Get(question)
		expectedHits := int64(i + 2) // +2 porque já teve 1 hit inicial
		if entry.HitCount != expectedHits {
			t.Errorf("HitCount na iteração %d = %d, esperado %d", i, entry.HitCount, expectedHits)
		}
	}

	finalEntry := cache.Get(question)
	if finalEntry.HitCount != 7 { // 1 inicial + 5 iterações + 1 final
		t.Errorf("HitCount final = %d, esperado 7", finalEntry.HitCount)
	}
}

func TestKnowledgeCache_Cleanup(t *testing.T) {
	config := DefaultCacheConfig()
	config.TTL = 50 * time.Millisecond
	config.CleanupInterval = 100 * time.Millisecond
	cache := NewKnowledgeCache(config)

	// Adicionar entradas
	cache.Set("Q1", "R1", 0.9)
	cache.Set("Q2", "R2", 0.9)

	// Aguardar expiração + cleanup
	time.Sleep(200 * time.Millisecond)

	// Verificar que foram limpas
	metrics := cache.GetMetrics()
	if metrics.ExpiredCount == 0 {
		t.Error("ExpiredCount deveria ser > 0 após cleanup")
	}

	// Verificar que cache está vazio
	if cache.Get("Q1") != nil || cache.Get("Q2") != nil {
		t.Error("Entradas ainda presentes após cleanup automático")
	}
}

func TestKnowledgeCache_GetCacheEntries(t *testing.T) {
	cache := NewKnowledgeCache(DefaultCacheConfig())

	// Adicionar algumas entradas
	cache.Set("Q1", "R1", 0.9)
	cache.Set("Q2", "R2", 0.8)

	entries := cache.GetCacheEntries()

	if len(entries) != 2 {
		t.Errorf("Número de entradas = %d, esperado 2", len(entries))
	}

	// Verificar que são cópias (não referências)
	for _, entry := range entries {
		entry.HitCount = 999 // Modificar cópia
	}

	// Verificar que original não foi modificado
	originalEntry := cache.Get("Q1")
	if originalEntry.HitCount == 999 {
		t.Error("GetCacheEntries não retornou cópias das entradas")
	}
}

// BenchmarkKnowledgeCache_SetGet mede performance do cache
func BenchmarkKnowledgeCache_SetGet(b *testing.B) {
	cache := NewKnowledgeCache(DefaultCacheConfig())
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		question := fmt.Sprintf("Pergunta %d", i%1000) // Reusar algumas perguntas
		cache.Set(question, "resposta", 0.9)
		cache.Get(question)
	}
}

// BenchmarkKnowledgeCache_GetHit mede performance apenas de cache hits
func BenchmarkKnowledgeCache_GetHit(b *testing.B) {
	cache := NewKnowledgeCache(DefaultCacheConfig())
	
	// Pré-popular cache
	for i := 0; i < 100; i++ {
		question := fmt.Sprintf("Pergunta %d", i)
		cache.Set(question, "resposta", 0.9)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		question := fmt.Sprintf("Pergunta %d", i%100)
		cache.Get(question)
	}
}

func TestKnowledgeCache_NormalizationEdgeCases(t *testing.T) {
	cache := NewKnowledgeCache(DefaultCacheConfig())

	testCases := []struct {
		question1 string
		question2 string
		shouldMatch bool
	}{
		{
			"Como calcular VR?",
			"como calcular vr?",
			true, // Case insensitive
		},
		{
			"Qual é o valor?",
			"Qual é o valor???",
			true, // Pontuação removida
		},
		{
			"Quantos colaboradores processados?",
			"Quantos    colaboradores     processados?",
			true, // Espaços normalizados
		},
		{
			"Como calcular dias úteis?",
			"Como calcular dias uteis?", 
			true, // Acentos removidos
		},
		{
			"Pergunta completamente diferente",
			"Outra pergunta totalmente distinta",
			false, // Não devem fazer match
		},
	}

	for i, tc := range testCases {
		// Limpar cache entre tests
		cache.Clear()
		
		// Set com primeira pergunta
		cache.Set(tc.question1, fmt.Sprintf("Resposta %d", i), 0.9)
		
		// Get com segunda pergunta
		entry := cache.Get(tc.question2)
		
		if tc.shouldMatch && entry == nil {
			t.Errorf("Caso %d: esperava match entre '%s' e '%s'", 
				i, tc.question1, tc.question2)
		}
		
		if !tc.shouldMatch && entry != nil {
			t.Errorf("Caso %d: não esperava match entre '%s' e '%s'", 
				i, tc.question1, tc.question2)
		}
	}
}