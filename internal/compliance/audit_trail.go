package compliance

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AuditTrail gerencia a trilha de auditoria do sistema
type AuditTrail struct {
	entries       []AuditEntry
	maxEntries    int
	retentionDays int
}

// NewAuditTrail cria uma nova instância da trilha de auditoria
func NewAuditTrail(maxEntries int, retentionDays int) *AuditTrail {
	return &AuditTrail{
		entries:       []AuditEntry{},
		maxEntries:    maxEntries,
		retentionDays: retentionDays,
	}
}

// RecordAction registra uma ação na trilha de auditoria
func (at *AuditTrail) RecordAction(action, actionType, user, entityID, entityType string, details map[string]interface{}) *AuditEntry {
	entry := AuditEntry{
		ID:         uuid.New().String(),
		Timestamp:  time.Now(),
		Action:     action,
		ActionType: actionType,
		User:       user,
		EntityID:   entityID,
		EntityType: entityType,
		Details:    details,
		Success:    true,
	}

	at.addEntry(entry)
	return &entry
}

// RecordActionWithChanges registra uma ação com detalhes de mudanças
func (at *AuditTrail) RecordActionWithChanges(action, actionType, user, entityID, entityType string,
	details map[string]interface{}, changes []ChangeLog) *AuditEntry {

	entry := AuditEntry{
		ID:         uuid.New().String(),
		Timestamp:  time.Now(),
		Action:     action,
		ActionType: actionType,
		User:       user,
		EntityID:   entityID,
		EntityType: entityType,
		Details:    details,
		Changes:    changes,
		Success:    true,
	}

	at.addEntry(entry)
	return &entry
}

// RecordFailedAction registra uma ação que falhou
func (at *AuditTrail) RecordFailedAction(action, actionType, user, entityID, entityType string,
	details map[string]interface{}, errorMsg string) *AuditEntry {

	entry := AuditEntry{
		ID:         uuid.New().String(),
		Timestamp:  time.Now(),
		Action:     action,
		ActionType: actionType,
		User:       user,
		EntityID:   entityID,
		EntityType: entityType,
		Details:    details,
		Success:    false,
		ErrorMsg:   errorMsg,
	}

	at.addEntry(entry)
	return &entry
}

// RecordComplianceCheck registra uma verificação de compliance
func (at *AuditTrail) RecordComplianceCheck(user, entityID string, result *ComplianceResult) *AuditEntry {
	details := map[string]interface{}{
		"compliance_status": result.Status,
		"score":             result.Score,
		"violations_count":  len(result.Violations),
		"rules_checked":     len(result.RulesChecked),
		"rules_passed":      len(result.RulesPassed),
		"rules_failed":      len(result.RulesFailed),
	}

	return at.RecordAction(
		"Verificação de compliance executada",
		"COMPLIANCE_CHECK",
		user,
		entityID,
		result.EntityType,
		details,
	)
}

// RecordViolationResolution registra resolução de violação
func (at *AuditTrail) RecordViolationResolution(user, violationID string, resolution map[string]interface{}) *AuditEntry {
	details := map[string]interface{}{
		"violation_id": violationID,
		"resolution":   resolution,
		"resolved_at":  time.Now(),
	}

	return at.RecordAction(
		"Violação de compliance resolvida",
		"VIOLATION_RESOLVED",
		user,
		violationID,
		"VIOLATION",
		details,
	)
}

// RecordSystemAccess registra acesso ao sistema
func (at *AuditTrail) RecordSystemAccess(user, userRole, ipAddress, userAgent, sessionID string, success bool) *AuditEntry {
	details := map[string]interface{}{
		"ip_address": ipAddress,
		"user_agent": userAgent,
		"session_id": sessionID,
		"user_role":  userRole,
	}

	entry := AuditEntry{
		ID:         uuid.New().String(),
		Timestamp:  time.Now(),
		Action:     "Acesso ao sistema",
		ActionType: "SYSTEM_ACCESS",
		User:       user,
		UserRole:   userRole,
		Details:    details,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		SessionID:  sessionID,
		Success:    success,
	}

	if !success {
		entry.ErrorMsg = "Falha na autenticação"
	}

	at.addEntry(entry)
	return &entry
}

