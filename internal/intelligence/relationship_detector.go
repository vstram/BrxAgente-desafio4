package intelligence

import (
	"fmt"
	"strings"
)

// RelationshipAnomalyDetector detecta anomalias de relacionamento entre dados
type RelationshipAnomalyDetector struct {
	config *AnalysisConfig
}

// NewRelationshipAnomalyDetector cria novo detector de relacionamentos
func NewRelationshipAnomalyDetector(config *AnalysisConfig) *RelationshipAnomalyDetector {
	return &RelationshipAnomalyDetector{
		config: config,
	}
}

// DetectDuplicateMatriculas detecta matrículas duplicadas
func (d *RelationshipAnomalyDetector) DetectDuplicateMatriculas(ctx *AnalysisContext) []Anomaly {
	anomalies := make([]Anomaly, 0)
	
	// Mapa para rastrear matrículas já vistas
	matriculasVistas := make(map[string][]string)
	
	// Primeiro passo: coletar todas as matrículas e suas fontes
	for _, data := range ctx.Colaboradores {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Extrair matrícula dos dados (não da chave)
		matriculaValue, exists := dataMap["matricula"]
		if !exists {
			continue
		}
		
		matriculaStr, ok := matriculaValue.(string)
		if !ok {
			continue
		}
		
		// Normalizar matrícula (remover espaços, zeros à esquerda, etc.)
		normalizedMatricula := d.normalizeMatricula(matriculaStr)
		
		// Extrair fonte dos dados se disponível
		fonte := "unknown"
		if f, exists := dataMap["fonte"]; exists {
			if fStr, ok := f.(string); ok {
				fonte = fStr
			}
		}
		
		// Adicionar ao mapa de matrículas vistas
		if fontes, exists := matriculasVistas[normalizedMatricula]; exists {
			matriculasVistas[normalizedMatricula] = append(fontes, fonte)
		} else {
			matriculasVistas[normalizedMatricula] = []string{fonte}
		}
		
		ctx.IncrementProcessed()
	}
	
	// Segundo passo: identificar duplicatas
	for matricula, fontes := range matriculasVistas {
		if len(fontes) > 1 {
			anomaly := NewAnomaly(
				AnomalyTypeRelationship,
				SeverityCritical,
				95.0,
				fmt.Sprintf("Matrícula duplicada: %s", matricula),
				fmt.Sprintf("Matrícula %s aparece em múltiplas fontes: %s",
					matricula, strings.Join(fontes, ", ")),
				matricula,
				matricula,
				"matricula",
				"relationship_duplicate_detector",
				"duplicate_matricula",
			)
			
			anomaly.AddData("fontes", fontes)
			anomaly.AddData("numero_duplicatas", len(fontes))
			
			anomaly.AddSuggestion("Consolidar dados da matrícula em uma única fonte")
			anomaly.AddSuggestion("Verificar se são colaboradores diferentes com matrícula similar")
			anomaly.AddSuggestion("Revisar processo de importação de dados")
			
			anomalies = append(anomalies, anomaly)
		}
	}
	
	return anomalies
}

