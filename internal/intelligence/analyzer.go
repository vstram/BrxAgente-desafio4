package intelligence

import (
	"fmt"
	"log"
	"time"

	"BrxAgente-desafio4/internal/workflows"
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
		"outlier_threshold":     a.config.OutlierThreshold,
		"confidence_threshold":  a.config.ConfidenceThreshold,
		"vr_min_value":         a.config.VRMinValue,
		"vr_max_value":         a.config.VRMaxValue,
		"max_work_days":        a.config.MaxWorkDaysPerMonth,
		"parallel_processing":   a.config.EnableParallelProcessing,
	}
	
	return stats
}

// AnomalyDetectionStep implementa um step de workflow para detecção de anomalias
type AnomalyDetectionStep struct {
	*workflows.BaseStep
	analyzer *Analyzer
}

// NewAnomalyDetectionStep cria novo step de detecção de anomalias
func NewAnomalyDetectionStep(analyzer *Analyzer) *AnomalyDetectionStep {
	return &AnomalyDetectionStep{
		BaseStep: workflows.NewBaseStep(
			"anomaly-detection",
			"Detecta anomalias nos dados de colaboradores",
			30*time.Second,
		),
		analyzer: analyzer,
	}
}

// Execute executa o step de detecção de anomalias
func (s *AnomalyDetectionStep) Execute(ctx *workflows.WorkflowContext) error {
	ctx.Logger.Info("Executando detecção de anomalias")
	
	// Extrair dados dos colaboradores do contexto
	colaboradores, exists := ctx.Get("colaboradores")
	if !exists {
		return fmt.Errorf("dados de colaboradores não encontrados no contexto")
	}
	
	colaboradoresMap, ok := colaboradores.(map[string]interface{})
	if !ok {
		return fmt.Errorf("formato de dados de colaboradores inválido")
	}
	
	// Extrair parâmetros de análise
	params := make(map[string]interface{})
	if p, exists := ctx.Get("analysis_params"); exists {
		if pMap, ok := p.(map[string]interface{}); ok {
			params = pMap
		}
	}
	
	// Executar análise
	report, err := s.analyzer.AnalyzeData(colaboradoresMap, params)
	if err != nil {
		return fmt.Errorf("erro na análise de anomalias: %w", err)
	}
	
	// Armazenar resultado no contexto
	ctx.SetResult(s.Name(), report)
	ctx.Set("anomaly_report", report)
	
	// Log de resultados
	ctx.Logger.Info("Detecção de anomalias concluída: %d anomalias encontradas (score: %.1f)",
		report.TotalAnomalies, report.Summary.OverallScore)
	
	// Se há anomalias críticas, adicionar warning
	if report.Summary.CriticalIssues > 0 {
		ctx.Logger.Warn("Encontradas %d anomalias críticas que requerem atenção imediata",
			report.Summary.CriticalIssues)
	}
	
	return nil
}

// CanSkip verifica se o step pode ser pulado
func (s *AnomalyDetectionStep) CanSkip(ctx *workflows.WorkflowContext) bool {
	// Não pular se não há dados de colaboradores
	_, exists := ctx.Get("colaboradores")
	return !exists
}

// Rollback desfaz as ações do step (não aplicável para detecção de anomalias)
func (s *AnomalyDetectionStep) Rollback(ctx *workflows.WorkflowContext) error {
	ctx.Logger.Debug("Rollback do step de detecção de anomalias (removendo resultados)")
	ctx.Delete("anomaly_report")
	return nil
}

// ValidatedVRWorkflow implementa workflow com detecção de anomalias
type ValidatedVRWorkflow struct {
	*workflows.BaseWorkflow
	analyzer *Analyzer
}

// NewValidatedVRWorkflow cria workflow com validação de anomalias
func NewValidatedVRWorkflow(analyzer *Analyzer) *ValidatedVRWorkflow {
	if analyzer == nil {
		analyzer = NewAnalyzer(nil)
	}
	
	steps := []workflows.WorkflowStep{
		workflows.NewLoadDataStep(),                    // Carregar dados
		NewAnomalyDetectionStep(analyzer),             // Detectar anomalias
		workflows.NewValidateDataStep(),               // Validação adicional
		NewAnomalyReportStep(),                        // Gerar relatório de anomalias
		workflows.NewReportResultsStep(),              // Relatório final
	}
	
	baseWorkflow := workflows.NewBaseWorkflow(
		"validated-vr-processing",
		"Workflow de processamento de VR com detecção de anomalias",
		steps,
	)
	
	return &ValidatedVRWorkflow{
		BaseWorkflow: baseWorkflow,
		analyzer:     analyzer,
	}
}

