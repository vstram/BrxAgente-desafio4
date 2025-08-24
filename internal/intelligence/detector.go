package intelligence

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

// AnomalyDetector é o detector principal de anomalias
type AnomalyDetector struct {
	config               *AnalysisConfig
	rules                []AnomalyRule
	valueDetector        *ValueAnomalyDetector
	temporalDetector     *TemporalAnomalyDetector
	relationshipDetector *RelationshipAnomalyDetector
	logger               *log.Logger
	mutex                sync.RWMutex
}

// NewAnomalyDetector cria novo detector de anomalias
func NewAnomalyDetector(config *AnalysisConfig) *AnomalyDetector {
	if config == nil {
		config = DefaultAnalysisConfig()
	}

	detector := &AnomalyDetector{
		config: config,
		rules:  make([]AnomalyRule, 0),
		logger: log.Default(),
	}

	// Inicializar detectores específicos
	detector.valueDetector = NewValueAnomalyDetector(config)
	detector.temporalDetector = NewTemporalAnomalyDetector(config)
	detector.relationshipDetector = NewRelationshipAnomalyDetector(config)

	// Registrar regras padrão
	detector.registerDefaultRules()

	return detector
}

// registerDefaultRules registra regras padrão de detecção
func (d *AnomalyDetector) registerDefaultRules() {
	// Regras de valor
	if d.config.EnableValueDetection {
		d.AddRule(AnomalyRule{
			Name:        "vr_outlier",
			Description: "Detecta valores de VR muito fora do padrão",
			Type:        AnomalyTypeValue,
			Severity:    SeverityMedium,
			Enabled:     true,
			Detector:    d.valueDetector.DetectVROutliers,
		})

		d.AddRule(AnomalyRule{
			Name:        "vr_zero_or_negative",
			Description: "Detecta valores de VR zero ou negativos",
			Type:        AnomalyTypeValue,
			Severity:    SeverityHigh,
			Enabled:     true,
			Detector:    d.valueDetector.DetectZeroOrNegativeVR,
		})
	}

	// Regras temporais
	if d.config.EnableTemporalDetection {
		d.AddRule(AnomalyRule{
			Name:        "future_dates",
			Description: "Detecta datas futuras inválidas",
			Type:        AnomalyTypeTemporal,
			Severity:    SeverityCritical,
			Enabled:     true,
			Detector:    d.temporalDetector.DetectFutureDates,
		})

		d.AddRule(AnomalyRule{
			Name:        "invalid_date_sequence",
			Description: "Detecta sequências de datas inválidas",
			Type:        AnomalyTypeTemporal,
			Severity:    SeverityHigh,
			Enabled:     true,
			Detector:    d.temporalDetector.DetectInvalidDateSequences,
		})
	}

	// Regras de relacionamento
	if d.config.EnableRelationshipDetection {
		d.AddRule(AnomalyRule{
			Name:        "duplicate_matricula",
			Description: "Detecta matrículas duplicadas",
			Type:        AnomalyTypeRelationship,
			Severity:    SeverityCritical,
			Enabled:     true,
			Detector:    d.relationshipDetector.DetectDuplicateMatriculas,
		})
	}
}

// AddRule adiciona nova regra de detecção
func (d *AnomalyDetector) AddRule(rule AnomalyRule) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.rules = append(d.rules, rule)
}

// RemoveRule remove regra de detecção
func (d *AnomalyDetector) RemoveRule(ruleName string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for i, rule := range d.rules {
		if rule.Name == ruleName {
			d.rules = append(d.rules[:i], d.rules[i+1:]...)
			break
		}
	}
}

