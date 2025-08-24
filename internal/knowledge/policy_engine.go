package knowledge

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// PolicyEngine executa raciocínio lógico sobre políticas e regras
type PolicyEngine struct {
	knowledgeBase *KnowledgeBaseManager
}

// NewPolicyEngine cria uma nova instância do motor de políticas
func NewPolicyEngine(kb *KnowledgeBaseManager) *PolicyEngine {
	return &PolicyEngine{
		knowledgeBase: kb,
	}
}

// AnalyzeScenario analisa um cenário complexo aplicando múltiplas regras
func (pe *PolicyEngine) AnalyzeScenario(scenario map[string]interface{}) (*ConsultationResult, error) {
	startTime := time.Now()

	// Extrair informações do cenário
	query := pe.extractQuery(scenario)

	// Identificar regras aplicáveis
	applicableRules, err := pe.findApplicableRules(scenario)
	if err != nil {
		return nil, fmt.Errorf("erro ao encontrar regras aplicáveis: %w", err)
	}

	// Executar raciocínio passo a passo
	reasoningSteps, err := pe.executeReasoning(scenario, applicableRules)
	if err != nil {
		return nil, fmt.Errorf("erro durante raciocínio: %w", err)
	}

	// Detectar conflitos
	conflicts := pe.detectConflicts(applicableRules)

	// Gerar resposta final
	answer := pe.generateAnswer(reasoningSteps, conflicts)

	// Calcular confiança
	confidence := pe.calculateConfidence(reasoningSteps, conflicts)

	// Extrair citações
	citations := pe.extractCitations(applicableRules)

	// Gerar recomendações
	recommendations := pe.generateRecommendations(reasoningSteps, conflicts)

	// Identificar ambiguidades
	ambiguities := pe.identifyAmbiguities(conflicts, reasoningSteps)

	result := &ConsultationResult{
		Query:           query,
		Answer:          answer,
		Confidence:      confidence,
		ReasoningSteps:  reasoningSteps,
		Sources:         citations,
		RelatedTopics:   pe.findRelatedTopics(applicableRules),
		Ambiguities:     ambiguities,
		Recommendations: recommendations,
		ProcessingTime:  time.Since(startTime),
	}

	return result, nil
}

// extractQuery extrai a query do cenário
func (pe *PolicyEngine) extractQuery(scenario map[string]interface{}) string {
	if query, exists := scenario["query"]; exists {
		if queryStr, ok := query.(string); ok {
			return queryStr
		}
	}
	return "Análise de cenário de Vale Refeição"
}

// findApplicableRules encontra as regras aplicáveis ao cenário
func (pe *PolicyEngine) findApplicableRules(scenario map[string]interface{}) ([]KnowledgeItem, error) {
	var applicableRules []KnowledgeItem

	// Buscar por palavras-chave do cenário
	keywords := pe.extractScenarioKeywords(scenario)

	for _, keyword := range keywords {
		results, err := pe.knowledgeBase.Search(keyword, 10)
		if err != nil {
			continue
		}

		for _, result := range results {
			if pe.isRuleApplicable(result.Item, scenario) {
				applicableRules = append(applicableRules, result.Item)
			}
		}
	}

	// Remover duplicatas
	applicableRules = pe.deduplicateRules(applicableRules)

	// Ordenar por prioridade
	sort.Slice(applicableRules, func(i, j int) bool {
		return pe.getRulePriority(applicableRules[i]) < pe.getRulePriority(applicableRules[j])
	})

	return applicableRules, nil
}

// extractScenarioKeywords extrai palavras-chave do cenário
func (pe *PolicyEngine) extractScenarioKeywords(scenario map[string]interface{}) []string {
	var keywords []string

	// Adicionar palavras-chave específicas baseadas no tipo de cenário
	if query, exists := scenario["query"]; exists {
		if queryStr, ok := query.(string); ok {
			words := strings.Fields(strings.ToLower(queryStr))
			keywords = append(keywords, words...)
		}
	}

	// Adicionar palavras-chave baseadas em campos específicos
	if _, exists := scenario["data_admissao"]; exists {
		keywords = append(keywords, "admissao", "data-quebrada")
	}

	if _, exists := scenario["data_desligamento"]; exists {
		keywords = append(keywords, "desligamento", "data-quebrada")
	}

	if _, exists := scenario["afastamento"]; exists {
		keywords = append(keywords, "afastamento", "dias-trabalhados")
	}

	if _, exists := scenario["ferias"]; exists {
		keywords = append(keywords, "ferias", "dias-trabalhados")
	}

	if _, exists := scenario["tipo_colaborador"]; exists {
		keywords = append(keywords, "elegibilidade", "tipos-colaborador")
	}

	// Sempre incluir palavras-chave base de VR
	keywords = append(keywords, "vr", "vale", "refeicao", "calculo")

	return pe.deduplicateStrings(keywords)
}

