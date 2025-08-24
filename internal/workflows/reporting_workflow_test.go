package workflows

import (
	"context"
	"fmt"
	"testing"
)

// MockReportGenerator implementa ReportGenerator para testes
type MockReportGenerator struct {
	shouldFail bool
}

func (m *MockReportGenerator) GenerateExecutiveReport(data interface{}, format string) (interface{}, error) {
	if m.shouldFail {
		return nil, fmt.Errorf("mock error")
	}

	return map[string]interface{}{
		"type":   "executive",
		"format": format,
		"data":   data,
	}, nil
}

func (m *MockReportGenerator) GenerateDetailedReport(data interface{}, format string) (interface{}, error) {
	if m.shouldFail {
		return nil, fmt.Errorf("mock error")
	}

	return map[string]interface{}{
		"type":   "detailed",
		"format": format,
		"data":   data,
	}, nil
}

func (m *MockReportGenerator) GenerateAnomalyReport(data interface{}, format string) (interface{}, error) {
	if m.shouldFail {
		return nil, fmt.Errorf("mock error")
	}

	return map[string]interface{}{
		"type":   "anomaly",
		"format": format,
		"data":   data,
	}, nil
}

// MockLogger implementa Logger para testes
type MockLogger struct {
	logs []string
}

func (l *MockLogger) Debug(msg string, args ...interface{}) {
	l.logs = append(l.logs, fmt.Sprintf("DEBUG: "+msg, args...))
}

func (l *MockLogger) Info(msg string, args ...interface{}) {
	l.logs = append(l.logs, fmt.Sprintf("INFO: "+msg, args...))
}

func (l *MockLogger) Warn(msg string, args ...interface{}) {
	l.logs = append(l.logs, fmt.Sprintf("WARN: "+msg, args...))
}

func (l *MockLogger) Error(msg string, args ...interface{}) {
	l.logs = append(l.logs, fmt.Sprintf("ERROR: "+msg, args...))
}

func TestNewReportingWorkflow(t *testing.T) {
	mockGenerator := &MockReportGenerator{}
	config := DefaultReportingWorkflowConfig()

	workflow := NewReportingWorkflow(mockGenerator, config)

	if workflow == nil {
		t.Fatal("Workflow não foi criado")
	}

	if workflow.Name() != "intelligent-reporting" {
		t.Errorf("Nome do workflow incorreto: %s", workflow.Name())
	}

	steps := workflow.Steps()
	if len(steps) != 6 {
		t.Errorf("Número de steps incorreto: %d", len(steps))
	}

	// Verificar se todos os steps foram criados
	expectedSteps := []string{
		"data-preparation",
		"executive-report",
		"detailed-report",
		"anomaly-report",
		"report-summary",
		"notification",
	}

	for i, expectedName := range expectedSteps {
		if steps[i].Name() != expectedName {
			t.Errorf("Step %d tem nome incorreto: esperado %s, obtido %s",
				i, expectedName, steps[i].Name())
		}
	}
}

func TestDataPreparationStep_Execute(t *testing.T) {
	step := NewDataPreparationStep()

	// Criar contexto de teste
	logger := &MockLogger{}
	ctx := NewWorkflowContext(context.Background(), "test", "test-123", nil, logger)

	// Adicionar dados de processamento
	processingData := map[string]interface{}{
		"total_colaboradores": 100,
		"total_vr":            50000.0,
	}
	ctx.Set("processing_results", processingData)

	// Executar step
	err := step.Execute(ctx)
	if err != nil {
		t.Fatalf("Erro ao executar step: %v", err)
	}

	// Verificar se dados preparados foram armazenados
	preparedData, exists := ctx.Get("prepared_data")
	if !exists {
		t.Fatal("Dados preparados não foram armazenados")
	}

	preparedMap, ok := preparedData.(map[string]interface{})
	if !ok {
		t.Fatal("Dados preparados têm formato incorreto")
	}

	// Verificar campos esperados
	if preparedMap["workflow_id"] != "test-123" {
		t.Error("workflow_id não foi definido corretamente")
	}

	if preparedMap["workflow_name"] != "test" {
		t.Error("workflow_name não foi definido corretamente")
	}

	if preparedMap["original_data"] == nil {
		t.Error("original_data não foi preservado")
	}
}

func TestExecutiveReportStep_Execute(t *testing.T) {
	mockGenerator := &MockReportGenerator{}
	config := DefaultReportingWorkflowConfig()
	config.AutoGenerateExecutive = true

	step := NewExecutiveReportStep(mockGenerator, config)

	// Criar contexto de teste
	logger := &MockLogger{}
	ctx := NewWorkflowContext(context.Background(), "test", "test-123", nil, logger)

	// Adicionar dados preparados
	preparedData := map[string]interface{}{
		"workflow_id":         "test-123",
		"total_colaboradores": 100,
	}
	ctx.Set("prepared_data", preparedData)

	// Executar step
	err := step.Execute(ctx)
	if err != nil {
		t.Fatalf("Erro ao executar step: %v", err)
	}

	// Verificar se relatórios foram gerados
	reports, exists := ctx.Get("executive_reports")
	if !exists {
		t.Fatal("Relatórios executivos não foram gerados")
	}

	reportsList, ok := reports.([]interface{})
	if !ok {
		t.Fatal("Formato de relatórios incorreto")
	}

	// Deve gerar relatórios para todos os formatos configurados
	expectedCount := len(config.DefaultFormats)
	if len(reportsList) != expectedCount {
		t.Errorf("Número de relatórios incorreto: esperado %d, obtido %d",
			expectedCount, len(reportsList))
	}
}