// DetectInconsistentStatus detecta inconsistências de status
func (d *RelationshipAnomalyDetector) DetectInconsistentStatus(ctx *AnalysisContext) []Anomaly {
	anomalies := make([]Anomaly, 0)
	
	for matricula, data := range ctx.Colaboradores {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Extrair status e data de desligamento
		status, hasStatus := dataMap["status"]
		dataDesligamento, hasDesligamento := dataMap["data_desligamento"]
		
		if hasStatus && hasDesligamento {
			statusStr, statusOk := status.(string)
			if !statusOk {
				continue
			}
			
			// Verificar inconsistências
			if strings.ToLower(statusStr) == "ativo" && dataDesligamento != nil {
				// Verificar se data de desligamento não é vazia/zero
				if dateStr, ok := dataDesligamento.(string); ok && dateStr != "" && dateStr != "0000-00-00" {
					anomaly := NewAnomaly(
						AnomalyTypeRelationship,
						SeverityHigh,
						90.0,
						"Status ativo com data de desligamento",
						fmt.Sprintf("Colaborador %s está marcado como ativo mas tem data de desligamento: %v",
							matricula, dataDesligamento),
						matricula,
						statusStr,
						"status",
						"relationship_status_detector",
						"inconsistent_status",
					)
					
					anomaly.AddData("status", status)
					anomaly.AddData("data_desligamento", dataDesligamento)
					
					anomaly.AddSuggestion("Atualizar status para 'desligado' ou remover data de desligamento")
					anomaly.AddSuggestion("Verificar se houve recontratação")
					
					anomalies = append(anomalies, anomaly)
				}
			} else if strings.ToLower(statusStr) == "desligado" && (dataDesligamento == nil || dataDesligamento == "") {
				anomaly := NewAnomaly(
					AnomalyTypeRelationship,
					SeverityHigh,
					90.0,
					"Status desligado sem data de desligamento",
					fmt.Sprintf("Colaborador %s está marcado como desligado mas não tem data de desligamento",
						matricula),
					matricula,
					statusStr,
					"status",
					"relationship_status_detector",
					"inconsistent_status",
				)
				
				anomaly.AddData("status", status)
				anomaly.AddData("data_desligamento", dataDesligamento)
				
				anomaly.AddSuggestion("Adicionar data de desligamento")
				anomaly.AddSuggestion("Verificar se status deveria ser 'ativo'")
				
				anomalies = append(anomalies, anomaly)
			}
		}
		
		ctx.IncrementProcessed()
	}
	
	return anomalies
}

// DetectOrphanedRecords detecta registros órfãos
func (d *RelationshipAnomalyDetector) DetectOrphanedRecords(ctx *AnalysisContext) []Anomaly {
	anomalies := make([]Anomaly, 0)
	
	// Extrair dados de referência se disponíveis no contexto
	sindicatosValidos := d.extractValidSindicatos(ctx)
	
	for matricula, data := range ctx.Colaboradores {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Verificar se sindicato existe
		if sindicato, exists := dataMap["sindicato"]; exists {
			if sindicatoStr, ok := sindicato.(string); ok {
				if len(sindicatosValidos) > 0 && !d.contains(sindicatosValidos, sindicatoStr) {
					anomaly := NewAnomaly(
						AnomalyTypeRelationship,
						SeverityMedium,
						80.0,
						fmt.Sprintf("Sindicato não reconhecido: %s", sindicatoStr),
						fmt.Sprintf("Colaborador %s tem sindicato '%s' que não está na lista de sindicatos válidos",
							matricula, sindicatoStr),
						matricula,
						sindicatoStr,
						"sindicato",
						"relationship_orphan_detector",
						"orphaned_sindicate",
					)
					
					anomaly.AddData("sindicato", sindicato)
					anomaly.AddData("sindicatos_validos", sindicatosValidos)
					
					anomaly.AddSuggestion("Verificar se sindicato foi digitado corretamente")
					anomaly.AddSuggestion("Atualizar lista de sindicatos válidos")
					anomaly.AddSuggestion("Corrigir sindicato do colaborador")
					
					anomalies = append(anomalies, anomaly)
				}
			}
		}
		
		// Verificar campos obrigatórios
		anomalies = append(anomalies, d.detectMissingRequiredFields(matricula, dataMap)...)
		
		ctx.IncrementProcessed()
	}
	
	return anomalies
}

