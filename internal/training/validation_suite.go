package training

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TestCase representa um caso de teste individual
type TestCase struct {
	ID                    string            `json:"id"`
	Question              string            `json:"question"`
	ExpectedAnswerContains []string         `json:"expected_answer_contains"`
	ExpectedConfidence    float64           `json:"expected_confidence"`
	Category              string            `json:"category"`
	Priority              string            `json:"priority"`
	InputData             map[string]interface{} `json:"input_data,omitempty"`
	ExpectedCalculation   *ExpectedCalculation   `json:"expected_calculation,omitempty"`
}

// ExpectedCalculation representa um cálculo esperado
type ExpectedCalculation struct {
	Valor     float64 `json:"valor"`
	Reasoning string  `json:"reasoning"`
	Formula   string  `json:"formula"`
}

// ConsistencyTest representa um teste de consistência
type ConsistencyTest struct {
	ID                 string      `json:"id"`
	Description        string      `json:"description"`
	TestQuestion       string      `json:"test_question"`
	TestScenario       interface{} `json:"test_scenario,omitempty"`
	RepeatCount        int         `json:"repeat_count"`
	ExpectedConsistency float64    `json:"expected_consistency"`
	ExpectedResult     interface{} `json:"expected_result,omitempty"`
	Tolerance          float64     `json:"tolerance"`
	Category           string      `json:"category"`
}

// QualityTest representa um teste de qualidade
type QualityTest struct {
	ID               string                 `json:"id"`
	Question         string                 `json:"question"`
	QualityCriteria  map[string]bool        `json:"quality_criteria"`
	MinScore         float64                `json:"min_score"`
	Category         string                 `json:"category"`
}

// PerformanceTest representa um teste de performance
type PerformanceTest struct {
	ID                    string  `json:"id"`
	TestType              string  `json:"test_type"`
	Question              string  `json:"question"`
	MaxResponseTimeSeconds float64 `json:"max_response_time_seconds"`
	ExpectedAvgTime       float64 `json:"expected_avg_time"`
	Category              string  `json:"category"`
}

// TestResult representa o resultado de um teste
type TestResult struct {
	TestID      string    `json:"test_id"`
	TestType    string    `json:"test_type"`
	Passed      bool      `json:"passed"`
	Score       float64   `json:"score"`
	ExpectedScore float64 `json:"expected_score"`
	ActualResponse string `json:"actual_response"`
	ExpectedElements []string `json:"expected_elements"`
	MissingElements []string `json:"missing_elements"`
	ResponseTime    float64   `json:"response_time_seconds"`
	ErrorMessage    string    `json:"error_message,omitempty"`
	Timestamp       time.Time `json:"timestamp"`
}

// TestSuite representa uma suíte completa de testes
type TestSuite struct {
	BasicEligibilityTests []TestCase         `json:"basic_eligibility"`
	CalculationTests      []TestCase         `json:"calculation_scenarios"`
	EdgeCaseTests         []TestCase         `json:"edge_cases"`
	ConsistencyTests      []ConsistencyTest  `json:"consistency_checks"`
	QualityTests          []QualityTest      `json:"response_quality"`
	PerformanceTests      []PerformanceTest  `json:"performance_benchmarks"`
}

// ValidationSuite gerencia e executa testes de validação
type ValidationSuite struct {
	testSuitePath string
	testSuite     *TestSuite
	results       []TestResult
	knowledgeManager *KnowledgeManager
	feedbackSystem   *FeedbackSystem
}

// NewValidationSuite cria uma nova suíte de validação
func NewValidationSuite(testSuitePath string, km *KnowledgeManager, fs *FeedbackSystem) *ValidationSuite {
	return &ValidationSuite{
		testSuitePath:    testSuitePath,
		results:          make([]TestResult, 0),
		knowledgeManager: km,
		feedbackSystem:   fs,
	}
}

// LoadTestSuite carrega a suíte de testes
func (vs *ValidationSuite) LoadTestSuite() error {
	testCasesPath := filepath.Join(vs.testSuitePath, "test_cases.json")
	data, err := os.ReadFile(testCasesPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar casos de teste: %v", err)
	}

	var testData struct {
		TestCases TestSuite `json:"test_cases"`
	}

	if err := json.Unmarshal(data, &testData); err != nil {
		return fmt.Errorf("erro ao parsear casos de teste: %v", err)
	}

	vs.testSuite = &testData.TestCases
	return nil
}

