// Package performance gerencia todas as otimizações de performance do sistema
package performance

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"BrxAgente-desafio4/internal/cache"
	"BrxAgente-desafio4/internal/monitoring"
	"BrxAgente-desafio4/internal/parallel"
)

// PerformanceManager gerencia todos os aspectos de performance do sistema
type PerformanceManager struct {
	// Componentes core
	smartCache       *cache.SmartCache
	metricsCollector *monitoring.MetricsCollector
	alertManager     *monitoring.AlertManager
	profiler         *monitoring.AutoProfiler
	autoTrigger      *monitoring.AutoTrigger

	// Processamento paralelo
	workerPool     *parallel.WorkerPool
	batchProcessor *parallel.BatchProcessor
	pipeline       *parallel.ProcessingPipeline

	// Configuração e estado
	config  PerformanceConfig
	running bool
	stopCh  chan bool
	mutex   sync.RWMutex
}

// PerformanceConfig configuração completa de performance
type PerformanceConfig struct {
	// Cache
	CacheConfig cache.SmartCacheConfig

	// Processamento Paralelo
	WorkerPoolSize int
	BatchSize      int
	MaxConcurrent  int
	PipelineBuffer int

	// Monitoramento
	MetricsInterval    time.Duration
	ProfilingEnabled   bool
	ProfilingOutputDir string
	AlertsEnabled      bool

	// Otimização automática
	AutoOptimization     bool
	OptimizationInterval time.Duration

	// Thresholds
	CPUThreshold    float64
	MemoryThreshold float64
}

// DefaultPerformanceConfig retorna configuração padrão
func DefaultPerformanceConfig() PerformanceConfig {
	return PerformanceConfig{
		CacheConfig: cache.SmartCacheConfig{
			LLMMaxSize:       1000,
			LLMTTLHours:      24,
			DataMaxSizeBytes: 512 * 1024 * 1024, // 512MB
			DataTTLHours:     12,
		},
		WorkerPoolSize:       runtime.NumCPU() * 2,
		BatchSize:            100,
		MaxConcurrent:        runtime.NumCPU(),
		PipelineBuffer:       200,
		MetricsInterval:      10 * time.Second,
		ProfilingEnabled:     false,
		ProfilingOutputDir:   "./profiles",
		AlertsEnabled:        true,
		AutoOptimization:     true,
		OptimizationInterval: 15 * time.Minute,
		CPUThreshold:         80.0,
		MemoryThreshold:      500.0,
	}
}

// NewPerformanceManager cria um novo gerenciador de performance
func NewPerformanceManager(config PerformanceConfig) *PerformanceManager {
	pm := &PerformanceManager{
		config: config,
		stopCh: make(chan bool),
	}

	pm.initializeComponents()
	return pm
}

// initializeComponents inicializa todos os componentes
func (pm *PerformanceManager) initializeComponents() {
	// Cache inteligente
	pm.smartCache = cache.NewSmartCache(pm.config.CacheConfig)

	// Métricas
	pm.metricsCollector = monitoring.NewMetricsCollector(int(pm.config.MetricsInterval.Seconds()))

	// Alertas
	pm.alertManager = monitoring.NewAlertManager()
	if !pm.config.AlertsEnabled {
		pm.alertManager.Disable()
	}

	// Profiler automático
	profilerConfig := monitoring.ProfilerConfig{
		Enabled:            pm.config.ProfilingEnabled,
		OutputDir:          pm.config.ProfilingOutputDir,
		MemProfileInterval: 10 * time.Minute,
		AutoCPUThreshold:   pm.config.CPUThreshold,
		AutoMemThreshold:   pm.config.MemoryThreshold,
	}
	pm.profiler = monitoring.NewAutoProfiler(profilerConfig)

	// Auto trigger para profiling baseado em métricas
	pm.autoTrigger = monitoring.NewAutoTrigger(
		pm.profiler,
		pm.metricsCollector,
		pm.config.CPUThreshold,
		pm.config.MemoryThreshold,
	)

	// Worker pool
	pm.workerPool = parallel.NewWorkerPool(pm.config.WorkerPoolSize, pm.config.BatchSize*2)

	// Batch processor
	batchConfig := parallel.BatchConfig{
		WorkerPoolSize: pm.config.WorkerPoolSize,
		BatchSize:      pm.config.BatchSize,
		MaxConcurrent:  pm.config.MaxConcurrent,
		TimeoutSeconds: 300, // 5 minutos
	}
	pm.batchProcessor = parallel.NewBatchProcessor(batchConfig)

	// Pipeline de processamento
	pm.pipeline = parallel.NewProcessingPipeline(pm.config.PipelineBuffer)

	// Handlers de alerta padrão
	logHandler := monitoring.NewLogHandler()
	pm.alertManager.AddHandler(logHandler)
}

