# VRAgent - Agente de IA para Processamento de VR

Este README documenta o uso do VRAgent, o agente inteligente para automação de processamento de Vale Refeição (VR) baseado em LangChainGo.

## 📋 Visão Geral

O VRAgent é um agente de IA que oferece:
- **Memory Management**: Contexto mantido entre conversações
- **Tool Integration**: Ferramentas especializadas para Excel, Cálculo e Validação
- **Workflow Automation**: Execução automatizada de processos complexos
- **Chat Integration**: Integração transparente com sistema de chat existente

## 🚀 Uso Básico

### Criar um Agente

```go
import (
    "BrxAgente-desafio4/internal/agent"
    "BrxAgente-desafio4/internal/chat"
    "BrxAgente-desafio4/internal/config"
)

// Configurar chat service
cfg := &config.Config{
    OpenAIKey: "sua-chave-openai",
    OllamaConfig: config.OllamaConfig{
        BaseURL: "http://localhost:11434",
        Model:   "llama2",
    },
}
chatSvc := chat.NewChat(cfg)

// Criar agente com configuração padrão
agent, err := agent.NewVRAgent(nil, chatSvc)
if err != nil {
    log.Fatalf("Erro ao criar agente: %v", err)
}
```

### Configuração Personalizada

```go
// Configuração customizada
agentConfig := &agent.AgentConfig{
    Enabled:          true,
    Model:           "gpt-4",
    Temperature:     0.7,
    MaxTokens:       4000,
    Timeout:         60 * time.Second,
    MemorySize:      200,
    MemoryTTL:       24 * time.Hour,
    ContextWindow:   8000,
    MaxMemoryTokens: 4000,
    WorkerPoolSize:  4,
    CacheEnabled:    true,
    CacheSize:       1000,
    CacheTTL:        2 * time.Hour,
    LogLevel:        "debug",
    DebugMode:       true,
    ToolsEnabled:    []string{"excel", "calculation", "validation"},
}

agent, err := agent.NewVRAgent(agentConfig, chatSvc)
```

## 🛠️ Ferramentas Disponíveis

### 1. ReadExcel Tool
Lê planilhas Excel e retorna dados estruturados.

```go
// Exemplo de uso
result, err := agent.ExecuteTool("read_excel", "files/ATIVOS.xlsx")
if err != nil {
    log.Printf("Erro ao ler Excel: %v", err)
} else {
    log.Printf("Dados lidos: %s", result)
}
```

**Input esperado:**
```json
{
    "file_path": "caminho/para/arquivo.xlsx",
    "sheet_name": "planilha1",
    "header_row": 1
}
```

**Output:**
```json
{
    "success": true,
    "row_count": 150,
    "columns": ["Matricula", "Sindicato", "Empresa"],
    "sample_data": [...]
}
```

### 2. CalculateVR Tool
Calcula valores de VR para colaboradores.

```go
input := `{
    "colaborador": {
        "matricula": "12345",
        "sindicato": "SINDPD"
    },
    "valor_por_sindicato": {"SINDPD": 21.25},
    "dias_uteis_por_sindicato": {"SINDPD": 22},
    "mes_referencia": "2025-09"
}`

result, err := agent.ExecuteTool("calculate_vr", input)
```

**Output:**
```json
{
    "success": true,
    "valor_total": 467.50,
    "valor_empresa": 374.00,
    "valor_colaborador": 93.50,
    "dias_efetivos": 22
}
```

### 3. ValidateData Tool
Valida dados de colaboradores e planilhas.

```go
input := `{
    "tipo_validacao": "colaborador",
    "colaborador": {
        "matricula": "12345",
        "sindicato": "SINDPD",
        "empresa": "Empresa ABC"
    },
    "validar_campos_obrigatorios": true
}`

result, err := agent.ExecuteTool("validate_data", input)
```

## 📊 Workflows Automatizados

### Workflow de Processamento VR Mensal

```go
// Executar workflow completo
err := agent.ExecuteWorkflow("processar-vr-mensal")
if err != nil {
    log.Printf("Erro no workflow: %v", err)
} else {
    log.Println("Processamento VR concluído com sucesso!")
}
```

**Etapas do workflow:**
1. Validação de dados
2. Cálculo de VR
3. Geração de relatórios
4. Notificação de resultados

### Workflow de Validação

```go
err := agent.ExecuteWorkflow("validar-dados")
```

**Etapas:**
1. Verificação de arquivos Excel
2. Validação de dados de colaboradores  
3. Verificação de datas
4. Análise de inconsistências

## 🧠 Memory Management

### Obter Histórico de Conversação

```go
memory, err := agent.GetMemory()
if err != nil {
    log.Printf("Erro ao obter memória: %v", err)
} else {
    for i, item := range memory {
        log.Printf("Memória %d: %s", i, item)
    }
}
```

### Limpar Memória

