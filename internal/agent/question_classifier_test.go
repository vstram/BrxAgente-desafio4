package agent

import (
	"fmt"
	"testing"
)

func TestNewQuestionClassifier(t *testing.T) {
	qc := NewQuestionClassifier()
	
	if qc == nil {
		t.Fatal("NewQuestionClassifier() retornou nil")
	}
	
	// Verificar se padrões foram inicializados
	if len(qc.policyPatterns) == 0 {
		t.Error("policyPatterns não foram inicializados")
	}
	if len(qc.calculationPatterns) == 0 {
		t.Error("calculationPatterns não foram inicializados")
	}
	if len(qc.dataPatterns) == 0 {
		t.Error("dataPatterns não foram inicializados")
	}
	if len(qc.whatIfPatterns) == 0 {
		t.Error("whatIfPatterns não foram inicializados")
	}
	if len(qc.compliancePatterns) == 0 {
		t.Error("compliancePatterns não foram inicializados")
	}
}

func TestQuestionType_String(t *testing.T) {
	tests := []struct {
		questionType QuestionType
		expected     string
	}{
		{PolicyQuestion, "PolicyQuestion"},
		{CalculationQuestion, "CalculationQuestion"},
		{ProcessedDataQuestion, "ProcessedDataQuestion"},
		{WhatIfQuestion, "WhatIfQuestion"},
		{ComplianceQuestion, "ComplianceQuestion"},
		{UnknownQuestion, "UnknownQuestion"},
		{QuestionType(999), "UnknownQuestion"}, // Valor inválido
	}
	
	for _, test := range tests {
		if result := test.questionType.String(); result != test.expected {
			t.Errorf("QuestionType(%d).String() = %s, esperado %s", 
				test.questionType, result, test.expected)
		}
	}
}

func TestClassify_PolicyQuestions(t *testing.T) {
	qc := NewQuestionClassifier()
	
	testCases := []struct {
		question           string
		expectedType       QuestionType
		minConfidence      float64
		shouldHaveMatches  bool
	}{
		{
			question:          "Diretores têm direito a VR?",
			expectedType:      PolicyQuestion,
			minConfidence:     0.8,
			shouldHaveMatches: true,
		},
		{
			question:          "Estagiários podem receber VR?",
			expectedType:      PolicyQuestion,
			minConfidence:     0.7,
			shouldHaveMatches: true,
		},
		{
			question:          "Qual a política de VR para terceirizados?",
			expectedType:      PolicyQuestion,
			minConfidence:     0.8,
			shouldHaveMatches: true,
		},
		{
			question:          "Colaborador com licença médica tem direito ao benefício?",
			expectedType:      PolicyQuestion,
			minConfidence:     0.7,
			shouldHaveMatches: true,
		},
		{
			question:          "Quem é elegível para receber vale refeição?",
			expectedType:      PolicyQuestion,
			minConfidence:     0.8,
			shouldHaveMatches: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.question, func(t *testing.T) {
			result := qc.Classify(tc.question)
			
			if result.QuestionType != tc.expectedType {
				t.Errorf("Pergunta: '%s'\nEsperado: %s\nObtido: %s\nScore: %.2f\nConfiança: %.2f\nMatches: %v",
					tc.question, tc.expectedType.String(), result.QuestionType.String(),
					result.Score, result.Confidence, result.Matches)
			}
			
			if result.Confidence < tc.minConfidence {
				t.Errorf("Confiança muito baixa para '%s': %.2f (esperado >= %.2f)",
					tc.question, result.Confidence, tc.minConfidence)
			}
			
			if tc.shouldHaveMatches && len(result.Matches) == 0 {
				t.Errorf("Esperava matches para '%s', mas nenhum foi encontrado", tc.question)
			}
		})
	}
}

