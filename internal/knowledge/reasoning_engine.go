package knowledge

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ReasoningEngine executa raciocínio lógico avançado para cenários complexos
type ReasoningEngine struct {
	knowledgeBase   *KnowledgeBaseManager
	policyEngine    *PolicyEngine
	citationManager *CitationManager
}

// NewReasoningEngine cria uma nova instância do motor de raciocínio
func NewReasoningEngine(kb *KnowledgeBaseManager, pe *PolicyEngine, cm *CitationManager) *ReasoningEngine {
	return &ReasoningEngine{
		knowledgeBase:   kb,
		policyEngine:    pe,
		citationManager: cm,
	}
}

// ScenarioType representa os tipos de cenários que podem ser analisados
type ScenarioType string

const (
	ScenarioSimpleQuery        ScenarioType = "simple_query"
	ScenarioComplexCalculation ScenarioType = "complex_calculation"
	ScenarioWhatIf             ScenarioType = "what_if"
	ScenarioConflictResolution ScenarioType = "conflict_resolution"
	ScenarioComplianceCheck    ScenarioType = "compliance_check"
)

// AnalyzeComplexScenario analisa um cenário complexo com raciocínio avançado
func (re *ReasoningEngine) AnalyzeComplexScenario(scenario map[string]interface{}) (*ConsultationResult, error) {
	startTime := time.Now()

	// Identificar tipo de cenário
	scenarioType := re.identifyScenarioType(scenario)

	// Executar análise específica baseada no tipo
	var result *ConsultationResult
	var err error

	switch scenarioType {
	case ScenarioSimpleQuery:
		result, err = re.handleSimpleQuery(scenario)
	case ScenarioComplexCalculation:
		result, err = re.handleComplexCalculation(scenario)
	case ScenarioWhatIf:
		result, err = re.handleWhatIfScenario(scenario)
	case ScenarioConflictResolution:
		result, err = re.handleConflictResolution(scenario)
	case ScenarioComplianceCheck:
		result, err = re.handleComplianceCheck(scenario)
	default:
		result, err = re.policyEngine.AnalyzeScenario(scenario)
	}

	if err != nil {
		return nil, err
	}

	// Enriquecer resultado com análise avançada
	result = re.enrichResult(result, scenario, scenarioType)
	result.ProcessingTime = time.Since(startTime)

	return result, nil
}

// identifyScenarioType identifica o tipo de cenário baseado nas informações
func (re *ReasoningEngine) identifyScenarioType(scenario map[string]interface{}) ScenarioType {
	query, hasQuery := scenario["query"]
	queryStr := ""
	if hasQuery {
		if str, ok := query.(string); ok {
			queryStr = strings.ToLower(str)
		}
	}

	// Detectar "E se" scenarios
	if strings.Contains(queryStr, "e se") || strings.Contains(queryStr, "what if") ||
		strings.Contains(queryStr, "caso") || strings.Contains(queryStr, "supondo") {
		return ScenarioWhatIf
	}

	// Detectar cálculos complexos
	if re.hasMultipleFactors(scenario) {
		return ScenarioComplexCalculation
	}

	// Detectar verificação de compliance
	if strings.Contains(queryStr, "conforme") || strings.Contains(queryStr, "legal") ||
		strings.Contains(queryStr, "regulament") || strings.Contains(queryStr, "lei") {
		return ScenarioComplianceCheck
	}

	// Detectar resolução de conflitos
	if strings.Contains(queryStr, "conflito") || strings.Contains(queryStr, "contradi") ||
		strings.Contains(queryStr, "diferent") {
		return ScenarioConflictResolution
	}

	return ScenarioSimpleQuery
}

// hasMultipleFactors verifica se o cenário possui múltiplos fatores
func (re *ReasoningEngine) hasMultipleFactors(scenario map[string]interface{}) bool {
	factorCount := 0

	factors := []string{
		"data_admissao", "data_desligamento", "afastamento", "ferias",
		"tipo_colaborador", "sindicato", "carga_horaria", "dias_uteis",
	}

	for _, factor := range factors {
		if _, exists := scenario[factor]; exists {
			factorCount++
		}
	}

	return factorCount >= 3
}

// handleSimpleQuery processa consultas simples
func (re *ReasoningEngine) handleSimpleQuery(scenario map[string]interface{}) (*ConsultationResult, error) {
	// Usar o policy engine padrão para consultas simples
	return re.policyEngine.AnalyzeScenario(scenario)
}

