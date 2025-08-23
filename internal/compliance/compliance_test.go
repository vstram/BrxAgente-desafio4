package compliance

import (
	"testing"
	"time"
)

func TestComplianceChecker_CheckCompliance(t *testing.T) {
	// Criar um ComplianceChecker para teste
	checker := NewComplianceChecker("../data/regulations")
	
	// Dados de teste válidos
	validData := map[string]interface{}{
		"matricula":                        "12345",
		"nome":                            "João da Silva",
		"valor_vr":                        45.0,
		"percentual_desconto":             15.0,
		"carga_horaria":                   8.0,
		"tipo_colaborador":                "efetivo",
		"tem_termo_adesao":                true,
		"tem_comprovantes_entrega":        true,
		"uso_exclusivo_alimentacao":       true,
		"controle_utilizacao":             true,
		"inscricao_pat":                   true,
		"renovacao_pat_atualizada":        true,
		"controle_beneficiarios_atualizado": true,
		"prestacao_contas_pat_em_dia":     true,
		"estabelecimentos_credenciados":   true,
		"qualidade_nutricional_ok":        true,
		"criterios_objetivos":             true,
		"aplicacao_uniforme":              true,
	}

	result := checker.CheckCompliance(validData)

	if result.Status != StatusCompliant {
		t.Errorf("Esperado status %s, obtido %s", StatusCompliant, result.Status)
	}

	if len(result.Violations) != 0 {
		t.Errorf("Esperado 0 violações, obtido %d", len(result.Violations))
	}

	if result.ComplianceScore != 100.0 {
		t.Errorf("Esperado score 100.0, obtido %.2f", result.ComplianceScore)
	}
}

func TestComplianceChecker_CheckCompliance_Violations(t *testing.T) {
	checker := NewComplianceChecker("../data/regulations")
	
	// Dados de teste com violações
	violationData := map[string]interface{}{
		"matricula":                        "12346",
		"nome":                            "Maria da Silva",
		"valor_vr":                        60.0, // Acima do limite
		"percentual_desconto":             25.0, // Acima do limite
		"carga_horaria":                   4.0,  // Abaixo do mínimo
		"tipo_colaborador":                "diretor", // Não elegível
		"tem_termo_adesao":                false,
		"tem_comprovantes_entrega":        false,
		"uso_exclusivo_alimentacao":       false,
		"controle_utilizacao":             false,
		"inscricao_pat":                   false,
		"renovacao_pat_atualizada":        false,
		"controle_beneficiarios_atualizado": false,
		"prestacao_contas_pat_em_dia":     false,
		"estabelecimentos_credenciados":   false,
		"qualidade_nutricional_ok":        false,
		"criterios_objetivos":             false,
		"aplicacao_uniforme":              false,
	}

	result := checker.CheckCompliance(violationData)

	if result.Status == StatusCompliant {
		t.Errorf("Esperado status não conforme, obtido %s", result.Status)
	}

	if len(result.Violations) == 0 {
		t.Errorf("Esperado pelo menos 1 violação, obtido %d", len(result.Violations))
	}

	if result.ComplianceScore >= 100.0 {
		t.Errorf("Esperado score menor que 100, obtido %.2f", result.ComplianceScore)
	}
}

func TestAuditTrail_RecordAction(t *testing.T) {
	auditTrail := NewAuditTrail(100)
	
	// Gravar uma ação
	metadata := map[string]interface{}{
		"test_key": "test_value",
	}
	
	auditTrail.RecordSimpleAction("testuser", "test_action", "Test description", "INFO", metadata)
	
	// Verificar que a entrada foi gravada
	entries := auditTrail.GetEntriesByDateRange(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour))
	
	if len(entries) != 1 {
		t.Errorf("Esperado 1 entrada, obtido %d", len(entries))
	}
	
	entry := entries[0]
	if entry.User != "testuser" {
		t.Errorf("Esperado usuário 'testuser', obtido '%s'", entry.User)
	}
	
	if entry.Action != "Test description" {
		t.Errorf("Esperado ação 'Test description', obtido '%s'", entry.Action)
	}
	
	if entry.Level != "INFO" {
		t.Errorf("Esperado nível 'INFO', obtido '%s'", entry.Level)
	}
}

