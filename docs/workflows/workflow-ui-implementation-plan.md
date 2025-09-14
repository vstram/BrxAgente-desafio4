# 🔧 Plano de Implementação - Workflow UI

## 📋 Visão Geral

Este documento apresenta o plano detalhado para implementar a funcionalidade de **Orquestrador de Workflows** na interface do usuário, baseado no sistema existente de Chat popup e seguindo as especificações do user guide.

## 🎯 Objetivo

Criar um sistema de interface para workflows que permita:
- ✅ **Seleção e execução de workflows**
- ✅ **Monitoramento em tempo real do progresso**
- ✅ **Controle de execução** (start, stop, cancel)
- ✅ **Visualização detalhada dos steps**
- ✅ **Dashboard com métricas e logs**

## 🔍 Análise da Implementação Atual do Chat

### **Estrutura Base Identificada:**
```typescript
// Chat.tsx: Popup com overlay
- useState para controle do popup (isOpen)
- Overlay com click-outside para fechar
- Header com título e botões de ação
- Área de conteúdo principal
- Botão flutuante para abrir

// App.tsx: Integração no footer
- Botão de configurações no footer
- Component Chat renderizado no final
- Modal de configuração existente
```

### **Padrões a Reutilizar:**
- ✅ **Overlay + Popup Pattern**: `chat-overlay` + `chat-window`
- ✅ **Header com Actions**: Título + botões de controle
- ✅ **Estado de Loading**: Indicadores visuais + disabled states
- ✅ **Botão Flutuante**: Posicionamento fixo no footer
- ✅ **CSS Classes**: Padrão `.btn`, `.icon`, `.overlay`, `.window`

## 📐 Arquitetura da Solução

### **1. Estrutura de Arquivos**
```
frontend/src/
├── components/
│   ├── WorkflowUI.tsx           # Componente principal (similar ao Chat.tsx)
│   ├── WorkflowUI.css          # Estilos específicos
│   ├── WorkflowSelector.tsx    # Seleção de workflows
│   ├── WorkflowControls.tsx    # Controles (start/stop/cancel)
│   ├── WorkflowProgress.tsx    # ✅ JÁ EXISTE - melhorar integração
│   ├── WorkflowStepsList.tsx   # Lista detalhada de steps
│   └── WorkflowMetrics.tsx     # Métricas e estatísticas
├── hooks/
│   └── useWorkflowManager.ts   # Hook para gerenciar workflows
├── types/
│   └── workflow.ts             # ✅ JÁ EXISTE - expandir se necessário
└── assets/icons/
    ├── workflow.svg            # Ícone do botão principal
    ├── play-workflow.svg       # Executar workflow
    ├── pause-workflow.svg      # Pausar workflow
    └── stop-workflow.svg       # Parar workflow
```

### **2. Componente Principal - WorkflowUI.tsx**