func TestClassify_CalculationQuestions(t *testing.T) {
	qc := NewQuestionClassifier()
	
	testCases := []struct {
		question           string
		expectedType       QuestionType
		minConfidence      float64
	}{
		{
			question:      "Como calcular VR para licença médica de 20 dias?",
			expectedType:  CalculationQuestion,
			minConfidence: 0.8,
		},
		{
			question:      "Qual valor de VR para admissão dia 20?",
			expectedType:  CalculationQuestion,
			minConfidence: 0.7,
		},
		{
			question:      "Quanto vale o VR proporcional este mês?",
			expectedType:  CalculationQuestion,
			minConfidence: 0.8,
		},
		{
			question:      "Como é feito o rateio 80%/20% do VR?",
			expectedType:  CalculationQuestion,
			minConfidence: 0.7,
		},
		{
			question:      "Fórmula para calcular dias úteis de VR?",
			expectedType:  CalculationQuestion,
			minConfidence: 0.8,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.question, func(t *testing.T) {
			result := qc.Classify(tc.question)
			
			if result.QuestionType != tc.expectedType {
				t.Errorf("Pergunta: '%s'\nEsperado: %s\nObtido: %s\nScore: %.2f\nConfiança: %.2f",
					tc.question, tc.expectedType.String(), result.QuestionType.String(),
					result.Score, result.Confidence)
			}
			
			if result.Confidence < tc.minConfidence {
				t.Errorf("Confiança muito baixa para '%s': %.2f (esperado >= %.2f)",
					tc.question, result.Confidence, tc.minConfidence)
			}
		})
	}
}

func TestClassify_ProcessedDataQuestions(t *testing.T) {
	qc := NewQuestionClassifier()
	
	testCases := []struct {
		question           string
		expectedType       QuestionType
		minConfidence      float64
	}{
		{
			question:      "Quantos colaboradores foram processados?",
			expectedType:  ProcessedDataQuestion,
			minConfidence: 0.8,
		},
		{
			question:      "Total de VR processado este mês?",
			expectedType:  ProcessedDataQuestion,
			minConfidence: 0.7,
		},
		{
			question:      "Lista de colaboradores processados no sistema",
			expectedType:  ProcessedDataQuestion,
			minConfidence: 0.8,
		},
		{
			question:      "Mostrar dados processados da planilha",
			expectedType:  ProcessedDataQuestion,
			minConfidence: 0.7,
		},
		{
			question:      "Buscar colaborador na base processada",
			expectedType:  ProcessedDataQuestion,
			minConfidence: 0.6,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.question, func(t *testing.T) {
			result := qc.Classify(tc.question)
			
			if result.QuestionType != tc.expectedType {
				t.Errorf("Pergunta: '%s'\nEsperado: %s\nObtido: %s\nScore: %.2f\nConfiança: %.2f",
					tc.question, tc.expectedType.String(), result.QuestionType.String(),
					result.Score, result.Confidence)
			}
			
			if result.Confidence < tc.minConfidence {
				t.Errorf("Confiança muito baixa para '%s': %.2f (esperado >= %.2f)",
					tc.question, result.Confidence, tc.minConfidence)
			}
		})
	}
}

func TestClassify_WhatIfQuestions(t *testing.T) {
	qc := NewQuestionClassifier()
	
	testCases := []struct {
		question           string
		expectedType       QuestionType
		minConfidence      float64
	}{
		{
			question:      "E se o colaborador fosse admitido dia 10?",
			expectedType:  WhatIfQuestion,
			minConfidence: 0.8,
		},
		{
			question:      "Caso o funcionário entrasse de férias, qual seria o VR?",
			expectedType:  WhatIfQuestion,
			minConfidence: 0.7,
		},
		{
			question:      "Supondo que houvesse desligamento no dia 15",
			expectedType:  WhatIfQuestion,
			minConfidence: 0.8,
		},
		{
			question:      "Em um cenário hipotético de licença médica",
			expectedType:  WhatIfQuestion,
			minConfidence: 0.7,
		},
		{
			question:      "Simulação: se mudasse o sindicato do colaborador",
			expectedType:  WhatIfQuestion,
			minConfidence: 0.8,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.question, func(t *testing.T) {
			result := qc.Classify(tc.question)
			
			if result.QuestionType != tc.expectedType {
				t.Errorf("Pergunta: '%s'\nEsperado: %s\nObtido: %s\nScore: %.2f\nConfiança: %.2f",
					tc.question, tc.expectedType.String(), result.QuestionType.String(),
					result.Score, result.Confidence)
			}
			
			if result.Confidence < tc.minConfidence {
				t.Errorf("Confiança muito baixa para '%s': %.2f (esperado >= %.2f)",
					tc.question, result.Confidence, tc.minConfidence)
			}
		})
	}
}

