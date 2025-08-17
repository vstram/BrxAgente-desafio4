// Package calculo provides functionality for calculating VR values
package calculo

import (
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

// CalcularDiasProporcionais calcula o número de dias proporcionais para um colaborador
// com base na data de admissão ou desligamento no meio do mês
func CalcularDiasProporcionais(colaborador *modelo.Colaborador, diasUteisMes int, mesReferencia time.Time) int {
	// Se o colaborador foi admitido no meio do mês, calcular dias proporcionais
	if !colaborador.DataAdmissao.IsZero() {
		// Verificar se a data de admissão é no mês de referência
		if colaborador.DataAdmissao.Month() == mesReferencia.Month() && 
		   colaborador.DataAdmissao.Year() == mesReferencia.Year() {
			// Calcular dias restantes no mês após a admissão
			diasNoMes := diasNoMes(colaborador.DataAdmissao)
			diasTrabalhados := diasNoMes - colaborador.DataAdmissao.Day() + 1
			return (diasUteisMes * diasTrabalhados) / diasNoMes
		}
	}
	
	// Se o colaborador foi desligado no meio do mês, calcular dias proporcionais
	if colaborador.DataDesligamento != nil && !colaborador.DataDesligamento.IsZero() {
		// Verificar se a data de desligamento é no mês de referência
		if colaborador.DataDesligamento.Month() == mesReferencia.Month() && 
		   colaborador.DataDesligamento.Year() == mesReferencia.Year() {
			// Verificar a regra de desligamento (até dia 15 não considerar, após dia 15 calcular proporcional)
			if colaborador.DataDesligamento.Day() <= 15 {
				// Não considerar para pagamento
				return 0
			} else {
				// Calcular dias trabalhados no mês até a data de desligamento
				diasNoMes := diasNoMes(*colaborador.DataDesligamento)
				diasTrabalhados := colaborador.DataDesligamento.Day()
				return (diasUteisMes * diasTrabalhados) / diasNoMes
			}
		}
	}
	
	// Se não houver datas quebradas, retornar todos os dias úteis do mês
	return diasUteisMes
}

// diasNoMes retorna o número de dias no mês da data fornecida
func diasNoMes(data time.Time) int {
	// Obter o primeiro dia do próximo mês
	primeiroProximoMes := time.Date(data.Year(), data.Month()+1, 1, 0, 0, 0, 0, data.Location())
	
	// Obter o último dia do mês atual (último dia do mês anterior ao próximo mês)
	ultimoDiaMes := primeiroProximoMes.AddDate(0, 0, -1).Day()
	
	return ultimoDiaMes
}

// CalcularDiasProporcionaisParaPeriodo calcula dias proporcionais para um período específico
func CalcularDiasProporcionaisParaPeriodo(inicio, fim time.Time, diasUteisMes int) int {
	// Calcular a diferença em dias entre início e fim
	duracao := fim.Sub(inicio)
	diasPeriodo := int(duracao.Hours() / 24)
	
	// Obter o número de dias no mês
	diasNoMes := diasNoMes(inicio)
	
	// Calcular dias proporcionais
	diasProporcionais := (diasUteisMes * diasPeriodo) / diasNoMes
	
	return diasProporcionais
}