package intelligence

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestAnomalyDetector testa funcionalidades básicas do detector
func TestAnomalyDetector(t *testing.T) {
	t.Run("CreateDetector", testCreateDetector)
	t.Run("RegisterRules", testRegisterRules)
	t.Run("DetectValueAnomalies", testDetectValueAnomalies)
	t.Run("DetectTemporalAnomalies", testDetectTemporalAnomalies)
	t.Run("DetectRelationshipAnomalies", testDetectRelationshipAnomalies)
	t.Run("GenerateReport", testGenerateReport)
}

func testCreateDetector(t *testing.T) {
	config := DefaultAnalysisConfig()
	detector := NewAnomalyDetector(config)

	if detector == nil {
		t.Fatal("Detector não deveria ser nil")
	}

	stats := detector.GetStats()
	if stats["total_rules"].(int) == 0 {
		t.Error("Detector deveria ter regras padrão registradas")
	}

	if stats["enabled_rules"].(int) == 0 {
		t.Error("Detector deveria ter regras habilitadas")
	}

	t.Logf("Detector criado com %d regras", stats["total_rules"])
}

func testRegisterRules(t *testing.T) {
	detector := NewAnomalyDetector(nil)
	initialRules := detector.GetStats()["total_rules"].(int)

	// Adicionar regra customizada
	customRule := AnomalyRule{
		Name:        "test_rule",
		Description: "Regra de teste",
		Type:        AnomalyTypeValue,
		Severity:    SeverityMedium,
		Enabled:     true,
		Detector: func(ctx *AnalysisContext) []Anomaly {
			return []Anomaly{}
		},
	}

	detector.AddRule(customRule)

	newStats := detector.GetStats()
	if newStats["total_rules"].(int) != initialRules+1 {
		t.Errorf("Esperava %d regras, mas tem %d", initialRules+1, newStats["total_rules"])
	}

	// Remover regra
	detector.RemoveRule("test_rule")

	finalStats := detector.GetStats()
	if finalStats["total_rules"].(int) != initialRules {
		t.Errorf("Esperava %d regras após remoção, mas tem %d", initialRules, finalStats["total_rules"])
	}
}

func testDetectValueAnomalies(t *testing.T) {
	detector := NewAnomalyDetector(nil)

	// Criar dados de teste com anomalias de valor
	colaboradores := map[string]interface{}{
		"001": map[string]interface{}{
			"matricula":  "001",
			"sindicato":  "SINDICATO_A",
			"vr_total":   500.0, // Valor normal
			"dias_uteis": 22,
		},
		"002": map[string]interface{}{
			"matricula":  "002",
			"sindicato":  "SINDICATO_A",
			"vr_total":   520.0, // Valor normal
			"dias_uteis": 23,
		},
		"003": map[string]interface{}{
			"matricula":  "003",
			"sindicato":  "SINDICATO_A",
			"vr_total":   510.0, // Valor normal
			"dias_uteis": 22,
		},
		"004": map[string]interface{}{
			"matricula":  "004",
			"sindicato":  "SINDICATO_A",
			"vr_total":   530.0, // Valor normal
			"dias_uteis": 23,
		},
		"005": map[string]interface{}{
			"matricula":  "005",
			"sindicato":  "SINDICATO_A",
			"vr_total":   1500.0, // Outlier alto
			"dias_uteis": 22,
		},
		"006": map[string]interface{}{
			"matricula":  "006",
			"sindicato":  "SINDICATO_A",
			"vr_total":   0.0, // Valor zero (anomalia)
			"dias_uteis": 22,
		},
		"007": map[string]interface{}{
			"matricula":  "007",
			"sindicato":  "SINDICATO_A",
			"vr_total":   -100.0, // Valor negativo (anomalia crítica)
			"dias_uteis": 22,
		},
	}

	ctx := NewAnalysisContext(colaboradores, nil)
	report, err := detector.DetectAnomalies(ctx)

	if err != nil {
		t.Fatalf("Erro na detecção de anomalias: %v", err)
	}

	if report.TotalAnomalies == 0 {
		t.Error("Deveria ter detectado anomalias nos dados de teste")
	}

	// Verificar se detectou VR zero
	foundZeroVR := false
	foundNegativeVR := false
	foundOutlier := false

	for _, anomaly := range report.Anomalies {
		t.Logf("Anomalia detectada: Entity=%s, Type=%s, Rule=%s, Field=%s, Title=%s",
			anomaly.Entity, anomaly.Type, anomaly.RuleName, anomaly.FieldName, anomaly.Title)

		if anomaly.Entity == "006" && anomaly.FieldName == "vr_value" {
			foundZeroVR = true
		}
		if anomaly.Entity == "007" && anomaly.FieldName == "vr_value" {
			foundNegativeVR = true
			if anomaly.Severity != SeverityCritical {
				t.Error("VR negativo deveria ter severidade crítica")
			}
		}
		if anomaly.Entity == "005" && anomaly.Type == AnomalyTypeValue {
			foundOutlier = true
		}
	}

	if !foundZeroVR {
		t.Error("Deveria ter detectado VR zero")
	}

	if !foundNegativeVR {
		t.Error("Deveria ter detectado VR negativo")
	}

	if !foundOutlier {
		t.Error("Deveria ter detectado outlier de VR")
	}

	t.Logf("Detectadas %d anomalias de valor", report.AnomaliesByType[AnomalyTypeValue])
}