```go
// Limpar apenas a memória
err := agent.ClearMemory()

// Ou reset completo (memória + estatísticas)
err := agent.Reset()
```

## 💬 Integração com Chat

### Configurar Agente no Chat

```go
// Configurar agente no chat service
chatSvc.SetAgent(agent)

// Fazer perguntas via chat (com fallback automático)
response, err := chatSvc.Ask(
    "Quantos colaboradores ativos temos?", 
    "Você é um assistente de RH especializado", 
    []chat.Message{},
)
```

### Fazer Perguntas Diretamente ao Agente

```go
response, err := agent.Ask("Verifique se há erros na planilha ATIVOS.xlsx")
if err != nil {
    log.Printf("Erro: %v", err)
} else {
    log.Printf("Resposta: %s", response)
}
```

## 📈 Monitoramento e Status

### Verificar Status do Agente

```go
status := agent.GetStatus()
fmt.Printf("Estado: %s\n", status.State)
fmt.Printf("Uptime: %v\n", status.Uptime)
fmt.Printf("Total Requests: %d\n", status.TotalRequests)
fmt.Printf("Error Count: %d\n", status.ErrorCount)
fmt.Printf("Tarefa Atual: %s\n", status.CurrentTask)
```

### Listar Ferramentas Disponíveis

```go
tools := agent.GetAvailableTools()
fmt.Printf("Ferramentas: %v\n", tools)

// Obter informações detalhadas
info, err := agent.GetToolInfo("read_excel")
allInfo := agent.GetAllToolsInfo()
```

## 🎯 Exemplos Práticos

### Exemplo 1: Consulta Básica

```go
// Pergunta que requer uso de ferramenta
response, err := agent.Ask("Quantos colaboradores ativos temos na planilha?")
// O agente automaticamente:
// 1. Identifica que precisa ler Excel
// 2. Usa ReadExcelTool
// 3. Conta os registros
// 4. Retorna resposta em linguagem natural
```

### Exemplo 2: Validação de Dados

```go
response, err := agent.Ask("Verifique se há erros na planilha ATIVOS.xlsx")
// O agente:
// 1. Usa ReadExcelTool para ler arquivo
// 2. Usa ValidateDataTool para verificar dados
// 3. Relata problemas encontrados
```

### Exemplo 3: Cálculo de VR

```go
response, err := agent.Ask("Calcule o VR para o colaborador matrícula 12345")
// O agente:
// 1. Busca dados do colaborador
// 2. Obtém parâmetros de cálculo
// 3. Usa CalculateVRTool
// 4. Retorna valor calculado
```

## 🔧 Troubleshooting

### Problemas Comuns

1. **Agente não responde**: Verificar se está habilitado
```go
if !agent.IsEnabled() {
    agent.Enable()
}
```

2. **Ferramentas não funcionam**: Verificar configuração
```go
tools := agent.GetAvailableTools()
if len(tools) == 0 {
    log.Println("Nenhuma ferramenta disponível")
}
```

3. **Erros de LLM**: Verificar configuração do chat service
```go
// Verificar se OpenAI ou Ollama está configurado
if cfg.OpenAIKey == "" && cfg.OllamaConfig.BaseURL == "" {
    log.Println("Nenhum LLM configurado")
}
```

### Debug Mode

```go
// Habilitar modo debug
agentConfig.DebugMode = true
agentConfig.LogLevel = "debug"

// Logs detalhados serão exibidos
```

## 🧪 Executar Testes

```bash
# Testes básicos
go test ./internal/agent/... -v

# Testes de integração
go test ./internal/agent/... -tags=integration

# Benchmark de performance  
go test ./internal/agent/... -bench=.

# Pular testes demorados
go test ./internal/agent/... -short
```

## 📝 Configurações de Ambiente

### Variáveis de Ambiente Suportadas

```bash
# OpenAI
export OPENAI_API_KEY="sua-chave"

# Ollama
export OLLAMA_BASE_URL="http://localhost:11434"
export OLLAMA_MODEL="llama2"

# Agente
export AGENT_DEBUG="true"
export AGENT_LOG_LEVEL="debug"
export AGENT_MEMORY_SIZE="200"
```

### Arquivo de Configuração

Crie um arquivo `agent_config.json`:

```json
{
    "enabled": true,
    "model": "gpt-3.5-turbo",
    "temperature": 0.7,
    "max_tokens": 2000,
    "memory_size": 100,
    "memory_ttl": "24h",
    "cache_enabled": true,
    "debug_mode": false
}
```

## 📚 Próximos Passos

1. **Workflows Avançados**: Implementação de workflows mais complexos
2. **Análise Preditiva**: Detecção de padrões e anomalias
3. **Relatórios Automáticos**: Geração inteligente de insights
4. **Integração com APIs**: Conectores para sistemas externos

Para mais informações, consulte o PRD Agente IA e a documentação das issues relacionadas.