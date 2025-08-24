package tools

import (
	"encoding/json"
	"fmt"
	"time"

	"BrxAgente-desafio4/internal/modelo"
	"BrxAgente-desafio4/internal/validacao"
)

// ValidateDataTool implementa ferramenta para validação de dados
type ValidateDataTool struct {
	*BaseTool
}

// ValidateDataInput representa o input esperado para ValidateDataTool
type ValidateDataInput struct {
	// Tipo de validação a ser executada
	TipoValidacao string `json:"tipo_validacao"` // "colaborador", "lista_colaboradores", "planilha"

	// Dados para validação
	Colaborador   *modelo.Colaborador    `json:"colaborador,omitempty"`
	Colaboradores []modelo.Colaborador   `json:"colaboradores,omitempty"`
	DadosPlanilha map[string]interface{} `json:"dados_planilha,omitempty"`

	// Parâmetros de validação
	ValidarCamposObrigatorios bool `json:"validar_campos_obrigatorios"`
	ValidarDatas              bool `json:"validar_datas"`
	ValidarConsistencia       bool `json:"validar_consistencia"`
	ValidarDuplicatas         bool `json:"validar_duplicatas"`
}

// ValidateDataOutput representa o output da ValidateDataTool
type ValidateDataOutput struct {
	Success          bool               `json:"success"`
	TipoValidacao    string             `json:"tipo_validacao"`
	TotalRegistros   int                `json:"total_registros"`
	RegistrosValidos int                `json:"registros_validos"`
	RegistrosComErro int                `json:"registros_com_erro"`
	Erros            []ErroValidacao    `json:"erros"`
	Warnings         []WarningValidacao `json:"warnings,omitempty"`
	Resumo           ResumoValidacao    `json:"resumo"`
	Error            string             `json:"error,omitempty"`
}

// ErroValidacao representa um erro encontrado na validação
type ErroValidacao struct {
	Matricula       string `json:"matricula,omitempty"`
	Campo           string `json:"campo"`
	TipoErro        string `json:"tipo_erro"`
	Mensagem        string `json:"mensagem"`
	ValorEncontrado string `json:"valor_encontrado,omitempty"`
	ValorEsperado   string `json:"valor_esperado,omitempty"`
	Severidade      string `json:"severidade"` // "critico", "alto", "medio", "baixo"
}

// WarningValidacao representa um warning encontrado na validação
type WarningValidacao struct {
	Matricula    string `json:"matricula,omitempty"`
	Campo        string `json:"campo,omitempty"`
	Mensagem     string `json:"mensagem"`
	Recomendacao string `json:"recomendacao,omitempty"`
}

// ResumoValidacao fornece um resumo dos resultados da validação
type ResumoValidacao struct {
	PercentualValidos float64  `json:"percentual_validos"`
	ErrosCriticos     int      `json:"erros_criticos"`
	ErrosAltos        int      `json:"erros_altos"`
	ErrosMedios       int      `json:"erros_medios"`
	ErrosBaixos       int      `json:"erros_baixos"`
	TotalWarnings     int      `json:"total_warnings"`
	StatusGeral       string   `json:"status_geral"` // "aprovado", "aprovado_com_ressalvas", "rejeitado"
	Observacoes       []string `json:"observacoes,omitempty"`
}

// NewValidateDataTool cria uma nova instância da ValidateDataTool
func NewValidateDataTool() *ValidateDataTool {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"tipo_validacao": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"colaborador", "lista_colaboradores", "planilha"},
				"description": "Tipo de validação: colaborador único, lista de colaboradores ou dados de planilha",
			},
			"colaborador": map[string]interface{}{
				"type":        "object",
				"description": "Dados do colaborador para validação (usado com tipo_validacao='colaborador')",
			},
			"colaboradores": map[string]interface{}{
				"type":        "array",
				"description": "Lista de colaboradores para validação (usado com tipo_validacao='lista_colaboradores')",
			},
			"dados_planilha": map[string]interface{}{
				"type":        "object",
				"description": "Dados da planilha para validação (usado com tipo_validacao='planilha')",
			},
			"validar_campos_obrigatorios": map[string]interface{}{
				"type":        "boolean",
				"description": "Se deve validar campos obrigatórios",
				"default":     true,
			},
			"validar_datas": map[string]interface{}{
				"type":        "boolean",
				"description": "Se deve validar consistência de datas",
				"default":     true,
			},
			"validar_consistencia": map[string]interface{}{
				"type":        "boolean",
				"description": "Se deve validar consistência de dados",
				"default":     true,
			},
			"validar_duplicatas": map[string]interface{}{
				"type":        "boolean",
				"description": "Se deve verificar duplicatas",
				"default":     false,
			},
		},
		"required": []string{"tipo_validacao"},
	}

	baseTool := NewBaseTool(
		"validate_data",
		"Valida dados de colaboradores verificando campos obrigatórios, consistência de datas, duplicatas e outras regras de negócio.",
		schema,
	)

	return &ValidateDataTool{
		BaseTool: baseTool,
	}
}

