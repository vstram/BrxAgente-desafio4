package compliance

import (
	"fmt"
	"math"
	"time"
)

type RiskAssessor struct {
	complianceChecker *ComplianceChecker
	auditTrail       *AuditTrail
}

func NewRiskAssessor(checker *ComplianceChecker, audit *AuditTrail) *RiskAssessor {
	return &RiskAssessor{
		complianceChecker: checker,
		auditTrail:       audit,
	}
}

type RiskAssessment struct {
	EntityID          string    `json:"entity_id"`
	EntityType        string    `json:"entity_type"`
	OverallRisk       RiskLevel `json:"overall_risk"`
	RiskScore         float64   `json:"risk_score"`
	RiskFactors       []RiskFactor `json:"risk_factors"`
	Recommendations   []string  `json:"recommendations"`
	AssessmentDate    time.Time `json:"assessment_date"`
	NextReviewDate    time.Time `json:"next_review_date"`
}

type RiskFactor struct {
	Category     string    `json:"category"`
	Description  string    `json:"description"`
	Impact       float64   `json:"impact"`
	Probability  float64   `json:"probability"`
	RiskScore    float64   `json:"risk_score"`
	Severity     RiskLevel `json:"severity"`
	Mitigation   []string  `json:"mitigation"`
}

type RiskMatrix struct {
	Probability map[string]float64 `json:"probability"`
	Impact      map[string]float64 `json:"impact"`
	Thresholds  map[string]float64 `json:"thresholds"`
}

func (r *RiskAssessor) GetRiskMatrix() RiskMatrix {
	return RiskMatrix{
		Probability: map[string]float64{
			"very_low":  0.1,
			"low":       0.3,
			"medium":    0.5,
			"high":      0.7,
			"very_high": 0.9,
		},
		Impact: map[string]float64{
			"negligible": 0.1,
			"minor":      0.3,
			"moderate":   0.5,
			"major":      0.7,
			"critical":   0.9,
		},
		Thresholds: map[string]float64{
			"low":      0.3,
			"medium":   0.6,
			"high":     0.8,
			"critical": 0.95,
		},
	}
}

func (r *RiskAssessor) AssessEmployeeRisk(employeeData map[string]interface{}) (*RiskAssessment, error) {
	entityID, ok := employeeData["matricula"].(string)
	if !ok {
		return nil, fmt.Errorf("matricula não encontrada nos dados do funcionário")
	}

	assessment := &RiskAssessment{
		EntityID:       entityID,
		EntityType:     "employee",
		AssessmentDate: time.Now(),
		NextReviewDate: time.Now().AddDate(0, 3, 0), // Review every 3 months
	}

	riskFactors := []RiskFactor{}

	// Check compliance violations
	complianceResult, _ := r.complianceChecker.CheckCompliance(employeeData)
	if len(complianceResult.Violations) > 0 {
		factor := r.assessComplianceRisk(complianceResult.Violations)
		riskFactors = append(riskFactors, factor)
	}

	// Assess VR value risk
	vrRisk := r.assessVRValueRisk(employeeData)
	if vrRisk.RiskScore > 0 {
		riskFactors = append(riskFactors, vrRisk)
	}

	// Assess documentation risk
	docRisk := r.assessDocumentationRisk(employeeData)
	if docRisk.RiskScore > 0 {
		riskFactors = append(riskFactors, docRisk)
	}

	// Assess historical risk based on audit trail
	histRisk := r.assessHistoricalRisk(entityID)
	if histRisk.RiskScore > 0 {
		riskFactors = append(riskFactors, histRisk)
	}

	assessment.RiskFactors = riskFactors
	assessment.RiskScore = r.calculateOverallRisk(riskFactors)
	assessment.OverallRisk = r.determineRiskLevel(assessment.RiskScore)
	assessment.Recommendations = r.generateRecommendations(riskFactors)

	return assessment, nil
}

