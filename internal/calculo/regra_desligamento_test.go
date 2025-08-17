// Package calculo provides functionality for calculating VR values
package calculo

import (
	"testing"
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

func TestAplicarRegraDesligamento(t *testing.T) {
	// Data de referência para testes (maio de 2025)
	mesReferencia := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	diasUteisMes := 22 // Exemplo de dias úteis em maio

	t.Run("SemDataComunicacaoDesligamento", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
		}

		dias := AplicarRegraDesligamento(colaborador, diasUteisMes, mesReferencia)

		// Deve retornar todos os dias úteis do mês
		if dias != diasUteisMes {
			t.Errorf("Esperava %d dias, mas obteve %d", diasUteisMes, dias)
		}
	})

	t.Run("DataComunicacaoAntesDoMesReferencia", func(t *testing.T) {
		// Data de comunicação no dia 15 de abril (mês anterior)
		dataComunicacao := time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC)
		colaborador := &modelo.Colaborador{
			Matricula:                   "001",
			DataComunicacaoDesligamento: &dataComunicacao,
		}

		dias := AplicarRegraDesligamento(colaborador, diasUteisMes, mesReferencia)

		// Deve retornar todos os dias úteis do mês
		if dias != diasUteisMes {
			t.Errorf("Esperava %d dias, mas obteve %d", diasUteisMes, dias)
		}
	})

	t.Run("DataComunicacaoDia15", func(t *testing.T) {
		// Data de comunicação no dia 15 de maio
		dataComunicacao := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)
		colaborador := &modelo.Colaborador{
			Matricula:                   "001",
			DataComunicacaoDesligamento: &dataComunicacao,
		}

		dias := AplicarRegraDesligamento(colaborador, diasUteisMes, mesReferencia)

		// Deve retornar 0 dias (regra: até dia 15 não considerar)
		if dias != 0 {
			t.Errorf("Esperava 0 dias (comunicação até dia 15), mas obteve %d", dias)
		}
	})

	t.Run("DataComunicacaoAposDia15", func(t *testing.T) {
		// Data de comunicação no dia 20 de maio
		dataComunicacao := time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC)
		colaborador := &modelo.Colaborador{
			Matricula:                   "001",
			DataComunicacaoDesligamento: &dataComunicacao,
		}

		dias := AplicarRegraDesligamento(colaborador, diasUteisMes, mesReferencia)

		// Dias trabalhados = 20
		// Dias proporcionais = (22 * 20) / 31 = 14.19... ≈ 14
		diasEsperados := (diasUteisMes * 20) / 31

		if dias != diasEsperados {
			t.Errorf("Esperava %d dias proporcionais, mas obteve %d", diasEsperados, dias)
		}
	})

	t.Run("DataComunicacaoDia1", func(t *testing.T) {
		// Data de comunicação no dia 1 de maio
		dataComunicacao := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
		colaborador := &modelo.Colaborador{
			Matricula:                   "001",
			DataComunicacaoDesligamento: &dataComunicacao,
		}

		dias := AplicarRegraDesligamento(colaborador, diasUteisMes, mesReferencia)

		// Deve retornar 0 dias (regra: até dia 15 não considerar)
		if dias != 0 {
			t.Errorf("Esperava 0 dias (comunicação até dia 15), mas obteve %d", dias)
		}
	})

	t.Run("DataComunicacaoDia31", func(t *testing.T) {
		// Data de comunicação no dia 31 de maio
		dataComunicacao := time.Date(2025, 5, 31, 0, 0, 0, 0, time.UTC)
		colaborador := &modelo.Colaborador{
			Matricula:                   "001",
			DataComunicacaoDesligamento: &dataComunicacao,
		}

		dias := AplicarRegraDesligamento(colaborador, diasUteisMes, mesReferencia)

		// Dias trabalhados = 31
		// Dias proporcionais = (22 * 31) / 31 = 22
		diasEsperados := (diasUteisMes * 31) / 31

		if dias != diasEsperados {
			t.Errorf("Esperava %d dias proporcionais, mas obteve %d", diasEsperados, dias)
		}
	})
}