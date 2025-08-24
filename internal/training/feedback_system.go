package training

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// ResponseFeedback representa feedback de uma resposta específica
type ResponseFeedback struct {
	ID          string    `json:"id"`
	Question    string    `json:"question"`
	Response    string    `json:"response"`
	UserRating  int       `json:"user_rating"`  // 1-5
	Corrections string    `json:"corrections"`
	Feedback    string    `json:"feedback"`
	Timestamp   time.Time `json:"timestamp"`
	Source      string    `json:"source"`       // "user", "automated", "expert"
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
}

// QualityMetrics representa métricas de qualidade de respostas
type QualityMetrics struct {
	TotalResponses    int     `json:"total_responses"`
	AverageRating     float64 `json:"average_rating"`
	AccuracyScore     float64 `json:"accuracy_score"`
	ConsistencyScore  float64 `json:"consistency_score"`
	CompletenessScore float64 `json:"completeness_score"`
	ResponseTime      float64 `json:"response_time_avg_seconds"`
	LastUpdated       time.Time `json:"last_updated"`
}

// LearningPattern representa um padrão aprendido do feedback
type LearningPattern struct {
	Pattern     string    `json:"pattern"`
	Frequency   int       `json:"frequency"`
	Confidence  float64   `json:"confidence"`
	Category    string    `json:"category"`
	Examples    []string  `json:"examples"`
	LastSeen    time.Time `json:"last_seen"`
	Improvement string    `json:"improvement"`
}

// FeedbackSystem gerencia o sistema de feedback e aprendizado
type FeedbackSystem struct {
	feedbackPath string
	feedback     []ResponseFeedback
	metrics      QualityMetrics
	patterns     []LearningPattern
	mutex        sync.RWMutex
	knowledgeManager *KnowledgeManager
}

// NewFeedbackSystem cria um novo sistema de feedback
func NewFeedbackSystem(feedbackPath string, km *KnowledgeManager) *FeedbackSystem {
	return &FeedbackSystem{
		feedbackPath:     feedbackPath,
		feedback:         make([]ResponseFeedback, 0),
		patterns:         make([]LearningPattern, 0),
		knowledgeManager: km,
	}
}

// LoadFeedbackData carrega dados de feedback existentes
func (fs *FeedbackSystem) LoadFeedbackData() error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// Criar diretório se não existir
	if err := os.MkdirAll(fs.feedbackPath, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de feedback: %v", err)
	}

	// Carregar feedback
	feedbackFile := filepath.Join(fs.feedbackPath, "feedback.json")
	if data, err := os.ReadFile(feedbackFile); err == nil {
		if err := json.Unmarshal(data, &fs.feedback); err != nil {
			return fmt.Errorf("erro ao carregar feedback: %v", err)
		}
	}

	// Carregar métricas
	metricsFile := filepath.Join(fs.feedbackPath, "metrics.json")
	if data, err := os.ReadFile(metricsFile); err == nil {
		if err := json.Unmarshal(data, &fs.metrics); err != nil {
			return fmt.Errorf("erro ao carregar métricas: %v", err)
		}
	}

	// Carregar padrões
	patternsFile := filepath.Join(fs.feedbackPath, "patterns.json")
	if data, err := os.ReadFile(patternsFile); err == nil {
		if err := json.Unmarshal(data, &fs.patterns); err != nil {
			return fmt.Errorf("erro ao carregar padrões: %v", err)
		}
	}

	return nil
}

// SaveFeedbackData salva todos os dados de feedback
func (fs *FeedbackSystem) SaveFeedbackData() error {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	// Salvar feedback
	feedbackData, err := json.MarshalIndent(fs.feedback, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar feedback: %v", err)
	}
	
	feedbackFile := filepath.Join(fs.feedbackPath, "feedback.json")
	if err := os.WriteFile(feedbackFile, feedbackData, 0644); err != nil {
		return fmt.Errorf("erro ao salvar feedback: %v", err)
	}

	// Salvar métricas
	metricsData, err := json.MarshalIndent(fs.metrics, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar métricas: %v", err)
	}
	
	metricsFile := filepath.Join(fs.feedbackPath, "metrics.json")
	if err := os.WriteFile(metricsFile, metricsData, 0644); err != nil {
		return fmt.Errorf("erro ao salvar métricas: %v", err)
	}

	// Salvar padrões
	patternsData, err := json.MarshalIndent(fs.patterns, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar padrões: %v", err)
	}
	
	patternsFile := filepath.Join(fs.feedbackPath, "patterns.json")
	if err := os.WriteFile(patternsFile, patternsData, 0644); err != nil {
		return fmt.Errorf("erro ao salvar padrões: %v", err)
	}

	return nil
}

