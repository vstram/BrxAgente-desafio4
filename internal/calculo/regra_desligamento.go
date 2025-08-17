// Package calculo provides functionality for calculating VR values
package calculo

import (
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

// AplicarRegraDesligamento aplica a regra específica para desligamentos:
// se comunicado até o dia 15, não considerar; se após o dia 15, calcular proporcionalmente
func AplicarRegraDesligamento(colaborador *modelo.Colaborador, diasUteisMes int, mesReferencia time.Time) int {
	// Verificar se o colaborador foi desligado
	if colaborador.DataComunicacaoDesligamento == nil || colaborador.DataComunicacaoDesligamento.IsZero() {
		// Se não foi desligado, retornar todos os dias úteis do mês
		return diasUteisMes
	}

	// Verificar se a data de comunicação do desligamento é no mês de referência
	if colaborador.DataComunicacaoDesligamento.Month() == mesReferencia.Month() &&
		colaborador.DataComunicacaoDesligamento.Year() == mesReferencia.Year() {
		// Verificar a regra de desligamento (até dia 15 não considerar, após dia 15 calcular proporcional)
		if colaborador.DataComunicacaoDesligamento.Day() <= 15 {
			// Não considerar para pagamento
			return 0
		} else {
			// Calcular dias trabalhados no mês até a data de desligamento
			diasNoMes := diasNoMes(*colaborador.DataComunicacaoDesligamento)
			diasTrabalhados := colaborador.DataComunicacaoDesligamento.Day()
			return (diasUteisMes * diasTrabalhados) / diasNoMes
		}
	}

	// Se a data de comunicação não é no mês de referência, retornar todos os dias úteis do mês
	return diasUteisMes
}