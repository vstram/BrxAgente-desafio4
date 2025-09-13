package agent

import "time"

// RealWorldTestScenario representa um cenário de teste real
type RealWorldTestScenario struct {
	ID           string            `json:"id"`
	Category     string            `json:"category"`
	Question     string            `json:"question"`
	Expected     ExpectedResponse  `json:"expected"`
	Metadata     map[string]string `json:"metadata"`
	Difficulty   int               `json:"difficulty"` // 1-5 escala
	Tags         []string          `json:"tags"`
	CreatedAt    time.Time         `json:"created_at"`
}

// ExpectedResponse representa a resposta esperada para um cenário
type ExpectedResponse struct {
	Type           string            `json:"type"`
	Answer         string            `json:"answer"`
	KeyPoints      []string          `json:"key_points"`
	MustContain    []string          `json:"must_contain"`
	MustNotContain []string          `json:"must_not_contain"`
	MinConfidence  float64           `json:"min_confidence"`
	MaxResponseTime time.Duration    `json:"max_response_time"`
	PolicyRefs     []string          `json:"policy_refs"`
	Calculations   map[string]string `json:"calculations"`
}

// LoadRealWorldDataset carrega o dataset completo de cenários de teste
func LoadRealWorldDataset() []RealWorldTestScenario {
	return []RealWorldTestScenario{
		// 1. POLÍTICAS DE ELEGIBILIDADE
		{
			ID:       "eligibility_001",
			Category: "Políticas de Elegibilidade",
			Question: "Diretores têm direito a VR?",
			Expected: ExpectedResponse{
				Type:   "policy",
				Answer: "Não. Diretores são excluídos do benefício de Vale Refeição conforme política da empresa.",
				KeyPoints: []string{
					"Exclusão de diretores",
					"Política empresarial",
					"Sem direito ao benefício",
				},
				MustContain:     []string{"diretor", "não", "excluíd", "polític"},
				MustNotContain:  []string{"tem direito", "recebe"},
				MinConfidence:   0.9,
				MaxResponseTime: 2 * time.Second,
				PolicyRefs:      []string{"política_elegibilidade"},
			},
			Metadata: map[string]string{
				"source": "manual_corporativo",
				"section": "elegibilidade",
			},
			Difficulty: 1,
			Tags:      []string{"elegibilidade", "diretor", "exclusão"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "eligibility_002", 
			Category: "Políticas de Elegibilidade",
			Question: "Estagiários podem receber Vale Refeição?",
			Expected: ExpectedResponse{
				Type:   "policy",
				Answer: "Não. Estagiários são excluídos do benefício de Vale Refeição.",
				KeyPoints: []string{
					"Estagiários excluídos",
					"Sem direito ao VR",
					"Política de exclusão",
				},
				MustContain:     []string{"estagiário", "não", "excluíd"},
				MustNotContain:  []string{"tem direito", "recebe VR"},
				MinConfidence:   0.9,
				MaxResponseTime: 2 * time.Second,
				PolicyRefs:      []string{"política_elegibilidade"},
			},
			Metadata: map[string]string{
				"source": "manual_corporativo",
				"section": "elegibilidade",
			},
			Difficulty: 1,
			Tags:      []string{"elegibilidade", "estagiário", "exclusão"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "eligibility_003",
			Category: "Políticas de Elegibilidade", 
			Question: "Terceirizados recebem VR?",
			Expected: ExpectedResponse{
				Type:   "policy",
				Answer: "Não. Colaboradores terceirizados não recebem Vale Refeição da empresa.",
				KeyPoints: []string{
					"Terceirizados excluídos",
					"Sem VR da empresa",
					"Responsabilidade da empresa terceirizada",
				},
				MustContain:     []string{"terceirizado", "não", "empresa"},
				MustNotContain:  []string{"recebe", "tem direito"},
				MinConfidence:   0.85,
				MaxResponseTime: 2 * time.Second,
				PolicyRefs:      []string{"política_elegibilidade"},
			},
			Metadata: map[string]string{
				"source": "manual_corporativo",
				"section": "elegibilidade",
			},
			Difficulty: 2,
			Tags:      []string{"elegibilidade", "terceirizado", "exclusão"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "eligibility_004",
			Category: "Políticas de Elegibilidade",
			Question: "Aprendizes menores de 18 anos têm direito?",
			Expected: ExpectedResponse{
				Type:   "policy",
				Answer: "Não. Aprendizes são excluídos do benefício de Vale Refeição.",
				KeyPoints: []string{
					"Aprendizes excluídos",
					"Idade irrelevante para exclusão",
					"Categoria profissional determina exclusão",
				},
				MustContain:     []string{"aprendiz", "não", "excluíd"},
				MustNotContain:  []string{"tem direito", "idade"},
				MinConfidence:   0.9,
				MaxResponseTime: 2 * time.Second,
				PolicyRefs:      []string{"política_elegibilidade"},
			},
			Metadata: map[string]string{
				"source": "manual_corporativo",
				"section": "elegibilidade",
			},
			Difficulty: 2,
			Tags:      []string{"elegibilidade", "aprendiz", "menor", "exclusão"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "eligibility_005",
			Category: "Políticas de Elegibilidade",
			Question: "Colaborador com 4 horas diárias tem direito?",
			Expected: ExpectedResponse{
				Type:   "policy",
				Answer: "Sim. Colaboradores com jornada de 4 horas ou mais têm direito ao Vale Refeição.",
				KeyPoints: []string{
					"4 horas = direito ao VR",
					"Jornada mínima atendida",
					"Elegível para benefício",
				},
				MustContain:     []string{"4 horas", "direito", "sim"},
				MustNotContain:  []string{"não tem", "excluído"},
				MinConfidence:   0.9,
				MaxResponseTime: 2 * time.Second,
				PolicyRefs:      []string{"política_jornada"},
			},
			Metadata: map[string]string{
				"source": "manual_corporativo", 
				"section": "jornada_trabalho",
			},
			Difficulty: 2,
			Tags:      []string{"elegibilidade", "jornada", "4_horas", "direito"},
			CreatedAt: time.Now(),
		},

		// 2. CÁLCULOS ESPECÍFICOS
		{
			ID:       "calculation_001",
			Category: "Cálculos Específicos",
			Question: "Como calcular VR para licença médica de 20 dias?",
			Expected: ExpectedResponse{
				Type:   "calculation",
				Answer: "Durante licença médica, o colaborador mantém direito ao VR. Para 20 dias de afastamento, calcular VR proporcional aos dias trabalhados no mês.",
				KeyPoints: []string{
					"Licença médica mantém direito",
					"Cálculo proporcional",
					"Dias trabalhados vs dias afastados",
				},
				MustContain:     []string{"licença médica", "mantém", "proporcional", "dias"},
				MustNotContain:  []string{"perde direito", "sem VR"},
				MinConfidence:   0.85,
				MaxResponseTime: 3 * time.Second,
				PolicyRefs:      []string{"política_licença_médica", "cálculo_vr"},
				Calculations: map[string]string{
					"formula": "VR = (dias_trabalhados / dias_úteis_mês) * valor_integral",
				},
			},
			Metadata: map[string]string{
				"source": "manual_cálculos",
				"complexity": "medium",
			},
			Difficulty: 3,
			Tags:      []string{"cálculo", "licença_médica", "proporcional"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "calculation_002",
			Category: "Cálculos Específicos",
			Question: "Colaborador admitido dia 25, qual valor?",
			Expected: ExpectedResponse{
				Type:   "calculation",
				Answer: "Colaborador admitido no dia 25 receberá VR proporcional aos dias trabalhados no mês de admissão.",
				KeyPoints: []string{
					"Admissão tardia no mês",
					"VR proporcional",
					"Contar apenas dias trabalhados",
				},
				MustContain:     []string{"admitido", "proporcional", "dias trabalhados"},
				MustNotContain:  []string{"valor integral", "mês completo"},
				MinConfidence:   0.9,
				MaxResponseTime: 2 * time.Second,
				PolicyRefs:      []string{"política_admissão", "cálculo_vr"},
				Calculations: map[string]string{
					"exemplo": "Se há 22 dias úteis no mês e trabalhou 5 dias (26-30), VR = (5/22) * valor_base",
				},
			},
			Metadata: map[string]string{
				"source": "manual_cálculos",
				"scenario": "admissão_tardía",
			},
			Difficulty: 3,
			Tags:      []string{"cálculo", "admissão", "proporcional", "data_quebrada"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "calculation_003",
			Category: "Cálculos Específicos", 
			Question: "Desligamento comunicado dia 20, tem VR?",
			Expected: ExpectedResponse{
				Type:   "policy_calculation",
				Answer: "Não. Desligamento comunicado após o dia 15 não dá direito ao VR do mês seguinte.",
				KeyPoints: []string{
					"Regra dos 15 dias",
					"Comunicação tardia",
					"Sem VR no mês seguinte",
				},
				MustContain:     []string{"dia 15", "não", "mês seguinte", "comunicado"},
				MustNotContain:  []string{"tem direito", "recebe VR"},
				MinConfidence:   0.95,
				MaxResponseTime: 2 * time.Second,
				PolicyRefs:      []string{"política_desligamento"},
			},
			Metadata: map[string]string{
				"source": "manual_corporativo",
				"rule": "regra_15_dias",
			},
			Difficulty: 2,
			Tags:      []string{"desligamento", "regra_15_dias", "comunicação"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "calculation_004",
			Category: "Cálculos Específicos",
			Question: "Férias de 30 dias, mantém VR?",
			Expected: ExpectedResponse{
				Type:   "policy",
				Answer: "Sim. Durante férias, o colaborador mantém o direito integral ao Vale Refeição.",
				KeyPoints: []string{
					"Férias mantém VR integral",
					"Não há desconto",
					"Direito preservado",
				},
				MustContain:     []string{"férias", "mantém", "integral"},
				MustNotContain:  []string{"proporcional", "desconto"},
				MinConfidence:   0.95,
				MaxResponseTime: 2 * time.Second,
				PolicyRefs:      []string{"política_férias"},
			},
			Metadata: map[string]string{
				"source": "manual_corporativo",
				"section": "férias",
			},
			Difficulty: 2,
			Tags:      []string{"férias", "vr_integral", "direito_preservado"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "calculation_005",
			Category: "Cálculos Específicos",
			Question: "Afastamento INSS 10 dias, como calcular?",
			Expected: ExpectedResponse{
				Type:   "calculation",
				Answer: "Afastamento pelo INSS mantém direito ao VR. Calcular valor proporcional aos dias trabalhados no mês.",
				KeyPoints: []string{
					"INSS mantém direito",
					"Cálculo proporcional",
					"Considerar apenas dias trabalhados",
				},
				MustContain:     []string{"INSS", "mantém", "proporcional"},
				MustNotContain:  []string{"perde direito", "sem VR"},
				MinConfidence:   0.9,
				MaxResponseTime: 3 * time.Second,
				PolicyRefs:      []string{"política_inss", "cálculo_vr"},
			},
			Metadata: map[string]string{
				"source": "manual_cálculos",
				"type": "afastamento_inss",
			},
			Difficulty: 3,
			Tags:      []string{"inss", "afastamento", "proporcional"},
			CreatedAt: time.Now(),
		},

		// 3. CENÁRIOS COMPLEXOS
		{
			ID:       "complex_001",
			Category: "Cenários Complexos",
			Question: "Colaborador admitido dia 10 e afastado 20 dias",
			Expected: ExpectedResponse{
				Type:   "complex_calculation",
				Answer: "Cenário complexo: colaborador trabalhou poucos dias devido à admissão tardia e afastamento. VR será proporcional aos dias efetivamente trabalhados.",
				KeyPoints: []string{
					"Dupla situação: admissão + afastamento",
					"Contar apenas dias efetivamente trabalhados",
					"Cálculo proporcional complexo",
				},
				MustContain:     []string{"admissão", "afastamento", "proporcional", "dias trabalhados"},
				MustNotContain:  []string{"valor integral"},
				MinConfidence:   0.8,
				MaxResponseTime: 4 * time.Second,
				PolicyRefs:      []string{"política_admissão", "política_afastamento", "cálculo_complexo"},
			},
			Metadata: map[string]string{
				"complexity": "high",
				"multiple_events": "true",
			},
			Difficulty: 4,
			Tags:      []string{"complexo", "admissão", "afastamento", "múltiplos_eventos"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "complex_002",
			Category: "Cenários Complexos",
			Question: "Licença maternidade no meio do mês",
			Expected: ExpectedResponse{
				Type:   "policy_calculation",
				Answer: "Licença maternidade mantém direito integral ao VR durante todo o período, independente de quando iniciar no mês.",
				KeyPoints: []string{
					"Maternidade = VR integral",
					"Data de início irrelevante",
					"Direito preservado totalmente",
				},
				MustContain:     []string{"maternidade", "integral", "mantém"},
				MustNotContain:  []string{"proporcional", "desconto"},
				MinConfidence:   0.95,
				MaxResponseTime: 3 * time.Second,
				PolicyRefs:      []string{"política_maternidade"},
			},
			Metadata: map[string]string{
				"source": "legislação_trabalhista",
				"protection_level": "high",
			},
			Difficulty: 3,
			Tags:      []string{"maternidade", "vr_integral", "proteção_legal"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "complex_003",
			Category: "Cenários Complexos",
			Question: "Acidente de trabalho + férias programadas",
			Expected: ExpectedResponse{
				Type:   "complex_policy",
				Answer: "Tanto acidente de trabalho quanto férias mantêm direito integral ao VR. Colaborador receberá valor completo.",
				KeyPoints: []string{
					"Acidente de trabalho = VR integral",
					"Férias = VR integral",
					"Dupla proteção do direito",
				},
				MustContain:     []string{"acidente", "férias", "integral", "mantém"},
				MustNotContain:  []string{"proporcional", "redução"},
				MinConfidence:   0.9,
				MaxResponseTime: 3 * time.Second,
				PolicyRefs:      []string{"política_acidente", "política_férias"},
			},
			Metadata: map[string]string{
				"complexity": "medium",
				"multiple_protections": "true",
			},
			Difficulty: 3,
			Tags:      []string{"acidente_trabalho", "férias", "proteção_dupla"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "complex_004",
			Category: "Cenários Complexos", 
			Question: "Desligamento + licença médica simultâneos",
			Expected: ExpectedResponse{
				Type:   "complex_policy",
				Answer: "Situação complexa que requer análise da data de comunicação do desligamento vs início da licença. Regras de desligamento podem se sobrepor à licença.",
				KeyPoints: []string{
					"Conflito entre regras",
					"Analisar datas específicas",
					"Precedência do desligamento",
				},
				MustContain:     []string{"desligamento", "licença", "datas", "análise"},
				MustNotContain:  []string{"simples", "automático"},
				MinConfidence:   0.7,
				MaxResponseTime: 5 * time.Second,
				PolicyRefs:      []string{"política_desligamento", "política_licença", "resolução_conflitos"},
			},
			Metadata: map[string]string{
				"complexity": "very_high",
				"conflict_resolution": "true",
			},
			Difficulty: 5,
			Tags:      []string{"conflito_regras", "desligamento", "licença", "análise_complexa"},
			CreatedAt: time.Now(),
		},

		// 4. DADOS PROCESSADOS
		{
			ID:       "data_001",
			Category: "Dados Processados",
			Question: "Quantos colaboradores foram processados?",
			Expected: ExpectedResponse{
				Type:   "data_query",
				Answer: "Foram processados X colaboradores no último processamento.",
				KeyPoints: []string{
					"Número total de colaboradores",
					"Dados do último processamento",
					"Estatística consolidada",
				},
				MustContain:     []string{"processados", "colaboradores", "total"},
				MustNotContain:  []string{"erro", "falha"},
				MinConfidence:   0.9,
				MaxResponseTime: 2 * time.Second,
			},
			Metadata: map[string]string{
				"query_type": "count",
				"source": "processed_data",
			},
			Difficulty: 1,
			Tags:      []string{"dados", "contagem", "estatística"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "data_002",
			Category: "Dados Processados",
			Question: "Qual o total de VR este mês?",
			Expected: ExpectedResponse{
				Type:   "data_calculation",
				Answer: "O valor total de VR processado neste mês é R$ X,XX.",
				KeyPoints: []string{
					"Valor financeiro total",
					"Soma de todos os VRs",
					"Período mensal",
				},
				MustContain:     []string{"total", "VR", "R$", "mês"},
				MustNotContain:  []string{"erro", "indisponível"},
				MinConfidence:   0.95,
				MaxResponseTime: 2 * time.Second,
			},
			Metadata: map[string]string{
				"query_type": "sum",
				"financial": "true",
			},
			Difficulty: 2,
			Tags:      []string{"dados", "valor_total", "financeiro"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "data_003",
			Category: "Dados Processados",
			Question: "Maior valor de VR individual?",
			Expected: ExpectedResponse{
				Type:   "data_analysis",
				Answer: "O maior valor individual de VR processado é R$ X,XX.",
				KeyPoints: []string{
					"Valor máximo individual",
					"Análise de extremos",
					"Identificação de outliers",
				},
				MustContain:     []string{"maior", "individual", "R$"},
				MustNotContain:  []string{"média", "total"},
				MinConfidence:   0.9,
				MaxResponseTime: 2 * time.Second,
			},
			Metadata: map[string]string{
				"query_type": "max",
				"analysis": "extremes",
			},
			Difficulty: 2,
			Tags:      []string{"dados", "máximo", "análise"},
			CreatedAt: time.Now(),
		},

		{
			ID:       "data_004",
			Category: "Dados Processados",
			Question: "Distribuição por sindicato",
			Expected: ExpectedResponse{
				Type:   "data_distribution",
				Answer: "A distribuição dos colaboradores por sindicato é: SINDPD: X colaboradores, Outros: Y colaboradores.",
				KeyPoints: []string{
					"Agrupamento por sindicato",
					"Distribuição quantitativa",
					"Análise demográfica",
				},
				MustContain:     []string{"sindicato", "distribuição", "colaboradores"},
				MustNotContain:  []string{"erro", "não encontrado"},
				MinConfidence:   0.85,
				MaxResponseTime: 3 * time.Second,
			},
			Metadata: map[string]string{
				"query_type": "group_by",
				"demographic": "true",
			},
			Difficulty: 3,
			Tags:      []string{"dados", "distribuição", "sindicato", "agrupamento"},
			CreatedAt: time.Now(),
		},
	}
}

// GetScenariosByCategory retorna cenários filtrados por categoria
func GetScenariosByCategory(category string) []RealWorldTestScenario {
	allScenarios := LoadRealWorldDataset()
	var filtered []RealWorldTestScenario
	
	for _, scenario := range allScenarios {
		if scenario.Category == category {
			filtered = append(filtered, scenario)
		}
	}
	
	return filtered
}

// GetScenariosByDifficulty retorna cenários filtrados por nível de dificuldade
func GetScenariosByDifficulty(difficulty int) []RealWorldTestScenario {
	allScenarios := LoadRealWorldDataset()
	var filtered []RealWorldTestScenario
	
	for _, scenario := range allScenarios {
		if scenario.Difficulty == difficulty {
			filtered = append(filtered, scenario)
		}
	}
	
	return filtered
}

// GetScenariosByTags retorna cenários que contêm todas as tags especificadas
func GetScenariosByTags(tags []string) []RealWorldTestScenario {
	allScenarios := LoadRealWorldDataset()
	var filtered []RealWorldTestScenario
	
	for _, scenario := range allScenarios {
		hasAllTags := true
		for _, requiredTag := range tags {
			found := false
			for _, scenarioTag := range scenario.Tags {
				if scenarioTag == requiredTag {
					found = true
					break
				}
			}
			if !found {
				hasAllTags = false
				break
			}
		}
		if hasAllTags {
			filtered = append(filtered, scenario)
		}
	}
	
	return filtered
}