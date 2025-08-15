import { useState, useCallback } from 'react';
import { NotificationProps } from '../components/ErrorNotification';

export interface NotificationOptions {
  type?: 'success' | 'warning' | 'error' | 'info';
  title?: string;
  duration?: number;
  persistent?: boolean;
  actions?: Array<{
    label: string;
    action: () => void;
    style?: 'primary' | 'secondary';
  }>;
}

interface MIASuggestion {
  label: string;
  command?: string;
}

// Función para generar sugerencias específicas de MIA
const generateMIASuggestions = (command: string, error: string, baseSuggestions?: string[]): MIASuggestion[] => {
  const suggestions: MIASuggestion[] = [];
  const errorLower = error.toLowerCase();
  const commandLower = command.toLowerCase();

  // Errores de FDISK
  if (commandLower.includes('fdisk')) {
    if (errorLower.includes('no hay espacio suficiente')) {
      suggestions.push(
        { label: 'Verificar espacio del disco', command: 'fdisk -path="/ruta/disco.mia"' },
        { label: 'Crear disco más grande', command: 'mkdisk -size=5000 -unit=M -path="/ruta/disco_grande.mia"' },
        { label: 'Eliminar particiones no usadas' }
      );
    }
    if (errorLower.includes('partición extendida')) {
      suggestions.push(
        { label: 'Crear partición primaria', command: 'fdisk -size=300 -unit=M -path="/ruta/disco.mia" -name="Primaria1"' },
        { label: 'Reducir tamaño de partición', command: 'fdisk -size=100 -unit=M -path="/ruta/disco.mia" -name="Pequeña1"' }
      );
    }
  }

  // Errores de MOUNT
  if (commandLower.includes('mount')) {
    if (errorLower.includes('no se encontró una partición')) {
      const partitionName = extractPartitionName(error);
      suggestions.push(
        { label: 'Listar particiones disponibles', command: 'fdisk -path="/ruta/disco.mia"' },
        { label: 'Verificar nombre exacto de partición' },
        { label: 'Crear la partición faltante', command: `fdisk -size=300 -unit=M -path="/ruta/disco.mia" -name="${partitionName}"` }
      );
    }
    if (errorLower.includes('ya está montada')) {
      suggestions.push(
        { label: 'Ver particiones montadas', command: 'mounted' },
        { label: 'Desmontar primero', command: 'umount -id="A1"' },
        { label: 'Usar diferente ID de montaje' }
      );
    }
    if (errorLower.includes('solo se pueden montar particiones primarias')) {
      suggestions.push(
        { label: 'Crear partición primaria', command: 'fdisk -size=300 -unit=M -path="/ruta/disco.mia" -name="Primaria1"' },
        { label: 'Verificar tipo de partición', command: 'fdisk -path="/ruta/disco.mia"' }
      );
    }
  }

  // Errores de MKDISK
  if (commandLower.includes('mkdisk')) {
    if (errorLower.includes('ya existe') || errorLower.includes('file exists')) {
      suggestions.push(
        { label: 'Usar diferente nombre', command: 'mkdisk -size=1000 -unit=M -path="/ruta/nuevo_disco.mia"' },
        { label: 'Eliminar disco existente', command: 'rmdisk -path="/ruta/disco.mia"' }
      );
    }
    if (errorLower.includes('espacio insuficiente') || errorLower.includes('no space')) {
      suggestions.push(
        { label: 'Crear disco más pequeño', command: 'mkdisk -size=500 -unit=M -path="/ruta/disco.mia"' },
        { label: 'Usar diferente ubicación', command: 'mkdisk -size=1000 -unit=M -path="/tmp/disco.mia"' }
      );
    }
  }

  // Errores generales de ruta
  if (errorLower.includes('no such file') || errorLower.includes('no existe')) {
    suggestions.push(
      { label: 'Crear directorios faltantes', command: 'mkdir -p /ruta/completa' },
      { label: 'Verificar ruta absoluta' },
      { label: 'Usar ruta relativa', command: './Discos/disco.mia' }
    );
  }

  // Errores de permisos
  if (errorLower.includes('permission denied') || errorLower.includes('permisos')) {
    suggestions.push(
      { label: 'Verificar permisos del directorio' },
      { label: 'Usar ruta con permisos de escritura', command: '/tmp/disco.mia' }
    );
  }

  // Agregar sugerencias base si se proporcionaron
  if (baseSuggestions) {
    baseSuggestions.forEach(suggestion => {
      suggestions.push({ label: suggestion });
    });
  }

  return suggestions.slice(0, 4); // Limitar a 4 sugerencias máximo
};

