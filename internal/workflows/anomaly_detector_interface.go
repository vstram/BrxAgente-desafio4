package workflows

import "time"

// AnomalyDetectorInterface define a interface para detectores de anomalias
// Esta interface permite que o pacote workflows use detectores de anomalias
// sem importar diretamente o pacote intelligence
type AnomalyDetectorInterface interface {
	// AnalyzeData executa análise de anomalias nos dados fornecidos
	AnalyzeData(colaboradores map[string]interface{}, params map[string]interface{}) (AnomalyReportInterface, error)
	
	// GetStats retorna estatísticas do detector
	GetStats() map[string]interface{}
}

// AnomalyReportInterface representa um relatório de anomalias (interface)
// Esta é uma versão simplificada que evita dependências circulares
type AnomalyReportInterface interface {
	// GetTotalAnomalies retorna o número total de anomalias
	GetTotalAnomalies() int
	
	// GetOverallScore retorna o score geral da análise
	GetOverallScore() float64
	
	// GetRiskLevel retorna o nível de risco
	GetRiskLevel() string
	
	// GetSummary retorna um resumo das anomalias encontradas
	GetSummary() map[string]interface{}
	
	// FormatForHuman retorna uma representação legível para humanos
	FormatForHuman() string
}

// AnomalyDetectionStep implementa um step de workflow para detecção de anomalias
type AnomalyDetectionStep struct {
	*BaseStep
	detector AnomalyDetectorInterface
}

// NewAnomalyDetectionStep cria novo step de detecção de anomalias
func NewAnomalyDetectionStep(detector AnomalyDetectorInterface) *AnomalyDetectionStep {
	return &AnomalyDetectionStep{
		BaseStep: NewBaseStep(
			"anomaly-detection",
			"Detecta anomalias nos dados de colaboradores",
			30*time.Second,
		),
		detector: detector,
	}
}

// Execute executa o step de detecção de anomalias  
func (s *AnomalyDetectionStep) Execute(ctx *WorkflowContext) error {
	if ctx.Logger != nil {
		ctx.Logger.Info("Executando detecção de anomalias")
	}
	
	// Extrair dados dos colaboradores do contexto
	colaboradores, exists := ctx.Get("colaboradores")
	if !exists {
		return NewWorkflowError("anomaly-detection", "", "dados de colaboradores não encontrados no contexto", nil)
	}
	
	colaboradoresMap, ok := colaboradores.(map[string]interface{})
	if !ok {
		return NewWorkflowError("anomaly-detection", "", "formato de dados de colaboradores inválido", nil)
	}
	
	// Extrair parâmetros de análise
	params := make(map[string]interface{})
	if p, exists := ctx.Get("analysis_params"); exists {
		if pMap, ok := p.(map[string]interface{}); ok {
			params = pMap
		}
	}
	
	// Executar análise
	report, err := s.detector.AnalyzeData(colaboradoresMap, params)
	if err != nil {
		return NewWorkflowError("anomaly-detection", "", "erro na análise de anomalias", err)
	}
	
	// Armazenar resultado no contexto
	ctx.Set("anomaly_report", report)
	
	// Log de resultados
	if ctx.Logger != nil {
		ctx.Logger.Info("Detecção de anomalias concluída: %d anomalias encontradas (score: %.1f)",
			report.GetTotalAnomalies(), report.GetOverallScore())
	}
	
	return nil
}

// CanSkip verifica se o step pode ser pulado
func (s *AnomalyDetectionStep) CanSkip(ctx *WorkflowContext) bool {
	// Não pular se não há dados de colaboradores
	_, exists := ctx.Get("colaboradores")
	return !exists
}

// Rollback desfaz as ações do step
func (s *AnomalyDetectionStep) Rollback(ctx *WorkflowContext) error {
	if ctx.Logger != nil {
		ctx.Logger.Info("Rollback do step de detecção de anomalias")
	}
	ctx.Set("anomaly_report", nil)
	return nil
}

