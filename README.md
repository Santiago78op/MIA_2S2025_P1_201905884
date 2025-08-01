# Proyecto MIA - Sistema de Archivos EXT2

## 📋 **Resumen del Proyecto**

Simulador de **sistema de archivos EXT2** con:
- **Backend en Go (Golang)** - obligatorio ✅ **IMPLEMENTADO**
- **Frontend en React + TypeScript** ✅ **IMPLEMENTADO**
- **WebSockets y SSE para logs en tiempo real** ✅ **IMPLEMENTADO**
- **Sistema de comandos completo** 🔄 **EN DESARROLLO**
- **Reportes con Graphviz** 🔄 **PENDIENTE**
- **Documentación técnica completa** ✅ **ACTUALIZADA**

**Fecha de entrega:** 7 de septiembre de 2025, 23:59 horas

## 🚀 **Instalación y Ejecución**

### En Windows (desarrollo)
```bash
# Backend
cd backend
go mod tidy
go run main.go

# Frontend (terminal separado)
cd frontend
npm install
npm start
```

### En Pop!_OS/Linux (producción)
```bash
# Clonar repositorio
git clone https://github.com/Santiago78op/MIA_2S2025_P1_201905884.git
cd MIA_2S2025_P1_201905884

# Seguir instrucciones en README_POPOS.md
```

### 🌐 **URLs de Acceso**
- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080
- **WebSocket:** ws://localhost:8080/ws
- **Health Check:** http://localhost:8080/api/health

---

## 🎯 **Componentes Obligatorios Mínimos**

### ✅ **Requisitos para Calificación**
1. ✅ Aplicación Web funcional
2. ✅ Creación de particiones con ajustes y Mount
3. ✅ Ejecución completa del script
4. ✅ Reportes para validar funcionamiento
5. ✅ Documentación completa

### 🔧 **Comandos Implementados**

#### ✅ **Gestión de discos:**
- `mkdisk` - ✅ **IMPLEMENTADO** - Crear discos virtuales (.mia)
- `rmdisk` - 🔄 **PENDIENTE** - Eliminar discos
- `fdisk` - 🔄 **PENDIENTE** - Gestión de particiones  
- `mount` - 🔄 **PENDIENTE** - Montar particiones

#### 🔄 **Sistema de archivos:**
- `mkfs` - 🔄 **PENDIENTE** - Formatear particiones con EXT2
- `login` - 🔄 **PENDIENTE** - Iniciar sesión
- `logout` - 🔄 **PENDIENTE** - Cerrar sesión

#### 🔄 **Usuarios y grupos:**
- `mkgrp` - 🔄 **PENDIENTE** - Crear grupos
- `rmgrp` - 🔄 **PENDIENTE** - Eliminar grupos
- `mkusr` - 🔄 **PENDIENTE** - Crear usuarios
- `rmusr` - 🔄 **PENDIENTE** - Eliminar usuarios

#### 🔄 **Archivos y carpetas:**
- `mkfile` - 🔄 **PENDIENTE** - Crear archivos
- `mkdir` - 🔄 **PENDIENTE** - Crear directorios
- `cat` - 🔄 **PENDIENTE** - Mostrar contenido

### 📡 **API Endpoints Disponibles**

| Método | Endpoint | Descripción | Estado |
|--------|----------|-------------|--------|
| GET | `/api/health` | Estado del servidor | ✅ |
| GET | `/api/filesystems` | Listar sistemas de archivos | ✅ |
| POST | `/api/partition` | Crear partición | ✅ |
| POST | `/api/execute` | Ejecutar comandos | ✅ |
| GET | `/api/logs` | Obtener logs (polling) | ✅ |
| GET | `/api/logs/stream` | Stream de logs (SSE) | ✅ |
| GET | `/ws` | WebSocket para logs en tiempo real | ✅ |

---

## 🗂️ **Arquitectura del Proyecto (Actualizada)**