// AnomalyReportStep gera relatório detalhado de anomalias
type AnomalyReportStep struct {
	*workflows.BaseStep
}

// NewAnomalyReportStep cria step de relatório de anomalias
func NewAnomalyReportStep() *AnomalyReportStep {
	return &AnomalyReportStep{
		BaseStep: workflows.NewBaseStep(
			"anomaly-report",
			"Gera relatório detalhado das anomalias encontradas",
			10*time.Second,
		),
	}
}

// Execute executa o step de relatório de anomalias
func (s *AnomalyReportStep) Execute(ctx *workflows.WorkflowContext) error {
	ctx.Logger.Info("Gerando relatório de anomalias")
	
	// Recuperar relatório de anomalias
	anomalyData, exists := ctx.Get("anomaly_report")
	if !exists {
		ctx.Logger.Warn("Relatório de anomalias não encontrado, pulando geração de relatório")
		return nil
	}
	
	report, ok := anomalyData.(*AnomalyReport)
	if !ok {
		return fmt.Errorf("formato de relatório de anomalias inválido")
	}
	
	// Gerar relatório executivo
	executiveReport := map[string]interface{}{
		"workflow_name":    ctx.WorkflowName,
		"execution_id":     ctx.ExecutionID,
		"analysis_date":    report.GeneratedAt,
		"total_records":    report.TotalRecords,
		"total_anomalies":  report.TotalAnomalies,
		"overall_score":    report.Summary.OverallScore,
		"risk_level":       report.Summary.RiskLevel,
		"critical_issues":  report.Summary.CriticalIssues,
		"recommendations":  report.Summary.RecommendedActions,
		"anomalies_by_type": report.AnomaliesByType,
		"anomalies_by_severity": report.AnomaliesBySeverity,
	}
	
	// Adicionar top problemas por categoria
	if len(report.Summary.TopValueIssues) > 0 {
		executiveReport["top_value_issues"] = report.Summary.TopValueIssues
	}
	
	if len(report.Summary.TopTemporalIssues) > 0 {
		executiveReport["top_temporal_issues"] = report.Summary.TopTemporalIssues
	}
	
	if len(report.Summary.TopRelationshipIssues) > 0 {
		executiveReport["top_relationship_issues"] = report.Summary.TopRelationshipIssues
	}
	
	// Armazenar resultado
	ctx.SetResult(s.Name(), executiveReport)
	ctx.Set("executive_report", executiveReport)
	
	ctx.Logger.Info("Relatório executivo de anomalias gerado com score %.1f (%s)",
		report.Summary.OverallScore, report.Summary.RiskLevel)
	
	return nil
}

// CanSkip verifica se o step pode ser pulado
func (s *AnomalyReportStep) CanSkip(ctx *workflows.WorkflowContext) bool {
	// Pular se não há relatório de anomalias
	_, exists := ctx.Get("anomaly_report")
	return !exists
}

// Rollback desfaz as ações do step
func (s *AnomalyReportStep) Rollback(ctx *workflows.WorkflowContext) error {
	ctx.Logger.Debug("Rollback do step de relatório de anomalias")
	ctx.Delete("executive_report")
	return nil
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
	
	// Resumo por severidade
	output += fmt.Sprintf("📊 DISTRIBUIÇÃO POR SEVERIDADE\n")
	output += fmt.Sprintf("─────────────────────────────\n")
	for severity, count := range report.AnomaliesBySeverity {
		if count > 0 {
			emoji := getSeverityEmoji(severity)
			output += fmt.Sprintf("%s %s: %d\n", emoji, severity, count)
		}
	}
	
	// Resumo por tipo
	output += fmt.Sprintf("\n📋 DISTRIBUIÇÃO POR TIPO\n")
	output += fmt.Sprintf("─────────────────────────\n")
	for anomalyType, count := range report.AnomaliesByType {
		if count > 0 {
			emoji := getTypeEmoji(string(anomalyType))
			output += fmt.Sprintf("%s %s: %d\n", emoji, anomalyType, count)
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