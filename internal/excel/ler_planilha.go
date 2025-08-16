// Package excel provides functionality for manipulating Excel spreadsheets
package excel

import (
	"github.com/xuri/excelize/v2"
	
	"BrxAgente-desafio4/internal/modelo"
)

// LerPlanilha abre e lê um arquivo Excel, retornando um ponteiro para o arquivo e um erro.
// Esta função é a base para o processamento dos dados das planilhas.
//
// Parâmetros:
//   - caminho: O caminho completo para o arquivo Excel a ser lido
//
// Retornos:
//   - *excelize.File: Ponteiro para o objeto File que representa a planilha Excel
//   - error: Erro ocorrido durante a leitura do arquivo, se houver
//
// Exemplo de uso:
//   f, err := LerPlanilha("./files/ATIVOS.xlsx")
//   if err != nil {
//       log.Fatal(err)
//   }
//   defer f.Close()
func LerPlanilha(caminho string) (*excelize.File, error) {
	// Tentando abrir o arquivo Excel
	f, err := excelize.OpenFile(caminho)
	if err != nil {
		// Criando um erro customizado usando o tipo ErroProcessamento
		return nil, modelo.NovoErroProcessamentoCompleto(
			"Erro ao abrir o arquivo Excel",
			caminho,
			0, // Linha não se aplica neste caso
		)
	}
	
	return f, nil
}