// RunAllTests executa todos os testes da suíte
func (vs *ValidationSuite) RunAllTests() (map[string]interface{}, error) {
	if vs.testSuite == nil {
		return nil, fmt.Errorf("suíte de testes não carregada")
	}

	results := make(map[string]interface{})
	vs.results = make([]TestResult, 0)

	// Executar testes de elegibilidade básica
	eligibilityResults, err := vs.runEligibilityTests()
	if err != nil {
		return nil, fmt.Errorf("erro nos testes de elegibilidade: %v", err)
	}
	results["eligibility_tests"] = eligibilityResults

	// Executar testes de cálculo
	calculationResults, err := vs.runCalculationTests()
	if err != nil {
		return nil, fmt.Errorf("erro nos testes de cálculo: %v", err)
	}
	results["calculation_tests"] = calculationResults

	// Executar testes de casos extremos
	edgeResults, err := vs.runEdgeCaseTests()
	if err != nil {
		return nil, fmt.Errorf("erro nos testes de casos extremos: %v", err)
	}
	results["edge_case_tests"] = edgeResults

	// Executar testes de consistência
	consistencyResults, err := vs.runConsistencyTests()
	if err != nil {
		return nil, fmt.Errorf("erro nos testes de consistência: %v", err)
	}
	results["consistency_tests"] = consistencyResults

	// Executar testes de qualidade
	qualityResults, err := vs.runQualityTests()
	if err != nil {
		return nil, fmt.Errorf("erro nos testes de qualidade: %v", err)
	}
	results["quality_tests"] = qualityResults

	// Executar testes de performance
	performanceResults, err := vs.runPerformanceTests()
	if err != nil {
		return nil, fmt.Errorf("erro nos testes de performance: %v", err)
	}
	results["performance_tests"] = performanceResults

	// Gerar resumo geral
	summary := vs.generateTestSummary()
	results["summary"] = summary

	return results, nil
}

// runEligibilityTests executa testes de elegibilidade
func (vs *ValidationSuite) runEligibilityTests() (map[string]interface{}, error) {
	results := make(map[string]interface{})
	passed := 0
	total := len(vs.testSuite.BasicEligibilityTests)

	for _, test := range vs.testSuite.BasicEligibilityTests {
		result := vs.executeBasicTest(test)
		vs.results = append(vs.results, result)
		
		if result.Passed {
			passed++
		}
	}

	results["passed"] = passed
	results["total"] = total
	results["pass_rate"] = float64(passed) / float64(total)
	results["category"] = "eligibility"

	return results, nil
}

// runCalculationTests executa testes de cálculo
func (vs *ValidationSuite) runCalculationTests() (map[string]interface{}, error) {
	results := make(map[string]interface{})
	passed := 0
	total := len(vs.testSuite.CalculationTests)

	for _, test := range vs.testSuite.CalculationTests {
		result := vs.executeCalculationTest(test)
		vs.results = append(vs.results, result)
		
		if result.Passed {
			passed++
		}
	}

	results["passed"] = passed
	results["total"] = total
	results["pass_rate"] = float64(passed) / float64(total)
	results["category"] = "calculation"

	return results, nil
}

// runEdgeCaseTests executa testes de casos extremos
func (vs *ValidationSuite) runEdgeCaseTests() (map[string]interface{}, error) {
	results := make(map[string]interface{})
	passed := 0
	total := len(vs.testSuite.EdgeCaseTests)

	for _, test := range vs.testSuite.EdgeCaseTests {
		result := vs.executeBasicTest(test)
		vs.results = append(vs.results, result)
		
		if result.Passed {
			passed++
		}
	}

	results["passed"] = passed
	results["total"] = total
	results["pass_rate"] = float64(passed) / float64(total)
	results["category"] = "edge_cases"

	return results, nil
}

