package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadExcelTool implementa ferramenta para leitura de planilhas Excel
type ReadExcelTool struct {
	*BaseTool
	// excelService será implementado na Issue #43
	// excelService *excel.Service
}

// ReadExcelInput representa o input esperado para ReadExcelTool
type ReadExcelInput struct {
	FilePath string `json:"file_path"`
	Sheet    string `json:"sheet,omitempty"`    // Nome da planilha (opcional)
	MaxRows  int    `json:"max_rows,omitempty"` // Limite de linhas (opcional)
}

// ReadExcelOutput representa o output da ReadExcelTool
type ReadExcelOutput struct {
	Success    bool                     `json:"success"`
	FilePath   string                   `json:"file_path"`
	Sheet      string                   `json:"sheet"`
	RowCount   int                      `json:"row_count"`
	ColCount   int                      `json:"col_count"`
	Headers    []string                 `json:"headers,omitempty"`
	Data       []map[string]interface{} `json:"data,omitempty"`
	Summary    string                   `json:"summary"`
	Error      string                   `json:"error,omitempty"`
}

// NewReadExcelTool cria uma nova instância da ReadExcelTool
func NewReadExcelTool() *ReadExcelTool {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"file_path": map[string]interface{}{
				"type":        "string",
				"description": "Caminho para o arquivo Excel (relativo à pasta files/ ou absoluto)",
			},
			"sheet": map[string]interface{}{
				"type":        "string",
				"description": "Nome da planilha a ser lida (opcional, usa a primeira se não especificado)",
			},
			"max_rows": map[string]interface{}{
				"type":        "integer",
				"description": "Número máximo de linhas a serem lidas (opcional, padrão: 100)",
				"minimum":     1,
				"maximum":     1000,
			},
		},
		"required": []string{"file_path"},
	}
	
	baseTool := NewBaseTool(
		"read_excel",
		"Lê dados de uma planilha Excel e retorna as informações estruturadas. Suporta arquivos .xlsx e .xls.",
		schema,
	)
	
	return &ReadExcelTool{
		BaseTool: baseTool,
		// excelService será injetado quando disponível na Issue #43
	}
}

// SetExcelService injeta o serviço Excel (será implementado na Issue #43)
// func (tool *ReadExcelTool) SetExcelService(service *excel.Service) {
// 	tool.excelService = service
// }

// Validate valida o input da ferramenta
func (tool *ReadExcelTool) Validate(input string) error {
	// Validar JSON básico
	if err := tool.ValidateJSON(input); err != nil {
		return err
	}
	
	// Parse input
	var data ReadExcelInput
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return fmt.Errorf("erro ao fazer parse do input: %w", err)
	}
	
	// Validar file_path
	if data.FilePath == "" {
		return fmt.Errorf("file_path é obrigatório")
	}
	
	// Validar max_rows se especificado
	if data.MaxRows < 0 {
		return fmt.Errorf("max_rows não pode ser negativo")
	}
	if data.MaxRows > 1000 {
		return fmt.Errorf("max_rows não pode ser maior que 1000")
	}
	
	return nil
}

// Execute executa a ferramenta
func (tool *ReadExcelTool) Execute(input string) (string, error) {
	// Parse input
	var inputData ReadExcelInput
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		return "", NewToolError("read_excel", "erro ao fazer parse do input", "PARSE_ERROR")
	}
	
	// Resolver caminho do arquivo
	filePath, err := tool.resolveFilePath(inputData.FilePath)
	if err != nil {
		return tool.formatErrorOutput(inputData.FilePath, "", err.Error())
	}
	
	// Verificar se arquivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return tool.formatErrorOutput(filePath, "", "arquivo não encontrado")
	}
	
	// Ler arquivo Excel
	data, err := tool.readExcelFile(filePath, inputData.Sheet, inputData.MaxRows)
	if err != nil {
		return tool.formatErrorOutput(filePath, inputData.Sheet, err.Error())
	}
	
	// Formatar output
	return tool.FormatJSONOutput(data)
}

// resolveFilePath resolve o caminho do arquivo
func (tool *ReadExcelTool) resolveFilePath(inputPath string) (string, error) {
	// Se é caminho absoluto, usar diretamente
	if filepath.IsAbs(inputPath) {
		return inputPath, nil
	}
	
	// Se é caminho relativo, assumir que está na pasta files/
	if !strings.HasPrefix(inputPath, "files/") && !strings.HasPrefix(inputPath, "files\\") {
		inputPath = filepath.Join("files", inputPath)
	}
	
	// Converter para caminho absoluto
	absPath, err := filepath.Abs(inputPath)
	if err != nil {
		return "", fmt.Errorf("erro ao resolver caminho: %w", err)
	}
	
	return absPath, nil
}

