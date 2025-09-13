package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"time"

	"BrxAgente-desafio4/internal/chat"
	"BrxAgente-desafio4/internal/config"
)

// RealWorldTestSuite executa suite completa de testes com cenários reais
type RealWorldTestSuite struct {
	agent               *VRAgent
	validator           *RealWorldValidator
	metricsCollector    *QualityMetricsCollector
	performanceBenchmark *PerformanceBenchmark
	logger              *log.Logger
	config              *TestSuiteConfig
}

// TestSuiteConfig configuração da suite de testes
type TestSuiteConfig struct {
	EnabledCategories   []string      `json:"enabled_categories"`
	AccuracyThreshold   float64       `json:"accuracy_threshold"`   // 0.95
	SatisfactionThreshold float64     `json:"satisfaction_threshold"` // 4.2
	MaxResponseTime     time.Duration `json:"max_response_time"`    // 2s
	StrictMode          bool          `json:"strict_mode"`
	RunBenchmarks       bool          `json:"run_benchmarks"`
	GenerateReports     bool          `json:"generate_reports"`
	OutputDir           string        `json:"output_dir"`
}

// DefaultTestSuiteConfig retorna configuração padrão
func DefaultTestSuiteConfig() *TestSuiteConfig {
	return &TestSuiteConfig{
		EnabledCategories:     []string{}, // Todas as categorias
		AccuracyThreshold:     0.95,
		SatisfactionThreshold: 4.2,
		MaxResponseTime:       2 * time.Second,
		StrictMode:            true,
		RunBenchmarks:         true,
		GenerateReports:       true,
		OutputDir:            "./test_reports",
	}
}

// NewRealWorldTestSuite cria nova suite de testes
func NewRealWorldTestSuite() (*RealWorldTestSuite, error) {
	// Setup básico para testes
	cfg := &config.Config{
		OpenAIKey: os.Getenv("OPENAI_KEY"),
		OllamaConfig: config.OllamaConfig{
			BaseURL: "http://localhost:11434",
			Model:   "llama2",
		},
	}

	chatSvc := chat.NewChat(cfg)
	
	agentConfig := &AgentConfig{
		Enabled:         true,
		Model:           "test-model",
		Temperature:     0.7,
		MaxTokens:       2000,
		Timeout:         30 * time.Second,
		MemorySize:      100,
		MemoryTTL:       24 * time.Hour,
		ContextWindow:   4000,
		MaxMemoryTokens: 2000,
		WorkerPoolSize:  4,
		CacheEnabled:    true,
		CacheSize:       1000,
		CacheTTL:        1 * time.Hour,
		LogLevel:        "info",
		DebugMode:       true,
		ToolsEnabled:    []string{"policy_consultant", "excel", "calculation", "validation"},
	}

	agent, err := NewVRAgent(agentConfig, chatSvc)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	logger := log.New(os.Stdout, "[RealWorldTest] ", log.LstdFlags)

	metricsCollector := NewQualityMetricsCollector(logger)
	validator := NewRealWorldValidator(agent, logger)
	benchmark := NewPerformanceBenchmark(agent, metricsCollector, logger)

	return &RealWorldTestSuite{
		agent:               agent,
		validator:           validator,
		metricsCollector:    metricsCollector,
		performanceBenchmark: benchmark,
		logger:              logger,
		config:              DefaultTestSuiteConfig(),
	}, nil
}

// SetConfig define configuração da suite
func (suite *RealWorldTestSuite) SetConfig(config *TestSuiteConfig) {
	suite.config = config
	
	// Aplicar configurações aos componentes
	suite.validator.SetCategories(config.EnabledCategories)
	suite.validator.SetTimeout(config.MaxResponseTime * 2) // Timeout maior para testes
	suite.validator.SetStrictMode(config.StrictMode)
	
	// Configurar thresholds para métricas
	thresholds := QualityThresholds{
		MinAccuracy:         config.AccuracyThreshold,
		MaxResponseTime:     config.MaxResponseTime,
		MinSatisfaction:     config.SatisfactionThreshold,
		MinCacheHitRate:     0.70,
		MaxErrorRate:        0.05,
		MinPolicyAdherence:  0.90,
	}
	suite.metricsCollector.SetThresholds(thresholds)
}