// isRuleApplicable verifica se uma regra é aplicável ao cenário
func (pe *PolicyEngine) isRuleApplicable(rule KnowledgeItem, scenario map[string]interface{}) bool {
	// Regras sempre aplicáveis para VR
	vrCategories := []string{"vr", "vale-refeicao", "calculo", "elegibilidade"}

	for _, category := range rule.Categories {
		for _, vrCat := range vrCategories {
			if strings.Contains(strings.ToLower(category), vrCat) {
				return true
			}
		}
	}

	// Verificar aplicabilidade baseada no conteúdo da regra
	content := strings.ToLower(rule.Content)
	title := strings.ToLower(rule.Title)

	// Cenários específicos
	if _, exists := scenario["data_admissao"]; exists {
		if strings.Contains(content, "admiss") || strings.Contains(title, "admiss") {
			return true
		}
	}

	if _, exists := scenario["data_desligamento"]; exists {
		if strings.Contains(content, "desligament") || strings.Contains(title, "desligament") {
			return true
		}
	}

	return false
}

// executeReasoning executa o raciocínio passo a passo
func (pe *PolicyEngine) executeReasoning(scenario map[string]interface{}, rules []KnowledgeItem) ([]ReasoningStep, error) {
	var steps []ReasoningStep
	stepNum := 1

	// Primeiro passo: Identificar o contexto
	contextStep := ReasoningStep{
		Step:        stepNum,
		Description: "Identificação do contexto do cenário",
		Logic:       "Analisar informações fornecidas para determinar o escopo da consulta",
		Result:      pe.describeScenario(scenario),
		Confidence:  0.9,
		Citations:   []Citation{},
	}
	steps = append(steps, contextStep)
	stepNum++

	// Aplicar regras por prioridade
	for _, rule := range rules {
		if pe.shouldApplyRule(rule, scenario) {
			step := pe.applyRule(rule, scenario, stepNum)
			steps = append(steps, step)
			stepNum++
		}
	}

	// Passo final: Consolidar resultado
	if len(steps) > 1 {
		finalStep := ReasoningStep{
			Step:        stepNum,
			Description: "Consolidação do resultado final",
			Logic:       "Combinar resultados de todas as regras aplicadas",
			Result:      pe.consolidateResults(steps),
			Confidence:  pe.calculateStepsConfidence(steps),
			Citations:   pe.consolidateCitations(steps),
		}
		steps = append(steps, finalStep)
	}

	return steps, nil
}

// applyRule aplica uma regra específica ao cenário
func (pe *PolicyEngine) applyRule(rule KnowledgeItem, scenario map[string]interface{}, stepNum int) ReasoningStep {
	var rulesUsed []string
	rulesUsed = append(rulesUsed, rule.ID)

	// Determinar a lógica aplicada baseada no tipo de regra
	logic := pe.determineLogic(rule, scenario)

	// Calcular resultado da aplicação da regra
	result := pe.calculateRuleResult(rule, scenario)

	// Calcular confiança baseada na clareza da regra
	confidence := pe.calculateRuleConfidence(rule, scenario)

	// Criar citação
	citation := Citation{
		Source:      rule.Source,
		Date:        rule.EffectiveDate,
		Reliability: "high",
	}

	return ReasoningStep{
		Step:        stepNum,
		Description: fmt.Sprintf("Aplicação da regra: %s", rule.Title),
		RulesUsed:   rulesUsed,
		Logic:       logic,
		Result:      result,
		Confidence:  confidence,
		Citations:   []Citation{citation},
	}
}

// determineLogic determina a lógica a ser aplicada
func (pe *PolicyEngine) determineLogic(rule KnowledgeItem, scenario map[string]interface{}) string {
	switch rule.Type {
	case "policy":
		return fmt.Sprintf("Aplicar política '%s': %s", rule.Title, rule.Content)
	case "regulation":
		return fmt.Sprintf("Verificar conformidade com regulamentação '%s': %s", rule.Title, rule.Content)
	case "business_rule":
		return fmt.Sprintf("Executar regra de negócio '%s': %s", rule.Title, rule.Content)
	default:
		return fmt.Sprintf("Aplicar regra '%s': %s", rule.Title, rule.Content)
	}
}

