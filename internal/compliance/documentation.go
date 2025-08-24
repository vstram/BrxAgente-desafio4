package compliance

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type DocumentationGenerator struct {
	complianceChecker *ComplianceChecker
	riskAssessor      *RiskAssessor
	auditTrail        *AuditTrail
}

func NewDocumentationGenerator(checker *ComplianceChecker, assessor *RiskAssessor, audit *AuditTrail) *DocumentationGenerator {
	return &DocumentationGenerator{
		complianceChecker: checker,
		riskAssessor:      assessor,
		auditTrail:        audit,
	}
}

type ComplianceReport struct {
	GeneratedAt     time.Time             `json:"generated_at"`
	Period          DateRange             `json:"period"`
	Summary         ComplianceSummary     `json:"summary"`
	Violations      []ComplianceViolation `json:"violations"`
	RiskAssessments []RiskAssessment      `json:"risk_assessments"`
	Recommendations []string              `json:"recommendations"`
	AuditTrail      []AuditEntry          `json:"audit_trail"`
	Metadata        ReportMetadata        `json:"metadata"`
}

type ComplianceSummary struct {
	TotalEmployees       int                       `json:"total_employees"`
	ComplianceRate       float64                   `json:"compliance_rate"`
	ViolationsBySeverity map[ViolationSeverity]int `json:"violations_by_severity"`
	RisksByLevel         map[RiskLevel]int         `json:"risks_by_level"`
	TopRiskCategories    []string                  `json:"top_risk_categories"`
	ImprovementTrends    []TrendData               `json:"improvement_trends"`
}

type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type ReportMetadata struct {
	Version        string             `json:"version"`
	GeneratedBy    string             `json:"generated_by"`
	ReportType     string             `json:"report_type"`
	Scope          string             `json:"scope"`
	DataSources    []string           `json:"data_sources"`
	QualityMetrics map[string]float64 `json:"quality_metrics"`
}

type TrendData struct {
	Period string  `json:"period"`
	Value  float64 `json:"value"`
	Metric string  `json:"metric"`
}

func (d *DocumentationGenerator) GenerateComplianceReport(employeeData []map[string]interface{}, period DateRange) (*ComplianceReport, error) {
	report := &ComplianceReport{
		GeneratedAt: time.Now(),
		Period:      period,
		Metadata: ReportMetadata{
			Version:        "1.0",
			GeneratedBy:    "BrxAgente Compliance Assistant",
			ReportType:     "Compliance Assessment Report",
			Scope:          "Vale Refeição (VR) Compliance",
			DataSources:    []string{"Employee Database", "Audit Trail", "Risk Assessment"},
			QualityMetrics: make(map[string]float64),
		},
	}

	// Process employee data for compliance and risk assessment
	violations := []ComplianceViolation{}
	riskAssessments := []RiskAssessment{}

	for _, employee := range employeeData {
		// Check compliance
		complianceResult, _ := d.complianceChecker.CheckCompliance(employee)
		violations = append(violations, complianceResult.Violations...)

		// Assess risk
		riskAssessment, err := d.riskAssessor.AssessEmployeeRisk(employee)
		if err == nil {
			riskAssessments = append(riskAssessments, *riskAssessment)
		}
	}

	report.Violations = violations
	report.RiskAssessments = riskAssessments

	// Generate summary
	report.Summary = d.generateSummary(employeeData, violations, riskAssessments)

	// Get audit trail for the period
	auditEntries := d.auditTrail.GetEntriesByDateRange(period.StartDate, period.EndDate)
	report.AuditTrail = auditEntries

	// Generate recommendations
	report.Recommendations = d.generateRecommendations(violations, riskAssessments)

	// Calculate quality metrics
	report.Metadata.QualityMetrics = d.calculateQualityMetrics(employeeData, violations)

	return report, nil
}