// RunCompleteValidation executa validação completa conforme issue #71
func (suite *RealWorldTestSuite) RunCompleteValidation() (*ComprehensiveReport, error) {
	suite.logger.Println("🚀 Iniciando validação completa com casos reais...")
	
	startTime := time.Now()
	
	// Criar diretório de saída se necessário
	if suite.config.GenerateReports {
		if err := os.MkdirAll(suite.config.OutputDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create output directory: %w", err)
		}
	}
	
	report := &ComprehensiveReport{
		TestSuite:   "Real World Validation - Issue #71",
		StartTime:   startTime,
		Config:      *suite.config,
	}
	
	// 1. Executar testes de validação
	suite.logger.Println("📋 Executando testes de validação...")
	validationReport, err := suite.validator.RunAllTests()
	if err != nil {
		return nil, fmt.Errorf("validation tests failed: %w", err)
	}
	report.ValidationReport = validationReport
	
	// 2. Executar benchmarks de performance se habilitado
	if suite.config.RunBenchmarks {
		suite.logger.Println("⚡ Executando benchmarks de performance...")
		benchmarkReport, err := suite.performanceBenchmark.RunAllBenchmarks()
		if err != nil {
			suite.logger.Printf("⚠️ Benchmark failed (non-critical): %v", err)
			// Continuar mesmo se benchmark falhar
		} else {
			report.BenchmarkReport = benchmarkReport
		}
	}
	
	// 3. Coletar métricas de qualidade
	suite.logger.Println("📊 Coletando métricas de qualidade...")
	report.QualityReport = suite.metricsCollector.GenerateReport()
	
	// 4. Executar testes específicos por categoria
	suite.logger.Println("🎯 Executando testes por categoria...")
	categoryResults, err := suite.runCategorySpecificTests()
	if err != nil {
		return nil, fmt.Errorf("category tests failed: %w", err)
	}
	report.CategoryResults = categoryResults
	
	// 5. Análise de casos edge
	suite.logger.Println("🔍 Testando casos edge...")
	edgeResults, err := suite.runEdgeCaseTests()
	if err != nil {
		return nil, fmt.Errorf("edge case tests failed: %w", err)
	}
	report.EdgeCaseResults = edgeResults
	
	// 6. Finalizar relatório
	report.EndTime = time.Now()
	report.TotalDuration = report.EndTime.Sub(report.StartTime)
	
	// 7. Gerar análise final e recomendações
	suite.generateFinalAnalysis(report)
	
	// 8. Exportar relatórios se habilitado
	if suite.config.GenerateReports {
		if err := suite.exportReports(report); err != nil {
			suite.logger.Printf("⚠️ Failed to export reports: %v", err)
		}
	}
	
	// 9. Log de resumo final
	suite.logFinalSummary(report)
	
	return report, nil
}

// runCategorySpecificTests executa testes específicos por categoria
func (suite *RealWorldTestSuite) runCategorySpecificTests() (map[string]CategoryTestResult, error) {
	categories := []string{
		"Políticas de Elegibilidade",
		"Cálculos Específicos",
		"Cenários Complexos",
		"Dados Processados",
	}
	
	results := make(map[string]CategoryTestResult)
	
	for _, category := range categories {
		suite.logger.Printf("Testing category: %s", category)
		
		scenarios := GetScenariosByCategory(category)
		if len(scenarios) == 0 {
			suite.logger.Printf("No scenarios found for category: %s", category)
			continue
		}
		
		result := CategoryTestResult{
			Category:      category,
			TotalScenarios: len(scenarios),
			StartTime:     time.Now(),
		}
		
		var passedTests int
		var totalScore float64
		var responseTimes []time.Duration
		var errors []string
		
		for _, scenario := range scenarios {
			startTime := time.Now()
			
			response, err := suite.agent.Ask(scenario.Question)
			responseTime := time.Since(startTime)
			responseTimes = append(responseTimes, responseTime)
			
			if err != nil {
				errors = append(errors, fmt.Sprintf("Scenario %s: %v", scenario.ID, err))
				continue
			}
			
			// Avaliar resposta
			score := suite.evaluateScenarioResponse(response, scenario)
			totalScore += score
			
			// Registrar métrica
			passed := score >= 4.0
			if passed {
				passedTests++
			}
			
			suite.metricsCollector.RecordQuestion(
				scenario.Question,
				response,
				responseTime,
				category,
				score,
				passed,
			)
		}
		
		result.PassedScenarios = passedTests
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		result.Accuracy = float64(passedTests) / float64(len(scenarios))
		result.AverageScore = totalScore / float64(len(scenarios))
		result.Errors = errors
		
		// Calcular estatísticas de tempo
		if len(responseTimes) > 0 {
			result.AverageResponseTime = suite.calculateAverageTime(responseTimes)
			result.MaxResponseTime = suite.calculateMaxTime(responseTimes)
		}
		
		results[category] = result
		
		suite.logger.Printf("Category %s completed: %.1f%% accuracy, %.2f avg score",
			category, result.Accuracy*100, result.AverageScore)
	}
	
	return results, nil
}

// runEdgeCaseTests testa casos edge e cenários difíceis
func (suite *RealWorldTestSuite) runEdgeCaseTests() (*EdgeCaseResults, error) {
	suite.logger.Println("Running edge case tests...")
	
	results := &EdgeCaseResults{
		StartTime: time.Now(),
		TestCases: []EdgeCaseTest{},
	}
	
	// 1. Teste de perguntas muito longas
	longQuestionResult := suite.testLongQuestion()
	results.TestCases = append(results.TestCases, longQuestionResult)
	
	// 2. Teste de perguntas ambíguas
	ambiguousResult := suite.testAmbiguousQuestions()
	results.TestCases = append(results.TestCases, ambiguousResult)
	
	// 3. Teste de timeout
	timeoutResult := suite.testTimeoutScenarios()
	results.TestCases = append(results.TestCases, timeoutResult)
	
	// 4. Teste de concorrência
	concurrencyResult := suite.testConcurrentRequests()
	results.TestCases = append(results.TestCases, concurrencyResult)
	
	// 5. Teste de memória
	memoryResult := suite.testMemoryUsage()
	results.TestCases = append(results.TestCases, memoryResult)
	
	results.EndTime = time.Now()
	results.Duration = results.EndTime.Sub(results.StartTime)
	
	// Calcular estatísticas gerais
	passed := 0
	for _, test := range results.TestCases {
		if test.Passed {
			passed++
		}
	}
	results.OverallSuccess = float64(passed) / float64(len(results.TestCases))
	
	return results, nil
}