func TestClassify_ComplianceQuestions(t *testing.T) {
	qc := NewQuestionClassifier()
	
	testCases := []struct {
		question           string
		expectedType       QuestionType
		minConfidence      float64
	}{
		{
			question:      "Está conforme a CLT?",
			expectedType:  ComplianceQuestion,
			minConfidence: 0.8,
		},
		{
			question:      "A política está legal perante a legislação trabalhista?",
			expectedType:  ComplianceQuestion,
			minConfidence: 0.7,
		},
		{
			question:      "Cumpre as normas do Ministério do Trabalho?",
			expectedType:  ComplianceQuestion,
			minConfidence: 0.8,
		},
		{
			question:      "De acordo com a convenção coletiva do sindicato",
			expectedType:  ComplianceQuestion,
			minConfidence: 0.7,
		},
		{
			question:      "Para fins de auditoria fiscal, está correto?",
			expectedType:  ComplianceQuestion,
			minConfidence: 0.7,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.question, func(t *testing.T) {
			result := qc.Classify(tc.question)
			
			if result.QuestionType != tc.expectedType {
				t.Errorf("Pergunta: '%s'\nEsperado: %s\nObtido: %s\nScore: %.2f\nConfiança: %.2f",
					tc.question, tc.expectedType.String(), result.QuestionType.String(),
					result.Score, result.Confidence)
			}
			
			if result.Confidence < tc.minConfidence {
				t.Errorf("Confiança muito baixa para '%s': %.2f (esperado >= %.2f)",
					tc.question, result.Confidence, tc.minConfidence)
			}
		})
	}
}

func TestClassify_EdgeCases(t *testing.T) {
	qc := NewQuestionClassifier()
	
	testCases := []struct {
		question           string
		expectedType       QuestionType
		maxConfidence      float64 // Para casos de baixa confiança
	}{
		{
			question:      "",
			expectedType:  UnknownQuestion,
			maxConfidence: 0.0,
		},
		{
			question:      "   ",
			expectedType:  UnknownQuestion,
			maxConfidence: 0.0,
		},
		{
			question:      "Olá, como vai?",
			expectedType:  UnknownQuestion,
			maxConfidence: 0.3,
		},
		{
			question:      "abc def ghi",
			expectedType:  UnknownQuestion,
			maxConfidence: 0.3,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.question+"_edge", func(t *testing.T) {
			result := qc.Classify(tc.question)
			
			if result.QuestionType != tc.expectedType {
				t.Errorf("Pergunta: '%s'\nEsperado: %s\nObtido: %s\nScore: %.2f\nConfiança: %.2f",
					tc.question, tc.expectedType.String(), result.QuestionType.String(),
					result.Score, result.Confidence)
			}
			
			if result.Confidence > tc.maxConfidence {
				t.Errorf("Confiança muito alta para '%s': %.2f (esperado <= %.2f)",
					tc.question, result.Confidence, tc.maxConfidence)
			}
		})
	}
}

