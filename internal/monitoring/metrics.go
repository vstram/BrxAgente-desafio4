// Package monitoring implementa monitoramento e métricas de performance
package monitoring

import (
	"runtime"
	"sync"
	"time"
)

// PerformanceMetrics contém métricas de performance do sistema
type PerformanceMetrics struct {
	// Métricas de processamento
	TotalProcessingTime   time.Duration
	AverageProcessingTime time.Duration
	ItemsProcessed        int64
	ItemsPerSecond        float64
	ErrorCount            int64
	ErrorRate             float64

	// Métricas de sistema
	MemoryUsage    MemoryMetrics
	CPUUsage       CPUMetrics
	GoroutineCount int

	// Métricas de cache
	CacheHitRatio float64
	CacheSize     int64
	CacheMisses   int64

	// Timestamps
	StartTime    time.Time
	LastUpdate   time.Time
	Measurements int64
}

// MemoryMetrics métricas de memória
type MemoryMetrics struct {
	AllocMB      float64 // Memória alocada em MB
	TotalAllocMB float64 // Total de memória alocada em MB
	SysMB        float64 // Memória do sistema em MB
	GCCycles     uint32  // Número de ciclos de GC
	NextGCMB     float64 // Próximo GC em MB
}

// CPUMetrics métricas de CPU (simplificado)
type CPUMetrics struct {
	Usage         float64 // % de uso da CPU (estimado)
	GoroutineTime time.Duration
}

// MetricsCollector coletor de métricas
type MetricsCollector struct {
	metrics *PerformanceMetrics
	mutex   sync.RWMutex
	ticker  *time.Ticker
	done    chan bool
}

// NewMetricsCollector cria um novo coletor de métricas
func NewMetricsCollector(intervalSeconds int) *MetricsCollector {
	collector := &MetricsCollector{
		metrics: &PerformanceMetrics{
			StartTime: time.Now(),
		},
		done: make(chan bool),
	}

	// Inicia coleta periódica
	if intervalSeconds > 0 {
		collector.ticker = time.NewTicker(time.Duration(intervalSeconds) * time.Second)
		go collector.collectPeriodically()
	}

	return collector
}

// RecordProcessingTime registra tempo de processamento
func (mc *MetricsCollector) RecordProcessingTime(duration time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics.ItemsProcessed++
	mc.metrics.TotalProcessingTime += duration
	mc.metrics.LastUpdate = time.Now()

	// Calcula média
	mc.metrics.AverageProcessingTime = time.Duration(
		int64(mc.metrics.TotalProcessingTime) / mc.metrics.ItemsProcessed,
	)

	// Calcula items por segundo
	elapsed := time.Since(mc.metrics.StartTime)
	if elapsed.Seconds() > 0 {
		mc.metrics.ItemsPerSecond = float64(mc.metrics.ItemsProcessed) / elapsed.Seconds()
	}
}

// RecordError registra um erro
func (mc *MetricsCollector) RecordError() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics.ErrorCount++
	mc.metrics.LastUpdate = time.Now()

	// Calcula taxa de erro
	if mc.metrics.ItemsProcessed > 0 {
		mc.metrics.ErrorRate = float64(mc.metrics.ErrorCount) / float64(mc.metrics.ItemsProcessed) * 100
	}
}

// UpdateCacheMetrics atualiza métricas de cache
func (mc *MetricsCollector) UpdateCacheMetrics(hitRatio float64, size int64, misses int64) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics.CacheHitRatio = hitRatio
	mc.metrics.CacheSize = size
	mc.metrics.CacheMisses = misses
	mc.metrics.LastUpdate = time.Now()
}

