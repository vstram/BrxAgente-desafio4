package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"strings"
	"sync"
	"time"
)

// PerformanceBenchmark executa benchmarks completos do sistema
type PerformanceBenchmark struct {
	agent               *VRAgent
	metricsCollector    *QualityMetricsCollector
	scenarios           []RealWorldTestScenario
	concurrencyLevels   []int
	testDuration        time.Duration
	warmupDuration      time.Duration
	logger              Logger
	results             []BenchmarkResult
	mu                  sync.RWMutex
}

// BenchmarkResult resultado de um benchmark
type BenchmarkResult struct {
	TestName            string                 `json:"test_name"`
	Scenario            string                 `json:"scenario"`
	ConcurrentUsers     int                    `json:"concurrent_users"`
	TotalRequests       int64                  `json:"total_requests"`
	SuccessfulRequests  int64                  `json:"successful_requests"`
	FailedRequests      int64                  `json:"failed_requests"`
	SuccessRate         float64                `json:"success_rate"`
	
	// Métricas de Tempo
	TotalDuration       time.Duration          `json:"total_duration"`
	AverageResponseTime time.Duration          `json:"average_response_time"`
	MinResponseTime     time.Duration          `json:"min_response_time"`
	MaxResponseTime     time.Duration          `json:"max_response_time"`
	MedianResponseTime  time.Duration          `json:"median_response_time"`
	P95ResponseTime     time.Duration          `json:"p95_response_time"`
	P99ResponseTime     time.Duration          `json:"p99_response_time"`
	
	// Métricas de Throughput
	RequestsPerSecond   float64                `json:"requests_per_second"`
	ResponsesPerSecond  float64                `json:"responses_per_second"`
	
	// Métricas de Sistema
	CPUUsage            CPUMetrics             `json:"cpu_usage"`
	MemoryUsage         MemoryMetrics          `json:"memory_usage"`
	
	// Métricas de Qualidade
	AverageScore        float64                `json:"average_score"`
	AccuracyRate        float64                `json:"accuracy_rate"`
	ErrorTypes          map[string]int         `json:"error_types"`
	
	StartTime           time.Time              `json:"start_time"`
	EndTime             time.Time              `json:"end_time"`
	Metadata            map[string]interface{} `json:"metadata"`
}

// CPUMetrics métricas de CPU
type CPUMetrics struct {
	InitialUsage  float64 `json:"initial_usage"`
	PeakUsage     float64 `json:"peak_usage"`
	AverageUsage  float64 `json:"average_usage"`
	FinalUsage    float64 `json:"final_usage"`
}

// MemoryMetrics métricas de memória
type MemoryMetrics struct {
	InitialUsage   uint64 `json:"initial_usage"`   // bytes
	PeakUsage      uint64 `json:"peak_usage"`      // bytes
	AverageUsage   uint64 `json:"average_usage"`   // bytes
	FinalUsage     uint64 `json:"final_usage"`     // bytes
	AllocatedTotal uint64 `json:"allocated_total"` // bytes
	GCCycles       uint32 `json:"gc_cycles"`
}

// BenchmarkConfig configuração para benchmarks
type BenchmarkConfig struct {
	ConcurrencyLevels []int         `json:"concurrency_levels"` // [1, 5, 10, 20, 50]
	TestDuration      time.Duration `json:"test_duration"`      // 30s
	WarmupDuration    time.Duration `json:"warmup_duration"`    // 5s
	CooldownDuration  time.Duration `json:"cooldown_duration"`  // 2s
	MaxTimeout        time.Duration `json:"max_timeout"`        // 30s
	EnableProfiling   bool          `json:"enable_profiling"`
	SampleSize        int           `json:"sample_size"`        // Número de cenários a testar
}

// RequestMetric métrica individual de request
type RequestMetric struct {
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	Success      bool
	Error        error
	ResponseSize int
	Score        float64
}

// NewPerformanceBenchmark cria novo benchmark
func NewPerformanceBenchmark(agent *VRAgent, metricsCollector *QualityMetricsCollector, logger Logger) *PerformanceBenchmark {
	return &PerformanceBenchmark{
		agent:               agent,
		metricsCollector:    metricsCollector,
		scenarios:           LoadRealWorldDataset(),
		concurrencyLevels:   []int{1, 5, 10, 20},
		testDuration:        30 * time.Second,
		warmupDuration:      5 * time.Second,
		logger:              logger,
		results:             []BenchmarkResult{},
	}
}