### **Estructura Actual**
```
MIA_2S2025_P1_201905884/
├── backend/
│   ├── main.go                     # ✅ Servidor principal con WebSockets/SSE
│   ├── go.mod                      # ✅ Dependencias Go
│   ├── go.sum                      # ✅ Lock de dependencias
│   ├── command/                    # ✅ Comandos del sistema
│   │   └── disk/
│   │       └── mkdisk.go          # ✅ Implementado
│   ├── struct/                     # ✅ Estructuras de datos
│   │   ├── strMBR.go              # ✅ Master Boot Record
│   │   └── strPartition.go        # ✅ Estructuras de particiones
│   ├── Utils/                      # ✅ Utilidades
│   │   └── logger.go              # ✅ Logger con WebSocket/SSE
│   ├── action/                     # 🔄 Acciones auxiliares
│   └── logs/                       # ✅ Directorio de logs
├── frontend/                       # ✅ React + TypeScript
│   ├── src/
│   │   ├── components/             # ✅ Componentes React
│   │   │   ├── CommandExecutor.tsx # ✅ Ejecutor de comandos
│   │   │   ├── ConnectionStatus.tsx # ✅ Estado de conexión
│   │   │   └── FileSystemList.tsx  # ✅ Lista de sistemas
│   │   ├── services/               # ✅ Servicios API
│   │   │   └── apiService.ts       # ✅ Cliente HTTP
│   │   ├── hooks/                  # ✅ Hooks personalizados
│   │   │   └── useApi.ts           # ✅ Hook para API
│   │   └── App.tsx                 # ✅ Aplicación principal
│   ├── package.json                # ✅ Dependencias Node
│   └── public/                     # ✅ Archivos públicos
├── docs/                           # ✅ Documentación
│   └── ESTRUCTURAS_GUIA.md         # ✅ Guía de estructuras
├── README.md                       # ✅ Este archivo
├── INTEGRATION_STATUS.md           # ✅ Estado de integración
├── docker-compose.yml              # ✅ Configuración Docker
└── run-*.sh                        # ✅ Scripts de ejecución
```

---

## � **Documentación Completa**

### **📋 Guías Principales:**
- 📖 **[README.md](README.md)** - Este archivo (documentación principal)
- 🗺️ **[ROADMAP.md](docs/ROADMAP.md)** - Cronograma y progreso del proyecto
- ✅ **[INTEGRATION_STATUS.md](INTEGRATION_STATUS.md)** - Estado de integración Frontend-Backend

### **🔧 Documentación Técnica:**
- 🛠️ **[COMANDOS_IMPLEMENTADOS.md](docs/COMANDOS_IMPLEMENTADOS.md)** - Comandos del sistema
- 📚 **[ESTRUCTURAS_GUIA.md](docs/ESTRUCTURAS_GUIA.md)** - Estructuras de datos EXT2
- 📡 **[API_ENDPOINTS.md](docs/API_ENDPOINTS.md)** - Documentación de API REST

### **🚀 Guías de Instalación:**
- 🐧 **[README_POPOS.md](README_POPOS.md)** - Instalación en Pop!_OS/Linux
- 🔧 **[GITHUB_GUIDE.md](GITHUB_GUIDE.md)** - Configuración de GitHub

---

## 🧪 **Testing y Desarrollo**

### **🔧 Comandos de Desarrollo:**
```bash
# Compilar backend
cd backend && go build -o backend.exe .

# Ejecutar tests
cd backend && go test ./...

# Ejecutar en modo desarrollo
cd backend && go run main.go

# Frontend en desarrollo
cd frontend && npm run start
```

### **📊 Estado de Testing:**
- ✅ **Backend:** Compilación exitosa
- ✅ **API Endpoints:** Todos funcionales
- ✅ **WebSockets/SSE:** Implementados y probados
- ✅ **Frontend:** Conectado al backend
- 🔄 **Comandos:** mkdisk probado, otros pendientes

---

## 🏆 **Características Destacadas**

### **✨ Features Implementadas:**
- 🔄 **Logs en Tiempo Real:** WebSockets + SSE
- 🌐 **API REST Completa:** 7 endpoints funcionales
- 📱 **Frontend Reactivo:** Componentes modernos en TypeScript
- 🏗️ **Arquitectura Escalable:** Modular y bien documentada
- 📊 **Sistema de Logging:** 4 niveles (INFO, WARNING, ERROR, SUCCESS)

