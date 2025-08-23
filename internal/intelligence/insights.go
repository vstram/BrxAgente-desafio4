package intelligence

import (
	"fmt"
	"log"
	"sort"
	"time"
)

// InsightType define os tipos de insights que podem ser gerados
type InsightType string

const (
	InsightTypeFinancial    InsightType = "financial"
	InsightTypeOperational  InsightType = "operational"
	InsightTypeQuality      InsightType = "quality"
	InsightTypeAnomaly      InsightType = "anomaly"
	InsightTypeTrend        InsightType = "trend"
)

// InsightPriority define a prioridade dos insights
type InsightPriority int

const (
	PriorityLow InsightPriority = iota
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

func (p InsightPriority) String() string {
	switch p {
	case PriorityLow:
		return "baixa"
	case PriorityMedium:
		return "média"
	case PriorityHigh:
		return "alta"
	case PriorityCritical:
		return "crítica"
	default:
		return "desconhecida"
	}
}

// Insight representa um insight gerado automaticamente
type Insight struct {
	Type        InsightType     `json:"type"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Impact      string          `json:"impact"`
	Action      string          `json:"action"`
	Priority    InsightPriority `json:"priority"`
	Confidence  float64         `json:"confidence"`
	GeneratedAt time.Time       `json:"generated_at"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// InsightGenerator gera insights automáticos a partir dos dados processados
type InsightGenerator struct {
	config *InsightConfig
	logger *log.Logger
}

// InsightConfig configurações para geração de insights
type InsightConfig struct {
	MinConfidenceThreshold float64
	MaxInsightsPerType     int
	EnableTrendAnalysis    bool
	EnableAnomalyInsights  bool
	ComparisonPeriodDays   int
}

// DefaultInsightConfig retorna configuração padrão para insights
func DefaultInsightConfig() *InsightConfig {
	return &InsightConfig{
		MinConfidenceThreshold: 0.7,
		MaxInsightsPerType:     5,
		EnableTrendAnalysis:    true,
		EnableAnomalyInsights:  true,
		ComparisonPeriodDays:   30,
	}
}

// NewInsightGenerator cria novo gerador de insights
func NewInsightGenerator(config *InsightConfig) *InsightGenerator {
	if config == nil {
		config = DefaultInsightConfig()
	}
	
	return &InsightGenerator{
		config: config,
		logger: log.Default(),
	}
}

// ProcessingData representa dados de um processamento
type ProcessingData struct {
	TotalCollaborators    int                    `json:"total_collaborators"`
	TotalVRValue         float64                `json:"total_vr_value"`
	ProcessingTime       time.Duration          `json:"processing_time"`
	ErrorCount           int                    `json:"error_count"`
	WarningCount         int                    `json:"warning_count"`
	AnomalyReport        *AnomalyReport         `json:"anomaly_report,omitempty"`
	SindicatoDistribution map[string]int        `json:"sindicato_distribution"`
	ValueDistribution    map[string]float64     `json:"value_distribution"`
	ProcessedAt          time.Time              `json:"processed_at"`
	Metadata             map[string]interface{} `json:"metadata"`
}

// GenerateInsights gera insights automáticos a partir dos dados de processamento
func (g *InsightGenerator) GenerateInsights(data *ProcessingData) ([]*Insight, error) {
	g.logger.Printf("Gerando insights para processamento com %d colaboradores", data.TotalCollaborators)
	
	var insights []*Insight
	
	// Insights financeiros
	financialInsights := g.generateFinancialInsights(data)
	insights = append(insights, financialInsights...)
	
	// Insights operacionais
	operationalInsights := g.generateOperationalInsights(data)
	insights = append(insights, operationalInsights...)
	
	// Insights de qualidade
	qualityInsights := g.generateQualityInsights(data)
	insights = append(insights, qualityInsights...)
	
	// Insights de anomalias (se configurado e disponível)
	if g.config.EnableAnomalyInsights && data.AnomalyReport != nil {
		anomalyInsights := g.generateAnomalyInsights(data)
		insights = append(insights, anomalyInsights...)
	}
	
	// Filtrar por confiança e ordenar por prioridade
	insights = g.filterAndSort(insights)
	
	g.logger.Printf("Gerados %d insights com confiança >= %.1f", len(insights), g.config.MinConfidenceThreshold)
	
	return insights, nil
}

// generateFinancialInsights gera insights financeiros
func (g *InsightGenerator) generateFinancialInsights(data *ProcessingData) []*Insight {
	var insights []*Insight
	
	// Insight sobre valor total
	if data.TotalVRValue > 0 {
		averageVR := data.TotalVRValue / float64(data.TotalCollaborators)
		
		insight := &Insight{
			Type:        InsightTypeFinancial,
			Title:       "Análise do Valor Total de VR",
			Description: fmt.Sprintf("Valor total de R$ %.2f para %d colaboradores (média de R$ %.2f por colaborador)", 
				data.TotalVRValue, data.TotalCollaborators, averageVR),
			Impact:      fmt.Sprintf("Representa um investimento significativo em benefícios alimentares"),
			Action:      "Monitorar variações mensais e otimizar distribuição",
			Priority:    PriorityMedium,
			Confidence:  0.95,
			GeneratedAt: time.Now(),
			Data: map[string]interface{}{
				"total_value":     data.TotalVRValue,
				"average_value":   averageVR,
				"total_employees": data.TotalCollaborators,
			},
		}
		insights = append(insights, insight)
	}
	
	// Insight sobre distribuição por sindicato
	if len(data.SindicatoDistribution) > 1 {
		maxSindicato := ""
		maxCount := 0
		for sindicato, count := range data.SindicatoDistribution {
			if count > maxCount {
				maxSindicato = sindicato
				maxCount = count
			}
		}
		
		if maxCount > 0 {
			percentage := float64(maxCount) / float64(data.TotalCollaborators) * 100
			
			insight := &Insight{
				Type:        InsightTypeFinancial,
				Title:       "Concentração por Sindicato",
				Description: fmt.Sprintf("Sindicato '%s' representa %.1f%% dos colaboradores (%d de %d)", 
					maxSindicato, percentage, maxCount, data.TotalCollaborators),
				Impact:      "Alta concentração pode indicar necessidade de revisão de políticas",
				Action:      "Avaliar equilíbrio na distribuição entre sindicatos",
				Priority:    PriorityMedium,
				Confidence:  0.90,
				GeneratedAt: time.Now(),
				Data: map[string]interface{}{
					"dominant_sindicato": maxSindicato,
					"dominant_count":     maxCount,
					"percentage":         percentage,
					"distribution":       data.SindicatoDistribution,
				},
			}
			insights = append(insights, insight)
		}
	}
	
	return insights
}

// generateOperationalInsights gera insights operacionais
func (g *InsightGenerator) generateOperationalInsights(data *ProcessingData) []*Insight {
	var insights []*Insight
	
	// Insight sobre tempo de processamento
	processingMinutes := data.ProcessingTime.Minutes()
	
	var priority InsightPriority
	var action string
	
	if processingMinutes < 60 {
		priority = PriorityLow
		action = "Manter eficiência atual do processamento"
	} else if processingMinutes < 90 {
		priority = PriorityMedium
		action = "Monitorar performance e identificar oportunidades de otimização"
	} else {
		priority = PriorityHigh
		action = "Investigar gargalos e otimizar processo de cálculo"
	}
	
	insight := &Insight{
		Type:        InsightTypeOperational,
		Title:       "Performance de Processamento",
		Description: fmt.Sprintf("Processamento concluído em %.1f minutos", processingMinutes),
		Impact:      "Impacta diretamente na produtividade da equipe de RH",
		Action:      action,
		Priority:    priority,
		Confidence:  0.95,
		GeneratedAt: time.Now(),
		Data: map[string]interface{}{
			"processing_time_minutes": processingMinutes,
			"processing_time_seconds": data.ProcessingTime.Seconds(),
			"records_per_minute":      float64(data.TotalCollaborators) / processingMinutes,
		},
	}
	insights = append(insights, insight)
	
	// Insight sobre taxa de erros
	if data.TotalCollaborators > 0 {
		errorRate := float64(data.ErrorCount) / float64(data.TotalCollaborators) * 100
		
		var errorPriority InsightPriority
		var errorAction string
		
		if errorRate == 0 {
			errorPriority = PriorityLow
			errorAction = "Excelente qualidade dos dados - manter práticas atuais"
		} else if errorRate < 5 {
			errorPriority = PriorityMedium
			errorAction = "Taxa de erros baixa - revisar casos específicos"
		} else {
			errorPriority = PriorityHigh
			errorAction = "Taxa de erros elevada - investigar qualidade dos dados de entrada"
		}
		
		errorInsight := &Insight{
			Type:        InsightTypeOperational,
			Title:       "Taxa de Erros",
			Description: fmt.Sprintf("%.1f%% de taxa de erros (%d erros em %d registros)", 
				errorRate, data.ErrorCount, data.TotalCollaborators),
			Impact:      "Erros podem resultar em cálculos incorretos e retrabalho",
			Action:      errorAction,
			Priority:    errorPriority,
			Confidence:  0.90,
			GeneratedAt: time.Now(),
			Data: map[string]interface{}{
				"error_rate":    errorRate,
				"error_count":   data.ErrorCount,
				"warning_count": data.WarningCount,
			},
		}
		insights = append(insights, errorInsight)
	}
	
	return insights
}

// generateQualityInsights gera insights de qualidade
func (g *InsightGenerator) generateQualityInsights(data *ProcessingData) []*Insight {
	var insights []*Insight
	
	// Insight sobre qualidade geral
	totalIssues := data.ErrorCount + data.WarningCount
	qualityScore := 100.0
	
	if data.TotalCollaborators > 0 {
		issueRate := float64(totalIssues) / float64(data.TotalCollaborators)
		qualityScore = 100.0 - (issueRate * 100)
	}
	
	var qualityLevel string
	var qualityPriority InsightPriority
	
	if qualityScore >= 95 {
		qualityLevel = "Excelente"
		qualityPriority = PriorityLow
	} else if qualityScore >= 85 {
		qualityLevel = "Boa"
		qualityPriority = PriorityMedium
	} else if qualityScore >= 70 {
		qualityLevel = "Regular"
		qualityPriority = PriorityHigh
	} else {
		qualityLevel = "Ruim"
		qualityPriority = PriorityCritical
	}
	
	insight := &Insight{
		Type:        InsightTypeQuality,
		Title:       "Qualidade Geral dos Dados",
		Description: fmt.Sprintf("Score de qualidade: %.1f%% (%s)", qualityScore, qualityLevel),
		Impact:      "Qualidade dos dados impacta diretamente na confiabilidade dos cálculos",
		Action:      "Implementar validações adicionais na entrada de dados",
		Priority:    qualityPriority,
		Confidence:  0.85,
		GeneratedAt: time.Now(),
		Data: map[string]interface{}{
			"quality_score": qualityScore,
			"quality_level": qualityLevel,
			"total_issues":  totalIssues,
		},
	}
	insights = append(insights, insight)
	
	return insights
}

// generateAnomalyInsights gera insights baseados no relatório de anomalias
func (g *InsightGenerator) generateAnomalyInsights(data *ProcessingData) []*Insight {
	var insights []*Insight
	
	if data.AnomalyReport == nil {
		return insights
	}
	
	report := data.AnomalyReport
	
	// Insight sobre anomalias gerais
	if report.TotalAnomalies > 0 {
		anomalyRate := float64(report.TotalAnomalies) / float64(report.TotalRecords) * 100
		
		var anomalyPriority InsightPriority
		if report.Summary.CriticalIssues > 0 {
			anomalyPriority = PriorityCritical
		} else if anomalyRate > 10 {
			anomalyPriority = PriorityHigh
		} else {
			anomalyPriority = PriorityMedium
		}
		
		insight := &Insight{
			Type:        InsightTypeAnomaly,
			Title:       "Detecção de Anomalias",
			Description: fmt.Sprintf("%.1f%% de anomalias detectadas (%d de %d registros)", 
				anomalyRate, report.TotalAnomalies, report.TotalRecords),
			Impact:      "Anomalias podem indicar problemas nos dados ou processos",
			Action:      "Revisar e corrigir anomalias identificadas",
			Priority:    anomalyPriority,
			Confidence:  0.90,
			GeneratedAt: time.Now(),
			Data: map[string]interface{}{
				"anomaly_rate":     anomalyRate,
				"total_anomalies":  report.TotalAnomalies,
				"critical_issues":  report.Summary.CriticalIssues,
				"overall_score":    report.Summary.OverallScore,
				"risk_level":       report.Summary.RiskLevel,
			},
		}
		insights = append(insights, insight)
	}
	
	// Insight sobre tipos de anomalias mais frequentes
	if len(report.AnomaliesByType) > 0 {
		type AnomalyTypeCount struct {
			Type  string
			Count int
		}
		
		var typesCounts []AnomalyTypeCount
		for anomalyType, count := range report.AnomaliesByType {
			typesCounts = append(typesCounts, AnomalyTypeCount{
				Type:  string(anomalyType),
				Count: count,
			})
		}
		
		// Ordenar por count (decrescente)
		sort.Slice(typesCounts, func(i, j int) bool {
			return typesCounts[i].Count > typesCounts[j].Count
		})
		
		if len(typesCounts) > 0 && typesCounts[0].Count > 0 {
			mostFrequentType := typesCounts[0]
			
			insight := &Insight{
				Type:        InsightTypeAnomaly,
				Title:       "Tipo de Anomalia Mais Frequente",
				Description: fmt.Sprintf("Anomalias do tipo '%s' são as mais frequentes (%d ocorrências)", 
					mostFrequentType.Type, mostFrequentType.Count),
				Impact:      "Padrão recorrente pode indicar problema sistemático",
				Action:      fmt.Sprintf("Investigar causa raiz das anomalias de tipo '%s'", mostFrequentType.Type),
				Priority:    PriorityHigh,
				Confidence:  0.85,
				GeneratedAt: time.Now(),
				Data: map[string]interface{}{
					"most_frequent_type":  mostFrequentType.Type,
					"frequency":           mostFrequentType.Count,
					"types_distribution":  report.AnomaliesByType,
				},
			}
			insights = append(insights, insight)
		}
	}
	
	return insights
}

// filterAndSort filtra insights por confiança e ordena por prioridade
func (g *InsightGenerator) filterAndSort(insights []*Insight) []*Insight {
	// Filtrar por confiança mínima
	var filtered []*Insight
	for _, insight := range insights {
		if insight.Confidence >= g.config.MinConfidenceThreshold {
			filtered = append(filtered, insight)
		}
	}
	
	// Ordenar por prioridade (decrescente) e depois por confiança (decrescente)
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Priority == filtered[j].Priority {
			return filtered[i].Confidence > filtered[j].Confidence
		}
		return filtered[i].Priority > filtered[j].Priority
	})
	
	// Limitar número de insights por tipo se configurado
	if g.config.MaxInsightsPerType > 0 {
		typeCount := make(map[InsightType]int)
		var limited []*Insight
		
		for _, insight := range filtered {
			if typeCount[insight.Type] < g.config.MaxInsightsPerType {
				limited = append(limited, insight)
				typeCount[insight.Type]++
			}
		}
		
		filtered = limited
	}
	
	return filtered
}

