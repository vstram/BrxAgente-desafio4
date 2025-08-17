// Package modelo provides data structures for the application
package modelo

import (
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
	
	// Nome do colaborador
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