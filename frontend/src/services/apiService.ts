import axios from 'axios';

// Configuración base de axios
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const apiClient = axios.create({
  baseURL: `${API_BASE_URL}/api`,
  timeout: 30000, // Aumentado para scripts largos
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor para manejo de respuestas
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('Error en API:', error.response?.data || error.message);
    
    // Manejar diferentes tipos de errores
    if (error.code === 'ECONNREFUSED') {
      throw new Error('No se puede conectar al servidor. Verifique que el backend esté ejecutándose.');
    }
    
    if (error.response?.status === 400) {
      throw new Error(error.response.data?.message || 'Comando inválido');
    }
    
    if (error.response?.status === 500) {
      throw new Error(error.response.data?.message || 'Error interno del servidor');
    }
    
    if (error.response?.status === 404) {
      throw new Error('Endpoint no encontrado');
    }
    
    return Promise.reject(error);
  }
);

// Tipos TypeScript
export interface ApiResponse<T = any> {
  message: string;
  data?: T;
  status: string;
  timestamp?: string;
}

export interface HealthResponse {
  status: string;
  timestamp: string;
  version: string;
  uptime?: string;
}

export interface FileSystemInfo {
  name: string;
  type: string;
  size: number;
  mountPoint: string;
  status?: 'mounted' | 'unmounted';
  path?: string; // Nueva propiedad para la ruta completa
}

export interface CommandRequest {
  command: string;
  script?: string;
}

// Servicios API
class ApiService {
  // Verificar estado del servidor
  async checkHealth(): Promise<ApiResponse<HealthResponse>> {
    const response = await apiClient.get<ApiResponse<HealthResponse>>('/health');
    return response.data;
  }

  // Obtener sistemas de archivos (FUNCIONAL)
  async getFileSystems(searchPath?: string): Promise<ApiResponse<FileSystemInfo[]>> {
    const params = searchPath ? { path: searchPath } : {};
    const response = await apiClient.get<ApiResponse<FileSystemInfo[]>>('/filesystems', { params });
    return response.data;
  }

  // Obtener comandos soportados (FUNCIONAL)
  async getSupportedCommands(): Promise<ApiResponse<{commands: string[], total: number}>> {
    const response = await apiClient.get<ApiResponse<{commands: string[], total: number}>>('/commands');
    return response.data;
  }

  // Ejecutar comando individual (FUNCIONAL)
  async executeCommand(command: string): Promise<ApiResponse> {
    const response = await apiClient.post<ApiResponse>('/execute', { 
      command: command.trim()
    });
    return response.data;
  }

  // Validar comando sin ejecutar (FUNCIONAL)
  async validateCommand(command: string): Promise<ApiResponse<{valid_commands: string[]}>> {
    const response = await apiClient.post<ApiResponse<{valid_commands: string[]}>>('/validate', { 
      command: command.trim() 
    });
    return response.data;
  }

  // Ejecutar script completo (FUNCIONAL)
  async executeScript(script: string): Promise<ApiResponse> {
    const response = await apiClient.post<ApiResponse>('/execute', {
      script: script.trim()
    });
    return response.data;
  }

  // Método genérico para testing de conectividad
  async testConnection(): Promise<boolean> {
    try {
      const response = await this.checkHealth();
      return response.status === 'success';
    } catch (error) {
      console.error('Error de conexión:', error);
      return false;
    }
  }
}

// Instancia singleton del servicio
export const apiService = new ApiService();

// Funciones de utilidad
export const formatFileSize = (bytes: number): string => {
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  if (bytes === 0) return '0 Bytes';
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i];
};

export const isValidPath = (path: string): boolean => {
  // Validación mejorada de rutas - permitir rutas relativas y absolutas
  if (!path || path.length === 0) return false;
  return path.startsWith('/') || path.startsWith('./') || path.startsWith('../');
};

