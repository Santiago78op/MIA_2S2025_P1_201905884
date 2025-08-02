package estructuras

import (
	utils "backend/Utils"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
	"time"
	"unsafe"
)

// MBR representa la estructura del Master Boot Record
// Esta estructura deberá estar en el primer sector del disco
/*
| Nombre             | Tipo   | Descripción                                         |
|--------------------|--------|-----------------------------------------------------|
| mbr_tamanio        | int64  | Tamaño total del disco en bytes                     |
| mbr_fecha_creacion | int64  | Fecha y hora de creación del disco                  |
| mbr_dsk_signature  | int64  | Número random, que identifica de forma única a cada disco |
| dsk_fit            | char   | Tipo de ajuste de la partición. B (Best), F (First), W (Worst) |
| mbr_partitions     | partition[4] | Estructura con información de las 4 particiones           |
*/
type MBR struct {
	MbrTamanio       int64 `binary:"little"`
	MbrFechaCreacion int64 `binary:"little"`
	MbrDiskSignature int64 `binary:"little"`
	MbrFit           byte  `binary:"little"`
	MbrParticiones   [4]Partition
}

// Constantes para el tipo de ajuste de partición
const (
	PartitionFitBest  byte = 'B' // Mejor ajuste
	PartitionFitFirst byte = 'F' // Primer ajuste
	PartitionFitWorst byte = 'W' // Peor ajuste
)

// Constantes para estado
const (
	StatusInactiva byte = 0
	StatusActiva   byte = 1
)

// Constantes para tipo de partición
const (
	PartitionTypePrimaria  byte = 'P'
	PartitionTypeExtendida byte = 'E'
	PartitionTypeLogica    byte = 'L'
)

// Tamaño del MBR en bytes
const MBR_SIZE = int(unsafe.Sizeof(MBR{}))

// NewMBR crea un nuevo MBR con valores iniciales
func NewMBR(tamanio int64, fit byte) *MBR {
	mbr := &MBR{
		MbrTamanio:       tamanio,
		MbrFechaCreacion: time.Now().Unix(),
		MbrDiskSignature: generateRandomSignature(),
		MbrFit:           fit,
		MbrParticiones:   [4]Partition{},
	}

	// Inicializar particiones como inactivas
	for i := range mbr.MbrParticiones {
		mbr.MbrParticiones[i].PartStatus = StatusInactiva // Partición inactiva
		mbr.MbrParticiones[i].PartCorrelativo = -1        // Correlativo no asignado
	}

	// Retornar el MBR recién creado
	return mbr
}

// WriteMBR escribe el MBR en el disco especificado
func WriteMBR(path string, sizeInBytes int64, fit string) error {
	// Validar y convertir el fit
	fitByte := ValidateFit(fit)

	if fitByte == 0 {
		utils.LogError("WriteMBR", "Ajuste no válido, use B, F o W")
		return fmt.Errorf("tipo de ajuste no válido: %s", fit)
	}

	// Crear un nuevo MBR
	mbr := NewMBR(sizeInBytes, fitByte)

	// Serializar el MBR
	mbrData, err := SerializeMBR(mbr)
	if err != nil {
		utils.LogError("WriteMBR", fmt.Sprintf("Error al serializar el MBR: %v", err))
		return fmt.Errorf("error al serializar el MBR: %v", err)
	}

	// Escribir el MBR en el disco
	err = WriteToDisk(path, mbrData, 0)
	if err != nil {
		utils.LogError("WriteMBR", fmt.Sprintf("Error al escribir el MBR en el disco: %v", err))
		return fmt.Errorf("error al escribir el MBR en el disco: %v", err)
	}

	utils.LogSuccess("WriteMBR", fmt.Sprintf("MBR escrito exitosamente en %s", path))
	return nil
}

