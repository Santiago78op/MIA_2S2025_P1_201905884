import { useState, useEffect, useCallback } from 'react';
import apiService, { ApiResponse, FileSystemInfo, HealthResponse } from '../services/apiService';

// Hook para verificar la conectividad
export const useApiConnection = () => {
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const checkConnection = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const connected = await apiService.testConnection();
      setIsConnected(connected);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error de conexión');
      setIsConnected(false);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    checkConnection();
    // Verificar conexión cada 30 segundos
    const interval = setInterval(checkConnection, 30000);
    return () => clearInterval(interval);
  }, [checkConnection]);

  return {
    isConnected,
    isLoading,
    error,
    checkConnection,
  };
};

// Hook para obtener el estado del servidor
export const useServerHealth = () => {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchHealth = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await apiService.checkHealth();
      setHealth(response.data || null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al obtener estado del servidor');
      setHealth(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchHealth();
  }, [fetchHealth]);

  return {
    health,
    isLoading,
    error,
    refetch: fetchHealth,
  };
};

// Hook para manejar sistemas de archivos con búsqueda
export const useFileSystems = () => {
  const [fileSystems, setFileSystems] = useState<FileSystemInfo[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentPath, setCurrentPath] = useState<string>('/home/julian/Documents/MIA_2S2025_P1_201905884/backend/Discos/mis discos');

  const fetchFileSystems = useCallback(async (searchPath?: string) => {
    try {
      setIsLoading(true);
      setError(null);
      const pathToSearch = searchPath || currentPath;
      console.log('Buscando en ruta:', pathToSearch); // Debug
      const response = await apiService.getFileSystems(pathToSearch);
      setFileSystems(response.data || []);
      if (searchPath) {
        setCurrentPath(searchPath);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Error al obtener sistemas de archivos';
      console.error('Error en fetchFileSystems:', errorMessage); // Debug
      setError(errorMessage);
      setFileSystems([]);
    } finally {
      setIsLoading(false);
    }
  }, [currentPath]);

  useEffect(() => {
    fetchFileSystems();
  }, [fetchFileSystems]);

  return {
    fileSystems,
    isLoading,
    error,
    currentPath,
    refetch: fetchFileSystems,
    searchByPath: (path: string) => fetchFileSystems(path),
  };
};

// Hook para ejecutar comandos con manejo mejorado de errores
export const useCommandExecution = () => {
  const [isExecuting, setIsExecuting] = useState(false);
  const [lastResult, setLastResult] = useState<ApiResponse | null>(null);
  const [error, setError] = useState<any | null>(null);

  const executeCommand = useCallback(async (command: string) => {
    try {
      setIsExecuting(true);
      setError(null);
      const response = await apiService.executeCommand(command);
      setLastResult(response);
      return response;
    } catch (err: any) {
      // Preservar toda la información del error enriquecido
      const enrichedError = {
        message: err.message || 'Error al ejecutar comando',
        title: err.title || 'Error',
        code: err.code || 'UNKNOWN',
        status: err.status,
        suggestions: err.suggestions || [],
        timestamp: err.timestamp || new Date().toISOString(),
        originalError: err.originalError || err
      };
      
      setError(enrichedError);
      throw enrichedError;
    } finally {
      setIsExecuting(false);
    }
  }, []);

  return {
    executeCommand,
    isExecuting,
    lastResult,
    error,
    clearError: () => setError(null),
  };
};
