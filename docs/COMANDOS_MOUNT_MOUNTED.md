# 🔧 Separación de Comandos MOUNT y MOUNTED

## 📋 Resumen de Cambios

Se han separado exitosamente los comandos MOUNT y MOUNTED que anteriormente estaban juntos en el archivo `mounting.go` en comandos independientes para mejorar la organización y mantenibilidad del código.

## 🗂️ Archivos Creados

### 1. `backend/command/disk/mount.go`
Contiene toda la funcionalidad relacionada con el montaje de particiones:

**Funciones principales:**
- `Mount(path, name string) error` - Monta una partición en el sistema
- `Unmount(id string) error` - Desmonta una partición del sistema
- `GetMountedPartitionByID(id string) (*MountedPartition, error)` - Busca una partición montada por ID
- `SetCarnetSuffix(suffix string)` - Configura el sufijo del carnet para generar IDs

**Estructuras:**
- `MountedPartition` - Representa una partición montada
- `MountSystem` - Maneja el sistema de montaje global

**Funciones auxiliares:**
- `validateMountParams()` - Valida parámetros del comando mount
- `findAndMountPartition()` - Busca y prepara el montaje de una partición
- `generatePartitionID()` - Genera IDs únicos para particiones
- `updatePartitionInDisk()` - Actualiza información de montaje en disco
- `unmountPartitionInDisk()` - Limpia información de montaje en disco

### 2. `backend/command/disk/mounted.go`
Contiene toda la funcionalidad para mostrar particiones montadas:

**Funciones principales:**
- `GetMountedPartitions() []*MountedPartition` - Obtiene lista de particiones montadas
- `ShowMountedPartitions() string` - Muestra particiones en formato texto legible
- `ShowMountedPartitionsTable() string` - Muestra particiones en formato tabla ASCII
- `ShowMountedPartitionsDetailed() string` - Muestra reporte detallado
- `ShowMountedPartitionsJSON() map[string]interface{}` - Formato JSON para API

**Funciones auxiliares:**
- `groupPartitionsByDisk()` - Agrupa particiones por disco
- `formatSize()` - Formatea tamaños de bytes a formato legible
- `formatTimestamp()` - Convierte timestamps a formato legible
- `truncateString()` - Trunca cadenas para tablas
- `getFileName()` - Extrae nombre de archivo de ruta

## 🔄 Integración con el Parser

El parser de comandos (`backend/command/commandResult.go`) ya estaba configurado para manejar ambos comandos por separado:

```go
case "mount":
    return cp.executeMount(params)
case "mounted":
    return cp.executeMounted(params)
```

Las funciones `executeMount()` y `executeMounted()` ya llamaban a las funciones correctas, por lo que no fue necesario modificar el parser.

## 📋 Parámetros de los Comandos

### Comando MOUNT
| Parámetro | Obligatorio | Descripción |
|-----------|-------------|-------------|
| `-path` | ✅ SÍ | Ruta del disco que se montará |
| `-name` | ✅ SÍ | Nombre de la partición a montar |

**Ejemplo de uso:**
```bash
mount -path="/ruta/disco.mia" -name="Particion1"
```

### Comando MOUNTED
No requiere parámetros. Muestra todas las particiones montadas.

**Ejemplo de uso:**
```bash
mounted
```

## 🆔 Sistema de IDs de Particiones

El sistema genera IDs únicos basados en:
- **Formato:** `[últimos 2 dígitos carnet][número partición][letra disco]`
- **Ejemplo:** Para carnet 201905884: `841A`, `842A`, `841B`, etc.
- **Configuración:** Usar `SetCarnetSuffix("84")` para configurar

## 🏗️ Arquitectura del Sistema de Montaje

```
MountSystem (Singleton Global)
├── mountedPartitions: map[string]*MountedPartition
├── diskPartitionCount: map[string]int
├── nextLetter: byte
└── carnetSuffix: string

MountedPartition
├── ID: string (generado automáticamente)
├── Name: string
├── Path: string
├── Type: string (Primary/Extended/Logical)
├── Size: int64
├── PartitionIndex: int
├── EBRPosition: int64
├── Correlative: int64
├── MountTime: string
└── DiskSignature: int64
```

