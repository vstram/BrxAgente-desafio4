package intelligence

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ReportFormat define os formatos de relatório disponíveis
type ReportFormat string

const (
	FormatExcel ReportFormat = "excel"
	FormatPDF   ReportFormat = "pdf"
	FormatJSON  ReportFormat = "json"
	FormatHTML  ReportFormat = "html"
	FormatText  ReportFormat = "text"
)

// ReportType define os tipos de relatório que podem ser gerados
type ReportType string

const (
	ReportTypeExecutive ReportType = "executive"
	ReportTypeDetailed  ReportType = "detailed"
	ReportTypeAnomaly   ReportType = "anomaly"
	ReportTypeInsight   ReportType = "insight"
)

// ReportTemplate define a estrutura de um template de relatório
type ReportTemplate struct {
	Name        string                 `json:"name"`
	Type        ReportType             `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Sections    []ReportSection        `json:"sections"`
	Style       ReportStyle            `json:"style"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ReportSection define uma seção do relatório
type ReportSection struct {
	Name       string                 `json:"name"`
	Title      string                 `json:"title"`
	Type       string                 `json:"type"` // summary, table, chart, text, insights
	Priority   int                    `json:"priority"`
	Visible    bool                   `json:"visible"`
	Data       map[string]interface{} `json:"data"`
	Formatting map[string]interface{} `json:"formatting"`
}

// ReportStyle define o estilo visual do relatório
type ReportStyle struct {
	Theme       string            `json:"theme"` // light, dark, corporate
	ColorScheme []string          `json:"color_scheme"`
	FontFamily  string            `json:"font_family"`
	FontSize    int               `json:"font_size"`
	HeaderStyle map[string]string `json:"header_style"`
	TableStyle  map[string]string `json:"table_style"`
	ChartStyle  map[string]string `json:"chart_style"`
}

// Report representa um relatório gerado
type Report struct {
	ID          string                 `json:"id"`
	Type        ReportType             `json:"type"`
	Format      ReportFormat           `json:"format"`
	Title       string                 `json:"title"`
	GeneratedAt time.Time              `json:"generated_at"`
	GeneratedBy string                 `json:"generated_by"`
	Data        *ProcessingData        `json:"data"`
	Insights    []*Insight             `json:"insights"`
	Metadata    map[string]interface{} `json:"metadata"`
	Content     interface{}            `json:"content"`
	Template    *ReportTemplate        `json:"template,omitempty"`
}

// ReportGenerator gera relatórios em múltiplos formatos
type ReportGenerator struct {
	config           *ReportConfig
	insightGenerator *InsightGenerator
	templates        map[string]*ReportTemplate
	logger           *log.Logger
}

// ReportConfig configurações para geração de relatórios
type ReportConfig struct {
	OutputPath      string            `json:"output_path"`
	DefaultFormat   ReportFormat      `json:"default_format"`
	EnableInsights  bool              `json:"enable_insights"`
	CustomTemplates []string          `json:"custom_templates"`
	DefaultStyle    ReportStyle       `json:"default_style"`
	MaxReportSize   int64             `json:"max_report_size"`
	CompanyInfo     map[string]string `json:"company_info"`
}

// DefaultReportConfig retorna configuração padrão para relatórios
func DefaultReportConfig() *ReportConfig {
	return &ReportConfig{
		OutputPath:     "./reports",
		DefaultFormat:  FormatExcel,
		EnableInsights: true,
		MaxReportSize:  50 * 1024 * 1024, // 50MB
		CompanyInfo: map[string]string{
			"name":    "Empresa",
			"address": "",
			"phone":   "",
			"email":   "",
		},
		DefaultStyle: ReportStyle{
			Theme:       "corporate",
			ColorScheme: []string{"#2E86AB", "#A23B72", "#F18F01", "#C73E1D"},
			FontFamily:  "Arial",
			FontSize:    11,
			HeaderStyle: map[string]string{
				"background": "#2E86AB",
				"color":      "white",
				"bold":       "true",
			},
			TableStyle: map[string]string{
				"border":          "thin",
				"header_bg":       "#F0F8FF",
				"alternate_rows":  "true",
				"alternate_color": "#F9F9F9",
			},
		},
	}
}

// NewReportGenerator cria novo gerador de relatórios
func NewReportGenerator(config *ReportConfig, insightGenerator *InsightGenerator) *ReportGenerator {
	if config == nil {
		config = DefaultReportConfig()
	}

	if insightGenerator == nil {
		insightGenerator = NewInsightGenerator(nil)
	}

	generator := &ReportGenerator{
		config:           config,
		insightGenerator: insightGenerator,
		templates:        make(map[string]*ReportTemplate),
		logger:           log.Default(),
	}

	// Carregar templates padrão
	generator.loadDefaultTemplates()

	return generator
}

// GenerateExecutiveReport gera relatório executivo
func (g *ReportGenerator) GenerateExecutiveReport(data *ProcessingData, format ReportFormat) (*Report, error) {
	g.logger.Printf("Gerando relatório executivo em formato %s", format)

	// Gerar insights automáticos
	insights, err := g.insightGenerator.GenerateInsights(data)
	if err != nil {
		g.logger.Printf("Erro ao gerar insights: %v", err)
		insights = []*Insight{} // Continuar sem insights
	}

	// Criar relatório base
	report := &Report{
		ID:          generateReportID("exec", time.Now()),
		Type:        ReportTypeExecutive,
		Format:      format,
		Title:       "Relatório Executivo - Processamento de VR",
		GeneratedAt: time.Now(),
		GeneratedBy: "BrxAgente IA",
		Data:        data,
		Insights:    insights,
		Metadata: map[string]interface{}{
			"version":        "1.0.0",
			"template":       "executive_default",
			"auto_generated": true,
		},
		Template: g.templates["executive"],
	}

	// Gerar conteúdo baseado no formato
	switch format {
	case FormatExcel:
		content, err := g.generateExecutiveExcel(report)
		if err != nil {
			return nil, fmt.Errorf("erro ao gerar Excel executivo: %w", err)
		}
		report.Content = content

	case FormatJSON:
		content, err := g.generateExecutiveJSON(report)
		if err != nil {
			return nil, fmt.Errorf("erro ao gerar JSON executivo: %w", err)
		}
		report.Content = content

	case FormatHTML:
		content, err := g.generateExecutiveHTML(report)
		if err != nil {
			return nil, fmt.Errorf("erro ao gerar HTML executivo: %w", err)
		}
		report.Content = content

	case FormatText:
		content := g.generateExecutiveText(report)
		report.Content = content

	default:
		return nil, fmt.Errorf("formato de relatório não suportado: %s", format)
	}

	g.logger.Printf("Relatório executivo gerado com sucesso (formato: %s)", format)
	return report, nil
}

// GenerateDetailedReport gera relatório detalhado
func (g *ReportGenerator) GenerateDetailedReport(data *ProcessingData, format ReportFormat) (*Report, error) {
	g.logger.Printf("Gerando relatório detalhado em formato %s", format)

	// Gerar insights automáticos
	insights, err := g.insightGenerator.GenerateInsights(data)
	if err != nil {
		g.logger.Printf("Erro ao gerar insights: %v", err)
		insights = []*Insight{}
	}

	report := &Report{
		ID:          generateReportID("det", time.Now()),
		Type:        ReportTypeDetailed,
		Format:      format,
		Title:       "Relatório Detalhado - Processamento de VR",
		GeneratedAt: time.Now(),
		GeneratedBy: "BrxAgente IA",
		Data:        data,
		Insights:    insights,
		Metadata: map[string]interface{}{
			"version":          "1.0.0",
			"template":         "detailed_default",
			"auto_generated":   true,
			"include_raw_data": true,
		},
		Template: g.templates["detailed"],
	}

	// Gerar conteúdo baseado no formato
	switch format {
	case FormatExcel:
		content, err := g.generateDetailedExcel(report)
		if err != nil {
			return nil, fmt.Errorf("erro ao gerar Excel detalhado: %w", err)
		}
		report.Content = content

	case FormatJSON:
		content, err := g.generateDetailedJSON(report)
		if err != nil {
			return nil, fmt.Errorf("erro ao gerar JSON detalhado: %w", err)
		}
		report.Content = content

	default:
		return nil, fmt.Errorf("formato '%s' não implementado para relatório detalhado", format)
	}

	g.logger.Printf("Relatório detalhado gerado com sucesso")
	return report, nil
}

// GenerateAnomalyReport gera relatório de anomalias
func (g *ReportGenerator) GenerateAnomalyReport(data *ProcessingData, format ReportFormat) (*Report, error) {
	g.logger.Printf("Gerando relatório de anomalias em formato %s", format)

	if data.AnomalyReport == nil {
		return nil, fmt.Errorf("dados de anomalias não disponíveis")
	}

	report := &Report{
		ID:          generateReportID("anom", time.Now()),
		Type:        ReportTypeAnomaly,
		Format:      format,
		Title:       "Relatório de Anomalias - Detecção Automática",
		GeneratedAt: time.Now(),
		GeneratedBy: "BrxAgente IA",
		Data:        data,
		Insights:    []*Insight{}, // Anomalias são os próprios insights aqui
		Metadata: map[string]interface{}{
			"version":         "1.0.0",
			"template":        "anomaly_default",
			"focus":           "anomalies",
			"total_anomalies": data.AnomalyReport.TotalAnomalies,
		},
		Template: g.templates["anomaly"],
	}

	// Gerar conteúdo baseado no formato
	switch format {
	case FormatText:
		content := FormatAnomalyReportForHuman(data.AnomalyReport)
		report.Content = content

	case FormatJSON:
		content, err := json.MarshalIndent(data.AnomalyReport, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("erro ao serializar anomalias para JSON: %w", err)
		}
		report.Content = string(content)

	default:
		return nil, fmt.Errorf("formato '%s' não implementado para relatório de anomalias", format)
	}

	g.logger.Printf("Relatório de anomalias gerado com sucesso")
	return report, nil
}

// generateExecutiveExcel gera relatório executivo em Excel
func (g *ReportGenerator) generateExecutiveExcel(report *Report) (string, error) {
	file := excelize.NewFile()
	defer file.Close()

	// Configurar primeira planilha
	sheetName := "Resumo Executivo"
	file.SetSheetName("Sheet1", sheetName)

	// Estilo do cabeçalho
	headerStyle, err := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{g.config.DefaultStyle.HeaderStyle["background"]}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err != nil {
		return "", fmt.Errorf("erro ao criar estilo de cabeçalho: %w", err)
	}

	// Título principal
	file.SetCellValue(sheetName, "A1", report.Title)
	file.SetCellStyle(sheetName, "A1", "A1", headerStyle)
	file.MergeCell(sheetName, "A1", "D1")

	// Informações básicas
	row := 3
	file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Data do Relatório:")
	file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.GeneratedAt.Format("02/01/2006 15:04"))
	row++

	file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Colaboradores Processados:")
	file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), report.Data.TotalCollaborators)
	row++

	file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Valor Total de VR:")
	file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("R$ %.2f", report.Data.TotalVRValue))
	row++

	file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Tempo de Processamento:")
	file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("%.1f minutos", report.Data.ProcessingTime.Minutes()))
	row++

	// Seção de insights
	if len(report.Insights) > 0 {
		row += 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "INSIGHTS PRINCIPAIS")
		file.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), headerStyle)
		row++

		for i, insight := range report.Insights {
			if i >= 10 { // Limitar a 10 insights no Excel
				break
			}
			file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("%d. %s", i+1, insight.Title))
			file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), insight.Description)
			file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), insight.Priority.String())
			file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("%.0f%%", insight.Confidence*100))
			row++
		}
	}

	// Salvar arquivo temporário
	fileName := fmt.Sprintf("relatorio_executivo_%s.xlsx", time.Now().Format("20060102_150405"))
	filePath := filepath.Join(g.config.OutputPath, fileName)

	if err := file.SaveAs(filePath); err != nil {
		return "", fmt.Errorf("erro ao salvar arquivo Excel: %w", err)
	}

	return filePath, nil
}

