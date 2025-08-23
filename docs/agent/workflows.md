# ⚙️ Workflows - Agente de IA BrxAgente

Este documento detalha todos os workflows disponíveis na aplicação, como configurá-los e personalizá-los para suas necessidades específicas.

## 📑 Índice

1. [Visão Geral dos Workflows](#-visão-geral-dos-workflows)
2. [Workflows Predefinidos](#-workflows-predefinidos)
3. [Configuração e Parâmetros](#-configuração-e-parâmetros)
4. [Workflows Personalizados](#-workflows-personalizados)
5. [Monitoramento e Controle](#-monitoramento-e-controle)
6. [Melhores Práticas](#-melhores-práticas)

## 🔄 Visão Geral dos Workflows

### **O que são Workflows?**

Workflows são sequências automatizadas de operações que o Agente de IA executa para processar dados de VR. Cada workflow tem:

- 📋 **Etapas Definidas**: Sequência de operações específicas
- ⚙️ **Parâmetros Configuráveis**: Opções para customizar o comportamento  
- 📊 **Monitoramento em Tempo Real**: Acompanhamento visual do progresso
- 🔄 **Controles de Execução**: Pausar, continuar, cancelar
- 📈 **Métricas de Performance**: Tempo, sucesso, erros

### **Como Funciona na Interface:**

```
┌─────────────────────────────────────────────────────┐
│ 🔄 Seletor de Workflow                             │
│ ┌─────────────────────────────────────────────────┐ │
│ │ Workflow: [Processamento Completo ▼]           │ │
│ └─────────────────────────────────────────────────┘ │
│                                                     │
│ ⚙️ Parâmetros:                                     │
│ • Diretório: /path/to/planilhas  [📁 Selecionar]   │
│ • Sindicatos: [✓] SINDPD [✓] SINDAC [✓] SINDMET   │
│ • Notificações: [✓] Email [✓] Dashboard           │
│                                                     │
│ 📊 Preview das Etapas:                             │
│ 1. ✅ Validação do diretório                       │
│ 2. ⏳ Leitura das planilhas                        │
│ 3. ⏳ Consolidação dos dados                       │
│ 4. ⏳ Aplicação das regras                         │
│ 5. ⏳ Geração de relatórios                        │
│                                                     │
│ [▶️ Iniciar] [⏸️ Pausar] [❌ Cancelar] [🔄 Reset]   │
└─────────────────────────────────────────────────────┘
```

## 📋 Workflows Predefinidos

### **1. Processamento Completo de VR**
**Nome:** `analise-vr-mensal`
**Descrição:** Executa todo o fluxo de processamento do VR mensal

#### **Etapas:**
1. **📂 Validação do Diretório** (15-30s)
   - Verifica se o diretório existe
   - Confirma presença de arquivos .xlsx
   - Valida permissões de leitura

2. **📊 Leitura das Planilhas** (30s-2min)
   - Lê colaboradores.xlsx
   - Lê afastamentos.xlsx
   - Lê feriados.xlsx
   - Valida formato e integridade

3. **🔄 Consolidação dos Dados** (1-3min)
   - Combina dados das planilhas
   - Aplica relacionamentos (matrícula)
   - Cria estrutura consolidada

4. **⚙️ Aplicação das Regras de Negócio** (2-5min)
   - Calcula dias úteis por sindicato
   - Aplica regras de elegibilidade
   - Processa afastamentos e férias
   - Calcula valores de VR

5. **📋 Geração da Planilha Final** (30s-1min)
   - Cria planilha resultado
   - Formata conforme padrão
   - Gera relatórios auxiliares

6. **🤖 Configuração do Contexto do Chat** (10-15s)
   - Carrega dados no chat inteligente
   - Prepara para consultas

#### **Parâmetros Configuráveis:**
```typescript
{
  directory: string,           // Diretório das planilhas
  sindicatos?: string[],       // Filtrar sindicatos específicos
  includeInactive?: boolean,   // Incluir colaboradores inativos
  sendNotifications?: boolean, // Enviar notificações
  generateBackup?: boolean,    // Criar backup automático
  strictValidation?: boolean   // Validação rigorosa
}
```

#### **Outputs:**
- 📊 Planilha final de VR
- 📋 Relatório de exceções
- 📈 Métricas de processamento
- 🗂️ Logs detalhados

---

### **2. Validação de Planilhas**
**Nome:** `validacao-planilhas`
**Descrição:** Valida formato e integridade das planilhas sem processamento

#### **Etapas:**
1. **🔍 Verificação de Arquivos** (10s)
   - Lista arquivos .xlsx encontrados
   - Verifica tamanhos e datas

2. **📊 Validação de Formato** (30s-1min)
   - Confirma estrutura das planilhas
   - Valida colunas obrigatórias
   - Verifica tipos de dados

3. **✅ Verificação de Consistência** (1-2min)
   - Valida relacionamentos entre planilhas
   - Identifica duplicatas
   - Verifica integridade referencial

#### **Relatório de Validação:**
```
✅ VALIDAÇÃO CONCLUÍDA

📂 Arquivos Encontrados:
• colaboradores.xlsx - 2.847 registros
• afastamentos.xlsx - 423 registros  
• feriados.xlsx - 18 registros

📊 Estrutura:
✅ Todas as colunas obrigatórias presentes
✅ Tipos de dados corretos
✅ Formato de datas válido

🔗 Consistência:
✅ Matrículas únicas em colaboradores
⚠️  3 matrículas em afastamentos não encontradas
✅ Datas de feriados válidas

💡 Recomendações:
• Verificar matrículas: MAT001, MAT045, MAT892
• Arquivo pronto para processamento
```

---

### **3. Detecção de Anomalias**
**Nome:** `deteccao-anomalias`
**Descrição:** Identifica padrões anômalos sem executar cálculos completos

#### **Etapas:**
1. **📊 Análise de Padrões** (1-2min)
   - Analisa distribuições estatísticas
   - Identifica outliers por sindicato
   - Compara com histórico (se disponível)

2. **🔍 Detecção de Outliers** (1-3min)
   - Aplica algoritmos de detecção
   - Calcula scores de anomalia
   - Classifica por severidade

3. **⚠️ Geração de Alertas** (30s)
   - Cria relatório estruturado
   - Define níveis de prioridade
   - Sugere ações corretivas

#### **Configurações de Detecção:**
```typescript
{
  sensitivity: 'low' | 'medium' | 'high',
  thresholdMultiplier: number,     // 1.5x, 2.0x, 2.5x
  includeHistorical: boolean,      // Comparar com histórico
  anomalyTypes: string[]           // ['value', 'pattern', 'temporal']
}
```

---

### **4. Geração de Relatórios**
**Nome:** `geracao-relatorios`  
**Descrição:** Gera relatórios baseados em dados já processados

#### **Etapas:**
1. **📊 Coleta de Métricas** (30s)
   - Agrega dados processados
   - Calcula estatísticas descritivas
   - Gera breakdowns por dimensão

2. **📈 Processamento de Dados** (1min)
   - Aplica transformações
   - Gera visualizações
   - Calcula tendências

3. **📋 Geração do Relatório** (1-2min)
   - Cria documento estruturado
   - Insere gráficos e tabelas
   - Formata para apresentação

#### **Tipos de Relatório:**
- 📊 **Executivo**: Resumo alto nível
- 📈 **Analítico**: Detalhes e breakdowns
- 🔍 **Anomalias**: Foco em exceções
- 📅 **Temporal**: Evolução por período
- 👥 **Por Sindicato**: Análise segmentada

## ⚙️ Configuração e Parâmetros

### **Parâmetros Globais**

Aplicam-se a todos os workflows:

```typescript
interface GlobalWorkflowParams {
  // Diretórios
  inputDirectory: string;
  outputDirectory?: string;
  backupDirectory?: string;
  
  // Processamento
  parallelProcessing?: boolean;
  maxWorkers?: number;
  timeout?: number; // segundos
  
  // Validação
  strictMode?: boolean;
  skipValidation?: boolean;
  continueOnError?: boolean;
  
  // Notificações
  notifications?: {
    email?: string[];
    slack?: string;
    desktop?: boolean;
  };
  
  // Cache
  useCache?: boolean;
  cacheInvalidation?: boolean;
}
```

### **Parâmetros Específicos por Workflow**

#### **Processamento Completo:**
```typescript
interface ProcessamentoCompletoParams extends GlobalWorkflowParams {
  sindicatos?: string[];          // Filtrar sindicatos
  includeInactive?: boolean;      // Incluir inativos
  dateRange?: {                   // Período de processamento
    start: Date;
    end: Date;
  };
  calculationRules?: {            // Regras de cálculo customizadas
    vrBase?: number;
    proportionalThreshold?: number;
  };
  outputFormats?: string[];       // ['xlsx', 'csv', 'json']
}
```

#### **Detecção de Anomalias:**
```typescript
interface DeteccaoAnomaliasParams extends GlobalWorkflowParams {
  sensitivity: 'low' | 'medium' | 'high';
  thresholdMultiplier: number;    // 1.5, 2.0, 2.5
  anomalyTypes: AnomalyType[];    // Tipos a detectar
  historicalComparison?: {        // Comparação histórica
    enabled: boolean;
    months: number;
  };
  actionThresholds?: {            // Limites para ação
    critical: number;
    warning: number;
    info: number;
  };
}
```

## 🔧 Workflows Personalizados

### **Criando Workflow Personalizado**

Você pode criar workflows customizados combinando etapas existentes:

```yaml
# config/workflows/custom-sindpd.yaml
name: "processamento-sindpd"
description: "Processamento específico para SINDPD com validações extras"
version: "1.0"

parameters:
  sindicato: "SINDPD"
  validacao_rigorosa: true
  backup_obrigatorio: true

steps:
  - name: "validar_diretorio"
    type: "validation"
    params:
      strict_mode: true
      
  - name: "ler_planilhas"
    type: "data_input"
    params:
      filter_sindicato: "${parameters.sindicato}"
      validate_each_row: true
      
  - name: "validacao_extra_sindpd"
    type: "custom_validation"
    params:
      rules: ["sindpd_specific", "double_check_values"]
      
  - name: "calcular_vr"
    type: "calculation"
    params:
      vr_base: 460.00
      regras_especiais: "sindpd_2025"
      
  - name: "gerar_relatorio"
    type: "report_generation"
    params:
      template: "sindpd_template.xlsx"
      include_breakdown: true
      
  - name: "backup_final"
    type: "backup"
    params:
      mandatory: true
      retention: "12_months"

conditions:
  - if: "errors.critical > 0"
    action: "pause_for_review"
  - if: "anomalies.count > 5"
    action: "send_alert"

notifications:
  on_completion: ["sindpd-manager@empresa.com"]
  on_error: ["it-support@empresa.com"]
  on_anomaly: ["rh-team@empresa.com"]
```

### **Usando Workflow Personalizado na Interface:**

```typescript
// Frontend React - Carregar workflow personalizado
const customWorkflow = await loadCustomWorkflow('custom-sindpd.yaml');

const startCustomWorkflow = async () => {
  await window.go.main.App.StartWorkflow({
    WorkflowName: 'processamento-sindpd',
    Parameters: {
      directory: selectedDirectory,
      sindicato: 'SINDPD',
      validacao_rigorosa: true
    }
  });
};
```

## 📊 Monitoramento e Controle

### **Interface de Monitoramento**

Durante a execução, a interface mostra:

```
┌─────────────────────────────────────────────────────┐
│ 🔄 Workflow: Processamento Completo                │
│ Status: ▶️  Executando • Etapa 3 de 6             │
├─────────────────────────────────────────────────────┤
│ 📊 Progresso Geral: ██████████████░░░░░ 75%        │
│                                                     │
│ 📋 Etapas:                                         │
│ 1. ✅ Validação do diretório          (15s)       │
│ 2. ✅ Leitura das planilhas           (1m 30s)    │
│ 3. 🔄 Consolidação dos dados          (2m 15s)    │  
│ 4. ⏳ Aplicação das regras            (est. 3m)   │
│ 5. ⏳ Geração da planilha             (est. 1m)   │
│ 6. ⏳ Configuração do contexto        (est. 15s)  │
│                                                     │
│ 📈 Métricas Atual:                                 │
│ • Colaboradores processados: 2.134 / 2.847        │
│ • Throughput: 18.5 colaboradores/min              │
│ • Anomalias detectadas: 3                          │
│ • Tempo estimado restante: 4m 32s                  │
│                                                     │
│ [⏸️ Pausar] [❌ Cancelar] [📊 Ver Logs]            │
└─────────────────────────────────────────────────────┘
```

### **Controles Disponíveis**

#### **Durante Execução:**
- ⏸️ **Pausar**: Para workflow atual (estado salvo)
- ▶️ **Continuar**: Retoma execução do ponto pausado
- ❌ **Cancelar**: Para workflow e limpa estado
- 📊 **Ver Logs**: Abre janela com logs detalhados
- ⏭️ **Pular Etapa**: Avança para próxima etapa (se seguro)

#### **Configurações de Execução:**
- 🔄 **Auto-restart**: Reinicia automaticamente em caso de erro
- ⏰ **Timeout**: Define tempo máximo por etapa
- 🚨 **Stop on Error**: Para execução no primeiro erro crítico
- 📧 **Notify Progress**: Envia atualizações periódicas

### **Sistema de Logs**

```
┌─────────────────────────────────────────────────────┐
│ 📋 Logs do Workflow - Processamento Completo       │
├─────────────────────────────────────────────────────┤
│ [INFO ] 14:32:15 - Workflow iniciado               │
│ [INFO ] 14:32:16 - Diretório validado: 3 planilhas │
│ [INFO ] 14:32:45 - Leitura concluída: 2847 linhas  │
│ [WARN ] 14:34:12 - Anomalia detectada: MAT001234   │
│ [INFO ] 14:34:15 - Consolidação: 75% concluído     │
│ [ERROR] 14:34:22 - Erro na matrícula: MAT005678    │
│ [INFO ] 14:34:25 - Erro corrigido automaticamente  │
│                                                     │
│ Filtros: [✓] INFO [✓] WARN [✓] ERROR              │
│ [📥 Exportar] [🗑️ Limpar] [🔍 Buscar]              │
└─────────────────────────────────────────────────────┘
```

## 🎯 Melhores Práticas

### **Escolha do Workflow Apropriado**

| Cenário | Workflow Recomendado | Motivo |
|---------|---------------------|---------|
| **Processamento mensal regular** | Processamento Completo | Fluxo completo A-Z |
| **Verificar qualidade dos dados** | Validação de Planilhas | Rápido, sem processamento |
| **Investigar problemas** | Detecção de Anomalias | Foco na identificação |
| **Gerar apresentações** | Geração de Relatórios | Outputs profissionais |
| **Processo customizado** | Workflow Personalizado | Necessidades específicas |

### **Configurações por Ambiente**

#### **Desenvolvimento:**
```yaml
environment: "development"
parameters:
  strict_mode: false
  timeout: 300
  notifications: false
  backup: false
logging:
  level: "debug"
  console: true
```

#### **Produção:**
```yaml
environment: "production"  
parameters:
  strict_mode: true
  timeout: 1800
  notifications: true
  backup: true
  auto_restart: true
logging:
  level: "info"
  file: true
  retention: "30d"
```

### **Otimização de Performance**

#### **Para Grandes Volumes:**
- 🔧 Use `parallelProcessing: true`
- 💾 Ative cache (`useCache: true`)
- 📊 Configure `maxWorkers` adequadamente
- ⚡ Desabilite validações não críticas

#### **Para Máxima Precisão:**
- ✅ Ative `strictMode: true`
- 🔍 Use `sensitivity: 'high'` para anomalias
- 📝 Configure logs detalhados
- 🚨 Ative `stopOnError: true`

### **Monitoramento Proativo**

```typescript
// Hook para monitoramento contínuo
const useWorkflowMonitoring = () => {
  const [metrics, setMetrics] = useState<WorkflowMetrics>();
  
  useEffect(() => {
    const interval = setInterval(async () => {
      const status = await window.go.main.App.GetAgentStatus();
      if (status.CurrentWorkflow) {
        setMetrics({
          progress: status.CurrentWorkflow.Progress,
          currentStep: status.CurrentWorkflow.Steps.find(s => s.Status === 'running'),
          anomaliesDetected: getAnomaliesFromLogs(status.RecentLogs),
          performance: calculatePerformance(status.Metrics)
        });
      }
    }, 2000); // Atualiza a cada 2 segundos
    
    return () => clearInterval(interval);
  }, []);
  
  return metrics;
};
```

### **Tratamento de Erros**

```typescript
// Padrão para recuperação de erros
const handleWorkflowError = async (error: string) => {
  // 1. Log do erro
  console.error('Workflow error:', error);
  
  // 2. Avaliar se pode recuperar
  const canRecover = await assessRecoveryOptions(error);
  
  if (canRecover.autoFix) {
    // 3a. Tentar correção automática
    await applyAutoFix(canRecover.fixType);
    await window.go.main.App.StartWorkflow(lastWorkflowParams);
  } else {
    // 3b. Pausar para intervenção manual
    await window.go.main.App.StopWorkflow();
    showManualInterventionDialog(error, canRecover.suggestions);
  }
};
```

---

**Os workflows são o coração da automação do BrxAgente. Use-os para criar processos eficientes e confiáveis! 🚀**