// calculateRuleResult calcula o resultado da aplicação de uma regra
func (pe *PolicyEngine) calculateRuleResult(rule KnowledgeItem, scenario map[string]interface{}) string {
	// Lógica específica baseada no ID da regra
	switch rule.ID {
	case "vr_001":
		return pe.checkEligibility(scenario)
	case "vr_002":
		return pe.calculateVRValue(scenario)
	case "vr_004":
		return pe.calculateAdmissionVR(scenario)
	case "vr_005":
		return pe.calculateTerminationVR(scenario)
	case "vr_006":
		return pe.calculateAbsenceVR(scenario)
	case "vr_007":
		return pe.calculateVacationVR(scenario)
	default:
		return fmt.Sprintf("Regra '%s' aplicada ao cenário", rule.Title)
	}
}

// Métodos específicos de cálculo
func (pe *PolicyEngine) checkEligibility(scenario map[string]interface{}) string {
	if tipoColaborador, exists := scenario["tipo_colaborador"]; exists {
		if tipo, ok := tipoColaborador.(string); ok {
			excludedTypes := []string{"diretor", "estagiario", "aprendiz", "terceirizado"}
			for _, excluded := range excludedTypes {
				if strings.Contains(strings.ToLower(tipo), excluded) {
					return fmt.Sprintf("Colaborador do tipo '%s' NÃO tem direito ao Vale Refeição", tipo)
				}
			}
			return fmt.Sprintf("Colaborador do tipo '%s' TEM direito ao Vale Refeição", tipo)
		}
	}
	return "Elegibilidade depende do tipo de colaborador (não especificado no cenário)"
}

func (pe *PolicyEngine) calculateVRValue(scenario map[string]interface{}) string {
	baseValue := 30.0
	companyPercentage := 0.8
	employeePercentage := 0.2

	if dias, exists := scenario["dias_uteis"]; exists {
		if diasUteis, ok := dias.(float64); ok {
			totalValue := baseValue * diasUteis
			companyValue := totalValue * companyPercentage
			employeeDiscount := totalValue * employeePercentage

			return fmt.Sprintf("Valor total: R$ %.2f (Empresa: R$ %.2f, Desconto colaborador: R$ %.2f)",
				totalValue, companyValue, employeeDiscount)
		}
	}
	return "Cálculo do valor depende do número de dias úteis trabalhados"
}

func (pe *PolicyEngine) calculateAdmissionVR(scenario map[string]interface{}) string {
	if dataAdmissao, exists := scenario["data_admissao"]; exists {
		if dataStr, ok := dataAdmissao.(string); ok {
			// Parse da data (assumindo formato YYYY-MM-DD)
			if len(dataStr) >= 10 {
				day := dataStr[8:10]
				if dayInt, err := time.Parse("02", day); err == nil {
					if dayInt.Day() <= 15 {
						return "Admissão até dia 15: Direito a VR INTEGRAL do mês"
					} else {
						return "Admissão após dia 15: Direito a 50% do VR mensal"
					}
				}
			}
		}
	}
	return "Cálculo depende da data de admissão (não especificada ou inválida)"
}

func (pe *PolicyEngine) calculateTerminationVR(scenario map[string]interface{}) string {
	if dataDesligamento, exists := scenario["data_desligamento"]; exists {
		if dataStr, ok := dataDesligamento.(string); ok {
			return fmt.Sprintf("Desligamento em %s: VR proporcional aos dias trabalhados. Verificar se comunicado foi feito até dia 15", dataStr)
		}
	}
	return "Cálculo de desligamento depende da data (não especificada)"
}

func (pe *PolicyEngine) calculateAbsenceVR(scenario map[string]interface{}) string {
	if afastamento, exists := scenario["afastamento"]; exists {
		if afastamentoMap, ok := afastamento.(map[string]interface{}); ok {
			if dias, ok := afastamentoMap["dias"].(float64); ok {
				if dias > 15 {
					return fmt.Sprintf("Afastamento de %.0f dias: SEM direito ao VR (superior a 15 dias)", dias)
				} else {
					return fmt.Sprintf("Afastamento de %.0f dias: COM direito proporcional ao VR", dias)
				}
			}
		}
	}
	return "Tratamento de afastamento depende da duração (não especificada)"
}

func (pe *PolicyEngine) calculateVacationVR(scenario map[string]interface{}) string {
	if ferias, exists := scenario["ferias"]; exists {
		if feriasMap, ok := ferias.(map[string]interface{}); ok {
			if periodo, ok := feriasMap["periodo"].(string); ok {
				return fmt.Sprintf("Período de férias (%s): SEM direito ao VR durante as férias", periodo)
			}
		}
	}
	return "Colaborador em férias não recebe VR no período"
}

