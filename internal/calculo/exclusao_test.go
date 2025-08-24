package calculo

import (
	"testing"

	"BrxAgente-desafio4/internal/modelo"
)

func TestAplicarRegrasExclusao(t *testing.T) {
	t.Run("NenhumColaboradorExcluido", func(t *testing.T) {
		// Criar colaboradores elegíveis
		colaboradores := map[string]*modelo.Colaborador{
			"001": {
				Matricula: "001",
				Cargo:     "Analista",
				Situacao:  "Trabalhando",
			},
			"002": {
				Matricula: "002",
				Cargo:     "Coordenador",
				Situacao:  "Trabalhando",
			},
		}

		// Mapas vazios (nenhum colaborador para excluir)
		afastamentos := make(map[string]bool)
		aprendizes := make(map[string]bool)
		estagios := make(map[string]bool)
		exterior := make(map[string]bool)

		// Aplicar regras de exclusão
		resultado := AplicarRegrasExclusao(colaboradores, afastamentos, aprendizes, estagios, exterior)

		// Verificar que todos os colaboradores foram mantidos
		if len(resultado) != len(colaboradores) {
			t.Errorf("Esperava %d colaboradores, mas obteve %d", len(colaboradores), len(resultado))
		}
	})

	t.Run("ExcluirDiretor", func(t *testing.T) {
		// Criar colaboradores, incluindo um diretor
		colaboradores := map[string]*modelo.Colaborador{
			"001": {
				Matricula: "001",
				Cargo:     "Analista",
				Situacao:  "Trabalhando",
			},
			"002": {
				Matricula: "002",
				Cargo:     "Diretor de Tecnologia",
				Situacao:  "Trabalhando",
			},
		}

		// Mapas vazios (nenhum colaborador para excluir)
		afastamentos := make(map[string]bool)
		aprendizes := make(map[string]bool)
		estagios := make(map[string]bool)
		exterior := make(map[string]bool)

		// Aplicar regras de exclusão
		resultado := AplicarRegrasExclusao(colaboradores, afastamentos, aprendizes, estagios, exterior)

		// Verificar que apenas 1 colaborador foi mantido (o diretor foi excluído)
		if len(resultado) != 1 {
			t.Errorf("Esperava 1 colaborador, mas obteve %d", len(resultado))
		}

		// Verificar que o colaborador mantido não é o diretor
		if colaborador, existe := resultado["002"]; existe {
			t.Errorf("Esperava que o diretor fosse excluído, mas foi mantido: %+v", colaborador)
		}
	})

	t.Run("ExcluirAfastado", func(t *testing.T) {
		// Criar colaboradores
		colaboradores := map[string]*modelo.Colaborador{
			"001": {
				Matricula: "001",
				Cargo:     "Analista",
				Situacao:  "Trabalhando",
			},
			"002": {
				Matricula: "002",
				Cargo:     "Coordenador",
				Situacao:  "Trabalhando",
			},
		}

		// Mapas com um colaborador afastado
		afastamentos := map[string]bool{"002": true}
		aprendizes := make(map[string]bool)
		estagios := make(map[string]bool)
		exterior := make(map[string]bool)

		// Aplicar regras de exclusão
		resultado := AplicarRegrasExclusao(colaboradores, afastamentos, aprendizes, estagios, exterior)

		// Verificar que apenas 1 colaborador foi mantido (o afastado foi excluído)
		if len(resultado) != 1 {
			t.Errorf("Esperava 1 colaborador, mas obteve %d", len(resultado))
		}

		// Verificar que o colaborador mantido não é o afastado
		if colaborador, existe := resultado["002"]; existe {
			t.Errorf("Esperava que o afastado fosse excluído, mas foi mantido: %+v", colaborador)
		}
	})

	t.Run("ExcluirAprendiz", func(t *testing.T) {
		// Criar colaboradores
		colaboradores := map[string]*modelo.Colaborador{
			"001": {
				Matricula: "001",
				Cargo:     "Analista",
				Situacao:  "Trabalhando",
			},
			"002": {
				Matricula: "002",
				Cargo:     "Aprendiz",
				Situacao:  "Trabalhando",
			},
		}

		// Mapas com um aprendiz
		afastamentos := make(map[string]bool)
		aprendizes := map[string]bool{"002": true}
		estagios := make(map[string]bool)
		exterior := make(map[string]bool)

		// Aplicar regras de exclusão
		resultado := AplicarRegrasExclusao(colaboradores, afastamentos, aprendizes, estagios, exterior)

		// Verificar que apenas 1 colaborador foi mantido (o aprendiz foi excluído)
		if len(resultado) != 1 {
			t.Errorf("Esperava 1 colaborador, mas obteve %d", len(resultado))
		}

		// Verificar que o colaborador mantido não é o aprendiz
		if colaborador, existe := resultado["002"]; existe {
			t.Errorf("Esperava que o aprendiz fosse excluído, mas foi mantido: %+v", colaborador)
		}
	})

	t.Run("ExcluirEstagiario", func(t *testing.T) {
		// Criar colaboradores
		colaboradores := map[string]*modelo.Colaborador{
			"001": {
				Matricula: "001",
				Cargo:     "Analista",
				Situacao:  "Trabalhando",
			},
			"002": {
				Matricula: "002",
				Cargo:     "Estagiário",
				Situacao:  "Trabalhando",
			},
		}

		// Mapas com um estagiário
		afastamentos := make(map[string]bool)
		aprendizes := make(map[string]bool)
		estagios := map[string]bool{"002": true}
		exterior := make(map[string]bool)

		// Aplicar regras de exclusão
		resultado := AplicarRegrasExclusao(colaboradores, afastamentos, aprendizes, estagios, exterior)

		// Verificar que apenas 1 colaborador foi mantido (o estagiário foi excluído)
		if len(resultado) != 1 {
			t.Errorf("Esperava 1 colaborador, mas obteve %d", len(resultado))
		}

		// Verificar que o colaborador mantido não é o estagiário
		if colaborador, existe := resultado["002"]; existe {
			t.Errorf("Esperava que o estagiário fosse excluído, mas foi mantido: %+v", colaborador)
		}
	})

	t.Run("ExcluirExterior", func(t *testing.T) {
		// Criar colaboradores
		colaboradores := map[string]*modelo.Colaborador{
			"001": {
				Matricula: "001",
				Cargo:     "Analista",
				Situacao:  "Trabalhando",
			},
			"002": {
				Matricula: "002",
				Cargo:     "Coordenador",
				Situacao:  "Trabalhando",
			},
		}

		// Mapas com um colaborador no exterior
		afastamentos := make(map[string]bool)
		aprendizes := make(map[string]bool)
		estagios := make(map[string]bool)
		exterior := map[string]bool{"002": true}

		// Aplicar regras de exclusão
		resultado := AplicarRegrasExclusao(colaboradores, afastamentos, aprendizes, estagios, exterior)

		// Verificar que apenas 1 colaborador foi mantido (o do exterior foi excluído)
		if len(resultado) != 1 {
			t.Errorf("Esperava 1 colaborador, mas obteve %d", len(resultado))
		}

		// Verificar que o colaborador mantido não é o do exterior
		if colaborador, existe := resultado["002"]; existe {
			t.Errorf("Esperava que o colaborador do exterior fosse excluído, mas foi mantido: %+v", colaborador)
		}
	})

	t.Run("ExcluirMultiplos", func(t *testing.T) {
		// Criar colaboradores
		colaboradores := map[string]*modelo.Colaborador{
			"001": {
				Matricula: "001",
				Cargo:     "Analista",
				Situacao:  "Trabalhando",
			},
			"002": {
				Matricula: "002",
				Cargo:     "Diretor",
				Situacao:  "Trabalhando",
			},
			"003": {
				Matricula: "003",
				Cargo:     "Aprendiz",
				Situacao:  "Trabalhando",
			},
			"004": {
				Matricula: "004",
				Cargo:     "Estagiário",
				Situacao:  "Trabalhando",
			},
			"005": {
				Matricula: "005",
				Cargo:     "Coordenador",
				Situacao:  "Trabalhando",
			},
		}

		// Mapas com vários colaboradores para excluir
		afastamentos := make(map[string]bool)
		aprendizes := map[string]bool{"003": true}
		estagios := map[string]bool{"004": true}
		exterior := make(map[string]bool)

		// Aplicar regras de exclusão
		resultado := AplicarRegrasExclusao(colaboradores, afastamentos, aprendizes, estagios, exterior)

		// Verificar que apenas 2 colaboradores foram mantidos (todos os outros foram excluídos)
		if len(resultado) != 2 {
			t.Errorf("Esperava 2 colaboradores, mas obteve %d", len(resultado))
		}

		// Verificar que os colaboradores excluídos não estão no resultado
		if _, existe := resultado["002"]; existe {
			t.Error("Esperava que o diretor fosse excluído, mas foi mantido")
		}

		if _, existe := resultado["003"]; existe {
			t.Error("Esperava que o aprendiz fosse excluído, mas foi mantido")
		}

		if _, existe := resultado["004"]; existe {
			t.Error("Esperava que o estagiário fosse excluído, mas foi mantido")
		}
	})
}