// ReadMBR lee el MBR desde el disco especificado
func ReadMBR(path string) (*MBR, error) {
	// Leer los primeros bytes del archivo (tamaño del MBR)
	data, err := ReadFromDisk(path, 0, MBR_SIZE)
	if err != nil {
		utils.LogError("ReadMBR", fmt.Sprintf("Error al leer el MBR desde el disco: %v", err))
		return nil, fmt.Errorf("error al leer el MBR desde el disco: %v", err)
	}

	// Deserializar los datos leídos en un MBR
	mbr, err := DeserializeMBR(data)
	if err != nil {
		utils.LogError("ReadMBR", fmt.Sprintf("Error al deserializar el MBR: %v", err))
		return nil, fmt.Errorf("error al deserializar el MBR: %v", err)
	}

	return mbr, nil
}

// generateRandomSignature genera una firma única para el disco
func generateRandomSignature() int64 {
	max := big.NewInt(1<<63 - 1)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return time.Now().UnixNano()
	}
	return n.Int64()
}

// ValidateFit asigna el valor a fit validando que sea el parametro correcto
func ValidateFit(fit string) byte {
	switch fit {
	case "BF", "B":
		return PartitionFitBest
	case "FF", "F":
		return PartitionFitFirst
	case "WF", "W":
		return PartitionFitWorst
	default:
		return 0 // Valor inválido
	}
}

// SerializeMBR convierte MBR a bytes
func SerializeMBR(mbr *MBR) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, mbr)
	if err != nil {
		return nil, fmt.Errorf("error al serializar MBR: %v", err)
	}
	return buf.Bytes(), nil
}

// DeserializeMBR convierte bytes a MBR
func DeserializeMBR(data []byte) (*MBR, error) {
	if len(data) < MBR_SIZE {
		return nil, fmt.Errorf("datos insuficientes para MBR")
	}
	mbr := &MBR{}
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, mbr)
	if err != nil {
		return nil, fmt.Errorf("error al deserializar MBR: %v", err)
	}
	return mbr, nil
}

// ValidarFit verifica si el tipo de ajuste del MBR es valido
func (m *MBR) ValidarFit() bool {
	return m.MbrFit == PartitionFitBest || m.MbrFit == PartitionFitFirst || m.MbrFit == PartitionFitWorst
}

// GetParticionLibre encuentra la primera partición disponible en el MBR
func (m *MBR) GetParticionLibre() *Partition {
	for i := range m.MbrParticiones {
		if m.MbrParticiones[i].PartStatus == StatusInactiva {
			return &m.MbrParticiones[i]
		}
	}
	return nil
}

// GetParticionByName busca una partición por nombre
func (m *MBR) GetParticionByName(name string) *Partition {
	for i := range m.MbrParticiones {
		if m.MbrParticiones[i].GetName() == name && m.MbrParticiones[i].PartStatus == StatusActiva {
			return &m.MbrParticiones[i]
		}
	}
	return nil
}

// HasExtendedPartition verifica si ya existe una partición extendida
func (m *MBR) HasExtendedPartition() bool {
	for i := range m.MbrParticiones {
		if m.MbrParticiones[i].PartStatus == StatusActiva &&
			m.MbrParticiones[i].PartType == PartitionTypeExtendida {
			return true
		}
	}
	return false
}

// CountActivePartitions cuenta las particiones activas (primarias y extendidas)
func (m *MBR) CountActivePartitions() int {
	count := 0
	for i := range m.MbrParticiones {
		if m.MbrParticiones[i].PartStatus == StatusActiva &&
			(m.MbrParticiones[i].PartType == PartitionTypePrimaria ||
				m.MbrParticiones[i].PartType == PartitionTypeExtendida) {
			count++
		}
	}
	return count
}

// GetFreeSpace calcula el espacio libre disponible en el disco
func (m *MBR) GetFreeSpace() int64 {
	usedSpace := int64(MBR_SIZE) // Espacio usado por el MBR

	for i := range m.MbrParticiones {
		if m.MbrParticiones[i].PartStatus == StatusActiva {
			usedSpace += m.MbrParticiones[i].PartSize
		}
	}

	return m.MbrTamanio - usedSpace
}

