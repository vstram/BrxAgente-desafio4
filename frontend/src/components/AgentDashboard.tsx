import React, { useState } from 'react';
import { AgentDashboardProps, AgentStatus } from '../types/agent';
import { useAgentStatus } from '../hooks/useAgentStatus';
import { usePerformanceMonitor } from '../hooks/usePerformanceMonitor';
import WorkflowProgress from './WorkflowProgress';
import AgentLogs from './AgentLogs';
import './AgentDashboard.css';

// Importando ícones existentes
import PlayIcon from '../assets/icons/play.svg';
import RefreshIcon from '../assets/icons/refresh.svg';
import CancelIcon from '../assets/icons/cancel.svg';
import SpinnerIcon from '../assets/icons/spinner.svg';
import InfoIcon from '../assets/icons/info.svg';

const AgentDashboard: React.FC<AgentDashboardProps> = ({
  className = '',
  compactMode = false,
  refreshInterval = 2000
}) => {
  // Monitor de performance
  const { metrics } = usePerformanceMonitor({
    enabled: true,
    componentName: 'AgentDashboard',
    logInterval: 3000
  });
  
  const {
    dashboardData,
    isConnected,
    isLoading,
    error,
    startWorkflow,
    stopWorkflow,
    cancelWorkflow,
    refreshData,
    clearLogs,
    filterLogs
  } = useAgentStatus();

  const [selectedWorkflow, setSelectedWorkflow] = useState<string>('');
  const [showLogs, setShowLogs] = useState(false);

  const getStatusColor = (status: AgentStatus): string => {
    switch (status) {
      case 'idle':
        return 'var(--success-color)';
      case 'running':
        return 'var(--primary-color)';
      case 'error':
        return 'var(--error-color)';
      case 'cancelled':
        return 'var(--warning-color)';
      default:
        return 'var(--text-secondary)';
    }
  };

  const getStatusIcon = (status: AgentStatus): string => {
    switch (status) {
      case 'idle':
        return '✅';
      case 'running':
        return '⏳';
      case 'error':
        return '❌';
      case 'cancelled':
        return '⚠️';
      default:
        return '❓';
    }
  };

  const getStatusText = (status: AgentStatus): string => {
    switch (status) {
      case 'idle':
        return 'Aguardando comandos';
      case 'running':
        return 'Executando workflow';
      case 'error':
        return 'Erro na execução';
      case 'cancelled':
        return 'Workflow cancelado';
      default:
        return 'Status desconhecido';
    }
  };

  const handleStartWorkflow = async () => {
    if (!selectedWorkflow) {
      alert('Por favor, selecione um workflow para executar');
      return;
    }
    
    try {
      let parameters = {};
      
      // Se for o workflow de análise VR, precisamos pedir o diretório
      if (selectedWorkflow === 'analise-vr-mensal') {
        const directory = prompt('Digite o caminho do diretório das planilhas:');
        if (!directory) {
          alert('Diretório é obrigatório para este workflow');
          return;
        }
        parameters = { directory };
      }
      
      await startWorkflow({ 
        workflowName: selectedWorkflow,
        parameters 
      });
      setSelectedWorkflow('');
    } catch (error) {
      console.error('Erro ao iniciar workflow:', error);
    }
  };

  const handleStopWorkflow = async () => {
    try {
      await stopWorkflow();
    } catch (error) {
      console.error('Erro ao parar workflow:', error);
    }
  };

  const handleCancelWorkflow = async () => {
    try {
      await cancelWorkflow();
    } catch (error) {
      console.error('Erro ao cancelar workflow:', error);
    }
  };

  const formatUptime = (ms: number): string => {
    const seconds = Math.floor(ms / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);
    
    if (days > 0) return `${days}d ${hours % 24}h`;
    if (hours > 0) return `${hours}h ${minutes % 60}m`;
    if (minutes > 0) return `${minutes}m ${seconds % 60}s`;
    return `${seconds}s`;
  };

  // Loading state
  if (isLoading) {
    return (
      <div className={`agent-dashboard loading ${className}`}>
        <div className="loading-content">
          <img src={SpinnerIcon} alt="Carregando" className="loading-spinner" />
          <span>Carregando status do agente...</span>
        </div>
      </div>
    );
  }

  // Error state
  if (error || !isConnected) {
    const isTimeoutError = error && error.includes('carregamento dos dados está demorando');
    
    return (
      <div className={`agent-dashboard error ${className}`}>
        <div className="error-content">
          <span className="error-icon">{isTimeoutError ? '⏳' : '⚠️'}</span>
          <h3>{isTimeoutError ? 'Carregamento Lento' : 'Erro de Conexão'}</h3>
          <p>{error || 'Não foi possível conectar ao agente'}</p>
          <div className="error-actions">
            <button className="btn primary-btn" onClick={refreshData}>
              <img src={RefreshIcon} alt="Tentar novamente" className="btn-icon" />
              Tentar Novamente
            </button>
            {isTimeoutError && (
              <p className="error-tip">
                💡 <strong>Dica:</strong> Com muitos dados processados, o carregamento pode demorar. 
                O dashboard funcionará normalmente após o processamento ser concluído.
              </p>
            )}
          </div>
        </div>
      </div>
    );
  }

  // No data state
  if (!dashboardData) {
    return (
      <div className={`agent-dashboard no-data ${className}`}>
        <div className="no-data-content">
          <span className="no-data-icon">📊</span>
          <h3>Dados não disponíveis</h3>
          <p>Aguardando dados do agente...</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`agent-dashboard ${compactMode ? 'compact' : ''} ${className}`}>
      {/* Cabeçalho do dashboard */}
      <div className="dashboard-header">
        <div className="agent-status">
          <div className="status-indicator">
            <span className="status-icon">{getStatusIcon(dashboardData.status)}</span>
            <span 
              className="status-text"
              style={{ color: getStatusColor(dashboardData.status) }}
            >
              {getStatusText(dashboardData.status)}
            </span>
          </div>
          <div className="last-updated">
            Atualizado: {dashboardData.lastUpdated.toLocaleTimeString()}
          </div>
        </div>
        
        <div className="dashboard-actions">
          <button 
            className="btn secondary-btn" 
            onClick={refreshData}
            disabled={isLoading}
            title="Atualizar status"
          >
            <img src={RefreshIcon} alt="Atualizar" className="btn-icon" />
            {!compactMode && 'Atualizar'}
          </button>
        </div>
      </div>

      {/* Seção de controles */}
      <div className="control-section">
        <h3>Controles do Agente</h3>
        <div className="control-content">
          {dashboardData.status === 'idle' ? (
            <div className="workflow-controls">
              <select 
                value={selectedWorkflow}
                onChange={(e) => setSelectedWorkflow(e.target.value)}
                className="workflow-select"
              >
                <option value="">Selecionar workflow...</option>
                {dashboardData.availableWorkflows.map(workflow => (
                  <option key={workflow} value={workflow}>
                    {workflow}
                  </option>
                ))}
              </select>
              <button 
                className="btn primary-btn"
                onClick={handleStartWorkflow}
                disabled={!selectedWorkflow}
              >
                <img src={PlayIcon} alt="Executar" className="btn-icon" />
                Executar
              </button>
            </div>
          ) : (
            <div className="workflow-controls">
              <button 
                className="btn warning-btn"
                onClick={handleStopWorkflow}
                disabled={dashboardData.status === 'error'}
              >
                <img src={CancelIcon} alt="Parar" className="btn-icon" />
                Parar Workflow
              </button>
              <button 
                className="btn error-btn"
                onClick={handleCancelWorkflow}
                disabled={dashboardData.status === 'error'}
              >
                <img src={CancelIcon} alt="Cancelar" className="btn-icon" />
                Cancelar
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Progresso do workflow atual */}
      {dashboardData.currentWorkflow && (
        <WorkflowProgress
          workflow={dashboardData.currentWorkflow}
          onCancel={handleCancelWorkflow}
          showDetails={!compactMode}
        />
      )}

      {/* Métricas do agente */}
      {!compactMode && (
        <div className="metrics-section">
          <h3>Métricas do Sistema</h3>
          <div className="metrics-grid">
            <div className="metric-card">
              <span className="metric-value">{dashboardData.metrics.totalWorkflowsExecuted}</span>
              <span className="metric-label">Workflows Executados</span>
            </div>
            <div className="metric-card">
              <span className="metric-value">{dashboardData.metrics.collaboratorsProcessed}</span>
              <span className="metric-label">Colaboradores Processados</span>
            </div>
            <div className="metric-card">
              <span className="metric-value">{dashboardData.metrics.reportsGenerated}</span>
              <span className="metric-label">Relatórios Gerados</span>
            </div>
            <div className="metric-card">
              <span className="metric-value">{dashboardData.metrics.anomaliesDetected}</span>
              <span className="metric-label">Anomalias Detectadas</span>
            </div>
            <div className="metric-card">
              <span className="metric-value">{formatUptime(dashboardData.metrics.uptime)}</span>
              <span className="metric-label">Tempo Ativo</span>
            </div>
            <div className="metric-card">
              <span className="metric-value">
                {dashboardData.metrics.totalWorkflowsExecuted > 0 ? 
                  `${Math.round((dashboardData.metrics.successfulWorkflows / dashboardData.metrics.totalWorkflowsExecuted) * 100)}%` : 
                  '0%'
                }
              </span>
              <span className="metric-label">Taxa de Sucesso</span>
            </div>
          </div>
        </div>
      )}

      {/* Logs do sistema */}
      <div className="logs-section">
        <div className="section-header">
          <h3>Logs Recentes</h3>
          <div className="section-actions">
            <button 
              className="btn secondary-btn"
              onClick={() => setShowLogs(!showLogs)}
            >
              <img src={InfoIcon} alt="Toggle logs" className="btn-icon" />
              {showLogs ? 'Ocultar' : 'Mostrar'} Logs
            </button>
            <button 
              className="btn secondary-btn"
              onClick={clearLogs}
            >
              <img src={CancelIcon} alt="Limpar" className="btn-icon" />
              Limpar
            </button>
          </div>
        </div>
        
        {showLogs && (
          <AgentLogs
            logs={dashboardData.recentLogs}
            maxHeight={compactMode ? '200px' : '300px'}
            showTimestamps={true}
            showSources={true}
          />
        )}
      </div>

      {/* Resumo rápido dos logs (sempre visível) */}
      {!showLogs && (
        <div className="logs-summary">
          {dashboardData.recentLogs.slice(0, 3).map(log => (
            <div key={log.id} className={`log-summary-item level-${log.level}`}>
              <span className="log-time">
                {log.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
              </span>
              <span className="log-content">{log.message}</span>
            </div>
          ))}
          {dashboardData.recentLogs.length === 0 && (
            <div className="no-logs-summary">Nenhum log recente</div>
          )}
        </div>
      )}
    </div>
  );
};

export default AgentDashboard;