package intelligence

import (
	"fmt"
	"math"
	"sort"
	"time"

	"BrxAgente-desafio4/internal/models"
	"BrxAgente-desafio4/internal/modelo"
)

// TrendPredictor implementa modelos preditivos para tendências
type TrendPredictor struct {
	historicalData []models.HistoricalVRData
	models         map[string]models.PredictionModel
	config         TrendPredictorConfig
}

// TrendPredictorConfig configuração do preditor de tendências
type TrendPredictorConfig struct {
	MinDataPoints       int     `json:"min_data_points"`       // mínimo de pontos para predição
	ConfidenceThreshold float64 `json:"confidence_threshold"`  // limite mínimo de confiança
	ForecastPeriods     int     `json:"forecast_periods"`      // períodos futuros a prever
	MovingAverageWindow int     `json:"moving_average_window"` // janela para média móvel
}

// ConsumptionPredictor especializado em prever consumo de VR
type ConsumptionPredictor struct {
	*TrendPredictor
	seasonalPatterns map[string]*models.SeasonalityInfo
}

// AnomalyPredictor especializado em prever anomalias
type AnomalyPredictor struct {
	*TrendPredictor
	riskThresholds map[string]float64
}

// ProcessPredictor especializado em otimização de processos
type ProcessPredictor struct {
	*TrendPredictor
	processMetrics map[string]models.ProcessState
}

// NewTrendPredictor cria um novo preditor de tendências
func NewTrendPredictor(config TrendPredictorConfig) *TrendPredictor {
	return &TrendPredictor{
		historicalData: make([]models.HistoricalVRData, 0),
		models:         make(map[string]models.PredictionModel),
		config:         config,
	}
}

// NewConsumptionPredictor cria um preditor de consumo
func NewConsumptionPredictor(config TrendPredictorConfig) *ConsumptionPredictor {
	return &ConsumptionPredictor{
		TrendPredictor:   NewTrendPredictor(config),
		seasonalPatterns: make(map[string]*models.SeasonalityInfo),
	}
}

// NewAnomalyPredictor cria um preditor de anomalias
func NewAnomalyPredictor(config TrendPredictorConfig) *AnomalyPredictor {
	return &AnomalyPredictor{
		TrendPredictor: NewTrendPredictor(config),
		riskThresholds: map[string]float64{
			"variacao_alta":    2.0,  // 2 desvios padrão
			"ausencia_padrao":  0.3,  // 30% acima do normal
			"valor_discrepante": 1.5, // 1.5x a média
		},
	}
}

// NewProcessPredictor cria um preditor de processos
func NewProcessPredictor(config TrendPredictorConfig) *ProcessPredictor {
	return &ProcessPredictor{
		TrendPredictor: NewTrendPredictor(config),
		processMetrics: make(map[string]models.ProcessState),
	}
}

// LoadHistoricalData carrega dados históricos para predição
func (tp *TrendPredictor) LoadHistoricalData(data []models.HistoricalVRData) error {
	if len(data) < tp.config.MinDataPoints {
		return fmt.Errorf("dados insuficientes: %d pontos, mínimo requerido: %d", 
			len(data), tp.config.MinDataPoints)
	}

	// Ordenar dados por data
	sort.Slice(data, func(i, j int) bool {
		return data[i].Month.Before(data[j].Month)
	})

	tp.historicalData = data
	return nil
}