// FindBestFitPosition encuentra la mejor posición para una nueva partición según el algoritmo especificado
func (m *MBR) FindBestFitPosition(size int64) (int64, error) {
	// Crear lista de espacios libres
	freeSpaces := m.getFreeSpaces()

	if len(freeSpaces) == 0 {
		return -1, fmt.Errorf("no hay espacios libres disponibles")
	}

	switch m.MbrFit {
	case PartitionFitFirst: // First Fit
		for _, space := range freeSpaces {
			if space.Size >= size {
				return space.Start, nil
			}
		}
	case PartitionFitBest: // Best Fit
		bestSpace := FreeSpace{Size: m.MbrTamanio + 1} // Inicializar con tamaño imposible
		found := false
		for _, space := range freeSpaces {
			if space.Size >= size && space.Size < bestSpace.Size {
				bestSpace = space
				found = true
			}
		}
		if found {
			return bestSpace.Start, nil
		}
	case PartitionFitWorst: // Worst Fit
		worstSpace := FreeSpace{Size: -1} // Inicializar con tamaño inválido
		for _, space := range freeSpaces {
			if space.Size >= size && space.Size > worstSpace.Size {
				worstSpace = space
			}
		}
		if worstSpace.Size >= size {
			return worstSpace.Start, nil
		}
	}

	return -1, fmt.Errorf("no se encontró espacio suficiente para la partición")
}

// FreeSpace representa un espacio libre en el disco
type FreeSpace struct {
	Start int64
	Size  int64
}

// getFreeSpaces obtiene todos los espacios libres en el disco
func (m *MBR) getFreeSpaces() []FreeSpace {
	var spaces []FreeSpace
	var occupiedSpaces []FreeSpace

	// Agregar el espacio ocupado por el MBR
	occupiedSpaces = append(occupiedSpaces, FreeSpace{Start: 0, Size: int64(MBR_SIZE)})

	// Agregar espacios ocupados por particiones activas
	for i := range m.MbrParticiones {
		if m.MbrParticiones[i].PartStatus == StatusActiva {
			occupiedSpaces = append(occupiedSpaces, FreeSpace{
				Start: m.MbrParticiones[i].PartStart,
				Size:  m.MbrParticiones[i].PartSize,
			})
		}
	}

	// Ordenar espacios ocupados por posición de inicio
	for i := 0; i < len(occupiedSpaces)-1; i++ {
		for j := i + 1; j < len(occupiedSpaces); j++ {
			if occupiedSpaces[i].Start > occupiedSpaces[j].Start {
				occupiedSpaces[i], occupiedSpaces[j] = occupiedSpaces[j], occupiedSpaces[i]
			}
		}
	}

	// Encontrar espacios libres entre espacios ocupados
	currentPos := int64(0)
	for _, occupied := range occupiedSpaces {
		if currentPos < occupied.Start {
			spaces = append(spaces, FreeSpace{
				Start: currentPos,
				Size:  occupied.Start - currentPos,
			})
		}
		currentPos = occupied.Start + occupied.Size
	}

	// Verificar si hay espacio libre al final del disco
	if currentPos < m.MbrTamanio {
		spaces = append(spaces, FreeSpace{
			Start: currentPos,
			Size:  m.MbrTamanio - currentPos,
		})
	}

	return spaces
}

// WriteToDisk escribe datos en el disco en la posición especificada
func WriteToDisk(path string, data []byte, offset int64) error {
	// Abrir el archivo del disco
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error al abrir el disco: %v", err)
	}
	defer file.Close()

	// Posicionarse en el offset especificado
	_, err = file.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("error al posicionar el disco: %v", err)
	}

	// Escribir los datos en el disco
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error al escribir en el disco: %v", err)
	}

	return nil
}

// ReadFromDisk lee datos del disco desde la posición especificada
func ReadFromDisk(path string, offset int64, size int) ([]byte, error) {
	// Abrir el archivo del disco
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el disco: %v", err)
	}
	defer file.Close()

	// Posicionarse en el offset especificado
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, fmt.Errorf("error al posicionar el disco: %v", err)
	}

	// Leer los datos del disco
	data := make([]byte, size)
	n, err := file.Read(data)
	if err != nil {
		return nil, fmt.Errorf("error al leer del disco: %v", err)
	}

	if n != size {
		return nil, fmt.Errorf("no se leyeron suficientes bytes del disco: %d de %d", n, size)
	}

	return data, nil
}
