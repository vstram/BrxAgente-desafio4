package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"BrxAgente-desafio4/internal/cache"
	"BrxAgente-desafio4/internal/knowledge"
	"github.com/tmc/langchaingo/tools"
)

// PolicyConsultantTool é a ferramenta que integra o consultor de políticas com LangChain
type PolicyConsultantTool struct {
	knowledgeBase   *knowledge.KnowledgeBaseManager
	policyEngine    *knowledge.PolicyEngine
	reasoningEngine *knowledge.ReasoningEngine
	citationManager *knowledge.CitationManager
	cache           *cache.KnowledgeCache
	initialized     bool
}

// NewPolicyConsultantTool cria uma nova instância da ferramenta
func NewPolicyConsultantTool(dataDir string) *PolicyConsultantTool {
	kb := knowledge.NewKnowledgeBaseManager()
	cm := knowledge.NewCitationManager()
	pe := knowledge.NewPolicyEngine(kb)
	re := knowledge.NewReasoningEngine(kb, pe, cm)

	// Criar cache com configuração padrão
	cacheConfig := cache.DefaultCacheConfig()
	knowledgeCache := cache.NewKnowledgeCache(cacheConfig)

	tool := &PolicyConsultantTool{
		knowledgeBase:   kb,
		policyEngine:    pe,
		reasoningEngine: re,
		citationManager: cm,
		cache:           knowledgeCache,
		initialized:     false,
	}

	// Carregar base de conhecimento
	if err := kb.LoadFromFiles(dataDir); err == nil {
		tool.initialized = true
	}

	return tool
}

// GetCache retorna o cache da ferramenta
func (pct *PolicyConsultantTool) GetCache() *cache.KnowledgeCache {
	return pct.cache
}

// EnableCache habilita o cache
func (pct *PolicyConsultantTool) EnableCache() {
	if pct.cache != nil {
		pct.cache.Enable()
	}
}

// DisableCache desabilita o cache
func (pct *PolicyConsultantTool) DisableCache() {
	if pct.cache != nil {
		pct.cache.Disable()
	}
}

// ClearCache limpa o cache
func (pct *PolicyConsultantTool) ClearCache() {
	if pct.cache != nil {
		pct.cache.Clear()
	}
}

// GetCacheMetrics retorna as métricas do cache
func (pct *PolicyConsultantTool) GetCacheMetrics() cache.KnowledgeCacheMetrics {
	if pct.cache != nil {
		return pct.cache.GetMetrics()
	}
	return cache.KnowledgeCacheMetrics{}
}

// Name retorna o nome da ferramenta
func (pct *PolicyConsultantTool) Name() string {
	return "policy_consultant"
}

// Description retorna a descrição da ferramenta
func (pct *PolicyConsultantTool) Description() string {
	return `Consultor inteligente de políticas de Vale Refeição. 

Esta ferramenta pode:
- Responder perguntas sobre políticas de VR e regulamentações
- Analisar cenários complexos com múltiplos fatores
- Fornecer citações precisas das fontes
- Executar raciocínio lógico sobre regras de negócio
- Resolver conflitos entre políticas
- Verificar compliance regulatório

Entrada: JSON com a consulta e parâmetros opcionais.
Exemplo: {"query": "Estagiários têm direito a VR?", "type": "simple"}
Exemplo: {"query": "Colaborador admitido dia 20 com afastamento", "type": "complex", "data_admissao": "2024-08-20", "afastamento": {"dias": 10}}

Tipos de consulta:
- "simple": Consulta direta sobre políticas
- "complex": Análise de cenários com múltiplos fatores  
- "whatif": Análise de cenários hipotéticos
- "compliance": Verificação de conformidade regulatória
- "conflict": Resolução de conflitos entre regras`
}

// Call implementa a interface tools.Tool do LangChain
func (pct *PolicyConsultantTool) Call(ctx context.Context, input string) (string, error) {
	return pct.Execute(ctx, input)
}

// Execute executa a ferramenta com o input fornecido, usando cache quando possível
func (pct *PolicyConsultantTool) Execute(ctx context.Context, input string) (string, error) {
	if !pct.initialized {
		return "", fmt.Errorf("ferramenta não foi inicializada corretamente - verifique se os arquivos de dados estão disponíveis")
	}

	// Parsear input JSON
	var requestData map[string]interface{}
	if err := json.Unmarshal([]byte(input), &requestData); err != nil {
		// Se não é JSON, tratar como consulta simples
		requestData = map[string]interface{}{
			"query": input,
			"type":  "simple",
		}
	}

	// Validar se tem query
	query, hasQuery := requestData["query"]
	if !hasQuery {
		return "", fmt.Errorf("campo 'query' é obrigatório")
	}

	queryStr, ok := query.(string)
	if !ok {
		return "", fmt.Errorf("campo 'query' deve ser uma string")
	}

	// Verificar cache primeiro
	if pct.cache != nil && pct.cache.IsEnabled() {
		if cached := pct.cache.Get(queryStr); cached != nil {
			return cached.Response, nil
		}
	}

	// Executar consulta normal
	result, err := pct.executeWithoutCache(ctx, input, requestData)
	if err != nil {
		return "", err
	}

	// Armazenar no cache se executou com sucesso
	if pct.cache != nil && pct.cache.IsEnabled() && result != "" {
		pct.cache.Set(queryStr, result, 0.9) // Confidence padrão de 90%
	}

	return result, nil
}

