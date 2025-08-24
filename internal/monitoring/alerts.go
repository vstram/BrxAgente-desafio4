package monitoring

import (
	"fmt"
	"sync"
	"time"
)

// Alert representa um alerta do sistema
type Alert struct {
	ID          string
	Type        AlertType
	Severity    AlertSeverity
	Title       string
	Description string
	Timestamp   time.Time
	Source      string
	Metadata    map[string]interface{}
	Resolved    bool
	ResolvedAt  *time.Time
}

// AlertType tipos de alerta
type AlertType string

const (
	AlertTypePerformance AlertType = "performance"
	AlertTypeMemory      AlertType = "memory"
	AlertTypeCPU         AlertType = "cpu"
	AlertTypeCache       AlertType = "cache"
	AlertTypeError       AlertType = "error"
	AlertTypeSystem      AlertType = "system"
)

// AlertSeverity severidade do alerta
type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "info"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityCritical AlertSeverity = "critical"
)

// AlertRule regra para geração de alertas
type AlertRule struct {
	ID          string
	Name        string
	Type        AlertType
	Enabled     bool
	Condition   func(metrics PerformanceMetrics) bool
	Threshold   float64
	Duration    time.Duration // Duração que a condição deve persistir
	Cooldown    time.Duration // Tempo entre alertas similares
	Severity    AlertSeverity
	Description string
	lastAlert   time.Time
}

// AlertManager gerencia alertas do sistema
type AlertManager struct {
	alerts   map[string]*Alert
	rules    map[string]*AlertRule
	handlers []AlertHandler
	mutex    sync.RWMutex
	enabled  bool
	history  []Alert // Histórico de alertas
}

// AlertHandler interface para handlers de alertas
type AlertHandler interface {
	HandleAlert(alert Alert) error
	GetName() string
}

// NewAlertManager cria um novo gerenciador de alertas
func NewAlertManager() *AlertManager {
	manager := &AlertManager{
		alerts:   make(map[string]*Alert),
		rules:    make(map[string]*AlertRule),
		handlers: make([]AlertHandler, 0),
		enabled:  true,
		history:  make([]Alert, 0),
	}

	// Adiciona regras padrão
	manager.addDefaultRules()

	return manager
}

// addDefaultRules adiciona regras de alerta padrão
func (am *AlertManager) addDefaultRules() {
	// Regra de alta utilização de memória
	am.AddRule(&AlertRule{
		ID:          "high_memory_usage",
		Name:        "Alta utilização de memória",
		Type:        AlertTypeMemory,
		Enabled:     true,
		Threshold:   500, // 500MB
		Duration:    2 * time.Minute,
		Cooldown:    10 * time.Minute,
		Severity:    AlertSeverityWarning,
		Description: "Uso de memória acima de 500MB",
		Condition: func(metrics PerformanceMetrics) bool {
			return metrics.MemoryUsage.AllocMB > 500
		},
	})

	// Regra de alto uso de CPU
	am.AddRule(&AlertRule{
		ID:          "high_cpu_usage",
		Name:        "Alto uso de CPU",
		Type:        AlertTypeCPU,
		Enabled:     true,
		Threshold:   80, // 80%
		Duration:    1 * time.Minute,
		Cooldown:    5 * time.Minute,
		Severity:    AlertSeverityWarning,
		Description: "Uso de CPU acima de 80%",
		Condition: func(metrics PerformanceMetrics) bool {
			return metrics.CPUUsage.Usage > 80
		},
	})

	// Regra de baixo hit ratio do cache
	am.AddRule(&AlertRule{
		ID:          "low_cache_hit_ratio",
		Name:        "Baixo hit ratio do cache",
		Type:        AlertTypeCache,
		Enabled:     true,
		Threshold:   0.5, // 50%
		Duration:    5 * time.Minute,
		Cooldown:    15 * time.Minute,
		Severity:    AlertSeverityInfo,
		Description: "Hit ratio do cache abaixo de 50%",
		Condition: func(metrics PerformanceMetrics) bool {
			return metrics.ItemsProcessed > 100 && metrics.CacheHitRatio < 0.5
		},
	})

	// Regra de alta taxa de erro
	am.AddRule(&AlertRule{
		ID:          "high_error_rate",
		Name:        "Alta taxa de erro",
		Type:        AlertTypeError,
		Enabled:     true,
		Threshold:   5, // 5%
		Duration:    1 * time.Minute,
		Cooldown:    5 * time.Minute,
		Severity:    AlertSeverityCritical,
		Description: "Taxa de erro acima de 5%",
		Condition: func(metrics PerformanceMetrics) bool {
			return metrics.ItemsProcessed > 50 && metrics.ErrorRate > 5
		},
	})

	// Regra de processamento lento
	am.AddRule(&AlertRule{
		ID:          "slow_processing",
		Name:        "Processamento lento",
		Type:        AlertTypePerformance,
		Enabled:     true,
		Threshold:   10, // 10 segundos
		Duration:    3 * time.Minute,
		Cooldown:    10 * time.Minute,
		Severity:    AlertSeverityWarning,
		Description: "Tempo médio de processamento acima de 10 segundos",
		Condition: func(metrics PerformanceMetrics) bool {
			return metrics.AverageProcessingTime > 10*time.Second
		},
	})
}