// evaluateScenarioResponse avalia resposta para um cenário específico
func (suite *RealWorldTestSuite) evaluateScenarioResponse(response string, scenario RealWorldTestScenario) float64 {
	score := 3.0 // Score base
	
	response = strings.ToLower(response)
	
	// 1. Verificar palavras obrigatórias
	mustContainScore := 0.0
	if len(scenario.Expected.MustContain) > 0 {
		found := 0
		for _, term := range scenario.Expected.MustContain {
			if strings.Contains(response, strings.ToLower(term)) {
				found++
			}
		}
		mustContainScore = float64(found) / float64(len(scenario.Expected.MustContain))
		score += mustContainScore * 1.5 // Max 1.5 pontos
	}
	
	// 2. Verificar palavras proibidas (penalizar)
	if len(scenario.Expected.MustNotContain) > 0 {
		violations := 0
		for _, term := range scenario.Expected.MustNotContain {
			if strings.Contains(response, strings.ToLower(term)) {
				violations++
			}
		}
		if violations > 0 {
			penalty := float64(violations) / float64(len(scenario.Expected.MustNotContain))
			score -= penalty * 2.0 // Penalizar pesado
		}
	}
	
	// 3. Verificar pontos-chave
	if len(scenario.Expected.KeyPoints) > 0 {
		covered := 0
		for _, point := range scenario.Expected.KeyPoints {
			pointWords := strings.Fields(strings.ToLower(point))
			if len(pointWords) > 0 {
				found := false
				for _, word := range pointWords {
					if len(word) > 3 && strings.Contains(response, word) {
						found = true
						break
					}
				}
				if found {
					covered++
				}
			}
		}
		keyPointsScore := float64(covered) / float64(len(scenario.Expected.KeyPoints))
		score += keyPointsScore * 0.5 // Max 0.5 pontos
	}
	
	return math.Min(5.0, math.Max(0.0, score))
}

// testLongQuestion testa perguntas muito longas
func (suite *RealWorldTestSuite) testLongQuestion() EdgeCaseTest {
	longQuestion := strings.Repeat("Como calcular VR para colaborador com licença médica de 30 dias, afastado pelo INSS, admitido no meio do mês, com férias programadas e sindicato SINDPD? ", 10)
	
	start := time.Now()
	response, err := suite.agent.Ask(longQuestion)
	duration := time.Since(start)
	
	passed := err == nil && len(response) > 0 && duration < 10*time.Second
	
	return EdgeCaseTest{
		Name:        "Long Question Test",
		Description: "Testa processamento de perguntas muito longas",
		Input:       longQuestion,
		Expected:    "Resposta estruturada mesmo para pergunta longa",
		Actual:      response,
		Passed:      passed,
		Duration:    duration,
		Error:       err,
	}
}

// testAmbiguousQuestions testa perguntas ambíguas
func (suite *RealWorldTestSuite) testAmbiguousQuestions() EdgeCaseTest {
	ambiguousQuestions := []string{
		"Qual o valor?",
		"Como funciona?",
		"E se?",
		"Tem direito?",
	}
	
	start := time.Now()
	var allPassed = true
	var responses []string
	
	for _, question := range ambiguousQuestions {
		response, err := suite.agent.Ask(question)
		responses = append(responses, response)
		
		// Para perguntas ambíguas, esperamos que o agente peça esclarecimento
		if err != nil || len(response) < 10 {
			allPassed = false
		}
	}
	
	duration := time.Since(start)
	
	return EdgeCaseTest{
		Name:        "Ambiguous Questions Test",
		Description: "Testa como o agente lida com perguntas ambíguas",
		Input:       fmt.Sprintf("%v", ambiguousQuestions),
		Expected:    "Solicitação de esclarecimento ou resposta informativa",
		Actual:      fmt.Sprintf("%v", responses),
		Passed:      allPassed,
		Duration:    duration,
	}
}

// testTimeoutScenarios testa cenários de timeout
func (suite *RealWorldTestSuite) testTimeoutScenarios() EdgeCaseTest {
	// Simular pergunta que pode causar timeout
	complexQuestion := "Análise completa de todos os colaboradores processados, incluindo cálculos detalhados de VR, distribuição por sindicato, análise de anomalias, relatórios por categoria e sugestões de otimização para os próximos 12 meses"
	
	start := time.Now()
	response, err := suite.agent.Ask(complexQuestion)
	duration := time.Since(start)
	
	// Passou se respondeu dentro de tempo razoável
	passed := duration < 30*time.Second && (err == nil || len(response) > 0)
	
	return EdgeCaseTest{
		Name:        "Timeout Test",
		Description: "Testa comportamento em perguntas complexas que podem causar timeout",
		Input:       complexQuestion,
		Expected:    "Resposta dentro do tempo limite",
		Actual:      response,
		Passed:      passed,
		Duration:    duration,
		Error:       err,
	}
}

