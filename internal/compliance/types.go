package compliance

import (
	"time"
)

// ComplianceStatus representa o status de compliance
type ComplianceStatus string

const (
	StatusCompliant    ComplianceStatus = "COMPLIANT"
	StatusNonCompliant ComplianceStatus = "NON_COMPLIANT"
	StatusWarning      ComplianceStatus = "WARNING"
	StatusUnderReview  ComplianceStatus = "UNDER_REVIEW"
)

// ViolationSeverity representa a severidade de uma violação
type ViolationSeverity string

const (
	SeverityCritical ViolationSeverity = "CRITICAL"
	SeverityHigh     ViolationSeverity = "HIGH"
	SeverityMedium   ViolationSeverity = "MEDIUM"
	SeverityLow      ViolationSeverity = "LOW"
)

// RiskLevel representa o nível de risco
type RiskLevel string

const (
	RiskCritical RiskLevel = "CRITICAL"
	RiskHigh     RiskLevel = "HIGH"
	RiskMedium   RiskLevel = "MEDIUM"
	RiskLow      RiskLevel = "LOW"
	RiskNone     RiskLevel = "NONE"
)

// ComplianceRule representa uma regra de compliance
type ComplianceRule struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Category      string                 `json:"category"` // "CLT", "PAT", "MTE", "INTERNAL"
	Type          string                 `json:"type"`     // "REGULATORY", "POLICY", "PROCEDURAL"
	Mandatory     bool                   `json:"mandatory"`
	EffectiveDate time.Time              `json:"effective_date"`
	ExpiryDate    *time.Time             `json:"expiry_date,omitempty"`
	Parameters    map[string]interface{} `json:"parameters"`
	Evidence      []string               `json:"evidence"`   // Tipos de evidência necessária
	Penalties     []string               `json:"penalties"`  // Penalidades por não conformidade
	References    []string               `json:"references"` // Referências legais
}

// ComplianceViolation representa uma violação de compliance
type ComplianceViolation struct {
	ID              string                 `json:"id"`
	RuleID          string                 `json:"rule_id"`
	RuleName        string                 `json:"rule_name"`
	Description     string                 `json:"description"`
	Severity        ViolationSeverity      `json:"severity"`
	Category        string                 `json:"category"`
	EntityID        string                 `json:"entity_id"` // ID da entidade que violou
	Details         map[string]interface{} `json:"details"`
	Evidence        []string               `json:"evidence"`
	Impact          string                 `json:"impact"`
	Remediation     string                 `json:"remediation"`
	DetectedAt      time.Time              `json:"detected_at"`
	DueDate         *time.Time             `json:"due_date,omitempty"`
	Resolved        bool                   `json:"resolved"`
	ResolvedAt      *time.Time             `json:"resolved_at,omitempty"`
	ResolvedBy      string                 `json:"resolved_by,omitempty"`
	Recommendations []string               `json:"recommendations,omitempty"`
}

// Risk representa um risco identificado
type Risk struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Level       RiskLevel              `json:"level"`
	Category    string                 `json:"category"`
	Probability float64                `json:"probability"` // 0-1
	Impact      float64                `json:"impact"`      // 0-1
	Score       float64                `json:"score"`       // Probability * Impact
	Triggers    []string               `json:"triggers"`    // Condições que ativam o risco
	Indicators  map[string]interface{} `json:"indicators"`  // Indicadores do risco
	Mitigation  []string               `json:"mitigation"`  // Ações de mitigação
	Owner       string                 `json:"owner"`       // Responsável
	ReviewDate  time.Time              `json:"review_date"`
	Status      string                 `json:"status"` // "ACTIVE", "MITIGATED", "ACCEPTED"
}

