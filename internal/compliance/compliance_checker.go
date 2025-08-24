package compliance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ComplianceChecker é o verificador principal de compliance
type ComplianceChecker struct {
	rules           []ComplianceRule
	rulesByID       map[string]ComplianceRule
	rulesByCategory map[string][]ComplianceRule
	loaded          bool
}

// NewComplianceChecker cria uma nova instância do verificador de compliance
func NewComplianceChecker(rulesDir string) *ComplianceChecker {
	checker := &ComplianceChecker{
		rules:           []ComplianceRule{},
		rulesByID:       make(map[string]ComplianceRule),
		rulesByCategory: make(map[string][]ComplianceRule),
		loaded:          false,
	}

	// Tentar carregar regras automaticamente
	checker.LoadRules(rulesDir)

	return checker
}

// LoadRules carrega as regras de compliance dos arquivos JSON
func (cc *ComplianceChecker) LoadRules(rulesDir string) error {
	// Arquivos de regras a serem carregados
	ruleFiles := []string{
		"clt_rules.json",
		"mte_norms.json",
		"internal_policies.json",
	}

	var allRules []ComplianceRule

	for _, filename := range ruleFiles {
		filePath := filepath.Join(rulesDir, filename)
		rules, err := cc.loadRulesFromFile(filePath)
		if err != nil {
			return fmt.Errorf("erro ao carregar %s: %w", filename, err)
		}
		allRules = append(allRules, rules...)
	}

	// Indexar as regras
	cc.rules = allRules
	cc.indexRules()
	cc.loaded = true

	return nil
}