func TestDeveExcluirColaborador(t *testing.T) {
	t.Run("NaoExcluir", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Cargo:     "Analista",
			Situacao:  "Trabalhando",
		}

		afastamentos := make(map[string]bool)
		aprendizes := make(map[string]bool)
		estagios := make(map[string]bool)
		exterior := make(map[string]bool)

		excluir := deveExcluirColaborador(colaborador, afastamentos, aprendizes, estagios, exterior)

		if excluir {
			t.Error("Esperava que o colaborador não fosse excluído, mas foi marcado para exclusão")
		}
	})

	t.Run("ExcluirDiretor", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Cargo:     "Diretor de Tecnologia",
			Situacao:  "Trabalhando",
		}

		afastamentos := make(map[string]bool)
		aprendizes := make(map[string]bool)
		estagios := make(map[string]bool)
		exterior := make(map[string]bool)

		excluir := deveExcluirColaborador(colaborador, afastamentos, aprendizes, estagios, exterior)

		if !excluir {
			t.Error("Esperava que o diretor fosse excluído, mas não foi marcado para exclusão")
		}
	})

	t.Run("ExcluirAfastado", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Cargo:     "Analista",
			Situacao:  "Trabalhando",
		}

		afastamentos := map[string]bool{"001": true}
		aprendizes := make(map[string]bool)
		estagios := make(map[string]bool)
		exterior := make(map[string]bool)

		excluir := deveExcluirColaborador(colaborador, afastamentos, aprendizes, estagios, exterior)

		if !excluir {
			t.Error("Esperava que o afastado fosse excluído, mas não foi marcado para exclusão")
		}
	})

	t.Run("ExcluirAprendiz", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Cargo:     "Aprendiz",
			Situacao:  "Trabalhando",
		}

		afastamentos := make(map[string]bool)
		aprendizes := map[string]bool{"001": true}
		estagios := make(map[string]bool)
		exterior := make(map[string]bool)

		excluir := deveExcluirColaborador(colaborador, afastamentos, aprendizes, estagios, exterior)

		if !excluir {
			t.Error("Esperava que o aprendiz fosse excluído, mas não foi marcado para exclusão")
		}
	})

	t.Run("ExcluirEstagiario", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Cargo:     "Estagiário",
			Situacao:  "Trabalhando",
		}

		afastamentos := make(map[string]bool)
		aprendizes := make(map[string]bool)
		estagios := map[string]bool{"001": true}
		exterior := make(map[string]bool)

		excluir := deveExcluirColaborador(colaborador, afastamentos, aprendizes, estagios, exterior)

		if !excluir {
			t.Error("Esperava que o estagiário fosse excluído, mas não foi marcado para exclusão")
		}
	})

	t.Run("ExcluirExterior", func(t *testing.T) {
		colaborador := &modelo.Colaborador{
			Matricula: "001",
			Cargo:     "Analista",
			Situacao:  "Trabalhando",
		}

		afastamentos := make(map[string]bool)
		aprendizes := make(map[string]bool)
		estagios := make(map[string]bool)
		exterior := map[string]bool{"001": true}

		excluir := deveExcluirColaborador(colaborador, afastamentos, aprendizes, estagios, exterior)

		if !excluir {
			t.Error("Esperava que o colaborador do exterior fosse excluído, mas não foi marcado para exclusão")
		}
	})
}

