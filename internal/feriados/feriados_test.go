// Package feriados provides functionality for handling Brazilian holidays
package feriados

import (
	"testing"
	"time"
)

func TestObterFeriadosNacionais(t *testing.T) {
	ano := 2025
	feriados := ObterFeriadosNacionais(ano)

	// Verificar se todos os feriados nacionais estão presentes
	expectedFeriados := []string{
		"Ano Novo",
		"Carnaval",
		"Sexta-feira Santa",
		"Páscoa",
		"Tiradentes",
		"Dia do Trabalho",
		"Corpus Christi",
		"Independência do Brasil",
		"Nossa Senhora Aparecida",
		"Finados",
		"Proclamação da República",
		"Natal",
	}

	if len(feriados) < len(expectedFeriados) {
		t.Errorf("Esperava ao menos %d feriados nacionais, obteve %d", len(expectedFeriados), len(feriados))
	}

	// Verificar algumas datas específicas
	for _, feriado := range feriados {
		switch feriado.Descricao {
		case "Ano Novo":
			if feriado.Data.Year() != ano || feriado.Data.Month() != time.January || feriado.Data.Day() != 1 {
				t.Errorf("Data incorreta para Ano Novo: %v", feriado.Data)
			}
		case "Tiradentes":
			if feriado.Data.Year() != ano || feriado.Data.Month() != time.April || feriado.Data.Day() != 21 {
				t.Errorf("Data incorreta para Tiradentes: %v", feriado.Data)
			}
		case "Dia do Trabalho":
			if feriado.Data.Year() != ano || feriado.Data.Month() != time.May || feriado.Data.Day() != 1 {
				t.Errorf("Data incorreta para Dia do Trabalho: %v", feriado.Data)
			}
		case "Independência do Brasil":
			if feriado.Data.Year() != ano || feriado.Data.Month() != time.September || feriado.Data.Day() != 7 {
				t.Errorf("Data incorreta para Independência do Brasil: %v", feriado.Data)
			}
		case "Natal":
			if feriado.Data.Year() != ano || feriado.Data.Month() != time.December || feriado.Data.Day() != 25 {
				t.Errorf("Data incorreta para Natal: %v", feriado.Data)
			}
		}
	}
}

func TestObterFeriadosEstaduais(t *testing.T) {
	ano := 2025

	// Testar feriados estaduais de São Paulo
	spFeriados := ObterFeriadosEstaduais("SP", ano)
	if len(spFeriados) == 0 {
		t.Error("Deveria haver feriados estaduais para São Paulo")
	}

	// Testar feriados estaduais do Rio de Janeiro
	rjFeriados := ObterFeriadosEstaduais("RJ", ano)
	if len(rjFeriados) == 0 {
		t.Error("Deveria haver feriados estaduais para Rio de Janeiro")
	}

	// Testar estado sem feriados estaduais específicos
	xxFeriados := ObterFeriadosEstaduais("XX", ano)
	if len(xxFeriados) != 0 {
		t.Error("Não deveria haver feriados estaduais para estado desconhecido")
	}
}

func TestEhFeriado(t *testing.T) {
	ano := 2025
	feriados := ObterFeriadosNacionais(ano)

	// Testar uma data que é feriado (Natal)
	natal := time.Date(ano, time.December, 25, 0, 0, 0, 0, time.UTC)
	if !EhFeriado(natal, feriados) {
		t.Error("25/12/2025 deveria ser considerado feriado")
	}

	// Testar uma data que não é feriado
	diaQualquer := time.Date(ano, time.June, 15, 0, 0, 0, 0, time.UTC)
	if EhFeriado(diaQualquer, feriados) {
		t.Error("15/06/2025 não deveria ser considerado feriado")
	}
}

func TestContarFeriadosNoPeriodo(t *testing.T) {
	ano := 2025
	feriados := ObterFeriadosNacionais(ano)

	// Testar contagem de feriados em janeiro
	inicio := time.Date(ano, time.January, 1, 0, 0, 0, 0, time.UTC)
	fim := time.Date(ano, time.January, 31, 0, 0, 0, 0, time.UTC)
	contagem := ContarFeriadosNoPeriodo(inicio, fim, feriados)

	// Em janeiro temos Ano Novo (01/01) e Carnaval (que varia)
	if contagem < 1 {
		t.Errorf("Deveria haver ao menos 1 feriado em janeiro, obteve %d", contagem)
	}
}

func TestCalcularPascoa(t *testing.T) {
	// Testar alguns anos conhecidos
	testes := []struct {
		ano      int
		expected time.Time
	}{
		{2025, time.Date(2025, time.April, 20, 0, 0, 0, 0, time.UTC)},
		{2026, time.Date(2026, time.April, 5, 0, 0, 0, 0, time.UTC)},
		{2027, time.Date(2027, time.March, 28, 0, 0, 0, 0, time.UTC)},
	}

	for _, teste := range testes {
		resultado := calcularPascoa(teste.ano)
		if !resultado.Equal(teste.expected) {
			t.Errorf("Para o ano %d, esperava %v, obteve %v", teste.ano, teste.expected, resultado)
		}
	}
}
