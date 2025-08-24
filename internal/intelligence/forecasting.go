package intelligence

import (
	"fmt"
	"math"
	"sort"
	"time"

	"BrxAgente-desafio4/internal/predicoes"
)

// Forecaster implementa modelos de previsão para VR
type Forecaster struct {
	config    ForecastConfig
	models    map[string]ForecastModel
	predictor *TrendPredictor
	analyzer  *PatternAnalyzer
	detector  *TrendDetector
}

// ForecastConfig configuração do sistema de previsão
type ForecastConfig struct {
	DefaultHorizon      int     `json:"default_horizon"`      // horizonte padrão em meses
	MinAccuracy         float64 `json:"min_accuracy"`         // acurácia mínima aceitável
	EnsembleModels      bool    `json:"ensemble_models"`      // usar ensemble de modelos
	SeasonalAdjustment  bool    `json:"seasonal_adjustment"`  // ajuste sazonal
	ConfidenceIntervals bool    `json:"confidence_intervals"` // calcular intervalos de confiança
	ValidationPeriod    int     `json:"validation_period"`    // períodos para validação
}

// ForecastModel interface para modelos de previsão
type ForecastModel interface {
	Train(data []predicoes.HistoricalVRData) error
	Predict(horizon int) (*predicoes.ConsumptionForecast, error)
	GetAccuracy() float64
	GetName() string
}

// SimpleMovingAverageModel modelo de média móvel simples
type SimpleMovingAverageModel struct {
	name       string
	window     int
	accuracy   float64
	lastValues []float64
	sindicato  string
}

// ExponentialSmoothingModel modelo de suavização exponencial
type ExponentialSmoothingModel struct {
	name      string
	alpha     float64 // parâmetro de suavização
	accuracy  float64
	level     float64         // nível atual
	trend     float64         // tendência atual
	seasonal  map[int]float64 // componentes sazonais
	sindicato string
}

// LinearRegressionModel modelo de regressão linear
type LinearRegressionModel struct {
	name      string
	slope     float64
	intercept float64
	accuracy  float64
	sindicato string
}

// EnsembleForecast resultado de previsão ensemble
type EnsembleForecast struct {
	WeightedForecast    *predicoes.ConsumptionForecast            `json:"weighted_forecast"`
	IndividualForecasts map[string]*predicoes.ConsumptionForecast `json:"individual_forecasts"`
	ModelWeights        map[string]float64                        `json:"model_weights"`
	EnsembleAccuracy    float64                                   `json:"ensemble_accuracy"`
	ConsensusLevel      float64                                   `json:"consensus_level"`
}

// ForecastValidation resultado da validação de previsões
type ForecastValidation struct {
	Model            string             `json:"model"`
	Accuracy         float64            `json:"accuracy"`
	MAPE             float64            `json:"mape"` // Mean Absolute Percentage Error
	RMSE             float64            `json:"rmse"` // Root Mean Square Error
	MAE              float64            `json:"mae"`  // Mean Absolute Error
	Predictions      []ValidationPoint  `json:"predictions"`
	ActualVsForecast map[string]float64 `json:"actual_vs_forecast"`
}

// ValidationPoint ponto de validação
type ValidationPoint struct {
	Date      time.Time `json:"date"`
	Actual    float64   `json:"actual"`
	Predicted float64   `json:"predicted"`
	Error     float64   `json:"error"`
	ErrorPct  float64   `json:"error_pct"`
}

// NewForecaster cria novo sistema de previsão
func NewForecaster(config ForecastConfig) *Forecaster {
	return &Forecaster{
		config: config,
		models: make(map[string]ForecastModel),
		predictor: NewTrendPredictor(TrendPredictorConfig{
			MinDataPoints:       6,
			ConfidenceThreshold: 0.6,
			ForecastPeriods:     config.DefaultHorizon,
			MovingAverageWindow: 3,
		}),
		analyzer: NewPatternAnalyzer(PatternAnalyzerConfig{
			MinDataPoints:        6,
			CorrelationThreshold: 0.7,
			SeasonalityWindow:    12,
			AnomalyThreshold:     2.0,
			ConfidenceLevel:      0.6,
		}),
		detector: NewTrendDetector(TrendDetectorConfig{
			MinDataPoints:       6,
			SensitivityLevel:    0.5,
			TrendThreshold:      0.05,
			VolatilityThreshold: 0.3,
			SeasonalWindow:      12,
		}),
	}
}

