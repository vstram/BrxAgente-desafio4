package agent

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"
)

// QualityMetricsCollector coleta e analisa métricas de qualidade em tempo real
type QualityMetricsCollector struct {
	mu                sync.RWMutex
	metrics           *QualityMetricsData
	historyWindow     time.Duration
	maxHistorySize    int
	thresholds        QualityThresholds
	alertCallbacks    []AlertCallback
	enabled           bool
	startTime         time.Time
	logger            Logger
}

// QualityMetricsData armazena dados de métricas de qualidade
type QualityMetricsData struct {
	// Métricas de Precisão
	TotalQuestions        int64                 `json:"total_questions"`
	CorrectAnswers        int64                 `json:"correct_answers"`
	PartialAnswers        int64                 `json:"partial_answers"`
	IncorrectAnswers      int64                 `json:"incorrect_answers"`
	CurrentAccuracy       float64               `json:"current_accuracy"`
	
	// Métricas de Satisfação
	SatisfactionRatings   []float64             `json:"satisfaction_ratings"`
	AverageSatisfaction   float64               `json:"average_satisfaction"`
	SatisfactionTrend     []TimestampedValue    `json:"satisfaction_trend"`
	
	// Métricas de Performance
	ResponseTimes         []time.Duration       `json:"response_times"`
	AverageResponseTime   time.Duration         `json:"average_response_time"`
	MedianResponseTime    time.Duration         `json:"median_response_time"`
	P95ResponseTime       time.Duration         `json:"p95_response_time"`
	P99ResponseTime       time.Duration         `json:"p99_response_time"`
	
	// Cache Performance
	CacheHits             int64                 `json:"cache_hits"`
	CacheMisses           int64                 `json:"cache_misses"`
	CacheHitRate          float64               `json:"cache_hit_rate"`
	
	// Métricas por Categoria
	CategoryMetrics       map[string]CategoryMetrics `json:"category_metrics"`
	
	// Métricas de Qualidade da Resposta
	AvgResponseLength     float64               `json:"avg_response_length"`
	PolicyAdherenceRate   float64               `json:"policy_adherence_rate"`
	FormattingQualityAvg  float64               `json:"formatting_quality_avg"`
	
	// Dados históricos
	HourlyStats           []HourlyMetrics       `json:"hourly_stats"`
	DailyStats            []DailyMetrics        `json:"daily_stats"`
	
	// Alertas e Anomalias
	ActiveAlerts          []QualityAlert        `json:"active_alerts"`
	AnomaliesDetected     int64                 `json:"anomalies_detected"`
	
	LastUpdated           time.Time             `json:"last_updated"`
}