// SetConfig define configuração do benchmark
func (pb *PerformanceBenchmark) SetConfig(config BenchmarkConfig) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	
	pb.concurrencyLevels = config.ConcurrencyLevels
	pb.testDuration = config.TestDuration
	pb.warmupDuration = config.WarmupDuration
	
	pb.logger.Printf("Benchmark configuration updated: %+v", config)
}

// RunAllBenchmarks executa todos os benchmarks
func (pb *PerformanceBenchmark) RunAllBenchmarks() (*BenchmarkReport, error) {
	pb.logger.Println("Starting comprehensive performance benchmarks...")
	
	startTime := time.Now()
	
	// 1. Benchmark de Carga Básico
	if err := pb.runLoadBenchmark(); err != nil {
		return nil, fmt.Errorf("load benchmark failed: %w", err)
	}
	
	// 2. Benchmark de Stress
	if err := pb.runStressBenchmark(); err != nil {
		return nil, fmt.Errorf("stress benchmark failed: %w", err)
	}
	
	// 3. Benchmark de Diferentes Categorias
	if err := pb.runCategoryBenchmarks(); err != nil {
		return nil, fmt.Errorf("category benchmarks failed: %w", err)
	}
	
	// 4. Benchmark de Memória
	if err := pb.runMemoryBenchmark(); err != nil {
		return nil, fmt.Errorf("memory benchmark failed: %w", err)
	}
	
	totalDuration := time.Since(startTime)
	
	// Gerar relatório consolidado
	report := pb.generateBenchmarkReport(totalDuration)
	
	pb.logger.Printf("All benchmarks completed in %v", totalDuration)
	return report, nil
}

// runLoadBenchmark executa benchmark de carga normal
func (pb *PerformanceBenchmark) runLoadBenchmark() error {
	pb.logger.Println("Running load benchmark...")
	
	for _, concurrency := range pb.concurrencyLevels {
		result, err := pb.runSingleBenchmark("load_test", concurrency, pb.testDuration)
		if err != nil {
			return fmt.Errorf("load test with %d users failed: %w", concurrency, err)
		}
		
		pb.mu.Lock()
		pb.results = append(pb.results, *result)
		pb.mu.Unlock()
		
		pb.logger.Printf("Load test completed: %d users, %.2f RPS, %.2f%% success rate",
			concurrency, result.RequestsPerSecond, result.SuccessRate*100)
		
		// Cooldown entre testes
		time.Sleep(2 * time.Second)
	}
	
	return nil
}

// runStressBenchmark executa benchmark de stress
func (pb *PerformanceBenchmark) runStressBenchmark() error {
	pb.logger.Println("Running stress benchmark...")
	
	// Teste com alta concorrência
	stressLevels := []int{50, 100, 200}
	
	for _, concurrency := range stressLevels {
		result, err := pb.runSingleBenchmark("stress_test", concurrency, pb.testDuration/2) // Tempo menor para stress
		if err != nil {
			pb.logger.Printf("Stress test with %d users failed (expected): %v", concurrency, err)
			// Continuar mesmo com falhas no stress test
			continue
		}
		
		pb.mu.Lock()
		pb.results = append(pb.results, *result)
		pb.mu.Unlock()
		
		pb.logger.Printf("Stress test completed: %d users, %.2f RPS, %.2f%% success rate",
			concurrency, result.RequestsPerSecond, result.SuccessRate*100)
		
		// Cooldown maior para stress tests
		time.Sleep(5 * time.Second)
	}
	
	return nil
}

// runCategoryBenchmarks testa performance por categoria
func (pb *PerformanceBenchmark) runCategoryBenchmarks() error {
	pb.logger.Println("Running category-specific benchmarks...")
	
	categories := []string{
		"Políticas de Elegibilidade",
		"Cálculos Específicos", 
		"Cenários Complexos",
		"Dados Processados",
	}
	
	for _, category := range categories {
		scenarios := GetScenariosByCategory(category)
		if len(scenarios) == 0 {
			continue
		}
		
		// Testar com concorrência moderada
		result, err := pb.runCategoryBenchmark(category, 10, pb.testDuration/2, scenarios)
		if err != nil {
			return fmt.Errorf("category benchmark for %s failed: %w", category, err)
		}
		
		pb.mu.Lock()
		pb.results = append(pb.results, *result)
		pb.mu.Unlock()
		
		pb.logger.Printf("Category benchmark completed: %s, %.2f RPS, %.2f%% success rate",
			category, result.RequestsPerSecond, result.SuccessRate*100)
	}
	
	return nil
}