func TestRiskAssessor_AssessEmployeeRisk(t *testing.T) {
	auditTrail := NewAuditTrail(100)
	checker := NewComplianceChecker("../data/regulations")
	riskAssessor := NewRiskAssessor(checker, auditTrail)
	
	// Dados de funcionário com baixo risco
	lowRiskData := map[string]interface{}{
		"matricula":    "12345",
		"nome":        "João da Silva",
		"valor_vr":    45.0,
		"data_admissao": time.Now().AddDate(-1, 0, 0),
	}
	
	assessment, err := riskAssessor.AssessEmployeeRisk(lowRiskData)
	
	if err != nil {
		t.Errorf("Erro inesperado: %v", err)
	}
	
	if assessment.EntityID != "12345" {
		t.Errorf("Esperado EntityID '12345', obtido '%s'", assessment.EntityID)
	}
	
	if assessment.OverallRisk == RiskCritical {
		t.Errorf("Não esperado risco crítico para dados válidos")
	}
}

func TestRiskAssessor_AssessEmployeeRisk_HighRisk(t *testing.T) {
	auditTrail := NewAuditTrail(100)
	checker := NewComplianceChecker("../data/regulations")
	riskAssessor := NewRiskAssessor(checker, auditTrail)
	
	// Dados de funcionário com alto risco
	highRiskData := map[string]interface{}{
		"matricula":    "12346",
		"valor_vr":    60.0, // Acima do limite
		// Dados obrigatórios faltando
	}
	
	assessment, err := riskAssessor.AssessEmployeeRisk(highRiskData)
	
	if err != nil {
		t.Errorf("Erro inesperado: %v", err)
	}
	
	if assessment.EntityID != "12346" {
		t.Errorf("Esperado EntityID '12346', obtido '%s'", assessment.EntityID)
	}
	
	if assessment.RiskScore == 0 {
		t.Errorf("Esperado score de risco > 0 para dados com problemas")
	}
	
	if len(assessment.RiskFactors) == 0 {
		t.Errorf("Esperado pelo menos 1 fator de risco")
	}
}

func TestDocumentationGenerator_GenerateComplianceReport(t *testing.T) {
	auditTrail := NewAuditTrail(100)
	checker := NewComplianceChecker("../data/regulations")
	riskAssessor := NewRiskAssessor(checker, auditTrail)
	docGen := NewDocumentationGenerator(checker, riskAssessor, auditTrail)
	
	// Dados de exemplo
	employeeData := []map[string]interface{}{
		{
			"matricula":    "12345",
			"nome":        "João da Silva",
			"valor_vr":    45.0,
			"data_admissao": time.Now().AddDate(-1, 0, 0),
		},
		{
			"matricula":    "12346",
			"nome":        "Maria dos Santos",
			"valor_vr":    60.0, // Violação
			"data_admissao": time.Now().AddDate(-2, 0, 0),
		},
	}
	
	period := DateRange{
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
	}
	
	report, err := docGen.GenerateComplianceReport(employeeData, period)
	
	if err != nil {
		t.Errorf("Erro inesperado: %v", err)
	}
	
	if report.Summary.TotalEmployees != 2 {
		t.Errorf("Esperado 2 funcionários, obtido %d", report.Summary.TotalEmployees)
	}
	
	if len(report.Violations) == 0 {
		t.Errorf("Esperado pelo menos 1 violação nos dados de teste")
	}
	
	if len(report.Recommendations) == 0 {
		t.Errorf("Esperado pelo menos 1 recomendação")
	}
}

func TestDocumentationGenerator_GenerateMarkdownReport(t *testing.T) {
	auditTrail := NewAuditTrail(100)
	checker := NewComplianceChecker("../data/regulations")
	riskAssessor := NewRiskAssessor(checker, auditTrail)
	docGen := NewDocumentationGenerator(checker, riskAssessor, auditTrail)
	
	// Criar um relatório básico
	report := &ComplianceReport{
		GeneratedAt: time.Now(),
		Period: DateRange{
			StartDate: time.Now().AddDate(0, -1, 0),
			EndDate:   time.Now(),
		},
		Summary: ComplianceSummary{
			TotalEmployees: 10,
			ComplianceRate: 85.5,
		},
		Violations: []ComplianceViolation{
			{
				RuleID:      "test_rule",
				EntityID:    "12345",
				Description: "Violação de teste",
				Severity:    SeverityMedium,
			},
		},
		Recommendations: []string{
			"Recomendação de teste",
		},
		Metadata: ReportMetadata{
			Version:     "1.0",
			GeneratedBy: "Test Suite",
			QualityMetrics: map[string]float64{
				"data_completeness": 95.0,
			},
		},
	}
	
	markdown := docGen.GenerateMarkdownReport(report)
	
	if len(markdown) == 0 {
		t.Errorf("Relatório em markdown não deve estar vazio")
	}
	
	// Verificar se contém elementos esperados
	if !contains(markdown, "Relatório de Compliance") {
		t.Errorf("Markdown deve conter título do relatório")
	}
	
	if !contains(markdown, "85.5%") {
		t.Errorf("Markdown deve conter taxa de compliance")
	}
	
	if !contains(markdown, "Recomendação de teste") {
		t.Errorf("Markdown deve conter recomendações")
	}
}

