package intelligence

import (
	"fmt"
	"math"
	"sort"
	"time"

	"BrxAgente-desafio4/internal/models"
)

// RecommendationEngine gera recomendações baseadas em análise preditiva
type RecommendationEngine struct {
	config      RecommendationConfig
	predictor   *TrendPredictor
	analyzer    *PatternAnalyzer
	detector    *TrendDetector
	forecaster  *Forecaster
}

// RecommendationConfig configuração do motor de recomendações
type RecommendationConfig struct {
	MaxRecommendations   int     `json:"max_recommendations"`   // máximo de recomendações por análise
	MinConfidence        float64 `json:"min_confidence"`        // confiança mínima para gerar recomendação
	PriorityThreshold    float64 `json:"priority_threshold"`    // limite para prioridade alta
	AutoApproveThreshold float64 `json:"auto_approve_threshold"` // limite para aprovação automática
	ContextWindow        int     `json:"context_window"`        // janela de contexto em meses
}

// RecommendationSuite conjunto completo de recomendações
type RecommendationSuite struct {
	Summary         RecommendationSummary    `json:"summary"`
	Recommendations []EnhancedRecommendation `json:"recommendations"`
	ActionPlan      ActionPlan               `json:"action_plan"`
	RiskFactors     []RiskFactor            `json:"risk_factors"`
	Opportunities   []Opportunity           `json:"opportunities"`
	GeneratedAt     time.Time               `json:"generated_at"`
	ValidUntil      time.Time               `json:"valid_until"`
}

// EnhancedRecommendation recomendação aprimorada com mais contexto
type EnhancedRecommendation struct {
	models.ActionItem
	Context         RecommendationContext    `json:"context"`
	Evidence        []Evidence               `json:"evidence"`
	ExpectedImpact  ImpactAssessment         `json:"expected_impact"`
	Implementation  ImplementationGuide      `json:"implementation"`
	Dependencies    []string                 `json:"dependencies"`
	Alternatives    []Alternative            `json:"alternatives"`
	Confidence      float64                  `json:"confidence"`
	AutoApproved    bool                     `json:"auto_approved"`
}

// RecommendationContext contexto da recomendação
type RecommendationContext struct {
	SourceAnalysis   string                 `json:"source_analysis"`   // de onde veio a recomendação
	TriggerCondition string                 `json:"trigger_condition"` // condição que disparou
	AffectedEntities []string               `json:"affected_entities"` // entidades afetadas
	TimeHorizon      string                 `json:"time_horizon"`      // horizonte temporal
	Seasonality      *models.SeasonalityInfo `json:"seasonality"`      // contexto sazonal
	TrendInfo        *models.VRTrend         `json:"trend_info"`       // informação de tendência
}

// Evidence evidência que suporta a recomendação
type Evidence struct {
	Type        EvidenceType               `json:"type"`
	Source      string                     `json:"source"`
	Description string                     `json:"description"`
	Value       float64                    `json:"value"`
	Confidence  float64                    `json:"confidence"`
	Metadata    map[string]interface{}     `json:"metadata"`
}

// EvidenceType tipo de evidência
type EvidenceType string

const (
	EvidenceTrend         EvidenceType = "trend"
	EvidenceSeasonality   EvidenceType = "seasonality"
	EvidenceAnomaly       EvidenceType = "anomaly"
	EvidenceCorrelation   EvidenceType = "correlation"
	EvidenceForecast      EvidenceType = "forecast"
	EvidencePattern       EvidenceType = "pattern"
	EvidenceRisk          EvidenceType = "risk"
)

// ImpactAssessment avaliação de impacto
type ImpactAssessment struct {
	Financial    ImpactDetail `json:"financial"`
	Operational  ImpactDetail `json:"operational"`
	Strategic    ImpactDetail `json:"strategic"`
	Risk         ImpactDetail `json:"risk"`
	Timeline     string       `json:"timeline"`
	Probability  float64      `json:"probability"`
}

// ImpactDetail detalhes do impacto
type ImpactDetail struct {
	Value       float64 `json:"value"`
	Unit        string  `json:"unit"`
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence"`
}

// ImplementationGuide guia de implementação
type ImplementationGuide struct {
	Steps        []ImplementationStep `json:"steps"`
	Resources    []RequiredResource   `json:"resources"`
	Timeline     string               `json:"timeline"`
	Risks        []string             `json:"risks"`
	Success      []string             `json:"success_criteria"`
	Rollback     []string             `json:"rollback_plan"`
}