func testDetectTemporalAnomalies(t *testing.T) {
	detector := NewAnomalyDetector(nil)

	futureDate := time.Now().AddDate(1, 0, 0) // 1 ano no futuro
	pastDate := time.Date(1985, 1, 1, 0, 0, 0, 0, time.UTC)

	// Criar dados de teste com anomalias temporais
	colaboradores := map[string]interface{}{
		"001": map[string]interface{}{
			"matricula":     "001",
			"sindicato":     "SINDICATO_A",
			"data_admissao": futureDate.Format("2006-01-02"), // Data futura
			"vr_total":      500.0,
		},
		"002": map[string]interface{}{
			"matricula":     "002",
			"sindicato":     "SINDICATO_A",
			"data_admissao": pastDate.Format("2006-01-02"), // Data muito antiga
			"vr_total":      500.0,
		},
		"003": map[string]interface{}{
			"matricula":         "003",
			"sindicato":         "SINDICATO_A",
			"data_admissao":     "2023-01-01",
			"data_desligamento": "2022-12-31", // Desligamento antes da admissão
			"vr_total":          500.0,
		},
		"004": map[string]interface{}{
			"matricula":     "004",
			"sindicato":     "SINDICATO_A",
			"data_admissao": "2023-01-01",
			"dias_uteis":    -5, // Dias úteis negativos
			"vr_total":      500.0,
		},
	}

	ctx := NewAnalysisContext(colaboradores, nil)
	report, err := detector.DetectAnomalies(ctx)

	if err != nil {
		t.Fatalf("Erro na detecção de anomalias: %v", err)
	}

	temporalAnomalies := report.AnomaliesByType[AnomalyTypeTemporal]
	if temporalAnomalies == 0 {
		t.Error("Deveria ter detectado anomalias temporais")
	}

	// Verificar tipos específicos de anomalias temporais
	foundFutureDate := false
	foundInvalidSequence := false

	for _, anomaly := range report.Anomalies {
		if anomaly.Type == AnomalyTypeTemporal {
			if anomaly.Entity == "001" && anomaly.RuleName == "future_dates" {
				foundFutureDate = true
				if anomaly.Severity != SeverityCritical {
					t.Error("Data futura deveria ter severidade crítica")
				}
			}
			if anomaly.Entity == "003" && anomaly.RuleName == "invalid_date_sequence" {
				foundInvalidSequence = true
			}
		}
	}

	if !foundFutureDate {
		t.Error("Deveria ter detectado data futura")
	}

	if !foundInvalidSequence {
		t.Error("Deveria ter detectado sequência de datas inválida")
	}

	t.Logf("Detectadas %d anomalias temporais", temporalAnomalies)
}

