package knowledge

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// CitationManager gerencia citações e referências para as consultas
type CitationManager struct {
	citations map[string]Citation // Cache de citações por ID
}

// NewCitationManager cria uma nova instância do gerenciador de citações
func NewCitationManager() *CitationManager {
	return &CitationManager{
		citations: make(map[string]Citation),
	}
}

// CreateCitation cria uma citação a partir de um item de conhecimento
func (cm *CitationManager) CreateCitation(item KnowledgeItem) Citation {
	citation := Citation{
		Source:      item.Source,
		Date:        item.EffectiveDate,
		Reliability: cm.assessReliability(item),
	}

	// Adicionar informações específicas baseadas no tipo
	switch item.Type {
	case "regulation":
		citation.Authority = cm.extractAuthority(item)
		citation.URL = cm.generateRegulatoryURL(item)
	case "policy":
		citation.Section = cm.extractPolicySection(item)
	case "business_rule":
		citation.Section = fmt.Sprintf("Regra %s", item.ID)
	}

	// Armazenar no cache
	cm.citations[item.ID] = citation

	return citation
}

// assessReliability avalia a confiabilidade de uma fonte
func (cm *CitationManager) assessReliability(item KnowledgeItem) string {
	// Regulamentações governamentais têm alta confiabilidade
	if item.Type == "regulation" {
		return "high"
	}

	// Políticas internas recentes têm alta confiabilidade
	if item.Type == "policy" && time.Since(item.EffectiveDate) < 365*24*time.Hour {
		return "high"
	}

	// Políticas antigas têm confiabilidade média
	if item.Type == "policy" && time.Since(item.EffectiveDate) < 2*365*24*time.Hour {
		return "medium"
	}

	// Regras de negócio têm confiabilidade média por padrão
	if item.Type == "business_rule" {
		return "medium"
	}

	// Fontes muito antigas têm baixa confiabilidade
	if time.Since(item.EffectiveDate) > 3*365*24*time.Hour {
		return "low"
	}

	return "medium"
}

// extractAuthority extrai a autoridade responsável pela regulamentação
func (cm *CitationManager) extractAuthority(item KnowledgeItem) string {
	source := strings.ToLower(item.Source)
	
	if strings.Contains(source, "clt") {
		return "Consolidação das Leis do Trabalho"
	}
	
	if strings.Contains(source, "receita federal") {
		return "Receita Federal do Brasil"
	}
	
	if strings.Contains(source, "ministério") {
		if strings.Contains(source, "trabalho") {
			return "Ministério do Trabalho e Emprego"
		}
		return "Ministério Federal"
	}
	
	if strings.Contains(source, "lei") {
		return "Governo Federal"
	}
	
	// Extrair autoridade dos metadados se disponível
	if authority, exists := item.Metadata["authority"]; exists {
		return authority
	}
	
	return "Autoridade não especificada"
}

// generateRegulatoryURL gera URL para regulamentações quando possível
func (cm *CitationManager) generateRegulatoryURL(item KnowledgeItem) string {
	source := strings.ToLower(item.Source)
	
	// URLs específicas para diferentes tipos de regulamentação
	if strings.Contains(source, "clt") {
		return "http://www.planalto.gov.br/ccivil_03/decreto-lei/del5452.htm"
	}
	
	if strings.Contains(source, "lei 6.321") {
		return "http://www.planalto.gov.br/ccivil_03/leis/l6321.htm"
	}
	
	// URL genérica para legislação federal
	if strings.Contains(source, "lei") || strings.Contains(source, "decreto") {
		return "http://www.planalto.gov.br/ccivil_03/"
	}
	
	return ""
}

// extractPolicySection extrai a seção da política
func (cm *CitationManager) extractPolicySection(item KnowledgeItem) string {
	// Tentar extrair seção do título ou ID
	if strings.Contains(item.Title, "Elegibilidade") {
		return "Seção 1 - Elegibilidade"
	}
	
	if strings.Contains(item.Title, "Valor") {
		return "Seção 2 - Valores e Cálculos"
	}
	
	if strings.Contains(item.Title, "Exclusão") {
		return "Seção 3 - Exclusões"
	}
	
	if strings.Contains(item.Title, "Data") || strings.Contains(item.Title, "Quebrada") {
		return "Seção 4 - Datas Quebradas"
	}
	
	if strings.Contains(item.Title, "Afastamento") {
		return "Seção 5 - Afastamentos"
	}
	
	if strings.Contains(item.Title, "Férias") {
		return "Seção 6 - Férias"
	}
	
	if strings.Contains(item.Title, "Feriado") {
		return "Seção 7 - Feriados"
	}
	
	// Seção baseada no ID
	if len(item.ID) >= 3 {
		return fmt.Sprintf("Item %s", strings.ToUpper(item.ID))
	}
	
	return "Seção geral"
}