func (r *RiskAssessor) assessComplianceRisk(violations []ComplianceViolation) RiskFactor {
	factor := RiskFactor{
		Category:    "compliance",
		Description: "Violações de compliance identificadas",
		Mitigation:  []string{},
	}

	var totalImpact float64
	var criticalCount int

	for _, violation := range violations {
		switch violation.Severity {
		case SeverityCritical:
			totalImpact += 0.9
			criticalCount++
			factor.Mitigation = append(factor.Mitigation, fmt.Sprintf("Resolver imediatamente: %s", violation.Description))
		case SeverityHigh:
			totalImpact += 0.7
			factor.Mitigation = append(factor.Mitigation, fmt.Sprintf("Resolver com prioridade: %s", violation.Description))
		case SeverityMedium:
			totalImpact += 0.5
			factor.Mitigation = append(factor.Mitigation, fmt.Sprintf("Programar correção: %s", violation.Description))
		case SeverityLow:
			totalImpact += 0.3
			factor.Mitigation = append(factor.Mitigation, fmt.Sprintf("Monitorar: %s", violation.Description))
		}
	}

	factor.Impact = math.Min(totalImpact/float64(len(violations)), 0.9)
	factor.Probability = 0.8 // High probability since violations were already detected
	factor.RiskScore = factor.Impact * factor.Probability

	if criticalCount > 0 {
		factor.Severity = RiskCritical
	} else if factor.RiskScore >= 0.6 {
		factor.Severity = RiskHigh
	} else if factor.RiskScore >= 0.3 {
		factor.Severity = RiskMedium
	} else {
		factor.Severity = RiskLow
	}

	return factor
}

func (r *RiskAssessor) assessVRValueRisk(employeeData map[string]interface{}) RiskFactor {
	factor := RiskFactor{
		Category:    "vr_value",
		Description: "Análise de risco relacionada aos valores de Vale Refeição",
		Mitigation:  []string{},
	}

	vrValue, ok := employeeData["valor_vr"].(float64)
	if !ok {
		factor.Impact = 0.5
		factor.Probability = 0.3
		factor.Mitigation = append(factor.Mitigation, "Verificar e validar valor do VR")
	} else {
		// Check if value exceeds legal limits
		if vrValue > 55.00 {
			factor.Impact = 0.8
			factor.Probability = 0.9
			factor.Description = "Valor do VR acima do limite legal (R$ 55,00)"
			factor.Mitigation = append(factor.Mitigation, "Ajustar valor do VR para dentro dos limites legais")
		} else if vrValue < 2.00 {
			factor.Impact = 0.6
			factor.Probability = 0.7
			factor.Description = "Valor do VR abaixo do mínimo recomendado"
			factor.Mitigation = append(factor.Mitigation, "Revisar valor mínimo do VR")
		} else {
			factor.Impact = 0.1
			factor.Probability = 0.1
		}
	}

	factor.RiskScore = factor.Impact * factor.Probability
	factor.Severity = r.determineRiskLevel(factor.RiskScore)

	return factor
}

