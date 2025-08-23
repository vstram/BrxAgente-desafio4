// Package parallel implementa processamento paralelo para otimização de performance
package parallel

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Task representa uma tarefa a ser executada
type Task interface {
	Execute(ctx context.Context) (Result, error)
	GetID() string
	GetPriority() int
}

// Result representa o resultado de uma tarefa
type Result interface {
	GetID() string
	GetData() interface{}
	GetDuration() time.Duration
}

// WorkerPool implementa um pool de workers para processamento paralelo
type WorkerPool struct {
	workers     int
	taskQueue   chan Task
	resultQueue chan Result
	errorQueue  chan error
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	stats       *PoolStats
	mutex       sync.RWMutex
}

// PoolStats estatísticas do worker pool
type PoolStats struct {
	TasksProcessed   int64
	TasksSuccessful  int64
	TasksFailed      int64
	TotalDuration    time.Duration
	AverageDuration  time.Duration
	ActiveWorkers    int
	QueuedTasks      int
	MaxQueueSize     int
	WorkersStartTime time.Time
}

// NewWorkerPool cria um novo worker pool
func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		workers:     workers,
		taskQueue:   make(chan Task, queueSize),
		resultQueue: make(chan Result, queueSize),
		errorQueue:  make(chan error, queueSize),
		ctx:         ctx,
		cancel:      cancel,
		stats: &PoolStats{
			MaxQueueSize:     queueSize,
			WorkersStartTime: time.Now(),
		},
	}

	// Inicia workers
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	// Inicia goroutine para coleta de estatísticas
	go pool.statsCollector()

	return pool
}

// worker executa tarefas do pool
func (wp *WorkerPool) worker(_ int) {
	defer wp.wg.Done()

	wp.mutex.Lock()
	wp.stats.ActiveWorkers++
	wp.mutex.Unlock()

	defer func() {
		wp.mutex.Lock()
		wp.stats.ActiveWorkers--
		wp.mutex.Unlock()
	}()

	for {
		select {
		case task := <-wp.taskQueue:
			if task == nil {
				return
			}

			startTime := time.Now()
			result, err := task.Execute(wp.ctx)
			duration := time.Since(startTime)

			wp.mutex.Lock()
			wp.stats.TasksProcessed++
			wp.stats.TotalDuration += duration
			wp.stats.AverageDuration = time.Duration(int64(wp.stats.TotalDuration) / wp.stats.TasksProcessed)

			if err != nil {
				wp.stats.TasksFailed++
				wp.mutex.Unlock()
				select {
				case wp.errorQueue <- err:
				default:
					// Error queue cheio, ignora erro
				}
			} else {
				wp.stats.TasksSuccessful++
				wp.mutex.Unlock()
				select {
				case wp.resultQueue <- result:
				default:
					// Result queue cheio, ignora resultado
				}
			}

		case <-wp.ctx.Done():
			return
		}
	}
}

// SubmitTask envia uma tarefa para o pool
func (wp *WorkerPool) SubmitTask(task Task) error {
	select {
	case wp.taskQueue <- task:
		wp.mutex.Lock()
		wp.stats.QueuedTasks = len(wp.taskQueue)
		wp.mutex.Unlock()
		return nil
	case <-wp.ctx.Done():
		return fmt.Errorf("worker pool foi fechado")
	default:
		return fmt.Errorf("queue de tarefas está cheia")
	}
}

// ProcessBatch processa um lote de tarefas
func (wp *WorkerPool) ProcessBatch(tasks []Task) ([]Result, []error) {
	var results []Result
	var errors []error
	var resultMutex sync.Mutex

	// Submete todas as tarefas
	for _, task := range tasks {
		if err := wp.SubmitTask(task); err != nil {
			errors = append(errors, err)
		}
	}

	// Coleta resultados
	expectedResults := len(tasks) - len(errors)
	received := 0

	for received < expectedResults {
		select {
		case result := <-wp.resultQueue:
			resultMutex.Lock()
			results = append(results, result)
			received++
			resultMutex.Unlock()

		case err := <-wp.errorQueue:
			resultMutex.Lock()
			errors = append(errors, err)
			received++
			resultMutex.Unlock()

		case <-time.After(30 * time.Second):
			// Timeout para evitar espera infinita
			return results, append(errors, fmt.Errorf("timeout aguardando resultados"))
		}
	}

	return results, errors
}