// ImplementationStep passo de implementação
type ImplementationStep struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Duration    string            `json:"duration"`
	Order       int               `json:"order"`
	Dependencies []string         `json:"dependencies"`
	Resources   []RequiredResource `json:"resources"`
	Status      string            `json:"status"`
	Progress    float64           `json:"progress"`
}

// RequiredResource recurso necessário
type RequiredResource struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
	Cost        float64 `json:"cost,omitempty"`
}

// Alternative alternativa para a recomendação
type Alternative struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Impact      ImpactAssessment `json:"impact"`
	Effort      float64          `json:"effort"`
	Risk        float64          `json:"risk"`
}

// RiskFactor fator de risco identificado
type RiskFactor struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Probability float64                `json:"probability"`
	Impact      models.ImpactLevel     `json:"impact"`
	Mitigation  []models.ActionItem    `json:"mitigation"`
	Source      string                 `json:"source"`
}

// Opportunity oportunidade identificada
type Opportunity struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Value       float64          `json:"value"`
	Effort      float64          `json:"effort"`
	Timeline    string           `json:"timeline"`
	Actions     []models.ActionItem `json:"actions"`
}

// ActionPlan plano de ação consolidado
type ActionPlan struct {
	QuickWins     []EnhancedRecommendation `json:"quick_wins"`
	ShortTerm     []EnhancedRecommendation `json:"short_term"`
	MediumTerm    []EnhancedRecommendation `json:"medium_term"`
	LongTerm      []EnhancedRecommendation `json:"long_term"`
	Dependencies  map[string][]string      `json:"dependencies"`
	Timeline      string                   `json:"timeline"`
}

// RecommendationSummary resumo das recomendações
type RecommendationSummary struct {
	TotalRecommendations int                    `json:"total_recommendations"`
	ByPriority          map[string]int          `json:"by_priority"`
	ByCategory          map[string]int          `json:"by_category"`
	ExpectedValue       float64                `json:"expected_value"`
	TotalEffort         float64                `json:"total_effort"`
	AverageConfidence   float64                `json:"average_confidence"`
	TopOpportunities    []string               `json:"top_opportunities"`
	CriticalRisks       []string               `json:"critical_risks"`
}

// NewRecommendationEngine cria nova engine de recomendações
func NewRecommendationEngine(config RecommendationConfig) *RecommendationEngine {
	return &RecommendationEngine{
		config: config,
		predictor: NewTrendPredictor(TrendPredictorConfig{
			MinDataPoints:       6,
			ConfidenceThreshold: config.MinConfidence,
		}),
		analyzer: NewPatternAnalyzer(PatternAnalyzerConfig{
			MinDataPoints:   6,
			ConfidenceLevel: config.MinConfidence,
		}),
		detector: NewTrendDetector(TrendDetectorConfig{
			MinDataPoints: 6,
		}),
		forecaster: NewForecaster(ForecastConfig{
			DefaultHorizon: config.ContextWindow,
			MinAccuracy:    config.MinConfidence,
		}),
	}
}

