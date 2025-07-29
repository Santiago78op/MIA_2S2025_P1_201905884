package estructuras

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"
	"unsafe"
)

// MBR representa la estructura del Master Boot Record
// Esta estructura deberá estar en el primer sector del disco
// Nota: Para más detalles sobre la estructura MBR, leer la página 7 de la documentación.
type MBR struct {
	// Tamaño total del disco en bytes
	MbrTamanio int64 `binary:"little"`
	// Fecha y hora de creación del Disco
	MbrFechaCreacion int64 `binary:"little"`
	// Número  random, que identifica de forma única a cada disco
	MbrDiskSignature int64 `binary:"little"`
	// Tipo de ajuste de la partición, -> Tendrá los valores de
	// B (Best), F (First) o W (Worst)
	MbrFit byte
	// Estructura con información de las particiones
	MbrParticiones [4]Partition
}

// Constantes para el tipo de ajuste de partición
const (
	PartitionFitBest  byte = 'B' // Mejor ajuste
	PartitionFitFirst byte = 'F' // Primer ajuste
	PartitionFitWorst byte = 'W' // Peor ajuste
)

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

// Tamaño del MBR en bytes
const MBR_SIZE = int(unsafe.Sizeof(MBR{}))

// NewMBR crea un nuevo MBR con valores iniciales
func NewMBR(tamanio int64, fit byte) *MBR {
	return &MBR{
		MbrTamanio:       tamanio,
		MbrFechaCreacion: time.Now().Unix(),
		MbrDiskSignature: generateRandomSignature(),
		MbrFit:           fit,
		MbrParticiones:   [4]Partition{},
	}
}

// NewPartition crea una nueva partición
func NewPartition(tipo byte, fit byte, start, size int64, name string) *Partition {
	p := &Partition{
		PartStatus: StatusActiva,
		PartType:   tipo,
		PartFit:    fit,
		PartStart:  start,
		PartSize:   size,
	}
	// Copiar nombre limitando a 16 bytes
	nameBytes := []byte(name)
	if len(nameBytes) > 16 {
		nameBytes = nameBytes[:16]
	}
	copy(p.PartName[:], nameBytes)
	return p
}

// GetParticionLibre encuentra la primera partición disponible
func (m *MBR) GetParticionLibre() *Partition {
	for i := range m.MbrParticiones {
		if m.MbrParticiones[i].PartStatus == StatusInactiva {
			return &m.MbrParticiones[i]
		}
	}
	return nil
}

// ValidarFit verifica si el tipo de ajuste es válido
func (m *MBR) ValidarFit() bool {
	return m.MbrFit == PartitionFitBest || m.MbrFit == PartitionFitFirst || m.MbrFit == PartitionFitWorst
}

// GetName obtiene el nombre de la partición como string
func (p *Partition) GetName() string {
	nameBytes := p.PartName[:]
	for i, b := range nameBytes {
		if b == 0 {
			nameBytes = nameBytes[:i]
			break
		}
	}
	return string(nameBytes)
}

// IsEmpty verifica si la partición está vacía
func (p *Partition) IsEmpty() bool {
	return p.PartStatus == StatusInactiva
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
