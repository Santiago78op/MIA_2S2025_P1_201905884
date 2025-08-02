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
