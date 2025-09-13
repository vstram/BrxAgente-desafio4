// Package agent implements Issue #71: Real World Testing and Validation
// This file contains the main test implementation according to the requirements
package agent

import (
	"strings"
	"testing"
	"time"

	"BrxAgente-desafio4/internal/chat"
	"BrxAgente-desafio4/internal/config"
)

// TestIssue71_RealWorldValidation implementa a Issue #71 conforme especificação
// Este é o teste principal que executa todos os cenários reais e coleta métricas de qualidade
func TestIssue71_RealWorldValidation(t *testing.T) {
	// Skip em testes rápidos pois este é um teste abrangente
	if testing.Short() {
		t.Skip("Skipping comprehensive real-world validation in short mode")
	}

	t.Log("🚀 Iniciando implementação da Issue #71: Teste com Casos Reais e Validação")

	// Criar suite de testes
	suite, err := NewRealWorldTestSuite()
	if err != nil {
		t.Fatalf("❌ Erro ao criar suite de testes: %v", err)
	}

	// Configurar suite com parâmetros da Issue #71
	config := &TestSuiteConfig{
		EnabledCategories:     []string{}, // Todas as categorias
		AccuracyThreshold:     0.95,       // Meta: > 95%
		SatisfactionThreshold: 4.2,        // Meta: > 4.2/5.0
		MaxResponseTime:       2 * time.Second, // Meta: < 2 segundos
		StrictMode:            true,
		RunBenchmarks:         true,
		GenerateReports:       true,
		OutputDir:            "./test_reports/issue71",
	}
	suite.SetConfig(config)

	// Executar validação completa
	t.Log("📋 Executando validação completa com dataset de cenários reais...")
	report, err := suite.RunCompleteValidation()
	if err != nil {
		t.Fatalf("❌ Falha na validação completa: %v", err)
	}

	// Validar resultados conforme critérios da Issue #71
	t.Run("ValidateAccuracyTarget", func(t *testing.T) {
		validateAccuracyTarget(t, report, config.AccuracyThreshold)
	})

	t.Run("ValidatePerformanceTarget", func(t *testing.T) {
		validatePerformanceTarget(t, report, config.MaxResponseTime)
	})

	t.Run("ValidateSatisfactionTarget", func(t *testing.T) {
		validateSatisfactionTarget(t, report, config.SatisfactionThreshold)
	})

	t.Run("ValidateCategoryPerformance", func(t *testing.T) {
		validateCategoryPerformance(t, report)
	})

	t.Run("ValidateEdgeCases", func(t *testing.T) {
		validateEdgeCases(t, report)
	})

	// Log de resumo final
	logTestSummary(t, report)

	// Verificar se passou nos critérios gerais da Issue #71
	if !passedOverallCriteria(report, config) {
		t.Errorf("❌ Sistema não atendeu aos critérios de aceite da Issue #71")
		logFailureDetails(t, report, config)
	} else {
		t.Log("✅ Todos os critérios de aceite da Issue #71 foram atendidos!")
	}
}

// validateAccuracyTarget valida se a meta de accuracy foi atingida
func validateAccuracyTarget(t *testing.T, report *ComprehensiveReport, threshold float64) {
	if report.ValidationReport == nil {
		t.Fatal("Relatório de validação não disponível")
	}

	accuracy := report.ValidationReport.OverallAccuracy
	t.Logf("📊 Accuracy obtida: %.2f%% (Meta: %.2f%%)", accuracy*100, threshold*100)

	if accuracy < threshold {
		t.Errorf("❌ Accuracy abaixo da meta: %.2f%% < %.2f%%", accuracy*100, threshold*100)
		
		// Detalhar falhas por categoria
		for category, categoryResult := range report.ValidationReport.CategoryResults {
			if categoryResult.Accuracy < threshold {
				t.Logf("   📋 Categoria '%s': %.1f%% (abaixo da meta)", category, categoryResult.Accuracy*100)
			}
		}
	} else {
		t.Logf("✅ Meta de accuracy atingida: %.2f%% >= %.2f%%", accuracy*100, threshold*100)
	}
}

