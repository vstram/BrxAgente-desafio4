package models

import (
	"time"
)

// Prediction representa uma predição gerada pelo sistema
type Prediction struct {
	ID          string           `json:"id"`
	Type        PredictionType   `json:"type"`
	Target      string           `json:"target"`      // o que está sendo previsto
	Value       interface{}      `json:"value"`       // valor previsto
	Confidence  float64          `json:"confidence"`  // 0.0 a 1.0
	Timeframe   string           `json:"timeframe"`   // período da previsão
	CreatedAt   time.Time        `json:"created_at"`
	ValidUntil  time.Time        `json:"valid_until"`
	Method      string           `json:"method"`      // método usado para predição
	Description string           `json:"description"`
	Metadata    PredictionMeta   `json:"metadata"`
	Actions     []ActionItem     `json:"actions"`     // ações recomendadas
}

// PredictionType enumera tipos de predições
type PredictionType string

const (
	PredictionConsumption    PredictionType = "consumption"     // previsão de consumo
	PredictionAnomaly        PredictionType = "anomaly"         // previsão de anomalia
	PredictionTrend          PredictionType = "trend"           // previsão de tendência
	PredictionRisk           PredictionType = "risk"            // avaliação de risco
	PredictionOptimization   PredictionType = "optimization"    // otimização de processo
	PredictionResource       PredictionType = "resource"        // necessidade de recursos
)

// PredictionMeta contém metadados sobre a predição
type PredictionMeta struct {
	Model          string                 `json:"model"`           // modelo usado
	Features       []string               `json:"features"`        // features consideradas
	DataPoints     int                    `json:"data_points"`     // pontos de dados usados
	TrainingPeriod string                 `json:"training_period"` // período de treino
	Accuracy       float64                `json:"accuracy"`        // acurácia estimada
	Parameters     map[string]interface{} `json:"parameters"`      // parâmetros do modelo
}

// ActionItem representa uma ação recomendada baseada em predição
type ActionItem struct {
	ID          string       `json:"id"`
	Priority    Priority     `json:"priority"`
	Category    string       `json:"category"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	DueDate     *time.Time   `json:"due_date,omitempty"`
	Owner       string       `json:"owner,omitempty"`
	Status      ActionStatus `json:"status"`
}

// Priority enumera níveis de prioridade
type Priority string

const (
	PriorityLow       Priority = "low"
	PriorityMedium    Priority = "medium"
	PriorityHigh      Priority = "high"
	PriorityCritical  Priority = "critical"
)

// ActionStatus enumera status das ações
type ActionStatus string

const (
	ActionPending    ActionStatus = "pending"
	ActionInProgress ActionStatus = "in_progress"
	ActionCompleted  ActionStatus = "completed"
	ActionCancelled  ActionStatus = "cancelled"
)

// ConsumptionForecast representa uma previsão de consumo de VR
type ConsumptionForecast struct {
	Sindicato     string              `json:"sindicato"`
	Month         time.Time           `json:"month"`
	PredictedVR   float64             `json:"predicted_vr"`
	Confidence    float64             `json:"confidence"`
	Range         ForecastRange       `json:"range"`           // intervalo de confiança
	Factors       []string            `json:"factors"`         // fatores considerados
	Seasonality   *SeasonalityInfo    `json:"seasonality"`     // info sazonal
	Trend         *VRTrend            `json:"trend"`           // tendência identificada
	Assumptions   []string            `json:"assumptions"`     // premissas da previsão
}

// ForecastRange representa intervalo de confiança da previsão
type ForecastRange struct {
	Lower float64 `json:"lower"`  // limite inferior
	Upper float64 `json:"upper"`  // limite superior
}

// RiskAssessment representa avaliação de risco para colaborador
type RiskAssessment struct {
	Matricula     string        `json:"matricula"`
	Sindicato     string        `json:"sindicato"`
	RiskScore     float64       `json:"risk_score"`     // 0-100
	RiskLevel     RiskLevel     `json:"risk_level"`
	RiskFactors   []RiskFactor  `json:"risk_factors"`
	Probability   float64       `json:"probability"`    // probabilidade de problema
	Impact        ImpactLevel   `json:"impact"`         // impacto potencial
	Mitigation    []ActionItem  `json:"mitigation"`     // ações de mitigação
	LastEvaluated time.Time     `json:"last_evaluated"`
}

// RiskLevel enumera níveis de risco
type RiskLevel string

const (
	RiskVeryLow  RiskLevel = "very_low"
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

// RiskFactor representa um fator de risco identificado
type RiskFactor struct {
	Factor      string  `json:"factor"`      // nome do fator
	Weight      float64 `json:"weight"`      // peso do fator (0-1)
	Value       float64 `json:"value"`       // valor atual do fator
	Threshold   float64 `json:"threshold"`   // limite de risco
	Description string  `json:"description"`
}

// ProcessOptimization representa otimização de processo sugerida
type ProcessOptimization struct {
	ProcessID     string                    `json:"process_id"`
	Month         time.Time                 `json:"month"`
	CurrentState  ProcessState              `json:"current_state"`
	OptimalState  ProcessState              `json:"optimal_state"`
	Improvements  []ImprovementSuggestion   `json:"improvements"`
	ExpectedGains OptimizationGains         `json:"expected_gains"`
	Implementation ImplementationPlan       `json:"implementation"`
}

// ProcessState representa estado de um processo
type ProcessState struct {
	Duration        time.Duration          `json:"duration"`
	ResourceUsage   map[string]float64     `json:"resource_usage"`
	Efficiency      float64                `json:"efficiency"`
	ErrorRate       float64                `json:"error_rate"`
	Throughput      int                    `json:"throughput"`
	Metrics         map[string]interface{} `json:"metrics"`
}

// ImprovementSuggestion representa sugestão de melhoria
type ImprovementSuggestion struct {
	Area        string    `json:"area"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Impact      float64   `json:"impact"`      // impacto esperado (0-1)
	Effort      float64   `json:"effort"`      // esforço necessário (0-1)
	Priority    Priority  `json:"priority"`
}