// RecordDataProcessing registra processamento de dados
func (at *AuditTrail) RecordDataProcessing(user, processType string, recordsProcessed int, duration time.Duration) *AuditEntry {
	details := map[string]interface{}{
		"process_type":      processType,
		"records_processed": recordsProcessed,
		"duration_ms":       duration.Milliseconds(),
		"processing_rate":   float64(recordsProcessed) / duration.Seconds(),
	}

	entry := AuditEntry{
		ID:         uuid.New().String(),
		Timestamp:  time.Now(),
		Action:     fmt.Sprintf("Processamento de dados: %s", processType),
		ActionType: "DATA_PROCESSING",
		User:       user,
		Details:    details,
		Success:    true,
		Duration:   duration,
	}

	at.addEntry(entry)
	return &entry
}

// addEntry adiciona uma entrada na trilha
func (at *AuditTrail) addEntry(entry AuditEntry) {
	at.entries = append(at.entries, entry)

	// Limpar entradas antigas se necessário
	at.cleanup()
}

// cleanup remove entradas antigas baseado na política de retenção
func (at *AuditTrail) cleanup() {
	now := time.Now()
	cutoffDate := now.AddDate(0, 0, -at.retentionDays)

	// Filtrar entradas dentro do período de retenção
	validEntries := []AuditEntry{}
	for _, entry := range at.entries {
		if entry.Timestamp.After(cutoffDate) {
			validEntries = append(validEntries, entry)
		}
	}

	// Limitar número máximo de entradas
	if len(validEntries) > at.maxEntries {
		// Manter as mais recentes
		sort.Slice(validEntries, func(i, j int) bool {
			return validEntries[i].Timestamp.After(validEntries[j].Timestamp)
		})
		validEntries = validEntries[:at.maxEntries]
	}

	at.entries = validEntries
}

// GetEntries retorna entradas filtradas
func (at *AuditTrail) GetEntries(filter AuditFilter) []AuditEntry {
	var filtered []AuditEntry

	for _, entry := range at.entries {
		if at.matchesFilter(entry, filter) {
			filtered = append(filtered, entry)
		}
	}

	// Ordenar por timestamp (mais recente primeiro)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.After(filtered[j].Timestamp)
	})

	// Aplicar limite se especificado
	if filter.Limit > 0 && len(filtered) > filter.Limit {
		filtered = filtered[:filter.Limit]
	}

	return filtered
}

// AuditFilter define filtros para busca de entradas
type AuditFilter struct {
	StartDate  *time.Time
	EndDate    *time.Time
	User       string
	ActionType string
	EntityID   string
	EntityType string
	Success    *bool
	Limit      int
}

// matchesFilter verifica se uma entrada corresponde aos filtros
func (at *AuditTrail) matchesFilter(entry AuditEntry, filter AuditFilter) bool {
	// Filtro por data de início
	if filter.StartDate != nil && entry.Timestamp.Before(*filter.StartDate) {
		return false
	}

	// Filtro por data fim
	if filter.EndDate != nil && entry.Timestamp.After(*filter.EndDate) {
		return false
	}

	// Filtro por usuário
	if filter.User != "" && !strings.Contains(strings.ToLower(entry.User), strings.ToLower(filter.User)) {
		return false
	}

	// Filtro por tipo de ação
	if filter.ActionType != "" && entry.ActionType != filter.ActionType {
		return false
	}

	// Filtro por ID da entidade
	if filter.EntityID != "" && entry.EntityID != filter.EntityID {
		return false
	}

	// Filtro por tipo da entidade
	if filter.EntityType != "" && entry.EntityType != filter.EntityType {
		return false
	}

	// Filtro por sucesso
	if filter.Success != nil && entry.Success != *filter.Success {
		return false
	}

	return true
}

// GenerateReport gera um relatório de auditoria para um período
func (at *AuditTrail) GenerateReport(period DateRange, reportType, generatedBy string) *AuditReport {
	filter := AuditFilter{
		StartDate: &period.StartDate,
		EndDate:   &period.EndDate,
	}

	entries := at.GetEntries(filter)
	summary := at.generateAuditSummary(entries)

	report := &AuditReport{
		GeneratedAt:  time.Now(),
		Period:       period,
		TotalEntries: len(entries),
		Summary:      summary,
		Entries:      entries,
	}

	return report
}

// generateAuditSummary gera resumo estatístico das entradas
func (at *AuditTrail) generateAuditSummary(entries []AuditEntry) AuditSummary {
	summary := AuditSummary{
		EntriesByLevel:  make(map[string]int),
		EntriesByAction: make(map[string]int),
		TopUsers:        []string{},
	}

	userCount := make(map[string]int)

	for _, entry := range entries {
		// Contar por nível
		summary.EntriesByLevel[entry.Level]++

		// Contar por ação
		summary.EntriesByAction[entry.Action]++

		// Contar usuários
		userCount[entry.User]++
	}

	// Get top users
	type userActivity struct {
		user  string
		count int
	}
	users := []userActivity{}
	for user, count := range userCount {
		users = append(users, userActivity{user, count})
	}

	// Sort by activity
	for i := 0; i < len(users)-1; i++ {
		for j := i + 1; j < len(users); j++ {
			if users[i].count < users[j].count {
				users[i], users[j] = users[j], users[i]
			}
		}
	}

	// Top 5 users
	for i, user := range users {
		if i >= 5 {
			break
		}
		summary.TopUsers = append(summary.TopUsers, fmt.Sprintf("%s (%d)", user.user, user.count))
	}

	return summary
}