// validatePerformanceTarget valida se a meta de performance foi atingida
func validatePerformanceTarget(t *testing.T, report *ComprehensiveReport, maxResponseTime time.Duration) {
	// Verificar tempo médio de resposta
	if report.ValidationReport != nil {
		avgTime := report.ValidationReport.AverageTime
		t.Logf("⚡ Tempo médio de resposta: %v (Meta: < %v)", avgTime, maxResponseTime)

		if avgTime > maxResponseTime {
			t.Errorf("❌ Tempo médio acima da meta: %v > %v", avgTime, maxResponseTime)
		} else {
			t.Logf("✅ Meta de tempo de resposta atingida: %v <= %v", avgTime, maxResponseTime)
		}
	}

	// Verificar performance dos benchmarks
	if report.BenchmarkReport != nil {
		rps := report.BenchmarkReport.Summary.AverageRPS
		grade := report.BenchmarkReport.Summary.PerformanceGrade
		
		t.Logf("🏎️ Performance: %.1f RPS (Nota: %s)", rps, grade)

		// Meta mínima: 10 RPS para nota >= C
		if rps < 10.0 {
			t.Errorf("❌ Performance abaixo do esperado: %.1f RPS < 10.0 RPS", rps)
		} else {
			t.Logf("✅ Performance adequada: %.1f RPS >= 10.0 RPS", rps)
		}
	}
}

// validateSatisfactionTarget valida se a meta de satisfação foi atingida
func validateSatisfactionTarget(t *testing.T, report *ComprehensiveReport, threshold float64) {
	satisfaction := report.QualityReport.AverageSatisfaction
	t.Logf("😊 Satisfação média: %.1f/5.0 (Meta: >= %.1f)", satisfaction, threshold)

	if satisfaction < threshold {
		t.Errorf("❌ Satisfação abaixo da meta: %.1f < %.1f", satisfaction, threshold)
	} else {
		t.Logf("✅ Meta de satisfação atingida: %.1f >= %.1f", satisfaction, threshold)
	}
}

// validateCategoryPerformance valida performance por categoria
func validateCategoryPerformance(t *testing.T, report *ComprehensiveReport) {
	t.Log("📋 Validando performance por categoria:")

	requiredCategories := []string{
		"Políticas de Elegibilidade",
		"Cálculos Específicos", 
		"Cenários Complexos",
		"Dados Processados",
	}

	for _, category := range requiredCategories {
		result, found := report.CategoryResults[category]
		if !found {
			t.Errorf("❌ Categoria '%s' não foi testada", category)
			continue
		}

		t.Logf("   📊 %s:", category)
		t.Logf("      - Accuracy: %.1f%% (%d/%d cenários)", 
			result.Accuracy*100, result.PassedScenarios, result.TotalScenarios)
		t.Logf("      - Score médio: %.1f/5.0", result.AverageScore)
		t.Logf("      - Tempo médio: %v", result.AverageResponseTime)

		// Meta por categoria: pelo menos 85% de accuracy
		if result.Accuracy < 0.85 {
			t.Errorf("❌ Categoria '%s' com accuracy baixa: %.1f%% < 85%%", 
				category, result.Accuracy*100)
		}

		// Meta de score: pelo menos 3.5/5.0
		if result.AverageScore < 3.5 {
			t.Errorf("❌ Categoria '%s' com score baixo: %.1f < 3.5", 
				category, result.AverageScore)
		}
	}
}

// validateEdgeCases valida se casos edge foram tratados adequadamente
func validateEdgeCases(t *testing.T, report *ComprehensiveReport) {
	if report.EdgeCaseResults == nil {
		t.Error("❌ Resultados de casos edge não disponíveis")
		return
	}

	success := report.EdgeCaseResults.OverallSuccess
	t.Logf("🔧 Casos edge: %.1f%% de sucesso", success*100)

	for _, test := range report.EdgeCaseResults.TestCases {
		status := "❌"
		if test.Passed {
			status = "✅"
		}
		t.Logf("   %s %s: %s (%.2fs)", status, test.Name, test.Description, test.Duration.Seconds())
	}

	// Meta: pelo menos 70% dos casos edge devem passar
	if success < 0.7 {
		t.Errorf("❌ Taxa de sucesso em casos edge baixa: %.1f%% < 70%%", success*100)
	} else {
		t.Logf("✅ Casos edge tratados adequadamente: %.1f%% >= 70%%", success*100)
	}
}

