package agent

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"
)

// TestResult representa o resultado de um teste individual
type TestResult struct {
	ScenarioID      string                 `json:"scenario_id"`
	Category        string                 `json:"category"`
	Question        string                 `json:"question"`
	Response        string                 `json:"response"`
	ResponseTime    time.Duration          `json:"response_time"`
	Score           float64                `json:"score"`
	Accuracy        bool                   `json:"accuracy"`
	Confidence      float64                `json:"confidence"`
	Passed          bool                   `json:"passed"`
	Errors          []string               `json:"errors"`
	Warnings        []string               `json:"warnings"`
	Metrics         map[string]interface{} `json:"metrics"`
	Timestamp       time.Time              `json:"timestamp"`
	ExpectedAnswer  string                 `json:"expected_answer"`
	ActualAnswer    string                 `json:"actual_answer"`
}

// ValidationReport representa o relatório completo de validação
type ValidationReport struct {
	TestSuite       string                 `json:"test_suite"`
	TotalScenarios  int                    `json:"total_scenarios"`
	PassedScenarios int                    `json:"passed_scenarios"`
	FailedScenarios int                    `json:"failed_scenarios"`
	OverallAccuracy float64                `json:"overall_accuracy"`
	AverageScore    float64                `json:"average_score"`
	AverageTime     time.Duration          `json:"average_response_time"`
	CategoryResults map[string]CategoryResult `json:"category_results"`
	QualityMetrics  QualityMetrics         `json:"quality_metrics"`
	PerformanceData PerformanceData        `json:"performance_data"`
	TestResults     []TestResult           `json:"test_results"`
	GeneratedAt     time.Time              `json:"generated_at"`
	Recommendations []string               `json:"recommendations"`
}

// CategoryResult representa resultados por categoria
type CategoryResult struct {
	Category        string    `json:"category"`
	TotalScenarios  int       `json:"total_scenarios"`
	PassedScenarios int       `json:"passed_scenarios"`
	Accuracy        float64   `json:"accuracy"`
	AverageScore    float64   `json:"average_score"`
	AverageTime     time.Duration `json:"average_time"`
	CommonErrors    []string  `json:"common_errors"`
}

// QualityMetrics representa métricas de qualidade
type QualityMetrics struct {
	Precision           float64 `json:"precision"`
	Recall              float64 `json:"recall"`
	F1Score             float64 `json:"f1_score"`
	ResponseCompleteness float64 `json:"response_completeness"`
	PolicyAdherence     float64 `json:"policy_adherence"`
	FormattingQuality   float64 `json:"formatting_quality"`
	UserSatisfaction    float64 `json:"user_satisfaction"`
}

// PerformanceData representa dados de performance
type PerformanceData struct {
	MinResponseTime   time.Duration `json:"min_response_time"`
	MaxResponseTime   time.Duration `json:"max_response_time"`
	MedianResponseTime time.Duration `json:"median_response_time"`
	P95ResponseTime   time.Duration `json:"p95_response_time"`
	TimeoutCount      int           `json:"timeout_count"`
	ErrorRate         float64       `json:"error_rate"`
	ThroughputQPS     float64       `json:"throughput_qps"`
}

// RealWorldValidator implementa validação de cenários reais
type RealWorldValidator struct {
	agent             *VRAgent
	scenarios         []RealWorldTestScenario
	enabledCategories []string
	maxTimeout        time.Duration
	strictMode        bool
	logger            Logger
}

// Logger interface para logging flexível
type Logger interface {
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}

// NewRealWorldValidator cria um novo validador
func NewRealWorldValidator(agent *VRAgent, logger Logger) *RealWorldValidator {
	return &RealWorldValidator{
		agent:             agent,
		scenarios:         LoadRealWorldDataset(),
		enabledCategories: []string{}, // Vazio = todas as categorias
		maxTimeout:        30 * time.Second,
		strictMode:        true,
		logger:            logger,
	}
}

// SetCategories define quais categorias testar
func (v *RealWorldValidator) SetCategories(categories []string) {
	v.enabledCategories = categories
}

// SetTimeout define timeout máximo para testes
func (v *RealWorldValidator) SetTimeout(timeout time.Duration) {
	v.maxTimeout = timeout
}