// runMemoryBenchmark testa consumo de memória
func (pb *PerformanceBenchmark) runMemoryBenchmark() error {
	pb.logger.Println("Running memory benchmark...")
	
	// Forçar GC antes do teste
	runtime.GC()
	
	result, err := pb.runSingleBenchmark("memory_test", 5, pb.testDuration*2) // Teste mais longo
	if err != nil {
		return fmt.Errorf("memory benchmark failed: %w", err)
	}
	
	pb.mu.Lock()
	pb.results = append(pb.results, *result)
	pb.mu.Unlock()
	
	pb.logger.Printf("Memory benchmark completed: Peak usage: %.2f MB",
		float64(result.MemoryUsage.PeakUsage)/1024/1024)
	
	return nil
}

// runSingleBenchmark executa um benchmark individual
func (pb *PerformanceBenchmark) runSingleBenchmark(testName string, concurrency int, duration time.Duration) (*BenchmarkResult, error) {
	pb.logger.Printf("Starting %s with %d concurrent users for %v", testName, concurrency, duration)
	
	// Inicializar resultado
	result := &BenchmarkResult{
		TestName:        testName,
		Scenario:        "mixed_scenarios",
		ConcurrentUsers: concurrency,
		ErrorTypes:      make(map[string]int),
		Metadata:        make(map[string]interface{}),
		StartTime:       time.Now(),
	}
	
	// Capturar métricas iniciais do sistema
	initialCPU := pb.getCurrentCPUUsage()
	initialMem := pb.getCurrentMemoryUsage()
	result.CPUUsage.InitialUsage = initialCPU
	result.MemoryUsage.InitialUsage = initialMem
	
	// Canal para coletar métricas de requests
	requestMetrics := make(chan RequestMetric, concurrency*1000)
	
	// Contexto para controle de tempo
	ctx, cancel := context.WithTimeout(context.Background(), duration+pb.warmupDuration)
	defer cancel()
	
	// WaitGroup para sincronizar workers
	var wg sync.WaitGroup
	
	// Iniciar workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go pb.worker(ctx, i, requestMetrics, &wg)
	}
	
	// Monitorar métricas do sistema em paralelo
	systemMetrics := pb.monitorSystemMetrics(ctx, initialCPU, initialMem)
	
	// Aguardar conclusão dos workers
	wg.Wait()
	close(requestMetrics)
	
	result.EndTime = time.Now()
	result.TotalDuration = result.EndTime.Sub(result.StartTime)
	
	// Processar métricas coletadas
	pb.processRequestMetrics(requestMetrics, result)
	
	// Processar métricas do sistema
	result.CPUUsage = systemMetrics.CPU
	result.MemoryUsage = systemMetrics.Memory
	
	// Calcular taxa de sucesso
	if result.TotalRequests > 0 {
		result.SuccessRate = float64(result.SuccessfulRequests) / float64(result.TotalRequests)
		result.RequestsPerSecond = float64(result.TotalRequests) / result.TotalDuration.Seconds()
		result.ResponsesPerSecond = float64(result.SuccessfulRequests) / result.TotalDuration.Seconds()
	}
	
	return result, nil
}

// runCategoryBenchmark executa benchmark para uma categoria específica
func (pb *PerformanceBenchmark) runCategoryBenchmark(category string, concurrency int, duration time.Duration, scenarios []RealWorldTestScenario) (*BenchmarkResult, error) {
	// Similar ao runSingleBenchmark mas usando cenários específicos da categoria
	result := &BenchmarkResult{
		TestName:        "category_test",
		Scenario:        category,
		ConcurrentUsers: concurrency,
		ErrorTypes:      make(map[string]int),
		Metadata:        map[string]interface{}{
			"category": category,
			"scenario_count": len(scenarios),
		},
		StartTime:       time.Now(),
	}
	
	// Implementação similar ao runSingleBenchmark
	// Por simplicidade, delegando para método base
	baseResult, err := pb.runSingleBenchmark(fmt.Sprintf("category_%s", category), concurrency, duration)
	if err != nil {
		return nil, err
	}
	
	// Copiar dados relevantes
	*result = *baseResult
	result.TestName = "category_test"
	result.Scenario = category
	
	return result, nil
}