export const sanitizeCommand = (command: string): string => {
  // Sanitizar comando para evitar inyecciones, pero preservar comillas para rutas con espacios
  return command.trim().replace(/[;&|`$()]/g, '');
};

// Función para validar y formatear comandos con espacios
export const validateAndFormatCommand = (command: string): {isValid: boolean, formattedCommand: string, errors: string[]} => {
  const errors: string[] = [];
  let formattedCommand = command.trim();
  
  // Verificar que las comillas estén balanceadas
  const quoteCount = (formattedCommand.match(/"/g) || []).length;
  if (quoteCount % 2 !== 0) {
    errors.push('Las comillas deben estar balanceadas');
    return { isValid: false, formattedCommand, errors };
  }
  
  // SOLUCIÓN: Convertir comillas a formato que el backend actual puede manejar
  // Escapar espacios en lugar de usar comillas
  formattedCommand = formattedCommand.replace(/-(\w+)="([^"]*\s[^"]*)"/g, (match, param, path) => {
    // Convertir espacios a %20 temporalmente para que el backend los maneje como un solo parámetro
    const escapedPath = path.replace(/\s/g, '%20');
    errors.push(`ℹ️ INFO: Ruta con espacios convertida para compatibilidad`);
    return `-${param}=${escapedPath}`;
  });
  
  // Remover comillas restantes
  formattedCommand = formattedCommand.replace(/"/g, '');
  
  return {
    isValid: true,
    formattedCommand,
    errors
  };
};

// Función para normalizar rutas de disco
export const normalizeDiskPath = (relativePath: string): string => {
  if (relativePath.startsWith('./Documents/')) {
    return `/home/julian/Documents/${relativePath.substring(13)}`;
  }
  if (relativePath.startsWith('./')) {
    return `/home/julian/Documents/MIA_2S2025_P1_201905884/backend${relativePath.substring(1)}`;
  }
  return relativePath;
};

// Función para formatear información de disco creado
export const formatDiskCreationResult = (result: any): string => {
  if (!result.success) {
    return `❌ Error: ${result.message}`;
  }
  
  const data = result.data;
  const fullPath = normalizeDiskPath(data.path);
  
  return `✅ Disco creado exitosamente:
📁 Ruta: ${fullPath}
📏 Tamaño: ${data.size} ${data.unit}B
⚙️ Ajuste: ${data.fit}
🔍 Verifica con: ls -la "${fullPath}"`;
};

// Función para generar comandos de ejemplo seguros
export const getExampleCommands = (): string[] => {
  return [
    'mkdisk -size=3000 -path=./Discos/disco1.mia',
    'mkdisk -size=5 -unit=M -path="./Discos/mis discos/disco2.mia"',
    'mkdisk -size=10 -path="/tmp/disco con espacios.mia"',
    'mkdisk -size=1000 -path=./Discos/mi_disco.mia -fit=BF',
    'fdisk -path=./Discos/disco1.mia',
    'mount -path="./Discos/mis discos/disco2.mia" -name=disco2',
    '# Para rutas sin espacios:',
    'mkdisk -size=5 -unit=M -path=./Discos/disco_sin_espacios.mia'
  ];
};

// Función para diagnosticar problemas de creación de archivos
export const getDiagnosticInfo = (): string => {
  return `🔧 Diagnóstico de problemas comunes:

1. ✅ SOLUCIONADO: Backend ahora soporta espacios en rutas con comillas
   
2. Usar comillas para rutas con espacios:
   ✅ -path="./Discos/mis discos/disco.mia"
   ✅ -path="/tmp/disco con espacios.mia"

3. Crear directorios si no existen:
   mkdir -p /home/julian/Documents/MIA_2S2025_P1_201905884/backend/Discos

4. Comandos que funcionan:
   mkdisk -size=5 -unit=M -path="./Discos/mis discos/disco.mia"
   mkdisk -size=10 -path="/tmp/disco temporal.mia"

5. Para debugging, verificar logs del backend en:
   /home/julian/Documents/MIA_2S2025_P1_201905884/backend/logs/`;
};

// Función para convertir nombres con espacios a nombres seguros
export const sanitizePathName = (path: string): string => {
  return path
    .replace(/\s+/g, '_')           // Espacios por guiones bajos
    .replace(/[^a-zA-Z0-9_.//-]/g, '') // Quitar caracteres especiales
    .replace(/_{2,}/g, '_');        // Múltiples guiones bajos por uno solo
};

// Función para generar sugerencias de nombres seguros
export const getSafePathSuggestion = (originalPath: string): string => {
  const safePath = sanitizePathName(originalPath);
  if (safePath !== originalPath) {
    return `Sugerencia de ruta segura: ${safePath}`;
  }
  return '';
};

// Hook personalizado para React
export { apiClient };
export default apiService;