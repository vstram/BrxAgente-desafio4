package validacao

import (
	"testing"
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

func TestValidarCamposObrigatorios(t *testing.T) {
	t.Run("TodosOsCamposPreenchidos", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "12345",
			Empresa:   "Empresa Teste",
			Cargo:     "Analista",
			Sindicato: "Sindicato Teste",
			Situacao:  "Ativo",
		}

		err := ValidarCamposObrigatorios(colaborador)
		if err != nil {
			t.Errorf("Não esperava erro, mas obteve: %v", err)
		}
	})

	t.Run("MatriculaVazia", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Empresa:   "Empresa Teste",
			Cargo:     "Analista",
			Sindicato: "Sindicato Teste",
			Situacao:  "Ativo",
		}

		err := ValidarCamposObrigatorios(colaborador)
		if err == nil {
			t.Error("Esperava erro por matrícula vazia, mas não obteve")
		}
	})

	t.Run("EmpresaVazia", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "12345",
			Cargo:     "Analista",
			Sindicato: "Sindicato Teste",
			Situacao:  "Ativo",
		}

		err := ValidarCamposObrigatorios(colaborador)
		if err == nil {
			t.Error("Esperava erro por empresa vazia, mas não obteve")
		}
	})

	t.Run("CargoVazio", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "12345",
			Empresa:   "Empresa Teste",
			Sindicato: "Sindicato Teste",
			Situacao:  "Ativo",
		}

		err := ValidarCamposObrigatorios(colaborador)
		if err == nil {
			t.Error("Esperava erro por cargo vazio, mas não obteve")
		}
	})

	t.Run("SindicatoVazio", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "12345",
			Empresa:   "Empresa Teste",
			Cargo:     "Analista",
			Situacao:  "Ativo",
		}

		err := ValidarCamposObrigatorios(colaborador)
		if err == nil {
			t.Error("Esperava erro por sindicato vazio, mas não obteve")
		}
	})

	t.Run("SituacaoVazia", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "12345",
			Empresa:   "Empresa Teste",
			Cargo:     "Analista",
			Sindicato: "Sindicato Teste",
		}

		err := ValidarCamposObrigatorios(colaborador)
		if err == nil {
			t.Error("Esperava erro por situação vazia, mas não obteve")
		}
	})
}

func TestValidarDatas(t *testing.T) {
	dataAdmissao := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	dataDesligamento := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	dataDesligamentoAnterior := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("DatasValidas", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula:        "12345",
			DataAdmissao:     dataAdmissao,
			DataDesligamento: &dataDesligamento,
		}

		err := ValidarDatas(colaborador)
		if err != nil {
			t.Errorf("Não esperava erro, mas obteve: %v", err)
		}
	})

	t.Run("DataAdmissaoZero", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "12345",
		}

		err := ValidarDatas(colaborador)
		if err != nil {
			t.Errorf("Não esperava erro por data de admissão zero, mas obteve: %v", err)
		}
	})

	t.Run("DataDesligamentoAnteriorAAdmissao", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula:        "12345",
			DataAdmissao:     dataAdmissao,
			DataDesligamento: &dataDesligamentoAnterior,
		}

		err := ValidarDatas(colaborador)
		if err == nil {
			t.Error("Esperava erro por data de desligamento anterior à admissão, mas não obteve")
		}
	})

	t.Run("AfastamentoComDatasInvalidas", func(t *testing.T) {
		dataInicio := time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC)
		dataFim := time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC) // Fim antes do início

		colaborador := &modelo.Colaborador{
			Matricula:    "12345",
			DataAdmissao: dataAdmissao,
			Afastamentos: []modelo.Periodo{
				{
					Inicio: dataInicio,
					Fim:    dataFim,
				},
			},
		}

		err := ValidarDatas(colaborador)
		if err == nil {
			t.Error("Esperava erro por datas de afastamento inválidas, mas não obteve")
		}
	})

	t.Run("FeriasComDatasInvalidas", func(t *testing.T) {
		dataInicio := time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC)
		dataFim := time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC) // Fim antes do início

		colaborador := &modelo.Colaborador{
			Matricula:    "12345",
			DataAdmissao: dataAdmissao,
			Ferias: []modelo.Periodo{
				{
					Inicio: dataInicio,
					Fim:    dataFim,
				},
			},
		}

		err := ValidarDatas(colaborador)
		if err == nil {
			t.Error("Esperava erro por datas de férias inválidas, mas não obteve")
		}
	})
}

func TestValidarFormatacao(t *testing.T) {
	t.Run("MatriculaVazia", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "",
		}

		err := ValidarFormatacao(colaborador)
		if err == nil {
			t.Error("Esperava erro por matrícula vazia, mas não obteve")
		}
	})

	t.Run("MatriculaPreenchida", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "12345",
		}

		err := ValidarFormatacao(colaborador)
		if err != nil {
			t.Errorf("Não esperava erro, mas obteve: %v", err)
		}
	})
}

func TestValidarRelacionamentos(t *testing.T) {
	t.Run("SindicatoExistente", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "12345",
			Sindicato: "Sindicato Teste",
		}

		sindicatos := map[string]float64{
			"Sindicato Teste": 25.50,
		}

		diasUteis := map[string]int{
			"Sindicato Teste": 22,
		}

		err := ValidarRelacionamentos(colaborador, sindicatos, diasUteis)
		if err != nil {
			t.Errorf("Não esperava erro, mas obteve: %v", err)
		}
	})

	t.Run("SindicatoInexistente", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "12345",
			Sindicato: "Sindicato Inexistente",
		}

		sindicatos := map[string]float64{
			"Sindicato Teste": 25.50,
		}

		diasUteis := map[string]int{
			"Sindicato Teste": 22,
		}

		err := ValidarRelacionamentos(colaborador, sindicatos, diasUteis)
		if err == nil {
			t.Error("Esperava erro por sindicato inexistente, mas não obteve")
		}
	})
}

func TestValidarColaborador(t *testing.T) {
	dataAdmissao := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	dataDesligamento := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	colaborador := &modelo.Colaborador{
		Matricula:        "12345",
		Empresa:          "Empresa Teste",
		Cargo:            "Analista",
		Sindicato:        "Sindicato Teste",
		Situacao:         "Ativo",
		DataAdmissao:     dataAdmissao,
		DataDesligamento: &dataDesligamento,
	}

	sindicatos := map[string]float64{
		"Sindicato Teste": 25.50,
	}

	diasUteis := map[string]int{
		"Sindicato Teste": 22,
	}

	erros := ValidarColaborador(colaborador, sindicatos, diasUteis)
	if len(erros) > 0 {
		t.Errorf("Não esperava erros, mas obteve %d erros: %v", len(erros), erros)
	}
}