```typescript
interface WorkflowUIProps {}

function WorkflowUI() {
  // Estados similares ao Chat
  const [isOpen, setIsOpen] = useState(false);
  const [selectedWorkflow, setSelectedWorkflow] = useState<string>('');
  const [currentExecution, setCurrentExecution] = useState<WorkflowExecution | null>(null);
  const [availableWorkflows, setAvailableWorkflows] = useState<string[]>([]);
  const [isExecuting, setIsExecuting] = useState(false);

  // Hook personalizado para gerenciar workflows
  const {
    executeWorkflow,
    cancelWorkflow,
    stopWorkflow,
    getWorkflowStatus,
    subscribeToUpdates
  } = useWorkflowManager();

  return (
    <>
      {/* Popup Principal */}
      {isOpen && (
        <div className="workflow-overlay" onClick={() => setIsOpen(false)}>
          <div className="workflow-window" onClick={(e) => e.stopPropagation()}>

            {/* Header */}
            <div className="workflow-header">
              <h3>🔧 Orquestrador de Workflows</h3>
              <div className="workflow-header-actions">
                <button className="workflow-refresh-btn" onClick={refreshWorkflows}>
                  <img src={RefreshIcon} alt="Atualizar" className="icon" />
                </button>
                <button className="workflow-close-btn" onClick={() => setIsOpen(false)}>
                  <img src={XIcon} alt="Fechar" className="icon" />
                </button>
              </div>
            </div>

            {/* Área Principal */}
            <div className="workflow-content">

              {/* Seleção de Workflow */}
              <WorkflowSelector
                workflows={availableWorkflows}
                selected={selectedWorkflow}
                onSelect={setSelectedWorkflow}
                disabled={isExecuting}
              />

              {/* Controles de Execução */}
              <WorkflowControls
                selectedWorkflow={selectedWorkflow}
                isExecuting={isExecuting}
                onExecute={handleExecute}
                onCancel={handleCancel}
                onStop={handleStop}
              />

              {/* Progresso em Tempo Real */}
              {currentExecution && (
                <WorkflowProgress
                  workflow={currentExecution}
                  onCancel={handleCancel}
                  showDetails={true}
                />
              )}

              {/* Lista de Steps Detalhada */}
              {currentExecution && (
                <WorkflowStepsList
                  steps={currentExecution.steps}
                  currentStep={currentExecution.currentStepIndex}
                />
              )}

              {/* Métricas e Logs */}
              <WorkflowMetrics
                execution={currentExecution}
                showLogs={true}
              />

            </div>
          </div>
        </div>
      )}

      {/* Botão Flutuante */}
      <button className="workflow-toggle-btn" onClick={() => setIsOpen(!isOpen)}>
        <img src={WorkflowIcon} alt="Workflows" className="icon" />
      </button>
    </>
  );
}
```

### **3. Hook de Gerenciamento - useWorkflowManager.ts**

```typescript
interface UseWorkflowManagerReturn {
  // Estados
  availableWorkflows: string[];
  currentExecution: WorkflowExecution | null;
  isLoading: boolean;
  error: string | null;

  // Ações
  executeWorkflow: (name: string, params?: any) => Promise<string>; // returns executionId
  cancelWorkflow: (executionId?: string) => Promise<void>;
  stopWorkflow: (executionId?: string) => Promise<void>;
  getWorkflowStatus: (executionId: string) => Promise<WorkflowExecution>;
  refreshWorkflows: () => Promise<void>;

  // Subscription para updates em tempo real
  subscribeToUpdates: (executionId: string) => () => void; // returns unsubscribe
}

export function useWorkflowManager(): UseWorkflowManagerReturn {
  // Implementação usando as funções do Wails
  // Similar ao padrão usado no Chat.tsx

  const executeWorkflow = async (name: string, params = {}) => {
    // Call: ExecuteWorkflowAsync(name, params)
    // Returns execution ID
  };

  const getWorkflowStatus = async (executionId: string) => {
    // Call: GetWorkflowStatus(executionId)
    // Returns WorkflowExecution
  };

  // Polling para updates em tempo real
  useEffect(() => {
    if (currentExecution && currentExecution.status === 'running') {
      const interval = setInterval(async () => {
        const status = await getWorkflowStatus(currentExecution.id);
        setCurrentExecution(status);
      }, 1000); // Update a cada 1 segundo

      return () => clearInterval(interval);
    }
  }, [currentExecution]);
}
```

## 🔌 Integração com Backend (Go)

### **Endpoints Necessários no app.go:**

```go
// Novos métodos para o contexto Wails
func (a *App) GetAvailableWorkflows() []string {
    if a.orchestrator == nil {
        return []string{}
    }
    return a.orchestrator.ListWorkflows()
}

func (a *App) ExecuteWorkflowAsync(name string, params map[string]interface{}) (string, error) {
    if a.orchestrator == nil {
        return "", fmt.Errorf("orchestrator not available")
    }
    return a.orchestrator.ExecuteWorkflowAsync(name, params)
}

func (a *App) GetWorkflowExecution(executionId string) (*WorkflowExecution, error) {
    if a.orchestrator == nil {
        return nil, fmt.Errorf("orchestrator not available")
    }
    return a.orchestrator.GetExecution(executionId)
}

func (a *App) CancelWorkflowExecution(executionId string) error {
    if a.orchestrator == nil {
        return fmt.Errorf("orchestrator not available")
    }
    return a.orchestrator.CancelExecution(executionId)
}

func (a *App) GetWorkflowStats() map[string]interface{} {
    if a.orchestrator == nil {
        return map[string]interface{}{}
    }
    return a.orchestrator.GetStats()
}
```

