# 🔧 Referência da API - Agente de IA BrxAgente

Esta documentação detalha todos os métodos disponíveis no backend Go da aplicação Wails, que são expostos para o frontend React.

## 📑 Índice

1. [Arquitetura da API](#-arquitetura-da-api)
2. [Métodos de Configuração](#️-métodos-de-configuração)
3. [Métodos de Processamento](#-métodos-de-processamento)
4. [Métodos do Chat](#-métodos-do-chat)
5. [Métodos de Monitoramento](#-métodos-de-monitoramento)
6. [Métodos de Análise Preditiva](#-métodos-de-análise-preditiva)
7. [Tipos de Dados](#-tipos-de-dados)
8. [Tratamento de Erros](#-tratamento-de-erros)

## 🏗️ Arquitetura da API

A aplicação BrxAgente usa **Wails v2** para expor métodos Go para o frontend React através de bindings automáticos.

### **Como Funciona:**
1. **Backend Go**: Métodos definidos na struct `App`
2. **Auto-binding**: Wails expõe automaticamente os métodos
3. **Frontend React**: Chama métodos via `window.go.main.App.NomeDoMetodo()`
4. **TypeScript**: Tipos gerados automaticamente pelo Wails

### **Estrutura Principal:**
```go
// app.go - Struct principal
type App struct {
    ctx                  context.Context
    cfg                  *config.Config
    chat                 *chat.Chat
    colaboradores        map[string]*modelo.Colaborador
    vrAgent              *agent.VRAgent
    agentStatus          string
    currentWorkflow      *WorkflowInfo
    
    // Intelligence systems
    trendPredictor       *intelligence.TrendPredictor
    patternAnalyzer      *intelligence.PatternAnalyzer
    forecaster           *intelligence.Forecaster
    recommendationEngine *intelligence.RecommendationEngine
    historicalData       []predicoes.HistoricalVRData
}
```

## ⚙️ Métodos de Configuração

### **SetOpenAIKey**
Define a chave da API do OpenAI.

```go
func (a *App) SetOpenAIKey(key string) error
```

**Parâmetros:**
- `key` (string): Chave da API OpenAI (formato: sk-...)

**Retorno:**
- `error`: Erro se a chave for inválida

**Exemplo (Frontend React):**
```typescript
try {
  await window.go.main.App.SetOpenAIKey('sk-proj-...');
  console.log('Chave configurada com sucesso');
} catch (error) {
  console.error('Erro ao configurar chave:', error);
}
```

---

### **SetOllamaConfig**
Configura conexão com Ollama local.

```go
func (a *App) SetOllamaConfig(ollamaConfig config.OllamaConfig) error
```

**Parâmetros:**
- `ollamaConfig` (OllamaConfig): Configuração do Ollama
  - `URL` (string): URL do servidor Ollama
  - `Model` (string): Nome do modelo
  - `Timeout` (int): Timeout em segundos

**Exemplo (Frontend):**
```typescript
const ollamaConfig = {
  URL: 'http://localhost:11434',
  Model: 'llama2',
  Timeout: 30
};

await window.go.main.App.SetOllamaConfig(ollamaConfig);
```

---

### **TestOpenAIKey**
Testa a validade de uma chave OpenAI.

```go
func (a *App) TestOpenAIKey(key string) (bool, error)
```

**Retorno:**
- `bool`: true se a chave for válida
- `error`: Erro se houver problemas na validação

---

### **GetConfig**
Retorna a configuração atual da aplicação.

```go
func (a *App) GetConfig() (*config.Config, error)
```

**Retorno:**
- `*config.Config`: Objeto com todas as configurações

## 🔄 Métodos de Processamento

### **SetDiretorioPlanilhas**
Valida e define o diretório das planilhas Excel.

```go
func (a *App) SetDiretorioPlanilhas(caminho string) (bool, error)
```

**Parâmetros:**
- `caminho` (string): Caminho absoluto para o diretório

**Validações executadas:**
- ✅ Caminho existe
- ✅ É um diretório
- ✅ Contém arquivos .xlsx
- ✅ Permissões de leitura

**Exemplo:**
```typescript
const directory = '/Users/user/planilhas-vr';
const [isValid, error] = await window.go.main.App.SetDiretorioPlanilhas(directory);

if (isValid) {
  console.log('Diretório válido!');
} else {
  console.error('Diretório inválido:', error);
}
```

---

### **RealizarAnaliseOrquestrada**
Executa todo o processo de análise VR de forma automatizada.

```go
func (a *App) RealizarAnaliseOrquestrada(diretorioPlanilhas string) (string, error)
```

**Parâmetros:**
- `diretorioPlanilhas` (string): Caminho para o diretório das planilhas

**Processo executado:**
1. Validação do diretório
2. Extração do mês/ano da planilha DESLIGADOS
3. Consolidação das bases de dados
4. Aplicação das regras de negócio
5. Geração da planilha final
6. Configuração do contexto do chat

**Exemplo:**
```typescript
try {
  const resultado = await window.go.main.App.RealizarAnaliseOrquestrada('/path/to/planilhas');
  console.log(resultado); // "Análise concluída com sucesso! X colaboradores processados..."
} catch (error) {
  console.error('Erro no processamento:', error);
}
```

---

### **SelecionarDiretorio**
Abre diálogo nativo para seleção de diretório.

```go
func (a *App) SelecionarDiretorio() (string, error)
```

**Retorno:**
- `string`: Caminho do diretório selecionado
- `error`: Erro se usuário cancelar ou falha no diálogo

**Exemplo:**
```typescript
try {
  const selectedPath = await window.go.main.App.SelecionarDiretorio();
  console.log('Diretório selecionado:', selectedPath);
} catch (error) {
  console.log('Usuário cancelou a seleção');
}
```

---

### **TestExcelReading**
Testa a leitura de um arquivo Excel de exemplo.

```go
func (a *App) TestExcelReading() (string, error)
```

**Retorno:**
- `string`: Informações sobre o arquivo lido
- `error`: Erro se houver problemas na leitura

## 💬 Métodos do Chat

### **AskAI**
Envia uma pergunta para o agente de IA e retorna a resposta.

```go
func (a *App) AskAI(question string) (string, error)
```

**Parâmetros:**
- `question` (string): Pergunta ou comando para o agente

**Exemplo:**
```typescript
const question = "Quantos colaboradores do SINDPD temos ativos?";
try {
  const response = await window.go.main.App.AskAI(question);
  console.log('Resposta do agente:', response);
} catch (error) {
  console.error('Erro na consulta:', error);
}
```

---

### **GetSystemPrompt**
Retorna o prompt do sistema com contexto atual.

```go
func (a *App) GetSystemPrompt() (string, error)
```

**Retorno:**
- `string`: Prompt do sistema completo com contexto

---

### **SetChatContext**
Define o contexto do chat com dados processados.

```go
func (a *App) SetChatContext() error
```

**Uso:**
Chame após processar dados para atualizar o contexto do chat.

## 📊 Métodos de Monitoramento

### **GetAgentStatus**
Retorna status completo do agente e métricas.

```go
func (a *App) GetAgentStatus() (*AgentStatus, error)
```

**Retorno:**
- `*AgentStatus`: Objeto com status completo

**Estrutura AgentStatus:**
```go
type AgentStatus struct {
    Status              string                `json:"status"`
    LastUpdated         time.Time             `json:"lastUpdated"`
    CurrentWorkflow     *WorkflowInfo         `json:"currentWorkflow"`
    AvailableWorkflows  []string              `json:"availableWorkflows"`
    Metrics             AgentMetrics          `json:"metrics"`
    RecentLogs          []LogEntry            `json:"recentLogs"`
}
```

**Exemplo:**
```typescript
const status = await window.go.main.App.GetAgentStatus();
console.log('Status do agente:', status.Status);
console.log('Workflows disponíveis:', status.AvailableWorkflows);
console.log('Uptime:', status.Metrics.Uptime + 'ms');
```

---

### **StartWorkflow**
Inicia execução de um workflow específico.

```go
func (a *App) StartWorkflow(request WorkflowStartRequest) error
```

**Parâmetros:**
- `request` (WorkflowStartRequest):
  - `WorkflowName` (string): Nome do workflow
  - `Parameters` (map[string]interface{}): Parâmetros

**Workflows Disponíveis:**
- `"analise-vr-mensal"`: Processamento completo
- `"validacao-planilhas"`: Apenas validação
- `"deteccao-anomalias"`: Detecção de anomalias
- `"geracao-relatorios"`: Geração de relatórios

**Exemplo:**
```typescript
const workflowRequest = {
  WorkflowName: 'analise-vr-mensal',
  Parameters: {
    directory: '/path/to/planilhas',
    sindicatos: ['SINDPD', 'SINDAC'],
    sendNotification: true
  }
};

await window.go.main.App.StartWorkflow(workflowRequest);
```

---

### **StopWorkflow**
Para graciosamente o workflow atual.

```go
func (a *App) StopWorkflow() error
```

---

### **CancelWorkflow**
Cancela forçadamente o workflow atual.

```go
func (a *App) CancelWorkflow() error
```

---

### **GetWorkflowHistory**
Retorna histórico de execuções de workflows.

```go
func (a *App) GetWorkflowHistory() ([]WorkflowExecution, error)
```

**Retorno:**
- `[]WorkflowExecution`: Array com histórico de execuções

---

### **ClearAgentLogs**
Limpa os logs do sistema.

```go
func (a *App) ClearAgentLogs() error
```

## 🔮 Métodos de Análise Preditiva

### **AddHistoricalData**
Adiciona dados históricos ao sistema de predição.

```go
func (a *App) AddHistoricalData(data predicoes.HistoricalVRData) error
```

**Parâmetros:**
- `data` (HistoricalVRData): Dados históricos estruturados

---

### **PredictTrends**
Executa análise de tendências para um sindicato.

```go
func (a *App) PredictTrends(sindicato string) (*predicoes.Prediction, error)
```

**Parâmetros:**
- `sindicato` (string): Nome do sindicato (ex: "SINDPD")

**Retorno:**
- `*predicoes.Prediction`: Objeto com predição e confiança

**Exemplo:**
```typescript
const prediction = await window.go.main.App.PredictTrends('SINDPD');
console.log('Tendência:', prediction.Type);
console.log('Confiança:', prediction.Confidence);
console.log('Descrição:', prediction.Description);
```

---

### **GenerateForecast**
Gera previsão de consumo para um sindicato.

```go
func (a *App) GenerateForecast(sindicato string, horizon int) (*predicoes.ConsumptionForecast, error)
```

**Parâmetros:**
- `sindicato` (string): Nome do sindicato
- `horizon` (int): Horizonte de previsão em meses

**Retorno:**
- `*predicoes.ConsumptionForecast`: Previsão estruturada

---

### **GetPredictiveAnalysisSummary**
Gera resumo completo de análise preditiva.

```go
func (a *App) GetPredictiveAnalysisSummary(sindicato string) (map[string]interface{}, error)
```

**Retorno:**
- `map[string]interface{}`: Objeto com análise completa
  - `trends`: Análise de tendências
  - `patterns`: Padrões identificados
  - `forecast`: Previsões
  - `recommendations`: Recomendações

## 📋 Tipos de Dados

### **WorkflowInfo**
```go
type WorkflowInfo struct {
    ID          string         `json:"id"`
    Name        string         `json:"name"`
    Status      string         `json:"status"`
    StartTime   time.Time      `json:"startTime"`
    EndTime     *time.Time     `json:"endTime"`
    Steps       []WorkflowStep `json:"steps"`
    Progress    float64        `json:"progress"`
    ErrorMsg    string         `json:"errorMsg"`
}
```

### **WorkflowStep**
```go
type WorkflowStep struct {
    ID          string     `json:"id"`
    Name        string     `json:"name"`
    Status      string     `json:"status"`
    StartTime   *time.Time `json:"startTime"`
    EndTime     *time.Time `json:"endTime"`
    Duration    int64      `json:"duration"`
    ErrorMsg    string     `json:"errorMsg"`
}
```

### **AgentMetrics**
```go
type AgentMetrics struct {
    TotalWorkflowsExecuted int   `json:"totalWorkflowsExecuted"`
    SuccessfulWorkflows    int   `json:"successfulWorkflows"`
    CollaboratorsProcessed int   `json:"collaboratorsProcessed"`
    ReportsGenerated       int   `json:"reportsGenerated"`
    AnomaliesDetected      int   `json:"anomaliesDetected"`
    Uptime                 int64 `json:"uptime"`
}
```

### **LogEntry**
```go
type LogEntry struct {
    ID        string    `json:"id"`
    Timestamp time.Time `json:"timestamp"`
    Level     string    `json:"level"`
    Message   string    `json:"message"`
    Source    string    `json:"source"`
}
```

### **Config.OllamaConfig**
```go
type OllamaConfig struct {
    URL     string `json:"url"`
    Model   string `json:"model"`
    Timeout int    `json:"timeout"`
}
```

### **HistoricalVRData**
```go
type HistoricalVRData struct {
    Month            time.Time              `json:"month"`
    Sindicato        string                 `json:"sindicato"`
    TotalVR          float64                `json:"totalVR"`
    NumColaboradores int                    `json:"numColaboradores"`
    MediaPorPessoa   float64                `json:"mediaPorPessoa"`
    DaysProcessed    int                    `json:"daysProcessed"`
    Anomalies        []string               `json:"anomalies"`
    Metadata         map[string]interface{} `json:"metadata"`
}
```

## ❌ Tratamento de Erros

### **Tipos de Erro Comuns**

**Configuração:**
- `"chave da API do OpenAI inválida"`: Chave malformada
- `"falha ao salvar a configuração"`: Problema de permissão
- `"configuração do Ollama inválida"`: URL ou modelo inválido

**Processamento:**
- `"diretório não encontrado"`: Caminho não existe
- `"nenhum arquivo .xlsx encontrado"`: Pasta sem planilhas
- `"erro ao acessar o diretório"`: Problemas de permissão

**Workflow:**
- `"agent is currently busy"`: Workflow já em execução
- `"no workflow is currently running"`: Tentativa de parar workflow inexistente
- `"Directory parameter not provided"`: Parâmetros incompletos

**Chat:**
- `"falha ao obter resposta da IA"`: Problema na API LLM
- `"nenhum dado de colaborador disponível"`: Context vazio

### **Padrão de Tratamento no Frontend**

```typescript
// Padrão recomendado para tratamento de erros
try {
  const result = await window.go.main.App.SomeMethod(params);
  // Sucesso
  handleSuccess(result);
} catch (error) {
  // Erro
  console.error('Erro na operação:', error);
  showErrorNotification(error.toString());
}

// Para métodos que retornam tuplas (bool, error)
const [success, error] = await window.go.main.App.SomeValidation(params);
if (!success) {
  console.error('Validação falhou:', error);
}
```

### **Códigos de Status de Workflow**

| Status | Descrição | Ação Recomendada |
|--------|-----------|-------------------|
| `"idle"` | Agente disponível | Pode iniciar workflows |
| `"running"` | Executando workflow | Aguardar conclusão |
| `"paused"` | Workflow pausado | Continuar ou cancelar |
| `"error"` | Erro na execução | Verificar logs, corrigir, reiniciar |
| `"completed"` | Concluído com sucesso | Verificar resultados |
| `"cancelled"` | Cancelado pelo usuário | Pode reiniciar se necessário |

---

## 🔗 Integração com Frontend React

### **Exemplo Completo de Integração**

```typescript
// hooks/useAgent.ts
import { useState, useEffect } from 'react';

export const useAgent = () => {
  const [status, setStatus] = useState<AgentStatus | null>(null);
  const [loading, setLoading] = useState(false);

  const updateStatus = async () => {
    try {
      const currentStatus = await window.go.main.App.GetAgentStatus();
      setStatus(currentStatus);
    } catch (error) {
      console.error('Erro ao obter status:', error);
    }
  };

  const startWorkflow = async (workflowName: string, params: any) => {
    setLoading(true);
    try {
      await window.go.main.App.StartWorkflow({
        WorkflowName: workflowName,
        Parameters: params
      });
      await updateStatus();
    } catch (error) {
      console.error('Erro ao iniciar workflow:', error);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  const askAI = async (question: string) => {
    setLoading(true);
    try {
      const response = await window.go.main.App.AskAI(question);
      return response;
    } catch (error) {
      console.error('Erro no chat:', error);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    updateStatus();
    // Atualizar status a cada 5 segundos
    const interval = setInterval(updateStatus, 5000);
    return () => clearInterval(interval);
  }, []);

  return {
    status,
    loading,
    startWorkflow,
    askAI,
    updateStatus
  };
};
```

---

**Essa API fornece todas as funcionalidades necessárias para uma integração completa entre frontend e backend na aplicação BrxAgente! 🚀**