// SetStrictMode habilita/desabilita modo strict
func (v *RealWorldValidator) SetStrictMode(strict bool) {
	v.strictMode = strict
}

// RunAllTests executa todos os testes do dataset
func (v *RealWorldValidator) RunAllTests() (*ValidationReport, error) {
	v.logger.Println("Iniciando validação com cenários reais...")
	
	// Filtrar cenários se necessário
	scenariosToTest := v.getFilteredScenarios()
	
	var allResults []TestResult
	categoryStats := make(map[string]*CategoryResult)
	
	startTime := time.Now()
	
	for _, scenario := range scenariosToTest {
		v.logger.Printf("Testando cenário %s: %s", scenario.ID, scenario.Question)
		
		result := v.runSingleTest(scenario)
		allResults = append(allResults, result)
		
		// Atualizar estatísticas por categoria
		if _, exists := categoryStats[scenario.Category]; !exists {
			categoryStats[scenario.Category] = &CategoryResult{
				Category:        scenario.Category,
				CommonErrors:    []string{},
			}
		}
		
		cat := categoryStats[scenario.Category]
		cat.TotalScenarios++
		if result.Passed {
			cat.PassedScenarios++
		} else {
			// Coletar erros comuns
			for _, err := range result.Errors {
				cat.CommonErrors = append(cat.CommonErrors, err)
			}
		}
	}
	
	totalDuration := time.Since(startTime)
	
	// Calcular métricas globais
	report := v.generateValidationReport(allResults, categoryStats, totalDuration)
	
	v.logger.Printf("Validação concluída: %d/%d cenários passaram (%.1f%% accuracy)",
		report.PassedScenarios, report.TotalScenarios, report.OverallAccuracy*100)
	
	return report, nil
}

// runSingleTest executa um teste individual
func (v *RealWorldValidator) runSingleTest(scenario RealWorldTestScenario) TestResult {
	result := TestResult{
		ScenarioID:     scenario.ID,
		Category:       scenario.Category,
		Question:       scenario.Question,
		ExpectedAnswer: scenario.Expected.Answer,
		Errors:         []string{},
		Warnings:       []string{},
		Metrics:        make(map[string]interface{}),
		Timestamp:      time.Now(),
	}
	
	// Executar pergunta com timeout
	startTime := time.Now()
	
	// Canal para controlar timeout
	responseChan := make(chan string, 1)
	errorChan := make(chan error, 1)
	
	go func() {
		response, err := v.agent.Ask(scenario.Question)
		if err != nil {
			errorChan <- err
		} else {
			responseChan <- response
		}
	}()
	
	// Aguardar resposta ou timeout
	select {
	case response := <-responseChan:
		result.Response = response
		result.ResponseTime = time.Since(startTime)
		result.ActualAnswer = response
		
	case err := <-errorChan:
		result.ResponseTime = time.Since(startTime)
		result.Errors = append(result.Errors, fmt.Sprintf("Erro na consulta: %v", err))
		result.Passed = false
		return result
		
	case <-time.After(v.maxTimeout):
		result.ResponseTime = v.maxTimeout
		result.Errors = append(result.Errors, "Timeout na resposta")
		result.Passed = false
		return result
	}
	
	// Avaliar qualidade da resposta
	v.evaluateResponse(&result, scenario)
	
	return result
}

