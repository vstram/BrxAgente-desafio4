// Package modelo provides data structures for the application
// ⚠️ NOTA DE CONFIDENCIALIDADE: Os nomes dos colaboradores não devem ser usados em saídas
// devido aos requisitos de sigilo. Todas as referências aos colaboradores devem ser feitas
// exclusivamente através do campo MATRICULA como identificador único.
package modelo

import (
	"strings"
	"time"
)

// Periodo representa um período com data de início e fim
type Periodo struct {
	// Inicio é a data de início do período
	Inicio time.Time

	// Fim é a data de término do período
	Fim time.Time
}

// Colaborador representa um funcionário com seus dados e períodos de afastamento/férias
type Colaborador struct {
	// Matricula é o identificador único do colaborador
	Matricula string

	// Nome do colaborador - Campo mantido para processamento interno, mas não deve ser usado em saídas
	// devido à Nota de Confidencialidade que requer o uso exclusivo da MATRICULA
	Nome string

	// DataAdmissao é a data em que o colaborador foi admitido
	DataAdmissao time.Time

	// DataDesligamento é a data em que o colaborador foi desligado (pode ser nulo)
	DataDesligamento *time.Time

	// DataComunicacaoDesligamento é a data em que o desligamento foi comunicado (pode ser nulo)
	DataComunicacaoDesligamento *time.Time

	// Sindicato ao qual o colaborador pertence
	Sindicato string

	// Afastamentos é uma lista de períodos em que o colaborador esteve afastado
	Afastamentos []Periodo

	// Ferias é uma lista de períodos em que o colaborador esteve de férias
	Ferias []Periodo

	// Empresa associada ao colaborador
	Empresa string

	// Cargo do colaborador
	Cargo string

	// Situacao atual do colaborador (Trabalhando, Férias, etc.)
	Situacao string

	// ValorTotalVR é o valor total de VR a ser concedido ao colaborador
	ValorTotalVR float64

	// ValorEmpresa é a parcela paga pela empresa (80% do valor total)
	ValorEmpresa float64

	// ValorColaborador é a parcela descontada do colaborador (20% do valor total)
	ValorColaborador float64

	// DiasUteisEfetivos é o número de dias úteis efetivos calculados para o colaborador
	DiasUteisEfetivos int
}

// EstaAtivo verifica se o colaborador está ativo em uma determinada data
func (c *Colaborador) EstaAtivo(data time.Time) bool {
	// Se não tem data de desligamento, está ativo
	if c.DataDesligamento == nil {
		return true
	}

	// Se a data é anterior à data de desligamento, está ativo
	return data.Before(*c.DataDesligamento)
}

// EhElegivel verifica se o colaborador é elegível para receber o benefício de VR
func (c *Colaborador) EhElegivel() bool {
	// Verificar se é diretor (não elegível)
	if c.Cargo != "" && strings.Contains(strings.ToUpper(c.Cargo), "DIRETOR") {
		return false
	}

	// Verificar se é estagiário (não elegível)
	if c.Cargo != "" && (strings.Contains(strings.ToUpper(c.Cargo), "ESTAGIÁRIO") || strings.Contains(strings.ToUpper(c.Cargo), "ESTAGIARIO")) {
		return false
	}

	// Verificar se é aprendiz (não elegível)
	if c.Cargo != "" && strings.Contains(strings.ToUpper(c.Cargo), "APRENDIZ") {
		return false
	}

	// Verificar se está afastado (não elegível)
	// Esta verificação pode ser feita em outro lugar onde temos as informações completas

	// Verificar se está no exterior (não elegível)
	// Esta verificação pode ser feita em outro lugar onde temos as informações completas

	// Se não se enquadra em nenhuma das categorias acima, é elegível
	return true
}
