package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"BrxAgente-desafio4/internal/compliance"
	"github.com/tmc/langchaingo/tools"
)

type ComplianceAssistant struct {
	complianceChecker    *compliance.ComplianceChecker
	riskAssessor        *compliance.RiskAssessor
	auditTrail          *compliance.AuditTrail
	documentationGen    *compliance.DocumentationGenerator
	name                string
	description         string
}

func NewComplianceAssistant() *ComplianceAssistant {
	auditTrail := compliance.NewAuditTrail(100000, 365) // 100k entries capacity, 365 days retention
	complianceChecker := compliance.NewComplianceChecker("./internal/data/regulations")
	riskAssessor := compliance.NewRiskAssessor(complianceChecker, auditTrail)
	documentationGen := compliance.NewDocumentationGenerator(complianceChecker, riskAssessor, auditTrail)

	return &ComplianceAssistant{
		complianceChecker: complianceChecker,
		riskAssessor:     riskAssessor,
		auditTrail:       auditTrail,
		documentationGen: documentationGen,
		name:             "compliance_assistant",
		description: `Assistente de compliance e auditoria para Vale Refeição. 

Funcionalidades:
- check_compliance: Verifica compliance de funcionários
- assess_risk: Avalia riscos de compliance
- generate_report: Gera relatórios de compliance
- audit_trail: Registra e consulta trilha de auditoria
- get_regulations: Consulta regulamentações específicas

Formato de entrada (JSON):
{
  "action": "check_compliance|assess_risk|generate_report|audit_trail|get_regulations",
  "data": {
    // dados específicos para cada ação
  }
}

Exemplos:
1. Verificar compliance:
   {"action": "check_compliance", "data": {"matricula": "12345", "valor_vr": 45.0}}

2. Avaliar risco:
   {"action": "assess_risk", "data": {"matricula": "12345", "nome": "João", "valor_vr": 45.0}}

3. Gerar relatório:
   {"action": "generate_report", "data": {"employees": [...], "start_date": "2024-01-01", "end_date": "2024-01-31"}}`,
	}
}

func (c *ComplianceAssistant) Name() string {
	return c.name
}

func (c *ComplianceAssistant) Description() string {
	return c.description
}

func (c *ComplianceAssistant) Call(ctx context.Context, input string) (string, error) {
	return c.Execute(input)
}

func (c *ComplianceAssistant) Execute(input string) (string, error) {
	var request struct {
		Action string                 `json:"action"`
		Data   map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal([]byte(input), &request); err != nil {
		return "", fmt.Errorf("erro ao fazer parse da entrada JSON: %v", err)
	}

	// Record audit entry
	c.auditTrail.RecordSimpleAction("system", request.Action, fmt.Sprintf("Executando ação: %s", request.Action), "INFO", request.Data)

	switch request.Action {
	case "check_compliance":
		return c.handleCheckCompliance(request.Data)
	case "assess_risk":
		return c.handleAssessRisk(request.Data)
	case "generate_report":
		return c.handleGenerateReport(request.Data)
	case "audit_trail":
		return c.handleAuditTrail(request.Data)
	case "get_regulations":
		return c.handleGetRegulations(request.Data)
	default:
		return "", fmt.Errorf("ação não reconhecida: %s", request.Action)
	}
}

func (c *ComplianceAssistant) handleCheckCompliance(data map[string]interface{}) (string, error) {
	result, err := c.complianceChecker.CheckCompliance(data)
	if err != nil {
		return "", fmt.Errorf("erro na verificação de compliance: %v", err)
	}

	response := map[string]interface{}{
		"status":              result.Status,
		"compliance_score":    result.Score,
		"total_violations":    len(result.Violations),
		"violations":          result.Violations,
		"recommendations":     result.Recommendations,
		"checked_at":          result.CheckedAt,
	}

	// Log compliance check result
	level := "INFO"
	if result.Status != compliance.StatusCompliant {
		level = "WARN"
	}
	
	entityID, _ := data["matricula"].(string)
	c.auditTrail.RecordSimpleAction(entityID, "compliance_check", 
		fmt.Sprintf("Verificação de compliance - Status: %s", result.Status), level, response)

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao serializar resposta: %v", err)
	}

	return string(jsonResponse), nil
}

