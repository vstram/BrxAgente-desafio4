package workflows

import (
	"fmt"
	"time"
)

// Workflow define a interface base para workflows
type Workflow interface {
	Name() string
	Description() string
	Steps() []WorkflowStep
	Execute(ctx *WorkflowContext) error
	Validate() error
}

// WorkflowStep define a interface para steps individuais
type WorkflowStep interface {
	Name() string
	Description() string
	Execute(ctx *WorkflowContext) error
	Rollback(ctx *WorkflowContext) error
	CanSkip(ctx *WorkflowContext) bool
	EstimatedDuration() time.Duration
}

// WorkflowStatus representa os possíveis status de um workflow
type WorkflowStatus int

const (
	StatusPending WorkflowStatus = iota
	StatusRunning
	StatusCompleted
	StatusFailed
	StatusCanceled
	StatusRolledBack
)

func (s WorkflowStatus) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusRunning:
		return "running"
	case StatusCompleted:
		return "completed"
	case StatusFailed:
		return "failed"
	case StatusCanceled:
		return "canceled"
	case StatusRolledBack:
		return "rolledback"
	default:
		return "unknown"
	}
}

// WorkflowResult contém o resultado da execução de um workflow
type WorkflowResult struct {
	WorkflowName   string                 `json:"workflow_name"`
	Status         WorkflowStatus         `json:"status"`
	StartTime      time.Time              `json:"start_time"`
	EndTime        time.Time              `json:"end_time"`
	Duration       time.Duration          `json:"duration"`
	CompletedSteps int                    `json:"completed_steps"`
	TotalSteps     int                    `json:"total_steps"`
	Error          error                  `json:"error,omitempty"`
	StepResults    []StepResult           `json:"step_results"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// StepResult contém o resultado da execução de um step
type StepResult struct {
	StepName  string        `json:"step_name"`
	Status    string        `json:"status"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Error     error         `json:"error,omitempty"`
	Output    interface{}   `json:"output,omitempty"`
	Skipped   bool          `json:"skipped"`
}

// BaseWorkflow fornece implementação base comum para workflows
type BaseWorkflow struct {
	name        string
	description string
	steps       []WorkflowStep
}

// NewBaseWorkflow cria um novo BaseWorkflow
func NewBaseWorkflow(name, description string, steps []WorkflowStep) *BaseWorkflow {
	return &BaseWorkflow{
		name:        name,
		description: description,
		steps:       steps,
	}
}

// Name retorna o nome do workflow
func (w *BaseWorkflow) Name() string {
	return w.name
}

// Description retorna a descrição do workflow
func (w *BaseWorkflow) Description() string {
	return w.description
}

// Steps retorna os steps do workflow
func (w *BaseWorkflow) Steps() []WorkflowStep {
	return w.steps
}

// Validate valida se o workflow está corretamente configurado
func (w *BaseWorkflow) Validate() error {
	if w.name == "" {
		return fmt.Errorf("workflow name cannot be empty")
	}

	if len(w.steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}

	// Validar cada step
	stepNames := make(map[string]bool)
	for i, step := range w.steps {
		if step == nil {
			return fmt.Errorf("step %d is nil", i)
		}

		stepName := step.Name()
		if stepName == "" {
			return fmt.Errorf("step %d has empty name", i)
		}

		if stepNames[stepName] {
			return fmt.Errorf("duplicate step name: %s", stepName)
		}
		stepNames[stepName] = true
	}

	return nil
}

// Execute executa o workflow (implementação padrão usando orchestrator)
func (w *BaseWorkflow) Execute(ctx *WorkflowContext) error {
	// Esta implementação será sobrescrita pelo orchestrator
	// Cada workflow concreto pode sobrescrever se necessário
	return fmt.Errorf("execute method should be called via orchestrator")
}

// BaseStep fornece implementação base comum para steps
type BaseStep struct {
	name              string
	description       string
	estimatedDuration time.Duration
}

// NewBaseStep cria um novo BaseStep
func NewBaseStep(name, description string, duration time.Duration) *BaseStep {
	return &BaseStep{
		name:              name,
		description:       description,
		estimatedDuration: duration,
	}
}

// Name retorna o nome do step
func (s *BaseStep) Name() string {
	return s.name
}

// Description retorna a descrição do step
func (s *BaseStep) Description() string {
	return s.description
}

// EstimatedDuration retorna a duração estimada do step
func (s *BaseStep) EstimatedDuration() time.Duration {
	return s.estimatedDuration
}

// CanSkip implementação padrão (nunca pular)
func (s *BaseStep) CanSkip(ctx *WorkflowContext) bool {
	return false
}

// Rollback implementação padrão (sem rollback)
func (s *BaseStep) Rollback(ctx *WorkflowContext) error {
	// Por padrão, steps não fazem rollback
	return nil
}

// WorkflowError representa um erro específico de workflow
type WorkflowError struct {
	WorkflowName string
	StepName     string
	Message      string
	Cause        error
}

func (e *WorkflowError) Error() string {
	if e.StepName != "" {
		return fmt.Sprintf("workflow '%s' step '%s': %s", e.WorkflowName, e.StepName, e.Message)
	}
	return fmt.Sprintf("workflow '%s': %s", e.WorkflowName, e.Message)
}

func (e *WorkflowError) Unwrap() error {
	return e.Cause
}

// NewWorkflowError cria um novo WorkflowError
func NewWorkflowError(workflowName, stepName, message string, cause error) *WorkflowError {
	return &WorkflowError{
		WorkflowName: workflowName,
		StepName:     stepName,
		Message:      message,
		Cause:        cause,
	}
}
