package agent

import (
	"testing"
	"time"

	"BrxAgente-desafio4/internal/chat"
	"BrxAgente-desafio4/internal/config"
	"BrxAgente-desafio4/internal/workflows"
)

// TestWorkflowIntegration testa a integração do sistema de workflows com o VRAgent
func TestWorkflowIntegration(t *testing.T) {
	t.Run("AgentWorkflowInitialization", testAgentWorkflowInitialization)
	t.Run("ExecuteSimpleValidationWorkflow", testExecuteSimpleValidationWorkflow)
	t.Run("ExecuteWorkflowAsync", testExecuteWorkflowAsync)
	t.Run("WorkflowCancellation", testWorkflowCancellation)
	t.Run("ListAvailableWorkflows", testListAvailableWorkflows)
	t.Run("WorkflowWithDisabledAgent", testWorkflowWithDisabledAgent)
	t.Run("WorkflowErrorHandling", testWorkflowErrorHandling)
}

func testAgentWorkflowInitialization(t *testing.T) {
	agent := createTestAgentForWorkflow(t)
	
	// Verificar se orquestrador foi inicializado
	orchestrator := agent.GetWorkflowOrchestrator()
	if orchestrator == nil {
		t.Fatal("Orquestrador não foi inicializado")
	}
	
	// Verificar se workflows padrão foram registrados
	workflows := agent.ListAvailableWorkflows()
	if len(workflows) == 0 {
		t.Error("Nenhum workflow foi registrado")
	}
	
	expectedWorkflows := map[string]bool{
		"simple-validation": false,
	}
	
	for _, workflowName := range workflows {
		if _, exists := expectedWorkflows[workflowName]; exists {
			expectedWorkflows[workflowName] = true
		}
	}
	
	for workflowName, found := range expectedWorkflows {
		if !found {
			t.Errorf("Workflow esperado não encontrado: %s", workflowName)
		}
	}
	
	t.Logf("Workflows registrados: %v", workflows)
}

func testExecuteSimpleValidationWorkflow(t *testing.T) {
	agent := createTestAgentForWorkflow(t)
	
	// Executar workflow de validação simples
	params := map[string]interface{}{
		"file_path": "test_data.xlsx",
	}
	
	result, err := agent.ExecuteWorkflowByName("simple-validation", params)
	if err != nil {
		t.Fatalf("Erro ao executar workflow: %v", err)
	}
	
	// Verificar resultado
	if result == nil {
		t.Fatal("Resultado não deveria ser nil")
	}
	
	if result.WorkflowName != "simple-validation" {
		t.Errorf("Nome do workflow incorreto: %s", result.WorkflowName)
	}
	
	if result.Status != workflows.StatusCompleted {
		t.Errorf("Status deveria ser Completed, got %s", result.Status)
	}
	
	if result.CompletedSteps != 3 {
		t.Errorf("Deveria ter completado 3 steps, got %d", result.CompletedSteps)
	}
	
	// Verificar step results
	if len(result.StepResults) != 3 {
		t.Errorf("Deveria haver 3 step results, got %d", len(result.StepResults))
	}
	
	expectedSteps := []string{"load-data", "validate-data", "report-results"}
	for i, stepResult := range result.StepResults {
		if stepResult.StepName != expectedSteps[i] {
			t.Errorf("Step %d deveria ser %s, got %s", i, expectedSteps[i], stepResult.StepName)
		}
		
		if stepResult.Status != "completed" {
			t.Errorf("Step %s deveria estar completed, got %s", stepResult.StepName, stepResult.Status)
		}
	}
	
	// Verificar se status do agente foi atualizado
	status := agent.GetStatus()
	if status.TotalRequests == 0 {
		t.Error("Total de requests deveria ter aumentado")
	}
	
	t.Logf("Workflow executado com sucesso em %v", result.Duration)
}

func testExecuteWorkflowAsync(t *testing.T) {
	agent := createTestAgentForWorkflow(t)
	
	// Executar workflow assincronamente
	params := map[string]interface{}{
		"file_path": "async_test_data.xlsx",
	}
	
	executionID, err := agent.ExecuteWorkflowAsync("simple-validation", params)
	if err != nil {
		t.Fatalf("Erro ao executar workflow async: %v", err)
	}
	
	if executionID == "" {
		t.Fatal("ExecutionID não deveria estar vazio")
	}
	
	// Aguardar um pouco para a execução iniciar
	time.Sleep(50 * time.Millisecond)
	
	// Tentar obter informações da execução
	execution, err := agent.GetWorkflowExecution(executionID)
	if err != nil {
		t.Logf("Execução pode ter terminado: %v", err)
	} else {
		t.Logf("Execução encontrada: %s, Status: %s", execution.WorkflowName, execution.Status)
	}
	
	// Aguardar execução terminar
	time.Sleep(1 * time.Second)
	
	t.Logf("Workflow async executado com ID: %s", executionID)
}

