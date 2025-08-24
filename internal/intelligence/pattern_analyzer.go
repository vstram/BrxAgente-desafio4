package intelligence

import (
	"fmt"
	"math"
	"sort"
	"time"

	"BrxAgente-desafio4/internal/predicoes"
)

// PatternAnalyzer analisa padrões em dados históricos de VR
type PatternAnalyzer struct {
	config       PatternAnalyzerConfig
	patterns     []predicoes.Pattern
	cache        map[string]interface{}
	lastAnalysis time.Time
}

// PatternAnalyzerConfig configuração do analisador de padrões
type PatternAnalyzerConfig struct {
	MinDataPoints        int     `json:"min_data_points"`       // mínimo de pontos para análise
	CorrelationThreshold float64 `json:"correlation_threshold"` // limite para correlação significativa
	SeasonalityWindow    int     `json:"seasonality_window"`    // janela para análise sazonal
	AnomalyThreshold     float64 `json:"anomaly_threshold"`     // limite para detecção de anomalias
	ConfidenceLevel      float64 `json:"confidence_level"`      // nível de confiança mínimo
}

// ConsumptionPattern representa padrão de consumo identificado
type ConsumptionPattern struct {
	Sindicato       string                     `json:"sindicato"`
	PatternType     predicoes.PatternType      `json:"pattern_type"`
	Characteristics map[string]float64         `json:"characteristics"`
	Trend           *predicoes.VRTrend         `json:"trend"`
	Seasonality     *predicoes.SeasonalityInfo `json:"seasonality"`
	Stability       float64                    `json:"stability"`      // 0-1, estabilidade do padrão
	Predictability  float64                    `json:"predictability"` // 0-1, previsibilidade
}

// BehaviorPattern representa padrão comportamental
type BehaviorPattern struct {
	Type        string                  `json:"type"`
	Description string                  `json:"description"`
	Frequency   float64                 `json:"frequency"`   // frequência de ocorrência
	Entities    []string                `json:"entities"`    // colaboradores/sindicatos afetados
	Correlation map[string]float64      `json:"correlation"` // correlação com outras variáveis
	Impact      predicoes.PatternImpact `json:"impact"`
}

// NewPatternAnalyzer cria novo analisador de padrões
func NewPatternAnalyzer(config PatternAnalyzerConfig) *PatternAnalyzer {
	return &PatternAnalyzer{
		config:   config,
		patterns: make([]predicoes.Pattern, 0),
		cache:    make(map[string]interface{}),
	}
}

// AnalyzeConsumptionPatterns analisa padrões de consumo de VR
func (pa *PatternAnalyzer) AnalyzeConsumptionPatterns(data []predicoes.HistoricalVRData) ([]ConsumptionPattern, error) {
	if len(data) < pa.config.MinDataPoints {
		return nil, fmt.Errorf("dados insuficientes: %d pontos, mínimo: %d", len(data), pa.config.MinDataPoints)
	}

	// Agrupar dados por sindicato
	sindicatoData := pa.groupBySindicato(data)

	patterns := make([]ConsumptionPattern, 0)

	for sindicato, vrData := range sindicatoData {
		pattern, err := pa.analyzeConsumptionPattern(sindicato, vrData)
		if err != nil {
			continue // pular sindicatos com dados insuficientes
		}
		patterns = append(patterns, *pattern)
	}

	// Armazenar padrões para análise posterior
	pa.updatePatternsCache(patterns)

	return patterns, nil
}

// DetectSeasonality detecta sazonalidade em séries temporais
func (pa *PatternAnalyzer) DetectSeasonality(timeSeries predicoes.TimeSeries) (*predicoes.SeasonalityInfo, error) {
	if len(timeSeries.Points) < 12 {
		return &predicoes.SeasonalityInfo{
			IsDetected: false,
			Confidence: 0.0,
		}, nil
	}

	// Extrair componente sazonal usando decomposição simples
	seasonality := pa.extractSeasonalComponent(timeSeries)

	// Testar significância estatística da sazonalidade
	significance := pa.testSeasonalSignificance(timeSeries, seasonality)

	// Identificar padrões sazonais específicos
	patterns := pa.identifySeasonalPatterns(seasonality)

	info := &predicoes.SeasonalityInfo{
		IsDetected:   significance > pa.config.ConfidenceLevel,
		Period:       12, // mensal
		Amplitude:    pa.calculateSeasonalAmplitude(seasonality),
		PeakMonths:   patterns.Peaks,
		TroughMonths: patterns.Troughs,
		Confidence:   significance,
		Pattern:      seasonality,
	}

	return info, nil
}

// IdentifyOutlierPatterns identifica padrões de outliers
func (pa *PatternAnalyzer) IdentifyOutlierPatterns(data []predicoes.HistoricalVRData) ([]predicoes.OutlierPattern, error) {
	outlierPatterns := make([]predicoes.OutlierPattern, 0)

	// Agrupar por sindicato para análise individual
	sindicatoData := pa.groupBySindicato(data)

	for sindicato, vrData := range sindicatoData {
		// Detectar outliers usando método IQR
		outliers := pa.detectOutliers(vrData)

		if len(outliers) > 0 {
			// Analisar padrões nos outliers
			patterns := pa.analyzeOutlierPatterns(sindicato, outliers, vrData)
			outlierPatterns = append(outlierPatterns, patterns...)
		}
	}

	return outlierPatterns, nil
}