func testDetectRelationshipAnomalies(t *testing.T) {
	detector := NewAnomalyDetector(nil)

	// Criar dados de teste com anomalias de relacionamento
	colaboradores := map[string]interface{}{
		"001_planilha1": map[string]interface{}{
			"matricula": "001", // Matrícula que será duplicada
			"sindicato": "SINDICATO_A",
			"vr_total":  500.0,
			"fonte":     "planilha1",
		},
		"001_planilha2": map[string]interface{}{ // Matrícula duplicada
			"matricula": "001", // Mesma matrícula
			"sindicato": "SINDICATO_B",
			"vr_total":  600.0,
			"fonte":     "planilha2",
		},
		"002": map[string]interface{}{
			"matricula":         "002",
			"sindicato":         "SINDICATO_INEXISTENTE", // Sindicato não reconhecido
			"vr_total":          500.0,
			"status":            "ativo",
			"data_desligamento": "2023-12-31", // Status inconsistente
		},
		"003": map[string]interface{}{
			// "matricula" ausente - campo obrigatório
			"sindicato": "SINDICATO_A",
			"vr_total":  500.0,
		},
		"004": map[string]interface{}{
			"matricula":    "004",
			"sindicato":    "SINDICATO_A",
			"vr_calculado": 500.0,
			"vr_esperado":  600.0, // Inconsistência de valores
		},
	}

	ctx := NewAnalysisContext(colaboradores, nil)
	report, err := detector.DetectAnomalies(ctx)

	if err != nil {
		t.Fatalf("Erro na detecção de anomalias: %v", err)
	}

	relationshipAnomalies := report.AnomaliesByType[AnomalyTypeRelationship]
	if relationshipAnomalies == 0 {
		t.Error("Deveria ter detectado anomalias de relacionamento")
	}

	// Verificar tipos específicos
	foundDuplicate := false

	for _, anomaly := range report.Anomalies {
		if anomaly.Type == AnomalyTypeRelationship {
			if anomaly.RuleName == "duplicate_matricula" {
				foundDuplicate = true
				if anomaly.Severity != SeverityCritical {
					t.Error("Matrícula duplicada deveria ter severidade crítica")
				}
			}
		}
	}

	if !foundDuplicate {
		t.Error("Deveria ter detectado matrícula duplicada")
	}

	t.Logf("Detectadas %d anomalias de relacionamento", relationshipAnomalies)
}

func testGenerateReport(t *testing.T) {
	detector := NewAnomalyDetector(nil)

	// Dados simples para teste de relatório
	colaboradores := map[string]interface{}{
		"001": map[string]interface{}{
			"matricula": "001",
			"sindicato": "SINDICATO_A",
			"vr_total":  500.0,
		},
	}

	ctx := NewAnalysisContext(colaboradores, nil)
	report, err := detector.DetectAnomalies(ctx)

	if err != nil {
		t.Fatalf("Erro na geração do relatório: %v", err)
	}

	// Verificar estrutura do relatório
	if report.GeneratedAt.IsZero() {
		t.Error("Relatório deveria ter data de geração")
	}

	if report.TotalRecords != 1 {
		t.Errorf("Esperava 1 registro, mas relatório mostra %d", report.TotalRecords)
	}

	if report.Summary.OverallScore < 0 || report.Summary.OverallScore > 100 {
		t.Errorf("Score geral deveria estar entre 0-100, mas é %.2f", report.Summary.OverallScore)
	}

	// Verificar níveis de risco válidos
	validRiskLevels := map[string]bool{
		"low": true, "medium": true, "high": true, "critical": true,
	}

	if !validRiskLevels[report.Summary.RiskLevel] {
		t.Errorf("Nível de risco inválido: %s", report.Summary.RiskLevel)
	}

	// Verificar se recomendações foram geradas
	if len(report.Summary.RecommendedActions) == 0 {
		t.Error("Relatório deveria ter pelo menos uma recomendação")
	}

	t.Logf("Relatório gerado: %d anomalias, score %.1f, risco %s",
		report.TotalAnomalies, report.Summary.OverallScore, report.Summary.RiskLevel)
}