// DetectAnomalies executa detecção de anomalias
func (d *AnomalyDetector) DetectAnomalies(ctx *AnalysisContext) (*AnomalyReport, error) {
	startTime := time.Now()

	d.logger.Printf("Iniciando detecção de anomalias com %d regras", len(d.rules))

	// Inicializar relatório
	report := &AnomalyReport{
		GeneratedAt:         startTime,
		TotalRecords:        len(ctx.Colaboradores),
		Anomalies:           make([]Anomaly, 0),
		AnomaliesByType:     make(map[AnomalyType]int),
		AnomaliesBySeverity: make(map[string]int),
	}

	// Canal para coletar anomalias de detectores paralelos
	anomalyChan := make(chan []Anomaly, len(d.rules))
	var wg sync.WaitGroup

	// Executar regras de detecção
	if d.config.EnableParallelProcessing {
		d.executeRulesParallel(ctx, &wg, anomalyChan)
	} else {
		d.executeRulesSequential(ctx, anomalyChan)
	}

	// Aguardar conclusão e coletar resultados
	go func() {
		wg.Wait()
		close(anomalyChan)
	}()

	for anomalies := range anomalyChan {
		for _, anomaly := range anomalies {
			if anomaly.Confidence >= d.config.ConfidenceThreshold {
				report.Anomalies = append(report.Anomalies, anomaly)
				report.AnomaliesByType[anomaly.Type]++
				report.AnomaliesBySeverity[anomaly.Severity.String()]++
			}
		}
	}

	// Finalizar relatório
	report.TotalAnomalies = len(report.Anomalies)
	report.Summary = d.generateSummary(report)

	duration := time.Since(startTime)
	d.logger.Printf("Detecção de anomalias concluída em %v - %d anomalias encontradas",
		duration, report.TotalAnomalies)

	return report, nil
}

// executeRulesParallel executa regras em paralelo
func (d *AnomalyDetector) executeRulesParallel(ctx *AnalysisContext, wg *sync.WaitGroup, anomalyChan chan<- []Anomaly) {
	semaphore := make(chan struct{}, d.config.MaxConcurrentDetectors)

	for _, rule := range d.rules {
		if !rule.Enabled {
			continue
		}

		wg.Add(1)
		go func(r AnomalyRule) {
			defer wg.Done()

			// Controlar concorrência
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			d.executeRule(r, ctx, anomalyChan)
		}(rule)
	}
}

// executeRulesSequential executa regras sequencialmente
func (d *AnomalyDetector) executeRulesSequential(ctx *AnalysisContext, anomalyChan chan<- []Anomaly) {
	for _, rule := range d.rules {
		if rule.Enabled {
			d.executeRule(rule, ctx, anomalyChan)
		}
	}
	close(anomalyChan)
}

// executeRule executa uma regra específica
func (d *AnomalyDetector) executeRule(rule AnomalyRule, ctx *AnalysisContext, anomalyChan chan<- []Anomaly) {
	defer func() {
		if r := recover(); r != nil {
			d.logger.Printf("Erro na execução da regra %s: %v", rule.Name, r)
			ctx.AddError(fmt.Errorf("regra %s falhou: %v", rule.Name, r))
		}
	}()

	anomalies := rule.Detector(ctx)
	anomalyChan <- anomalies
}

// generateSummary gera resumo executivo do relatório
func (d *AnomalyDetector) generateSummary(report *AnomalyReport) AnomalySummary {
	summary := AnomalySummary{
		RecommendedActions:    make([]string, 0),
		TopValueIssues:        make([]string, 0),
		TopTemporalIssues:     make([]string, 0),
		TopRelationshipIssues: make([]string, 0),
	}

	// Calcular score geral (0-100, onde 100 = sem problemas)
	if report.TotalRecords == 0 {
		summary.OverallScore = 100.0
	} else {
		// Penalizar baseado na severidade das anomalias
		penalty := 0.0
		for _, anomaly := range report.Anomalies {
			switch anomaly.Severity {
			case SeverityCritical:
				penalty += 10.0
			case SeverityHigh:
				penalty += 5.0
			case SeverityMedium:
				penalty += 2.0
			case SeverityLow:
				penalty += 1.0
			default:
				penalty += 0.5
			}
		}

		// Normalizar penalty pelo número de registros
		normalizedPenalty := (penalty / float64(report.TotalRecords)) * 100
		summary.OverallScore = math.Max(0, 100-normalizedPenalty)
	}

	// Determinar nível de risco
	criticalCount := report.AnomaliesBySeverity["critical"]
	highCount := report.AnomaliesBySeverity["high"]

	summary.CriticalIssues = criticalCount

	switch {
	case criticalCount > 0:
		summary.RiskLevel = "critical"
	case highCount > 0:
		summary.RiskLevel = "high"
	case summary.OverallScore < 80:
		summary.RiskLevel = "medium"
	default:
		summary.RiskLevel = "low"
	}

	// Gerar recomendações
	summary.RecommendedActions = d.generateRecommendations(report)

	// Agrupar problemas por tipo
	summary.TopValueIssues = d.getTopIssuesByType(report, AnomalyTypeValue, 5)
	summary.TopTemporalIssues = d.getTopIssuesByType(report, AnomalyTypeTemporal, 5)
	summary.TopRelationshipIssues = d.getTopIssuesByType(report, AnomalyTypeRelationship, 5)

	return summary
}