func TestClassify_MultipleClassification(t *testing.T) {
	qc := NewQuestionClassifier()
	
	// Pergunta que pode ter múltiplas classificações
	question := "Como calcular VR conforme a política da empresa para diretores?"
	result := qc.Classify(question)
	
	// Deve classificar como uma das duas opções principais
	if result.QuestionType != PolicyQuestion && result.QuestionType != CalculationQuestion {
		t.Errorf("Pergunta com múltiplas classificações não foi classificada adequadamente: %s",
			result.QuestionType.String())
	}
	
	// Deve ter classificação alternativa
	if result.Alternative == nil {
		t.Error("Esperava classificação alternativa para pergunta com múltiplos aspectos")
	}
	
	// A diferença de score deve ser pequena
	if result.Alternative != nil {
		scoreDiff := result.Score - result.Alternative.Score
		if scoreDiff > qc.multiClassThreshold*2 { // Tolerância para teste
			t.Errorf("Diferença de score muito grande entre classificações: %.2f", scoreDiff)
		}
	}
}

func TestClassify_PerformanceRequirement(t *testing.T) {
	qc := NewQuestionClassifier()
	question := "Como calcular VR para licença médica de 20 dias?"
	
	// Executar várias classificações para verificar performance
	for i := 0; i < 100; i++ {
		result := qc.Classify(question)
		if result.QuestionType == UnknownQuestion {
			t.Errorf("Classificação inconsistente na iteração %d", i)
		}
	}
	
	// Este teste verifica consistência, não tempo exato, pois o tempo
	// pode variar dependendo do ambiente de teste
}

func TestGetStats(t *testing.T) {
	qc := NewQuestionClassifier()
	stats := qc.GetStats()
	
	requiredKeys := []string{
		"policy_patterns",
		"calculation_patterns", 
		"data_patterns",
		"whatif_patterns",
		"compliance_patterns",
		"min_confidence_threshold",
		"multi_class_threshold",
	}
	
	for _, key := range requiredKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("Estatística '%s' não encontrada", key)
		}
	}
	
	// Verificar se há padrões carregados
	if stats["policy_patterns"].(int) == 0 {
		t.Error("Nenhum padrão de política carregado")
	}
}

func TestSetConfidenceThresholds(t *testing.T) {
	qc := NewQuestionClassifier()
	
	// Testar configuração de threshold de confiança mínima
	qc.SetMinConfidenceThreshold(0.5)
	stats := qc.GetStats()
	if stats["min_confidence_threshold"].(float64) != 0.5 {
		t.Error("Threshold de confiança mínima não foi definido corretamente")
	}
	
	// Testar configuração de threshold de classificação múltipla
	qc.SetMultiClassThreshold(0.2)
	stats = qc.GetStats()
	if stats["multi_class_threshold"].(float64) != 0.2 {
		t.Error("Threshold de classificação múltipla não foi definido corretamente")
	}
}