### **Modificação na Inicialização:**

```go
// No NewApp() ou similar, garantir que o orchestrator seja inicializado
func NewApp() *App {
    // ... código existente ...

    // Inicializar orchestrator
    orchestrator := workflows.NewOrchestrator(logger, workflows.OrchestratorConfig{
        MaxConcurrentWorkflows: 5,
        DefaultTimeout:         30 * time.Minute,
        EnableRollback:         true,
        DetailedLogging:        true,
    })

    // Registrar workflows padrão
    orchestrator.RegisterWorkflow(workflows.NewVRWorkflow())
    orchestrator.RegisterWorkflow(workflows.NewReportingWorkflow())
    orchestrator.RegisterWorkflow(workflows.NewSimpleValidationWorkflow())

    app.orchestrator = orchestrator

    return app
}
```

## 🎨 Interface do Usuário

### **1. Layout Principal**
```
┌─────────────────────────────────────────────────────────────┐
│  🔧 Orquestrador de Workflows                    [🔄] [✕]   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  📋 Seleção de Workflow                                     │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │ [📋 Processamento Completo de VR        ▼] [▶️ Executar] │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                             │
│  ⚡ Controles de Execução                                   │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │ Status: 🔄 Executando (Step 3/7)                        │ │
│  │ [⏸️ Pausar]  [⏹️ Parar]  [❌ Cancelar]                    │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                             │
│  📊 Progresso Detalhado                                     │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │ ████████████░░░░░░░░░░░ 60% (3/7 steps)                │ │
│  │                                                         │ │
│  │ ✅ 1. Validação do Diretório         (✅ 15s)          │ │
│  │ ✅ 2. Leitura das Planilhas          (✅ 45s)          │ │
│  │ 🔄 3. Auditoria de Dados             (⏳ 2m 30s)       │ │
│  │ ⏳ 4. Cálculos de VR                 (⏳ Aguardando)   │ │
│  │ ⏳ 5. Geração de Relatórios          (⏳ Aguardando)   │ │
│  │ ⏳ 6. Detecção de Anomalias          (⏳ Aguardando)   │ │
│  │ ⏳ 7. Notificações                   (⏳ Aguardando)   │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                             │
│  📈 Métricas                                                │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │ Tempo de Execução: 3m 45s    Colaboradores: 2.847     │ │
│  │ Steps Concluídos: 3/7         Taxa de Sucesso: 100%    │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### **2. Estados da Interface**

**Estado Idle (Aguardando):**
- ✅ Dropdown de workflows habilitado
- ✅ Botão "Executar" habilitado
- ❌ Controles de execução ocultos
- ❌ Progresso oculto

**Estado Executando:**
- ❌ Dropdown de workflows desabilitado
- ❌ Botão "Executar" desabilitado
- ✅ Controles de execução visíveis
- ✅ Progresso em tempo real
- ✅ Steps com status detalhado

**Estado Concluído/Erro:**
- ✅ Dropdown habilitado para nova execução
- ✅ Botão "Executar" habilitado
- ❌ Controles de execução ocultos
- ✅ Progresso final/erro exibido

## 📱 Integração no App Principal

### **Modificações no App.tsx:**
```typescript
// Adicionar import
import WorkflowUI from './components/WorkflowUI';