// AddFeedback adiciona novo feedback ao sistema
func (fs *FeedbackSystem) AddFeedback(feedback ResponseFeedback) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	// Gerar ID se não fornecido
	if feedback.ID == "" {
		feedback.ID = fmt.Sprintf("fb_%d", time.Now().UnixNano())
	}

	// Definir timestamp se não fornecido
	if feedback.Timestamp.IsZero() {
		feedback.Timestamp = time.Now()
	}

	// Validar rating
	if feedback.UserRating < 1 || feedback.UserRating > 5 {
		return fmt.Errorf("rating deve estar entre 1 e 5")
	}

	fs.feedback = append(fs.feedback, feedback)
	
	// Atualizar métricas
	fs.updateMetrics()
	
	// Identificar padrões
	fs.identifyPatterns()

	// Salvar dados
	return fs.SaveFeedbackData()
}

// updateMetrics atualiza as métricas de qualidade
func (fs *FeedbackSystem) updateMetrics() {
	if len(fs.feedback) == 0 {
		return
	}

	var totalRating float64
	var accuracySum float64
	var consistencySum float64
	var completenessSum float64

	for _, fb := range fs.feedback {
		totalRating += float64(fb.UserRating)
		
		// Calcular accuracy baseado no rating (simplificado)
		if fb.UserRating >= 4 {
			accuracySum += 1.0
		} else if fb.UserRating == 3 {
			accuracySum += 0.6
		}

		// Outros cálculos de qualidade...
	}

	fs.metrics.TotalResponses = len(fs.feedback)
	fs.metrics.AverageRating = totalRating / float64(len(fs.feedback))
	fs.metrics.AccuracyScore = accuracySum / float64(len(fs.feedback))
	fs.metrics.LastUpdated = time.Now()

	// Calcular consistência (respostas similares para perguntas similares)
	fs.metrics.ConsistencyScore = fs.calculateConsistencyScore()
}

// calculateConsistencyScore calcula score de consistência
func (fs *FeedbackSystem) calculateConsistencyScore() float64 {
	// Implementação básica - na prática seria mais sofisticada
	questionGroups := make(map[string][]ResponseFeedback)
	
	for _, fb := range fs.feedback {
		// Agrupar por pergunta normalizada
		normalizedQ := fs.normalizeQuestion(fb.Question)
		questionGroups[normalizedQ] = append(questionGroups[normalizedQ], fb)
	}

	var consistencySum float64
	groupCount := 0
	
	for _, group := range questionGroups {
		if len(group) > 1 {
			// Calcular variação das respostas no grupo
			consistency := fs.calculateGroupConsistency(group)
			consistencySum += consistency
			groupCount++
		}
	}

	if groupCount == 0 {
		return 1.0 // Se não há grupos para comparar, assume consistência perfeita
	}

	return consistencySum / float64(groupCount)
}

// normalizeQuestion normaliza uma pergunta para agrupamento
func (fs *FeedbackSystem) normalizeQuestion(question string) string {
	// Implementação básica de normalização
	normalized := strings.ToLower(question)
	// Remover pontuação, palavras comuns, etc.
	// Na prática, seria mais sofisticado
	return normalized
}

// calculateGroupConsistency calcula consistência dentro de um grupo
func (fs *FeedbackSystem) calculateGroupConsistency(group []ResponseFeedback) float64 {
	if len(group) < 2 {
		return 1.0
	}

	// Calcular similaridade das respostas (implementação básica)
	var similarities []float64
	
	for i := 0; i < len(group); i++ {
		for j := i + 1; j < len(group); j++ {
			similarity := fs.calculateResponseSimilarity(group[i].Response, group[j].Response)
			similarities = append(similarities, similarity)
		}
	}

	if len(similarities) == 0 {
		return 1.0
	}

	var sum float64
	for _, sim := range similarities {
		sum += sim
	}

	return sum / float64(len(similarities))
}