// generateRecommendations gera recomendações baseadas nas anomalias
func (d *AnomalyDetector) generateRecommendations(report *AnomalyReport) []string {
	recommendations := make([]string, 0)

	criticalCount := report.AnomaliesBySeverity["critical"]
	highCount := report.AnomaliesBySeverity["high"]

	if criticalCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Resolver imediatamente %d problemas críticos antes de prosseguir", criticalCount))
	}

	if highCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Revisar e corrigir %d problemas de alta severidade", highCount))
	}

	// Recomendações específicas por tipo
	if count := report.AnomaliesByType[AnomalyTypeValue]; count > 0 {
		recommendations = append(recommendations, "Verificar planilhas de valores de VR e fórmulas de cálculo")
	}

	if count := report.AnomaliesByType[AnomalyTypeTemporal]; count > 0 {
		recommendations = append(recommendations, "Validar datas de admissão, desligamento e períodos")
	}

	if count := report.AnomaliesByType[AnomalyTypeRelationship]; count > 0 {
		recommendations = append(recommendations, "Eliminar duplicatas e inconsistências entre planilhas")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Dados estão dentro dos padrões esperados")
	}

	return recommendations
}

// getTopIssuesByType retorna os principais problemas de um tipo específico
func (d *AnomalyDetector) getTopIssuesByType(report *AnomalyReport, anomalyType AnomalyType, limit int) []string {
	var anomalies []Anomaly

	// Filtrar por tipo
	for _, anomaly := range report.Anomalies {
		if anomaly.Type == anomalyType {
			anomalies = append(anomalies, anomaly)
		}
	}

	// Ordenar por severidade e confidence
	sort.Slice(anomalies, func(i, j int) bool {
		if anomalies[i].Severity != anomalies[j].Severity {
			return anomalies[i].Severity > anomalies[j].Severity
		}
		return anomalies[i].Confidence > anomalies[j].Confidence
	})

	// Extrair títulos dos top problemas
	issues := make([]string, 0)
	for i, anomaly := range anomalies {
		if i >= limit {
			break
		}
		issues = append(issues, anomaly.Title)
	}

	return issues
}

// GetStats retorna estatísticas do detector
func (d *AnomalyDetector) GetStats() map[string]interface{} {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	enabledRules := 0
	rulesByType := make(map[AnomalyType]int)

	for _, rule := range d.rules {
		if rule.Enabled {
			enabledRules++
			rulesByType[rule.Type]++
		}
	}

	return map[string]interface{}{
		"total_rules":              len(d.rules),
		"enabled_rules":            enabledRules,
		"rules_by_type":            rulesByType,
		"parallel_processing":      d.config.EnableParallelProcessing,
		"max_concurrent_detectors": d.config.MaxConcurrentDetectors,
		"confidence_threshold":     d.config.ConfidenceThreshold,
		"outlier_threshold":        d.config.OutlierThreshold,
	}
}

// Helper function para converter interface{} para Colaborador
func convertToColaborador(data interface{}) (*modelo.Colaborador, error) {
	if colaborador, ok := data.(*modelo.Colaborador); ok {
		return colaborador, nil
	}

	if colaboradorMap, ok := data.(map[string]interface{}); ok {
		// Converter map para struct (implementação simplificada)
		colaborador := &modelo.Colaborador{}

		if matricula, exists := colaboradorMap["matricula"]; exists {
			if str, ok := matricula.(string); ok {
				colaborador.Matricula = str
			}
		}

		if sindicato, exists := colaboradorMap["sindicato"]; exists {
			if str, ok := sindicato.(string); ok {
				colaborador.Sindicato = str
			}
		}

		// Adicionar outros campos conforme necessário
		return colaborador, nil
	}

	return nil, fmt.Errorf("não foi possível converter dados para Colaborador")
}

// Helper function para extrair valor numérico
func extractNumericValue(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		// Tentar converter string para número
		cleanValue := strings.ReplaceAll(v, ",", ".")
		cleanValue = strings.ReplaceAll(cleanValue, "R$", "")
		cleanValue = strings.TrimSpace(cleanValue)
		return strconv.ParseFloat(cleanValue, 64)
	default:
		return 0, fmt.Errorf("valor não numérico: %T", value)
	}
}