// ProcessBatchWithContext processa lote com context
func (wp *WorkerPool) ProcessBatchWithContext(ctx context.Context, tasks []Task) ([]Result, []error) {
	var results []Result
	var errors []error
	var resultMutex sync.Mutex

	// Canal para sincronização
	done := make(chan struct{})
	var wg sync.WaitGroup

	// Submete tarefas em goroutine separada
	go func() {
		defer close(done)
		for _, task := range tasks {
			select {
			case <-ctx.Done():
				resultMutex.Lock()
				errors = append(errors, ctx.Err())
				resultMutex.Unlock()
				return
			default:
				if err := wp.SubmitTask(task); err != nil {
					resultMutex.Lock()
					errors = append(errors, err)
					resultMutex.Unlock()
				} else {
					wg.Add(1)
				}
			}
		}
	}()

	// Coleta resultados
	go func() {
		for {
			select {
			case result := <-wp.resultQueue:
				resultMutex.Lock()
				results = append(results, result)
				resultMutex.Unlock()
				wg.Done()

			case err := <-wp.errorQueue:
				resultMutex.Lock()
				errors = append(errors, err)
				resultMutex.Unlock()
				wg.Done()

			case <-ctx.Done():
				return

			case <-done:
				return
			}
		}
	}()

	// Aguarda conclusão ou timeout
	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		return results, errors
	case <-ctx.Done():
		return results, append(errors, ctx.Err())
	}
}

// GetStats retorna estatísticas atuais do pool
func (wp *WorkerPool) GetStats() PoolStats {
	wp.mutex.RLock()
	defer wp.mutex.RUnlock()

	stats := *wp.stats
	stats.QueuedTasks = len(wp.taskQueue)
	return stats
}

// Close fecha o worker pool
func (wp *WorkerPool) Close() error {
	wp.cancel()

	// Fecha task queue
	close(wp.taskQueue)

	// Aguarda workers terminarem
	wp.wg.Wait()

	// Fecha result e error queues
	close(wp.resultQueue)
	close(wp.errorQueue)

	return nil
}

// Resize altera o número de workers (experimental)
func (wp *WorkerPool) Resize(newSize int) error {
	if newSize <= 0 {
		return fmt.Errorf("tamanho do pool deve ser positivo")
	}

	wp.mutex.Lock()
	currentSize := wp.workers
	wp.mutex.Unlock()

	if newSize > currentSize {
		// Adiciona workers
		for i := currentSize; i < newSize; i++ {
			wp.wg.Add(1)
			go wp.worker(i)
		}
	} else if newSize < currentSize {
		// Remove workers (fecha context)
		// Nota: implementação simplificada
		wp.cancel()
		wp.ctx, wp.cancel = context.WithCancel(context.Background())
	}

	wp.mutex.Lock()
	wp.workers = newSize
	wp.mutex.Unlock()

	return nil
}

// statsCollector coleta estatísticas periodicamente
func (wp *WorkerPool) statsCollector() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			wp.mutex.Lock()
			wp.stats.QueuedTasks = len(wp.taskQueue)
			wp.mutex.Unlock()

		case <-wp.ctx.Done():
			return
		}
	}
}

// WaitForCompletion aguarda até que todas as tarefas sejam processadas
func (wp *WorkerPool) WaitForCompletion(timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timer.C:
			return fmt.Errorf("timeout aguardando conclusão das tarefas")
		case <-ticker.C:
			if len(wp.taskQueue) == 0 {
				return nil
			}
		case <-wp.ctx.Done():
			return wp.ctx.Err()
		}
	}
}