// Métodos auxiliares
func (pe *PolicyEngine) shouldApplyRule(rule KnowledgeItem, scenario map[string]interface{}) bool {
	// Regras sempre aplicáveis para cenários de VR
	alwaysApplicable := []string{"vr_001", "vr_002"}
	for _, id := range alwaysApplicable {
		if rule.ID == id {
			return true
		}
	}

	// Aplicar baseado no contexto do cenário
	if _, exists := scenario["data_admissao"]; exists && rule.ID == "vr_004" {
		return true
	}

	if _, exists := scenario["data_desligamento"]; exists && rule.ID == "vr_005" {
		return true
	}

	if _, exists := scenario["afastamento"]; exists && rule.ID == "vr_006" {
		return true
	}

	if _, exists := scenario["ferias"]; exists && rule.ID == "vr_007" {
		return true
	}

	return false
}

func (pe *PolicyEngine) describeScenario(scenario map[string]interface{}) string {
	var parts []string

	if query, exists := scenario["query"]; exists {
		if queryStr, ok := query.(string); ok {
			parts = append(parts, fmt.Sprintf("Consulta: %s", queryStr))
		}
	}

	if tipoColaborador, exists := scenario["tipo_colaborador"]; exists {
		if tipo, ok := tipoColaborador.(string); ok {
			parts = append(parts, fmt.Sprintf("Tipo de colaborador: %s", tipo))
		}
	}

	if dataAdmissao, exists := scenario["data_admissao"]; exists {
		if data, ok := dataAdmissao.(string); ok {
			parts = append(parts, fmt.Sprintf("Data de admissão: %s", data))
		}
	}

	if dataDesligamento, exists := scenario["data_desligamento"]; exists {
		if data, ok := dataDesligamento.(string); ok {
			parts = append(parts, fmt.Sprintf("Data de desligamento: %s", data))
		}
	}

	if len(parts) == 0 {
		return "Cenário genérico de consulta sobre Vale Refeição"
	}

	return strings.Join(parts, "; ")
}

// Outros métodos auxiliares
func (pe *PolicyEngine) deduplicateRules(rules []KnowledgeItem) []KnowledgeItem {
	seen := make(map[string]bool)
	var result []KnowledgeItem

	for _, rule := range rules {
		if !seen[rule.ID] {
			seen[rule.ID] = true
			result = append(result, rule)
		}
	}

	return result
}

func (pe *PolicyEngine) deduplicateStrings(strs []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, str := range strs {
		if !seen[str] && str != "" {
			seen[str] = true
			result = append(result, str)
		}
	}

	return result
}

func (pe *PolicyEngine) getRulePriority(rule KnowledgeItem) int {
	// Prioridades baseadas no tipo e ID
	if rule.Type == "regulation" {
		return 1 // Regulamentações têm prioridade máxima
	}

	if rule.Type == "policy" {
		return 2 // Políticas em segundo
	}

	return 3 // Regras de negócio por último
}

func (pe *PolicyEngine) calculateRuleConfidence(rule KnowledgeItem, scenario map[string]interface{}) float64 {
	confidence := 0.8 // Base confidence

	// Aumentar confiança para regulamentações
	if rule.Type == "regulation" {
		confidence = 0.95
	}

	// Diminuir confiança para regras muito antigas
	if time.Since(rule.EffectiveDate) > 2*365*24*time.Hour {
		confidence *= 0.9
	}

	return confidence
}

func (pe *PolicyEngine) calculateStepsConfidence(steps []ReasoningStep) float64 {
	if len(steps) == 0 {
		return 0.0
	}

	total := 0.0
	for _, step := range steps {
		total += step.Confidence
	}

	return total / float64(len(steps))
}

func (pe *PolicyEngine) consolidateResults(steps []ReasoningStep) string {
	var results []string

	for i, step := range steps {
		if i == 0 || i == len(steps)-1 {
			continue // Skip context and final steps
		}
		if step.Result != "" {
			results = append(results, step.Result)
		}
	}

	if len(results) == 0 {
		return "Análise concluída com base nas regras aplicáveis"
	}

	return strings.Join(results, "; ")
}

func (pe *PolicyEngine) consolidateCitations(steps []ReasoningStep) []Citation {
	var citations []Citation
	seen := make(map[string]bool)

	for _, step := range steps {
		for _, citation := range step.Citations {
			key := citation.Source + citation.Date.String()
			if !seen[key] {
				seen[key] = true
				citations = append(citations, citation)
			}
		}
	}

	return citations
}