// AnalyzeCrossSindicatoPatterns analisa padrões entre sindicatos
func (pa *PatternAnalyzer) AnalyzeCrossSindicatoPatterns(data []predicoes.HistoricalVRData) ([]BehaviorPattern, error) {
	behaviorPatterns := make([]BehaviorPattern, 0)

	// Análise de correlação entre sindicatos
	correlations := pa.calculateCrossSindicatoCorrelations(data)

	// Identificar grupos de sindicatos com comportamento similar
	clusters := pa.clusterSimilarSindicatos(correlations)

	// Analisar padrões temporais sincronizados
	syncPatterns := pa.findSynchronizedPatterns(data)

	// Converter em padrões comportamentais
	for _, cluster := range clusters {
		pattern := BehaviorPattern{
			Type:        "cluster_comportamental",
			Description: fmt.Sprintf("Grupo de sindicatos com comportamento similar: %v", cluster.Members),
			Frequency:   cluster.Stability,
			Entities:    cluster.Members,
			Correlation: cluster.Correlations,
			Impact: predicoes.PatternImpact{
				Level:       pa.assessClusterImpact(cluster),
				Areas:       []string{"planejamento", "orçamento"},
				Magnitude:   cluster.Strength,
				Description: "Sindicatos que seguem padrões similares de consumo",
			},
		}
		behaviorPatterns = append(behaviorPatterns, pattern)
	}

	for _, sync := range syncPatterns {
		pattern := BehaviorPattern{
			Type:        "sincronizacao_temporal",
			Description: fmt.Sprintf("Padrão sincronizado: %s", sync.Description),
			Frequency:   sync.Frequency,
			Entities:    sync.Entities,
			Impact: predicoes.PatternImpact{
				Level:       sync.Impact,
				Areas:       []string{"processamento", "recursos"},
				Magnitude:   sync.Strength,
				Description: "Eventos que ocorrem simultaneamente entre sindicatos",
			},
		}
		behaviorPatterns = append(behaviorPatterns, pattern)
	}

	return behaviorPatterns, nil
}

// DetectAnomalousPatterns detecta padrões anômalos
func (pa *PatternAnalyzer) DetectAnomalousPatterns(data []predicoes.HistoricalVRData) ([]predicoes.Pattern, error) {
	anomalousPatterns := make([]predicoes.Pattern, 0)

	// Detectar mudanças estruturais
	structuralChanges := pa.detectStructuralChanges(data)

	// Detectar comportamentos atípicos
	atypicalBehaviors := pa.detectAtypicalBehaviors(data)

	// Detectar correlações inusitadas
	unusualCorrelations := pa.detectUnusualCorrelations(data)

	// Converter em padrões
	for _, change := range structuralChanges {
		pattern := predicoes.Pattern{
			ID:          pa.generatePatternID("structural"),
			Type:        predicoes.PatternAnomaly,
			Description: change.Description,
			Confidence:  change.Confidence,
			StartDate:   change.StartDate,
			EndDate:     change.EndDate,
			Entities:    change.AffectedEntities,
			Impact: predicoes.PatternImpact{
				Level:       change.Impact,
				Areas:       []string{"estrutural", "planejamento"},
				Magnitude:   change.Magnitude,
				Description: "Mudança estrutural detectada nos padrões",
			},
		}
		anomalousPatterns = append(anomalousPatterns, pattern)
	}

	for _, behavior := range atypicalBehaviors {
		pattern := predicoes.Pattern{
			ID:          pa.generatePatternID("atypical"),
			Type:        predicoes.PatternBehavior,
			Description: behavior.Description,
			Confidence:  behavior.Confidence,
			StartDate:   behavior.StartDate,
			Entities:    behavior.Entities,
			Impact:      behavior.Impact,
		}
		anomalousPatterns = append(anomalousPatterns, pattern)
	}

	for _, correlation := range unusualCorrelations {
		pattern := predicoes.Pattern{
			ID:          pa.generatePatternID("correlation"),
			Type:        predicoes.PatternAnomaly,
			Description: correlation.Description,
			Confidence:  correlation.Strength,
			StartDate:   time.Now().AddDate(0, -3, 0), // últimos 3 meses
			Entities:    correlation.Entities,
			Impact: predicoes.PatternImpact{
				Level:       predicoes.ImpactMedium,
				Areas:       []string{"correlação", "predição"},
				Magnitude:   correlation.Strength,
				Description: "Correlação inusitada entre entidades",
			},
		}
		anomalousPatterns = append(anomalousPatterns, pattern)
	}

	pa.patterns = append(pa.patterns, anomalousPatterns...)
	pa.lastAnalysis = time.Now()

	return anomalousPatterns, nil
}

// Métodos auxiliares internos

func (pa *PatternAnalyzer) groupBySindicato(data []predicoes.HistoricalVRData) map[string][]predicoes.HistoricalVRData {
	groups := make(map[string][]predicoes.HistoricalVRData)

	for _, d := range data {
		groups[d.Sindicato] = append(groups[d.Sindicato], d)
	}

	// Ordenar cada grupo por data
	for sindicato := range groups {
		sort.Slice(groups[sindicato], func(i, j int) bool {
			return groups[sindicato][i].Month.Before(groups[sindicato][j].Month)
		})
	}

	return groups
}

