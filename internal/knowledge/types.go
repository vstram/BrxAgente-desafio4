package knowledge

import (
	"time"
)

// KnowledgeItem representa um item de conhecimento genérico
type KnowledgeItem struct {
	ID           string            `json:"id"`
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	Source       string            `json:"source"`
	EffectiveDate time.Time        `json:"effective_date"`
	Categories   []string          `json:"categories"`
	Metadata     map[string]string `json:"metadata"`
	Type         string            `json:"type"` // "policy", "regulation", "business_rule"
}

// Policy representa uma política da empresa
type Policy struct {
	KnowledgeItem
	Version      string    `json:"version"`
	ApprovedBy   string    `json:"approved_by"`
	ReviewDate   time.Time `json:"review_date"`
	Status       string    `json:"status"` // "active", "deprecated", "draft"
	Scope        []string  `json:"scope"`  // "VR", "benefits", "HR"
}

// Regulation representa uma regulamentação externa
type Regulation struct {
	KnowledgeItem
	Authority    string `json:"authority"`    // "CLT", "Receita Federal", etc.
	LawNumber    string `json:"law_number"`
	Jurisdiction string `json:"jurisdiction"` // "federal", "estadual", "municipal"
	Mandatory    bool   `json:"mandatory"`
}

// BusinessRule representa uma regra de negócio
type BusinessRule struct {
	KnowledgeItem
	Priority    int      `json:"priority"`    // 1-10, onde 1 é mais alta
	Conditions  []string `json:"conditions"`  // Condições para aplicar a regra
	Actions     []string `json:"actions"`     // Ações a serem tomadas
	Exceptions  []string `json:"exceptions"`  // Exceções à regra
	RelatedRules []string `json:"related_rules"` // IDs de regras relacionadas
}

// KnowledgeBase representa toda a base de conhecimento
type KnowledgeBase struct {
	Policies    []Policy      `json:"policies"`
	Regulations []Regulation  `json:"regulations"`
	Rules       []BusinessRule `json:"business_rules"`
	Version     string        `json:"version"`
	LastUpdated time.Time     `json:"last_updated"`
}

// SearchResult representa um resultado de busca
type SearchResult struct {
	Item       KnowledgeItem `json:"item"`
	Score      float64       `json:"score"`      // Relevância do resultado (0-1)
	Highlights []string      `json:"highlights"` // Trechos que matcharam
	Context    string        `json:"context"`    // Contexto adicional
}

// Citation representa uma citação de fonte
type Citation struct {
	Source      string    `json:"source"`
	Section     string    `json:"section,omitempty"`
	Page        string    `json:"page,omitempty"`
	Date        time.Time `json:"date"`
	URL         string    `json:"url,omitempty"`
	Authority   string    `json:"authority,omitempty"`
	Reliability string    `json:"reliability"` // "high", "medium", "low"
}

// ReasoningStep representa um passo no raciocínio lógico
type ReasoningStep struct {
	Step        int          `json:"step"`
	Description string       `json:"description"`
	RulesUsed   []string     `json:"rules_used"`   // IDs das regras aplicadas
	Logic       string       `json:"logic"`        // Lógica aplicada
	Result      string       `json:"result"`       // Resultado do passo
	Confidence  float64      `json:"confidence"`   // Confiança no resultado (0-1)
	Citations   []Citation   `json:"citations"`
}

// ConsultationResult representa o resultado de uma consulta
type ConsultationResult struct {
	Query           string          `json:"query"`
	Answer          string          `json:"answer"`
	Confidence      float64         `json:"confidence"`
	ReasoningSteps  []ReasoningStep `json:"reasoning_steps"`
	Sources         []Citation      `json:"sources"`
	RelatedTopics   []string        `json:"related_topics"`
	Ambiguities     []string        `json:"ambiguities,omitempty"`
	Recommendations []string        `json:"recommendations,omitempty"`
	ProcessingTime  time.Duration   `json:"processing_time"`
}

// ConflictDetection representa um conflito entre regras
type ConflictDetection struct {
	Type        string   `json:"type"`         // "contradiction", "ambiguity", "gap"
	Description string   `json:"description"`
	RulesInvolved []string `json:"rules_involved"` // IDs das regras em conflito
	Severity    string   `json:"severity"`     // "high", "medium", "low"
	Resolution  string   `json:"resolution,omitempty"` // Sugestão de resolução
	Priority    int      `json:"priority"`     // Prioridade para resolução
}