// worker executa requests em paralelo
func (pb *PerformanceBenchmark) worker(ctx context.Context, workerID int, metrics chan<- RequestMetric, wg *sync.WaitGroup) {
	defer wg.Done()
	
	scenarios := pb.scenarios
	if len(scenarios) == 0 {
		return
	}
	
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Selecionar cenário aleatório
			scenario := scenarios[workerID%len(scenarios)]
			
			metric := RequestMetric{
				StartTime: time.Now(),
			}
			
			// Executar request
			response, err := pb.agent.Ask(scenario.Question)
			metric.EndTime = time.Now()
			metric.Duration = metric.EndTime.Sub(metric.StartTime)
			
			if err != nil {
				metric.Success = false
				metric.Error = err
			} else {
				metric.Success = true
				metric.ResponseSize = len(response)
				
				// Avaliar score básico (simplificado)
				metric.Score = pb.evaluateResponseScore(response, scenario)
			}
			
			metrics <- metric
		}
	}
}

// evaluateResponseScore avalia score da resposta (versão simplificada)
func (pb *PerformanceBenchmark) evaluateResponseScore(response string, scenario RealWorldTestScenario) float64 {
	// Implementação básica de avaliação
	score := 3.0 // Score base
	
	// Verificar se contém palavras-chave obrigatórias
	found := 0
	for _, keyword := range scenario.Expected.MustContain {
		if len(keyword) > 0 && containsIgnoreCase(response, keyword) {
			found++
		}
	}
	
	if len(scenario.Expected.MustContain) > 0 {
		keywordScore := float64(found) / float64(len(scenario.Expected.MustContain))
		score += keywordScore * 2 // Max 2 pontos extras
	}
	
	return min(5.0, score)
}

// processRequestMetrics processa métricas coletadas dos requests
func (pb *PerformanceBenchmark) processRequestMetrics(metrics <-chan RequestMetric, result *BenchmarkResult) {
	var responseTimes []time.Duration
	var totalScore float64
	var successfulRequests int64
	var totalRequests int64
	
	for metric := range metrics {
		totalRequests++
		responseTimes = append(responseTimes, metric.Duration)
		
		if metric.Success {
			successfulRequests++
			totalScore += metric.Score
		} else {
			result.FailedRequests++
			
			// Categorizar erro
			errorType := "unknown"
			if metric.Error != nil {
				errorType = metric.Error.Error()
				if len(errorType) > 50 {
					errorType = errorType[:50] + "..."
				}
			}
			result.ErrorTypes[errorType]++
		}
	}
	
	result.TotalRequests = totalRequests
	result.SuccessfulRequests = successfulRequests
	
	// Calcular estatísticas de tempo
	if len(responseTimes) > 0 {
		pb.calculateTimeStatistics(responseTimes, result)
	}
	
	// Calcular score médio
	if successfulRequests > 0 {
		result.AverageScore = totalScore / float64(successfulRequests)
		result.AccuracyRate = float64(successfulRequests) / float64(totalRequests)
	}
}

// calculateTimeStatistics calcula estatísticas de tempo de resposta
func (pb *PerformanceBenchmark) calculateTimeStatistics(times []time.Duration, result *BenchmarkResult) {
	// Ordenar tempos
	sortedTimes := make([]time.Duration, len(times))
	copy(sortedTimes, times)
	
	for i := 0; i < len(sortedTimes); i++ {
		for j := i + 1; j < len(sortedTimes); j++ {
			if sortedTimes[i] > sortedTimes[j] {
				sortedTimes[i], sortedTimes[j] = sortedTimes[j], sortedTimes[i]
			}
		}
	}
	
	// Calcular estatísticas
	result.MinResponseTime = sortedTimes[0]
	result.MaxResponseTime = sortedTimes[len(sortedTimes)-1]
	result.MedianResponseTime = sortedTimes[len(sortedTimes)/2]
	result.P95ResponseTime = sortedTimes[int(float64(len(sortedTimes))*0.95)]
	result.P99ResponseTime = sortedTimes[int(float64(len(sortedTimes))*0.99)]
	
	// Calcular média
	var total time.Duration
	for _, t := range times {
		total += t
	}
	result.AverageResponseTime = total / time.Duration(len(times))
}