// GetStatistics retorna estatísticas da trilha de auditoria
func (at *AuditTrail) GetStatistics() map[string]interface{} {
	totalEntries := len(at.entries)

	if totalEntries == 0 {
		return map[string]interface{}{
			"total_entries": 0,
			"oldest_entry":  nil,
			"newest_entry":  nil,
		}
	}

	// Encontrar entrada mais antiga e mais nova
	oldest := at.entries[0].Timestamp
	newest := at.entries[0].Timestamp

	actionTypes := make(map[string]int)
	users := make(map[string]int)
	successCount := 0

	for _, entry := range at.entries {
		if entry.Timestamp.Before(oldest) {
			oldest = entry.Timestamp
		}
		if entry.Timestamp.After(newest) {
			newest = entry.Timestamp
		}

		actionTypes[entry.ActionType]++
		users[entry.User]++

		if entry.Success {
			successCount++
		}
	}

	stats := map[string]interface{}{
		"total_entries":  totalEntries,
		"oldest_entry":   oldest,
		"newest_entry":   newest,
		"success_rate":   float64(successCount) / float64(totalEntries) * 100,
		"action_types":   actionTypes,
		"active_users":   len(users),
		"users":          users,
		"retention_days": at.retentionDays,
		"max_entries":    at.maxEntries,
	}

	return stats
}

// ExportEntries exporta entradas para análise externa
func (at *AuditTrail) ExportEntries(filter AuditFilter, format string) ([]byte, error) {
	entries := at.GetEntries(filter)

	switch strings.ToLower(format) {
	case "json":
		return at.exportJSON(entries)
	case "csv":
		return at.exportCSV(entries)
	default:
		return nil, fmt.Errorf("formato não suportado: %s", format)
	}
}

// exportJSON exporta entradas em formato JSON
func (at *AuditTrail) exportJSON(entries []AuditEntry) ([]byte, error) {
	// Implementação seria feita usando encoding/json
	return []byte("[]"), nil // Placeholder
}

// exportCSV exporta entradas em formato CSV
func (at *AuditTrail) exportCSV(entries []AuditEntry) ([]byte, error) {
	var csv strings.Builder

	// Cabeçalho CSV
	csv.WriteString("Timestamp,Action,ActionType,User,EntityID,EntityType,Success,ErrorMsg\n")

	// Dados
	for _, entry := range entries {
		csv.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%t,%s\n",
			entry.Timestamp.Format(time.RFC3339),
			entry.Action,
			entry.ActionType,
			entry.User,
			entry.EntityID,
			entry.EntityType,
			entry.Success,
			entry.ErrorMsg,
		))
	}

	return []byte(csv.String()), nil
}

// RecordSimpleAction is a simplified method for recording audit actions
func (at *AuditTrail) RecordSimpleAction(userID, action, description, level string, metadata map[string]interface{}) {
	entry := AuditEntry{
		ID:         uuid.New().String(),
		Timestamp:  time.Now(),
		Action:     action,
		ActionType: "USER_ACTION",
		User:       userID,
		EntityID:   userID,
		EntityType: "user",
		Details:    metadata,
		Success:    true,
		ErrorMsg:   "",
		Level:      level,
	}

	if description != "" {
		entry.Action = description
	}

	if level == "ERROR" {
		entry.Success = false
		entry.ErrorMsg = description
	}

	at.addEntry(entry)
}

// GetEntriesForEntity with simplified signature for compatibility
func (at *AuditTrail) GetEntriesForEntity(entityID string, since, until time.Time) []AuditEntry {
	filter := AuditFilter{
		StartDate: &since,
		EndDate:   &until,
	}

	if entityID != "" {
		filter.EntityID = entityID
	}

	return at.GetEntries(filter)
}

// GetEntriesByDateRange returns entries within a date range
func (at *AuditTrail) GetEntriesByDateRange(startDate, endDate time.Time) []AuditEntry {
	filter := AuditFilter{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	return at.GetEntries(filter)
}