// RegisterModel registra modelo de previsão
func (f *Forecaster) RegisterModel(model ForecastModel) {
	f.models[model.GetName()] = model
}

// ForecastConsumption prevê consumo de VR
func (f *Forecaster) ForecastConsumption(data []predicoes.HistoricalVRData, sindicato string, horizon int) (*EnsembleForecast, error) {
	if len(data) < 3 {
		return nil, fmt.Errorf("dados insuficientes para previsão")
	}

	// Filtrar dados para o sindicato específico
	sindicatoData := f.filterBySindicato(data, sindicato)
	if len(sindicatoData) < 3 {
		return nil, fmt.Errorf("dados insuficientes para sindicato %s", sindicato)
	}

	// Registrar modelos padrão se necessário
	f.ensureDefaultModels(sindicato)

	// Treinar todos os modelos
	modelForecasts := make(map[string]*predicoes.ConsumptionForecast)
	modelWeights := make(map[string]float64)

	for name, model := range f.models {
		err := model.Train(sindicatoData)
		if err != nil {
			continue
		}

		forecast, err := model.Predict(horizon)
		if err != nil {
			continue
		}

		modelForecasts[name] = forecast
		modelWeights[name] = model.GetAccuracy()
	}

	if len(modelForecasts) == 0 {
		return nil, fmt.Errorf("nenhum modelo conseguiu gerar previsão")
	}

	// Calcular previsão ensemble
	ensembleForecast := f.calculateEnsembleForecast(modelForecasts, modelWeights, sindicato, horizon)

	// Calcular nível de consenso
	consensusLevel := f.calculateConsensusLevel(modelForecasts)

	result := &EnsembleForecast{
		WeightedForecast:    ensembleForecast,
		IndividualForecasts: modelForecasts,
		ModelWeights:        modelWeights,
		EnsembleAccuracy:    f.calculateEnsembleAccuracy(modelWeights),
		ConsensusLevel:      consensusLevel,
	}

	return result, nil
}

// ValidateModels valida modelos usando dados históricos
func (f *Forecaster) ValidateModels(data []predicoes.HistoricalVRData, sindicato string) ([]ForecastValidation, error) {
	if len(data) < f.config.ValidationPeriod+3 {
		return nil, fmt.Errorf("dados insuficientes para validação")
	}

	sindicatoData := f.filterBySindicato(data, sindicato)
	if len(sindicatoData) < f.config.ValidationPeriod+3 {
		return nil, fmt.Errorf("dados insuficientes para validação do sindicato %s", sindicato)
	}

	// Dividir dados: treino e teste
	splitPoint := len(sindicatoData) - f.config.ValidationPeriod
	trainingData := sindicatoData[:splitPoint]
	testData := sindicatoData[splitPoint:]

	validations := make([]ForecastValidation, 0)

	for name, model := range f.models {
		validation := f.validateModel(model, trainingData, testData)
		validation.Model = name
		validations = append(validations, validation)
	}

	return validations, nil
}

// GenerateScenarios gera cenários de previsão
func (f *Forecaster) GenerateScenarios(data []predicoes.HistoricalVRData, sindicato string) (*ForecastScenarios, error) {
	// Cenário base
	baseForecast, err := f.ForecastConsumption(data, sindicato, f.config.DefaultHorizon)
	if err != nil {
		return nil, err
	}

	scenarios := &ForecastScenarios{
		Base:         baseForecast.WeightedForecast,
		Optimistic:   f.generateOptimisticScenario(baseForecast.WeightedForecast),
		Pessimistic:  f.generatePessimisticScenario(baseForecast.WeightedForecast),
		Conservative: f.generateConservativeScenario(baseForecast.WeightedForecast),
	}

	return scenarios, nil
}

// Implementação dos modelos

// SimpleMovingAverageModel

func NewSimpleMovingAverageModel(window int, sindicato string) *SimpleMovingAverageModel {
	return &SimpleMovingAverageModel{
		name:      fmt.Sprintf("sma_%d", window),
		window:    window,
		sindicato: sindicato,
	}
}

