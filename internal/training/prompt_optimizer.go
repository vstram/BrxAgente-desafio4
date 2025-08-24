package training

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// PromptTemplate representa um template de prompt otimizado
type PromptTemplate struct {
	Role         string `yaml:"role"`
	Persona      string `yaml:"persona"`
	Instructions string `yaml:"instructions"`
	ContextInstructions string `yaml:"context_instructions,omitempty"`
}

// PromptConfig representa a configuração completa de prompts
type PromptConfig struct {
	SystemPrompts map[string]PromptTemplate `yaml:"system_prompts"`
	ToolPrompts   map[string]ToolPrompt     `yaml:"tool_prompts"`
	WorkflowPrompts map[string]WorkflowPrompt `yaml:"workflow_prompts"`
	OptimizationTechniques OptimizationTechniques `yaml:"prompt_optimization"`
	ValidationPrompts ValidationPrompts `yaml:"validation_prompts"`
}

// ToolPrompt representa prompts específicos para ferramentas
type ToolPrompt struct {
	Instruction string `yaml:"instruction"`
}

// WorkflowPrompt representa prompts para workflows
type WorkflowPrompt struct {
	Name               string `yaml:"name"`
	Description        string `yaml:"description"`
	OrchestratorPrompt string `yaml:"orchestrator_prompt,omitempty"`
	ValidatorPrompt    string `yaml:"validator_prompt,omitempty"`
	DetectorPrompt     string `yaml:"detector_prompt,omitempty"`
	ReporterPrompt     string `yaml:"reporter_prompt,omitempty"`
	PredictorPrompt    string `yaml:"predictor_prompt,omitempty"`
}

// OptimizationTechniques representa técnicas de otimização de prompts
type OptimizationTechniques struct {
	Techniques map[string]OptimizationTechnique `yaml:"techniques"`
}

// OptimizationTechnique representa uma técnica específica
type OptimizationTechnique struct {
	Description    string   `yaml:"description"`
	ExamplesCount  int      `yaml:"examples_count,omitempty"`
	Format         string   `yaml:"format,omitempty"`
	Keywords       []string `yaml:"keywords,omitempty"`
	Personas       []string `yaml:"personas,omitempty"`
	MaxTokens      int      `yaml:"max_tokens,omitempty"`
	PriorityOrder  []string `yaml:"priority_order,omitempty"`
}

// ValidationPrompts representa prompts para validação de qualidade
type ValidationPrompts struct {
	ResponseQuality ResponseQuality `yaml:"response_quality"`
	ConsistencyCheck ConsistencyCheck `yaml:"consistency_check"`
}

// ResponseQuality representa critérios de qualidade
type ResponseQuality struct {
	Criteria []string           `yaml:"criteria"`
	Scoring  map[string]string  `yaml:"scoring"`
}

// ConsistencyCheck representa verificação de consistência
type ConsistencyCheck struct {
	TestQuestions     []string `yaml:"test_questions"`
	ExpectedPatterns  []string `yaml:"expected_patterns"`
}

// PromptOptimizer otimiza prompts para melhor performance
type PromptOptimizer struct {
	configPath   string
	config       *PromptConfig
	knowledgeManager *KnowledgeManager
	loadedAt     time.Time
}

// NewPromptOptimizer cria um novo otimizador de prompts
func NewPromptOptimizer(configPath string, km *KnowledgeManager) *PromptOptimizer {
	return &PromptOptimizer{
		configPath: configPath,
		knowledgeManager: km,
	}
}

// LoadPromptConfig carrega a configuração de prompts
func (po *PromptOptimizer) LoadPromptConfig() error {
	// Carregar prompts de sistema
	systemPromptsPath := filepath.Join(po.configPath, "prompts", "system_prompts.yaml")
	systemData, err := ioutil.ReadFile(systemPromptsPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar prompts de sistema: %v", err)
	}

	var systemConfig struct {
		SystemPrompts map[string]PromptTemplate `yaml:"system_prompts"`
		PromptOptimization OptimizationTechniques `yaml:"prompt_optimization"`
		ValidationPrompts ValidationPrompts `yaml:"validation_prompts"`
	}
	
	if err := yaml.Unmarshal(systemData, &systemConfig); err != nil {
		return fmt.Errorf("erro ao parsear prompts de sistema: %v", err)
	}

	// Carregar prompts de ferramentas
	toolPromptsPath := filepath.Join(po.configPath, "prompts", "tool_prompts.yaml")
	toolData, err := ioutil.ReadFile(toolPromptsPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar prompts de ferramentas: %v", err)
	}

	var toolConfig struct {
		ToolPrompts map[string]ToolPrompt `yaml:"tool_prompts"`
	}
	
	if err := yaml.Unmarshal(toolData, &toolConfig); err != nil {
		return fmt.Errorf("erro ao parsear prompts de ferramentas: %v", err)
	}

	// Carregar prompts de workflows
	workflowPromptsPath := filepath.Join(po.configPath, "prompts", "workflow_prompts.yaml")
	workflowData, err := ioutil.ReadFile(workflowPromptsPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar prompts de workflows: %v", err)
	}

	var workflowConfig struct {
		WorkflowPrompts map[string]WorkflowPrompt `yaml:"workflow_prompts"`
	}
	
	if err := yaml.Unmarshal(workflowData, &workflowConfig); err != nil {
		return fmt.Errorf("erro ao parsear prompts de workflows: %v", err)
	}

	// Montar configuração completa
	po.config = &PromptConfig{
		SystemPrompts:          systemConfig.SystemPrompts,
		ToolPrompts:           toolConfig.ToolPrompts,
		WorkflowPrompts:       workflowConfig.WorkflowPrompts,
		OptimizationTechniques: systemConfig.PromptOptimization,
		ValidationPrompts:     systemConfig.ValidationPrompts,
	}

	po.loadedAt = time.Now()
	return nil
}

