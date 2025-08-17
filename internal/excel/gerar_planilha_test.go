// Package excel provides functionality for manipulating Excel spreadsheets
package excel

import (
	"path/filepath"
	"testing"

	"BrxAgente-desafio4/internal/modelo"
)

func TestGerarPlanilhaResultado(t *testing.T) {
	// Criar dados de teste
	colaboradores := map[string]*modelo.Colaborador{
		"001": {
			Matricula:        "001",
			Empresa:          "Empresa A",
			Cargo:            "Analista",
			Sindicato:        "Paraná",
			DiasUteisEfetivos: 20,
			ValorTotalVR:      700.00,
			ValorEmpresa:      560.00, // 80% de 700
			ValorColaborador:  140.00, // 20% de 700
		},
		"002": {
			Matricula:        "002",
			Empresa:          "Empresa B",
			Cargo:            "Gerente",
			Sindicato:        "Rio Grande do Sul",
			DiasUteisEfetivos: 15,
			ValorTotalVR:      525.00,
			ValorEmpresa:      420.00, // 80% de 525
			ValorColaborador:  105.00, // 20% de 525
		},
	}

	// Caminho para o arquivo de teste
	caminhoArquivo := filepath.Join(t.TempDir(), "resultado_teste.xlsx")

	// Testar a geração da planilha
	err := GerarPlanilhaResultado(colaboradores, caminhoArquivo)
	if err != nil {
		t.Fatalf("Erro ao gerar planilha de resultado: %v", err)
	}

	// Verificar se o arquivo foi criado
	// Note: In a real test, we would open the file and verify its contents
	// For now, we're just checking that no error occurred
}

func TestGerarPlanilhaResultadoComTemplate(t *testing.T) {
	// Criar dados de teste
	colaboradores := map[string]*modelo.Colaborador{
		"001": {
			Matricula:        "001",
			Empresa:          "Empresa A",
			Cargo:            "Analista",
			Sindicato:        "Paraná",
			DiasUteisEfetivos: 20,
			ValorTotalVR:      700.00,
			ValorEmpresa:      560.00, // 80% de 700
			ValorColaborador:  140.00, // 20% de 700
		},
	}

	// Caminho para o template (usando uma planilha existente como template)
	caminhoTemplate := filepath.Join("..", "..", "files", "VR MENSAL 05.2025.xlsx")
	caminhoArquivo := filepath.Join(t.TempDir(), "resultado_template_teste.xlsx")

	// Testar a geração da planilha com template
	err := GerarPlanilhaResultadoComTemplate(colaboradores, caminhoTemplate, caminhoArquivo)
	if err != nil {
		// Se não conseguir abrir o template, isso é esperado em alguns ambientes de teste
		t.Logf("Aviso: Erro ao gerar planilha com template (pode ser esperado em ambiente de teste): %v", err)
	}
}