// GenerateRecommendations gera recomendações completas
func (re *RecommendationEngine) GenerateRecommendations(data []models.HistoricalVRData, sindicato string) (*RecommendationSuite, error) {
	if len(data) < 3 {
		return nil, fmt.Errorf("dados insuficientes para gerar recomendações")
	}

	// Filtrar dados do sindicato
	sindicatoData := re.filterBySindicato(data, sindicato)
	if len(sindicatoData) < 3 {
		return nil, fmt.Errorf("dados insuficientes para sindicato %s", sindicato)
	}

	recommendations := make([]EnhancedRecommendation, 0)
	riskFactors := make([]RiskFactor, 0)
	opportunities := make([]Opportunity, 0)

	// 1. Análise de tendências
	trendRecs, trendRisks, trendOpps := re.analyzeTrendsForRecommendations(sindicatoData, sindicato)
	recommendations = append(recommendations, trendRecs...)
	riskFactors = append(riskFactors, trendRisks...)
	opportunities = append(opportunities, trendOpps...)

	// 2. Análise de padrões sazonais
	seasonalRecs, seasonalOpps := re.analyzeSeasonalityForRecommendations(sindicatoData, sindicato)
	recommendations = append(recommendations, seasonalRecs...)
	opportunities = append(opportunities, seasonalOpps...)

	// 3. Análise preditiva
	forecastRecs, forecastRisks := re.analyzeForecastForRecommendations(sindicatoData, sindicato)
	recommendations = append(recommendations, forecastRecs...)
	riskFactors = append(riskFactors, forecastRisks...)

	// 4. Análise de anomalias
	anomalyRecs, anomalyRisks := re.analyzeAnomaliesForRecommendations(sindicatoData, sindicato)
	recommendations = append(recommendations, anomalyRecs...)
	riskFactors = append(riskFactors, anomalyRisks...)

	// 5. Otimização de processos
	processRecs, processOpps := re.analyzeProcessOptimization(sindicatoData, sindicato)
	recommendations = append(recommendations, processRecs...)
	opportunities = append(opportunities, processOpps...)

	// Priorizar e filtrar recomendações
	recommendations = re.prioritizeRecommendations(recommendations)
	if len(recommendations) > re.config.MaxRecommendations {
		recommendations = recommendations[:re.config.MaxRecommendations]
	}

	// Criar plano de ação
	actionPlan := re.createActionPlan(recommendations)

	// Criar resumo
	summary := re.createSummary(recommendations, riskFactors, opportunities)

	suite := &RecommendationSuite{
		Summary:         summary,
		Recommendations: recommendations,
		ActionPlan:      actionPlan,
		RiskFactors:     riskFactors,
		Opportunities:   opportunities,
		GeneratedAt:     time.Now(),
		ValidUntil:      time.Now().AddDate(0, 1, 0), // válido por 1 mês
	}

	return suite, nil
}

// Métodos de análise específicos

func (re *RecommendationEngine) analyzeTrendsForRecommendations(data []models.HistoricalVRData, sindicato string) ([]EnhancedRecommendation, []RiskFactor, []Opportunity) {
	recommendations := make([]EnhancedRecommendation, 0)
	riskFactors := make([]RiskFactor, 0)
	opportunities := make([]Opportunity, 0)

	// Analisar tendências
	result, err := re.detector.DetectTrends(data, sindicato)
	if err != nil {
		return recommendations, riskFactors, opportunities
	}

	trend := result.PrimaryTrend
	if trend.Confidence < re.config.MinConfidence {
		return recommendations, riskFactors, opportunities
	}

	// Recomendações baseadas em tendências de crescimento
	if trend.Type == models.TrendUpward && trend.Strength > 0.3 {
		rec := EnhancedRecommendation{
			ActionItem: models.ActionItem{
				ID:          fmt.Sprintf("trend-growth-%s", sindicato),
				Priority:    models.PriorityHigh,
				Category:    "budget_planning",
				Title:       "Ajustar Orçamento para Crescimento",
				Description: fmt.Sprintf("Tendência de crescimento forte detectada (%.1f%% confiança)", trend.Confidence*100),
				Status:      models.ActionPending,
			},
			Context: RecommendationContext{
				SourceAnalysis:   "trend_analysis",
				TriggerCondition: "strong_upward_trend",
				AffectedEntities: []string{sindicato},
				TimeHorizon:      "3-6 meses",
				TrendInfo:        trend,
			},
			Evidence: []Evidence{
				{
					Type:        EvidenceTrend,
					Source:      "trend_detector",
					Description: fmt.Sprintf("Crescimento de %.2f%% por período", trend.Strength*100),
					Value:       trend.Strength,
					Confidence:  trend.Confidence,
				},
			},
			ExpectedImpact: ImpactAssessment{
				Financial: ImpactDetail{
					Value:       trend.Strength * 10000, // estimativa
					Unit:        "R$",
					Description: "Aumento estimado no orçamento mensal",
					Confidence:  trend.Confidence,
				},
				Timeline:    "2-3 meses",
				Probability: trend.Confidence,
			},
			Implementation: ImplementationGuide{
				Steps: []ImplementationStep{
					{
						ID:          "budget-review",
						Title:       "Revisar orçamento atual",
						Description: "Analisar orçamento atual e projetar necessidades",
						Duration:    "1 semana",
					},
					{
						ID:          "budget-adjust",
						Title:       "Ajustar provisão",
						Description: "Aumentar provisão baseado na tendência",
						Duration:    "1 semana",
					},
				},
				Timeline: "2-3 semanas",
				Resources: []RequiredResource{
					{
						Type:        "human",
						Description: "Analista financeiro",
						Quantity:    1,
						Unit:        "pessoa",
					},
				},
			},
			Confidence:   trend.Confidence,
			AutoApproved: trend.Confidence > re.config.AutoApproveThreshold,
		}
		recommendations = append(recommendations, rec)

		// Oportunidade de otimização
		opportunity := Opportunity{
			ID:          fmt.Sprintf("growth-opp-%s", sindicato),
			Title:       "Aproveitar Crescimento",
			Description: "Usar crescimento para melhorar eficiência operacional",
			Value:       trend.Strength * 5000, // estimativa
			Effort:      0.6,
			Timeline:    "3-4 meses",
		}
		opportunities = append(opportunities, opportunity)
	}

	// Recomendações para tendências de declínio
	if trend.Type == models.TrendDownward && trend.Strength > 0.2 {
		rec := EnhancedRecommendation{
			ActionItem: models.ActionItem{
				ID:          fmt.Sprintf("trend-decline-%s", sindicato),
				Priority:    models.PriorityHigh,
				Category:    "investigation",
				Title:       "Investigar Causa do Declínio",
				Description: fmt.Sprintf("Tendência de declínio detectada (%.1f%% confiança)", trend.Confidence*100),
				Status:      models.ActionPending,
			},
			Context: RecommendationContext{
				SourceAnalysis:   "trend_analysis",
				TriggerCondition: "downward_trend",
				AffectedEntities: []string{sindicato},
				TimeHorizon:      "1-2 meses",
				TrendInfo:        trend,
			},
			Confidence: trend.Confidence,
		}
		recommendations = append(recommendations, rec)

		// Fator de risco
		risk := RiskFactor{
			ID:          fmt.Sprintf("decline-risk-%s", sindicato),
			Type:        "operational",
			Description: "Declínio continuado pode afetar orçamento",
			Probability: trend.Confidence,
			Impact:      models.ImpactMedium,
			Source:      "trend_analysis",
		}
		riskFactors = append(riskFactors, risk)
	}

	return recommendations, riskFactors, opportunities
}