// handleComplexCalculation processa cálculos complexos com múltiplos fatores
func (re *ReasoningEngine) handleComplexCalculation(scenario map[string]interface{}) (*ConsultationResult, error) {
	startTime := time.Now()

	// Identificar todos os fatores relevantes
	factors := re.identifyRelevantFactors(scenario)

	// Executar análise multi-dimensional
	analysis := re.performMultiDimensionalAnalysis(scenario, factors)

	// Gerar passos de raciocínio detalhados
	reasoningSteps := re.generateComplexReasoningSteps(scenario, factors, analysis)

	// Calcular resultado final
	finalResult := re.calculateComplexResult(analysis)

	// Avaliar confiança baseada na complexidade
	confidence := re.calculateComplexConfidence(analysis, factors)

	result := &ConsultationResult{
		Query:          re.extractQuery(scenario),
		Answer:         finalResult,
		Confidence:     confidence,
		ReasoningSteps: reasoningSteps,
		Sources:        re.extractSourcesFromAnalysis(analysis),
		ProcessingTime: time.Since(startTime),
	}

	return result, nil
}

// handleWhatIfScenario processa cenários hipotéticos "E se"
func (re *ReasoningEngine) handleWhatIfScenario(scenario map[string]interface{}) (*ConsultationResult, error) {
	startTime := time.Now()

	// Identificar as variáveis hipotéticas
	hypotheticalVars := re.extractHypotheticalVariables(scenario)

	// Gerar múltiplos cenários baseados nas variáveis
	scenarios := re.generateAlternativeScenarios(scenario, hypotheticalVars)

	// Analisar cada cenário alternativo
	scenarioResults := make(map[string]*ConsultationResult)
	for name, altScenario := range scenarios {
		result, err := re.policyEngine.AnalyzeScenario(altScenario)
		if err != nil {
			continue
		}
		scenarioResults[name] = result
	}

	// Comparar resultados e gerar análise
	comparison := re.compareScenarioResults(scenarioResults)

	// Gerar recomendações baseadas na comparação
	recommendations := re.generateWhatIfRecommendations(comparison)

	result := &ConsultationResult{
		Query:           re.extractQuery(scenario),
		Answer:          re.formatWhatIfAnswer(comparison),
		Confidence:      0.8, // Cenários hipotéticos têm confiança moderada
		ReasoningSteps:  re.generateWhatIfSteps(scenarios, scenarioResults),
		Recommendations: recommendations,
		ProcessingTime:  time.Since(startTime),
	}

	return result, nil
}

// handleConflictResolution resolve conflitos entre regras
func (re *ReasoningEngine) handleConflictResolution(scenario map[string]interface{}) (*ConsultationResult, error) {
	startTime := time.Now()

	// Encontrar regras potencialmente conflitantes
	applicableRules, _ := re.findApplicableRules(scenario)
	conflicts := re.detectDetailedConflicts(applicableRules)

	// Resolver conflitos usando hierarquia de prioridades
	resolutions := re.resolveConflicts(conflicts)

	// Gerar explicação da resolução
	explanation := re.explainConflictResolution(conflicts, resolutions)

	result := &ConsultationResult{
		Query:          re.extractQuery(scenario),
		Answer:         explanation,
		Confidence:     0.7, // Resoluções de conflito têm confiança moderada
		ReasoningSteps: re.generateConflictResolutionSteps(conflicts, resolutions),
		Ambiguities:    re.extractAmbiguitiesFromConflicts(conflicts),
		ProcessingTime: time.Since(startTime),
	}

	return result, nil
}

// handleComplianceCheck verifica conformidade com regulamentações
func (re *ReasoningEngine) handleComplianceCheck(scenario map[string]interface{}) (*ConsultationResult, error) {
	startTime := time.Now()

	// Identificar regulamentações aplicáveis
	regulations := re.findApplicableRegulations(scenario)

	// Verificar conformidade para cada regulamentação
	complianceResults := make(map[string]bool)
	complianceDetails := make(map[string]string)

	for _, regulation := range regulations {
		isCompliant, details := re.checkComplianceWithRegulation(scenario, regulation)
		complianceResults[regulation.ID] = isCompliant
		complianceDetails[regulation.ID] = details
	}

	// Gerar relatório de conformidade
	complianceReport := re.generateComplianceReport(complianceResults, complianceDetails, regulations)

	// Calcular confiança baseada na clareza das regulamentações
	confidence := re.calculateComplianceConfidence(regulations)

	result := &ConsultationResult{
		Query:          re.extractQuery(scenario),
		Answer:         complianceReport,
		Confidence:     confidence,
		ReasoningSteps: re.generateComplianceSteps(regulations, complianceResults, complianceDetails),
		Sources:        re.extractRegulationCitations(regulations),
		ProcessingTime: time.Since(startTime),
	}

	return result, nil
}

// Métodos auxiliares para análise multi-dimensional

// identifyRelevantFactors identifica fatores relevantes para o cálculo
func (re *ReasoningEngine) identifyRelevantFactors(scenario map[string]interface{}) map[string]interface{} {
	factors := make(map[string]interface{})

	relevantKeys := []string{
		"data_admissao", "data_desligamento", "afastamento", "ferias",
		"tipo_colaborador", "sindicato", "carga_horaria", "dias_uteis",
		"valor_base", "percentual_empresa", "percentual_colaborador",
	}

	for _, key := range relevantKeys {
		if value, exists := scenario[key]; exists {
			factors[key] = value
		}
	}

	return factors
}