func TestRiskMatrix_GetRiskMatrix(t *testing.T) {
	auditTrail := NewAuditTrail(100)
	checker := NewComplianceChecker("../data/regulations")
	riskAssessor := NewRiskAssessor(checker, auditTrail)
	
	matrix := riskAssessor.GetRiskMatrix()
	
	// Verificar se contém os valores esperados
	if matrix.Probability["high"] != 0.7 {
		t.Errorf("Esperado probabilidade 'high' = 0.7, obtido %.1f", matrix.Probability["high"])
	}
	
	if matrix.Impact["critical"] != 0.9 {
		t.Errorf("Esperado impacto 'critical' = 0.9, obtido %.1f", matrix.Impact["critical"])
	}
	
	if matrix.Thresholds["medium"] != 0.6 {
		t.Errorf("Esperado threshold 'medium' = 0.6, obtido %.1f", matrix.Thresholds["medium"])
	}
}

func TestAuditTrail_Cleanup(t *testing.T) {
	// Criar um audit trail com capacidade baixa para testar limpeza
	auditTrail := NewAuditTrail(5)
	
	// Adicionar mais entradas que a capacidade
	for i := 0; i < 10; i++ {
		auditTrail.RecordSimpleAction("testuser", "test_action", "Test description", "INFO", nil)
		time.Sleep(1 * time.Millisecond) // Garantir timestamps diferentes
	}
	
	// Verificar que só mantém o número máximo
	entries := auditTrail.GetEntriesByDateRange(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour))
	
	if len(entries) > 5 {
		t.Errorf("Esperado máximo 5 entradas após limpeza, obtido %d", len(entries))
	}
}

// Função auxiliar para verificar se uma string contém uma substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || contains(s[1:], substr) || (len(s) > 0 && s[:len(substr)] == substr))
}

// Benchmark para verificar performance do compliance checker
func BenchmarkComplianceChecker_CheckCompliance(b *testing.B) {
	checker := NewComplianceChecker("../data/regulations")
	
	data := map[string]interface{}{
		"matricula":                        "12345",
		"nome":                            "João da Silva",
		"valor_vr":                        45.0,
		"percentual_desconto":             15.0,
		"carga_horaria":                   8.0,
		"tipo_colaborador":                "efetivo",
		"tem_termo_adesao":                true,
		"tem_comprovantes_entrega":        true,
		"uso_exclusivo_alimentacao":       true,
		"controle_utilizacao":             true,
		"inscricao_pat":                   true,
		"renovacao_pat_atualizada":        true,
		"controle_beneficiarios_atualizado": true,
		"prestacao_contas_pat_em_dia":     true,
		"estabelecimentos_credenciados":   true,
		"qualidade_nutricional_ok":        true,
		"criterios_objetivos":             true,
		"aplicacao_uniforme":              true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.CheckCompliance(data)
	}
}

// Benchmark para verificar performance do risk assessor
func BenchmarkRiskAssessor_AssessEmployeeRisk(b *testing.B) {
	auditTrail := NewAuditTrail(1000)
	checker := NewComplianceChecker("../data/regulations")
	riskAssessor := NewRiskAssessor(checker, auditTrail)
	
	data := map[string]interface{}{
		"matricula":    "12345",
		"nome":        "João da Silva",
		"valor_vr":    45.0,
		"data_admissao": time.Now().AddDate(-1, 0, 0),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		riskAssessor.AssessEmployeeRisk(data)
	}
}