func TestClassificationAccuracy(t *testing.T) {
	qc := NewQuestionClassifier()
	
	// Conjunto de teste abrangente para medir precisão geral
	testSuite := []struct {
		question     string
		expectedType QuestionType
	}{
		// Políticas (10 casos)
		{"Diretores têm direito a VR?", PolicyQuestion},
		{"Estagiários podem receber vale refeição?", PolicyQuestion},
		{"Qual a política para terceirizados?", PolicyQuestion},
		{"Quem é elegível para o benefício?", PolicyQuestion},
		{"Colaborador com licença tem direito?", PolicyQuestion},
		{"Regras de VR para aprendizes", PolicyQuestion},
		{"Critério para receber vale alimentação", PolicyQuestion},
		{"Política de afastamento médico", PolicyQuestion},
		{"Funcionários CLT podem receber?", PolicyQuestion},
		{"Perfil excluído do benefício", PolicyQuestion},
		
		// Cálculos (10 casos)
		{"Como calcular VR proporcional?", CalculationQuestion},
		{"Qual valor para admissão dia 15?", CalculationQuestion},
		{"Fórmula de rateio 80/20", CalculationQuestion},
		{"Quanto vale o VR este mês?", CalculationQuestion},
		{"Como calcular dias úteis?", CalculationQuestion},
		{"Valor proporcional para férias", CalculationQuestion},
		{"Cálculo para desligamento", CalculationQuestion},
		{"Total de VR do colaborador", CalculationQuestion},
		{"Percentual empresa/funcionário", CalculationQuestion},
		{"Método de cálculo mensal", CalculationQuestion},
		
		// Dados processados (10 casos)
		{"Quantos colaboradores processados?", ProcessedDataQuestion},
		{"Lista de dados processados", ProcessedDataQuestion},
		{"Total processado no sistema", ProcessedDataQuestion},
		{"Buscar colaborador processado", ProcessedDataQuestion},
		{"Mostrar resultado do processamento", ProcessedDataQuestion},
		{"Consultar planilha final", ProcessedDataQuestion},
		{"Dados do último processamento", ProcessedDataQuestion},
		{"Relatório de colaboradores", ProcessedDataQuestion},
		{"Verificar processamento", ProcessedDataQuestion},
		{"Arquivo de saída gerado", ProcessedDataQuestion},
		
		// What-if (5 casos)
		{"E se fosse admitido dia 10?", WhatIfQuestion},
		{"Caso entrasse de férias", WhatIfQuestion},
		{"Supondo desligamento dia 20", WhatIfQuestion},
		{"Em cenário hipotético", WhatIfQuestion},
		{"Simulação de licença", WhatIfQuestion},
		
		// Compliance (5 casos)
		{"Está conforme a CLT?", ComplianceQuestion},
		{"Legal perante legislação?", ComplianceQuestion},
		{"Cumpre normas trabalhistas?", ComplianceQuestion},
		{"De acordo com sindicato", ComplianceQuestion},
		{"Para auditoria fiscal", ComplianceQuestion},
	}
	
	correct := 0
	total := len(testSuite)
	var errors []string
	
	for _, test := range testSuite {
		result := qc.Classify(test.question)
		if result.QuestionType == test.expectedType {
			correct++
		} else {
			errors = append(errors, 
				fmt.Sprintf("'%s': esperado %s, obtido %s (score: %.2f)",
					test.question, test.expectedType.String(), 
					result.QuestionType.String(), result.Score))
		}
	}
	
	accuracy := float64(correct) / float64(total)
	
	t.Logf("Precisão da classificação: %.1f%% (%d/%d corretos)", 
		accuracy*100, correct, total)
	
	// Exigir pelo menos 85% de precisão conforme especificado na issue
	if accuracy < 0.85 {
		t.Errorf("Precisão abaixo do esperado: %.1f%% (esperado >= 85%%)", accuracy*100)
		for _, err := range errors {
			t.Logf("Erro: %s", err)
		}
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a very long string", 10, "this is a ..."},
		{"exactly10c", 10, "exactly10c"},
		{"", 5, ""},
	}
	
	for _, test := range tests {
		result := truncateString(test.input, test.maxLen)
		if result != test.expected {
			t.Errorf("truncateString('%s', %d) = '%s', esperado '%s'",
				test.input, test.maxLen, result, test.expected)
		}
	}
}

// BenchmarkClassify mede performance da classificação
func BenchmarkClassify(b *testing.B) {
	qc := NewQuestionClassifier()
	question := "Como calcular VR para licença médica de 20 dias?"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qc.Classify(question)
	}
}

func TestClassify_CaseSensitivity(t *testing.T) {
	qc := NewQuestionClassifier()
	
	// Testar que classificação é case-insensitive
	questions := []string{
		"DIRETORES TÊM DIREITO A VR?",
		"diretores têm direito a vr?",
		"Diretores Têm Direito A VR?",
		"DiReToReS tÊm DiReItO a Vr?",
	}
	
	expectedType := PolicyQuestion
	for _, question := range questions {
		result := qc.Classify(question)
		if result.QuestionType != expectedType {
			t.Errorf("Case sensitivity falhou para '%s': esperado %s, obtido %s",
				question, expectedType.String(), result.QuestionType.String())
		}
	}
}