// runConsistencyTests executa testes de consistência
func (vs *ValidationSuite) runConsistencyTests() (map[string]interface{}, error) {
	results := make(map[string]interface{})
	passed := 0
	total := len(vs.testSuite.ConsistencyTests)

	for _, test := range vs.testSuite.ConsistencyTests {
		result := vs.executeConsistencyTest(test)
		vs.results = append(vs.results, result)
		
		if result.Passed {
			passed++
		}
	}

	results["passed"] = passed
	results["total"] = total
	results["pass_rate"] = float64(passed) / float64(total)
	results["category"] = "consistency"

	return results, nil
}

// runQualityTests executa testes de qualidade
func (vs *ValidationSuite) runQualityTests() (map[string]interface{}, error) {
	results := make(map[string]interface{})
	passed := 0
	total := len(vs.testSuite.QualityTests)

	for _, test := range vs.testSuite.QualityTests {
		result := vs.executeQualityTest(test)
		vs.results = append(vs.results, result)
		
		if result.Passed {
			passed++
		}
	}

	results["passed"] = passed
	results["total"] = total
	results["pass_rate"] = float64(passed) / float64(total)
	results["category"] = "quality"

	return results, nil
}

// runPerformanceTests executa testes de performance
func (vs *ValidationSuite) runPerformanceTests() (map[string]interface{}, error) {
	results := make(map[string]interface{})
	passed := 0
	total := len(vs.testSuite.PerformanceTests)

	for _, test := range vs.testSuite.PerformanceTests {
		result := vs.executePerformanceTest(test)
		vs.results = append(vs.results, result)
		
		if result.Passed {
			passed++
		}
	}

	results["passed"] = passed
	results["total"] = total
	results["pass_rate"] = float64(passed) / float64(total)
	results["category"] = "performance"

	return results, nil
}

// executeBasicTest executa um teste básico
func (vs *ValidationSuite) executeBasicTest(test TestCase) TestResult {
	start := time.Now()
	
	// Simular resposta do agente (na implementação real, seria chamada ao agente)
	response := vs.simulateAgentResponse(test.Question)
	
	responseTime := time.Since(start).Seconds()
	
	// Verificar se a resposta contém elementos esperados
	missingElements := make([]string, 0)
	for _, expected := range test.ExpectedAnswerContains {
		if !strings.Contains(strings.ToLower(response), strings.ToLower(expected)) {
			missingElements = append(missingElements, expected)
		}
	}

	passed := len(missingElements) == 0
	score := float64(len(test.ExpectedAnswerContains)-len(missingElements)) / float64(len(test.ExpectedAnswerContains))

	return TestResult{
		TestID:           test.ID,
		TestType:         "basic",
		Passed:           passed,
		Score:            score,
		ExpectedScore:    1.0,
		ActualResponse:   response,
		ExpectedElements: test.ExpectedAnswerContains,
		MissingElements:  missingElements,
		ResponseTime:     responseTime,
		Timestamp:        time.Now(),
	}
}

// executeCalculationTest executa um teste de cálculo
func (vs *ValidationSuite) executeCalculationTest(test TestCase) TestResult {
	start := time.Now()
	
	response := vs.simulateAgentResponse(test.Question)
	responseTime := time.Since(start).Seconds()
	
	// Verificar elementos esperados
	missingElements := make([]string, 0)
	for _, expected := range test.ExpectedAnswerContains {
		if !strings.Contains(strings.ToLower(response), strings.ToLower(expected)) {
			missingElements = append(missingElements, expected)
		}
	}

	// Verificar cálculo se especificado
	calculationCorrect := true
	if test.ExpectedCalculation != nil {
		// Na implementação real, verificaria se o valor calculado está correto
		expectedValueStr := fmt.Sprintf("%.2f", test.ExpectedCalculation.Valor)
		calculationCorrect = strings.Contains(response, expectedValueStr)
	}

	passed := len(missingElements) == 0 && calculationCorrect
	score := float64(len(test.ExpectedAnswerContains)-len(missingElements)) / float64(len(test.ExpectedAnswerContains))
	
	if !calculationCorrect {
		score *= 0.5 // Penalizar erro de cálculo
	}

	return TestResult{
		TestID:           test.ID,
		TestType:         "calculation",
		Passed:           passed,
		Score:            score,
		ExpectedScore:    1.0,
		ActualResponse:   response,
		ExpectedElements: test.ExpectedAnswerContains,
		MissingElements:  missingElements,
		ResponseTime:     responseTime,
		Timestamp:        time.Now(),
	}
}