// testConcurrentRequests testa requests concorrentes
func (suite *RealWorldTestSuite) testConcurrentRequests() EdgeCaseTest {
	questions := []string{
		"Quantos colaboradores foram processados?",
		"Qual o total de VR?",
		"Diretores têm direito a VR?",
		"Como calcular VR proporcional?",
	}
	
	start := time.Now()
	
	type result struct {
		response string
		err      error
		duration time.Duration
	}
	
	results := make(chan result, len(questions))
	
	// Executar requests em paralelo
	for _, question := range questions {
		go func(q string) {
			reqStart := time.Now()
			resp, err := suite.agent.Ask(q)
			results <- result{
				response: resp,
				err:      err,
				duration: time.Since(reqStart),
			}
		}(question)
	}
	
	// Coletar resultados
	var allPassed = true
	var responses []string
	for i := 0; i < len(questions); i++ {
		res := <-results
		responses = append(responses, res.response)
		if res.err != nil || res.duration > 10*time.Second {
			allPassed = false
		}
	}
	
	totalDuration := time.Since(start)
	
	return EdgeCaseTest{
		Name:        "Concurrent Requests Test",
		Description: "Testa processamento concorrente de múltiplas perguntas",
		Input:       fmt.Sprintf("%v", questions),
		Expected:    "Todas as perguntas respondidas corretamente",
		Actual:      fmt.Sprintf("%v", responses),
		Passed:      allPassed,
		Duration:    totalDuration,
	}
}

// testMemoryUsage testa uso de memória
func (suite *RealWorldTestSuite) testMemoryUsage() EdgeCaseTest {
	// Capturar uso inicial de memória
	var initialMem runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&initialMem)
	
	start := time.Now()
	
	// Executar muitas perguntas para testar vazamentos de memória
	for i := 0; i < 50; i++ {
		question := fmt.Sprintf("Teste de memória iteração %d - quantos colaboradores foram processados?", i)
		suite.agent.Ask(question)
	}
	
	// Forçar GC e medir memória final
	runtime.GC()
	var finalMem runtime.MemStats
	runtime.ReadMemStats(&finalMem)
	
	duration := time.Since(start)
	memoryGrowth := finalMem.Alloc - initialMem.Alloc
	
	// Passou se crescimento de memória for razoável (< 10MB)
	passed := memoryGrowth < 10*1024*1024
	
	return EdgeCaseTest{
		Name:        "Memory Usage Test",
		Description: "Testa vazamentos de memória durante múltiplas operações",
		Input:       "50 perguntas sequenciais",
		Expected:    "Crescimento de memória < 10MB",
		Actual:      fmt.Sprintf("Crescimento: %.2f MB", float64(memoryGrowth)/1024/1024),
		Passed:      passed,
		Duration:    duration,
		Metadata: map[string]any{
			"initial_memory":  initialMem.Alloc,
			"final_memory":    finalMem.Alloc,
			"memory_growth":   memoryGrowth,
		},
	}
}

// generateFinalAnalysis gera análise final e recomendações
func (suite *RealWorldTestSuite) generateFinalAnalysis(report *ComprehensiveReport) {
	suite.logger.Println("Gerando análise final...")
	
	analysis := FinalAnalysis{
		OverallGrade: "F", // Será calculado
		Recommendations: []string{},
		CriticalIssues: []string{},
		Achievements: []string{},
	}
	
	// Analisar accuracy geral
	if report.ValidationReport != nil {
		accuracy := report.ValidationReport.OverallAccuracy
		analysis.AccuracyAnalysis = fmt.Sprintf("Accuracy geral: %.1f%% (Meta: %.1f%%)",
			accuracy*100, suite.config.AccuracyThreshold*100)
			
		if accuracy >= suite.config.AccuracyThreshold {
			analysis.Achievements = append(analysis.Achievements, "✅ Meta de accuracy atingida")
		} else {
			analysis.CriticalIssues = append(analysis.CriticalIssues, 
				fmt.Sprintf("❌ Accuracy abaixo da meta: %.1f%% < %.1f%%", 
					accuracy*100, suite.config.AccuracyThreshold*100))
		}
	}
	
	// Analisar performance
	if report.BenchmarkReport != nil {
		avgRPS := report.BenchmarkReport.Summary.AverageRPS
		analysis.PerformanceAnalysis = fmt.Sprintf("Performance média: %.1f RPS", avgRPS)
		
		if avgRPS >= 10.0 {
			analysis.Achievements = append(analysis.Achievements, "✅ Performance adequada")
		} else {
			analysis.CriticalIssues = append(analysis.CriticalIssues, "❌ Performance abaixo do esperado")
		}
	}
	
	// Analisar qualidade
	satisfaction := report.QualityReport.AverageSatisfaction
	if satisfaction >= suite.config.SatisfactionThreshold {
		analysis.Achievements = append(analysis.Achievements, "✅ Meta de satisfação atingida")
	} else {
		analysis.CriticalIssues = append(analysis.CriticalIssues, 
			fmt.Sprintf("❌ Satisfação abaixo da meta: %.1f < %.1f", 
				satisfaction, suite.config.SatisfactionThreshold))
	}
	
	// Calcular nota geral
	analysis.OverallGrade = suite.calculateOverallGrade(report)
	
	// Gerar recomendações baseadas na análise
	analysis.Recommendations = suite.generateDetailedRecommendations(report)
	
	report.FinalAnalysis = analysis
}

