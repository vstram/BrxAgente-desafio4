package workflows

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Orchestrator é responsável por executar workflows
type Orchestrator struct {
	// Workflows registrados
	workflows map[string]Workflow
	mutex     sync.RWMutex

	// Logger
	logger Logger

	// Execuções ativas
	activeExecutions map[string]*WorkflowExecution
	executionMutex   sync.RWMutex

	// Configurações
	config OrchestratorConfig
}

// OrchestratorConfig contém configurações do orquestrador
type OrchestratorConfig struct {
	MaxConcurrentWorkflows int
	DefaultTimeout         time.Duration
	EnableRollback         bool
	DetailedLogging        bool
}

// WorkflowExecution representa uma execução ativa de workflow
type WorkflowExecution struct {
	ID           string
	WorkflowName string
	Status       WorkflowStatus
	StartTime    time.Time
	Context      *WorkflowContext
	Result       *WorkflowResult
	CancelFunc   context.CancelFunc
}

// DefaultLogger implementação simples de Logger
type DefaultLogger struct {
	*log.Logger
}

func (l *DefaultLogger) Debug(msg string, args ...interface{}) {
	l.Printf("[DEBUG] "+msg, args...)
}

func (l *DefaultLogger) Info(msg string, args ...interface{}) {
	l.Printf("[INFO] "+msg, args...)
}

func (l *DefaultLogger) Warn(msg string, args ...interface{}) {
	l.Printf("[WARN] "+msg, args...)
}

func (l *DefaultLogger) Error(msg string, args ...interface{}) {
	l.Printf("[ERROR] "+msg, args...)
}

// NewOrchestrator cria um novo orquestrador
func NewOrchestrator(logger Logger, config OrchestratorConfig) *Orchestrator {
	if logger == nil {
		logger = &DefaultLogger{Logger: log.Default()}
	}

	if config.MaxConcurrentWorkflows == 0 {
		config.MaxConcurrentWorkflows = 10
	}

	if config.DefaultTimeout == 0 {
		config.DefaultTimeout = 30 * time.Minute
	}

	return &Orchestrator{
		workflows:        make(map[string]Workflow),
		logger:           logger,
		activeExecutions: make(map[string]*WorkflowExecution),
		config:           config,
	}
}

// RegisterWorkflow registra um novo workflow
func (o *Orchestrator) RegisterWorkflow(workflow Workflow) error {
	if workflow == nil {
		return fmt.Errorf("workflow cannot be nil")
	}

	// Validar workflow
	if err := workflow.Validate(); err != nil {
		return fmt.Errorf("workflow validation failed: %w", err)
	}

	o.mutex.Lock()
	defer o.mutex.Unlock()

	name := workflow.Name()
	if _, exists := o.workflows[name]; exists {
		return fmt.Errorf("workflow '%s' already registered", name)
	}

	o.workflows[name] = workflow
	o.logger.Info("Workflow registered: %s", name)

	return nil
}

// UnregisterWorkflow remove um workflow
func (o *Orchestrator) UnregisterWorkflow(name string) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	if _, exists := o.workflows[name]; !exists {
		return fmt.Errorf("workflow '%s' not found", name)
	}

	delete(o.workflows, name)
	o.logger.Info("Workflow unregistered: %s", name)

	return nil
}

// GetWorkflow recupera um workflow registrado
func (o *Orchestrator) GetWorkflow(name string) (Workflow, error) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	workflow, exists := o.workflows[name]
	if !exists {
		return nil, fmt.Errorf("workflow '%s' not found", name)
	}

	return workflow, nil
}

// ListWorkflows retorna lista de workflows registrados
func (o *Orchestrator) ListWorkflows() []string {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	names := make([]string, 0, len(o.workflows))
	for name := range o.workflows {
		names = append(names, name)
	}

	return names
}

// ExecuteWorkflow executa um workflow
func (o *Orchestrator) ExecuteWorkflow(name string, params map[string]interface{}) (*WorkflowResult, error) {
	return o.ExecuteWorkflowWithTimeout(name, params, o.config.DefaultTimeout)
}