// passedOverallCriteria verifica se todos os critérios principais foram atendidos
func passedOverallCriteria(report *ComprehensiveReport, config *TestSuiteConfig) bool {
	// 1. Accuracy >= threshold
	if report.ValidationReport == nil || report.ValidationReport.OverallAccuracy < config.AccuracyThreshold {
		return false
	}

	// 2. Performance adequada (tempo de resposta e RPS)
	if report.ValidationReport.AverageTime > config.MaxResponseTime {
		return false
	}

	if report.BenchmarkReport != nil && report.BenchmarkReport.Summary.AverageRPS < 5.0 {
		return false // Performance mínima aceitável
	}

	// 3. Satisfação >= threshold  
	if report.QualityReport.AverageSatisfaction < config.SatisfactionThreshold {
		return false
	}

	// 4. Casos edge com taxa mínima de sucesso
	if report.EdgeCaseResults != nil && report.EdgeCaseResults.OverallSuccess < 0.5 {
		return false
	}

	// 5. Pelo menos 3 das 4 categorias principais com accuracy >= 80%
	goodCategories := 0
	requiredCategories := []string{
		"Políticas de Elegibilidade",
		"Cálculos Específicos", 
		"Cenários Complexos", 
		"Dados Processados",
	}

	for _, category := range requiredCategories {
		if result, found := report.CategoryResults[category]; found && result.Accuracy >= 0.8 {
			goodCategories++
		}
	}

	if goodCategories < 3 {
		return false
	}

	return true
}

// logTestSummary registra resumo completo do teste
func logTestSummary(t *testing.T, report *ComprehensiveReport) {
	t.Log("\n" + strings.Repeat("=", 80))
	t.Log("📊 RESUMO FINAL - ISSUE #71: TESTE COM CASOS REAIS E VALIDAÇÃO")
	t.Log(strings.Repeat("=", 80))

	// Informações gerais
	t.Logf("⏱️  Duração total: %v", report.TotalDuration)
	t.Logf("🏆 Nota final: %s", report.FinalAnalysis.OverallGrade)

	// Métricas principais
	if report.ValidationReport != nil {
		t.Logf("📈 Accuracy geral: %.1f%% (%d/%d cenários)", 
			report.ValidationReport.OverallAccuracy*100,
			report.ValidationReport.PassedScenarios,
			report.ValidationReport.TotalScenarios)
		t.Logf("⚡ Tempo médio de resposta: %v", report.ValidationReport.AverageTime)
	}

	if report.BenchmarkReport != nil {
		t.Logf("🚀 Performance: %.1f RPS (Nota: %s)",
			report.BenchmarkReport.Summary.AverageRPS,
			report.BenchmarkReport.Summary.PerformanceGrade)
	}

	t.Logf("😊 Satisfação média: %.1f/5.0", report.QualityReport.AverageSatisfaction)

	if report.EdgeCaseResults != nil {
		t.Logf("🔧 Casos edge: %.1f%% sucesso", report.EdgeCaseResults.OverallSuccess*100)
	}

	// Resultados por categoria
	t.Log("\n📋 RESULTADOS POR CATEGORIA:")
	for category, result := range report.CategoryResults {
		t.Logf("   %s: %.1f%% accuracy, %.1f score, %v tempo médio",
			category, result.Accuracy*100, result.AverageScore, result.AverageResponseTime)
	}

	// Conquistas e issues
	if len(report.FinalAnalysis.Achievements) > 0 {
		t.Log("\n🎉 CONQUISTAS:")
		for _, achievement := range report.FinalAnalysis.Achievements {
			t.Logf("   %s", achievement)
		}
	}

	if len(report.FinalAnalysis.CriticalIssues) > 0 {
		t.Log("\n❗ ISSUES CRÍTICAS:")
		for _, issue := range report.FinalAnalysis.CriticalIssues {
			t.Logf("   %s", issue)
		}
	}

	// Recomendações
	t.Log("\n📝 RECOMENDAÇÕES:")
	for _, recommendation := range report.FinalAnalysis.Recommendations {
		t.Logf("   %s", recommendation)
	}

	t.Log(strings.Repeat("=", 80))
}