// generateExecutiveJSON gera relatório executivo em JSON
func (g *ReportGenerator) generateExecutiveJSON(report *Report) (string, error) {
	executiveSummary := map[string]interface{}{
		"metadata": map[string]interface{}{
			"report_id":    report.ID,
			"title":        report.Title,
			"generated_at": report.GeneratedAt,
			"generated_by": report.GeneratedBy,
			"report_type":  string(report.Type),
			"format":       string(report.Format),
		},
		"summary": map[string]interface{}{
			"total_collaborators": report.Data.TotalCollaborators,
			"total_vr_value":      report.Data.TotalVRValue,
			"processing_time":     report.Data.ProcessingTime.String(),
			"error_count":         report.Data.ErrorCount,
			"warning_count":       report.Data.WarningCount,
		},
		"key_metrics": map[string]interface{}{
			"average_vr_per_employee": report.Data.TotalVRValue / float64(report.Data.TotalCollaborators),
			"processing_efficiency":   float64(report.Data.TotalCollaborators) / report.Data.ProcessingTime.Minutes(),
			"data_quality_score":      calculateQualityScore(report.Data),
		},
		"insights":       report.Insights,
		"sindicato_data": report.Data.SindicatoDistribution,
	}

	// Adicionar dados de anomalias se disponível
	if report.Data.AnomalyReport != nil {
		executiveSummary["anomaly_summary"] = map[string]interface{}{
			"total_anomalies": report.Data.AnomalyReport.TotalAnomalies,
			"overall_score":   report.Data.AnomalyReport.Summary.OverallScore,
			"risk_level":      report.Data.AnomalyReport.Summary.RiskLevel,
			"critical_issues": report.Data.AnomalyReport.Summary.CriticalIssues,
		}
	}

	jsonBytes, err := json.MarshalIndent(executiveSummary, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao serializar JSON: %w", err)
	}

	return string(jsonBytes), nil
}

