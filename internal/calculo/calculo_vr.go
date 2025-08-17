// Package calculo provides functionality for calculating VR values
package calculo

import (
	"strings"
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

// CalcularVRPorColaborador calcula o valor total de VR a ser concedido a um colaborador
// com base nos dias úteis calculados e no valor definido para o seu sindicato
func CalcularVRPorColaborador(colaborador *modelo.Colaborador, valorPorSindicato map[string]float64, diasUteisPorSindicato map[string]int, mesReferencia time.Time) (float64, error) {
	// Mapear o sindicato do colaborador para o estado correspondente
	estadoSindicato := mapearSindicatoParaEstado(colaborador.Sindicato)
	if estadoSindicato == "" {
		// Se não conseguir mapear, usar o nome do sindicato diretamente
		estadoSindicato = colaborador.Sindicato
	}

	// Obter o valor de VR por dia para o sindicato do colaborador
	valorPorDia, existe := valorPorSindicato[estadoSindicato]
	if !existe {
		// Tentar encontrar o valor usando o nome do sindicato diretamente
		valorPorDia, existe = valorPorSindicato[colaborador.Sindicato]
		if !existe {
			// Se não encontrar o valor, retornar 0
			return 0, nil
		}
	}

	// Obter o número de dias úteis para o sindicato do colaborador
	diasUteisSindicato, existe := diasUteisPorSindicato[estadoSindicato]
	if !existe {
		// Tentar encontrar os dias úteis usando o nome do sindicato diretamente
		diasUteisSindicato, existe = diasUteisPorSindicato[colaborador.Sindicato]
		if !existe {
			// Se não encontrar os dias úteis, usar 22 como padrão
			diasUteisSindicato = 22
		}
	}

	// Calcular os dias úteis efetivos para o colaborador
	diasUteisEfetivos := CalcularDiasUteisPorSindicato(colaborador, diasUteisSindicato, mesReferencia)

	// Calcular o valor total de VR
	valorTotal := float64(diasUteisEfetivos) * valorPorDia

	return valorTotal, nil
}

// mapearSindicatoParaEstado mapeia o nome do sindicato do colaborador para o estado correspondente
func mapearSindicatoParaEstado(sindicato string) string {
	// Mapear os sindicatos para os estados
	// This function needs to handle both formats:
	// 1. From "Base dias uteis.xlsx": 
	//    - "SITEPD PR - SIND DOS TRAB EM EMPR PRIVADAS DE PROC DE DADOS DE CURITIBA E REGIAO METROPOLITANA 22"
	//    - "SINDPPD RS - SINDICATO DOS TRAB. EM PROC. DE DADOS RIO GRANDE DO SUL 21"
	//    - "SINDPD SP - SIND.TRAB.EM PROC DADOS E EMPR.EMPRESAS PROC DADOS ESTADO DE SP. 22"
	//    - "SINDPD RJ - SINDICATO PROFISSIONAIS DE PROC DADOS DO RIO DE JANEIRO 21"
	// 2. From "ATIVOS.xlsx": Direct state names like "Paraná", "Rio Grande do Sul", etc.
	
	// Handle direct state names first
	switch strings.ToUpper(sindicato) {
	case "PARANÁ", "PARANA":
		return "Paraná"
	case "RIO GRANDE DO SUL":
		return "Rio Grande do Sul"
	case "SÃO PAULO", "SAO PAULO":
		return "São Paulo"
	case "RIO DE JANEIRO":
		return "Rio de Janeiro"
	}
	
	// Handle sindicato names with acronyms
	switch {
	case strings.Contains(strings.ToUpper(sindicato), "PR") && strings.Contains(strings.ToUpper(sindicato), "CURITIBA"):
		return "Paraná"
	case strings.Contains(strings.ToUpper(sindicato), "RS") && strings.Contains(strings.ToUpper(sindicato), "RIO GRANDE"):
		return "Rio Grande do Sul"
	case strings.Contains(strings.ToUpper(sindicato), "SP") && strings.Contains(strings.ToUpper(sindicato), "SÃO PAULO"):
		return "São Paulo"
	case strings.Contains(strings.ToUpper(sindicato), "RJ") && strings.Contains(strings.ToUpper(sindicato), "RIO DE JANEIRO"):
		return "Rio de Janeiro"
	case strings.Contains(strings.ToUpper(sindicato), "PR -"):
		return "Paraná"
	case strings.Contains(strings.ToUpper(sindicato), "RS -"):
		return "Rio Grande do Sul"
	case strings.Contains(strings.ToUpper(sindicato), "SP -"):
		return "São Paulo"
	case strings.Contains(strings.ToUpper(sindicato), "RJ -"):
		return "Rio de Janeiro"
	default:
		// Check if the sindicato name directly matches a state name
		if sindicato == "Paraná" || sindicato == "Rio Grande do Sul" || sindicato == "São Paulo" || sindicato == "Rio de Janeiro" {
			return sindicato
		}
		// Se não conseguir identificar, retornar vazio
		return ""
	}
}