func (r *RiskAssessor) assessDocumentationRisk(employeeData map[string]interface{}) RiskFactor {
	factor := RiskFactor{
		Category:    "documentation",
		Description: "Risco relacionado à documentação e conformidade documental",
		Mitigation:  []string{},
	}

	requiredFields := []string{"matricula", "data_admissao", "sindicato"}
	missingFields := []string{}

	for _, field := range requiredFields {
		if _, ok := employeeData[field]; !ok {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		factor.Impact = float64(len(missingFields)) * 0.2
		factor.Probability = 0.8
		factor.Description = fmt.Sprintf("Campos obrigatórios faltando: %v", missingFields)
		factor.Mitigation = append(factor.Mitigation, "Completar documentação obrigatória")
	} else {
		factor.Impact = 0.1
		factor.Probability = 0.2
	}

	factor.RiskScore = factor.Impact * factor.Probability
	factor.Severity = r.determineRiskLevel(factor.RiskScore)

	return factor
}

func (r *RiskAssessor) assessHistoricalRisk(entityID string) RiskFactor {
	factor := RiskFactor{
		Category:    "historical",
		Description: "Análise de risco baseada no histórico de auditoria",
		Mitigation:  []string{},
	}

	// Get last 6 months of audit entries
	since := time.Now().AddDate(0, -6, 0)
	entries := r.auditTrail.GetEntriesForEntity(entityID, since, time.Now())

	if len(entries) == 0 {
		factor.Impact = 0.3
		factor.Probability = 0.4
		factor.Description = "Falta de histórico de auditoria"
		factor.Mitigation = append(factor.Mitigation, "Implementar auditoria regular")
	} else {
		errorCount := 0
		warningCount := 0

		for _, entry := range entries {
			if entry.Level == "ERROR" {
				errorCount++
			} else if entry.Level == "WARN" {
				warningCount++
			}
		}

		if errorCount > 3 {
			factor.Impact = 0.8
			factor.Probability = 0.7
			factor.Description = fmt.Sprintf("Múltiplos erros no histórico (%d erros)", errorCount)
			factor.Mitigation = append(factor.Mitigation, "Investigar e resolver erros recorrentes")
		} else if warningCount > 5 {
			factor.Impact = 0.5
			factor.Probability = 0.6
			factor.Description = fmt.Sprintf("Múltiplos avisos no histórico (%d avisos)", warningCount)
			factor.Mitigation = append(factor.Mitigation, "Revisar processos para reduzir avisos")
		} else {
			factor.Impact = 0.2
			factor.Probability = 0.3
		}
	}

	factor.RiskScore = factor.Impact * factor.Probability
	factor.Severity = r.determineRiskLevel(factor.RiskScore)

	return factor
}

func (r *RiskAssessor) calculateOverallRisk(factors []RiskFactor) float64 {
	if len(factors) == 0 {
		return 0.0
	}

	var totalWeightedRisk float64
	var totalWeight float64

	for _, factor := range factors {
		weight := r.getCategoryWeight(factor.Category)
		totalWeightedRisk += factor.RiskScore * weight
		totalWeight += weight
	}

	return totalWeightedRisk / totalWeight
}

func (r *RiskAssessor) getCategoryWeight(category string) float64 {
	weights := map[string]float64{
		"compliance":     1.0, // Highest weight
		"vr_value":       0.8,
		"documentation":  0.6,
		"historical":     0.4, // Lowest weight
	}

	if weight, ok := weights[category]; ok {
		return weight
	}
	return 0.5 // Default weight
}

func (r *RiskAssessor) determineRiskLevel(riskScore float64) RiskLevel {
	matrix := r.GetRiskMatrix()

	if riskScore >= matrix.Thresholds["critical"] {
		return RiskCritical
	} else if riskScore >= matrix.Thresholds["high"] {
		return RiskHigh
	} else if riskScore >= matrix.Thresholds["medium"] {
		return RiskMedium
	} else if riskScore >= matrix.Thresholds["low"] {
		return RiskLow
	}

	return RiskNone
}

func (r *RiskAssessor) generateRecommendations(factors []RiskFactor) []string {
	recommendations := []string{}

	for _, factor := range factors {
		recommendations = append(recommendations, factor.Mitigation...)
	}

	// Add general recommendations based on risk level
	hasHighRisk := false
	hasCriticalRisk := false

	for _, factor := range factors {
		if factor.Severity == RiskHigh {
			hasHighRisk = true
		}
		if factor.Severity == RiskCritical {
			hasCriticalRisk = true
		}
	}

	if hasCriticalRisk {
		recommendations = append(recommendations, "Solicitar revisão urgente da diretoria")
		recommendations = append(recommendations, "Implementar plano de contingência")
	}

	if hasHighRisk {
		recommendations = append(recommendations, "Agendar auditoria extraordinária")
		recommendations = append(recommendations, "Revisar controles internos")
	}

	return recommendations
}

func (r *RiskAssessor) GenerateRiskReport(assessments []RiskAssessment) *RiskReport {
	report := &RiskReport{
		GeneratedAt:    time.Now(),
		TotalEntities:  len(assessments),
		Summary:        make(map[RiskLevel]int),
		Assessments:    assessments,
	}

	for _, assessment := range assessments {
		report.Summary[assessment.OverallRisk]++
	}

	// Calculate aggregate metrics
	if len(assessments) > 0 {
		var totalScore float64
		for _, assessment := range assessments {
			totalScore += assessment.RiskScore
		}
		report.AverageRiskScore = totalScore / float64(len(assessments))
	}

	return report
}

type RiskReport struct {
	GeneratedAt       time.Time              `json:"generated_at"`
	TotalEntities     int                    `json:"total_entities"`
	AverageRiskScore  float64               `json:"average_risk_score"`
	Summary           map[RiskLevel]int      `json:"summary"`
	Assessments       []RiskAssessment      `json:"assessments"`
}