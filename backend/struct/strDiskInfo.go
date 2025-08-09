package estructuras

import (
	utils "backend/Utils"
	"fmt"
	"os"
	"time"
)

// DiskInfo contiene información sobre un disco
type DiskInfo struct {
	Path             string    `json:"path"`
	Size             int64     `json:"size"`
	CreationDate     time.Time `json:"creation_date"`
	DiskSignature    int64     `json:"disk_signature"`
	Fit              string    `json:"fit"`
	ActivePartitions int       `json:"active_partitions"`
	FreeSpace        int64     `json:"free_space"`
}

// GetDiskInfo obtiene información completa sobre un disco
func GetDiskInfo(path string) (*DiskInfo, error) {
	// Verificar que el archivo existe
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("error al acceder al disco: %v", err)
	}

	// Leer el MBR
	mbr, err := ReadMBR(path)
	if err != nil {
		return nil, fmt.Errorf("error al leer MBR: %v", err)
	}

	// Crear información del disco
	diskInfo := &DiskInfo{
		Path:             path,
		Size:             fileInfo.Size(),
		CreationDate:     time.Unix(mbr.MbrFechaCreacion, 0),
		DiskSignature:    mbr.MbrDiskSignature,
		Fit:              string(mbr.MbrFit),
		ActivePartitions: mbr.CountActivePartitions(),
		FreeSpace:        mbr.GetFreeSpace(),
	}

	return diskInfo, nil
}

// ValidateDiskIntegrity valida la integridad de un disco
func ValidateDiskIntegrity(path string) error {
	utils.LogInfo("ValidateDisk", fmt.Sprintf("Validando integridad del disco: %s", path))

	// Verificar que el archivo existe
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("el archivo del disco no existe: %v", err)
	}

	// Verificar que no es un directorio
	if fileInfo.IsDir() {
		return fmt.Errorf("la ruta especificada es un directorio, no un archivo")
	}

	// Leer y validar el MBR
	mbr, err := ReadMBR(path)
	if err != nil {
		return fmt.Errorf("error al leer MBR: %v", err)
	}

	// Validar que el tamaño del archivo coincide con el MBR
	if fileInfo.Size() != mbr.MbrTamanio {
		return fmt.Errorf("el tamaño del archivo (%d) no coincide con el tamaño en el MBR (%d)",
			fileInfo.Size(), mbr.MbrTamanio)
	}

	// Validar el tipo de ajuste
	if !mbr.ValidarFit() {
		return fmt.Errorf("tipo de ajuste inválido en el MBR: %c", mbr.MbrFit)
	}

	// Validar particiones
	for i, partition := range mbr.MbrParticiones {
		if !partition.IsEmpty() {
			if err := partition.ValidatePartition(); err != nil {
				return fmt.Errorf("partición %d inválida: %v", i, err)
			}

			// Verificar que la partición no excede el tamaño del disco
			if partition.GetEndPosition() > mbr.MbrTamanio {
				return fmt.Errorf("la partición %d excede el tamaño del disco", i)
			}
		}
	}

	// Verificar solapamientos entre particiones
	for i := 0; i < len(mbr.MbrParticiones); i++ {
		if mbr.MbrParticiones[i].IsEmpty() {
			continue
		}
		for j := i + 1; j < len(mbr.MbrParticiones); j++ {
			if mbr.MbrParticiones[j].IsEmpty() {
				continue
			}
			if mbr.MbrParticiones[i].Overlaps(&mbr.MbrParticiones[j]) {
				return fmt.Errorf("las particiones %d y %d se solapan", i, j)
			}
		}
	}

	utils.LogSuccess("ValidateDisk", "El disco pasó todas las validaciones de integridad")
	return nil
}

// BackupMBR crea una copia de seguridad del MBR
func BackupMBR(diskPath, backupPath string) error {
	utils.LogInfo("BackupMBR", fmt.Sprintf("Creando backup del MBR: %s -> %s", diskPath, backupPath))

	// Leer el MBR
	mbr, err := ReadMBR(diskPath)
	if err != nil {
		return fmt.Errorf("error al leer MBR para backup: %v", err)
	}

	// Serializar el MBR
	mbrData, err := SerializeMBR(mbr)
	if err != nil {
		return fmt.Errorf("error al serializar MBR para backup: %v", err)
	}

	// Escribir el backup
	err = os.WriteFile(backupPath, mbrData, 0644)
	if err != nil {
		return fmt.Errorf("error al escribir backup del MBR: %v", err)
	}

	utils.LogSuccess("BackupMBR", fmt.Sprintf("Backup del MBR creado exitosamente en: %s", backupPath))
	return nil
}

