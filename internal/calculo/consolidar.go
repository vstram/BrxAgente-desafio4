// Package calculo provides functionality for calculating VR values
package calculo

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	
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
	
	// Mapas para armazenar dados das planilhas de exclusão
	afastamentos := make(map[string]string)
	aprendizes := make(map[string]string)
	estagios := make(map[string]string)
	exterior := make(map[string]float64)
	
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
	
	// 6. Ler a planilha de AFASTAMENTOS
	caminhoAfastamentos := filepath.Join(diretorioPlanilhas, "AFASTAMENTOS.xlsx")
	fAfastamentos, err := excel.LerPlanilha(caminhoAfastamentos)
	if err != nil {
		// Se não encontrar a planilha, continuamos sem ela
		fmt.Println("Aviso: Planilha de AFASTAMENTOS não encontrada")
	} else {
		defer fAfastamentos.Close()
		
		// Processar dados da planilha de afastamentos
		err = processarAfastamentos(fAfastamentos, afastamentos)
		if err != nil {
			return nil, err
		}
	}
	
	// 7. Ler a planilha de APRENDIZ
	caminhoAprendizes := filepath.Join(diretorioPlanilhas, "APRENDIZ.xlsx")
	fAprendizes, err := excel.LerPlanilha(caminhoAprendizes)
	if err != nil {
		// Se não encontrar a planilha, continuamos sem ela
		fmt.Println("Aviso: Planilha de APRENDIZ não encontrada")
	} else {
		defer fAprendizes.Close()
		
		// Processar dados da planilha de aprendizes
		err = processarAprendizes(fAprendizes, aprendizes)
		if err != nil {
			return nil, err
		}
	}
	
	// 8. Ler a planilha de ESTÁGIO
	caminhoEstagios := filepath.Join(diretorioPlanilhas, "ESTÁGIO.xlsx")
	fEstagios, err := excel.LerPlanilha(caminhoEstagios)
	if err != nil {
		// Se não encontrar a planilha, continuamos sem ela
		fmt.Println("Aviso: Planilha de ESTÁGIO não encontrada")
	} else {
		defer fEstagios.Close()
		
		// Processar dados da planilha de estagiários
		err = processarEstagios(fEstagios, estagios)
		if err != nil {
			return nil, err
		}
	}
	
	// 9. Ler a planilha de EXTERIOR
	caminhoExterior := filepath.Join(diretorioPlanilhas, "EXTERIOR.xlsx")
	fExterior, err := excel.LerPlanilha(caminhoExterior)
	if err != nil {
		// Se não encontrar a planilha, continuamos sem ela
		fmt.Println("Aviso: Planilha de EXTERIOR não encontrada")
	} else {
		defer fExterior.Close()
		
		// Processar dados da planilha de exterior
		err = processarExterior(fExterior, exterior)
		if err != nil {
			return nil, err
		}
	}
	
	// 10. Ler a planilha de ADMISSÃO
	caminhoAdmissao := filepath.Join(diretorioPlanilhas, "ADMISSÃO ABRIL.xlsx")
	fAdmissao, err := excel.LerPlanilha(caminhoAdmissao)
	if err != nil {
		// Se não encontrar a planilha, continuamos sem ela
		fmt.Println("Aviso: Planilha de ADMISSÃO não encontrada")
	} else {
		defer fAdmissao.Close()
		
		// Processar dados da planilha de admissões
		err = processarAdmissao(fAdmissao, colaboradores)
		if err != nil {
			return nil, err
		}
	}
	
	// 11. Aplicar regras de exclusão
	afastamentosMap := CriarMapaAfastamentos(afastamentos)
	aprendizesMap := CriarMapaAprendizes(aprendizes)
	estagiosMap := CriarMapaEstagios(estagios)
	exteriorMap := CriarMapaExterior(exterior)
	
	colaboradores = AplicarRegrasExclusao(colaboradores, afastamentosMap, aprendizesMap, estagiosMap, exteriorMap)
	
	// 12. Validar os dados consolidados
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
		
		// Verificar se a linha tem dados suficientes (pelo menos 4 colunas)
		if len(row) < 4 {
			continue
		}
		
		// Extrair os dados da linha
		matricula := strings.TrimSpace(row[0])
		situacao := strings.TrimSpace(row[1])
		dataInicioStr := strings.TrimSpace(row[2])
		dataFimStr := strings.TrimSpace(row[3])
		
		// Parse das datas
		dataInicio, err := time.Parse("01-02-06", dataInicioStr) // MM-DD-YY
		if err != nil {
			// Se não conseguir converter, continua para o próximo
			continue
		}
		
		// Ajustar o ano para 2025 (assumindo que YY=25 é 2025)
		dataInicio = time.Date(2000+dataInicio.Year()%100, dataInicio.Month(), dataInicio.Day(), 0, 0, 0, 0, time.UTC)
		
		dataFim, err := time.Parse("01-02-06", dataFimStr) // MM-DD-YY
		if err != nil {
			// Se não conseguir converter, continua para o próximo
			continue
		}
		
		// Ajustar o ano para 2025 (assumindo que YY=25 é 2025)
		dataFim = time.Date(2000+dataFim.Year()%100, dataFim.Month(), dataFim.Day(), 0, 0, 0, 0, time.UTC)
		
		// Verificar se o colaborador existe
		if colaborador, existe := colaboradores[matricula]; existe {
			// Atualizar situação do colaborador
			colaborador.Situacao = situacao
			
			// Adicionar período de férias ao colaborador
			ferias := modelo.Periodo{
				Inicio: dataInicio,
				Fim:    dataFim,
			}
			colaborador.Ferias = append(colaborador.Ferias, ferias)
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
		
		// Extrair data de desligamento se disponível (coluna 1)
		var dataDesligamento *time.Time
		if len(row) >= 2 {
			dataDesligamentoStr := strings.TrimSpace(row[1])
			if dataDesligamentoStr != "" {
				// Parse da data de desligamento
				data, err := parseDataDesligamento(dataDesligamentoStr)
				if err == nil {
					dataDesligamento = &data
				}
			}
		}
		
		// Extrair data de comunicação de desligamento se disponível (coluna 2)
		var dataComunicacao *time.Time
		if len(row) >= 3 {
			dataComunicacaoStr := strings.TrimSpace(row[2])
			if dataComunicacaoStr != "" {
				// Parse da data de comunicação de desligamento
				data, err := parseDataDesligamento(dataComunicacaoStr)
				if err == nil {
					dataComunicacao = &data
				}
			}
		}
		
		// Verificar se o colaborador existe
		if colaborador, existe := colaboradores[matricula]; existe {
			// Atualizar situação do colaborador como desligado
			colaborador.Situacao = "Desligado"
			
			// Atualizar data de desligamento do colaborador
			colaborador.DataDesligamento = dataDesligamento
			
			// Atualizar data de comunicação de desligamento do colaborador
			colaborador.DataComunicacaoDesligamento = dataComunicacao
		}
		// Se o colaborador não existe, ignoramos (pode ser de outro mês)
	}
	
	return nil
}