// ExecuteWorkflowWithTimeout executa um workflow com timeout específico
func (o *Orchestrator) ExecuteWorkflowWithTimeout(name string, params map[string]interface{}, timeout time.Duration) (*WorkflowResult, error) {
	// Verificar limite de execuções concorrentes
	if o.getActiveExecutionsCount() >= o.config.MaxConcurrentWorkflows {
		return nil, fmt.Errorf("maximum concurrent workflows exceeded (%d)", o.config.MaxConcurrentWorkflows)
	}

	// Obter workflow
	workflow, err := o.GetWorkflow(name)
	if err != nil {
		return nil, err
	}

	// Criar contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	// Gerar ID único para execução
	executionID := uuid.New().String()

	// Criar contexto do workflow
	workflowCtx := NewWorkflowContext(ctx, name, executionID, params, o.logger)

	// Criar execução
	execution := &WorkflowExecution{
		ID:           executionID,
		WorkflowName: name,
		Status:       StatusPending,
		StartTime:    time.Now(),
		Context:      workflowCtx,
		CancelFunc:   cancel,
	}

	// Registrar execução ativa
	o.registerExecution(execution)
	defer o.unregisterExecution(executionID)
	defer cancel()

	o.logger.Info("Starting workflow execution: %s (ID: %s)", name, executionID)

	// Executar workflow
	result := o.executeWorkflowInternal(workflow, workflowCtx)

	// Atualizar execução com resultado
	execution.Result = result
	execution.Status = result.Status

	o.logger.Info("Workflow execution completed: %s (ID: %s) Status: %s Duration: %v",
		name, executionID, result.Status, result.Duration)

	return result, nil
}

// ExecuteWorkflowAsync executa um workflow de forma assíncrona
func (o *Orchestrator) ExecuteWorkflowAsync(name string, params map[string]interface{}) (string, error) {
	// Verificar limite de execuções concorrentes
	if o.getActiveExecutionsCount() >= o.config.MaxConcurrentWorkflows {
		return "", fmt.Errorf("maximum concurrent workflows exceeded (%d)", o.config.MaxConcurrentWorkflows)
	}

	// Obter workflow
	workflow, err := o.GetWorkflow(name)
	if err != nil {
		return "", err
	}

	// Gerar ID único para execução
	executionID := uuid.New().String()

	// Criar contexto
	ctx, cancel := context.WithTimeout(context.Background(), o.config.DefaultTimeout)
	workflowCtx := NewWorkflowContext(ctx, name, executionID, params, o.logger)

	// Criar execução
	execution := &WorkflowExecution{
		ID:           executionID,
		WorkflowName: name,
		Status:       StatusPending,
		StartTime:    time.Now(),
		Context:      workflowCtx,
		CancelFunc:   cancel,
	}

	// Registrar execução ativa
	o.registerExecution(execution)

	// Executar assincronamente
	go func() {
		defer o.unregisterExecution(executionID)
		defer cancel()

		o.logger.Info("Starting async workflow execution: %s (ID: %s)", name, executionID)

		result := o.executeWorkflowInternal(workflow, workflowCtx)

		execution.Result = result
		execution.Status = result.Status

		o.logger.Info("Async workflow execution completed: %s (ID: %s) Status: %s Duration: %v",
			name, executionID, result.Status, result.Duration)
	}()

	return executionID, nil
}

// executeWorkflowInternal executa o workflow internamente
func (o *Orchestrator) executeWorkflowInternal(workflow Workflow, ctx *WorkflowContext) *WorkflowResult {
	startTime := time.Now()

	result := &WorkflowResult{
		WorkflowName: workflow.Name(),
		Status:       StatusRunning,
		StartTime:    startTime,
		StepResults:  make([]StepResult, 0),
		Metadata:     make(map[string]interface{}),
	}

	steps := workflow.Steps()
	result.TotalSteps = len(steps)

	o.logger.Info("Executing workflow '%s' with %d steps", workflow.Name(), len(steps))

	// Executar cada step
	completedSteps := make([]WorkflowStep, 0)

	for i, step := range steps {
		if ctx.IsCanceled() {
			o.logger.Warn("Workflow execution canceled at step %d", i+1)
			result.Status = StatusCanceled
			break
		}

		stepResult := o.executeStep(step, ctx)
		result.StepResults = append(result.StepResults, stepResult)

		if stepResult.Error != nil {
			o.logger.Error("Step '%s' failed: %v", step.Name(), stepResult.Error)
			result.Status = StatusFailed
			result.Error = NewWorkflowError(workflow.Name(), step.Name(), "step execution failed", stepResult.Error)

			// Executar rollback se habilitado
			if o.config.EnableRollback && len(completedSteps) > 0 {
				o.logger.Info("Starting rollback for %d completed steps", len(completedSteps))
				o.rollbackSteps(completedSteps, ctx)
				result.Status = StatusRolledBack
			}

			break
		}

		if !stepResult.Skipped {
			completedSteps = append(completedSteps, step)
			result.CompletedSteps++
		}

		o.logger.Info("Step '%s' completed successfully in %v", step.Name(), stepResult.Duration)
	}

	// Finalizar resultado
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	if result.Status == StatusRunning && !ctx.IsCanceled() {
		result.Status = StatusCompleted
	}

	return result
}

