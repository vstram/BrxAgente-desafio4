package intelligence

// AnalyzerAdapter implementa as interfaces definidas em workflows
// permitindo que o Analyzer seja usado sem dependência circular
type AnalyzerAdapter struct {
	analyzer *Analyzer
}

// NewAnalyzerAdapter cria um novo adapter para o Analyzer
func NewAnalyzerAdapter(analyzer *Analyzer) *AnalyzerAdapter {
	return &AnalyzerAdapter{
		analyzer: analyzer,
	}
}

// AnalyzeData implementa a interface AnomalyDetector do workflows
func (a *AnalyzerAdapter) AnalyzeData(colaboradores map[string]interface{}, params map[string]interface{}) (AnomalyReportInterface, error) {
	report, err := a.analyzer.AnalyzeData(colaboradores, params)
	if err != nil {
		return nil, err
	}

	// Retornar um wrapper que implementa a interface
	return &AnomalyReportWrapper{report: report}, nil
}

// GetStats implementa a interface AnomalyDetector do workflows
func (a *AnalyzerAdapter) GetStats() map[string]interface{} {
	return a.analyzer.GetStats()
}

// AnomalyReportInterface define a interface que o workflow espera
type AnomalyReportInterface interface {
	GetTotalAnomalies() int
	GetOverallScore() float64
	GetRiskLevel() string
	GetSummary() map[string]interface{}
	FormatForHuman() string
}

// AnomalyReportWrapper envolve AnomalyReport para implementar a interface
type AnomalyReportWrapper struct {
	report *AnomalyReport
}

// GetTotalAnomalies implementa a interface
func (w *AnomalyReportWrapper) GetTotalAnomalies() int {
	if w.report == nil {
		return 0
	}
	return w.report.TotalAnomalies
}

// GetOverallScore implementa a interface
func (w *AnomalyReportWrapper) GetOverallScore() float64 {
	if w.report == nil {
		return 0.0
	}
	return w.report.Summary.OverallScore
}

// GetRiskLevel implementa a interface
func (w *AnomalyReportWrapper) GetRiskLevel() string {
	if w.report == nil {
		return "unknown"
	}
	return w.report.Summary.RiskLevel
}

// GetSummary implementa a interface
func (w *AnomalyReportWrapper) GetSummary() map[string]interface{} {
	if w.report == nil {
		return make(map[string]interface{})
	}

	summary := make(map[string]interface{})
	summary["generated_at"] = w.report.GeneratedAt
	summary["total_records"] = w.report.TotalRecords
	summary["total_anomalies"] = w.report.TotalAnomalies
	summary["anomalies_by_type"] = w.report.AnomaliesByType
	summary["anomalies_by_severity"] = w.report.AnomaliesBySeverity

	// Summary sempre existe pois é uma struct, não ponteiro
	summary["overall_score"] = w.report.Summary.OverallScore
	summary["risk_level"] = w.report.Summary.RiskLevel
	summary["critical_issues"] = w.report.Summary.CriticalIssues
	summary["recommended_actions"] = w.report.Summary.RecommendedActions
	summary["top_value_issues"] = w.report.Summary.TopValueIssues
	summary["top_temporal_issues"] = w.report.Summary.TopTemporalIssues
	summary["top_relationship_issues"] = w.report.Summary.TopRelationshipIssues

	return summary
}

// FormatForHuman implementa a interface
func (w *AnomalyReportWrapper) FormatForHuman() string {
	if w.report == nil {
		return "Nenhum relatório de anomalias disponível"
	}
	return FormatAnomalyReportForHuman(w.report)
}

// GetOriginalReport retorna o relatório original (método adicional)
func (w *AnomalyReportWrapper) GetOriginalReport() *AnomalyReport {
	return w.report
}
