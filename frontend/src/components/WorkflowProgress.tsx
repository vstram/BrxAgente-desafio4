import React from 'react';
import { WorkflowProgressProps, WorkflowStep } from '../types/agent';
import './WorkflowProgress.css';

// Importando ícones existentes
import CancelIcon from '../assets/icons/cancel.svg';
import SpinnerIcon from '../assets/icons/spinner.svg';

const WorkflowProgress: React.FC<WorkflowProgressProps> = ({
  workflow,
  onCancel,
  showDetails = true,
  className = ''
}) => {
  const formatDuration = (ms?: number): string => {
    if (!ms) return '0s';
    const seconds = Math.floor(ms / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    
    if (hours > 0) {
      return `${hours}h ${minutes % 60}m ${seconds % 60}s`;
    } else if (minutes > 0) {
      return `${minutes}m ${seconds % 60}s`;
    }
    return `${seconds}s`;
  };

  const formatTimestamp = (date?: Date): string => {
    if (!date) return '--:--';
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  const getStepStatusIcon = (step: WorkflowStep): string => {
    switch (step.status) {
      case 'running':
        return '⏳';
      case 'completed':
        return '✅';
      case 'failed':
        return '❌';
      case 'skipped':
        return '⏭️';
      default:
        return '⏸️';
    }
  };

  const getStepStatusColor = (step: WorkflowStep): string => {
    switch (step.status) {
      case 'running':
        return 'var(--primary-color)';
      case 'completed':
        return 'var(--success-color)';
      case 'failed':
        return 'var(--error-color)';
      case 'skipped':
        return 'var(--warning-color)';
      default:
        return 'var(--text-secondary)';
    }
  };

  const currentStepDescription = workflow.steps[workflow.currentStepIndex]?.description || 'Workflow finalizado';
  const elapsedTime = workflow.startTime ? Date.now() - workflow.startTime.getTime() : 0;

  return (
    <div className={`workflow-progress status-${workflow.status} ${className}`}>
      {/* Cabeçalho do workflow */}
      <div className="workflow-header">
        <div className="workflow-info">
          <h3>{workflow.name}</h3>
          <p>{workflow.description}</p>
        </div>
        {onCancel && workflow.status === 'running' && (
          <div className="workflow-actions">
            <button className="btn" onClick={onCancel}>
              <img src={CancelIcon} alt="Cancelar" className="btn-icon" />
              Cancelar
            </button>
          </div>
        )}
      </div>

      {/* Barra de progresso */}
      <div className="progress-section">
        <div className="progress-info">
          <div className="current-step">
            {workflow.status === 'running' && (
              <img src={SpinnerIcon} alt="Executando" className="btn-icon spinner" />
            )}
            <span>
              {workflow.status === 'running' 
                ? `Executando: ${currentStepDescription}`
                : workflow.status === 'error'
                ? 'Erro na execução'
                : workflow.status === 'cancelled'
                ? 'Workflow cancelado'
                : 'Workflow concluído'
              }
            </span>
          </div>
          <span>{workflow.progress}%</span>
        </div>
        
        <div className="progress-bar-container">
          <div 
            className={`progress-bar ${workflow.status === 'running' ? 'animated' : ''}`}
            style={{ width: `${workflow.progress}%` }}
          />
        </div>

        <div className="progress-info">
          <span>Step {workflow.currentStepIndex + 1} de {workflow.totalSteps}</span>
          <span>Tempo: {formatDuration(elapsedTime)}</span>
        </div>
      </div>

      {/* Estatísticas */}
      <div className="stats">
        <div className="stat-item">
          <span className="stat-value">{workflow.currentStepIndex + (workflow.status === 'idle' ? 0 : 1)}</span>
          <span className="stat-label">Steps Executados</span>
        </div>
        <div className="stat-item">
          <span className="stat-value">{formatTimestamp(workflow.startTime)}</span>
          <span className="stat-label">Iniciado em</span>
        </div>
        <div className="stat-item">
          <span className="stat-value">
            {workflow.totalDuration ? formatDuration(workflow.totalDuration) : formatDuration(elapsedTime)}
          </span>
          <span className="stat-label">Duração</span>
        </div>
        <div className="stat-item">
          <span className="stat-value status-indicator">
            {workflow.status.toUpperCase()}
          </span>
          <span className="stat-label">Status</span>
        </div>
      </div>

      {/* Lista detalhada de steps */}
      {showDetails && (
        <div className="steps-list">
          <div className="steps-header">Progresso Detalhado:</div>
          {workflow.steps.map((step, index) => (
            <div 
              key={step.name} 
              className={`step-item ${index === workflow.currentStepIndex ? 'current' : ''}`}
            >
              <div className="step-icon">{getStepStatusIcon(step)}</div>
              <div className="step-info">
                <div className="step-name" style={{ color: getStepStatusColor(step) }}>
                  {step.description || step.name}
                </div>
                {step.error && (
                  <div className="step-description error-text">
                    Erro: {step.error}
                  </div>
                )}
              </div>
              <div className="step-duration">
                {step.duration ? formatDuration(step.duration) : 
                 step.status === 'running' ? formatDuration(elapsedTime) : '--'}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Mensagem de erro geral */}
      {workflow.error && (
        <div className="error-message">
          <strong>Erro no Workflow:</strong> {workflow.error}
        </div>
      )}
    </div>
  );
};

export default WorkflowProgress;