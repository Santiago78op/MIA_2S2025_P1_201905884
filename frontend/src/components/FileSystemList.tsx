import React from 'react';
import { useFileSystems } from '../hooks/useApi';
import { FileSystemInfo } from '../services/apiService';
import './FileSystemList.css';

const FileSystemList: React.FC = () => {
  const { fileSystems, isLoading, error, refetch } = useFileSystems();

  const formatSize = (bytes: number): string => {
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    if (bytes === 0) return '0 Bytes';
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i];
  };

  const getTypeIcon = (type: string): string => {
    switch (type.toUpperCase()) {
      case 'EXT2':
        return '💾';
      case 'EXT3':
        return '🗃️';
      case 'EXT4':
        return '📁';
      default:
        return '📋';
    }
  };

  if (isLoading) {
    return (
      <div className="filesystem-list">
        <h3>Sistemas de Archivos</h3>
        <div className="loading">
          <span className="loading-spinner">🔄</span>
          Cargando sistemas de archivos...
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="filesystem-list">
        <h3>Sistemas de Archivos</h3>
        <div className="error-message">
          <span>❌ {error}</span>
          <button onClick={refetch} className="retry-button">
            Reintentar
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="filesystem-list">
      <div className="list-header">
        <h3>Sistemas de Archivos</h3>
        <button onClick={refetch} className="refresh-list-button">
          🔄 Actualizar
        </button>
      </div>

      {fileSystems.length === 0 ? (
        <div className="empty-state">
          <span>📂</span>
          <p>No hay sistemas de archivos disponibles</p>
          <small>Crea una partición para comenzar</small>
        </div>
      ) : (
        <div className="filesystem-grid">
          {fileSystems.map((fs: FileSystemInfo, index: number) => (
            <div key={index} className="filesystem-card">
              <div className="filesystem-header">
                <span className="filesystem-icon">{getTypeIcon(fs.type)}</span>
                <div className="filesystem-info">
                  <h4 className="filesystem-name">{fs.name}</h4>
                  <span className="filesystem-type">{fs.type}</span>
                </div>
              </div>
              
              <div className="filesystem-details">
                <div className="detail-row">
                  <span className="detail-label">Tamaño:</span>
                  <span className="detail-value">{formatSize(fs.size)}</span>
                </div>
                <div className="detail-row">
                  <span className="detail-label">Punto de montaje:</span>
                  <span className="detail-value">{fs.mountPoint}</span>
                </div>
              </div>

              <div className="filesystem-actions">
                <button className="action-button primary">
                  📊 Ver detalles
                </button>
                <button className="action-button secondary">
                  ⚙️ Administrar
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default FileSystemList;