// Validate valida o input da ferramenta
func (tool *ValidateDataTool) Validate(input string) error {
	// Validar JSON básico
	if err := tool.ValidateJSON(input); err != nil {
		return err
	}

	// Parse input
	var data ValidateDataInput
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return fmt.Errorf("erro ao fazer parse do input: %w", err)
	}

	// Validar tipo de validação
	validTipos := map[string]bool{
		"colaborador":         true,
		"lista_colaboradores": true,
		"planilha":            true,
	}

	if !validTipos[data.TipoValidacao] {
		return fmt.Errorf("tipo de validação inválido: %s", data.TipoValidacao)
	}

	// Validar se os dados necessários foram fornecidos
	switch data.TipoValidacao {
	case "colaborador":
		if data.Colaborador == nil {
			return fmt.Errorf("dados do colaborador são obrigatórios para tipo_validacao='colaborador'")
		}
	case "lista_colaboradores":
		if len(data.Colaboradores) == 0 {
			return fmt.Errorf("lista de colaboradores é obrigatória para tipo_validacao='lista_colaboradores'")
		}
	case "planilha":
		if data.DadosPlanilha == nil {
			return fmt.Errorf("dados da planilha são obrigatórios para tipo_validacao='planilha'")
		}
	}

	return nil
}

// Execute executa a ferramenta
func (tool *ValidateDataTool) Execute(input string) (string, error) {
	// Parse input
	var inputData ValidateDataInput
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		return "", NewToolError("validate_data", "erro ao fazer parse do input", "PARSE_ERROR")
	}

	// Executar validação baseada no tipo
	var output *ValidateDataOutput
	var err error

	switch inputData.TipoValidacao {
	case "colaborador":
		output, err = tool.validateColaborador(inputData)
	case "lista_colaboradores":
		output, err = tool.validateListaColaboradores(inputData)
	case "planilha":
		output, err = tool.validatePlanilha(inputData)
	default:
		return tool.formatErrorOutput(inputData.TipoValidacao, "tipo de validação não suportado")
	}

	if err != nil {
		return tool.formatErrorOutput(inputData.TipoValidacao, fmt.Sprintf("erro na validação: %v", err))
	}

	return tool.FormatJSONOutput(output)
}

// validateColaborador valida um único colaborador
func (tool *ValidateDataTool) validateColaborador(input ValidateDataInput) (*ValidateDataOutput, error) {
	colaborador := input.Colaborador
	erros := []ErroValidacao{}
	warnings := []WarningValidacao{}

	// Validar campos obrigatórios
	if input.ValidarCamposObrigatorios {
		if err := validacao.ValidarCamposObrigatorios(colaborador); err != nil {
			erros = append(erros, ErroValidacao{
				Matricula:  colaborador.Matricula,
				Campo:      "campos_obrigatorios",
				TipoErro:   "campo_vazio",
				Mensagem:   err.Error(),
				Severidade: "critico",
			})
		}
	}

	// Validar datas
	if input.ValidarDatas {
		if err := validacao.ValidarDatas(colaborador); err != nil {
			erros = append(erros, ErroValidacao{
				Matricula:  colaborador.Matricula,
				Campo:      "datas",
				TipoErro:   "data_inconsistente",
				Mensagem:   err.Error(),
				Severidade: "alto",
			})
		}
	}

	// Validações adicionais específicas
	tool.validateColaboradorEspecifico(colaborador, &erros, &warnings)

	// Montar resumo
	registrosValidosCount := 0
	if len(erros) == 0 {
		registrosValidosCount = 1
	}
	resumo := tool.buildResumo(1, registrosValidosCount, erros, warnings)

	output := &ValidateDataOutput{
		Success:          true,
		TipoValidacao:    input.TipoValidacao,
		TotalRegistros:   1,
		RegistrosValidos: registrosValidosCount,
		RegistrosComErro: len(erros),
		Erros:            erros,
		Warnings:         warnings,
		Resumo:           resumo,
	}

	return output, nil
}

