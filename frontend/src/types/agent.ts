// Agent monitoring types for the VR/VA automation system

export type AgentStatus = 'idle' | 'running' | 'error' | 'cancelled';

export type LogLevel = 'debug' | 'info' | 'warn' | 'error';

export interface LogEntry {
  id: string;
  timestamp: Date;
  level: LogLevel;
  message: string;
  source?: string; // Nome do step/componente que gerou o log
  context?: Record<string, any>; // Dados contextuais adicionais
}

export interface WorkflowStep {
  name: string;
  description: string;
  status: 'pending' | 'running' | 'completed' | 'failed' | 'skipped';
  startTime?: Date;
  endTime?: Date;
  duration?: number; // em millisegundos
  error?: string;
}

export interface WorkflowStatus {
  id: string;
  name: string;
  description: string;
  status: AgentStatus;
  currentStepIndex: number;
  totalSteps: number;
  steps: WorkflowStep[];
  startTime?: Date;
  endTime?: Date;
  totalDuration?: number; // em millisegundos
  progress: number; // 0-100
  error?: string;
  metadata?: Record<string, any>;
}

export interface AgentMetrics {
  totalWorkflowsExecuted: number;
  successfulWorkflows: number;
  failedWorkflows: number;
  averageExecutionTime: number; // em millisegundos
  totalProcessingTime: number; // em millisegundos
  collaboratorsProcessed: number;
  reportsGenerated: number;
  anomaliesDetected: number;
  uptime: number; // em millisegundos
  memoryUsage?: number; // em MB
  cpuUsage?: number; // percentual
}

export interface AgentDashboardData {
  status: AgentStatus;
  currentWorkflow?: WorkflowStatus;
  recentLogs: LogEntry[];
  metrics: AgentMetrics;
  availableWorkflows: string[];
  lastUpdated: Date;
}

export interface WorkflowExecutionRequest {
  workflowName: string;
  parameters?: Record<string, any>;
}

export interface WorkflowControlCommand {
  action: 'start' | 'stop' | 'cancel' | 'pause' | 'resume';
  workflowId?: string;
  parameters?: Record<string, any>;
}

// Eventos em tempo real
export type AgentEventType = 
  | 'workflow_started'
  | 'workflow_completed' 
  | 'workflow_failed'
  | 'step_started'
  | 'step_completed'
  | 'step_failed'
  | 'log_entry'
  | 'metrics_updated'
  | 'status_changed';

export interface AgentEvent {
  type: AgentEventType;
  timestamp: Date;
  data: any;
  workflowId?: string;
  stepName?: string;
}

// Configurações do dashboard
export interface DashboardConfig {
  refreshInterval: number; // em millisegundos
  maxLogs: number; // máximo de logs para manter
  logLevels: LogLevel[]; // níveis de log a exibir
  showMetrics: boolean;
  compactMode: boolean;
  soundNotifications: boolean;
}

// Filtros para logs
export interface LogFilter {
  levels?: LogLevel[];
  source?: string;
  dateFrom?: Date;
  dateTo?: Date;
  searchText?: string;
}

// Estado do hook useAgentStatus
export interface AgentStatusState {
  dashboardData: AgentDashboardData | null;
  isConnected: boolean;
  isLoading: boolean;
  error: string | null;
  config: DashboardConfig;
}

// Ações do hook useAgentStatus
export interface AgentStatusActions {
  startWorkflow: (request: WorkflowExecutionRequest) => Promise<void>;
  stopWorkflow: (workflowId?: string) => Promise<void>;
  cancelWorkflow: (workflowId?: string) => Promise<void>;
  refreshData: () => Promise<void>;
  clearLogs: () => void;
  filterLogs: (filter: LogFilter) => LogEntry[];
  updateConfig: (config: Partial<DashboardConfig>) => void;
  exportLogs: (format: 'json' | 'csv' | 'txt') => void;
}

// Props dos componentes
export interface AgentDashboardProps {
  className?: string;
  compactMode?: boolean;
  refreshInterval?: number;
}

export interface WorkflowProgressProps {
  workflow: WorkflowStatus;
  onCancel?: () => void;
  showDetails?: boolean;
  className?: string;
}

export interface AgentLogsProps {
  logs: LogEntry[];
  filter?: LogFilter;
  onFilterChange?: (filter: LogFilter) => void;
  maxHeight?: string;
  showTimestamps?: boolean;
  showSources?: boolean;
  className?: string;
}

export interface AgentControlsProps {
  status: AgentStatus;
  availableWorkflows: string[];
  onStartWorkflow: (workflowName: string) => void;
  onStopWorkflow: () => void;
  onRefresh: () => void;
  disabled?: boolean;
  className?: string;
}

export interface AgentMetricsProps {
  metrics: AgentMetrics;
  showDetailed?: boolean;
  className?: string;
}

// Utility functions types
export type LogLevelColor = Record<LogLevel, string>;
export type StatusColor = Record<AgentStatus, string>;

// Mock data para desenvolvimento
export interface MockDataConfig {
  enabled: boolean;
  workflowDuration: number; // duração simulada do workflow em ms
  logFrequency: number; // frequência de logs simulados em ms
  randomErrors: boolean; // simular erros aleatórios
}