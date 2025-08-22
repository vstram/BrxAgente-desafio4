package intelligence

import (
	"fmt"
	"time"
)

// AnomalyType define os tipos de anomalia
type AnomalyType string

const (
	AnomalyTypeValue        AnomalyType = "value"
	AnomalyTypeTemporal     AnomalyType = "temporal"
	AnomalyTypeRelationship AnomalyType = "relationship"
	AnomalyTypeStructural   AnomalyType = "structural"
)

// AnomalySeverity define níveis de severidade
type AnomalySeverity int

const (
	SeverityInfo     AnomalySeverity = 1  // 1-2: Informativo
	SeverityLow      AnomalySeverity = 3  // 3-4: Baixa
	SeverityMedium   AnomalySeverity = 5  // 5-6: Média
	SeverityHigh     AnomalySeverity = 7  // 7-8: Alta
	SeverityCritical AnomalySeverity = 9  // 9-10: Crítica
)

func (s AnomalySeverity) String() string {
	switch {
	case s <= 2:
		return "info"
	case s <= 4:
		return "low"
	case s <= 6:
		return "medium"
	case s <= 8:
		return "high"
	default:
		return "critical"
	}
}

// Anomaly representa uma anomalia detectada
type Anomaly struct {
	ID          string                 `json:"id"`
	Type        AnomalyType            `json:"type"`
	Severity    AnomalySeverity        `json:"severity"`
	Confidence  float64                `json:"confidence"`  // 0-100%
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	DetectedAt  time.Time              `json:"detected_at"`
	
	// Dados específicos da anomalia
	Entity      string                 `json:"entity"`      // MATRICULA, SINDICATO, etc.
	EntityValue string                 `json:"entity_value"` // Valor do campo
	FieldName   string                 `json:"field_name"`  // Nome do campo problemático
	
	// Contexto adicional
	Data        map[string]interface{} `json:"data"`
	Suggestions []string               `json:"suggestions"`
	
	// Metadados
	DetectorName string `json:"detector_name"`
	RuleName     string `json:"rule_name"`
}

// AnomalyRule define uma regra de detecção
type AnomalyRule struct {
	Name        string
	Description string
	Type        AnomalyType
	Severity    AnomalySeverity
	Enabled     bool
	
	// Função de detecção
	Detector func(ctx *AnalysisContext) []Anomaly
}

// AnomalyReport representa um relatório completo de anomalias
type AnomalyReport struct {
	GeneratedAt time.Time `json:"generated_at"`
	TotalRecords int      `json:"total_records"`
	
	// Estatísticas de anomalias
	TotalAnomalies  int `json:"total_anomalies"`
	AnomaliesByType map[AnomalyType]int `json:"anomalies_by_type"`
	AnomaliesBySeverity map[string]int `json:"anomalies_by_severity"`
	
	// Anomalias detectadas
	Anomalies []Anomaly `json:"anomalies"`
	
	// Resumo executivo
	Summary AnomalySummary `json:"summary"`
}

// AnomalySummary fornece resumo executivo do relatório
type AnomalySummary struct {
	OverallScore     float64 `json:"overall_score"`     // 0-100 (100 = sem problemas)
	RiskLevel        string  `json:"risk_level"`        // low, medium, high, critical
	CriticalIssues   int     `json:"critical_issues"`
	RecommendedActions []string `json:"recommended_actions"`
	
	// Top problemas por categoria
	TopValueIssues        []string `json:"top_value_issues"`
	TopTemporalIssues     []string `json:"top_temporal_issues"`
	TopRelationshipIssues []string `json:"top_relationship_issues"`
}

// AnalysisContext fornece contexto para análise de anomalias
type AnalysisContext struct {
	// Dados de entrada
	Colaboradores map[string]interface{} `json:"colaboradores"`
	Parameters    map[string]interface{} `json:"parameters"`
	
	// Configurações
	Config *AnalysisConfig `json:"config"`
	
	// Metadados
	AnalysisDate time.Time `json:"analysis_date"`
	DataSources  []string  `json:"data_sources"`
	
	// Estado da análise
	ProcessedRecords int              `json:"processed_records"`
	Errors          []error           `json:"errors"`
	Warnings        []string          `json:"warnings"`
	
	// Cache para performance
	cache map[string]interface{}
}

