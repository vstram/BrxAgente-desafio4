package training

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// KnowledgeItem representa um item da base de conhecimento
type KnowledgeItem struct {
	ID          string            `json:"id"`
	Category    string            `json:"category"`
	Question    string            `json:"question"`
	Answer      string            `json:"answer"`
	Context     []string          `json:"context"`
	Examples    []string          `json:"examples"`
	Confidence  float64           `json:"confidence"`
	Metadata    map[string]string `json:"metadata"`
}

// KnowledgeBase representa toda a base de conhecimento estruturada
type KnowledgeBase struct {
	Policies      []PolicyItem      `json:"policies"`
	Calculations  []CalculationItem `json:"calculations"`
	Syndicates    []SyndicateItem   `json:"syndicates"`
	Holidays      []HolidayItem     `json:"holidays"`
	FAQ           []FAQItem         `json:"faq"`
	Examples      []ExampleItem     `json:"examples"`
	LastUpdated   time.Time         `json:"last_updated"`
}

// PolicyItem representa uma política específica
type PolicyItem struct {
	ID          string            `json:"id"`
	Category    string            `json:"category"`
	Question    string            `json:"question"`
	Answer      string            `json:"answer"`
	Context     []string          `json:"context"`
	Examples    []string          `json:"examples"`
	Confidence  float64           `json:"confidence"`
	Metadata    map[string]string `json:"metadata"`
}

// CalculationItem representa uma regra de cálculo
type CalculationItem struct {
	ID       string    `json:"id"`
	Scenario string    `json:"scenario"`
	Rule     string    `json:"rule"`
	Formula  string    `json:"formula"`
	Context  []string  `json:"context"`
	Examples []Example `json:"examples"`
	Confidence float64 `json:"confidence"`
}

// Example representa um exemplo de cálculo
type Example struct {
	Case        string  `json:"case"`
	Calculation string  `json:"calculation"`
	Reasoning   string  `json:"reasoning"`
}

// SyndicateItem representa informações de um sindicato
type SyndicateItem struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	VRValue      float64  `json:"vr_value"`
	Description  string   `json:"description"`
	SpecialRules []string `json:"special_rules"`
}

// HolidayItem representa informações sobre feriados
type HolidayItem struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Examples    []string `json:"examples,omitempty"`
	Note        string   `json:"note,omitempty"`
}

// FAQItem representa uma pergunta frequente
type FAQItem struct {
	ID               string   `json:"id"`
	Question         string   `json:"question"`
	Answer           string   `json:"answer"`
	Category         string   `json:"category"`
	Frequency        string   `json:"frequency"`
	RelatedPolicies  []string `json:"related_policies"`
	Keywords         []string `json:"keywords"`
	RequiresManualReview bool `json:"requires_manual_review,omitempty"`
}

// ExampleItem representa um exemplo prático completo
type ExampleItem struct {
	ID                   string               `json:"id"`
	Title                string               `json:"title"`
	Scenario             string               `json:"scenario"`
	InputData            map[string]interface{} `json:"input_data"`
	StepByStepCalculation map[string]CalculationStep `json:"step_by_step_calculation"`
	ExpectedResult       ExpectedResult       `json:"expected_result"`
	KeyLearning          string               `json:"key_learning"`
}

// CalculationStep representa um passo do cálculo
type CalculationStep struct {
	Description string `json:"description"`
	Rule        string `json:"rule,omitempty"`
	Calculation string `json:"calculation,omitempty"`
	Result      string `json:"result"`
}

// ExpectedResult representa o resultado esperado
type ExpectedResult struct {
	ValorVR     float64 `json:"valor_vr"`
	Observacoes string  `json:"observacoes"`
	Confidence  float64 `json:"confidence"`
}

// KnowledgeManager gerencia a base de conhecimento
type KnowledgeManager struct {
	basePath     string
	knowledgeBase *KnowledgeBase
	loadedAt     time.Time
}

// NewKnowledgeManager cria um novo gerenciador de conhecimento
func NewKnowledgeManager(basePath string) *KnowledgeManager {
	return &KnowledgeManager{
		basePath: basePath,
	}
}

