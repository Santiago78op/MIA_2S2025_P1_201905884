package disk

/*
 * MOUNT - Este comando montará una partición del disco en el sistema.
 * Cada partición se identificará por un ID único basado en el número de carnet.
 *
 * MOUNTED - Este comando mostrará todas las particiones montadas en memoria.
 */

import (
	utils "backend/Utils"
	estructuras "backend/struct"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

/*
| PARÁMETRO | CATEGORÍA    | DESCRIPCIÓN                                                                                                                                                                                                                                                                                                                                                      |
|-----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| -path     | Obligatorio  | Ruta del disco que se montará en el sistema. Este archivo ya debe existir.                                                                                                                                                                                                                                                                                      |
| -name     | Obligatorio  | Indica el nombre de la partición a cargar. Si no existe debe mostrar error.                                                                                                                                                                                                                                                                                     |
*/

// MountedPartition representa una partición montada en el sistema
type MountedPartition struct {
	ID             string `json:"id"`              // ID generado (ej: 341A, 342B)
	Name           string `json:"name"`            // Nombre de la partición
	Path           string `json:"path"`            // Ruta del disco
	Type           string `json:"type"`            // Tipo: Primary, Extended, Logical
	Size           int64  `json:"size"`            // Tamaño en bytes
	PartitionIndex int    `json:"partition_index"` // Índice en el MBR (para primarias/extendidas)
	EBRPosition    int64  `json:"ebr_position"`    // Posición del EBR (para lógicas)
	Correlative    int64  `json:"correlative"`     // Número correlativo de montaje
	MountTime      string `json:"mount_time"`      // Timestamp de montaje
	DiskSignature  int64  `json:"disk_signature"`  // Firma del disco
}

// MountSystem maneja el sistema de montaje de particiones
type MountSystem struct {
	mutex              sync.RWMutex
	mountedPartitions  map[string]*MountedPartition // Key: ID de partición
	diskPartitionCount map[string]int               // Key: firma del disco, Value: contador de particiones
	nextLetter         byte                         // Próxima letra a usar (A, B, C, ...)
	carnetSuffix       string                       // Últimos dos dígitos del carnet
}

// Instancia global del sistema de montaje
var mountSystem *MountSystem

// init inicializa el sistema de montaje
func init() {
	mountSystem = &MountSystem{
		mountedPartitions:  make(map[string]*MountedPartition),
		diskPartitionCount: make(map[string]int),
		nextLetter:         'A',
		carnetSuffix:       "34", // Últimos dos dígitos del carnet (ajustar según corresponda)
	}
}

// SetCarnetSuffix permite configurar los últimos dos dígitos del carnet
func SetCarnetSuffix(suffix string) {
	mountSystem.mutex.Lock()
	defer mountSystem.mutex.Unlock()

	if len(suffix) == 2 {
		mountSystem.carnetSuffix = suffix
		utils.LogInfo("MOUNT", fmt.Sprintf("Sufijo de carnet configurado: %s", suffix))
	} else {
		utils.LogWarning("MOUNT", "El sufijo del carnet debe tener exactamente 2 dígitos")
	}
}

// Mount monta una partición en el sistema
func Mount(path, name string) error {
	utils.LogInfo("MOUNT", fmt.Sprintf("Iniciando montaje de partición: path=%s, name=%s", path, name))

	// Validar parámetros
	if err := validateMountParams(path, name); err != nil {
		return err
	}

	// Validar que el disco existe y es válido
	if err := estructuras.ValidateDiskIntegrity(path); err != nil {
		utils.LogError("MOUNT", fmt.Sprintf("Error de integridad del disco: %v", err))
		return fmt.Errorf("error de integridad del disco: %v", err)
	}

	// Leer el MBR
	mbr, err := estructuras.ReadMBR(path)
	if err != nil {
		return fmt.Errorf("error al leer MBR: %v", err)
	}

	// Buscar la partición por nombre
	mountedPartition, err := findAndMountPartition(path, name, mbr)
	if err != nil {
		return err
	}

	// Agregar al sistema de montaje
	mountSystem.mutex.Lock()
	defer mountSystem.mutex.Unlock()

	// Verificar que no esté ya montada
	if isPartitionAlreadyMounted(path, name) {
		utils.LogError("MOUNT", fmt.Sprintf("La partición '%s' en '%s' ya está montada", name, path))
		return fmt.Errorf("la partición '%s' ya está montada", name)
	}

	// Generar ID único
	id := generatePartitionID(mbr.MbrDiskSignature)
	mountedPartition.ID = id

	// Actualizar la partición en el disco con el correlativo y ID
	if err := updatePartitionInDisk(path, name, mountedPartition, mbr); err != nil {
		return fmt.Errorf("error al actualizar partición en disco: %v", err)
	}

	// Agregar al sistema de montaje
	mountSystem.mountedPartitions[id] = mountedPartition

	utils.LogSuccess("MOUNT", "Partición montada exitosamente:")
	utils.LogSuccess("MOUNT", fmt.Sprintf("  → ID: %s", id))
	utils.LogSuccess("MOUNT", fmt.Sprintf("  → Nombre: %s", name))
	utils.LogSuccess("MOUNT", fmt.Sprintf("  → Tipo: %s", mountedPartition.Type))
	utils.LogSuccess("MOUNT", fmt.Sprintf("  → Tamaño: %d bytes", mountedPartition.Size))
	utils.LogSuccess("MOUNT", fmt.Sprintf("  → Correlativo: %d", mountedPartition.Correlative))

	return nil
}

// Unmount desmonta una partición del sistema
func Unmount(id string) error {
	utils.LogInfo("UNMOUNT", fmt.Sprintf("Desmontando partición con ID: %s", id))

	mountSystem.mutex.Lock()
	defer mountSystem.mutex.Unlock()

	// Buscar la partición montada
	mountedPartition, exists := mountSystem.mountedPartitions[id]
	if !exists {
		utils.LogError("UNMOUNT", fmt.Sprintf("No se encontró una partición montada con ID: %s", id))
		return fmt.Errorf("no se encontró una partición montada con ID: %s", id)
	}

	// Actualizar la partición en el disco (desmontar)
	if err := unmountPartitionInDisk(mountedPartition); err != nil {
		utils.LogWarning("UNMOUNT", fmt.Sprintf("Advertencia al desmontar en disco: %v", err))
	}

	// Remover del sistema de montaje
	delete(mountSystem.mountedPartitions, id)

	// Decrementar contador del disco
	diskKey := strconv.FormatInt(mountedPartition.DiskSignature, 10)
	if count, exists := mountSystem.diskPartitionCount[diskKey]; exists && count > 0 {
		mountSystem.diskPartitionCount[diskKey] = count - 1
	}

	utils.LogSuccess("UNMOUNT", fmt.Sprintf("Partición %s desmontada exitosamente", id))
	return nil
}

// GetMountedPartitions retorna todas las particiones montadas
func GetMountedPartitions() []*MountedPartition {
	mountSystem.mutex.RLock()
	defer mountSystem.mutex.RUnlock()

	var partitions []*MountedPartition
	for _, partition := range mountSystem.mountedPartitions {
		partitions = append(partitions, partition)
	}

	// Ordenar por ID para consistencia
	sort.Slice(partitions, func(i, j int) bool {
		return partitions[i].ID < partitions[j].ID
	})

	return partitions
}

// GetMountedPartitionByID busca una partición montada por su ID
func GetMountedPartitionByID(id string) (*MountedPartition, error) {
	mountSystem.mutex.RLock()
	defer mountSystem.mutex.RUnlock()

	if partition, exists := mountSystem.mountedPartitions[id]; exists {
		return partition, nil
	}

	return nil, fmt.Errorf("no se encontró una partición montada con ID: %s", id)
}

// validateMountParams valida los parámetros del comando mount
func validateMountParams(path, name string) error {
	if path == "" {
		utils.LogError("MOUNT", "El parámetro -path es obligatorio")
		return fmt.Errorf("el parámetro -path es obligatorio")
	}

	if name == "" {
		utils.LogError("MOUNT", "El parámetro -name es obligatorio")
		return fmt.Errorf("el parámetro -name es obligatorio")
	}

	return nil
}

// findAndMountPartition busca una partición por nombre y prepara el montaje
func findAndMountPartition(path, name string, mbr *estructuras.MBR) (*MountedPartition, error) {
	// Buscar en particiones primarias y extendidas
	for i, partition := range mbr.MbrParticiones {
		if partition.IsActive() && partition.GetName() == name {
			// Verificar que solo se monten particiones primarias (según especificaciones)
			if !partition.IsPrimary() {
				utils.LogError("MOUNT", "Solo se pueden montar particiones primarias")
				return nil, fmt.Errorf("solo se pueden montar particiones primarias")
			}

			return &MountedPartition{
				Name:           name,
				Path:           path,
				Type:           partition.GetTypeString(),
				Size:           partition.PartSize,
				PartitionIndex: i,
				EBRPosition:    -1, // No aplica para primarias
				DiskSignature:  mbr.MbrDiskSignature,
				MountTime:      fmt.Sprintf("%d", getCurrentTimestamp()),
			}, nil
		}
	}

	// Si no se encuentra en primarias, buscar en lógicas (para referencia futura)
	extendedPartition, _, err := estructuras.GetExtendedPartition(path)
	if err == nil {
		ebr, ebrPosition, err := estructuras.FindEBRByName(path, extendedPartition.PartStart, name)
		if err == nil {
			utils.LogWarning("MOUNT", "Particiones lógicas no se implementarán en este proyecto según especificaciones")
			return &MountedPartition{
				Name:           name,
				Path:           path,
				Type:           "Logical",
				Size:           ebr.PartSize,
				PartitionIndex: -1, // No aplica para lógicas
				EBRPosition:    ebrPosition,
				DiskSignature:  mbr.MbrDiskSignature,
				MountTime:      fmt.Sprintf("%d", getCurrentTimestamp()),
			}, nil
		}
	}

	utils.LogError("MOUNT", fmt.Sprintf("No se encontró una partición con el nombre '%s'", name))
	return nil, fmt.Errorf("no se encontró una partición con el nombre '%s'", name)
}

// generatePartitionID genera un ID único para la partición
func generatePartitionID(diskSignature int64) string {
	diskKey := strconv.FormatInt(diskSignature, 10)

	// Obtener o inicializar contador para este disco
	count, exists := mountSystem.diskPartitionCount[diskKey]
	if !exists {
		count = 0
	}

	// Incrementar contador
	count++
	mountSystem.diskPartitionCount[diskKey] = count

	// Si es una partición del mismo disco, incrementar número
	// Si es de otro disco, usar la siguiente letra y reiniciar en 1
	if count == 1 {
		// Primera partición de este disco, puede ser nueva letra
		if len(mountSystem.diskPartitionCount) > 1 {
			// Hay otros discos, usar siguiente letra
			mountSystem.nextLetter++
			if mountSystem.nextLetter > 'Z' {
				mountSystem.nextLetter = 'A' // Reiniciar si se acaban las letras
			}
		}
	}

	// Generar ID: últimos 2 dígitos del carnet + número + letra
	id := fmt.Sprintf("%s%d%c", mountSystem.carnetSuffix, count, mountSystem.nextLetter)

	utils.LogInfo("MOUNT", fmt.Sprintf("ID generado: %s para disco %d", id, diskSignature))
	return id
}

// isPartitionAlreadyMounted verifica si una partición ya está montada
func isPartitionAlreadyMounted(path, name string) bool {
	for _, partition := range mountSystem.mountedPartitions {
		if partition.Path == path && partition.Name == name {
			return true
		}
	}
	return false
}

// updatePartitionInDisk actualiza la partición en el disco con información de montaje
func updatePartitionInDisk(path, name string, mountedPartition *MountedPartition, mbr *estructuras.MBR) error {
	// Encontrar la partición en el MBR
	for i := range mbr.MbrParticiones {
		if mbr.MbrParticiones[i].IsActive() && mbr.MbrParticiones[i].GetName() == name {
			// Asignar correlativo y ID
			mbr.MbrParticiones[i].PartCorrelativo = int64(len(mountSystem.mountedPartitions) + 1)
			mbr.MbrParticiones[i].SetID(mountedPartition.ID)

			// Actualizar la estructura montada
			mountedPartition.Correlative = mbr.MbrParticiones[i].PartCorrelativo

			// Escribir MBR actualizado
			mbrData, err := estructuras.SerializeMBR(mbr)
			if err != nil {
				return fmt.Errorf("error al serializar MBR: %v", err)
			}

			return estructuras.WriteToDisk(path, mbrData, 0)
		}
	}

	return fmt.Errorf("no se pudo actualizar la partición en el disco")
}

// unmountPartitionInDisk actualiza la partición en el disco para desmontar
func unmountPartitionInDisk(mountedPartition *MountedPartition) error {
	// Leer el MBR
	mbr, err := estructuras.ReadMBR(mountedPartition.Path)
	if err != nil {
		return fmt.Errorf("error al leer MBR: %v", err)
	}

	// Encontrar y actualizar la partición
	for i := range mbr.MbrParticiones {
		if mbr.MbrParticiones[i].IsActive() &&
			mbr.MbrParticiones[i].GetName() == mountedPartition.Name {

			// Limpiar información de montaje
			mbr.MbrParticiones[i].PartCorrelativo = -1
			mbr.MbrParticiones[i].SetID("")

			// Escribir MBR actualizado
			mbrData, err := estructuras.SerializeMBR(mbr)
			if err != nil {
				return fmt.Errorf("error al serializar MBR: %v", err)
			}

			return estructuras.WriteToDisk(mountedPartition.Path, mbrData, 0)
		}
	}

	return fmt.Errorf("no se pudo actualizar la partición en el disco")
}

// getCurrentTimestamp retorna el timestamp actual
func getCurrentTimestamp() int64 {
	return 1672531200 // Timestamp ejemplo, usar time.Now().Unix() en implementación real
}

// GetMountSystemStats retorna estadísticas del sistema de montaje
func GetMountSystemStats() map[string]interface{} {
	mountSystem.mutex.RLock()
	defer mountSystem.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_mounted": len(mountSystem.mountedPartitions),
		"unique_disks":  len(mountSystem.diskPartitionCount),
		"next_letter":   string(mountSystem.nextLetter),
		"carnet_suffix": mountSystem.carnetSuffix,
	}

	// Información por disco
	diskStats := make(map[string]int)
	for diskKey, count := range mountSystem.diskPartitionCount {
		diskStats[diskKey] = count
	}
	stats["partitions_per_disk"] = diskStats

	// Lista de IDs montados
	var mountedIDs []string
	for id := range mountSystem.mountedPartitions {
		mountedIDs = append(mountedIDs, id)
	}
	sort.Strings(mountedIDs)
	stats["mounted_ids"] = mountedIDs

	return stats
}