// generateExecutiveHTML gera relatório executivo em HTML
func (g *ReportGenerator) generateExecutiveHTML(report *Report) (string, error) {
	htmlTemplate := `
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { background-color: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #2E86AB; border-bottom: 2px solid #2E86AB; padding-bottom: 10px; }
        h2 { color: #A23B72; margin-top: 30px; }
        .metric-card { background-color: #f8f9fa; padding: 15px; margin: 10px 0; border-radius: 5px; border-left: 4px solid #2E86AB; }
        .insight { background-color: #fff3cd; padding: 15px; margin: 10px 0; border-radius: 5px; border-left: 4px solid #F18F01; }
        .priority-critical { border-left-color: #C73E1D; }
        .priority-high { border-left-color: #F18F01; }
        .priority-medium { border-left-color: #A23B72; }
        .priority-low { border-left-color: #2E86AB; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #f8f9fa; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <h1>{{.Title}}</h1>
        <p><strong>Gerado em:</strong> {{.GeneratedAt.Format "02/01/2006 15:04"}}</p>
        
        <h2>Resumo Executivo</h2>
        <div class="metric-card">
            <h3>Métricas Principais</h3>
            <p><strong>Colaboradores Processados:</strong> {{.Data.TotalCollaborators}}</p>
            <p><strong>Valor Total de VR:</strong> R$ {{printf "%.2f" .Data.TotalVRValue}}</p>
            <p><strong>Tempo de Processamento:</strong> {{printf "%.1f" .Data.ProcessingTime.Minutes}} minutos</p>
            <p><strong>Erros Encontrados:</strong> {{.Data.ErrorCount}}</p>
        </div>
        
        {{if .Insights}}
        <h2>Insights Automáticos</h2>
        {{range .Insights}}
        <div class="insight priority-{{.Priority.String}}">
            <h4>{{.Title}}</h4>
            <p><strong>Descrição:</strong> {{.Description}}</p>
            {{if .Impact}}<p><strong>Impacto:</strong> {{.Impact}}</p>{{end}}
            {{if .Action}}<p><strong>Ação Recomendada:</strong> {{.Action}}</p>{{end}}
            <p><strong>Confiança:</strong> {{printf "%.0f" (.Confidence | multiply 100)}}% | <strong>Prioridade:</strong> {{.Priority.String}}</p>
        </div>
        {{end}}
        {{end}}
        
        <h2>Distribuição por Sindicato</h2>
        <table>
            <thead>
                <tr><th>Sindicato</th><th>Colaboradores</th><th>Percentual</th></tr>
            </thead>
            <tbody>
                {{range $sindicato, $count := .Data.SindicatoDistribution}}
                <tr>
                    <td>{{$sindicato}}</td>
                    <td>{{$count}}</td>
                    <td>{{printf "%.1f" (divide $count $.Data.TotalCollaborators | multiply 100)}}%</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>
</body>
</html>
`

	// Funções auxiliares para o template
	funcMap := template.FuncMap{
		"multiply": func(a, b float64) float64 { return a * b },
		"divide": func(a, b int) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
	}

	tmpl, err := template.New("executive").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("erro ao criar template HTML: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, report); err != nil {
		return "", fmt.Errorf("erro ao executar template HTML: %w", err)
	}

	return buf.String(), nil
}

