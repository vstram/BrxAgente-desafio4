import { useState, useEffect, useCallback, useRef } from 'react';
import {
  AgentStatus,
  AgentDashboardData,
  AgentStatusState,
  AgentStatusActions,
  WorkflowExecutionRequest,
  LogEntry,
  LogFilter,
  DashboardConfig,
  LogLevel,
  WorkflowStatus,
  AgentMetrics
} from '../types/agent';
import { GetAgentStatus, StartWorkflow, StopWorkflow, CancelWorkflow, ClearAgentLogs } from "../../wailsjs/go/main/App";

// Mock data para desenvolvimento - será substituído por chamadas Wails reais
const createMockDashboardData = (): AgentDashboardData => ({
  status: 'idle' as AgentStatus,
  recentLogs: [
    {
      id: '1',
      timestamp: new Date(Date.now() - 5000),
      level: 'info',
      message: 'Sistema iniciado com sucesso',
      source: 'system'
    },
    {
      id: '2', 
      timestamp: new Date(Date.now() - 3000),
      level: 'info',
      message: 'Aguardando comandos do usuário',
      source: 'agent'
    }
  ],
  metrics: {
    totalWorkflowsExecuted: 0,
    successfulWorkflows: 0,
    failedWorkflows: 0,
    averageExecutionTime: 0,
    totalProcessingTime: 0,
    collaboratorsProcessed: 0,
    reportsGenerated: 0,
    anomaliesDetected: 0,
    uptime: Date.now()
  },
  availableWorkflows: [
    'processar-vr',
    'validar-dados', 
    'gerar-relatorios',
    'detectar-anomalias'
  ],
  lastUpdated: new Date()
});

const defaultConfig: DashboardConfig = {
  refreshInterval: 2000, // 2 segundos
  maxLogs: 100,
  logLevels: ['info', 'warn', 'error'],
  showMetrics: true,
  compactMode: false,
  soundNotifications: false
};