func (re *RecommendationEngine) analyzeSeasonalityForRecommendations(data []models.HistoricalVRData, sindicato string) ([]EnhancedRecommendation, []Opportunity) {
	recommendations := make([]EnhancedRecommendation, 0)
	opportunities := make([]Opportunity, 0)

	// Criar série temporal
	ts := re.createTimeSeries(data)
	seasonality, err := re.analyzer.DetectSeasonality(*ts)
	if err != nil || !seasonality.IsDetected {
		return recommendations, opportunities
	}

	if seasonality.Confidence > re.config.MinConfidence {
		rec := EnhancedRecommendation{
			ActionItem: models.ActionItem{
				ID:          fmt.Sprintf("seasonal-%s", sindicato),
				Priority:    models.PriorityMedium,
				Category:    "planning",
				Title:       "Implementar Planejamento Sazonal",
				Description: fmt.Sprintf("Padrão sazonal detectado (picos: %v, vales: %v)", seasonality.PeakMonths, seasonality.TroughMonths),
				Status:      models.ActionPending,
			},
			Context: RecommendationContext{
				SourceAnalysis:   "seasonality_analysis",
				TriggerCondition: "seasonal_pattern",
				AffectedEntities: []string{sindicato},
				TimeHorizon:      "12 meses",
				Seasonality:      seasonality,
			},
			Confidence: seasonality.Confidence,
		}
		recommendations = append(recommendations, rec)

		// Oportunidade de otimização sazonal
		opportunity := Opportunity{
			ID:          fmt.Sprintf("seasonal-opt-%s", sindicato),
			Title:       "Otimização Sazonal",
			Description: "Otimizar recursos baseado em padrão sazonal",
			Value:       seasonality.Amplitude * 0.1, // 10% da amplitude
			Effort:      0.4,
			Timeline:    "6 meses",
		}
		opportunities = append(opportunities, opportunity)
	}

	return recommendations, opportunities
}

