// Package calculo provides functionality for calculating VR values
package calculo

import (
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

// CalcularDiasUteisPorSindicato calcula os dias úteis relevantes para um colaborador,
// considerando o calendário de dias úteis do seu sindicato, férias, outros afastamentos
// e data de desligamento
func CalcularDiasUteisPorSindicato(colaborador *modelo.Colaborador, diasUteisSindicato int, mesReferencia time.Time) int {
	// Começar com todos os dias úteis do sindicato
	diasUteis := diasUteisSindicato

	// Aplicar regra de datas quebradas (admissão/desligamento no meio do mês)
	diasUteis = CalcularDiasProporcionais(colaborador, diasUteis, mesReferencia)

	// Descontar dias de férias
	diasUteis -= calcularDiasFerias(colaborador, diasUteisSindicato, mesReferencia)

	// Descontar dias de outros afastamentos
	diasUteis -= calcularDiasAfastamentos(colaborador, diasUteisSindicato, mesReferencia)

	// Garantir que não tenha dias negativos
	if diasUteis < 0 {
		diasUteis = 0
	}

	return diasUteis
}

// calcularDiasFerias calcula os dias a serem descontados por férias
func calcularDiasFerias(colaborador *modelo.Colaborador, diasUteisSindicato int, mesReferencia time.Time) int {
	diasDescontados := 0

	// Para cada período de férias do colaborador
	for _, ferias := range colaborador.Ferias {
		// Verificar se o período de férias intersecta com o mês de referência
		if periodoIntersectaMes(ferias, mesReferencia) {
			// Calcular dias úteis proporcionais ao período de férias
			diasPeriodo := calcularDiasProporcionaisParaPeriodo(ferias, mesReferencia)
			diasDescontados += (diasUteisSindicato * diasPeriodo) / diasNoMes(mesReferencia)
		}
	}

	return diasDescontados
}

// calcularDiasAfastamentos calcula os dias a serem descontados por outros afastamentos
func calcularDiasAfastamentos(colaborador *modelo.Colaborador, diasUteisSindicato int, mesReferencia time.Time) int {
	diasDescontados := 0

	// Para cada período de afastamento do colaborador
	for _, afastamento := range colaborador.Afastamentos {
		// Verificar se o período de afastamento intersecta com o mês de referência
		if periodoIntersectaMes(afastamento, mesReferencia) {
			// Calcular dias úteis proporcionais ao período de afastamento
			diasPeriodo := calcularDiasProporcionaisParaPeriodo(afastamento, mesReferencia)
			diasDescontados += (diasUteisSindicato * diasPeriodo) / diasNoMes(mesReferencia)
		}
	}

	return diasDescontados
}

// periodoIntersectaMes verifica se um período intersecta com um mês/ano específico
func periodoIntersectaMes(periodo modelo.Periodo, mesReferencia time.Time) bool {
	// Criar datas para o início e fim do mês de referência
	inicioMes := time.Date(mesReferencia.Year(), mesReferencia.Month(), 1, 0, 0, 0, 0, time.UTC)
	fimMes := time.Date(mesReferencia.Year(), mesReferencia.Month()+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)

	// Verificar se há intersecção entre os períodos
	// Dois períodos se intersectam se:
	// - O início do primeiro é menor ou igual ao fim do segundo E
	// - O fim do primeiro é maior ou igual ao início do segundo
	return !periodo.Inicio.After(fimMes) && !periodo.Fim.Before(inicioMes)
}

// calcularDiasProporcionaisParaPeriodo calcula os dias proporcionais de um período
// que caem no mês de referência
func calcularDiasProporcionaisParaPeriodo(periodo modelo.Periodo, mesReferencia time.Time) int {
	// Criar datas para o início e fim do mês de referência
	inicioMes := time.Date(mesReferencia.Year(), mesReferencia.Month(), 1, 0, 0, 0, 0, time.UTC)
	fimMes := time.Date(mesReferencia.Year(), mesReferencia.Month()+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)

	// Determinar o período efetivo dentro do mês
	inicioEfetivo := periodo.Inicio
	if inicioEfetivo.Before(inicioMes) {
		inicioEfetivo = inicioMes
	}

	fimEfetivo := periodo.Fim
	if fimEfetivo.After(fimMes) {
		fimEfetivo = fimMes
	}

	// Calcular a diferença em dias
	duracao := fimEfetivo.Sub(inicioEfetivo)
	dias := int(duracao.Hours()/24) + 1 // +1 porque incluímos o dia inicial

	// Garantir que não tenha dias negativos
	if dias < 0 {
		dias = 0
	}

	return dias
}