// logFailureDetails registra detalhes das falhas
func logFailureDetails(t *testing.T, report *ComprehensiveReport, config *TestSuiteConfig) {
	t.Log("\n❌ DETALHES DAS FALHAS:")

	if report.ValidationReport != nil && report.ValidationReport.OverallAccuracy < config.AccuracyThreshold {
		t.Logf("   📊 Accuracy insuficiente: %.1f%% < %.1f%%",
			report.ValidationReport.OverallAccuracy*100, config.AccuracyThreshold*100)
	}

	if report.ValidationReport != nil && report.ValidationReport.AverageTime > config.MaxResponseTime {
		t.Logf("   ⚡ Tempo de resposta alto: %v > %v", 
			report.ValidationReport.AverageTime, config.MaxResponseTime)
	}

	if report.QualityReport.AverageSatisfaction < config.SatisfactionThreshold {
		t.Logf("   😞 Satisfação baixa: %.1f < %.1f",
			report.QualityReport.AverageSatisfaction, config.SatisfactionThreshold)
	}

	// Categorias com problemas
	for category, result := range report.CategoryResults {
		if result.Accuracy < 0.8 {
			t.Logf("   📋 Categoria '%s' com problemas: %.1f%% accuracy", 
				category, result.Accuracy*100)
		}
	}
}

// TestRealWorldDatasetLoad testa carregamento do dataset
func TestRealWorldDatasetLoad(t *testing.T) {
	t.Log("📚 Testando carregamento do dataset de cenários reais...")

	scenarios := LoadRealWorldDataset()
	if len(scenarios) == 0 {
		t.Fatal("❌ Dataset de cenários reais está vazio")
	}

	t.Logf("✅ Dataset carregado com %d cenários", len(scenarios))

	// Verificar se todas as categorias estão presentes
	categories := make(map[string]int)
	for _, scenario := range scenarios {
		categories[scenario.Category]++
	}

	expectedCategories := []string{
		"Políticas de Elegibilidade",
		"Cálculos Específicos",
		"Cenários Complexos",
		"Dados Processados",
	}

	for _, expectedCat := range expectedCategories {
		count, found := categories[expectedCat]
		if !found {
			t.Errorf("❌ Categoria '%s' não encontrada no dataset", expectedCat)
		} else {
			t.Logf("   📋 %s: %d cenários", expectedCat, count)
		}
	}

	// Verificar estrutura dos cenários
	for i, scenario := range scenarios[:3] { // Verificar apenas os primeiros 3
		if scenario.ID == "" {
			t.Errorf("❌ Cenário %d sem ID", i)
		}
		if scenario.Question == "" {
			t.Errorf("❌ Cenário %d sem pergunta", i)
		}
		if scenario.Expected.Answer == "" {
			t.Errorf("❌ Cenário %d sem resposta esperada", i)
		}
	}

	t.Log("✅ Estrutura do dataset validada")
}

// TestAgentBasicFunctionality testa funcionalidade básica do agente
func TestAgentBasicFunctionality(t *testing.T) {
	t.Log("🤖 Testando funcionalidade básica do agente...")

	// Setup básico
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("❌ Erro ao criar agente: %v", err)
	}

	// Teste básico de pergunta
	question := "Diretores têm direito a VR?"
	t.Logf("❓ Pergunta: %s", question)

	start := time.Now()
	response, err := agent.Ask(question)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("❌ Erro na consulta: %v", err)
	} else {
		t.Logf("✅ Resposta recebida em %v", duration)
		t.Logf("💬 Resposta: %.100s...", response)
	}

	// Verificar se resposta não está vazia
	if len(response) == 0 {
		t.Error("❌ Resposta vazia")
	}

	// Verificar tempo de resposta razoável
	if duration > 10*time.Second {
		t.Errorf("❌ Tempo de resposta muito alto: %v", duration)
	}

	t.Log("✅ Funcionalidade básica do agente validada")
}