func (sma *SimpleMovingAverageModel) Train(data []predicoes.HistoricalVRData) error {
	if len(data) < sma.window {
		return fmt.Errorf("dados insuficientes para janela de %d", sma.window)
	}

	// Extrair valores
	values := make([]float64, len(data))
	for i, d := range data {
		values[i] = d.TotalVR
	}

	// Guardar últimos valores para predição
	sma.lastValues = values[len(values)-sma.window:]

	// Calcular acurácia usando validação cruzada simples
	sma.accuracy = sma.calculateAccuracy(values)

	return nil
}

func (sma *SimpleMovingAverageModel) Predict(horizon int) (*predicoes.ConsumptionForecast, error) {
	if len(sma.lastValues) == 0 {
		return nil, fmt.Errorf("modelo não foi treinado")
	}

	// Calcular média móvel
	sum := 0.0
	for _, v := range sma.lastValues {
		sum += v
	}
	predictedValue := sum / float64(len(sma.lastValues))

	// Próximo mês
	nextMonth := time.Now().AddDate(0, 1, 0)
	nextMonth = time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())

	forecast := &predicoes.ConsumptionForecast{
		Sindicato:   sma.sindicato,
		Month:       nextMonth,
		PredictedVR: predictedValue,
		Confidence:  sma.accuracy,
		Range: predicoes.ForecastRange{
			Lower: predictedValue * 0.9, // ±10%
			Upper: predictedValue * 1.1,
		},
		Factors: []string{
			fmt.Sprintf("Média móvel de %d períodos", sma.window),
			fmt.Sprintf("Acurácia: %.1f%%", sma.accuracy*100),
		},
		Assumptions: []string{
			"Comportamento futuro similar à média recente",
			"Não há mudanças estruturais",
		},
	}

	return forecast, nil
}

func (sma *SimpleMovingAverageModel) GetAccuracy() float64 {
	return sma.accuracy
}

func (sma *SimpleMovingAverageModel) GetName() string {
	return sma.name
}

func (sma *SimpleMovingAverageModel) calculateAccuracy(values []float64) float64 {
	if len(values) <= sma.window {
		return 0.5
	}

	// Validação usando últimos períodos
	errors := make([]float64, 0)
	for i := sma.window; i < len(values); i++ {
		// Calcular média móvel para posição i
		sum := 0.0
		for j := i - sma.window; j < i; j++ {
			sum += values[j]
		}
		predicted := sum / float64(sma.window)
		actual := values[i]

		if actual != 0 {
			error := math.Abs((predicted - actual) / actual)
			errors = append(errors, error)
		}
	}

	if len(errors) == 0 {
		return 0.5
	}

	// Calcular MAPE
	mape := 0.0
	for _, e := range errors {
		mape += e
	}
	mape /= float64(len(errors))

	// Converter MAPE em acurácia
	accuracy := 1.0 - math.Min(mape, 1.0)
	return math.Max(accuracy, 0.1)
}

// ExponentialSmoothingModel

func NewExponentialSmoothingModel(alpha float64, sindicato string) *ExponentialSmoothingModel {
	return &ExponentialSmoothingModel{
		name:      "exponential_smoothing",
		alpha:     alpha,
		seasonal:  make(map[int]float64),
		sindicato: sindicato,
	}
}

func (esm *ExponentialSmoothingModel) Train(data []predicoes.HistoricalVRData) error {
	if len(data) < 2 {
		return fmt.Errorf("dados insuficientes")
	}

	values := make([]float64, len(data))
	for i, d := range data {
		values[i] = d.TotalVR
	}

	// Inicializar nível e tendência
	esm.level = values[0]
	if len(values) > 1 {
		esm.trend = values[1] - values[0]
	}

	// Aplicar suavização exponencial
	for i := 1; i < len(values); i++ {
		prevLevel := esm.level
		esm.level = esm.alpha*values[i] + (1-esm.alpha)*(esm.level+esm.trend)
		esm.trend = esm.alpha*(esm.level-prevLevel) + (1-esm.alpha)*esm.trend
	}

	// Calcular componentes sazonais se há dados suficientes
	if len(data) >= 12 {
		esm.calculateSeasonalComponents(data)
	}

	// Calcular acurácia
	esm.accuracy = esm.calculateAccuracy(values)

	return nil
}

