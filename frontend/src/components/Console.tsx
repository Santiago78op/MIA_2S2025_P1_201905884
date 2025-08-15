import React, { useState, useRef, useEffect } from 'react';
import './Console.css';

export interface LogEntry {
  id: string;
  timestamp: string;
  type: 'INFO' | 'WARNING' | 'ERROR' | 'SUCCESS' | 'SYSTEM';
  command: string;
  message: string;
  data?: any;
}

interface ConsoleProps {
  logs: LogEntry[];
  onClear?: () => void;
  maxEntries?: number;
  autoScroll?: boolean;
}

const Console: React.FC<ConsoleProps> = ({ 
  logs, 
  onClear, 
  maxEntries = 1000,
  autoScroll = true 
}) => {
  const [isVisible, setIsVisible] = useState(true);
  const [filter, setFilter] = useState<string>('ALL');
  const [searchTerm, setSearchTerm] = useState('');
  const [isAutoScroll, setIsAutoScroll] = useState(autoScroll);
  const [selectedEntry, setSelectedEntry] = useState<LogEntry | null>(null);
  
  const consoleRef = useRef<HTMLDivElement>(null);
  const endRef = useRef<HTMLDivElement>(null);

  // Filtros disponibles
  const logFilters = ['ALL', 'INFO', 'WARNING', 'ERROR', 'SUCCESS', 'SYSTEM'];

  // Auto-scroll al final cuando hay nuevos logs
  useEffect(() => {
    if (isAutoScroll && endRef.current) {
      endRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [logs, isAutoScroll]);

  // Filtrar logs
  const filteredLogs = logs.filter(log => {
    const matchesFilter = filter === 'ALL' || log.type === filter;
    const matchesSearch = searchTerm === '' || 
      log.message.toLowerCase().includes(searchTerm.toLowerCase()) ||
      log.command.toLowerCase().includes(searchTerm.toLowerCase());
    
    return matchesFilter && matchesSearch;
  }).slice(-maxEntries); // Limitar número de entradas

  // Obtener icono según el tipo de log
  const getLogIcon = (type: LogEntry['type']) => {
    switch (type) {
      case 'INFO': return '📘';
      case 'WARNING': return '⚠️';
      case 'ERROR': return '❌';
      case 'SUCCESS': return '✅';
      case 'SYSTEM': return '🖥️';
      default: return '📝';
    }
  };

  // Obtener color según el tipo de log
  const getLogColor = (type: LogEntry['type']) => {
    switch (type) {
      case 'INFO': return '#2196F3';
      case 'WARNING': return '#FF9800';
      case 'ERROR': return '#f44336';
      case 'SUCCESS': return '#4CAF50';
      case 'SYSTEM': return '#9C27B0';
      default: return '#666';
    }
  };

  // Formatear timestamp
  const formatTimestamp = (timestamp: string) => {
    try {
      const date = new Date(timestamp);
      return date.toLocaleTimeString();
    } catch {
      return timestamp;
    }
  };

  // Exportar logs
  const exportLogs = () => {
    const logsText = filteredLogs.map(log => 
      `[${formatTimestamp(log.timestamp)}] [${log.type}] ${log.command}: ${log.message}`
    ).join('\n');

    const blob = new Blob([logsText], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `mia-logs-${new Date().toISOString().split('T')[0]}.txt`;
    a.click();
    URL.revokeObjectURL(url);
  };

  // Copiar log al portapapeles
  const copyLog = (log: LogEntry) => {
    const logText = `[${formatTimestamp(log.timestamp)}] [${log.type}] ${log.command}: ${log.message}`;
    navigator.clipboard.writeText(logText).then(() => {
      // Mostrar feedback visual
      setSelectedEntry(log);
      setTimeout(() => setSelectedEntry(null), 1000);
    });
  };

  // Contar logs por tipo
  const logCounts = logs.reduce((acc, log) => {
    acc[log.type] = (acc[log.type] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  if (!isVisible) {
    return (
      <div className="console-minimized" onClick={() => setIsVisible(true)}>
        <span>📊 Console ({logs.length} logs)</span>
        {logs.filter(log => log.type === 'ERROR').length > 0 && (
          <span className="error-badge">
            {logs.filter(log => log.type === 'ERROR').length} errores
          </span>
        )}
      </div>
    );
  }

  return (
    <div className="console-container">
      <div className="console-header">
        <div className="console-title">
          <span className="console-icon">📊</span>
          <span>Console ({filteredLogs.length}/{logs.length})</span>
        </div>
        
        <div className="console-stats">
          {logFilters.slice(1).map(type => (
            <span 
              key={type} 
              className={`stat-badge ${type.toLowerCase()}`}
              title={`${type}: ${logCounts[type] || 0} entradas`}
            >
              {getLogIcon(type as LogEntry['type'])} {logCounts[type] || 0}
            </span>
          ))}
        </div>

        <div className="console-controls">
          <button 
            className={`console-control auto-scroll ${isAutoScroll ? 'active' : ''}`}
            onClick={() => setIsAutoScroll(!isAutoScroll)}
            title="Auto-scroll"
          >
            📜
          </button>
          <button 
            className="console-control export"
            onClick={exportLogs}
            title="Exportar logs"
          >
            💾
          </button>
          <button 
            className="console-control clear"
            onClick={onClear}
            title="Limpiar console"
          >
            🗑️
          </button>
          <button 
            className="console-control minimize"
            onClick={() => setIsVisible(false)}
            title="Minimizar"
          >
            —
          </button>
        </div>
      </div>

      <div className="console-filters">
        <div className="filter-group">
          <label>Filtro:</label>
          <select 
            value={filter} 
            onChange={(e) => setFilter(e.target.value)}
            className="filter-select"
          >
            {logFilters.map(filterType => (
              <option key={filterType} value={filterType}>
                {filterType === 'ALL' ? 'Todos' : filterType}
              </option>
            ))}
          </select>
        </div>

        <div className="search-group">
          <label>Buscar:</label>
          <input
            type="text"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder="Buscar en logs..."
            className="search-input"
          />
          {searchTerm && (
            <button 
              onClick={() => setSearchTerm('')}
              className="clear-search"
              title="Limpiar búsqueda"
            >
              ✕
            </button>
          )}
        </div>
      </div>

      <div className="console-body" ref={consoleRef}>
        {filteredLogs.length === 0 ? (
          <div className="console-empty">
            <div className="empty-icon">📝</div>
            <div className="empty-message">
              {logs.length === 0 
                ? 'No hay logs disponibles'
                : 'No hay logs que coincidan con los filtros'
              }
            </div>
            {searchTerm && (
              <button 
                onClick={() => setSearchTerm('')}
                className="clear-filters-btn"
              >
                Limpiar filtros
              </button>
            )}
          </div>
        ) : (
          <>
            {filteredLogs.map((log, index) => (
              <div 
                key={log.id || index} 
                className={`console-entry ${log.type.toLowerCase()} ${
                  selectedEntry?.id === log.id ? 'copied' : ''
                }`}
                onClick={() => copyLog(log)}
                title="Click para copiar"
              >
                <div className="entry-header">
                  <span className="entry-icon">
                    {getLogIcon(log.type)}
                  </span>
                  <span className="entry-timestamp">
                    {formatTimestamp(log.timestamp)}
                  </span>
                  <span 
                    className="entry-type"
                    style={{ color: getLogColor(log.type) }}
                  >
                    {log.type}
                  </span>
                  <span className="entry-command">
                    {log.command}
                  </span>
                </div>
                
                <div className="entry-message">
                  {log.message}
                </div>

                {log.data && (
                  <details className="entry-data">
                    <summary>Ver datos adicionales</summary>
                    <pre className="data-content">
                      {typeof log.data === 'string' 
                        ? log.data 
                        : JSON.stringify(log.data, null, 2)
                      }
                    </pre>
                  </details>
                )}

                {selectedEntry?.id === log.id && (
                  <div className="copy-feedback">
                    ✅ Copiado al portapapeles
                  </div>
                )}
              </div>
            ))}
            <div ref={endRef} />
          </>
        )}
      </div>

      <div className="console-footer">
        <div className="footer-info">
          <span>
            📈 {filteredLogs.length} entradas mostradas
          </span>
          {searchTerm && (
            <span>
              🔍 Filtrado por: "{searchTerm}"
            </span>
          )}
        </div>
        <div className="footer-actions">
          <button 
            onClick={() => endRef.current?.scrollIntoView({ behavior: 'smooth' })}
            className="scroll-to-bottom"
            title="Ir al final"
          >
            ⬇️
          </button>
        </div>
      </div>
    </div>
  );
};

export default Console;