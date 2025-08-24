package tools

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

func TestCalculateVRTool_Name(t *testing.T) {
	tool := NewCalculateVRTool()

	expected := "calculate_vr"
	if tool.Name() != expected {
		t.Errorf("Expected name '%s', got '%s'", expected, tool.Name())
	}
}

func TestCalculateVRTool_Description(t *testing.T) {
	tool := NewCalculateVRTool()

	description := tool.Description()
	if description == "" {
		t.Error("Description should not be empty")
	}

	if !strings.Contains(description, "Vale Refeição") {
		t.Error("Description should mention Vale Refeição")
	}
}

func TestCalculateVRTool_Schema(t *testing.T) {
	tool := NewCalculateVRTool()

	schema := tool.Schema()
	if schema == nil {
		t.Fatal("Schema should not be nil")
	}

	// Verificar propriedades obrigatórias
	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema should have properties")
	}

	requiredFields := []string{"colaborador", "valor_por_sindicato", "dias_uteis_por_sindicato", "mes_referencia"}
	for _, field := range requiredFields {
		if _, exists := props[field]; !exists {
			t.Errorf("Schema should have %s property", field)
		}
	}
}

func TestCalculateVRTool_Validate(t *testing.T) {
	tool := NewCalculateVRTool()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "input válido completo",
			input: `{
				"colaborador": {
					"matricula": "12345",
					"sindicato": "SINDPD"
				},
				"valor_por_sindicato": {"SINDPD": 21.25},
				"dias_uteis_por_sindicato": {"SINDPD": 22},
				"mes_referencia": "2025-09"
			}`,
			wantErr: false,
		},
		{
			name: "JSON inválido",
			input: `{
				"colaborador": {
					"matricula": "12345"
				}`,
			wantErr: true,
		},
		{
			name: "matrícula vazia",
			input: `{
				"colaborador": {
					"matricula": "",
					"sindicato": "SINDPD"
				},
				"valor_por_sindicato": {"SINDPD": 21.25},
				"dias_uteis_por_sindicato": {"SINDPD": 22},
				"mes_referencia": "2025-09"
			}`,
			wantErr: true,
		},
		{
			name: "sindicato vazio",
			input: `{
				"colaborador": {
					"matricula": "12345",
					"sindicato": ""
				},
				"valor_por_sindicato": {"SINDPD": 21.25},
				"dias_uteis_por_sindicato": {"SINDPD": 22},
				"mes_referencia": "2025-09"
			}`,
			wantErr: true,
		},
		{
			name: "mês referência inválido",
			input: `{
				"colaborador": {
					"matricula": "12345",
					"sindicato": "SINDPD"
				},
				"valor_por_sindicato": {"SINDPD": 21.25},
				"dias_uteis_por_sindicato": {"SINDPD": 22},
				"mes_referencia": "2025-13"
			}`,
			wantErr: true,
		},
		{
			name: "valor_por_sindicato vazio",
			input: `{
				"colaborador": {
					"matricula": "12345",
					"sindicato": "SINDPD"
				},
				"valor_por_sindicato": {},
				"dias_uteis_por_sindicato": {"SINDPD": 22},
				"mes_referencia": "2025-09"
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tool.Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalculateVRTool_Execute_Success(t *testing.T) {
	tool := NewCalculateVRTool()

	// Criar colaborador de teste
	// ⚠️ CONFIDENCIALIDADE: Remove campo Nome conforme PRD.md
	colaborador := modelo.Colaborador{
		Matricula:    "12345",
		DataAdmissao: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Sindicato:    "SINDPD",
		Empresa:      "Empresa Teste",
	}

	input := map[string]interface{}{
		"colaborador": colaborador,
		"valor_por_sindicato": map[string]float64{
			"SINDPD": 21.25,
		},
		"dias_uteis_por_sindicato": map[string]int{
			"SINDPD": 22,
		},
		"mes_referencia": "2025-09",
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	result, err := tool.Execute(string(inputJSON))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Parse result
	var output CalculateVROutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verificar resultados básicos
	if !output.Success {
		t.Errorf("Output should indicate success, got error: %s", output.Error)
	}

	if output.Matricula != "12345" {
		t.Errorf("Expected matricula '12345', got '%s'", output.Matricula)
	}

	if output.Sindicato != "SINDPD" {
		t.Errorf("Expected sindicato 'SINDPD', got '%s'", output.Sindicato)
	}

	if output.MesReferencia != "2025-09" {
		t.Errorf("Expected mes_referencia '2025-09', got '%s'", output.MesReferencia)
	}

	// Os valores calculados podem ser diferentes devido às regras de negócio (afastamentos, etc.)
	// Apenas verificar se são valores razoáveis
	if output.ValorTotal <= 0 {
		t.Error("ValorTotal should be greater than 0")
	}

	if output.ValorEmpresa <= 0 {
		t.Error("ValorEmpresa should be greater than 0")
	}

	if output.ValorColaborador <= 0 {
		t.Error("ValorColaborador should be greater than 0")
	}

	// Verificar proporção 80/20
	expectedEmpresa := output.ValorTotal * 0.8
	expectedColaborador := output.ValorTotal * 0.2

	if output.ValorEmpresa != expectedEmpresa {
		t.Errorf("Expected valor_empresa %.2f, got %.2f", expectedEmpresa, output.ValorEmpresa)
	}

	if output.ValorColaborador != expectedColaborador {
		t.Errorf("Expected valor_colaborador %.2f, got %.2f", expectedColaborador, output.ValorColaborador)
	}

	if output.DiasUteisEfetivos == 0 {
		t.Error("DiasUteisEfetivos should not be zero")
	}

	if output.ValorPorDia != 21.25 {
		t.Errorf("Expected valor_por_dia 21.25, got %.2f", output.ValorPorDia)
	}
}

func TestCalculateVRTool_Execute_WithAfastamentos(t *testing.T) {
	tool := NewCalculateVRTool()

	// Colaborador com afastamentos
	// ⚠️ CONFIDENCIALIDADE: Remove campo Nome conforme PRD.md
	colaborador := modelo.Colaborador{
		Matricula:    "12346",
		DataAdmissao: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Sindicato:    "SINDPD",
		Empresa:      "Empresa Teste",
		Afastamentos: []modelo.Periodo{
			{
				Inicio: time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC),
				Fim:    time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	input := map[string]interface{}{
		"colaborador": colaborador,
		"valor_por_sindicato": map[string]float64{
			"SINDPD": 21.25,
		},
		"dias_uteis_por_sindicato": map[string]int{
			"SINDPD": 22,
		},
		"mes_referencia": "2025-09",
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	result, err := tool.Execute(string(inputJSON))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Parse result
	var output CalculateVROutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verificar que detectou afastamentos
	if !output.Detalhes.TemAfastamentos {
		t.Error("Should detect afastamentos")
	}

	if output.Detalhes.QuantidadeAfastamentos != 1 {
		t.Errorf("Expected 1 afastamento, got %d", output.Detalhes.QuantidadeAfastamentos)
	}

	if len(output.Detalhes.Observacoes) == 0 {
		t.Error("Should have observacoes about afastamentos")
	}
}

func TestCalculateVRTool_Execute_InvalidInput(t *testing.T) {
	tool := NewCalculateVRTool()

	input := `{"invalid": "json structure"}`

	result, err := tool.Execute(input)
	if err != nil {
		t.Fatalf("Execute() should not return error, got %v", err)
	}

	// Parse result
	var output CalculateVROutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verificar que marcou como falha
	if output.Success {
		t.Error("Output should indicate failure for invalid input")
	}

	if output.Error == "" {
		t.Error("Output should contain error message")
	}
}

func TestCalculateVRTool_Execute_InvalidMesReferencia(t *testing.T) {
	tool := NewCalculateVRTool()

	input := `{
		"colaborador": {
			"matricula": "12345",
			"sindicato": "SINDPD"
		},
		"valor_por_sindicato": {"SINDPD": 21.25},
		"dias_uteis_por_sindicato": {"SINDPD": 22},
		"mes_referencia": "invalid-date"
	}`

	result, err := tool.Execute(input)
	if err != nil {
		t.Fatalf("Execute() should not return error, got %v", err)
	}

	// Parse result
	var output CalculateVROutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verificar que marcou como falha
	if output.Success {
		t.Error("Output should indicate failure for invalid date")
	}

	if !strings.Contains(output.Error, "formato") {
		t.Error("Error should mention invalid format")
	}
}
