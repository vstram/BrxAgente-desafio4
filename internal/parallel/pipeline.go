package parallel

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// PipelineStage representa um estágio do pipeline
type PipelineStage interface {
	Process(ctx context.Context, input <-chan interface{}) <-chan interface{}
	GetName() string
	GetStats() StageStats
}

// StageStats estatísticas de um estágio
type StageStats struct {
	Name            string
	ItemsProcessed  int64
	AverageDuration time.Duration
	ErrorCount      int64
	LastProcessed   time.Time
}

// ProcessingPipeline implementa pipeline de processamento otimizado
type ProcessingPipeline struct {
	stages     []PipelineStage
	bufferSize int
	ctx        context.Context
	cancel     context.CancelFunc
	stats      *PipelineStats
	mutex      sync.RWMutex
}

// PipelineStats estatísticas do pipeline
type PipelineStats struct {
	TotalItemsProcessed int64
	StartTime           time.Time
	LastItemTime        time.Time
	StageStats          map[string]StageStats
	Throughput          float64 // items por segundo
}

// NewProcessingPipeline cria um novo pipeline
func NewProcessingPipeline(bufferSize int) *ProcessingPipeline {
	ctx, cancel := context.WithCancel(context.Background())

	return &ProcessingPipeline{
		stages:     make([]PipelineStage, 0),
		bufferSize: bufferSize,
		ctx:        ctx,
		cancel:     cancel,
		stats: &PipelineStats{
			StartTime:  time.Now(),
			StageStats: make(map[string]StageStats),
		},
	}
}

// AddStage adiciona um estágio ao pipeline
func (pp *ProcessingPipeline) AddStage(stage PipelineStage) {
	pp.mutex.Lock()
	defer pp.mutex.Unlock()

	pp.stages = append(pp.stages, stage)
	pp.stats.StageStats[stage.GetName()] = StageStats{
		Name: stage.GetName(),
	}
}

// Process executa o pipeline completo
func (pp *ProcessingPipeline) Process(input <-chan interface{}) <-chan interface{} {
	if len(pp.stages) == 0 {
		// Se não há estágios, retorna input inalterado
		return input
	}

	// Conecta estágios sequencialmente
	current := input
	for i, stage := range pp.stages {
		// Cria canal intermediário com buffer
		next := make(chan interface{}, pp.bufferSize)

		// Inicia processamento do estágio atual
		go pp.runStage(stage, current, next, i == len(pp.stages)-1)

		current = next
	}

	return current
}

// runStage executa um estágio específico com monitoramento
func (pp *ProcessingPipeline) runStage(stage PipelineStage, input <-chan interface{}, output chan interface{}, isLast bool) {
	defer func() {
		if isLast {
			close(output)
		}
	}()

	stageOutput := stage.Process(pp.ctx, input)

	for {
		select {
		case item, ok := <-stageOutput:
			if !ok {
				return
			}

			pp.updateStats(stage.GetName())

			select {
			case output <- item:
			case <-pp.ctx.Done():
				return
			}

		case <-pp.ctx.Done():
			return
		}
	}
}

// ProcessBatch processa um lote de dados através do pipeline
func (pp *ProcessingPipeline) ProcessBatch(data []interface{}) ([]interface{}, error) {
	input := make(chan interface{}, len(data))

	// Envia dados para o pipeline
	go func() {
		defer close(input)
		for _, item := range data {
			select {
			case input <- item:
			case <-pp.ctx.Done():
				return
			}
		}
	}()

	// Processa através do pipeline
	output := pp.Process(input)

	// Coleta resultados
	var results []interface{}
	timeout := time.NewTimer(30 * time.Second)
	defer timeout.Stop()

	for {
		select {
		case result, ok := <-output:
			if !ok {
				return results, nil
			}
			results = append(results, result)

		case <-timeout.C:
			return results, fmt.Errorf("timeout processando lote")

		case <-pp.ctx.Done():
			return results, pp.ctx.Err()
		}
	}
}

