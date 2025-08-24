// Package calculo provides functionality for calculating VR values
package calculo

import (
	"strings"

	"BrxAgente-desafio4/internal/modelo"
)

// AplicarRegrasExclusao filtra e remove colaboradores não elegíveis da base consolidada
func AplicarRegrasExclusao(colaboradores map[string]*modelo.Colaborador,
	afastamentos, aprendizes, estagios, exterior map[string]bool) map[string]*modelo.Colaborador {
	// Criar um novo mapa para os colaboradores elegíveis
	elegiveis := make(map[string]*modelo.Colaborador)

	// Verificar cada colaborador
	for matricula, colaborador := range colaboradores {
		// Verificar se o colaborador deve ser excluído
		if deveExcluirColaborador(colaborador, afastamentos, aprendizes, estagios, exterior) {
			continue // Pular este colaborador (excluído)
		}

		// Adicionar colaborador elegível ao novo mapa
		elegiveis[matricula] = colaborador
	}

	return elegiveis
}

// deveExcluirColaborador verifica se um colaborador deve ser excluído com base nas regras de negócio
func deveExcluirColaborador(colaborador *modelo.Colaborador,
	afastamentos, aprendizes, estagios, exterior map[string]bool) bool {
	// 1. Excluir diretores (verificar pelo cargo)
	if strings.Contains(strings.ToUpper(colaborador.Cargo), "DIRETOR") {
		return true
	}

	// 2. Excluir estagiários
	if estagios[colaborador.Matricula] {
		return true
	}

	// 3. Excluir aprendizes
	if aprendizes[colaborador.Matricula] {
		return true
	}

	// 4. Excluir profissionais afastados (licença maternidade, etc.)
	if afastamentos[colaborador.Matricula] {
		return true
	}

	// 5. Excluir profissionais que atuam no exterior
	if exterior[colaborador.Matricula] {
		return true
	}

	// Não excluir
	return false
}

// CriarMapaAfastamentos cria um mapa de matrículas de colaboradores afastados
func CriarMapaAfastamentos(afastamentos map[string]string) map[string]bool {
	mapa := make(map[string]bool)

	for matricula := range afastamentos {
		mapa[matricula] = true
	}

	return mapa
}

// CriarMapaAprendizes cria um mapa de matrículas de aprendizes
func CriarMapaAprendizes(aprendizes map[string]string) map[string]bool {
	mapa := make(map[string]bool)

	for matricula := range aprendizes {
		mapa[matricula] = true
	}

	return mapa
}

// CriarMapaEstagios cria um mapa de matrículas de estagiários
func CriarMapaEstagios(estagios map[string]string) map[string]bool {
	mapa := make(map[string]bool)

	for matricula := range estagios {
		mapa[matricula] = true
	}

	return mapa
}

// CriarMapaExterior cria um mapa de matrículas de colaboradores que atuam no exterior
func CriarMapaExterior(exterior map[string]float64) map[string]bool {
	mapa := make(map[string]bool)

	for matricula := range exterior {
		mapa[matricula] = true
	}

	return mapa
}
