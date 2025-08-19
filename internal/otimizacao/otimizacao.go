// Package otimizacao provides functionality for optimizing performance in data processing
package otimizacao

import (
	"sync"
	"time"

	"BrxAgente-desafio4/internal/feriados"
)

// Cache para armazenar dados já calculados e evitar recálculos
type Cache struct {
	mu              sync.RWMutex
	feriadosCache   map[int][]feriados.Feriado
	calculoCache    map[string]interface{}
	expirationTimes map[string]time.Time
	defaultTTL      time.Duration
}

// NovoCache cria uma nova instância de Cache
func NovoCache(defaultTTL time.Duration) *Cache {
	return &Cache{
		feriadosCache:   make(map[int][]feriados.Feriado),
		calculoCache:    make(map[string]interface{}),
		expirationTimes: make(map[string]time.Time),
		defaultTTL:      defaultTTL,
	}
}

// ObterFeriadosNacionaisCached obtém feriados nacionais com cache
func (c *Cache) ObterFeriadosNacionaisCached(ano int) []feriados.Feriado {
	c.mu.RLock()
	if feriados, existe := c.feriadosCache[ano]; existe {
		c.mu.RUnlock()
		return feriados
	}
	c.mu.RUnlock()

	// Se não estiver em cache, calcular e armazenar
	feriados := feriados.ObterFeriadosNacionais(ano)
	
	c.mu.Lock()
	c.feriadosCache[ano] = feriados
	c.mu.Unlock()
	
	return feriados
}

// ObterFeriadosEstaduaisCached obtém feriados estaduais com cache
func (c *Cache) ObterFeriadosEstaduaisCached(estado string, ano int) []feriados.Feriado {
	chave := estado + "_" + string(rune(ano))
	
	c.mu.RLock()
	if item, existe := c.calculoCache[chave]; existe {
		if time.Now().Before(c.expirationTimes[chave]) {
			c.mu.RUnlock()
			return item.([]feriados.Feriado)
		}
	}
	c.mu.RUnlock()

	// Se não estiver em cache ou expirado, calcular e armazenar
	feriados := feriados.ObterFeriadosEstaduais(estado, ano)
	
	c.mu.Lock()
	c.calculoCache[chave] = feriados
	c.expirationTimes[chave] = time.Now().Add(c.defaultTTL)
	c.mu.Unlock()
	
	return feriados
}