// executeStep executa um step individual
func (o *Orchestrator) executeStep(step WorkflowStep, ctx *WorkflowContext) StepResult {
	stepResult := StepResult{
		StepName:  step.Name(),
		StartTime: time.Now(),
		Status:    "running",
	}

	// Verificar se o step pode ser pulado
	if step.CanSkip(ctx) {
		o.logger.Debug("Skipping step: %s", step.Name())
		stepResult.Status = "skipped"
		stepResult.Skipped = true
		stepResult.EndTime = time.Now()
		stepResult.Duration = stepResult.EndTime.Sub(stepResult.StartTime)
		return stepResult
	}

	o.logger.Debug("Executing step: %s", step.Name())

	// Executar step
	if err := step.Execute(ctx); err != nil {
		stepResult.Status = "failed"
		stepResult.Error = err
		ctx.SetError(step.Name(), err)
	} else {
		stepResult.Status = "completed"
	}

	stepResult.EndTime = time.Now()
	stepResult.Duration = stepResult.EndTime.Sub(stepResult.StartTime)

	return stepResult
}

// rollbackSteps executa rollback dos steps concluídos
func (o *Orchestrator) rollbackSteps(completedSteps []WorkflowStep, ctx *WorkflowContext) {
	// Rollback em ordem reversa
	for i := len(completedSteps) - 1; i >= 0; i-- {
		step := completedSteps[i]
		o.logger.Debug("Rolling back step: %s", step.Name())

		if err := step.Rollback(ctx); err != nil {
			o.logger.Error("Rollback failed for step '%s': %v", step.Name(), err)
		}
	}
}

// CancelExecution cancela uma execução ativa
func (o *Orchestrator) CancelExecution(executionID string) error {
	o.executionMutex.RLock()
	execution, exists := o.activeExecutions[executionID]
	o.executionMutex.RUnlock()

	if !exists {
		return fmt.Errorf("execution '%s' not found", executionID)
	}

	execution.Context.Cancel()
	execution.CancelFunc()

	o.logger.Info("Execution canceled: %s", executionID)

	return nil
}

// GetExecution retorna informações sobre uma execução
func (o *Orchestrator) GetExecution(executionID string) (*WorkflowExecution, error) {
	o.executionMutex.RLock()
	defer o.executionMutex.RUnlock()

	execution, exists := o.activeExecutions[executionID]
	if !exists {
		return nil, fmt.Errorf("execution '%s' not found", executionID)
	}

	return execution, nil
}

// ListActiveExecutions retorna lista de execuções ativas
func (o *Orchestrator) ListActiveExecutions() []string {
	o.executionMutex.RLock()
	defer o.executionMutex.RUnlock()

	ids := make([]string, 0, len(o.activeExecutions))
	for id := range o.activeExecutions {
		ids = append(ids, id)
	}

	return ids
}

// GetStats retorna estatísticas do orquestrador
func (o *Orchestrator) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"registered_workflows":     len(o.workflows),
		"active_executions":        o.getActiveExecutionsCount(),
		"max_concurrent_workflows": o.config.MaxConcurrentWorkflows,
		"rollback_enabled":         o.config.EnableRollback,
	}
}

// registerExecution registra uma execução ativa
func (o *Orchestrator) registerExecution(execution *WorkflowExecution) {
	o.executionMutex.Lock()
	defer o.executionMutex.Unlock()
	o.activeExecutions[execution.ID] = execution
}

// unregisterExecution remove uma execução ativa
func (o *Orchestrator) unregisterExecution(executionID string) {
	o.executionMutex.Lock()
	defer o.executionMutex.Unlock()
	delete(o.activeExecutions, executionID)
}

// getActiveExecutionsCount retorna o número de execuções ativas
func (o *Orchestrator) getActiveExecutionsCount() int {
	o.executionMutex.RLock()
	defer o.executionMutex.RUnlock()
	return len(o.activeExecutions)
}
