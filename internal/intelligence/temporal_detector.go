package intelligence

import (
	"fmt"
	"time"
)

// TemporalAnomalyDetector detecta anomalias temporais
type TemporalAnomalyDetector struct {
	config *AnalysisConfig
}

// NewTemporalAnomalyDetector cria novo detector temporal
func NewTemporalAnomalyDetector(config *AnalysisConfig) *TemporalAnomalyDetector {
	return &TemporalAnomalyDetector{
		config: config,
	}
}

// DetectFutureDates detecta datas futuras inválidas
func (d *TemporalAnomalyDetector) DetectFutureDates(ctx *AnalysisContext) []Anomaly {
	anomalies := make([]Anomaly, 0)
	now := time.Now()
	
	for matricula, data := range ctx.Colaboradores {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Verificar data de admissão
		if admissaoDate, err := d.extractDate(dataMap, "data_admissao"); err == nil {
			if admissaoDate.After(now) {
				anomaly := NewAnomaly(
					AnomalyTypeTemporal,
					SeverityCritical,
					95.0,
					"Data de admissão no futuro",
					fmt.Sprintf("Colaborador %s tem data de admissão %s, que está no futuro",
						matricula, admissaoDate.Format("02/01/2006")),
					matricula,
					admissaoDate.Format("2006-01-02"),
					"data_admissao",
					"temporal_future_detector",
					"future_dates",
				)
				
				anomaly.AddData("data_admissao", admissaoDate)
				anomaly.AddData("data_atual", now)
				anomaly.AddData("dias_no_futuro", int(admissaoDate.Sub(now).Hours()/24))
				
				anomaly.AddSuggestion("Corrigir data de admissão para data válida no passado")
				anomaly.AddSuggestion("Verificar se houve erro de digitação no ano")
				
				anomalies = append(anomalies, anomaly)
			}
		}
		
		// Verificar data de desligamento
		if desligamentoDate, err := d.extractDate(dataMap, "data_desligamento"); err == nil {
			if desligamentoDate.After(now) {
				anomaly := NewAnomaly(
					AnomalyTypeTemporal,
					SeverityHigh,
					90.0,
					"Data de desligamento no futuro",
					fmt.Sprintf("Colaborador %s tem data de desligamento %s, que está no futuro",
						matricula, desligamentoDate.Format("02/01/2006")),
					matricula,
					desligamentoDate.Format("2006-01-02"),
					"data_desligamento",
					"temporal_future_detector",
					"future_dates",
				)
				
				anomaly.AddData("data_desligamento", desligamentoDate)
				anomaly.AddData("data_atual", now)
				anomaly.AddData("dias_no_futuro", int(desligamentoDate.Sub(now).Hours()/24))
				
				anomaly.AddSuggestion("Corrigir data de desligamento")
				anomaly.AddSuggestion("Verificar se é uma previsão que deveria estar marcada diferentemente")
				
				anomalies = append(anomalies, anomaly)
			}
		}
		
		ctx.IncrementProcessed()
	}
	
	return anomalies
}

// DetectInvalidDateSequences detecta sequências de datas inválidas
func (d *TemporalAnomalyDetector) DetectInvalidDateSequences(ctx *AnalysisContext) []Anomaly {
	anomalies := make([]Anomaly, 0)
	
	for matricula, data := range ctx.Colaboradores {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Extrair datas principais
		admissaoDate, errAdmissao := d.extractDate(dataMap, "data_admissao")
		desligamentoDate, errDesligamento := d.extractDate(dataMap, "data_desligamento")
		
		// Verificar sequência admissão -> desligamento
		if errAdmissao == nil && errDesligamento == nil {
			if desligamentoDate.Before(admissaoDate) {
				anomaly := NewAnomaly(
					AnomalyTypeTemporal,
					SeverityCritical,
					95.0,
					"Data de desligamento antes da admissão",
					fmt.Sprintf("Colaborador %s tem data de desligamento (%s) anterior à data de admissão (%s)",
						matricula, desligamentoDate.Format("02/01/2006"), admissaoDate.Format("02/01/2006")),
					matricula,
					fmt.Sprintf("%s > %s", desligamentoDate.Format("2006-01-02"), admissaoDate.Format("2006-01-02")),
					"sequencia_datas",
					"temporal_sequence_detector",
					"invalid_date_sequence",
				)
				
				anomaly.AddData("data_admissao", admissaoDate)
				anomaly.AddData("data_desligamento", desligamentoDate)
				anomaly.AddData("diferenca_dias", int(admissaoDate.Sub(desligamentoDate).Hours()/24))
				
				anomaly.AddSuggestion("Corrigir sequência de datas")
				anomaly.AddSuggestion("Verificar se as datas foram invertidas")
				
				anomalies = append(anomalies, anomaly)
			}
		}
		
		// Verificar datas muito antigas
		if errAdmissao == nil {
			if admissaoDate.Before(d.config.MinAdmissionDate) {
				anomaly := NewAnomaly(
					AnomalyTypeTemporal,
					SeverityMedium,
					80.0,
					"Data de admissão muito antiga",
					fmt.Sprintf("Colaborador %s tem data de admissão muito antiga: %s",
						matricula, admissaoDate.Format("02/01/2006")),
					matricula,
					admissaoDate.Format("2006-01-02"),
					"data_admissao",
					"temporal_ancient_detector",
					"ancient_date",
				)
				
				anomaly.AddData("data_admissao", admissaoDate)
				anomaly.AddData("data_minima", d.config.MinAdmissionDate)
				anomaly.AddData("anos_atras", int(time.Since(admissaoDate).Hours()/24/365))
				
				anomaly.AddSuggestion("Verificar se a data está correta")
				anomaly.AddSuggestion("Confirmar se colaborador ainda está ativo")
				
				anomalies = append(anomalies, anomaly)
			}
		}
		
		// Detectar períodos de férias/afastamentos problemáticos
		anomalies = append(anomalies, d.detectPeriodAnomalies(matricula, dataMap)...)
		
		ctx.IncrementProcessed()
	}
	
	return anomalies
}