// calculateResponseSimilarity calcula similaridade entre duas respostas
func (fs *FeedbackSystem) calculateResponseSimilarity(response1, response2 string) float64 {
	// Implementação muito básica - na prática usaria algoritmos mais sofisticados
	words1 := strings.Fields(strings.ToLower(response1))
	words2 := strings.Fields(strings.ToLower(response2))
	
	common := 0
	total := len(words1) + len(words2)
	
	if total == 0 {
		return 1.0
	}

	wordSet1 := make(map[string]bool)
	for _, word := range words1 {
		wordSet1[word] = true
	}

	for _, word := range words2 {
		if wordSet1[word] {
			common++
		}
	}

	return float64(common*2) / float64(total)
}

// identifyPatterns identifica padrões nos feedbacks
func (fs *FeedbackSystem) identifyPatterns() {
	// Identificar padrões de problemas comuns
	problemPatterns := make(map[string]int)
	improvementSuggestions := make(map[string][]string)

	for _, fb := range fs.feedback {
		if fb.UserRating <= 2 { // Ratings baixos
			// Analisar correções e feedback para identificar padrões
			if fb.Corrections != "" {
				pattern := fs.extractPattern(fb.Corrections)
				problemPatterns[pattern]++
				improvementSuggestions[pattern] = append(improvementSuggestions[pattern], fb.Corrections)
			}
		}
	}

	// Atualizar padrões identificados
	for pattern, frequency := range problemPatterns {
		fs.updatePattern(pattern, frequency, improvementSuggestions[pattern])
	}
}

// extractPattern extrai padrão de problema do feedback
func (fs *FeedbackSystem) extractPattern(correction string) string {
	correctionLower := strings.ToLower(correction)
	
	// Padrões comuns identificados
	patterns := map[string]string{
		"fonte":       "Falta citação de fonte",
		"cálculo":     "Erro no cálculo",
		"valor":       "Valor incorreto", 
		"política":    "Política não citada",
		"exemplo":     "Falta exemplo prático",
		"confuso":     "Resposta confusa",
		"incompleto":  "Resposta incompleta",
	}

	for keyword, patternName := range patterns {
		if strings.Contains(correctionLower, keyword) {
			return patternName
		}
	}

	return "Outros problemas"
}

// updatePattern atualiza ou cria um padrão
func (fs *FeedbackSystem) updatePattern(patternName string, frequency int, examples []string) {
	// Procurar padrão existente
	for i, pattern := range fs.patterns {
		if pattern.Pattern == patternName {
			fs.patterns[i].Frequency += frequency
			fs.patterns[i].LastSeen = time.Now()
			
			// Adicionar novos exemplos (limitar a 5)
			for _, example := range examples {
				if len(fs.patterns[i].Examples) < 5 {
					fs.patterns[i].Examples = append(fs.patterns[i].Examples, example)
				}
			}
			return
		}
	}

	// Criar novo padrão
	newPattern := LearningPattern{
		Pattern:     patternName,
		Frequency:   frequency,
		Confidence:  float64(frequency) / float64(len(fs.feedback)),
		Category:    fs.categorizePattern(patternName),
		Examples:    examples,
		LastSeen:    time.Now(),
		Improvement: fs.generateImprovement(patternName),
	}

	fs.patterns = append(fs.patterns, newPattern)
}

// categorizePattern categoriza um padrão
func (fs *FeedbackSystem) categorizePattern(pattern string) string {
	categories := map[string]string{
		"Falta citação de fonte": "Qualidade",
		"Erro no cálculo":       "Precisão",
		"Valor incorreto":       "Precisão", 
		"Política não citada":   "Qualidade",
		"Falta exemplo prático": "Completude",
		"Resposta confusa":      "Clareza",
		"Resposta incompleta":   "Completude",
	}

	if category, exists := categories[pattern]; exists {
		return category
	}
	return "Geral"
}