// calculateOverallGrade calcula nota geral do sistema
func (suite *RealWorldTestSuite) calculateOverallGrade(report *ComprehensiveReport) string {
	score := 0.0
	maxScore := 0.0
	
	// Accuracy (peso 40%)
	if report.ValidationReport != nil {
		accuracy := report.ValidationReport.OverallAccuracy
		score += accuracy * 40
	}
	maxScore += 40
	
	// Performance (peso 30%)
	if report.BenchmarkReport != nil {
		rps := report.BenchmarkReport.Summary.AverageRPS
		perfScore := math.Min(30, rps/10.0 * 30) // 10 RPS = nota máxima
		score += perfScore
	}
	maxScore += 30
	
	// Satisfação (peso 20%)
	satisfaction := report.QualityReport.AverageSatisfaction
	satScore := (satisfaction / 5.0) * 20
	score += satScore
	maxScore += 20
	
	// Edge cases (peso 10%)
	if report.EdgeCaseResults != nil {
		edgeScore := report.EdgeCaseResults.OverallSuccess * 10
		score += edgeScore
	}
	maxScore += 10
	
	percentage := (score / maxScore) * 100
	
	switch {
	case percentage >= 90:
		return "A+"
	case percentage >= 85:
		return "A"
	case percentage >= 80:
		return "B+"
	case percentage >= 75:
		return "B"
	case percentage >= 70:
		return "C+"
	case percentage >= 65:
		return "C"
	case percentage >= 60:
		return "D"
	default:
		return "F"
	}
}

// generateDetailedRecommendations gera recomendações detalhadas
func (suite *RealWorldTestSuite) generateDetailedRecommendations(report *ComprehensiveReport) []string {
	var recommendations []string
	
	// Baseado na accuracy
	if report.ValidationReport != nil && report.ValidationReport.OverallAccuracy < 0.95 {
		recommendations = append(recommendations, 
			"🎯 Melhorar base de conhecimento para aumentar accuracy")
		recommendations = append(recommendations,
			"📚 Revisar e expandir templates de resposta")
	}
	
	// Baseado na performance
	if report.BenchmarkReport != nil && report.BenchmarkReport.Summary.AverageRPS < 10 {
		recommendations = append(recommendations,
			"⚡ Otimizar performance do sistema - considere cache e otimização de queries")
	}
	
	// Baseado na satisfação
	if report.QualityReport.AverageSatisfaction < 4.0 {
		recommendations = append(recommendations,
			"😊 Melhorar qualidade e formato das respostas para aumentar satisfação")
	}
	
	// Baseado nos edge cases
	if report.EdgeCaseResults != nil && report.EdgeCaseResults.OverallSuccess < 0.8 {
		recommendations = append(recommendations,
			"🔧 Melhorar robustez do sistema para casos edge")
	}
	
	// Baseado nas categorias com pior performance
	worstCategory := ""
	lowestAccuracy := 1.0
	for category, result := range report.CategoryResults {
		if result.Accuracy < lowestAccuracy {
			lowestAccuracy = result.Accuracy
			worstCategory = category
		}
	}
	
	if worstCategory != "" && lowestAccuracy < 0.8 {
		recommendations = append(recommendations,
			fmt.Sprintf("📋 Focar melhorias na categoria '%s' (accuracy: %.1f%%)", 
				worstCategory, lowestAccuracy*100))
	}
	
	if len(recommendations) == 0 {
		recommendations = append(recommendations,
			"🎉 Sistema atendendo todos os critérios! Manter monitoramento contínuo.")
	}
	
	return recommendations
}

// Helper functions
func (suite *RealWorldTestSuite) calculateAverageTime(times []time.Duration) time.Duration {
	if len(times) == 0 {
		return 0
	}
	
	var total time.Duration
	for _, t := range times {
		total += t
	}
	return total / time.Duration(len(times))
}

func (suite *RealWorldTestSuite) calculateMaxTime(times []time.Duration) time.Duration {
	if len(times) == 0 {
		return 0
	}
	
	max := times[0]
	for _, t := range times {
		if t > max {
			max = t
		}
	}
	return max
}

// exportReports exporta todos os relatórios
func (suite *RealWorldTestSuite) exportReports(report *ComprehensiveReport) error {
	suite.logger.Println("Exportando relatórios...")
	
	// Relatório principal em JSON
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}
	
	filename := fmt.Sprintf("%s/comprehensive_report_%s.json", 
		suite.config.OutputDir, time.Now().Format("20060102_150405"))
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write report file: %w", err)
	}
	
	suite.logger.Printf("Relatório exportado: %s", filename)
	
	// Exportar relatório de métricas
	if metricsData, err := suite.metricsCollector.ExportMetrics(); err == nil {
		metricsFile := fmt.Sprintf("%s/metrics_report_%s.json", 
			suite.config.OutputDir, time.Now().Format("20060102_150405"))
		os.WriteFile(metricsFile, metricsData, 0644)
	}
	
	return nil
}