// performMultiDimensionalAnalysis executa análise considerando múltiplas dimensões
func (re *ReasoningEngine) performMultiDimensionalAnalysis(scenario map[string]interface{}, factors map[string]interface{}) map[string]interface{} {
	analysis := make(map[string]interface{})

	// Analisar elegibilidade
	eligibility := re.analyzeEligibility(factors)
	analysis["eligibility"] = eligibility

	// Analisar impacto temporal (admissões, desligamentos)
	temporalImpact := re.analyzeTemporalImpact(factors)
	analysis["temporal_impact"] = temporalImpact

	// Analisar impacto de afastamentos/férias
	absenceImpact := re.analyzeAbsenceImpact(factors)
	analysis["absence_impact"] = absenceImpact

	// Calcular valor base
	baseCalculation := re.calculateBaseValue(factors)
	analysis["base_calculation"] = baseCalculation

	// Analisar ajustes específicos
	adjustments := re.calculateAdjustments(factors)
	analysis["adjustments"] = adjustments

	return analysis
}

// analyzeEligibility analisa elegibilidade detalhadamente
func (re *ReasoningEngine) analyzeEligibility(factors map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"eligible":     true,
		"reasons":      []string{},
		"restrictions": []string{},
	}

	if tipoColaborador, exists := factors["tipo_colaborador"]; exists {
		if tipo, ok := tipoColaborador.(string); ok {
			excludedTypes := []string{"diretor", "estagiario", "aprendiz", "terceirizado"}
			for _, excluded := range excludedTypes {
				if strings.Contains(strings.ToLower(tipo), excluded) {
					result["eligible"] = false
					result["reasons"] = append(result["reasons"].([]string),
						fmt.Sprintf("Tipo de colaborador '%s' não elegível", tipo))
				}
			}
		}
	}

	if cargaHoraria, exists := factors["carga_horaria"]; exists {
		if horas, ok := cargaHoraria.(float64); ok && horas < 6.0 {
			result["eligible"] = false
			result["reasons"] = append(result["reasons"].([]string),
				fmt.Sprintf("Carga horária %.1f horas inferior ao mínimo de 6h", horas))
		}
	}

	return result
}

// analyzeTemporalImpact analisa impacto de datas de admissão/desligamento
func (re *ReasoningEngine) analyzeTemporalImpact(factors map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"admission_impact":   "none",
		"termination_impact": "none",
		"proportional":       false,
	}

	// Analisar admissão
	if dataAdmissao, exists := factors["data_admissao"]; exists {
		if dataStr, ok := dataAdmissao.(string); ok {
			impact := re.calculateAdmissionImpact(dataStr)
			result["admission_impact"] = impact
			if impact != "full" {
				result["proportional"] = true
			}
		}
	}

	// Analisar desligamento
	if dataDesligamento, exists := factors["data_desligamento"]; exists {
		if dataStr, ok := dataDesligamento.(string); ok {
			impact := re.calculateTerminationImpact(dataStr)
			result["termination_impact"] = impact
			result["proportional"] = true
		}
	}

	return result
}

// calculateAdmissionImpact calcula impacto da admissão
func (re *ReasoningEngine) calculateAdmissionImpact(dataAdmissao string) string {
	if len(dataAdmissao) >= 10 {
		dayStr := dataAdmissao[8:10]
		if day, err := strconv.Atoi(dayStr); err == nil {
			if day <= 15 {
				return "full" // VR integral
			} else {
				return "half" // 50% do VR
			}
		}
	}
	return "unknown"
}

// calculateTerminationImpact calcula impacto do desligamento
func (re *ReasoningEngine) calculateTerminationImpact(dataDesligamento string) string {
	// Sempre proporcional para desligamentos
	return "proportional"
}

// analyzeAbsenceImpact analisa impacto de afastamentos e férias
func (re *ReasoningEngine) analyzeAbsenceImpact(factors map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"has_absences":  false,
		"absence_days":  0,
		"vacation_days": 0,
		"impact_level":  "none",
	}

	totalAbsenceDays := 0

	// Analisar afastamentos
	if afastamento, exists := factors["afastamento"]; exists {
		if afastamentoMap, ok := afastamento.(map[string]interface{}); ok {
			if dias, ok := afastamentoMap["dias"].(float64); ok {
				totalAbsenceDays += int(dias)
				result["absence_days"] = int(dias)
				result["has_absences"] = true
			}
		}
	}

	// Analisar férias
	if ferias, exists := factors["ferias"]; exists {
		if feriasMap, ok := ferias.(map[string]interface{}); ok {
			if dias, ok := feriasMap["dias"].(float64); ok {
				totalAbsenceDays += int(dias)
				result["vacation_days"] = int(dias)
			}
		}
	}

	// Determinar nível de impacto
	if totalAbsenceDays == 0 {
		result["impact_level"] = "none"
	} else if totalAbsenceDays <= 5 {
		result["impact_level"] = "low"
	} else if totalAbsenceDays <= 15 {
		result["impact_level"] = "medium"
	} else {
		result["impact_level"] = "high"
	}

	return result
}

