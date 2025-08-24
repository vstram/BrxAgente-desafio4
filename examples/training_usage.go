package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"BrxAgente-desafio4/internal/training"
)

// TrainingUsageExample demonstra como usar o sistema de treinamento
func TrainingUsageExample() {
	// 1. Inicializar sistema de treinamento
	fmt.Println("🤖 Inicializando Sistema de Treinamento do Agente IA")
	fmt.Println("=" + strings.Repeat("=", 60))

	// Configurar caminhos
	basePath := "./data/training"

	// Criar gerenciador de conhecimento
	km := training.NewKnowledgeManager(basePath)

	// Carregar base de conhecimento
	fmt.Println("📚 Carregando base de conhecimento...")
	if err := km.LoadKnowledgeBase(); err != nil {
		log.Fatalf("Erro ao carregar base de conhecimento: %v", err)
	}

	// Obter estatísticas da base de conhecimento
	stats, err := km.GetKnowledgeStatistics()
	if err != nil {
		log.Fatalf("Erro ao obter estatísticas: %v", err)
	}

	fmt.Printf("✅ Base de conhecimento carregada:\n")
	fmt.Printf("   • Políticas: %v\n", stats["total_policies"])
	fmt.Printf("   • FAQ: %v\n", stats["total_faq"])
	fmt.Printf("   • Exemplos: %v\n", stats["total_examples"])
	fmt.Printf("   • Confiança média: %.2f\n", stats["average_confidence"])

	// 2. Inicializar otimizador de prompts
	fmt.Println("\n🔧 Configurando otimizador de prompts...")
	po := training.NewPromptOptimizer(basePath, km)

	if err := po.LoadPromptConfig(); err != nil {
		log.Fatalf("Erro ao carregar configuração de prompts: %v", err)
	}

	promptStats := po.GetPromptStatistics()
	fmt.Printf("✅ Prompts configurados:\n")
	fmt.Printf("   • Prompts de sistema: %v\n", promptStats["system_prompts_count"])
	fmt.Printf("   • Prompts de ferramentas: %v\n", promptStats["tool_prompts_count"])
	fmt.Printf("   • Prompts de workflows: %v\n", promptStats["workflow_prompts_count"])

	// 3. Demonstrar busca de conhecimento relevante
	fmt.Println("\n🔍 Testando busca de conhecimento relevante...")

	testQuestions := []string{
		"Estagiários têm direito a VR?",
		"Como calcular VR para admissão no dia 20?",
		"Qual o valor do VR para SINDPD?",
	}

	for _, question := range testQuestions {
		fmt.Printf("\n❓ Pergunta: %s\n", question)

		relevantKnowledge, err := km.FindRelevantKnowledge(question)
		if err != nil {
			fmt.Printf("   ❌ Erro: %v\n", err)
			continue
		}

		fmt.Printf("   📋 Conhecimento relevante encontrado: %d itens\n", len(relevantKnowledge))

		for i, item := range relevantKnowledge {
			if i >= 2 { // Mostrar apenas os 2 primeiros
				break
			}
			fmt.Printf("     %d. %s (confiança: %.0f%%)\n", i+1, item.Question, item.Confidence*100)
		}

		// Construir prompt contextualizado
		prompt, err := po.BuildContextualPrompt("vr_expert", question, relevantKnowledge)
		if err != nil {
			fmt.Printf("   ❌ Erro ao construir prompt: %v\n", err)
			continue
		}

		fmt.Printf("   📝 Prompt contextualizado gerado (%d caracteres)\n", len(prompt))
	}

	// 4. Inicializar sistema de feedback
	fmt.Println("\n💬 Inicializando sistema de feedback...")
	fs := training.NewFeedbackSystem(basePath+"/feedback", km)

	if err := fs.LoadFeedbackData(); err != nil {
		fmt.Printf("⚠️  Dados de feedback não encontrados, criando novo: %v\n", err)
	}

	// Simular alguns feedbacks
	fmt.Println("📊 Adicionando feedback de exemplo...")

	sampleFeedbacks := []training.ResponseFeedback{
		{
			Question:   "Estagiários têm direito a VR?",
			Response:   "Não. Estagiários são excluídos conforme Política VR-003.",
			UserRating: 5,
			Feedback:   "Resposta perfeita, clara e com fonte citada",
			Source:     "user",
			Category:   "elegibilidade",
		},
		{
			Question:    "Como calcular VR para SINDPD?",
			Response:    "Valor base R$ 467,00 para colaboradores SINDPD",
			UserRating:  3,
			Corrections: "Faltou explicar o cálculo proporcional",
			Source:      "user",
			Category:    "cálculos",
		},
		{
			Question:    "Qual a política para diretores?",
			Response:    "Diretores não recebem VR",
			UserRating:  2,
			Corrections: "Não citou a política específica (VR-004)",
			Source:      "expert",
			Category:    "elegibilidade",
		},
	}

	for _, feedback := range sampleFeedbacks {
		if err := fs.AddFeedback(feedback); err != nil {
			fmt.Printf("   ❌ Erro ao adicionar feedback: %v\n", err)
		}
	}

	// Obter métricas de qualidade
	metrics := fs.GetQualityMetrics()
	fmt.Printf("✅ Métricas de qualidade atualizadas:\n")
	fmt.Printf("   • Total de respostas: %d\n", metrics.TotalResponses)
	fmt.Printf("   • Rating médio: %.2f/5\n", metrics.AverageRating)
	fmt.Printf("   • Score de precisão: %.1f%%\n", metrics.AccuracyScore*100)
	fmt.Printf("   • Score de consistência: %.1f%%\n", metrics.ConsistencyScore*100)

	// Obter padrões de aprendizado
	patterns := fs.GetLearningPatterns()
	fmt.Printf("   • Padrões identificados: %d\n", len(patterns))

	for i, pattern := range patterns {
		if i >= 3 { // Mostrar apenas os 3 primeiros
			break
		}
		fmt.Printf("     %d. %s (freq: %d, cat: %s)\n", i+1, pattern.Pattern, pattern.Frequency, pattern.Category)
	}

	// 5. Executar suíte de validação
	fmt.Println("\n🧪 Executando suíte de validação...")
	vs := training.NewValidationSuite(basePath+"/validation", km, fs)

	if err := vs.LoadTestSuite(); err != nil {
		log.Fatalf("Erro ao carregar suíte de testes: %v", err)
	}

	fmt.Println("⏳ Executando todos os testes (pode demorar alguns segundos)...")

	testResults, err := vs.RunAllTests()
	if err != nil {
		log.Fatalf("Erro ao executar testes: %v", err)
	}

	// Mostrar resultados dos testes
	fmt.Printf("✅ Testes concluídos:\n")

	if summary, ok := testResults["summary"].(map[string]interface{}); ok {
		fmt.Printf("   • Status geral: %v\n", summary["status"])
		fmt.Printf("   • Taxa de aprovação: %.1f%%\n", summary["overall_pass_rate"].(float64)*100)
		fmt.Printf("   • Score médio: %.2f\n", summary["average_score"])
		fmt.Printf("   • Tempo médio de resposta: %.2fs\n", summary["average_response_time"])
		fmt.Printf("   • Total de testes: %v\n", summary["total_tests"])
	}

	// Mostrar breakdown por categoria
	fmt.Printf("\n📊 Breakdown por categoria:\n")
	categories := []string{"eligibility_tests", "calculation_tests", "consistency_tests", "quality_tests", "performance_tests"}

	for _, category := range categories {
		if categoryResult, ok := testResults[category].(map[string]interface{}); ok {
			passed := categoryResult["passed"].(int)
			total := categoryResult["total"].(int)
			passRate := categoryResult["pass_rate"].(float64)

			status := "✅"
			if passRate < 0.9 {
				status = "⚠️"
			}
			if passRate < 0.8 {
				status = "❌"
			}

			fmt.Printf("   %s %s: %d/%d (%.1f%%)\n", status, category, passed, total, passRate*100)
		}
	}

	// 6. Salvar resultados
	fmt.Println("\n💾 Salvando resultados...")

	if err := vs.SaveTestResults(); err != nil {
		fmt.Printf("⚠️  Erro ao salvar resultados dos testes: %v\n", err)
	}

	// 7. Gerar relatório de melhorias
	fmt.Println("\n📈 Gerando relatório de melhorias...")

	improvementReport := fs.GenerateImprovementReport()

	if improvements, ok := improvementReport["priority_improvements"].([]map[string]interface{}); ok {
		fmt.Printf("🎯 Principais áreas de melhoria:\n")

		for i, improvement := range improvements {
			if i >= 3 { // Mostrar apenas as 3 principais
				break
			}
			area := improvement["area"].(string)
			frequency := improvement["frequency"].(int)
			suggestion := improvement["suggestion"].(string)

			fmt.Printf("   %d. %s (freq: %d)\n", i+1, area, frequency)
			fmt.Printf("      💡 %s\n", suggestion)
		}
	}

	// 8. Demonstrar otimização de prompts
	fmt.Println("\n🎨 Demonstrando otimização de prompts...")

	// Construir prompt para ferramenta
	toolPrompt, err := po.BuildToolPrompt("calculate_vr_tool", map[string]interface{}{
		"matricula":        "MAT001234",
		"sindicato":        "SINDPD",
		"dias_trabalhados": 15,
	})
	if err != nil {
		fmt.Printf("❌ Erro ao construir prompt de ferramenta: %v\n", err)
	} else {
		fmt.Printf("✅ Prompt de ferramenta gerado (%d caracteres)\n", len(toolPrompt))
	}

	// Construir prompt de workflow
	workflowPrompt, err := po.BuildWorkflowPrompt("complete_vr_processing", "orchestrator")
	if err != nil {
		fmt.Printf("❌ Erro ao construir prompt de workflow: %v\n", err)
	} else {
		fmt.Printf("✅ Prompt de workflow gerado (%d caracteres)\n", len(workflowPrompt))
	}

	// Finalizar
	fmt.Println("\n🎉 Sistema de treinamento demonstrado com sucesso!")
	fmt.Println("=" + strings.Repeat("=", 60))

	// Mostrar próximos passos
	fmt.Println("\n📋 Próximos passos recomendados:")
	fmt.Println("   1. Integrar com agente LangChain real")
	fmt.Println("   2. Configurar execução automática dos testes")
	fmt.Println("   3. Implementar coleta contínua de feedback")
	fmt.Println("   4. Expandir base de conhecimento com casos reais")
	fmt.Println("   5. Otimizar prompts com base nos resultados")

	fmt.Printf("\n⏰ Demonstração concluída em: %v\n", time.Now().Format("15:04:05"))
}

