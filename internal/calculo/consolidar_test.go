package calculo

import (
	"path/filepath"
	"testing"
)

func TestConsolidarBases(t *testing.T) {
	// Teste com um diretório válido contendo planilhas
	t.Run("DiretorioValido", func(t *testing.T) {
		// Usando o diretório de arquivos do projeto
		diretorio := filepath.Join("..", "..", "files")

		colaboradores, err := ConsolidarBases(diretorio)

		// Verificando se não houve erro
		if err != nil {
			t.Errorf("Erro inesperado ao consolidar bases: %v", err)
			return
		}

		// Verificando se o mapa de colaboradores foi retornado
		if colaboradores == nil {
			t.Error("Mapa de colaboradores retornado é nil")
			return
		}

		// Verificando se foram encontrados colaboradores
		if len(colaboradores) == 0 {
			t.Error("Nenhum colaborador foi encontrado na consolidação")
			return
		}

		// Verificando se alguns colaboradores esperados estão presentes
		// (usando matrículas de exemplo das planilhas)
		matriculasExemplo := []string{"34941", "24401", "32104"}
		for _, matricula := range matriculasExemplo {
			if _, existe := colaboradores[matricula]; !existe {
				t.Logf("Aviso: Colaborador com matrícula %s não encontrado", matricula)
			}
		}
	})

	// Teste com um diretório inexistente
	t.Run("DiretorioInexistente", func(t *testing.T) {
		diretorio := filepath.Join("..", "..", "diretorio_inexistente")

		colaboradores, err := ConsolidarBases(diretorio)

		// Verificando se houve erro (como esperado)
		if err == nil {
			t.Error("Esperava-se um erro ao consolidar bases de diretório inexistente, mas não houve erro")
			return
		}

		// Verificando se o mapa de colaboradores retornado é nil (como esperado)
		if colaboradores != nil {
			t.Error("Esperava-se nil ao consolidar bases de diretório inexistente, mas foi retornado um mapa")
		}
	})
}

func TestProcessarAtivos(t *testing.T) {
	// TODO: Implementar testes específicos para a função processarAtivos
}

func TestProcessarFerias(t *testing.T) {
	// TODO: Implementar testes específicos para a função processarFerias
}

func TestProcessarDesligados(t *testing.T) {
	// TODO: Implementar testes específicos para a função processarDesligados
}
