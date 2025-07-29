# 📚 Guía de Estructuras - Sistema de Archivos EXT2

## 🎯 **Objetivo**
Esta guía define todas las estructuras de datos necesarias para implementar un simulador de sistema de archivos EXT2 compatible con el proyecto MIA.

---

## 📋 **Índice de Estructuras**

1. [MBR (Master Boot Record)](#1-mbr-master-boot-record)
2. [Partition (Partición)](#2-partition-partición)
3. [EBR (Extended Boot Record)](#3-ebr-extended-boot-record)
4. [SuperBloque](#4-superbloque)
5. [Inodo](#5-inodo)
6. [Bloque de Datos](#6-bloque-de-datos)
7. [Journaling](#7-journaling)
8. [Estructuras de Reportes](#8-estructuras-de-reportes)

---

## 1. MBR (Master Boot Record)

### 📝 **Descripción:**
Estructura principal que debe estar en el **primer sector** (primeros 512 bytes) del disco. Contiene información del disco y tabla de particiones.

### 🏗️ **Estructura:**
```go
package estructuras

import (
    "encoding/binary"
    "fmt"
    "time"
    "unsafe"
)

// MBR representa la estructura del Master Boot Record
type MBR struct {
    // Tamaño total del disco en bytes
    MbrTamanio int64 `binary:"little"`
    
    // Fecha y hora de creación (timestamp Unix)
    MbrFechaCreacion int64 `binary:"little"`
    
    // Número random único que identifica el disco
    MbrDiskSignature int64 `binary:"little"`
    
    // Tipo de ajuste: 'B'=Best, 'F'=First, 'W'=Worst
    MbrFit byte
    
    // Tabla de particiones (máximo 4)
    Particiones [4]Partition
}

// Constantes para tipos de ajuste
const (
    FitBest  byte = 'B'
    FitFirst byte = 'F' 
    FitWorst byte = 'W'
)

// Tamaño fijo del MBR en bytes
const MBR_SIZE = int(unsafe.Sizeof(MBR{}))
```

### 🛠️ **Métodos recomendados:**
```go
// NewMBR crea un nuevo MBR
func NewMBR(tamanio int64, fit byte) *MBR {
    return &MBR{
        MbrTamanio:       tamanio,
        MbrFechaCreacion: time.Now().Unix(),
        MbrDiskSignature: generateRandomSignature(),
        MbrFit:          fit,
        Particiones:     [4]Partition{}, // Inicializar vacío
    }
}

// GetParticionLibre encuentra la primera partición disponible
func (m *MBR) GetParticionLibre() *Partition {
    for i := range m.Particiones {
        if m.Particiones[i].PartStatus == 0 {
            return &m.Particiones[i]
        }
    }
    return nil
}

// ValidarFit verifica si el tipo de ajuste es válido
func (m *MBR) ValidarFit() bool {
    return m.MbrFit == FitBest || m.MbrFit == FitFirst || m.MbrFit == FitWorst
}
```

---

## 2. Partition (Partición)

### 📝 **Descripción:**
Cada partición en la tabla del MBR. Puede ser primaria, extendida o lógica.

### 🏗️ **Estructura:**
```go
// Partition representa una entrada en la tabla de particiones
type Partition struct {
    // Estado: 0=inactiva, 1=activa
    PartStatus byte
    
    // Tipo: 'P'=Primaria, 'E'=Extendida, 'L'=Lógica
    PartType byte
    
    // Ajuste: 'B'=Best, 'F'=First, 'W'=Worst
    PartFit byte
    
    // Byte donde inicia la partición
    PartStart int64 `binary:"little"`
    
    // Tamaño en bytes de la partición
    PartSize int64 `binary:"little"`
    
    // Nombre de la partición (16 caracteres max)
    PartName [16]byte
}

// Constantes para tipos de partición
const (
    PartPrimaria  byte = 'P'
    PartExtendida byte = 'E'
    PartLogica    byte = 'L'
)

// Constantes para estado
const (
    StatusInactiva byte = 0
    StatusActiva   byte = 1
)
```

### 🛠️ **Métodos recomendados:**
```go
// NewPartition crea una nueva partición
func NewPartition(tipo byte, fit byte, start, size int64, name string) *Partition {
    p := &Partition{
        PartStatus: StatusActiva,
        PartType:   tipo,
        PartFit:    fit,
        PartStart:  start,
        PartSize:   size,
    }
    copy(p.PartName[:], []byte(name))
    return p
}

// GetName obtiene el nombre como string
func (p *Partition) GetName() string {
    return string(p.PartName[:])
}

// IsEmpty verifica si la partición está vacía
func (p *Partition) IsEmpty() bool {
    return p.PartStatus == StatusInactiva
}
```

---

## 3. EBR (Extended Boot Record)

### 📝 **Descripción:**
Estructura para particiones lógicas dentro de una partición extendida.

### 🏗️ **Estructura:**
```go
// EBR representa el Extended Boot Record
type EBR struct {
    // Estado: 0=inactiva, 1=activa
    PartStatus byte
    
    // Ajuste: 'B'=Best, 'F'=First, 'W'=Worst
    PartFit byte
    
    // Byte donde inicia la partición lógica
    PartStart int64 `binary:"little"`
    
    // Tamaño en bytes de la partición lógica
    PartSize int64 `binary:"little"`
    
    // Byte donde está el siguiente EBR (-1 si es el último)
    PartNext int64 `binary:"little"`
    
    // Nombre de la partición lógica
    PartName [16]byte
}

const EBR_SIZE = int(unsafe.Sizeof(EBR{}))
```

---

## 4. SuperBloque

### 📝 **Descripción:**
Contiene información sobre el sistema de archivos EXT2/EXT3.

### 🏗️ **Estructura:**
```go
// SuperBloque contiene metadatos del sistema de archivos
type SuperBloque struct {
    // Número de inodos en el sistema
    SInodosCount int64 `binary:"little"`
    
    // Número de bloques en el sistema  
    SBloquesCount int64 `binary:"little"`
    
    // Número de bloques libres
    SFreeBloquesCount int64 `binary:"little"`
    
    // Número de inodos libres
    SFreeInodosCount int64 `binary:"little"`
    
    // Fecha de montaje (timestamp Unix)
    SMtime int64 `binary:"little"`
    
    // Fecha de desmontaje (timestamp Unix) 
    SUmtime int64 `binary:"little"`
    
    // Número de veces que se ha montado
    SMntCount int64 `binary:"little"`
    
    // Tamaño del bloque en bytes
    SBlockSize int64 `binary:"little"`
    
    // Tamaño del inodo en bytes
    SInodeSize int64 `binary:"little"`
    
    // Byte donde inicia la tabla de inodos
    SInodeStart int64 `binary:"little"`
    
    // Byte donde inicia el bitmap de inodos
    SBminodeStart int64 `binary:"little"`
    
    // Byte donde inicia el bitmap de bloques
    SBmblockStart int64 `binary:"little"`
    
    // Byte donde inician los bloques
    SBlockStart int64 `binary:"little"`
    
    // Indica si tiene journaling: 0=EXT2, 1=EXT3
    SMagic int64 `binary:"little"`
}

// Constantes para el sistema de archivos
const (
    SUPERBLOCK_SIZE = int(unsafe.Sizeof(SuperBloque{}))
    EXT2_MAGIC      = 0
    EXT3_MAGIC      = 1
    DEFAULT_BLOCK_SIZE = 1024
    DEFAULT_INODE_SIZE = 128
)
```

---

## 5. Inodo

### 📝 **Descripción:**
Estructura que contiene metadatos de archivos y directorios.

### 🏗️ **Estructura:**
```go
// Inodo contiene metadatos de archivos y directorios
type Inodo struct {
    // Identificador único del usuario propietario
    IUid int64 `binary:"little"`
    
    // Identificador único del grupo propietario  
    IGid int64 `binary:"little"`
    
    // Tamaño del archivo en bytes
    ISize int64 `binary:"little"`
    
    // Última fecha de acceso
    IAtime int64 `binary:"little"`
    
    // Fecha de creación
    ICtime int64 `binary:"little"`
    
    // Última fecha de modificación
    IMtime int64 `binary:"little"`
    
    // Array de bloques directos e indirectos
    IBlock [15]int64 `binary:"little"`
    
    // Tipo: 0=archivo, 1=directorio
    IType int64 `binary:"little"`
    
    // Permisos del archivo/directorio
    IPerm int64 `binary:"little"`
}

// Constantes para tipos de inodo
const (
    INODE_FILE int64 = 0
    INODE_DIR  int64 = 1
    INODE_SIZE = int(unsafe.Sizeof(Inodo{}))
)

// Constantes para permisos (formato octal)
const (
    PERM_664 int64 = 0664
    PERM_775 int64 = 0775
)
```

---

## 6. Bloque de Datos

### 📝 **Descripción:**
Estructuras para diferentes tipos de bloques de datos.

### 🏗️ **Estructura:**

#### **6.1 Bloque de Carpeta:**
```go
// BloqueCarpeta contiene entradas de directorio
type BloqueCarpeta struct {
    // Array de contenidos del directorio
    BContent [4]Content
}

// Content representa una entrada en el directorio  
type Content struct {
    // Nombre del archivo/carpeta
    BName [12]byte
    
    // Apuntador al inodo
    BInodo int64 `binary:"little"`
}
```

#### **6.2 Bloque de Archivo:**
```go
// BloqueArchivo contiene datos de archivo
type BloqueArchivo struct {
    // Contenido del archivo (64 caracteres)
    BContent [64]byte
}
```

#### **6.3 Bloque de Apuntadores:**
```go
// BloqueApuntadores para indirección
type BloqueApuntadores struct {
    // Array de apuntadores a otros bloques
    BPointers [16]int64 `binary:"little"`
}
```

### 🛠️ **Constantes de tamaños:**
```go
const (
    BLOQUE_SIZE     = 1024
    CONTENT_SIZE    = int(unsafe.Sizeof(Content{}))
    ENTRIES_PER_DIR = 4
    POINTERS_PER_BLOCK = 16
)
```

---

## 7. Journaling

### 📝 **Descripción:**
Estructura para el registro de transacciones (EXT3).

### 🏗️ **Estructura:**
```go
// Journaling registra las operaciones para EXT3
type Journaling struct {
    // Tipo de operación
    TipoOperacion [10]byte
    
    // Tipo: 0=archivo, 1=carpeta
    Tipo int64 `binary:"little"`
    
    // Nombre del archivo/carpeta
    Nombre [100]byte
    
    // Contenido si es archivo
    Contenido [100]byte
    
    // Fecha de la operación
    Fecha int64 `binary:"little"`
    
    // Propietario
    Propietario [10]byte
    
    // Permisos
    Permisos int64 `binary:"little"`
}

// Constantes para operaciones
const (
    OP_MKDIR  = "mkdir"
    OP_MKFILE = "mkfile" 
    OP_REMOVE = "remove"
    OP_RENAME = "rename"
    OP_CHMOD  = "chmod"
    OP_CHOWN  = "chown"
)
```

---

## 8. Estructuras de Reportes

### 📝 **Descripción:**
Estructuras auxiliares para generar reportes con Graphviz.

### 🏗️ **Estructuras:**
```go
// ReporteMBR información para reporte MBR
type ReporteMBR struct {
    TamanioMBR       int64
    FechaCreacion    string
    DiskSignature    int64
    Fit              string
    Particiones      []ReporteParticion
}

// ReporteParticion información de partición para reporte
type ReporteParticion struct {
    Status    string
    Tipo      string
    Fit       string
    Start     int64
    Size      int64
    Nombre    string
}

// ReporteDisk información del disco para reporte
type ReporteDisk struct {
    MBR        ReporteMBR
    Partitions []ReporteParticionDetalle
    FreeSpace  []ReporteEspacioLibre
}

// ReporteParticionDetalle detalle de partición
type ReporteParticionDetalle struct {
    Nombre       string
    Tipo         string
    Start        int64
    Size         int64
    PorcentajeUso float64
}

// ReporteEspacioLibre espacio no utilizado
type ReporteEspacioLibre struct {
    Start int64
    Size  int64
}
```

---

## 🔧 **Utilidades de Serialización**

### 📝 **Funciones auxiliares:**
```go
// SerializeMBR convierte MBR a bytes
func SerializeMBR(mbr *MBR) ([]byte, error) {
    buf := new(bytes.Buffer)
    err := binary.Write(buf, binary.LittleEndian, mbr)
    return buf.Bytes(), err
}

// DeserializeMBR convierte bytes a MBR
func DeserializeMBR(data []byte) (*MBR, error) {
    mbr := &MBR{}
    buf := bytes.NewReader(data)
    err := binary.Read(buf, binary.LittleEndian, mbr)
    return mbr, err
}

// WriteToFile escribe estructura al archivo
func WriteToFile(filename string, offset int64, data []byte) error {
    file, err := os.OpenFile(filename, os.O_RDWR, 0644)
    if err != nil {
        return err
    }
    defer file.Close()
    
    _, err = file.WriteAt(data, offset)
    return err
}

// ReadFromFile lee estructura del archivo
func ReadFromFile(filename string, offset int64, size int) ([]byte, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    data := make([]byte, size)
    _, err = file.ReadAt(data, offset)
    return data, err
}
```

---

## 📊 **Mapa de Memoria del Disco**

```
┌─────────────────────────────────────────────────────────────┐
│  SECTOR 0: MBR (512 bytes)                                  │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ Tamaño │ Fecha │ Signature │ Fit │ Particiones[4] │    │
│  └─────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────┤
│  PARTICIÓN 1 (Primaria)                                    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ SuperBloque │ Bitmap Inodos │ Bitmap Bloques │...   │    │
│  └─────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────┤
│  PARTICIÓN 2 (Extendida)                                   │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ EBR → Partición Lógica 1 → EBR → Partición Lógica 2│    │
│  └─────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────┤
│  PARTICIÓN 3 (Primaria)                                    │
├─────────────────────────────────────────────────────────────┤
│  PARTICIÓN 4 (Libre)                                       │
└─────────────────────────────────────────────────────────────┘
```

---

## 🎯 **Mejores Prácticas Implementadas**

### ✅ **Nomenclatura Go:**
- PascalCase para campos exportados
- camelCase para funciones y variables locales

### ✅ **Serialización:**
- Tags `binary:"little"` para consistencia
- Funciones de serialización/deserialización

### ✅ **Constantes:**
- Tamaños fijos definidos
- Valores mágicos como constantes
- Tipos enumerados

### ✅ **Validaciones:**
- Métodos de validación en cada estructura
- Verificación de rangos y tipos

### ✅ **Documentación:**
- Comentarios explicativos
- Referencias a documentación del proyecto

---

## 🚀 **Ejemplo de Uso Completo**

```go
// Crear un nuevo disco
mbr := NewMBR(1024*1024*100, FitFirst) // 100MB

// Crear partición primaria
partition := NewPartition(PartPrimaria, FitFirst, 1024, 1024*1024*30, "Particion1")
mbr.Particiones[0] = *partition

// Serializar y escribir al disco
mbrBytes, _ := SerializeMBR(mbr)
WriteToFile("disco1.dk", 0, mbrBytes)

// Leer MBR del disco
data, _ := ReadFromFile("disco1.dk", 0, MBR_SIZE)
mbrLeido, _ := DeserializeMBR(data)
```

---

## 📚 **Referencias**
- Documentación oficial del proyecto MIA
- Especificación EXT2/EXT3 filesystem
- Manual de buenas prácticas en Go

**Fecha de creación:** 29 de julio de 2025  
**Versión:** 1.0
