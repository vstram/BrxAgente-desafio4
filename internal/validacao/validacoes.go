// Package validacao provides functionality for validating data integrity
package validacao

import (
	"fmt"

	"BrxAgente-desafio4/internal/modelo"
)

// ValidarCamposObrigatorios verifica se todos os campos obrigatórios do colaborador estão preenchidos
func ValidarCamposObrigatorios(c *modelo.Colaborador) error {
	if c.Matricula == "" {
		return modelo.NovoErroProcessamento("Matrícula é obrigatória")
	}

	if c.Empresa == "" {
		return modelo.NovoErroProcessamento("Empresa é obrigatória")
	}

	if c.Cargo == "" {
		return modelo.NovoErroProcessamento("Cargo é obrigatório")
	}

	if c.Sindicato == "" {
		return modelo.NovoErroProcessamento("Sindicato é obrigatório")
	}

	if c.Situacao == "" {
		return modelo.NovoErroProcessamento("Situação é obrigatória")
	}

	return nil
}

// ValidarDatas verifica a consistência das datas do colaborador
func ValidarDatas(c *modelo.Colaborador) error {
	// Se tem data de admissão, verificar se é válida
	if !c.DataAdmissao.IsZero() {
		// Se tem data de desligamento, verificar se é posterior à data de admissão
		if c.DataDesligamento != nil && !c.DataDesligamento.IsZero() {
			if c.DataDesligamento.Before(c.DataAdmissao) {
				return modelo.NovoErroProcessamento("Data de desligamento não pode ser anterior à data de admissão")
			}
		}
	}

	// Validar períodos de afastamento
	for _, afastamento := range c.Afastamentos {
		if afastamento.Inicio.IsZero() || afastamento.Fim.IsZero() {
			return modelo.NovoErroProcessamento("Datas de início e fim do afastamento são obrigatórias")
		}

		if afastamento.Fim.Before(afastamento.Inicio) {
			return modelo.NovoErroProcessamento("Data de fim do afastamento não pode ser anterior à data de início")
		}
	}

	// Validar períodos de férias
	for _, ferias := range c.Ferias {
		if ferias.Inicio.IsZero() || ferias.Fim.IsZero() {
			return modelo.NovoErroProcessamento("Datas de início e fim das férias são obrigatórias")
		}

		if ferias.Fim.Before(ferias.Inicio) {
			return modelo.NovoErroProcessamento("Data de fim das férias não pode ser anterior à data de início")
		}
	}

	return nil
}

// ValidarRelacionamentos verifica se o colaborador está corretamente relacionado com os dados das outras planilhas
func ValidarRelacionamentos(c *modelo.Colaborador, sindicatos map[string]float64, diasUteis map[string]int) error {
	// Verificar se o sindicato do colaborador existe na planilha de valores por sindicato
	if _, existe := sindicatos[c.Sindicato]; !existe && c.Sindicato != "" {
		return modelo.NovoErroProcessamento(fmt.Sprintf("Sindicato '%s' não encontrado na base de valores", c.Sindicato))
	}

	// Verificar se o sindicato do colaborador existe na planilha de dias úteis
	if _, existe := diasUteis[c.Sindicato]; !existe && c.Sindicato != "" {
		return modelo.NovoErroProcessamento(fmt.Sprintf("Sindicato '%s' não encontrado na base de dias úteis", c.Sindicato))
	}

	return nil
}

// ValidarFormatacao verifica se os dados estão formatados corretamente
func ValidarFormatacao(c *modelo.Colaborador) error {
	// Verificar se a matrícula contém apenas caracteres válidos (letras e números)
	// Esta validação pode ser expandida conforme necessário
	if c.Matricula == "" {
		return modelo.NovoErroProcessamento("Matrícula não pode ser vazia")
	}

	// Verificar se a data de admissão está no formato correto (já foi convertida para time.Time)
	// Esta verificação é feita implicitamente pelo fato de ser um time.Time

	return nil
}

// ValidarColaborador executa todas as validações para um colaborador
func ValidarColaborador(c *modelo.Colaborador, sindicatos map[string]float64, diasUteis map[string]int) []error {
	var erros []error

	// Validar campos obrigatórios
	if err := ValidarCamposObrigatorios(c); err != nil {
		erros = append(erros, err)
	}

	// Validar datas
	if err := ValidarDatas(c); err != nil {
		erros = append(erros, err)
	}

	// Validar formatação
	if err := ValidarFormatacao(c); err != nil {
		erros = append(erros, err)
	}

	// Validar relacionamentos (se as bases estiverem disponíveis)
	if sindicatos != nil && diasUteis != nil {
		if err := ValidarRelacionamentos(c, sindicatos, diasUteis); err != nil {
			erros = append(erros, err)
		}
	}

	return erros
}