// TestWorkflowIntegration testa integração com workflows
func TestWorkflowIntegration(t *testing.T) {
	analyzer := NewAnalyzer(nil)
	// TODO: Implementar NewValidatedVRWorkflow quando disponível
	// workflow := NewValidatedVRWorkflow(analyzer)

	if analyzer == nil {
		t.Fatal("Analyzer não deveria ser nil")
	}

	// TODO: Continuar teste quando workflow estiver implementado
	// if workflow.Name() != "validated-vr-processing" {
	//     t.Errorf("Nome do workflow incorreto: %s", workflow.Name())
	// }

	// TODO: Testar steps quando workflow estiver implementado
	// steps := workflow.Steps()
	// if len(steps) == 0 {
	//     t.Error("Workflow deveria ter steps")
	// }

	// TODO: Verificar steps quando workflow estiver implementado
	// hasAnomalyDetection := false
	// for _, step := range steps {
	//     if step.Name() == "anomaly-detection" {
	//         hasAnomalyDetection = true
	//         break
	//     }
	// }
	//
	// if !hasAnomalyDetection {
	//     t.Error("Workflow deveria ter step de detecção de anomalias")
	// }

	// TODO: Log steps quando workflow estiver implementado
	// t.Logf("Workflow criado com %d steps", len(steps))
	t.Log("Teste de integração concluído - pendente implementação do workflow")
}

// TestFormatting testa formatação de relatórios
func TestFormatting(t *testing.T) {
	// Criar relatório de teste
	report := &AnomalyReport{
		GeneratedAt:    time.Now(),
		TotalRecords:   100,
		TotalAnomalies: 5,
		AnomaliesByType: map[AnomalyType]int{
			AnomalyTypeValue:        2,
			AnomalyTypeTemporal:     2,
			AnomalyTypeRelationship: 1,
		},
		AnomaliesBySeverity: map[string]int{
			"critical": 1,
			"high":     2,
			"medium":   2,
		},
		Summary: AnomalySummary{
			OverallScore:   75.5,
			RiskLevel:      "medium",
			CriticalIssues: 1,
			RecommendedActions: []string{
				"Resolver problema crítico",
				"Revisar dados de entrada",
			},
		},
	}

	formatted := FormatAnomalyReportForHuman(report)

	if len(formatted) == 0 {
		t.Error("Relatório formatado não deveria estar vazio")
	}

	// Verificar se contém informações essenciais
	requiredContent := []string{
		"RELATÓRIO DE ANOMALIAS",
		"100",      // Total de registros
		"5",        // Total de anomalias
		"75.5",     // Score
		"medium",   // Risk level
		"critical", // Severidade
		"RECOMENDAÇÕES",
	}

	for _, content := range requiredContent {
		if !strings.Contains(formatted, content) {
			t.Errorf("Relatório formatado deveria conter '%s'", content)
		}
	}

	t.Logf("Relatório formatado gerado com %d caracteres", len(formatted))
}

// Benchmark para performance
func BenchmarkAnomalyDetection(b *testing.B) {
	detector := NewAnomalyDetector(nil)

	// Criar dataset de teste maior
	colaboradores := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		matricula := fmt.Sprintf("%04d", i)
		colaboradores[matricula] = map[string]interface{}{
			"matricula":  matricula,
			"sindicato":  "SINDICATO_A",
			"vr_total":   500.0 + float64(i%100), // Variar valores
			"dias_uteis": 22,
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx := NewAnalysisContext(colaboradores, nil)
		_, err := detector.DetectAnomalies(ctx)
		if err != nil {
			b.Fatalf("Erro no benchmark: %v", err)
		}
	}
}