func (d *DocumentationGenerator) generateSummary(employeeData []map[string]interface{}, violations []ComplianceViolation, riskAssessments []RiskAssessment) ComplianceSummary {
	summary := ComplianceSummary{
		TotalEmployees:       len(employeeData),
		ViolationsBySeverity: make(map[ViolationSeverity]int),
		RisksByLevel:         make(map[RiskLevel]int),
		TopRiskCategories:    []string{},
		ImprovementTrends:    []TrendData{},
	}

	// Count violations by severity
	for _, violation := range violations {
		summary.ViolationsBySeverity[violation.Severity]++
	}

	// Count risks by level
	for _, risk := range riskAssessments {
		summary.RisksByLevel[risk.OverallRisk]++
	}

	// Calculate compliance rate
	employeesWithViolations := make(map[string]bool)
	for _, violation := range violations {
		employeesWithViolations[violation.EntityID] = true
	}

	compliantEmployees := len(employeeData) - len(employeesWithViolations)
	if len(employeeData) > 0 {
		summary.ComplianceRate = (float64(compliantEmployees) / float64(len(employeeData))) * 100
	}

	// Identify top risk categories
	categoryCount := make(map[string]int)
	for _, risk := range riskAssessments {
		for _, factor := range risk.RiskFactors {
			categoryCount[factor.Category]++
		}
	}

	type categoryRisk struct {
		category string
		count    int
	}
	categories := []categoryRisk{}
	for cat, count := range categoryCount {
		categories = append(categories, categoryRisk{cat, count})
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].count > categories[j].count
	})

	for i, cat := range categories {
		if i >= 5 { // Top 5 categories
			break
		}
		summary.TopRiskCategories = append(summary.TopRiskCategories, cat.category)
	}

	return summary
}

func (d *DocumentationGenerator) generateRecommendations(violations []ComplianceViolation, riskAssessments []RiskAssessment) []string {
	recommendations := []string{}

	// Priority recommendations based on violations
	criticalViolations := 0
	highViolations := 0

	for _, violation := range violations {
		if violation.Severity == SeverityCritical {
			criticalViolations++
		} else if violation.Severity == SeverityHigh {
			highViolations++
		}
	}

	if criticalViolations > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🚨 URGENTE: Resolver %d violação(ões) crítica(s) imediatamente", criticalViolations))
		recommendations = append(recommendations,
			"Implementar controles de emergência para violações críticas")
	}

	if highViolations > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("⚠️ ALTA PRIORIDADE: Planejar resolução de %d violação(ões) de alta severidade", highViolations))
	}

	// Risk-based recommendations
	highRiskCount := 0
	criticalRiskCount := 0

	for _, risk := range riskAssessments {
		if risk.OverallRisk == RiskHigh {
			highRiskCount++
		} else if risk.OverallRisk == RiskCritical {
			criticalRiskCount++
		}
	}

	if criticalRiskCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🔴 Revisar %d funcionário(s) com risco crítico", criticalRiskCount))
	}

	if highRiskCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🟡 Monitorar %d funcionário(s) com risco alto", highRiskCount))
	}

	// General recommendations
	recommendations = append(recommendations, "📊 Implementar monitoramento contínuo de compliance")
	recommendations = append(recommendations, "📚 Realizar treinamento de compliance para equipe")
	recommendations = append(recommendations, "🔍 Estabelecer auditorias regulares mensais")

	if len(violations) > 0 {
		recommendations = append(recommendations, "📝 Documentar procedimentos de correção de não conformidades")
	}

	return recommendations
}