// validateListaColaboradores valida uma lista de colaboradores
func (tool *ValidateDataTool) validateListaColaboradores(input ValidateDataInput) (*ValidateDataOutput, error) {
	colaboradores := input.Colaboradores
	erros := []ErroValidacao{}
	warnings := []WarningValidacao{}
	registrosValidos := 0
	matriculasVistas := make(map[string]bool)

	for _, colaborador := range colaboradores {
		colaboradorValido := true

		// Verificar duplicatas
		if input.ValidarDuplicatas {
			if matriculasVistas[colaborador.Matricula] {
				erros = append(erros, ErroValidacao{
					Matricula:  colaborador.Matricula,
					Campo:      "matricula",
					TipoErro:   "duplicata",
					Mensagem:   "Matrícula duplicada na lista",
					Severidade: "critico",
				})
				colaboradorValido = false
			}
			matriculasVistas[colaborador.Matricula] = true
		}

		// Validar campos obrigatórios
		if input.ValidarCamposObrigatorios {
			if err := validacao.ValidarCamposObrigatorios(&colaborador); err != nil {
				erros = append(erros, ErroValidacao{
					Matricula:  colaborador.Matricula,
					Campo:      "campos_obrigatorios",
					TipoErro:   "campo_vazio",
					Mensagem:   err.Error(),
					Severidade: "critico",
				})
				colaboradorValido = false
			}
		}

		// Validar datas
		if input.ValidarDatas {
			if err := validacao.ValidarDatas(&colaborador); err != nil {
				erros = append(erros, ErroValidacao{
					Matricula:  colaborador.Matricula,
					Campo:      "datas",
					TipoErro:   "data_inconsistente",
					Mensagem:   err.Error(),
					Severidade: "alto",
				})
				colaboradorValido = false
			}
		}

		// Validações específicas
		tool.validateColaboradorEspecifico(&colaborador, &erros, &warnings)

		if colaboradorValido {
			registrosValidos++
		}
	}

	// Montar resumo
	resumo := tool.buildResumo(len(colaboradores), registrosValidos, erros, warnings)

	output := &ValidateDataOutput{
		Success:          true,
		TipoValidacao:    input.TipoValidacao,
		TotalRegistros:   len(colaboradores),
		RegistrosValidos: registrosValidos,
		RegistrosComErro: len(colaboradores) - registrosValidos,
		Erros:            erros,
		Warnings:         warnings,
		Resumo:           resumo,
	}

	return output, nil
}

// validatePlanilha valida dados de planilha
func (tool *ValidateDataTool) validatePlanilha(input ValidateDataInput) (*ValidateDataOutput, error) {
	// Implementação simplificada para validação de planilha
	dados := input.DadosPlanilha
	erros := []ErroValidacao{}
	warnings := []WarningValidacao{}

	// Verificar estrutura básica da planilha
	// ⚠️ CONFIDENCIALIDADE: Remove campo Nome conforme PRD.md
	if headers, ok := dados["headers"].([]interface{}); ok {
		expectedHeaders := []string{"Matricula", "DataAdmissao", "Empresa", "Sindicato"}
		for _, expected := range expectedHeaders {
			found := false
			for _, header := range headers {
				if header.(string) == expected {
					found = true
					break
				}
			}
			if !found {
				erros = append(erros, ErroValidacao{
					Campo:      "headers",
					TipoErro:   "header_ausente",
					Mensagem:   fmt.Sprintf("Header obrigatório '%s' não encontrado", expected),
					Severidade: "critico",
				})
			}
		}
	}

	// Verificar se há dados
	totalRegistros := 0
	if rowCount, ok := dados["row_count"].(float64); ok {
		totalRegistros = int(rowCount)
		if totalRegistros == 0 {
			warnings = append(warnings, WarningValidacao{
				Campo:        "dados",
				Mensagem:     "Planilha não contém dados",
				Recomendacao: "Verificar se a planilha foi preenchida corretamente",
			})
		}
	}

	// Montar resumo
	registrosValidosCount := totalRegistros
	if len(erros) > 0 {
		registrosValidosCount = 0
	}
	resumo := tool.buildResumo(totalRegistros, registrosValidosCount, erros, warnings)

	output := &ValidateDataOutput{
		Success:          true,
		TipoValidacao:    input.TipoValidacao,
		TotalRegistros:   totalRegistros,
		RegistrosValidos: totalRegistros,
		RegistrosComErro: len(erros),
		Erros:            erros,
		Warnings:         warnings,
		Resumo:           resumo,
	}

	return output, nil
}

