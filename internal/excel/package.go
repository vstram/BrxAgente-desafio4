// Package excel provides functionality for manipulating Excel spreadsheets
package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

// Service representa um serviço para manipulação de planilhas Excel
type Service struct {
	// Cache para arquivos abertos
	openFiles map[string]*excelize.File
}

// NewService cria uma nova instância do serviço Excel
func NewService() *Service {
	return &Service{
		openFiles: make(map[string]*excelize.File),
	}
}

// ReadFile lê um arquivo Excel e retorna suas informações
func (s *Service) ReadFile(filePath string, sheetName string, maxRows int) (map[string]interface{}, error) {
	// Usar a função existente LerPlanilha
	f, err := LerPlanilha(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Se não especificou a planilha, usar a primeira
	if sheetName == "" {
		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			return nil, fmt.Errorf("nenhuma planilha encontrada no arquivo")
		}
		sheetName = sheets[0]
	}

	// Obter dados da planilha
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler planilha '%s': %w", sheetName, err)
	}

	if len(rows) == 0 {
		return map[string]interface{}{
			"file_path": filePath,
			"sheet":     sheetName,
			"row_count": 0,
			"col_count": 0,
			"headers":   []string{},
			"data":      []map[string]interface{}{},
		}, nil
	}

	// Primeira linha como headers
	headers := rows[0]
	colCount := len(headers)

	// Preparar dados
	data := []map[string]interface{}{}
	rowCount := len(rows) - 1 // Excluir header

	// Aplicar limitação de linhas se especificada
	endRow := len(rows)
	if maxRows > 0 && maxRows < rowCount {
		endRow = maxRows + 1 // +1 para incluir header
		rowCount = maxRows
	}

	// Processar dados (excluir header)
	for i := 1; i < endRow; i++ {
		row := rows[i]
		rowData := make(map[string]interface{})

		for j, header := range headers {
			if j < len(row) {
				rowData[header] = row[j]
			} else {
				rowData[header] = ""
			}
		}
		data = append(data, rowData)
	}

	return map[string]interface{}{
		"file_path": filePath,
		"sheet":     sheetName,
		"row_count": rowCount,
		"col_count": colCount,
		"headers":   headers,
		"data":      data,
	}, nil
}