// executeWithoutCache executa a consulta sem usar cache
func (pct *PolicyConsultantTool) executeWithoutCache(ctx context.Context, input string, requestData map[string]interface{}) (string, error) {
	// Determinar tipo de consulta
	consultationType := pct.determineConsultationType(requestData)

	// Executar consulta baseada no tipo
	result, err := pct.executeConsultation(consultationType, requestData)
	if err != nil {
		return "", fmt.Errorf("erro durante consulta: %w", err)
	}

	// Formatar resposta
	response := pct.formatResponse(result, consultationType)

	return response, nil
}

// determineConsultationType determina o tipo de consulta
func (pct *PolicyConsultantTool) determineConsultationType(requestData map[string]interface{}) string {
	// Verificar se foi especificado explicitamente
	if typeValue, exists := requestData["type"]; exists {
		if typeStr, ok := typeValue.(string); ok {
			return typeStr
		}
	}

	// Determinar automaticamente baseado na query e parâmetros
	query, _ := requestData["query"].(string)
	queryLower := strings.ToLower(query)

	// Detectar "E se" ou cenários hipotéticos
	if strings.Contains(queryLower, "e se") || strings.Contains(queryLower, "what if") ||
		strings.Contains(queryLower, "supondo") || strings.Contains(queryLower, "caso") {
		return "whatif"
	}

	// Detectar verificação de compliance
	if strings.Contains(queryLower, "conforme") || strings.Contains(queryLower, "legal") ||
		strings.Contains(queryLower, "regulament") || strings.Contains(queryLower, "lei") {
		return "compliance"
	}

	// Detectar resolução de conflitos
	if strings.Contains(queryLower, "conflito") || strings.Contains(queryLower, "contradi") {
		return "conflict"
	}

	// Verificar se tem múltiplos parâmetros para análise complexa
	paramCount := 0
	complexParams := []string{"data_admissao", "data_desligamento", "afastamento", "ferias", "tipo_colaborador", "sindicato"}

	for _, param := range complexParams {
		if _, exists := requestData[param]; exists {
			paramCount++
		}
	}

	if paramCount >= 2 {
		return "complex"
	}

	return "simple"
}

// executeConsultation executa a consulta baseada no tipo
func (pct *PolicyConsultantTool) executeConsultation(consultationType string, requestData map[string]interface{}) (*knowledge.ConsultationResult, error) {
	switch consultationType {
	case "simple":
		return pct.executeSimpleQuery(requestData)
	case "complex":
		return pct.executeComplexAnalysis(requestData)
	case "whatif":
		return pct.executeWhatIfAnalysis(requestData)
	case "compliance":
		return pct.executeComplianceCheck(requestData)
	case "conflict":
		return pct.executeConflictResolution(requestData)
	default:
		return pct.executeSimpleQuery(requestData)
	}
}

// executeSimpleQuery executa uma consulta simples
func (pct *PolicyConsultantTool) executeSimpleQuery(requestData map[string]interface{}) (*knowledge.ConsultationResult, error) {
	query := requestData["query"].(string)

	// Buscar na base de conhecimento
	searchResults, err := pct.knowledgeBase.Search(query, 5)
	if err != nil {
		return nil, fmt.Errorf("erro na busca: %w", err)
	}

	if len(searchResults) == 0 {
		return &knowledge.ConsultationResult{
			Query:      query,
			Answer:     "Não foram encontradas informações relevantes sobre esta consulta.",
			Confidence: 0.1,
		}, nil
	}

	// Usar o primeiro resultado mais relevante
	bestResult := searchResults[0]

	// Criar citação
	citation := pct.citationManager.CreateCitation(bestResult.Item)

	// Formatar resposta baseada no conteúdo encontrado
	answer := pct.formatSimpleAnswer(bestResult, query)

	return &knowledge.ConsultationResult{
		Query:         query,
		Answer:        answer,
		Confidence:    bestResult.Score,
		Sources:       []knowledge.Citation{citation},
		RelatedTopics: bestResult.Item.Categories,
	}, nil
}