func (d *DocumentationGenerator) calculateQualityMetrics(employeeData []map[string]interface{}, violations []ComplianceViolation) map[string]float64 {
	metrics := make(map[string]float64)

	// Data completeness
	requiredFields := []string{"matricula", "nome", "data_admissao", "sindicato", "valor_vr"}
	totalFields := len(employeeData) * len(requiredFields)
	completeFields := 0

	for _, employee := range employeeData {
		for _, field := range requiredFields {
			if _, ok := employee[field]; ok {
				completeFields++
			}
		}
	}

	if totalFields > 0 {
		metrics["data_completeness"] = (float64(completeFields) / float64(totalFields)) * 100
	}

	// Accuracy (inverse of violation rate)
	if len(employeeData) > 0 {
		violationRate := (float64(len(violations)) / float64(len(employeeData))) * 100
		metrics["data_accuracy"] = 100 - violationRate
	}

	// Timeliness (based on audit trail activity)
	recentEntries := d.auditTrail.GetEntriesByDateRange(time.Now().AddDate(0, -1, 0), time.Now())
	if len(employeeData) > 0 {
		metrics["data_timeliness"] = (float64(len(recentEntries)) / float64(len(employeeData))) * 100
	}

	return metrics
}

func (d *DocumentationGenerator) GenerateMarkdownReport(report *ComplianceReport) string {
	var md strings.Builder

	// Header
	md.WriteString(fmt.Sprintf("# Relatório de Compliance - Vale Refeição\n\n"))
	md.WriteString(fmt.Sprintf("**Data de Geração:** %s\n", report.GeneratedAt.Format("02/01/2006 15:04:05")))
	md.WriteString(fmt.Sprintf("**Período:** %s a %s\n\n",
		report.Period.StartDate.Format("02/01/2006"),
		report.Period.EndDate.Format("02/01/2006")))

	// Executive Summary
	md.WriteString("## 📊 Resumo Executivo\n\n")
	md.WriteString(fmt.Sprintf("- **Total de Funcionários Analisados:** %d\n", report.Summary.TotalEmployees))
	md.WriteString(fmt.Sprintf("- **Taxa de Compliance:** %.2f%%\n", report.Summary.ComplianceRate))
	md.WriteString(fmt.Sprintf("- **Violações Identificadas:** %d\n", len(report.Violations)))
	md.WriteString(fmt.Sprintf("- **Avaliações de Risco:** %d\n\n", len(report.RiskAssessments)))

	// Violations Summary
	if len(report.Violations) > 0 {
		md.WriteString("## ⚠️ Violações por Severidade\n\n")
		md.WriteString("| Severidade | Quantidade |\n")
		md.WriteString("|------------|------------|\n")
		for severity, count := range report.Summary.ViolationsBySeverity {
			if count > 0 {
				md.WriteString(fmt.Sprintf("| %s | %d |\n", d.severityToString(severity), count))
			}
		}
		md.WriteString("\n")
	}

	// Risk Summary
	if len(report.RiskAssessments) > 0 {
		md.WriteString("## 🎯 Distribuição de Riscos\n\n")
		md.WriteString("| Nível de Risco | Quantidade |\n")
		md.WriteString("|----------------|------------|\n")
		for level, count := range report.Summary.RisksByLevel {
			if count > 0 {
				md.WriteString(fmt.Sprintf("| %s | %d |\n", d.riskLevelToString(level), count))
			}
		}
		md.WriteString("\n")
	}

	// Recommendations
	if len(report.Recommendations) > 0 {
		md.WriteString("## 💡 Recomendações\n\n")
		for i, rec := range report.Recommendations {
			md.WriteString(fmt.Sprintf("%d. %s\n", i+1, rec))
		}
		md.WriteString("\n")
	}

	// Quality Metrics
	if len(report.Metadata.QualityMetrics) > 0 {
		md.WriteString("## 📈 Métricas de Qualidade\n\n")
		md.WriteString("| Métrica | Valor |\n")
		md.WriteString("|---------|-------|\n")
		for metric, value := range report.Metadata.QualityMetrics {
			md.WriteString(fmt.Sprintf("| %s | %.2f%% |\n", d.metricToString(metric), value))
		}
		md.WriteString("\n")
	}

	// Detailed Violations (if any)
	if len(report.Violations) > 0 {
		md.WriteString("## 📋 Detalhamento das Violações\n\n")
		for i, violation := range report.Violations {
			if i >= 10 { // Limit to first 10 violations
				md.WriteString(fmt.Sprintf("... e mais %d violação(ões)\n\n", len(report.Violations)-10))
				break
			}
			md.WriteString(fmt.Sprintf("### Violação %d - %s\n", i+1, violation.RuleID))
			md.WriteString(fmt.Sprintf("- **Funcionário:** %s\n", violation.EntityID))
			md.WriteString(fmt.Sprintf("- **Descrição:** %s\n", violation.Description))
			md.WriteString(fmt.Sprintf("- **Severidade:** %s\n", d.severityToString(violation.Severity)))
			if len(violation.Recommendations) > 0 {
				md.WriteString("- **Ações Recomendadas:**\n")
				for _, rec := range violation.Recommendations {
					md.WriteString(fmt.Sprintf("  - %s\n", rec))
				}
			}
			md.WriteString("\n")
		}
	}

	// Footer
	md.WriteString("---\n")
	md.WriteString(fmt.Sprintf("*Relatório gerado por: %s v%s*\n",
		report.Metadata.GeneratedBy, report.Metadata.Version))

	return md.String()
}

