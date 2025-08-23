import React, { useState, useRef, useEffect } from 'react';
import { AgentLogsProps, LogEntry, LogLevel, LogFilter } from '../types/agent';
import './AgentLogs.css';

// Importando ícones existentes
import ArrowDownIcon from '../assets/icons/arrow-down.svg';
import ClearIcon from '../assets/icons/clear.svg';
import SaveIcon from '../assets/icons/save.svg';

const AgentLogs: React.FC<AgentLogsProps> = ({
  logs,
  filter,
  onFilterChange,
  maxHeight = '400px',
  showTimestamps = true,
  showSources = true,
  className = ''
}) => {
  const [localFilter, setLocalFilter] = useState<LogFilter>(filter || {});
  const [autoScroll, setAutoScroll] = useState(true);
  const [searchText, setSearchText] = useState('');
  const logsEndRef = useRef<HTMLDivElement>(null);
  const logsContainerRef = useRef<HTMLDivElement>(null);

  // Auto-scroll para o final quando novos logs chegam
  useEffect(() => {
    if (autoScroll && logsEndRef.current) {
      logsEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [logs, autoScroll]);

  // Detectar se o usuário scrollou para cima (desabilitar auto-scroll)
  const handleScroll = () => {
    if (!logsContainerRef.current) return;
    
    const { scrollTop, scrollHeight, clientHeight } = logsContainerRef.current;
    const isNearBottom = scrollTop + clientHeight >= scrollHeight - 50;
    setAutoScroll(isNearBottom);
  };

  const getLevelColor = (level: LogLevel): string => {
    switch (level) {
      case 'error':
        return 'var(--error-color)';
      case 'warn':
        return 'var(--warning-color)';
      case 'info':
        return 'var(--info-color)';
      case 'debug':
        return 'var(--text-secondary)';
      default:
        return 'var(--text-primary)';
    }
  };

  const getLevelIcon = (level: LogLevel): string => {
    switch (level) {
      case 'error':
        return '❌';
      case 'warn':
        return '⚠️';
      case 'info':
        return 'ℹ️';
      case 'debug':
        return '🔍';
      default:
        return '📝';
    }
  };

  const formatTimestamp = (date: Date): string => {
    return date.toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit', 
      second: '2-digit',
      fractionalSecondDigits: 3
    });
  };

  const formatDate = (date: Date): string => {
    return date.toLocaleDateString();
  };

  const handleLevelToggle = (level: LogLevel) => {
    const currentLevels = localFilter.levels || ['info', 'warn', 'error'];
    const newLevels = currentLevels.includes(level)
      ? currentLevels.filter(l => l !== level)
      : [...currentLevels, level];
    
    const newFilter = { ...localFilter, levels: newLevels };
    setLocalFilter(newFilter);
    onFilterChange?.(newFilter);
  };

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value;
    setSearchText(text);
    
    const newFilter = { ...localFilter, searchText: text || undefined };
    setLocalFilter(newFilter);
    onFilterChange?.(newFilter);
  };

  const clearLogs = () => {
    // Esta função seria chamada pelo componente pai
    console.log('Clear logs requested');
  };

  const exportLogs = () => {
    const dataStr = JSON.stringify(filteredLogs, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(dataBlob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `agent-logs-${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  };

  const scrollToBottom = () => {
    setAutoScroll(true);
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  // Filtrar logs
  const filteredLogs = logs.filter(log => {
    if (localFilter.levels && !localFilter.levels.includes(log.level)) return false;
    if (localFilter.source && log.source !== localFilter.source) return false;
    if (localFilter.searchText && !log.message.toLowerCase().includes(localFilter.searchText.toLowerCase())) return false;
    if (localFilter.dateFrom && log.timestamp < localFilter.dateFrom) return false;
    if (localFilter.dateTo && log.timestamp > localFilter.dateTo) return false;
    return true;
  });

  // Agrupar logs por data
  const logsByDate = filteredLogs.reduce((acc, log) => {
    const dateKey = formatDate(log.timestamp);
    if (!acc[dateKey]) acc[dateKey] = [];
    acc[dateKey].push(log);
    return acc;
  }, {} as Record<string, LogEntry[]>);

  const allLevels: LogLevel[] = ['debug', 'info', 'warn', 'error'];
  const activeLevels = localFilter.levels || ['info', 'warn', 'error'];

  return (
    <div className={`agent-logs ${className}`}>
      {/* Cabeçalho dos logs */}
      <div className="logs-header">
        <h3 className="logs-title">Logs do Sistema</h3>
        
        <div className="logs-controls">
          <div className="search-box">
            <input
              type="text"
              placeholder="Buscar nos logs..."
              value={searchText}
              onChange={handleSearchChange}
              className="search-input"
            />
          </div>

          <div className="level-filters">
            {allLevels.map(level => (
              <button
                key={level}
                className={`level-filter ${activeLevels.includes(level) ? 'active' : ''}`}
                onClick={() => handleLevelToggle(level)}
              >
                <span>{getLevelIcon(level)}</span>
                <span>{level.toUpperCase()}</span>
              </button>
            ))}
          </div>

          <div className="logs-actions">
            <button className="btn" onClick={exportLogs} title="Exportar logs">
              <img src={SaveIcon} alt="Exportar" className="btn-icon" />
            </button>
            <button className="btn" onClick={clearLogs} title="Limpar logs">
              <img src={ClearIcon} alt="Limpar" className="btn-icon" />
            </button>
          </div>
        </div>
      </div>

      {/* Estatísticas */}
      <div className="logs-stats">
        <span>
          {filteredLogs.length} de {logs.length} logs exibidos
        </span>
        <span>
          Auto-scroll: {autoScroll ? 'Ativo' : 'Inativo'}
        </span>
      </div>

      {/* Container dos logs */}
      <div 
        className="logs-container"
        style={{ maxHeight }}
        ref={logsContainerRef}
        onScroll={handleScroll}
      >
        {filteredLogs.length === 0 ? (
          <div className="no-logs">
            <div className="no-logs-icon">📋</div>
            <div>Nenhum log encontrado</div>
            {localFilter.searchText && (
              <div>Tente ajustar os filtros ou o termo de busca</div>
            )}
          </div>
        ) : (
          Object.entries(logsByDate).map(([date, dateLogs]) => (
            <div key={date} className="date-group">
              <div className="date-header">{date}</div>
              {dateLogs.map(log => (
                <div key={log.id} className={`log-entry level-${log.level}`}>
                  <span className="log-icon">{getLevelIcon(log.level)}</span>
                  
                  {showTimestamps && (
                    <span className="log-timestamp">
                      {formatTimestamp(log.timestamp)}
                    </span>
                  )}
                  
                  <span 
                    className={`log-level level-${log.level}`}
                    style={{ color: getLevelColor(log.level) }}
                  >
                    {log.level}
                  </span>
                  
                  {showSources && log.source && (
                    <span className="log-source">[{log.source}]</span>
                  )}
                  
                  <span className="log-message">{log.message}</span>
                </div>
              ))}
            </div>
          ))
        )}
        
        <div ref={logsEndRef} />
      </div>

      {/* Botão scroll to bottom */}
      <button 
        className={`scroll-to-bottom ${autoScroll ? 'hidden' : ''}`}
        onClick={scrollToBottom}
        title="Ir para o final"
      >
        <img src={ArrowDownIcon} alt="Ir para baixo" className="btn-icon" />
      </button>
    </div>
  );
};

export default AgentLogs;