// ComplianceResult representa o resultado de uma verificação
type ComplianceResult struct {
	ID              string                 `json:"id"`
	EntityID        string                 `json:"entity_id"`   // ID da entidade verificada
	EntityType      string                 `json:"entity_type"` // "EMPLOYEE", "PROCESS", "SYSTEM"
	Status          ComplianceStatus       `json:"status"`
	Score           float64                `json:"score"` // 0-100
	CheckedAt       time.Time              `json:"checked_at"`
	CheckedBy       string                 `json:"checked_by"`
	RulesChecked    []string               `json:"rules_checked"` // IDs das regras verificadas
	RulesPassed     []string               `json:"rules_passed"`  // IDs das regras aprovadas
	RulesFailed     []string               `json:"rules_failed"`  // IDs das regras reprovadas
	Violations      []ComplianceViolation  `json:"violations"`
	Risks           []Risk                 `json:"risks"`
	Recommendations []string               `json:"recommendations"`
	Evidence        []Evidence             `json:"evidence"`
	NextReview      time.Time              `json:"next_review"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// Evidence representa uma evidência coletada
type Evidence struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // "DOCUMENT", "RECORD", "SCREENSHOT", "LOG"
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Source      string                 `json:"source"`
	CollectedAt time.Time              `json:"collected_at"`
	CollectedBy string                 `json:"collected_by"`
	Hash        string                 `json:"hash,omitempty"` // Hash para integridade
	Size        int64                  `json:"size,omitempty"` // Tamanho em bytes
	MimeType    string                 `json:"mime_type,omitempty"`
	Path        string                 `json:"path,omitempty"` // Caminho do arquivo
	Metadata    map[string]interface{} `json:"metadata"`
	RelatedTo   []string               `json:"related_to"` // IDs relacionados
}

// AuditEntry representa uma entrada na trilha de auditoria
type AuditEntry struct {
	ID         string                 `json:"id"`
	Timestamp  time.Time              `json:"timestamp"`
	Action     string                 `json:"action"`
	ActionType string                 `json:"action_type"` // "CREATE", "UPDATE", "DELETE", "VIEW", "PROCESS"
	User       string                 `json:"user"`
	UserRole   string                 `json:"user_role,omitempty"`
	EntityID   string                 `json:"entity_id,omitempty"`
	EntityType string                 `json:"entity_type,omitempty"`
	Details    map[string]interface{} `json:"details"`
	Changes    []ChangeLog            `json:"changes,omitempty"`
	Evidence   []string               `json:"evidence,omitempty"` // IDs de evidências
	IPAddress  string                 `json:"ip_address,omitempty"`
	UserAgent  string                 `json:"user_agent,omitempty"`
	SessionID  string                 `json:"session_id,omitempty"`
	Success    bool                   `json:"success"`
	ErrorMsg   string                 `json:"error_msg,omitempty"`
	Duration   time.Duration          `json:"duration,omitempty"`
	Level      string                 `json:"level,omitempty"` // "INFO", "WARN", "ERROR"
}

// ChangeLog representa uma mudança específica
type ChangeLog struct {
	Field     string      `json:"field"`
	OldValue  interface{} `json:"old_value"`
	NewValue  interface{} `json:"new_value"`
	Timestamp time.Time   `json:"timestamp"`
}

// AuditFinding representa um achado de auditoria
type AuditFinding struct {
	ID             string            `json:"id"`
	Type           string            `json:"type"` // "VIOLATION", "RISK", "IMPROVEMENT"
	Category       string            `json:"category"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Severity       ViolationSeverity `json:"severity"`
	Evidence       []string          `json:"evidence"` // IDs de evidências
	Impact         string            `json:"impact"`
	Recommendation string            `json:"recommendation"`
	Priority       int               `json:"priority"` // 1-5
	AssignedTo     string            `json:"assigned_to,omitempty"`
	DueDate        *time.Time        `json:"due_date,omitempty"`
	Status         string            `json:"status"` // "OPEN", "IN_PROGRESS", "RESOLVED", "CLOSED"
}