// RestoreMBR restaura un MBR desde una copia de seguridad
func RestoreMBR(diskPath, backupPath string) error {
	utils.LogInfo("RestoreMBR", fmt.Sprintf("Restaurando MBR: %s <- %s", diskPath, backupPath))

	// Leer el backup
	mbrData, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("error al leer backup del MBR: %v", err)
	}

	// Verificar que los datos son válidos
	_, err = DeserializeMBR(mbrData)
	if err != nil {
		return fmt.Errorf("el backup del MBR está corrupto: %v", err)
	}

	// Escribir el MBR al disco
	err = WriteToDisk(diskPath, mbrData, 0)
	if err != nil {
		return fmt.Errorf("error al restaurar MBR al disco: %v", err)
	}

	// Validar que la restauración fue exitosa
	err = ValidateDiskIntegrity(diskPath)
	if err != nil {
		return fmt.Errorf("la restauración del MBR falló la validación: %v", err)
	}

	utils.LogSuccess("RestoreMBR", "MBR restaurado exitosamente")
	return nil
}

// CleanDisk limpia un disco eliminando todas las particiones
func CleanDisk(path string) error {
	utils.LogInfo("CleanDisk", fmt.Sprintf("Limpiando disco: %s", path))

	// Leer el MBR actual
	mbr, err := ReadMBR(path)
	if err != nil {
		return fmt.Errorf("error al leer MBR: %v", err)
	}

	// Mantener información básica del disco pero limpiar particiones
	originalSize := mbr.MbrTamanio
	originalSignature := mbr.MbrDiskSignature
	originalFit := mbr.MbrFit

	// Crear nuevo MBR limpio conservando la firma original
	cleanMBR := NewMBR(originalSize, originalFit)
	cleanMBR.MbrDiskSignature = originalSignature // Mantener la firma original

	// Serializar y escribir el MBR limpio para no generar una nueva firma
	mbrData, err := SerializeMBR(cleanMBR)
	if err != nil {
		return fmt.Errorf("error al serializar MBR limpio: %v", err)
	}

	if err = WriteToDisk(path, mbrData, 0); err != nil {
		return fmt.Errorf("error al escribir MBR limpio: %v", err)
	}

	utils.LogSuccess("CleanDisk", "Disco limpiado exitosamente - todas las particiones eliminadas")
	return nil
}

// GetDiskUsage calcula el uso del disco en porcentaje
func GetDiskUsage(path string) (float64, error) {
	mbr, err := ReadMBR(path)
	if err != nil {
		return 0, fmt.Errorf("error al leer MBR: %v", err)
	}

	usedSpace := int64(MBR_SIZE) // Espacio del MBR

	for _, partition := range mbr.MbrParticiones {
		if partition.IsActive() {
			usedSpace += partition.PartSize
		}
	}

	usage := (float64(usedSpace) / float64(mbr.MbrTamanio)) * 100
	return usage, nil
}

// ListPartitions lista todas las particiones activas de un disco
func ListPartitions(path string) ([]*Partition, error) {
	mbr, err := ReadMBR(path)
	if err != nil {
		return nil, fmt.Errorf("error al leer MBR: %v", err)
	}

	var partitions []*Partition

	for i := range mbr.MbrParticiones {
		if mbr.MbrParticiones[i].IsActive() {
			partitions = append(partitions, &mbr.MbrParticiones[i])
		}
	}

	return partitions, nil
}

// FindPartition busca una partición por nombre
func FindPartition(path, name string) (*Partition, int, error) {
	mbr, err := ReadMBR(path)
	if err != nil {
		return nil, -1, fmt.Errorf("error al leer MBR: %v", err)
	}

	for i := range mbr.MbrParticiones {
		if mbr.MbrParticiones[i].IsActive() && mbr.MbrParticiones[i].GetName() == name {
			return &mbr.MbrParticiones[i], i, nil
		}
	}

	return nil, -1, fmt.Errorf("no se encontró una partición con el nombre: %s", name)
}

// GetExtendedPartition obtiene la partición extendida si existe
func GetExtendedPartition(path string) (*Partition, int, error) {
	mbr, err := ReadMBR(path)
	if err != nil {
		return nil, -1, fmt.Errorf("error al leer MBR: %v", err)
	}

	for i := range mbr.MbrParticiones {
		if mbr.MbrParticiones[i].IsActive() && mbr.MbrParticiones[i].IsExtended() {
			return &mbr.MbrParticiones[i], i, nil
		}
	}

	return nil, -1, fmt.Errorf("no se encontró una partición extendida")
}