// No footer, adicionar botão (similar ao botão de configurações)
<footer className="app-footer">
    <button className="btn footer-btn" onClick={() => setIsConfigModalOpen(true)}>
        <img src={SettingsIcon} alt="Configurações" className="btn-icon" />
        Configurações
    </button>

    {/* NOVO: Botão de Workflows */}
    <button className="btn footer-btn workflow-footer-btn">
        <img src={WorkflowIcon} alt="Workflows" className="btn-icon" />
        Workflows
    </button>
</footer>

{/* No final do componente */}
<Chat />
<WorkflowUI />  {/* NOVO */}
```

## 🎯 Workflows Predefinidos

### **Workflows a Implementar (baseado no user guide):**

1. **📈 Processamento Completo**:
   - Steps: Validação → Leitura → Auditoria → Cálculos → Relatórios → Notificações
   - Parâmetros: Diretório das planilhas

2. **🔍 Apenas Auditoria**:
   - Steps: Validação → Leitura → Auditoria → Relatório de Anomalias
   - Parâmetros: Diretório + Sensibilidade

3. **⚙️ Cálculos Isolados**:
   - Steps: Validação → Leitura → Cálculos → Planilha Final
   - Parâmetros: Diretório

4. **📊 Geração de Relatórios**:
   - Steps: Leitura de dados → Análise → Geração → Exportação
   - Parâmetros: Tipo de relatório + Formato

## 🔧 Implementação por Fases

### **🚀 Fase 1 - MVP (1-2 semanas)**
- ✅ Botão flutuante no footer
- ✅ Popup básico com estrutura
- ✅ Seleção de workflow (dropdown)
- ✅ Execução básica (start/cancel)
- ✅ Progresso simples (barra de progresso)
- ✅ Integração básica com backend

### **📈 Fase 2 - Funcionalidades Avançadas (2-3 semanas)**
- ✅ Progresso detalhado com steps
- ✅ Controles avançados (pause/resume)
- ✅ Métricas em tempo real
- ✅ Logs de execução
- ✅ Preview de workflows
- ✅ Configuração de parâmetros

### **🎨 Fase 3 - Refinamentos (1 semana)**
- ✅ Melhorias visuais e animações
- ✅ Tratamento robusto de erros
- ✅ Testes automatizados
- ✅ Performance e otimizações
- ✅ Documentação final

## 📋 Checklist de Implementação

### **Backend (Go):**
- [ ] Adicionar métodos Wails para workflow API
- [ ] Garantir inicialização do orchestrator no app startup
- [ ] Registrar workflows predefinidos
- [ ] Implementar polling/websocket para updates em tempo real (opcional)
- [ ] Testes de integração

### **Frontend (React/TypeScript):**
- [ ] Criar componente WorkflowUI.tsx (estrutura base)
- [ ] Implementar hook useWorkflowManager.ts
- [ ] Criar componentes filhos (Selector, Controls, etc.)
- [ ] Adicionar ícones e assets necessários
- [ ] Integrar no App.tsx (botão + renderização)
- [ ] Estilos CSS (baseado no padrão do Chat)
- [ ] Testes unitários

### **Integração:**
- [ ] Testar execução end-to-end
- [ ] Validar cancelamento e controles
- [ ] Verificar updates em tempo real
- [ ] Testar com múltiplos workflows
- [ ] Validação de erros e edge cases

## 📚 Referências

- **Componente Chat Existente**: `/frontend/src/Chat.tsx`
- **Documentação User Guide**: `/docs/agent/user-guide.md` (seções 101-173)
- **Backend Orchestrator**: `/internal/workflows/orchestrator.go`
- **Testes de Workflows**: `/internal/workflows/*_test.go`
- **Componentes Existentes**: `/frontend/src/components/WorkflowProgress.tsx`

## 💡 Considerações Finais

Esta implementação seguirá o padrão já estabelecido pelo sistema de Chat, garantindo consistência visual e de experiência do usuário. O backend já está tecnicamente pronto e testado, necessitando apenas da exposição via API Wails e da interface frontend.

**A funcionalidade estará completamente integrada ao sistema existente, proporcionando aos usuários controle total sobre os workflows de processamento de VR/VA através de uma interface intuitiva e moderna.**