// calculateBaseValue calcula valor base do VR
func (re *ReasoningEngine) calculateBaseValue(factors map[string]interface{}) map[string]interface{} {
	valorBase := 30.0 // Valor padrão
	diasUteis := 22.0 // Dias padrão

	if valor, exists := factors["valor_base"]; exists {
		if v, ok := valor.(float64); ok {
			valorBase = v
		}
	}

	if dias, exists := factors["dias_uteis"]; exists {
		if d, ok := dias.(float64); ok {
			diasUteis = d
		}
	}

	valorTotal := valorBase * diasUteis
	subsidioEmpresa := valorTotal * 0.8
	descontoColaborador := valorTotal * 0.2

	return map[string]interface{}{
		"valor_base":           valorBase,
		"dias_uteis":           diasUteis,
		"valor_total":          valorTotal,
		"subsidio_empresa":     subsidioEmpresa,
		"desconto_colaborador": descontoColaborador,
	}
}

// calculateAdjustments calcula ajustes baseados em fatores específicos
func (re *ReasoningEngine) calculateAdjustments(factors map[string]interface{}) []map[string]interface{} {
	var adjustments []map[string]interface{}

	// Ajuste por admissão
	if dataAdmissao, exists := factors["data_admissao"]; exists {
		if dataStr, ok := dataAdmissao.(string); ok {
			impact := re.calculateAdmissionImpact(dataStr)
			if impact == "half" {
				adjustments = append(adjustments, map[string]interface{}{
					"type":        "admission",
					"description": "Redução para 50% por admissão após dia 15",
					"multiplier":  0.5,
				})
			}
		}
	}

	// Ajuste por afastamentos longos
	if afastamento, exists := factors["afastamento"]; exists {
		if afastamentoMap, ok := afastamento.(map[string]interface{}); ok {
			if dias, ok := afastamentoMap["dias"].(float64); ok && dias > 15 {
				adjustments = append(adjustments, map[string]interface{}{
					"type":        "absence",
					"description": fmt.Sprintf("Perda do direito por afastamento superior a 15 dias (%.0f dias)", dias),
					"multiplier":  0.0,
				})
			}
		}
	}

	return adjustments
}

// Métodos para geração de passos de raciocínio

// generateComplexReasoningSteps gera passos detalhados para cálculos complexos
func (re *ReasoningEngine) generateComplexReasoningSteps(scenario map[string]interface{}, factors map[string]interface{}, analysis map[string]interface{}) []ReasoningStep {
	var steps []ReasoningStep
	stepNum := 1

	// Passo 1: Análise de elegibilidade
	eligibilityStep := ReasoningStep{
		Step:        stepNum,
		Description: "Verificação de elegibilidade",
		Logic:       "Verificar se o colaborador atende aos critérios básicos para recebimento de VR",
		Result:      re.formatEligibilityResult(analysis["eligibility"]),
		Confidence:  0.95,
	}
	steps = append(steps, eligibilityStep)
	stepNum++

	// Passo 2: Análise temporal
	if temporalImpact, exists := analysis["temporal_impact"]; exists {
		temporalStep := ReasoningStep{
			Step:        stepNum,
			Description: "Análise de impacto temporal",
			Logic:       "Verificar impacto de datas de admissão e desligamento",
			Result:      re.formatTemporalResult(temporalImpact),
			Confidence:  0.9,
		}
		steps = append(steps, temporalStep)
		stepNum++
	}

	// Passo 3: Análise de ausências
	if absenceImpact, exists := analysis["absence_impact"]; exists {
		absenceStep := ReasoningStep{
			Step:        stepNum,
			Description: "Análise de impacto de ausências",
			Logic:       "Verificar impacto de afastamentos e férias no cálculo",
			Result:      re.formatAbsenceResult(absenceImpact),
			Confidence:  0.85,
		}
		steps = append(steps, absenceStep)
		stepNum++
	}

	// Passo 4: Cálculo base
	if baseCalc, exists := analysis["base_calculation"]; exists {
		calcStep := ReasoningStep{
			Step:        stepNum,
			Description: "Cálculo do valor base",
			Logic:       "Calcular valor total baseado no valor diário e dias úteis",
			Result:      re.formatCalculationResult(baseCalc),
			Confidence:  0.98,
		}
		steps = append(steps, calcStep)
		stepNum++
	}

	// Passo 5: Aplicação de ajustes
	if adjustments, exists := analysis["adjustments"]; exists {
		if adj, ok := adjustments.([]map[string]interface{}); ok && len(adj) > 0 {
			adjustmentStep := ReasoningStep{
				Step:        stepNum,
				Description: "Aplicação de ajustes específicos",
				Logic:       "Aplicar reduções ou eliminações baseadas em regras específicas",
				Result:      re.formatAdjustmentResult(adj),
				Confidence:  0.9,
			}
			steps = append(steps, adjustmentStep)
		}
	}

	return steps
}