### **🔧 Tecnologías Utilizadas:**
- **Backend:** Go 1.21 + Gorilla Mux + Gorilla WebSocket
- **Frontend:** React 18 + TypeScript + Hooks personalizados
- **Base de Datos:** Archivos binarios (.mia)
- **Comunicación:** REST API + WebSockets + SSE
- **Build:** Go Modules + npm

---

## 📞 **Soporte y Contacto**

### **🐛 Reportar Issues:**
1. Crear issue en GitHub con descripción detallada
2. Incluir logs del backend/frontend
3. Especificar pasos para reproducir

### **💡 Contribuir:**
1. Fork del repositorio
2. Crear branch para feature/bugfix
3. Pull request con descripción clara

### **📧 Contacto:**
- **Repositorio:** [Santiago78op/MIA_2S2025_P1_201905884](https://github.com/Santiago78op/MIA_2S2025_P1_201905884)
- **Documentación:** Carpeta `/docs`

---

**📅 Última actualización:** 1 de agosto de 2025  
**👨‍💻 Desarrollador:** Santiago78op  
**📚 Curso:** MIA - Segundo Semestre 2025  
**🎯 Proyecto:** Sistema de Archivos EXT2
   - EBR (Extended Boot Record)

3. **Comandos básicos de disco**
   - `mkdisk`: Crear disco virtual (.mia)
   - `rmdisk`: Eliminar disco
   - `fdisk`: Gestión de particiones
   - `mount`: Montar particiones

#### Puntos de Validación:
- ✅ Crear discos de diferentes tamaños
- ✅ Crear particiones primarias, extendidas y lógicas
- ✅ Reportes MBR y DISK funcionando

### **Sprint 2: Sistema EXT2 (Semana 2)**
#### Tareas:
1. **Estructuras EXT2**
   - Superbloque
   - Bitmap de inodos
   - Bitmap de bloques
   - Tabla de inodos
   - Tabla de bloques

2. **Comando mkfs**
   - Formatear partición a EXT2
   - Crear archivo users.txt en raíz
   - Inicializar bitmaps

#### Puntos de Validación:
- ✅ Particiones formateadas correctamente
- ✅ Reportes bm_inode, bm_block, inode, block
- ✅ Reporte tree mostrando estructura

### **Sprint 3: Gestión de Usuarios (Semana 3)**
#### Tareas:
1. **Sistema de autenticación**
   - `login`: Iniciar sesión
   - `logout`: Cerrar sesión

2. **Gestión de grupos y usuarios**
   - `mkgrp`: Crear grupo
   - `rmgrp`: Eliminar grupo  
   - `mkusr`: Crear usuario
   - `rmusr`: Eliminar usuario

#### Puntos de Validación:
- ✅ Archivo users.txt actualizado correctamente
- ✅ Validación de sesiones activas
- ✅ Reporte cat del archivo users.txt

### **Sprint 4: Archivos y Carpetas (Semana 4)**
#### Tareas:
1. **Gestión de archivos**
   - `mkfile`: Crear archivos
   - `mkdir`: Crear directorios
   - `cat`: Mostrar contenido

2. **Frontend web**
   - Área de entrada de comandos
   - Área de salida
   - Carga de scripts .smia
   - Botón ejecutar

#### Puntos de Validación:
- ✅ Crear archivos y carpetas
- ✅ Reporte tree completo
- ✅ Frontend funcional

### **Sprint 5: Reportes y Documentación (Semana 5)**
#### Tareas:
1. **Completar todos los reportes**
   - MBR, DISK, INODE, BLOCK
   - BM_INODE, BM_BLOCK
   - TREE, SB

2. **Documentación**
   - Manual técnico
   - Manual de usuario
   - Arquitectura del sistema

---

## 💻 **Implementación Técnica**

### **1. Estructuras Clave**

#### **MBR (Master Boot Record)**
```go
type MBR struct {
    Mbr_size       int32
    Mbr_date_create [19]byte
    Mbr_disk_signature int32
    Disk_fit       byte
    Mbr_partitions [4]Partition
}
```

#### **Superbloque EXT2**
```go
type Superblock struct {
    S_filesystem_type   int32    // 2 = EXT2
    S_inodes_count      int32    // Total de inodos
    S_blocks_count      int32    // Total de bloques
    S_free_blocks_count int32    // Bloques libres
    S_free_inodes_count int32    // Inodos libres
    S_mtime             [19]byte // Fecha montaje
    S_umtime            [19]byte // Fecha desmontaje
    S_bm_inode_start    int32    // Inicio bitmap inodos
    S_bm_block_start    int32    // Inicio bitmap bloques
    S_inode_start       int32    // Inicio tabla inodos
    S_block_start       int32    // Inicio tabla bloques
    // ... más campos
}
```

### **2. Fórmulas Importantes**

#### **Cálculo de inodos y bloques:**
```
n = floor((particion_size - sizeof(superblock)) / (4 + 3*sizeof(block) + sizeof(inode)))

inodos = n
bloques = 3*n
```

#### **Distribución en partición:**
```
| Superbloque | Bitmap Inodos | Bitmap Bloques | Tabla Inodos | Tabla Bloques |
|    92 bytes |     n bytes   |    3n bytes    |  n*124 bytes |   3n*64 bytes |
```

### **3. Sistema de IDs para Mount**
- Formato: `[últimos 2 dígitos carnet][número partición][letra disco]`
- Ejemplo carnet 202401234: `341A`, `341B`, `342A`, etc.

---

## 📊 **Reportes Obligatorios**

### **Para Calificación se requieren:**

1. **Comandos mkdisk y fdisk:** Reportes `mbr` y `disk`
2. **Comando mount:** Comando `mounted`
3. **Comando mkfs:** Reportes `bm_inode`, `bm_block`, `inode`, `block` (o `tree`)
4. **Gestión usuarios:** Comando `cat` o reporte `file` del users.txt
5. **mkfile y mkdir:** Reporte `tree`

### **Implementación con Graphviz:**
```go
// Ejemplo estructura reporte
cad := "digraph G {\n"
cad += "  node [shape=record];\n"
cad += "  MBR [label=\"{MBR|Size: " + fmt.Sprintf("%d", mbr.Mbr_size) + "}\"]\n"
cad += "}\n"

// Generar imagen
exec.Command("dot", "-Tpng", "-o", path).Run()
```

---

## 🌐 **Frontend Requirements**

### **Componentes Obligatorios:**
1. **Área de entrada:** TextArea para comandos
2. **Área de salida:** Mostrar respuestas del servidor
3. **Botón cargar script:** Cargar archivos .smia
4. **Botón ejecutar:** Enviar comandos al backend

### **Ejemplo React:**
```jsx
function App() {
  const [commands, setCommands] = useState('');
  const [output, setOutput] = useState('');

  const executeCommands = async () => {
    const response = await fetch('/api/execute', {
      method: 'POST',
      body: JSON.stringify({commands}),
      headers: {'Content-Type': 'application/json'}
    });
    const result = await response.text();
    setOutput(result);
  };

  return (
    <div>
      <textarea value={commands} onChange={(e) => setCommands(e.target.value)} />
      <button onClick={executeCommands}>Ejecutar</button>
      <textarea value={output} readOnly />
    </div>
  );
}
```

---

## 🔍 **Mejoras y Optimizaciones del Código**

### **1. Manejo de Errores**
```go
// Usar logger consistente
logger := utils.NewLogger("comando")
if err != nil {
    logger.LogError("ERROR: %s", err.Error())
    return logger.GetErrors()
}
```

### **2. Validación de Parámetros**
```go
// Validar parámetros obligatorios
if id == "" {
    logger.LogError("ERROR: Parámetro -id es obligatorio")
    return logger.GetErrors()
}
```

### **3. Gestión de Archivos**
```go
// Siempre cerrar archivos
defer file.Close()

// Usar funciones auxiliares
if err := Acciones.ReadObject(file, &structure, offset); err != nil {
    return err
}
```

### **4. Optimización de Memoria**
- No usar estructuras en memoria para archivos/carpetas
- Leer/escribir directamente al disco
- Usar offsets precisos para navegación

---

## 📚 **Documentación Requerida**

### **Manual Técnico debe incluir:**
1. **Arquitectura del sistema**
   - Diagramas de componentes
   - Comunicación frontend-backend
   
2. **Estructuras de datos**
   - MBR, EBR, Superbloque
   - Inodos, bloques, bitmaps
   
3. **Comandos implementados**
   - Descripción detallada
   - Parámetros y ejemplos
   - Efectos en el sistema

### **Manual de Usuario debe incluir:**
1. **Instalación y configuración**
2. **Capturas de pantalla**
3. **Ejemplos de uso**
4. **Resolución de problemas**

---

## ⚠️ **Puntos Críticos a Evitar**

### **Errores Comunes:**
1. ❌ **Usar otro lenguaje que no sea Go para backend**
2. ❌ **Modificar estructuras definidas (agregar/quitar atributos)**
3. ❌ **Usar listas/árboles en memoria para archivos**
4. ❌ **Permitir que el archivo .mia crezca**
5. ❌ **No implementar todos los reportes obligatorios**
6. ❌ **Copiar código (detección automática = 0 puntos)**

### **Mejores Prácticas:**
1. ✅ **Commits frecuentes en GitHub**
2. ✅ **Dar acceso temprano a auxiliares**
3. ✅ **Probar cada comando antes de continuar**
4. ✅ **Validar reportes paso a paso**
5. ✅ **Documentar mientras desarrollas**

---

## 🎯 **Distribución de Puntos**

| Área | Puntos | Descripción |
|------|--------|-------------|
| **Funcionalidad** | 80 | Comandos + reportes funcionando |
| **Procedimiento** | 10 | Código bien estructurado |
| **Preguntas** | 5 | Conocimiento del proyecto |
| **Documentación** | 5 | Manuales completos |
| **TOTAL** | **100** | |

### **Desglose Funcionalidad (80 pts):**
- Aplicación web: 5 pts
- Manejo discos (mkdisk, rmdisk, fdisk, mount): 32 pts
- Sistema EXT2 (mkfs): 18 pts  
- Usuarios (login, logout, mkgrp, etc.): 15 pts
- Archivos (mkfile, mkdir, cat): 10 pts

---

## 📝 **Checklist Final Pre-Entrega**

### **Funcionalidad:**
- [ ] Todos los comandos obligatorios implementados
- [ ] Todos los reportes obligatorios funcionando
- [ ] Frontend carga scripts y ejecuta comandos
- [ ] Backend responde correctamente via API
- [ ] Sistema maneja errores apropiadamente

### **Calidad:**
- [ ] Código bien comentado y estructurado
- [ ] No hay memory leaks o archivos abiertos
- [ ] Validación de parámetros completa
- [ ] Manejo consistente de errores

### **Entrega:**
- [ ] Repositorio GitHub privado configurado
- [ ] Acceso habilitado para auxiliar de tu sección
- [ ] Documentación completa en /docs
- [ ] README con instrucciones de instalación
- [ ] Último commit antes de 23:59 del 7/9/2025

### **Auxiliares por Sección:**
- **Sección A:** joshi20022021
- **Sección B:** melladodaniel  
- **Sección C:** SaulCerezo
- **Sección D:** kmsu

---

## 🚀 **Siguientes Pasos Inmediatos**

1. **Configura tu entorno de desarrollo**
   - Instala Go y tu framework frontend preferido
   - Crea el repositorio GitHub privado
   - Configura la estructura de carpetas

2. **Implementa la funcionalidad básica**
   - Empieza con mkdisk y las estructuras MBR
   - Prueba con reportes MBR y DISK
   - Avanza gradualmente según el plan de sprints

3. **Prueba continuamente**
   - Cada comando debe probarse antes de continuar
   - Usa los reportes para validar que todo funciona
   - Mantén un script de prueba actualizado