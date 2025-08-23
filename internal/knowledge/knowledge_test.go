package knowledge

import (
	"testing"
	"time"
)

// TestKnowledgeBaseCreation testa a criação de uma base de conhecimento
func TestKnowledgeBaseCreation(t *testing.T) {
	kb := NewKnowledgeBaseManager()
	
	if kb == nil {
		t.Fatal("KnowledgeBaseManager não foi criado corretamente")
	}
	
	if kb.loaded {
		t.Error("KnowledgeBaseManager deveria começar como não carregado")
	}
	
	stats := kb.GetStats()
	if loaded, exists := stats["loaded"]; !exists || loaded.(bool) {
		t.Error("Stats deveria mostrar que a base não foi carregada")
	}
}

// TestKnowledgeItemCreation testa a criação de itens de conhecimento
func TestKnowledgeItemCreation(t *testing.T) {
	item := KnowledgeItem{
		ID:           "test_001",
		Title:        "Teste de Conhecimento",
		Content:      "Este é um item de teste para validação.",
		Source:       "Teste Unitário",
		EffectiveDate: time.Now(),
		Categories:   []string{"teste", "validacao"},
		Metadata:     map[string]string{"tipo": "teste"},
		Type:         "test",
	}
	
	if item.ID != "test_001" {
		t.Errorf("Esperado ID 'test_001', obtido '%s'", item.ID)
	}
	
	if len(item.Categories) != 2 {
		t.Errorf("Esperado 2 categorias, obtido %d", len(item.Categories))
	}
}

// TestPolicyEngineCreation testa a criação do motor de políticas
func TestPolicyEngineCreation(t *testing.T) {
	kb := NewKnowledgeBaseManager()
	pe := NewPolicyEngine(kb)
	
	if pe == nil {
		t.Fatal("PolicyEngine não foi criado corretamente")
	}
	
	if pe.knowledgeBase == nil {
		t.Error("PolicyEngine deveria ter uma referência para KnowledgeBase")
	}
}

// TestCitationManagerCreation testa a criação do gerenciador de citações
func TestCitationManagerCreation(t *testing.T) {
	cm := NewCitationManager()
	
	if cm == nil {
		t.Fatal("CitationManager não foi criado corretamente")
	}
	
	if cm.citations == nil {
		t.Error("CitationManager deveria ter um mapa de citações inicializado")
	}
}

// TestReasoningEngineCreation testa a criação do motor de raciocínio
func TestReasoningEngineCreation(t *testing.T) {
	kb := NewKnowledgeBaseManager()
	pe := NewPolicyEngine(kb)
	cm := NewCitationManager()
	re := NewReasoningEngine(kb, pe, cm)
	
	if re == nil {
		t.Fatal("ReasoningEngine não foi criado corretamente")
	}
	
	if re.knowledgeBase == nil {
		t.Error("ReasoningEngine deveria ter uma referência para KnowledgeBase")
	}
	
	if re.policyEngine == nil {
		t.Error("ReasoningEngine deveria ter uma referência para PolicyEngine")
	}
	
	if re.citationManager == nil {
		t.Error("ReasoningEngine deveria ter uma referência para CitationManager")
	}
}

// TestCitationCreation testa a criação de citações
func TestCitationCreation(t *testing.T) {
	cm := NewCitationManager()
	
	item := KnowledgeItem{
		ID:           "cite_001",
		Title:        "Teste de Citação",
		Content:      "Conteúdo para teste de citação.",
		Source:       "Manual de Teste v1.0",
		EffectiveDate: time.Now(),
		Type:         "policy",
	}
	
	citation := cm.CreateCitation(item)
	
	if citation.Source != item.Source {
		t.Errorf("Esperado fonte '%s', obtido '%s'", item.Source, citation.Source)
	}
	
	if citation.Reliability == "" {
		t.Error("Citation deveria ter um nível de confiabilidade definido")
	}
}

// TestFormatCitation testa a formatação de citações
func TestFormatCitation(t *testing.T) {
	cm := NewCitationManager()
	
	citation := Citation{
		Source:      "Manual de Teste v1.0",
		Date:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Reliability: "high",
		Section:     "Seção 1",
	}
	
	formatted := cm.FormatCitation(citation)
	
	if formatted == "" {
		t.Error("Citação formatada não deveria estar vazia")
	}
	
	if !contains(formatted, citation.Source) {
		t.Errorf("Citação formatada deveria conter a fonte '%s'", citation.Source)
	}
}

// TestSearchEmptyBase testa busca em base vazia
func TestSearchEmptyBase(t *testing.T) {
	kb := NewKnowledgeBaseManager()
	
	results, err := kb.Search("teste", 10)
	if err == nil {
		t.Error("Busca em base não carregada deveria retornar erro")
	}
	
	if len(results) != 0 {
		t.Errorf("Esperado 0 resultados, obtido %d", len(results))
	}
}

// TestScenarioTypeIdentification testa identificação de tipos de cenário
func TestScenarioTypeIdentification(t *testing.T) {
	kb := NewKnowledgeBaseManager()
	pe := NewPolicyEngine(kb)
	cm := NewCitationManager()
	re := NewReasoningEngine(kb, pe, cm)
	
	// Cenário simples
	simpleScenario := map[string]interface{}{
		"query": "Qual é o valor base do VR?",
	}
	scenarioType := re.identifyScenarioType(simpleScenario)
	if scenarioType != ScenarioSimpleQuery {
		t.Errorf("Esperado ScenarioSimpleQuery, obtido %v", scenarioType)
	}
	
	// Cenário complexo
	complexScenario := map[string]interface{}{
		"query":          "Calcular VR com múltiplos fatores",
		"data_admissao":  "2024-08-20",
		"data_desligamento": "2024-09-15",
		"afastamento":    map[string]interface{}{"dias": 5},
	}
	scenarioType = re.identifyScenarioType(complexScenario)
	if scenarioType != ScenarioComplexCalculation {
		t.Errorf("Esperado ScenarioComplexCalculation, obtido %v", scenarioType)
	}
}

// Função auxiliar para verificar se uma string contém outra
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		 s[len(s)-len(substr):] == substr || 
		 containsInMiddle(s, substr))))
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}