// DetectWorkDayAnomalies detecta anomalias nos dias úteis
func (d *TemporalAnomalyDetector) DetectWorkDayAnomalies(ctx *AnalysisContext) []Anomaly {
	anomalies := make([]Anomaly, 0)
	
	for matricula, data := range ctx.Colaboradores {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Extrair dias úteis
		diasUteis, err := d.extractNumericField(dataMap, "dias_uteis")
		if err != nil {
			continue
		}
		
		// Verificar limites de dias úteis
		if diasUteis > float64(d.config.MaxWorkDaysPerMonth) {
			anomaly := NewAnomaly(
				AnomalyTypeTemporal,
				SeverityHigh,
				90.0,
				fmt.Sprintf("Dias úteis excessivos: %.0f", diasUteis),
				fmt.Sprintf("Colaborador %s tem %.0f dias úteis, acima do máximo de %d",
					matricula, diasUteis, d.config.MaxWorkDaysPerMonth),
				matricula,
				fmt.Sprintf("%.0f", diasUteis),
				"dias_uteis",
				"temporal_workday_detector",
				"excessive_workdays",
			)
			
			anomaly.AddData("dias_uteis", diasUteis)
			anomaly.AddData("maximo_permitido", d.config.MaxWorkDaysPerMonth)
			anomaly.AddData("excesso", diasUteis-float64(d.config.MaxWorkDaysPerMonth))
			
			anomaly.AddSuggestion("Verificar cálculo de dias úteis")
			anomaly.AddSuggestion("Confirmar período de referência")
			anomaly.AddSuggestion("Revisar feriados e finais de semana")
			
			anomalies = append(anomalies, anomaly)
		}
		
		if diasUteis < 0 {
			anomaly := NewAnomaly(
				AnomalyTypeTemporal,
				SeverityCritical,
				95.0,
				"Dias úteis negativos",
				fmt.Sprintf("Colaborador %s tem dias úteis negativos: %.0f", matricula, diasUteis),
				matricula,
				fmt.Sprintf("%.0f", diasUteis),
				"dias_uteis",
				"temporal_workday_detector",
				"negative_workdays",
			)
			
			anomaly.AddData("dias_uteis", diasUteis)
			
			anomaly.AddSuggestion("Corrigir erro de cálculo que resultou em valor negativo")
			anomaly.AddSuggestion("Verificar fórmulas de cálculo de dias úteis")
			
			anomalies = append(anomalies, anomaly)
		}
		
		ctx.IncrementProcessed()
	}
	
	return anomalies
}

