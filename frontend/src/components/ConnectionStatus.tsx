import React from 'react';
import { useApiConnection } from '../hooks/useApi';
import './ConnectionStatus.css';

const ConnectionStatus: React.FC = () => {
  const { isConnected, isLoading, error, checkConnection } = useApiConnection();

  const getStatusIcon = () => {
    if (isLoading) return '🔄';
    if (isConnected) return '🟢';
    return '🔴';
  };

  const getStatusText = () => {
    if (isLoading) return 'Conectando...';
    if (isConnected) return 'Conectado al servidor';
    return 'Desconectado del servidor';
  };

  const getStatusClass = () => {
    if (isLoading) return 'status-loading';
    if (isConnected) return 'status-connected';
    return 'status-disconnected';
  };

  return (
    <div className={`connection-status ${getStatusClass()}`}>
      <div className="status-main">
        <span className="status-icon">{getStatusIcon()}</span>
        <span className="status-text">{getStatusText()}</span>
        <button 
          className="refresh-button" 
          onClick={checkConnection}
          disabled={isLoading}
          title="Verificar conexión"
        >
          🔄
        </button>
      </div>
      
      {error && (
        <div className="status-error">
          <span>❌ {error}</span>
        </div>
      )}
      
      <div className="status-details">
        <small>Backend: http://localhost:8080</small>
      </div>
    </div>
  );
};

export default ConnectionStatus;