func (re *RecommendationEngine) analyzeForecastForRecommendations(data []models.HistoricalVRData, sindicato string) ([]EnhancedRecommendation, []RiskFactor) {
	recommendations := make([]EnhancedRecommendation, 0)
	riskFactors := make([]RiskFactor, 0)

	// Gerar previsão
	forecast, err := re.forecaster.ForecastConsumption(data, sindicato, 3)
	if err != nil {
		return recommendations, riskFactors
	}

	baseForecast := forecast.WeightedForecast
	if baseForecast.Confidence < re.config.MinConfidence {
		return recommendations, riskFactors
	}

	// Comparar com valor atual
	currentValue := data[len(data)-1].TotalVR
	change := (baseForecast.PredictedVR - currentValue) / currentValue

	if change > 0.1 { // aumento > 10%
		rec := EnhancedRecommendation{
			ActionItem: models.ActionItem{
				ID:          fmt.Sprintf("forecast-increase-%s", sindicato),
				Priority:    models.PriorityMedium,
				Category:    "forecast_planning",
				Title:       "Preparar para Aumento Previsto",
				Description: fmt.Sprintf("Previsão indica aumento de %.1f%% no próximo período", change*100),
				Status:      models.ActionPending,
			},
			Context: RecommendationContext{
				SourceAnalysis:   "forecast_analysis",
				TriggerCondition: "forecast_increase",
				AffectedEntities: []string{sindicato},
				TimeHorizon:      "1-3 meses",
			},
			Confidence: baseForecast.Confidence,
		}
		recommendations = append(recommendations, rec)
	}

	if change < -0.1 { // diminuição > 10%
		risk := RiskFactor{
			ID:          fmt.Sprintf("forecast-decline-risk-%s", sindicato),
			Type:        "forecast",
			Description: fmt.Sprintf("Previsão indica declínio de %.1f%%", math.Abs(change)*100),
			Probability: baseForecast.Confidence,
			Impact:      models.ImpactMedium,
			Source:      "forecast_analysis",
		}
		riskFactors = append(riskFactors, risk)
	}

	return recommendations, riskFactors
}

func (re *RecommendationEngine) analyzeAnomaliesForRecommendations(data []models.HistoricalVRData, sindicato string) ([]EnhancedRecommendation, []RiskFactor) {
	recommendations := make([]EnhancedRecommendation, 0)
	riskFactors := make([]RiskFactor, 0)

	// Contar anomalias recentes (últimos 3 meses)
	recentAnomalies := 0
	for i := len(data) - 3; i < len(data) && i >= 0; i++ {
		recentAnomalies += len(data[i].Anomalies)
	}

	if recentAnomalies > 2 { // mais de 2 anomalias recentes
		rec := EnhancedRecommendation{
			ActionItem: models.ActionItem{
				ID:          fmt.Sprintf("anomaly-review-%s", sindicato),
				Priority:    models.PriorityHigh,
				Category:    "quality_control",
				Title:       "Revisar Controle de Qualidade",
				Description: fmt.Sprintf("%d anomalias detectadas nos últimos 3 meses", recentAnomalies),
				Status:      models.ActionPending,
			},
			Context: RecommendationContext{
				SourceAnalysis:   "anomaly_analysis",
				TriggerCondition: "frequent_anomalies",
				AffectedEntities: []string{sindicato},
				TimeHorizon:      "1 mês",
			},
			Confidence: 0.8,
		}
		recommendations = append(recommendations, rec)

		risk := RiskFactor{
			ID:          fmt.Sprintf("anomaly-risk-%s", sindicato),
			Type:        "quality",
			Description: "Anomalias frequentes podem indicar problemas sistêmicos",
			Probability: 0.7,
			Impact:      models.ImpactMedium,
			Source:      "anomaly_analysis",
		}
		riskFactors = append(riskFactors, risk)
	}

	return recommendations, riskFactors
}

func (re *RecommendationEngine) analyzeProcessOptimization(data []models.HistoricalVRData, sindicato string) ([]EnhancedRecommendation, []Opportunity) {
	recommendations := make([]EnhancedRecommendation, 0)
	opportunities := make([]Opportunity, 0)

	// Analisar eficiência baseado no volume de dados
	avgCollaborators := 0.0
	for _, d := range data {
		avgCollaborators += float64(d.NumColaboradores)
	}
	avgCollaborators /= float64(len(data))

	// Se processamento de muitos colaboradores, recomendar otimização
	if avgCollaborators > 500 {
		rec := EnhancedRecommendation{
			ActionItem: models.ActionItem{
				ID:          fmt.Sprintf("process-opt-%s", sindicato),
				Priority:    models.PriorityMedium,
				Category:    "optimization",
				Title:       "Otimizar Processamento",
				Description: fmt.Sprintf("Volume alto de colaboradores (%.0f em média) pode ser otimizado", avgCollaborators),
				Status:      models.ActionPending,
			},
			Context: RecommendationContext{
				SourceAnalysis:   "process_analysis",
				TriggerCondition: "high_volume",
				AffectedEntities: []string{sindicato},
				TimeHorizon:      "2-4 meses",
			},
			Confidence: 0.7,
		}
		recommendations = append(recommendations, rec)

		opportunity := Opportunity{
			ID:          fmt.Sprintf("process-eff-%s", sindicato),
			Title:       "Melhoria de Eficiência",
			Description: "Otimizar processamento para grandes volumes",
			Value:       1000, // estimativa de economia
			Effort:      0.8,
			Timeline:    "3-4 meses",
		}
		opportunities = append(opportunities, opportunity)
	}

	return recommendations, opportunities
}