// BuildContextualPrompt constrói um prompt contextualizado para uma pergunta específica
func (po *PromptOptimizer) BuildContextualPrompt(role string, question string, context []KnowledgeItem) (string, error) {
	if po.config == nil {
		return "", fmt.Errorf("configuração de prompts não carregada")
	}

	template, exists := po.config.SystemPrompts[role]
	if !exists {
		return "", fmt.Errorf("template de prompt não encontrado para role: %s", role)
	}

	var promptBuilder strings.Builder

	// Adicionar persona e instruções
	promptBuilder.WriteString(fmt.Sprintf("ROLE: %s\n", template.Role))
	promptBuilder.WriteString(fmt.Sprintf("PERSONA: %s\n\n", template.Persona))
	promptBuilder.WriteString(fmt.Sprintf("INSTRUÇÕES:\n%s\n\n", template.Instructions))

	// Adicionar contexto relevante da base de conhecimento
	if len(context) > 0 {
		promptBuilder.WriteString("CONTEXTO DA BASE DE CONHECIMENTO:\n")
		for i, item := range context {
			if i >= 3 { // Limitar a 3 itens para não sobrecarregar
				break
			}
			promptBuilder.WriteString(fmt.Sprintf("---\n"))
			promptBuilder.WriteString(fmt.Sprintf("ID: %s\n", item.ID))
			promptBuilder.WriteString(fmt.Sprintf("Categoria: %s\n", item.Category))
			promptBuilder.WriteString(fmt.Sprintf("Pergunta: %s\n", item.Question))
			promptBuilder.WriteString(fmt.Sprintf("Resposta: %s\n", item.Answer))
			if len(item.Context) > 0 {
				promptBuilder.WriteString(fmt.Sprintf("Fontes: %s\n", strings.Join(item.Context, ", ")))
			}
			promptBuilder.WriteString("---\n\n")
		}
	}

	// Adicionar instruções específicas de contexto
	if template.ContextInstructions != "" {
		promptBuilder.WriteString(fmt.Sprintf("INSTRUÇÕES DE CONTEXTO:\n%s\n\n", template.ContextInstructions))
	}

	// Adicionar a pergunta atual
	promptBuilder.WriteString(fmt.Sprintf("PERGUNTA DO USUÁRIO:\n%s\n\n", question))

	// Adicionar instruções de formato de resposta
	promptBuilder.WriteString("Forneça uma resposta precisa e bem fundamentada seguindo as instruções acima.")

	return promptBuilder.String(), nil
}

// BuildToolPrompt constrói um prompt específico para uso de ferramentas
func (po *PromptOptimizer) BuildToolPrompt(toolName string, parameters map[string]interface{}) (string, error) {
	if po.config == nil {
		return "", fmt.Errorf("configuração de prompts não carregada")
	}

	toolPrompt, exists := po.config.ToolPrompts[toolName]
	if !exists {
		return "", fmt.Errorf("prompt de ferramenta não encontrado: %s", toolName)
	}

	var promptBuilder strings.Builder
	
	// Instrução base da ferramenta
	promptBuilder.WriteString(toolPrompt.Instruction)
	promptBuilder.WriteString("\n\n")

	// Adicionar parâmetros específicos se fornecidos
	if len(parameters) > 0 {
		promptBuilder.WriteString("PARÂMETROS:\n")
		for key, value := range parameters {
			promptBuilder.WriteString(fmt.Sprintf("- %s: %v\n", key, value))
		}
		promptBuilder.WriteString("\n")
	}

	return promptBuilder.String(), nil
}

// BuildWorkflowPrompt constrói um prompt para orquestração de workflow
func (po *PromptOptimizer) BuildWorkflowPrompt(workflowName string, stage string) (string, error) {
	if po.config == nil {
		return "", fmt.Errorf("configuração de prompts não carregada")
	}

	workflowPrompt, exists := po.config.WorkflowPrompts[workflowName]
	if !exists {
		return "", fmt.Errorf("prompt de workflow não encontrado: %s", workflowName)
	}

	var prompt string
	switch stage {
	case "orchestrator":
		prompt = workflowPrompt.OrchestratorPrompt
	case "validator":
		prompt = workflowPrompt.ValidatorPrompt
	case "detector":
		prompt = workflowPrompt.DetectorPrompt
	case "reporter":
		prompt = workflowPrompt.ReporterPrompt
	case "predictor":
		prompt = workflowPrompt.PredictorPrompt
	default:
		return "", fmt.Errorf("estágio de workflow inválido: %s", stage)
	}

	if prompt == "" {
		return "", fmt.Errorf("prompt não definido para workflow %s estágio %s", workflowName, stage)
	}

	return prompt, nil
}