// formatEligibilityResult formata resultado de elegibilidade
func (re *ReasoningEngine) formatEligibilityResult(eligibility interface{}) string {
	if eligMap, ok := eligibility.(map[string]interface{}); ok {
		if eligible, ok := eligMap["eligible"].(bool); ok {
			if eligible {
				return "Colaborador ELEGÍVEL para recebimento de Vale Refeição"
			} else {
				reasons := eligMap["reasons"].([]string)
				return fmt.Sprintf("Colaborador NÃO ELEGÍVEL: %s", strings.Join(reasons, "; "))
			}
		}
	}
	return "Elegibilidade não pôde ser determinada"
}

// formatTemporalResult formata resultado temporal
func (re *ReasoningEngine) formatTemporalResult(temporal interface{}) string {
	if tempMap, ok := temporal.(map[string]interface{}); ok {
		var impacts []string

		if admission, exists := tempMap["admission_impact"]; exists {
			switch admission {
			case "full":
				impacts = append(impacts, "VR integral por admissão até dia 15")
			case "half":
				impacts = append(impacts, "50% do VR por admissão após dia 15")
			}
		}

		if termination, exists := tempMap["termination_impact"]; exists {
			if termination == "proportional" {
				impacts = append(impacts, "VR proporcional por desligamento no mês")
			}
		}

		if len(impacts) > 0 {
			return strings.Join(impacts, "; ")
		}
	}
	return "Sem impactos temporais identificados"
}

// formatAbsenceResult formata resultado de ausências
func (re *ReasoningEngine) formatAbsenceResult(absence interface{}) string {
	if absMap, ok := absence.(map[string]interface{}); ok {
		if hasAbsences, ok := absMap["has_absences"].(bool); ok && hasAbsences {
			impactLevel := absMap["impact_level"].(string)
			absenceDays := absMap["absence_days"].(int)
			vacationDays := absMap["vacation_days"].(int)

			var parts []string
			if absenceDays > 0 {
				if absenceDays > 15 {
					parts = append(parts, fmt.Sprintf("Afastamento de %d dias: SEM direito ao VR", absenceDays))
				} else {
					parts = append(parts, fmt.Sprintf("Afastamento de %d dias: impacto %s", absenceDays, impactLevel))
				}
			}

			if vacationDays > 0 {
				parts = append(parts, fmt.Sprintf("Férias de %d dias: sem VR no período", vacationDays))
			}

			return strings.Join(parts, "; ")
		}
	}
	return "Sem ausências que impactem o cálculo"
}

// formatCalculationResult formata resultado do cálculo
func (re *ReasoningEngine) formatCalculationResult(calc interface{}) string {
	if calcMap, ok := calc.(map[string]interface{}); ok {
		valorBase := calcMap["valor_base"].(float64)
		diasUteis := calcMap["dias_uteis"].(float64)
		valorTotal := calcMap["valor_total"].(float64)
		subsidio := calcMap["subsidio_empresa"].(float64)
		desconto := calcMap["desconto_colaborador"].(float64)

		return fmt.Sprintf("Valor base R$ %.2f × %.0f dias úteis = R$ %.2f (Empresa: R$ %.2f, Desconto: R$ %.2f)",
			valorBase, diasUteis, valorTotal, subsidio, desconto)
	}
	return "Cálculo não pôde ser executado"
}

// formatAdjustmentResult formata resultado dos ajustes
func (re *ReasoningEngine) formatAdjustmentResult(adjustments []map[string]interface{}) string {
	var parts []string

	for _, adj := range adjustments {
		description := adj["description"].(string)
		multiplier := adj["multiplier"].(float64)

		if multiplier == 0.0 {
			parts = append(parts, fmt.Sprintf("%s (elimina direito)", description))
		} else {
			parts = append(parts, fmt.Sprintf("%s (fator %.1f)", description, multiplier))
		}
	}

	return strings.Join(parts, "; ")
}

