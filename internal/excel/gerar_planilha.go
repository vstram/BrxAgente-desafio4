// Package excel provides functionality for manipulating Excel spreadsheets
package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"BrxAgente-desafio4/internal/modelo"
)

// GerarPlanilhaResultado cria uma planilha Excel com os resultados do cálculo de VR
// seguindo o formato especificado no modelo VR MENSAL 05.2025.xlsx
// Respeitando a Nota de Confidencialidade, apenas a MATRICULA é usada para identificar o colaborador
func GerarPlanilhaResultado(colaboradores map[string]*modelo.Colaborador, caminhoArquivo string) error {
	// Criar um novo arquivo Excel
	f := excelize.NewFile()

	// Definir o nome da planilha
	sheetName := "Planilha1"
	f.SetSheetName("Sheet1", sheetName)

	// Criar os cabeçalhos da planilha conforme modelo VR MENSAL 05.2025.xlsx
	cabecalhos := []string{
		"Matricula", "Admissão", "Sindicato do Colaborador", "Competência",
		"Dias", "VALOR DIÁRIO VR", "TOTAL", "Custo empresa", "Desconto profissional", "OBS GERAL",
	}

	// Preencher os cabeçalhos
	for i, cabecalho := range cabecalhos {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, cabecalho)
	}

	// Preencher os dados dos colaboradores
	linha := 2
	for _, colaborador := range colaboradores {
		// Preencher os dados do colaborador conforme modelo da planilha
		dados := []interface{}{
			colaborador.Matricula,                                // Matricula
			colaborador.DataAdmissao.Format("02-01-06"),         // Admissão (formato dd-mm-yy)
			colaborador.Sindicato,                               // Sindicato do Colaborador
			"05/2025",                                           // Competência
			fmt.Sprintf("%.2f", float64(colaborador.DiasUteisEfetivos)), // Dias
			37.50,                                               // VALOR DIÁRIO VR (valor fixo conforme planilha)
			colaborador.ValorTotalVR,                           // TOTAL
			colaborador.ValorEmpresa,                           // Custo empresa
			colaborador.ValorColaborador,                       // Desconto profissional
			"",                                                 // OBS GERAL (vazio por padrão)
		}

		// Preencher os dados na linha correspondente
		for i, dado := range dados {
			cell, _ := excelize.CoordinatesToCellName(i+1, linha)
			f.SetCellValue(sheetName, cell, dado)
		}

		linha++
	}

	// Salvar o arquivo
	if err := f.SaveAs(caminhoArquivo); err != nil {
		return fmt.Errorf("erro ao salvar arquivo: %w", err)
	}

	return nil
}

// GerarPlanilhaResultadoComTemplate cria uma planilha Excel com os resultados do cálculo de VR
// usando um template existente como base
// Respeitando a Nota de Confidencialidade, apenas a MATRICULA é usada para identificar o colaborador
func GerarPlanilhaResultadoComTemplate(colaboradores map[string]*modelo.Colaborador, caminhoTemplate, caminhoArquivo string) error {
	// Abrir o template
	f, err := excelize.OpenFile(caminhoTemplate)
	if err != nil {
		return fmt.Errorf("erro ao abrir template: %w", err)
	}

	// Definir o nome da planilha (assumindo que é a primeira)
	sheetName := f.GetSheetName(0)

	// Limpar os dados existentes (mantendo apenas o cabeçalho)
	// Assumindo que os dados começam na linha 2
	linha := 2
	for {
		cell, _ := excelize.CoordinatesToCellName(1, linha)
		val, err := f.GetCellValue(sheetName, cell)
		if err != nil || val == "" {
			break
		}

		// Limpar a linha
		for col := 1; col <= 10; col++ { // 10 colunas conforme novo modelo
			cell, _ := excelize.CoordinatesToCellName(col, linha)
			f.SetCellValue(sheetName, cell, "")
		}

		linha++
	}

	// Preencher os dados dos colaboradores começando da linha 2
	linha = 2
	for _, colaborador := range colaboradores {
		// Preencher os dados do colaborador conforme modelo da planilha
		dados := []interface{}{
			colaborador.Matricula,                                // Matricula
			colaborador.DataAdmissao.Format("02-01-06"),         // Admissão (formato dd-mm-yy)
			colaborador.Sindicato,                               // Sindicato do Colaborador
			"05/2025",                                           // Competência
			fmt.Sprintf("%.2f", float64(colaborador.DiasUteisEfetivos)), // Dias
			37.50,                                               // VALOR DIÁRIO VR (valor fixo conforme planilha)
			colaborador.ValorTotalVR,                           // TOTAL
			colaborador.ValorEmpresa,                           // Custo empresa
			colaborador.ValorColaborador,                       // Desconto profissional
			"",                                                 // OBS GERAL (vazio por padrão)
		}

		// Preencher os dados na linha correspondente
		for i, dado := range dados {
			cell, _ := excelize.CoordinatesToCellName(i+1, linha)
			f.SetCellValue(sheetName, cell, dado)
		}

		linha++
	}

	// Salvar o arquivo
	if err := f.SaveAs(caminhoArquivo); err != nil {
		return fmt.Errorf("erro ao salvar arquivo: %w", err)
	}

	return nil
}
