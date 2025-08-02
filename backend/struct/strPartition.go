package estructuras

import (
	"fmt"
	"strings"
)

// Una partición es una división lógica de un disco, que el
// sistema operativo puede tratar como un disco independiente.
/*
Nombre            | Tipo       | Descripción
------------------|------------|---------------------------------------------------------------------------------------------------------------------------------
part_status       | char       | Indica si la partición está montada o no
part_type         | char       | Indica el tipo de partición, primaria o extendida. Tendrá los valores P o E
part_fit          | char       | Tipo de ajuste de la partición. Tendrá los valores B (Best), F (First) o W (Worst)
part_start        | int        | Indica en qué byte del disco inicia la partición
part_s            | int        | Contiene el tamaño total de la partición en bytes
part_name         | char[16]   | Nombre de la partición
part_correlative  | int        | Indica el correlativo de la partición, inicialmente -1 hasta que sea montado (luego la primera partición montada empezará en 1 e irán incrementando)
part_id           | char[4]    | Indica el ID de la partición generada al montar esta partición
*/
type Partition struct {
	PartStatus      byte     `binary:"little"` // Indica si la partición está montada o no
	PartType        byte     `binary:"little"` // Tipo de partición: P (Primaria), E (Extendida), L (Lógica)
	PartFit         byte     `binary:"little"` // Tipo de ajuste: B (Best), F (First), W (Worst)
	PartStart       int64    `binary:"little"` // Byte donde inicia la partición
	PartSize        int64    `binary:"little"` // Tamaño total de la partición en bytes
	PartName        [16]byte `binary:"little"` // Nombre de la partición (máximo 16 caracteres)
	PartCorrelativo int64    `binary:"little"` // Correlativo de la partición (-1 hasta que sea montada)
	PartID          [4]byte  `binary:"little"` // ID de la partición generada al montar
}

// NewPartition crea una nueva partición con valores iniciales
func NewPartition(tipo byte, fit byte, start, size int64, name string) *Partition {
	p := &Partition{
		PartStatus:      StatusActiva,
		PartType:        tipo,
		PartFit:         fit,
		PartStart:       start,
		PartSize:        size,
		PartCorrelativo: -1, // Inicialmente -1 hasta que sea montada
	}

	// Copiar nombre limitando a 16 bytes
	nameBytes := []byte(name)
	if len(nameBytes) > 16 {
		nameBytes = nameBytes[:16]
	}
	copy(p.PartName[:], nameBytes)

	// Inicializar PartID como vacío
	for i := range p.PartID {
		p.PartID[i] = 0
	}

	return p
}

// NewEmptyPartition crea una partición vacía/inactiva
func NewEmptyPartition() *Partition {
	p := &Partition{
		PartStatus:      StatusInactiva,
		PartType:        0,
		PartFit:         0,
		PartStart:       0,
		PartSize:        0,
		PartCorrelativo: -1,
	}

	// Limpiar nombre
	for i := range p.PartName {
		p.PartName[i] = 0
	}

	// Limpiar ID
	for i := range p.PartID {
		p.PartID[i] = 0
	}

	return p
}

// GetName obtiene el nombre de la partición como string
func (p *Partition) GetName() string {
	nameBytes := p.PartName[:]

	// Encontrar el primer byte nulo para terminar la cadena
	for i, b := range nameBytes {
		if b == 0 {
			nameBytes = nameBytes[:i]
			break
		}
	}

	return strings.TrimSpace(string(nameBytes))
}

// SetName establece el nombre de la partición
func (p *Partition) SetName(name string) {
	// Limpiar el array primero
	for i := range p.PartName {
		p.PartName[i] = 0
	}

	// Copiar el nuevo nombre limitando a 16 bytes
	nameBytes := []byte(name)
	if len(nameBytes) > 16 {
		nameBytes = nameBytes[:16]
	}
	copy(p.PartName[:], nameBytes)
}

// GetID obtiene el ID de la partición como string
func (p *Partition) GetID() string {
	idBytes := p.PartID[:]

	// Encontrar el primer byte nulo para terminar la cadena
	for i, b := range idBytes {
		if b == 0 {
			idBytes = idBytes[:i]
			break
		}
	}

	return strings.TrimSpace(string(idBytes))
}

// SetID establece el ID de la partición
func (p *Partition) SetID(id string) {
	// Limpiar el array primero
	for i := range p.PartID {
		p.PartID[i] = 0
	}

	// Copiar el nuevo ID limitando a 4 bytes
	idBytes := []byte(id)
	if len(idBytes) > 4 {
		idBytes = idBytes[:4]
	}
	copy(p.PartID[:], idBytes)
}

// IsEmpty verifica si la partición está vacía/inactiva
func (p *Partition) IsEmpty() bool {
	return p.PartStatus == StatusInactiva
}

// IsActive verifica si la partición está activa
func (p *Partition) IsActive() bool {
	return p.PartStatus == StatusActiva
}

// IsMounted verifica si la partición está montada (correlativo >= 1)
func (p *Partition) IsMounted() bool {
	return p.PartCorrelativo >= 1
}