// Start inicia todos os componentes de performance
func (pm *PerformanceManager) Start() error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if pm.running {
		return fmt.Errorf("performance manager já está executando")
	}

	// Inicia profiler se habilitado
	if pm.config.ProfilingEnabled {
		if err := pm.profiler.StartContinuousProfiling(); err != nil {
			return fmt.Errorf("erro iniciando profiler: %w", err)
		}
		pm.autoTrigger.StartMonitoring()
	}

	// Inicia loop de otimização automática
	if pm.config.AutoOptimization {
		go pm.optimizationLoop()
	}

	// Inicia loop de monitoramento
	go pm.monitoringLoop()

	pm.running = true
	return nil
}

// Stop para todos os componentes
func (pm *PerformanceManager) Stop() error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if !pm.running {
		return nil
	}

	// Sinaliza parada
	close(pm.stopCh)

	// Para componentes
	pm.metricsCollector.Stop()
	pm.profiler.Stop()
	pm.autoTrigger.Stop()
	pm.workerPool.Close()
	pm.batchProcessor.Close()
	pm.pipeline.Close()

	pm.running = false
	return nil
}

// GetSmartCache retorna o cache inteligente
func (pm *PerformanceManager) GetSmartCache() *cache.SmartCache {
	return pm.smartCache
}

// GetBatchProcessor retorna o processador de lotes
func (pm *PerformanceManager) GetBatchProcessor() *parallel.BatchProcessor {
	return pm.batchProcessor
}

// GetWorkerPool retorna o worker pool
func (pm *PerformanceManager) GetWorkerPool() *parallel.WorkerPool {
	return pm.workerPool
}

// GetPipeline retorna o pipeline de processamento
func (pm *PerformanceManager) GetPipeline() *parallel.ProcessingPipeline {
	return pm.pipeline
}

// GetMetrics retorna métricas atuais
func (pm *PerformanceManager) GetMetrics() monitoring.PerformanceMetrics {
	return pm.metricsCollector.GetMetrics()
}

// GetPerformanceReport gera relatório completo de performance
func (pm *PerformanceManager) GetPerformanceReport() PerformanceReport {
	metricsReport := pm.metricsCollector.GetPerformanceReport()
	cacheStats := pm.smartCache.GetCombinedStats()
	alertStats := pm.alertManager.GetAlertStats()
	workerStats := pm.workerPool.GetStats()
	batchStats := pm.batchProcessor.GetStats()

	return PerformanceReport{
		Timestamp:       time.Now(),
		MetricsReport:   metricsReport,
		CacheStats:      cacheStats,
		AlertStats:      alertStats,
		WorkerPoolStats: workerStats,
		BatchStats:      batchStats,
		SystemHealth:    pm.assessSystemHealth(),
		Recommendations: pm.generateRecommendations(),
	}
}

// PerformanceReport relatório completo de performance
type PerformanceReport struct {
	Timestamp       time.Time
	MetricsReport   monitoring.PerformanceReport
	CacheStats      cache.CombinedCacheStats
	AlertStats      monitoring.AlertStats
	WorkerPoolStats parallel.PoolStats
	BatchStats      parallel.BatchProcessorStats
	SystemHealth    SystemHealth
	Recommendations []string
}

// SystemHealth avaliação da saúde do sistema
type SystemHealth struct {
	Overall    string // "healthy", "warning", "critical"
	CPU        string
	Memory     string
	Cache      string
	Processing string
	Errors     string
}