// generateExecutiveText gera relatório executivo em texto
func (g *ReportGenerator) generateExecutiveText(report *Report) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("📊 %s\n", strings.ToUpper(report.Title)))
	builder.WriteString(strings.Repeat("=", len(report.Title)+4))
	builder.WriteString("\n\n")

	// Informações básicas
	builder.WriteString("📅 INFORMAÇÕES BÁSICAS\n")
	builder.WriteString("─────────────────────\n")
	builder.WriteString(fmt.Sprintf("Data: %s\n", report.GeneratedAt.Format("02/01/2006 15:04")))
	builder.WriteString(fmt.Sprintf("Colaboradores: %d\n", report.Data.TotalCollaborators))
	builder.WriteString(fmt.Sprintf("Valor Total: R$ %.2f\n", report.Data.TotalVRValue))
	builder.WriteString(fmt.Sprintf("Tempo: %.1f minutos\n", report.Data.ProcessingTime.Minutes()))
	builder.WriteString(fmt.Sprintf("Erros: %d\n\n", report.Data.ErrorCount))

	// Métricas calculadas
	avgVR := report.Data.TotalVRValue / float64(report.Data.TotalCollaborators)
	builder.WriteString("📈 MÉTRICAS CALCULADAS\n")
	builder.WriteString("──────────────────────\n")
	builder.WriteString(fmt.Sprintf("VR Médio por Colaborador: R$ %.2f\n", avgVR))
	builder.WriteString(fmt.Sprintf("Eficiência: %.1f colaboradores/minuto\n",
		float64(report.Data.TotalCollaborators)/report.Data.ProcessingTime.Minutes()))
	builder.WriteString(fmt.Sprintf("Score de Qualidade: %.1f%%\n\n", calculateQualityScore(report.Data)))

	// Insights
	if len(report.Insights) > 0 {
		builder.WriteString(FormatInsightsForHuman(report.Insights))
	}

	// Distribuição por sindicato
	if len(report.Data.SindicatoDistribution) > 0 {
		builder.WriteString("\n👥 DISTRIBUIÇÃO POR SINDICATO\n")
		builder.WriteString("─────────────────────────────\n")

		// Ordenar por quantidade
		type SindicatoCount struct {
			Nome  string
			Count int
		}

		var sindicatos []SindicatoCount
		for nome, count := range report.Data.SindicatoDistribution {
			sindicatos = append(sindicatos, SindicatoCount{Nome: nome, Count: count})
		}

		sort.Slice(sindicatos, func(i, j int) bool {
			return sindicatos[i].Count > sindicatos[j].Count
		})

		for _, s := range sindicatos {
			percentage := float64(s.Count) / float64(report.Data.TotalCollaborators) * 100
			builder.WriteString(fmt.Sprintf("%-20s: %4d colaboradores (%.1f%%)\n",
				s.Nome, s.Count, percentage))
		}
	}

	return builder.String()
}