// Función para extraer nombre de partición del mensaje de error
const extractPartitionName = (error: string): string => {
  const match = error.match(/'([^']+)'/);
  return match ? match[1] : 'NuevaParticion';
};

// Función para obtener nombre display del comando
const getCommandDisplayName = (command: string): string => {
  const commandMap: { [key: string]: string } = {
    'mkdisk': 'Crear Disco',
    'rmdisk': 'Eliminar Disco', 
    'fdisk': 'Gestionar Particiones',
    'mount': 'Montar Partición',
    'umount': 'Desmontar Partición',
    'mkfs': 'Formatear Partición',
    'login': 'Iniciar Sesión',
    'logout': 'Cerrar Sesión',
    'mkgrp': 'Crear Grupo',
    'rmgrp': 'Eliminar Grupo',
    'mkusr': 'Crear Usuario',
    'rmusr': 'Eliminar Usuario',
    'chmod': 'Cambiar Permisos',
    'mkfile': 'Crear Archivo',
    'cat': 'Ver Archivo',
    'remove': 'Eliminar Archivo',
    'edit': 'Editar Archivo',
    'rename': 'Renombrar',
    'mkdir': 'Crear Directorio',
    'rmdir': 'Eliminar Directorio',
    'copy': 'Copiar',
    'move': 'Mover',
    'find': 'Buscar',
    'chown': 'Cambiar Propietario',
    'chgrp': 'Cambiar Grupo',
    'recovery': 'Recuperar',
    'loss': 'Simular Pérdida',
    'rep': 'Generar Reporte'
  };
  
  const baseCommand = command.split(' ')[0].toLowerCase();
  return commandMap[baseCommand] || command.toUpperCase();
};

// Función para extraer comando del mensaje de log
const extractCommandFromLog = (logMessage: string): string => {
  // Buscar patrón [COMMAND] en el mensaje
  const match = logMessage.match(/\[(\w+)\]|\b(\w+):/);
  if (match) {
    return match[1] || match[2];
  }
  
  // Buscar comandos conocidos en el mensaje
  const knownCommands = ['mkdisk', 'rmdisk', 'fdisk', 'mount', 'umount', 'mkfs', 'login', 'logout'];
  for (const cmd of knownCommands) {
    if (logMessage.toLowerCase().includes(cmd)) {
      return cmd;
    }
  }
  
  return 'UNKNOWN';
};

// Función para limpiar mensaje de log
const cleanLogMessage = (logMessage: string): string => {
  // Remover timestamp y prefijos de log
  return logMessage
    .replace(/^\[\d+\]\s*/, '') // Remover timestamp
    .replace(/^\[ERROR\]\s*/, '') // Remover tipo de log
    .replace(/^ERROR:\s*\w+\s*-\s*/, '') // Remover prefijo ERROR: COMMAND -
    .replace(/^\w+:\s*/, '') // Remover prefijo COMMAND:
    .trim();
};

// Función para determinar si se debe mostrar notificación para un log
const shouldShowNotificationForLog = (log: any): boolean => {
  // No mostrar notificaciones para logs del sistema o WebSocket
  if (log.command === 'WEBSOCKET' || log.command === 'SYSTEM') {
    return false;
  }
  
  // Solo mostrar para errores de comandos específicos
  const commandsToNotify = ['FDISK', 'MOUNT', 'MKDISK', 'MKFS', 'LOGIN', 'LOGOUT'];
  return commandsToNotify.includes(log.command?.toUpperCase());
};

