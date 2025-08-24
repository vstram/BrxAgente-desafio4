package parallel

import (
	"context"
	"fmt"
	"sync"
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

// BatchProcessor processador otimizado para grandes lotes de dados
type BatchProcessor struct {
	workerPool    *WorkerPool
	batchSize     int
	maxConcurrent int
	timeout       time.Duration
}

// BatchConfig configuração do processador de lotes
type BatchConfig struct {
	WorkerPoolSize int
	BatchSize      int
	MaxConcurrent  int
	TimeoutSeconds int
}

// NewBatchProcessor cria um novo processador de lotes
func NewBatchProcessor(config BatchConfig) *BatchProcessor {
	timeout := time.Duration(config.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &BatchProcessor{
		workerPool:    NewWorkerPool(config.WorkerPoolSize, config.BatchSize*2),
		batchSize:     config.BatchSize,
		maxConcurrent: config.MaxConcurrent,
		timeout:       timeout,
	}
}

// ProcessColaboradoresBatch processa lote de colaboradores em paralelo
func (bp *BatchProcessor) ProcessColaboradoresBatch(
	colaboradores map[string]*modelo.Colaborador,
	processor func(*modelo.Colaborador) (*modelo.Colaborador, error),
) (map[string]*modelo.Colaborador, []error) {

	// Converte map em slice para processamento
	colaboradoresSlice := make([]*modelo.Colaborador, 0, len(colaboradores))
	for _, col := range colaboradores {
		colaboradoresSlice = append(colaboradoresSlice, col)
	}

	// Divide em lotes
	batches := bp.dividirEmLotes(colaboradoresSlice)

	// Processa lotes em paralelo
	resultados := make(map[string]*modelo.Colaborador)
	var erros []error
	var mutex sync.Mutex

	semaphore := make(chan struct{}, bp.maxConcurrent)
	var wg sync.WaitGroup

	for _, batch := range batches {
		wg.Add(1)
		go func(batch []*modelo.Colaborador) {
			defer wg.Done()
			semaphore <- struct{}{}        // Adquire semáforo
			defer func() { <-semaphore }() // Libera semáforo

			batchResults, batchErrors := bp.processColaboradorBatch(batch, processor)

			mutex.Lock()
			for k, v := range batchResults {
				resultados[k] = v
			}
			erros = append(erros, batchErrors...)
			mutex.Unlock()
		}(batch)
	}

	wg.Wait()
	return resultados, erros
}

// ProcessFilesBatch processa múltiplos arquivos em paralelo
func (bp *BatchProcessor) ProcessFilesBatch(
	filePaths []string,
	processor func(string) (interface{}, error),
) ([]interface{}, []error) {

	tasks := make([]Task, len(filePaths))
	for i, filePath := range filePaths {
		tasks[i] = &FileTask{
			ID:        fmt.Sprintf("file_%d", i),
			FilePath:  filePath,
			Processor: processor,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), bp.timeout)
	defer cancel()

	results, errors := bp.workerPool.ProcessBatchWithContext(ctx, tasks)

	// Converte results para interface{} slice
	processedResults := make([]interface{}, 0, len(results))
	for _, result := range results {
		if fileResult, ok := result.(*FileResult); ok {
			processedResults = append(processedResults, fileResult.Data)
		}
	}

	return processedResults, errors
}

// ProcessGenericBatch processa lote genérico de dados
func (bp *BatchProcessor) ProcessGenericBatch(
	data []interface{},
	processor func(interface{}) (interface{}, error),
) ([]interface{}, []error) {

	tasks := make([]Task, len(data))
	for i, item := range data {
		tasks[i] = &GenericTask{
			ID:        fmt.Sprintf("generic_%d", i),
			Data:      item,
			Processor: processor,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), bp.timeout)
	defer cancel()

	results, errors := bp.workerPool.ProcessBatchWithContext(ctx, tasks)

	processedResults := make([]interface{}, 0, len(results))
	for _, result := range results {
		if genericResult, ok := result.(*GenericResult); ok {
			processedResults = append(processedResults, genericResult.Data)
		}
	}

	return processedResults, errors
}

// dividirEmLotes divide slice de colaboradores em lotes menores
func (bp *BatchProcessor) dividirEmLotes(colaboradores []*modelo.Colaborador) [][]*modelo.Colaborador {
	var batches [][]*modelo.Colaborador

	for i := 0; i < len(colaboradores); i += bp.batchSize {
		end := i + bp.batchSize
		if end > len(colaboradores) {
			end = len(colaboradores)
		}
		batches = append(batches, colaboradores[i:end])
	}

	return batches
}

// processColaboradorBatch processa um lote de colaboradores
func (bp *BatchProcessor) processColaboradorBatch(
	batch []*modelo.Colaborador,
	processor func(*modelo.Colaborador) (*modelo.Colaborador, error),
) (map[string]*modelo.Colaborador, []error) {

	resultados := make(map[string]*modelo.Colaborador)
	var erros []error

	for _, colaborador := range batch {
		result, err := processor(colaborador)
		if err != nil {
			erros = append(erros, fmt.Errorf("erro processando colaborador %s: %w", colaborador.Matricula, err))
		} else if result != nil {
			resultados[result.Matricula] = result
		}
	}

	return resultados, erros
}

// GetStats retorna estatísticas do processador
func (bp *BatchProcessor) GetStats() BatchProcessorStats {
	poolStats := bp.workerPool.GetStats()

	return BatchProcessorStats{
		WorkerPoolStats: poolStats,
		BatchSize:       bp.batchSize,
		MaxConcurrent:   bp.maxConcurrent,
		Timeout:         bp.timeout,
	}
}

// BatchProcessorStats estatísticas do processador de lotes
type BatchProcessorStats struct {
	WorkerPoolStats PoolStats
	BatchSize       int
	MaxConcurrent   int
	Timeout         time.Duration
}

// Close fecha o processador de lotes
func (bp *BatchProcessor) Close() error {
	return bp.workerPool.Close()
}

// FileTask tarefa para processamento de arquivo
type FileTask struct {
	ID        string
	FilePath  string
	Processor func(string) (interface{}, error)
}

func (ft *FileTask) Execute(ctx context.Context) (Result, error) {
	start := time.Now()
	data, err := ft.Processor(ft.FilePath)
	duration := time.Since(start)

	if err != nil {
		return nil, err
	}

	return &FileResult{
		ID:       ft.ID,
		Data:     data,
		Duration: duration,
	}, nil
}

func (ft *FileTask) GetID() string {
	return ft.ID
}

func (ft *FileTask) GetPriority() int {
	return 0 // Prioridade padrão
}

// FileResult resultado do processamento de arquivo
type FileResult struct {
	ID       string
	Data     interface{}
	Duration time.Duration
}

func (fr *FileResult) GetID() string {
	return fr.ID
}

func (fr *FileResult) GetData() interface{} {
	return fr.Data
}

func (fr *FileResult) GetDuration() time.Duration {
	return fr.Duration
}

// GenericTask tarefa genérica
type GenericTask struct {
	ID        string
	Data      interface{}
	Processor func(interface{}) (interface{}, error)
}

func (gt *GenericTask) Execute(ctx context.Context) (Result, error) {
	start := time.Now()
	result, err := gt.Processor(gt.Data)
	duration := time.Since(start)

	if err != nil {
		return nil, err
	}

	return &GenericResult{
		ID:       gt.ID,
		Data:     result,
		Duration: duration,
	}, nil
}

func (gt *GenericTask) GetID() string {
	return gt.ID
}

func (gt *GenericTask) GetPriority() int {
	return 0
}

// GenericResult resultado genérico
type GenericResult struct {
	ID       string
	Data     interface{}
	Duration time.Duration
}

func (gr *GenericResult) GetID() string {
	return gr.ID
}

func (gr *GenericResult) GetData() interface{} {
	return gr.Data
}

func (gr *GenericResult) GetDuration() time.Duration {
	return gr.Duration
}
