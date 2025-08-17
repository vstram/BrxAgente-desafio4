// Package calculo provides functionality for calculating VR values
package calculo

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	
	"github.com/xuri/excelize/v2"
	
	"BrxAgente-desafio4/internal/excel"
	"BrxAgente-desafio4/internal/modelo"
)

// ConsolidarBases consolida os dados das 5 planilhas separadas em uma única estrutura de dados
// As planilhas esperadas são:
// 1. ATIVOS.xlsx - Lista de colaboradores ativos
// 2. FÉRIAS.xlsx - Períodos de férias dos colaboradores
// 3. DESLIGADOS.xlsx - Colaboradores desligados
// 4. Base sindicato x valor.xlsx - Valores de VR por sindicato
// 5. Base dias uteis.xlsx - Dias úteis por sindicato
//
// Parâmetros:
//   - diretorioPlanilhas: Caminho para o diretório onde estão as planilhas
//
// Retornos:
//   - map[string]*modelo.Colaborador: Mapa de colaboradores com matrícula como chave
//   - error: Erro ocorrido durante a consolidação, se houver
func ConsolidarBases(diretorioPlanilhas string) (map[string]*modelo.Colaborador, error) {
	// Mapa para armazenar os colaboradores consolidados
	colaboradores := make(map[string]*modelo.Colaborador)
	
	// 1. Ler a planilha de ATIVOS
	caminhoAtivos := filepath.Join(diretorioPlanilhas, "ATIVOS.xlsx")
	fAtivos, err := excel.LerPlanilha(caminhoAtivos)
	if err != nil {
		return nil, err
	}
	defer fAtivos.Close()
	
	// Processar dados da planilha de ATIVOS
	err = processarAtivos(fAtivos, colaboradores)
	if err != nil {
		return nil, err
	}
	
	// 2. Ler a planilha de FÉRIAS
	caminhoFerias := filepath.Join(diretorioPlanilhas, "FÉRIAS.xlsx")
	fFerias, err := excel.LerPlanilha(caminhoFerias)
	if err != nil {
		return nil, err
	}
	defer fFerias.Close()
	
	// Processar dados da planilha de FÉRIAS
	err = processarFerias(fFerias, colaboradores)
	if err != nil {
		return nil, err
	}
	
	// 3. Ler a planilha de DESLIGADOS
	caminhoDesligados := filepath.Join(diretorioPlanilhas, "DESLIGADOS.xlsx")
	fDesligados, err := excel.LerPlanilha(caminhoDesligados)
	if err != nil {
		return nil, err
	}
	defer fDesligados.Close()
	
	// Processar dados da planilha de DESLIGADOS
	err = processarDesligados(fDesligados, colaboradores)
	if err != nil {
		return nil, err
	}
	
	// 4. Ler a planilha Base sindicato x valor
	caminhoSindicato := filepath.Join(diretorioPlanilhas, "Base sindicato x valor.xlsx")
	fSindicato, err := excel.LerPlanilha(caminhoSindicato)
	if err != nil {
		return nil, err
	}
	defer fSindicato.Close()
	
	// Processar dados da planilha de valores por sindicato
	// TODO: Implementar processamento da planilha de valores por sindicato
	
	// 5. Ler a planilha Base dias uteis
	caminhoDiasUteis := filepath.Join(diretorioPlanilhas, "Base dias uteis.xlsx")
	fDiasUteis, err := excel.LerPlanilha(caminhoDiasUteis)
	if err != nil {
		return nil, err
	}
	defer fDiasUteis.Close()
	
	// Processar dados da planilha de dias úteis
	// TODO: Implementar processamento da planilha de dias úteis
	
	return colaboradores, nil
}