// ZeroDisk llena el disco con ceros (excepto el MBR)
func ZeroDisk(path string) error {
	utils.LogInfo("ZeroDisk", fmt.Sprintf("Limpiando contenido del disco con ceros: %s", path))

	// Leer el MBR para preservarlo
	mbr, err := ReadMBR(path)
	if err != nil {
		return fmt.Errorf("error al leer MBR: %v", err)
	}

	// Abrir el archivo para escritura
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error al abrir archivo para limpiar: %v", err)
	}
	defer file.Close()

	// Posicionarse después del MBR
	_, err = file.Seek(int64(MBR_SIZE), 0)
	if err != nil {
		return fmt.Errorf("error al posicionarse en el archivo: %v", err)
	}

	// Crear buffer de ceros
	bufferSize := 1024 * 1024 // 1MB buffer
	zeroBuffer := make([]byte, bufferSize)

	// Calcular cuántos bytes limpiar
	remainingBytes := mbr.MbrTamanio - int64(MBR_SIZE)

	// Escribir ceros en chunks
	for remainingBytes > 0 {
		writeSize := bufferSize
		if remainingBytes < int64(bufferSize) {
			writeSize = int(remainingBytes)
		}

		_, err = file.Write(zeroBuffer[:writeSize])
		if err != nil {
			return fmt.Errorf("error al escribir ceros: %v", err)
		}

		remainingBytes -= int64(writeSize)
	}

	utils.LogSuccess("ZeroDisk", "Contenido del disco limpiado con ceros exitosamente")
	return nil
}

// CompareMBR compara dos MBRs y retorna las diferencias
func CompareMBR(path1, path2 string) ([]string, error) {
	mbr1, err := ReadMBR(path1)
	if err != nil {
		return nil, fmt.Errorf("error al leer MBR de %s: %v", path1, err)
	}

	mbr2, err := ReadMBR(path2)
	if err != nil {
		return nil, fmt.Errorf("error al leer MBR de %s: %v", path2, err)
	}

	var differences []string

	if mbr1.MbrTamanio != mbr2.MbrTamanio {
		differences = append(differences, fmt.Sprintf("Tamaño diferente: %d vs %d", mbr1.MbrTamanio, mbr2.MbrTamanio))
	}

	if mbr1.MbrDiskSignature != mbr2.MbrDiskSignature {
		differences = append(differences, fmt.Sprintf("Firma diferente: %d vs %d", mbr1.MbrDiskSignature, mbr2.MbrDiskSignature))
	}

	if mbr1.MbrFit != mbr2.MbrFit {
		differences = append(differences, fmt.Sprintf("Ajuste diferente: %c vs %c", mbr1.MbrFit, mbr2.MbrFit))
	}

	// Comparar particiones
	for i := 0; i < 4; i++ {
		p1 := &mbr1.MbrParticiones[i]
		p2 := &mbr2.MbrParticiones[i]

		if p1.IsEmpty() && p2.IsEmpty() {
			continue
		}

		if p1.IsEmpty() != p2.IsEmpty() {
			differences = append(differences, fmt.Sprintf("Partición %d: una existe y la otra no", i))
			continue
		}

		if p1.GetName() != p2.GetName() {
			differences = append(differences, fmt.Sprintf("Partición %d nombre: %s vs %s", i, p1.GetName(), p2.GetName()))
		}

		if p1.PartSize != p2.PartSize {
			differences = append(differences, fmt.Sprintf("Partición %d tamaño: %d vs %d", i, p1.PartSize, p2.PartSize))
		}

		if p1.PartStart != p2.PartStart {
			differences = append(differences, fmt.Sprintf("Partición %d inicio: %d vs %d", i, p1.PartStart, p2.PartStart))
		}
	}

	return differences, nil
}

// GetDiskStatistics obtiene estadísticas detalladas del disco
func GetDiskStatistics(path string) (map[string]interface{}, error) {
	mbr, err := ReadMBR(path)
	if err != nil {
		return nil, fmt.Errorf("error al leer MBR: %v", err)
	}

	stats := make(map[string]interface{})

	// Información básica
	stats["total_size"] = mbr.MbrTamanio
	stats["creation_date"] = time.Unix(mbr.MbrFechaCreacion, 0).Format(time.RFC3339)
	stats["disk_signature"] = mbr.MbrDiskSignature
	stats["fit_type"] = string(mbr.MbrFit)

	// Estadísticas de particiones
	stats["total_partitions"] = 4
	stats["active_partitions"] = mbr.CountActivePartitions()
	stats["free_partitions"] = 4 - mbr.CountActivePartitions()

	// Estadísticas de espacio
	stats["free_space"] = mbr.GetFreeSpace()
	stats["used_space"] = mbr.MbrTamanio - mbr.GetFreeSpace()

	usage, _ := GetDiskUsage(path)
	stats["usage_percentage"] = usage

	// Información de particiones activas
	var partitionInfo []map[string]interface{}
	for i, partition := range mbr.MbrParticiones {
		if partition.IsActive() {
			pInfo := map[string]interface{}{
				"index":       i,
				"name":        partition.GetName(),
				"type":        partition.GetTypeString(),
				"status":      partition.GetStatusString(),
				"start":       partition.PartStart,
				"size":        partition.PartSize,
				"end":         partition.GetEndPosition(),
				"fit":         partition.GetFitString(),
				"mounted":     partition.IsMounted(),
				"correlative": partition.PartCorrelativo,
				"id":          partition.GetID(),
			}
			partitionInfo = append(partitionInfo, pInfo)
		}
	}
	stats["partitions"] = partitionInfo

	return stats, nil
}
