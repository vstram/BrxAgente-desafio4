// Package modelo provides data structures for the application
package modelo

// GetColaboradoresElegiveis filtra e retorna apenas os colaboradores elegíveis para o benefício de VR
func GetColaboradoresElegiveis(colaboradores map[string]*Colaborador) map[string]*Colaborador {
	elegiveis := make(map[string]*Colaborador)

	for matricula, colaborador := range colaboradores {
		if colaborador.EhElegivel() && colaborador.ValorTotalVR > 0 {
			elegiveis[matricula] = colaborador
		}
	}

	return elegiveis
}