func (pa *PatternAnalyzer) analyzeConsumptionPattern(sindicato string, data []predicoes.HistoricalVRData) (*ConsumptionPattern, error) {
	if len(data) < 3 {
		return nil, fmt.Errorf("dados insuficientes para sindicato %s", sindicato)
	}

	// Extrair série temporal de VR total
	timeSeries := pa.createTimeSeriesFromVRData(data)

	// Detectar tendência
	trend := pa.detectTrendInVRData(data)

	// Detectar sazonalidade
	seasonality, _ := pa.DetectSeasonality(*timeSeries)

	// Calcular estabilidade do padrão
	stability := pa.calculatePatternStability(data)

	// Calcular previsibilidade
	predictability := pa.calculatePredictability(data, trend, seasonality)

	// Extrair características do padrão
	characteristics := pa.extractConsumptionCharacteristics(data)

	pattern := &ConsumptionPattern{
		Sindicato:       sindicato,
		PatternType:     pa.classifyConsumptionPattern(trend, seasonality, stability),
		Characteristics: characteristics,
		Trend:           trend,
		Seasonality:     seasonality,
		Stability:       stability,
		Predictability:  predictability,
	}

	return pattern, nil
}

func (pa *PatternAnalyzer) createTimeSeriesFromVRData(data []predicoes.HistoricalVRData) *predicoes.TimeSeries {
	ts := predicoes.NewTimeSeries("VR Total", "R$")

	for _, d := range data {
		ts.AddPoint(d.Month, d.TotalVR, map[string]interface{}{
			"num_colaboradores": d.NumColaboradores,
			"media_pessoa":      d.MediaPorPessoa,
		})
	}

	return ts
}