// FormatInsightsForHuman formata insights para leitura humana
func FormatInsightsForHuman(insights []*Insight) string {
	if len(insights) == 0 {
		return "Nenhum insight foi gerado para este processamento."
	}
	
	output := fmt.Sprintf("🧠 INSIGHTS AUTOMÁTICOS (%d)\n", len(insights))
	output += "══════════════════════════════\n\n"
	
	// Agrupar por tipo
	typeGroups := make(map[InsightType][]*Insight)
	for _, insight := range insights {
		typeGroups[insight.Type] = append(typeGroups[insight.Type], insight)
	}
	
	// Ordem de exibição dos tipos
	typeOrder := []InsightType{
		InsightTypeFinancial,
		InsightTypeOperational,
		InsightTypeQuality,
		InsightTypeAnomaly,
		InsightTypeTrend,
	}
	
	for _, insightType := range typeOrder {
		if insights, exists := typeGroups[insightType]; exists && len(insights) > 0 {
			output += fmt.Sprintf("📊 %s\n", getInsightTypeTitle(insightType))
			output += "────────────────────────────────\n"
			
			for i, insight := range insights {
				priorityEmoji := getPriorityEmoji(insight.Priority)
				confidenceBar := getConfidenceBar(insight.Confidence)
				
				output += fmt.Sprintf("%s %d. %s %s\n", priorityEmoji, i+1, insight.Title, confidenceBar)
				output += fmt.Sprintf("   📝 %s\n", insight.Description)
				if insight.Impact != "" {
					output += fmt.Sprintf("   💥 %s\n", insight.Impact)
				}
				if insight.Action != "" {
					output += fmt.Sprintf("   🎯 %s\n", insight.Action)
				}
				output += "\n"
			}
		}
	}
	
	return output
}

