// Package calculo provides functionality for calculating VR values
package calculo

import (
	"testing"
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

func TestCalcularDiasUteisPorSindicato(t *testing.T) {
	// Data de referência para testes (maio de 2025)
	mesReferencia := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)

	// Test case 1: Colaborador sem férias ou afastamentos
	colaborador1 := &modelo.Colaborador{
		Matricula:    "001",
		Ferias:       []modelo.Periodo{},
		Afastamentos: []modelo.Periodo{},
		Sindicato:    "Paraná", // Specify sindicato to test holiday calculation
	}

	diasUteisSindicato := 22 // Exemplo para Paraná
	expectedDiasUteis := 21  // 22 - 1 feriado nacional (Dia do Trabalho)

	resultado1 := CalcularDiasUteisPorSindicato(colaborador1, diasUteisSindicato, mesReferencia)

	// Deve retornar todos os dias úteis do sindicato menos os feriados
	if resultado1 != expectedDiasUteis {
		t.Errorf("Esperava %d dias úteis (22 - 1 feriado), obteve %d", expectedDiasUteis, resultado1)
	}

	// Test case 2: Colaborador com férias parciais
	dataInicioFerias := time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC)
	dataFimFerias := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)

	colaborador2 := &modelo.Colaborador{
		Matricula: "002",
		Sindicato: "Paraná",
		Ferias: []modelo.Periodo{
			{Inicio: dataInicioFerias, Fim: dataFimFerias},
		},
		Afastamentos: []modelo.Periodo{},
	}

	// 6 dias de férias em um mês de 31 dias
	// Total de dias úteis: 22 - 1 feriado = 21
	// Dias descontados por férias: (21 * 6) / 31 ≈ 4.06 => 4 dias descontados
	// Resultado esperado: 21 - 4 = 17
	resultado2 := CalcularDiasUteisPorSindicato(colaborador2, diasUteisSindicato, mesReferencia)

	// Verificando se o resultado está próximo do esperado
	if resultado2 < 16 || resultado2 > 18 {
		t.Errorf("Esperava aproximadamente 17 dias úteis, obteve %d", resultado2)
	}
}

func TestCalcularDiasFerias(t *testing.T) {
	// Data de referência para testes (maio de 2025)
	mesReferencia := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	diasUteisSindicato := 22

	// Test case 1: Sem férias no mês
	colaborador1 := &modelo.Colaborador{
		Matricula: "001",
		Ferias:    []modelo.Periodo{},
	}

	resultado1 := calcularDiasFerias(colaborador1, diasUteisSindicato, mesReferencia)

	if resultado1 != 0 {
		t.Errorf("Esperava 0 dias de férias, obteve %d", resultado1)
	}

	// Test case 2: Com férias no mês
	dataInicioFerias := time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC)
	dataFimFerias := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)

	colaborador2 := &modelo.Colaborador{
		Matricula: "002",
		Ferias: []modelo.Periodo{
			{Inicio: dataInicioFerias, Fim: dataFimFerias},
		},
	}

	// 6 dias de férias em um mês de 31 dias
	// Proporção: (22 * 6) / 31 ≈ 4.26 => 4 dias descontados
	resultado2 := calcularDiasFerias(colaborador2, diasUteisSindicato, mesReferencia)

	// Verificando se o resultado está próximo do esperado
	if resultado2 < 3 || resultado2 > 5 {
		t.Errorf("Esperava aproximadamente 4 dias de férias descontados, obteve %d", resultado2)
	}
}

func TestCalcularDiasAfastamentos(t *testing.T) {
	// Data de referência para testes (maio de 2025)
	mesReferencia := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	diasUteisSindicato := 22

	// Test case 1: Sem afastamentos no mês
	colaborador1 := &modelo.Colaborador{
		Matricula:    "001",
		Afastamentos: []modelo.Periodo{},
	}

	resultado1 := calcularDiasAfastamentos(colaborador1, diasUteisSindicato, mesReferencia)

	if resultado1 != 0 {
		t.Errorf("Esperava 0 dias de afastamento, obteve %d", resultado1)
	}

	// Test case 2: Com afastamentos no mês
	dataInicioAfastamento := time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC)
	dataFimAfastamento := time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC)

	colaborador2 := &modelo.Colaborador{
		Matricula: "002",
		Afastamentos: []modelo.Periodo{
			{Inicio: dataInicioAfastamento, Fim: dataFimAfastamento},
		},
	}

	// 6 dias de afastamento em um mês de 31 dias
	// Proporção: (22 * 6) / 31 ≈ 4.26 => 4 dias descontados
	resultado2 := calcularDiasAfastamentos(colaborador2, diasUteisSindicato, mesReferencia)

	// Verificando se o resultado está próximo do esperado
	if resultado2 < 3 || resultado2 > 5 {
		t.Errorf("Esperava aproximadamente 4 dias de afastamento descontados, obteve %d", resultado2)
	}
}

func TestPeriodoIntersectaMes(t *testing.T) {
	// Data de referência para testes (maio de 2025)
	mesReferencia := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)

	// Test case 1: Período antes do mês de referência
	periodo1 := modelo.Periodo{
		Inicio: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		Fim:    time.Date(2025, 4, 30, 0, 0, 0, 0, time.UTC),
	}

	if periodoIntersectaMes(periodo1, mesReferencia) {
		t.Error("Período antes do mês de referência não deveria intersectar")
	}

	// Test case 2: Período depois do mês de referência
	periodo2 := modelo.Periodo{
		Inicio: time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		Fim:    time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	}

	if periodoIntersectaMes(periodo2, mesReferencia) {
		t.Error("Período depois do mês de referência não deveria intersectar")
	}

	// Test case 3: Período durante o mês de referência
	periodo3 := modelo.Periodo{
		Inicio: time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC),
		Fim:    time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC),
	}

	if !periodoIntersectaMes(periodo3, mesReferencia) {
		t.Error("Período durante o mês de referência deveria intersectar")
	}

	// Test case 4: Período que começa antes e termina durante o mês
	periodo4 := modelo.Periodo{
		Inicio: time.Date(2025, 4, 25, 0, 0, 0, 0, time.UTC),
		Fim:    time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
	}

	if !periodoIntersectaMes(periodo4, mesReferencia) {
		t.Error("Período que começa antes e termina durante o mês deveria intersectar")
	}

	// Test case 5: Período que começa durante e termina depois do mês
	periodo5 := modelo.Periodo{
		Inicio: time.Date(2025, 5, 25, 0, 0, 0, 0, time.UTC),
		Fim:    time.Date(2025, 6, 5, 0, 0, 0, 0, time.UTC),
	}

	if !periodoIntersectaMes(periodo5, mesReferencia) {
		t.Error("Período que começa durante e termina depois do mês deveria intersectar")
	}
}