// processarAtivos processa os dados da planilha de colaboradores ativos
func processarAtivos(f *excelize.File, colaboradores map[string]*modelo.Colaborador) error {
	// Obter todas as linhas da primeira sheet
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return modelo.NovoErroProcessamentoCompleto(
			"Erro ao ler linhas da planilha de ativos",
			"ATIVOS.xlsx",
			0,
		)
	}
	
	fmt.Printf("Total de linhas na planilha ATIVOS: %d\n", len(rows))
	
	// Processar cada linha (ignorando o cabeçalho)
	for i, row := range rows {
		// Ignorar a primeira linha (cabeçalho)
		if i == 0 {
			continue
		}
		
		// Verificar se a linha tem dados suficientes (pelo menos 5 colunas)
		if len(row) < 5 {
			continue
		}
		
		// Extrair os dados da linha
		matricula := strings.TrimSpace(row[0])
		empresa := strings.TrimSpace(row[1])
		cargo := strings.TrimSpace(row[2])
		situacao := strings.TrimSpace(row[3])
		sindicato := strings.TrimSpace(row[4])
		
		// Criar um novo colaborador
		colaborador := &modelo.Colaborador{
			Matricula: matricula,
			Empresa:   empresa,
			Cargo:     cargo,
			Situacao:  situacao,
			Sindicato: sindicato,
		}
		
		// Adicionar ao mapa de colaboradores
		colaboradores[matricula] = colaborador
	}
	
	return nil
}

// processarFerias processa os dados da planilha de férias
func processarFerias(f *excelize.File, colaboradores map[string]*modelo.Colaborador) error {
	// Obter todas as linhas da primeira sheet
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return modelo.NovoErroProcessamentoCompleto(
			"Erro ao ler linhas da planilha de férias",
			"FÉRIAS.xlsx",
			0,
		)
	}
	
	fmt.Printf("Total de linhas na planilha FÉRIAS: %d\n", len(rows))
	
	// Processar cada linha (ignorando o cabeçalho)
	for i, row := range rows {
		// Ignorar a primeira linha (cabeçalho)
		if i == 0 {
			continue
		}
		
		// Verificar se a linha tem dados suficientes (pelo menos 3 colunas)
		if len(row) < 3 {
			continue
		}
		
		// Extrair os dados da linha
		matricula := strings.TrimSpace(row[0])
		situacao := strings.TrimSpace(row[1])
		diasStr := strings.TrimSpace(row[2])
		
		// Converter dias para inteiro
		dias, err := strconv.Atoi(diasStr)
		if err != nil {
			// Se não conseguir converter, continua para o próximo
			continue
		}
		
		// Verificar se o colaborador existe
		if colaborador, existe := colaboradores[matricula]; existe {
			// Atualizar situação do colaborador
			colaborador.Situacao = situacao
			
			// TODO: Adicionar período de férias ao colaborador
			// Por enquanto, estamos apenas atualizando a situação e usando a variável dias
			_ = dias // Ignorando o valor por enquanto, mas usando a variável para evitar erro
		}
		// Se o colaborador não existe, ignoramos (pode ser de outro mês)
	}
	
	return nil
}

// processarDesligados processa os dados da planilha de desligados
func processarDesligados(f *excelize.File, colaboradores map[string]*modelo.Colaborador) error {
	// Obter todas as linhas da primeira sheet
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return modelo.NovoErroProcessamentoCompleto(
			"Erro ao ler linhas da planilha de desligados",
			"DESLIGADOS.xlsx",
			0,
		)
	}
	
	fmt.Printf("Total de linhas na planilha DESLIGADOS: %d\n", len(rows))
	
	// Processar cada linha (ignorando o cabeçalho)
	for i, row := range rows {
		// Ignorar a primeira linha (cabeçalho)
		if i == 0 {
			continue
		}
		
		// Verificar se a linha tem dados suficientes (pelo menos 1 coluna)
		if len(row) < 1 {
			continue
		}
		
		// Extrair os dados da linha
		matricula := strings.TrimSpace(row[0])
		
		// Verificar se o colaborador existe
		if colaborador, existe := colaboradores[matricula]; existe {
			// Atualizar situação do colaborador como desligado
			colaborador.Situacao = "Desligado"
			
			// TODO: Adicionar data de desligamento ao colaborador
			// Por enquanto, estamos apenas atualizando a situação
		}
		// Se o colaborador não existe, ignoramos (pode ser de outro mês)
	}
	
	return nil
}