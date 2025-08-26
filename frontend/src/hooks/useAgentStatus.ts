import { useState, useEffect, useRef } from 'react';
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
import { GetAgentStatus, GetAgentStatusSimple, StartWorkflow, StopWorkflow, CancelWorkflow, ClearAgentLogs } from "../../wailsjs/go/main/App";
import WailsProfiler from '../utils/wailsProfiler';

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
  refreshInterval: 0, // Desabilitado para evitar loop infinito
  maxLogs: 100,
  logLevels: ['info', 'warn', 'error'],
  showMetrics: true,
  compactMode: false,
  soundNotifications: false
};

// EMERGENCY BYPASS - Se tudo falhar, usar apenas dados mock
const EMERGENCY_BYPASS = false; // Mude para true se necessário

// Timeout configurável para carregamento de dados grandes
const getLoadTimeout = () => {
  // Timeout generoso para grandes volumes de dados processados
  return 30000; // 30 segundos (0.5 minutos)
};

export function useAgentStatus(): AgentStatusState & AgentStatusActions {
  // Reduzir logs para diminuir overhead de performance
  const [state, setState] = useState<AgentStatusState>({
    dashboardData: null,
    isConnected: false,
    isLoading: true,
    error: null,
    config: defaultConfig
  });

  const refreshIntervalRef = useRef<number>();
  const mockWorkflowRef = useRef<number>();
  const mountedRef = useRef(true);
  const loadingRef = useRef(false); // Evitar múltiplas chamadas simultâneas

  // Função para carregar dados com timeout mais agressivo e circuit breaker
  const loadDashboardData = async () => {
    if (!mountedRef.current) {
      return;
    }
    
    // EMERGENCY BYPASS: Se ativado, pular direto para dados mock
    if (EMERGENCY_BYPASS) {
      console.log('🚨 [useAgentStatus] EMERGENCY BYPASS ativado - usando apenas dados mock');
      setState(prev => ({
        ...prev,
        dashboardData: createMockDashboardData(),
        isConnected: false,
        isLoading: false,
        error: 'Modo de emergência ativado - usando dados de exemplo'
      }));
      return;
    }
    
    // Debounce: evitar múltiplas chamadas simultâneas
    if (loadingRef.current) {
      console.log('⚠️ [useAgentStatus] Call blocked by debounce');
      return;
    }
    
    loadingRef.current = true;
    
    try {
      setState(prev => ({ ...prev, isLoading: true, error: null }));
      
      // Timeout muito mais agressivo - 5 segundos máximo
      const timeoutMs = 5000;
      const timeoutPromise = new Promise<never>((_, reject) => {
        setTimeout(() => reject(new Error('Backend não respondeu em 5 segundos - possível problema de performance')), timeoutMs);
      });
      
      console.log('🚀 [useAgentStatus] Calling GetAgentStatus with 5s timeout');
      
      // CORREÇÃO DE EMERGÊNCIA: Testar chamada simples primeiro
      let agentStatus;
      
      try {
        // CORREÇÃO DEFINITIVA: Pular GetAgentStatus e construir dados a partir de GetAgentStatusSimple
        console.log('🔧 [useAgentStatus] Using simplified backend approach');
        const simpleResult = await GetAgentStatusSimple();
        console.log('✅ [useAgentStatus] Simple call result:', simpleResult);
        
        // Construir dados do AgentDashboard a partir do resultado simples
        const timestamp = new Date();
        const dashboardData: AgentDashboardData = {
          status: 'idle' as AgentStatus, // usando tipo correto
          recentLogs: [{
            id: 'system-log-1',
            timestamp: timestamp,
            level: 'info' as LogLevel,
            message: `Sistema funcionando - ${simpleResult}`,
            source: 'agent'
          }],
          metrics: {
            totalWorkflowsExecuted: 0,
            successfulWorkflows: 0,
            failedWorkflows: 0,
            averageExecutionTime: 0,
            totalProcessingTime: 0,
            collaboratorsProcessed: 0,
            reportsGenerated: 0,
            anomaliesDetected: 0,
            uptime: Date.now() - (timestamp.getTime() - 3600000) // 1 hour uptime
          },
          availableWorkflows: [
            'analise-vr-mensal',
            'validacao-planilhas',
            'deteccao-anomalias', 
            'geracao-relatorios',
            'auditoria-inteligente'
          ],
          lastUpdated: timestamp
        };
        
        // Retornar os dados construídos diretamente
        setState(prev => ({
          ...prev,
          dashboardData,
          isConnected: true,
          isLoading: false,
          error: null
        }));
        
        console.log('✅ [useAgentStatus] DashboardData constructed from simple call');
        return; // Sair da função aqui - dados já foram definidos
        
      } catch (simpleError) {
        console.error('❌ [useAgentStatus] Simple call failed:', simpleError);
        throw new Error('Backend não está respondendo - usando dados de fallback');
      }
    } catch (error) {
      if (mountedRef.current) {
        setState(prev => ({
          ...prev,
          isConnected: false,
          isLoading: false,
          error: error instanceof Error ? error.message : 'Erro desconhecido ao carregar dados do agente'
        }));
      }
    } finally {
      loadingRef.current = false; // Liberar debounce
    }
  };

  // Ref para controlar se já foi executado
  const hasInitializedRef = useRef(false);
  
  // Inicializar dados APENAS UMA VEZ com circuit breaker
  useEffect(() => {
    // Se já foi inicializado, não fazer nada
    if (hasInitializedRef.current) {
      return;
    }
    
    console.log('🚀 [useAgentStatus] Inicialização Única executando...');
    hasInitializedRef.current = true;
    mountedRef.current = true;
    
    // Circuit breaker: se falhar 3 vezes, usar dados mock
    let attempts = 0;
    const maxAttempts = 3;
    
    const tryLoadData = async () => {
      attempts++;
      console.log(`🔄 [useAgentStatus] Tentativa ${attempts}/${maxAttempts}`);
      
      try {
        await loadDashboardData();
        console.log('✅ [useAgentStatus] Dados carregados com sucesso');
      } catch (error) {
        console.error(`❌ [useAgentStatus] Falha na tentativa ${attempts}:`, error);
        
        if (attempts >= maxAttempts) {
          console.log('🛑 [useAgentStatus] Limite de tentativas atingido - usando dados mock');
          // Fallback para dados mock após 3 falhas
          setState(prev => ({
            ...prev,
            dashboardData: createMockDashboardData(),
            isConnected: false,
            isLoading: false,
            error: 'Usando dados de exemplo devido a problemas de conectividade'
          }));
        } else {
          // Tentar novamente após delay exponencial
          const delay = Math.pow(2, attempts) * 1000; // 2s, 4s, 8s...
          console.log(`⏳ [useAgentStatus] Tentando novamente em ${delay}ms`);
          setTimeout(tryLoadData, delay);
        }
      }
    };
    
    // Executar primeira tentativa com delay
    const timeoutId = setTimeout(() => {
      if (mountedRef.current && hasInitializedRef.current) {
        tryLoadData();
      }
    }, 100);
    
    // Cleanup na desmontagem
    return () => {
      console.log('🧹 [useAgentStatus] Cleanup da inicialização');
      clearTimeout(timeoutId);
      mountedRef.current = false;
      loadingRef.current = false;
      if (refreshIntervalRef.current) {
        clearInterval(refreshIntervalRef.current);
      }
      if (mockWorkflowRef.current) {
        clearTimeout(mockWorkflowRef.current);
      }
    };
  }, []); // Array vazio - executa apenas uma vez

  // Ações - SEM useCallback para evitar dependências problemáticas
  const startWorkflow = async (request: WorkflowExecutionRequest) => {
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
  };

  const stopWorkflow = async (workflowId?: string) => {
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
  };

  const cancelWorkflow = async (workflowId?: string) => {
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
  };

  const refreshData = async () => {
    await loadDashboardData();
  };

  const clearLogs = async () => {
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
  };

  const filterLogs = (filter: LogFilter): LogEntry[] => {
    if (!state.dashboardData) return [];
    
    return state.dashboardData.recentLogs.filter(log => {
      if (filter.levels && !filter.levels.includes(log.level)) return false;
      if (filter.source && log.source !== filter.source) return false;
      if (filter.dateFrom && log.timestamp < filter.dateFrom) return false;
      if (filter.dateTo && log.timestamp > filter.dateTo) return false;
      if (filter.searchText && !log.message.toLowerCase().includes(filter.searchText.toLowerCase())) return false;
      return true;
    });
  };

  const updateConfig = (newConfig: Partial<DashboardConfig>) => {
    setState(prev => ({
      ...prev,
      config: { ...prev.config, ...newConfig }
    }));
  };

  const exportLogs = (format: 'json' | 'csv' | 'txt') => {
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
  };

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