package knowledge

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// KnowledgeBaseManager gerencia toda a base de conhecimento
type KnowledgeBaseManager struct {
	base     *KnowledgeBase
	index    map[string][]KnowledgeItem // Índice invertido para busca
	keywords map[string][]string        // Cache de palavras-chave por item
	loaded   bool
}

// NewKnowledgeBaseManager cria um novo gerenciador de base de conhecimento
func NewKnowledgeBaseManager() *KnowledgeBaseManager {
	return &KnowledgeBaseManager{
		base:     &KnowledgeBase{},
		index:    make(map[string][]KnowledgeItem),
		keywords: make(map[string][]string),
		loaded:   false,
	}
}

// LoadFromFiles carrega a base de conhecimento dos arquivos JSON
func (kb *KnowledgeBaseManager) LoadFromFiles(dataDir string) error {
	// Carregar políticas
	policiesData, err := kb.loadPolicies(filepath.Join(dataDir, "vr_policies.json"))
	if err != nil {
		return fmt.Errorf("erro ao carregar políticas: %w", err)
	}

	// Carregar regulamentações
	regulationsData, err := kb.loadRegulations(filepath.Join(dataDir, "regulations.json"))
	if err != nil {
		return fmt.Errorf("erro ao carregar regulamentações: %w", err)
	}

	// Carregar regras de negócio
	rulesData, err := kb.loadBusinessRules(filepath.Join(dataDir, "business_rules.json"))
	if err != nil {
		return fmt.Errorf("erro ao carregar regras de negócio: %w", err)
	}

	// Montar base de conhecimento
	kb.base = &KnowledgeBase{
		Policies:    policiesData,
		Regulations: regulationsData,
		Rules:       rulesData,
		Version:     "1.0.0",
		LastUpdated: time.Now(),
	}

	// Construir índice
	if err := kb.buildIndex(); err != nil {
		return fmt.Errorf("erro ao construir índice: %w", err)
	}

	kb.loaded = true
	return nil
}

// loadPolicies carrega políticas de um arquivo JSON
func (kb *KnowledgeBaseManager) loadPolicies(filepath string) ([]Policy, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Policies []Policy `json:"policies"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	return wrapper.Policies, nil
}

// loadRegulations carrega regulamentações de um arquivo JSON
func (kb *KnowledgeBaseManager) loadRegulations(filepath string) ([]Regulation, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Regulations []Regulation `json:"regulations"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	return wrapper.Regulations, nil
}