// generateDetailedExcel gera relatório detalhado em Excel
func (g *ReportGenerator) generateDetailedExcel(report *Report) (string, error) {
	// Por simplicidade, usar mesmo formato do executivo mas com mais detalhes
	return g.generateExecutiveExcel(report)
}

// generateDetailedJSON gera relatório detalhado em JSON
func (g *ReportGenerator) generateDetailedJSON(report *Report) (string, error) {
	// Incluir todos os dados disponíveis
	detailedData := map[string]interface{}{
		"metadata":        report.Metadata,
		"processing_data": report.Data,
		"insights":        report.Insights,
		"generation_info": map[string]interface{}{
			"generated_at": report.GeneratedAt,
			"generated_by": report.GeneratedBy,
			"report_id":    report.ID,
		},
	}

	jsonBytes, err := json.MarshalIndent(detailedData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao serializar JSON detalhado: %w", err)
	}

	return string(jsonBytes), nil
}

// loadDefaultTemplates carrega templates padrão
func (g *ReportGenerator) loadDefaultTemplates() {
	// Template executivo
	g.templates["executive"] = &ReportTemplate{
		Name:        "executive_default",
		Type:        ReportTypeExecutive,
		Title:       "Relatório Executivo",
		Description: "Resumo executivo com métricas principais e insights",
		Sections: []ReportSection{
			{Name: "summary", Title: "Resumo", Type: "summary", Priority: 1, Visible: true},
			{Name: "metrics", Title: "Métricas", Type: "table", Priority: 2, Visible: true},
			{Name: "insights", Title: "Insights", Type: "insights", Priority: 3, Visible: true},
		},
		Style: g.config.DefaultStyle,
	}

	// Template detalhado
	g.templates["detailed"] = &ReportTemplate{
		Name:        "detailed_default",
		Type:        ReportTypeDetailed,
		Title:       "Relatório Detalhado",
		Description: "Relatório completo com todos os dados e análises",
		Sections: []ReportSection{
			{Name: "summary", Title: "Resumo", Type: "summary", Priority: 1, Visible: true},
			{Name: "raw_data", Title: "Dados Brutos", Type: "table", Priority: 2, Visible: true},
			{Name: "analysis", Title: "Análise", Type: "text", Priority: 3, Visible: true},
			{Name: "insights", Title: "Insights", Type: "insights", Priority: 4, Visible: true},
		},
		Style: g.config.DefaultStyle,
	}

	// Template de anomalias
	g.templates["anomaly"] = &ReportTemplate{
		Name:        "anomaly_default",
		Type:        ReportTypeAnomaly,
		Title:       "Relatório de Anomalias",
		Description: "Relatório focado em anomalias e problemas detectados",
		Sections: []ReportSection{
			{Name: "anomaly_summary", Title: "Resumo de Anomalias", Type: "summary", Priority: 1, Visible: true},
			{Name: "anomaly_list", Title: "Lista de Anomalias", Type: "table", Priority: 2, Visible: true},
			{Name: "recommendations", Title: "Recomendações", Type: "text", Priority: 3, Visible: true},
		},
		Style: g.config.DefaultStyle,
	}
}

// Helper functions
func generateReportID(prefix string, timestamp time.Time) string {
	return fmt.Sprintf("%s_%s", prefix, timestamp.Format("20060102_150405"))
}

func calculateQualityScore(data *ProcessingData) float64 {
	if data.TotalCollaborators == 0 {
		return 100.0
	}

	totalIssues := data.ErrorCount + data.WarningCount
	issueRate := float64(totalIssues) / float64(data.TotalCollaborators)
	return 100.0 - (issueRate * 100)
}