// CollectSystemMetrics coleta métricas do sistema
func (mc *MetricsCollector) CollectSystemMetrics() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Atualiza métricas de memória
	mc.metrics.MemoryUsage = MemoryMetrics{
		AllocMB:      float64(m.Alloc) / 1024 / 1024,
		TotalAllocMB: float64(m.TotalAlloc) / 1024 / 1024,
		SysMB:        float64(m.Sys) / 1024 / 1024,
		GCCycles:     m.NumGC,
		NextGCMB:     float64(m.NextGC) / 1024 / 1024,
	}

	// Conta goroutines
	mc.metrics.GoroutineCount = runtime.NumGoroutine()

	// Estimativa simples de uso de CPU baseada em goroutines
	// (implementação mais sofisticada requereria bibliotecas específicas)
	mc.metrics.CPUUsage.Usage = float64(mc.metrics.GoroutineCount) / float64(runtime.NumCPU()) * 10
	if mc.metrics.CPUUsage.Usage > 100 {
		mc.metrics.CPUUsage.Usage = 100
	}

	mc.metrics.LastUpdate = time.Now()
	mc.metrics.Measurements++
}

// GetMetrics retorna cópia das métricas atuais
func (mc *MetricsCollector) GetMetrics() PerformanceMetrics {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	// Cria cópia das métricas
	metrics := PerformanceMetrics{
		TotalProcessingTime:   mc.metrics.TotalProcessingTime,
		AverageProcessingTime: mc.metrics.AverageProcessingTime,
		ItemsProcessed:        mc.metrics.ItemsProcessed,
		ItemsPerSecond:        mc.metrics.ItemsPerSecond,
		ErrorCount:            mc.metrics.ErrorCount,
		ErrorRate:             mc.metrics.ErrorRate,
		MemoryUsage:           mc.metrics.MemoryUsage,
		CPUUsage:              mc.metrics.CPUUsage,
		GoroutineCount:        mc.metrics.GoroutineCount,
		CacheHitRatio:         mc.metrics.CacheHitRatio,
		CacheSize:             mc.metrics.CacheSize,
		CacheMisses:           mc.metrics.CacheMisses,
		StartTime:             mc.metrics.StartTime,
		LastUpdate:            mc.metrics.LastUpdate,
		Measurements:          mc.metrics.Measurements,
	}
	return metrics
}

// GetSummary retorna resumo das métricas
func (mc *MetricsCollector) GetSummary() MetricsSummary {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	uptime := time.Since(mc.metrics.StartTime)

	return MetricsSummary{
		Uptime:                uptime,
		ItemsProcessed:        mc.metrics.ItemsProcessed,
		AverageProcessingTime: mc.metrics.AverageProcessingTime,
		ItemsPerSecond:        mc.metrics.ItemsPerSecond,
		ErrorRate:             mc.metrics.ErrorRate,
		MemoryUsageMB:         mc.metrics.MemoryUsage.AllocMB,
		CacheHitRatio:         mc.metrics.CacheHitRatio,
		GoroutineCount:        mc.metrics.GoroutineCount,
		CPUUsagePercent:       mc.metrics.CPUUsage.Usage,
		Measurements:          mc.metrics.Measurements,
	}
}

// MetricsSummary resumo das métricas principais
type MetricsSummary struct {
	Uptime                time.Duration
	ItemsProcessed        int64
	AverageProcessingTime time.Duration
	ItemsPerSecond        float64
	ErrorRate             float64
	MemoryUsageMB         float64
	CacheHitRatio         float64
	GoroutineCount        int
	CPUUsagePercent       float64
	Measurements          int64
}

// ResetMetrics reseta as métricas
func (mc *MetricsCollector) ResetMetrics() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics = &PerformanceMetrics{
		StartTime: time.Now(),
	}
}

// Stop para o coletor de métricas
func (mc *MetricsCollector) Stop() {
	if mc.ticker != nil {
		mc.ticker.Stop()
	}
	close(mc.done)
}

// collectPeriodically coleta métricas periodicamente
func (mc *MetricsCollector) collectPeriodically() {
	for {
		select {
		case <-mc.ticker.C:
			mc.CollectSystemMetrics()
		case <-mc.done:
			return
		}
	}
}

