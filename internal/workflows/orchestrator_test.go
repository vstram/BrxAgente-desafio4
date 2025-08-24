package workflows

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

// TestOrchestrator testa funcionalidades básicas do orquestrador
func TestOrchestrator(t *testing.T) {
	t.Run("CreateOrchestrator", testCreateOrchestrator)
	t.Run("RegisterWorkflow", testRegisterWorkflow)
	t.Run("ExecuteWorkflow", testExecuteWorkflow)
	t.Run("ExecuteWorkflowAsync", testExecuteWorkflowAsync)
	t.Run("CancelExecution", testCancelExecution)
	t.Run("ConcurrentExecutions", testConcurrentExecutions)
	t.Run("WorkflowValidation", testWorkflowValidation)
	t.Run("ErrorHandling", testErrorHandling)
}

func testCreateOrchestrator(t *testing.T) {
	logger := &DefaultLogger{Logger: log.New(os.Stdout, "[TEST] ", log.LstdFlags)}
	config := OrchestratorConfig{
		MaxConcurrentWorkflows: 5,
		DefaultTimeout:         30 * time.Second,
		EnableRollback:         true,
		DetailedLogging:        true,
	}

	orchestrator := NewOrchestrator(logger, config)

	if orchestrator == nil {
		t.Fatal("Orquestrador não deveria ser nil")
	}

	stats := orchestrator.GetStats()
	if stats["registered_workflows"] != 0 {
		t.Errorf("Workflows registrados deveria ser 0, got %v", stats["registered_workflows"])
	}

	if stats["max_concurrent_workflows"] != 5 {
		t.Errorf("Max concurrent workflows deveria ser 5, got %v", stats["max_concurrent_workflows"])
	}
}

func testRegisterWorkflow(t *testing.T) {
	orchestrator := createTestOrchestrator()

	// Registrar workflow válido
	workflow := NewSimpleValidationWorkflow()
	err := orchestrator.RegisterWorkflow(workflow)
	if err != nil {
		t.Fatalf("Erro ao registrar workflow: %v", err)
	}

	// Verificar se foi registrado
	registeredWorkflow, err := orchestrator.GetWorkflow("simple-validation")
	if err != nil {
		t.Fatalf("Erro ao recuperar workflow registrado: %v", err)
	}

	if registeredWorkflow.Name() != "simple-validation" {
		t.Errorf("Nome do workflow incorreto: %s", registeredWorkflow.Name())
	}

	// Tentar registrar o mesmo workflow novamente (deve falhar)
	err = orchestrator.RegisterWorkflow(workflow)
	if err == nil {
		t.Error("Deveria falhar ao registrar workflow duplicado")
	}

	if !strings.Contains(err.Error(), "already registered") {
		t.Errorf("Erro deveria mencionar workflow já registrado: %v", err)
	}

	// Verificar lista de workflows
	workflows := orchestrator.ListWorkflows()
	if len(workflows) != 1 {
		t.Errorf("Deveria haver 1 workflow registrado, got %d", len(workflows))
	}

	if workflows[0] != "simple-validation" {
		t.Errorf("Nome do workflow na lista incorreto: %s", workflows[0])
	}
}

func testExecuteWorkflow(t *testing.T) {
	orchestrator := createTestOrchestrator()

	// Registrar workflow
	workflow := NewSimpleValidationWorkflow()
	err := orchestrator.RegisterWorkflow(workflow)
	if err != nil {
		t.Fatalf("Erro ao registrar workflow: %v", err)
	}

	// Executar workflow
	params := map[string]interface{}{
		"file_path": "test_file.xlsx",
	}

	result, err := orchestrator.ExecuteWorkflow("simple-validation", params)
	if err != nil {
		t.Fatalf("Erro ao executar workflow: %v", err)
	}

	// Verificar resultado
	if result == nil {
		t.Fatal("Resultado não deveria ser nil")
	}

	if result.WorkflowName != "simple-validation" {
		t.Errorf("Nome do workflow no resultado incorreto: %s", result.WorkflowName)
	}

	if result.Status != StatusCompleted {
		t.Errorf("Status deveria ser Completed, got %s", result.Status)
	}

	if result.CompletedSteps != 3 {
		t.Errorf("Deveria ter completado 3 steps, got %d", result.CompletedSteps)
	}

	if result.TotalSteps != 3 {
		t.Errorf("Total de steps deveria ser 3, got %d", result.TotalSteps)
	}

	// Verificar step results
	if len(result.StepResults) != 3 {
		t.Errorf("Deveria haver 3 step results, got %d", len(result.StepResults))
	}

	for i, stepResult := range result.StepResults {
		if stepResult.Status != "completed" {
			t.Errorf("Step %d deveria estar completed, got %s", i, stepResult.Status)
		}

		if stepResult.Error != nil {
			t.Errorf("Step %d não deveria ter erro: %v", i, stepResult.Error)
		}
	}
}

func testExecuteWorkflowAsync(t *testing.T) {
	orchestrator := createTestOrchestrator()

	// Registrar workflow
	workflow := NewSimpleValidationWorkflow()
	err := orchestrator.RegisterWorkflow(workflow)
	if err != nil {
		t.Fatalf("Erro ao registrar workflow: %v", err)
	}

	// Executar workflow assincronamente
	params := map[string]interface{}{
		"file_path": "async_test_file.xlsx",
	}

	executionID, err := orchestrator.ExecuteWorkflowAsync("simple-validation", params)
	if err != nil {
		t.Fatalf("Erro ao executar workflow async: %v", err)
	}

	if executionID == "" {
		t.Fatal("ExecutionID não deveria estar vazio")
	}

	// Aguardar um pouco para a execução assíncrona
	time.Sleep(100 * time.Millisecond)

	// Verificar se a execução ainda existe (pode ter terminado)
	activeExecutions := orchestrator.ListActiveExecutions()
	t.Logf("Execuções ativas: %v", activeExecutions)

	// Aguardar mais um pouco para garantir que terminou
	time.Sleep(1 * time.Second)

	// Verificar se terminou
	activeExecutions = orchestrator.ListActiveExecutions()
	if len(activeExecutions) > 0 {
		t.Logf("Ainda há execuções ativas: %v", activeExecutions)
	}
}