// PredictNextMonth prevê o VR para o próximo mês (ConsumptionPredictor)
func (cp *ConsumptionPredictor) PredictNextMonth(sindicato string) (*models.ConsumptionForecast, error) {
	// Filtrar dados para o sindicato específico
	sindicatoData := cp.filterBySindicato(sindicato)
	if len(sindicatoData) < cp.config.MinDataPoints {
		return nil, fmt.Errorf("dados insuficientes para sindicato %s", sindicato)
	}

	// Extrair série temporal de VR
	timeSeries := cp.extractVRTimeSeries(sindicatoData)
	
	// Detectar tendência
	trend := cp.detectTrend(timeSeries)
	
	// Detectar sazonalidade
	seasonality := cp.detectSeasonality(timeSeries)
	cp.seasonalPatterns[sindicato] = seasonality

	// Calcular previsão usando modelo híbrido
	prediction := cp.calculateForecast(timeSeries, trend, seasonality)

	// Próximo mês
	nextMonth := time.Now().AddDate(0, 1, 0)
	nextMonth = time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())

	forecast := &models.ConsumptionForecast{
		Sindicato:   sindicato,
		Month:       nextMonth,
		PredictedVR: prediction.Value,
		Confidence:  prediction.Confidence,
		Range: models.ForecastRange{
			Lower: prediction.Value * (1 - prediction.Uncertainty),
			Upper: prediction.Value * (1 + prediction.Uncertainty),
		},
		Factors:     prediction.Factors,
		Seasonality: seasonality,
		Trend:       trend,
		Assumptions: []string{
			"Padrões históricos se mantêm",
			"Não há mudanças estruturais significativas",
			"Fatores sazonais se repetem",
		},
	}

	return forecast, nil
}

// AssessRisk avalia risco de anomalias para colaborador (AnomalyPredictor)
func (ap *AnomalyPredictor) AssessRisk(colaborador *modelo.Colaborador) (*models.RiskAssessment, error) {
	if colaborador == nil {
		return nil, fmt.Errorf("colaborador não fornecido")
	}

	// Calcular fatores de risco
	riskFactors := ap.calculateRiskFactors(colaborador)
	
	// Calcular score de risco agregado
	riskScore := ap.aggregateRiskScore(riskFactors)
	
	// Determinar nível de risco
	riskLevel := models.GetRiskLevel(riskScore)
	
	// Calcular probabilidade baseada em histórico
	probability := ap.calculateProbability(colaborador, riskFactors)
	
	// Determinar impacto potencial
	impact := ap.assessImpact(colaborador, riskScore)
	
	// Gerar ações de mitigação
	mitigation := ap.generateMitigationActions(riskLevel, riskFactors)

	assessment := &models.RiskAssessment{
		Matricula:     colaborador.Matricula,
		Sindicato:     colaborador.Sindicato,
		RiskScore:     riskScore,
		RiskLevel:     riskLevel,
		RiskFactors:   riskFactors,
		Probability:   probability,
		Impact:        impact,
		Mitigation:    mitigation,
		LastEvaluated: time.Now(),
	}

	return assessment, nil
}

// OptimizeSchedule otimiza cronograma de processamento (ProcessPredictor)
func (pp *ProcessPredictor) OptimizeSchedule(month time.Time) (*models.ProcessOptimization, error) {
	// Prever volume de dados para o mês
	dataVolume := pp.predictDataVolume(month)
	
	// Estimar tempo de processamento atual
	currentState := pp.estimateCurrentProcessing(dataVolume)
	
	// Calcular estado otimizado
	optimalState := pp.calculateOptimalState(currentState, dataVolume)
	
	// Gerar sugestões de melhoria
	improvements := pp.generateImprovements(currentState, optimalState)
	
	// Calcular ganhos esperados
	gains := pp.calculateGains(currentState, optimalState)
	
	// Criar plano de implementação
	implementation := pp.createImplementationPlan(improvements)

	optimization := &models.ProcessOptimization{
		ProcessID:      fmt.Sprintf("vr-process-%s", month.Format("2006-01")),
		Month:          month,
		CurrentState:   currentState,
		OptimalState:   optimalState,
		Improvements:   improvements,
		ExpectedGains:  gains,
		Implementation: implementation,
	}

	return optimization, nil
}

// Métodos auxiliares internos

type forecastResult struct {
	Value       float64
	Confidence  float64
	Uncertainty float64
	Factors     []string
}

func (cp *ConsumptionPredictor) filterBySindicato(sindicato string) []models.HistoricalVRData {
	var result []models.HistoricalVRData
	for _, data := range cp.historicalData {
		if data.Sindicato == sindicato {
			result = append(result, data)
		}
	}
	return result
}

func (cp *ConsumptionPredictor) extractVRTimeSeries(data []models.HistoricalVRData) *models.TimeSeries {
	ts := models.NewTimeSeries("VR Total", "R$")
	for _, d := range data {
		ts.AddPoint(d.Month, d.TotalVR, map[string]interface{}{
			"colaboradores": d.NumColaboradores,
			"media":         d.MediaPorPessoa,
		})
	}
	return ts
}