// detectPeriodAnomalies detecta anomalias em períodos de férias/afastamentos
func (d *TemporalAnomalyDetector) detectPeriodAnomalies(matricula string, dataMap map[string]interface{}) []Anomaly {
	anomalies := make([]Anomaly, 0)
	
	// Extrair data de admissão para validar períodos
	admissaoDate, errAdmissao := d.extractDate(dataMap, "data_admissao")
	if errAdmissao != nil {
		return anomalies
	}
	
	// Verificar período de férias se disponível
	if feriasInicio, err1 := d.extractDate(dataMap, "ferias_inicio"); err1 == nil {
		if feriasFim, err2 := d.extractDate(dataMap, "ferias_fim"); err2 == nil {
			// Verificar se férias começam antes da admissão
			if feriasInicio.Before(admissaoDate) {
				anomaly := NewAnomaly(
					AnomalyTypeTemporal,
					SeverityHigh,
					90.0,
					"Férias antes da admissão",
					fmt.Sprintf("Colaborador %s tem férias (%s a %s) que começam antes da admissão (%s)",
						matricula, feriasInicio.Format("02/01/2006"), feriasFim.Format("02/01/2006"), admissaoDate.Format("02/01/2006")),
					matricula,
					fmt.Sprintf("%s-%s", feriasInicio.Format("2006-01-02"), feriasFim.Format("2006-01-02")),
					"periodo_ferias",
					"temporal_period_detector",
					"vacation_before_admission",
				)
				
				anomaly.AddData("ferias_inicio", feriasInicio)
				anomaly.AddData("ferias_fim", feriasFim)
				anomaly.AddData("data_admissao", admissaoDate)
				
				anomaly.AddSuggestion("Corrigir datas de férias")
				anomaly.AddSuggestion("Verificar se são férias prêmio de emprego anterior")
				
				anomalies = append(anomalies, anomaly)
			}
			
			// Verificar duração das férias
			duracaoFerias := feriasFim.Sub(feriasInicio).Hours() / 24
			tempoTrabalhado := time.Since(admissaoDate).Hours() / 24
			
			if duracaoFerias > 30 && tempoTrabalhado < 365 {
				anomaly := NewAnomaly(
					AnomalyTypeTemporal,
					SeverityMedium,
					75.0,
					fmt.Sprintf("Férias longas para novo colaborador: %.0f dias", duracaoFerias),
					fmt.Sprintf("Colaborador %s tem %.0f dias de férias, mas trabalha há apenas %.0f dias",
						matricula, duracaoFerias, tempoTrabalhado),
					matricula,
					fmt.Sprintf("%.0f", duracaoFerias),
					"duracao_ferias",
					"temporal_period_detector",
					"excessive_vacation_new_employee",
				)
				
				anomaly.AddData("duracao_ferias", duracaoFerias)
				anomaly.AddData("tempo_trabalhado", tempoTrabalhado)
				anomaly.AddData("data_admissao", admissaoDate)
				
				anomaly.AddSuggestion("Verificar se colaborador tem direito a tantos dias de férias")
				anomaly.AddSuggestion("Confirmar cálculo proporcional de férias")
				
				anomalies = append(anomalies, anomaly)
			}
		}
	}
	
	return anomalies
}

// extractDate extrai data de um campo específico
func (d *TemporalAnomalyDetector) extractDate(dataMap map[string]interface{}, fieldName string) (time.Time, error) {
	value, exists := dataMap[fieldName]
	if !exists {
		return time.Time{}, fmt.Errorf("campo %s não encontrado", fieldName)
	}
	
	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		// Tentar vários formatos de data
		formats := []string{
			"2006-01-02",
			"02/01/2006",
			"2006-01-02 15:04:05",
			"02/01/2006 15:04:05",
			"2006/01/02",
			"01-02-2006",
		}
		
		for _, format := range formats {
			if date, err := time.Parse(format, v); err == nil {
				return date, nil
			}
		}
		
		return time.Time{}, fmt.Errorf("formato de data inválido: %s", v)
	default:
		return time.Time{}, fmt.Errorf("tipo de data inválido: %T", value)
	}
}

// extractNumericField extrai valor numérico de um campo
func (d *TemporalAnomalyDetector) extractNumericField(dataMap map[string]interface{}, fieldName string) (float64, error) {
	value, exists := dataMap[fieldName]
	if !exists {
		return 0, fmt.Errorf("campo %s não encontrado", fieldName)
	}
	
	return extractNumericValue(value)
}

// IsBusinessDay verifica se uma data é dia útil (segunda a sexta, excluindo feriados)
func (d *TemporalAnomalyDetector) IsBusinessDay(date time.Time) bool {
	// Verificar se é final de semana
	weekday := date.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	
	// TODO: Implementar verificação de feriados nacionais
	// Por enquanto, apenas considera segunda a sexta como dias úteis
	
	return true
}

// CalculateBusinessDays calcula número de dias úteis entre duas datas
func (d *TemporalAnomalyDetector) CalculateBusinessDays(start, end time.Time) int {
	if start.After(end) {
		return 0
	}
	
	businessDays := 0
	current := start
	
	for current.Before(end) || current.Equal(end) {
		if d.IsBusinessDay(current) {
			businessDays++
		}
		current = current.AddDate(0, 0, 1)
	}
	
	return businessDays
}

// GetQuarter retorna o trimestre de uma data (1-4)
func (d *TemporalAnomalyDetector) GetQuarter(date time.Time) int {
	month := int(date.Month())
	return ((month - 1) / 3) + 1
}

// IsValidWorkPeriod verifica se um período de trabalho é válido
func (d *TemporalAnomalyDetector) IsValidWorkPeriod(start, end time.Time, expectedBusinessDays int) bool {
	actualBusinessDays := d.CalculateBusinessDays(start, end)
	tolerance := 2 // 2 dias de tolerância
	
	return abs(actualBusinessDays-expectedBusinessDays) <= tolerance
}

// Helper function para valor absoluto
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}