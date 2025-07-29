import axios from 'axios';

// Configuración base de axios
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const apiClient = axios.create({
  baseURL: `${API_BASE_URL}/api`,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor para manejo de respuestas
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('Error en API:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

// Tipos TypeScript
export interface ApiResponse<T = any> {
  message: string;
  data?: T;
  status: string;
}

export interface HealthResponse {
  status: string;
  timestamp: string;
  version: string;
}

export interface FileSystemInfo {
  name: string;
  type: string;
  size: number;
  mountPoint: string;
}

export interface PartitionRequest {
  name: string;
  size: number;
  type: string;
  path: string;
}

export interface CommandRequest {
  command: string;
}

// Servicios API
class ApiService {
  // Verificar estado del servidor
  async checkHealth(): Promise<ApiResponse<HealthResponse>> {
    const response = await apiClient.get<ApiResponse<HealthResponse>>('/health');
    return response.data;
  }

  // Obtener sistemas de archivos
  async getFileSystems(): Promise<ApiResponse<FileSystemInfo[]>> {
    const response = await apiClient.get<ApiResponse<FileSystemInfo[]>>('/filesystems');
    return response.data;
  }

  // Crear partición
  async createPartition(partitionData: PartitionRequest): Promise<ApiResponse> {
    const response = await apiClient.post<ApiResponse>('/partition', partitionData);
    return response.data;
  }

  // Ejecutar comando
  async executeCommand(command: string): Promise<ApiResponse> {
    const response = await apiClient.post<ApiResponse>('/execute', { command });
    return response.data;
  }

  // Método genérico para testing de conectividad
  async testConnection(): Promise<boolean> {
    try {
      await this.checkHealth();
      return true;
    } catch (error) {
      console.error('Error de conexión:', error);
      return false;
    }
  }
}

// Instancia singleton del servicio
export const apiService = new ApiService();

// Hook personalizado para React
export { apiClient };
export default apiService;
