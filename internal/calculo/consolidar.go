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
	"BrxAgente-desafio4/internal/validacao"
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
	
	// Mapas para armazenar dados das planilhas de sindicato e dias úteis
	sindicatos := make(map[string]float64)
	diasUteis := make(map[string]int)
	
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
	err = processarValoresSindicato(fSindicato, sindicatos)
	if err != nil {
		return nil, err
	}
	
	// 5. Ler a planilha Base dias uteis
	caminhoDiasUteis := filepath.Join(diretorioPlanilhas, "Base dias uteis.xlsx")
	fDiasUteis, err := excel.LerPlanilha(caminhoDiasUteis)
	if err != nil {
		return nil, err
	}
	defer fDiasUteis.Close()
	
	// Processar dados da planilha de dias úteis
	err = processarDiasUteis(fDiasUteis, diasUteis)
	if err != nil {
		return nil, err
	}
	
	// 6. Validar os dados consolidados
	erros := validarDadosConsolidados(colaboradores, sindicatos, diasUteis)
	if len(erros) > 0 {
		// Retornar o primeiro erro encontrado
		return nil, erros[0]
	}
	
	return colaboradores, nil
}

// validarDadosConsolidados executa as validações nos dados consolidados
func validarDadosConsolidados(colaboradores map[string]*modelo.Colaborador, sindicatos map[string]float64, diasUteis map[string]int) []error {
	var erros []error
	
	// Criar mapas de mapeamento de sindicatos para estados
	mapaSindicatos := criarMapaSindicatos(sindicatos)
	
	// Validar cada colaborador
	for _, colaborador := range colaboradores {
		// Mapear o sindicato do colaborador para o estado correspondente
		estadoSindicato := mapearSindicatoParaEstado(colaborador.Sindicato, mapaSindicatos)
		if estadoSindicato != "" {
			// Criar uma cópia do colaborador com o sindicato mapeado
			colaboradorMapeado := *colaborador
			colaboradorMapeado.Sindicato = estadoSindicato
			
			// Validar com os sindicatos mapeados
			errs := validacao.ValidarColaborador(&colaboradorMapeado, sindicatos, diasUteis)
			erros = append(erros, errs...)
		} else {
			// Se não conseguir mapear, validar com os dados originais
			errs := validacao.ValidarColaborador(colaborador, sindicatos, diasUteis)
			erros = append(erros, errs...)
		}
	}
	
	return erros
}

// criarMapaSindicatos cria um mapa de sindicatos para estados
func criarMapaSindicatos(sindicatos map[string]float64) map[string]string {
	mapa := make(map[string]string)
	
	// Mapear os sindicatos para os estados
	for sindicato := range sindicatos {
		mapa[sindicato] = sindicato
	}
	
	return mapa
}

// mapearSindicatoParaEstado mapeia o nome do sindicato do colaborador para o estado correspondente
func mapearSindicatoParaEstado(sindicato string, mapaSindicatos map[string]string) string {
	// Mapear os sindicatos para os estados
	switch {
	case strings.Contains(sindicato, "PR") || strings.Contains(sindicato, "CURITIBA"):
		return "Paraná"
	case strings.Contains(sindicato, "RS"):
		return "Rio Grande do Sul"
	case strings.Contains(sindicato, "SP"):
		return "São Paulo"
	case strings.Contains(sindicato, "RJ"):
		return "Rio de Janeiro"
	default:
		// Se não conseguir identificar, retornar vazio
		return ""
	}
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

// processarValoresSindicato processa os dados da planilha de valores por sindicato
func processarValoresSindicato(f *excelize.File, sindicatos map[string]float64) error {
	// Obter todas as linhas da primeira sheet
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return modelo.NovoErroProcessamentoCompleto(
			"Erro ao ler linhas da planilha de valores por sindicato",
			"Base sindicato x valor.xlsx",
			0,
		)
	}
	
	fmt.Printf("Total de linhas na planilha Base sindicato x valor: %d\n", len(rows))
	
	// Processar cada linha (ignorando o cabeçalho)
	for i, row := range rows {
		// Ignorar a primeira linha (cabeçalho)
		if i == 0 {
			continue
		}
		
		// Verificar se a linha tem dados suficientes (pelo menos 2 colunas)
		if len(row) < 2 {
			continue
		}
		
		// Extrair os dados da linha
		sindicato := strings.TrimSpace(row[0])
		valorStr := strings.TrimSpace(row[1])
		
		// Extrair apenas o nome do estado do sindicato
		// Ex: "Paraná R$ 35.00" -> "Paraná"
		if strings.Contains(sindicato, "R$") {
			parts := strings.Split(sindicato, "R$")
			sindicato = strings.TrimSpace(parts[0])
		}
		
		// Converter valor para float64
		valorStr = strings.ReplaceAll(valorStr, "R$", "")
		valorStr = strings.ReplaceAll(valorStr, ",", ".")
		valorStr = strings.TrimSpace(valorStr)
		
		valor, err := strconv.ParseFloat(valorStr, 64)
		if err != nil {
			// Se não conseguir converter, continua para o próximo
			continue
		}
		
		// Adicionar ao mapa de sindicatos
		sindicatos[sindicato] = valor
	}
	
	return nil
}

