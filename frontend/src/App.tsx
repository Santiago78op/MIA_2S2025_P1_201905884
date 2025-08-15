import React, { useState, useCallback } from 'react';
import ConnectionStatus from './components/ConnectionStatus';
import FileSystemList from './components/FileSystemList';
import CommandExecutor from './components/CommandExecutor';
import Terminal from './components/Terminal';
import Console from './components/Console';
import { useWebSocket } from './hooks/useWebSocket';
import './App.css';
import './AppNew.css';

function App() {
  const [activeView, setActiveView] = useState<'classic' | 'terminal'>('terminal');
  const [layout, setLayout] = useState<'horizontal' | 'vertical'>('horizontal');
  
  // WebSocket para logs en tiempo real
  const {
    isConnected,
    connectionStatus,
    logs,
    clearLogs,
    addLog
  } = useWebSocket({
    autoConnect: true,
    maxLogs: 1000
  });

  // Manejar comando ejecutado desde terminal
  const handleCommandExecuted = useCallback((command: string, result: any) => {
    addLog('SUCCESS', 'TERMINAL', `Comando ejecutado: ${command}`, result);
  }, [addLog]);

  // Manejar errores desde terminal
  const handleError = useCallback((error: string) => {
    addLog('ERROR', 'TERMINAL', error);
  }, [addLog]);

  // Cambiar vista
  const toggleView = () => {
    setActiveView(prev => prev === 'classic' ? 'terminal' : 'classic');
  };

  // Cambiar layout
  const toggleLayout = () => {
    setLayout(prev => prev === 'horizontal' ? 'vertical' : 'horizontal');
  };

  return (
    <div className="App">
      <header className="App-header">
        <div className="header-content">
          <div className="header-title">
            <h1>🗃️ Sistema de Archivos EXT2</h1>
            <p>Simulador de sistema de archivos - MIA 2S2025</p>
          </div>
          
          <div className="header-controls">
            <button 
              onClick={toggleView}
              className={`view-toggle ${activeView}`}
              title={`Cambiar a vista ${activeView === 'classic' ? 'terminal' : 'clásica'}`}
            >
              {activeView === 'classic' ? '📟' : '🖥️'} 
              {activeView === 'classic' ? 'Terminal' : 'Clásica'}
            </button>
            
            {activeView === 'terminal' && (
              <button 
                onClick={toggleLayout}
                className={`layout-toggle ${layout}`}
                title={`Cambiar a layout ${layout === 'horizontal' ? 'vertical' : 'horizontal'}`}
              >
                {layout === 'horizontal' ? '⚏' : '⚋'} 
                {layout === 'horizontal' ? 'Vertical' : 'Horizontal'}
              </button>
            )}
            
            <div className="connection-indicator">
              <span className={`connection-dot ${connectionStatus}`}></span>
              <span className="connection-text">
                {connectionStatus === 'connected' ? 'Conectado' : 
                 connectionStatus === 'connecting' ? 'Conectando...' : 
                 connectionStatus === 'error' ? 'Error' : 'Desconectado'}
              </span>
            </div>
          </div>
        </div>
      </header>
      
      <main className={`App-main ${activeView} ${layout}`}>
        {activeView === 'classic' ? (
          // Vista clásica original
          <div className="classic-view">
            <ConnectionStatus />
            <CommandExecutor />
            <FileSystemList />
          </div>
        ) : (
          // Nueva vista terminal
          <div className={`terminal-view ${layout}`}>
            <div className="left-panel">
              <Terminal 
                onCommandExecuted={handleCommandExecuted}
                onError={handleError}
              />
              
              {/* Panel de información del sistema */}
              <div className="system-info">
                <ConnectionStatus />
                <div className="quick-stats">
                  <div className="stat-item">
                    <span className="stat-label">📊 Logs:</span>
                    <span className="stat-value">{logs.length}</span>
                  </div>
                  <div className="stat-item">
                    <span className="stat-label">🔌 WebSocket:</span>
                    <span className={`stat-value ${connectionStatus}`}>
                      {isConnected ? 'Conectado' : 'Desconectado'}
                    </span>
                  </div>
                  <div className="stat-item">
                    <span className="stat-label">❌ Errores:</span>
                    <span className="stat-value error">
                      {logs.filter(log => log.type === 'ERROR').length}
                    </span>
                  </div>
                </div>
              </div>
            </div>
            
            <div className="right-panel">
              <Console 
                logs={logs}
                onClear={clearLogs}
                maxEntries={1000}
                autoScroll={true}
              />
              
              <FileSystemList />
            </div>
          </div>
        )}
      </main>
      
      <footer className="App-footer">
        <div className="footer-content">
          <div className="footer-info">
            <p>Proyecto MIA - Sistema de Archivos EXT2</p>
            <small>Frontend: React + TypeScript | Backend: Go | WebSocket: Tiempo Real</small>
          </div>
          
          <div className="footer-stats">
            <span className="footer-stat">
              🖥️ Vista: {activeView === 'classic' ? 'Clásica' : 'Terminal'}
            </span>
            {activeView === 'terminal' && (
              <span className="footer-stat">
                📐 Layout: {layout === 'horizontal' ? 'Horizontal' : 'Vertical'}
              </span>
            )}
            <span className="footer-stat">
              📡 WebSocket: {connectionStatus}
            </span>
          </div>
        </div>
      </footer>

      {/* Botón flotante para cambio rápido de vista */}
      <button 
        className="floating-view-toggle"
        onClick={toggleView}
        title={`Cambiar a vista ${activeView === 'classic' ? 'terminal' : 'clásica'}`}
      >
        {activeView === 'classic' ? '📟' : '🖥️'}
      </button>
    </div>
  );
}

export default App;