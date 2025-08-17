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
	
	// Criar os cabeçalhos da planilha (sem NOME, conforme Nota de Confidencialidade)
	cabecalhos := []string{
		"EMPRESA", "MATRICULA", "CARGO", "SINDICATO", 
		"REF", "DIAS", "VR", "DESCONTO", "VALOR",
	}
	
	// Preencher os cabeçalhos
	for i, cabecalho := range cabecalhos {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, cabecalho)
	}
	
	// Preencher os dados dos colaboradores
	linha := 2
	for _, colaborador := range colaboradores {
		// Preencher os dados do colaborador (sem NOME, conforme Nota de Confidencialidade)
		dados := []interface{}{
			colaborador.Empresa,
			colaborador.Matricula,
			colaborador.Cargo,
			colaborador.Sindicato,
			"05/2025", // REF (mês de referência)
			colaborador.DiasUteisEfetivos,
			fmt.Sprintf("R$ %.2f", colaborador.ValorTotalVR),     // VR (valor total)
			fmt.Sprintf("R$ %.2f", colaborador.ValorColaborador), // DESCONTO (20% do valor total)
			fmt.Sprintf("R$ %.2f", colaborador.ValorEmpresa),     // VALOR (80% do valor total)
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
		for col := 1; col <= 9; col++ { // 9 colunas agora (sem NOME)
			cell, _ := excelize.CoordinatesToCellName(col, linha)
			f.SetCellValue(sheetName, cell, "")
		}
		
		linha++
	}
	
	// Preencher os dados dos colaboradores começando da linha 2
	linha = 2
	for _, colaborador := range colaboradores {
		// Preencher os dados do colaborador (sem NOME, conforme Nota de Confidencialidade)
		dados := []interface{}{
			colaborador.Empresa,
			colaborador.Matricula,
			colaborador.Cargo,
			colaborador.Sindicato,
			"05/2025", // REF (mês de referência)
			colaborador.DiasUteisEfetivos,
			colaborador.ValorTotalVR,     // VR (valor total)
			colaborador.ValorColaborador, // DESCONTO (20% do valor total)
			colaborador.ValorEmpresa,     // VALOR (80% do valor total)
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