// GetPerformanceReport gera relatório de performance
func (mc *MetricsCollector) GetPerformanceReport() PerformanceReport {
	summary := mc.GetSummary()
	metrics := mc.GetMetrics()

	report := PerformanceReport{
		Timestamp:       time.Now(),
		Summary:         summary,
		Details:         metrics,
		Status:          "healthy",
		Warnings:        []string{},
		Recommendations: []string{},
	}

	// Analisa métricas e adiciona avisos/recomendações
	if summary.MemoryUsageMB > 500 { // 500MB
		report.Status = "warning"
		report.Warnings = append(report.Warnings, "Alto uso de memória")
		report.Recommendations = append(report.Recommendations, "Considere otimizar uso de memória")
	}

	if summary.ErrorRate > 5 { // 5%
		report.Status = "warning"
		report.Warnings = append(report.Warnings, "Taxa de erro alta")
		report.Recommendations = append(report.Recommendations, "Investigue causas dos erros")
	}

	if summary.CPUUsagePercent > 80 {
		report.Status = "warning"
		report.Warnings = append(report.Warnings, "Alto uso de CPU")
		report.Recommendations = append(report.Recommendations, "Considere otimizar processamento")
	}

	if summary.CacheHitRatio < 0.7 && summary.ItemsProcessed > 100 {
		report.Warnings = append(report.Warnings, "Baixa taxa de acerto do cache")
		report.Recommendations = append(report.Recommendations, "Otimize configuração do cache")
	}

	if len(report.Warnings) > 2 {
		report.Status = "critical"
	}

	return report
}

// PerformanceReport relatório de performance
type PerformanceReport struct {
	Timestamp       time.Time
	Summary         MetricsSummary
	Details         PerformanceMetrics
	Status          string // "healthy", "warning", "critical"
	Warnings        []string
	Recommendations []string
}

// Benchmark executa benchmark de performance
type Benchmark struct {
	Name       string
	Function   func() error
	Iterations int
	Warmup     int
	metrics    *MetricsCollector
}

// NewBenchmark cria um novo benchmark
func NewBenchmark(name string, function func() error, iterations int) *Benchmark {
	return &Benchmark{
		Name:       name,
		Function:   function,
		Iterations: iterations,
		Warmup:     iterations / 10,        // 10% para warmup
		metrics:    NewMetricsCollector(0), // Sem coleta periódica
	}
}

// Run executa o benchmark
func (b *Benchmark) Run() BenchmarkResult {
	result := BenchmarkResult{
		Name:       b.Name,
		Iterations: b.Iterations,
		StartTime:  time.Now(),
	}

	var totalDuration time.Duration
	var errorCount int

	// Warmup
	for i := 0; i < b.Warmup; i++ {
		b.Function()
	}

	// Execução real
	for i := 0; i < b.Iterations; i++ {
		start := time.Now()
		err := b.Function()
		duration := time.Since(start)

		totalDuration += duration
		b.metrics.RecordProcessingTime(duration)

		if err != nil {
			errorCount++
			b.metrics.RecordError()
		}
	}

	result.EndTime = time.Now()
	result.TotalDuration = totalDuration
	result.AverageDuration = totalDuration / time.Duration(b.Iterations)
	result.ErrorCount = errorCount
	result.ErrorRate = float64(errorCount) / float64(b.Iterations) * 100
	result.Throughput = float64(b.Iterations) / result.EndTime.Sub(result.StartTime).Seconds()

	// Coleta métricas finais
	b.metrics.CollectSystemMetrics()
	result.FinalMetrics = b.metrics.GetMetrics()

	return result
}

// BenchmarkResult resultado do benchmark
type BenchmarkResult struct {
	Name            string
	Iterations      int
	StartTime       time.Time
	EndTime         time.Time
	TotalDuration   time.Duration
	AverageDuration time.Duration
	ErrorCount      int
	ErrorRate       float64
	Throughput      float64 // operações por segundo
	FinalMetrics    PerformanceMetrics
}