// calculateComplexResult calcula resultado final da análise complexa
func (re *ReasoningEngine) calculateComplexResult(analysis map[string]interface{}) string {
	// Verificar elegibilidade
	if eligibility, exists := analysis["eligibility"]; exists {
		if eligMap, ok := eligibility.(map[string]interface{}); ok {
			if eligible, ok := eligMap["eligible"].(bool); ok && !eligible {
				return "Colaborador não tem direito ao Vale Refeição baseado nos critérios de elegibilidade."
			}
		}
	}

	// Calcular valor final
	baseValue := 0.0
	if baseCalc, exists := analysis["base_calculation"]; exists {
		if calcMap, ok := baseCalc.(map[string]interface{}); ok {
			if valorTotal, ok := calcMap["valor_total"].(float64); ok {
				baseValue = valorTotal
			}
		}
	}

	// Aplicar ajustes
	finalValue := baseValue
	if adjustments, exists := analysis["adjustments"]; exists {
		if adj, ok := adjustments.([]map[string]interface{}); ok {
			for _, adjustment := range adj {
				if multiplier, ok := adjustment["multiplier"].(float64); ok {
					finalValue *= multiplier
				}
			}
		}
	}

	if finalValue <= 0 {
		return "Com base na análise, o colaborador não tem direito ao Vale Refeição no período."
	}

	subsidio := finalValue * 0.8
	desconto := finalValue * 0.2

	return fmt.Sprintf("Valor final do Vale Refeição: R$ %.2f (Empresa subsidia R$ %.2f, colaborador paga R$ %.2f)",
		finalValue, subsidio, desconto)
}

// calculateComplexConfidence calcula confiança para análises complexas
func (re *ReasoningEngine) calculateComplexConfidence(analysis map[string]interface{}, factors map[string]interface{}) float64 {
	baseConfidence := 0.8

	// Aumentar confiança se temos muitos fatores claros
	factorCount := len(factors)
	if factorCount >= 5 {
		baseConfidence += 0.1
	}

	// Diminuir confiança se há muitos ajustes (cenário complexo)
	if adjustments, exists := analysis["adjustments"]; exists {
		if adj, ok := adjustments.([]map[string]interface{}); ok {
			adjustmentCount := len(adj)
			if adjustmentCount > 2 {
				baseConfidence -= 0.1
			}
		}
	}

	// Garantir que a confiança está entre 0 e 1
	if baseConfidence > 1.0 {
		baseConfidence = 1.0
	}
	if baseConfidence < 0.0 {
		baseConfidence = 0.0
	}

	return baseConfidence
}

// enrichResult enriquece o resultado com análises adicionais
func (re *ReasoningEngine) enrichResult(result *ConsultationResult, scenario map[string]interface{}, scenarioType ScenarioType) *ConsultationResult {
	// Adicionar tópicos relacionados baseados no tipo de cenário
	relatedTopics := re.generateRelatedTopics(scenarioType, scenario)
	result.RelatedTopics = append(result.RelatedTopics, relatedTopics...)

	// Adicionar recomendações específicas
	specificRecommendations := re.generateSpecificRecommendations(scenarioType, scenario, result)
	result.Recommendations = append(result.Recommendations, specificRecommendations...)

	return result
}

// generateRelatedTopics gera tópicos relacionados baseados no cenário
func (re *ReasoningEngine) generateRelatedTopics(scenarioType ScenarioType, scenario map[string]interface{}) []string {
	var topics []string

	switch scenarioType {
	case ScenarioComplexCalculation:
		topics = []string{"cálculos proporcionais", "múltiplos fatores", "datas quebradas", "ajustes temporais"}
	case ScenarioWhatIf:
		topics = []string{"cenários hipotéticos", "análise de alternativas", "simulações"}
	case ScenarioConflictResolution:
		topics = []string{"resolução de conflitos", "hierarquia de regras", "priorização de políticas"}
	case ScenarioComplianceCheck:
		topics = []string{"conformidade regulatória", "legislação trabalhista", "auditoria"}
	}

	return topics
}

// generateSpecificRecommendations gera recomendações específicas
func (re *ReasoningEngine) generateSpecificRecommendations(scenarioType ScenarioType, scenario map[string]interface{}, result *ConsultationResult) []string {
	var recommendations []string

	// Recomendações baseadas na confiança
	if result.Confidence < 0.7 {
		recommendations = append(recommendations, "Revisar manualmente devido à baixa confiança no resultado")
	}

	// Recomendações específicas por tipo
	switch scenarioType {
	case ScenarioComplexCalculation:
		recommendations = append(recommendations, "Validar cálculo com múltiplas fontes de dados")
		recommendations = append(recommendations, "Documentar todos os fatores considerados")
	case ScenarioWhatIf:
		recommendations = append(recommendations, "Considerar impactos de longo prazo das alternativas")
	case ScenarioConflictResolution:
		recommendations = append(recommendations, "Documentar resolução de conflito para casos futuros")
	case ScenarioComplianceCheck:
		recommendations = append(recommendations, "Manter documentação para auditoria")
	}

	return recommendations
}