// parseDataDesligamento converte uma string de data no formato MM-DD-YY para time.Time
func parseDataDesligamento(dataStr string) (time.Time, error) {
	// Formato esperado: MM-DD-YY (ex: 05-01-25)
	// Precisamos converter para 2025 (assumindo que YY=25 é 2025)
	
	parts := strings.Split(dataStr, "-")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("formato de data inválido: %s", dataStr)
	}
	
	mes, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("mês inválido: %s", parts[0])
	}
	
	dia, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("dia inválido: %s", parts[1])
	}
	
	ano, err := strconv.Atoi(parts[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("ano inválido: %s", parts[2])
	}
	
	// Assumindo que o ano é 2025 (para YY=25)
	anoCompleto := 2000 + ano
	
	// Criar a data
	data := time.Date(anoCompleto, time.Month(mes), dia, 0, 0, 0, 0, time.UTC)
	
	return data, nil
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

// processarExterior processa os dados da planilha de colaboradores no exterior
func processarExterior(f *excelize.File, exterior map[string]float64) error {
	// Obter todas as linhas da primeira sheet
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return modelo.NovoErroProcessamentoCompleto(
			"Erro ao ler linhas da planilha de exterior",
			"EXTERIOR.xlsx",
			0,
		)
	}
	
	fmt.Printf("Total de linhas na planilha EXTERIOR: %d\n", len(rows))
	
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
		
		// Extrair valor se disponível
		var valor float64
		if len(row) >= 2 {
			valorStr := strings.TrimSpace(row[1])
			// Remover caracteres não numéricos
			valorStr = strings.ReplaceAll(valorStr, "R$", "")
			valorStr = strings.ReplaceAll(valorStr, ",", ".")
			valorStr = strings.TrimSpace(valorStr)
			
			if valorStr != "" && valorStr != "#N/A" {
				valor, _ = strconv.ParseFloat(valorStr, 64)
			}
		}
		
		// Adicionar ao mapa de exterior
		exterior[matricula] = valor
	}
	
	return nil
}

