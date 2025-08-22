package workflows

import (
	"context"
	"sync"
	"time"
)

// WorkflowContext fornece contexto compartilhado durante a execução do workflow
type WorkflowContext struct {
	// Context padrão do Go para cancelamento e timeouts
	ctx context.Context
	
	// Dados compartilhados entre steps
	data  map[string]interface{}
	mutex sync.RWMutex
	
	// Metadados de execução
	WorkflowName string
	ExecutionID  string
	StartTime    time.Time
	
	// Parâmetros de entrada
	Parameters map[string]interface{}
	
	// Canal para cancelamento
	cancelChan chan struct{}
	
	// Logger para este contexto específico
	Logger Logger
}

// Logger interface para logging no contexto do workflow
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// NewWorkflowContext cria um novo contexto de workflow
func NewWorkflowContext(ctx context.Context, workflowName, executionID string, params map[string]interface{}, logger Logger) *WorkflowContext {
	if params == nil {
		params = make(map[string]interface{})
	}
	
	return &WorkflowContext{
		ctx:          ctx,
		data:         make(map[string]interface{}),
		WorkflowName: workflowName,
		ExecutionID:  executionID,
		StartTime:    time.Now(),
		Parameters:   params,
		cancelChan:   make(chan struct{}, 1),
		Logger:       logger,
	}
}

// Context retorna o context.Context subjacente
func (wc *WorkflowContext) Context() context.Context {
	return wc.ctx
}

// Set armazena um valor no contexto
func (wc *WorkflowContext) Set(key string, value interface{}) {
	wc.mutex.Lock()
	defer wc.mutex.Unlock()
	wc.data[key] = value
}

// Get recupera um valor do contexto
func (wc *WorkflowContext) Get(key string) (interface{}, bool) {
	wc.mutex.RLock()
	defer wc.mutex.RUnlock()
	value, exists := wc.data[key]
	return value, exists
}

// GetString recupera um valor como string
func (wc *WorkflowContext) GetString(key string) (string, bool) {
	value, exists := wc.Get(key)
	if !exists {
		return "", false
	}
	
	if str, ok := value.(string); ok {
		return str, true
	}
	return "", false
}

// GetInt recupera um valor como int
func (wc *WorkflowContext) GetInt(key string) (int, bool) {
	value, exists := wc.Get(key)
	if !exists {
		return 0, false
	}
	
	if i, ok := value.(int); ok {
		return i, true
	}
	return 0, false
}

// GetBool recupera um valor como bool
func (wc *WorkflowContext) GetBool(key string) (bool, bool) {
	value, exists := wc.Get(key)
	if !exists {
		return false, false
	}
	
	if b, ok := value.(bool); ok {
		return b, true
	}
	return false, false
}

// GetFloat64 recupera um valor como float64
func (wc *WorkflowContext) GetFloat64(key string) (float64, bool) {
	value, exists := wc.Get(key)
	if !exists {
		return 0, false
	}
	
	if f, ok := value.(float64); ok {
		return f, true
	}
	return 0, false
}

// GetAll retorna todos os dados do contexto (cópia)
func (wc *WorkflowContext) GetAll() map[string]interface{} {
	wc.mutex.RLock()
	defer wc.mutex.RUnlock()
	
	result := make(map[string]interface{})
	for k, v := range wc.data {
		result[k] = v
	}
	return result
}

// Has verifica se uma chave existe no contexto
func (wc *WorkflowContext) Has(key string) bool {
	_, exists := wc.Get(key)
	return exists
}

// Delete remove uma chave do contexto
func (wc *WorkflowContext) Delete(key string) {
	wc.mutex.Lock()
	defer wc.mutex.Unlock()
	delete(wc.data, key)
}

// Clear limpa todos os dados do contexto
func (wc *WorkflowContext) Clear() {
	wc.mutex.Lock()
	defer wc.mutex.Unlock()
	wc.data = make(map[string]interface{})
}

// Size retorna o número de itens no contexto
func (wc *WorkflowContext) Size() int {
	wc.mutex.RLock()
	defer wc.mutex.RUnlock()
	return len(wc.data)
}

// IsCanceled verifica se o workflow foi cancelado
func (wc *WorkflowContext) IsCanceled() bool {
	select {
	case <-wc.cancelChan:
		return true
	case <-wc.ctx.Done():
		return true
	default:
		return false
	}
}

// Cancel cancela a execução do workflow
func (wc *WorkflowContext) Cancel() {
	select {
	case wc.cancelChan <- struct{}{}:
	default:
		// Canal já tem um valor ou está fechado
	}
}

// Elapsed retorna o tempo decorrido desde o início da execução
func (wc *WorkflowContext) Elapsed() time.Duration {
	return time.Since(wc.StartTime)
}

// GetParameter recupera um parâmetro de entrada
func (wc *WorkflowContext) GetParameter(key string) (interface{}, bool) {
	value, exists := wc.Parameters[key]
	return value, exists
}

// GetParameterString recupera um parâmetro como string
func (wc *WorkflowContext) GetParameterString(key string) (string, bool) {
	value, exists := wc.GetParameter(key)
	if !exists {
		return "", false
	}
	
	if str, ok := value.(string); ok {
		return str, true
	}
	return "", false
}

// GetParameterInt recupera um parâmetro como int
func (wc *WorkflowContext) GetParameterInt(key string) (int, bool) {
	value, exists := wc.GetParameter(key)
	if !exists {
		return 0, false
	}
	
	if i, ok := value.(int); ok {
		return i, true
	}
	return 0, false
}

// SetResult armazena o resultado de um step
func (wc *WorkflowContext) SetResult(stepName string, result interface{}) {
	wc.Set("result_"+stepName, result)
}

// GetResult recupera o resultado de um step
func (wc *WorkflowContext) GetResult(stepName string) (interface{}, bool) {
	return wc.Get("result_" + stepName)
}

// SetError armazena um erro de um step
func (wc *WorkflowContext) SetError(stepName string, err error) {
	wc.Set("error_"+stepName, err)
}

// GetError recupera o erro de um step
func (wc *WorkflowContext) GetError(stepName string) (error, bool) {
	value, exists := wc.Get("error_" + stepName)
	if !exists {
		return nil, false
	}
	
	if err, ok := value.(error); ok {
		return err, true
	}
	return nil, false
}

// Clone cria uma cópia do contexto (shallow copy dos dados)
func (wc *WorkflowContext) Clone() *WorkflowContext {
	wc.mutex.RLock()
	defer wc.mutex.RUnlock()
	
	newCtx := &WorkflowContext{
		ctx:          wc.ctx,
		data:         make(map[string]interface{}),
		WorkflowName: wc.WorkflowName,
		ExecutionID:  wc.ExecutionID,
		StartTime:    wc.StartTime,
		Parameters:   make(map[string]interface{}),
		cancelChan:   make(chan struct{}, 1),
		Logger:       wc.Logger,
	}
	
	// Copiar dados
	for k, v := range wc.data {
		newCtx.data[k] = v
	}
	
	// Copiar parâmetros
	for k, v := range wc.Parameters {
		newCtx.Parameters[k] = v
	}
	
	return newCtx
}