// executeConsistencyTest executa um teste de consistência
func (vs *ValidationSuite) executeConsistencyTest(test ConsistencyTest) TestResult {
	start := time.Now()
	
	responses := make([]string, 0)
	
	// Executar a pergunta múltiplas vezes
	for i := 0; i < test.RepeatCount; i++ {
		response := vs.simulateAgentResponse(test.TestQuestion)
		responses = append(responses, response)
	}
	
	responseTime := time.Since(start).Seconds()
	
	// Calcular consistência
	consistency := vs.calculateConsistency(responses)
	
	passed := consistency >= test.ExpectedConsistency
	
	return TestResult{
		TestID:        test.ID,
		TestType:      "consistency",
		Passed:        passed,
		Score:         consistency,
		ExpectedScore: test.ExpectedConsistency,
		ActualResponse: fmt.Sprintf("Consistência: %.2f%% (%d execuções)", consistency*100, test.RepeatCount),
		ResponseTime:  responseTime,
		Timestamp:     time.Now(),
	}
}

// executeQualityTest executa um teste de qualidade
func (vs *ValidationSuite) executeQualityTest(test QualityTest) TestResult {
	start := time.Now()
	
	response := vs.simulateAgentResponse(test.Question)
	responseTime := time.Since(start).Seconds()
	
	// Avaliar critérios de qualidade
	qualityScore := vs.evaluateQualityCriteria(response, test.QualityCriteria)
	
	passed := qualityScore >= test.MinScore
	
	return TestResult{
		TestID:         test.ID,
		TestType:       "quality",
		Passed:         passed,
		Score:          qualityScore,
		ExpectedScore:  test.MinScore,
		ActualResponse: response,
		ResponseTime:   responseTime,
		Timestamp:      time.Now(),
	}
}

// executePerformanceTest executa um teste de performance
func (vs *ValidationSuite) executePerformanceTest(test PerformanceTest) TestResult {
	start := time.Now()
	
	response := vs.simulateAgentResponse(test.Question)
	responseTime := time.Since(start).Seconds()
	
	passed := responseTime <= test.MaxResponseTimeSeconds
	score := 1.0
	if responseTime > test.ExpectedAvgTime {
		score = test.ExpectedAvgTime / responseTime
	}

	return TestResult{
		TestID:         test.ID,
		TestType:       "performance",
		Passed:         passed,
		Score:          score,
		ExpectedScore:  1.0,
		ActualResponse: response,
		ResponseTime:   responseTime,
		Timestamp:      time.Now(),
	}
}

// simulateAgentResponse simula uma resposta do agente (para testes)
func (vs *ValidationSuite) simulateAgentResponse(question string) string {
	// Esta é uma simulação básica para testes
	// Na implementação real, seria uma chamada ao agente de IA

	questionLower := strings.ToLower(question)
	
	if strings.Contains(questionLower, "estagiário") {
		return "Não. Estagiários são expressamente excluídos do benefício do Vale Refeição conforme Política VR-003, independente da carga horária."
	}
	
	if strings.Contains(questionLower, "diretor") {
		return "Não. Colaboradores em cargos de diretoria são excluídos do benefício conforme Política VR-004."
	}
	
	if strings.Contains(questionLower, "admitido dia 20") {
		return "Para admissão após o dia 15, aplica-se cálculo proporcional (Política VR-005). Valor: R$ 191,05 para SINDPD."
	}
	
	if strings.Contains(questionLower, "20 dias de licença") {
		return "Licenças médicas superiores a 15 dias resultam em VR zero conforme Política VR-008. Valor: R$ 0,00."
	}
	
	return "Resposta simulada para teste: " + question
}