// AddRule adiciona uma regra de alerta
func (am *AlertManager) AddRule(rule *AlertRule) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	am.rules[rule.ID] = rule
}

// RemoveRule remove uma regra de alerta
func (am *AlertManager) RemoveRule(ruleID string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	delete(am.rules, ruleID)
}

// AddHandler adiciona um handler de alertas
func (am *AlertManager) AddHandler(handler AlertHandler) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	am.handlers = append(am.handlers, handler)
}

// CheckMetrics verifica métricas contra regras e gera alertas
func (am *AlertManager) CheckMetrics(metrics PerformanceMetrics) []Alert {
	if !am.enabled {
		return nil
	}

	am.mutex.Lock()
	defer am.mutex.Unlock()

	var newAlerts []Alert

	for _, rule := range am.rules {
		if !rule.Enabled {
			continue
		}

		// Verifica cooldown
		if time.Since(rule.lastAlert) < rule.Cooldown {
			continue
		}

		// Avalia condição
		if rule.Condition(metrics) {
			alert := am.createAlert(rule, metrics)
			newAlerts = append(newAlerts, alert)

			// Atualiza timestamp do último alerta
			rule.lastAlert = time.Now()

			// Adiciona aos alertas ativos
			am.alerts[alert.ID] = &alert

			// Adiciona ao histórico
			am.history = append(am.history, alert)

			// Limita tamanho do histórico
			if len(am.history) > 1000 {
				am.history = am.history[100:]
			}

			// Processa handlers
			go am.processHandlers(alert)
		}
	}

	return newAlerts
}

// createAlert cria um novo alerta baseado na regra
func (am *AlertManager) createAlert(rule *AlertRule, metrics PerformanceMetrics) Alert {
	alertID := fmt.Sprintf("%s_%d", rule.ID, time.Now().Unix())

	metadata := map[string]interface{}{
		"rule_id":             rule.ID,
		"threshold":           rule.Threshold,
		"current_memory_mb":   metrics.MemoryUsage.AllocMB,
		"current_cpu_pct":     metrics.CPUUsage.Usage,
		"current_error_rate":  metrics.ErrorRate,
		"cache_hit_ratio":     metrics.CacheHitRatio,
		"items_processed":     metrics.ItemsProcessed,
		"avg_processing_time": metrics.AverageProcessingTime.Seconds(),
	}

	return Alert{
		ID:          alertID,
		Type:        rule.Type,
		Severity:    rule.Severity,
		Title:       rule.Name,
		Description: rule.Description,
		Timestamp:   time.Now(),
		Source:      "monitoring_system",
		Metadata:    metadata,
		Resolved:    false,
	}
}