## 🔄 Estados de Montaje

1. **Partición Desmontada:** `PartCorrelativo = -1`, `PartID = ""`
2. **Partición Montada:** `PartCorrelativo >= 1`, `PartID = "841A"`
3. **Persistencia:** Los cambios se escriben directamente al MBR del disco

## 🛡️ Validaciones Implementadas

### Comando MOUNT
- ✅ Validación de parámetros obligatorios
- ✅ Verificación de integridad del disco
- ✅ Validación de existencia de la partición
- ✅ Verificación de que la partición no esté ya montada
- ✅ Solo permite montar particiones primarias (según especificaciones)

### Comando MOUNTED
- ✅ Manejo de caso sin particiones montadas
- ✅ Ordenación consistente por ID
- ✅ Agrupación por disco para mejor presentación
- ✅ Múltiples formatos de salida (texto, tabla, detallado, JSON)

## 📊 Formatos de Salida del Comando MOUNTED

### 1. Formato Estándar (`ShowMountedPartitions()`)
```
=== PARTICIONES MONTADAS ===
Total: 2 particiones

📁 DISCO: /ruta/disco1.mia
--------------------------------------------------
  🔹 ID: 841A
     Nombre: Particion1
     Tipo: Primaria
     Tamaño: 10.00 MB
     Correlativo: 1
     Montado: Unix: 1672531200
```

### 2. Formato Tabla (`ShowMountedPartitionsTable()`)
```
┌──────────┬────────────────┬──────────────┬──────────────┬────────────────┬─────────────┐
│    ID    │     NOMBRE     │     TIPO     │   TAMAÑO     │     DISCO      │ CORRELATIVO │
├──────────┼────────────────┼──────────────┼──────────────┼────────────────┼─────────────┤
│ 841A     │ Particion1     │ Primaria     │ 10.00 MB     │ disco1.mia     │ 1           │
└──────────┴────────────────┴──────────────┴──────────────┴────────────────┴─────────────┘
```

### 3. Formato JSON para API (`ShowMountedPartitionsJSON()`)
```json
{
  "partitions": [...],
  "total": 2,
  "stats": {
    "total_mounted": 2,
    "unique_disks": 1,
    "next_letter": "A",
    "carnet_suffix": "84"
  }
}
```

## 🧪 Testing y Verificación

### Compilación Exitosa
✅ El backend compila sin errores después de la separación
✅ No hay conflictos de funciones duplicadas
✅ Todas las dependencias están correctamente resueltas

### Funcionalidades Preservadas
✅ Sistema de montaje global funcional
✅ Generación de IDs única
✅ Persistencia en disco de información de montaje
✅ Validaciones de integridad
✅ Múltiples formatos de salida

## 🔧 Mantenimiento

### Para Agregar Nuevas Funcionalidades
- **Montaje:** Modificar `mount.go`
- **Visualización:** Modificar `mounted.go`
- **Parser:** Solo si se agregan nuevos comandos relacionados

### Para Debugging
- Logs disponibles en cada operación con nivel apropiado (INFO, WARNING, ERROR, SUCCESS)
- Sistema de estadísticas accesible via `GetMountSystemStats()`
- Función `ClearMountSystem()` para limpiar estado en testing

## 📈 Mejoras Implementadas

1. **Separación de Responsabilidades:** Cada archivo tiene una responsabilidad específica
2. **Múltiples Formatos:** El comando MOUNTED ahora soporta varios formatos de salida
3. **Mejor Documentación:** Cada función tiene documentación clara
4. **Validaciones Robustas:** Manejo de errores mejorado
5. **Estadísticas del Sistema:** Información detallada del estado de montaje
6. **Compatibilidad API:** Formatos JSON para integración con frontend

---

**Fecha de implementación:** Agosto 2025  
**Archivos modificados:** `mount.go`, `mounted.go`, `mounting.go` (eliminado)  
**Compatibilidad:** Totalmente compatible con sistema existente