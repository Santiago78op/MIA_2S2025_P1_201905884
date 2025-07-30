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