// DetectDataInconsistencies detecta inconsistências entre campos relacionados
func (d *RelationshipAnomalyDetector) DetectDataInconsistencies(ctx *AnalysisContext) []Anomaly {
	anomalies := make([]Anomaly, 0)
	
	for matricula, data := range ctx.Colaboradores {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Verificar consistência entre VR calculado e valor esperado
		if vrCalculado, err1 := d.extractNumeric(dataMap, "vr_calculado"); err1 == nil {
			if vrEsperado, err2 := d.extractNumeric(dataMap, "vr_esperado"); err2 == nil {
				diferenca := abs64(vrCalculado - vrEsperado)
				tolerancia := vrEsperado * 0.05 // 5% de tolerância
				
				if diferenca > tolerancia {
					percentualDiferenca := (diferenca / vrEsperado) * 100
					
					anomaly := NewAnomaly(
						AnomalyTypeRelationship,
						d.getSeverityForDifference(percentualDiferenca),
						85.0,
						fmt.Sprintf("VR calculado difere do esperado: %.1f%%", percentualDiferenca),
						fmt.Sprintf("Colaborador %s tem VR calculado (R$ %.2f) diferente do esperado (R$ %.2f)",
							matricula, vrCalculado, vrEsperado),
						matricula,
						fmt.Sprintf("%.2f vs %.2f", vrCalculado, vrEsperado),
						"vr_consistencia",
						"relationship_consistency_detector",
						"vr_inconsistency",
					)
					
					anomaly.AddData("vr_calculado", vrCalculado)
					anomaly.AddData("vr_esperado", vrEsperado)
					anomaly.AddData("diferenca", diferenca)
					anomaly.AddData("percentual_diferenca", percentualDiferenca)
					
					anomaly.AddSuggestion("Revisar cálculo de VR")
					anomaly.AddSuggestion("Verificar parâmetros utilizados no cálculo")
					anomaly.AddSuggestion("Confirmar valor esperado")
					
					anomalies = append(anomalies, anomaly)
				}
			}
		}
		
		// Verificar consistência entre dias úteis e período
		anomalies = append(anomalies, d.detectPeriodConsistency(matricula, dataMap)...)
		
		ctx.IncrementProcessed()
	}
	
	return anomalies
}

// detectMissingRequiredFields detecta campos obrigatórios ausentes
func (d *RelationshipAnomalyDetector) detectMissingRequiredFields(matricula string, dataMap map[string]interface{}) []Anomaly {
	anomalies := make([]Anomaly, 0)
	
	requiredFields := []string{"matricula", "sindicato", "data_admissao"}
	
	for _, field := range requiredFields {
		if value, exists := dataMap[field]; !exists || d.isEmpty(value) {
			anomaly := NewAnomaly(
				AnomalyTypeRelationship,
				SeverityHigh,
				90.0,
				fmt.Sprintf("Campo obrigatório ausente: %s", field),
				fmt.Sprintf("Colaborador %s não tem valor para o campo obrigatório '%s'",
					matricula, field),
				matricula,
				"null",
				field,
				"relationship_required_detector",
				"missing_required_field",
			)
			
			anomaly.AddData("campo_ausente", field)
			anomaly.AddData("campos_obrigatorios", requiredFields)
			
			anomaly.AddSuggestion(fmt.Sprintf("Preencher campo obrigatório '%s'", field))
			anomaly.AddSuggestion("Verificar fonte dos dados")
			
			anomalies = append(anomalies, anomaly)
		}
	}
	
	return anomalies
}

