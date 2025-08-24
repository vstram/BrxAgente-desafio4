package feriados

import (
	"testing"
	"time"
)

func TestObterFeriadosMunicipaisAdicional(t *testing.T) {
	ano := 2025

	// Testar feriados municipais de São Paulo
	spFeriados := ObterFeriadosMunicipais("São Paulo", "SP", ano)
	// A função retorna 1 feriado para São Paulo
	if len(spFeriados) != 1 {
		t.Errorf("Esperava 1 feriado municipal para São Paulo, obteve %d", len(spFeriados))
	}

	// Testar feriados municipais do Rio de Janeiro
	rjFeriados := ObterFeriadosMunicipais("Rio de Janeiro", "RJ", ano)
	// A função retorna 1 feriado para Rio de Janeiro
	if len(rjFeriados) != 1 {
		t.Errorf("Esperava 1 feriado municipal para Rio de Janeiro, obteve %d", len(rjFeriados))
	}

	// Testar municípios sem feriados definidos
	outrosFeriados := ObterFeriadosMunicipais("Outro Município", "SP", ano)
	// A função retorna 0 feriados para municípios sem definição
	if len(outrosFeriados) != 0 {
		t.Errorf("Esperava 0 feriados municipais para município não definido, obteve %d", len(outrosFeriados))
	}
}

func TestCalcularPascoaAdicional(t *testing.T) {
	// Testar alguns anos conhecidos
	testes := []struct {
		ano      int
		expected time.Time
	}{
		{2025, time.Date(2025, time.April, 20, 0, 0, 0, 0, time.UTC)},
		{2026, time.Date(2026, time.April, 5, 0, 0, 0, 0, time.UTC)},
		{2027, time.Date(2027, time.March, 28, 0, 0, 0, 0, time.UTC)},
		{2028, time.Date(2028, time.April, 16, 0, 0, 0, 0, time.UTC)},
	}

	for _, teste := range testes {
		resultado := calcularPascoa(teste.ano)
		if !resultado.Equal(teste.expected) {
			t.Errorf("Para o ano %d, esperava %v, obteve %v", teste.ano, teste.expected, resultado)
		}
	}
}

func TestContarFeriadosNoPeriodoAdicional(t *testing.T) {
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

	// Testar contagem de feriados em junho (tem Corpus Christi)
	inicioJunho := time.Date(ano, time.June, 1, 0, 0, 0, 0, time.UTC)
	fimJunho := time.Date(ano, time.June, 30, 0, 0, 0, 0, time.UTC)
	contagemJunho := ContarFeriadosNoPeriodo(inicioJunho, fimJunho, feriados)

	// Em junho temos Corpus Christi (19/06)
	if contagemJunho != 1 {
		t.Errorf("Deveria haver 1 feriado em junho (Corpus Christi), obteve %d", contagemJunho)
	}
}
