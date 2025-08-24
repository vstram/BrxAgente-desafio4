package intelligence

import (
	"fmt"
	"log"
)

// Analyzer é o analisador principal que coordena a detecção de anomalias
type Analyzer struct {
	detector *AnomalyDetector
	config   *AnalysisConfig
	logger   *log.Logger
}

// NewAnalyzer cria novo analisador
func NewAnalyzer(config *AnalysisConfig) *Analyzer {
	if config == nil {
		config = DefaultAnalysisConfig()
	}

	return &Analyzer{
		detector: NewAnomalyDetector(config),
		config:   config,
		logger:   log.Default(),
	}
}

// AnalyzeData executa análise completa dos dados
func (a *Analyzer) AnalyzeData(colaboradores map[string]interface{}, params map[string]interface{}) (*AnomalyReport, error) {
	a.logger.Printf("Iniciando análise de anomalias para %d registros", len(colaboradores))

	// Criar contexto de análise
	ctx := NewAnalysisContext(colaboradores, a.config)
	if params != nil {
		ctx.Parameters = params
	}

	// Executar detecção de anomalias
	report, err := a.detector.DetectAnomalies(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro na detecção de anomalias: %w", err)
	}

	// Adicionar metadados do contexto ao relatório
	report = a.enrichReport(report, ctx)

	a.logger.Printf("Análise concluída: %d anomalias detectadas", report.TotalAnomalies)

	return report, nil
}

// enrichReport enriquece o relatório com informações do contexto
func (a *Analyzer) enrichReport(report *AnomalyReport, ctx *AnalysisContext) *AnomalyReport {
	// Adicionar informações de processamento
	if len(ctx.Errors) > 0 {
		report.Summary.RecommendedActions = append(report.Summary.RecommendedActions,
			fmt.Sprintf("Resolver %d erros de processamento encontrados", len(ctx.Errors)))
	}

	if len(ctx.Warnings) > 0 {
		report.Summary.RecommendedActions = append(report.Summary.RecommendedActions,
			fmt.Sprintf("Revisar %d avisos de processamento", len(ctx.Warnings)))
	}

	return report
}

// GetStats retorna estatísticas do analisador
func (a *Analyzer) GetStats() map[string]interface{} {
	stats := a.detector.GetStats()
	stats["analyzer_version"] = "1.0.0"
	stats["config"] = map[string]interface{}{
		"outlier_threshold":    a.config.OutlierThreshold,
		"confidence_threshold": a.config.ConfidenceThreshold,
		"vr_min_value":         a.config.VRMinValue,
		"vr_max_value":         a.config.VRMaxValue,
		"max_work_days":        a.config.MaxWorkDaysPerMonth,
		"parallel_processing":  a.config.EnableParallelProcessing,
	}

	return stats
}

// FormatAnomalyReportForHuman formata relatório para leitura humana
func FormatAnomalyReportForHuman(report *AnomalyReport) string {
	if report == nil {
		return "Nenhum relatório de anomalias disponível"
	}

	output := fmt.Sprintf("📊 RELATÓRIO DE ANOMALIAS\n")
	output += fmt.Sprintf("═══════════════════════════\n\n")
	output += fmt.Sprintf("📅 Data: %s\n", report.GeneratedAt.Format("02/01/2006 15:04"))
	output += fmt.Sprintf("📁 Registros analisados: %d\n", report.TotalRecords)
	output += fmt.Sprintf("🚨 Anomalias encontradas: %d\n", report.TotalAnomalies)
	output += fmt.Sprintf("📈 Score geral: %.1f/100\n", report.Summary.OverallScore)
	output += fmt.Sprintf("⚠️  Nível de risco: %s\n\n", report.Summary.RiskLevel)

	if report.TotalAnomalies == 0 {
		output += "✅ Nenhuma anomalia detectada - dados estão dentro dos padrões esperados.\n"
		return output
	}

	// Distribuição por severidade
	if len(report.AnomaliesBySeverity) > 0 {
		output += fmt.Sprintf("🔴 DISTRIBUIÇÃO POR SEVERIDADE\n")
		output += fmt.Sprintf("─────────────────────────\n")
		for severity, count := range report.AnomaliesBySeverity {
			if count > 0 {
				emoji := getSeverityEmoji(severity)
				output += fmt.Sprintf("%s %s: %d\n", emoji, severity, count)
			}
		}
		output += fmt.Sprintf("\n")
	}

	// Distribuição por tipo
	if len(report.AnomaliesByType) > 0 {
		output += fmt.Sprintf("📊 DISTRIBUIÇÃO POR TIPO\n")
		output += fmt.Sprintf("─────────────────────────\n")
		for anomalyType, count := range report.AnomaliesByType {
			if count > 0 {
				emoji := getTypeEmoji(string(anomalyType))
				output += fmt.Sprintf("%s %s: %d\n", emoji, anomalyType, count)
			}
		}
	}

	// Recomendações
	if len(report.Summary.RecommendedActions) > 0 {
		output += fmt.Sprintf("\n💡 RECOMENDAÇÕES\n")
		output += fmt.Sprintf("────────────────\n")
		for i, action := range report.Summary.RecommendedActions {
			output += fmt.Sprintf("%d. %s\n", i+1, action)
		}
	}

	return output
}

// Helper functions para formatação
func getSeverityEmoji(severity string) string {
	switch severity {
	case "critical":
		return "🔴"
	case "high":
		return "🟠"
	case "medium":
		return "🟡"
	case "low":
		return "🔵"
	default:
		return "⚪"
	}
}

func getTypeEmoji(anomalyType string) string {
	switch anomalyType {
	case "value":
		return "💰"
	case "temporal":
		return "📅"
	case "relationship":
		return "🔗"
	case "structural":
		return "🏗️"
	default:
		return "❓"
	}
}
