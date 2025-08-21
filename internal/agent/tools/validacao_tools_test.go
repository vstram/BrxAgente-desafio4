package tools

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
	
	"BrxAgente-desafio4/internal/modelo"
)

func TestValidateDataTool_Name(t *testing.T) {
	tool := NewValidateDataTool()
	
	expected := "validate_data"
	if tool.Name() != expected {
		t.Errorf("Expected name '%s', got '%s'", expected, tool.Name())
	}
}

func TestValidateDataTool_Description(t *testing.T) {
	tool := NewValidateDataTool()
	
	description := tool.Description()
	if description == "" {
		t.Error("Description should not be empty")
	}
	
	if !strings.Contains(description, "Valida") {
		t.Error("Description should mention validation")
	}
}

func TestValidateDataTool_Schema(t *testing.T) {
	tool := NewValidateDataTool()
	
	schema := tool.Schema()
	if schema == nil {
		t.Fatal("Schema should not be nil")
	}
	
	// Verificar propriedades obrigatórias
	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema should have properties")
	}
	
	if _, exists := props["tipo_validacao"]; !exists {
		t.Error("Schema should have tipo_validacao property")
	}
}

func TestValidateDataTool_Validate(t *testing.T) {
	tool := NewValidateDataTool()
	
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "validação de colaborador válida",
			input: `{
				"tipo_validacao": "colaborador",
				"colaborador": {
					"matricula": "12345",
					"sindicato": "SINDPD"
				},
				"validar_campos_obrigatorios": true
			}`,
			wantErr: false,
		},
		{
			name: "validação de lista válida",
			input: `{
				"tipo_validacao": "lista_colaboradores",
				"colaboradores": [
					{"matricula": "12345", "sindicato": "SINDPD"}
				],
				"validar_campos_obrigatorios": true
			}`,
			wantErr: false,
		},
		{
			name: "validação de planilha válida",
			input: `{
				"tipo_validacao": "planilha",
				"dados_planilha": {
					"headers": ["Matricula", "Nome"],
					"row_count": 10
				}
			}`,
			wantErr: false,
		},
		{
			name: "JSON inválido",
			input: `{
				"tipo_validacao": "colaborador"`,
			wantErr: true,
		},
		{
			name: "tipo de validação inválido",
			input: `{
				"tipo_validacao": "invalid_type"
			}`,
			wantErr: true,
		},
		{
			name: "colaborador ausente para tipo colaborador",
			input: `{
				"tipo_validacao": "colaborador",
				"validar_campos_obrigatorios": true
			}`,
			wantErr: true,
		},
		{
			name: "lista vazia para tipo lista_colaboradores",
			input: `{
				"tipo_validacao": "lista_colaboradores",
				"colaboradores": []
			}`,
			wantErr: true,
		},
		{
			name: "dados_planilha ausente para tipo planilha",
			input: `{
				"tipo_validacao": "planilha"
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

func TestValidateDataTool_Execute_ColaboradorValido(t *testing.T) {
	tool := NewValidateDataTool()
	
	// Colaborador válido
	colaborador := modelo.Colaborador{
		Matricula:    "12345",
		Nome:         "Teste Colaborador",
		DataAdmissao: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Sindicato:    "SINDPD",
		Empresa:      "Empresa Teste",
		Cargo:        "Analista",
		Situacao:     "Ativo",
	}
	
	input := map[string]interface{}{
		"tipo_validacao":             "colaborador",
		"colaborador":                colaborador,
		"validar_campos_obrigatorios": true,
		"validar_datas":              true,
		"validar_consistencia":       true,
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
	var output ValidateDataOutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}
	
	// Verificar resultados
	if !output.Success {
		t.Errorf("Output should indicate success, got error: %s", output.Error)
	}
	
	if output.TipoValidacao != "colaborador" {
		t.Errorf("Expected tipo_validacao 'colaborador', got '%s'", output.TipoValidacao)
	}
	
	if output.TotalRegistros != 1 {
		t.Errorf("Expected 1 total registro, got %d", output.TotalRegistros)
	}
	
	if output.RegistrosValidos != 1 {
		t.Errorf("Expected 1 valid registro, got %d", output.RegistrosValidos)
	}
	
	if output.RegistrosComErro != 0 {
		t.Errorf("Expected 0 error registros, got %d", output.RegistrosComErro)
	}
	
	if output.Resumo.StatusGeral != "aprovado" {
		t.Errorf("Expected status 'aprovado', got '%s'", output.Resumo.StatusGeral)
	}
}

func TestValidateDataTool_Execute_ColaboradorComErros(t *testing.T) {
	tool := NewValidateDataTool()
	
	// Colaborador com erros
	colaborador := modelo.Colaborador{
		Matricula:    "", // Erro: matrícula vazia
		Nome:         "Teste Colaborador",
		DataAdmissao: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC), // Erro: data futura
		Sindicato:    "", // Erro: sindicato vazio
		Empresa:      "", // Erro: empresa vazia
	}
	
	input := map[string]interface{}{
		"tipo_validacao":             "colaborador",
		"colaborador":                colaborador,
		"validar_campos_obrigatorios": true,
		"validar_datas":              true,
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
	var output ValidateDataOutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}
	
	// Verificar que encontrou erros
	if output.RegistrosComErro == 0 {
		t.Error("Should have found errors in invalid colaborador")
	}
	
	if len(output.Erros) == 0 {
		t.Error("Should have error details")
	}
	
	if output.Resumo.StatusGeral == "aprovado" {
		t.Error("Should not be approved with errors")
	}
	
	// Verificar se encontrou erros críticos
	if output.Resumo.ErrosCriticos == 0 {
		t.Error("Should have found critical errors")
	}
}

func TestValidateDataTool_Execute_ListaColaboradores(t *testing.T) {
	tool := NewValidateDataTool()
	
	// Lista com colaboradores válidos e inválidos
	colaboradores := []modelo.Colaborador{
		{
			Matricula:    "12345",
			Nome:         "Colaborador Válido",
			DataAdmissao: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Sindicato:    "SINDPD",
			Empresa:      "Empresa Teste",
			Cargo:        "Analista",
			Situacao:     "Ativo",
		},
		{
			Matricula:    "", // Inválido
			Nome:         "Colaborador Inválido",
			Sindicato:    "SINDPD",
		},
		{
			Matricula:    "12345", // Duplicata
			Nome:         "Colaborador Duplicado",
			Sindicato:    "SINDPD",
		},
	}
	
	input := map[string]interface{}{
		"tipo_validacao":             "lista_colaboradores",
		"colaboradores":              colaboradores,
		"validar_campos_obrigatorios": true,
		"validar_duplicatas":         true,
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
	var output ValidateDataOutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}
	
	// Verificar resultados
	if output.TotalRegistros != 3 {
		t.Errorf("Expected 3 total registros, got %d", output.TotalRegistros)
	}
	
	if output.RegistrosComErro == 0 {
		t.Error("Should have found errors")
	}
	
	// Verificar se detectou duplicata
	foundDuplicate := false
	for _, erro := range output.Erros {
		if erro.TipoErro == "duplicata" {
			foundDuplicate = true
			break
		}
	}
	if !foundDuplicate {
		t.Error("Should have detected duplicate matricula")
	}
}

func TestValidateDataTool_Execute_Planilha(t *testing.T) {
	tool := NewValidateDataTool()
	
	// Dados de planilha com headers válidos
	dadosPlanilha := map[string]interface{}{
		"headers": []interface{}{
			"Matricula", "Nome", "DataAdmissao", "Empresa", "Sindicato",
		},
		"row_count": 100.0,
		"col_count": 5.0,
	}
	
	input := map[string]interface{}{
		"tipo_validacao": "planilha",
		"dados_planilha": dadosPlanilha,
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
	var output ValidateDataOutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}
	
	// Verificar resultados
	if !output.Success {
		t.Errorf("Output should indicate success, got error: %s", output.Error)
	}
	
	if output.TotalRegistros != 100 {
		t.Errorf("Expected 100 total registros, got %d", output.TotalRegistros)
	}
	
	// Como todos os headers obrigatórios estão presentes, não deve ter erros críticos
	if output.Resumo.ErrosCriticos > 0 {
		t.Errorf("Should not have critical errors with valid headers")
	}
}

func TestValidateDataTool_Execute_PlanilhaComHeadersAusentes(t *testing.T) {
	tool := NewValidateDataTool()
	
	// Dados de planilha com headers ausentes
	dadosPlanilha := map[string]interface{}{
		"headers": []interface{}{
			"Matricula", "Nome", // Faltam headers obrigatórios
		},
		"row_count": 50.0,
	}
	
	input := map[string]interface{}{
		"tipo_validacao": "planilha",
		"dados_planilha": dadosPlanilha,
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
	var output ValidateDataOutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}
	
	// Verificar que encontrou erros
	if len(output.Erros) == 0 {
		t.Error("Should have found missing header errors")
	}
	
	// Verificar se encontrou erros de headers ausentes
	foundHeaderError := false
	for _, erro := range output.Erros {
		if erro.TipoErro == "header_ausente" {
			foundHeaderError = true
			break
		}
	}
	if !foundHeaderError {
		t.Error("Should have detected missing headers")
	}
}

func TestValidateDataTool_Execute_InvalidInput(t *testing.T) {
	tool := NewValidateDataTool()
	
	input := `{"invalid": "structure"}`
	
	result, err := tool.Execute(input)
	if err != nil {
		t.Fatalf("Execute() should not return error, got %v", err)
	}
	
	// Parse result
	var output ValidateDataOutput
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