// detectPeriodConsistency verifica consistência entre períodos e dias úteis
func (d *RelationshipAnomalyDetector) detectPeriodConsistency(matricula string, dataMap map[string]interface{}) []Anomaly {
	anomalies := make([]Anomaly, 0)
	
	// Extrair dados do período
	diasUteis, errDias := d.extractNumeric(dataMap, "dias_uteis")
	if errDias != nil {
		return anomalies
	}
	
	// Se temos informações de período, verificar consistência
	if periodoInicio, err1 := d.extractString(dataMap, "periodo_inicio"); err1 == nil {
		if periodoFim, err2 := d.extractString(dataMap, "periodo_fim"); err2 == nil {
			// Aqui poderíamos calcular dias úteis esperados para o período
			// Por simplicidade, vamos verificar apenas valores óbvios
			
			if diasUteis > 31 {
				anomaly := NewAnomaly(
					AnomalyTypeRelationship,
					SeverityMedium,
					75.0,
					fmt.Sprintf("Dias úteis inconsistentes: %.0f", diasUteis),
					fmt.Sprintf("Colaborador %s tem %.0f dias úteis no período %s a %s, que parece excessivo",
						matricula, diasUteis, periodoInicio, periodoFim),
					matricula,
					fmt.Sprintf("%.0f", diasUteis),
					"dias_uteis",
					"relationship_period_detector",
					"period_inconsistency",
				)
				
				anomaly.AddData("dias_uteis", diasUteis)
				anomaly.AddData("periodo_inicio", periodoInicio)
				anomaly.AddData("periodo_fim", periodoFim)
				
				anomaly.AddSuggestion("Verificar cálculo de dias úteis")
				anomaly.AddSuggestion("Confirmar período de referência")
				
				anomalies = append(anomalies, anomaly)
			}
		}
	}
	
	return anomalies
}

// Helper functions

func (d *RelationshipAnomalyDetector) normalizeMatricula(matricula string) string {
	// Remover espaços e converter para maiúsculas
	normalized := strings.TrimSpace(strings.ToUpper(matricula))
	
	// Remover zeros à esquerda se for numérico
	if len(normalized) > 1 && normalized[0] == '0' {
		// Verificar se é todo numérico
		isNumeric := true
		for _, r := range normalized {
			if r < '0' || r > '9' {
				isNumeric = false
				break
			}
		}
		
		if isNumeric {
			// Remover zeros à esquerda
			for len(normalized) > 1 && normalized[0] == '0' {
				normalized = normalized[1:]
			}
		}
	}
	
	return normalized
}

func (d *RelationshipAnomalyDetector) extractValidSindicatos(ctx *AnalysisContext) []string {
	// Tentar extrair lista de sindicatos válidos do contexto
	if sindicatos, exists := ctx.Parameters["sindicatos_validos"]; exists {
		if sindicatosList, ok := sindicatos.([]string); ok {
			return sindicatosList
		}
	}
	
	// Se não há lista específica, extrair sindicatos únicos dos dados
	sindicatosMap := make(map[string]bool)
	for _, data := range ctx.Colaboradores {
		if dataMap, ok := data.(map[string]interface{}); ok {
			if sindicato, exists := dataMap["sindicato"]; exists {
				if sindicatoStr, ok := sindicato.(string); ok && sindicatoStr != "" {
					sindicatosMap[sindicatoStr] = true
				}
			}
		}
	}
	
	// Converter para slice
	sindicatos := make([]string, 0, len(sindicatosMap))
	for sindicato := range sindicatosMap {
		sindicatos = append(sindicatos, sindicato)
	}
	
	return sindicatos
}

func (d *RelationshipAnomalyDetector) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (d *RelationshipAnomalyDetector) isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case int, int64, float64:
		return false // Números não são considerados vazios
	default:
		return false
	}
}

func (d *RelationshipAnomalyDetector) extractNumeric(dataMap map[string]interface{}, fieldName string) (float64, error) {
	value, exists := dataMap[fieldName]
	if !exists {
		return 0, fmt.Errorf("campo %s não encontrado", fieldName)
	}
	
	return extractNumericValue(value)
}

func (d *RelationshipAnomalyDetector) extractString(dataMap map[string]interface{}, fieldName string) (string, error) {
	value, exists := dataMap[fieldName]
	if !exists {
		return "", fmt.Errorf("campo %s não encontrado", fieldName)
	}
	
	if str, ok := value.(string); ok {
		return str, nil
	}
	
	return fmt.Sprintf("%v", value), nil
}

func (d *RelationshipAnomalyDetector) getSeverityForDifference(percentual float64) AnomalySeverity {
	switch {
	case percentual > 50:
		return SeverityCritical
	case percentual > 25:
		return SeverityHigh
	case percentual > 10:
		return SeverityMedium
	default:
		return SeverityLow
	}
}

func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}