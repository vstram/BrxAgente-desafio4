package agent

import (
	"strings"
	"testing"
	"time"
)

// TestNewResponseFormatter testa a criação de um novo formatador
func TestNewResponseFormatter(t *testing.T) {
	// Testar com configuração nil (usar padrão)
	formatter := NewResponseFormatter(nil)
	if formatter == nil {
		t.Fatal("NewResponseFormatter should not return nil")
	}

	if formatter.config == nil {
		t.Fatal("FormatterConfig should not be nil")
	}

	// Verificar configuração padrão
	if !formatter.config.UseEmojis {
		t.Error("Default config should use emojis")
	}

	if formatter.config.DetailLevel != "normal" {
		t.Error("Default detail level should be 'normal'")
	}

	if !formatter.config.IncludeFooter {
		t.Error("Default config should include footer")
	}

	// Testar com configuração customizada
	customConfig := &FormatterConfig{
		UseEmojis:     false,
		DetailLevel:   "minimal",
		IncludeFooter: false,
		CompactMode:   true,
	}

	formatter2 := NewResponseFormatter(customConfig)
	if formatter2.config.UseEmojis {
		t.Error("Custom config should not use emojis")
	}

	if formatter2.config.DetailLevel != "minimal" {
		t.Error("Custom detail level should be 'minimal'")
	}
}