// evaluateResponse avalia a qualidade da resposta
func (v *RealWorldValidator) evaluateResponse(result *TestResult, scenario RealWorldTestScenario) {
	response := strings.ToLower(result.Response)
	
	var score float64
	var errors []string
	var warnings []string
	
	// 1. Verificar timeout
	if result.ResponseTime > scenario.Expected.MaxResponseTime {
		warnings = append(warnings, fmt.Sprintf("Tempo de resposta acima do esperado: %v > %v",
			result.ResponseTime, scenario.Expected.MaxResponseTime))
		score -= 0.1
	}
	
	// 2. Verificar conteúdo obrigatório
	mustContainScore := v.evaluateMustContain(response, scenario.Expected.MustContain)
	score += mustContainScore * 0.4
	
	if mustContainScore < 0.8 {
		errors = append(errors, "Resposta não contém palavras-chave obrigatórias")
	}
	
	// 3. Verificar conteúdo proibido
	mustNotContainScore := v.evaluateMustNotContain(response, scenario.Expected.MustNotContain)
	score += mustNotContainScore * 0.3
	
	if mustNotContainScore < 1.0 {
		errors = append(errors, "Resposta contém informações inadequadas")
	}
	
	// 4. Avaliar pontos-chave
	keyPointsScore := v.evaluateKeyPoints(response, scenario.Expected.KeyPoints)
	score += keyPointsScore * 0.3
	
	if keyPointsScore < 0.6 {
		warnings = append(warnings, "Resposta não aborda todos os pontos-chave")
	}
	
	// 5. Validar formato da resposta
	formattingScore := v.evaluateFormatting(result.Response)
	result.Metrics["formatting_score"] = formattingScore
	
	// 6. Calcular confiança baseada na consistência
	confidence := v.calculateConfidence(response, scenario)
	result.Confidence = confidence
	
	// Normalizar score (0-5 scale)
	result.Score = math.Max(0, math.Min(5, score*5))
	
	// Determinar se passou
	minScore := 4.0
	if !v.strictMode {
		minScore = 3.5
	}
	
	result.Accuracy = result.Score >= minScore && confidence >= scenario.Expected.MinConfidence
	result.Passed = result.Accuracy && len(errors) == 0
	
	result.Errors = errors
	result.Warnings = warnings
	
	// Métricas adicionais
	result.Metrics["must_contain_score"] = mustContainScore
	result.Metrics["must_not_contain_score"] = mustNotContainScore
	result.Metrics["key_points_score"] = keyPointsScore
	result.Metrics["confidence"] = confidence
}

// evaluateMustContain verifica se a resposta contém palavras obrigatórias
func (v *RealWorldValidator) evaluateMustContain(response string, mustContain []string) float64 {
	if len(mustContain) == 0 {
		return 1.0
	}
	
	found := 0
	for _, term := range mustContain {
		if strings.Contains(response, strings.ToLower(term)) {
			found++
		}
	}
	
	return float64(found) / float64(len(mustContain))
}

// evaluateMustNotContain verifica se a resposta não contém palavras proibidas
func (v *RealWorldValidator) evaluateMustNotContain(response string, mustNotContain []string) float64 {
	if len(mustNotContain) == 0 {
		return 1.0
	}
	
	violations := 0
	for _, term := range mustNotContain {
		if strings.Contains(response, strings.ToLower(term)) {
			violations++
		}
	}
	
	return math.Max(0, 1.0 - float64(violations)/float64(len(mustNotContain)))
}

// evaluateKeyPoints avalia se pontos-chave são abordados
func (v *RealWorldValidator) evaluateKeyPoints(response string, keyPoints []string) float64 {
	if len(keyPoints) == 0 {
		return 1.0
	}
	
	covered := 0
	for _, point := range keyPoints {
		// Busca por palavras-chave do ponto na resposta
		pointWords := strings.Fields(strings.ToLower(point))
		pointCoverage := 0
		
		for _, word := range pointWords {
			if len(word) > 3 && strings.Contains(response, word) {
				pointCoverage++
			}
		}
		
		// Se pelo menos 50% das palavras do ponto foram encontradas
		if float64(pointCoverage)/float64(len(pointWords)) >= 0.5 {
			covered++
		}
	}
	
	return float64(covered) / float64(len(keyPoints))
}

// evaluateFormatting avalia a qualidade da formatação
func (v *RealWorldValidator) evaluateFormatting(response string) float64 {
	score := 1.0
	
	// Verificar se resposta não está vazia
	if len(strings.TrimSpace(response)) == 0 {
		return 0.0
	}
	
	// Verificar comprimento adequado (não muito curto, não muito longo)
	length := len(response)
	if length < 20 {
		score -= 0.3
	} else if length > 1000 {
		score -= 0.2
	}
	
	// Verificar se tem pontuação adequada
	if !strings.Contains(response, ".") && !strings.Contains(response, "?") && !strings.Contains(response, "!") {
		score -= 0.2
	}
	
	return math.Max(0, score)
}