// IsPrimary verifica si es una partición primaria
func (p *Partition) IsPrimary() bool {
	return p.PartType == PartitionTypePrimaria
}

// IsExtended verifica si es una partición extendida
func (p *Partition) IsExtended() bool {
	return p.PartType == PartitionTypeExtendida
}

// IsLogical verifica si es una partición lógica
func (p *Partition) IsLogical() bool {
	return p.PartType == PartitionTypeLogica
}

// GetTypeString obtiene el tipo de partición como string
func (p *Partition) GetTypeString() string {
	switch p.PartType {
	case PartitionTypePrimaria:
		return "Primaria"
	case PartitionTypeExtendida:
		return "Extendida"
	case PartitionTypeLogica:
		return "Lógica"
	default:
		return "Desconocida"
	}
}

// GetFitString obtiene el tipo de ajuste como string
func (p *Partition) GetFitString() string {
	switch p.PartFit {
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

// GetStatusString obtiene el estado como string
func (p *Partition) GetStatusString() string {
	switch p.PartStatus {
	case StatusActiva:
		return "Activa"
	case StatusInactiva:
		return "Inactiva"
	default:
		return "Desconocido"
	}
}

// GetEndPosition calcula la posición final de la partición
func (p *Partition) GetEndPosition() int64 {
	return p.PartStart + p.PartSize
}

// Contains verifica si una posición está dentro de la partición
func (p *Partition) Contains(position int64) bool {
	return position >= p.PartStart && position < p.GetEndPosition()
}

// Overlaps verifica si esta partición se solapa con otra
func (p *Partition) Overlaps(other *Partition) bool {
	if p.IsEmpty() || other.IsEmpty() {
		return false
	}

	return !(p.GetEndPosition() <= other.PartStart || other.GetEndPosition() <= p.PartStart)
}

// ValidatePartition verifica si la partición es válida
func (p *Partition) ValidatePartition() error {
	if p.IsEmpty() {
		return nil // Una partición vacía es válida
	}

	// Validar tipo de partición
	if p.PartType != PartitionTypePrimaria &&
		p.PartType != PartitionTypeExtendida &&
		p.PartType != PartitionTypeLogica {
		return fmt.Errorf("tipo de partición inválido: %c", p.PartType)
	}

	// Validar tipo de ajuste
	if p.PartFit != PartitionFitBest &&
		p.PartFit != PartitionFitFirst &&
		p.PartFit != PartitionFitWorst {
		return fmt.Errorf("tipo de ajuste inválido: %c", p.PartFit)
	}

	// Validar tamaño
	if p.PartSize <= 0 {
		return fmt.Errorf("el tamaño de la partición debe ser mayor que 0")
	}

	// Validar posición de inicio
	if p.PartStart < 0 {
		return fmt.Errorf("la posición de inicio no puede ser negativa")
	}

	// Validar nombre
	name := p.GetName()
	if name == "" {
		return fmt.Errorf("la partición debe tener un nombre")
	}

	return nil
}

// Mount monta la partición asignándole un correlativo y ID
func (p *Partition) Mount(correlativo int64, id string) error {
	if p.IsEmpty() {
		return fmt.Errorf("no se puede montar una partición vacía")
	}

	if p.IsMounted() {
		return fmt.Errorf("la partición ya está montada")
	}

	p.PartCorrelativo = correlativo
	p.SetID(id)

	return nil
}

// Unmount desmonta la partición
func (p *Partition) Unmount() {
	p.PartCorrelativo = -1
	// Limpiar ID
	for i := range p.PartID {
		p.PartID[i] = 0
	}
}

// Delete elimina la partición (la marca como inactiva)
func (p *Partition) Delete() {
	p.PartStatus = StatusInactiva
	p.PartType = 0
	p.PartFit = 0
	p.PartStart = 0
	p.PartSize = 0
	p.PartCorrelativo = -1

	// Limpiar nombre
	for i := range p.PartName {
		p.PartName[i] = 0
	}

	// Limpiar ID
	for i := range p.PartID {
		p.PartID[i] = 0
	}
}

// Clone crea una copia de la partición
func (p *Partition) Clone() *Partition {
	clone := &Partition{
		PartStatus:      p.PartStatus,
		PartType:        p.PartType,
		PartFit:         p.PartFit,
		PartStart:       p.PartStart,
		PartSize:        p.PartSize,
		PartCorrelativo: p.PartCorrelativo,
	}

	// Copiar arrays
	copy(clone.PartName[:], p.PartName[:])
	copy(clone.PartID[:], p.PartID[:])

	return clone
}

// String implementa la interfaz Stringer para debugging
func (p *Partition) String() string {
	if p.IsEmpty() {
		return "Partition{Empty}"
	}

	return fmt.Sprintf("Partition{Name: %s, Type: %s, Status: %s, Start: %d, Size: %d, ID: %s, Correlativo: %d}",
		p.GetName(),
		p.GetTypeString(),
		p.GetStatusString(),
		p.PartStart,
		p.PartSize,
		p.GetID(),
		p.PartCorrelativo,
	)
}
