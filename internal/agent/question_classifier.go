package agent

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"
)

// QuestionType representa os diferentes tipos de perguntas que o sistema pode classificar
type QuestionType int

const (
	// PolicyQuestion - Perguntas sobre políticas, elegibilidade e regras de VR
	PolicyQuestion QuestionType = iota
	
	// CalculationQuestion - Perguntas sobre cálculos específicos e valores
	CalculationQuestion
	
	// ProcessedDataQuestion - Perguntas sobre dados já processados pelo sistema
	ProcessedDataQuestion
	
	// WhatIfQuestion - Perguntas sobre cenários hipotéticos
	WhatIfQuestion
	
	// ComplianceQuestion - Perguntas sobre conformidade legal e regulamentos
	ComplianceQuestion
	
	// UnknownQuestion - Perguntas que não se encaixam nas categorias conhecidas
	UnknownQuestion
)

// String implementa fmt.Stringer para QuestionType
func (qt QuestionType) String() string {
	switch qt {
	case PolicyQuestion:
		return "PolicyQuestion"
	case CalculationQuestion:
		return "CalculationQuestion"
	case ProcessedDataQuestion:
		return "ProcessedDataQuestion"
	case WhatIfQuestion:
		return "WhatIfQuestion"
	case ComplianceQuestion:
		return "ComplianceQuestion"
	default:
		return "UnknownQuestion"
	}
}

// ToResponseType converte QuestionType para ResponseType
func (qt QuestionType) ToResponseType() ResponseType {
	switch qt {
	case PolicyQuestion, ComplianceQuestion:
		return PolicyResponse
	case CalculationQuestion:
		return CalculationResponse
	case ProcessedDataQuestion:
		return DataResponse
	case WhatIfQuestion:
		return WhatIfResponse
	default:
		return DataResponse // Default para dados processados
	}
}

// PatternMatcher representa um padrão de detecção para classificação
type PatternMatcher struct {
	Keywords []string  // Palavras-chave simples
	Phrases  []string  // Frases exatas
	Regexes  []*regexp.Regexp // Expressões regulares compiladas
	Weight   float64   // Peso deste padrão no score final
}

// ClassificationResult representa o resultado da classificação de uma pergunta
type ClassificationResult struct {
	QuestionType QuestionType `json:"question_type"`
	Confidence   float64      `json:"confidence"`
	Score        float64      `json:"score"`
	Matches      []string     `json:"matches"`      // Padrões que fizeram match
	Alternative  *AlternativeClassification `json:"alternative,omitempty"` // Classificação alternativa se houver
}

// AlternativeClassification representa uma classificação alternativa com menor score
type AlternativeClassification struct {
	QuestionType QuestionType `json:"question_type"`
	Confidence   float64      `json:"confidence"`
	Score        float64      `json:"score"`
}

// QuestionClassifier é responsável por classificar perguntas automaticamente
type QuestionClassifier struct {
	policyPatterns     []PatternMatcher
	calculationPatterns []PatternMatcher
	dataPatterns       []PatternMatcher
	whatIfPatterns     []PatternMatcher
	compliancePatterns []PatternMatcher
	
	// Configuração
	minConfidenceThreshold float64
	multiClassThreshold    float64
	
	logger *log.Logger
}

// NewQuestionClassifier cria uma nova instância do classificador com padrões pré-definidos
func NewQuestionClassifier() *QuestionClassifier {
	qc := &QuestionClassifier{
		minConfidenceThreshold: 0.3,  // Confiança mínima para uma classificação válida
		multiClassThreshold:    0.15, // Diferença mínima entre scores para classificação única
		logger:                log.Default(),
	}
	
	qc.initializePatterns()
	return qc
}

// initializePatterns inicializa todos os padrões de detecção
func (qc *QuestionClassifier) initializePatterns() {
	qc.initializePolicyPatterns()
	qc.initializeCalculationPatterns()
	qc.initializeDataPatterns()
	qc.initializeWhatIfPatterns()
	qc.initializeCompliancePatterns()
}

// initializePolicyPatterns define padrões para perguntas sobre políticas
func (qc *QuestionClassifier) initializePolicyPatterns() {
	qc.policyPatterns = []PatternMatcher{
		{
			Keywords: []string{
				"direito", "pode", "elegível", "permitido", "política", "regra",
				"diretores", "estagiários", "aprendizes", "terceirizados",
				"licença", "afastamento", "férias", "admissão", "desligamento",
				"benefício", "vale refeição", "vale alimentação", "excluído",
				"incluído", "tem direito", "não tem direito", "critério",
			},
			Phrases: []string{
				"tem direito", "não tem direito", "pode receber", "não pode receber",
				"é elegível", "não é elegível", "política de", "regras de",
				"critério para", "requisito para",
			},
			Weight: 0.8,
		},
		{
			Keywords: []string{
				"perfil", "categoria", "cargo", "função", "contrato",
				"clt", "pj", "temporário", "efetivo",
			},
			Weight: 0.6,
		},
	}
}

