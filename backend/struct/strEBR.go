package estructuras

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"unsafe"
)

// EBR representa la estructura del Extended Boot Record
// Es un descriptor de una unidad lógica que contiene información y datos de la misma
// y apunta hacia el espacio donde se escribirá el siguiente EBR (como una lista enlazada)
/*
| Nombre     | Tipo     | Descripción                                                                    |
|------------|----------|--------------------------------------------------------------------------------|
| part_mount | char     | Indica si la partición está montada o no                                      |
| part_fit   | char     | Tipo de ajuste de la partición. B (Best), F (First), W (Worst)               |
| part_start | int      | Indica en qué byte del disco inicia la partición                             |
| part_s     | int      | Contiene el tamaño total de la partición en bytes                            |
| part_next  | int      | Byte en el que está el próximo EBR. -1 si no hay siguiente                   |
| part_name  | char[16] | Nombre de la partición                                                         |
*/
type EBR struct {
	PartMount byte     `binary:"little"` // Indica si la partición está montada o no
	PartFit   byte     `binary:"little"` // Tipo de ajuste: B (Best), F (First), W (Worst)
	PartStart int64    `binary:"little"` // Byte donde inicia la partición
	PartSize  int64    `binary:"little"` // Tamaño total de la partición en bytes
	PartNext  int64    `binary:"little"` // Posición del próximo EBR (-1 si no hay siguiente)
	PartName  [16]byte `binary:"little"` // Nombre de la partición (máximo 16 caracteres)
}

// Tamaño del EBR en bytes
const EBR_SIZE = int(unsafe.Sizeof(EBR{}))

// NewEBR crea un nuevo EBR con valores iniciales
func NewEBR(fit byte, start, size int64, name string, next int64) *EBR {
	ebr := &EBR{
		PartMount: StatusInactiva, // Por defecto inactiva hasta que se monte
		PartFit:   fit,
		PartStart: start,
		PartSize:  size,
		PartNext:  next,
	}

	// Copiar nombre limitando a 16 bytes
	nameBytes := []byte(name)
	if len(nameBytes) > 16 {
		nameBytes = nameBytes[:16]
	}
	copy(ebr.PartName[:], nameBytes)

	return ebr
}

// NewEmptyEBR crea un EBR vacío
func NewEmptyEBR() *EBR {
	ebr := &EBR{
		PartMount: StatusInactiva,
		PartFit:   0,
		PartStart: 0,
		PartSize:  0,
		PartNext:  -1,
	}

	// Limpiar nombre
	for i := range ebr.PartName {
		ebr.PartName[i] = 0
	}

	return ebr
}

// GetName obtiene el nombre de la partición lógica como string
func (e *EBR) GetName() string {
	nameBytes := e.PartName[:]

	// Encontrar el primer byte nulo para terminar la cadena
	for i, b := range nameBytes {
		if b == 0 {
			nameBytes = nameBytes[:i]
			break
		}
	}

	return strings.TrimSpace(string(nameBytes))
}

// SetName establece el nombre de la partición lógica
func (e *EBR) SetName(name string) {
	// Limpiar el array primero
	for i := range e.PartName {
		e.PartName[i] = 0
	}

	// Copiar el nuevo nombre limitando a 16 bytes
	nameBytes := []byte(name)
	if len(nameBytes) > 16 {
		nameBytes = nameBytes[:16]
	}
	copy(e.PartName[:], nameBytes)
}

// IsEmpty verifica si el EBR está vacío
func (e *EBR) IsEmpty() bool {
	return e.PartSize == 0 || e.GetName() == ""
}

// IsMounted verifica si la partición lógica está montada
func (e *EBR) IsMounted() bool {
	return e.PartMount == StatusActiva
}

// HasNext verifica si hay un siguiente EBR en la cadena
func (e *EBR) HasNext() bool {
	return e.PartNext != -1 && e.PartNext > 0
}

// GetEndPosition calcula la posición final de la partición lógica
func (e *EBR) GetEndPosition() int64 {
	return e.PartStart + e.PartSize
}

// GetFitString obtiene el tipo de ajuste como string
func (e *EBR) GetFitString() string {
	switch e.PartFit {
	case PartitionFitBest:
		return "Best Fit"
	case PartitionFitFirst:
		return "First Fit"
	case PartitionFitWorst:
		return "Worst Fit"
	default:
		return "Desconocido"
	}
}

// Mount monta la partición lógica
func (e *EBR) Mount() {
	e.PartMount = StatusActiva
}

// Unmount desmonta la partición lógica
func (e *EBR) Unmount() {
	e.PartMount = StatusInactiva
}