// processHandlers processa todos os handlers para um alerta
func (am *AlertManager) processHandlers(alert Alert) {
	for _, handler := range am.handlers {
		go func(h AlertHandler) {
			if err := h.HandleAlert(alert); err != nil {
				// Log do erro (em produção usaria logger proper)
				fmt.Printf("Erro processando alerta com handler %s: %v\n", h.GetName(), err)
			}
		}(handler)
	}
}

// ResolveAlert resolve um alerta
func (am *AlertManager) ResolveAlert(alertID string) bool {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if alert, exists := am.alerts[alertID]; exists {
		now := time.Now()
		alert.Resolved = true
		alert.ResolvedAt = &now
		return true
	}

	return false
}

// GetActiveAlerts retorna alertas ativos
func (am *AlertManager) GetActiveAlerts() []Alert {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	alerts := make([]Alert, 0, len(am.alerts))
	for _, alert := range am.alerts {
		if !alert.Resolved {
			alerts = append(alerts, *alert)
		}
	}

	return alerts
}

// GetAlertHistory retorna histórico de alertas
func (am *AlertManager) GetAlertHistory(limit int) []Alert {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	if limit <= 0 || limit > len(am.history) {
		limit = len(am.history)
	}

	// Retorna os mais recentes primeiro
	recent := make([]Alert, limit)
	for i := 0; i < limit; i++ {
		recent[i] = am.history[len(am.history)-1-i]
	}

	return recent
}

// GetAlertStats retorna estatísticas de alertas
func (am *AlertManager) GetAlertStats() AlertStats {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	stats := AlertStats{
		TotalAlerts:  len(am.history),
		ActiveAlerts: 0,
		ByType:       make(map[AlertType]int),
		BySeverity:   make(map[AlertSeverity]int),
		RulesEnabled: len(am.rules),
	}

	// Conta alertas ativos
	for _, alert := range am.alerts {
		if !alert.Resolved {
			stats.ActiveAlerts++
		}
	}

	// Estatísticas por tipo e severidade
	for _, alert := range am.history {
		stats.ByType[alert.Type]++
		stats.BySeverity[alert.Severity]++
	}

	return stats
}

// AlertStats estatísticas de alertas
type AlertStats struct {
	TotalAlerts  int
	ActiveAlerts int
	ByType       map[AlertType]int
	BySeverity   map[AlertSeverity]int
	RulesEnabled int
}

// Enable habilita o sistema de alertas
func (am *AlertManager) Enable() {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	am.enabled = true
}

// Disable desabilita o sistema de alertas
func (am *AlertManager) Disable() {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	am.enabled = false
}

// ClearHistory limpa o histórico de alertas
func (am *AlertManager) ClearHistory() {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	am.history = make([]Alert, 0)
}

// LogHandler handler que faz log dos alertas
type LogHandler struct {
	name string
}

// NewLogHandler cria um novo handler de log
func NewLogHandler() *LogHandler {
	return &LogHandler{name: "log_handler"}
}

func (lh *LogHandler) HandleAlert(alert Alert) error {
	fmt.Printf("[ALERT] %s - %s: %s (Severity: %s)\n",
		alert.Timestamp.Format("2006-01-02 15:04:05"),
		alert.Type,
		alert.Title,
		alert.Severity)
	return nil
}

func (lh *LogHandler) GetName() string {
	return lh.name
}

// CallbackHandler handler que executa callback customizado
type CallbackHandler struct {
	name     string
	callback func(Alert) error
}

// NewCallbackHandler cria um novo handler de callback
func NewCallbackHandler(name string, callback func(Alert) error) *CallbackHandler {
	return &CallbackHandler{
		name:     name,
		callback: callback,
	}
}

func (ch *CallbackHandler) HandleAlert(alert Alert) error {
	return ch.callback(alert)
}

func (ch *CallbackHandler) GetName() string {
	return ch.name
}