// assessSystemHealth avalia a saúde geral do sistema
func (pm *PerformanceManager) assessSystemHealth() SystemHealth {
	metrics := pm.metricsCollector.GetMetrics()
	cacheHealth := pm.smartCache.GetCacheHealth()

	health := SystemHealth{
		Overall: "healthy",
	}

	// Avalia CPU
	if metrics.CPUUsage.Usage > 85 {
		health.CPU = "critical"
		health.Overall = "critical"
	} else if metrics.CPUUsage.Usage > 70 {
		health.CPU = "warning"
		if health.Overall == "healthy" {
			health.Overall = "warning"
		}
	} else {
		health.CPU = "healthy"
	}

	// Avalia Memória
	if metrics.MemoryUsage.AllocMB > 800 {
		health.Memory = "critical"
		health.Overall = "critical"
	} else if metrics.MemoryUsage.AllocMB > 500 {
		health.Memory = "warning"
		if health.Overall == "healthy" {
			health.Overall = "warning"
		}
	} else {
		health.Memory = "healthy"
	}

	// Avalia Cache
	health.Cache = cacheHealth.Overall

	// Avalia Processamento
	if metrics.AverageProcessingTime > 15*time.Second {
		health.Processing = "critical"
		health.Overall = "critical"
	} else if metrics.AverageProcessingTime > 5*time.Second {
		health.Processing = "warning"
		if health.Overall == "healthy" {
			health.Overall = "warning"
		}
	} else {
		health.Processing = "healthy"
	}

	// Avalia Erros
	if metrics.ErrorRate > 5 {
		health.Errors = "critical"
		health.Overall = "critical"
	} else if metrics.ErrorRate > 2 {
		health.Errors = "warning"
		if health.Overall == "healthy" {
			health.Overall = "warning"
		}
	} else {
		health.Errors = "healthy"
	}

	return health
}

// generateRecommendations gera recomendações de otimização
func (pm *PerformanceManager) generateRecommendations() []string {
	var recommendations []string

	metrics := pm.metricsCollector.GetMetrics()
	cacheStats := pm.smartCache.GetCombinedStats()
	workerStats := pm.workerPool.GetStats()

	// Recomendações baseadas em métricas
	if metrics.CPUUsage.Usage > 80 {
		recommendations = append(recommendations, "Considere aumentar o número de workers ou otimizar algoritmos")
	}

	if metrics.MemoryUsage.AllocMB > 500 {
		recommendations = append(recommendations, "Considere reduzir tamanhos de cache ou implementar cleanup mais agressivo")
	}

	if cacheStats.LLMStats.HitRatio < 0.7 {
		recommendations = append(recommendations, "Otimize prompts para melhorar hit ratio do cache LLM")
		recommendations = append(recommendations, "Considere aumentar TTL do cache LLM")
	}

	if metrics.ErrorRate > 2 {
		recommendations = append(recommendations, "Investigue causas dos erros e implemente retry logic")
	}

	if metrics.AverageProcessingTime > 10*time.Second {
		recommendations = append(recommendations, "Implemente processamento paralelo para operações longas")
	}

	// Recomendações baseadas em utilização de workers
	if float64(workerStats.ActiveWorkers) > float64(pm.config.WorkerPoolSize)*0.8 {
		recommendations = append(recommendations, "Considere aumentar o tamanho do worker pool")
	}

	return recommendations
}

// optimizationLoop executa otimizações automáticas periodicamente
func (pm *PerformanceManager) optimizationLoop() {
	ticker := time.NewTicker(pm.config.OptimizationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pm.performAutoOptimization()
		case <-pm.stopCh:
			return
		}
	}
}

// monitoringLoop monitora métricas e gera alertas
func (pm *PerformanceManager) monitoringLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pm.checkAlertsAndMetrics()
		case <-pm.stopCh:
			return
		}
	}
}