// calculateConsistency calcula consistência entre respostas
func (vs *ValidationSuite) calculateConsistency(responses []string) float64 {
	if len(responses) < 2 {
		return 1.0
	}

	// Implementação básica - na prática seria mais sofisticada
	identical := 0
	total := 0
	
	for i := 0; i < len(responses); i++ {
		for j := i + 1; j < len(responses); j++ {
			total++
			if responses[i] == responses[j] {
				identical++
			}
		}
	}

	if total == 0 {
		return 1.0
	}

	return float64(identical) / float64(total)
}

// evaluateQualityCriteria avalia critérios de qualidade de uma resposta
func (vs *ValidationSuite) evaluateQualityCriteria(response string, criteria map[string]bool) float64 {
	responseLower := strings.ToLower(response)
	passed := 0
	total := len(criteria)
	
	for criterion, expected := range criteria {
		switch criterion {
		case "cites_source":
			hasSource := strings.Contains(responseLower, "política") || strings.Contains(responseLower, "vr-")
			if hasSource == expected {
				passed++
			}
		case "provides_calculation":
			hasCalc := strings.Contains(responseLower, "r$") || strings.Contains(responseLower, "cálculo")
			if hasCalc == expected {
				passed++
			}
		case "uses_example":
			hasExample := strings.Contains(responseLower, "exemplo") || strings.Contains(responseLower, "caso")
			if hasExample == expected {
				passed++
			}
		case "maintains_confidentiality":
			hasName := strings.Contains(responseLower, "nome") && !strings.Contains(responseLower, "matrícula")
			if (!hasName) == expected {
				passed++
			}
		case "professional_language":
			// Verificação básica de profissionalismo
			isProfessional := !strings.Contains(responseLower, "tipo assim") && len(response) > 20
			if isProfessional == expected {
				passed++
			}
		}
	}

	return float64(passed) / float64(total)
}

// generateTestSummary gera resumo dos testes executados
func (vs *ValidationSuite) generateTestSummary() map[string]interface{} {
	summary := make(map[string]interface{})
	
	totalTests := len(vs.results)
	passedTests := 0
	totalScore := 0.0
	totalResponseTime := 0.0
	
	categoryResults := make(map[string]map[string]int)
	
	for _, result := range vs.results {
		if result.Passed {
			passedTests++
		}
		totalScore += result.Score
		totalResponseTime += result.ResponseTime
		
		// Agrupar por categoria
		if categoryResults[result.TestType] == nil {
			categoryResults[result.TestType] = make(map[string]int)
		}
		categoryResults[result.TestType]["total"]++
		if result.Passed {
			categoryResults[result.TestType]["passed"]++
		}
	}

	summary["total_tests"] = totalTests
	summary["passed_tests"] = passedTests
	summary["overall_pass_rate"] = float64(passedTests) / float64(totalTests)
	summary["average_score"] = totalScore / float64(totalTests)
	summary["average_response_time"] = totalResponseTime / float64(totalTests)
	summary["category_breakdown"] = categoryResults
	summary["execution_timestamp"] = time.Now()

	// Determinar status geral
	passRate := float64(passedTests) / float64(totalTests)
	if passRate >= 0.95 {
		summary["status"] = "EXCELLENT"
	} else if passRate >= 0.90 {
		summary["status"] = "GOOD"
	} else if passRate >= 0.80 {
		summary["status"] = "ACCEPTABLE"
	} else {
		summary["status"] = "NEEDS_IMPROVEMENT"
	}

	return summary
}

// SaveTestResults salva os resultados dos testes
func (vs *ValidationSuite) SaveTestResults() error {
	resultsPath := filepath.Join(vs.testSuitePath, "test_results.json")
	
	resultsData := map[string]interface{}{
		"execution_time": time.Now(),
		"summary":       vs.generateTestSummary(),
		"detailed_results": vs.results,
	}

	data, err := json.MarshalIndent(resultsData, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar resultados: %v", err)
	}

	return os.WriteFile(resultsPath, data, 0644)
}

// GetFailedTests retorna testes que falharam
func (vs *ValidationSuite) GetFailedTests() []TestResult {
	failed := make([]TestResult, 0)
	for _, result := range vs.results {
		if !result.Passed {
			failed = append(failed, result)
		}
	}
	return failed
}