func (esm *ExponentialSmoothingModel) Predict(horizon int) (*predicoes.ConsumptionForecast, error) {
	if esm.level == 0 {
		return nil, fmt.Errorf("modelo não foi treinado")
	}

	// Predição usando nível + tendência
	predictedValue := esm.level + esm.trend

	// Ajustar por sazonalidade se disponível
	nextMonth := time.Now().AddDate(0, 1, 0)
	if seasonal, exists := esm.seasonal[int(nextMonth.Month())]; exists {
		predictedValue *= seasonal
	}

	nextMonth = time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())

	forecast := &predicoes.ConsumptionForecast{
		Sindicato:   esm.sindicato,
		Month:       nextMonth,
		PredictedVR: predictedValue,
		Confidence:  esm.accuracy,
		Range: predicoes.ForecastRange{
			Lower: predictedValue * 0.85, // ±15%
			Upper: predictedValue * 1.15,
		},
		Factors: []string{
			"Suavização exponencial com tendência",
			fmt.Sprintf("Alpha: %.2f", esm.alpha),
		},
		Assumptions: []string{
			"Tendências recentes continuam",
			"Padrões sazonais se mantêm",
		},
	}

	return forecast, nil
}

func (esm *ExponentialSmoothingModel) GetAccuracy() float64 {
	return esm.accuracy
}

func (esm *ExponentialSmoothingModel) GetName() string {
	return esm.name
}

func (esm *ExponentialSmoothingModel) calculateSeasonalComponents(data []predicoes.HistoricalVRData) {
	monthlyData := make(map[int][]float64)

	for _, d := range data {
		month := int(d.Month.Month())
		monthlyData[month] = append(monthlyData[month], d.TotalVR)
	}

	// Calcular fator sazonal para cada mês
	overallSum := 0.0
	overallCount := 0
	for _, values := range monthlyData {
		for _, v := range values {
			overallSum += v
			overallCount++
		}
	}
	overallAvg := overallSum / float64(overallCount)

	for month, values := range monthlyData {
		if len(values) > 0 {
			monthSum := 0.0
			for _, v := range values {
				monthSum += v
			}
			monthAvg := monthSum / float64(len(values))
			esm.seasonal[month] = monthAvg / overallAvg
		}
	}
}

func (esm *ExponentialSmoothingModel) calculateAccuracy(values []float64) float64 {
	// Implementação simplificada - na prática usaria validação cruzada
	return 0.75 // placeholder
}

// LinearRegressionModel

func NewLinearRegressionModel(sindicato string) *LinearRegressionModel {
	return &LinearRegressionModel{
		name:      "linear_regression",
		sindicato: sindicato,
	}
}

func (lrm *LinearRegressionModel) Train(data []predicoes.HistoricalVRData) error {
	if len(data) < 2 {
		return fmt.Errorf("dados insuficientes")
	}

	values := make([]float64, len(data))
	for i, d := range data {
		values[i] = d.TotalVR
	}

	// Calcular regressão linear
	n := float64(len(values))
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0

	for i, y := range values {
		x := float64(i + 1)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	lrm.slope = (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	lrm.intercept = (sumY - lrm.slope*sumX) / n

	// Calcular R² como medida de acurácia
	meanY := sumY / n
	ssRes, ssTot := 0.0, 0.0

	for i, y := range values {
		x := float64(i + 1)
		predicted := lrm.intercept + lrm.slope*x
		ssRes += (y - predicted) * (y - predicted)
		ssTot += (y - meanY) * (y - meanY)
	}

	if ssTot > 0 {
		lrm.accuracy = 1.0 - (ssRes / ssTot)
		lrm.accuracy = math.Max(lrm.accuracy, 0.1)
	} else {
		lrm.accuracy = 0.1
	}

	return nil
}

func (lrm *LinearRegressionModel) Predict(horizon int) (*predicoes.ConsumptionForecast, error) {
	if lrm.accuracy == 0 {
		return nil, fmt.Errorf("modelo não foi treinado")
	}

	// Predição para próximo período (assumindo sequência contínua)
	nextX := float64(100) // valor placeholder
	predictedValue := lrm.intercept + lrm.slope*nextX

	// Garantir valor positivo
	predictedValue = math.Max(predictedValue, 0)

	nextMonth := time.Now().AddDate(0, 1, 0)
	nextMonth = time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())

	forecast := &predicoes.ConsumptionForecast{
		Sindicato:   lrm.sindicato,
		Month:       nextMonth,
		PredictedVR: predictedValue,
		Confidence:  lrm.accuracy,
		Range: predicoes.ForecastRange{
			Lower: predictedValue * 0.8, // ±20%
			Upper: predictedValue * 1.2,
		},
		Factors: []string{
			"Regressão linear",
			fmt.Sprintf("Slope: %.3f", lrm.slope),
		},
		Assumptions: []string{
			"Relação linear continua",
			"Não há mudanças estruturais",
		},
	}

	return forecast, nil
}

func (lrm *LinearRegressionModel) GetAccuracy() float64 {
	return lrm.accuracy
}

func (lrm *LinearRegressionModel) GetName() string {
	return lrm.name
}

// Métodos auxiliares do Forecaster

func (f *Forecaster) filterBySindicato(data []predicoes.HistoricalVRData, sindicato string) []predicoes.HistoricalVRData {
	filtered := make([]predicoes.HistoricalVRData, 0)
	for _, d := range data {
		if d.Sindicato == sindicato {
			filtered = append(filtered, d)
		}
	}

	// Ordenar por data
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Month.Before(filtered[j].Month)
	})

	return filtered
}

