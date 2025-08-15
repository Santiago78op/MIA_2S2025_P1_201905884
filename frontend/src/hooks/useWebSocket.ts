import { useState, useEffect, useRef, useCallback } from 'react';
import { LogEntry } from '../components/Console';

interface WebSocketState {
  isConnected: boolean;
  logs: LogEntry[];
  connectionStatus: 'connecting' | 'connected' | 'disconnected' | 'error';
  lastMessage: any;
  error: string | null;
}

interface UseWebSocketOptions {
  url?: string;
  autoConnect?: boolean;
  maxLogs?: number;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
}

export const useWebSocket = (options: UseWebSocketOptions = {}) => {
  const {
    url = 'ws://localhost:8080/api/ws',
    autoConnect = true,
    maxLogs = 1000,
    reconnectInterval = 3000,
    maxReconnectAttempts = 5
  } = options;

  const [state, setState] = useState<WebSocketState>({
    isConnected: false,
    logs: [],
    connectionStatus: 'disconnected',
    lastMessage: null,
    error: null
  });

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const isManuallyClosedRef = useRef(false);

  // Generar ID único para logs
  const generateLogId = useCallback(() => {
    return `log_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }, []);

  // Limpiar timeouts
  const clearReconnectTimeout = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
  }, []);

  // Conectar WebSocket
  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      return;
    }

    try {
      setState(prev => ({ 
        ...prev, 
        connectionStatus: 'connecting',
        error: null 
      }));

      wsRef.current = new WebSocket(url);

      wsRef.current.onopen = () => {
        console.log('WebSocket conectado');
        setState(prev => ({
          ...prev,
          isConnected: true,
          connectionStatus: 'connected',
          error: null
        }));
        
        reconnectAttemptsRef.current = 0;
        clearReconnectTimeout();

        // Agregar log de conexión
        const connectionLog: LogEntry = {
          id: generateLogId(),
          timestamp: new Date().toISOString(),
          type: 'SYSTEM',
          command: 'WEBSOCKET',
          message: 'Conexión WebSocket establecida correctamente'
        };

        setState(prev => ({
          ...prev,
          logs: [...prev.logs.slice(-maxLogs + 1), connectionLog]
        }));
      };

      wsRef.current.onmessage = (event) => {
        try {
          let data: any;
          
          // Intentar parsear como JSON primero
          try {
            data = JSON.parse(event.data);
          } catch (jsonError) {
            // Si no es JSON válido, tratar como mensaje de texto
            // Intentar extraer información del formato: [timestamp] [type] command: message
            const textMessage = event.data;
            const match = textMessage.match(/\[(\d+)\] \[(\w+)\] ([^:]+): (.+)/);
            
            if (match) {
              data = {
                time: match[1],
                type: match[2],
                command: match[3],
                message: match[4]
              };
            } else {
              // Si no coincide con el patrón, crear un log genérico
              data = {
                type: 'INFO',
                command: 'SYSTEM',
                message: textMessage,
                time: Math.floor(Date.now() / 1000).toString()
              };
            }
          }
          
          setState(prev => ({
            ...prev,
            lastMessage: data
          }));

          // Convertir mensaje a LogEntry
          const logEntry: LogEntry = {
            id: generateLogId(),
            timestamp: data.time ? new Date(parseInt(data.time) * 1000).toISOString() : new Date().toISOString(),
            type: data.type || 'INFO',
            command: data.command || 'SYSTEM',
            message: data.message || event.data,
            data: data.data || null
          };

          setState(prev => ({
            ...prev,
            logs: [...prev.logs.slice(-maxLogs + 1), logEntry]
          }));

        } catch (error) {
          console.error('Error processing WebSocket message:', error);
          
          // Agregar mensaje como log de error solo si realmente hay un error grave
          const errorLog: LogEntry = {
            id: generateLogId(),
            timestamp: new Date().toISOString(),
            type: 'ERROR',
            command: 'WEBSOCKET',
            message: `Error processing message: ${error instanceof Error ? error.message : 'Unknown error'}`,
            data: { originalMessage: event.data, error: error instanceof Error ? error.message : 'Unknown error' }
          };

          setState(prev => ({
            ...prev,
            logs: [...prev.logs.slice(-maxLogs + 1), errorLog]
          }));
        }
      };

      wsRef.current.onclose = (event) => {
        console.log('WebSocket desconectado', event.code, event.reason);
        
        setState(prev => ({
          ...prev,
          isConnected: false,
          connectionStatus: 'disconnected'
        }));

        // Agregar log de desconexión
        const disconnectionLog: LogEntry = {
          id: generateLogId(),
          timestamp: new Date().toISOString(),
          type: 'WARNING',
          command: 'WEBSOCKET',
          message: `Conexión WebSocket cerrada (código: ${event.code}${event.reason ? `, razón: ${event.reason}` : ''})`
        };

        setState(prev => ({
          ...prev,
          logs: [...prev.logs.slice(-maxLogs + 1), disconnectionLog]
        }));

        // Intentar reconectar si no fue cerrado manualmente
        if (!isManuallyClosedRef.current && reconnectAttemptsRef.current < maxReconnectAttempts) {
          reconnectAttemptsRef.current++;
          
          const reconnectLog: LogEntry = {
            id: generateLogId(),
            timestamp: new Date().toISOString(),
            type: 'INFO',
            command: 'WEBSOCKET',
            message: `Intentando reconectar... (intento ${reconnectAttemptsRef.current}/${maxReconnectAttempts})`
          };

          setState(prev => ({
            ...prev,
            logs: [...prev.logs.slice(-maxLogs + 1), reconnectLog]
          }));

          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, reconnectInterval);
        } else if (reconnectAttemptsRef.current >= maxReconnectAttempts) {
          setState(prev => ({
            ...prev,
            connectionStatus: 'error',
            error: 'Máximo número de intentos de reconexión alcanzado'
          }));

          const maxAttemptsLog: LogEntry = {
            id: generateLogId(),
            timestamp: new Date().toISOString(),
            type: 'ERROR',
            command: 'WEBSOCKET',
            message: `Máximo número de intentos de reconexión alcanzado (${maxReconnectAttempts})`
          };

          setState(prev => ({
            ...prev,
            logs: [...prev.logs.slice(-maxLogs + 1), maxAttemptsLog]
          }));
        }
      };

      wsRef.current.onerror = (error) => {
        console.error('WebSocket error:', error);
        
        setState(prev => ({
          ...prev,
          connectionStatus: 'error',
          error: 'Error de conexión WebSocket'
        }));

        const errorLog: LogEntry = {
          id: generateLogId(),
          timestamp: new Date().toISOString(),
          type: 'ERROR',
          command: 'WEBSOCKET',
          message: 'Error en la conexión WebSocket'
        };

        setState(prev => ({
          ...prev,
          logs: [...prev.logs.slice(-maxLogs + 1), errorLog]
        }));
      };

    } catch (error) {
      console.error('Error creating WebSocket:', error);
      setState(prev => ({
        ...prev,
        connectionStatus: 'error',
        error: 'Error al crear conexión WebSocket'
      }));
    }
  }, [url, maxLogs, reconnectInterval, maxReconnectAttempts, generateLogId, clearReconnectTimeout]);

  // Desconectar WebSocket
  const disconnect = useCallback(() => {
    isManuallyClosedRef.current = true;
    clearReconnectTimeout();
    
    if (wsRef.current) {
      wsRef.current.close(1000, 'Manual disconnect');
      wsRef.current = null;
    }

    setState(prev => ({
      ...prev,
      isConnected: false,
      connectionStatus: 'disconnected'
    }));
  }, [clearReconnectTimeout]);

  // Reconectar manualmente
  const reconnect = useCallback(() => {
    isManuallyClosedRef.current = false;
    reconnectAttemptsRef.current = 0;
    disconnect();
    setTimeout(connect, 1000);
  }, [connect, disconnect]);

  // Enviar mensaje
  const sendMessage = useCallback((message: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      try {
        const messageStr = typeof message === 'string' ? message : JSON.stringify(message);
        wsRef.current.send(messageStr);
        return true;
      } catch (error) {
        console.error('Error sending message:', error);
        return false;
      }
    }
    return false;
  }, []);

  // Limpiar logs
  const clearLogs = useCallback(() => {
    setState(prev => ({ ...prev, logs: [] }));
  }, []);

  // Agregar log manualmente
  const addLog = useCallback((type: LogEntry['type'], command: string, message: string, data?: any) => {
    const logEntry: LogEntry = {
      id: generateLogId(),
      timestamp: new Date().toISOString(),
      type,
      command,
      message,
      data
    };

    setState(prev => ({
      ...prev,
      logs: [...prev.logs.slice(-maxLogs + 1), logEntry]
    }));
  }, [generateLogId, maxLogs]);

  // Efecto para auto-conectar
  useEffect(() => {
    if (autoConnect) {
      isManuallyClosedRef.current = false;
      connect();
    }

    return () => {
      isManuallyClosedRef.current = true;
      clearReconnectTimeout();
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [autoConnect, connect, clearReconnectTimeout]);

  // Efecto para cleanup
  useEffect(() => {
    return () => {
      clearReconnectTimeout();
    };
  }, [clearReconnectTimeout]);

  return {
    isConnected: state.isConnected,
    connectionStatus: state.connectionStatus,
    logs: state.logs,
    lastMessage: state.lastMessage,
    error: state.error,
    connect,
    disconnect,
    reconnect,
    sendMessage,
    clearLogs,
    addLog
  };
};