// processarEstagios processa os dados da planilha de estagiários
func processarEstagios(f *excelize.File, estagios map[string]string) error {
	// Obter todas as linhas da primeira sheet
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return modelo.NovoErroProcessamentoCompleto(
			"Erro ao ler linhas da planilha de estagiários",
			"ESTÁGIO.xlsx",
			0,
		)
	}
	
	fmt.Printf("Total de linhas na planilha ESTÁGIO: %d\n", len(rows))
	
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
		
		// Extrair cargo se disponível
		var cargo string
		if len(row) >= 2 {
			cargo = strings.TrimSpace(row[1])
		}
		
		// Adicionar ao mapa de estagiários
		estagios[matricula] = cargo
	}
	
	return nil
}

// processarAprendizes processa os dados da planilha de aprendizes
func processarAprendizes(f *excelize.File, aprendizes map[string]string) error {
	// Obter todas as linhas da primeira sheet
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return modelo.NovoErroProcessamentoCompleto(
			"Erro ao ler linhas da planilha de aprendizes",
			"APRENDIZ.xlsx",
			0,
		)
	}
	
	fmt.Printf("Total de linhas na planilha APRENDIZ: %d\n", len(rows))
	
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
		
		// Extrair cargo se disponível
		var cargo string
		if len(row) >= 2 {
			cargo = strings.TrimSpace(row[1])
		}
		
		// Adicionar ao mapa de aprendizes
		aprendizes[matricula] = cargo
	}
	
	return nil
}

// processarAfastamentos processa os dados da planilha de afastamentos
func processarAfastamentos(f *excelize.File, afastamentos map[string]string) error {
	// Obter todas as linhas da primeira sheet
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return modelo.NovoErroProcessamentoCompleto(
			"Erro ao ler linhas da planilha de afastamentos",
			"AFASTAMENTOS.xlsx",
			0,
		)
	}
	
	fmt.Printf("Total de linhas na planilha AFASTAMENTOS: %d\n", len(rows))
	
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
		
		// Extrair situação se disponível
		var situacao string
		if len(row) >= 2 {
			situacao = strings.TrimSpace(row[1])
		}
		
		// Adicionar ao mapa de afastamentos
		afastamentos[matricula] = situacao
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

// processarAdmissao processa os dados da planilha de admissões
func processarAdmissao(f *excelize.File, colaboradores map[string]*modelo.Colaborador) error {
	// Obter todas as linhas da primeira sheet
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return modelo.NovoErroProcessamentoCompleto(
			"Erro ao ler linhas da planilha de admissões",
			"ADMISSÃO ABRIL.xlsx",
			0,
		)
	}
	
	fmt.Printf("Total de linhas na planilha ADMISSÃO: %d\n", len(rows))
	
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
		matricula := strings.TrimSpace(row[0])
		dataAdmissaoStr := strings.TrimSpace(row[1])
		
		// Parse da data de admissão
		// Formato esperado: MM-DD-YY (ex: 04-07-25)
		dataAdmissao, err := parseDataAdmissao(dataAdmissaoStr)
		if err != nil {
			// Se não conseguir converter, continua para o próximo
			continue
		}
		
		// Verificar se o colaborador existe
		if colaborador, existe := colaboradores[matricula]; existe {
			// Atualizar data de admissão do colaborador
			colaborador.DataAdmissao = dataAdmissao
		}
		// Se o colaborador não existe, ignoramos (pode ser de outro mês)
	}
	
	return nil
}

// parseDataAdmissao converte uma string de data no formato MM-DD-YY para time.Time
func parseDataAdmissao(dataStr string) (time.Time, error) {
	// Formato esperado: MM-DD-YY (ex: 04-07-25)
	// Precisamos converter para 2025 (assumindo que YY=25 é 2025)
	
	parts := strings.Split(dataStr, "-")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("formato de data inválido: %s", dataStr)
	}
	
	mes, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("mês inválido: %s", parts[0])
	}
	
	dia, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("dia inválido: %s", parts[1])
	}
	
	ano, err := strconv.Atoi(parts[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("ano inválido: %s", parts[2])
	}
	
	// Assumindo que o ano é 2025 (para YY=25)
	anoCompleto := 2000 + ano
	
	// Criar a data
	data := time.Date(anoCompleto, time.Month(mes), dia, 0, 0, 0, 0, time.UTC)
	
	return data, nil
}