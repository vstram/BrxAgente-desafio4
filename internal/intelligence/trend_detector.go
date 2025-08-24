package intelligence

import (
	"fmt"
	"math"
	"sort"
	"time"

	"BrxAgente-desafio4/internal/predicoes"
)

// TrendDetector detecta tendências em dados de séries temporais
type TrendDetector struct {
	config TrendDetectorConfig
	cache  map[string]*TrendCache
}

// TrendDetectorConfig configuração do detector de tendências
type TrendDetectorConfig struct {
	MinDataPoints       int     `json:"min_data_points"`      // mínimo de pontos para análise
	SensitivityLevel    float64 `json:"sensitivity_level"`    // sensibilidade para detecção (0-1)
	TrendThreshold      float64 `json:"trend_threshold"`      // limite para considerar tendência significativa
	VolatilityThreshold float64 `json:"volatility_threshold"` // limite para considerar dados voláteis
	SeasonalWindow      int     `json:"seasonal_window"`      // janela para análise sazonal
}

// TrendCache armazena tendências calculadas
type TrendCache struct {
	Trend          *predicoes.VRTrend
	CalculatedAt   time.Time
	ValidityPeriod time.Duration
}

// TrendAnalysisResult resultado da análise de tendências
type TrendAnalysisResult struct {
	PrimaryTrend    *predicoes.VRTrend         `json:"primary_trend"`
	SecondaryTrends []predicoes.VRTrend        `json:"secondary_trends"`
	Seasonality     *predicoes.SeasonalityInfo `json:"seasonality"`
	Volatility      float64                    `json:"volatility"`
	Confidence      float64                    `json:"confidence"`
	Recommendations []TrendRecommendation      `json:"recommendations"`
	Metadata        TrendAnalysisMetadata      `json:"metadata"`
}