// initializeCalculationPatterns define padrões para perguntas sobre cálculos
func (qc *QuestionClassifier) initializeCalculationPatterns() {
	qc.calculationPatterns = []PatternMatcher{
		{
			Keywords: []string{
				"como calcular", "valor", "quanto", "proporcional", "cálculo",
				"fórmula", "resultado", "total", "parcial", "dias úteis",
				"período", "mês quebrado", "rateio", "desconto",
			},
			Phrases: []string{
				"como calcular", "qual valor", "quanto vale", "como é calculado",
				"fórmula para", "método de cálculo", "valor proporcional",
			},
			Weight: 0.9,
		},
		{
			Keywords: []string{
				"80%", "20%", "empresa", "colaborador", "empregado", "funcionário",
				"percentual", "porcentagem", "divisão", "split",
			},
			Weight: 0.7,
		},
	}
}

// initializeDataPatterns define padrões para perguntas sobre dados processados
func (qc *QuestionClassifier) initializeDataPatterns() {
	qc.dataPatterns = []PatternMatcher{
		{
			Keywords: []string{
				"quantos", "total", "lista", "colaboradores processados",
				"dados processados", "resultado do processamento", "relatório",
				"planilha", "arquivo", "resultado", "output",
			},
			Phrases: []string{
				"quantos colaboradores", "total de", "lista de", "colaboradores processados",
				"dados do sistema", "resultado do cálculo", "planilha final",
			},
			Weight: 0.85,
		},
		{
			Keywords: []string{
				"consultar", "verificar", "mostrar", "exibir", "buscar",
				"encontrar", "localizar", "filtrar",
			},
			Weight: 0.6,
		},
	}
}

// initializeWhatIfPatterns define padrões para perguntas hipotéticas
func (qc *QuestionClassifier) initializeWhatIfPatterns() {
	qc.whatIfPatterns = []PatternMatcher{
		{
			Keywords: []string{
				"e se", "caso", "supondo", "hipotético", "cenário",
				"simulação", "exemplo", "suponha", "imagine",
			},
			Phrases: []string{
				"e se", "caso o colaborador", "supondo que", "se fosse",
				"em um cenário", "como exemplo", "para simular",
			},
			Weight: 0.9,
		},
		{
			Keywords: []string{
				"fosse", "seria", "aconteceria", "resultaria", "ficaria",
				"mudasse", "alterasse", "modificasse",
			},
			Weight: 0.7,
		},
	}
}

// initializeCompliancePatterns define padrões para perguntas sobre conformidade
func (qc *QuestionClassifier) initializeCompliancePatterns() {
	qc.compliancePatterns = []PatternMatcher{
		{
			Keywords: []string{
				"conforme", "legal", "regulamento", "lei", "clt",
				"trabalhista", "norma", "compliance", "auditoria",
				"fiscal", "tributário", "legislação",
			},
			Phrases: []string{
				"conforme a lei", "está legal", "é legal", "cumpre a", "atende a",
				"de acordo com", "seguindo a", "respeitando a",
			},
			Weight: 0.85,
		},
		{
			Keywords: []string{
				"ministério do trabalho", "receita federal", "inss",
				"sindicato", "convenção coletiva", "acordo coletivo",
			},
			Weight: 0.9,
		},
	}
}

// Classify classifica uma pergunta e retorna o resultado com score de confiança
func (qc *QuestionClassifier) Classify(question string) ClassificationResult {
	if strings.TrimSpace(question) == "" {
		return ClassificationResult{
			QuestionType: UnknownQuestion,
			Confidence:   0.0,
			Score:        0.0,
			Matches:      []string{},
		}
	}
	
	questionLower := strings.ToLower(strings.TrimSpace(question))
	
	// Calcular scores para cada tipo
	scores := map[QuestionType]float64{
		PolicyQuestion:        qc.calculateScore(questionLower, qc.policyPatterns),
		CalculationQuestion:   qc.calculateScore(questionLower, qc.calculationPatterns),
		ProcessedDataQuestion: qc.calculateScore(questionLower, qc.dataPatterns),
		WhatIfQuestion:        qc.calculateScore(questionLower, qc.whatIfPatterns),
		ComplianceQuestion:    qc.calculateScore(questionLower, qc.compliancePatterns),
	}
	
	// Encontrar matches detalhados para debugging
	matches := qc.findMatches(questionLower, scores)
	
	// Encontrar o tipo com maior score
	bestType := UnknownQuestion
	bestScore := 0.0
	secondBestType := UnknownQuestion
	secondBestScore := 0.0
	
	for questionType, score := range scores {
		if score > bestScore {
			secondBestType = bestType
			secondBestScore = bestScore
			bestType = questionType
			bestScore = score
		} else if score > secondBestScore {
			secondBestType = questionType
			secondBestScore = score
		}
	}
	
	// Calcular confiança baseada no score normalizado
	confidence := qc.calculateConfidence(bestScore)
	
	result := ClassificationResult{
		QuestionType: bestType,
		Confidence:   confidence,
		Score:        bestScore,
		Matches:      matches,
	}
	
	// Adicionar classificação alternativa se a diferença for pequena
	if secondBestScore > 0 && (bestScore-secondBestScore) < qc.multiClassThreshold {
		result.Alternative = &AlternativeClassification{
			QuestionType: secondBestType,
			Confidence:   qc.calculateConfidence(secondBestScore),
			Score:        secondBestScore,
		}
	}
	
	// Log da classificação para debugging
	qc.logger.Printf("Question classified: '%s' -> %s (confidence: %.2f, score: %.2f)",
		truncateString(question, 50), result.QuestionType.String(), confidence, bestScore)
	
	return result
}