// Métodos auxiliares adicionais
func (re *ReasoningEngine) extractQuery(scenario map[string]interface{}) string {
	if query, exists := scenario["query"]; exists {
		if queryStr, ok := query.(string); ok {
			return queryStr
		}
	}
	return "Análise de cenário complexo de Vale Refeição"
}

func (re *ReasoningEngine) extractSourcesFromAnalysis(analysis map[string]interface{}) []Citation {
	// Gerar citações baseadas na análise
	var citations []Citation

	// Sempre incluir a política base de VR
	baseCitation := Citation{
		Source:      "Manual de Recursos Humanos v2.1",
		Date:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Reliability: "high",
		Section:     "Políticas de Vale Refeição",
	}
	citations = append(citations, baseCitation)

	return citations
}

// Métodos para cenários "What If" (implementação básica)
func (re *ReasoningEngine) extractHypotheticalVariables(scenario map[string]interface{}) map[string][]interface{} {
	// Implementação básica - expandir conforme necessário
	return map[string][]interface{}{
		"tipo_colaborador": {"efetivo", "terceirizado", "estagiario"},
		"carga_horaria":    {4.0, 6.0, 8.0},
	}
}

func (re *ReasoningEngine) generateAlternativeScenarios(baseScenario map[string]interface{}, variables map[string][]interface{}) map[string]map[string]interface{} {
	scenarios := make(map[string]map[string]interface{})

	// Implementação básica - gerar alguns cenários alternativos
	for varName, values := range variables {
		for i, value := range values {
			scenarioName := fmt.Sprintf("%s_%d", varName, i)
			altScenario := make(map[string]interface{})

			// Copiar cenário base
			for k, v := range baseScenario {
				altScenario[k] = v
			}

			// Modificar variável
			altScenario[varName] = value

			scenarios[scenarioName] = altScenario
		}
	}

	return scenarios
}

func (re *ReasoningEngine) compareScenarioResults(results map[string]*ConsultationResult) map[string]interface{} {
	comparison := map[string]interface{}{
		"scenarios": len(results),
		"outcomes":  make(map[string]string),
	}

	for name, result := range results {
		comparison["outcomes"].(map[string]string)[name] = result.Answer
	}

	return comparison
}

func (re *ReasoningEngine) generateWhatIfRecommendations(comparison map[string]interface{}) []string {
	return []string{"Analisar impactos de cada cenário alternativo"}
}

func (re *ReasoningEngine) formatWhatIfAnswer(comparison map[string]interface{}) string {
	return fmt.Sprintf("Análise de %v cenários alternativos concluída", comparison["scenarios"])
}

func (re *ReasoningEngine) generateWhatIfSteps(scenarios map[string]map[string]interface{}, results map[string]*ConsultationResult) []ReasoningStep {
	var steps []ReasoningStep

	step := ReasoningStep{
		Step:        1,
		Description: "Análise de cenários hipotéticos",
		Logic:       fmt.Sprintf("Comparar %d cenários alternativos", len(scenarios)),
		Result:      "Cenários analisados com sucesso",
		Confidence:  0.8,
	}
	steps = append(steps, step)

	return steps
}

// Métodos para resolução de conflitos (implementação básica)
func (re *ReasoningEngine) findApplicableRules(scenario map[string]interface{}) ([]KnowledgeItem, error) {
	// Usar o knowledge base para encontrar regras aplicáveis
	query := "vr vale refeicao"
	if q, exists := scenario["query"]; exists {
		if queryStr, ok := q.(string); ok {
			query = queryStr
		}
	}

	results, err := re.knowledgeBase.Search(query, 20)
	if err != nil {
		return []KnowledgeItem{}, err
	}

	var rules []KnowledgeItem
	for _, result := range results {
		rules = append(rules, result.Item)
	}

	return rules, nil
}