// TrendRecommendation recomendação baseada em tendência
type TrendRecommendation struct {
	Type        string                 `json:"type"`
	Priority    predicoes.Priority     `json:"priority"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Actions     []predicoes.ActionItem `json:"actions"`
	Impact      float64                `json:"impact"`    // impacto esperado (0-1)
	Timeframe   string                 `json:"timeframe"` // prazo para implementação
}

// TrendAnalysisMetadata metadados da análise
type TrendAnalysisMetadata struct {
	DataPoints     int                    `json:"data_points"`
	AnalysisPeriod string                 `json:"analysis_period"`
	Methods        []string               `json:"methods"`
	Parameters     map[string]interface{} `json:"parameters"`
	QualityScore   float64                `json:"quality_score"`
}

// NewTrendDetector cria novo detector de tendências
func NewTrendDetector(config TrendDetectorConfig) *TrendDetector {
	return &TrendDetector{
		config: config,
		cache:  make(map[string]*TrendCache),
	}
}

// DetectTrends detecta tendências em dados históricos
func (td *TrendDetector) DetectTrends(data []predicoes.HistoricalVRData, entity string) (*TrendAnalysisResult, error) {
	if len(data) < td.config.MinDataPoints {
		return nil, fmt.Errorf("dados insuficientes: %d pontos, mínimo: %d", len(data), td.config.MinDataPoints)
	}

	// Verificar cache
	if cached := td.getCachedTrend(entity); cached != nil {
		if time.Since(cached.CalculatedAt) < cached.ValidityPeriod {
			return td.buildResultFromCache(cached), nil
		}
	}

	// Ordenar dados por data
	sort.Slice(data, func(i, j int) bool {
		return data[i].Month.Before(data[j].Month)
	})

	// Criar série temporal
	timeSeries := td.createTimeSeries(data)

	// Detectar tendência principal
	primaryTrend := td.detectPrimaryTrend(timeSeries)

	// Detectar tendências secundárias (curto/médio prazo)
	secondaryTrends := td.detectSecondaryTrends(timeSeries)

	// Analisar sazonalidade
	seasonality := td.analyzeSeasonality(timeSeries)

	// Calcular volatilidade
	volatility := td.calculateVolatility(timeSeries)

	// Calcular confiança geral
	confidence := td.calculateOverallConfidence(primaryTrend, seasonality, volatility)

	// Gerar recomendações
	recommendations := td.generateTrendRecommendations(primaryTrend, secondaryTrends, seasonality)

	// Metadados
	metadata := TrendAnalysisMetadata{
		DataPoints:     len(data),
		AnalysisPeriod: fmt.Sprintf("%s to %s", data[0].Month.Format("2006-01"), data[len(data)-1].Month.Format("2006-01")),
		Methods:        []string{"linear_regression", "moving_average", "seasonal_decomposition"},
		Parameters: map[string]interface{}{
			"sensitivity":          td.config.SensitivityLevel,
			"trend_threshold":      td.config.TrendThreshold,
			"volatility_threshold": td.config.VolatilityThreshold,
		},
		QualityScore: td.assessDataQuality(data),
	}

	result := &TrendAnalysisResult{
		PrimaryTrend:    primaryTrend,
		SecondaryTrends: secondaryTrends,
		Seasonality:     seasonality,
		Volatility:      volatility,
		Confidence:      confidence,
		Recommendations: recommendations,
		Metadata:        metadata,
	}

	// Atualizar cache
	td.updateCache(entity, primaryTrend)

	return result, nil
}

// DetectTrendShifts detecta mudanças bruscas de tendência
func (td *TrendDetector) DetectTrendShifts(data []predicoes.HistoricalVRData) ([]TrendShift, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("dados insuficientes para detectar mudanças de tendência")
	}

	shifts := make([]TrendShift, 0)
	windowSize := 6 // janela de 6 meses

	for i := windowSize; i < len(data)-windowSize; i++ {
		// Analisar tendência antes e depois do ponto
		beforeData := data[i-windowSize : i]
		afterData := data[i : i+windowSize]

		trendBefore := td.calculateTrendSlope(beforeData)
		trendAfter := td.calculateTrendSlope(afterData)

		// Detectar mudança significativa
		if td.isSignificantShift(trendBefore, trendAfter) {
			shift := TrendShift{
				Date:           data[i].Month,
				TrendBefore:    trendBefore,
				TrendAfter:     trendAfter,
				Magnitude:      math.Abs(trendAfter - trendBefore),
				Confidence:     td.calculateShiftConfidence(beforeData, afterData),
				Type:           td.classifyShiftType(trendBefore, trendAfter),
				AffectedEntity: data[i].Sindicato,
			}
			shifts = append(shifts, shift)
		}
	}

	return shifts, nil
}

// PredictTrendContinuation prevê continuação da tendência atual
func (td *TrendDetector) PredictTrendContinuation(data []predicoes.HistoricalVRData, periodsAhead int) (*TrendPrediction, error) {
	if len(data) < td.config.MinDataPoints {
		return nil, fmt.Errorf("dados insuficientes para predição")
	}

	// Detectar tendência atual
	result, err := td.DetectTrends(data, "prediction")
	if err != nil {
		return nil, err
	}

	// Calcular predição baseada na tendência
	lastValue := data[len(data)-1].TotalVR
	trendSlope := td.extractTrendSlope(result.PrimaryTrend)

	predictions := make([]TrendPredictionPoint, periodsAhead)
	lastDate := data[len(data)-1].Month

	for i := 0; i < periodsAhead; i++ {
		// Data da predição
		predDate := lastDate.AddDate(0, i+1, 0)

		// Valor base com tendência
		baseValue := lastValue * (1 + trendSlope*float64(i+1))

		// Ajustar por sazonalidade se detectada
		if result.Seasonality != nil && result.Seasonality.IsDetected {
			seasonalFactor := td.getSeasonalFactor(result.Seasonality, predDate)
			baseValue *= seasonalFactor
		}

		// Calcular intervalo de confiança
		uncertainty := td.calculatePredictionUncertainty(result, i+1)

		predictions[i] = TrendPredictionPoint{
			Date:       predDate,
			Value:      baseValue,
			Confidence: result.Confidence * (1 - float64(i)*0.1), // decai com tempo
			LowerBound: baseValue * (1 - uncertainty),
			UpperBound: baseValue * (1 + uncertainty),
			Factors:    td.identifyPredictionFactors(result),
		}
	}

	prediction := &TrendPrediction{
		BaseTrend:   result.PrimaryTrend,
		Predictions: predictions,
		Seasonality: result.Seasonality,
		Methodology: "trend_continuation",
		Assumptions: td.generatePredictionAssumptions(result),
		Reliability: td.assessPredictionReliability(result),
	}

	return prediction, nil
}

// Métodos auxiliares internos

func (td *TrendDetector) createTimeSeries(data []predicoes.HistoricalVRData) *predicoes.TimeSeries {
	ts := predicoes.NewTimeSeries("VR Total", "R$")
	for _, d := range data {
		ts.AddPoint(d.Month, d.TotalVR, map[string]interface{}{
			"sindicato":         d.Sindicato,
			"num_colaboradores": d.NumColaboradores,
		})
	}
	return ts
}

func (td *TrendDetector) detectPrimaryTrend(ts *predicoes.TimeSeries) *predicoes.VRTrend {
	values := ts.GetValues()

	// Regressão linear para tendência de longo prazo
	slope, _, r2 := td.linearRegression(values)

	// Classificar tipo de tendência
	trendType := td.classifyTrendType(slope, values)

	// Calcular força da tendência
	strength := td.calculateTrendStrength(slope, values)

	// Ajustar confiança baseada no R²
	confidence := math.Min(r2*1.2, 1.0) // R² ajustado

	return &predicoes.VRTrend{
		Type:        trendType,
		Strength:    strength,
		Period:      len(values),
		Confidence:  confidence,
		StartDate:   ts.Points[0].Timestamp,
		Description: td.describeTrend(trendType, strength, confidence),
	}
}

func (td *TrendDetector) detectSecondaryTrends(ts *predicoes.TimeSeries) []predicoes.VRTrend {
	trends := make([]predicoes.VRTrend, 0)

	// Tendência de curto prazo (últimos 3 meses)
	if len(ts.Points) >= 3 {
		shortTermValues := ts.GetValues()[len(ts.Points)-3:]
		slope, _, r2 := td.linearRegression(shortTermValues)

		if r2 > 0.5 && math.Abs(slope) > td.config.TrendThreshold {
			trend := predicoes.VRTrend{
				Type:        td.classifyTrendType(slope, shortTermValues),
				Strength:    td.calculateTrendStrength(slope, shortTermValues),
				Period:      3,
				Confidence:  r2,
				StartDate:   ts.Points[len(ts.Points)-3].Timestamp,
				Description: "Tendência de curto prazo (3 meses)",
			}
			trends = append(trends, trend)
		}
	}

	// Tendência de médio prazo (últimos 6 meses)
	if len(ts.Points) >= 6 {
		mediumTermValues := ts.GetValues()[len(ts.Points)-6:]
		slope, _, r2 := td.linearRegression(mediumTermValues)

		if r2 > 0.5 && math.Abs(slope) > td.config.TrendThreshold {
			trend := predicoes.VRTrend{
				Type:        td.classifyTrendType(slope, mediumTermValues),
				Strength:    td.calculateTrendStrength(slope, mediumTermValues),
				Period:      6,
				Confidence:  r2,
				StartDate:   ts.Points[len(ts.Points)-6].Timestamp,
				Description: "Tendência de médio prazo (6 meses)",
			}
			trends = append(trends, trend)
		}
	}

	return trends
}

func (td *TrendDetector) analyzeSeasonality(ts *predicoes.TimeSeries) *predicoes.SeasonalityInfo {
	if len(ts.Points) < 12 {
		return &predicoes.SeasonalityInfo{
			IsDetected: false,
			Confidence: 0.0,
		}
	}

	// Extrair componente sazonal
	seasonal := td.extractSeasonalComponent(ts)

	// Testar significância
	significance := td.testSeasonalSignificance(seasonal)

	// Identificar padrões
	peaks, troughs := td.identifySeasonalPeaks(seasonal)

	return &predicoes.SeasonalityInfo{
		IsDetected:   significance > 0.6, // 60% de confiança mínima
		Period:       12,
		Amplitude:    td.calculateSeasonalAmplitude(seasonal),
		PeakMonths:   peaks,
		TroughMonths: troughs,
		Confidence:   significance,
		Pattern:      seasonal,
	}
}

func (td *TrendDetector) calculateVolatility(ts *predicoes.TimeSeries) float64 {
	values := ts.GetValues()
	if len(values) < 2 {
		return 0.0
	}

	// Calcular retornos percentuais
	returns := make([]float64, len(values)-1)
	for i := 1; i < len(values); i++ {
		if values[i-1] != 0 {
			returns[i-1] = (values[i] - values[i-1]) / values[i-1]
		}
	}

	// Calcular desvio padrão dos retornos
	mean := 0.0
	for _, r := range returns {
		mean += r
	}
	mean /= float64(len(returns))

	variance := 0.0
	for _, r := range returns {
		diff := r - mean
		variance += diff * diff
	}
	variance /= float64(len(returns))

	return math.Sqrt(variance)
}

func (td *TrendDetector) calculateOverallConfidence(trend *predicoes.VRTrend, seasonality *predicoes.SeasonalityInfo, volatility float64) float64 {
	confidence := trend.Confidence * 0.6 // peso da tendência

	if seasonality.IsDetected {
		confidence += seasonality.Confidence * 0.3 // peso da sazonalidade
	}

	// Penalizar por alta volatilidade
	volatilityPenalty := math.Min(volatility*2, 0.3) // máximo 30% de penalidade
	confidence -= volatilityPenalty

	// Ajustar baseado na qualidade dos dados
	if trend.Period < 6 {
		confidence *= 0.8 // penalizar poucos dados
	}

	return math.Max(0, math.Min(confidence, 1.0))
}

func (td *TrendDetector) generateTrendRecommendations(primary *predicoes.VRTrend, secondary []predicoes.VRTrend, seasonality *predicoes.SeasonalityInfo) []TrendRecommendation {
	recommendations := make([]TrendRecommendation, 0)

	// Recomendações baseadas na tendência principal
	if primary.Type == predicoes.TrendUpward && primary.Confidence > 0.7 {
		rec := TrendRecommendation{
			Type:        "budget_planning",
			Priority:    predicoes.PriorityHigh,
			Title:       "Ajustar Orçamento para Crescimento",
			Description: fmt.Sprintf("Tendência de crescimento detectada (%.1f%% confiança). Considerar aumento no orçamento de VR.", primary.Confidence*100),
			Impact:      primary.Strength,
			Timeframe:   "2-3 meses",
			Actions: []predicoes.ActionItem{
				{
					ID:          "budget-adjust-1",
					Priority:    predicoes.PriorityHigh,
					Category:    "planning",
					Title:       "Revisar orçamento mensal",
					Description: "Aumentar provisão para VR baseado na tendência de crescimento",
					Status:      predicoes.ActionPending,
				},
			},
		}
		recommendations = append(recommendations, rec)
	}

	if primary.Type == predicoes.TrendDownward && primary.Confidence > 0.7 {
		rec := TrendRecommendation{
			Type:        "investigation",
			Priority:    predicoes.PriorityMedium,
			Title:       "Investigar Causa da Redução",
			Description: fmt.Sprintf("Tendência de declínio detectada (%.1f%% confiança). Investigar causas possíveis.", primary.Confidence*100),
			Impact:      primary.Strength,
			Timeframe:   "1-2 semanas",
			Actions: []predicoes.ActionItem{
				{
					ID:          "investigate-1",
					Priority:    predicoes.PriorityMedium,
					Category:    "analysis",
					Title:       "Analisar causas do declínio",
					Description: "Identificar fatores que contribuem para a redução no VR",
					Status:      predicoes.ActionPending,
				},
			},
		}
		recommendations = append(recommendations, rec)
	}

	// Recomendações baseadas na sazonalidade
	if seasonality.IsDetected && seasonality.Confidence > 0.6 {
		rec := TrendRecommendation{
			Type:        "seasonal_planning",
			Priority:    predicoes.PriorityMedium,
			Title:       "Planejamento Sazonal",
			Description: fmt.Sprintf("Padrão sazonal detectado. Ajustar planejamento para picos em %v e vales em %v.", seasonality.PeakMonths, seasonality.TroughMonths),
			Impact:      seasonality.Confidence,
			Timeframe:   "próximos 6 meses",
			Actions: []predicoes.ActionItem{
				{
					ID:          "seasonal-plan-1",
					Priority:    predicoes.PriorityMedium,
					Category:    "planning",
					Title:       "Criar cronograma sazonal",
					Description: "Desenvolver plano de ação específico para períodos sazonais",
					Status:      predicoes.ActionPending,
				},
			},
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations
}

// Estruturas auxiliares

type TrendShift struct {
	Date           time.Time      `json:"date"`
	TrendBefore    float64        `json:"trend_before"`
	TrendAfter     float64        `json:"trend_after"`
	Magnitude      float64        `json:"magnitude"`
	Confidence     float64        `json:"confidence"`
	Type           TrendShiftType `json:"type"`
	AffectedEntity string         `json:"affected_entity"`
}

type TrendShiftType string

const (
	ShiftAcceleration TrendShiftType = "acceleration" // aceleração da tendência
	ShiftDeceleration TrendShiftType = "deceleration" // desaceleração da tendência
	ShiftReversal     TrendShiftType = "reversal"     // reversão da tendência
	ShiftBreakout     TrendShiftType = "breakout"     // quebra de padrão
)

type TrendPrediction struct {
	BaseTrend   *predicoes.VRTrend         `json:"base_trend"`
	Predictions []TrendPredictionPoint     `json:"predictions"`
	Seasonality *predicoes.SeasonalityInfo `json:"seasonality"`
	Methodology string                     `json:"methodology"`
	Assumptions []string                   `json:"assumptions"`
	Reliability float64                    `json:"reliability"`
}

type TrendPredictionPoint struct {
	Date       time.Time `json:"date"`
	Value      float64   `json:"value"`
	Confidence float64   `json:"confidence"`
	LowerBound float64   `json:"lower_bound"`
	UpperBound float64   `json:"upper_bound"`
	Factors    []string  `json:"factors"`
}

// Implementação dos métodos auxiliares

func (td *TrendDetector) linearRegression(values []float64) (slope, intercept, r2 float64) {
	n := float64(len(values))
	if n < 2 {
		return 0, 0, 0
	}

	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0

	for i, y := range values {
		x := float64(i + 1)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	slope = (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept = (sumY - slope*sumX) / n

	// Calcular R²
	meanY := sumY / n
	ssRes, ssTot := 0.0, 0.0

	for i, y := range values {
		x := float64(i + 1)
		predicted := intercept + slope*x
		ssRes += (y - predicted) * (y - predicted)
		ssTot += (y - meanY) * (y - meanY)
	}

	if ssTot == 0 {
		r2 = 0
	} else {
		r2 = 1.0 - (ssRes / ssTot)
	}

	return slope, intercept, math.Max(0, r2)
}

func (td *TrendDetector) classifyTrendType(slope float64, values []float64) predicoes.TrendType {
	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	normalizedSlope := slope / mean
	threshold := td.config.TrendThreshold

	if math.Abs(normalizedSlope) < threshold {
		return predicoes.TrendStable
	} else if normalizedSlope > threshold {
		return predicoes.TrendUpward
	} else {
		return predicoes.TrendDownward
	}
}

func (td *TrendDetector) calculateTrendStrength(slope float64, values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	normalizedSlope := math.Abs(slope) / mean
	return math.Min(normalizedSlope, 1.0)
}

func (td *TrendDetector) describeTrend(trendType predicoes.TrendType, strength, confidence float64) string {
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
	case predicoes.TrendUpward:
		return fmt.Sprintf("Tendência de crescimento %s (confiança %s)", strengthDesc, confidenceDesc)
	case predicoes.TrendDownward:
		return fmt.Sprintf("Tendência de decréscimo %s (confiança %s)", strengthDesc, confidenceDesc)
	case predicoes.TrendStable:
		return fmt.Sprintf("Tendência estável (confiança %s)", confidenceDesc)
	default:
		return "Tendência não identificada"
	}
}

func (td *TrendDetector) extractSeasonalComponent(ts *predicoes.TimeSeries) map[int]float64 {
	monthlyData := make(map[int][]float64)

	for _, point := range ts.Points {
		month := int(point.Timestamp.Month())
		monthlyData[month] = append(monthlyData[month], point.Value)
	}

	seasonal := make(map[int]float64)
	for month := 1; month <= 12; month++ {
		if values, exists := monthlyData[month]; exists && len(values) > 0 {
			sum := 0.0
			for _, v := range values {
				sum += v
			}
			seasonal[month] = sum / float64(len(values))
		}
	}

	return seasonal
}

func (td *TrendDetector) testSeasonalSignificance(seasonal map[int]float64) float64 {
	if len(seasonal) < 4 {
		return 0.0
	}

	// Calcular média e variância
	sum := 0.0
	for _, v := range seasonal {
		sum += v
	}
	mean := sum / float64(len(seasonal))

	variance := 0.0
	for _, v := range seasonal {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(seasonal))

	// Significância baseada no coeficiente de variação
	cv := math.Sqrt(variance) / mean
	return math.Min(cv*2, 1.0) // converter CV em medida de significância
}

func (td *TrendDetector) identifySeasonalPeaks(seasonal map[int]float64) ([]int, []int) {
	if len(seasonal) == 0 {
		return nil, nil
	}

	// Calcular média
	sum := 0.0
	for _, v := range seasonal {
		sum += v
	}
	mean := sum / float64(len(seasonal))

	threshold := mean * 0.1 // 10% da média

	var peaks, troughs []int
	for month, value := range seasonal {
		if value > mean+threshold {
			peaks = append(peaks, month)
		} else if value < mean-threshold {
			troughs = append(troughs, month)
		}
	}

	return peaks, troughs
}

func (td *TrendDetector) calculateSeasonalAmplitude(seasonal map[int]float64) float64 {
	if len(seasonal) == 0 {
		return 0
	}

	min, max := math.Inf(1), math.Inf(-1)
	for _, v := range seasonal {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	return max - min
}

// Métodos para cache

func (td *TrendDetector) getCachedTrend(entity string) *TrendCache {
	return td.cache[entity]
}

func (td *TrendDetector) updateCache(entity string, trend *predicoes.VRTrend) {
	td.cache[entity] = &TrendCache{
		Trend:          trend,
		CalculatedAt:   time.Now(),
		ValidityPeriod: time.Hour * 24, // 24 horas de validade
	}
}

func (td *TrendDetector) buildResultFromCache(cached *TrendCache) *TrendAnalysisResult {
	return &TrendAnalysisResult{
		PrimaryTrend: cached.Trend,
		Confidence:   cached.Trend.Confidence,
		Metadata: TrendAnalysisMetadata{
			DataPoints: cached.Trend.Period,
			Methods:    []string{"cached"},
		},
	}
}

// Métodos auxiliares para detecção de mudanças

func (td *TrendDetector) calculateTrendSlope(data []predicoes.HistoricalVRData) float64 {
	values := make([]float64, len(data))
	for i, d := range data {
		values[i] = d.TotalVR
	}

	slope, _, _ := td.linearRegression(values)
	return slope
}

func (td *TrendDetector) isSignificantShift(before, after float64) bool {
	diff := math.Abs(after - before)
	avgSlope := (math.Abs(before) + math.Abs(after)) / 2

	if avgSlope == 0 {
		return diff > td.config.TrendThreshold
	}

	relativeDiff := diff / avgSlope
	return relativeDiff > 1.0 // mudança > 100% da inclinação média
}

func (td *TrendDetector) calculateShiftConfidence(before, after []predicoes.HistoricalVRData) float64 {
	// Calcular R² para ambos os segmentos
	valuesBefore := make([]float64, len(before))
	valuesAfter := make([]float64, len(after))

	for i, d := range before {
		valuesBefore[i] = d.TotalVR
	}
	for i, d := range after {
		valuesAfter[i] = d.TotalVR
	}

	_, _, r2Before := td.linearRegression(valuesBefore)
	_, _, r2After := td.linearRegression(valuesAfter)

	// Confiança baseada na qualidade do ajuste em ambos os segmentos
	return (r2Before + r2After) / 2
}

func (td *TrendDetector) classifyShiftType(before, after float64) TrendShiftType {
	// Determinar direção das tendências
	beforeDirection := 0.0
	if math.Abs(before) > td.config.TrendThreshold {
		if before > 0 {
			beforeDirection = 1.0
		} else {
			beforeDirection = -1.0
		}
	}

	afterDirection := 0.0
	if math.Abs(after) > td.config.TrendThreshold {
		if after > 0 {
			afterDirection = 1.0
		} else {
			afterDirection = -1.0
		}
	}

	// Classificar tipo de mudança
	if beforeDirection*afterDirection < 0 {
		return ShiftReversal // inversão de direção
	}

	if math.Abs(after) > math.Abs(before)*1.5 {
		return ShiftAcceleration // aceleração
	}

	if math.Abs(after) < math.Abs(before)*0.5 {
		return ShiftDeceleration // desaceleração
	}

	return ShiftBreakout // quebra de padrão
}

// Métodos para predição

func (td *TrendDetector) extractTrendSlope(trend *predicoes.VRTrend) float64 {
	// Converter força da tendência em slope normalizado
	baseSlope := trend.Strength * 0.1 // 10% máximo por período

	if trend.Type == predicoes.TrendDownward {
		return -baseSlope
	} else if trend.Type == predicoes.TrendUpward {
		return baseSlope
	}

	return 0.0
}

func (td *TrendDetector) getSeasonalFactor(seasonality *predicoes.SeasonalityInfo, date time.Time) float64 {
	month := int(date.Month())

	if factor, exists := seasonality.Pattern[month]; exists {
		// Calcular fator relativo à média anual
		totalSum := 0.0
		count := 0
		for _, v := range seasonality.Pattern {
			totalSum += v
			count++
		}

		if count > 0 {
			annualAvg := totalSum / float64(count)
			return factor / annualAvg
		}
	}

	return 1.0 // fator neutro
}

func (td *TrendDetector) calculatePredictionUncertainty(result *TrendAnalysisResult, periodsAhead int) float64 {
	baseUncertainty := 0.1 // 10% base

	// Aumentar incerteza com o tempo
	timeDecay := float64(periodsAhead) * 0.05 // 5% por período

	// Ajustar pela volatilidade
	volatilityFactor := result.Volatility * 0.5

	// Ajustar pela confiança
	confidenceFactor := (1 - result.Confidence) * 0.3

	totalUncertainty := baseUncertainty + timeDecay + volatilityFactor + confidenceFactor

	return math.Min(totalUncertainty, 0.5) // máximo 50% de incerteza
}

func (td *TrendDetector) identifyPredictionFactors(result *TrendAnalysisResult) []string {
	factors := []string{
		fmt.Sprintf("Tendência: %s", result.PrimaryTrend.Type),
		fmt.Sprintf("Confiança: %.1f%%", result.PrimaryTrend.Confidence*100),
	}

	if result.Seasonality.IsDetected {
		factors = append(factors, "Ajuste sazonal aplicado")
	}

	if len(result.SecondaryTrends) > 0 {
		factors = append(factors, "Tendências secundárias consideradas")
	}

	return factors
}

func (td *TrendDetector) generatePredictionAssumptions(result *TrendAnalysisResult) []string {
	assumptions := []string{
		"Padrões históricos se mantêm",
		"Não há mudanças estruturais significativas",
	}

	if result.Seasonality.IsDetected {
		assumptions = append(assumptions, "Padrão sazonal se repete")
	}

	if result.Volatility > td.config.VolatilityThreshold {
		assumptions = append(assumptions, "Alta volatilidade pode afetar precisão")
	}

	return assumptions
}

func (td *TrendDetector) assessPredictionReliability(result *TrendAnalysisResult) float64 {
	reliability := result.Confidence

	// Penalizar por alta volatilidade
	if result.Volatility > td.config.VolatilityThreshold {
		reliability *= 0.8
	}

	// Bonus por sazonalidade detectada
	if result.Seasonality.IsDetected && result.Seasonality.Confidence > 0.6 {
		reliability *= 1.1
	}

	return math.Min(reliability, 1.0)
}

func (td *TrendDetector) assessDataQuality(data []predicoes.HistoricalVRData) float64 {
	if len(data) == 0 {
		return 0.0
	}

	quality := 1.0

	// Penalizar por dados faltantes
	expectedPoints := 12 // esperamos dados mensais por 1 ano
	if len(data) < expectedPoints {
		quality *= float64(len(data)) / float64(expectedPoints)
	}

	// Penalizar por anomalias
	anomaliaCount := 0
	for _, d := range data {
		anomaliaCount += len(d.Anomalies)
	}

	if anomaliaCount > 0 {
		anomalyRate := float64(anomaliaCount) / float64(len(data))
		quality *= (1 - math.Min(anomalyRate, 0.3)) // máximo 30% de penalidade
	}

	return quality
}
