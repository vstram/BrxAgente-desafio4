package excel

import (
	"path/filepath"
	"testing"
)

func TestLerPlanilha(t *testing.T) {
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