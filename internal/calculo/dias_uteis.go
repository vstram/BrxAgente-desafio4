// Package calculo provides functionality for calculating VR values
package calculo

import (
	"time"

	"BrxAgente-desafio4/internal/feriados"
	"BrxAgente-desafio4/internal/modelo"
)

// FeriadosCache armazena os feriados calculados para evitar recálculos
var feriadosCache map[int][]feriados.Feriado
var feriadosEstaduaisCache map[string]map[int][]feriados.Feriado

func init() {
	feriadosCache = make(map[int][]feriados.Feriado)
	feriadosEstaduaisCache = make(map[string]map[int][]feriados.Feriado)
}

// ObterFeriadosNacionaisCached obtém os feriados nacionais com cache
func ObterFeriadosNacionaisCached(ano int) []feriados.Feriado {
	if feriados, existe := feriadosCache[ano]; existe {
		return feriados
	}

	feriados := feriados.ObterFeriadosNacionais(ano)
	feriadosCache[ano] = feriados
	return feriados
}

// ObterFeriadosEstaduaisCached obtém os feriados estaduais com cache
func ObterFeriadosEstaduaisCached(estado string, ano int) []feriados.Feriado {
	if feriadosEstaduaisCache[estado] == nil {
		feriadosEstaduaisCache[estado] = make(map[int][]feriados.Feriado)
	}

	if feriados, existe := feriadosEstaduaisCache[estado][ano]; existe {
		return feriados
	}

	feriados := feriados.ObterFeriadosEstaduais(estado, ano)
	feriadosEstaduaisCache[estado][ano] = feriados
	return feriados
}

// CalcularDiasUteisPorSindicato calcula os dias úteis relevantes para um colaborador,
// considerando o calendário de dias úteis de seu sindicato, férias, outros afastamentos,
// data de desligamento e feriados
func CalcularDiasUteisPorSindicato(colaborador *modelo.Colaborador, diasUteisSindicato int, mesReferencia time.Time) int {
	// Começar com todos os dias úteis do sindicato
	diasUteis := diasUteisSindicato

	// Obter feriados para o mês de referência com cache para evitar recálculos
	feriadosNacionais := ObterFeriadosNacionaisCached(mesReferencia.Year())
	
	// Determinar o estado do colaborador com base no sindicato
	estado := determinarEstadoPorSindicato(colaborador.Sindicato)
	feriadosEstaduais := ObterFeriadosEstaduaisCached(estado, mesReferencia.Year())
	
	// TODO: Implementar obtenção de feriados municipais quando disponível
	// feriadosMunicipais := feriados.ObterFeriadosMunicipais(municipio, estado, mesReferencia.Year())
	
	// Combinar todos os feriados
	todosFeriados := append(feriadosNacionais, feriadosEstaduais...)
	// todosFeriados = append(todosFeriados, feriadosMunicipais...)

	// Descontar feriados do total de dias úteis
	diasUteis -= contarFeriadosNoMes(todosFeriados, mesReferencia)

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

// determinarEstadoPorSindicato determina o estado com base no nome do sindicato
func determinarEstadoPorSindicato(sindicato string) string {
	// Mapear sindicatos para estados
	switch {
	case containsIgnoreCase(sindicato, "São Paulo") || containsIgnoreCase(sindicato, "SP"):
		return "SP"
	case containsIgnoreCase(sindicato, "Rio de Janeiro") || containsIgnoreCase(sindicato, "RJ"):
		return "RJ"
	case containsIgnoreCase(sindicato, "Paraná") || containsIgnoreCase(sindicato, "PR"):
		return "PR"
	case containsIgnoreCase(sindicato, "Rio Grande do Sul") || containsIgnoreCase(sindicato, "RS"):
		return "RS"
	default:
		return ""
	}
}

// containsIgnoreCase verifica se uma string contém outra string ignorando maiúsculas/minúsculas
func containsIgnoreCase(s, substr string) bool {
	return containsStringIgnoreCase(s, substr)
}

// contarFeriadosNoMes conta quantos feriados existem em um mês específico
func contarFeriadosNoMes(feriadosLista []feriados.Feriado, mesReferencia time.Time) int {
	inicioMes := time.Date(mesReferencia.Year(), mesReferencia.Month(), 1, 0, 0, 0, 0, time.UTC)
	fimMes := time.Date(mesReferencia.Year(), mesReferencia.Month()+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	
	return feriados.ContarFeriadosNoPeriodo(inicioMes, fimMes, feriadosLista)
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

// containsStringIgnoreCase verifica se uma string contém outra string ignorando case
func containsStringIgnoreCase(s, substr string) bool {
	// Convert both strings to uppercase for comparison
	sUpper := ""
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			sUpper += string(r - 'a' + 'A')
		} else {
			sUpper += string(r)
		}
	}
	
	substrUpper := ""
	for _, r := range substr {
		if r >= 'a' && r <= 'z' {
			substrUpper += string(r - 'a' + 'A')
		} else {
			substrUpper += string(r)
		}
	}
	
	// Simple contains implementation
	for i := 0; i <= len(sUpper)-len(substrUpper); i++ {
		if sUpper[i:i+len(substrUpper)] == substrUpper {
			return true
		}
	}
	return false
}