func (c *ComplianceAssistant) handleAssessRisk(data map[string]interface{}) (string, error) {
	assessment, err := c.riskAssessor.AssessEmployeeRisk(data)
	if err != nil {
		return "", fmt.Errorf("erro na avaliação de risco: %v", err)
	}

	response := map[string]interface{}{
		"entity_id":        assessment.EntityID,
		"overall_risk":     assessment.OverallRisk,
		"risk_score":       assessment.RiskScore,
		"risk_factors":     assessment.RiskFactors,
		"recommendations":  assessment.Recommendations,
		"assessment_date":  assessment.AssessmentDate,
		"next_review_date": assessment.NextReviewDate,
	}

	// Log risk assessment
	level := "INFO"
	if assessment.OverallRisk == compliance.RiskHigh || assessment.OverallRisk == compliance.RiskCritical {
		level = "WARN"
	}

	c.auditTrail.RecordSimpleAction(assessment.EntityID, "risk_assessment", 
		fmt.Sprintf("Avaliação de risco - Nível: %s (Score: %.2f)", assessment.OverallRisk, assessment.RiskScore), 
		level, response)

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao serializar resposta: %v", err)
	}

	return string(jsonResponse), nil
}

func (c *ComplianceAssistant) handleGenerateReport(data map[string]interface{}) (string, error) {
	// Parse employees data
	employeesData, ok := data["employees"].([]interface{})
	if !ok {
		return "", fmt.Errorf("dados de funcionários não fornecidos ou inválidos")
	}

	employees := make([]map[string]interface{}, len(employeesData))
	for i, emp := range employeesData {
		empMap, ok := emp.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("formato inválido para dados do funcionário %d", i)
		}
		employees[i] = empMap
	}

	// Parse dates
	startDateStr, _ := data["start_date"].(string)
	endDateStr, _ := data["end_date"].(string)

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		endDate = time.Now() // Default to now
	}

	period := compliance.DateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Generate report
	report, err := c.documentationGen.GenerateComplianceReport(employees, period)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar relatório: %v", err)
	}

	// Check if markdown format is requested
	format, _ := data["format"].(string)
	if format == "markdown" {
		markdownReport := c.documentationGen.GenerateMarkdownReport(report)
		
		response := map[string]interface{}{
			"format":    "markdown",
			"report":    markdownReport,
			"summary":   report.Summary,
			"generated_at": report.GeneratedAt,
		}

		jsonResponse, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return "", fmt.Errorf("erro ao serializar resposta: %v", err)
		}
		return string(jsonResponse), nil
	}

	// Log report generation
	c.auditTrail.RecordSimpleAction("system", "generate_report", 
		fmt.Sprintf("Relatório de compliance gerado - Período: %s a %s", 
			startDate.Format("02/01/2006"), endDate.Format("02/01/2006")), 
		"INFO", map[string]interface{}{
			"total_employees": report.Summary.TotalEmployees,
			"compliance_rate": report.Summary.ComplianceRate,
			"total_violations": len(report.Violations),
		})

	jsonResponse, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao serializar relatório: %v", err)
	}

	return string(jsonResponse), nil
}

func (c *ComplianceAssistant) handleAuditTrail(data map[string]interface{}) (string, error) {
	operation, _ := data["operation"].(string)
	
	switch operation {
	case "query":
		return c.handleAuditQuery(data)
	case "record":
		return c.handleAuditRecord(data)
	case "report":
		return c.handleAuditReport(data)
	default:
		return "", fmt.Errorf("operação de auditoria não reconhecida: %s", operation)
	}
}