// ClearMountSystem limpia el sistema de montaje (útil para testing)
func ClearMountSystem() {
	mountSystem.mutex.Lock()
	defer mountSystem.mutex.Unlock()

	mountSystem.mountedPartitions = make(map[string]*MountedPartition)
	mountSystem.diskPartitionCount = make(map[string]int)
	mountSystem.nextLetter = 'A'

	utils.LogInfo("MOUNT", "Sistema de montaje limpiado")
}

// ShowMountedPartitions muestra todas las particiones montadas (comando MOUNTED)
func ShowMountedPartitions() string {
	partitions := GetMountedPartitions()

	if len(partitions) == 0 {
		return "No hay particiones montadas en el sistema"
	}

	var result strings.Builder
	result.WriteString("=== PARTICIONES MONTADAS ===\n")
	result.WriteString(fmt.Sprintf("Total: %d particiones\n\n", len(partitions)))

	for _, partition := range partitions {
		result.WriteString(fmt.Sprintf("ID: %s\n", partition.ID))
		result.WriteString(fmt.Sprintf("  Nombre: %s\n", partition.Name))
		result.WriteString(fmt.Sprintf("  Disco: %s\n", partition.Path))
		result.WriteString(fmt.Sprintf("  Tipo: %s\n", partition.Type))
		result.WriteString(fmt.Sprintf("  Tamaño: %d bytes (%.2f MB)\n",
			partition.Size, float64(partition.Size)/(1024*1024)))
		result.WriteString(fmt.Sprintf("  Correlativo: %d\n", partition.Correlative))
		result.WriteString(fmt.Sprintf("  Montado: %s\n", partition.MountTime))
		result.WriteString("\n")
	}

	return result.String()
}
