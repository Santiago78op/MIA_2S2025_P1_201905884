package estructuras

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
	PartStatus      byte     `binary:"little"`
	PartType        byte     `binary:"little"`
	PartFit         byte     `binary:"little"`
	PartStart       int64    `binary:"little"`
	PartSize        int64    `binary:"little"`
	PartName        [16]byte `binary:"little"`
	PartCorrelativo int64    `binary:"little"`
	PartID          int64    `binary:"little"`
}

// NewPartition crea una nueva partición con valores iniciales
func NewPartition(tipo byte, fit byte, start, size int64, name string) *Partition {

}

/*
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
*/