export const useNotifications = () => {
  const [notifications, setNotifications] = useState<NotificationProps[]>([]);

  const generateId = useCallback(() => {
    return `notification_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }, []);

  const addNotification = useCallback((
    message: string, 
    options: NotificationOptions = {}
  ) => {
    const {
      type = 'info',
      title = 'Notificación',
      duration = 5000,
      persistent = false,
      actions = []
    } = options;

    const id = generateId();
    
    const notification: NotificationProps = {
      id,
      type,
      title,
      message,
      duration,
      persistent,
      actions
    };

    setNotifications(prev => [...prev, notification]);
    return id;
  }, [generateId]);

  const removeNotification = useCallback((id: string) => {
    setNotifications(prev => prev.filter(notification => notification.id !== id));
  }, []);

  const clearAll = useCallback(() => {
    setNotifications([]);
  }, []);

  // Métodos de conveniencia
  const showSuccess = useCallback((message: string, options?: Omit<NotificationOptions, 'type'>) => {
    return addNotification(message, { ...options, type: 'success', title: options?.title || '✅ Éxito' });
  }, [addNotification]);

  const showError = useCallback((message: string, options?: Omit<NotificationOptions, 'type'>) => {
    return addNotification(message, { 
      ...options, 
      type: 'error', 
      title: options?.title || '❌ Error',
      duration: options?.duration || 8000 // Errores duran más tiempo
    });
  }, [addNotification]);

  const showWarning = useCallback((message: string, options?: Omit<NotificationOptions, 'type'>) => {
    return addNotification(message, { 
      ...options, 
      type: 'warning', 
      title: options?.title || '⚠️ Advertencia',
      duration: options?.duration || 6000
    });
  }, [addNotification]);

  const showInfo = useCallback((message: string, options?: Omit<NotificationOptions, 'type'>) => {
    return addNotification(message, { ...options, type: 'info', title: options?.title || 'ℹ️ Información' });
  }, [addNotification]);

  // Método para errores de comando específicos de MIA
  const showCommandError = useCallback((
    command: string, 
    error: string, 
    suggestions?: string[]
  ) => {
    // Generar sugerencias específicas basadas en el tipo de error
    const enhancedSuggestions = generateMIASuggestions(command, error, suggestions);
    
    const actions = enhancedSuggestions.map((suggestion: MIASuggestion) => ({
      label: suggestion.label,
      action: () => {
        if (suggestion.command) {
          // Si la sugerencia incluye un comando, copiarlo al portapapeles
          navigator.clipboard.writeText(suggestion.command);
          showInfo(`Comando copiado: ${suggestion.command}`, {
            title: '📋 Comando Copiado',
            duration: 3000
          });
        }
        console.log('Aplicando sugerencia:', suggestion.label);
      },
      style: 'secondary' as const
    }));

    return showError(error, {
      title: `⚠️ Error: ${getCommandDisplayName(command)}`,
      persistent: true,
      duration: 10000,
      actions: [
        ...actions,
        {
          label: 'Ver ayuda del comando',
          action: () => {
            console.log('Mostrando ayuda para:', command);
          },
          style: 'primary' as const
        }
      ]
    });
  }, [showError, showInfo]);

  // Método para errores de validación
  const showValidationError = useCallback((
    field: string, 
    errors: string[]
  ) => {
    const errorMessage = errors.join(', ');
    return showError(`Errores en ${field}: ${errorMessage}`, {
      title: '📝 Error de Validación',
      duration: 6000
    });
  }, [showError]);

  // Método para errores de conexión
  const showConnectionError = useCallback((
    retryAction?: () => void
  ) => {
    const actions = retryAction ? [{
      label: 'Reintentar',
      action: retryAction,
      style: 'primary' as const
    }] : [];

    return showError('No se puede conectar al servidor. Verifique su conexión.', {
      title: '🌐 Error de Conexión',
      persistent: true,
      actions
    });
  }, [showError]);

  // Método para mostrar errores basados en logs del backend
  const showBackendLogError = useCallback((logMessage: string, command?: string) => {
    // Extraer comando del log si no se proporciona
    const extractedCommand = command || extractCommandFromLog(logMessage);
    const cleanError = cleanLogMessage(logMessage);
    
    return showCommandError(extractedCommand, cleanError);
  }, [showCommandError]);

  // Método para procesar logs en tiempo real y mostrar notificaciones
  const processLogForNotification = useCallback((log: any) => {
    if (log.type === 'ERROR') {
      // Solo mostrar notificación para errores críticos, no todos los logs
      if (shouldShowNotificationForLog(log)) {
        showBackendLogError(log.message, log.command);
      }
    } else if (log.type === 'SUCCESS' && log.command !== 'WEBSOCKET') {
      // Mostrar éxitos para comandos ejecutados
      showSuccess(log.message, {
        title: `✅ ${getCommandDisplayName(log.command)}`,
        duration: 4000
      });
    }
  }, [showBackendLogError, showSuccess]);

  return {
    notifications,
    addNotification,
    removeNotification,
    clearAll,
    showSuccess,
    showError,
    showWarning,
    showInfo,
    showCommandError,
    showValidationError,
    showConnectionError,
    showBackendLogError,
    processLogForNotification
  };
};