// calculateConfidence calcula confiança na resposta
func (v *RealWorldValidator) calculateConfidence(response string, scenario RealWorldTestScenario) float64 {
	confidence := 0.5 // Base confidence
	
	// Aumentar confiança se resposta contém referências a políticas
	if len(scenario.Expected.PolicyRefs) > 0 {
		for _, policyRef := range scenario.Expected.PolicyRefs {
			if strings.Contains(response, strings.ToLower(policyRef)) {
				confidence += 0.1
			}
		}
	}
	
	// Aumentar confiança para respostas específicas por tipo
	switch scenario.Expected.Type {
	case "policy":
		if strings.Contains(response, "política") || strings.Contains(response, "regra") {
			confidence += 0.2
		}
	case "calculation":
		if strings.Contains(response, "calcul") || strings.Contains(response, "fórmul") {
			confidence += 0.2
		}
	case "data_query":
		if strings.Contains(response, "processad") || strings.Contains(response, "dado") {
			confidence += 0.2
		}
	}
	
	return math.Min(1.0, confidence)
}

// getFilteredScenarios retorna cenários filtrados por categoria
func (v *RealWorldValidator) getFilteredScenarios() []RealWorldTestScenario {
	if len(v.enabledCategories) == 0 {
		return v.scenarios
	}
	
	var filtered []RealWorldTestScenario
	for _, scenario := range v.scenarios {
		for _, enabledCat := range v.enabledCategories {
			if scenario.Category == enabledCat {
				filtered = append(filtered, scenario)
				break
			}
		}
	}
	
	return filtered
}

// generateValidationReport gera relatório completo de validação
func (v *RealWorldValidator) generateValidationReport(results []TestResult, categoryStats map[string]*CategoryResult, totalDuration time.Duration) *ValidationReport {
	report := &ValidationReport{
		TestSuite:       "Real World Scenarios",
		TotalScenarios:  len(results),
		CategoryResults: make(map[string]CategoryResult),
		TestResults:     results,
		GeneratedAt:     time.Now(),
		Recommendations: []string{},
	}
	
	// Calcular estatísticas globais
	var totalScore float64
	var totalTime time.Duration
	var responseTimes []time.Duration
	
	for _, result := range results {
		if result.Passed {
			report.PassedScenarios++
		}
		totalScore += result.Score
		totalTime += result.ResponseTime
		responseTimes = append(responseTimes, result.ResponseTime)
	}
	
	report.FailedScenarios = report.TotalScenarios - report.PassedScenarios
	report.OverallAccuracy = float64(report.PassedScenarios) / float64(report.TotalScenarios)
	report.AverageScore = totalScore / float64(report.TotalScenarios)
	report.AverageTime = totalTime / time.Duration(report.TotalScenarios)
	
	// Calcular estatísticas por categoria
	for catName, catStats := range categoryStats {
		catStats.Accuracy = float64(catStats.PassedScenarios) / float64(catStats.TotalScenarios)
		
		// Calcular scores médios por categoria
		var catTotalScore float64
		var catTotalTime time.Duration
		var catCount int
		
		for _, result := range results {
			if result.Category == catName {
				catTotalScore += result.Score
				catTotalTime += result.ResponseTime
				catCount++
			}
		}
		
		if catCount > 0 {
			catStats.AverageScore = catTotalScore / float64(catCount)
			catStats.AverageTime = catTotalTime / time.Duration(catCount)
		}
		
		report.CategoryResults[catName] = *catStats
	}
	
	// Calcular métricas de qualidade
	report.QualityMetrics = v.calculateQualityMetrics(results)
	
	// Calcular dados de performance
	report.PerformanceData = v.calculatePerformanceData(results, responseTimes, totalDuration)
	
	// Gerar recomendações
	report.Recommendations = v.generateRecommendations(report)
	
	return report
}