// ValidateEBR verifica si el EBR es válido
func (e *EBR) ValidateEBR() error {
	if e.IsEmpty() {
		return nil // Un EBR vacío es válido
	}

	// Validar tipo de ajuste
	if e.PartFit != PartitionFitBest &&
		e.PartFit != PartitionFitFirst &&
		e.PartFit != PartitionFitWorst {
		return fmt.Errorf("tipo de ajuste inválido en EBR: %c", e.PartFit)
	}

	// Validar tamaño
	if e.PartSize <= 0 {
		return fmt.Errorf("el tamaño de la partición lógica debe ser mayor que 0")
	}

	// Validar posición de inicio
	if e.PartStart < 0 {
		return fmt.Errorf("la posición de inicio no puede ser negativa")
	}

	// Validar nombre
	name := e.GetName()
	if name == "" {
		return fmt.Errorf("la partición lógica debe tener un nombre")
	}

	return nil
}

// SerializeEBR convierte EBR a bytes
func SerializeEBR(ebr *EBR) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, ebr)
	if err != nil {
		return nil, fmt.Errorf("error al serializar EBR: %v", err)
	}
	return buf.Bytes(), nil
}

// DeserializeEBR convierte bytes a EBR
func DeserializeEBR(data []byte) (*EBR, error) {
	if len(data) < EBR_SIZE {
		return nil, fmt.Errorf("datos insuficientes para EBR: necesarios %d, recibidos %d", EBR_SIZE, len(data))
	}

	ebr := &EBR{}
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, ebr)
	if err != nil {
		return nil, fmt.Errorf("error al deserializar EBR: %v", err)
	}
	return ebr, nil
}

// WriteEBR escribe un EBR en el disco en la posición especificada
func WriteEBR(path string, ebr *EBR, position int64) error {
	// Serializar el EBR
	ebrData, err := SerializeEBR(ebr)
	if err != nil {
		return fmt.Errorf("error al serializar EBR: %v", err)
	}

	// Escribir el EBR en el disco
	err = WriteToDisk(path, ebrData, position)
	if err != nil {
		return fmt.Errorf("error al escribir EBR en el disco: %v", err)
	}

	return nil
}

// ReadEBR lee un EBR desde el disco en la posición especificada
func ReadEBR(path string, position int64) (*EBR, error) {
	// Leer los datos del EBR
	data, err := ReadFromDisk(path, position, EBR_SIZE)
	if err != nil {
		return nil, fmt.Errorf("error al leer EBR: %v", err)
	}

	// Deserializar el EBR
	ebr, err := DeserializeEBR(data)
	if err != nil {
		return nil, fmt.Errorf("error al deserializar EBR: %v", err)
	}

	return ebr, nil
}

// ReadAllEBRs lee todos los EBRs en una cadena desde una posición inicial
func ReadAllEBRs(path string, startPosition int64) ([]*EBR, error) {
	var ebrs []*EBR
	currentPosition := startPosition

	for currentPosition != -1 {
		// Leer el EBR actual
		ebr, err := ReadEBR(path, currentPosition)
		if err != nil {
			return nil, fmt.Errorf("error al leer EBR en posición %d: %v", currentPosition, err)
		}

		// Agregar el EBR a la lista si no está vacío
		if !ebr.IsEmpty() {
			ebrs = append(ebrs, ebr)
		}

		// Mover a la siguiente posición
		currentPosition = ebr.PartNext

		// Prevenir bucles infinitos
		if len(ebrs) > 100 {
			return nil, fmt.Errorf("demasiados EBRs en la cadena, posible corrupción")
		}
	}

	return ebrs, nil
}

// FindEBRByName busca un EBR por nombre en una cadena de EBRs
func FindEBRByName(path string, startPosition int64, name string) (*EBR, int64, error) {
	currentPosition := startPosition

	for currentPosition != -1 {
		// Leer el EBR actual
		ebr, err := ReadEBR(path, currentPosition)
		if err != nil {
			return nil, -1, fmt.Errorf("error al leer EBR en posición %d: %v", currentPosition, err)
		}

		// Verificar si es el EBR que buscamos
		if !ebr.IsEmpty() && ebr.GetName() == name {
			return ebr, currentPosition, nil
		}

		// Mover a la siguiente posición
		currentPosition = ebr.PartNext
	}

	return nil, -1, fmt.Errorf("no se encontró una partición lógica con el nombre: %s", name)
}

// String implementa la interfaz Stringer para debugging
func (e *EBR) String() string {
	if e.IsEmpty() {
		return "EBR{Empty}"
	}

	return fmt.Sprintf("EBR{Name: %s, Mount: %t, Start: %d, Size: %d, Next: %d, Fit: %s}",
		e.GetName(),
		e.IsMounted(),
		e.PartStart,
		e.PartSize,
		e.PartNext,
		e.GetFitString(),
	)
}