func testCancelExecution(t *testing.T) {
	orchestrator := createTestOrchestrator()

	// Registrar workflow
	workflow := NewSimpleValidationWorkflow()
	err := orchestrator.RegisterWorkflow(workflow)
	if err != nil {
		t.Fatalf("Erro ao registrar workflow: %v", err)
	}

	// Executar workflow assincronamente
	params := map[string]interface{}{
		"file_path": "cancel_test_file.xlsx",
	}

	executionID, err := orchestrator.ExecuteWorkflowAsync("simple-validation", params)
	if err != nil {
		t.Fatalf("Erro ao executar workflow async: %v", err)
	}

	// Cancelar imediatamente
	err = orchestrator.CancelExecution(executionID)
	if err != nil {
		t.Logf("Aviso: Erro ao cancelar execução (pode ter terminado): %v", err)
	}

	// Tentar cancelar execução inexistente
	err = orchestrator.CancelExecution("inexistente-id")
	if err == nil {
		t.Error("Deveria falhar ao cancelar execução inexistente")
	}
}

func testConcurrentExecutions(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de concorrência em modo short")
	}

	orchestrator := createTestOrchestrator()

	// Registrar workflow
	workflow := NewSimpleValidationWorkflow()
	err := orchestrator.RegisterWorkflow(workflow)
	if err != nil {
		t.Fatalf("Erro ao registrar workflow: %v", err)
	}

	numExecutions := 3
	params := map[string]interface{}{
		"file_path": "concurrent_test_file.xlsx",
	}

	// Executar múltiplos workflows concorrentemente
	executionIDs := make([]string, numExecutions)
	for i := 0; i < numExecutions; i++ {
		executionID, err := orchestrator.ExecuteWorkflowAsync("simple-validation", params)
		if err != nil {
			t.Fatalf("Erro ao executar workflow async %d: %v", i, err)
		}
		executionIDs[i] = executionID
	}

	// Verificar que todas as execuções foram iniciadas
	activeExecutions := orchestrator.ListActiveExecutions()
	t.Logf("Execuções ativas iniciadas: %d", len(activeExecutions))

	// Aguardar todas terminarem
	time.Sleep(2 * time.Second)

	// Verificar estatísticas
	stats := orchestrator.GetStats()
	t.Logf("Estatísticas finais: %+v", stats)
}

func testWorkflowValidation(t *testing.T) {
	orchestrator := createTestOrchestrator()

	// Tentar registrar workflow inválido (nil)
	err := orchestrator.RegisterWorkflow(nil)
	if err == nil {
		t.Error("Deveria falhar ao registrar workflow nil")
	}

	// Tentar executar workflow inexistente
	_, err = orchestrator.ExecuteWorkflow("workflow-inexistente", nil)
	if err == nil {
		t.Error("Deveria falhar ao executar workflow inexistente")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Erro deveria mencionar workflow não encontrado: %v", err)
	}
}

func testErrorHandling(t *testing.T) {
	orchestrator := createTestOrchestrator()

	// Criar workflow que falha
	failingWorkflow := NewFailingWorkflow()
	err := orchestrator.RegisterWorkflow(failingWorkflow)
	if err != nil {
		t.Fatalf("Erro ao registrar workflow que falha: %v", err)
	}

	// Executar workflow que falha
	result, err := orchestrator.ExecuteWorkflow("failing-workflow", nil)
	if err != nil {
		t.Fatalf("Erro ao executar workflow que falha: %v", err)
	}

	// Verificar que o resultado indica falha
	if result.Status != StatusFailed {
		t.Errorf("Status deveria ser Failed, got %s", result.Status)
	}

	if result.Error == nil {
		t.Error("Resultado deveria ter erro")
	}
}

// Helper functions

func createTestOrchestrator() *Orchestrator {
	logger := &DefaultLogger{Logger: log.New(os.Stdout, "[TEST] ", log.LstdFlags)}
	config := OrchestratorConfig{
		MaxConcurrentWorkflows: 3,
		DefaultTimeout:         10 * time.Second,
		EnableRollback:         true,
		DetailedLogging:        false, // Menos verbose para testes
	}

	return NewOrchestrator(logger, config)
}

// FailingWorkflow - Workflow de teste que sempre falha
type FailingWorkflow struct {
	*BaseWorkflow
}

func NewFailingWorkflow() *FailingWorkflow {
	steps := []WorkflowStep{
		NewFailingStep(),
	}

	baseWorkflow := NewBaseWorkflow(
		"failing-workflow",
		"Workflow que sempre falha para testes",
		steps,
	)

	return &FailingWorkflow{
		BaseWorkflow: baseWorkflow,
	}
}

// FailingStep - Step que sempre falha
type FailingStep struct {
	*BaseStep
}

func NewFailingStep() *FailingStep {
	return &FailingStep{
		BaseStep: NewBaseStep(
			"failing-step",
			"Step que sempre falha",
			1*time.Second,
		),
	}
}

func (s *FailingStep) Execute(ctx *WorkflowContext) error {
	ctx.Logger.Debug("Executando step que vai falhar: %s", s.Name())
	return context.DeadlineExceeded // Simular erro
}
