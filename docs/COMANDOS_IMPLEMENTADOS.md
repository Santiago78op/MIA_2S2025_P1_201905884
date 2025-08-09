# 🛠️ Comandos Implementados - Sistema de Archivos
## 🆕 Cambios recientes (agosto 2025)
- Se documenta el proceso interno de `mkdisk` y sus validaciones.
- Se agregan ejemplos de logs generados por comandos.
- Se detallan las estructuras MBR y Partition utilizadas.
- Se actualiza el roadmap de comandos y fases de desarrollo.

## 📋 **Resumen**
Documentación detallada de todos los comandos implementados en el simulador de sistema de archivos EXT2.

---

## ✅ **1. MKDISK - Crear Disco Virtual**

### 📝 **Descripción**
Crea un archivo binario que simula un disco duro. El archivo tendrá extensión `.mia` y será inicializado con ceros binarios.

### 🔧 **Sintaxis**
```bash
mkdisk -size [número] -unit [K|M] -fit [BF|FF|WF] -path [ruta]
```

### 📋 **Parámetros**

| Parámetro | Obligatorio | Descripción | Valores |
|-----------|-------------|-------------|---------|
| `-size` | ✅ **SÍ** | Tamaño del disco | Número positivo > 0 |
| `-unit` | ❌ **NO** | Unidades del tamaño | `K` (Kilobytes), `M` (Megabytes) |
| `-fit` | ❌ **NO** | Tipo de ajuste | `BF` (Best), `FF` (First), `WF` (Worst) |
| `-path` | ✅ **SÍ** | Ruta del archivo | Ruta completa con extensión `.mia` |

### 🔧 **Valores por Defecto**
- **unit:** `M` (Megabytes)
- **fit:** `FF` (First Fit)

### ✅ **Ejemplos de Uso**
```bash
# Crear disco de 10 MB con ajuste First Fit
mkdisk -size 10 -path /home/usuario/disco1.mia

# Crear disco de 2048 KB con ajuste Best Fit
mkdisk -size 2048 -unit K -fit BF -path /tmp/disco2.mia

# Crear disco de 50 MB con ajuste Worst Fit
mkdisk -size 50 -unit M -fit WF -path /discos/disco3.mia
```

### 🚨 **Validaciones**
- **size:** Debe ser un número positivo mayor que 0
- **unit:** Solo acepta `K` o `M`
- **fit:** Solo acepta `BF`, `FF` o `WF`
- **path:** Se crean los directorios automáticamente si no existen

### 📊 **Cálculo de Tamaños**
- **Kilobytes (K):** size × 1024 bytes
- **Megabytes (M):** size × 1024 × 1024 bytes

### 🔄 **Proceso Interno**
1. Validar todos los parámetros
2. Calcular tamaño en bytes según la unidad
3. Crear el archivo en la ruta especificada
4. Escribir ceros binarios hasta alcanzar el tamaño
5. Generar log de éxito/error

### 📤 **Logging**
```json
{
  "type": "INFO",
  "command": "MkDisk",
  "message": "Disco creado con éxito en /ruta/disco.mia de tamaño 10485760 bytes con ajuste FF",
  "time": "1693574400"
}
```

---

## 🔄 **Comandos en Desarrollo**

### 🚧 **2. RMDISK - Eliminar Disco** (Pendiente)
- **Función:** Eliminar archivos de disco `.mia`
- **Parámetros:** `-path [ruta]`

### 🚧 **3. FDISK - Gestión de Particiones** (Pendiente)
- **Función:** Crear, eliminar y gestionar particiones
- **Tipos:** Primaria, Extendida, Lógica

### 🚧 **4. MOUNT - Montar Partición** (Pendiente)
- **Función:** Montar particiones para uso del sistema
- **Parámetros:** `-path [disco]`, `-name [partición]`

---

## 🏗️ **Estructuras Relacionadas**

### **MBR (Master Boot Record)**
```go
type MBR struct {
    MbrTamanio       int64        // Tamaño total del disco
    MbrFechaCreacion int64        // Timestamp de creación
    MbrDiskSignature int64        // Identificador único
    MbrFit           byte         // Tipo de ajuste (B/F/W)
    MbrParticiones   [4]Partition // Tabla de particiones
}
```

### **Partition**
```go
type Partition struct {
    PartStatus byte      // Estado: activa (A) o inactiva (I)
    PartType   byte      // Tipo: primaria (P), extendida (E), lógica (L)
    PartFit    byte      // Ajuste: BF, FF, WF
    PartStart  int64     // Byte de inicio
    PartSize   int64     // Tamaño en bytes
    PartName   [16]byte  // Nombre de la partición
}
```

---

## 📈 **Roadmap de Implementación**

### **Fase 1 - Gestión de Discos** ✅ **COMPLETADA**
- [x] `mkdisk` - Crear discos

### **Fase 2 - Gestión de Particiones** 🔄 **EN DESARROLLO**
- [ ] `rmdisk` - Eliminar discos
- [ ] `fdisk` - Gestión de particiones
- [ ] `mount` - Montar particiones

### **Fase 3 - Sistema de Archivos** 🔄 **PENDIENTE**
- [ ] `mkfs` - Formatear con EXT2
- [ ] `login`/`logout` - Autenticación

### **Fase 4 - Gestión de Usuarios** 🔄 **PENDIENTE**
- [ ] `mkgrp`/`rmgrp` - Gestión de grupos
- [ ] `mkusr`/`rmusr` - Gestión de usuarios

### **Fase 5 - Archivos y Directorios** 🔄 **PENDIENTE**
- [ ] `mkdir` - Crear directorios
- [ ] `mkfile` - Crear archivos
- [ ] `cat` - Mostrar contenido

---

## 🔧 **Notas Técnicas**

### **Compilación**
```bash
cd backend
go build -o backend.exe .
```

### **Ejecución**
```bash
./backend.exe
```

### **Testing**
```bash
# Probar comando via API
curl -X POST http://localhost:8080/api/execute \
  -H "Content-Type: application/json" \
  -d '{"command": "mkdisk -size 10 -path /tmp/test.mia"}'
```

### **Logs en Tiempo Real**
- **WebSocket:** `ws://localhost:8080/ws`
- **SSE:** `http://localhost:8080/api/logs/stream`
- **HTTP:** `http://localhost:8080/api/logs`