// OptimizationGains representa ganhos esperados de otimização
type OptimizationGains struct {
	TimeReduction     float64 `json:"time_reduction"`     // % redução de tempo
	ResourceSaving    float64 `json:"resource_saving"`    // % economia de recursos
	AccuracyIncrease  float64 `json:"accuracy_increase"`  // % aumento de acurácia
	CostReduction     float64 `json:"cost_reduction"`     // % redução de custos
}

// ImplementationPlan representa plano de implementação
type ImplementationPlan struct {
	Steps       []ImplementationStep `json:"steps"`
	Timeline    string               `json:"timeline"`
	Resources   []string             `json:"resources"`
	Risks       []string             `json:"risks"`
	Success     []string             `json:"success_criteria"`
}

// ImplementationStep representa passo de implementação
type ImplementationStep struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Duration    string     `json:"duration"`
	Owner       string     `json:"owner"`
	Dependencies []string  `json:"dependencies"`
}

// PredictionModel representa um modelo de predição
type PredictionModel struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Type           ModelType              `json:"type"`
	Version        string                 `json:"version"`
	Algorithm      string                 `json:"algorithm"`
	Parameters     map[string]interface{} `json:"parameters"`
	TrainingData   ModelTrainingData      `json:"training_data"`
	Performance    ModelPerformance       `json:"performance"`
	LastTrained    time.Time              `json:"last_trained"`
	Status         ModelStatus            `json:"status"`
}

// ModelType enumera tipos de modelos
type ModelType string

const (
	ModelRegression    ModelType = "regression"
	ModelClassification ModelType = "classification"
	ModelTimeSeries    ModelType = "time_series"
	ModelClustering    ModelType = "clustering"
	ModelAnomaly       ModelType = "anomaly_detection"
)

// ModelTrainingData informações sobre dados de treino
type ModelTrainingData struct {
	Sources    []string  `json:"sources"`
	Size       int       `json:"size"`
	Period     string    `json:"period"`
	Features   []string  `json:"features"`
	LastUpdate time.Time `json:"last_update"`
}

// ModelPerformance métricas de performance do modelo
type ModelPerformance struct {
	Accuracy    float64            `json:"accuracy"`
	Precision   float64            `json:"precision"`
	Recall      float64            `json:"recall"`
	F1Score     float64            `json:"f1_score"`
	RMSE        float64            `json:"rmse,omitempty"`      // para regressão
	MAE         float64            `json:"mae,omitempty"`       // para regressão
	Custom      map[string]float64 `json:"custom,omitempty"`    // métricas customizadas
}

// ModelStatus enumera status dos modelos
type ModelStatus string

const (
	ModelActive    ModelStatus = "active"
	ModelTraining  ModelStatus = "training"
	ModelRetired   ModelStatus = "retired"
	ModelFailed    ModelStatus = "failed"
)

// NewPrediction cria uma nova predição
func NewPrediction(predType PredictionType, target string, value interface{}, confidence float64) *Prediction {
	return &Prediction{
		ID:          generatePredictionID(),
		Type:        predType,
		Target:      target,
		Value:       value,
		Confidence:  confidence,
		CreatedAt:   time.Now(),
		Actions:     make([]ActionItem, 0),
	}
}

// generatePredictionID gera um ID único para predição
func generatePredictionID() string {
	return "pred-" + time.Now().Format("20060102150405")
}

// IsValid verifica se a predição ainda é válida
func (p *Prediction) IsValid() bool {
	return time.Now().Before(p.ValidUntil)
}

// GetRiskLevel converte score de risco em nível
func GetRiskLevel(score float64) RiskLevel {
	switch {
	case score <= 20:
		return RiskVeryLow
	case score <= 40:
		return RiskLow
	case score <= 60:
		return RiskMedium
	case score <= 80:
		return RiskHigh
	default:
		return RiskCritical
	}
}

// GetPriorityFromRisk converte nível de risco em prioridade
func GetPriorityFromRisk(risk RiskLevel) Priority {
	switch risk {
	case RiskVeryLow, RiskLow:
		return PriorityLow
	case RiskMedium:
		return PriorityMedium
	case RiskHigh:
		return PriorityHigh
	case RiskCritical:
		return PriorityCritical
	default:
		return PriorityMedium
	}
}