func (d *DocumentationGenerator) severityToString(severity ViolationSeverity) string {
	switch severity {
	case SeverityCritical:
		return "🔴 Crítica"
	case SeverityHigh:
		return "🟠 Alta"
	case SeverityMedium:
		return "🟡 Média"
	case SeverityLow:
		return "🟢 Baixa"
	default:
		return "❓ Desconhecida"
	}
}

func (d *DocumentationGenerator) riskLevelToString(level RiskLevel) string {
	switch level {
	case RiskCritical:
		return "🔴 Crítico"
	case RiskHigh:
		return "🟠 Alto"
	case RiskMedium:
		return "🟡 Médio"
	case RiskLow:
		return "🟢 Baixo"
	case RiskNone:
		return "⚪ Nenhum"
	default:
		return "❓ Desconhecido"
	}
}

func (d *DocumentationGenerator) metricToString(metric string) string {
	switch metric {
	case "data_completeness":
		return "Completude dos Dados"
	case "data_accuracy":
		return "Precisão dos Dados"
	case "data_timeliness":
		return "Atualidade dos Dados"
	default:
		return metric
	}
}

func (d *DocumentationGenerator) GenerateAuditReport(startDate, endDate time.Time) *AuditReport {
	entries := d.auditTrail.GetEntriesByDateRange(startDate, endDate)

	report := &AuditReport{
		GeneratedAt:  time.Now(),
		Period:       DateRange{StartDate: startDate, EndDate: endDate},
		TotalEntries: len(entries),
		Entries:      entries,
		Summary:      d.generateAuditSummary(entries),
	}

	return report
}

func (d *DocumentationGenerator) generateAuditSummary(entries []AuditEntry) AuditSummary {
	summary := AuditSummary{
		EntriesByLevel:  make(map[string]int),
		EntriesByAction: make(map[string]int),
		TopUsers:        []string{},
	}

	userCount := make(map[string]int)

	for _, entry := range entries {
		summary.EntriesByLevel[entry.Level]++
		summary.EntriesByAction[entry.Action]++
		userCount[entry.User]++
	}

	// Get top users
	type userActivity struct {
		user  string
		count int
	}
	users := []userActivity{}
	for user, count := range userCount {
		users = append(users, userActivity{user, count})
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].count > users[j].count
	})

	for i, user := range users {
		if i >= 5 { // Top 5 users
			break
		}
		summary.TopUsers = append(summary.TopUsers, fmt.Sprintf("%s (%d)", user.user, user.count))
	}

	return summary
}

type AuditReport struct {
	GeneratedAt  time.Time    `json:"generated_at"`
	Period       DateRange    `json:"period"`
	TotalEntries int          `json:"total_entries"`
	Summary      AuditSummary `json:"summary"`
	Entries      []AuditEntry `json:"entries"`
}

type AuditSummary struct {
	EntriesByLevel  map[string]int `json:"entries_by_level"`
	EntriesByAction map[string]int `json:"entries_by_action"`
	TopUsers        []string       `json:"top_users"`
}
