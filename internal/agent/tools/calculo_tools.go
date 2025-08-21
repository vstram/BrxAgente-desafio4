package tools

import (
	"encoding/json"
	"fmt"
	"time"
	
	"BrxAgente-desafio4/internal/calculo"
	"BrxAgente-desafio4/internal/modelo"
)

// CalculateVRTool implementa ferramenta para cálculo de valores de VR
type CalculateVRTool struct {
	*BaseTool
}

// CalculateVRInput representa o input esperado para CalculateVRTool
type CalculateVRInput struct {
	Colaborador              modelo.Colaborador    `json:"colaborador"`
	ValorPorSindicato        map[string]float64    `json:"valor_por_sindicato"`
	DiasUteisPorSindicato    map[string]int        `json:"dias_uteis_por_sindicato"`
	MesReferencia            string                `json:"mes_referencia"` // Formato: "2025-09"
}

// CalculateVROutput representa o output da CalculateVRTool
type CalculateVROutput struct {
	Success                  bool                  `json:"success"`
	Matricula                string                `json:"matricula"`
	Sindicato                string                `json:"sindicato"`
	MesReferencia            string                `json:"mes_referencia"`
	ValorTotal               float64               `json:"valor_total"`
	ValorEmpresa             float64               `json:"valor_empresa"`     // 80%
	ValorColaborador         float64               `json:"valor_colaborador"` // 20%
	DiasUteisEfetivos        int                   `json:"dias_uteis_efetivos"`
	DiasUteisSindicato       int                   `json:"dias_uteis_sindicato"`
	ValorPorDia              float64               `json:"valor_por_dia"`
	Detalhes                 CalculoDetalhes       `json:"detalhes"`
	Error                    string                `json:"error,omitempty"`
}

// CalculoDetalhes fornece informações detalhadas sobre o cálculo
type CalculoDetalhes struct {
	DataAdmissao             string                `json:"data_admissao"`
	DataDesligamento         string                `json:"data_desligamento,omitempty"`
	TemAfastamentos          bool                  `json:"tem_afastamentos"`
	QuantidadeAfastamentos   int                   `json:"quantidade_afastamentos"`
	TemFerias                bool                  `json:"tem_ferias"`
	QuantidadeFerias         int                   `json:"quantidade_ferias"`
	RegraAplicada            string                `json:"regra_aplicada"`
	Observacoes              []string              `json:"observacoes,omitempty"`
}

// NewCalculateVRTool cria uma nova instância da CalculateVRTool
func NewCalculateVRTool() *CalculateVRTool {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"colaborador": map[string]interface{}{
				"type": "object",
				"description": "Dados do colaborador para cálculo",
				"properties": map[string]interface{}{
					"matricula": map[string]interface{}{
						"type": "string",
						"description": "Matrícula do colaborador",
					},
					"sindicato": map[string]interface{}{
						"type": "string", 
						"description": "Sindicato do colaborador",
					},
					"data_admissao": map[string]interface{}{
						"type": "string",
						"format": "date",
						"description": "Data de admissão (YYYY-MM-DD)",
					},
				},
				"required": []string{"matricula", "sindicato"},
			},
			"valor_por_sindicato": map[string]interface{}{
				"type": "object",
				"description": "Mapa com valor por dia para cada sindicato",
			},
			"dias_uteis_por_sindicato": map[string]interface{}{
				"type": "object", 
				"description": "Mapa com dias úteis para cada sindicato",
			},
			"mes_referencia": map[string]interface{}{
				"type": "string",
				"pattern": "^\\d{4}-\\d{2}$",
				"description": "Mês de referência no formato YYYY-MM",
			},
		},
		"required": []string{"colaborador", "valor_por_sindicato", "dias_uteis_por_sindicato", "mes_referencia"},
	}
	
	baseTool := NewBaseTool(
		"calculate_vr",
		"Calcula o valor de Vale Refeição para um colaborador com base nas regras de negócio e dados fornecidos.",
		schema,
	)
	
	return &CalculateVRTool{
		BaseTool: baseTool,
	}
}

// Validate valida o input da ferramenta
func (tool *CalculateVRTool) Validate(input string) error {
	// Validar JSON básico
	if err := tool.ValidateJSON(input); err != nil {
		return err
	}
	
	// Parse input
	var data CalculateVRInput
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return fmt.Errorf("erro ao fazer parse do input: %w", err)
	}
	
	// Validar campos obrigatórios
	if data.Colaborador.Matricula == "" {
		return fmt.Errorf("matrícula do colaborador é obrigatória")
	}
	
	if data.Colaborador.Sindicato == "" {
		return fmt.Errorf("sindicato do colaborador é obrigatório")
	}
	
	if data.MesReferencia == "" {
		return fmt.Errorf("mês de referência é obrigatório")
	}
	
	// Validar formato do mês de referência
	if _, err := time.Parse("2006-01", data.MesReferencia); err != nil {
		return fmt.Errorf("mês de referência deve estar no formato YYYY-MM")
	}
	
	// Validar se há dados de sindicatos
	if len(data.ValorPorSindicato) == 0 {
		return fmt.Errorf("dados de valor por sindicato são obrigatórios")
	}
	
	if len(data.DiasUteisPorSindicato) == 0 {
		return fmt.Errorf("dados de dias úteis por sindicato são obrigatórios")
	}
	
	return nil
}