// loadRulesFromFile carrega regras de um arquivo JSON específico
func (cc *ComplianceChecker) loadRulesFromFile(filePath string) ([]ComplianceRule, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		ComplianceRules []ComplianceRule `json:"compliance_rules"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	return wrapper.ComplianceRules, nil
}

// indexRules cria índices para acesso rápido às regras
func (cc *ComplianceChecker) indexRules() {
	cc.rulesByID = make(map[string]ComplianceRule)
	cc.rulesByCategory = make(map[string][]ComplianceRule)

	for _, rule := range cc.rules {
		cc.rulesByID[rule.ID] = rule
		cc.rulesByCategory[rule.Category] = append(cc.rulesByCategory[rule.Category], rule)
	}
}

// CheckCompliance executa verificação de compliance para um colaborador
func (cc *ComplianceChecker) CheckCompliance(colaboradorData map[string]interface{}) (*ComplianceResult, error) {
	if !cc.loaded {
		return nil, fmt.Errorf("regras de compliance não foram carregadas")
	}

	result := &ComplianceResult{
		ID:         uuid.New().String(),
		EntityID:   cc.extractEntityID(colaboradorData),
		EntityType: "EMPLOYEE",
		CheckedAt:  time.Now(),
		CheckedBy:  "compliance_checker",
		Status:     StatusCompliant,
		Score:      100.0,
		Violations: []ComplianceViolation{},
		Risks:      []Risk{},
		Evidence:   []Evidence{},
		Metadata:   make(map[string]interface{}),
	}

	// Verificar cada regra aplicável
	applicableRules := cc.getApplicableRules(colaboradorData)
	result.RulesChecked = make([]string, len(applicableRules))

	for i, rule := range applicableRules {
		result.RulesChecked[i] = rule.ID

		violation := cc.checkRule(rule, colaboradorData)
		if violation != nil {
			result.Violations = append(result.Violations, *violation)
			result.RulesFailed = append(result.RulesFailed, rule.ID)
		} else {
			result.RulesPassed = append(result.RulesPassed, rule.ID)
		}
	}

	// Calcular status e score final
	cc.calculateFinalStatus(result)
	cc.generateRecommendations(result)
	result.NextReview = cc.calculateNextReview(result)

	return result, nil
}

// getApplicableRules retorna as regras aplicáveis para o contexto
func (cc *ComplianceChecker) getApplicableRules(data map[string]interface{}) []ComplianceRule {
	var applicable []ComplianceRule

	// Para VR, aplicar regras CLT, MTE e internas
	categories := []string{"CLT", "MTE", "INTERNAL"}

	for _, category := range categories {
		if rules, exists := cc.rulesByCategory[category]; exists {
			for _, rule := range rules {
				if cc.isRuleApplicable(rule, data) {
					applicable = append(applicable, rule)
				}
			}
		}
	}

	return applicable
}

// isRuleApplicable verifica se uma regra é aplicável ao contexto
func (cc *ComplianceChecker) isRuleApplicable(rule ComplianceRule, data map[string]interface{}) bool {
	// Verificar se a regra está ativa (não expirada)
	if rule.ExpiryDate != nil && time.Now().After(*rule.ExpiryDate) {
		return false
	}

	// Todas as regras de VR são aplicáveis por padrão
	return true
}

// checkRule verifica uma regra específica
// Violation is an alias for ComplianceViolation for backward compatibility
type Violation = ComplianceViolation

func (cc *ComplianceChecker) checkRule(rule ComplianceRule, data map[string]interface{}) *ComplianceViolation {
	switch rule.ID {
	case "clt_vr_001":
		return cc.checkCLTValueLimit(rule, data)
	case "clt_vr_002":
		return cc.checkCLTDiscountLimit(rule, data)
	case "clt_vr_003":
		return cc.checkCLTEligibility(rule, data)
	case "clt_vr_004":
		return cc.checkCLTDocumentation(rule, data)
	case "clt_vr_005":
		return cc.checkCLTSalaryCharacterization(rule, data)
	case "mte_pat_001":
		return cc.checkMTEPATRegistration(rule, data)
	case "mte_pat_002":
		return cc.checkMTEBeneficiaryControl(rule, data)
	case "mte_pat_003":
		return cc.checkMTEAccountability(rule, data)
	case "mte_seg_001":
		return cc.checkMTEHealthSafety(rule, data)
	case "mte_trab_001":
		return cc.checkMTENonDiscrimination(rule, data)
	default:
		return cc.checkInternalPolicy(rule, data)
	}
}

// Verificações específicas CLT
func (cc *ComplianceChecker) checkCLTValueLimit(rule ComplianceRule, data map[string]interface{}) *ComplianceViolation {
	valorVR := cc.getFloatValue(data, "valor_vr")
	valorMaximo := rule.Parameters["valor_maximo_diario"].(float64)

	if valorVR > valorMaximo {
		return &ComplianceViolation{
			ID:          uuid.New().String(),
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Description: fmt.Sprintf("Valor do VR (R$ %.2f) excede o limite máximo PAT (R$ %.2f)", valorVR, valorMaximo),
			Severity:    SeverityHigh,
			Category:    rule.Category,
			EntityID:    cc.extractEntityID(data),
			Details: map[string]interface{}{
				"valor_praticado": valorVR,
				"valor_maximo":    valorMaximo,
				"excesso":         valorVR - valorMaximo,
			},
			Impact:      "Perda de incentivos fiscais, caracterização como salário",
			Remediation: "Ajustar valor do VR para dentro do limite PAT ou justificar exceção",
			DetectedAt:  time.Now(),
			Resolved:    false,
		}
	}
	return nil
}

func (cc *ComplianceChecker) checkCLTDiscountLimit(rule ComplianceRule, data map[string]interface{}) *ComplianceViolation {
	percentualDesconto := cc.getFloatValue(data, "percentual_desconto")
	percentualMaximo := rule.Parameters["percentual_maximo_desconto"].(float64)

	if percentualDesconto > percentualMaximo {
		return &ComplianceViolation{
			ID:          uuid.New().String(),
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Description: fmt.Sprintf("Desconto em folha (%.1f%%) excede limite legal (%.1f%%)", percentualDesconto, percentualMaximo),
			Severity:    SeverityHigh,
			Category:    rule.Category,
			EntityID:    cc.extractEntityID(data),
			Details: map[string]interface{}{
				"percentual_praticado": percentualDesconto,
				"percentual_maximo":    percentualMaximo,
				"excesso":              percentualDesconto - percentualMaximo,
			},
			Impact:      "Desconto indevido, possível ação trabalhista",
			Remediation: "Ajustar percentual de desconto para máximo de 20%",
			DetectedAt:  time.Now(),
			Resolved:    false,
		}
	}
	return nil
}

func (cc *ComplianceChecker) checkCLTEligibility(rule ComplianceRule, data map[string]interface{}) *ComplianceViolation {
	cargaHoraria := cc.getFloatValue(data, "carga_horaria")
	cargaMinima := rule.Parameters["jornada_minima_horas"].(float64)
	tipoColaborador := cc.getStringValue(data, "tipo_colaborador")

	// Verificar exclusões
	excecoes := rule.Parameters["excecoes"].([]interface{})
	for _, excecao := range excecoes {
		if strings.Contains(strings.ToLower(tipoColaborador), strings.ToLower(excecao.(string))) {
			return &ComplianceViolation{
				ID:          uuid.New().String(),
				RuleID:      rule.ID,
				RuleName:    rule.Name,
				Description: fmt.Sprintf("Tipo de colaborador '%s' não elegível para VR", tipoColaborador),
				Severity:    SeverityMedium,
				Category:    rule.Category,
				Details: map[string]interface{}{
					"tipo_colaborador": tipoColaborador,
					"motivo_exclusao":  "Categoria não elegível conforme CLT",
				},
				Impact:      "Concessão indevida de benefício",
				Remediation: "Remover benefício ou reclassificar colaborador",
				DetectedAt:  time.Now(),
				Resolved:    false,
			}
		}
	}

	// Verificar carga horária
	if cargaHoraria < cargaMinima {
		return &ComplianceViolation{
			ID:          uuid.New().String(),
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Description: fmt.Sprintf("Carga horária (%.1fh) inferior ao mínimo (%.1fh)", cargaHoraria, cargaMinima),
			Severity:    SeverityMedium,
			Category:    rule.Category,
			EntityID:    cc.extractEntityID(data),
			Details: map[string]interface{}{
				"carga_horaria_atual":  cargaHoraria,
				"carga_horaria_minima": cargaMinima,
			},
			Impact:      "Concessão indevida de benefício",
			Remediation: "Verificar jornada real ou remover benefício",
			DetectedAt:  time.Now(),
			Resolved:    false,
		}
	}

	return nil
}

func (cc *ComplianceChecker) checkCLTDocumentation(rule ComplianceRule, data map[string]interface{}) *ComplianceViolation {
	temTermoAdesao := cc.getBoolValue(data, "tem_termo_adesao")
	temComprovantesEntrega := cc.getBoolValue(data, "tem_comprovantes_entrega")

	var problemas []string
	if !temTermoAdesao {
		problemas = append(problemas, "Termo de adesão não localizado")
	}
	if !temComprovantesEntrega {
		problemas = append(problemas, "Comprovantes de entrega ausentes")
	}

	if len(problemas) > 0 {
		return &ComplianceViolation{
			ID:          uuid.New().String(),
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Description: fmt.Sprintf("Documentação incompleta: %s", strings.Join(problemas, ", ")),
			Severity:    SeverityMedium,
			Category:    rule.Category,
			EntityID:    cc.extractEntityID(data),
			Details: map[string]interface{}{
				"problemas_encontrados": problemas,
			},
			Impact:      "Dificuldade em defesa fiscal, presunção de irregularidade",
			Remediation: "Providenciar documentação faltante",
			DetectedAt:  time.Now(),
			Resolved:    false,
		}
	}

	return nil
}

func (cc *ComplianceChecker) checkCLTSalaryCharacterization(rule ComplianceRule, data map[string]interface{}) *ComplianceViolation {
	usoExclusivoAlimentacao := cc.getBoolValue(data, "uso_exclusivo_alimentacao")
	controleUtilizacao := cc.getBoolValue(data, "controle_utilizacao")

	var problemas []string
	if !usoExclusivoAlimentacao {
		problemas = append(problemas, "Uso não exclusivo para alimentação")
	}
	if !controleUtilizacao {
		problemas = append(problemas, "Controle de utilização inadequado")
	}

	if len(problemas) > 0 {
		return &ComplianceViolation{
			ID:          uuid.New().String(),
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Description: fmt.Sprintf("Risco de caracterização como salário: %s", strings.Join(problemas, ", ")),
			Severity:    SeverityCritical,
			Category:    rule.Category,
			EntityID:    cc.extractEntityID(data),
			Details: map[string]interface{}{
				"problemas_encontrados": problemas,
			},
			Impact:      "Caracterização como salário, incidência de encargos sociais",
			Remediation: "Implementar controles rigorosos de uso exclusivo para alimentação",
			DetectedAt:  time.Now(),
			Resolved:    false,
		}
	}

	return nil
}

// Verificações específicas MTE
func (cc *ComplianceChecker) checkMTEPATRegistration(rule ComplianceRule, data map[string]interface{}) *ComplianceViolation {
	inscricaoPAT := cc.getBoolValue(data, "inscricao_pat")
	renovacaoAtualizada := cc.getBoolValue(data, "renovacao_pat_atualizada")

	var problemas []string
	if !inscricaoPAT {
		problemas = append(problemas, "Não inscrito no PAT")
	}
	if !renovacaoAtualizada {
		problemas = append(problemas, "Renovação PAT em atraso")
	}

	if len(problemas) > 0 {
		return &ComplianceViolation{
			ID:          uuid.New().String(),
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Description: fmt.Sprintf("Problemas com inscrição PAT: %s", strings.Join(problemas, ", ")),
			Severity:    SeverityHigh,
			Category:    rule.Category,
			EntityID:    cc.extractEntityID(data),
			Details: map[string]interface{}{
				"problemas_encontrados": problemas,
			},
			Impact:      "Perda de incentivos fiscais, multa MTE",
			Remediation: "Regularizar situação junto ao PAT/MTE",
			DetectedAt:  time.Now(),
			Resolved:    false,
		}
	}

	return nil
}

func (cc *ComplianceChecker) checkMTEBeneficiaryControl(rule ComplianceRule, data map[string]interface{}) *ComplianceViolation {
	controleAtualizado := cc.getBoolValue(data, "controle_beneficiarios_atualizado")

	if !controleAtualizado {
		return &ComplianceViolation{
			ID:          uuid.New().String(),
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Description: "Controle de beneficiários PAT desatualizado",
			Severity:    SeverityMedium,
			Category:    rule.Category,
			EntityID:    cc.extractEntityID(data),
			Details: map[string]interface{}{
				"ultima_atualizacao": cc.getStringValue(data, "ultima_atualizacao_beneficiarios"),
			},
			Impact:      "Multa por desatualização, suspensão de benefícios fiscais",
			Remediation: "Atualizar controle de beneficiários no sistema PAT",
			DetectedAt:  time.Now(),
			Resolved:    false,
		}
	}

	return nil
}

func (cc *ComplianceChecker) checkMTEAccountability(rule ComplianceRule, data map[string]interface{}) *ComplianceViolation {
	prestacaoContasEmDia := cc.getBoolValue(data, "prestacao_contas_pat_em_dia")

	if !prestacaoContasEmDia {
		return &ComplianceViolation{
			ID:          uuid.New().String(),
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Description: "Prestação de contas PAT em atraso",
			Severity:    SeverityHigh,
			Category:    rule.Category,
			EntityID:    cc.extractEntityID(data),
			Details: map[string]interface{}{
				"prazo_vencimento": "31 de março",
				"status":           "em_atraso",
			},
			Impact:      "Cancelamento da inscrição PAT, devolução de incentivos",
			Remediation: "Enviar prestação de contas urgentemente",
			DetectedAt:  time.Now(),
			Resolved:    false,
		}
	}

	return nil
}

func (cc *ComplianceChecker) checkMTEHealthSafety(rule ComplianceRule, data map[string]interface{}) *ComplianceViolation {
	estabelecimentosCredenciados := cc.getBoolValue(data, "estabelecimentos_credenciados")
	qualidadeNutricional := cc.getBoolValue(data, "qualidade_nutricional_ok")

	var problemas []string
	if !estabelecimentosCredenciados {
		problemas = append(problemas, "Estabelecimentos não adequadamente credenciados")
	}
	if !qualidadeNutricional {
		problemas = append(problemas, "Qualidade nutricional não verificada")
	}

	if len(problemas) > 0 {
		return &ComplianceViolation{
			ID:          uuid.New().String(),
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Description: fmt.Sprintf("Problemas de segurança e saúde: %s", strings.Join(problemas, ", ")),
			Severity:    SeverityMedium,
			Category:    rule.Category,
			EntityID:    cc.extractEntityID(data),
			Details: map[string]interface{}{
				"problemas_encontrados": problemas,
			},
			Impact:      "Risco à saúde dos empregados, responsabilização",
			Remediation: "Melhorar controles de qualidade dos fornecedores",
			DetectedAt:  time.Now(),
			Resolved:    false,
		}
	}

	return nil
}

func (cc *ComplianceChecker) checkMTENonDiscrimination(rule ComplianceRule, data map[string]interface{}) *ComplianceViolation {
	criteriosObjetivos := cc.getBoolValue(data, "criterios_objetivos")
	aplicacaoUniforme := cc.getBoolValue(data, "aplicacao_uniforme")

	var problemas []string
	if !criteriosObjetivos {
		problemas = append(problemas, "Critérios não objetivos para concessão")
	}
	if !aplicacaoUniforme {
		problemas = append(problemas, "Aplicação não uniforme do benefício")
	}

	if len(problemas) > 0 {
		return &ComplianceViolation{
			ID:          uuid.New().String(),
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Description: fmt.Sprintf("Possível discriminação: %s", strings.Join(problemas, ", ")),
			Severity:    SeverityHigh,
			Category:    rule.Category,
			EntityID:    cc.extractEntityID(data),
			Details: map[string]interface{}{
				"problemas_encontrados": problemas,
			},
			Impact:      "Ação por discriminação, danos morais",
			Remediation: "Revisar e padronizar critérios de concessão",
			DetectedAt:  time.Now(),
			Resolved:    false,
		}
	}

	return nil
}

// checkInternalPolicy verifica políticas internas
func (cc *ComplianceChecker) checkInternalPolicy(rule ComplianceRule, data map[string]interface{}) *ComplianceViolation {
	// Implementação genérica para políticas internas
	complianceOK := cc.getBoolValue(data, fmt.Sprintf("compliance_%s", strings.ToLower(rule.ID)))

	if !complianceOK {
		return &ComplianceViolation{
			ID:          uuid.New().String(),
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Description: fmt.Sprintf("Não conformidade com política interna: %s", rule.Name),
			Severity:    SeverityLow,
			Category:    rule.Category,
			EntityID:    cc.extractEntityID(data),
			Details: map[string]interface{}{
				"politica": rule.Name,
			},
			Impact:      "Não conformidade com procedimentos internos",
			Remediation: "Seguir procedimentos estabelecidos na política",
			DetectedAt:  time.Now(),
			Resolved:    false,
		}
	}

	return nil
}

// Métodos auxiliares
func (cc *ComplianceChecker) extractEntityID(data map[string]interface{}) string {
	if id, exists := data["colaborador_id"]; exists {
		if idStr, ok := id.(string); ok {
			return idStr
		}
	}
	if matricula, exists := data["matricula"]; exists {
		if matStr, ok := matricula.(string); ok {
			return matStr
		}
	}
	return "unknown"
}

func (cc *ComplianceChecker) getFloatValue(data map[string]interface{}, key string) float64 {
	if value, exists := data[key]; exists {
		if fVal, ok := value.(float64); ok {
			return fVal
		}
		if iVal, ok := value.(int); ok {
			return float64(iVal)
		}
	}
	return 0.0
}

func (cc *ComplianceChecker) getStringValue(data map[string]interface{}, key string) string {
	if value, exists := data[key]; exists {
		if sVal, ok := value.(string); ok {
			return sVal
		}
	}
	return ""
}

func (cc *ComplianceChecker) getBoolValue(data map[string]interface{}, key string) bool {
	if value, exists := data[key]; exists {
		if bVal, ok := value.(bool); ok {
			return bVal
		}
	}
	return true // Default para verdadeiro para evitar falsos positivos
}

func (cc *ComplianceChecker) calculateFinalStatus(result *ComplianceResult) {
	if len(result.Violations) == 0 {
		result.Status = StatusCompliant
		result.Score = 100.0
		return
	}

	// Calcular score baseado nas violações
	totalDeductions := 0.0
	hasCritical := false
	hasHigh := false

	for _, violation := range result.Violations {
		switch violation.Severity {
		case SeverityCritical:
			totalDeductions += 30.0
			hasCritical = true
		case SeverityHigh:
			totalDeductions += 20.0
			hasHigh = true
		case SeverityMedium:
			totalDeductions += 10.0
		case SeverityLow:
			totalDeductions += 5.0
		}
	}

	result.Score = 100.0 - totalDeductions
	if result.Score < 0 {
		result.Score = 0
	}

	// Determinar status
	if hasCritical {
		result.Status = StatusNonCompliant
	} else if hasHigh || result.Score < 70 {
		result.Status = StatusWarning
	} else {
		result.Status = StatusCompliant
	}
}

func (cc *ComplianceChecker) generateRecommendations(result *ComplianceResult) {
	if len(result.Violations) == 0 {
		result.Recommendations = []string{"Manter práticas atuais de compliance"}
		return
	}

	recommendations := make(map[string]bool)

	for _, violation := range result.Violations {
		if violation.Remediation != "" && !recommendations[violation.Remediation] {
			result.Recommendations = append(result.Recommendations, violation.Remediation)
			recommendations[violation.Remediation] = true
		}
	}

	// Adicionar recomendações gerais
	if result.Status == StatusNonCompliant {
		result.Recommendations = append(result.Recommendations, "Solicitar revisão urgente por especialista em compliance")
	}
}

func (cc *ComplianceChecker) calculateNextReview(result *ComplianceResult) time.Time {
	baseReview := time.Now().AddDate(0, 1, 0) // 1 mês por padrão

	if result.Status == StatusNonCompliant {
		return time.Now().AddDate(0, 0, 7) // 1 semana para não conformidades críticas
	}

	if result.Status == StatusWarning {
		return time.Now().AddDate(0, 0, 15) // 2 semanas para warnings
	}

	return baseReview
}

// GetRuleByID retorna uma regra específica por ID
func (cc *ComplianceChecker) GetRuleByID(ruleID string) (ComplianceRule, error) {
	if !cc.loaded {
		return ComplianceRule{}, fmt.Errorf("regras não foram carregadas")
	}

	if rule, exists := cc.rulesByID[ruleID]; exists {
		return rule, nil
	}

	return ComplianceRule{}, fmt.Errorf("regra %s não encontrada", ruleID)
}

// GetRulesByCategory retorna regras de uma categoria específica
func (cc *ComplianceChecker) GetRulesByCategory(category string) ([]ComplianceRule, error) {
	if !cc.loaded {
		return nil, fmt.Errorf("regras não foram carregadas")
	}

	if rules, exists := cc.rulesByCategory[category]; exists {
		return rules, nil
	}

	return []ComplianceRule{}, nil
}

// GetStats retorna estatísticas das regras carregadas
func (cc *ComplianceChecker) GetStats() map[string]interface{} {
	stats := map[string]interface{}{
		"loaded":      cc.loaded,
		"total_rules": len(cc.rules),
	}

	if cc.loaded {
		categoryCounts := make(map[string]int)
		typeCounts := make(map[string]int)

		for _, rule := range cc.rules {
			categoryCounts[rule.Category]++
			typeCounts[rule.Type]++
		}

		stats["rules_by_category"] = categoryCounts
		stats["rules_by_type"] = typeCounts
	}

	return stats
}

// GetAllRules retorna todas as regras carregadas
func (cc *ComplianceChecker) GetAllRules() []ComplianceRule {
	return cc.rules
}

// LoadRegulations é um alias para LoadRules para compatibilidade
func (cc *ComplianceChecker) LoadRegulations() error {
	if cc.loaded {
		return nil
	}
	return fmt.Errorf("regras não foram carregadas durante a inicialização")
}
