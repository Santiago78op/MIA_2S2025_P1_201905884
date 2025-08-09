import React, { useState } from 'react';
import { useFileSystems } from '../hooks/useApi';
import { FileSystemInfo } from '../services/apiService';
import './FileSystemList.css';

const FileSystemList: React.FC = () => {
  const { fileSystems, isLoading, error, currentPath, searchByPath } = useFileSystems();
  const [searchPath, setSearchPath] = useState<string>('/home/julian/Documents/MIA_2S2025_P1_201905884/backend/Discos/mis discos');
  const [isSearching, setIsSearching] = useState(false);

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
      case 'DSK':
        return '💿';
      default:
        return '📋';
    }
  };

  const getStatusIcon = (status?: string): string => {
    return status === 'mounted' ? '🟢' : '🔴';
  };

  const handleSearch = async () => {
    if (!searchPath.trim()) return;
    
    setIsSearching(true);
    try {
      await searchByPath(searchPath.trim());
    } finally {
      setIsSearching(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
  };

  const getCommonPaths = () => [
    './Discos',
    '/home/julian/Documents/MIA_2S2025_P1_201905884/backend/Discos',
    '/home/julian/Documents/MIA_2S2025_P1_201905884/backend/Discos/mis discos',
    '/tmp',
    '/home/julian/Documents',
    './backend/Discos'
  ];

  if (isLoading && !isSearching) {
    return (
      <div className="filesystem-list">
        <h3>💾 Sistemas de Archivos</h3>
        <div className="loading">
          <span className="loading-spinner">🔄</span>
          Cargando sistemas de archivos...
        </div>
      </div>
    );
  }

  return (
    <div className="filesystem-list">
      <div className="list-header">
        <h3>💾 Sistemas de Archivos</h3>
        
        {/* Barra de búsqueda */}
        <div className="search-container">
          <div className="search-input-group">
            <input
              type="text"
              value={searchPath}
              onChange={(e) => setSearchPath(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Ruta donde buscar discos (ej: /home/julian/Documents/...)"
              className="search-input"
              disabled={isSearching}
            />
            <button 
              onClick={handleSearch} 
              className="search-button"
              disabled={isSearching || !searchPath.trim()}
            >
              {isSearching ? '🔍 Buscando...' : '🔍 Buscar Discos'}
            </button>
          </div>

          {/* Rutas comunes */}
          <div className="common-paths">
            <span className="common-paths-label">📁 Rutas frecuentes:</span>
            {getCommonPaths().map((path) => (
              <button
                key={path}
                onClick={() => {
                  setSearchPath(path);
                  searchByPath(path);
                }}
                className="common-path-button"
                disabled={isSearching}
                title={`Buscar en: ${path}`}
              >
                {path.length > 25 ? `...${path.slice(-25)}` : path}
              </button>
            ))}
          </div>
        </div>
      </div>

      {error && (
        <div className="error-message">
          <span>❌ {error}</span>
          <button onClick={handleSearch} className="retry-button">
            🔄 Reintentar
          </button>
        </div>
      )}

      {fileSystems.length === 0 ? (
        <div className="empty-state">
          <span>📂</span>
          <p>No se encontraron discos (.mia, .dsk) en esta ubicación</p>
          <small>📍 Última búsqueda: <code>{currentPath}</code></small>
          <div className="empty-state-suggestions">
            <p>💡 Sugerencias:</p>
            <ul>
              <li>✅ Verifica que la ruta existe y tiene permisos de lectura</li>
              <li>🔧 Crea un disco: <code>mkdisk -size=5 -unit=M -path="{currentPath}/disco1.mia"</code></li>
              <li>📁 Verifica rutas con espacios: asegúrate de escribir la ruta completa</li>
              <li>🔍 Prueba rutas como: /tmp, /home/julian/Documents</li>
            </ul>
          </div>
        </div>
      ) : (
        <>
          <div className="results-summary">
            <span className="results-count">
              📊 {fileSystems.length} disco{fileSystems.length !== 1 ? 's' : ''} encontrado{fileSystems.length !== 1 ? 's' : ''} en:
            </span>
            <code className="results-path">{currentPath}</code>
          </div>
          
          <div className="filesystem-grid">
            {fileSystems.map((fs: FileSystemInfo, index: number) => (
              <div key={index} className="filesystem-card">
                <div className="filesystem-header">
                  <span className="filesystem-icon">{getTypeIcon(fs.type)}</span>
                  <div className="filesystem-info">
                    <h4 className="filesystem-name" title={fs.name}>{fs.name}</h4>
                    <span className="filesystem-type">{fs.type}</span>
                  </div>
                  <span className="status-indicator" title={fs.status === 'mounted' ? 'Montado' : 'No montado'}>
                    {getStatusIcon(fs.status)}
                  </span>
                </div>
                
                <div className="filesystem-details">
                  <div className="detail-row">
                    <span className="detail-label">📏 Tamaño:</span>
                    <span className="detail-value">{formatSize(fs.size)}</span>
                  </div>
                  <div className="detail-row">
                    <span className="detail-label">📌 Montaje:</span>
                    <span className="detail-value">{fs.mountPoint || 'No montado'}</span>
                  </div>
                  {fs.path && (
                    <div className="detail-row">
                      <span className="detail-label">📂 Ubicación:</span>
                      <span className="detail-value path-value" title={fs.path}>
                        {fs.path.length > 30 ? `...${fs.path.slice(-30)}` : fs.path}
                      </span>
                    </div>
                  )}
                  <div className="detail-row">
                    <span className="detail-label">⚡ Estado:</span>
                    <span className={`status-badge ${fs.status}`}>
                      {fs.status === 'mounted' ? 'Montado' : 'No montado'}
                    </span>
                  </div>
                </div>

                <div className="filesystem-actions">
                  <button className="action-button primary" title="Ver información detallada">
                    📊 Detalles
                  </button>
                  <button className="action-button secondary" title="Opciones de administración">
                    ⚙️ Administrar
                  </button>
                </div>
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
};

export default FileSystemList;