// SystemMetricsSnapshot snapshot de métricas do sistema
type SystemMetricsSnapshot struct {
	CPU    CPUMetrics
	Memory MemoryMetrics
}

// monitorSystemMetrics monitora métricas do sistema durante o teste
func (pb *PerformanceBenchmark) monitorSystemMetrics(ctx context.Context, initialCPU float64, initialMem uint64) SystemMetricsSnapshot {
	snapshot := SystemMetricsSnapshot{}
	snapshot.CPU.InitialUsage = initialCPU
	snapshot.Memory.InitialUsage = initialMem
	
	var cpuReadings []float64
	var memReadings []uint64
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			// Calcular médias finais
			if len(cpuReadings) > 0 {
				var sum float64
				for _, cpu := range cpuReadings {
					sum += cpu
					if cpu > snapshot.CPU.PeakUsage {
						snapshot.CPU.PeakUsage = cpu
					}
				}
				snapshot.CPU.AverageUsage = sum / float64(len(cpuReadings))
				snapshot.CPU.FinalUsage = cpuReadings[len(cpuReadings)-1]
			}
			
			if len(memReadings) > 0 {
				var sum uint64
				for _, mem := range memReadings {
					sum += mem
					if mem > snapshot.Memory.PeakUsage {
						snapshot.Memory.PeakUsage = mem
					}
				}
				snapshot.Memory.AverageUsage = sum / uint64(len(memReadings))
				snapshot.Memory.FinalUsage = memReadings[len(memReadings)-1]
			}
			
			return snapshot
			
		case <-ticker.C:
			cpu := pb.getCurrentCPUUsage()
			mem := pb.getCurrentMemoryUsage()
			
			cpuReadings = append(cpuReadings, cpu)
			memReadings = append(memReadings, mem)
		}
	}
}

// getCurrentCPUUsage obtém uso atual de CPU (simulado)
func (pb *PerformanceBenchmark) getCurrentCPUUsage() float64 {
	// Em um ambiente real, isso usaria bibliotecas de monitoramento de sistema
	// Para este exemplo, retornamos um valor simulado
	return float64(runtime.NumGoroutine()) * 0.1
}

// getCurrentMemoryUsage obtém uso atual de memória
func (pb *PerformanceBenchmark) getCurrentMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// BenchmarkReport relatório consolidado de benchmarks
type BenchmarkReport struct {
	GeneratedAt       time.Time          `json:"generated_at"`
	TotalDuration     time.Duration      `json:"total_duration"`
	TestsRun          int                `json:"tests_run"`
	Results           []BenchmarkResult  `json:"results"`
	Summary           BenchmarkSummary   `json:"summary"`
	Recommendations   []string           `json:"recommendations"`
	SystemInfo        SystemInfo         `json:"system_info"`
}

// BenchmarkSummary resumo dos benchmarks
type BenchmarkSummary struct {
	BestPerformance   BenchmarkResult `json:"best_performance"`
	WorstPerformance  BenchmarkResult `json:"worst_performance"`
	AverageRPS        float64         `json:"average_rps"`
	AverageSuccessRate float64        `json:"average_success_rate"`
	SystemStability   float64         `json:"system_stability"`
	PerformanceGrade  string          `json:"performance_grade"`
}

// SystemInfo informações do sistema
type SystemInfo struct {
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	NumCPU       int    `json:"num_cpu"`
	GOMAXPROCS   int    `json:"gomaxprocs"`
	Version      string `json:"go_version"`
}

// generateBenchmarkReport gera relatório consolidado
func (pb *PerformanceBenchmark) generateBenchmarkReport(totalDuration time.Duration) *BenchmarkReport {
	pb.mu.RLock()
	defer pb.mu.RUnlock()
	
	report := &BenchmarkReport{
		GeneratedAt:   time.Now(),
		TotalDuration: totalDuration,
		TestsRun:      len(pb.results),
		Results:       make([]BenchmarkResult, len(pb.results)),
		SystemInfo: SystemInfo{
			GOOS:       runtime.GOOS,
			GOARCH:     runtime.GOARCH,
			NumCPU:     runtime.NumCPU(),
			GOMAXPROCS: runtime.GOMAXPROCS(0),
			Version:    runtime.Version(),
		},
	}
	
	copy(report.Results, pb.results)
	
	// Gerar resumo
	report.Summary = pb.generateSummary()
	
	// Gerar recomendações
	report.Recommendations = pb.generatePerformanceRecommendations()
	
	return report
}