// logFinalSummary registra resumo final no log
func (suite *RealWorldTestSuite) logFinalSummary(report *ComprehensiveReport) {
	suite.logger.Println("\n" + strings.Repeat("=", 60))
	suite.logger.Println("🏁 RESUMO FINAL DA VALIDAÇÃO")
	suite.logger.Println(strings.Repeat("=", 60))
	
	if report.ValidationReport != nil {
		suite.logger.Printf("📊 Accuracy Geral: %.1f%% (%d/%d cenários passaram)",
			report.ValidationReport.OverallAccuracy*100,
			report.ValidationReport.PassedScenarios,
			report.ValidationReport.TotalScenarios)
	}
	
	if report.BenchmarkReport != nil {
		suite.logger.Printf("⚡ Performance: %.1f RPS (Nota: %s)",
			report.BenchmarkReport.Summary.AverageRPS,
			report.BenchmarkReport.Summary.PerformanceGrade)
	}
	
	suite.logger.Printf("😊 Satisfação Média: %.1f/5.0",
		report.QualityReport.AverageSatisfaction)
	
	if report.EdgeCaseResults != nil {
		suite.logger.Printf("🔧 Casos Edge: %.1f%% sucesso",
			report.EdgeCaseResults.OverallSuccess*100)
	}
	
	suite.logger.Printf("🏆 Nota Final: %s", report.FinalAnalysis.OverallGrade)
	suite.logger.Printf("⏱️  Tempo Total: %v", report.TotalDuration)
	
	suite.logger.Println("\n📝 RECOMENDAÇÕES:")
	for _, rec := range report.FinalAnalysis.Recommendations {
		suite.logger.Printf("   %s", rec)
	}
	
	if len(report.FinalAnalysis.CriticalIssues) > 0 {
		suite.logger.Println("\n❗ ISSUES CRÍTICAS:")
		for _, issue := range report.FinalAnalysis.CriticalIssues {
			suite.logger.Printf("   %s", issue)
		}
	}
	
	suite.logger.Println(strings.Repeat("=", 60))
}

// FinalReportGenerator gera relatórios finais em múltiplos formatos
type FinalReportGenerator struct {
	outputDir string
	logger    Logger
}

// NewFinalReportGenerator cria novo gerador de relatórios
func NewFinalReportGenerator(outputDir string, logger Logger) *FinalReportGenerator {
	return &FinalReportGenerator{
		outputDir: outputDir,
		logger:    logger,
	}
}