// calculateQualityMetrics calcula métricas de qualidade
func (v *RealWorldValidator) calculateQualityMetrics(results []TestResult) QualityMetrics {
	var precision, recall, completeness, adherence, formatting, satisfaction float64
	
	validResults := 0
	for _, result := range results {
		if result.Score > 0 {
			validResults++
			precision += result.Score / 5.0 // Normalizar para 0-1
			
			if formattingScore, ok := result.Metrics["formatting_score"]; ok {
				if fs, ok := formattingScore.(float64); ok {
					formatting += fs
				}
			}
			
			if result.Confidence > 0.7 {
				adherence += 1.0
			}
			
			if len(result.Response) > 50 {
				completeness += 1.0
			}
			
			if result.Score >= 4.0 {
				satisfaction += 1.0
			}
		}
	}
	
	if validResults > 0 {
		precision /= float64(validResults)
		recall = precision // Simplificação - em um sistema real, seria calculado diferente
		completeness /= float64(validResults)
		adherence /= float64(validResults)
		formatting /= float64(validResults)
		satisfaction /= float64(validResults)
	}
	
	f1Score := 2 * (precision * recall) / (precision + recall)
	if math.IsNaN(f1Score) {
		f1Score = 0
	}
	
	return QualityMetrics{
		Precision:           precision,
		Recall:              recall,
		F1Score:             f1Score,
		ResponseCompleteness: completeness,
		PolicyAdherence:     adherence,
		FormattingQuality:   formatting,
		UserSatisfaction:    satisfaction,
	}
}

// calculatePerformanceData calcula dados de performance
func (v *RealWorldValidator) calculatePerformanceData(results []TestResult, responseTimes []time.Duration, totalDuration time.Duration) PerformanceData {
	if len(responseTimes) == 0 {
		return PerformanceData{}
	}
	
	// Ordenar tempos para calcular mediana e percentis
	for i := 0; i < len(responseTimes); i++ {
		for j := i + 1; j < len(responseTimes); j++ {
			if responseTimes[i] > responseTimes[j] {
				responseTimes[i], responseTimes[j] = responseTimes[j], responseTimes[i]
			}
		}
	}
	
	var errorCount int
	for _, result := range results {
		if len(result.Errors) > 0 {
			errorCount++
		}
	}
	
	return PerformanceData{
		MinResponseTime:    responseTimes[0],
		MaxResponseTime:    responseTimes[len(responseTimes)-1],
		MedianResponseTime: responseTimes[len(responseTimes)/2],
		P95ResponseTime:    responseTimes[int(float64(len(responseTimes))*0.95)],
		TimeoutCount:       0, // Será implementado se necessário
		ErrorRate:          float64(errorCount) / float64(len(results)),
		ThroughputQPS:      float64(len(results)) / totalDuration.Seconds(),
	}
}

// generateRecommendations gera recomendações baseadas nos resultados
func (v *RealWorldValidator) generateRecommendations(report *ValidationReport) []string {
	var recommendations []string
	
	// Recomendações baseadas em accuracy
	if report.OverallAccuracy < 0.95 {
		recommendations = append(recommendations, "Accuracy está abaixo da meta de 95%. Revisar classificação de perguntas e templates de resposta.")
	}
	
	// Recomendações baseadas em performance
	if report.AverageTime > 2*time.Second {
		recommendations = append(recommendations, "Tempo de resposta médio acima da meta. Considerar otimizações de cache e processamento.")
	}
	
	// Recomendações baseadas em quality metrics
	if report.QualityMetrics.F1Score < 0.8 {
		recommendations = append(recommendations, "F1 Score baixo indica necessidade de melhorar precisão e recall das respostas.")
	}
	
	if report.QualityMetrics.PolicyAdherence < 0.9 {
		recommendations = append(recommendations, "Baixa aderência a políticas. Expandir base de conhecimento e melhorar citações de fontes.")
	}
	
	// Recomendações por categoria
	for catName, catResult := range report.CategoryResults {
		if catResult.Accuracy < 0.9 {
			recommendations = append(recommendations, fmt.Sprintf("Categoria '%s' com baixa accuracy (%.1f%%). Revisar cenários específicos.", catName, catResult.Accuracy*100))
		}
	}
	
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Todos os critérios de qualidade foram atendidos. Monitorar performance continuamente.")
	}
	
	return recommendations
}

// ExportReport exporta relatório em formato JSON
func (v *RealWorldValidator) ExportReport(report *ValidationReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}