func (f *Forecaster) ensureDefaultModels(sindicato string) {
	if len(f.models) == 0 {
		// Registrar modelos padrão
		f.RegisterModel(NewSimpleMovingAverageModel(3, sindicato))
		f.RegisterModel(NewSimpleMovingAverageModel(6, sindicato))
		f.RegisterModel(NewExponentialSmoothingModel(0.3, sindicato))
		f.RegisterModel(NewLinearRegressionModel(sindicato))
	}
}

func (f *Forecaster) calculateEnsembleForecast(forecasts map[string]*predicoes.ConsumptionForecast, weights map[string]float64, sindicato string, horizon int) *predicoes.ConsumptionForecast {
	if len(forecasts) == 0 {
		return nil
	}

	// Normalizar pesos
	totalWeight := 0.0
	for _, w := range weights {
		totalWeight += w
	}

	if totalWeight == 0 {
		return nil
	}

	// Calcular média ponderada
	weightedSum := 0.0
	weightedConfidence := 0.0
	lowerSum := 0.0
	upperSum := 0.0

	for name, forecast := range forecasts {
		weight := weights[name] / totalWeight
		weightedSum += forecast.PredictedVR * weight
		weightedConfidence += forecast.Confidence * weight
		lowerSum += forecast.Range.Lower * weight
		upperSum += forecast.Range.Upper * weight
	}

	// Coletar fatores únicos
	factorSet := make(map[string]bool)
	for _, forecast := range forecasts {
		for _, factor := range forecast.Factors {
			factorSet[factor] = true
		}
	}

	factors := make([]string, 0, len(factorSet))
	for factor := range factorSet {
		factors = append(factors, factor)
	}

	// Próximo mês
	nextMonth := time.Now().AddDate(0, 1, 0)
	nextMonth = time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())

	ensembleForecast := &predicoes.ConsumptionForecast{
		Sindicato:   sindicato,
		Month:       nextMonth,
		PredictedVR: weightedSum,
		Confidence:  weightedConfidence,
		Range: predicoes.ForecastRange{
			Lower: lowerSum,
			Upper: upperSum,
		},
		Factors: factors,
		Assumptions: []string{
			fmt.Sprintf("Ensemble de %d modelos", len(forecasts)),
			"Pesos baseados na acurácia individual",
		},
	}

	return ensembleForecast
}

func (f *Forecaster) calculateConsensusLevel(forecasts map[string]*predicoes.ConsumptionForecast) float64 {
	if len(forecasts) <= 1 {
		return 1.0
	}

	values := make([]float64, 0, len(forecasts))
	for _, forecast := range forecasts {
		values = append(values, forecast.PredictedVR)
	}

	// Calcular coeficiente de variação
	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))

	if mean == 0 {
		return 0.5
	}

	cv := math.Sqrt(variance) / mean
	consensus := 1.0 / (1.0 + cv) // inverso do coeficiente de variação

	return math.Max(0, math.Min(consensus, 1.0))
}

func (f *Forecaster) calculateEnsembleAccuracy(weights map[string]float64) float64 {
	totalWeight := 0.0
	weightedAccuracy := 0.0

	for _, weight := range weights {
		totalWeight += weight
		weightedAccuracy += weight * weight // peso quadrático para dar mais importância a modelos melhores
	}

	if totalWeight == 0 {
		return 0.5
	}

	return weightedAccuracy / totalWeight
}