func (pa *PatternAnalyzer) detectTrendInVRData(data []predicoes.HistoricalVRData) *predicoes.VRTrend {
	values := make([]float64, len(data))
	for i, d := range data {
		values[i] = d.TotalVR
	}

	// Regressão linear simples
	n := float64(len(values))
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0

	for i, y := range values {
		x := float64(i + 1)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

	// Classificar tendência
	var trendType predicoes.TrendType
	avgValue := sumY / n
	normalizedSlope := math.Abs(slope) / avgValue

	if normalizedSlope < 0.02 { // menos de 2% de variação
		trendType = predicoes.TrendStable
	} else if slope > 0 {
		trendType = predicoes.TrendUpward
	} else {
		trendType = predicoes.TrendDownward
	}

	// Calcular confiança baseada no R²
	confidence := pa.calculateRSquared(values, slope, sumY/n)

	return &predicoes.VRTrend{
		Type:        trendType,
		Strength:    math.Min(normalizedSlope, 1.0),
		Period:      len(data),
		Confidence:  confidence,
		StartDate:   data[0].Month,
		Description: pa.describeTrend(trendType, normalizedSlope),
	}
}

func (pa *PatternAnalyzer) calculateRSquared(values []float64, slope, intercept float64) float64 {
	totalSumSquares := 0.0
	residualSumSquares := 0.0
	mean := 0.0

	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	for i, actual := range values {
		predicted := intercept + slope*float64(i+1)
		totalSumSquares += (actual - mean) * (actual - mean)
		residualSumSquares += (actual - predicted) * (actual - predicted)
	}

	if totalSumSquares == 0 {
		return 0
	}

	return 1.0 - (residualSumSquares / totalSumSquares)
}

func (pa *PatternAnalyzer) calculatePatternStability(data []predicoes.HistoricalVRData) float64 {
	if len(data) < 2 {
		return 0.0
	}

	// Calcular coeficiente de variação
	mean := 0.0
	for _, d := range data {
		mean += d.TotalVR
	}
	mean /= float64(len(data))

	variance := 0.0
	for _, d := range data {
		diff := d.TotalVR - mean
		variance += diff * diff
	}
	variance /= float64(len(data))

	stdDev := math.Sqrt(variance)
	cv := stdDev / mean

	// Converter CV em estabilidade (inverso)
	stability := 1.0 / (1.0 + cv)
	return math.Max(0, math.Min(stability, 1.0))
}

func (pa *PatternAnalyzer) calculatePredictability(data []predicoes.HistoricalVRData, trend *predicoes.VRTrend, seasonality *predicoes.SeasonalityInfo) float64 {
	basePredictability := trend.Confidence * 0.6 // peso da tendência

	if seasonality.IsDetected {
		basePredictability += seasonality.Confidence * 0.4 // peso da sazonalidade
	}

	// Penalizar por anomalias
	anomalyPenalty := float64(len(pa.countAnomalies(data))) * 0.1

	predictability := basePredictability - anomalyPenalty
	return math.Max(0, math.Min(predictability, 1.0))
}

func (pa *PatternAnalyzer) extractConsumptionCharacteristics(data []predicoes.HistoricalVRData) map[string]float64 {
	characteristics := make(map[string]float64)

	// Média de VR total
	totalVR := 0.0
	totalColaboradores := 0
	for _, d := range data {
		totalVR += d.TotalVR
		totalColaboradores += d.NumColaboradores
	}

	characteristics["media_vr_total"] = totalVR / float64(len(data))
	characteristics["media_colaboradores"] = float64(totalColaboradores) / float64(len(data))
	characteristics["vr_per_capita"] = totalVR / float64(totalColaboradores)

	// Crescimento month-over-month médio
	if len(data) > 1 {
		momGrowth := 0.0
		for i := 1; i < len(data); i++ {
			if data[i-1].TotalVR > 0 {
				growth := (data[i].TotalVR - data[i-1].TotalVR) / data[i-1].TotalVR
				momGrowth += growth
			}
		}
		characteristics["crescimento_mom"] = momGrowth / float64(len(data)-1)
	}

	// Volatilidade
	mean := characteristics["media_vr_total"]
	variance := 0.0
	for _, d := range data {
		diff := d.TotalVR - mean
		variance += diff * diff
	}
	characteristics["volatilidade"] = math.Sqrt(variance/float64(len(data))) / mean

	return characteristics
}

func (pa *PatternAnalyzer) classifyConsumptionPattern(trend *predicoes.VRTrend, seasonality *predicoes.SeasonalityInfo, stability float64) predicoes.PatternType {
	if seasonality.IsDetected && seasonality.Confidence > 0.7 {
		return predicoes.PatternSeasonal
	}

	if trend.Type != predicoes.TrendStable && trend.Confidence > 0.7 {
		return predicoes.PatternGrowth
	}

	if stability > 0.8 {
		return predicoes.PatternConsumption
	}

	return predicoes.PatternBehavior
}

func (pa *PatternAnalyzer) countAnomalies(data []predicoes.HistoricalVRData) []string {
	anomalies := make([]string, 0)
	for _, d := range data {
		anomalies = append(anomalies, d.Anomalies...)
	}
	return anomalies
}

func (pa *PatternAnalyzer) describeTrend(trendType predicoes.TrendType, strength float64) string {
	strengthDesc := "fraca"
	if strength > 0.05 {
		strengthDesc = "moderada"
	}
	if strength > 0.1 {
		strengthDesc = "forte"
	}

	switch trendType {
	case predicoes.TrendUpward:
		return fmt.Sprintf("Tendência de crescimento %s", strengthDesc)
	case predicoes.TrendDownward:
		return fmt.Sprintf("Tendência de decréscimo %s", strengthDesc)
	case predicoes.TrendStable:
		return "Comportamento estável"
	default:
		return "Padrão indefinido"
	}
}

// Métodos para análise de sazonalidade

type seasonalPatterns struct {
	Peaks   []int
	Troughs []int
}

func (pa *PatternAnalyzer) extractSeasonalComponent(ts predicoes.TimeSeries) map[int]float64 {
	monthlyData := make(map[int][]float64)

	// Agrupar por mês
	for _, point := range ts.Points {
		month := int(point.Timestamp.Month())
		monthlyData[month] = append(monthlyData[month], point.Value)
	}

	// Calcular média por mês
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

func (pa *PatternAnalyzer) testSeasonalSignificance(ts predicoes.TimeSeries, seasonal map[int]float64) float64 {
	if len(seasonal) < 4 {
		return 0.0
	}

	// Calcular média geral
	totalSum := 0.0
	count := 0
	for _, avg := range seasonal {
		totalSum += avg
		count++
	}
	overallAvg := totalSum / float64(count)

	// Calcular variabilidade sazonal
	seasonalVariance := 0.0
	for _, avg := range seasonal {
		diff := avg - overallAvg
		seasonalVariance += diff * diff
	}
	seasonalVariance /= float64(count)

	// Calcular variabilidade total dos dados
	totalVariance := 0.0
	for _, point := range ts.Points {
		diff := point.Value - overallAvg
		totalVariance += diff * diff
	}
	totalVariance /= float64(len(ts.Points))

	// Significância baseada na proporção da variância sazonal
	if totalVariance == 0 {
		return 0
	}

	significance := seasonalVariance / totalVariance
	return math.Min(significance, 1.0)
}

func (pa *PatternAnalyzer) identifySeasonalPatterns(seasonal map[int]float64) seasonalPatterns {
	if len(seasonal) == 0 {
		return seasonalPatterns{}
	}

	// Calcular média
	sum := 0.0
	for _, v := range seasonal {
		sum += v
	}
	avg := sum / float64(len(seasonal))

	// Identificar picos e vales
	var peaks, troughs []int
	threshold := avg * 0.1 // 10% da média

	for month, value := range seasonal {
		if value > avg+threshold {
			peaks = append(peaks, month)
		} else if value < avg-threshold {
			troughs = append(troughs, month)
		}
	}

	return seasonalPatterns{
		Peaks:   peaks,
		Troughs: troughs,
	}
}

func (pa *PatternAnalyzer) calculateSeasonalAmplitude(seasonal map[int]float64) float64 {
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

// Métodos para análise de outliers

func (pa *PatternAnalyzer) detectOutliers(data []predicoes.HistoricalVRData) []predicoes.HistoricalVRData {
	if len(data) < 4 {
		return nil
	}

	values := make([]float64, len(data))
	for i, d := range data {
		values[i] = d.TotalVR
	}

	// Calcular quartis usando método simples
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	q1Index := len(sorted) / 4
	q3Index := 3 * len(sorted) / 4

	q1 := sorted[q1Index]
	q3 := sorted[q3Index]
	iqr := q3 - q1

	// Definir limites
	lowerBound := q1 - 1.5*iqr
	upperBound := q3 + 1.5*iqr

	// Identificar outliers
	var outliers []predicoes.HistoricalVRData
	for i, value := range values {
		if value < lowerBound || value > upperBound {
			outliers = append(outliers, data[i])
		}
	}

	return outliers
}

func (pa *PatternAnalyzer) analyzeOutlierPatterns(sindicato string, outliers, allData []predicoes.HistoricalVRData) []predicoes.OutlierPattern {
	patterns := make([]predicoes.OutlierPattern, 0)

	if len(outliers) == 0 {
		return patterns
	}

	// Calcular frequência de outliers
	frequency := float64(len(outliers)) / float64(len(allData))

	// Calcular severidade média
	totalSeverity := 0.0
	mean := pa.calculateMean(allData)

	for _, outlier := range outliers {
		severity := math.Abs(outlier.TotalVR-mean) / mean
		totalSeverity += severity
	}
	avgSeverity := totalSeverity / float64(len(outliers))

	// Classificar tipo de outlier
	outType := pa.classifyOutlierType(outliers, allData)

	pattern := predicoes.OutlierPattern{
		Type:        outType,
		Frequency:   frequency,
		Severity:    avgSeverity,
		Affected:    []string{sindicato},
		Description: pa.describeOutlierPattern(outType, frequency, avgSeverity),
		Confidence:  math.Min(frequency*2, 1.0), // mais outliers = maior confiança
	}

	patterns = append(patterns, pattern)
	return patterns
}

func (pa *PatternAnalyzer) calculateMean(data []predicoes.HistoricalVRData) float64 {
	sum := 0.0
	for _, d := range data {
		sum += d.TotalVR
	}
	return sum / float64(len(data))
}

func (pa *PatternAnalyzer) classifyOutlierType(outliers, allData []predicoes.HistoricalVRData) predicoes.OutlierType {
	if len(outliers) == 1 {
		return predicoes.OutlierSpike
	}

	// Verificar se outliers são consecutivos (sustained)
	consecutive := 0
	maxConsecutive := 0

	outlierDates := make(map[time.Time]bool)
	for _, o := range outliers {
		outlierDates[o.Month] = true
	}

	for i := 1; i < len(allData); i++ {
		if outlierDates[allData[i].Month] && outlierDates[allData[i-1].Month] {
			consecutive++
			if consecutive > maxConsecutive {
				maxConsecutive = consecutive
			}
		} else {
			consecutive = 0
		}
	}

	if maxConsecutive >= 2 {
		return predicoes.OutlierSustained
	}

	// Verificar se há padrão cíclico
	if pa.detectCyclicOutlierPattern(outliers) {
		return predicoes.OutlierCyclic
	}

	return predicoes.OutlierSpike
}

func (pa *PatternAnalyzer) detectCyclicOutlierPattern(outliers []predicoes.HistoricalVRData) bool {
	if len(outliers) < 4 {
		return false
	}

	// Verificar se outliers ocorrem em intervalos regulares
	intervals := make([]int, len(outliers)-1)
	for i := 1; i < len(outliers); i++ {
		diff := int(outliers[i].Month.Sub(outliers[i-1].Month).Hours() / 24 / 30) // meses
		intervals[i-1] = diff
	}

	// Verificar consistência dos intervalos
	if len(intervals) < 2 {
		return false
	}

	tolerance := 1 // 1 mês de tolerância
	baseInterval := intervals[0]

	for _, interval := range intervals[1:] {
		if math.Abs(float64(interval-baseInterval)) > float64(tolerance) {
			return false
		}
	}

	return true
}

func (pa *PatternAnalyzer) describeOutlierPattern(outType predicoes.OutlierType, frequency, severity float64) string {
	freqDesc := "baixa"
	if frequency > 0.2 {
		freqDesc = "moderada"
	}
	if frequency > 0.4 {
		freqDesc = "alta"
	}

	sevDesc := "leve"
	if severity > 0.5 {
		sevDesc = "moderada"
	}
	if severity > 1.0 {
		sevDesc = "severa"
	}

	switch outType {
	case predicoes.OutlierSpike:
		return fmt.Sprintf("Picos isolados com frequência %s e severidade %s", freqDesc, sevDesc)
	case predicoes.OutlierSustained:
		return fmt.Sprintf("Valores sustentados fora do padrão com frequência %s e severidade %s", freqDesc, sevDesc)
	case predicoes.OutlierCyclic:
		return fmt.Sprintf("Outliers cíclicos com frequência %s e severidade %s", freqDesc, sevDesc)
	default:
		return fmt.Sprintf("Padrão de outliers com frequência %s e severidade %s", freqDesc, sevDesc)
	}
}

// Métodos para análise cross-sindicato

type sindicatoCluster struct {
	Members      []string
	Stability    float64
	Strength     float64
	Correlations map[string]float64
}

type synchronizedPattern struct {
	Description string
	Frequency   float64
	Entities    []string
	Strength    float64
	Impact      predicoes.ImpactLevel
}

type structuralChange struct {
	Description      string
	Confidence       float64
	StartDate        time.Time
	EndDate          *time.Time
	AffectedEntities []string
	Impact           predicoes.ImpactLevel
	Magnitude        float64
}

type atypicalBehavior struct {
	Description string
	Confidence  float64
	StartDate   time.Time
	Entities    []string
	Impact      predicoes.PatternImpact
}

type unusualCorrelation struct {
	Description string
	Strength    float64
	Entities    []string
}

func (pa *PatternAnalyzer) calculateCrossSindicatoCorrelations(data []predicoes.HistoricalVRData) map[string]map[string]float64 {
	// Agrupar dados por sindicato
	sindicatoData := pa.groupBySindicato(data)

	// Calcular matriz de correlação
	correlations := make(map[string]map[string]float64)
	sindicatos := make([]string, 0, len(sindicatoData))

	for sindicato := range sindicatoData {
		sindicatos = append(sindicatos, sindicato)
		correlations[sindicato] = make(map[string]float64)
	}

	// Calcular correlação entre cada par de sindicatos
	for i, sindA := range sindicatos {
		for j, sindB := range sindicatos {
			if i != j {
				correlation := pa.calculateCorrelation(sindicatoData[sindA], sindicatoData[sindB])
				correlations[sindA][sindB] = correlation
			} else {
				correlations[sindA][sindB] = 1.0
			}
		}
	}

	return correlations
}

func (pa *PatternAnalyzer) calculateCorrelation(dataA, dataB []predicoes.HistoricalVRData) float64 {
	// Encontrar períodos comuns
	commonPeriods := pa.findCommonPeriods(dataA, dataB)
	if len(commonPeriods) < 3 {
		return 0.0
	}

	// Extrair valores para períodos comuns
	valuesA := make([]float64, len(commonPeriods))
	valuesB := make([]float64, len(commonPeriods))

	for i, period := range commonPeriods {
		valuesA[i] = pa.getValueForPeriod(dataA, period)
		valuesB[i] = pa.getValueForPeriod(dataB, period)
	}

	// Calcular coeficiente de correlação de Pearson
	return pa.pearsonCorrelation(valuesA, valuesB)
}

func (pa *PatternAnalyzer) findCommonPeriods(dataA, dataB []predicoes.HistoricalVRData) []time.Time {
	periodsA := make(map[string]time.Time)
	for _, d := range dataA {
		key := d.Month.Format("2006-01")
		periodsA[key] = d.Month
	}

	var commonPeriods []time.Time
	for _, d := range dataB {
		key := d.Month.Format("2006-01")
		if _, exists := periodsA[key]; exists {
			commonPeriods = append(commonPeriods, d.Month)
		}
	}

	sort.Slice(commonPeriods, func(i, j int) bool {
		return commonPeriods[i].Before(commonPeriods[j])
	})

	return commonPeriods
}

func (pa *PatternAnalyzer) getValueForPeriod(data []predicoes.HistoricalVRData, period time.Time) float64 {
	key := period.Format("2006-01")
	for _, d := range data {
		if d.Month.Format("2006-01") == key {
			return d.TotalVR
		}
	}
	return 0.0
}

func (pa *PatternAnalyzer) pearsonCorrelation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0.0
	}

	// Calcular médias
	meanX, meanY := 0.0, 0.0
	for i := 0; i < len(x); i++ {
		meanX += x[i]
		meanY += y[i]
	}
	meanX /= float64(len(x))
	meanY /= float64(len(y))

	// Calcular correlação
	numerator := 0.0
	sumSqX := 0.0
	sumSqY := 0.0

	for i := 0; i < len(x); i++ {
		dx := x[i] - meanX
		dy := y[i] - meanY
		numerator += dx * dy
		sumSqX += dx * dx
		sumSqY += dy * dy
	}

	denominator := math.Sqrt(sumSqX * sumSqY)
	if denominator == 0 {
		return 0.0
	}

	return numerator / denominator
}

func (pa *PatternAnalyzer) clusterSimilarSindicatos(correlations map[string]map[string]float64) []sindicatoCluster {
	clusters := make([]sindicatoCluster, 0)

	// Implementação simples de clustering baseado em threshold
	processed := make(map[string]bool)

	for sindicato, correls := range correlations {
		if processed[sindicato] {
			continue
		}

		cluster := sindicatoCluster{
			Members:      []string{sindicato},
			Correlations: make(map[string]float64),
		}

		// Encontrar sindicatos similares
		for other, correlation := range correls {
			if !processed[other] && correlation > pa.config.CorrelationThreshold {
				cluster.Members = append(cluster.Members, other)
				cluster.Correlations[other] = correlation
				processed[other] = true
			}
		}

		if len(cluster.Members) > 1 {
			cluster.Stability = pa.calculateClusterStability(cluster.Correlations)
			cluster.Strength = pa.calculateClusterStrength(cluster.Correlations)
			clusters = append(clusters, cluster)
		}

		processed[sindicato] = true
	}

	return clusters
}

func (pa *PatternAnalyzer) calculateClusterStability(correlations map[string]float64) float64 {
	if len(correlations) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, corr := range correlations {
		sum += corr
	}

	return sum / float64(len(correlations))
}

func (pa *PatternAnalyzer) calculateClusterStrength(correlations map[string]float64) float64 {
	if len(correlations) == 0 {
		return 0.0
	}

	min := math.Inf(1)
	for _, corr := range correlations {
		if corr < min {
			min = corr
		}
	}

	return min // força do cluster é determinada pela correlação mais fraca
}

func (pa *PatternAnalyzer) assessClusterImpact(cluster sindicatoCluster) predicoes.ImpactLevel {
	if cluster.Strength > 0.8 && len(cluster.Members) > 3 {
		return predicoes.ImpactHigh
	} else if cluster.Strength > 0.6 && len(cluster.Members) > 2 {
		return predicoes.ImpactMedium
	}
	return predicoes.ImpactLow
}

func (pa *PatternAnalyzer) findSynchronizedPatterns(data []predicoes.HistoricalVRData) []synchronizedPattern {
	// Implementação simplificada - detectar mudanças simultâneas
	patterns := make([]synchronizedPattern, 0)

	// Agrupar por período
	periodData := make(map[string][]predicoes.HistoricalVRData)
	for _, d := range data {
		period := d.Month.Format("2006-01")
		periodData[period] = append(periodData[period], d)
	}

	// Detectar períodos com mudanças significativas em múltiplos sindicatos
	for period, periodVR := range periodData {
		if len(periodVR) < 3 {
			continue
		}

		// Verificar se houve mudança significativa comparado ao período anterior
		prevPeriod := pa.getPreviousPeriod(period)
		if prevData, exists := periodData[prevPeriod]; exists {
			syncPattern := pa.detectSynchronizedChange(periodVR, prevData)
			if syncPattern != nil {
				patterns = append(patterns, *syncPattern)
			}
		}
	}

	return patterns
}

func (pa *PatternAnalyzer) getPreviousPeriod(period string) string {
	t, err := time.Parse("2006-01", period)
	if err != nil {
		return ""
	}
	prev := t.AddDate(0, -1, 0)
	return prev.Format("2006-01")
}

func (pa *PatternAnalyzer) detectSynchronizedChange(current, previous []predicoes.HistoricalVRData) *synchronizedPattern {
	// Mapear dados por sindicato
	currentMap := make(map[string]float64)
	previousMap := make(map[string]float64)

	for _, d := range current {
		currentMap[d.Sindicato] = d.TotalVR
	}
	for _, d := range previous {
		previousMap[d.Sindicato] = d.TotalVR
	}

	// Calcular mudanças
	significantChanges := 0
	totalChanges := 0
	affectedSindicatos := make([]string, 0)

	for sindicato, currentVR := range currentMap {
		if prevVR, exists := previousMap[sindicato]; exists {
			change := math.Abs(currentVR-prevVR) / prevVR
			totalChanges++

			if change > 0.15 { // mudança > 15%
				significantChanges++
				affectedSindicatos = append(affectedSindicatos, sindicato)
			}
		}
	}

	// Verificar se é um padrão sincronizado significativo
	if significantChanges >= 3 && float64(significantChanges)/float64(totalChanges) > 0.5 {
		return &synchronizedPattern{
			Description: fmt.Sprintf("Mudança sincronizada detectada em %d sindicatos", significantChanges),
			Frequency:   float64(significantChanges) / float64(totalChanges),
			Entities:    affectedSindicatos,
			Strength:    float64(significantChanges) / 10.0, // normalizar
			Impact:      predicoes.ImpactMedium,
		}
	}

	return nil
}

// Métodos para detecção de padrões anômalos

func (pa *PatternAnalyzer) detectStructuralChanges(data []predicoes.HistoricalVRData) []structuralChange {
	changes := make([]structuralChange, 0)

	// Agrupar por sindicato
	sindicatoData := pa.groupBySindicato(data)

	for sindicato, vrData := range sindicatoData {
		if len(vrData) < 6 { // precisa de pelo menos 6 pontos
			continue
		}

		// Detectar mudanças estruturais usando changepoint detection simples
		changepoints := pa.detectChangepoints(vrData)

		for _, cp := range changepoints {
			change := structuralChange{
				Description:      fmt.Sprintf("Mudança estrutural detectada no sindicato %s", sindicato),
				Confidence:       cp.Confidence,
				StartDate:        cp.Date,
				AffectedEntities: []string{sindicato},
				Impact:           pa.assessStructuralChangeImpact(cp.Magnitude),
				Magnitude:        cp.Magnitude,
			}
			changes = append(changes, change)
		}
	}

	return changes
}

type changepoint struct {
	Date       time.Time
	Confidence float64
	Magnitude  float64
}

func (pa *PatternAnalyzer) detectChangepoints(data []predicoes.HistoricalVRData) []changepoint {
	if len(data) < 6 {
		return nil
	}

	changepoints := make([]changepoint, 0)
	minSegmentSize := 3

	// Implementação simples de detecção de changepoint
	for i := minSegmentSize; i < len(data)-minSegmentSize; i++ {
		// Calcular médias antes e depois do ponto
		meanBefore := pa.calculateSegmentMean(data[:i])
		meanAfter := pa.calculateSegmentMean(data[i:])

		// Calcular magnitude da mudança
		magnitude := math.Abs(meanAfter-meanBefore) / meanBefore

		if magnitude > 0.3 { // mudança > 30%
			// Calcular confiança baseada na consistência da mudança
			confidence := pa.calculateChangepointConfidence(data, i, meanBefore, meanAfter)

			if confidence > 0.7 {
				cp := changepoint{
					Date:       data[i].Month,
					Confidence: confidence,
					Magnitude:  magnitude,
				}
				changepoints = append(changepoints, cp)
			}
		}
	}

	return changepoints
}

func (pa *PatternAnalyzer) calculateSegmentMean(data []predicoes.HistoricalVRData) float64 {
	sum := 0.0
	for _, d := range data {
		sum += d.TotalVR
	}
	return sum / float64(len(data))
}

func (pa *PatternAnalyzer) calculateChangepointConfidence(data []predicoes.HistoricalVRData, changeIndex int, meanBefore, meanAfter float64) float64 {
	// Calcular variância dentro dos segmentos
	varBefore := pa.calculateSegmentVariance(data[:changeIndex], meanBefore)
	varAfter := pa.calculateSegmentVariance(data[changeIndex:], meanAfter)

	// Calcular variância entre segmentos
	overallMean := (meanBefore + meanAfter) / 2
	varBetween := (meanBefore-overallMean)*(meanBefore-overallMean) + (meanAfter-overallMean)*(meanAfter-overallMean)

	// Confiança baseada na razão entre variância entre segmentos e dentro dos segmentos
	withinVar := (varBefore + varAfter) / 2
	if withinVar == 0 {
		return 1.0
	}

	confidence := varBetween / withinVar
	return math.Min(confidence/10.0, 1.0) // normalizar
}

func (pa *PatternAnalyzer) calculateSegmentVariance(data []predicoes.HistoricalVRData, mean float64) float64 {
	if len(data) <= 1 {
		return 0
	}

	sumSquaredDiffs := 0.0
	for _, d := range data {
		diff := d.TotalVR - mean
		sumSquaredDiffs += diff * diff
	}

	return sumSquaredDiffs / float64(len(data)-1)
}

func (pa *PatternAnalyzer) assessStructuralChangeImpact(magnitude float64) predicoes.ImpactLevel {
	if magnitude > 0.8 {
		return predicoes.ImpactCritical
	} else if magnitude > 0.5 {
		return predicoes.ImpactHigh
	} else if magnitude > 0.3 {
		return predicoes.ImpactMedium
	}
	return predicoes.ImpactLow
}

func (pa *PatternAnalyzer) detectAtypicalBehaviors(data []predicoes.HistoricalVRData) []atypicalBehavior {
	behaviors := make([]atypicalBehavior, 0)

	// Detectar sindicatos com comportamento muito diferente da média
	overallStats := pa.calculateOverallStats(data)
	sindicatoData := pa.groupBySindicato(data)

	for sindicato, vrData := range sindicatoData {
		stats := pa.calculateSindicatoStats(vrData)

		// Verificar desvios significativos
		if pa.isAtypicalBehavior(stats, overallStats) {
			behavior := atypicalBehavior{
				Description: fmt.Sprintf("Comportamento atípico no sindicato %s", sindicato),
				Confidence:  pa.calculateBehaviorConfidence(stats, overallStats),
				StartDate:   vrData[0].Month,
				Entities:    []string{sindicato},
				Impact: predicoes.PatternImpact{
					Level:       predicoes.ImpactMedium,
					Areas:       []string{"predição", "planejamento"},
					Magnitude:   pa.calculateBehaviorMagnitude(stats, overallStats),
					Description: "Comportamento significativamente diferente da média",
				},
			}
			behaviors = append(behaviors, behavior)
		}
	}

	return behaviors
}

type statsData struct {
	Mean       float64
	Variance   float64
	GrowthRate float64
	Volatility float64
}

func (pa *PatternAnalyzer) calculateOverallStats(data []predicoes.HistoricalVRData) statsData {
	if len(data) == 0 {
		return statsData{}
	}

	sum := 0.0
	for _, d := range data {
		sum += d.TotalVR
	}
	mean := sum / float64(len(data))

	variance := 0.0
	for _, d := range data {
		diff := d.TotalVR - mean
		variance += diff * diff
	}
	variance /= float64(len(data))

	return statsData{
		Mean:       mean,
		Variance:   variance,
		GrowthRate: 0, // calcularia crescimento geral
		Volatility: math.Sqrt(variance) / mean,
	}
}

func (pa *PatternAnalyzer) calculateSindicatoStats(data []predicoes.HistoricalVRData) statsData {
	return pa.calculateOverallStats(data)
}

func (pa *PatternAnalyzer) isAtypicalBehavior(sindicatoStats, overallStats statsData) bool {
	// Verificar se está fora de 2 desvios padrão
	meanDiff := math.Abs(sindicatoStats.Mean - overallStats.Mean)
	threshold := 2 * math.Sqrt(overallStats.Variance)

	return meanDiff > threshold
}

func (pa *PatternAnalyzer) calculateBehaviorConfidence(sindicatoStats, overallStats statsData) float64 {
	meanDiff := math.Abs(sindicatoStats.Mean - overallStats.Mean)
	stdDev := math.Sqrt(overallStats.Variance)

	if stdDev == 0 {
		return 0.5
	}

	zScore := meanDiff / stdDev
	confidence := math.Min(zScore/3.0, 1.0) // normalizar z-score

	return confidence
}

func (pa *PatternAnalyzer) calculateBehaviorMagnitude(sindicatoStats, overallStats statsData) float64 {
	return math.Abs(sindicatoStats.Mean-overallStats.Mean) / overallStats.Mean
}

func (pa *PatternAnalyzer) detectUnusualCorrelations(data []predicoes.HistoricalVRData) []unusualCorrelation {
	correlations := make([]unusualCorrelation, 0)

	// Calcular correlações cross-sindicato
	crossCorrelations := pa.calculateCrossSindicatoCorrelations(data)

	// Identificar correlações inusitadas (muito altas ou inesperadas)
	for sindA, correls := range crossCorrelations {
		for sindB, correlation := range correls {
			if sindA != sindB && math.Abs(correlation) > 0.9 {
				unusual := unusualCorrelation{
					Description: fmt.Sprintf("Correlação inusitadamente alta entre %s e %s (%.3f)", sindA, sindB, correlation),
					Strength:    math.Abs(correlation),
					Entities:    []string{sindA, sindB},
				}
				correlations = append(correlations, unusual)
			}
		}
	}

	return correlations
}

// Métodos auxiliares

func (pa *PatternAnalyzer) updatePatternsCache(patterns []ConsumptionPattern) {
	pa.cache["consumption_patterns"] = patterns
	pa.cache["last_consumption_analysis"] = time.Now()
}

func (pa *PatternAnalyzer) generatePatternID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// GetCachedPatterns retorna padrões em cache
func (pa *PatternAnalyzer) GetCachedPatterns() []predicoes.Pattern {
	return pa.patterns
}

// GetLastAnalysisTime retorna timestamp da última análise
func (pa *PatternAnalyzer) GetLastAnalysisTime() time.Time {
	return pa.lastAnalysis
}
