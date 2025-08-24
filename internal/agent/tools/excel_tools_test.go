package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadExcelTool_Name(t *testing.T) {
	tool := NewReadExcelTool()

	expected := "read_excel"
	if tool.Name() != expected {
		t.Errorf("Expected name '%s', got '%s'", expected, tool.Name())
	}
}

func TestReadExcelTool_Description(t *testing.T) {
	tool := NewReadExcelTool()

	description := tool.Description()
	if description == "" {
		t.Error("Description should not be empty")
	}

	if !strings.Contains(description, "Excel") {
		t.Error("Description should mention Excel")
	}
}

func TestReadExcelTool_Schema(t *testing.T) {
	tool := NewReadExcelTool()

	schema := tool.Schema()
	if schema == nil {
		t.Fatal("Schema should not be nil")
	}

	// Verificar se tem propriedades obrigatórias
	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema should have properties")
	}

	if _, exists := props["file_path"]; !exists {
		t.Error("Schema should have file_path property")
	}
}

func TestReadExcelTool_Validate(t *testing.T) {
	tool := NewReadExcelTool()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "input válido",
			input:   `{"file_path": "test.xlsx"}`,
			wantErr: false,
		},
		{
			name:    "input válido com sheet",
			input:   `{"file_path": "test.xlsx", "sheet": "Sheet1"}`,
			wantErr: false,
		},
		{
			name:    "input válido com max_rows",
			input:   `{"file_path": "test.xlsx", "max_rows": 100}`,
			wantErr: false,
		},
		{
			name:    "JSON inválido",
			input:   `{"file_path": "test.xlsx"`,
			wantErr: true,
		},
		{
			name:    "file_path vazio",
			input:   `{"file_path": ""}`,
			wantErr: true,
		},
		{
			name:    "file_path ausente",
			input:   `{"sheet": "Sheet1"}`,
			wantErr: true,
		},
		{
			name:    "max_rows negativo",
			input:   `{"file_path": "test.xlsx", "max_rows": -1}`,
			wantErr: true,
		},
		{
			name:    "max_rows muito alto",
			input:   `{"file_path": "test.xlsx", "max_rows": 1001}`,
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

func TestReadExcelTool_Execute_FileNotFound(t *testing.T) {
	tool := NewReadExcelTool()

	input := `{"file_path": "arquivo_inexistente.xlsx"}`

	result, err := tool.Execute(input)
	if err != nil {
		t.Fatalf("Execute() should not return error, got %v", err)
	}

	// Parse result
	var output ReadExcelOutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verificar se marcou como falha
	if output.Success {
		t.Error("Output should indicate failure for non-existent file")
	}

	if output.Error == "" {
		t.Error("Output should contain error message")
	}
}

func TestReadExcelTool_Execute_WithTestFile(t *testing.T) {
	// Criar arquivo de teste temporário
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "ATIVOS.xlsx")

	// Criar arquivo vazio (simulando Excel)
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	tool := NewReadExcelTool()

	input := `{"file_path": "` + strings.ReplaceAll(testFile, `\`, `\\`) + `"}`

	result, err := tool.Execute(input)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Parse result
	var output ReadExcelOutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verificar output básico
	if !output.Success {
		t.Errorf("Output should indicate success, got error: %s", output.Error)
	}

	if output.FilePath != testFile {
		t.Errorf("Expected file_path '%s', got '%s'", testFile, output.FilePath)
	}

	if output.Summary == "" {
		t.Error("Summary should not be empty")
	}
}

func TestReadExcelTool_Execute_UnsupportedFormat(t *testing.T) {
	// Criar arquivo de teste com extensão não suportada
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "test.txt")

	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	tool := NewReadExcelTool()

	input := `{"file_path": "` + strings.ReplaceAll(testFile, `\`, `\\`) + `"}`

	result, err := tool.Execute(input)
	if err != nil {
		t.Fatalf("Execute() should not return error, got %v", err)
	}

	// Parse result
	var output ReadExcelOutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verificar se marcou como falha
	if output.Success {
		t.Error("Output should indicate failure for unsupported format")
	}

	if !strings.Contains(output.Error, "formato") {
		t.Error("Error should mention unsupported format")
	}
}

func TestReadExcelTool_Execute_SpecialFiles(t *testing.T) {
	// Criar arquivos de teste para casos especiais
	testDir := t.TempDir()

	testCases := []struct {
		fileName     string
		expectedRows int
		expectedCols int
	}{
		{"ATIVOS.xlsx", 1247, 7},     // ⚠️ CONFIDENCIALIDADE: Reduzido de 8 para 7 (removido Nome)
		{"DESLIGADOS.xlsx", 23, 8},   // ⚠️ CONFIDENCIALIDADE: Reduzido de 9 para 8 (removido Nome)
		{"FERIAS.xlsx", 156, 4},      // ⚠️ CONFIDENCIALIDADE: Reduzido de 5 para 4 (removido Nome)
		{"AFASTAMENTOS.xlsx", 89, 5}, // ⚠️ CONFIDENCIALIDADE: Reduzido de 6 para 5 (removido Nome)
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			testFile := filepath.Join(testDir, tc.fileName)

			// Criar arquivo
			file, err := os.Create(testFile)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
			file.Close()

			tool := NewReadExcelTool()
			input := `{"file_path": "` + strings.ReplaceAll(testFile, `\`, `\\`) + `"}`

			result, err := tool.Execute(input)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			// Parse result
			var output ReadExcelOutput
			if err := json.Unmarshal([]byte(result), &output); err != nil {
				t.Fatalf("Failed to parse output: %v", err)
			}

			// Verificar resultados específicos
			if !output.Success {
				t.Errorf("Output should indicate success, got error: %s", output.Error)
			}

			if output.RowCount != tc.expectedRows {
				t.Errorf("Expected %d rows, got %d", tc.expectedRows, output.RowCount)
			}

			if output.ColCount != tc.expectedCols {
				t.Errorf("Expected %d columns, got %d", tc.expectedCols, output.ColCount)
			}

			if len(output.Headers) != tc.expectedCols {
				t.Errorf("Expected %d headers, got %d", tc.expectedCols, len(output.Headers))
			}
		})
	}
}

func TestReadExcelTool_Execute_MaxRows(t *testing.T) {
	// Criar arquivo ATIVOS para teste
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "ATIVOS.xlsx")

	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	tool := NewReadExcelTool()

	// Testar com max_rows = 3
	input := `{"file_path": "` + strings.ReplaceAll(testFile, `\`, `\\`) + `", "max_rows": 3}`

	result, err := tool.Execute(input)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Parse result
	var output ReadExcelOutput
	if err := json.Unmarshal([]byte(result), &output); err != nil {
		t.Fatalf("Failed to parse output: %v", err)
	}

	// Verificar se limitou corretamente
	if len(output.Data) > 3 {
		t.Errorf("Expected max 3 data rows, got %d", len(output.Data))
	}
}