func (f *Forecaster) validateModel(model ForecastModel, trainingData, testData []predicoes.HistoricalVRData) ForecastValidation {
	// Treinar modelo
	err := model.Train(trainingData)
	if err != nil {
		return ForecastValidation{
			Accuracy: 0.0,
			MAPE:     100.0,
		}
	}

	// Fazer predições e comparar
	predictions := make([]ValidationPoint, 0)
	errors := make([]float64, 0)
	percentErrors := make([]float64, 0)

	for _, actual := range testData {
		// Para simplificar, usar uma predição por ponto
		forecast, err := model.Predict(1)
		if err != nil {
			continue
		}

		predicted := forecast.PredictedVR
		error := predicted - actual.TotalVR
		errorPct := 0.0
		if actual.TotalVR != 0 {
			errorPct = math.Abs(error/actual.TotalVR) * 100
		}

		predictions = append(predictions, ValidationPoint{
			Date:      actual.Month,
			Actual:    actual.TotalVR,
			Predicted: predicted,
			Error:     error,
			ErrorPct:  errorPct,
		})

		errors = append(errors, error)
		percentErrors = append(percentErrors, errorPct)
	}

	if len(errors) == 0 {
		return ForecastValidation{
			Accuracy: 0.0,
			MAPE:     100.0,
		}
	}

	// Calcular métricas
	mae := f.calculateMAE(errors)
	rmse := f.calculateRMSE(errors)
	mape := f.calculateMAPE(percentErrors)
	accuracy := math.Max(0, 1.0-mape/100.0)

	return ForecastValidation{
		Accuracy:    accuracy,
		MAPE:        mape,
		RMSE:        rmse,
		MAE:         mae,
		Predictions: predictions,
	}
}

func (f *Forecaster) calculateMAE(errors []float64) float64 {
	sum := 0.0
	for _, e := range errors {
		sum += math.Abs(e)
	}
	return sum / float64(len(errors))
}

func (f *Forecaster) calculateRMSE(errors []float64) float64 {
	sum := 0.0
	for _, e := range errors {
		sum += e * e
	}
	return math.Sqrt(sum / float64(len(errors)))
}

func (f *Forecaster) calculateMAPE(percentErrors []float64) float64 {
	sum := 0.0
	for _, e := range percentErrors {
		sum += e
	}
	return sum / float64(len(percentErrors))
}

// Estruturas para cenários

type ForecastScenarios struct {
	Base         *predicoes.ConsumptionForecast `json:"base"`
	Optimistic   *predicoes.ConsumptionForecast `json:"optimistic"`
	Pessimistic  *predicoes.ConsumptionForecast `json:"pessimistic"`
	Conservative *predicoes.ConsumptionForecast `json:"conservative"`
}

func (f *Forecaster) generateOptimisticScenario(baseForecast *predicoes.ConsumptionForecast) *predicoes.ConsumptionForecast {
	optimistic := *baseForecast
	optimistic.PredictedVR *= 1.15 // +15%
	optimistic.Range.Lower *= 1.10
	optimistic.Range.Upper *= 1.20
	optimistic.Assumptions = append(optimistic.Assumptions, "Cenário otimista: +15%")
	return &optimistic
}

func (f *Forecaster) generatePessimisticScenario(baseForecast *predicoes.ConsumptionForecast) *predicoes.ConsumptionForecast {
	pessimistic := *baseForecast
	pessimistic.PredictedVR *= 0.85 // -15%
	pessimistic.Range.Lower *= 0.80
	pessimistic.Range.Upper *= 0.90
	pessimistic.Assumptions = append(pessimistic.Assumptions, "Cenário pessimista: -15%")
	return &pessimistic
}

func (f *Forecaster) generateConservativeScenario(baseForecast *predicoes.ConsumptionForecast) *predicoes.ConsumptionForecast {
	conservative := *baseForecast
	conservative.PredictedVR *= 0.95 // -5% (conservador)
	conservative.Range.Lower *= 0.90
	conservative.Range.Upper *= 1.00
	conservative.Assumptions = append(conservative.Assumptions, "Cenário conservador: -5%")
	return &conservative
}