export function useAgentStatus(): AgentStatusState & AgentStatusActions {
  const [state, setState] = useState<AgentStatusState>({
    dashboardData: null,
    isConnected: false,
    isLoading: true,
    error: null,
    config: defaultConfig
  });

  const refreshIntervalRef = useRef<number>();
  const mockWorkflowRef = useRef<number>();

  // Simular logs em tempo real para desenvolvimento
  const addMockLog = useCallback((level: LogLevel, message: string, source?: string) => {
    setState(prev => {
      if (!prev.dashboardData) return prev;
      
      const newLog: LogEntry = {
        id: Date.now().toString(),
        timestamp: new Date(),
        level,
        message,
        source: source || 'system'
      };

      const updatedLogs = [newLog, ...prev.dashboardData.recentLogs]
        .slice(0, prev.config.maxLogs);

      return {
        ...prev,
        dashboardData: {
          ...prev.dashboardData,
          recentLogs: updatedLogs,
          lastUpdated: new Date()
        }
      };
    });
  }, []);

  // Simular execução de workflow
  const simulateWorkflow = useCallback((workflowName: string) => {
    const steps = [
      { name: 'load-data', description: 'Carregando dados das planilhas' },
      { name: 'validate-data', description: 'Validando consistência dos dados' },
      { name: 'calculate-vr', description: 'Calculando valores de VR' },
      { name: 'generate-reports', description: 'Gerando relatórios' },
      { name: 'finalize', description: 'Finalizando processamento' }
    ];

    const workflow: WorkflowStatus = {
      id: Date.now().toString(),
      name: workflowName,
      description: `Execução do workflow ${workflowName}`,
      status: 'running',
      currentStepIndex: 0,
      totalSteps: steps.length,
      steps: steps.map(step => ({ 
        ...step, 
        status: 'pending' as const 
      })),
      startTime: new Date(),
      progress: 0
    };

    setState(prev => ({
      ...prev,
      dashboardData: prev.dashboardData ? {
        ...prev.dashboardData,
        status: 'running',
        currentWorkflow: workflow,
        lastUpdated: new Date()
      } : null
    }));

    addMockLog('info', `Iniciando workflow: ${workflowName}`, 'orchestrator');

    let currentStep = 0;
    const stepDuration = 3000; // 3 segundos por step

    const executeStep = () => {
      if (currentStep >= steps.length) {
        // Workflow completo
        setState(prev => {
          const updatedMetrics: AgentMetrics = prev.dashboardData ? {
            ...prev.dashboardData.metrics,
            totalWorkflowsExecuted: prev.dashboardData.metrics.totalWorkflowsExecuted + 1,
            successfulWorkflows: prev.dashboardData.metrics.successfulWorkflows + 1,
            collaboratorsProcessed: prev.dashboardData.metrics.collaboratorsProcessed + Math.floor(Math.random() * 500 + 50),
            reportsGenerated: prev.dashboardData.metrics.reportsGenerated + Math.floor(Math.random() * 3 + 1),
          } : {
            totalWorkflowsExecuted: 1,
            successfulWorkflows: 1,
            failedWorkflows: 0,
            averageExecutionTime: stepDuration * steps.length,
            totalProcessingTime: stepDuration * steps.length,
            collaboratorsProcessed: Math.floor(Math.random() * 500 + 50),
            reportsGenerated: Math.floor(Math.random() * 3 + 1),
            anomaliesDetected: Math.floor(Math.random() * 10),
            uptime: Date.now()
          };

          return {
            ...prev,
            dashboardData: prev.dashboardData ? {
              ...prev.dashboardData,
              status: 'idle',
              currentWorkflow: {
                ...workflow,
                status: 'idle',
                currentStepIndex: steps.length,
                progress: 100,
                endTime: new Date(),
                totalDuration: stepDuration * steps.length
              },
              metrics: updatedMetrics,
              lastUpdated: new Date()
            } : null
          };
        });

        addMockLog('info', `Workflow concluído com sucesso: ${workflowName}`, 'orchestrator');
        return;
      }

      // Atualizar step atual
      setState(prev => ({
        ...prev,
        dashboardData: prev.dashboardData ? {
          ...prev.dashboardData,
          currentWorkflow: prev.dashboardData.currentWorkflow ? {
            ...prev.dashboardData.currentWorkflow,
            currentStepIndex: currentStep,
            progress: Math.round((currentStep / steps.length) * 100),
            steps: prev.dashboardData.currentWorkflow.steps.map((step, index) => ({
              ...step,
              status: index < currentStep ? 'completed' : 
                     index === currentStep ? 'running' : 'pending'
            }))
          } : undefined,
          lastUpdated: new Date()
        } : null
      }));

      addMockLog('info', `Executando: ${steps[currentStep].description}`, steps[currentStep].name);
      
      currentStep++;
      mockWorkflowRef.current = setTimeout(executeStep, stepDuration);
    };

    executeStep();
  }, [addMockLog]);

  // Carregar dados iniciais
  const loadDashboardData = useCallback(async () => {
    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }));
      
      // Chamar o backend para obter o status do agente
      const agentStatus = await GetAgentStatus();
      
      // Converter os dados do backend para o formato esperado pelo frontend
      const dashboardData: AgentDashboardData = {
        status: agentStatus.status as AgentStatus,
        recentLogs: (agentStatus.recentLogs || []).map(log => ({
          id: log.id,
          timestamp: new Date(log.timestamp),
          level: log.level as LogLevel,
          message: log.message,
          source: log.source
        })),
        metrics: {
          totalWorkflowsExecuted: agentStatus.metrics.totalWorkflowsExecuted,
          successfulWorkflows: agentStatus.metrics.successfulWorkflows,
          failedWorkflows: agentStatus.metrics.totalWorkflowsExecuted - agentStatus.metrics.successfulWorkflows,
          averageExecutionTime: 0, // Será calculado quando necessário
          totalProcessingTime: 0, // Será calculado quando necessário
          collaboratorsProcessed: agentStatus.metrics.collaboratorsProcessed,
          reportsGenerated: agentStatus.metrics.reportsGenerated,
          anomaliesDetected: agentStatus.metrics.anomaliesDetected,
          uptime: agentStatus.metrics.uptime
        },
        currentWorkflow: agentStatus.currentWorkflow ? {
          id: agentStatus.currentWorkflow.id,
          name: agentStatus.currentWorkflow.name,
          description: agentStatus.currentWorkflow.name,
          status: agentStatus.currentWorkflow.status as AgentStatus,
          currentStepIndex: agentStatus.currentWorkflow.steps.findIndex(step => step.status === 'running'),
          totalSteps: agentStatus.currentWorkflow.steps.length,
          steps: agentStatus.currentWorkflow.steps.map(step => ({
            name: step.id,
            description: step.name,
            status: step.status === 'error' ? 'failed' : step.status as 'pending' | 'running' | 'completed' | 'failed',
            startTime: step.startTime ? new Date(step.startTime) : undefined,
            endTime: step.endTime ? new Date(step.endTime) : undefined,
            duration: step.duration,
            error: step.errorMsg
          })),
          startTime: new Date(agentStatus.currentWorkflow.startTime),
          endTime: agentStatus.currentWorkflow.endTime ? new Date(agentStatus.currentWorkflow.endTime) : undefined,
          totalDuration: agentStatus.currentWorkflow.endTime ? 
            new Date(agentStatus.currentWorkflow.endTime).getTime() - new Date(agentStatus.currentWorkflow.startTime).getTime() : undefined,
          progress: agentStatus.currentWorkflow.progress
        } : undefined,
        availableWorkflows: agentStatus.availableWorkflows,
        lastUpdated: agentStatus.lastUpdated
      };
      
      setState(prev => ({
        ...prev,
        dashboardData,
        isConnected: true,
        isLoading: false,
        error: null
      }));

    } catch (error) {
      console.error('Error loading dashboard data:', error);
      setState(prev => ({
        ...prev,
        isConnected: false,
        isLoading: false,
        error: error instanceof Error ? error.message : 'Erro desconhecido'
      }));
    }
  }, []);

  // Inicializar dados e configurar refresh automático
  useEffect(() => {
    loadDashboardData();

    if (state.config.refreshInterval > 0) {
      refreshIntervalRef.current = setInterval(() => {
        if (state.isConnected && !state.isLoading) {
          loadDashboardData();
        }
      }, state.config.refreshInterval);
    }

    return () => {
      if (refreshIntervalRef.current) {
        clearInterval(refreshIntervalRef.current);
      }
      if (mockWorkflowRef.current) {
        clearTimeout(mockWorkflowRef.current);
      }
    };
  }, [state.config.refreshInterval, state.isConnected, state.isLoading, loadDashboardData]);

  // Ações
  const startWorkflow = useCallback(async (request: WorkflowExecutionRequest) => {
    try {
      // Preparar parâmetros para o backend
      const workflowRequest = {
        workflowName: request.workflowName,
        parameters: request.parameters || {}
      };
      
      // Chamar o backend para iniciar o workflow
      await StartWorkflow(workflowRequest);
      
      // Atualizar dados após iniciar workflow
      await loadDashboardData();
      
    } catch (error) {
      setState(prev => ({
        ...prev,
        error: error instanceof Error ? error.message : 'Erro ao iniciar workflow'
      }));
      console.error('Error starting workflow:', error);
    }
  }, [loadDashboardData]);

  const stopWorkflow = useCallback(async (workflowId?: string) => {
    try {
      // Chamar o backend para parar o workflow
      await StopWorkflow();
      
      // Atualizar dados após parar workflow
      await loadDashboardData();

    } catch (error) {
      setState(prev => ({
        ...prev,
        error: error instanceof Error ? error.message : 'Erro ao parar workflow'
      }));
      console.error('Error stopping workflow:', error);
    }
  }, [loadDashboardData]);

  const cancelWorkflow = useCallback(async (workflowId?: string) => {
    try {
      // Chamar o backend para cancelar o workflow
      await CancelWorkflow();
      
      // Atualizar dados após cancelar workflow
      await loadDashboardData();

    } catch (error) {
      setState(prev => ({
        ...prev,
        error: error instanceof Error ? error.message : 'Erro ao cancelar workflow'
      }));
      console.error('Error cancelling workflow:', error);
    }
  }, [loadDashboardData]);

  const refreshData = useCallback(async () => {
    await loadDashboardData();
  }, [loadDashboardData]);

  const clearLogs = useCallback(async () => {
    try {
      // Chamar o backend para limpar os logs
      await ClearAgentLogs();
      
      // Atualizar dados após limpar logs
      await loadDashboardData();
      
    } catch (error) {
      setState(prev => ({
        ...prev,
        error: error instanceof Error ? error.message : 'Erro ao limpar logs'
      }));
      console.error('Error clearing logs:', error);
    }
  }, [loadDashboardData]);

  const filterLogs = useCallback((filter: LogFilter): LogEntry[] => {
    if (!state.dashboardData) return [];
    
    return state.dashboardData.recentLogs.filter(log => {
      if (filter.levels && !filter.levels.includes(log.level)) return false;
      if (filter.source && log.source !== filter.source) return false;
      if (filter.dateFrom && log.timestamp < filter.dateFrom) return false;
      if (filter.dateTo && log.timestamp > filter.dateTo) return false;
      if (filter.searchText && !log.message.toLowerCase().includes(filter.searchText.toLowerCase())) return false;
      return true;
    });
  }, [state.dashboardData]);

  const updateConfig = useCallback((newConfig: Partial<DashboardConfig>) => {
    setState(prev => ({
      ...prev,
      config: { ...prev.config, ...newConfig }
    }));
  }, []);

  const exportLogs = useCallback((format: 'json' | 'csv' | 'txt') => {
    if (!state.dashboardData) return;

    const logs = state.dashboardData.recentLogs;
    let content: string;
    let mimeType: string;
    let extension: string;

    switch (format) {
      case 'json':
        content = JSON.stringify(logs, null, 2);
        mimeType = 'application/json';
        extension = 'json';
        break;
      case 'csv':
        const header = 'Timestamp,Level,Source,Message\n';
        const rows = logs.map(log => 
          `"${log.timestamp.toISOString()}","${log.level}","${log.source || ''}","${log.message.replace(/"/g, '""')}"`
        ).join('\n');
        content = header + rows;
        mimeType = 'text/csv';
        extension = 'csv';
        break;
      case 'txt':
        content = logs.map(log => 
          `[${log.timestamp.toISOString()}] [${log.level.toUpperCase()}] ${log.source ? `[${log.source}] ` : ''}${log.message}`
        ).join('\n');
        mimeType = 'text/plain';
        extension = 'txt';
        break;
    }

    const blob = new Blob([content], { type: mimeType });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `agent-logs-${new Date().toISOString().split('T')[0]}.${extension}`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }, [state.dashboardData]);

  return {
    ...state,
    startWorkflow,
    stopWorkflow,
    cancelWorkflow,
    refreshData,
    clearLogs,
    filterLogs,
    updateConfig,
    exportLogs
  };
}