// validateColaboradorEspecifico executa validações específicas do domínio
func (tool *ValidateDataTool) validateColaboradorEspecifico(colaborador *modelo.Colaborador, erros *[]ErroValidacao, warnings *[]WarningValidacao) {
	// Validar se data de admissão não é futura
	if !colaborador.DataAdmissao.IsZero() && colaborador.DataAdmissao.After(time.Now()) {
		*erros = append(*erros, ErroValidacao{
			Matricula:       colaborador.Matricula,
			Campo:           "data_admissao",
			TipoErro:        "data_futura",
			Mensagem:        "Data de admissão não pode ser futura",
			ValorEncontrado: colaborador.DataAdmissao.Format("2006-01-02"),
			Severidade:      "alto",
		})
	}

	// Validar formato da matrícula (exemplo: deve ter pelo menos 4 caracteres)
	if len(colaborador.Matricula) < 4 {
		*warnings = append(*warnings, WarningValidacao{
			Matricula:    colaborador.Matricula,
			Campo:        "matricula",
			Mensagem:     "Matrícula muito curta (menos de 4 caracteres)",
			Recomendacao: "Verificar se a matrícula está correta",
		})
	}

	// Validar sindicatos conhecidos
	sindicatosConhecidos := map[string]bool{
		"SINDPD":      true,
		"SINDICATO_A": true,
		"SINDICATO_B": true,
	}

	if !sindicatosConhecidos[colaborador.Sindicato] {
		*warnings = append(*warnings, WarningValidacao{
			Matricula:    colaborador.Matricula,
			Campo:        "sindicato",
			Mensagem:     fmt.Sprintf("Sindicato '%s' não reconhecido", colaborador.Sindicato),
			Recomendacao: "Verificar se o nome do sindicato está correto",
		})
	}
}

// buildResumo constrói o resumo da validação
func (tool *ValidateDataTool) buildResumo(total, registrosValidos int, erros []ErroValidacao, warnings []WarningValidacao) ResumoValidacao {
	var percentualValidos float64
	if total > 0 {
		percentualValidos = float64(registrosValidos) / float64(total) * 100
	}

	// Contar erros por severidade
	errosCriticos, errosAltos, errosMedios, errosBaixos := 0, 0, 0, 0
	for _, erro := range erros {
		switch erro.Severidade {
		case "critico":
			errosCriticos++
		case "alto":
			errosAltos++
		case "medio":
			errosMedios++
		case "baixo":
			errosBaixos++
		}
	}

	// Determinar status geral
	statusGeral := "aprovado"
	if errosCriticos > 0 {
		statusGeral = "rejeitado"
	} else if errosAltos > 0 || len(warnings) > 5 {
		statusGeral = "aprovado_com_ressalvas"
	}

	// Observações
	observacoes := []string{}
	if errosCriticos > 0 {
		observacoes = append(observacoes, fmt.Sprintf("%d erro(s) crítico(s) encontrado(s)", errosCriticos))
	}
	if len(warnings) > 0 {
		observacoes = append(observacoes, fmt.Sprintf("%d warning(s) encontrado(s)", len(warnings)))
	}
	if percentualValidos == 100 {
		observacoes = append(observacoes, "Todos os registros passaram na validação")
	}

	return ResumoValidacao{
		PercentualValidos: percentualValidos,
		ErrosCriticos:     errosCriticos,
		ErrosAltos:        errosAltos,
		ErrosMedios:       errosMedios,
		ErrosBaixos:       errosBaixos,
		TotalWarnings:     len(warnings),
		StatusGeral:       statusGeral,
		Observacoes:       observacoes,
	}
}

// formatErrorOutput formata um output de erro
func (tool *ValidateDataTool) formatErrorOutput(tipoValidacao, errorMsg string) (string, error) {
	output := &ValidateDataOutput{
		Success:       false,
		TipoValidacao: tipoValidacao,
		Error:         errorMsg,
	}

	return tool.FormatJSONOutput(output)
}