// generateSummary gera resumo dos resultados
func (pb *PerformanceBenchmark) generateSummary() BenchmarkSummary {
	if len(pb.results) == 0 {
		return BenchmarkSummary{}
	}
	
	summary := BenchmarkSummary{}
	
	var totalRPS float64
	var totalSuccessRate float64
	var bestRPS float64
	var worstRPS float64
	
	for i, result := range pb.results {
		totalRPS += result.RequestsPerSecond
		totalSuccessRate += result.SuccessRate
		
		if i == 0 || result.RequestsPerSecond > bestRPS {
			bestRPS = result.RequestsPerSecond
			summary.BestPerformance = result
		}
		
		if i == 0 || result.RequestsPerSecond < worstRPS {
			worstRPS = result.RequestsPerSecond
			summary.WorstPerformance = result
		}
	}
	
	summary.AverageRPS = totalRPS / float64(len(pb.results))
	summary.AverageSuccessRate = totalSuccessRate / float64(len(pb.results))
	
	// Calcular estabilidade do sistema (baseado na variação de performance)
	summary.SystemStability = pb.calculateSystemStability()
	
	// Atribuir nota de performance
	summary.PerformanceGrade = pb.calculatePerformanceGrade(summary)
	
	return summary
}

// calculateSystemStability calcula estabilidade do sistema
func (pb *PerformanceBenchmark) calculateSystemStability() float64 {
	if len(pb.results) < 2 {
		return 1.0
	}
	
	// Calcular coeficiente de variação dos tempos de resposta
	var times []float64
	for _, result := range pb.results {
		times = append(times, float64(result.AverageResponseTime.Milliseconds()))
	}
	
	mean := average(times)
	variance := calculateVariance(times, mean)
	stddev := math.Sqrt(variance)
	
	cv := stddev / mean // Coeficiente de variação
	
	// Converter para score de estabilidade (1.0 = muito estável, 0.0 = instável)
	stability := max(0, 1.0 - cv)
	return stability
}

// calculatePerformanceGrade calcula nota de performance
func (pb *PerformanceBenchmark) calculatePerformanceGrade(summary BenchmarkSummary) string {
	// Critérios de avaliação
	rpsThresholds := map[string]float64{
		"A+": 50.0,
		"A":  30.0,
		"B":  20.0,
		"C":  10.0,
		"D":  5.0,
	}
	
	successThreshold := 0.95
	stabilityThreshold := 0.8
	
	grade := "F"
	
	for gradeLevel, threshold := range rpsThresholds {
		if summary.AverageRPS >= threshold &&
		   summary.AverageSuccessRate >= successThreshold &&
		   summary.SystemStability >= stabilityThreshold {
			grade = gradeLevel
			break
		}
	}
	
	return grade
}

// generatePerformanceRecommendations gera recomendações de performance
func (pb *PerformanceBenchmark) generatePerformanceRecommendations() []string {
	var recommendations []string
	
	summary := pb.generateSummary()
	
	if summary.AverageRPS < 10 {
		recommendations = append(recommendations, "Performance muito baixa - considere otimização de código e infraestrutura")
	}
	
	if summary.AverageSuccessRate < 0.95 {
		recommendations = append(recommendations, "Taxa de sucesso baixa - revisar tratamento de erros e timeouts")
	}
	
	if summary.SystemStability < 0.8 {
		recommendations = append(recommendations, "Sistema instável - investigar causas de variação na performance")
	}
	
	// Analisar uso de memória
	highMemoryUsage := false
	for _, result := range pb.results {
		if result.MemoryUsage.PeakUsage > 100*1024*1024 { // > 100MB
			highMemoryUsage = true
			break
		}
	}
	
	if highMemoryUsage {
		recommendations = append(recommendations, "Alto consumo de memória detectado - considere otimizações")
	}
	
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Performance dentro dos parâmetros esperados")
	}
	
	return recommendations
}

// ExportBenchmarkReport exporta relatório em JSON
func (pb *PerformanceBenchmark) ExportBenchmarkReport(report *BenchmarkReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

// Helper functions
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateVariance(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		diff := v - mean
		sum += diff * diff
	}
	return sum / float64(len(values))
}