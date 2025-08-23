package intelligence

import (
	"testing"
	"time"
)

func TestInsightGenerator_GenerateInsights(t *testing.T) {
	generator := NewInsightGenerator(nil)
	
	// Dados de teste
	processingData := &ProcessingData{
		TotalCollaborators: 100,
		TotalVRValue:      50000.0,
		ProcessingTime:    45 * time.Minute,
		ErrorCount:        2,
		WarningCount:      1,
		ProcessedAt:       time.Now(),
		SindicatoDistribution: map[string]int{
			"Sindicato A": 60,
			"Sindicato B": 40,
		},
		Metadata: make(map[string]interface{}),
	}
	
	// Gerar insights
	insights, err := generator.GenerateInsights(processingData)
	if err != nil {
		t.Fatalf("Erro ao gerar insights: %v", err)
	}
	
	// Verificações
	if len(insights) == 0 {
		t.Error("Nenhum insight foi gerado")
	}
	
	// Verificar se há insights de diferentes tipos
	types := make(map[InsightType]int)
	for _, insight := range insights {
		types[insight.Type]++
		
		// Verificar campos obrigatórios
		if insight.Title == "" {
			t.Error("Insight sem título")
		}
		if insight.Description == "" {
			t.Error("Insight sem descrição")
		}
		if insight.Confidence < 0 || insight.Confidence > 1 {
			t.Errorf("Confiança inválida: %f", insight.Confidence)
		}
	}
	
	// Deve ter pelo menos insights financeiros e operacionais
	if types[InsightTypeFinancial] == 0 {
		t.Error("Nenhum insight financeiro gerado")
	}
	if types[InsightTypeOperational] == 0 {
		t.Error("Nenhum insight operacional gerado")
	}
}

func TestInsightGenerator_WithAnomalyReport(t *testing.T) {
	config := DefaultInsightConfig()
	config.EnableAnomalyInsights = true
	generator := NewInsightGenerator(config)
	
	// Criar dados com relatório de anomalias
	anomalyReport := &AnomalyReport{
		TotalRecords:   100,
		TotalAnomalies: 5,
		GeneratedAt:    time.Now(),
		Summary: AnomalySummary{
			OverallScore:   85.0,
			RiskLevel:      "medium",
			CriticalIssues: 1,
		},
		AnomaliesByType: map[AnomalyType]int{
			"value":    3,
			"temporal": 2,
		},
	}
	
	processingData := &ProcessingData{
		TotalCollaborators: 100,
		TotalVRValue:      50000.0,
		ProcessingTime:    45 * time.Minute,
		ErrorCount:        0,
		WarningCount:      0,
		AnomalyReport:     anomalyReport,
		ProcessedAt:       time.Now(),
		Metadata:          make(map[string]interface{}),
	}
	
	insights, err := generator.GenerateInsights(processingData)
	if err != nil {
		t.Fatalf("Erro ao gerar insights: %v", err)
	}
	
	// Verificar se há insights de anomalias
	hasAnomalyInsight := false
	for _, insight := range insights {
		if insight.Type == InsightTypeAnomaly {
			hasAnomalyInsight = true
			break
		}
	}
	
	if !hasAnomalyInsight {
		t.Error("Nenhum insight de anomalia gerado")
	}
}

func TestFormatInsightsForHuman(t *testing.T) {
	insights := []*Insight{
		{
			Type:        InsightTypeFinancial,
			Title:       "Teste Financeiro",
			Description: "Descrição teste",
			Impact:      "Impacto teste",
			Action:      "Ação teste",
			Priority:    PriorityHigh,
			Confidence:  0.90,
			GeneratedAt: time.Now(),
		},
		{
			Type:        InsightTypeOperational,
			Title:       "Teste Operacional",
			Description: "Descrição teste 2",
			Priority:    PriorityMedium,
			Confidence:  0.75,
			GeneratedAt: time.Now(),
		},
	}
	
	output := FormatInsightsForHuman(insights)
	
	if output == "" {
		t.Error("Output vazio")
	}
	
	// Verificar se contém elementos esperados
	if !contains(output, "INSIGHTS AUTOMÁTICOS") {
		t.Error("Título principal não encontrado")
	}
	
	if !contains(output, "Teste Financeiro") {
		t.Error("Título do insight não encontrado")
	}
}

func TestInsightPriority_String(t *testing.T) {
	tests := []struct {
		priority InsightPriority
		expected string
	}{
		{PriorityLow, "baixa"},
		{PriorityMedium, "média"},
		{PriorityHigh, "alta"},
		{PriorityCritical, "crítica"},
	}
	
	for _, test := range tests {
		if test.priority.String() != test.expected {
			t.Errorf("Prioridade %d deveria retornar '%s', mas retornou '%s'", 
				test.priority, test.expected, test.priority.String())
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && 
		   len(s) >= len(substr) && 
		   findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}