// LoadKnowledgeBase carrega toda a base de conhecimento
func (km *KnowledgeManager) LoadKnowledgeBase() error {
	kb := &KnowledgeBase{}
	
	// Carregar políticas
	policies, err := km.loadVRPolicies()
	if err != nil {
		return fmt.Errorf("erro ao carregar políticas: %v", err)
	}
	kb.Policies = policies

	// Carregar regulamentações
	regulations, err := km.loadRegulations()
	if err != nil {
		return fmt.Errorf("erro ao carregar regulamentações: %v", err)
	}
	kb.Policies = append(kb.Policies, regulations...)

	// Carregar FAQ
	faq, err := km.loadFAQ()
	if err != nil {
		return fmt.Errorf("erro ao carregar FAQ: %v", err)
	}
	kb.FAQ = faq

	// Carregar exemplos
	examples, err := km.loadExamples()
	if err != nil {
		return fmt.Errorf("erro ao carregar exemplos: %v", err)
	}
	kb.Examples = examples

	kb.LastUpdated = time.Now()
	km.knowledgeBase = kb
	km.loadedAt = time.Now()

	return nil
}

// loadVRPolicies carrega as políticas de VR
func (km *KnowledgeManager) loadVRPolicies() ([]PolicyItem, error) {
	filePath := filepath.Join(km.basePath, "knowledge_base", "vr_policies.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var vrData struct {
		VRKnowledge struct {
			Policies []PolicyItem `json:"policies"`
		} `json:"vr_knowledge"`
	}

	if err := json.Unmarshal(data, &vrData); err != nil {
		return nil, err
	}

	return vrData.VRKnowledge.Policies, nil
}

