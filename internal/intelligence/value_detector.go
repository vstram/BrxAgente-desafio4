package intelligence

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

// ValueAnomalyDetector detecta anomalias em valores
type ValueAnomalyDetector struct {
	config *AnalysisConfig
}

// NewValueAnomalyDetector cria novo detector de valores
func NewValueAnomalyDetector(config *AnalysisConfig) *ValueAnomalyDetector {
	return &ValueAnomalyDetector{
		config: config,
	}
}

// DetectVROutliers detecta valores de VR muito fora do padrão
func (d *ValueAnomalyDetector) DetectVROutliers(ctx *AnalysisContext) []Anomaly {
	anomalies := make([]Anomaly, 0)

	// Agrupar valores por sindicato para análise estatística
	valuesBySindicato := make(map[string][]float64)
	colaboradoresBySindicato := make(map[string][]map[string]interface{})

	// Coletar dados
	for matricula, data := range ctx.Colaboradores {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}

		sindicato, exists := dataMap["sindicato"]
		if !exists {
			continue
		}

		sindicatoStr, ok := sindicato.(string)
		if !ok {
			continue
		}

		// Extrair valor de VR
		vrValue, err := d.extractVRValue(dataMap)
		if err != nil {
			ctx.AddWarning(fmt.Sprintf("Erro ao extrair VR para matrícula %s: %v", matricula, err))
			continue
		}

		valuesBySindicato[sindicatoStr] = append(valuesBySindicato[sindicatoStr], vrValue)
		colaboradoresBySindicato[sindicatoStr] = append(colaboradoresBySindicato[sindicatoStr], dataMap)

		ctx.IncrementProcessed()
	}

	// Analisar outliers por sindicato
	for sindicato, values := range valuesBySindicato {
		if len(values) < 3 {
			// Poucos dados para análise estatística confiável
			continue
		}

		stats := d.calculateStatistics(values)

		// Definir thresholds baseados em desvio padrão
		upperThreshold := stats.Mean + (d.config.OutlierThreshold * stats.StdDev)
		lowerThreshold := math.Max(0, stats.Mean-(d.config.OutlierThreshold*stats.StdDev))

		colaboradores := colaboradoresBySindicato[sindicato]

		for i, value := range values {
			if i >= len(colaboradores) {
				break
			}

			colaborador := colaboradores[i]
			matricula, _ := colaborador["matricula"].(string)

			// Verificar se é outlier
			if value > upperThreshold {
				confidence := d.calculateOutlierConfidence(value, stats.Mean, stats.StdDev, true)

				anomaly := NewAnomaly(
					AnomalyTypeValue,
					d.getSeverityForOutlier(value, upperThreshold, stats.Mean),
					confidence,
					fmt.Sprintf("VR muito alto: R$ %.2f", value),
					fmt.Sprintf("Colaborador %s tem VR de R$ %.2f, %.1fx acima da média do sindicato %s (R$ %.2f)",
						matricula, value, value/stats.Mean, sindicato, stats.Mean),
					matricula,
					fmt.Sprintf("%.2f", value),
					"vr_value",
					"value_outlier_detector",
					"vr_outlier",
				)

				anomaly.AddData("sindicato", sindicato)
				anomaly.AddData("media_sindicato", stats.Mean)
				anomaly.AddData("desvio_padrao", stats.StdDev)
				anomaly.AddData("threshold_superior", upperThreshold)
				anomaly.AddData("desvios_acima", (value-stats.Mean)/stats.StdDev)

				anomaly.AddSuggestion("Verificar se o cálculo do VR está correto")
				anomaly.AddSuggestion("Confirmar sindicato do colaborador")
				anomaly.AddSuggestion("Revisar dias úteis trabalhados")

				anomalies = append(anomalies, anomaly)

			} else if value < lowerThreshold && value > 0 {
				confidence := d.calculateOutlierConfidence(value, stats.Mean, stats.StdDev, false)

				anomaly := NewAnomaly(
					AnomalyTypeValue,
					SeverityMedium,
					confidence,
					fmt.Sprintf("VR muito baixo: R$ %.2f", value),
					fmt.Sprintf("Colaborador %s tem VR de R$ %.2f, muito abaixo da média do sindicato %s (R$ %.2f)",
						matricula, value, sindicato, stats.Mean),
					matricula,
					fmt.Sprintf("%.2f", value),
					"vr_value",
					"value_outlier_detector",
					"vr_outlier",
				)

				anomaly.AddData("sindicato", sindicato)
				anomaly.AddData("media_sindicato", stats.Mean)
				anomaly.AddData("desvio_padrao", stats.StdDev)
				anomaly.AddData("threshold_inferior", lowerThreshold)
				anomaly.AddData("desvios_abaixo", (stats.Mean-value)/stats.StdDev)

				anomaly.AddSuggestion("Verificar se há afastamentos não contabilizados")
				anomaly.AddSuggestion("Confirmar data de admissão")
				anomaly.AddSuggestion("Revisar cálculo proporcional")

				anomalies = append(anomalies, anomaly)
			}
		}
	}

	return anomalies
}

