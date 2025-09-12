package agent

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

// ResponseType representa os diferentes tipos de resposta
type ResponseType int

const (
	PolicyResponse ResponseType = iota
	DataResponse
	CalculationResponse
	ErrorResponse
	WhatIfResponse
)

// String implementa fmt.Stringer para ResponseType
func (rt ResponseType) String() string {
	switch rt {
	case PolicyResponse:
		return "policy"
	case DataResponse:
		return "data"
	case CalculationResponse:
		return "calculation"
	case ErrorResponse:
		return "error"
	case WhatIfResponse:
		return "whatif"
	default:
		return "unknown"
	}
}

// ResponseData contém os dados para formatação de resposta
type ResponseData struct {
	Question      string                 `json:"question"`
	Answer        string                 `json:"answer"`
	Source        string                 `json:"source,omitempty"`
	Confidence    float64                `json:"confidence,omitempty"`
	Data          string                 `json:"data,omitempty"`
	Stats         string                 `json:"stats,omitempty"`
	ProcessedAt   time.Time              `json:"processed_at,omitempty"`
	Result        string                 `json:"result,omitempty"`
	Rule          string                 `json:"rule,omitempty"`
	Calculation   string                 `json:"calculation,omitempty"`
	PolicyRef     string                 `json:"policy_reference,omitempty"`
	ErrorMessage  string                 `json:"error_message,omitempty"`
	Suggestions   []string               `json:"suggestions,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
}

// ResponseTemplate define a estrutura de um template de resposta
type ResponseTemplate struct {
	Header   string
	Body     string
	Footer   string
	Metadata bool
}

// ResponseContext fornece contexto adicional para formatação
type ResponseContext struct {
	OutputMode   string // compact, detailed, technical
	UseEmojis    bool
	DetailLevel  string // minimal, normal, verbose
	IncludeFooter bool
	CompactMode  bool
}

// FormatterConfig configurações do formatador
type FormatterConfig struct {
	UseEmojis     bool   `json:"use_emojis"`
	DetailLevel   string `json:"detail_level"`
	IncludeFooter bool   `json:"include_footer"`
	CompactMode   bool   `json:"compact_mode"`
}

// DefaultFormatterConfig retorna configuração padrão
func DefaultFormatterConfig() *FormatterConfig {
	return &FormatterConfig{
		UseEmojis:     true,
		DetailLevel:   "normal",
		IncludeFooter: true,
		CompactMode:   false,
	}
}

// ResponseFormatter sistema de formatação de respostas
type ResponseFormatter struct {
	templates map[ResponseType]ResponseTemplate
	config    *FormatterConfig
}

// NewResponseFormatter cria uma nova instância do formatador
func NewResponseFormatter(config *FormatterConfig) *ResponseFormatter {
	if config == nil {
		config = DefaultFormatterConfig()
	}

	rf := &ResponseFormatter{
		config:    config,
		templates: make(map[ResponseType]ResponseTemplate),
	}

	// Inicializar templates padrão
	rf.initializeDefaultTemplates()

	return rf
}

// initializeDefaultTemplates inicializa os templates padrão
func (rf *ResponseFormatter) initializeDefaultTemplates() {
	// Template para respostas de política
	rf.templates[PolicyResponse] = ResponseTemplate{
		Header: "## 📋 Consulta de Política\n\n**Pergunta:** {{.Question}}\n\n",
		Body: "**Resposta:**\n{{.Answer}}\n\n{{if .Source}}**Fonte:** {{.Source}}\n{{end}}{{if .Confidence}}**Confiança:** {{printf \"%.0f\" (mul .Confidence 100)}}% ({{.ConfidenceLevel}})\n{{end}}\n",
		Footer: "💡 *Esta resposta é baseada nas políticas oficiais da empresa.*\n",
		Metadata: true,
	}

	// Template para respostas de dados
	rf.templates[DataResponse] = ResponseTemplate{
		Header: "## 📊 Consulta de Dados Processados\n\n**Pergunta:** {{.Question}}\n\n",
		Body: "**Resultado:**\n{{.Data}}\n\n{{if .Stats}}**Estatísticas:**\n{{.Stats}}\n\n{{end}}",
		Footer: "*Dados baseados no último processamento{{if .ProcessedAt}}: {{.ProcessedAt.Format \"02/01/2006 15:04\"}}{{end}}*\n",
		Metadata: true,
	}

	// Template para respostas de cálculo
	rf.templates[CalculationResponse] = ResponseTemplate{
		Header: "## 🧮 Cálculo de VR\n\n**Cenário:** {{.Question}}\n\n",
		Body: "**Resultado:** {{.Result}}\n\n{{if .Rule}}**Aplicação da Regra:**\n{{.Rule}}\n\n{{end}}{{if .Calculation}}**Cálculo:**\n{{.Calculation}}\n\n{{end}}",
		Footer: "{{if .PolicyRef}}📖 **Política aplicada:** {{.PolicyRef}}\n{{end}}",
		Metadata: true,
	}

	// Template para respostas de erro
	rf.templates[ErrorResponse] = ResponseTemplate{
		Header: "## ❌ Erro no Processamento\n\n",
		Body: "**Erro:** {{.ErrorMessage}}\n\n{{if .Suggestions}}**Sugestões:**\n{{range .Suggestions}}- {{.}}\n{{end}}\n{{end}}",
		Footer: "💭 *Se o problema persistir, verifique os dados de entrada ou entre em contato com o suporte.*\n",
		Metadata: false,
	}

	// Template para cenários hipotéticos
	rf.templates[WhatIfResponse] = ResponseTemplate{
		Header: "## 🤔 Análise Hipotética\n\n**Cenário:** {{.Question}}\n\n",
		Body: "**Resultado da Simulação:**\n{{.Answer}}\n\n{{if .Calculation}}**Detalhes do Cálculo:**\n{{.Calculation}}\n\n{{end}}",
		Footer: "🔮 *Esta é uma simulação baseada nas políticas atuais. Resultados reais podem variar.*\n",
		Metadata: true,
	}
}

// Format formata uma resposta usando o template apropriado
func (rf *ResponseFormatter) Format(responseType ResponseType, data ResponseData) string {
	template, exists := rf.templates[responseType]
	if !exists {
		return rf.formatPlainText(data)
	}

	var result strings.Builder

	// Processar header
	if template.Header != "" {
		result.WriteString(rf.processTemplate(template.Header, data))
	}

	// Processar body com formatação específica
	if template.Body != "" {
		body := rf.processTemplate(template.Body, data)
		result.WriteString(rf.enhanceFormatting(body, responseType))
	}

	// Processar footer se habilitado
	if template.Footer != "" && rf.config.IncludeFooter {
		result.WriteString(rf.processTemplate(template.Footer, data))
	}

	return result.String()
}

// FormatForContext aplica formatação específica para o contexto
func (rf *ResponseFormatter) FormatForContext(response string, context ResponseContext) string {
	switch context.OutputMode {
	case "compact":
		return rf.compactFormat(response)
	case "detailed":
		return rf.detailedFormat(response)
	case "technical":
		return rf.technicalFormat(response)
	default:
		return response
	}
}

// processTemplate processa um template string com os dados fornecidos
func (rf *ResponseFormatter) processTemplate(templateStr string, data ResponseData) string {
	// Adicionar funções helper
	funcMap := template.FuncMap{
		"mul": func(a, b float64) float64 { return a * b },
		"add": func(a, b int) int { return a + b },
		"formatTime": func(t time.Time) string { return t.Format("02/01/2006 15:04") },
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
	}

	// Criar template
	tmpl, err := template.New("response").Funcs(funcMap).Parse(templateStr)
	if err != nil {
		return templateStr // Retornar template original em caso de erro
	}

	// Preparar dados com campos auxiliares
	templateData := rf.prepareTemplateData(data)

	// Executar template
	var buf strings.Builder
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return templateStr // Retornar template original em caso de erro
	}

	return buf.String()
}

// prepareTemplateData prepara os dados para o template, adicionando campos auxiliares
func (rf *ResponseFormatter) prepareTemplateData(data ResponseData) map[string]interface{} {
	templateData := map[string]interface{}{
		"Question":      data.Question,
		"Answer":        data.Answer,
		"Source":        data.Source,
		"Confidence":    data.Confidence,
		"Data":          data.Data,
		"Stats":         data.Stats,
		"ProcessedAt":   data.ProcessedAt,
		"Result":        data.Result,
		"Rule":          data.Rule,
		"Calculation":   data.Calculation,
		"PolicyRef":     data.PolicyRef,
		"ErrorMessage":  data.ErrorMessage,
		"Suggestions":   data.Suggestions,
		"Metadata":      data.Metadata,
		"Timestamp":     data.Timestamp,
	}

	// Adicionar campos auxiliares
	templateData["ConfidenceLevel"] = rf.getConfidenceLevel(data.Confidence)
	templateData["FormattedTimestamp"] = data.Timestamp.Format("02/01/2006 15:04")
	
	if !data.ProcessedAt.IsZero() {
		templateData["FormattedProcessedAt"] = data.ProcessedAt.Format("02/01/2006 15:04")
	}

	return templateData
}

// getConfidenceLevel converte confiança numérica em texto
func (rf *ResponseFormatter) getConfidenceLevel(confidence float64) string {
	if confidence >= 0.9 {
		return "Muito Alta"
	} else if confidence >= 0.8 {
		return "Alta"
	} else if confidence >= 0.7 {
		return "Boa"
	} else if confidence >= 0.5 {
		return "Moderada"
	} else {
		return "Baixa"
	}
}

// enhanceFormatting aplica melhorias de formatação específicas por tipo
func (rf *ResponseFormatter) enhanceFormatting(body string, responseType ResponseType) string {
	switch responseType {
	case CalculationResponse:
		return rf.enhanceCalculationFormatting(body)
	case DataResponse:
		return rf.enhanceDataFormatting(body)
	case PolicyResponse:
		return rf.enhancePolicyFormatting(body)
	default:
		return body
	}
}

// enhanceCalculationFormatting melhora formatação de cálculos
func (rf *ResponseFormatter) enhanceCalculationFormatting(text string) string {
	// Destacar valores monetários
	text = rf.highlightMonetaryValues(text)
	
	// Formatar listas numeradas para passos de cálculo
	text = rf.formatCalculationSteps(text)
	
	return text
}

// enhanceDataFormatting melhora formatação de dados
func (rf *ResponseFormatter) enhanceDataFormatting(text string) string {
	// Destacar números importantes
	text = rf.highlightNumbers(text)
	
	// Formatar percentuais
	text = rf.formatPercentages(text)
	
	return text
}

// enhancePolicyFormatting melhora formatação de políticas
func (rf *ResponseFormatter) enhancePolicyFormatting(text string) string {
	// Destacar referências de políticas
	text = rf.highlightPolicyReferences(text)
	
	return text
}

// highlightMonetaryValues destaca valores monetários no texto
func (rf *ResponseFormatter) highlightMonetaryValues(text string) string {
	// Não aplicar destaque automático - pode quebrar a formatação
	// Deixar para o template decidir a formatação
	return text
}

// formatCalculationSteps formata passos de cálculo como lista numerada
func (rf *ResponseFormatter) formatCalculationSteps(text string) string {
	// Por enquanto, manter formatação original sem modificações automáticas
	// Evitar adicionar numeração que pode confundir
	return text
}

// highlightNumbers destaca números importantes
func (rf *ResponseFormatter) highlightNumbers(text string) string {
	// Implementação simples - pode ser expandida com regex mais sofisticadas
	return text
}

// formatPercentages formata percentuais
func (rf *ResponseFormatter) formatPercentages(text string) string {
	// Destacar percentuais
	return strings.ReplaceAll(text, "%", "**%**")
}

// highlightPolicyReferences destaca referências de políticas
func (rf *ResponseFormatter) highlightPolicyReferences(text string) string {
	// Destacar códigos de políticas (ex: VR_003)
	return text
}

// compactFormat aplica formatação compacta
func (rf *ResponseFormatter) compactFormat(text string) string {
	// Remove emojis se configurado
	if !rf.config.UseEmojis {
		text = rf.removeEmojis(text)
	}
	
	// Remove linhas extras
	text = rf.removeExtraLines(text)
	
	// Remove seções opcionais
	text = rf.removeOptionalSections(text)
	
	return text
}

// detailedFormat aplica formatação detalhada
func (rf *ResponseFormatter) detailedFormat(text string) string {
	// Adiciona seções extras de contexto
	return text
}

// technicalFormat aplica formatação técnica
func (rf *ResponseFormatter) technicalFormat(text string) string {
	// Adiciona metadados técnicos
	return text
}

// formatPlainText formata texto simples quando não há template específico
func (rf *ResponseFormatter) formatPlainText(data ResponseData) string {
	var result strings.Builder
	
	if data.Question != "" {
		result.WriteString(fmt.Sprintf("**Pergunta:** %s\n\n", data.Question))
	}
	
	if data.Answer != "" {
		result.WriteString(fmt.Sprintf("**Resposta:** %s\n\n", data.Answer))
	}
	
	return result.String()
}

// removeEmojis remove emojis do texto
func (rf *ResponseFormatter) removeEmojis(text string) string {
	// Lista básica de emojis comuns usados nos templates
	emojis := []string{"📋", "📊", "🧮", "❌", "🤔", "💡", "📖", "💭", "🔮"}
	
	for _, emoji := range emojis {
		text = strings.ReplaceAll(text, emoji+" ", "")
		text = strings.ReplaceAll(text, emoji, "")
	}
	
	return text
}

// removeExtraLines remove linhas extras para formato compacto
func (rf *ResponseFormatter) removeExtraLines(text string) string {
	// Remove múltiplas linhas vazias consecutivas
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}
	
	return text
}

// removeOptionalSections remove seções opcionais em modo compacto
func (rf *ResponseFormatter) removeOptionalSections(text string) string {
	// Remove seções iniciadas com asterisco (geralmente metadados)
	lines := strings.Split(text, "\n")
	var result []string
	
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "*") {
			result = append(result, line)
		}
	}
	
	return strings.Join(result, "\n")
}

// SetTemplate permite personalizar um template específico
func (rf *ResponseFormatter) SetTemplate(responseType ResponseType, template ResponseTemplate) {
	rf.templates[responseType] = template
}

// GetTemplate retorna o template para um tipo específico
func (rf *ResponseFormatter) GetTemplate(responseType ResponseType) (ResponseTemplate, bool) {
	template, exists := rf.templates[responseType]
	return template, exists
}

// UpdateConfig atualiza a configuração do formatador
func (rf *ResponseFormatter) UpdateConfig(config *FormatterConfig) {
	if config != nil {
		rf.config = config
	}
}

// GetConfig retorna a configuração atual
func (rf *ResponseFormatter) GetConfig() *FormatterConfig {
	return rf.config
}