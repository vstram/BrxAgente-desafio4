package workflows

import (
	"fmt"
	"time"
)

// SimpleValidationWorkflow implementa um workflow básico de validação para teste
type SimpleValidationWorkflow struct {
	*BaseWorkflow
}

// NewSimpleValidationWorkflow cria um novo workflow de validação simples
func NewSimpleValidationWorkflow() *SimpleValidationWorkflow {
	steps := []WorkflowStep{
		NewLoadDataStep(),
		NewValidateDataStep(),
		NewReportResultsStep(),
	}
	
	baseWorkflow := NewBaseWorkflow(
		"simple-validation",
		"Workflow de validação simples para demonstrar o orquestrador",
		steps,
	)
	
	return &SimpleValidationWorkflow{
		BaseWorkflow: baseWorkflow,
	}
}

// Execute executa o workflow através do orquestrador (sobrescreve BaseWorkflow.Execute)
func (w *SimpleValidationWorkflow) Execute(ctx *WorkflowContext) error {
	ctx.Logger.Info("Executando workflow de validação simples")
	
	// Esta implementação seria chamada pelo orquestrador
	// Em uma implementação real, isto seria tratado pelo orchestrator.executeWorkflowInternal
	for _, step := range w.Steps() {
		if ctx.IsCanceled() {
			return fmt.Errorf("workflow cancelado")
		}
		
		if err := step.Execute(ctx); err != nil {
			return fmt.Errorf("erro no step '%s': %w", step.Name(), err)
		}
	}
	
	return nil
}

// LoadDataStep - Step para carregar dados
type LoadDataStep struct {
	*BaseStep
}

func NewLoadDataStep() *LoadDataStep {
	return &LoadDataStep{
		BaseStep: NewBaseStep(
			"load-data",
			"Carrega dados para validação",
			5*time.Second,
		),
	}
}

func (s *LoadDataStep) Execute(ctx *WorkflowContext) error {
	ctx.Logger.Debug("Executando step: %s", s.Name())
	
	// Simular carregamento de dados
	filePath, hasFile := ctx.GetParameterString("file_path")
	if !hasFile {
		filePath = "files/ATIVOS.xlsx" // Padrão
	}
	
	ctx.Logger.Info("Carregando dados de: %s", filePath)
	
	// Simular dados carregados
	loadedData := map[string]interface{}{
		"source_file": filePath,
		"records_count": 150,
		"columns": []string{"MATRICULA", "SINDICATO", "DATA_ADMISSAO"},
		"load_time": time.Now(),
	}
	
	ctx.SetResult(s.Name(), loadedData)
	ctx.Set("data_loaded", true)
	
	return nil
}

func (s *LoadDataStep) CanSkip(ctx *WorkflowContext) bool {
	// Pular se dados já foram carregados
	return ctx.Has("data_loaded")
}

func (s *LoadDataStep) Rollback(ctx *WorkflowContext) error {
	ctx.Logger.Debug("Fazendo rollback de %s", s.Name())
	ctx.Delete("data_loaded")
	return nil
}

// ValidateDataStep - Step para validar dados
type ValidateDataStep struct {
	*BaseStep
}

func NewValidateDataStep() *ValidateDataStep {
	return &ValidateDataStep{
		BaseStep: NewBaseStep(
			"validate-data",
			"Valida consistência dos dados carregados",
			3*time.Second,
		),
	}
}

func (s *ValidateDataStep) Execute(ctx *WorkflowContext) error {
	ctx.Logger.Debug("Executando step: %s", s.Name())
	
	// Verificar se dados foram carregados
	if !ctx.Has("data_loaded") {
		return fmt.Errorf("dados não foram carregados antes da validação")
	}
	
	// Recuperar dados do step anterior
	loadResult, exists := ctx.GetResult("load-data")
	if !exists {
		return fmt.Errorf("resultado do carregamento não encontrado")
	}
	
	loadData, ok := loadResult.(map[string]interface{})
	if !ok {
		return fmt.Errorf("formato de dados inválido")
	}
	
	ctx.Logger.Info("Validando dados de: %v", loadData["source_file"])
	
	// Simular validações
	validationResult := map[string]interface{}{
		"total_records": loadData["records_count"],
		"valid_records": 145,
		"invalid_records": 5,
		"errors": []string{
			"Matrícula duplicada: 12345",
			"Data de admissão inválida: colaborador 67890",
			"Sindicato não reconhecido: colaborador 11111",
		},
		"validation_time": time.Now(),
		"is_valid": false, // Há erros
	}
	
	ctx.SetResult(s.Name(), validationResult)
	ctx.Set("validation_completed", true)
	
	return nil
}

func (s *ValidateDataStep) CanSkip(ctx *WorkflowContext) bool {
	// Pular se validação já foi feita
	return ctx.Has("validation_completed")
}

func (s *ValidateDataStep) Rollback(ctx *WorkflowContext) error {
	ctx.Logger.Debug("Fazendo rollback de %s", s.Name())
	ctx.Delete("validation_completed")
	return nil
}

// ReportResultsStep - Step para gerar relatório
type ReportResultsStep struct {
	*BaseStep
}

func NewReportResultsStep() *ReportResultsStep {
	return &ReportResultsStep{
		BaseStep: NewBaseStep(
			"report-results",
			"Gera relatório dos resultados da validação",
			2*time.Second,
		),
	}
}

func (s *ReportResultsStep) Execute(ctx *WorkflowContext) error {
	ctx.Logger.Debug("Executando step: %s", s.Name())
	
	// Verificar se validação foi feita
	if !ctx.Has("validation_completed") {
		return fmt.Errorf("validação não foi completada antes do relatório")
	}
	
	// Recuperar resultados da validação
	validationResult, exists := ctx.GetResult("validate-data")
	if !exists {
		return fmt.Errorf("resultado da validação não encontrado")
	}
	
	validationData, ok := validationResult.(map[string]interface{})
	if !ok {
		return fmt.Errorf("formato de resultado de validação inválido")
	}
	
	// Gerar relatório
	report := map[string]interface{}{
		"workflow_name": ctx.WorkflowName,
		"execution_id": ctx.ExecutionID,
		"total_duration": ctx.Elapsed(),
		"summary": fmt.Sprintf(
			"Validação completada: %d registros válidos, %d inválidos de %d total",
			validationData["valid_records"], 
			validationData["invalid_records"], 
			validationData["total_records"],
		),
		"details": validationData,
		"report_time": time.Now(),
	}
	
	ctx.SetResult(s.Name(), report)
	ctx.Logger.Info("Relatório gerado: %s", report["summary"])
	
	return nil
}

func (s *ReportResultsStep) CanSkip(ctx *WorkflowContext) bool {
	// Nunca pular o relatório final
	return false
}

func (s *ReportResultsStep) Rollback(ctx *WorkflowContext) error {
	ctx.Logger.Debug("Fazendo rollback de %s", s.Name())
	// Não há muito o que desfazer em um relatório
	return nil
}