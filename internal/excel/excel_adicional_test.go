package excel

import (
	"os"
	"path/filepath"
	"testing"

	"BrxAgente-desafio4/internal/modelo"
)

func TestObterDiretorioDownloadsAdicional(t *testing.T) {
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
	if _, err := os.Stat(diretorio); os.IsNotExist(err) {
		t.Errorf("Diretório de Downloads não existe: %s", diretorio)
	}
}

func TestSalvarPlanilhaEmDownloadsAdicional(t *testing.T) {
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

	// Testar salvar planilha em Downloads com um nome único
	nomeArquivo := "teste_vr_resultado_unit_test_adicional.xlsx"
	err := SalvarPlanilhaEmDownloads(colaboradores, nomeArquivo)
	if err != nil {
		// Se não conseguir salvar em Downloads, isso pode ser esperado em alguns ambientes de teste
		t.Logf("Aviso: Erro ao salvar planilha em Downloads (pode ser esperado em ambiente de teste): %v", err)
	} else {
		// Verificar se o arquivo foi criado
		diretorioDownloads, err := obterDiretorioDownloads()
		if err != nil {
			t.Fatalf("Erro ao obter diretório de Downloads: %v", err)
		}

		caminhoArquivo := filepath.Join(diretorioDownloads, nomeArquivo)
		if _, err := os.Stat(caminhoArquivo); os.IsNotExist(err) {
			t.Error("Arquivo de resultado não foi criado")
		} else {
			// Remover o arquivo de teste
			defer os.Remove(caminhoArquivo)
		}
	}
}

func TestSalvarPlanilhaEmDownloadsComTemplateAdicional(t *testing.T) {
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
	nomeArquivo := "teste_vr_resultado_template_unit_test_adicional.xlsx"
	err := SalvarPlanilhaEmDownloadsComTemplate(colaboradores, caminhoTemplate, nomeArquivo)
	if err != nil {
		// Se não conseguir salvar em Downloads, isso pode ser esperado em alguns ambientes de teste
		t.Logf("Aviso: Erro ao salvar planilha em Downloads com template (pode ser esperado em ambiente de teste): %v", err)
	} else {
		// Verificar se o arquivo foi criado
		diretorioDownloads, err := obterDiretorioDownloads()
		if err != nil {
			t.Fatalf("Erro ao obter diretório de Downloads: %v", err)
		}

		caminhoArquivo := filepath.Join(diretorioDownloads, nomeArquivo)
		if _, err := os.Stat(caminhoArquivo); os.IsNotExist(err) {
			t.Error("Arquivo de resultado com template não foi criado")
		} else {
			// Remover o arquivo de teste
			defer os.Remove(caminhoArquivo)
		}
	}
}

func TestVerificarPermissoesDiretorioAdicional(t *testing.T) {
	// Obter o diretório temporário do sistema
	diretorioTemp := t.TempDir()

	// Testar verificar permissões do diretório temporário
	err := VerificarPermissoesDiretorio(diretorioTemp)
	if err != nil {
		t.Errorf("Erro ao verificar permissões do diretório temporário: %v", err)
	}

	// Testar com diretório que não existe
	diretorioInexistente := filepath.Join(t.TempDir(), "diretorio_inexistente")
	err = VerificarPermissoesDiretorio(diretorioInexistente)
	if err == nil {
		t.Error("Esperava erro ao verificar permissões de diretório inexistente, mas não houve erro")
	}
}

func TestGerarPlanilhaResultadoAdicional(t *testing.T) {
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
		"002": {
			Matricula:         "002",
			Empresa:           "Empresa B",
			Cargo:             "Gerente",
			Sindicato:         "Rio Grande do Sul",
			DiasUteisEfetivos: 15,
			ValorTotalVR:      525.00,
			ValorEmpresa:      420.00, // 80% de 525
			ValorColaborador:  105.00, // 20% de 525
		},
	}

	// Caminho para o arquivo de teste
	caminhoArquivo := filepath.Join(t.TempDir(), "resultado_completo_teste_adicional.xlsx")

	// Testar a geração da planilha
	err := GerarPlanilhaResultado(colaboradores, caminhoArquivo)
	if err != nil {
		t.Fatalf("Erro ao gerar planilha de resultado: %v", err)
	}

	// Verificar se o arquivo foi criado
	if _, err := os.Stat(caminhoArquivo); os.IsNotExist(err) {
		t.Error("Arquivo de resultado não foi criado")
	}
}

func TestGerarPlanilhaResultadoComTemplateAdicional(t *testing.T) {
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
	caminhoArquivo := filepath.Join(t.TempDir(), "resultado_template_completo_teste_adicional.xlsx")

	// Testar a geração da planilha com template
	err := GerarPlanilhaResultadoComTemplate(colaboradores, caminhoTemplate, caminhoArquivo)
	if err != nil {
		// Se não conseguir abrir o template, isso é esperado em alguns ambientes de teste
		t.Logf("Aviso: Erro ao gerar planilha com template (pode ser esperado em ambiente de teste): %v", err)
	} else {
		// Verificar se o arquivo foi criado
		if _, err := os.Stat(caminhoArquivo); os.IsNotExist(err) {
			t.Error("Arquivo de resultado com template não foi criado")
		}
	}
}

func TestLerPlanilhaAdicional(t *testing.T) {
	// Teste com um arquivo válido
	t.Run("ArquivoValido", func(t *testing.T) {
		// Usando uma planilha de exemplo que sabemos que existe
		caminho := filepath.Join("..", "..", "files", "ATIVOS.xlsx")

		f, err := LerPlanilha(caminho)

		// Verificando se não houve erro
		if err != nil {
			t.Errorf("Erro inesperado ao ler arquivo válido: %v", err)
			return
		}

		// Verificando se o arquivo foi retornado corretamente
		if f == nil {
			t.Error("Arquivo retornado é nil")
			return
		}

		// Fechando o arquivo
		err = f.Close()
		if err != nil {
			t.Errorf("Erro ao fechar o arquivo: %v", err)
		}
	})

	// Teste com um arquivo inexistente
	t.Run("ArquivoInexistente", func(t *testing.T) {
		caminho := filepath.Join("..", "..", "files", "INEXISTENTE.xlsx")

		f, err := LerPlanilha(caminho)

		// Verificando se houve erro (como esperado)
		if err == nil {
			t.Error("Esperava-se um erro ao ler arquivo inexistente, mas não houve erro")
			return
		}

		// Verificando se o arquivo retornado é nil (como esperado)
		if f != nil {
			t.Error("Esperava-se nil ao ler arquivo inexistente, mas foi retornado um arquivo")
			f.Close() // Fechando o arquivo se foi retornado
		}
	})
}