// Helper functions para formatação
func getInsightTypeTitle(insightType InsightType) string {
	switch insightType {
	case InsightTypeFinancial:
		return "INSIGHTS FINANCEIROS"
	case InsightTypeOperational:
		return "INSIGHTS OPERACIONAIS"
	case InsightTypeQuality:
		return "INSIGHTS DE QUALIDADE"
	case InsightTypeAnomaly:
		return "INSIGHTS DE ANOMALIAS"
	case InsightTypeTrend:
		return "INSIGHTS DE TENDÊNCIAS"
	default:
		return "INSIGHTS GERAIS"
	}
}

func getPriorityEmoji(priority InsightPriority) string {
	switch priority {
	case PriorityCritical:
		return "🔴"
	case PriorityHigh:
		return "🟠"
	case PriorityMedium:
		return "🟡"
	case PriorityLow:
		return "🟢"
	default:
		return "⚪"
	}
}

func getConfidenceBar(confidence float64) string {
	bars := int(confidence * 5) // 5 barras máximo
	filled := ""
	for i := 0; i < bars; i++ {
		filled += "█"
	}
	empty := ""
	for i := bars; i < 5; i++ {
		empty += "░"
	}
	return fmt.Sprintf("[%s%s %.0f%%]", filled, empty, confidence*100)
}