// processarDiasUteis processa os dados da planilha de dias úteis
func processarDiasUteis(f *excelize.File, diasUteis map[string]int) error {
	// Obter todas as linhas da primeira sheet
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return modelo.NovoErroProcessamentoCompleto(
			"Erro ao ler linhas da planilha de dias úteis",
			"Base dias uteis.xlsx",
			0,
		)
	}
	
	fmt.Printf("Total de linhas na planilha Base dias uteis: %d\n", len(rows))
	
	// Processar cada linha (ignorando o cabeçalho)
	for i, row := range rows {
		// Ignorar a primeira linha (cabeçalho)
		if i == 0 {
			continue
		}
		
		// Verificar se a linha tem dados suficientes (pelo menos 2 colunas)
		if len(row) < 2 {
			continue
		}
		
		// Extrair os dados da linha
		sindicato := strings.TrimSpace(row[0])
		diasStr := strings.TrimSpace(row[1])
		
		// Extrair apenas o nome do estado do sindicato
		// Ex: "SITEPD PR - SIND DOS TRAB EM EMPR PRIVADAS DE PROC DE DADOS DE CURITIBA E REGIAO METROPOLITANA 22" -> "Paraná"
		// Ex: "SINDPPD RS - SINDICATO DOS TRAB. EM PROC. DE DADOS RIO GRANDE DO SUL 21" -> "Rio Grande do Sul"
		// Ex: "SINDPD SP - SIND.TRAB.EM PROC DADOS E EMPR.EMPRESAS PROC DADOS ESTADO DE SP. 22" -> "São Paulo"
		// Ex: "SINDPD RJ - SINDICATO PROFISSIONAIS DE PROC DADOS DO RIO DE JANEIRO 21" -> "Rio de Janeiro"
		
		var estado string
		switch {
		case strings.Contains(sindicato, "PR") || strings.Contains(sindicato, "CURITIBA"):
			estado = "Paraná"
		case strings.Contains(sindicato, "RS"):
			estado = "Rio Grande do Sul"
		case strings.Contains(sindicato, "SP"):
			estado = "São Paulo"
		case strings.Contains(sindicato, "RJ"):
			estado = "Rio de Janeiro"
		default:
			// Se não conseguir identificar, continua para o próximo
			continue
		}
		
		// Extrair dias (último número da string)
		// Ex: "SITEPD PR - SIND DOS TRAB EM EMPR PRIVADAS DE PROC DE DADOS DE CURITIBA E REGIAO METROPOLITANA 22" -> 22
		var dias int
		fmt.Sscanf(sindicato, "%*s %d", &dias)
		if dias == 0 {
			// Tentar extrair o número do final da string
			for i := len(sindicato) - 1; i >= 0; i-- {
				if sindicato[i] >= '0' && sindicato[i] <= '9' {
					// Encontrar o início do número
					start := i
					for start > 0 && sindicato[start-1] >= '0' && sindicato[start-1] <= '9' {
						start--
					}
					diasStr := sindicato[start : i+1]
					dias, _ = strconv.Atoi(diasStr)
					break
				}
			}
		}
		
		// Se ainda não temos dias, tentar da segunda coluna
		if dias == 0 {
			dias, err = strconv.Atoi(diasStr)
			if err != nil {
				// Se não conseguir converter, continua para o próximo
				continue
			}
		}
		
		// Adicionar ao mapa de dias úteis
		diasUteis[estado] = dias
	}
	
	return nil
}