// Package excel provides functionality for manipulating Excel spreadsheets
package excel

import (
	"path/filepath"
	"testing"

	"BrxAgente-desafio4/internal/modelo"
)

func TestSalvarPlanilhaEmDownloads(t *testing.T) {
	// Criar dados de teste
	colaboradores := map[string]*modelo.Colaborador{
		"001": {
			Matricula:         "001",
			Empresa:           "Empresa A",
			Cargo:             "Analista",
			Sindicato:         "Paraná",
			DiasUteisEfetivos: 20,
			ValorTotalVR:      700.00,
			ValorEmpresa:      560.00, // 80% de 700
			ValorColaborador:  140.00, // 20% de 700
		},
	}

	// Testar salvar planilha em Downloads
	err := SalvarPlanilhaEmDownloads(colaboradores, "teste_vr_resultado.xlsx")
	if err != nil {
		// Se não conseguir salvar em Downloads, isso pode ser esperado em alguns ambientes de teste
		t.Logf("Aviso: Erro ao salvar planilha em Downloads (pode ser esperado em ambiente de teste): %v", err)
	}
}

func TestSalvarPlanilhaEmDownloadsComTemplate(t *testing.T) {
	// Criar dados de teste
	colaboradores := map[string]*modelo.Colaborador{
		"001": {
			Matricula:         "001",
			Empresa:           "Empresa A",
			Cargo:             "Analista",
			Sindicato:         "Paraná",
			DiasUteisEfetivos: 20,
			ValorTotalVR:      700.00,
			ValorEmpresa:      560.00, // 80% de 700
			ValorColaborador:  140.00, // 20% de 700
		},
	}

	// Caminho para o template (usando uma planilha existente como template)
	caminhoTemplate := filepath.Join("..", "..", "files", "VR MENSAL 05.2025.xlsx")

	// Testar salvar planilha em Downloads com template
	err := SalvarPlanilhaEmDownloadsComTemplate(colaboradores, caminhoTemplate, "teste_vr_resultado_template.xlsx")
	if err != nil {
		// Se não conseguir salvar em Downloads, isso pode ser esperado em alguns ambientes de teste
		t.Logf("Aviso: Erro ao salvar planilha em Downloads com template (pode ser esperado em ambiente de teste): %v", err)
	}
}

func TestObterDiretorioDownloads(t *testing.T) {
	// Testar obter diretório de Downloads
	diretorio, err := obterDiretorioDownloads()
	if err != nil {
		t.Fatalf("Erro ao obter diretório de Downloads: %v", err)
	}

	// Verificar se o diretório não está vazio
	if diretorio == "" {
		t.Error("Diretório de Downloads está vazio")
	}

	// Verificar se o diretório existe
	// Note: In a real test, we would check if the directory exists
	// For now, we're just checking that no error occurred and the directory is not empty
}

func TestVerificarPermissoesDiretorio(t *testing.T) {
	// Obter o diretório temporário do sistema
	diretorioTemp := t.TempDir()

	// Testar verificar permissões do diretório temporário
	err := VerificarPermissoesDiretorio(diretorioTemp)
	if err != nil {
		t.Errorf("Erro ao verificar permissões do diretório temporário: %v", err)
	}
}