func (c *ComplianceAssistant) handleAuditQuery(data map[string]interface{}) (string, error) {
	entityID, _ := data["entity_id"].(string)
	
	// Parse date range
	sinceStr, _ := data["since"].(string)
	untilStr, _ := data["until"].(string)
	
	since := time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	until := time.Now()
	
	if sinceStr != "" {
		if parsed, err := time.Parse("2006-01-02", sinceStr); err == nil {
			since = parsed
		}
	}
	
	if untilStr != "" {
		if parsed, err := time.Parse("2006-01-02", untilStr); err == nil {
			until = parsed
		}
	}

	entries := c.auditTrail.GetEntriesForEntity(entityID, since, until)
	
	response := map[string]interface{}{
		"total_entries": len(entries),
		"period": map[string]string{
			"since": since.Format("2006-01-02"),
			"until": until.Format("2006-01-02"),
		},
		"entries": entries,
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao serializar resposta: %v", err)
	}

	return string(jsonResponse), nil
}

func (c *ComplianceAssistant) handleAuditRecord(data map[string]interface{}) (string, error) {
	userID, _ := data["user_id"].(string)
	action, _ := data["action"].(string)
	description, _ := data["description"].(string)
	level, _ := data["level"].(string)
	
	if userID == "" || action == "" {
		return "", fmt.Errorf("user_id e action são obrigatórios para registrar entrada de auditoria")
	}
	
	if level == "" {
		level = "INFO"
	}

	metadata, _ := data["metadata"].(map[string]interface{})
	
	c.auditTrail.RecordSimpleAction(userID, action, description, level, metadata)
	
	response := map[string]interface{}{
		"status": "recorded",
		"timestamp": time.Now(),
		"user_id": userID,
		"action": action,
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao serializar resposta: %v", err)
	}

	return string(jsonResponse), nil
}

func (c *ComplianceAssistant) handleAuditReport(data map[string]interface{}) (string, error) {
	startDateStr, _ := data["start_date"].(string)
	endDateStr, _ := data["end_date"].(string)
	
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		endDate = time.Now() // Default to now
	}

	auditReport := c.documentationGen.GenerateAuditReport(startDate, endDate)
	
	jsonResponse, err := json.MarshalIndent(auditReport, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao serializar relatório de auditoria: %v", err)
	}

	return string(jsonResponse), nil
}

func (c *ComplianceAssistant) handleGetRegulations(data map[string]interface{}) (string, error) {
	category, _ := data["category"].(string)
	ruleID, _ := data["rule_id"].(string)
	ruleType, _ := data["type"].(string)
	
	// Load regulations if needed
	err := c.complianceChecker.LoadRegulations()
	if err != nil {
		return "", fmt.Errorf("erro ao carregar regulamentações: %v", err)
	}

	var filteredRules []compliance.ComplianceRule
	
	// Get all rules and filter based on criteria
	allRules := c.complianceChecker.GetAllRules()
	
	for _, rule := range allRules {
		matches := true
		
		if category != "" && rule.Category != category {
			matches = false
		}
		
		if ruleID != "" && rule.ID != ruleID {
			matches = false
		}
		
		if ruleType != "" && rule.Type != ruleType {
			matches = false
		}
		
		if matches {
			filteredRules = append(filteredRules, rule)
		}
	}

	response := map[string]interface{}{
		"total_rules": len(filteredRules),
		"filters": map[string]interface{}{
			"category": category,
			"rule_id":  ruleID,
			"type":     ruleType,
		},
		"rules": filteredRules,
	}

	// Log regulation query
	c.auditTrail.RecordSimpleAction("system", "query_regulations", 
		fmt.Sprintf("Consulta de regulamentações - Filtros: category=%s, rule_id=%s, type=%s", 
			category, ruleID, ruleType), "INFO", response)

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao serializar resposta: %v", err)
	}

	return string(jsonResponse), nil
}

// Ensure ComplianceAssistant implements the tools.Tool interface
var _ tools.Tool = (*ComplianceAssistant)(nil)