// calculateScore calcula o score para um conjunto de padrões
func (qc *QuestionClassifier) calculateScore(question string, patterns []PatternMatcher) float64 {
	totalScore := 0.0
	totalWeight := 0.0
	
	for _, pattern := range patterns {
		patternScore := qc.calculatePatternScore(question, pattern)
		totalScore += patternScore * pattern.Weight
		totalWeight += pattern.Weight
	}
	
	if totalWeight == 0 {
		return 0.0
	}
	
	// Normalizar o score pelo peso total
	normalizedScore := totalScore / totalWeight
	
	// Aplicar função sigmoid para suavizar o score
	return qc.applySigmoid(normalizedScore)
}

// calculatePatternScore calcula o score para um padrão específico
func (qc *QuestionClassifier) calculatePatternScore(question string, pattern PatternMatcher) float64 {
	score := 0.0
	matches := 0
	totalPatterns := len(pattern.Keywords) + len(pattern.Phrases) + len(pattern.Regexes)
	
	if totalPatterns == 0 {
		return 0.0
	}
	
	// Verificar keywords
	for _, keyword := range pattern.Keywords {
		if strings.Contains(question, strings.ToLower(keyword)) {
			score += 1.0
			matches++
		}
	}
	
	// Verificar frases exatas (peso maior)
	for _, phrase := range pattern.Phrases {
		if strings.Contains(question, strings.ToLower(phrase)) {
			score += 2.0
			matches++
		}
	}
	
	// Verificar regex patterns (peso ainda maior)
	for _, regex := range pattern.Regexes {
		if regex.MatchString(question) {
			score += 3.0
			matches++
		}
	}
	
	// Normalizar pelo número total de padrões
	normalizedScore := score / float64(totalPatterns)
	
	// Boost para múltiplos matches
	if matches > 1 {
		boostFactor := 1.0 + (float64(matches-1) * 0.2) // 20% de boost por match adicional
		normalizedScore *= boostFactor
	}
	
	return normalizedScore
}

// calculateConfidence calcula a confiança baseada no score
func (qc *QuestionClassifier) calculateConfidence(score float64) float64 {
	if score < qc.minConfidenceThreshold {
		return 0.0
	}
	
	// Mapear score [0,1] para confiança [0,1] com curva não-linear
	confidence := math.Tanh(score * 2.0) // Tanh para suavizar a curva
	
	// Garantir que está no range [0,1]
	if confidence < 0 {
		confidence = 0
	}
	if confidence > 1 {
		confidence = 1
	}
	
	return confidence
}

// applySigmoid aplica função sigmoid para suavizar scores
func (qc *QuestionClassifier) applySigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x*2.0)) // Sigmoid com inclinação ajustada
}

// findMatches encontra os padrões que fizeram match para debugging
func (qc *QuestionClassifier) findMatches(question string, scores map[QuestionType]float64) []string {
	var matches []string
	
	allPatterns := map[string][]PatternMatcher{
		"policy":     qc.policyPatterns,
		"calculation": qc.calculationPatterns,
		"data":       qc.dataPatterns,
		"whatif":     qc.whatIfPatterns,
		"compliance": qc.compliancePatterns,
	}
	
	for patternType, patterns := range allPatterns {
		for _, pattern := range patterns {
			for _, keyword := range pattern.Keywords {
				if strings.Contains(question, strings.ToLower(keyword)) {
					matches = append(matches, fmt.Sprintf("%s:%s", patternType, keyword))
				}
			}
			
			for _, phrase := range pattern.Phrases {
				if strings.Contains(question, strings.ToLower(phrase)) {
					matches = append(matches, fmt.Sprintf("%s:phrase:%s", patternType, phrase))
				}
			}
		}
	}
	
	return matches
}

// SetMinConfidenceThreshold permite ajustar o threshold mínimo de confiança
func (qc *QuestionClassifier) SetMinConfidenceThreshold(threshold float64) {
	qc.minConfidenceThreshold = threshold
}

// SetMultiClassThreshold permite ajustar o threshold para classificação múltipla
func (qc *QuestionClassifier) SetMultiClassThreshold(threshold float64) {
	qc.multiClassThreshold = threshold
}

// GetStats retorna estatísticas básicas do classificador
func (qc *QuestionClassifier) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"policy_patterns":           len(qc.policyPatterns),
		"calculation_patterns":      len(qc.calculationPatterns),
		"data_patterns":            len(qc.dataPatterns),
		"whatif_patterns":          len(qc.whatIfPatterns),
		"compliance_patterns":       len(qc.compliancePatterns),
		"min_confidence_threshold":  qc.minConfidenceThreshold,
		"multi_class_threshold":     qc.multiClassThreshold,
	}
}

// truncateString trunca uma string para debugging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}