// DemonstrateFeedbackLoop demonstra o ciclo de feedback e melhoria
func DemonstrateFeedbackLoop() {
	fmt.Println("\n🔄 Demonstrando Ciclo de Feedback e Melhoria")
	fmt.Println("-" + strings.Repeat("-", 50))

	// Simular evolução do sistema ao longo do tempo
	scenarios := []struct {
		day      int
		feedback string
		rating   int
	}{
		{1, "Ótima resposta, muito clara", 5},
		{2, "Faltou citar a política específica", 3},
		{3, "Cálculo correto, mas poderia ter exemplo", 4},
		{5, "Perfeito! Fonte citada e exemplo prático", 5},
		{7, "Inconsistente com resposta anterior", 2},
		{10, "Excelente melhoria na consistência", 5},
	}

	fmt.Println("📊 Simulando evolução do sistema:")

	for _, scenario := range scenarios {
		fmt.Printf("   Dia %2d: Rating %d/5 - %s\n", scenario.day, scenario.rating, scenario.feedback)
	}

	// Calcular tendência
	var totalRating float64
	for _, s := range scenarios {
		totalRating += float64(s.rating)
	}
	avgRating := totalRating / float64(len(scenarios))

	fmt.Printf("\n📈 Rating médio final: %.2f/5\n", avgRating)

	if avgRating >= 4.5 {
		fmt.Println("🎉 Sistema está performando excelentemente!")
	} else if avgRating >= 4.0 {
		fmt.Println("✅ Sistema está performando bem")
	} else {
		fmt.Println("⚠️  Sistema precisa de melhorias")
	}
}