// AnalysisConfig contém configurações para análise
type AnalysisConfig struct {
	// Thresholds gerais
	OutlierThreshold     float64 `json:"outlier_threshold"`      // Número de desvios padrão (default: 2.5)
	ConfidenceThreshold  float64 `json:"confidence_threshold"`   // Confidence mínimo para reportar (default: 70%)
	
	// Configurações específicas
	VRMinValue          float64 `json:"vr_min_value"`           // VR mínimo esperado
	VRMaxValue          float64 `json:"vr_max_value"`           // VR máximo esperado
	MaxWorkDaysPerMonth int     `json:"max_work_days_per_month"` // Dias úteis máximos por mês
	
	// Configurações temporais
	MaxFutureDate       time.Time `json:"max_future_date"`       // Data máxima aceita para admissão
	MinAdmissionDate    time.Time `json:"min_admission_date"`    // Data mínima aceita para admissão
	
	// Configurações de performance
	EnableParallelProcessing bool `json:"enable_parallel_processing"`
	MaxConcurrentDetectors   int  `json:"max_concurrent_detectors"`
	
	// Flags de funcionalidades
	EnableValueDetection        bool `json:"enable_value_detection"`
	EnableTemporalDetection     bool `json:"enable_temporal_detection"`
	EnableRelationshipDetection bool `json:"enable_relationship_detection"`
}

// DefaultAnalysisConfig retorna configuração padrão
func DefaultAnalysisConfig() *AnalysisConfig {
	return &AnalysisConfig{
		OutlierThreshold:     2.5,
		ConfidenceThreshold:  70.0,
		VRMinValue:          50.0,   // R$ 50,00
		VRMaxValue:          2000.0, // R$ 2.000,00
		MaxWorkDaysPerMonth: 31,
		MaxFutureDate:       time.Now().AddDate(0, 1, 0), // 1 mês no futuro
		MinAdmissionDate:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		
		EnableParallelProcessing:    true,
		MaxConcurrentDetectors:     3,
		EnableValueDetection:        true,
		EnableTemporalDetection:     true,
		EnableRelationshipDetection: true,
	}
}

// NewAnalysisContext cria novo contexto de análise
func NewAnalysisContext(colaboradores map[string]interface{}, config *AnalysisConfig) *AnalysisContext {
	if config == nil {
		config = DefaultAnalysisConfig()
	}
	
	return &AnalysisContext{
		Colaboradores:    colaboradores,
		Parameters:       make(map[string]interface{}),
		Config:          config,
		AnalysisDate:    time.Now(),
		DataSources:     []string{"excel"},
		ProcessedRecords: 0,
		Errors:          make([]error, 0),
		Warnings:        make([]string, 0),
		cache:           make(map[string]interface{}),
	}
}

// AddError adiciona erro ao contexto
func (ctx *AnalysisContext) AddError(err error) {
	ctx.Errors = append(ctx.Errors, err)
}

// AddWarning adiciona warning ao contexto
func (ctx *AnalysisContext) AddWarning(warning string) {
	ctx.Warnings = append(ctx.Warnings, warning)
}

// SetCache define valor no cache
func (ctx *AnalysisContext) SetCache(key string, value interface{}) {
	ctx.cache[key] = value
}

// GetCache recupera valor do cache
func (ctx *AnalysisContext) GetCache(key string) (interface{}, bool) {
	value, exists := ctx.cache[key]
	return value, exists
}

// IncrementProcessed incrementa contador de registros processados
func (ctx *AnalysisContext) IncrementProcessed() {
	ctx.ProcessedRecords++
}

// NewAnomaly cria nova anomalia
func NewAnomaly(
	anomalyType AnomalyType,
	severity AnomalySeverity,
	confidence float64,
	title, description string,
	entity, entityValue, fieldName string,
	detectorName, ruleName string,
) Anomaly {
	return Anomaly{
		ID:           fmt.Sprintf("%s_%d_%s", anomalyType, time.Now().UnixNano(), entity),
		Type:         anomalyType,
		Severity:     severity,
		Confidence:   confidence,
		Title:        title,
		Description:  description,
		DetectedAt:   time.Now(),
		Entity:       entity,
		EntityValue:  entityValue,
		FieldName:    fieldName,
		Data:         make(map[string]interface{}),
		Suggestions:  make([]string, 0),
		DetectorName: detectorName,
		RuleName:     ruleName,
	}
}

// AddData adiciona dados contextuais à anomalia
func (a *Anomaly) AddData(key string, value interface{}) {
	a.Data[key] = value
}

// AddSuggestion adiciona sugestão de correção
func (a *Anomaly) AddSuggestion(suggestion string) {
	a.Suggestions = append(a.Suggestions, suggestion)
}

// IsHighPriority verifica se a anomalia é de alta prioridade
func (a *Anomaly) IsHighPriority() bool {
	return a.Severity >= SeverityHigh && a.Confidence >= 80.0
}

// IsCritical verifica se a anomalia é crítica
func (a *Anomaly) IsCritical() bool {
	return a.Severity >= SeverityCritical
}