// main.go
package main

import (
	"fmt"
	"log"
	
	"BrxAgente-desafio4/internal/calculo"
)

func main() {
	// Diretório onde estão as planilhas
	diretorioPlanilhas := "./files"
	
	// Consolidar as bases de dados
	fmt.Println("Iniciando consolidação das bases de dados...")
	colaboradores, err := calculo.ConsolidarBases(diretorioPlanilhas)
	if err != nil {
		log.Fatalf("Erro ao consolidar bases: %v", err)
	}
	
	// Exibir informações sobre os colaboradores consolidados
	fmt.Printf("Total de colaboradores consolidados: %d\n", len(colaboradores))
	
	// Exibir os primeiros 5 colaboradores como exemplo
	i := 0
	for _, colaborador := range colaboradores {
		if i >= 5 {
			break
		}
		
		fmt.Printf("Matrícula: %s, Empresa: %s, Cargo: %s, Situação: %s, Sindicato: %s\n",
			colaborador.Matricula, colaborador.Empresa, colaborador.Cargo, 
			colaborador.Situacao, colaborador.Sindicato)
		
		i++
	}
	
	fmt.Println("Consolidação concluída com sucesso!")
}