// Execute executa a ferramenta
func (tool *CalculateVRTool) Execute(input string) (string, error) {
	// Parse input
	var inputData CalculateVRInput
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		return "", NewToolError("calculate_vr", "erro ao fazer parse do input", "PARSE_ERROR")
	}
	
	// Parse mês de referência
	mesReferencia, err := time.Parse("2006-01", inputData.MesReferencia)
	if err != nil {
		return tool.formatErrorOutput(inputData.Colaborador.Matricula, inputData.MesReferencia, 
			"formato de mês de referência inválido")
	}
	
	// Executar cálculo
	valorTotal, err := calculo.CalcularVRPorColaborador(
		&inputData.Colaborador,
		inputData.ValorPorSindicato,
		inputData.DiasUteisPorSindicato,
		mesReferencia,
	)
	if err != nil {
		return tool.formatErrorOutput(inputData.Colaborador.Matricula, inputData.MesReferencia, 
			fmt.Sprintf("erro no cálculo: %v", err))
	}
	
	// Calcular dias úteis efetivos
	diasUteisSindicato := inputData.DiasUteisPorSindicato[inputData.Colaborador.Sindicato]
	if diasUteisSindicato == 0 {
		diasUteisSindicato = 22 // Padrão
	}
	
	diasUteisEfetivos := calculo.CalcularDiasUteisPorSindicato(&inputData.Colaborador, diasUteisSindicato, mesReferencia)
	
	// Calcular valor por dia
	valorPorDia := inputData.ValorPorSindicato[inputData.Colaborador.Sindicato]
	
	// Calcular rateio 80/20
	valorEmpresa := valorTotal * 0.8
	valorColaborador := valorTotal * 0.2
	
	// Montar detalhes
	detalhes := tool.buildDetalhes(&inputData.Colaborador, diasUteisEfetivos, diasUteisSindicato)
	
	// Montar output
	output := &CalculateVROutput{
		Success:               true,
		Matricula:             inputData.Colaborador.Matricula,
		Sindicato:             inputData.Colaborador.Sindicato,
		MesReferencia:         inputData.MesReferencia,
		ValorTotal:            valorTotal,
		ValorEmpresa:          valorEmpresa,
		ValorColaborador:      valorColaborador,
		DiasUteisEfetivos:     diasUteisEfetivos,
		DiasUteisSindicato:    diasUteisSindicato,
		ValorPorDia:           valorPorDia,
		Detalhes:              detalhes,
	}
	
	return tool.FormatJSONOutput(output)
}

// buildDetalhes constrói os detalhes do cálculo
func (tool *CalculateVRTool) buildDetalhes(colaborador *modelo.Colaborador, diasUteisEfetivos, diasUteisSindicato int) CalculoDetalhes {
	detalhes := CalculoDetalhes{
		DataAdmissao:           colaborador.DataAdmissao.Format("2006-01-02"),
		TemAfastamentos:        len(colaborador.Afastamentos) > 0,
		QuantidadeAfastamentos: len(colaborador.Afastamentos),
		TemFerias:              len(colaborador.Ferias) > 0,
		QuantidadeFerias:       len(colaborador.Ferias),
		RegraAplicada:          "calculo_padrao",
		Observacoes:            []string{},
	}
	
	// Adicionar data de desligamento se existir
	if colaborador.DataDesligamento != nil && !colaborador.DataDesligamento.IsZero() {
		detalhes.DataDesligamento = colaborador.DataDesligamento.Format("2006-01-02")
		detalhes.RegraAplicada = "calculo_com_desligamento"
	}
	
	// Adicionar observações baseadas nos cálculos
	if diasUteisEfetivos != diasUteisSindicato {
		if diasUteisEfetivos < diasUteisSindicato {
			detalhes.Observacoes = append(detalhes.Observacoes, 
				fmt.Sprintf("Dias úteis reduzidos de %d para %d devido a afastamentos/férias/admissão tardia", 
					diasUteisSindicato, diasUteisEfetivos))
		}
	}
	
	if len(colaborador.Afastamentos) > 0 {
		detalhes.Observacoes = append(detalhes.Observacoes, 
			fmt.Sprintf("Colaborador possui %d período(s) de afastamento", len(colaborador.Afastamentos)))
	}
	
	if len(colaborador.Ferias) > 0 {
		detalhes.Observacoes = append(detalhes.Observacoes, 
			fmt.Sprintf("Colaborador possui %d período(s) de férias", len(colaborador.Ferias)))
	}
	
	return detalhes
}

// formatErrorOutput formata um output de erro
func (tool *CalculateVRTool) formatErrorOutput(matricula, mesReferencia, errorMsg string) (string, error) {
	output := &CalculateVROutput{
		Success:       false,
		Matricula:     matricula,
		MesReferencia: mesReferencia,
		Error:         errorMsg,
	}
	
	return tool.FormatJSONOutput(output)
}