// OptimizeForTokens otimiza um prompt para reduzir uso de tokens
func (po *PromptOptimizer) OptimizeForTokens(originalPrompt string, targetTokens int) (string, error) {
	// Implementação básica de otimização de tokens
	lines := strings.Split(originalPrompt, "\n")
	var optimizedLines []string
	
	currentTokens := len(strings.Fields(originalPrompt)) // Aproximação simples
	
	for _, line := range lines {
		lineTokens := len(strings.Fields(line))
		if currentTokens-lineTokens > targetTokens {
			// Pular linhas menos importantes para reduzir tokens
			if !strings.Contains(strings.ToLower(line), "importante") && 
			   !strings.Contains(strings.ToLower(line), "crítico") &&
			   !strings.Contains(strings.ToLower(line), "obrigatório") {
				currentTokens -= lineTokens
				continue
			}
		}
		optimizedLines = append(optimizedLines, line)
	}
	
	return strings.Join(optimizedLines, "\n"), nil
}

// ApplyFewShotLearning adiciona exemplos ao prompt para few-shot learning
func (po *PromptOptimizer) ApplyFewShotLearning(basePrompt string, examples []ExampleItem, maxExamples int) string {
	if len(examples) == 0 {
		return basePrompt
	}

	var promptBuilder strings.Builder
	promptBuilder.WriteString(basePrompt)
	promptBuilder.WriteString("\n\nEXEMPLOS DE REFERÊNCIA:\n")

	limit := len(examples)
	if maxExamples > 0 && maxExamples < len(examples) {
		limit = maxExamples
	}

	for i := 0; i < limit; i++ {
		example := examples[i]
		promptBuilder.WriteString(fmt.Sprintf("\nEXEMPLO %d:\n", i+1))
		promptBuilder.WriteString(fmt.Sprintf("Cenário: %s\n", example.Scenario))
		promptBuilder.WriteString(fmt.Sprintf("Resultado: VR = R$ %.2f\n", example.ExpectedResult.ValorVR))
		promptBuilder.WriteString(fmt.Sprintf("Observação: %s\n", example.ExpectedResult.Observacoes))
	}

	promptBuilder.WriteString("\nUse estes exemplos como referência para manter consistência nas suas respostas.\n")

	return promptBuilder.String()
}

// ValidatePromptQuality valida a qualidade de um prompt baseado nos critérios configurados
func (po *PromptOptimizer) ValidatePromptQuality(prompt string) (map[string]bool, float64, error) {
	if po.config == nil {
		return nil, 0, fmt.Errorf("configuração de prompts não carregada")
	}

	criteria := po.config.ValidationPrompts.ResponseQuality.Criteria
	results := make(map[string]bool)
	
	promptLower := strings.ToLower(prompt)
	
	for _, criterion := range criteria {
		switch criterion {
		case "Resposta cita fonte específica?":
			results[criterion] = strings.Contains(promptLower, "política") || 
								strings.Contains(promptLower, "fonte") ||
								strings.Contains(promptLower, "referência")
		case "Cálculo está completo e correto?":
			results[criterion] = strings.Contains(promptLower, "cálculo") && 
								strings.Contains(promptLower, "fórmula")
		case "Exemplo prático foi fornecido quando relevante?":
			results[criterion] = strings.Contains(promptLower, "exemplo")
		case "Linguagem é clara e profissional?":
			// Verificação básica - na prática seria mais sofisticada
			results[criterion] = !strings.Contains(promptLower, "gíria") && 
								len(strings.Fields(prompt)) > 10
		case "Confidencialidade foi mantida (só matrícula)?":
			results[criterion] = !strings.Contains(promptLower, "nome") || 
								strings.Contains(promptLower, "matrícula")
		}
	}

	// Calcular score
	totalCriteria := len(criteria)
	passedCriteria := 0
	for _, passed := range results {
		if passed {
			passedCriteria++
		}
	}

	score := float64(passedCriteria) / float64(totalCriteria)
	return results, score, nil
}

// GetPromptStatistics retorna estatísticas dos prompts configurados
func (po *PromptOptimizer) GetPromptStatistics() map[string]interface{} {
	if po.config == nil {
		return map[string]interface{}{"error": "configuração não carregada"}
	}

	stats := make(map[string]interface{})
	stats["system_prompts_count"] = len(po.config.SystemPrompts)
	stats["tool_prompts_count"] = len(po.config.ToolPrompts)
	stats["workflow_prompts_count"] = len(po.config.WorkflowPrompts)
	stats["optimization_techniques_count"] = len(po.config.OptimizationTechniques.Techniques)
	stats["loaded_at"] = po.loadedAt

	return stats
}