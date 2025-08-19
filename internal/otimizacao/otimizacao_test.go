// Package otimizacao provides functionality for optimizing performance in data processing
package otimizacao

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NovoCache(1 * time.Minute)

	// Testar criação de cache
	if cache == nil {
		t.Error("Falha ao criar cache")
	}

	// Testar propriedades do cache
	if cache.defaultTTL != 1*time.Minute {
		t.Errorf("Esperava TTL de 1 minuto, obteve %v", cache.defaultTTL)
	}
}