// executeComplexAnalysis executa análise complexa
func (pct *PolicyConsultantTool) executeComplexAnalysis(requestData map[string]interface{}) (*knowledge.ConsultationResult, error) {
	return pct.reasoningEngine.AnalyzeComplexScenario(requestData)
}

// executeWhatIfAnalysis executa análise de cenário hipotético
func (pct *PolicyConsultantTool) executeWhatIfAnalysis(requestData map[string]interface{}) (*knowledge.ConsultationResult, error) {
	return pct.reasoningEngine.AnalyzeComplexScenario(requestData)
}

// executeComplianceCheck executa verificação de compliance
func (pct *PolicyConsultantTool) executeComplianceCheck(requestData map[string]interface{}) (*knowledge.ConsultationResult, error) {
	return pct.reasoningEngine.AnalyzeComplexScenario(requestData)
}

// executeConflictResolution executa resolução de conflitos
func (pct *PolicyConsultantTool) executeConflictResolution(requestData map[string]interface{}) (*knowledge.ConsultationResult, error) {
	return pct.reasoningEngine.AnalyzeComplexScenario(requestData)
}

// formatSimpleAnswer formata resposta para consultas simples
func (pct *PolicyConsultantTool) formatSimpleAnswer(result knowledge.SearchResult, query string) string {
	var answer strings.Builder

	// Resposta principal baseada no conteúdo
	answer.WriteString(result.Item.Content)

	// Adicionar contexto se disponível
	if result.Context != "" {
		answer.WriteString(fmt.Sprintf("\n\n**Contexto:** %s", result.Context))
	}

	// Adicionar fonte
	answer.WriteString(fmt.Sprintf("\n\n**Fonte:** %s", result.Item.Source))

	// Adicionar destaques se disponíveis
	if len(result.Highlights) > 0 {
		answer.WriteString(fmt.Sprintf("\n\n**Trechos relevantes:** %s", strings.Join(result.Highlights, "; ")))
	}

	return answer.String()
}

// formatResponse formata a resposta final
func (pct *PolicyConsultantTool) formatResponse(result *knowledge.ConsultationResult, consultationType string) string {
	var response strings.Builder

	// Cabeçalho com tipo de consulta
	switch consultationType {
	case "simple":
		response.WriteString("## 📚 Consulta de Política\n\n")
	case "complex":
		response.WriteString("## 🧮 Análise Complexa\n\n")
	case "whatif":
		response.WriteString("## 🤔 Análise Hipotética\n\n")
	case "compliance":
		response.WriteString("## ⚖️ Verificação de Compliance\n\n")
	case "conflict":
		response.WriteString("## ⚡ Resolução de Conflitos\n\n")
	}

	// Pergunta original
	response.WriteString(fmt.Sprintf("**Pergunta:** %s\n\n", result.Query))

	// Resposta principal
	response.WriteString(fmt.Sprintf("**Resposta:**\n%s\n\n", result.Answer))

	// Nível de confiança
	confidenceLevel := pct.getConfidenceLevel(result.Confidence)
	response.WriteString(fmt.Sprintf("**Confiança:** %.0f%% (%s)\n\n", result.Confidence*100, confidenceLevel))

	// Passos de raciocínio (se disponível)
	if len(result.ReasoningSteps) > 0 {
		response.WriteString("**Raciocínio:**\n")
		for _, step := range result.ReasoningSteps {
			response.WriteString(fmt.Sprintf("%d. %s: %s\n", step.Step, step.Description, step.Result))
		}
		response.WriteString("\n")
	}

	// Fontes citadas
	if len(result.Sources) > 0 {
		response.WriteString("**Fontes:**\n")
		citationList := pct.citationManager.FormatCitationList(result.Sources)
		response.WriteString(citationList)
		response.WriteString("\n\n")
	}

	// Ambiguidades (se houver)
	if len(result.Ambiguities) > 0 {
		response.WriteString("**⚠️ Ambiguidades identificadas:**\n")
		for _, ambiguity := range result.Ambiguities {
			response.WriteString(fmt.Sprintf("- %s\n", ambiguity))
		}
		response.WriteString("\n")
	}

	// Recomendações (se houver)
	if len(result.Recommendations) > 0 {
		response.WriteString("**💡 Recomendações:**\n")
		for _, recommendation := range result.Recommendations {
			response.WriteString(fmt.Sprintf("- %s\n", recommendation))
		}
		response.WriteString("\n")
	}

	// Tópicos relacionados
	if len(result.RelatedTopics) > 0 {
		response.WriteString(fmt.Sprintf("**Tópicos relacionados:** %s\n\n", strings.Join(result.RelatedTopics, ", ")))
	}

	// Tempo de processamento
	if result.ProcessingTime > 0 {
		response.WriteString(fmt.Sprintf("*Processado em %v*", result.ProcessingTime))
	}

	return response.String()
}