// performAutoOptimization executa otimizações automáticas
func (pm *PerformanceManager) performAutoOptimization() {
	// Otimiza caches
	pm.smartCache.OptimizeCaches()

	// Coleta métricas do sistema
	pm.metricsCollector.CollectSystemMetrics()

	// Atualiza métricas de cache
	cacheStats := pm.smartCache.GetCombinedStats()
	pm.metricsCollector.UpdateCacheMetrics(
		cacheStats.LLMStats.HitRatio,
		cacheStats.DataStats.CurrentSize,
		cacheStats.LLMStats.MissCount,
	)
}

// checkAlertsAndMetrics verifica alertas baseado nas métricas atuais
func (pm *PerformanceManager) checkAlertsAndMetrics() {
	metrics := pm.metricsCollector.GetMetrics()
	pm.alertManager.CheckMetrics(metrics)
}

// RunStressTest executa teste de stress do sistema
func (pm *PerformanceManager) RunStressTest(duration time.Duration, concurrency int) StressTestResult {
	startTime := time.Now()

	result := StressTestResult{
		StartTime:      startTime,
		Duration:       duration,
		Concurrency:    concurrency,
		InitialMetrics: pm.metricsCollector.GetMetrics(),
	}

	// Cria snapshot inicial
	pm.profiler.CreateFullSnapshot()

	// Simula carga pesada
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			pm.stressWorker(ctx, workerID)
		}(i)
	}

	wg.Wait()

	// Coleta métricas finais
	result.FinalMetrics = pm.metricsCollector.GetMetrics()
	result.EndTime = time.Now()
	result.ActualDuration = result.EndTime.Sub(startTime)

	// Cria snapshot final
	pm.profiler.CreateFullSnapshot()

	return result
}

// stressWorker simula carga de trabalho
func (pm *PerformanceManager) stressWorker(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Simula operações típicas
			pm.simulateOperation(workerID)
			time.Sleep(time.Millisecond * 10)
		}
	}
}

// simulateOperation simula uma operação típica do sistema
func (pm *PerformanceManager) simulateOperation(workerID int) {
	start := time.Now()

	// Simula cache hit/miss
	prompt := fmt.Sprintf("worker_%d_prompt_%d", workerID, time.Now().UnixNano())
	if response, found := pm.smartCache.GetLLMResponse(prompt); !found {
		// Cache miss - simula processamento
		response = fmt.Sprintf("response_for_%s", prompt)
		metadata := cache.CacheMetadata{
			TokensUsed: 100,
			Duration:   time.Millisecond * 50,
			Model:      "stress-test",
			Quality:    0.9,
		}
		pm.smartCache.SetLLMResponse(prompt, response, metadata)
	}

	duration := time.Since(start)
	pm.metricsCollector.RecordProcessingTime(duration)

	// Simula erro ocasional (1%)
	if time.Now().UnixNano()%100 == 0 {
		pm.metricsCollector.RecordError()
	}
}

// StressTestResult resultado do teste de stress
type StressTestResult struct {
	StartTime      time.Time
	EndTime        time.Time
	Duration       time.Duration
	ActualDuration time.Duration
	Concurrency    int
	InitialMetrics monitoring.PerformanceMetrics
	FinalMetrics   monitoring.PerformanceMetrics
	Success        bool
}

// GetOptimizationReport gera relatório de otimização
func (pm *PerformanceManager) GetOptimizationReport() OptimizationReport {
	cacheHealth := pm.smartCache.GetCacheHealth()

	return OptimizationReport{
		Timestamp:            time.Now(),
		CacheHealth:          cacheHealth,
		SystemMetrics:        pm.metricsCollector.GetMetrics(),
		Recommendations:      pm.generateRecommendations(),
		OptimizationsApplied: pm.getRecentOptimizations(),
	}
}

// OptimizationReport relatório de otimização
type OptimizationReport struct {
	Timestamp            time.Time
	CacheHealth          cache.CacheHealthReport
	SystemMetrics        monitoring.PerformanceMetrics
	Recommendations      []string
	OptimizationsApplied []string
}

// getRecentOptimizations retorna otimizações aplicadas recentemente
func (pm *PerformanceManager) getRecentOptimizations() []string {
	return []string{
		"Cache LLM otimizado automaticamente",
		"Métricas do sistema coletadas",
		"Verificação de alertas executada",
	}
}
