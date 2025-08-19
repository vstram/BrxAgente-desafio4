// Package feriados provides functionality for handling Brazilian holidays
package feriados

import (
	"time"
)

// Feriado represents a holiday with its date and description
type Feriado struct {
	Data        time.Time
	Descricao   string
	Tipo        string // "nacional", "estadual", "municipal"
	Estado      string // For state holidays
	Municipio   string // For municipal holidays
}

// ObterFeriadosNacionais obtém os feriados nacionais para um determinado ano
func ObterFeriadosNacionais(ano int) []Feriado {
	feriados := []Feriado{}

	// Páscoa (calculada)
	pascoa := calcularPascoa(ano)
	feriados = append(feriados, Feriado{
		Data:      pascoa,
		Descricao: "Páscoa",
		Tipo:      "nacional",
	})

	// Carnaval (47 dias antes da Páscoa)
	carnaval := pascoa.AddDate(0, 0, -47)
	feriados = append(feriados, Feriado{
		Data:      carnaval,
		Descricao: "Carnaval",
		Tipo:      "nacional",
	})

	// Sexta-feira Santa (2 dias antes da Páscoa)
	sextaSanta := pascoa.AddDate(0, 0, -2)
	feriados = append(feriados, Feriado{
		Data:      sextaSanta,
		Descricao: "Sexta-feira Santa",
		Tipo:      "nacional",
	})

	// Corpus Christi (60 dias após a Páscoa)
	corpusChristi := pascoa.AddDate(0, 0, 60)
	feriados = append(feriados, Feriado{
		Data:      corpusChristi,
		Descricao: "Corpus Christi",
		Tipo:      "nacional",
	})

	// Ano Novo (1 de janeiro)
	feriados = append(feriados, Feriado{
		Data:      time.Date(ano, time.January, 1, 0, 0, 0, 0, time.UTC),
		Descricao: "Ano Novo",
		Tipo:      "nacional",
	})

	// Tiradentes (21 de abril)
	feriados = append(feriados, Feriado{
		Data:      time.Date(ano, time.April, 21, 0, 0, 0, 0, time.UTC),
		Descricao: "Tiradentes",
		Tipo:      "nacional",
	})

	// Dia do Trabalho (1 de maio)
	feriados = append(feriados, Feriado{
		Data:      time.Date(ano, time.May, 1, 0, 0, 0, 0, time.UTC),
		Descricao: "Dia do Trabalho",
		Tipo:      "nacional",
	})

	// Independência (7 de setembro)
	feriados = append(feriados, Feriado{
		Data:      time.Date(ano, time.September, 7, 0, 0, 0, 0, time.UTC),
		Descricao: "Independência do Brasil",
		Tipo:      "nacional",
	})

	// Nossa Senhora Aparecida (12 de outubro)
	feriados = append(feriados, Feriado{
		Data:      time.Date(ano, time.October, 12, 0, 0, 0, 0, time.UTC),
		Descricao: "Nossa Senhora Aparecida",
		Tipo:      "nacional",
	})

	// Finados (2 de novembro)
	feriados = append(feriados, Feriado{
		Data:      time.Date(ano, time.November, 2, 0, 0, 0, 0, time.UTC),
		Descricao: "Finados",
		Tipo:      "nacional",
	})

	// República (15 de novembro)
	feriados = append(feriados, Feriado{
		Data:      time.Date(ano, time.November, 15, 0, 0, 0, 0, time.UTC),
		Descricao: "Proclamação da República",
		Tipo:      "nacional",
	})

	// Natal (25 de dezembro)
	feriados = append(feriados, Feriado{
		Data:      time.Date(ano, time.December, 25, 0, 0, 0, 0, time.UTC),
		Descricao: "Natal",
		Tipo:      "nacional",
	})

	return feriados
}

// ObterFeriadosEstaduais obtém os feriados estaduais para um determinado estado e ano
func ObterFeriadosEstaduais(estado string, ano int) []Feriado {
	feriados := []Feriado{}

	// Feriados estaduais comuns
	switch estado {
	case "SP": // São Paulo
		feriados = append(feriados, Feriado{
			Data:      time.Date(ano, time.July, 9, 0, 0, 0, 0, time.UTC),
			Descricao: "Revolução Constitucionalista",
			Tipo:      "estadual",
			Estado:    "SP",
		})
	case "RJ": // Rio de Janeiro
		feriados = append(feriados, Feriado{
			Data:      time.Date(ano, time.April, 23, 0, 0, 0, 0, time.UTC),
			Descricao: "Dia de São Jorge",
			Tipo:      "estadual",
			Estado:    "RJ",
		})
	case "PR": // Paraná
		feriados = append(feriados, Feriado{
			Data:      time.Date(ano, time.December, 19, 0, 0, 0, 0, time.UTC),
			Descricao: "Emancipação Política do Estado do Paraná",
			Tipo:      "estadual",
			Estado:    "PR",
		})
	case "RS": // Rio Grande do Sul
		feriados = append(feriados, Feriado{
			Data:      time.Date(ano, time.September, 20, 0, 0, 0, 0, time.UTC),
			Descricao: "Dia do Gaúcho",
			Tipo:      "estadual",
			Estado:    "RS",
		})
	}

	return feriados
}

// ObterFeriadosMunicipais obtém os feriados municipais para um determinado município e ano
// Esta é uma implementação básica que pode ser expandida conforme necessário
func ObterFeriadosMunicipais(municipio string, estado string, ano int) []Feriado {
	feriados := []Feriado{}

	// Exemplos de feriados municipais (pode ser expandido)
	switch municipio {
	case "São Paulo":
		feriados = append(feriados, Feriado{
			Data:      time.Date(ano, time.January, 25, 0, 0, 0, 0, time.UTC),
			Descricao: "Aniversário da Cidade de São Paulo",
			Tipo:      "municipal",
			Estado:    "SP",
			Municipio: "São Paulo",
		})
	case "Rio de Janeiro":
		feriados = append(feriados, Feriado{
			Data:      time.Date(ano, time.January, 20, 0, 0, 0, 0, time.UTC),
			Descricao: "Dia de São Sebastião",
			Tipo:      "municipal",
			Estado:    "RJ",
			Municipio: "Rio de Janeiro",
		})
	}

	return feriados
}

// EhFeriado verifica se uma data específica é feriado
func EhFeriado(data time.Time, feriados []Feriado) bool {
	for _, feriado := range feriados {
		if feriado.Data.Year() == data.Year() && 
		   feriado.Data.Month() == data.Month() && 
		   feriado.Data.Day() == data.Day() {
			return true
		}
	}
	return false
}

// ContarFeriadosNoPeriodo conta quantos feriados existem em um período específico
func ContarFeriadosNoPeriodo(inicio, fim time.Time, feriados []Feriado) int {
	contagem := 0
	for _, feriado := range feriados {
		if !feriado.Data.Before(inicio) && !feriado.Data.After(fim) {
			contagem++
		}
	}
	return contagem
}

// calcularPascoa calcula a data da Páscoa usando o algoritmo de Meeus/Jones/Butcher
func calcularPascoa(ano int) time.Time {
	a := ano % 19
	b := ano / 100
	c := ano % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	mes := (h + l - 7*m + 114) / 31
	dia := ((h + l - 7*m + 114) % 31) + 1

	return time.Date(ano, time.Month(mes), dia, 0, 0, 0, 0, time.UTC)
}