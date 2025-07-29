package estructuras

// Paquete necesario para manejar tiempos

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