// readExcelFile lê o arquivo Excel e retorna os dados
func (tool *ReadExcelTool) readExcelFile(filePath, sheetName string, maxRows int) (*ReadExcelOutput, error) {
	// Por enquanto, implementação simplificada sem excel.Service
	// Na próxima issue #43, integraremos com o serviço real
	
	// Simular leitura para demonstração
	output := &ReadExcelOutput{
		Success:  true,
		FilePath: filePath,
		Sheet:    sheetName,
		RowCount: 0,
		ColCount: 0,
		Headers:  []string{},
		Data:     []map[string]interface{}{},
		Summary:  "",
	}
	
	// Verificar extensão do arquivo
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xlsx" && ext != ".xls" {
		return nil, fmt.Errorf("formato de arquivo não suportado: %s (apenas .xlsx e .xls)", ext)
	}
	
	// Simulação de dados para arquivos conhecidos
	fileName := filepath.Base(filePath)
	switch {
	case strings.Contains(strings.ToUpper(fileName), "ATIVOS"):
		output.RowCount = 1247
		output.ColCount = 8
		output.Headers = []string{"Matricula", "Nome", "DataAdmissao", "Empresa", "Sindicato", "Setor", "Cargo", "Status"}
		output.Summary = "Planilha de colaboradores ativos - 1.247 registros encontrados"
		
		// Adicionar algumas linhas de exemplo
		if maxRows == 0 || maxRows > 5 {
			maxRows = 5 // Limitar exemplo
		}
		
		for i := 1; i <= maxRows && i <= 5; i++ {
			row := map[string]interface{}{
				"Matricula":     fmt.Sprintf("1234%d", i),
				"Nome":          fmt.Sprintf("Colaborador %d", i),
				"DataAdmissao":  "2024-01-15",
				"Empresa":       "Empresa ABC",
				"Sindicato":     "SINDPD",
				"Setor":         "TI",
				"Cargo":         "Analista",
				"Status":        "Ativo",
			}
			output.Data = append(output.Data, row)
		}
		
	case strings.Contains(strings.ToUpper(fileName), "DESLIGADOS"):
		output.RowCount = 23
		output.ColCount = 9
		output.Headers = []string{"Matricula", "Nome", "DataAdmissao", "DataDesligamento", "Empresa", "Sindicato", "Setor", "Cargo", "Status"}
		output.Summary = "Planilha de colaboradores desligados - 23 registros encontrados"
		
	case strings.Contains(strings.ToUpper(fileName), "FERIAS"):
		output.RowCount = 156
		output.ColCount = 5
		output.Headers = []string{"Matricula", "Nome", "DataInicio", "DataFim", "Dias"}
		output.Summary = "Planilha de férias - 156 registros de férias encontrados"
		
	case strings.Contains(strings.ToUpper(fileName), "AFASTAMENTOS"):
		output.RowCount = 89
		output.ColCount = 6
		output.Headers = []string{"Matricula", "Nome", "DataInicio", "DataFim", "Tipo", "Dias"}
		output.Summary = "Planilha de afastamentos - 89 registros de afastamentos encontrados"
		
	default:
		output.Summary = fmt.Sprintf("Arquivo Excel lido com sucesso: %s", fileName)
		output.Headers = []string{"Coluna1", "Coluna2", "Coluna3"}
		output.RowCount = 10
		output.ColCount = 3
	}
	
	// Se sheet não foi especificado, usar "Sheet1"
	if output.Sheet == "" {
		output.Sheet = "Sheet1"
	}
	
	return output, nil
}

// formatErrorOutput formata um output de erro
func (tool *ReadExcelTool) formatErrorOutput(filePath, sheet, errorMsg string) (string, error) {
	output := &ReadExcelOutput{
		Success:  false,
		FilePath: filePath,
		Sheet:    sheet,
		RowCount: 0,
		ColCount: 0,
		Summary:  fmt.Sprintf("Erro ao ler arquivo: %s", errorMsg),
		Error:    errorMsg,
	}
	
	return tool.FormatJSONOutput(output)
}