// CompliancePackage representa um pacote completo de compliance
type CompliancePackage struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	Type             string                 `json:"type"` // "INTERNAL", "EXTERNAL", "REGULATORY"
	CreatedAt        time.Time              `json:"created_at"`
	CreatedBy        string                 `json:"created_by"`
	ComplianceResult *ComplianceResult      `json:"compliance_result"`
	AuditReport      *AuditReport           `json:"audit_report"`
	Documentation    []DocumentationItem    `json:"documentation"`
	Evidence         []Evidence             `json:"evidence"`
	Attachments      []string               `json:"attachments"` // Caminhos de arquivos
	Recipients       []string               `json:"recipients"`
	DeliveryDate     *time.Time             `json:"delivery_date,omitempty"`
	Status           string                 `json:"status"` // "PREPARING", "READY", "DELIVERED", "ACCEPTED"
	Metadata         map[string]interface{} `json:"metadata"`
}

// DocumentationItem representa um item de documentação
type DocumentationItem struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // "POLICY", "PROCEDURE", "REPORT", "CERTIFICATE"
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	Template    string                 `json:"template,omitempty"`
	Format      string                 `json:"format"` // "PDF", "DOCX", "HTML", "MARKDOWN"
	GeneratedAt time.Time              `json:"generated_at"`
	GeneratedBy string                 `json:"generated_by"`
	Version     string                 `json:"version"`
	Hash        string                 `json:"hash,omitempty"`
	Path        string                 `json:"path,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
	RelatedTo   []string               `json:"related_to"` // IDs relacionados
}

// ComplianceConfiguration representa configurações de compliance
type ComplianceConfiguration struct {
	ID                   string                 `json:"id"`
	Name                 string                 `json:"name"`
	Description          string                 `json:"description"`
	EnabledRules         []string               `json:"enabled_rules"`
	RuleConfigurations   map[string]interface{} `json:"rule_configurations"`
	ScheduledChecks      []ScheduledCheck       `json:"scheduled_checks"`
	NotificationSettings NotificationSettings   `json:"notification_settings"`
	ReportingSettings    ReportingSettings      `json:"reporting_settings"`
	RetentionPolicy      RetentionPolicy        `json:"retention_policy"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	CreatedBy            string                 `json:"created_by"`
	UpdatedBy            string                 `json:"updated_by"`
	Active               bool                   `json:"active"`
}

// ScheduledCheck representa uma verificação agendada
type ScheduledCheck struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Rules       []string               `json:"rules"`
	Schedule    string                 `json:"schedule"` // Cron expression
	Scope       []string               `json:"scope"`    // Entidades a verificar
	Enabled     bool                   `json:"enabled"`
	LastRun     *time.Time             `json:"last_run,omitempty"`
	NextRun     *time.Time             `json:"next_run,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// NotificationSettings representa configurações de notificação
type NotificationSettings struct {
	Enabled           bool              `json:"enabled"`
	EmailRecipients   []string          `json:"email_recipients"`
	SlackChannels     []string          `json:"slack_channels,omitempty"`
	WebhookURLs       []string          `json:"webhook_urls,omitempty"`
	NotifyOnViolation bool              `json:"notify_on_violation"`
	NotifyOnRisk      bool              `json:"notify_on_risk"`
	NotifyOnSuccess   bool              `json:"notify_on_success"`
	SeverityThreshold ViolationSeverity `json:"severity_threshold"`
}

// ReportingSettings representa configurações de relatórios
type ReportingSettings struct {
	AutoGenerate           bool     `json:"auto_generate"`
	Recipients             []string `json:"recipients"`
	Frequency              string   `json:"frequency"` // "DAILY", "WEEKLY", "MONTHLY"
	Format                 string   `json:"format"`    // "PDF", "HTML", "JSON"
	IncludeEvidence        bool     `json:"include_evidence"`
	IncludeRecommendations bool     `json:"include_recommendations"`
	Template               string   `json:"template,omitempty"`
}

// RetentionPolicy representa política de retenção
type RetentionPolicy struct {
	AuditEntries      int `json:"audit_entries_days"`      // Dias para manter entradas de auditoria
	ComplianceResults int `json:"compliance_results_days"` // Dias para manter resultados
	Evidence          int `json:"evidence_days"`           // Dias para manter evidências
	Reports           int `json:"reports_days"`            // Dias para manter relatórios
	Violations        int `json:"violations_days"`         // Dias para manter violações resolvidas
}