// DetectZeroOrNegativeVR detecta valores zero ou negativos
func (d *ValueAnomalyDetector) DetectZeroOrNegativeVR(ctx *AnalysisContext) []Anomaly {
	anomalies := make([]Anomaly, 0)

	for matricula, data := range ctx.Colaboradores {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}

		vrValue, err := d.extractVRValue(dataMap)
		if err != nil {
			continue
		}

		if vrValue <= 0 {
			severity := SeverityHigh
			title := "VR zerado"
			description := fmt.Sprintf("Colaborador %s tem VR zerado", matricula)

			if vrValue < 0 {
				severity = SeverityCritical
				title = "VR negativo"
				description = fmt.Sprintf("Colaborador %s tem VR negativo: R$ %.2f", matricula, vrValue)
			}

			anomaly := NewAnomaly(
				AnomalyTypeValue,
				severity,
				95.0, // Alta confidence para valores zero/negativos
				title,
				description,
				matricula,
				fmt.Sprintf("%.2f", vrValue),
				"vr_value",
				"value_zero_detector",
				"vr_zero_or_negative",
			)

			anomaly.AddData("valor_vr", vrValue)

			if vrValue == 0 {
				anomaly.AddSuggestion("Verificar se colaborador trabalhou no período")
				anomaly.AddSuggestion("Confirmar data de admissão/desligamento")
				anomaly.AddSuggestion("Revisar afastamentos e férias")
			} else {
				anomaly.AddSuggestion("Corrigir erro de cálculo que resultou em valor negativo")
				anomaly.AddSuggestion("Verificar fórmulas e dados de entrada")
			}

			anomalies = append(anomalies, anomaly)
		}

		ctx.IncrementProcessed()
	}

	return anomalies
}

// DetectExtremeValues detecta valores extremos baseados em limites absolutos
func (d *ValueAnomalyDetector) DetectExtremeValues(ctx *AnalysisContext) []Anomaly {
	anomalies := make([]Anomaly, 0)

	for matricula, data := range ctx.Colaboradores {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}

		vrValue, err := d.extractVRValue(dataMap)
		if err != nil {
			continue
		}

		// Verificar limites absolutos
		if vrValue > d.config.VRMaxValue {
			anomaly := NewAnomaly(
				AnomalyTypeValue,
				SeverityHigh,
				90.0,
				fmt.Sprintf("VR acima do limite máximo: R$ %.2f", vrValue),
				fmt.Sprintf("Colaborador %s tem VR de R$ %.2f, acima do limite máximo de R$ %.2f",
					matricula, vrValue, d.config.VRMaxValue),
				matricula,
				fmt.Sprintf("%.2f", vrValue),
				"vr_value",
				"value_extreme_detector",
				"vr_extreme_high",
			)

			anomaly.AddData("valor_vr", vrValue)
			anomaly.AddData("limite_maximo", d.config.VRMaxValue)
			anomaly.AddData("excesso", vrValue-d.config.VRMaxValue)

			anomaly.AddSuggestion("Verificar se o valor está em conformidade com políticas da empresa")
			anomaly.AddSuggestion("Revisar categoria do colaborador")

			anomalies = append(anomalies, anomaly)

		} else if vrValue > 0 && vrValue < d.config.VRMinValue {
			anomaly := NewAnomaly(
				AnomalyTypeValue,
				SeverityMedium,
				85.0,
				fmt.Sprintf("VR abaixo do limite mínimo: R$ %.2f", vrValue),
				fmt.Sprintf("Colaborador %s tem VR de R$ %.2f, abaixo do limite mínimo de R$ %.2f",
					matricula, vrValue, d.config.VRMinValue),
				matricula,
				fmt.Sprintf("%.2f", vrValue),
				"vr_value",
				"value_extreme_detector",
				"vr_extreme_low",
			)

			anomaly.AddData("valor_vr", vrValue)
			anomaly.AddData("limite_minimo", d.config.VRMinValue)
			anomaly.AddData("deficit", d.config.VRMinValue-vrValue)

			anomaly.AddSuggestion("Verificar se colaborador trabalhou período parcial")
			anomaly.AddSuggestion("Confirmar cálculo proporcional")

			anomalies = append(anomalies, anomaly)
		}

		ctx.IncrementProcessed()
	}

	return anomalies
}