func testWorkflowCancellation(t *testing.T) {
	agent := createTestAgentForWorkflow(t)
	
	// Executar workflow assincronamente
	params := map[string]interface{}{
		"file_path": "cancellation_test.xlsx",
	}
	
	executionID, err := agent.ExecuteWorkflowAsync("simple-validation", params)
	if err != nil {
		t.Fatalf("Erro ao executar workflow async: %v", err)
	}
	
	// Tentar cancelar imediatamente
	err = agent.CancelWorkflowExecution(executionID)
	if err != nil {
		t.Logf("Aviso: Erro ao cancelar (pode ter terminado): %v", err)
	}
	
	// Tentar cancelar execução inexistente
	err = agent.CancelWorkflowExecution("inexistente-id")
	if err == nil {
		t.Error("Deveria falhar ao cancelar execução inexistente")
	}
	
	t.Logf("Teste de cancelamento completado")
}

func testListAvailableWorkflows(t *testing.T) {
	agent := createTestAgentForWorkflow(t)
	
	workflows := agent.ListAvailableWorkflows()
	if len(workflows) == 0 {
		t.Error("Deveria haver workflows disponíveis")
	}
	
	found := false
	for _, workflow := range workflows {
		if workflow == "simple-validation" {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Workflow simple-validation deveria estar disponível")
	}
	
	t.Logf("Workflows disponíveis: %v", workflows)
}

func testWorkflowWithDisabledAgent(t *testing.T) {
	agent := createTestAgentForWorkflow(t)
	
	// Desabilitar agente
	agent.Disable()
	
	// Tentar executar workflow
	params := map[string]interface{}{
		"file_path": "disabled_test.xlsx",
	}
	
	_, err := agent.ExecuteWorkflowByName("simple-validation", params)
	if err == nil {
		t.Error("Deveria falhar com agente desabilitado")
	}
	
	// Tentar executar workflow async
	_, err = agent.ExecuteWorkflowAsync("simple-validation", params)
	if err == nil {
		t.Error("Workflow async deveria falhar com agente desabilitado")
	}
	
	// Reabilitar agente
	agent.Enable()
	
	// Verificar que funciona novamente
	_, err = agent.ExecuteWorkflowByName("simple-validation", params)
	if err != nil {
		t.Errorf("Workflow deveria funcionar após reabilitar agente: %v", err)
	}
}

func testWorkflowErrorHandling(t *testing.T) {
	agent := createTestAgentForWorkflow(t)
	
	// Tentar executar workflow inexistente
	_, err := agent.ExecuteWorkflowByName("workflow-inexistente", nil)
	if err == nil {
		t.Error("Deveria falhar para workflow inexistente")
	}
	
	// Verificar se contador de erro aumentou
	status := agent.GetStatus()
	if status.ErrorCount == 0 {
		t.Error("Contador de erros deveria ter aumentado")
	}
	
	// Executar workflow válido para verificar que agente ainda funciona
	params := map[string]interface{}{
		"file_path": "error_recovery_test.xlsx",
	}
	
	_, err = agent.ExecuteWorkflowByName("simple-validation", params)
	if err != nil {
		t.Errorf("Workflow válido falhou após erro: %v", err)
	}
	
	t.Log("Teste de tratamento de erros completado")
}

// TestWorkflowPerformance testa performance do sistema de workflows
func TestWorkflowPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de performance em modo short")
	}
	
	agent := createTestAgentForWorkflow(t)
	
	// Executar múltiplos workflows para medir performance
	numWorkflows := 10
	params := map[string]interface{}{
		"file_path": "performance_test.xlsx",
	}
	
	start := time.Now()
	
	for i := 0; i < numWorkflows; i++ {
		result, err := agent.ExecuteWorkflowByName("simple-validation", params)
		if err != nil {
			t.Errorf("Workflow %d falhou: %v", i, err)
			continue
		}
		
		if result.Status != workflows.StatusCompleted {
			t.Errorf("Workflow %d não completou: %s", i, result.Status)
		}
	}
	
	duration := time.Since(start)
	avgDuration := duration / time.Duration(numWorkflows)
	
	t.Logf("Performance Test:")
	t.Logf("- Workflows executados: %d", numWorkflows)
	t.Logf("- Tempo total: %v", duration)
	t.Logf("- Tempo médio por workflow: %v", avgDuration)
	t.Logf("- Workflows por segundo: %.2f", float64(numWorkflows)/duration.Seconds())
	
	// Verificar que performance está razoável (menos de 1 segundo por workflow em média)
	if avgDuration > 1*time.Second {
		t.Errorf("Performance muito baixa: %v por workflow", avgDuration)
	}
}

// Helper function
func createTestAgentForWorkflow(t *testing.T) *VRAgent {
	cfg := &config.Config{}
	chatSvc := chat.NewChat(cfg)
	
	agent, err := NewVRAgent(DefaultAgentConfig(), chatSvc)
	if err != nil {
		t.Fatalf("Erro ao criar agente: %v", err)
	}
	
	return agent
}