// getConfidenceLevel retorna o nível de confiança em texto
func (pct *PolicyConsultantTool) getConfidenceLevel(confidence float64) string {
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

// GetKnowledgeStats retorna estatísticas da base de conhecimento
func (pct *PolicyConsultantTool) GetKnowledgeStats() map[string]interface{} {
	return pct.knowledgeBase.GetStats()
}

// SearchKnowledge busca diretamente na base de conhecimento
func (pct *PolicyConsultantTool) SearchKnowledge(query string, limit int) ([]knowledge.SearchResult, error) {
	return pct.knowledgeBase.Search(query, limit)
}

// GetByCategory busca itens por categoria
func (pct *PolicyConsultantTool) GetByCategory(category string) ([]knowledge.KnowledgeItem, error) {
	return pct.knowledgeBase.GetByCategory(category)
}

// ReloadKnowledgeBase recarrega a base de conhecimento
func (pct *PolicyConsultantTool) ReloadKnowledgeBase(dataDir string) error {
	err := pct.knowledgeBase.LoadFromFiles(dataDir)
	if err != nil {
		pct.initialized = false
		return err
	}

	pct.initialized = true
	return nil
}

// AvailableCategories retorna categorias disponíveis na base de conhecimento
func (pct *PolicyConsultantTool) AvailableCategories() []string {
	// Esta é uma implementação básica - idealmente deveria extrair categorias reais
	return []string{
		"elegibilidade", "valores", "calculo", "admissao", "desligamento",
		"afastamento", "ferias", "feriados", "tipos-colaborador", "validacao",
		"clt", "pat", "regulamentacoes", "compliance",
	}
}

// ValidateInput valida se o input está no formato correto
func (pct *PolicyConsultantTool) ValidateInput(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input não pode estar vazio")
	}

	// Tentar parsear como JSON
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		// Se não é JSON válido, assumir que é texto simples (válido)
		return nil
	}

	// Se é JSON, validar se tem campo query
	if _, hasQuery := data["query"]; !hasQuery {
		return fmt.Errorf("JSON deve conter o campo 'query'")
	}

	return nil
}

// Implementar interface tools.Tool do LangChain
var _ tools.Tool = (*PolicyConsultantTool)(nil)

// CreatePolicyConsultantLangChainTool cria uma ferramenta compatível com LangChain
func CreatePolicyConsultantLangChainTool(dataDir string) tools.Tool {
	// Garantir que o caminho para os dados esteja correto
	if dataDir == "" {
		dataDir = filepath.Join("internal", "data", "policies")
	}

	pct := NewPolicyConsultantTool(dataDir)

	// Retornar como tools.Tool
	return pct
}

// Métodos adicionais para integração

// GetCapabilities retorna as capacidades da ferramenta
func (pct *PolicyConsultantTool) GetCapabilities() map[string]interface{} {
	return map[string]interface{}{
		"consultation_types": []string{"simple", "complex", "whatif", "compliance", "conflict"},
		"supported_formats":  []string{"text", "json"},
		"features": []string{
			"policy_search", "complex_reasoning", "citation_management",
			"conflict_resolution", "compliance_checking", "scenario_analysis",
		},
		"initialized":     pct.initialized,
		"knowledge_stats": pct.GetKnowledgeStats(),
	}
}

// GetExamples retorna exemplos de uso da ferramenta
func (pct *PolicyConsultantTool) GetExamples() []map[string]string {
	return []map[string]string{
		{
			"type":        "simple",
			"description": "Consulta simples sobre elegibilidade",
			"input":       `{"query": "Estagiários têm direito a Vale Refeição?", "type": "simple"}`,
		},
		{
			"type":        "complex",
			"description": "Análise de cenário complexo com múltiplos fatores",
			"input":       `{"query": "Calcular VR para colaborador admitido dia 20 com afastamento", "type": "complex", "data_admissao": "2024-08-20", "afastamento": {"dias": 10}, "tipo_colaborador": "efetivo"}`,
		},
		{
			"type":        "whatif",
			"description": "Análise de cenário hipotético",
			"input":       `{"query": "E se o colaborador fosse admitido dia 10 ao invés de dia 20?", "type": "whatif", "data_admissao": "2024-08-10"}`,
		},
		{
			"type":        "compliance",
			"description": "Verificação de conformidade regulatória",
			"input":       `{"query": "O cálculo está conforme com a CLT?", "type": "compliance", "valor_vr": 25.50}`,
		},
		{
			"type":        "text",
			"description": "Consulta em texto simples",
			"input":       "Como calcular VR para datas quebradas?",
		},
	}
}