// FormatCitation formata uma citação em texto legível
func (cm *CitationManager) FormatCitation(citation Citation) string {
	var parts []string
	
	// Fonte principal
	parts = append(parts, citation.Source)
	
	// Seção se disponível
	if citation.Section != "" {
		parts = append(parts, citation.Section)
	}
	
	// Autoridade se disponível
	if citation.Authority != "" && citation.Authority != "Autoridade não especificada" {
		parts = append(parts, fmt.Sprintf("(%s)", citation.Authority))
	}
	
	// Data
	if !citation.Date.IsZero() {
		year := citation.Date.Year()
		if year > 1900 {
			parts = append(parts, fmt.Sprintf("%d", year))
		}
	}
	
	formatted := strings.Join(parts, ", ")
	
	// Adicionar indicador de confiabilidade
	switch citation.Reliability {
	case "high":
		formatted += " [Alta Confiabilidade]"
	case "low":
		formatted += " [Baixa Confiabilidade]"
	}
	
	return formatted
}

// FormatCitationList formata uma lista de citações
func (cm *CitationManager) FormatCitationList(citations []Citation) string {
	if len(citations) == 0 {
		return "Nenhuma fonte citada."
	}
	
	// Remover duplicatas
	citations = cm.deduplicateCitations(citations)
	
	// Ordenar por confiabilidade e data
	sort.Slice(citations, func(i, j int) bool {
		// Priorizar por confiabilidade
		reliabilityOrder := map[string]int{"high": 1, "medium": 2, "low": 3}
		if reliabilityOrder[citations[i].Reliability] != reliabilityOrder[citations[j].Reliability] {
			return reliabilityOrder[citations[i].Reliability] < reliabilityOrder[citations[j].Reliability]
		}
		// Depois por data (mais recente primeiro)
		return citations[i].Date.After(citations[j].Date)
	})
	
	var formattedList []string
	for i, citation := range citations {
		formatted := fmt.Sprintf("%d. %s", i+1, cm.FormatCitation(citation))
		formattedList = append(formattedList, formatted)
	}
	
	return strings.Join(formattedList, "\n")
}

// deduplicateCitations remove citações duplicadas
func (cm *CitationManager) deduplicateCitations(citations []Citation) []Citation {
	seen := make(map[string]bool)
	var result []Citation
	
	for _, citation := range citations {
		// Criar chave única baseada em fonte e data
		key := citation.Source + citation.Date.String()
		if !seen[key] {
			seen[key] = true
			result = append(result, citation)
		}
	}
	
	return result
}

// ValidateCitation valida se uma citação está completa e correta
func (cm *CitationManager) ValidateCitation(citation Citation) []string {
	var issues []string
	
	// Verificar campos obrigatórios
	if citation.Source == "" {
		issues = append(issues, "Fonte não especificada")
	}
	
	if citation.Date.IsZero() {
		issues = append(issues, "Data não especificada")
	}
	
	if citation.Reliability == "" {
		issues = append(issues, "Confiabilidade não avaliada")
	}
	
	// Verificar validade da data
	if !citation.Date.IsZero() && citation.Date.After(time.Now()) {
		issues = append(issues, "Data no futuro")
	}
	
	// Verificar se a data não é muito antiga para certas fontes
	if !citation.Date.IsZero() && time.Since(citation.Date) > 10*365*24*time.Hour {
		if !strings.Contains(strings.ToLower(citation.Source), "clt") {
			issues = append(issues, "Fonte potencialmente desatualizada (>10 anos)")
		}
	}
	
	return issues
}

// GetCitationStats retorna estatísticas das citações
func (cm *CitationManager) GetCitationStats(citations []Citation) map[string]interface{} {
	stats := make(map[string]interface{})
	
	if len(citations) == 0 {
		stats["total"] = 0
		return stats
	}
	
	// Contar por confiabilidade
	reliabilityCount := make(map[string]int)
	authorityCount := make(map[string]int)
	typeCount := make(map[string]int)
	
	for _, citation := range citations {
		reliabilityCount[citation.Reliability]++
		
		if citation.Authority != "" {
			authorityCount[citation.Authority]++
		}
		
		// Determinar tipo baseado na fonte
		citationType := cm.determineCitationType(citation)
		typeCount[citationType]++
	}
	
	stats["total"] = len(citations)
	stats["by_reliability"] = reliabilityCount
	stats["by_authority"] = authorityCount
	stats["by_type"] = typeCount
	
	// Calcular idade média das fontes
	totalAge := 0.0
	validDates := 0
	for _, citation := range citations {
		if !citation.Date.IsZero() {
			age := time.Since(citation.Date).Hours() / (24 * 365) // Anos
			totalAge += age
			validDates++
		}
	}
	
	if validDates > 0 {
		stats["average_age_years"] = totalAge / float64(validDates)
	}
	
	return stats
}