// loadRegulations carrega as regulamentações
func (km *KnowledgeManager) loadRegulations() ([]PolicyItem, error) {
	filePath := filepath.Join(km.basePath, "knowledge_base", "regulations.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Implementação básica - converter regulamentações para PolicyItem
	var regulationsData struct {
		Regulations struct {
			CLTArticles       []interface{} `json:"clt_articles"`
			LaborLaws        []interface{} `json:"labor_laws"`
			CollectiveAgreements []interface{} `json:"collective_agreements"`
			InternalPolicies []interface{} `json:"internal_policies"`
		} `json:"regulations"`
	}

	if err := json.Unmarshal(data, &regulationsData); err != nil {
		return nil, err
	}

	// Por enquanto, retornar lista vazia - na implementação completa, 
	// converteria os dados para PolicyItem
	var policies []PolicyItem
	
	return policies, nil
}

// loadFAQ carrega as perguntas frequentes
func (km *KnowledgeManager) loadFAQ() ([]FAQItem, error) {
	filePath := filepath.Join(km.basePath, "knowledge_base", "faq.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var faqData struct {
		FAQ struct {
			CommonQuestions []FAQItem `json:"common_questions"`
		} `json:"faq"`
	}

	if err := json.Unmarshal(data, &faqData); err != nil {
		return nil, err
	}

	return faqData.FAQ.CommonQuestions, nil
}

// loadExamples carrega os exemplos práticos
func (km *KnowledgeManager) loadExamples() ([]ExampleItem, error) {
	filePath := filepath.Join(km.basePath, "knowledge_base", "examples.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var exampleData struct {
		PracticalExamples []ExampleItem `json:"practical_examples"`
	}

	if err := json.Unmarshal(data, &exampleData); err != nil {
		return nil, err
	}

	return exampleData.PracticalExamples, nil
}

// FindRelevantKnowledge encontra conhecimento relevante para uma pergunta
func (km *KnowledgeManager) FindRelevantKnowledge(question string) ([]KnowledgeItem, error) {
	if km.knowledgeBase == nil {
		return nil, fmt.Errorf("base de conhecimento não carregada")
	}

	var relevantItems []KnowledgeItem
	question = strings.ToLower(question)

	// Buscar em políticas
	for _, policy := range km.knowledgeBase.Policies {
		if km.isRelevant(question, policy.Question, policy.Answer, policy.Context) {
			item := KnowledgeItem{
				ID:         policy.ID,
				Category:   policy.Category,
				Question:   policy.Question,
				Answer:     policy.Answer,
				Context:    policy.Context,
				Examples:   policy.Examples,
				Confidence: policy.Confidence,
				Metadata:   policy.Metadata,
			}
			relevantItems = append(relevantItems, item)
		}
	}

	// Buscar em FAQ
	for _, faq := range km.knowledgeBase.FAQ {
		if km.isRelevant(question, faq.Question, faq.Answer, faq.Keywords) {
			item := KnowledgeItem{
				ID:         faq.ID,
				Category:   faq.Category,
				Question:   faq.Question,
				Answer:     faq.Answer,
				Context:    faq.RelatedPolicies,
				Confidence: 0.85, // FAQ tem confiança padrão menor
				Metadata: map[string]string{
					"frequency": faq.Frequency,
					"source":    "faq",
				},
			}
			relevantItems = append(relevantItems, item)
		}
	}

	return relevantItems, nil
}

// isRelevant verifica se um item é relevante para a pergunta
func (km *KnowledgeManager) isRelevant(question string, itemQuestion, itemAnswer string, keywords []string) bool {
	// Implementação básica de relevância usando palavras-chave
	questionWords := strings.Fields(strings.ToLower(question))
	
	// Verificar correspondência na pergunta do item
	itemQuestionWords := strings.Fields(strings.ToLower(itemQuestion))
	for _, qWord := range questionWords {
		for _, iWord := range itemQuestionWords {
			if strings.Contains(iWord, qWord) || strings.Contains(qWord, iWord) {
				return true
			}
		}
	}

	// Verificar correspondência nas palavras-chave
	for _, keyword := range keywords {
		for _, qWord := range questionWords {
			if strings.Contains(strings.ToLower(keyword), qWord) {
				return true
			}
		}
	}

	// Verificar termos específicos de VR
	vrTerms := []string{
		"vale", "refeição", "vr", "colaborador", "sindicato",
		"admissão", "desligamento", "férias", "licença", "afastamento",
		"estagiário", "diretor", "aprendiz", "cálculo", "proporcional",
	}

	for _, term := range vrTerms {
		if strings.Contains(question, term) {
			// Verificar se o item também contém termos relacionados
			fullText := strings.ToLower(itemQuestion + " " + itemAnswer + " " + strings.Join(keywords, " "))
			if strings.Contains(fullText, term) {
				return true
			}
		}
	}

	return false
}

// GetKnowledgeStatistics retorna estatísticas da base de conhecimento
func (km *KnowledgeManager) GetKnowledgeStatistics() (map[string]interface{}, error) {
	if km.knowledgeBase == nil {
		return nil, fmt.Errorf("base de conhecimento não carregada")
	}

	stats := make(map[string]interface{})
	
	stats["total_policies"] = len(km.knowledgeBase.Policies)
	stats["total_faq"] = len(km.knowledgeBase.FAQ)
	stats["total_examples"] = len(km.knowledgeBase.Examples)
	stats["last_updated"] = km.knowledgeBase.LastUpdated
	stats["loaded_at"] = km.loadedAt

	// Estatísticas por categoria
	categoryCount := make(map[string]int)
	for _, policy := range km.knowledgeBase.Policies {
		categoryCount[policy.Category]++
	}
	stats["policies_by_category"] = categoryCount

	// Estatísticas de confiança
	var totalConfidence float64
	confidenceCount := 0
	for _, policy := range km.knowledgeBase.Policies {
		if policy.Confidence > 0 {
			totalConfidence += policy.Confidence
			confidenceCount++
		}
	}
	if confidenceCount > 0 {
		stats["average_confidence"] = totalConfidence / float64(confidenceCount)
	}

	return stats, nil
}

// RefreshKnowledge recarrega a base de conhecimento
func (km *KnowledgeManager) RefreshKnowledge() error {
	return km.LoadKnowledgeBase()
}

// IsKnowledgeStale verifica se a base de conhecimento precisa ser atualizada
func (km *KnowledgeManager) IsKnowledgeStale(maxAge time.Duration) bool {
	if km.knowledgeBase == nil {
		return true
	}
	return time.Since(km.loadedAt) > maxAge
}