package estructuras

// Una partición es una división lógica de un disco, que el
// sistema operativo puede tratar como un disco independiente.
type Partition struct {
	// Indica si la partición está activa (1) o inactiva (0)
	PartStatus byte `binary:"little"`
	// Inidica el tipo de partición, puede ser primaria (P), extendida (E)
	// Tendrá los valores de P (Primary), E (Extended)
	PartType byte `binary:"little"`
	// Tipo de ajuste de la partición. Tendrá los valores B
    // (Best), F (First) o W (worst)
	PartFit byte `binary:"little"`
	// Indica en qué byte del disco inicia la partición
	PartStart int64 `binary:"little"`
	// Contiene el tamaño total de la partición en bytes
	PartSize int64 `binary:"little"`
	// Nombre de la partición
	PartName [16]byte `binary:"little"`
	// Indica el correlativo de la partición este valor será
    // inicialmente -1 hasta que sea montado (luego la
	// primera partición montada empezará en 1 e irán
	// incrementando)
	PartCorrelativo int64 `binary:"little"`
	// Indica el ID de la partición
	PartID int64 `binary:"little"`
}