func TestExecutiveReportStep_CanSkip(t *testing.T) {
	mockGenerator := &MockReportGenerator{}

	// Teste com geração habilitada
	config := DefaultReportingWorkflowConfig()
	config.AutoGenerateExecutive = true
	step := NewExecutiveReportStep(mockGenerator, config)

	logger := &MockLogger{}
	ctx := NewWorkflowContext(context.Background(), "test", "test-123", nil, logger)

	if step.CanSkip(ctx) {
		t.Error("Step deveria não pular quando geração está habilitada")
	}

	// Teste com geração desabilitada
	config.AutoGenerateExecutive = false
	step = NewExecutiveReportStep(mockGenerator, config)

	if !step.CanSkip(ctx) {
		t.Error("Step deveria pular quando geração está desabilitada")
	}
}

func TestReportSummaryStep_Execute(t *testing.T) {
	step := NewReportSummaryStep()

	// Criar contexto de teste
	logger := &MockLogger{}
	ctx := NewWorkflowContext(context.Background(), "test", "test-123", nil, logger)

	// Adicionar dados de relatórios
	execReports := []interface{}{
		map[string]interface{}{"type": "executive"},
		map[string]interface{}{"type": "executive"},
	}
	ctx.Set("executive_reports", execReports)

	anomReports := []interface{}{
		map[string]interface{}{"type": "anomaly"},
	}
	ctx.Set("anomaly_reports", anomReports)
	ctx.Set("critical_anomalies_detected", true)

	// Executar step
	err := step.Execute(ctx)
	if err != nil {
		t.Fatalf("Erro ao executar step: %v", err)
	}

	// Verificar se resumo foi gerado
	summaryData, exists := ctx.Get("reporting_summary")
	if !exists {
		t.Fatal("Resumo não foi gerado")
	}

	summary, ok := summaryData.(map[string]interface{})
	if !ok {
		t.Fatal("Formato de resumo incorreto")
	}

	// Verificar campos esperados
	if summary["total_reports_generated"] != 3 {
		t.Errorf("Total de relatórios incorreto: %v", summary["total_reports_generated"])
	}

	if summary["priority_level"] != "high" {
		t.Errorf("Nível de prioridade incorreto: %v", summary["priority_level"])
	}

	if summary["executive_reports_count"] != 2 {
		t.Errorf("Contagem de relatórios executivos incorreta: %v", summary["executive_reports_count"])
	}
}

func TestNotificationStep_Execute(t *testing.T) {
	config := DefaultReportingWorkflowConfig()
	config.NotifyStakeholders = true
	config.EmailRecipients = []string{"test@example.com"}

	step := NewNotificationStep(config)

	// Criar contexto de teste
	logger := &MockLogger{}
	ctx := NewWorkflowContext(context.Background(), "test", "test-123", nil, logger)

	// Adicionar resumo
	summary := map[string]interface{}{
		"total_reports_generated": 3,
		"priority_level":          "high",
	}
	ctx.Set("reporting_summary", summary)

	// Executar step
	err := step.Execute(ctx)
	if err != nil {
		t.Fatalf("Erro ao executar step: %v", err)
	}

	// Verificar se notificação foi registrada
	notification, exists := ctx.Get("notification_result")
	if !exists {
		t.Fatal("Notificação não foi registrada")
	}

	notificationMap, ok := notification.(map[string]interface{})
	if !ok {
		t.Fatal("Formato de notificação incorreto")
	}

	if notificationMap["status"] != "sent" {
		t.Errorf("Status de notificação incorreto: %v", notificationMap["status"])
	}

	if notificationMap["priority"] != "high" {
		t.Errorf("Prioridade de notificação incorreta: %v", notificationMap["priority"])
	}
}

func TestContainsFormat(t *testing.T) {
	formats := []string{"excel", "json", "text"}

	if !containsFormat(formats, "excel") {
		t.Error("containsFormat deveria encontrar 'excel'")
	}

	if containsFormat(formats, "pdf") {
		t.Error("containsFormat não deveria encontrar 'pdf'")
	}
}

func TestDefaultReportingWorkflowConfig(t *testing.T) {
	config := DefaultReportingWorkflowConfig()

	if !config.AutoGenerateExecutive {
		t.Error("AutoGenerateExecutive deveria ser true por padrão")
	}

	if config.AutoGenerateDetailed {
		t.Error("AutoGenerateDetailed deveria ser false por padrão")
	}

	if !config.AutoGenerateAnomaly {
		t.Error("AutoGenerateAnomaly deveria ser true por padrão")
	}

	if len(config.DefaultFormats) == 0 {
		t.Error("DefaultFormats não deveria ser vazio")
	}

	if config.OutputDirectory != "./reports" {
		t.Errorf("OutputDirectory incorreto: %s", config.OutputDirectory)
	}
}
