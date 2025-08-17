// Package calculo provides functionality for calculating VR values
package calculo

import (
	"testing"
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

func TestCalcularVRPorColaborador(t *testing.T) {
	// Data de referência para testes (maio de 2025)
	mesReferencia := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)

	// Mapas de teste
	valorPorSindicato := map[string]float64{
		"Paraná": 35.00,
	}

	diasUteisPorSindicato := map[string]int{
		"Paraná": 22,
	}

	t.Run("ColaboradorSemFeriasOuAfastamentos", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Sindicato: "Paraná",
			Ferias:    []modelo.Periodo{},
			Afastamentos: []modelo.Periodo{},
		}

		valor, err := CalcularVRPorColaborador(colaborador, valorPorSindicato, diasUteisPorSindicato, mesReferencia)
		
		// Valor esperado: 22 dias * R$ 35,00 = R$ 770,00
		valorEsperado := 22.0 * 35.0

		if err != nil {
			t.Errorf("Não esperava erro, mas obteve: %v", err)
		}

		if valor != valorEsperado {
			t.Errorf("Esperava valor R$ %.2f, mas obteve R$ %.2f", valorEsperado, valor)
		}
	})

	t.Run("ColaboradorComFerias", func(t *testing.T) {
		// Criar um período de férias de 5 dias em maio
		dataInicioFerias := time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC)
		dataFimFerias := time.Date(2025, 5, 14, 0, 0, 0, 0, time.UTC)
		
		colaborador := &modelo.Colaborador{
			Matricula: "002",
			Sindicato: "Paraná",
			Ferias: []modelo.Periodo{
				{Inicio: dataInicioFerias, Fim: dataFimFerias},
			},
			Afastamentos: []modelo.Periodo{},
		}

		valor, err := CalcularVRPorColaborador(colaborador, valorPorSindicato, diasUteisPorSindicato, mesReferencia)
		
		// Valor esperado: (22 - dias descontados) * R$ 35,00
		// Dias descontados: (22 * 5) / 31 ≈ 3.55 => 3 dias
		// Valor esperado: (22 - 3) * 35.0 = 19 * 35.0 = 665.0
		valorEsperado := 19.0 * 35.0

		if err != nil {
			t.Errorf("Não esperava erro, mas obteve: %v", err)
		}

		if valor != valorEsperado {
			t.Errorf("Esperava valor R$ %.2f, mas obteve R$ %.2f", valorEsperado, valor)
		}
	})

	t.Run("ColaboradorComSindicatoNaoMapeado", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "003",
			Sindicato: "Sindicato Inexistente",
			Ferias:    []modelo.Periodo{},
			Afastamentos: []modelo.Periodo{},
		}

		valor, err := CalcularVRPorColaborador(colaborador, valorPorSindicato, diasUteisPorSindicato, mesReferencia)
		
		// Como o sindicato não existe, deve retornar 0
		if err != nil {
			t.Errorf("Não esperava erro, mas obteve: %v", err)
		}

		if valor != 0 {
			t.Errorf("Esperava valor R$ 0.00, mas obteve R$ %.2f", valor)
		}
	})

	t.Run("ColaboradorComAdmissaoNoMeioDoMes", func(t *testing.T) {
		// Admissão no dia 15 de maio
		dataAdmissao := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)
		
		colaborador := &modelo.Colaborador{
			Matricula:    "004",
			Sindicato:    "Paraná",
			DataAdmissao: dataAdmissao,
			Ferias:       []modelo.Periodo{},
			Afastamentos: []modelo.Periodo{},
		}

		valor, err := CalcularVRPorColaborador(colaborador, valorPorSindicato, diasUteisPorSindicato, mesReferencia)
		
		// Dias trabalhados = 31 - 15 + 1 = 17
		// Dias proporcionais = (22 * 17) / 31 ≈ 12.06 => 12 dias
		// Valor esperado: 12 * 35.0 = 420.0
		valorEsperado := 12.0 * 35.0

		if err != nil {
			t.Errorf("Não esperava erro, mas obteve: %v", err)
		}

		if valor != valorEsperado {
			t.Errorf("Esperava valor R$ %.2f, mas obteve R$ %.2f", valorEsperado, valor)
		}
	})
}

func TestMapearSindicatoParaEstado(t *testing.T) {
	testCases := []struct {
		sindicato string
		estado    string
	}{
		{"SITEPD PR - SIND DOS TRAB EM EMPR PRIVADAS DE PROC DE DADOS DE CURITIBA E REGIAO METROPOLITANA 22", "Paraná"},
		{"SINDPPD RS - SINDICATO DOS TRAB. EM PROC. DE DADOS RIO GRANDE DO SUL 21", "Rio Grande do Sul"},
		{"SINDPD SP - SIND.TRAB.EM PROC DADOS E EMPR.EMPRESAS PROC DADOS ESTADO DE SP. 22", "São Paulo"},
		{"SINDPD RJ - SINDICATO PROFISSIONAIS DE PROC DADOS DO RIO DE JANEIRO 21", "Rio de Janeiro"},
		{"Paraná", "Paraná"},
		{"Rio Grande do Sul", "Rio Grande do Sul"},
		{"São Paulo", "São Paulo"},
		{"Rio de Janeiro", "Rio de Janeiro"},
		{"Sindicato Qualquer", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.sindicato, func(t *testing.T) {
			estado := mapearSindicatoParaEstado(tc.sindicato)
			if estado != tc.estado {
				t.Errorf("Esperava estado '%s' para sindicato '%s', mas obteve '%s'", tc.estado, tc.sindicato, estado)
			}
		})
	}
}