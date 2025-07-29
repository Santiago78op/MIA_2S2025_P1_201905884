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

// Hook para manejar sistemas de archivos
export const useFileSystems = () => {
  const [fileSystems, setFileSystems] = useState<FileSystemInfo[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchFileSystems = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await apiService.getFileSystems();
      setFileSystems(response.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al obtener sistemas de archivos');
      setFileSystems([]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchFileSystems();
  }, [fetchFileSystems]);

  return {
    fileSystems,
    isLoading,
    error,
    refetch: fetchFileSystems,
  };
};

// Hook para ejecutar comandos
export const useCommandExecution = () => {
  const [isExecuting, setIsExecuting] = useState(false);
  const [lastResult, setLastResult] = useState<ApiResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  const executeCommand = useCallback(async (command: string) => {
    try {
      setIsExecuting(true);
      setError(null);
      const response = await apiService.executeCommand(command);
      setLastResult(response);
      return response;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Error al ejecutar comando';
      setError(errorMessage);
      throw new Error(errorMessage);
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
