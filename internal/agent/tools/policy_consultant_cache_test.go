package tools

import (
	"context"
	"strings"
	"testing"
	"time"

	"BrxAgente-desafio4/internal/cache"
)

func TestPolicyConsultantTool_CacheIntegration(t *testing.T) {
	// Create a minimal test for cache integration
	// Note: This test focuses on cache functionality, not the full policy engine
	
	// Create cache configuration
	config := cache.DefaultCacheConfig()
	config.TTL = 1 * time.Hour // Set to 1 hour for testing
	
	// Create knowledge cache
	knowledgeCache := cache.NewKnowledgeCache(config)
	
	// Create basic PolicyConsultantTool structure for testing cache
	pct := &PolicyConsultantTool{
		cache: knowledgeCache,
		initialized: false, // Set to false to test error handling
	}
	
	// Test 1: Cache is properly initialized
	if pct.GetCache() != knowledgeCache {
		t.Error("Cache not properly initialized in PolicyConsultantTool")
	}
	
	// Test 2: Cache enable/disable functionality
	pct.DisableCache()
	if pct.cache.IsEnabled() {
		t.Error("Cache should be disabled")
	}
	
	pct.EnableCache()
	if !pct.cache.IsEnabled() {
		t.Error("Cache should be enabled")
	}
	
	// Test 3: Cache clear functionality
	pct.ClearCache()
	entries := pct.cache.GetCacheEntries()
	if len(entries) != 0 {
		t.Error("Cache should be empty after clear")
	}
	
	// Test 4: Cache metrics functionality
	metrics := pct.GetCacheMetrics()
	if metrics.TotalRequests < 0 {
		t.Error("Invalid cache metrics")
	}
}

func TestPolicyConsultantTool_CacheUsage(t *testing.T) {
	// Test cache usage with actual queries (without full knowledge base)
	
	config := cache.DefaultCacheConfig()
	knowledgeCache := cache.NewKnowledgeCache(config)
	
	pct := &PolicyConsultantTool{
		cache: knowledgeCache,
		initialized: false,
	}
	
	ctx := context.Background()
	
	// Test with tool not initialized - should return error without using cache
	query := `{"query": "Teste de cache", "type": "simple"}`
	
	result, err := pct.Execute(ctx, query)
	if err == nil {
		t.Error("Expected error when tool is not initialized")
	}
	
	if result != "" {
		t.Error("Expected empty result when tool is not initialized")
	}
	
	// Verify error message contains expected text
	expectedErrorText := "ferramenta não foi inicializada"
	if !strings.Contains(err.Error(), expectedErrorText) {
		t.Errorf("Error message should contain '%s', got: %s", expectedErrorText, err.Error())
	}
	
	// Verify no cache entries were created due to error
	entries := pct.cache.GetCacheEntries()
	if len(entries) != 0 {
		t.Error("Cache should remain empty when tool execution fails")
	}
	
	// Verify metrics show no successful operations
	metrics := pct.GetCacheMetrics()
	if metrics.HitCount > 0 {
		t.Error("Should have no cache hits when tool is not initialized")
	}
}

func TestPolicyConsultantTool_CacheKeyGeneration(t *testing.T) {
	// Test that similar queries generate the same cache key
	
	config := cache.DefaultCacheConfig()
	knowledgeCache := cache.NewKnowledgeCache(config)
	
	// Directly test cache key generation by using cache methods
	query1 := "Diretores têm direito a VR?"
	query2 := "DIRETORES TEM DIREITO A VR?"
	query3 := "Diretores  têm   direito   a   VR ???"
	
	// Set a test response for the first query
	knowledgeCache.Set(query1, "Resposta sobre diretores", 0.9)
	
	// Try to get with similar queries - should hit cache due to normalization
	entry2 := knowledgeCache.Get(query2)
	entry3 := knowledgeCache.Get(query3)
	
	if entry2 == nil {
		t.Error("Cache miss for normalized query 2")
	}
	
	if entry3 == nil {
		t.Error("Cache miss for normalized query 3")
	}
	
	if entry2 != nil && entry3 != nil {
		if entry2.Response != entry3.Response {
			t.Error("Normalized queries should return same cached response")
		}
	}
	
	// Verify hit counts increased (starts from 1 after Set, then increments)
	if entry2 != nil && entry2.HitCount < 1 {
		t.Errorf("Expected hit count >= 1 for query 2, got %d", entry2.HitCount)
	}
	
	if entry3 != nil && entry3.HitCount < 2 {
		t.Errorf("Expected hit count >= 2 for query 3, got %d", entry3.HitCount)
	}
}

func TestPolicyConsultantTool_CacheMetrics(t *testing.T) {
	// Test cache metrics collection
	
	config := cache.DefaultCacheConfig()
	knowledgeCache := cache.NewKnowledgeCache(config)
	
	pct := &PolicyConsultantTool{
		cache: knowledgeCache,
	}
	
	// Add some test data to cache
	knowledgeCache.Set("Query 1", "Response 1", 0.9)
	knowledgeCache.Set("Query 2", "Response 2", 0.8)
	
	// Generate some hits and misses
	knowledgeCache.Get("Query 1") // Hit
	knowledgeCache.Get("Query 1") // Hit
	knowledgeCache.Get("Query 3") // Miss
	knowledgeCache.Get("Query 4") // Miss
	
	// Check metrics
	metrics := pct.GetCacheMetrics()
	
	if metrics.HitCount != 2 {
		t.Errorf("Expected 2 hits, got %d", metrics.HitCount)
	}
	
	if metrics.MissCount != 2 {
		t.Errorf("Expected 2 misses, got %d", metrics.MissCount)
	}
	
	if metrics.TotalRequests != 4 {
		t.Errorf("Expected 4 total requests, got %d", metrics.TotalRequests)
	}
	
	if metrics.CacheSize != 2 {
		t.Errorf("Expected cache size 2, got %d", metrics.CacheSize)
	}
	
	// Calculate hit rate
	expectedHitRate := float64(2) / float64(4) // 2 hits out of 4 requests
	actualHitRate := float64(metrics.HitCount) / float64(metrics.TotalRequests)
	
	if actualHitRate != expectedHitRate {
		t.Errorf("Expected hit rate %.2f, got %.2f", expectedHitRate, actualHitRate)
	}
}

func TestPolicyConsultantTool_CacheConfiguration(t *testing.T) {
	// Test cache configuration
	
	config := cache.DefaultCacheConfig()
	config.MaxEntries = 100
	config.TTL = 30 * time.Minute
	config.Enabled = true
	
	knowledgeCache := cache.NewKnowledgeCache(config)
	
	pct := &PolicyConsultantTool{
		cache: knowledgeCache,
	}
	
	cacheConfig := pct.cache.GetConfig()
	
	if cacheConfig.MaxEntries != 100 {
		t.Errorf("Expected MaxEntries 100, got %d", cacheConfig.MaxEntries)
	}
	
	if cacheConfig.TTL != 30*time.Minute {
		t.Errorf("Expected TTL 30 minutes, got %v", cacheConfig.TTL)
	}
	
	if !cacheConfig.Enabled {
		t.Error("Cache should be enabled")
	}
}

func BenchmarkPolicyConsultantTool_CacheHit(b *testing.B) {
	// Benchmark cache hit performance
	
	config := cache.DefaultCacheConfig()
	knowledgeCache := cache.NewKnowledgeCache(config)
	
	// Pre-populate cache
	for i := 0; i < 100; i++ {
		query := "Pergunta frequente"
		knowledgeCache.Set(query, "Resposta frequente", 0.9)
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		knowledgeCache.Get("Pergunta frequente")
	}
}