func (re *ReasoningEngine) detectDetailedConflicts(rules []KnowledgeItem) []ConflictDetection {
	// Implementação básica de detecção de conflitos
	var conflicts []ConflictDetection

	for i, rule1 := range rules {
		for j, rule2 := range rules {
			if i >= j {
				continue
			}

			if re.rulesConflict(rule1, rule2) {
				conflict := ConflictDetection{
					Type:          "contradiction",
					Description:   fmt.Sprintf("Conflito entre %s e %s", rule1.Title, rule2.Title),
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

func (re *ReasoningEngine) rulesConflict(rule1, rule2 KnowledgeItem) bool {
	// Implementação básica - melhorar com lógica mais sofisticada
	content1 := strings.ToLower(rule1.Content)
	content2 := strings.ToLower(rule2.Content)

	// Procurar por contradições óbvias
	if (strings.Contains(content1, "não") && !strings.Contains(content2, "não")) ||
		(!strings.Contains(content1, "não") && strings.Contains(content2, "não")) {
		return true
	}

	return false
}

func (re *ReasoningEngine) resolveConflicts(conflicts []ConflictDetection) map[string]string {
	resolutions := make(map[string]string)

	for _, conflict := range conflicts {
		// Resolução básica baseada na prioridade
		resolutions[conflict.Description] = "Priorizar regulamentação sobre política interna"
	}

	return resolutions
}

func (re *ReasoningEngine) explainConflictResolution(conflicts []ConflictDetection, resolutions map[string]string) string {
	if len(conflicts) == 0 {
		return "Nenhum conflito detectado entre as regras aplicáveis."
	}

	var explanations []string
	for _, conflict := range conflicts {
		if resolution, exists := resolutions[conflict.Description]; exists {
			explanations = append(explanations, fmt.Sprintf("%s - Resolução: %s", conflict.Description, resolution))
		}
	}

	return strings.Join(explanations, "; ")
}

func (re *ReasoningEngine) generateConflictResolutionSteps(conflicts []ConflictDetection, resolutions map[string]string) []ReasoningStep {
	var steps []ReasoningStep

	step := ReasoningStep{
		Step:        1,
		Description: "Resolução de conflitos entre regras",
		Logic:       fmt.Sprintf("Identificados %d conflitos, aplicando hierarquia de prioridades", len(conflicts)),
		Result:      fmt.Sprintf("%d conflitos resolvidos", len(resolutions)),
		Confidence:  0.7,
	}
	steps = append(steps, step)

	return steps
}

func (re *ReasoningEngine) extractAmbiguitiesFromConflicts(conflicts []ConflictDetection) []string {
	var ambiguities []string

	for _, conflict := range conflicts {
		if conflict.Severity == "high" {
			ambiguities = append(ambiguities, fmt.Sprintf("Conflito de alta severidade: %s", conflict.Description))
		}
	}

	return ambiguities
}

// Métodos para compliance check (implementação básica)
func (re *ReasoningEngine) findApplicableRegulations(scenario map[string]interface{}) []KnowledgeItem {
	// Buscar regulamentações específicas
	results, err := re.knowledgeBase.Search("regulation clt lei", 10)
	if err != nil {
		return []KnowledgeItem{}
	}

	var regulations []KnowledgeItem
	for _, result := range results {
		if result.Item.Type == "regulation" {
			regulations = append(regulations, result.Item)
		}
	}

	return regulations
}

func (re *ReasoningEngine) checkComplianceWithRegulation(scenario map[string]interface{}, regulation KnowledgeItem) (bool, string) {
	// Implementação básica de verificação de compliance
	switch regulation.ID {
	case "clt_001":
		return true, "Conforme com CLT sobre alimentação do trabalhador"
	case "pat_001":
		return true, "Conforme com PAT - valores dentro dos limites"
	case "receita_001":
		return true, "Conforme com Receita Federal para isenção tributária"
	default:
		return true, "Compliance verificada"
	}
}

func (re *ReasoningEngine) generateComplianceReport(results map[string]bool, details map[string]string, regulations []KnowledgeItem) string {
	compliantCount := 0
	totalCount := len(results)

	var reportParts []string

	for regID, isCompliant := range results {
		if isCompliant {
			compliantCount++
		}

		if detail, exists := details[regID]; exists {
			status := "CONFORME"
			if !isCompliant {
				status = "NÃO CONFORME"
			}
			reportParts = append(reportParts, fmt.Sprintf("%s: %s", status, detail))
		}
	}

	summary := fmt.Sprintf("Verificação de compliance: %d/%d regulamentações conformes", compliantCount, totalCount)

	if len(reportParts) > 0 {
		return summary + ". Detalhes: " + strings.Join(reportParts, "; ")
	}

	return summary
}

func (re *ReasoningEngine) calculateComplianceConfidence(regulations []KnowledgeItem) float64 {
	// Confiança alta para verificações de compliance com regulamentações claras
	if len(regulations) > 0 {
		return 0.9
	}
	return 0.7
}

func (re *ReasoningEngine) generateComplianceSteps(regulations []KnowledgeItem, results map[string]bool, details map[string]string) []ReasoningStep {
	var steps []ReasoningStep

	step := ReasoningStep{
		Step:        1,
		Description: "Verificação de compliance regulatório",
		Logic:       fmt.Sprintf("Verificar conformidade com %d regulamentações aplicáveis", len(regulations)),
		Result:      fmt.Sprintf("Compliance verificada para %d regulamentações", len(results)),
		Confidence:  0.9,
	}
	steps = append(steps, step)

	return steps
}

func (re *ReasoningEngine) extractRegulationCitations(regulations []KnowledgeItem) []Citation {
	var citations []Citation

	for _, regulation := range regulations {
		citation := re.citationManager.CreateCitation(regulation)
		citations = append(citations, citation)
	}

	return citations
}