// updateStats atualiza estatísticas do pipeline
func (pp *ProcessingPipeline) updateStats(stageName string) {
	pp.mutex.Lock()
	defer pp.mutex.Unlock()

	pp.stats.TotalItemsProcessed++
	pp.stats.LastItemTime = time.Now()

	// Calcula throughput (items por segundo)
	elapsed := time.Since(pp.stats.StartTime)
	if elapsed.Seconds() > 0 {
		pp.stats.Throughput = float64(pp.stats.TotalItemsProcessed) / elapsed.Seconds()
	}

	// Atualiza stats do estágio
	if stageStats, exists := pp.stats.StageStats[stageName]; exists {
		stageStats.ItemsProcessed++
		stageStats.LastProcessed = time.Now()
		pp.stats.StageStats[stageName] = stageStats
	}
}

// GetStats retorna estatísticas atuais do pipeline
func (pp *ProcessingPipeline) GetStats() PipelineStats {
	pp.mutex.RLock()
	defer pp.mutex.RUnlock()

	// Cria cópia das estatísticas
	stats := *pp.stats
	stats.StageStats = make(map[string]StageStats)

	for name, stageStats := range pp.stats.StageStats {
		stats.StageStats[name] = stageStats
	}

	return stats
}

// Close fecha o pipeline
func (pp *ProcessingPipeline) Close() error {
	pp.cancel()
	return nil
}

// FilterStage estágio que filtra dados
type FilterStage struct {
	name      string
	predicate func(interface{}) bool
	stats     StageStats
	mutex     sync.Mutex
}

// NewFilterStage cria um novo estágio de filtro
func NewFilterStage(name string, predicate func(interface{}) bool) *FilterStage {
	return &FilterStage{
		name:      name,
		predicate: predicate,
		stats:     StageStats{Name: name},
	}
}

func (fs *FilterStage) Process(ctx context.Context, input <-chan interface{}) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)

		for {
			select {
			case item, ok := <-input:
				if !ok {
					return
				}

				start := time.Now()

				if fs.predicate(item) {
					select {
					case output <- item:
					case <-ctx.Done():
						return
					}
				}

				fs.updateStats(time.Since(start))

			case <-ctx.Done():
				return
			}
		}
	}()

	return output
}

func (fs *FilterStage) GetName() string {
	return fs.name
}

func (fs *FilterStage) GetStats() StageStats {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	return fs.stats
}

func (fs *FilterStage) updateStats(duration time.Duration) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	fs.stats.ItemsProcessed++
	fs.stats.LastProcessed = time.Now()

	// Atualiza média de duração
	if fs.stats.ItemsProcessed == 1 {
		fs.stats.AverageDuration = duration
	} else {
		// Média móvel simples
		fs.stats.AverageDuration = time.Duration(
			(int64(fs.stats.AverageDuration) + int64(duration)) / 2,
		)
	}
}

// TransformStage estágio que transforma dados
type TransformStage struct {
	name        string
	transformer func(interface{}) (interface{}, error)
	stats       StageStats
	mutex       sync.Mutex
}

// NewTransformStage cria um novo estágio de transformação
func NewTransformStage(name string, transformer func(interface{}) (interface{}, error)) *TransformStage {
	return &TransformStage{
		name:        name,
		transformer: transformer,
		stats:       StageStats{Name: name},
	}
}

func (ts *TransformStage) Process(ctx context.Context, input <-chan interface{}) <-chan interface{} {
	output := make(chan interface{})

	go func() {
		defer close(output)

		for {
			select {
			case item, ok := <-input:
				if !ok {
					return
				}

				start := time.Now()

				transformed, err := ts.transformer(item)
				duration := time.Since(start)

				if err != nil {
					ts.updateStatsWithError(duration)
					continue
				}

				select {
				case output <- transformed:
				case <-ctx.Done():
					return
				}

				ts.updateStats(duration)

			case <-ctx.Done():
				return
			}
		}
	}()

	return output
}

func (ts *TransformStage) GetName() string {
	return ts.name
}

func (ts *TransformStage) GetStats() StageStats {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	return ts.stats
}

func (ts *TransformStage) updateStats(duration time.Duration) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	ts.stats.ItemsProcessed++
	ts.stats.LastProcessed = time.Now()

	if ts.stats.ItemsProcessed == 1 {
		ts.stats.AverageDuration = duration
	} else {
		ts.stats.AverageDuration = time.Duration(
			(int64(ts.stats.AverageDuration) + int64(duration)) / 2,
		)
	}
}

func (ts *TransformStage) updateStatsWithError(_ time.Duration) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	ts.stats.ErrorCount++
	ts.stats.LastProcessed = time.Now()
}
