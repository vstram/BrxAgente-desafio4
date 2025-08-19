// Package calculo provides functionality for calculating VR values
package calculo

import (
	"testing"
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

// BenchmarkCalcularDiasUteisPorSindicato benchmarks the CalcularDiasUteisPorSindicato function
func BenchmarkCalcularDiasUteisPorSindicato(b *testing.B) {
	// Criar um colaborador de exemplo para o benchmark
	colaborador := &modelo.Colaborador{
		Matricula: "001",
		Empresa:   "Empresa A",
		Cargo:     "Analista",
		Sindicato: "Paraná",
		Ferias: []modelo.Periodo{
			{
				Inicio: time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC),
				Fim:    time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC),
			},
		},
		Afastamentos: []modelo.Periodo{
			{
				Inicio: time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC),
				Fim:    time.Date(2025, 5, 22, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	// Data de referência para cálculo (maio de 2025)
	mesReferencia := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)

	// Número de dias úteis do sindicato (exemplo para Paraná)
	diasUteisSindicato := 22

	// Executar o benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalcularDiasUteisPorSindicato(colaborador, diasUteisSindicato, mesReferencia)
	}
}

// BenchmarkCalcularDiasFerias benchmarks the calcularDiasFerias function
func BenchmarkCalcularDiasFerias(b *testing.B) {
	// Criar um colaborador de exemplo com várias férias para o benchmark
	colaborador := &modelo.Colaborador{
		Matricula: "001",
		Empresa:   "Empresa A",
		Cargo:     "Analista",
		Sindicato: "Paraná",
		Ferias: []modelo.Periodo{
			{
				Inicio: time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
				Fim:    time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			},
			{
				Inicio: time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC),
				Fim:    time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC),
			},
			{
				Inicio: time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC),
				Fim:    time.Date(2025, 5, 25, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	// Data de referência para cálculo (maio de 2025)
	mesReferencia := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)

	// Número de dias úteis do sindicato (exemplo para Paraná)
	diasUteisSindicato := 22

	// Executar o benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calcularDiasFerias(colaborador, diasUteisSindicato, mesReferencia)
	}
}

// BenchmarkCalcularDiasAfastamentos benchmarks the calcularDiasAfastamentos function
func BenchmarkCalcularDiasAfastamentos(b *testing.B) {
	// Criar um colaborador de exemplo com vários afastamentos para o benchmark
	colaborador := &modelo.Colaborador{
		Matricula: "001",
		Empresa:   "Empresa A",
		Cargo:     "Analista",
		Sindicato: "Paraná",
		Afastamentos: []modelo.Periodo{
			{
				Inicio: time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
				Fim:    time.Date(2025, 5, 3, 0, 0, 0, 0, time.UTC),
			},
			{
				Inicio: time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC),
				Fim:    time.Date(2025, 5, 12, 0, 0, 0, 0, time.UTC),
			},
			{
				Inicio: time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC),
				Fim:    time.Date(2025, 5, 22, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	// Data de referência para cálculo (maio de 2025)
	mesReferencia := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)

	// Número de dias úteis do sindicato (exemplo para Paraná)
	diasUteisSindicato := 22

	// Executar o benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calcularDiasAfastamentos(colaborador, diasUteisSindicato, mesReferencia)
	}
}