// loadBusinessRules carrega regras de negócio de um arquivo JSON
func (kb *KnowledgeBaseManager) loadBusinessRules(filepath string) ([]BusinessRule, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Rules []BusinessRule `json:"business_rules"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	return wrapper.Rules, nil
}

// buildIndex constrói o índice invertido para busca eficiente
func (kb *KnowledgeBaseManager) buildIndex() error {
	kb.index = make(map[string][]KnowledgeItem)
	kb.keywords = make(map[string][]string)

	// Indexar políticas
	for _, policy := range kb.base.Policies {
		item := KnowledgeItem{
			ID:           policy.ID,
			Title:        policy.Title,
			Content:      policy.Content,
			Source:       policy.Source,
			EffectiveDate: policy.EffectiveDate,
			Categories:   policy.Categories,
			Metadata:     policy.Metadata,
			Type:         policy.Type,
		}
		kb.indexItem(item)
	}

	// Indexar regulamentações
	for _, regulation := range kb.base.Regulations {
		item := KnowledgeItem{
			ID:           regulation.ID,
			Title:        regulation.Title,
			Content:      regulation.Content,
			Source:       regulation.Source,
			EffectiveDate: regulation.EffectiveDate,
			Categories:   regulation.Categories,
			Metadata:     regulation.Metadata,
			Type:         regulation.Type,
		}
		kb.indexItem(item)
	}

	// Indexar regras de negócio
	for _, rule := range kb.base.Rules {
		item := KnowledgeItem{
			ID:           rule.ID,
			Title:        rule.Title,
			Content:      rule.Content,
			Source:       rule.Source,
			EffectiveDate: rule.EffectiveDate,
			Categories:   rule.Categories,
			Metadata:     rule.Metadata,
			Type:         rule.Type,
		}
		kb.indexItem(item)
	}

	return nil
}

// indexItem adiciona um item ao índice invertido
func (kb *KnowledgeBaseManager) indexItem(item KnowledgeItem) {
	keywords := kb.extractKeywords(item)
	kb.keywords[item.ID] = keywords

	for _, keyword := range keywords {
		kb.index[strings.ToLower(keyword)] = append(kb.index[strings.ToLower(keyword)], item)
	}
}

// extractKeywords extrai palavras-chave de um item
func (kb *KnowledgeBaseManager) extractKeywords(item KnowledgeItem) []string {
	var keywords []string

	// Adicionar palavras do título
	keywords = append(keywords, kb.tokenize(item.Title)...)

	// Adicionar palavras do conteúdo
	keywords = append(keywords, kb.tokenize(item.Content)...)

	// Adicionar categorias
	keywords = append(keywords, item.Categories...)

	// Adicionar ID
	keywords = append(keywords, item.ID)

	// Adicionar tipo
	keywords = append(keywords, item.Type)

	// Remover duplicatas e filtrar palavras muito pequenas
	return kb.deduplicateAndFilter(keywords)
}

// tokenize quebra uma string em tokens (palavras)
func (kb *KnowledgeBaseManager) tokenize(text string) []string {
	// Converter para minúsculas e dividir por espaços e pontuação
	text = strings.ToLower(text)
	separators := []string{" ", ",", ".", ";", ":", "(", ")", "[", "]", "{", "}", "-", "_", "/", "\\"}
	
	words := []string{text}
	for _, sep := range separators {
		var newWords []string
		for _, word := range words {
			newWords = append(newWords, strings.Split(word, sep)...)
		}
		words = newWords
	}

	// Filtrar palavras vazias e muito pequenas
	var filtered []string
	stopWords := map[string]bool{
		"o": true, "a": true, "os": true, "as": true, "um": true, "uma": true,
		"de": true, "do": true, "da": true, "dos": true, "das": true,
		"em": true, "no": true, "na": true, "nos": true, "nas": true,
		"para": true, "por": true, "com": true, "sem": true, "até": true,
		"e": true, "ou": true, "mas": true, "que": true, "se": true,
	}

	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) >= 3 && !stopWords[word] {
			filtered = append(filtered, word)
		}
	}

	return filtered
}

// deduplicateAndFilter remove duplicatas e filtra palavras
func (kb *KnowledgeBaseManager) deduplicateAndFilter(words []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, word := range words {
		if !seen[word] && len(word) >= 2 {
			seen[word] = true
			result = append(result, word)
		}
	}

	return result
}

// Search busca na base de conhecimento por uma query
func (kb *KnowledgeBaseManager) Search(query string, limit int) ([]SearchResult, error) {
	if !kb.loaded {
		return nil, fmt.Errorf("base de conhecimento não foi carregada")
	}

	queryWords := kb.tokenize(query)
	if len(queryWords) == 0 {
		return []SearchResult{}, nil
	}

	// Buscar itens que contenham as palavras da query
	itemScores := make(map[string]float64)
	itemDetails := make(map[string]KnowledgeItem)

	for _, word := range queryWords {
		if items, exists := kb.index[strings.ToLower(word)]; exists {
			for _, item := range items {
				// Calcular score baseado na relevância
				score := kb.calculateRelevanceScore(item, word, query)
				itemScores[item.ID] += score
				itemDetails[item.ID] = item
			}
		}
	}

	// Converter para slice e ordenar por score
	var results []SearchResult
	for itemID, score := range itemScores {
		item := itemDetails[itemID]
		highlights := kb.extractHighlights(item, queryWords)
		
		result := SearchResult{
			Item:       item,
			Score:      score,
			Highlights: highlights,
			Context:    kb.generateContext(item, queryWords),
		}
		results = append(results, result)
	}

	// Ordenar por score (maior primeiro)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Limitar resultados
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// calculateRelevanceScore calcula o score de relevância para um item
func (kb *KnowledgeBaseManager) calculateRelevanceScore(item KnowledgeItem, matchedWord, query string) float64 {
	score := 1.0

	// Boost se a palavra aparece no título
	if strings.Contains(strings.ToLower(item.Title), strings.ToLower(matchedWord)) {
		score += 2.0
	}

	// Boost se a palavra aparece nas categorias
	for _, category := range item.Categories {
		if strings.Contains(strings.ToLower(category), strings.ToLower(matchedWord)) {
			score += 1.5
		}
	}

	// Boost se é uma correspondência exata no ID
	if strings.EqualFold(item.ID, matchedWord) {
		score += 3.0
	}

	// Penalizar itens muito antigos (mais de 2 anos)
	if time.Since(item.EffectiveDate) > 2*365*24*time.Hour {
		score *= 0.8
	}

	return score
}

// extractHighlights extrai trechos relevantes que matcharam a busca
func (kb *KnowledgeBaseManager) extractHighlights(item KnowledgeItem, queryWords []string) []string {
	var highlights []string

	for _, word := range queryWords {
		// Buscar no título
		if strings.Contains(strings.ToLower(item.Title), strings.ToLower(word)) {
			highlights = append(highlights, fmt.Sprintf("Título: %s", item.Title))
		}

		// Buscar no conteúdo (extrair contexto ao redor da palavra)
		content := strings.ToLower(item.Content)
		wordLower := strings.ToLower(word)
		if index := strings.Index(content, wordLower); index != -1 {
			start := max(0, index-50)
			end := min(len(content), index+len(word)+50)
			highlight := item.Content[start:end]
			highlights = append(highlights, fmt.Sprintf("...%s...", strings.TrimSpace(highlight)))
		}
	}

	return highlights
}

// generateContext gera contexto adicional para um resultado
func (kb *KnowledgeBaseManager) generateContext(item KnowledgeItem, queryWords []string) string {
	context := fmt.Sprintf("Tipo: %s | Fonte: %s", item.Type, item.Source)
	
	if len(item.Categories) > 0 {
		context += fmt.Sprintf(" | Categorias: %s", strings.Join(item.Categories, ", "))
	}

	return context
}

// GetByCategory busca itens por categoria
func (kb *KnowledgeBaseManager) GetByCategory(category string) ([]KnowledgeItem, error) {
	if !kb.loaded {
		return nil, fmt.Errorf("base de conhecimento não foi carregada")
	}

	var items []KnowledgeItem

	// Buscar em políticas
	for _, policy := range kb.base.Policies {
		for _, cat := range policy.Categories {
			if strings.EqualFold(cat, category) {
				item := KnowledgeItem{
					ID:           policy.ID,
					Title:        policy.Title,
					Content:      policy.Content,
					Source:       policy.Source,
					EffectiveDate: policy.EffectiveDate,
					Categories:   policy.Categories,
					Metadata:     policy.Metadata,
					Type:         policy.Type,
				}
				items = append(items, item)
				break
			}
		}
	}

	// Buscar em regulamentações
	for _, regulation := range kb.base.Regulations {
		for _, cat := range regulation.Categories {
			if strings.EqualFold(cat, category) {
				item := KnowledgeItem{
					ID:           regulation.ID,
					Title:        regulation.Title,
					Content:      regulation.Content,
					Source:       regulation.Source,
					EffectiveDate: regulation.EffectiveDate,
					Categories:   regulation.Categories,
					Metadata:     regulation.Metadata,
					Type:         regulation.Type,
				}
				items = append(items, item)
				break
			}
		}
	}

	// Buscar em regras de negócio
	for _, rule := range kb.base.Rules {
		for _, cat := range rule.Categories {
			if strings.EqualFold(cat, category) {
				item := KnowledgeItem{
					ID:           rule.ID,
					Title:        rule.Title,
					Content:      rule.Content,
					Source:       rule.Source,
					EffectiveDate: rule.EffectiveDate,
					Categories:   rule.Categories,
					Metadata:     rule.Metadata,
					Type:         rule.Type,
				}
				items = append(items, item)
				break
			}
		}
	}

	return items, nil
}

// GetByID busca um item específico por ID
func (kb *KnowledgeBaseManager) GetByID(id string) (KnowledgeItem, error) {
	if !kb.loaded {
		return KnowledgeItem{}, fmt.Errorf("base de conhecimento não foi carregada")
	}

	// Buscar em políticas
	for _, policy := range kb.base.Policies {
		if policy.ID == id {
			return KnowledgeItem{
				ID:           policy.ID,
				Title:        policy.Title,
				Content:      policy.Content,
				Source:       policy.Source,
				EffectiveDate: policy.EffectiveDate,
				Categories:   policy.Categories,
				Metadata:     policy.Metadata,
				Type:         policy.Type,
			}, nil
		}
	}

	// Buscar em regulamentações
	for _, regulation := range kb.base.Regulations {
		if regulation.ID == id {
			return KnowledgeItem{
				ID:           regulation.ID,
				Title:        regulation.Title,
				Content:      regulation.Content,
				Source:       regulation.Source,
				EffectiveDate: regulation.EffectiveDate,
				Categories:   regulation.Categories,
				Metadata:     regulation.Metadata,
				Type:         regulation.Type,
			}, nil
		}
	}

	// Buscar em regras de negócio
	for _, rule := range kb.base.Rules {
		if rule.ID == id {
			return KnowledgeItem{
				ID:           rule.ID,
				Title:        rule.Title,
				Content:      rule.Content,
				Source:       rule.Source,
				EffectiveDate: rule.EffectiveDate,
				Categories:   rule.Categories,
				Metadata:     rule.Metadata,
				Type:         rule.Type,
			}, nil
		}
	}

	return KnowledgeItem{}, fmt.Errorf("item com ID %s não encontrado", id)
}

// GetStats retorna estatísticas da base de conhecimento
func (kb *KnowledgeBaseManager) GetStats() map[string]interface{} {
	if !kb.loaded {
		return map[string]interface{}{
			"loaded": false,
		}
	}

	return map[string]interface{}{
		"loaded":                true,
		"policies_count":        len(kb.base.Policies),
		"regulations_count":     len(kb.base.Regulations),
		"business_rules_count":  len(kb.base.Rules),
		"total_items":          len(kb.base.Policies) + len(kb.base.Regulations) + len(kb.base.Rules),
		"index_size":           len(kb.index),
		"version":              kb.base.Version,
		"last_updated":         kb.base.LastUpdated,
	}
}

// Funções auxiliares
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}