func (cp *ConsumptionPredictor) detectTrend(ts *models.TimeSeries) *models.VRTrend {
	values := ts.GetValues()
	if len(values) < 3 {
		return &models.VRTrend{
			Type:        models.TrendStable,
			Strength:    0.0,
			Confidence:  0.5,
			Description: "Dados insuficientes para detectar tendência",
		}
	}

	// Calcular tendência usando regressão linear simples
	n := len(values)
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0
	
	for i, y := range values {
		x := float64(i + 1)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	
	// Coeficiente angular (slope)
	slope := (float64(n)*sumXY - sumX*sumY) / (float64(n)*sumX2 - sumX*sumX)
	
	// Determinar tipo e força da tendência
	var trendType models.TrendType
	strength := math.Abs(slope) / (sumY / float64(n)) // normalizar pelo valor médio
	
	if math.Abs(slope) < 0.01 {
		trendType = models.TrendStable
	} else if slope > 0 {
		trendType = models.TrendUpward
	} else {
		trendType = models.TrendDownward
	}

	// Calcular coeficiente de correlação para confiança
	meanX, meanY := sumX/float64(n), sumY/float64(n)
	ssXY, ssX, ssY := 0.0, 0.0, 0.0
	
	for i, y := range values {
		x := float64(i + 1)
		ssXY += (x - meanX) * (y - meanY)
		ssX += (x - meanX) * (x - meanX)
		ssY += (y - meanY) * (y - meanY)
	}
	
	correlation := ssXY / math.Sqrt(ssX*ssY)
	confidence := math.Abs(correlation)

	return &models.VRTrend{
		Type:        trendType,
		Strength:    math.Min(strength, 1.0),
		Period:      n,
		Confidence:  confidence,
		StartDate:   ts.Points[0].Timestamp,
		Description: cp.describeTrend(trendType, strength, confidence),
	}
}

func (cp *ConsumptionPredictor) detectSeasonality(ts *models.TimeSeries) *models.SeasonalityInfo {
	if len(ts.Points) < 12 {
		return &models.SeasonalityInfo{
			IsDetected: false,
			Confidence: 0.0,
		}
	}

	// Agrupar por mês do ano
	monthlyValues := make(map[int][]float64)
	for _, point := range ts.Points {
		month := int(point.Timestamp.Month())
		monthlyValues[month] = append(monthlyValues[month], point.Value)
	}

	// Calcular médias mensais
	monthlyAvgs := make(map[int]float64)
	overallSum, overallCount := 0.0, 0
	
	for month := 1; month <= 12; month++ {
		if values, exists := monthlyValues[month]; exists {
			sum := 0.0
			for _, v := range values {
				sum += v
				overallCount++
			}
			monthlyAvgs[month] = sum / float64(len(values))
			overallSum += monthlyAvgs[month]
		}
	}

	if len(monthlyAvgs) < 6 {
		return &models.SeasonalityInfo{
			IsDetected: false,
			Confidence: 0.0,
		}
	}

	overallAvg := overallSum / float64(len(monthlyAvgs))

	// Calcular variabilidade sazonal
	variance := 0.0
	for _, avg := range monthlyAvgs {
		variance += (avg - overallAvg) * (avg - overallAvg)
	}
	variance /= float64(len(monthlyAvgs))
	
	// Se variabilidade for significativa, detectar sazonalidade
	cv := math.Sqrt(variance) / overallAvg // coeficiente de variação
	isDetected := cv > 0.1 // 10% de variação

	// Identificar picos e vales
	var peakMonths, troughMonths []int
	threshold := overallAvg * 0.1 // 10% da média
	
	for month, avg := range monthlyAvgs {
		if avg > overallAvg+threshold {
			peakMonths = append(peakMonths, month)
		} else if avg < overallAvg-threshold {
			troughMonths = append(troughMonths, month)
		}
	}

	return &models.SeasonalityInfo{
		IsDetected:   isDetected,
		Period:       12,
		Amplitude:    math.Sqrt(variance),
		PeakMonths:   peakMonths,
		TroughMonths: troughMonths,
		Confidence:   math.Min(cv*2, 1.0), // converter CV em confiança
		Pattern:      monthlyAvgs,
	}
}

func (cp *ConsumptionPredictor) calculateForecast(ts *models.TimeSeries, trend *models.VRTrend, seasonality *models.SeasonalityInfo) *forecastResult {
	values := ts.GetValues()
	if len(values) == 0 {
		return &forecastResult{
			Value:      0,
			Confidence: 0,
			Factors:    []string{"Sem dados históricos"},
		}
	}

	// Média móvel simples como base
	windowSize := cp.config.MovingAverageWindow
	if windowSize > len(values) {
		windowSize = len(values)
	}

	baseValue := 0.0
	for i := len(values) - windowSize; i < len(values); i++ {
		baseValue += values[i]
	}
	baseValue /= float64(windowSize)

	// Ajustar por tendência
	if trend.Type == models.TrendUpward {
		baseValue *= (1 + trend.Strength*0.1) // crescimento moderado
	} else if trend.Type == models.TrendDownward {
		baseValue *= (1 - trend.Strength*0.1) // decrescimento moderado
	}

	// Ajustar por sazonalidade
	nextMonth := int(time.Now().AddDate(0, 1, 0).Month())
	if seasonality.IsDetected && len(seasonality.Pattern) > 0 {
		if seasonal, exists := seasonality.Pattern[nextMonth]; exists {
			overallAvg := 0.0
			for _, v := range seasonality.Pattern {
				overallAvg += v
			}
			overallAvg /= float64(len(seasonality.Pattern))
			
			seasonalFactor := seasonal / overallAvg
			baseValue *= seasonalFactor
		}
	}

	// Calcular confiança
	confidence := (trend.Confidence + seasonality.Confidence) / 2
	if confidence < cp.config.ConfidenceThreshold {
		confidence = cp.config.ConfidenceThreshold
	}

	// Calcular incerteza
	uncertainty := 0.2 // 20% base de incerteza
	if trend.Type == models.TrendVolatile {
		uncertainty += 0.1
	}
	if !seasonality.IsDetected {
		uncertainty += 0.1
	}

	factors := []string{
		fmt.Sprintf("Tendência: %s", trend.Type),
		fmt.Sprintf("Confiança da tendência: %.1f%%", trend.Confidence*100),
	}
	
	if seasonality.IsDetected {
		factors = append(factors, "Padrão sazonal detectado")
	}

	return &forecastResult{
		Value:       baseValue,
		Confidence:  confidence,
		Uncertainty: uncertainty,
		Factors:     factors,
	}
}

func (cp *ConsumptionPredictor) describeTrend(trendType models.TrendType, strength, confidence float64) string {
	strengthDesc := "fraca"
	if strength > 0.3 {
		strengthDesc = "moderada"
	}
	if strength > 0.6 {
		strengthDesc = "forte"
	}

	confidenceDesc := "baixa"
	if confidence > 0.5 {
		confidenceDesc = "média"
	}
	if confidence > 0.7 {
		confidenceDesc = "alta"
	}

	switch trendType {
	case models.TrendUpward:
		return fmt.Sprintf("Tendência de crescimento %s com confiança %s", strengthDesc, confidenceDesc)
	case models.TrendDownward:
		return fmt.Sprintf("Tendência de decréscimo %s com confiança %s", strengthDesc, confidenceDesc)
	case models.TrendStable:
		return fmt.Sprintf("Tendência estável com confiança %s", confidenceDesc)
	default:
		return "Tendência não identificada"
	}
}

// Métodos do AnomalyPredictor

func (ap *AnomalyPredictor) calculateRiskFactors(colaborador *modelo.Colaborador) []models.RiskFactor {
	factors := make([]models.RiskFactor, 0)

	// Fator: Padrão de ausências
	ausenciaRisk := ap.calculateAbsenceRisk(colaborador)
	factors = append(factors, models.RiskFactor{
		Factor:      "padrao_ausencias",
		Weight:      0.3,
		Value:       ausenciaRisk,
		Threshold:   ap.riskThresholds["ausencia_padrao"],
		Description: "Padrão de ausências acima do normal",
	})

	// Fator: Variação nos valores de VR
	variacaoRisk := ap.calculateVariationRisk(colaborador)
	factors = append(factors, models.RiskFactor{
		Factor:      "variacao_vr",
		Weight:      0.4,
		Value:       variacaoRisk,
		Threshold:   ap.riskThresholds["variacao_alta"],
		Description: "Variação alta nos valores de VR",
	})

	// Fator: Dados incompletos
	completudeRisk := ap.calculateCompletenessRisk(colaborador)
	factors = append(factors, models.RiskFactor{
		Factor:      "completude_dados",
		Weight:      0.3,
		Value:       completudeRisk,
		Threshold:   0.2,
		Description: "Dados incompletos ou inconsistentes",
	})

	return factors
}

func (ap *AnomalyPredictor) calculateAbsenceRisk(colaborador *modelo.Colaborador) float64 {
	// Simular cálculo baseado em dados de afastamentos e férias
	totalDays := 0
	absenceDays := 0

	// Contar dias de afastamento
	for _, afastamento := range colaborador.Afastamentos {
		days := int(afastamento.Fim.Sub(afastamento.Inicio).Hours() / 24)
		absenceDays += days
		totalDays += 30 // mês base
	}

	// Contar dias de férias
	for _, ferias := range colaborador.Ferias {
		days := int(ferias.Fim.Sub(ferias.Inicio).Hours() / 24)
		absenceDays += days
		totalDays += 30 // mês base
	}

	if totalDays == 0 {
		return 0.0
	}

	absenceRate := float64(absenceDays) / float64(totalDays)
	return absenceRate
}

func (ap *AnomalyPredictor) calculateVariationRisk(colaborador *modelo.Colaborador) float64 {
	// Simular análise de variação nos valores de VR
	// Em implementação real, analisaria histórico de valores
	return 0.5 // placeholder
}

func (ap *AnomalyPredictor) calculateCompletenessRisk(colaborador *modelo.Colaborador) float64 {
	completeness := 1.0

	// Verificar campos obrigatórios
	if colaborador.Nome == "" {
		completeness -= 0.2
	}
	if colaborador.Matricula == "" {
		completeness -= 0.3
	}
	if colaborador.Sindicato == "" {
		completeness -= 0.2
	}
	if colaborador.DataAdmissao.IsZero() {
		completeness -= 0.3
	}

	return math.Max(0, 1-completeness)
}

func (ap *AnomalyPredictor) aggregateRiskScore(factors []models.RiskFactor) float64 {
	totalScore := 0.0
	totalWeight := 0.0

	for _, factor := range factors {
		score := math.Min(factor.Value/factor.Threshold, 1.0) * 100
		totalScore += score * factor.Weight
		totalWeight += factor.Weight
	}

	if totalWeight == 0 {
		return 0
	}

	return totalScore / totalWeight
}

func (ap *AnomalyPredictor) calculateProbability(colaborador *modelo.Colaborador, factors []models.RiskFactor) float64 {
	// Calcular probabilidade baseada em fatores de risco
	baseProb := 0.1 // 10% base
	
	for _, factor := range factors {
		if factor.Value > factor.Threshold {
			baseProb += factor.Weight * 0.3
		}
	}

	return math.Min(baseProb, 0.95)
}

func (ap *AnomalyPredictor) assessImpact(colaborador *modelo.Colaborador, riskScore float64) models.ImpactLevel {
	if riskScore < 30 {
		return models.ImpactLow
	} else if riskScore < 60 {
		return models.ImpactMedium
	} else if riskScore < 80 {
		return models.ImpactHigh
	}
	return models.ImpactCritical
}

func (ap *AnomalyPredictor) generateMitigationActions(riskLevel models.RiskLevel, factors []models.RiskFactor) []models.ActionItem {
	actions := make([]models.ActionItem, 0)
	priority := models.GetPriorityFromRisk(riskLevel)

	for _, factor := range factors {
		if factor.Value > factor.Threshold {
			action := models.ActionItem{
				ID:          fmt.Sprintf("action-%s-%d", factor.Factor, time.Now().Unix()),
				Priority:    priority,
				Category:    "risk_mitigation",
				Title:       fmt.Sprintf("Mitigar risco: %s", factor.Factor),
				Description: factor.Description,
				Status:      models.ActionPending,
			}
			actions = append(actions, action)
		}
	}

	return actions
}

// Métodos do ProcessPredictor

func (pp *ProcessPredictor) predictDataVolume(month time.Time) int {
	// Simular predição de volume baseada em histórico
	baseVolume := 1000 // número base de colaboradores
	
	// Ajuste sazonal (dezembro tem mais dados)
	if month.Month() == 12 {
		baseVolume = int(float64(baseVolume) * 1.2)
	}

	return baseVolume
}

func (pp *ProcessPredictor) estimateCurrentProcessing(dataVolume int) models.ProcessState {
	return models.ProcessState{
		Duration:      time.Duration(dataVolume) * time.Millisecond * 50, // 50ms por colaborador
		ResourceUsage: map[string]float64{
			"cpu":    0.7,
			"memory": 0.5,
			"disk":   0.3,
		},
		Efficiency:    0.75,
		ErrorRate:     0.05,
		Throughput:    dataVolume,
		Metrics: map[string]interface{}{
			"parallel_workers": 4,
			"batch_size":      100,
		},
	}
}

func (pp *ProcessPredictor) calculateOptimalState(current models.ProcessState, dataVolume int) models.ProcessState {
	optimal := current
	
	// Otimizar duração (redução de 30%)
	optimal.Duration = time.Duration(float64(optimal.Duration) * 0.7)
	
	// Melhor uso de recursos
	optimal.ResourceUsage["cpu"] = 0.8
	optimal.ResourceUsage["memory"] = 0.6
	
	// Maior eficiência
	optimal.Efficiency = 0.9
	
	// Menor taxa de erro
	optimal.ErrorRate = 0.02

	// Otimizar configurações
	optimal.Metrics["parallel_workers"] = 8
	optimal.Metrics["batch_size"] = 200

	return optimal
}

func (pp *ProcessPredictor) generateImprovements(current, optimal models.ProcessState) []models.ImprovementSuggestion {
	improvements := []models.ImprovementSuggestion{
		{
			Area:        "paralelizacao",
			Type:        "performance",
			Description: "Aumentar número de workers paralelos",
			Impact:      0.7,
			Effort:      0.3,
			Priority:    models.PriorityHigh,
		},
		{
			Area:        "batch_processing",
			Type:        "optimization",
			Description: "Otimizar tamanho dos lotes de processamento",
			Impact:      0.5,
			Effort:      0.2,
			Priority:    models.PriorityMedium,
		},
		{
			Area:        "error_handling",
			Type:        "reliability",
			Description: "Melhorar tratamento de erros e retry logic",
			Impact:      0.6,
			Effort:      0.4,
			Priority:    models.PriorityMedium,
		},
	}

	return improvements
}

func (pp *ProcessPredictor) calculateGains(current, optimal models.ProcessState) models.OptimizationGains {
	timeDiff := current.Duration - optimal.Duration
	timeReduction := float64(timeDiff) / float64(current.Duration) * 100

	return models.OptimizationGains{
		TimeReduction:    timeReduction,
		ResourceSaving:   15.0, // 15% economia
		AccuracyIncrease: (optimal.Efficiency - current.Efficiency) * 100,
		CostReduction:    12.0, // 12% redução de custos
	}
}

func (pp *ProcessPredictor) createImplementationPlan(improvements []models.ImprovementSuggestion) models.ImplementationPlan {
	steps := make([]models.ImplementationStep, 0)

	for i, improvement := range improvements {
		step := models.ImplementationStep{
			ID:          fmt.Sprintf("step-%d", i+1),
			Title:       improvement.Description,
			Description: fmt.Sprintf("Implementar melhoria em %s", improvement.Area),
			Duration:    "2-3 dias",
			Owner:       "DevOps Team",
		}
		steps = append(steps, step)
	}

	return models.ImplementationPlan{
		Steps:    steps,
		Timeline: "2-3 semanas",
		Resources: []string{
			"Desenvolvedor Backend",
			"DevOps Engineer",
			"Ambiente de Teste",
		},
		Risks: []string{
			"Possível instabilidade temporária",
			"Necessidade de rollback",
		},
		Success: []string{
			"Redução de tempo de processamento > 20%",
			"Taxa de erro < 3%",
			"Utilização de recursos otimizada",
		},
	}
}