// generateImprovement gera sugestão de melhoria para um padrão
func (fs *FeedbackSystem) generateImprovement(pattern string) string {
	improvements := map[string]string{
		"Falta citação de fonte": "Sempre incluir referência à política específica (ex: Política VR-001)",
		"Erro no cálculo":       "Verificar fórmulas e valores antes de responder",
		"Valor incorreto":       "Conferir tabela de valores por sindicato",
		"Política não citada":   "Citar política específica relevante à pergunta",
		"Falta exemplo prático": "Incluir exemplo numérico quando relevante",
		"Resposta confusa":      "Estruturar resposta de forma mais clara e objetiva",
		"Resposta incompleta":   "Verificar se todos os aspectos da pergunta foram abordados",
	}

	if improvement, exists := improvements[pattern]; exists {
		return improvement
	}
	return "Analisar feedback específico para identificar melhoria"
}

// GetQualityMetrics retorna métricas de qualidade atuais
func (fs *FeedbackSystem) GetQualityMetrics() QualityMetrics {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()
	return fs.metrics
}

// GetLearningPatterns retorna padrões de aprendizado ordenados por frequência
func (fs *FeedbackSystem) GetLearningPatterns() []LearningPattern {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	// Ordenar por frequência (mais comum primeiro)
	patterns := make([]LearningPattern, len(fs.patterns))
	copy(patterns, fs.patterns)
	
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Frequency > patterns[j].Frequency
	})

	return patterns
}

// GetFeedbackSummary retorna resumo do feedback
func (fs *FeedbackSystem) GetFeedbackSummary() map[string]interface{} {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	summary := make(map[string]interface{})
	
	if len(fs.feedback) == 0 {
		summary["total_feedback"] = 0
		return summary
	}

	summary["total_feedback"] = len(fs.feedback)
	summary["average_rating"] = fs.metrics.AverageRating
	summary["accuracy_score"] = fs.metrics.AccuracyScore

	// Distribuição por rating
	ratingDist := make(map[int]int)
	for _, fb := range fs.feedback {
		ratingDist[fb.UserRating]++
	}
	summary["rating_distribution"] = ratingDist

	// Categorias mais problemáticas
	categoryProblems := make(map[string]int)
	for _, pattern := range fs.patterns {
		categoryProblems[pattern.Category] += pattern.Frequency
	}
	summary["problem_categories"] = categoryProblems

	summary["top_patterns"] = fs.getTopPatterns(5)
	summary["last_updated"] = fs.metrics.LastUpdated

	return summary
}

// getTopPatterns retorna os N padrões mais comuns
func (fs *FeedbackSystem) getTopPatterns(n int) []LearningPattern {
	patterns := fs.GetLearningPatterns()
	if len(patterns) <= n {
		return patterns
	}
	return patterns[:n]
}

// GenerateImprovementReport gera relatório de melhorias sugeridas
func (fs *FeedbackSystem) GenerateImprovementReport() map[string]interface{} {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	report := make(map[string]interface{})
	
	// Principais áreas de melhoria
	topPatterns := fs.getTopPatterns(5)
	improvements := make([]map[string]interface{}, 0)
	
	for _, pattern := range topPatterns {
		improvement := map[string]interface{}{
			"area":        pattern.Pattern,
			"frequency":   pattern.Frequency,
			"impact":      pattern.Category,
			"suggestion":  pattern.Improvement,
			"examples":    pattern.Examples,
			"confidence":  pattern.Confidence,
		}
		improvements = append(improvements, improvement)
	}

	report["priority_improvements"] = improvements
	report["overall_score"] = fs.metrics.AverageRating
	report["accuracy_trend"] = fs.calculateAccuracyTrend()
	report["generated_at"] = time.Now()

	return report
}

// calculateAccuracyTrend calcula tendência de precisão
func (fs *FeedbackSystem) calculateAccuracyTrend() string {
	if len(fs.feedback) < 10 {
		return "Dados insuficientes"
	}

	// Separar feedback em duas metades (mais antigo vs mais recente)
	mid := len(fs.feedback) / 2
	
	var oldSum, newSum float64
	for i := 0; i < mid; i++ {
		if fs.feedback[i].UserRating >= 4 {
			oldSum += 1.0
		}
	}
	for i := mid; i < len(fs.feedback); i++ {
		if fs.feedback[i].UserRating >= 4 {
			newSum += 1.0
		}
	}

	oldAvg := oldSum / float64(mid)
	newAvg := newSum / float64(len(fs.feedback)-mid)

	if newAvg > oldAvg+0.1 {
		return "Melhorando"
	} else if newAvg < oldAvg-0.1 {
		return "Piorando" 
	} else {
		return "Estável"
	}
}