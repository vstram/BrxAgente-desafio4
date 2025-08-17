package calculo

import (
	"testing"
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

func TestCalcularDiasProporcionais(t *testing.T) {
	// Data de referência para testes (maio de 2025)
	mesReferencia := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	diasUteisMes := 22 // Exemplo de dias úteis em maio

	t.Run("SemDatasQuebradas", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
		}

		dias := CalcularDiasProporcionais(colaborador, diasUteisMes, mesReferencia)

		// Deve retornar todos os dias úteis do mês
		if dias != diasUteisMes {
			t.Errorf("Esperava %d dias, mas obteve %d", diasUteisMes, dias)
		}
	})

	t.Run("AdmissaoNoMeioDoMes", func(t *testing.T) {
		// Admissão no dia 15 de maio
		dataAdmissao := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)
		colaborador := &modelo.Colaborador{
			Matricula:    "001",
			DataAdmissao: dataAdmissao,
		}

		dias := CalcularDiasProporcionais(colaborador, diasUteisMes, mesReferencia)

		// Maio tem 31 dias, então dias trabalhados = 31 - 15 + 1 = 17
		// Dias proporcionais = (22 * 17) / 31 = 12.06... ≈ 12
		diasEsperados := (diasUteisMes * (31 - 15 + 1)) / 31
		
		if dias != diasEsperados {
			t.Errorf("Esperava %d dias proporcionais, mas obteve %d", diasEsperados, dias)
		}
	})

	t.Run("AdmissaoForaDoMesReferencia", func(t *testing.T) {
		// Admissão no dia 15 de abril (mês anterior)
		dataAdmissao := time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC)
		colaborador := &modelo.Colaborador{
			Matricula:    "001",
			DataAdmissao: dataAdmissao,
		}

		dias := CalcularDiasProporcionais(colaborador, diasUteisMes, mesReferencia)

		// Deve retornar todos os dias úteis do mês
		if dias != diasUteisMes {
			t.Errorf("Esperava %d dias, mas obteve %d", diasUteisMes, dias)
		}
	})

	t.Run("DesligamentoDia15", func(t *testing.T) {
		// Desligamento no dia 15 de maio
		dataDesligamento := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)
		colaborador := &modelo.Colaborador{
			Matricula:        "001",
			DataDesligamento: &dataDesligamento,
		}

		dias := CalcularDiasProporcionais(colaborador, diasUteisMes, mesReferencia)

		// Deve retornar 0 dias (regra: até dia 15 não considerar)
		if dias != 0 {
			t.Errorf("Esperava 0 dias (desligamento até dia 15), mas obteve %d", dias)
		}
	})

	t.Run("DesligamentoAposDia15", func(t *testing.T) {
		// Desligamento no dia 20 de maio
		dataDesligamento := time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC)
		colaborador := &modelo.Colaborador{
			Matricula:        "001",
			DataDesligamento: &dataDesligamento,
		}

		dias := CalcularDiasProporcionais(colaborador, diasUteisMes, mesReferencia)

		// Dias trabalhados = 20
		// Dias proporcionais = (22 * 20) / 31 = 14.19... ≈ 14
		diasEsperados := (diasUteisMes * 20) / 31
		
		if dias != diasEsperados {
			t.Errorf("Esperava %d dias proporcionais, mas obteve %d", diasEsperados, dias)
		}
	})

	t.Run("DesligamentoForaDoMesReferencia", func(t *testing.T) {
		// Desligamento no dia 15 de abril (mês anterior)
		dataDesligamento := time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC)
		colaborador := &modelo.Colaborador{
			Matricula:        "001",
			DataDesligamento: &dataDesligamento,
		}

		dias := CalcularDiasProporcionais(colaborador, diasUteisMes, mesReferencia)

		// Deve retornar todos os dias úteis do mês
		if dias != diasUteisMes {
			t.Errorf("Esperava %d dias, mas obteve %d", diasUteisMes, dias)
		}
	})
}

func TestDiasNoMes(t *testing.T) {
	t.Run("Janeiro", func(t *testing.T) {
		data := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
		dias := diasNoMes(data)
		
		if dias != 31 {
			t.Errorf("Esperava 31 dias para janeiro, mas obteve %d", dias)
		}
	})

	t.Run("FevereiroAnoBissexto", func(t *testing.T) {
		data := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)
		dias := diasNoMes(data)
		
		if dias != 29 {
			t.Errorf("Esperava 29 dias para fevereiro de ano bissexto, mas obteve %d", dias)
		}
	})

	t.Run("FevereiroAnoNaoBissexto", func(t *testing.T) {
		data := time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC)
		dias := diasNoMes(data)
		
		if dias != 28 {
			t.Errorf("Esperava 28 dias para fevereiro de ano não bissexto, mas obteve %d", dias)
		}
	})

	t.Run("Abril", func(t *testing.T) {
		data := time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC)
		dias := diasNoMes(data)
		
		if dias != 30 {
			t.Errorf("Esperava 30 dias para abril, mas obteve %d", dias)
		}
	})
}

func TestCalcularDiasProporcionaisParaPeriodo(t *testing.T) {
	diasUteisMes := 22 // Exemplo de dias úteis em maio
	
	t.Run("PeriodoCurto", func(t *testing.T) {
		inicio := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
		fim := time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC)
		
		dias := CalcularDiasProporcionaisParaPeriodo(inicio, fim, diasUteisMes)
		
		// Período de 9 dias (1 a 10) em um mês de 31 dias
		// Dias proporcionais = (22 * 9) / 31 = 6.38... ≈ 6
		diasEsperados := (diasUteisMes * 9) / 31
		
		if dias != diasEsperados {
			t.Errorf("Esperava %d dias proporcionais, mas obteve %d", diasEsperados, dias)
		}
	})
	
	t.Run("PeriodoLongo", func(t *testing.T) {
		inicio := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
		fim := time.Date(2025, 5, 25, 0, 0, 0, 0, time.UTC)
		
		dias := CalcularDiasProporcionaisParaPeriodo(inicio, fim, diasUteisMes)
		
		// Período de 24 dias (1 a 25) em um mês de 31 dias
		// Dias proporcionais = (22 * 24) / 31 = 17.03... ≈ 17
		diasEsperados := (diasUteisMes * 24) / 31
		
		if dias != diasEsperados {
			t.Errorf("Esperava %d dias proporcionais, mas obteve %d", diasEsperados, dias)
		}
	})
}