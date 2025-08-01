# 📡 API y Endpoints - Documentación

## 🌐 **Servidor Backend**

**URL Base:** `http://localhost:8080`
**Tecnología:** Go (Golang) con Gorilla Mux
**CORS:** Habilitado para `localhost:3000`

---

## 📋 **Endpoints Disponibles**

### ✅ **1. Health Check**
```http
GET /api/health
```

**Descripción:** Verificar estado del servidor
**Respuesta:**
```json
{
  "message": "Backend funcionando correctamente",
  "data": {
    "status": "OK",
    "timestamp": "2025-08-01T10:30:00Z",
    "version": "1.0.0"
  },
  "status": "success"
}
```

---

### ✅ **2. Listar Sistemas de Archivos**
```http
GET /api/filesystems
```

**Descripción:** Obtener lista de sistemas de archivos disponibles
**Respuesta:**
```json
{
  "message": "Sistemas de archivos obtenidos",
  "data": [
    {
      "name": "disk1.dk",
      "type": "EXT2",
      "size": 1024000,
      "mountPoint": "/mnt/disk1"
    }
  ],
  "status": "success"
}
```

---

### ✅ **3. Crear Partición**
```http
POST /api/partition
```

**Descripción:** Crear una nueva partición
**Body (JSON):**
```json
{
  "disk": "/ruta/disco.mia",
  "size": 1024,
  "unit": "M",
  "fit": "FF",
  "type": "P",
  "name": "Particion1"
}
```

**Respuesta:**
```json
{
  "message": "Partición creada exitosamente",
  "data": { ...requestData },
  "status": "success"
}
```

---

### ✅ **4. Ejecutar Comando**
```http
POST /api/execute
```

**Descripción:** Ejecutar comandos del sistema de archivos
**Body (JSON):**
```json
{
  "command": "mkdisk -size 10 -unit M -fit FF -path /tmp/disco.mia"
}
```

**Respuesta:**
```json
{
  "message": "Comando ejecutado",
  "data": {
    "command": "mkdisk -size 10 -unit M -fit FF -path /tmp/disco.mia",
    "result": "Comando ejecutado exitosamente"
  },
  "status": "success"
}
```

---

### ✅ **5. Obtener Logs (HTTP Polling)**
```http
GET /api/logs
```

**Descripción:** Obtener logs del sistema via HTTP
**Respuesta:**
```json
{
  "logs": "contenido de los logs...",
  "status": "success"
}
```

---

### ✅ **6. Stream de Logs (Server-Sent Events)**
```http
GET /api/logs/stream
```

**Descripción:** Stream en tiempo real de logs via SSE
**Content-Type:** `text/event-stream`
**Ejemplo de evento:**
```
data: {"type":"INFO","command":"MkDisk","message":"Disco creado","time":"1693574400"}

```

---

### ✅ **7. WebSocket para Logs**
```http
GET /ws
```

**Descripción:** Conexión WebSocket para logs en tiempo real
**Protocolo:** WebSocket
**Ejemplo de mensaje:**
```json
{
  "type": "INFO",
  "command": "MkDisk", 
  "message": "Disco creado con éxito en /tmp/disco.mia",
  "time": "1693574400"
}
```

---

## 🔧 **Niveles de Log**

| Nivel | Descripción | Color Sugerido |
|-------|-------------|----------------|
| `INFO` | Información general | 🔵 Azul |
| `WARNING` | Advertencias | 🟡 Amarillo |
| `ERROR` | Errores | 🔴 Rojo |
| `SUCCESS` | Operaciones exitosas | 🟢 Verde |

---

## 🌐 **Ejemplos de Uso con Frontend**

### **JavaScript/TypeScript:**

#### **1. Ejecutar Comando**
```typescript
const response = await fetch('http://localhost:8080/api/execute', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    command: 'mkdisk -size 10 -unit M -path /tmp/test.mia'
  })
});

const result = await response.json();
console.log(result);
```

#### **2. WebSocket para Logs**
```typescript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = (event) => {
  const logMessage = JSON.parse(event.data);
  console.log(`[${logMessage.type}] ${logMessage.command}: ${logMessage.message}`);
};

ws.onopen = () => {
  console.log('Conectado al WebSocket');
};

ws.onclose = () => {
  console.log('Desconectado del WebSocket');
};
```

#### **3. Server-Sent Events**
```typescript
const eventSource = new EventSource('http://localhost:8080/api/logs/stream');

eventSource.onmessage = (event) => {
  const logMessage = JSON.parse(event.data);
  console.log(`[${logMessage.type}] ${logMessage.message}`);
};

eventSource.onerror = (error) => {
  console.error('Error en SSE:', error);
};
```

---

## 🛡️ **Manejo de Errores**

### **Códigos de Estado HTTP:**
- `200` - Éxito
- `400` - Error en request (JSON inválido, parámetros faltantes)
- `405` - Método no permitido
- `500` - Error interno del servidor

### **Formato de Error:**
```json
{
  "message": "Descripción del error",
  "status": "error"
}
```

---

## 🔧 **Configuración CORS**

```go
c := cors.New(cors.Options{
    AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"*"},
    AllowCredentials: true,
})
```

---

## 🚀 **Testing de API**

### **Con cURL:**
```bash
# Health check
curl http://localhost:8080/api/health

# Ejecutar comando
curl -X POST http://localhost:8080/api/execute \
  -H "Content-Type: application/json" \
  -d '{"command": "mkdisk -size 10 -path /tmp/test.mia"}'

# Obtener logs
curl http://localhost:8080/api/logs
```

### **Con Postman:**
1. Importar collection con los endpoints
2. Configurar `Base URL` como `http://localhost:8080`
3. Probar cada endpoint

### **Con Frontend React:**
Ya integrado en `src/services/apiService.ts`

---

## 📊 **Métricas y Monitoreo**

### **Logs del Servidor:**
- Todos los requests se registran en consola
- Errores se guardan en `backend/logs/error.log`
- Logs en tiempo real via WebSocket/SSE

### **Health Monitoring:**
- Endpoint `/api/health` para verificar estado
- Timestamp en cada respuesta
- Información de versión incluida