func (pe *PolicyEngine) detectConflicts(rules []KnowledgeItem) []ConflictDetection {
	// Implementação básica de detecção de conflitos
	var conflicts []ConflictDetection

	// Verificar por regras contraditórias (implementação simplificada)
	for i, rule1 := range rules {
		for j, rule2 := range rules {
			if i >= j {
				continue
			}

			if pe.rulesConflict(rule1, rule2) {
				conflict := ConflictDetection{
					Type:          "contradiction",
					Description:   fmt.Sprintf("Possível conflito entre %s e %s", rule1.Title, rule2.Title),
					RulesInvolved: []string{rule1.ID, rule2.ID},
					Severity:      "medium",
					Priority:      2,
				}
				conflicts = append(conflicts, conflict)
			}
		}
	}

	return conflicts
}

func (pe *PolicyEngine) rulesConflict(rule1, rule2 KnowledgeItem) bool {
	// Implementação básica de detecção de conflito
	// Regras do mesmo tipo com conteúdo contraditório
	if rule1.Type == rule2.Type {
		content1 := strings.ToLower(rule1.Content)
		content2 := strings.ToLower(rule2.Content)

		// Procurar por palavras contraditórias
		if (strings.Contains(content1, "não") && !strings.Contains(content2, "não")) ||
			(!strings.Contains(content1, "não") && strings.Contains(content2, "não")) {
			return true
		}
	}

	return false
}

func (pe *PolicyEngine) generateAnswer(steps []ReasoningStep, conflicts []ConflictDetection) string {
	if len(steps) == 0 {
		return "Não foi possível gerar uma resposta baseada nas informações disponíveis."
	}

	// Usar o resultado consolidado do último passo
	lastStep := steps[len(steps)-1]
	answer := lastStep.Result

	// Adicionar informação sobre conflitos se existirem
	if len(conflicts) > 0 {
		answer += fmt.Sprintf(" ATENÇÃO: Foram detectados %d possíveis conflitos entre as regras aplicadas.", len(conflicts))
	}

	return answer
}

func (pe *PolicyEngine) calculateConfidence(steps []ReasoningStep, conflicts []ConflictDetection) float64 {
	baseConfidence := pe.calculateStepsConfidence(steps)

	// Reduzir confiança baseado no número de conflitos
	conflictPenalty := float64(len(conflicts)) * 0.1

	confidence := baseConfidence - conflictPenalty
	if confidence < 0 {
		confidence = 0
	}
	if confidence > 1 {
		confidence = 1
	}

	return confidence
}

func (pe *PolicyEngine) extractCitations(rules []KnowledgeItem) []Citation {
	var citations []Citation
	seen := make(map[string]bool)

	for _, rule := range rules {
		key := rule.Source + rule.EffectiveDate.String()
		if !seen[key] {
			seen[key] = true
			citation := Citation{
				Source:      rule.Source,
				Date:        rule.EffectiveDate,
				Reliability: "high",
			}
			citations = append(citations, citation)
		}
	}

	return citations
}

func (pe *PolicyEngine) findRelatedTopics(rules []KnowledgeItem) []string {
	var topics []string
	seen := make(map[string]bool)

	for _, rule := range rules {
		for _, category := range rule.Categories {
			if !seen[category] {
				seen[category] = true
				topics = append(topics, category)
			}
		}
	}

	return topics
}

func (pe *PolicyEngine) generateRecommendations(steps []ReasoningStep, conflicts []ConflictDetection) []string {
	var recommendations []string

	if len(conflicts) > 0 {
		recommendations = append(recommendations, "Revisar possíveis conflitos entre regras identificados")
	}

	// Analisar confiança dos passos
	lowConfidenceSteps := 0
	for _, step := range steps {
		if step.Confidence < 0.7 {
			lowConfidenceSteps++
		}
	}

	if lowConfidenceSteps > 0 {
		recommendations = append(recommendations, "Verificar manualmente resultados com baixa confiança")
	}

	// Recomendações específicas baseadas no contexto
	recommendations = append(recommendations, "Consultar documentação atualizada para confirmação")

	return recommendations
}

func (pe *PolicyEngine) identifyAmbiguities(conflicts []ConflictDetection, steps []ReasoningStep) []string {
	var ambiguities []string

	// Ambiguidades baseadas em conflitos
	for _, conflict := range conflicts {
		ambiguities = append(ambiguities, conflict.Description)
	}

	// Ambiguidades baseadas em baixa confiança
	for _, step := range steps {
		if step.Confidence < 0.5 {
			ambiguities = append(ambiguities, fmt.Sprintf("Baixa confiança no passo: %s", step.Description))
		}
	}

	return ambiguities
}