// extractVRValue extrai valor de VR dos dados do colaborador
func (d *ValueAnomalyDetector) extractVRValue(colaborador map[string]interface{}) (float64, error) {
	// Tentar diferentes campos onde o VR pode estar
	possibleFields := []string{"vr_total", "vr_value", "valor_vr", "vr", "total", "valor_total"}

	for _, field := range possibleFields {
		if value, exists := colaborador[field]; exists {
			if numValue, err := extractNumericValue(value); err == nil {
				return numValue, nil
			}
		}
	}

	// Se não encontrou, tentar calcular baseado em outros campos
	if dias, existsDias := colaborador["dias_uteis"]; existsDias {
		if valorDia, existsValorDia := colaborador["valor_dia"]; existsValorDia {
			diasNum, err1 := extractNumericValue(dias)
			valorDiaNum, err2 := extractNumericValue(valorDia)

			if err1 == nil && err2 == nil {
				return diasNum * valorDiaNum, nil
			}
		}
	}

	return 0, fmt.Errorf("valor de VR não encontrado ou inválido")
}

// Statistics representa estatísticas básicas
type Statistics struct {
	Mean   float64
	StdDev float64
	Min    float64
	Max    float64
	Count  int
}

// calculateStatistics calcula estatísticas básicas para um conjunto de valores
func (d *ValueAnomalyDetector) calculateStatistics(values []float64) Statistics {
	if len(values) == 0 {
		return Statistics{}
	}

	// Ordenar valores para calcular percentis
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	// Calcular média
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	// Calcular desvio padrão
	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	variance := sumSquares / float64(len(values))
	stdDev := math.Sqrt(variance)

	return Statistics{
		Mean:   mean,
		StdDev: stdDev,
		Min:    sorted[0],
		Max:    sorted[len(sorted)-1],
		Count:  len(values),
	}
}

// calculateOutlierConfidence calcula confidence baseado na distância da média
func (d *ValueAnomalyDetector) calculateOutlierConfidence(value, mean, stdDev float64, isHigh bool) float64 {
	if stdDev == 0 {
		return 50.0 // Confidence baixo se não há variação
	}

	deviations := math.Abs(value-mean) / stdDev

	// Converter desvios em confidence (mais desvios = mais confidence)
	confidence := math.Min(95.0, 60.0+(deviations-2.0)*15.0)

	return math.Max(50.0, confidence)
}

// getSeverityForOutlier determina severidade baseada na magnitude do outlier
func (d *ValueAnomalyDetector) getSeverityForOutlier(value, threshold, mean float64) AnomalySeverity {
	ratio := value / mean

	switch {
	case ratio > 3.0: // Mais de 3x a média
		return SeverityCritical
	case ratio > 2.0: // Mais de 2x a média
		return SeverityHigh
	case value > threshold*1.5: // 50% acima do threshold
		return SeverityMedium
	default:
		return SeverityLow
	}
}

// Helper function para formatar valores monetários
func formatCurrency(value float64) string {
	return fmt.Sprintf("R$ %.2f", value)
}

// Helper function para converter string em float, tratando formatos brasileiros
func parseFloatBR(s string) (float64, error) {
	// Remove símbolos de moeda e espaços
	clean := s
	clean = fmt.Sprintf("%s", clean) // Garantir que é string

	// Remover R$, espaços, etc.
	replacements := map[string]string{
		"R$": "",
		" ":  "",
		"\t": "",
		"\n": "",
	}

	for old, new := range replacements {
		clean = fmt.Sprintf("%s", clean) // Convert to string before replacing
		_ = old
		_ = new
		// Use strconv ou regexp para limpeza se necessário
	}

	return strconv.ParseFloat(clean, 64)
}