func TestCriarMapaAfastamentos(t *testing.T) {
	afastamentos := map[string]string{
		"001": "Licença Maternidade",
		"002": "Auxílio Doença",
	}

	mapa := CriarMapaAfastamentos(afastamentos)

	if len(mapa) != 2 {
		t.Errorf("Esperava 2 entradas no mapa, mas obteve %d", len(mapa))
	}

	if !mapa["001"] {
		t.Error("Esperava que a matrícula 001 estivesse no mapa")
	}

	if !mapa["002"] {
		t.Error("Esperava que a matrícula 002 estivesse no mapa")
	}
}

func TestCriarMapaAprendizes(t *testing.T) {
	aprendizes := map[string]string{
		"001": "APRENDIZ",
		"002": "APRENDIZ TECH",
	}

	mapa := CriarMapaAprendizes(aprendizes)

	if len(mapa) != 2 {
		t.Errorf("Esperava 2 entradas no mapa, mas obteve %d", len(mapa))
	}

	if !mapa["001"] {
		t.Error("Esperava que a matrícula 001 estivesse no mapa")
	}

	if !mapa["002"] {
		t.Error("Esperava que a matrícula 002 estivesse no mapa")
	}
}

func TestCriarMapaEstagios(t *testing.T) {
	estagios := map[string]string{
		"001": "ESTAGIARIO",
		"002": "ESTAGIARIO #N/A",
	}

	mapa := CriarMapaEstagios(estagios)

	if len(mapa) != 2 {
		t.Errorf("Esperava 2 entradas no mapa, mas obteve %d", len(mapa))
	}

	if !mapa["001"] {
		t.Error("Esperava que a matrícula 001 estivesse no mapa")
	}

	if !mapa["002"] {
		t.Error("Esperava que a matrícula 002 estivesse no mapa")
	}
}

func TestCriarMapaExterior(t *testing.T) {
	exterior := map[string]float64{
		"001": 28.00,
		"002": 554.40,
	}

	mapa := CriarMapaExterior(exterior)

	if len(mapa) != 2 {
		t.Errorf("Esperava 2 entradas no mapa, mas obteve %d", len(mapa))
	}

	if !mapa["001"] {
		t.Error("Esperava que a matrícula 001 estivesse no mapa")
	}

	if !mapa["002"] {
		t.Error("Esperava que a matrícula 002 estivesse no mapa")
	}
}