// Métodos auxiliares

func (re *RecommendationEngine) filterBySindicato(data []models.HistoricalVRData, sindicato string) []models.HistoricalVRData {
	filtered := make([]models.HistoricalVRData, 0)
	for _, d := range data {
		if d.Sindicato == sindicato {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

func (re *RecommendationEngine) createTimeSeries(data []models.HistoricalVRData) *models.TimeSeries {
	ts := models.NewTimeSeries("VR Total", "R$")
	for _, d := range data {
		ts.AddPoint(d.Month, d.TotalVR, nil)
	}
	return ts
}

func (re *RecommendationEngine) prioritizeRecommendations(recommendations []EnhancedRecommendation) []EnhancedRecommendation {
	// Ordenar por prioridade e confiança
	sort.Slice(recommendations, func(i, j int) bool {
		// Primeiro por prioridade
		priorityOrder := map[models.Priority]int{
			models.PriorityCritical: 4,
			models.PriorityHigh:     3,
			models.PriorityMedium:   2,
			models.PriorityLow:      1,
		}
		
		if priorityOrder[recommendations[i].Priority] != priorityOrder[recommendations[j].Priority] {
			return priorityOrder[recommendations[i].Priority] > priorityOrder[recommendations[j].Priority]
		}
		
		// Depois por confiança
		return recommendations[i].Confidence > recommendations[j].Confidence
	})

	return recommendations
}

func (re *RecommendationEngine) createActionPlan(recommendations []EnhancedRecommendation) ActionPlan {
	plan := ActionPlan{
		QuickWins:    make([]EnhancedRecommendation, 0),
		ShortTerm:    make([]EnhancedRecommendation, 0),
		MediumTerm:   make([]EnhancedRecommendation, 0),
		LongTerm:     make([]EnhancedRecommendation, 0),
		Dependencies: make(map[string][]string),
	}

	for _, rec := range recommendations {
		// Classificar por timeline
		switch rec.Context.TimeHorizon {
		case "1 semana", "2 semanas":
			plan.QuickWins = append(plan.QuickWins, rec)
		case "1 mês", "1-2 meses", "2 meses":
			plan.ShortTerm = append(plan.ShortTerm, rec)
		case "3-6 meses", "6 meses":
			plan.MediumTerm = append(plan.MediumTerm, rec)
		default:
			plan.LongTerm = append(plan.LongTerm, rec)
		}
	}

	return plan
}

func (re *RecommendationEngine) createSummary(recommendations []EnhancedRecommendation, risks []RiskFactor, opportunities []Opportunity) RecommendationSummary {
	byPriority := make(map[string]int)
	byCategory := make(map[string]int)
	totalConfidence := 0.0

	for _, rec := range recommendations {
		byPriority[string(rec.Priority)]++
		byCategory[rec.Category]++
		totalConfidence += rec.Confidence
	}

	avgConfidence := 0.0
	if len(recommendations) > 0 {
		avgConfidence = totalConfidence / float64(len(recommendations))
	}

	// Top oportunidades
	topOpps := make([]string, 0)
	for i, opp := range opportunities {
		if i < 3 { // top 3
			topOpps = append(topOpps, opp.Title)
		}
	}

	// Riscos críticos
	criticalRisks := make([]string, 0)
	for _, risk := range risks {
		if risk.Impact == models.ImpactHigh || risk.Impact == models.ImpactCritical {
			criticalRisks = append(criticalRisks, risk.Description)
		}
	}

	return RecommendationSummary{
		TotalRecommendations: len(recommendations),
		ByPriority:          byPriority,
		ByCategory:          byCategory,
		ExpectedValue:       0, // seria calculado baseado nos impactos
		TotalEffort:         0, // seria calculado baseado nos esforços
		AverageConfidence:   avgConfidence,
		TopOpportunities:    topOpps,
		CriticalRisks:       criticalRisks,
	}
}