// TimestampedValue representa um valor com timestamp
type TimestampedValue struct {
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// CategoryMetrics métricas específicas por categoria de pergunta
type CategoryMetrics struct {
	Category            string        `json:"category"`
	TotalQuestions      int64         `json:"total_questions"`
	CorrectAnswers      int64         `json:"correct_answers"`
	Accuracy            float64       `json:"accuracy"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	AverageScore        float64       `json:"average_score"`
	CommonFailures      []string      `json:"common_failures"`
	LastUpdated         time.Time     `json:"last_updated"`
}

// HourlyMetrics estatísticas por hora
type HourlyMetrics struct {
	Hour                time.Time `json:"hour"`
	QuestionsCount      int64     `json:"questions_count"`
	Accuracy            float64   `json:"accuracy"`
	AvgResponseTime     time.Duration `json:"avg_response_time"`
	SatisfactionAvg     float64   `json:"satisfaction_avg"`
	ErrorRate           float64   `json:"error_rate"`
}

// DailyMetrics estatísticas diárias
type DailyMetrics struct {
	Date                time.Time `json:"date"`
	QuestionsCount      int64     `json:"questions_count"`
	Accuracy            float64   `json:"accuracy"`
	AvgResponseTime     time.Duration `json:"avg_response_time"`
	SatisfactionAvg     float64   `json:"satisfaction_avg"`
	CacheHitRate        float64   `json:"cache_hit_rate"`
	PolicyAdherence     float64   `json:"policy_adherence"`
}

// QualityThresholds define limites para alertas
type QualityThresholds struct {
	MinAccuracy         float64       `json:"min_accuracy"`          // 0.95
	MaxResponseTime     time.Duration `json:"max_response_time"`     // 2s
	MinSatisfaction     float64       `json:"min_satisfaction"`      // 4.2
	MinCacheHitRate     float64       `json:"min_cache_hit_rate"`    // 0.70
	MaxErrorRate        float64       `json:"max_error_rate"`        // 0.05
	MinPolicyAdherence  float64       `json:"min_policy_adherence"`  // 0.90
}

// QualityAlert representa um alerta de qualidade
type QualityAlert struct {
	ID          string              `json:"id"`
	Type        AlertType           `json:"type"`
	Severity    AlertSeverity       `json:"severity"`
	Message     string              `json:"message"`
	Metric      string              `json:"metric"`
	CurrentValue any        `json:"current_value"`
	ThresholdValue any      `json:"threshold_value"`
	TriggeredAt time.Time           `json:"triggered_at"`
	AcknowledgedAt *time.Time       `json:"acknowledged_at,omitempty"`
	ResolvedAt  *time.Time          `json:"resolved_at,omitempty"`
	Metadata    map[string]any `json:"metadata"`
}

// AlertType tipos de alerta
type AlertType string

const (
	AlertAccuracyDrop      AlertType = "accuracy_drop"
	AlertHighResponseTime  AlertType = "high_response_time"
	AlertLowSatisfaction   AlertType = "low_satisfaction"
	AlertCachePerformance  AlertType = "cache_performance"
	AlertErrorSpike        AlertType = "error_spike"
	AlertAnomalyDetected   AlertType = "anomaly_detected"
)

// AlertSeverity níveis de severidade
type AlertSeverity string

const (
	SeverityLow      AlertSeverity = "low"
	SeverityMedium   AlertSeverity = "medium"
	SeverityHigh     AlertSeverity = "high"
	SeverityCritical AlertSeverity = "critical"
)

// AlertCallback função callback para alertas
type AlertCallback func(alert QualityAlert)

// NewQualityMetricsCollector cria um novo coletor de métricas
func NewQualityMetricsCollector(logger Logger) *QualityMetricsCollector {
	return &QualityMetricsCollector{
		metrics: &QualityMetricsData{
			CategoryMetrics: make(map[string]CategoryMetrics),
			HourlyStats:     []HourlyMetrics{},
			DailyStats:      []DailyMetrics{},
			ActiveAlerts:    []QualityAlert{},
			LastUpdated:     time.Now(),
		},
		historyWindow:  24 * time.Hour,
		maxHistorySize: 1000,
		thresholds: QualityThresholds{
			MinAccuracy:         0.95,
			MaxResponseTime:     2 * time.Second,
			MinSatisfaction:     4.2,
			MinCacheHitRate:     0.70,
			MaxErrorRate:        0.05,
			MinPolicyAdherence:  0.90,
		},
		alertCallbacks: []AlertCallback{},
		enabled:        true,
		startTime:      time.Now(),
		logger:         logger,
	}
}

// Enable habilita coleta de métricas
func (qmc *QualityMetricsCollector) Enable() {
	qmc.mu.Lock()
	defer qmc.mu.Unlock()
	qmc.enabled = true
	qmc.logger.Println("Quality metrics collection enabled")
}

// Disable desabilita coleta de métricas
func (qmc *QualityMetricsCollector) Disable() {
	qmc.mu.Lock()
	defer qmc.mu.Unlock()
	qmc.enabled = false
	qmc.logger.Println("Quality metrics collection disabled")
}

// SetThresholds atualiza os thresholds para alertas
func (qmc *QualityMetricsCollector) SetThresholds(thresholds QualityThresholds) {
	qmc.mu.Lock()
	defer qmc.mu.Unlock()
	qmc.thresholds = thresholds
	qmc.logger.Printf("Quality thresholds updated: %+v", thresholds)
}

// AddAlertCallback adiciona callback para alertas
func (qmc *QualityMetricsCollector) AddAlertCallback(callback AlertCallback) {
	qmc.mu.Lock()
	defer qmc.mu.Unlock()
	qmc.alertCallbacks = append(qmc.alertCallbacks, callback)
}

// RecordQuestion registra uma nova pergunta e resposta
func (qmc *QualityMetricsCollector) RecordQuestion(question, response string, responseTime time.Duration, category string, score float64, correct bool) {
	if !qmc.enabled {
		return
	}
	
	qmc.mu.Lock()
	defer qmc.mu.Unlock()
	
	// Atualizar métricas gerais
	qmc.metrics.TotalQuestions++
	qmc.metrics.ResponseTimes = append(qmc.metrics.ResponseTimes, responseTime)
	
	// Categorizar resposta
	if correct {
		qmc.metrics.CorrectAnswers++
	} else if score >= 3.0 {
		qmc.metrics.PartialAnswers++
	} else {
		qmc.metrics.IncorrectAnswers++
	}
	
	// Atualizar accuracy
	qmc.metrics.CurrentAccuracy = float64(qmc.metrics.CorrectAnswers) / float64(qmc.metrics.TotalQuestions)
	
	// Atualizar métricas por categoria
	qmc.updateCategoryMetrics(category, responseTime, score, correct)
	
	// Calcular estatísticas de tempo de resposta
	qmc.updateResponseTimeStats()
	
	// Atualizar métricas de qualidade
	qmc.updateQualityMetrics(response)
	
	// Verificar thresholds e gerar alertas
	qmc.checkThresholds()
	
	qmc.metrics.LastUpdated = time.Now()
}

// RecordSatisfaction registra uma avaliação de satisfação
func (qmc *QualityMetricsCollector) RecordSatisfaction(rating float64) {
	if !qmc.enabled {
		return
	}
	
	qmc.mu.Lock()
	defer qmc.mu.Unlock()
	
	qmc.metrics.SatisfactionRatings = append(qmc.metrics.SatisfactionRatings, rating)
	qmc.metrics.SatisfactionTrend = append(qmc.metrics.SatisfactionTrend, TimestampedValue{
		Value:     rating,
		Timestamp: time.Now(),
	})
	
	// Calcular média móvel de satisfação
	qmc.updateSatisfactionAverage()
	
	// Manter apenas dados dentro da janela de tempo
	qmc.cleanupOldData()
}

// RecordCacheEvent registra evento de cache
func (qmc *QualityMetricsCollector) RecordCacheEvent(hit bool) {
	if !qmc.enabled {
		return
	}
	
	qmc.mu.Lock()
	defer qmc.mu.Unlock()
	
	if hit {
		qmc.metrics.CacheHits++
	} else {
		qmc.metrics.CacheMisses++
	}
	
	total := qmc.metrics.CacheHits + qmc.metrics.CacheMisses
	if total > 0 {
		qmc.metrics.CacheHitRate = float64(qmc.metrics.CacheHits) / float64(total)
	}
}

// updateCategoryMetrics atualiza métricas por categoria
func (qmc *QualityMetricsCollector) updateCategoryMetrics(category string, responseTime time.Duration, score float64, correct bool) {
	catMetrics, exists := qmc.metrics.CategoryMetrics[category]
	if !exists {
		catMetrics = CategoryMetrics{
			Category:       category,
			CommonFailures: []string{},
		}
	}
	
	catMetrics.TotalQuestions++
	if correct {
		catMetrics.CorrectAnswers++
	}
	
	// Atualizar accuracy
	catMetrics.Accuracy = float64(catMetrics.CorrectAnswers) / float64(catMetrics.TotalQuestions)
	
	// Atualizar tempo médio de resposta
	if catMetrics.AverageResponseTime == 0 {
		catMetrics.AverageResponseTime = responseTime
	} else {
		catMetrics.AverageResponseTime = time.Duration(
			(int64(catMetrics.AverageResponseTime) + int64(responseTime)) / 2,
		)
	}
	
	// Atualizar score médio
	if catMetrics.AverageScore == 0 {
		catMetrics.AverageScore = score
	} else {
		catMetrics.AverageScore = (catMetrics.AverageScore + score) / 2
	}
	
	catMetrics.LastUpdated = time.Now()
	qmc.metrics.CategoryMetrics[category] = catMetrics
}

// updateResponseTimeStats calcula estatísticas de tempo de resposta
func (qmc *QualityMetricsCollector) updateResponseTimeStats() {
	if len(qmc.metrics.ResponseTimes) == 0 {
		return
	}
	
	// Criar cópia para não alterar ordem original
	times := make([]time.Duration, len(qmc.metrics.ResponseTimes))
	copy(times, qmc.metrics.ResponseTimes)
	
	// Ordenar para calcular percentis
	slices.SortFunc(times, func(a, b time.Duration) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	
	// Calcular estatísticas
	var total time.Duration
	for _, t := range times {
		total += t
	}
	
	qmc.metrics.AverageResponseTime = total / time.Duration(len(times))
	qmc.metrics.MedianResponseTime = times[len(times)/2]
	qmc.metrics.P95ResponseTime = times[int(float64(len(times))*0.95)]
	qmc.metrics.P99ResponseTime = times[int(float64(len(times))*0.99)]
}

// updateSatisfactionAverage calcula média de satisfação
func (qmc *QualityMetricsCollector) updateSatisfactionAverage() {
	if len(qmc.metrics.SatisfactionRatings) == 0 {
		return
	}
	
	var sum float64
	for _, rating := range qmc.metrics.SatisfactionRatings {
		sum += rating
	}
	
	qmc.metrics.AverageSatisfaction = sum / float64(len(qmc.metrics.SatisfactionRatings))
}

// updateQualityMetrics atualiza métricas de qualidade da resposta
func (qmc *QualityMetricsCollector) updateQualityMetrics(response string) {
	responseLength := float64(len(response))
	
	// Atualizar comprimento médio de resposta
	if qmc.metrics.AvgResponseLength == 0 {
		qmc.metrics.AvgResponseLength = responseLength
	} else {
		qmc.metrics.AvgResponseLength = (qmc.metrics.AvgResponseLength + responseLength) / 2
	}
	
	// Verificar aderência a políticas (simplified)
	policyKeywords := []string{"política", "regra", "norma", "procedimento"}
	hasPolicy := false
	for _, keyword := range policyKeywords {
		if contains(response, keyword) {
			hasPolicy = true
			break
		}
	}
	
	if hasPolicy {
		if qmc.metrics.PolicyAdherenceRate == 0 {
			qmc.metrics.PolicyAdherenceRate = 1.0
		} else {
			qmc.metrics.PolicyAdherenceRate = (qmc.metrics.PolicyAdherenceRate + 1.0) / 2
		}
	} else {
		if qmc.metrics.PolicyAdherenceRate > 0 {
			qmc.metrics.PolicyAdherenceRate = qmc.metrics.PolicyAdherenceRate * 0.9
		}
	}
}

// checkThresholds verifica thresholds e gera alertas
func (qmc *QualityMetricsCollector) checkThresholds() {
	alerts := []QualityAlert{}
	
	// Verificar accuracy
	if qmc.metrics.CurrentAccuracy < qmc.thresholds.MinAccuracy {
		alerts = append(alerts, QualityAlert{
			ID:          fmt.Sprintf("accuracy_%d", time.Now().Unix()),
			Type:        AlertAccuracyDrop,
			Severity:    qmc.determineSeverity(qmc.metrics.CurrentAccuracy, qmc.thresholds.MinAccuracy),
			Message:     fmt.Sprintf("Accuracy dropped to %.2f%% (threshold: %.2f%%)", qmc.metrics.CurrentAccuracy*100, qmc.thresholds.MinAccuracy*100),
			Metric:      "accuracy",
			CurrentValue: qmc.metrics.CurrentAccuracy,
			ThresholdValue: qmc.thresholds.MinAccuracy,
			TriggeredAt: time.Now(),
		})
	}
	
	// Verificar tempo de resposta
	if qmc.metrics.AverageResponseTime > qmc.thresholds.MaxResponseTime {
		alerts = append(alerts, QualityAlert{
			ID:          fmt.Sprintf("response_time_%d", time.Now().Unix()),
			Type:        AlertHighResponseTime,
			Severity:    qmc.determineTimeSeverity(qmc.metrics.AverageResponseTime, qmc.thresholds.MaxResponseTime),
			Message:     fmt.Sprintf("Average response time is %v (threshold: %v)", qmc.metrics.AverageResponseTime, qmc.thresholds.MaxResponseTime),
			Metric:      "response_time",
			CurrentValue: qmc.metrics.AverageResponseTime,
			ThresholdValue: qmc.thresholds.MaxResponseTime,
			TriggeredAt: time.Now(),
		})
	}
	
	// Verificar satisfação
	if qmc.metrics.AverageSatisfaction < qmc.thresholds.MinSatisfaction {
		alerts = append(alerts, QualityAlert{
			ID:          fmt.Sprintf("satisfaction_%d", time.Now().Unix()),
			Type:        AlertLowSatisfaction,
			Severity:    qmc.determineSeverity(qmc.metrics.AverageSatisfaction, qmc.thresholds.MinSatisfaction),
			Message:     fmt.Sprintf("User satisfaction is %.1f/5.0 (threshold: %.1f)", qmc.metrics.AverageSatisfaction, qmc.thresholds.MinSatisfaction),
			Metric:      "satisfaction",
			CurrentValue: qmc.metrics.AverageSatisfaction,
			ThresholdValue: qmc.thresholds.MinSatisfaction,
			TriggeredAt: time.Now(),
		})
	}
	
	// Verificar cache hit rate
	if qmc.metrics.CacheHitRate < qmc.thresholds.MinCacheHitRate {
		alerts = append(alerts, QualityAlert{
			ID:          fmt.Sprintf("cache_%d", time.Now().Unix()),
			Type:        AlertCachePerformance,
			Severity:    qmc.determineSeverity(qmc.metrics.CacheHitRate, qmc.thresholds.MinCacheHitRate),
			Message:     fmt.Sprintf("Cache hit rate is %.1f%% (threshold: %.1f%%)", qmc.metrics.CacheHitRate*100, qmc.thresholds.MinCacheHitRate*100),
			Metric:      "cache_hit_rate",
			CurrentValue: qmc.metrics.CacheHitRate,
			ThresholdValue: qmc.thresholds.MinCacheHitRate,
			TriggeredAt: time.Now(),
		})
	}
	
	// Adicionar novos alertas e notificar callbacks
	for _, alert := range alerts {
		qmc.metrics.ActiveAlerts = append(qmc.metrics.ActiveAlerts, alert)
		for _, callback := range qmc.alertCallbacks {
			go callback(alert)
		}
	}
}

// determineSeverity determina severidade baseada na distância do threshold
func (qmc *QualityMetricsCollector) determineSeverity(current, threshold float64) AlertSeverity {
	ratio := current / threshold
	
	if ratio < 0.8 {
		return SeverityCritical
	} else if ratio < 0.9 {
		return SeverityHigh
	} else if ratio < 0.95 {
		return SeverityMedium
	}
	return SeverityLow
}

// determineTimeSeverity determina severidade para tempo de resposta
func (qmc *QualityMetricsCollector) determineTimeSeverity(current, threshold time.Duration) AlertSeverity {
	ratio := float64(current) / float64(threshold)
	
	if ratio > 3.0 {
		return SeverityCritical
	} else if ratio > 2.0 {
		return SeverityHigh
	} else if ratio > 1.5 {
		return SeverityMedium
	}
	return SeverityLow
}

// cleanupOldData remove dados antigos para evitar crescimento excessivo
func (qmc *QualityMetricsCollector) cleanupOldData() {
	cutoff := time.Now().Add(-qmc.historyWindow)
	
	// Limpar satisfaction trend
	var validTrend []TimestampedValue
	for _, item := range qmc.metrics.SatisfactionTrend {
		if item.Timestamp.After(cutoff) {
			validTrend = append(validTrend, item)
		}
	}
	qmc.metrics.SatisfactionTrend = validTrend
	
	// Limpar response times se muito grande
	if len(qmc.metrics.ResponseTimes) > qmc.maxHistorySize {
		qmc.metrics.ResponseTimes = qmc.metrics.ResponseTimes[len(qmc.metrics.ResponseTimes)-qmc.maxHistorySize:]
	}
	
	// Limpar satisfaction ratings se muito grande
	if len(qmc.metrics.SatisfactionRatings) > qmc.maxHistorySize {
		qmc.metrics.SatisfactionRatings = qmc.metrics.SatisfactionRatings[len(qmc.metrics.SatisfactionRatings)-qmc.maxHistorySize:]
	}
	
	// Limpar alertas resolvidos antigos
	var activeAlerts []QualityAlert
	for _, alert := range qmc.metrics.ActiveAlerts {
		if alert.ResolvedAt == nil || alert.ResolvedAt.After(cutoff) {
			activeAlerts = append(activeAlerts, alert)
		}
	}
	qmc.metrics.ActiveAlerts = activeAlerts
}

// GetMetrics retorna snapshot das métricas atuais
func (qmc *QualityMetricsCollector) GetMetrics() QualityMetricsData {
	qmc.mu.RLock()
	defer qmc.mu.RUnlock()
	
	// Criar cópia para evitar data races
	data := *qmc.metrics
	return data
}

// GetCategoryMetrics retorna métricas de uma categoria específica
func (qmc *QualityMetricsCollector) GetCategoryMetrics(category string) (CategoryMetrics, bool) {
	qmc.mu.RLock()
	defer qmc.mu.RUnlock()
	
	metrics, exists := qmc.metrics.CategoryMetrics[category]
	return metrics, exists
}

// GetActiveAlerts retorna alertas ativos
func (qmc *QualityMetricsCollector) GetActiveAlerts() []QualityAlert {
	qmc.mu.RLock()
	defer qmc.mu.RUnlock()
	
	alerts := make([]QualityAlert, len(qmc.metrics.ActiveAlerts))
	copy(alerts, qmc.metrics.ActiveAlerts)
	return alerts
}

// AcknowledgeAlert marca um alerta como reconhecido
func (qmc *QualityMetricsCollector) AcknowledgeAlert(alertID string) bool {
	qmc.mu.Lock()
	defer qmc.mu.Unlock()
	
	for i := range qmc.metrics.ActiveAlerts {
		if qmc.metrics.ActiveAlerts[i].ID == alertID {
			now := time.Now()
			qmc.metrics.ActiveAlerts[i].AcknowledgedAt = &now
			return true
		}
	}
	return false
}

// ResolveAlert marca um alerta como resolvido
func (qmc *QualityMetricsCollector) ResolveAlert(alertID string) bool {
	qmc.mu.Lock()
	defer qmc.mu.Unlock()
	
	for i := range qmc.metrics.ActiveAlerts {
		if qmc.metrics.ActiveAlerts[i].ID == alertID {
			now := time.Now()
			qmc.metrics.ActiveAlerts[i].ResolvedAt = &now
			return true
		}
	}
	return false
}

// ExportMetrics exporta métricas em formato JSON
func (qmc *QualityMetricsCollector) ExportMetrics() ([]byte, error) {
	qmc.mu.RLock()
	defer qmc.mu.RUnlock()
	
	return json.MarshalIndent(qmc.metrics, "", "  ")
}

// GenerateReport gera relatório de qualidade
func (qmc *QualityMetricsCollector) GenerateReport() QualityReport {
	qmc.mu.RLock()
	defer qmc.mu.RUnlock()
	
	report := QualityReport{
		GeneratedAt:       time.Now(),
		DataPeriod:        time.Since(qmc.startTime),
		TotalQuestions:    qmc.metrics.TotalQuestions,
		OverallAccuracy:   qmc.metrics.CurrentAccuracy,
		AverageResponseTime: qmc.metrics.AverageResponseTime,
		AverageSatisfaction: qmc.metrics.AverageSatisfaction,
		CacheHitRate:      qmc.metrics.CacheHitRate,
		ActiveAlertsCount: len(qmc.metrics.ActiveAlerts),
		Recommendations:   qmc.generateRecommendations(),
	}
	
	// Adicionar top 5 categorias por volume
	report.TopCategories = qmc.getTopCategories(5)
	
	return report
}

// QualityReport relatório de qualidade resumido
type QualityReport struct {
	GeneratedAt         time.Time                `json:"generated_at"`
	DataPeriod          time.Duration            `json:"data_period"`
	TotalQuestions      int64                    `json:"total_questions"`
	OverallAccuracy     float64                  `json:"overall_accuracy"`
	AverageResponseTime time.Duration            `json:"average_response_time"`
	AverageSatisfaction float64                  `json:"average_satisfaction"`
	CacheHitRate        float64                  `json:"cache_hit_rate"`
	ActiveAlertsCount   int                      `json:"active_alerts_count"`
	TopCategories       []CategoryMetrics        `json:"top_categories"`
	Recommendations     []string                 `json:"recommendations"`
}

// generateRecommendations gera recomendações baseadas nas métricas atuais
func (qmc *QualityMetricsCollector) generateRecommendations() []string {
	var recommendations []string
	
	if qmc.metrics.CurrentAccuracy < 0.95 {
		recommendations = append(recommendations, "Consider reviewing and expanding the knowledge base to improve accuracy")
	}
	
	if qmc.metrics.AverageResponseTime > 2*time.Second {
		recommendations = append(recommendations, "Response time is above target - consider performance optimizations")
	}
	
	if qmc.metrics.AverageSatisfaction < 4.0 {
		recommendations = append(recommendations, "User satisfaction is below target - review response quality and format")
	}
	
	if qmc.metrics.CacheHitRate < 0.7 {
		recommendations = append(recommendations, "Cache hit rate is low - review caching strategy")
	}
	
	if len(qmc.metrics.ActiveAlerts) > 3 {
		recommendations = append(recommendations, "Multiple active alerts - prioritize resolution of critical issues")
	}
	
	return recommendations
}

// getTopCategories retorna top N categorias por volume
func (qmc *QualityMetricsCollector) getTopCategories(n int) []CategoryMetrics {
	var categories []CategoryMetrics
	for _, cat := range qmc.metrics.CategoryMetrics {
		categories = append(categories, cat)
	}
	
	// Ordenar por número de perguntas
	slices.SortFunc(categories, func(a, b CategoryMetrics) int {
		if a.TotalQuestions > b.TotalQuestions {
			return -1
		} else if a.TotalQuestions < b.TotalQuestions {
			return 1
		}
		return 0
	})
	
	if len(categories) > n {
		categories = categories[:n]
	}
	
	return categories
}

// contains helper function para verificar se string contém substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || len(s) > 0 && (s[0:len(substr)] == substr || strings.Contains(s[1:], substr)))
}