// TestResponseTypeString testa a conversão de ResponseType para string
func TestResponseTypeString(t *testing.T) {
	tests := []struct {
		responseType ResponseType
		expected     string
	}{
		{PolicyResponse, "policy"},
		{DataResponse, "data"},
		{CalculationResponse, "calculation"},
		{ErrorResponse, "error"},
		{WhatIfResponse, "whatif"},
		{ResponseType(999), "unknown"},
	}

	for _, test := range tests {
		result := test.responseType.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

// TestFormatPolicyResponse testa formatação de resposta de política
func TestFormatPolicyResponse(t *testing.T) {
	formatter := NewResponseFormatter(nil)

	data := ResponseData{
		Question:   "Estagiários têm direito a VR?",
		Answer:     "Não, estagiários não têm direito ao Vale Refeição segundo a política VR_001.",
		Source:     "Manual de Recursos Humanos v2.1",
		Confidence: 0.95,
		Timestamp:  time.Now(),
	}

	result := formatter.Format(PolicyResponse, data)

	// Verificar elementos obrigatórios
	if !strings.Contains(result, "📋 Consulta de Política") {
		t.Error("Response should contain policy header")
	}

	if !strings.Contains(result, data.Question) {
		t.Error("Response should contain the question")
	}

	if !strings.Contains(result, data.Answer) {
		t.Error("Response should contain the answer")
	}

	if !strings.Contains(result, data.Source) {
		t.Error("Response should contain the source")
	}

	if !strings.Contains(result, "95% (Muito Alta)") {
		t.Error("Response should contain confidence level")
	}

	if !strings.Contains(result, "Esta resposta é baseada nas políticas oficiais") {
		t.Error("Response should contain footer")
	}
}

// TestFormatDataResponse testa formatação de resposta de dados
func TestFormatDataResponse(t *testing.T) {
	formatter := NewResponseFormatter(nil)

	data := ResponseData{
		Question:    "Quantos colaboradores foram processados?",
		Data:        "Total de 150 colaboradores processados",
		Stats:       "- Processados: 150\n- Pendentes: 5\n- Erro: 2",
		ProcessedAt: time.Now(),
		Timestamp:   time.Now(),
	}

	result := formatter.Format(DataResponse, data)

	// Verificar elementos obrigatórios
	if !strings.Contains(result, "📊 Consulta de Dados Processados") {
		t.Error("Response should contain data header")
	}

	if !strings.Contains(result, data.Question) {
		t.Error("Response should contain the question")
	}

	if !strings.Contains(result, data.Data) {
		t.Error("Response should contain the data")
	}

	if !strings.Contains(result, data.Stats) {
		t.Error("Response should contain the statistics")
	}

	if !strings.Contains(result, "Dados baseados no último processamento") {
		t.Error("Response should contain footer with timestamp")
	}
}

// TestFormatCalculationResponse testa formatação de resposta de cálculo
func TestFormatCalculationResponse(t *testing.T) {
	formatter := NewResponseFormatter(nil)

	data := ResponseData{
		Question:    "Calcular VR para colaborador admitido dia 20",
		Result:      "R$ 350,00",
		Rule:        "Aplicar regra de data quebrada: proporcional aos dias úteis",
		Calculation: "22 dias úteis - 15 dias trabalhados = R$ 450,00 * (15/22) = R$ 306,82",
		PolicyRef:   "Política VR_005 - Datas Quebradas",
		Timestamp:   time.Now(),
	}

	result := formatter.Format(CalculationResponse, data)

	// Verificar elementos obrigatórios
	if !strings.Contains(result, "🧮 Cálculo de VR") {
		t.Error("Response should contain calculation header")
	}

	if !strings.Contains(result, data.Question) {
		t.Error("Response should contain the question")
	}

	if !strings.Contains(result, data.Result) {
		t.Error("Response should contain the result")
	}

	if !strings.Contains(result, data.Rule) {
		t.Error("Response should contain the rule")
	}

	if !strings.Contains(result, data.Calculation) {
		t.Error("Response should contain the calculation")
	}

	if !strings.Contains(result, data.PolicyRef) {
		t.Error("Response should contain policy reference")
	}
}

// TestFormatErrorResponse testa formatação de resposta de erro
func TestFormatErrorResponse(t *testing.T) {
	formatter := NewResponseFormatter(nil)

	data := ResponseData{
		ErrorMessage: "Dados de colaborador não encontrados",
		Suggestions:  []string{"Verificar se a matrícula está correta", "Tentar novamente em alguns minutos"},
		Timestamp:    time.Now(),
	}

	result := formatter.Format(ErrorResponse, data)

	// Verificar elementos obrigatórios
	if !strings.Contains(result, "❌ Erro no Processamento") {
		t.Error("Response should contain error header")
	}

	if !strings.Contains(result, data.ErrorMessage) {
		t.Error("Response should contain error message")
	}

	if !strings.Contains(result, data.Suggestions[0]) {
		t.Error("Response should contain first suggestion")
	}

	if !strings.Contains(result, data.Suggestions[1]) {
		t.Error("Response should contain second suggestion")
	}

	if !strings.Contains(result, "Se o problema persistir") {
		t.Error("Response should contain footer")
	}
}

// TestFormatWhatIfResponse testa formatação de resposta hipotética
func TestFormatWhatIfResponse(t *testing.T) {
	formatter := NewResponseFormatter(nil)

	data := ResponseData{
		Question:    "E se o colaborador fosse admitido dia 10?",
		Answer:      "O valor seria R$ 450,00 (valor integral)",
		Calculation: "Admissão dia 10: 22 dias úteis = valor integral",
		Timestamp:   time.Now(),
	}

	result := formatter.Format(WhatIfResponse, data)

	// Verificar elementos obrigatórios
	if !strings.Contains(result, "🤔 Análise Hipotética") {
		t.Error("Response should contain whatif header")
	}

	if !strings.Contains(result, data.Question) {
		t.Error("Response should contain the question")
	}

	if !strings.Contains(result, data.Answer) {
		t.Error("Response should contain the answer")
	}

	if !strings.Contains(result, data.Calculation) {
		t.Error("Response should contain the calculation")
	}

	if !strings.Contains(result, "Esta é uma simulação") {
		t.Error("Response should contain simulation disclaimer")
	}
}

// TestFormatForContext testa formatação por contexto
func TestFormatForContext(t *testing.T) {
	formatter := NewResponseFormatter(nil)

	originalResponse := `## 📋 Consulta de Política

**Pergunta:** Teste

**Resposta:** Resposta de teste

💡 *Esta resposta é baseada nas políticas oficiais da empresa.*`

	// Testar modo compacto
	context := ResponseContext{
		OutputMode: "compact",
	}

	compactResult := formatter.FormatForContext(originalResponse, context)

	// Em modo compacto, deve ter menos linhas
	originalLines := len(strings.Split(originalResponse, "\n"))
	compactLines := len(strings.Split(compactResult, "\n"))

	if compactLines >= originalLines {
		t.Error("Compact format should have fewer lines")
	}
}

// TestGetConfidenceLevel testa conversão de confiança numérica para texto
func TestGetConfidenceLevel(t *testing.T) {
	formatter := NewResponseFormatter(nil)

	tests := []struct {
		confidence float64
		expected   string
	}{
		{0.95, "Muito Alta"},
		{0.85, "Alta"},
		{0.75, "Boa"},
		{0.65, "Moderada"},
		{0.35, "Baixa"},
	}

	for _, test := range tests {
		result := formatter.getConfidenceLevel(test.confidence)
		if result != test.expected {
			t.Errorf("For confidence %.2f, expected %s, got %s", test.confidence, test.expected, result)
		}
	}
}

// TestSetAndGetTemplate testa personalização de templates
func TestSetAndGetTemplate(t *testing.T) {
	formatter := NewResponseFormatter(nil)

	// Template customizado
	customTemplate := ResponseTemplate{
		Header: "## CUSTOM HEADER\n\n",
		Body:   "CUSTOM BODY: {{.Answer}}\n\n",
		Footer: "CUSTOM FOOTER\n",
	}

	// Definir template customizado
	formatter.SetTemplate(PolicyResponse, customTemplate)

	// Recuperar template
	retrieved, exists := formatter.GetTemplate(PolicyResponse)
	if !exists {
		t.Error("Template should exist")
	}

	if retrieved.Header != customTemplate.Header {
		t.Error("Custom template header should match")
	}

	// Testar formatação com template customizado
	data := ResponseData{
		Question:  "Teste",
		Answer:    "Resposta teste",
		Timestamp: time.Now(),
	}

	result := formatter.Format(PolicyResponse, data)

	if !strings.Contains(result, "CUSTOM HEADER") {
		t.Error("Response should use custom header")
	}

	if !strings.Contains(result, "CUSTOM BODY: Resposta teste") {
		t.Error("Response should use custom body")
	}

	if !strings.Contains(result, "CUSTOM FOOTER") {
		t.Error("Response should use custom footer")
	}
}

// TestUpdateConfig testa atualização de configuração
func TestUpdateConfig(t *testing.T) {
	formatter := NewResponseFormatter(nil)

	// Configuração inicial deve usar emojis
	if !formatter.config.UseEmojis {
		t.Error("Initial config should use emojis")
	}

	// Atualizar configuração
	newConfig := &FormatterConfig{
		UseEmojis:     false,
		DetailLevel:   "verbose",
		IncludeFooter: false,
		CompactMode:   true,
	}

	formatter.UpdateConfig(newConfig)

	// Verificar se configuração foi atualizada
	if formatter.config.UseEmojis {
		t.Error("Config should not use emojis after update")
	}

	if formatter.config.DetailLevel != "verbose" {
		t.Error("Detail level should be 'verbose' after update")
	}

	if formatter.config.IncludeFooter {
		t.Error("Config should not include footer after update")
	}
}

// TestRemoveEmojis testa remoção de emojis
func TestRemoveEmojis(t *testing.T) {
	formatter := NewResponseFormatter(nil)

	textWithEmojis := "📋 Consulta de Política\n🧮 Cálculo\n📊 Dados"
	result := formatter.removeEmojis(textWithEmojis)

	// Verificar se emojis foram removidos
	if strings.Contains(result, "📋") || strings.Contains(result, "🧮") || strings.Contains(result, "📊") {
		t.Error("Emojis should be removed from text")
	}

	// Verificar se texto restante está correto
	if !strings.Contains(result, "Consulta de Política") {
		t.Error("Text content should remain after emoji removal")
	}
}

// TestQuestionTypeToResponseType testa conversão de tipo de pergunta para tipo de resposta
func TestQuestionTypeToResponseType(t *testing.T) {
	tests := []struct {
		questionType QuestionType
		expected     ResponseType
	}{
		{PolicyQuestion, PolicyResponse},
		{ComplianceQuestion, PolicyResponse},
		{CalculationQuestion, CalculationResponse},
		{ProcessedDataQuestion, DataResponse},
		{WhatIfQuestion, WhatIfResponse},
		{UnknownQuestion, DataResponse}, // Default
	}

	for _, test := range tests {
		result := test.questionType.ToResponseType()
		if result != test.expected {
			t.Errorf("For question type %s, expected %s, got %s", 
				test.questionType.String(), test.expected.String(), result.String())
		}
	}
}

// BenchmarkFormat testa performance do formatador
func BenchmarkFormat(b *testing.B) {
	formatter := NewResponseFormatter(nil)

	data := ResponseData{
		Question:   "Teste de performance",
		Answer:     "Esta é uma resposta de teste para medir performance",
		Source:     "Fonte de teste",
		Confidence: 0.85,
		Timestamp:  time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatter.Format(PolicyResponse, data)
	}
}