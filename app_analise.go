package main

import (
	"fmt"

	"BrxAgente-desafio4/internal/calculo"
	"BrxAgente-desafio4/internal/excel"
)

// RealizarAnaliseOrquestrada realiza todo o processo de análise de forma orquestrada
// incluindo leitura, cálculo e geração da planilha de resultado
func (a *App) RealizarAnaliseOrquestrada(diretorioPlanilhas string) (string, error) {
	// Validar o diretório de planilhas
	valido, err := a.SetDiretorioPlanilhas(diretorioPlanilhas)
	if err != nil {
		return "", fmt.Errorf("diretório de planilhas inválido: %w", err)
	}
	if !valido {
		return "", fmt.Errorf("diretório de planilhas inválido")
	}

	// Extrair e validar o mês da planilha DESLIGADOS
	mes, ano, err := calculo.ExtrairEValidarMesDesligados(diretorioPlanilhas)
	if err != nil {
		return "", fmt.Errorf("erro ao extrair mês da planilha DESLIGADOS: %w", err)
	}

	// Consolidar bases de dados
	colaboradores, err := calculo.ConsolidarBases(diretorioPlanilhas)
	if err != nil {
		return "", fmt.Errorf("erro ao consolidar bases de dados: %w", err)
	}

	// Verificar se há colaboradores para processar
	if len(colaboradores) == 0 {
		return "", fmt.Errorf("nenhum colaborador encontrado para processar")
	}

	// Salvar planilha de resultado na pasta de Downloads usando o mês extraído da planilha DESLIGADOS
	nomeArquivo := fmt.Sprintf("VR_Mensal %02d.%d.xlsx", mes, ano)
	if err := excel.SalvarPlanilhaEmDownloads(colaboradores, nomeArquivo); err != nil {
		return "", fmt.Errorf("erro ao salvar planilha de resultado: %w", err)
	}

	// Armazenar os dados consolidados na instância do App
	a.mu.Lock()
	a.colaboradores = colaboradores
	a.mu.Unlock()

	// Enviar os dados consolidados para o chat
	if err := a.SetChatContext(); err != nil {
		// Logar o erro, mas não falhar a operação principal
		fmt.Printf("Aviso: Falha ao definir o contexto do chat: %v\n", err)
	}

	// Retornar mensagem de sucesso com o número de colaboradores processados
	return fmt.Sprintf("Análise concluída com sucesso! %d colaboradores processados. Planilha salva em Downloads como %s",
		len(colaboradores), nomeArquivo), nil
}