// determineCitationType determina o tipo de citação baseado na fonte
func (cm *CitationManager) determineCitationType(citation Citation) string {
	source := strings.ToLower(citation.Source)
	
	if strings.Contains(source, "manual") || strings.Contains(source, "política") {
		return "policy"
	}
	
	if strings.Contains(source, "clt") || strings.Contains(source, "lei") || 
	   strings.Contains(source, "decreto") || strings.Contains(source, "portaria") {
		return "regulation"
	}
	
	if strings.Contains(source, "sistema") || strings.Contains(source, "brxagente") {
		return "business_rule"
	}
	
	return "other"
}

// CreateBibliography cria uma bibliografia formatada
func (cm *CitationManager) CreateBibliography(citations []Citation) string {
	if len(citations) == 0 {
		return "# Bibliografia\n\nNenhuma fonte citada."
	}
	
	// Organizar por tipo
	policyRefs := []Citation{}
	regulationRefs := []Citation{}
	ruleRefs := []Citation{}
	otherRefs := []Citation{}
	
	for _, citation := range citations {
		switch cm.determineCitationType(citation) {
		case "policy":
			policyRefs = append(policyRefs, citation)
		case "regulation":
			regulationRefs = append(regulationRefs, citation)
		case "business_rule":
			ruleRefs = append(ruleRefs, citation)
		default:
			otherRefs = append(otherRefs, citation)
		}
	}
	
	var bibliography strings.Builder
	bibliography.WriteString("# Bibliografia\n\n")
	
	// Seção de Regulamentações
	if len(regulationRefs) > 0 {
		bibliography.WriteString("## Regulamentações e Leis\n\n")
		for _, citation := range regulationRefs {
			bibliography.WriteString(fmt.Sprintf("- %s\n", cm.FormatCitation(citation)))
		}
		bibliography.WriteString("\n")
	}
	
	// Seção de Políticas
	if len(policyRefs) > 0 {
		bibliography.WriteString("## Políticas Internas\n\n")
		for _, citation := range policyRefs {
			bibliography.WriteString(fmt.Sprintf("- %s\n", cm.FormatCitation(citation)))
		}
		bibliography.WriteString("\n")
	}
	
	// Seção de Regras de Negócio
	if len(ruleRefs) > 0 {
		bibliography.WriteString("## Regras de Negócio\n\n")
		for _, citation := range ruleRefs {
			bibliography.WriteString(fmt.Sprintf("- %s\n", cm.FormatCitation(citation)))
		}
		bibliography.WriteString("\n")
	}
	
	// Outras referências
	if len(otherRefs) > 0 {
		bibliography.WriteString("## Outras Referências\n\n")
		for _, citation := range otherRefs {
			bibliography.WriteString(fmt.Sprintf("- %s\n", cm.FormatCitation(citation)))
		}
		bibliography.WriteString("\n")
	}
	
	return bibliography.String()
}

// GenerateAcademicCitation gera citação no formato acadêmico
func (cm *CitationManager) GenerateAcademicCitation(citation Citation) string {
	var parts []string
	
	// Autor/Autoridade
	if citation.Authority != "" && citation.Authority != "Autoridade não especificada" {
		parts = append(parts, strings.ToUpper(citation.Authority))
	}
	
	// Título (fonte)
	parts = append(parts, fmt.Sprintf("\"%s\"", citation.Source))
	
	// Seção
	if citation.Section != "" {
		parts = append(parts, citation.Section)
	}
	
	// Data
	if !citation.Date.IsZero() {
		parts = append(parts, fmt.Sprintf("%d", citation.Date.Year()))
	}
	
	// URL se disponível
	if citation.URL != "" {
		parts = append(parts, fmt.Sprintf("Disponível em: %s", citation.URL))
	}
	
	return strings.Join(parts, ". ") + "."
}

// CheckCitationFreshness verifica se as citações estão atualizadas
func (cm *CitationManager) CheckCitationFreshness(citations []Citation) map[string][]Citation {
	result := map[string][]Citation{
		"current":   {},
		"aging":     {},
		"outdated":  {},
	}
	
	now := time.Now()
	
	for _, citation := range citations {
		if citation.Date.IsZero() {
			result["outdated"] = append(result["outdated"], citation)
			continue
		}
		
		age := now.Sub(citation.Date)
		
		// Regulamentações federais (CLT, leis) são consideradas sempre atuais
		if strings.Contains(strings.ToLower(citation.Source), "clt") ||
		   strings.Contains(strings.ToLower(citation.Source), "lei") {
			result["current"] = append(result["current"], citation)
			continue
		}
		
		// Classificar por idade
		if age < 1*365*24*time.Hour {
			result["current"] = append(result["current"], citation)
		} else if age < 3*365*24*time.Hour {
			result["aging"] = append(result["aging"], citation)
		} else {
			result["outdated"] = append(result["outdated"], citation)
		}
	}
	
	return result
}