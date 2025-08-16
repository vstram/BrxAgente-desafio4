// Package modelo provides data structures for the application
package modelo

import (
	"fmt"
)

// ErroProcessamento representa um erro customizado para o processamento de planilhas
type ErroProcessamento struct {
	// Mensagem é a descrição do erro
	Mensagem string
	
	// Arquivo é o nome do arquivo em que ocorreu o erro
	Arquivo string
	
	// Linha é o número da linha onde ocorreu o erro (se aplicável)
	Linha int
}

// Error implementa a interface error do Go para o tipo ErroProcessamento
func (e *ErroProcessamento) Error() string {
	if e.Arquivo != "" && e.Linha > 0 {
		return fmt.Sprintf("Erro no arquivo %s linha %d: %s", e.Arquivo, e.Linha, e.Mensagem)
	}
	
	if e.Arquivo != "" {
		return fmt.Sprintf("Erro no arquivo %s: %s", e.Arquivo, e.Mensagem)
	}
	
	return fmt.Sprintf("Erro: %s", e.Mensagem)
}

// NovoErroProcessamento cria uma nova instância de ErroProcessamento com mensagem
func NovoErroProcessamento(mensagem string) *ErroProcessamento {
	return &ErroProcessamento{
		Mensagem: mensagem,
	}
}

// NovoErroProcessamentoComArquivo cria uma nova instância de ErroProcessamento com mensagem e arquivo
func NovoErroProcessamentoComArquivo(mensagem, arquivo string) *ErroProcessamento {
	return &ErroProcessamento{
		Mensagem: mensagem,
		Arquivo:  arquivo,
	}
}

// NovoErroProcessamentoCompleto cria uma nova instância de ErroProcessamento com todos os campos
func NovoErroProcessamentoCompleto(mensagem, arquivo string, linha int) *ErroProcessamento {
	return &ErroProcessamento{
		Mensagem: mensagem,
		Arquivo:  arquivo,
		Linha:    linha,
	}
}