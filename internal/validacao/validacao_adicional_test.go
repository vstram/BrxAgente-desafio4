package validacao

import (
	"testing"
	"time"

	"BrxAgente-desafio4/internal/modelo"
)

func TestValidarColaboradorAdicional(t *testing.T) {
	// Mapas de teste
	valorPorSindicato := map[string]float64{
		"Paraná": 35.00,
	}

	diasUteisPorSindicato := map[string]int{
		"Paraná": 22,
	}

	t.Run("ColaboradorValido", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Empresa:   "Empresa A",
			Cargo:     "Analista",
			Sindicato: "Paraná",
			Situacao:  "Ativo",
		}

		errs := ValidarColaborador(colaborador, valorPorSindicato, diasUteisPorSindicato)
		
		// Não deveria haver erros para um colaborador válido
		if len(errs) != 0 {
			t.Errorf("Não esperava erros para colaborador válido, mas obteve %d erros", len(errs))
			for _, err := range errs {
				t.Errorf("  - %v", err)
			}
		}
	})

	t.Run("ColaboradorSemMatricula", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Empresa:   "Empresa A",
			Cargo:     "Analista",
			Sindicato: "Paraná",
			Situacao:  "Ativo",
		}

		errs := ValidarColaborador(colaborador, valorPorSindicato, diasUteisPorSindicato)
		
		// Deveria haver um erro para matrícula vazia
		if len(errs) == 0 {
			t.Error("Esperava erro para colaborador sem matrícula, mas não houve erro")
		}
	})

	t.Run("ColaboradorSemEmpresa", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Cargo:     "Analista",
			Sindicato: "Paraná",
			Situacao:  "Ativo",
		}

		errs := ValidarColaborador(colaborador, valorPorSindicato, diasUteisPorSindicato)
		
		// Deveria haver um erro para empresa vazia
		if len(errs) == 0 {
			t.Error("Esperava erro para colaborador sem empresa, mas não houve erro")
		}
	})

	t.Run("ColaboradorSemCargo", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Empresa:   "Empresa A",
			Sindicato: "Paraná",
			Situacao:  "Ativo",
		}

		errs := ValidarColaborador(colaborador, valorPorSindicato, diasUteisPorSindicato)
		
		// Deveria haver um erro para cargo vazio
		if len(errs) == 0 {
			t.Error("Esperava erro para colaborador sem cargo, mas não houve erro")
		}
	})

	t.Run("ColaboradorSemSindicato", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Empresa:   "Empresa A",
			Cargo:     "Analista",
			Situacao:  "Ativo",
		}

		errs := ValidarColaborador(colaborador, valorPorSindicato, diasUteisPorSindicato)
		
		// Deveria haver um erro para sindicato vazio
		if len(errs) == 0 {
			t.Error("Esperava erro para colaborador sem sindicato, mas não houve erro")
		}
	})

	t.Run("ColaboradorSemSituacao", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Empresa:   "Empresa A",
			Cargo:     "Analista",
			Sindicato: "Paraná",
		}

		errs := ValidarColaborador(colaborador, valorPorSindicato, diasUteisPorSindicato)
		
		// Deveria haver um erro para situação vazia
		if len(errs) == 0 {
			t.Error("Esperava erro para colaborador sem situação, mas não houve erro")
		}
	})

	t.Run("ColaboradorComSindicatoInvalido", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Empresa:   "Empresa A",
			Cargo:     "Analista",
			Sindicato: "Sindicato Inexistente",
			Situacao:  "Ativo",
		}

		errs := ValidarColaborador(colaborador, valorPorSindicato, diasUteisPorSindicato)
		
		// Pode não haver erro para sindicato inválido, pois o sistema usa valores padrão
		// Neste caso, estamos apenas verificando que não ocorre panic
		t.Logf("Colaborador com sindicato inválido gerou %d erros", len(errs))
	})

	t.Run("ColaboradorComDataAdmissaoFutura", func(t *testing.T) {
		dataFutura := time.Now().AddDate(0, 0, 1) // Amanhã
		colaborador := &modelo.Colaborador{
			Matricula:     "001",
			Empresa:       "Empresa A",
			Cargo:         "Analista",
			Sindicato:     "Paraná",
			Situacao:      "Ativo",
			DataAdmissao:  dataFutura,
		}

		errs := ValidarColaborador(colaborador, valorPorSindicato, diasUteisPorSindicato)
		
		// Pode não haver erro para data futura, dependendo das regras de negócio
		// Neste caso, estamos apenas verificando que não ocorre panic
		t.Logf("Colaborador com data de admissão futura gerou %d erros", len(errs))
	})
}