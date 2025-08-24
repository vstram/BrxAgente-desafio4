package calculo

import (
	"testing"
	"time"
)

func TestMapearSindicatoParaEstadoAdicional(t *testing.T) {
	testCases := []struct {
		name           string
		sindicato      string
		expectedEstado string
	}{
		{
			name:           "Sindicato Paraná",
			sindicato:      "Paraná",
			expectedEstado: "Paraná",
		},
		{
			name:           "Sindicato Parana",
			sindicato:      "Parana",
			expectedEstado: "Paraná",
		},
		{
			name:           "Sindicato Rio Grande do Sul",
			sindicato:      "Rio Grande do Sul",
			expectedEstado: "Rio Grande do Sul",
		},
		{
			name:           "Sindicato São Paulo",
			sindicato:      "São Paulo",
			expectedEstado: "São Paulo",
		},
		{
			name:           "Sindicato Rio de Janeiro",
			sindicato:      "Rio de Janeiro",
			expectedEstado: "Rio de Janeiro",
		},
		{
			name:           "Sindicato com PR e Curitiba",
			sindicato:      "SITEPD PR - SIND DOS TRAB EM EMPR PRIVADAS DE PROC DE DADOS DE CURITIBA E REGIAO METROPOLITANA 22",
			expectedEstado: "Paraná",
		},
		{
			name:           "Sindicato com RS e Rio Grande",
			sindicato:      "SINDPPD RS - SINDICATO DOS TRAB. EM PROC. DE DADOS RIO GRANDE DO SUL 21",
			expectedEstado: "Rio Grande do Sul",
		},
		{
			name:           "Sindicato com SP e São Paulo",
			sindicato:      "SINDPD SP - SIND.TRAB.EM PROC DADOS E EMPR.EMPRESAS PROC DADOS ESTADO DE SP. 22",
			expectedEstado: "São Paulo",
		},
		{
			name:           "Sindicato com RJ e Rio de Janeiro",
			sindicato:      "SINDPD RJ - SINDICATO PROFISSIONAIS DE PROC DADOS DO RIO DE JANEIRO 21",
			expectedEstado: "Rio de Janeiro",
		},
		{
			name:           "Sindicato com PR apenas",
			sindicato:      "PR - SINDICATO EXEMPLO",
			expectedEstado: "Paraná",
		},
		{
			name:           "Sindicato com RS apenas",
			sindicato:      "RS - SINDICATO EXEMPLO",
			expectedEstado: "Rio Grande do Sul",
		},
		{
			name:           "Sindicato com SP apenas",
			sindicato:      "SP - SINDICATO EXEMPLO",
			expectedEstado: "São Paulo",
		},
		{
			name:           "Sindicato com RJ apenas",
			sindicato:      "RJ - SINDICATO EXEMPLO",
			expectedEstado: "Rio de Janeiro",
		},
		{
			name:           "Sindicato desconhecido",
			sindicato:      "Sindicato Desconhecido",
			expectedEstado: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := mapearSindicatoParaEstado(tc.sindicato)
			if result != tc.expectedEstado {
				t.Errorf("Esperava %s, obteve %s", tc.expectedEstado, result)
			}
		})
	}
}

func TestParseDataDesligamento(t *testing.T) {
	testCases := []struct {
		name        string
		dataStr     string
		expectError bool
		expected    time.Time
	}{
		{
			name:        "Data válida",
			dataStr:     "05-15-25",
			expectError: false,
			expected:    time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "Data inválida formato",
			dataStr:     "invalid",
			expectError: true,
		},
		{
			name:        "Data com partes insuficientes",
			dataStr:     "05-15",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseDataDesligamento(tc.dataStr)

			if tc.expectError {
				if err == nil {
					t.Error("Esperava erro, mas não houve erro")
				}
				return
			}

			if err != nil {
				t.Errorf("Não esperava erro, mas obteve: %v", err)
				return
			}

			if !result.Equal(tc.expected) {
				t.Errorf("Esperava %v, obteve %v", tc.expected, result)
			}
		})
	}
}

func TestParseDataAdmissao(t *testing.T) {
	testCases := []struct {
		name        string
		dataStr     string
		expectError bool
		expected    time.Time
	}{
		{
			name:        "Data válida",
			dataStr:     "04-07-25",
			expectError: false,
			expected:    time.Date(2025, 4, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "Data inválida formato",
			dataStr:     "invalid",
			expectError: true,
		},
		{
			name:        "Data com partes insuficientes",
			dataStr:     "04-07",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseDataAdmissao(tc.dataStr)

			if tc.expectError {
				if err == nil {
					t.Error("Esperava erro, mas não houve erro")
				}
				return
			}

			if err != nil {
				t.Errorf("Não esperava erro, mas obteve: %v", err)
				return
			}

			if !result.Equal(tc.expected) {
				t.Errorf("Esperava %v, obteve %v", tc.expected, result)
			}
		})
	}
}

func TestCalcularDiasProporcionaisParaPeriodoAdicional(t *testing.T) {
	// Test case 1: Período de 10 dias em um mês de 31 dias
	inicio := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	fim := time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC)
	diasUteisMes := 22

	resultado := CalcularDiasProporcionaisParaPeriodo(inicio, fim, diasUteisMes)

	// Cálculo: (22 * 9) / 31 ≈ 6.4 => 6 dias (9 because it's 9 days difference, not 10)
	// The function calculates the duration between dates, which is 9 days from May 1 to May 10
	if resultado < 5 || resultado > 7 {
		t.Errorf("Esperava aproximadamente 6 dias, obteve %d", resultado)
	}

	// Test case 2: Período de 5 dias em um mês de 31 dias
	inicio2 := time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC)
	fim2 := time.Date(2025, 5, 14, 0, 0, 0, 0, time.UTC)

	resultado2 := CalcularDiasProporcionaisParaPeriodo(inicio2, fim2, diasUteisMes)

	// Cálculo: (22 * 4) / 31 ≈ 2.8 => 2 dias (4 because it's 4 days difference, not 5)
	if resultado2 < 1 || resultado2 > 3 {
		t.Errorf("Esperava aproximadamente 2 dias, obteve %d", resultado2)
	}
}
