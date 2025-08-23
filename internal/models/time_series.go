package models

import (
	"time"
)

// TimeSeries representa uma série temporal de dados
type TimeSeries struct {
	Points []TimePoint `json:"points"`
	Label  string      `json:"label"`
	Unit   string      `json:"unit"`
}

// TimePoint representa um ponto no tempo com valor
type TimePoint struct {
	Timestamp time.Time   `json:"timestamp"`
	Value     float64     `json:"value"`
	Metadata  interface{} `json:"metadata,omitempty"`
}

// HistoricalVRData representa dados históricos de VR
type HistoricalVRData struct {
	Month             time.Time          `json:"month"`
	Sindicato         string             `json:"sindicato"`
	TotalVR           float64            `json:"total_vr"`
	NumColaboradores  int                `json:"num_colaboradores"`
	MediaPorPessoa    float64            `json:"media_por_pessoa"`
	DaysProcessed     int                `json:"days_processed"`
	Anomalies         []string           `json:"anomalies"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// VRTrend representa uma tendência identificada nos dados de VR
type VRTrend struct {
	Type        TrendType `json:"type"`
	Strength    float64   `json:"strength"`    // -1.0 a 1.0
	Period      int       `json:"period"`      // em meses
	Confidence  float64   `json:"confidence"`  // 0.0 a 1.0
	StartDate   time.Time `json:"start_date"`
	Description string    `json:"description"`
}

// TrendType enumera os tipos de tendência
type TrendType string

const (
	TrendUpward    TrendType = "upward"
	TrendDownward  TrendType = "downward" 
	TrendStable    TrendType = "stable"
	TrendVolatile  TrendType = "volatile"
	TrendSeasonal  TrendType = "seasonal"
	TrendCyclical  TrendType = "cyclical"
)

// SeasonalityInfo contém informações sobre sazonalidade detectada
type SeasonalityInfo struct {
	IsDetected   bool                   `json:"is_detected"`
	Period       int                    `json:"period"`        // período em meses
	Amplitude    float64                `json:"amplitude"`     // amplitude da variação sazonal
	PeakMonths   []int                  `json:"peak_months"`   // meses de pico (1-12)
	TroughMonths []int                  `json:"trough_months"` // meses de vale (1-12)
	Confidence   float64                `json:"confidence"`
	Pattern      map[int]float64        `json:"pattern"`       // padrão sazonal por mês
}

// OutlierPattern representa um padrão de outliers identificado
type OutlierPattern struct {
	Type        OutlierType `json:"type"`
	Frequency   float64     `json:"frequency"`   // frequência de ocorrência
	Severity    float64     `json:"severity"`    // severidade média
	Affected    []string    `json:"affected"`    // entidades afetadas
	Description string      `json:"description"`
	Confidence  float64     `json:"confidence"`
}

// OutlierType enumera tipos de padrões de outliers
type OutlierType string

const (
	OutlierSpike     OutlierType = "spike"      // picos isolados
	OutlierDip       OutlierType = "dip"        // quedas isoladas
	OutlierSustained OutlierType = "sustained"  // valores sustentados fora do padrão
	OutlierCyclic    OutlierType = "cyclic"     // outliers cíclicos
)

// Pattern representa um padrão identificado nos dados
type Pattern struct {
	ID          string                 `json:"id"`
	Type        PatternType            `json:"type"`
	Description string                 `json:"description"`
	Confidence  float64                `json:"confidence"`
	StartDate   time.Time              `json:"start_date"`
	EndDate     *time.Time             `json:"end_date,omitempty"`
	Entities    []string               `json:"entities"`    // sindicatos, colaboradores afetados
	Attributes  map[string]interface{} `json:"attributes"`
	Impact      PatternImpact          `json:"impact"`
}

// PatternType enumera tipos de padrões
type PatternType string

const (
	PatternConsumption PatternType = "consumption"  // padrão de consumo
	PatternBehavior    PatternType = "behavior"     // padrão comportamental
	PatternSeasonal    PatternType = "seasonal"     // padrão sazonal
	PatternAnomaly     PatternType = "anomaly"      // padrão de anomalias
	PatternGrowth      PatternType = "growth"       // padrão de crescimento
)

// PatternImpact descreve o impacto de um padrão
type PatternImpact struct {
	Level       ImpactLevel `json:"level"`
	Areas       []string    `json:"areas"`       // áreas afetadas
	Magnitude   float64     `json:"magnitude"`   // magnitude do impacto
	Description string      `json:"description"`
}

// ImpactLevel enumera níveis de impacto
type ImpactLevel string

const (
	ImpactLow    ImpactLevel = "low"
	ImpactMedium ImpactLevel = "medium"
	ImpactHigh   ImpactLevel = "high"
	ImpactCritical ImpactLevel = "critical"
)

// MovingAverageConfig configuração para média móvel
type MovingAverageConfig struct {
	Window int     `json:"window"`  // janela de tempo (em períodos)
	Weight float64 `json:"weight"`  // peso para média ponderada
}

// TrendAnalysisConfig configuração para análise de tendências
type TrendAnalysisConfig struct {
	MinPeriods      int     `json:"min_periods"`      // mínimo de períodos para análise
	ConfidenceLevel float64 `json:"confidence_level"` // nível de confiança requerido
	Sensitivity     float64 `json:"sensitivity"`      // sensibilidade para detecção
}

// NewTimeSeries cria uma nova série temporal
func NewTimeSeries(label, unit string) *TimeSeries {
	return &TimeSeries{
		Points: make([]TimePoint, 0),
		Label:  label,
		Unit:   unit,
	}
}

// AddPoint adiciona um ponto à série temporal
func (ts *TimeSeries) AddPoint(timestamp time.Time, value float64, metadata interface{}) {
	point := TimePoint{
		Timestamp: timestamp,
		Value:     value,
		Metadata:  metadata,
	}
	ts.Points = append(ts.Points, point)
}

// Len retorna o número de pontos na série
func (ts *TimeSeries) Len() int {
	return len(ts.Points)
}

// GetRange retorna pontos dentro de um período específico
func (ts *TimeSeries) GetRange(start, end time.Time) []TimePoint {
	var result []TimePoint
	for _, point := range ts.Points {
		if point.Timestamp.After(start) && point.Timestamp.Before(end) {
			result = append(result, point)
		}
	}
	return result
}

// GetValues retorna apenas os valores da série
func (ts *TimeSeries) GetValues() []float64 {
	values := make([]float64, len(ts.Points))
	for i, point := range ts.Points {
		values[i] = point.Value
	}
	return values
}

// GetLatest retorna o último ponto da série
func (ts *TimeSeries) GetLatest() (*TimePoint, bool) {
	if len(ts.Points) == 0 {
		return nil, false
	}
	return &ts.Points[len(ts.Points)-1], true
}