// GenerateAllReports gera todos os tipos de relatório
func (frg *FinalReportGenerator) GenerateAllReports(report *ComprehensiveReport) error {
	frg.logger.Println("📄 Gerando relatórios finais...")

	// Criar diretório se não existir
	if err := os.MkdirAll(frg.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")

	// 1. Relatório JSON detalhado
	if err := frg.generateJSONReport(report, timestamp); err != nil {
		frg.logger.Printf("❌ Erro ao gerar relatório JSON: %v", err)
	} else {
		frg.logger.Println("✅ Relatório JSON gerado")
	}

	// 2. Relatório Markdown para documentação
	if err := frg.generateMarkdownReport(report, timestamp); err != nil {
		frg.logger.Printf("❌ Erro ao gerar relatório Markdown: %v", err)
	} else {
		frg.logger.Println("✅ Relatório Markdown gerado")
	}

	// 3. Relatório CSV com métricas
	if err := frg.generateCSVReport(report, timestamp); err != nil {
		frg.logger.Printf("❌ Erro ao gerar relatório CSV: %v", err)
	} else {
		frg.logger.Println("✅ Relatório CSV gerado")
	}

	// 4. Sumário executivo
	if err := frg.generateExecutiveSummary(report, timestamp); err != nil {
		frg.logger.Printf("❌ Erro ao gerar sumário executivo: %v", err)
	} else {
		frg.logger.Println("✅ Sumário executivo gerado")
	}

	frg.logger.Printf("📁 Todos os relatórios salvos em: %s", frg.outputDir)
	return nil
}

// generateJSONReport gera relatório JSON completo
func (frg *FinalReportGenerator) generateJSONReport(report *ComprehensiveReport, timestamp string) error {
	filename := fmt.Sprintf("%s/issue71_complete_report_%s.json", frg.outputDir, timestamp)

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// generateMarkdownReport gera relatório Markdown
func (frg *FinalReportGenerator) generateMarkdownReport(report *ComprehensiveReport, timestamp string) error {
	filename := fmt.Sprintf("%s/issue71_report_%s.md", frg.outputDir, timestamp)

	var md strings.Builder

	md.WriteString("# Issue #71: Teste com Casos Reais e Validação\n\n")
	md.WriteString("## 📊 Resumo Executivo\n\n")
	
	// Informações básicas
	md.WriteString(fmt.Sprintf("- **Data do Teste:** %s\n", report.StartTime.Format("02/01/2006 15:04:05")))
	md.WriteString(fmt.Sprintf("- **Duração Total:** %v\n", report.TotalDuration))
	md.WriteString(fmt.Sprintf("- **Nota Final:** %s\n", report.FinalAnalysis.OverallGrade))
	md.WriteString("\n")

	// Métricas principais
	md.WriteString("## 🎯 Métricas Principais\n\n")
	
	if report.ValidationReport != nil {
		status := "✅ ATINGIDA"
		if report.ValidationReport.OverallAccuracy < 0.95 {
			status = "❌ NÃO ATINGIDA"
		}
		md.WriteString(fmt.Sprintf("### Accuracy\n- **Resultado:** %.1f%%\n- **Meta:** ≥ 95%%\n- **Status:** %s\n\n",
			report.ValidationReport.OverallAccuracy*100, status))
	}

	if report.BenchmarkReport != nil {
		status := "✅ ATINGIDA"
		if report.BenchmarkReport.Summary.AverageRPS < 10 {
			status = "❌ NÃO ATINGIDA"
		}
		md.WriteString(fmt.Sprintf("### Performance\n- **RPS:** %.1f\n- **Nota:** %s\n- **Status:** %s\n\n",
			report.BenchmarkReport.Summary.AverageRPS,
			report.BenchmarkReport.Summary.PerformanceGrade, status))
	}

	status := "✅ ATINGIDA"
	if report.QualityReport.AverageSatisfaction < 4.2 {
		status = "❌ NÃO ATINGIDA"
	}
	md.WriteString(fmt.Sprintf("### Satisfação\n- **Score:** %.1f/5.0\n- **Meta:** ≥ 4.2\n- **Status:** %s\n\n",
		report.QualityReport.AverageSatisfaction, status))

	// Resultados por categoria
	md.WriteString("## 📋 Resultados por Categoria\n\n")
	for category, result := range report.CategoryResults {
		md.WriteString(fmt.Sprintf("### %s\n", category))
		md.WriteString(fmt.Sprintf("- **Accuracy:** %.1f%% (%d/%d cenários)\n",
			result.Accuracy*100, result.PassedScenarios, result.TotalScenarios))
		md.WriteString(fmt.Sprintf("- **Score Médio:** %.1f/5.0\n", result.AverageScore))
		md.WriteString(fmt.Sprintf("- **Tempo Médio:** %v\n", result.AverageResponseTime))
		md.WriteString("\n")
	}

	// Casos edge
	if report.EdgeCaseResults != nil {
		md.WriteString("## 🔧 Casos Edge\n\n")
		md.WriteString(fmt.Sprintf("**Taxa de Sucesso Geral:** %.1f%%\n\n", report.EdgeCaseResults.OverallSuccess*100))
		
		md.WriteString("| Teste | Status | Tempo | Descrição |\n")
		md.WriteString("|-------|--------|-------|-----------|\n")
		
		for _, test := range report.EdgeCaseResults.TestCases {
			status := "❌ FAILED"
			if test.Passed {
				status = "✅ PASSED"
			}
			md.WriteString(fmt.Sprintf("| %s | %s | %.2fs | %s |\n",
				test.Name, status, test.Duration.Seconds(), test.Description))
		}
		md.WriteString("\n")
	}

	// Conquistas
	if len(report.FinalAnalysis.Achievements) > 0 {
		md.WriteString("## 🎉 Conquistas\n\n")
		for _, achievement := range report.FinalAnalysis.Achievements {
			md.WriteString(fmt.Sprintf("- %s\n", achievement))
		}
		md.WriteString("\n")
	}

	// Issues críticas
	if len(report.FinalAnalysis.CriticalIssues) > 0 {
		md.WriteString("## ❗ Issues Críticas\n\n")
		for _, issue := range report.FinalAnalysis.CriticalIssues {
			md.WriteString(fmt.Sprintf("- %s\n", issue))
		}
		md.WriteString("\n")
	}

	// Recomendações
	md.WriteString("## 📝 Recomendações\n\n")
	for _, rec := range report.FinalAnalysis.Recommendations {
		md.WriteString(fmt.Sprintf("- %s\n", rec))
	}
	md.WriteString("\n")

	// Footer
	md.WriteString("---\n")
	md.WriteString(fmt.Sprintf("*Relatório gerado automaticamente em %s*\n", time.Now().Format("02/01/2006 15:04:05")))
	md.WriteString("*VRAgent Real World Validation Suite - Issue #71*\n")

	return os.WriteFile(filename, []byte(md.String()), 0644)
}

// generateCSVReport gera relatório CSV com métricas
func (frg *FinalReportGenerator) generateCSVReport(report *ComprehensiveReport, timestamp string) error {
	filename := fmt.Sprintf("%s/issue71_metrics_%s.csv", frg.outputDir, timestamp)

	var csv strings.Builder
	csv.WriteString("Categoria,Accuracy,Score_Medio,Tempo_Medio_ms,Cenarios_Total,Cenarios_Passaram\n")

	for category, result := range report.CategoryResults {
		csv.WriteString(fmt.Sprintf("%s,%.3f,%.2f,%d,%d,%d\n",
			strings.ReplaceAll(category, ",", ";"),
			result.Accuracy,
			result.AverageScore,
			result.AverageResponseTime.Milliseconds(),
			result.TotalScenarios,
			result.PassedScenarios))
	}

	return os.WriteFile(filename, []byte(csv.String()), 0644)
}

// generateExecutiveSummary gera sumário executivo
func (frg *FinalReportGenerator) generateExecutiveSummary(report *ComprehensiveReport, timestamp string) error {
	filename := fmt.Sprintf("%s/issue71_executive_summary_%s.md", frg.outputDir, timestamp)

	var summary strings.Builder

	summary.WriteString("# Sumário Executivo - Issue #71\n\n")
	summary.WriteString("## Objetivo\n")
	summary.WriteString("Validar que todas as melhorias implementadas funcionam adequadamente em cenários reais de produção.\n\n")

	// Status geral
	summary.WriteString("## Status Geral\n")
	summary.WriteString(fmt.Sprintf("**Nota Final: %s**\n\n", report.FinalAnalysis.OverallGrade))

	// Métricas vs Metas
	summary.WriteString("## Métricas vs Metas\n\n")
	summary.WriteString("| Métrica | Resultado | Meta | Status |\n")
	summary.WriteString("|---------|-----------|------|--------|\n")

	if report.ValidationReport != nil {
		status := "✅ ATINGIDA"
		if report.ValidationReport.OverallAccuracy < 0.95 {
			status = "❌ NÃO ATINGIDA"
		}
		summary.WriteString(fmt.Sprintf("| Accuracy | %.1f%% | ≥ 95%% | %s |\n",
			report.ValidationReport.OverallAccuracy*100, status))
	}

	if report.BenchmarkReport != nil {
		status := "✅ ATINGIDA"
		if report.BenchmarkReport.Summary.AverageRPS < 10 {
			status = "❌ NÃO ATINGIDA"
		}
		summary.WriteString(fmt.Sprintf("| Performance | %.1f RPS | ≥ 10 RPS | %s |\n",
			report.BenchmarkReport.Summary.AverageRPS, status))
	}

	status := "✅ ATINGIDA"
	if report.QualityReport.AverageSatisfaction < 4.2 {
		status = "❌ NÃO ATINGIDA"
	}
	summary.WriteString(fmt.Sprintf("| Satisfação | %.1f/5.0 | ≥ 4.2 | %s |\n",
		report.QualityReport.AverageSatisfaction, status))

	summary.WriteString("\n")

	// Próximos passos
	summary.WriteString("## Próximos Passos\n\n")
	if len(report.FinalAnalysis.CriticalIssues) > 0 {
		summary.WriteString("### Issues Críticas a Resolver:\n")
		for _, issue := range report.FinalAnalysis.CriticalIssues {
			summary.WriteString(fmt.Sprintf("- %s\n", issue))
		}
		summary.WriteString("\n")
	}

	summary.WriteString("### Recomendações Prioritárias:\n")
	maxRecs := 3
	if len(report.FinalAnalysis.Recommendations) < maxRecs {
		maxRecs = len(report.FinalAnalysis.Recommendations)
	}
	for i := 0; i < maxRecs; i++ {
		summary.WriteString(fmt.Sprintf("%d. %s\n", i+1, report.FinalAnalysis.Recommendations[i]))
	}

	// Conclusão
	summary.WriteString("\n## Conclusão\n")
	grade := report.FinalAnalysis.OverallGrade
	switch grade {
	case "A+", "A":
		summary.WriteString("✅ Sistema apresentou excelente performance e atendeu todos os critérios de qualidade.\n")
	case "B+", "B":
		summary.WriteString("✅ Sistema apresentou boa performance com pequenos pontos de melhoria.\n")
	case "C+", "C":
		summary.WriteString("⚠️ Sistema funcional mas requer melhorias para atingir critérios ideais.\n")
	case "D":
		summary.WriteString("⚠️ Sistema apresenta deficiências significativas que precisam ser endereçadas.\n")
	case "F":
		summary.WriteString("❌ Sistema não atendeu aos critérios mínimos e requer revisão completa.\n")
	}

	return os.WriteFile(filename, []byte(summary.String()), 0644)
}

// Estruturas de dados para os relatórios

// ComprehensiveReport relatório abrangente de toda a validação
type ComprehensiveReport struct {
	TestSuite        string                      `json:"test_suite"`
	StartTime        time.Time                   `json:"start_time"`
	EndTime          time.Time                   `json:"end_time"`
	TotalDuration    time.Duration               `json:"total_duration"`
	Config           TestSuiteConfig             `json:"config"`
	ValidationReport *ValidationReport           `json:"validation_report"`
	BenchmarkReport  *BenchmarkReport            `json:"benchmark_report"`
	QualityReport    QualityReport               `json:"quality_report"`
	CategoryResults  map[string]CategoryTestResult `json:"category_results"`
	EdgeCaseResults  *EdgeCaseResults            `json:"edge_case_results"`
	FinalAnalysis    FinalAnalysis               `json:"final_analysis"`
}

// CategoryTestResult resultado de testes por categoria
type CategoryTestResult struct {
	Category            string        `json:"category"`
	TotalScenarios      int           `json:"total_scenarios"`
	PassedScenarios     int           `json:"passed_scenarios"`
	Accuracy            float64       `json:"accuracy"`
	AverageScore        float64       `json:"average_score"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	MaxResponseTime     time.Duration `json:"max_response_time"`
	StartTime           time.Time     `json:"start_time"`
	EndTime             time.Time     `json:"end_time"`
	Duration            time.Duration `json:"duration"`
	Errors              []string      `json:"errors"`
}

// EdgeCaseResults resultados dos testes de casos edge
type EdgeCaseResults struct {
	StartTime       time.Time       `json:"start_time"`
	EndTime         time.Time       `json:"end_time"`
	Duration        time.Duration   `json:"duration"`
	TestCases       []EdgeCaseTest  `json:"test_cases"`
	OverallSuccess  float64         `json:"overall_success"`
}

// EdgeCaseTest teste de caso edge individual
type EdgeCaseTest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Input       string                 `json:"input"`
	Expected    string                 `json:"expected"`
	Actual      string                 `json:"actual"`
	Passed      bool                   `json:"passed"`
	Duration    time.Duration          `json:"duration"`
	Error       error                  `json:"error,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// FinalAnalysis análise final com recomendações
type FinalAnalysis struct {
	OverallGrade          string   `json:"overall_grade"`
	AccuracyAnalysis      string   `json:"accuracy_analysis"`
	PerformanceAnalysis   string   `json:"performance_analysis"`
	Recommendations       []string `json:"recommendations"`
	CriticalIssues        []string `json:"critical_issues"`
	Achievements          []string `json:"achievements"`
}