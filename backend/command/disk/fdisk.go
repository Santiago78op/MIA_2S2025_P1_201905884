package disk

/*
 * FDISK - Comando para administración de particiones en el archivo que representa al disco duro virtual.
 * Permite a los usuarios crear, eliminar o modificar particiones dentro del archivo de disco duro.
 *
 * Funcionalidades:
 * - Crear particiones primarias, extendidas y lógicas
 * - Aplicar algoritmos de ajuste (Best Fit, First Fit, Worst Fit)
 * - Validar restricciones de teoría de particiones
 * - Manejar particiones lógicas dentro de extendidas
 */

import (
	utils "backend/Utils"
	estructuras "backend/struct"
	"fmt"
	"strings"
)

/*
| PARÁMETRO | CATEGORÍA    | DESCRIPCIÓN                                                                                                                                                                                                                                                                                                                                                      |
|-----------|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| -size     | Obligatorio  | Número que indica el tamaño de la partición a crear. Debe ser positivo y mayor que cero.                                                                                                                                                                                                                                                                        |
| -unit     | Opcional     | Letra que indica las unidades del parámetro size. Valores: B (bytes), K (Kilobytes), M (Megabytes). Default: Kilobytes.                                                                                                                                                                                                                                        |
| -path     | Obligatorio  | Ruta del disco donde se creará la partición. El archivo debe existir.                                                                                                                                                                                                                                                                                           |
| -type     | Opcional     | Tipo de partición a crear. Valores: P (Primaria), E (Extendida), L (Lógica). Default: Primaria.                                                                                                                                                                                                                                                                |
| -fit      | Opcional     | Algoritmo de ajuste para asignar espacio. Valores: BF (Best Fit), FF (First Fit), WF (Worst Fit). Default: Worst Fit.                                                                                                                                                                                                                                          |
| -name     | Obligatorio  | Nombre de la partición. No debe repetirse dentro de las particiones de cada disco.                                                                                                                                                                                                                                                                              |
*/

// FdiskAction define las acciones posibles con FDISK
type FdiskAction int

const (
	ActionCreate FdiskAction = iota
	ActionDelete
	ActionAdd
)

// FdiskParams contiene los parámetros del comando FDISK
type FdiskParams struct {
	Size   int64
	Unit   string
	Path   string
	Type   string
	Fit    string
	Name   string
	Action FdiskAction
}

// FdiskResult contiene el resultado de la operación FDISK
type FdiskResult struct {
	Success       bool
	Message       string
	PartitionInfo map[string]interface{}
}

// Fdisk ejecuta el comando fdisk con los parámetros especificados
func Fdisk(size int64, unit, path, partType, fit, name string) error {
	utils.LogInfo("FDISK", fmt.Sprintf("Iniciando administración de particiones: path=%s, name=%s, type=%s, size=%d, unit=%s, fit=%s",
		path, name, partType, size, unit, fit))

	// Validar parámetros obligatorios
	if err := validateFdiskParams(size, path, name); err != nil {
		return err
	}

	// Validar que el disco existe
	if err := validateDiskExists(path); err != nil {
		return err
	}

	// Normalizar parámetros opcionales
	unit = normalizeUnit(unit)
	partType = normalizePartitionType(partType)
	fit = normalizeFit(fit)

	// Validar parámetros normalizados
	if err := validateNormalizedParams(unit, partType, fit); err != nil {
		return err
	}

	// Calcular tamaño en bytes
	sizeInBytes, err := calculateSizeInBytes(size, unit)
	if err != nil {
		return err
	}

	// Ejecutar la creación de partición según el tipo
	switch partType {
	case "P":
		return createPrimaryPartition(path, name, sizeInBytes, fit)
	case "E":
		return createExtendedPartition(path, name, sizeInBytes, fit)
	case "L":
		return createLogicalPartition(path, name, sizeInBytes, fit)
	default:
		return fmt.Errorf("tipo de partición no válido: %s", partType)
	}
}

// validateFdiskParams valida los parámetros obligatorios
func validateFdiskParams(size int64, path, name string) error {
	if size <= 0 {
		utils.LogError("FDISK", "El tamaño de la partición debe ser positivo y mayor que cero")
		return fmt.Errorf("el tamaño de la partición debe ser positivo y mayor que cero")
	}

	if path == "" {
		utils.LogError("FDISK", "El parámetro -path es obligatorio")
		return fmt.Errorf("el parámetro -path es obligatorio")
	}

	if name == "" {
		utils.LogError("FDISK", "El parámetro -name es obligatorio")
		return fmt.Errorf("el parámetro -name es obligatorio")
	}

	// Validar longitud del nombre
	if len(name) > 16 {
		utils.LogError("FDISK", "El nombre de la partición no puede exceder 16 caracteres")
		return fmt.Errorf("el nombre de la partición no puede exceder 16 caracteres")
	}

	return nil
}

// validateDiskExists verifica que el archivo del disco existe
func validateDiskExists(path string) error {
	// Intentar validar la integridad del disco
	if err := estructuras.ValidateDiskIntegrity(path); err != nil {
		utils.LogError("FDISK", fmt.Sprintf("Error de integridad del disco: %v", err))
		return fmt.Errorf("error de integridad del disco: %v", err)
	}

	utils.LogInfo("FDISK", "Disco validado correctamente")
	return nil
}

// normalizeUnit normaliza la unidad especificada
func normalizeUnit(unit string) string {
	unit = strings.ToUpper(strings.TrimSpace(unit))
	if unit == "" {
		return "K" // Default: Kilobytes
	}
	return unit
}

// normalizePartitionType normaliza el tipo de partición
func normalizePartitionType(partType string) string {
	partType = strings.ToUpper(strings.TrimSpace(partType))
	if partType == "" {
		return "P" // Default: Primaria
	}
	return partType
}

// normalizeFit normaliza el algoritmo de ajuste
func normalizeFit(fit string) string {
	fit = strings.ToUpper(strings.TrimSpace(fit))
	if fit == "" {
		return "WF" // Default: Worst Fit
	}
	return fit
}

// validateNormalizedParams valida los parámetros después de normalizar
func validateNormalizedParams(unit, partType, fit string) error {
	// Validar unidad
	validUnits := map[string]bool{"B": true, "K": true, "M": true}
	if !validUnits[unit] {
		utils.LogError("FDISK", fmt.Sprintf("Unidad no válida '%s', use B, K o M", unit))
		return fmt.Errorf("unidad no válida '%s', use B, K o M", unit)
	}

	// Validar tipo de partición
	validTypes := map[string]bool{"P": true, "E": true, "L": true}
	if !validTypes[partType] {
		utils.LogError("FDISK", fmt.Sprintf("Tipo de partición no válido '%s', use P, E o L", partType))
		return fmt.Errorf("tipo de partición no válido '%s', use P, E o L", partType)
	}

	// Validar algoritmo de ajuste
	validFits := map[string]bool{"BF": true, "FF": true, "WF": true}
	if !validFits[fit] {
		utils.LogError("FDISK", fmt.Sprintf("Algoritmo de ajuste no válido '%s', use BF, FF o WF", fit))
		return fmt.Errorf("algoritmo de ajuste no válido '%s', use BF, FF o WF", fit)
	}

	return nil
}

// calculateSizeInBytes calcula el tamaño en bytes según la unidad
func calculateSizeInBytes(size int64, unit string) (int64, error) {
	var multiplier int64
	var unitName string

	switch unit {
	case "B":
		multiplier = 1
		unitName = "bytes"
	case "K":
		multiplier = 1024
		unitName = "Kilobytes"
	case "M":
		multiplier = 1024 * 1024
		unitName = "Megabytes"
	default:
		return 0, fmt.Errorf("unidad no reconocida: %s", unit)
	}

	sizeInBytes := size * multiplier
	utils.LogInfo("FDISK", fmt.Sprintf("Tamaño calculado: %d %s = %d bytes", size, unitName, sizeInBytes))

	return sizeInBytes, nil
}

// createPrimaryPartition crea una partición primaria
func createPrimaryPartition(path, name string, sizeInBytes int64, fit string) error {
	utils.LogInfo("FDISK", fmt.Sprintf("Creando partición primaria: %s", name))

	// Leer el MBR
	mbr, err := estructuras.ReadMBR(path)
	if err != nil {
		return fmt.Errorf("error al leer MBR: %v", err)
	}

	// Verificar restricciones de particiones primarias y extendidas
	activeCount := mbr.CountActivePartitions()
	if activeCount >= 4 {
		utils.LogError("FDISK", "No se pueden crear más de 4 particiones primarias y extendidas")
		return fmt.Errorf("no se pueden crear más de 4 particiones primarias y extendidas (ya existen %d)", activeCount)
	}

	// Verificar que el nombre no se repita
	if mbr.GetParticionByName(name) != nil {
		utils.LogError("FDISK", fmt.Sprintf("Ya existe una partición con el nombre '%s'", name))
		return fmt.Errorf("ya existe una partición con el nombre '%s'", name)
	}

	// Encontrar una partición libre
	partition := mbr.GetParticionLibre()
	if partition == nil {
		utils.LogError("FDISK", "No hay espacios de partición disponibles")
		return fmt.Errorf("no hay espacios de partición disponibles")
	}

	// Obtener fit byte del MBR (usar el del disco si no se especifica diferente)
	fitByte := estructuras.ValidateFit(fit)
	if fitByte == 0 {
		fitByte = mbr.MbrFit // Usar el fit del disco
	}

	// Temporalmente cambiar el fit del MBR para usar el algoritmo especificado
	originalFit := mbr.MbrFit
	mbr.MbrFit = fitByte

	// Encontrar la mejor posición para la partición
	startPosition, err := mbr.FindBestFitPosition(sizeInBytes)
	if err != nil {
		utils.LogError("FDISK", fmt.Sprintf("No se pudo encontrar espacio para la partición: %v", err))
		return fmt.Errorf("no se pudo encontrar espacio para la partición: %v", err)
	}

	// Restaurar el fit original del MBR
	mbr.MbrFit = originalFit

	// Configurar la partición
	partition.PartStatus = estructuras.StatusActiva
	partition.PartType = estructuras.PartitionTypePrimaria
	partition.PartFit = fitByte
	partition.PartStart = startPosition
	partition.PartSize = sizeInBytes
	partition.SetName(name)
	partition.PartCorrelativo = -1 // No montada inicialmente

	// Validar la partición creada
	if err := partition.ValidatePartition(); err != nil {
		return fmt.Errorf("la partición creada no es válida: %v", err)
	}

	// Escribir el MBR actualizado al disco
	if err := writeUpdatedMBR(path, mbr); err != nil {
		return fmt.Errorf("error al escribir MBR actualizado: %v", err)
	}

	utils.LogSuccess("FDISK", "Partición primaria creada exitosamente:")
	logPartitionInfo(partition, "Primaria")

	return nil
}

// createExtendedPartition crea una partición extendida
func createExtendedPartition(path, name string, sizeInBytes int64, fit string) error {
	utils.LogInfo("FDISK", fmt.Sprintf("Creando partición extendida: %s", name))

	// Leer el MBR
	mbr, err := estructuras.ReadMBR(path)
	if err != nil {
		return fmt.Errorf("error al leer MBR: %v", err)
	}

	// Verificar que no exista ya una partición extendida
	if mbr.HasExtendedPartition() {
		utils.LogError("FDISK", "Solo puede haber una partición extendida por disco")
		return fmt.Errorf("solo puede haber una partición extendida por disco")
	}

	// Verificar restricciones de particiones primarias y extendidas
	activeCount := mbr.CountActivePartitions()
	if activeCount >= 4 {
		utils.LogError("FDISK", "No se pueden crear más de 4 particiones primarias y extendidas")
		return fmt.Errorf("no se pueden crear más de 4 particiones primarias y extendidas (ya existen %d)", activeCount)
	}

	// Verificar que el nombre no se repita
	if mbr.GetParticionByName(name) != nil {
		utils.LogError("FDISK", fmt.Sprintf("Ya existe una partición con el nombre '%s'", name))
		return fmt.Errorf("ya existe una partición con el nombre '%s'", name)
	}

	// Encontrar una partición libre
	partition := mbr.GetParticionLibre()
	if partition == nil {
		utils.LogError("FDISK", "No hay espacios de partición disponibles")
		return fmt.Errorf("no hay espacios de partición disponibles")
	}

	// Obtener fit byte
	fitByte := estructuras.ValidateFit(fit)
	if fitByte == 0 {
		fitByte = mbr.MbrFit
	}

	// Temporalmente cambiar el fit del MBR
	originalFit := mbr.MbrFit
	mbr.MbrFit = fitByte

	// Encontrar la mejor posición para la partición
	startPosition, err := mbr.FindBestFitPosition(sizeInBytes)
	if err != nil {
		utils.LogError("FDISK", fmt.Sprintf("No se pudo encontrar espacio para la partición extendida: %v", err))
		return fmt.Errorf("no se pudo encontrar espacio para la partición extendida: %v", err)
	}

	// Restaurar el fit original del MBR
	mbr.MbrFit = originalFit

	// Configurar la partición extendida
	partition.PartStatus = estructuras.StatusActiva
	partition.PartType = estructuras.PartitionTypeExtendida
	partition.PartFit = fitByte
	partition.PartStart = startPosition
	partition.PartSize = sizeInBytes
	partition.SetName(name)
	partition.PartCorrelativo = -1

	// Validar la partición creada
	if err := partition.ValidatePartition(); err != nil {
		return fmt.Errorf("la partición extendida creada no es válida: %v", err)
	}

	// Escribir el MBR actualizado al disco
	if err := writeUpdatedMBR(path, mbr); err != nil {
		return fmt.Errorf("error al escribir MBR actualizado: %v", err)
	}

	// Crear el primer EBR en la partición extendida (inicialmente vacío)
	if err := createInitialEBR(path, startPosition); err != nil {
		utils.LogWarning("FDISK", fmt.Sprintf("Advertencia al crear EBR inicial: %v", err))
	}

	utils.LogSuccess("FDISK", "Partición extendida creada exitosamente:")
	logPartitionInfo(partition, "Extendida")

	return nil
}

// createLogicalPartition crea una partición lógica dentro de la extendida
func createLogicalPartition(path, name string, sizeInBytes int64, fit string) error {
	utils.LogInfo("FDISK", fmt.Sprintf("Creando partición lógica: %s", name))

	// Leer el MBR
	mbr, err := estructuras.ReadMBR(path)
	if err != nil {
		return fmt.Errorf("error al leer MBR: %v", err)
	}

	// Buscar la partición extendida
	extendedPartition, extendedIndex, err := estructuras.GetExtendedPartition(path)
	if err != nil {
		utils.LogError("FDISK", "No se puede crear una partición lógica sin una partición extendida")
		return fmt.Errorf("no se puede crear una partición lógica sin una partición extendida")
	}

	utils.LogInfo("FDISK", fmt.Sprintf("Partición extendida encontrada en índice %d: %s", extendedIndex, extendedPartition.GetName()))

	// Verificar que el nombre no se repita (tanto en MBR como en EBRs existentes)
	if mbr.GetParticionByName(name) != nil {
		utils.LogError("FDISK", fmt.Sprintf("Ya existe una partición primaria/extendida con el nombre '%s'", name))
		return fmt.Errorf("ya existe una partición primaria/extendida con el nombre '%s'", name)
	}

	// Verificar nombres en particiones lógicas existentes
	if err := checkLogicalPartitionNameExists(path, extendedPartition.PartStart, name); err != nil {
		return err
	}

	// Encontrar espacio dentro de la partición extendida para la partición lógica
	logicalStart, err := findLogicalPartitionSpace(path, extendedPartition, sizeInBytes, fit)
	if err != nil {
		return fmt.Errorf("no se pudo encontrar espacio en la partición extendida: %v", err)
	}

	// Crear el EBR para la nueva partición lógica
	fitByte := estructuras.ValidateFit(fit)
	if fitByte == 0 {
		fitByte = extendedPartition.PartFit
	}

	// Crear el EBR
	nextEBRPosition := findNextEBRPosition(path, extendedPartition.PartStart, logicalStart+sizeInBytes+int64(estructuras.EBR_SIZE))

	newEBR := estructuras.NewEBR(fitByte, logicalStart+int64(estructuras.EBR_SIZE), sizeInBytes, name, nextEBRPosition)

	// Escribir el EBR al disco
	if err := estructuras.WriteEBR(path, newEBR, logicalStart); err != nil {
		return fmt.Errorf("error al escribir EBR: %v", err)
	}

	// Actualizar la cadena de EBRs si es necesario
	if err := updateEBRChain(path, extendedPartition.PartStart, logicalStart); err != nil {
		utils.LogWarning("FDISK", fmt.Sprintf("Advertencia al actualizar cadena de EBRs: %v", err))
	}

	utils.LogSuccess("FDISK", "Partición lógica creada exitosamente:")
	utils.LogSuccess("FDISK", fmt.Sprintf("  → Nombre: %s", name))
	utils.LogSuccess("FDISK", "  → Tipo: Lógica")
	utils.LogSuccess("FDISK", fmt.Sprintf("  → Tamaño: %d bytes", sizeInBytes))
	utils.LogSuccess("FDISK", fmt.Sprintf("  → Inicio: %d", logicalStart))
	utils.LogSuccess("FDISK", fmt.Sprintf("  → Ajuste: %s", fit))

	return nil
}

// Helper functions

// writeUpdatedMBR escribe el MBR actualizado al disco
func writeUpdatedMBR(path string, mbr *estructuras.MBR) error {
	mbrData, err := estructuras.SerializeMBR(mbr)
	if err != nil {
		return fmt.Errorf("error al serializar MBR: %v", err)
	}

	err = estructuras.WriteToDisk(path, mbrData, 0)
	if err != nil {
		return fmt.Errorf("error al escribir MBR: %v", err)
	}

	return nil
}

// createInitialEBR crea el primer EBR vacío en una partición extendida
func createInitialEBR(path string, extendedStart int64) error {
	emptyEBR := estructuras.NewEmptyEBR()
	return estructuras.WriteEBR(path, emptyEBR, extendedStart)
}

// checkLogicalPartitionNameExists verifica si ya existe una partición lógica con el nombre dado
func checkLogicalPartitionNameExists(path string, extendedStart int64, name string) error {
	_, _, err := estructuras.FindEBRByName(path, extendedStart, name)
	if err == nil {
		utils.LogError("FDISK", fmt.Sprintf("Ya existe una partición lógica con el nombre '%s'", name))
		return fmt.Errorf("ya existe una partición lógica con el nombre '%s'", name)
	}
	return nil // No existe, está bien
}

// findLogicalPartitionSpace encuentra espacio disponible dentro de la partición extendida
func findLogicalPartitionSpace(path string, extendedPartition *estructuras.Partition, sizeInBytes int64, fit string) (int64, error) {
	// Obtener todos los EBRs existentes
	ebrs, err := estructuras.ReadAllEBRs(path, extendedPartition.PartStart)
	if err != nil {
		return 0, fmt.Errorf("error al leer EBRs existentes: %v", err)
	}

	// Calcular espacios ocupados dentro de la partición extendida
	occupiedSpaces := []estructuras.FreeSpace{
		{Start: extendedPartition.PartStart, Size: int64(estructuras.EBR_SIZE)}, // Espacio del primer EBR
	}

	for _, ebr := range ebrs {
		if !ebr.IsEmpty() {
			// Cada partición lógica ocupa: EBR + datos
			totalSize := int64(estructuras.EBR_SIZE) + ebr.PartSize
			occupiedSpaces = append(occupiedSpaces, estructuras.FreeSpace{
				Start: ebr.PartStart - int64(estructuras.EBR_SIZE),
				Size:  totalSize,
			})
		}
	}

	// Encontrar espacios libres
	freeSpaces := findFreeSpacesInExtended(extendedPartition, occupiedSpaces)

	// Aplicar algoritmo de ajuste
	neededSize := sizeInBytes + int64(estructuras.EBR_SIZE) // EBR + datos

	fitByte := estructuras.ValidateFit(fit)
	switch fitByte {
	case estructuras.PartitionFitFirst:
		for _, space := range freeSpaces {
			if space.Size >= neededSize {
				return space.Start, nil
			}
		}
	case estructuras.PartitionFitBest:
		bestSpace := estructuras.FreeSpace{Size: extendedPartition.PartSize + 1}
		found := false
		for _, space := range freeSpaces {
			if space.Size >= neededSize && space.Size < bestSpace.Size {
				bestSpace = space
				found = true
			}
		}
		if found {
			return bestSpace.Start, nil
		}
	case estructuras.PartitionFitWorst:
		worstSpace := estructuras.FreeSpace{Size: -1}
		for _, space := range freeSpaces {
			if space.Size >= neededSize && space.Size > worstSpace.Size {
				worstSpace = space
			}
		}
		if worstSpace.Size >= neededSize {
			return worstSpace.Start, nil
		}
	}

	return 0, fmt.Errorf("no hay espacio suficiente en la partición extendida")
}

// findFreeSpacesInExtended encuentra espacios libres dentro de la partición extendida
func findFreeSpacesInExtended(extendedPartition *estructuras.Partition, occupiedSpaces []estructuras.FreeSpace) []estructuras.FreeSpace {
	var freeSpaces []estructuras.FreeSpace

	// Ordenar espacios ocupados por posición
	for i := 0; i < len(occupiedSpaces)-1; i++ {
		for j := i + 1; j < len(occupiedSpaces); j++ {
			if occupiedSpaces[i].Start > occupiedSpaces[j].Start {
				occupiedSpaces[i], occupiedSpaces[j] = occupiedSpaces[j], occupiedSpaces[i]
			}
		}
	}

	// Encontrar espacios libres entre ocupados
	currentPos := extendedPartition.PartStart
	extendedEnd := extendedPartition.PartStart + extendedPartition.PartSize

	for _, occupied := range occupiedSpaces {
		if currentPos < occupied.Start {
			freeSpaces = append(freeSpaces, estructuras.FreeSpace{
				Start: currentPos,
				Size:  occupied.Start - currentPos,
			})
		}
		currentPos = occupied.Start + occupied.Size
	}

	// Espacio libre al final
	if currentPos < extendedEnd {
		freeSpaces = append(freeSpaces, estructuras.FreeSpace{
			Start: currentPos,
			Size:  extendedEnd - currentPos,
		})
	}

	return freeSpaces
}

// findNextEBRPosition encuentra la posición del siguiente EBR en la cadena
func findNextEBRPosition(path string, extendedStart, afterPosition int64) int64 {
	// Buscar el siguiente EBR después de la posición especificada
	currentPos := extendedStart

	for currentPos != -1 {
		ebr, err := estructuras.ReadEBR(path, currentPos)
		if err != nil {
			break
		}

		if ebr.PartStart > afterPosition {
			return currentPos
		}

		currentPos = ebr.PartNext
	}

	return -1 // No hay siguiente EBR
}

// updateEBRChain actualiza la cadena de EBRs para mantener la consistencia
func updateEBRChain(path string, extendedStart, newEBRPosition int64) error {
	// Encontrar el EBR anterior que debe apuntar al nuevo EBR
	currentPos := extendedStart
	var previousPos int64 = -1

	for currentPos != -1 && currentPos < newEBRPosition {
		ebr, err := estructuras.ReadEBR(path, currentPos)
		if err != nil {
			return err
		}

		if ebr.PartNext == -1 || ebr.PartNext > newEBRPosition {
			// Este EBR debe apuntar al nuevo EBR
			previousPos = currentPos
			break
		}

		previousPos = currentPos
		currentPos = ebr.PartNext
	}

	// Actualizar el EBR anterior para que apunte al nuevo EBR
	if previousPos != -1 {
		prevEBR, err := estructuras.ReadEBR(path, previousPos)
		if err != nil {
			return err
		}

		oldNext := prevEBR.PartNext
		prevEBR.PartNext = newEBRPosition

		if err := estructuras.WriteEBR(path, prevEBR, previousPos); err != nil {
			return err
		}

		// Hacer que el nuevo EBR apunte al siguiente en la cadena original
		newEBR, err := estructuras.ReadEBR(path, newEBRPosition)
		if err != nil {
			return err
		}

		newEBR.PartNext = oldNext
		return estructuras.WriteEBR(path, newEBR, newEBRPosition)
	}

	return nil
}

// logPartitionInfo registra información de la partición creada
func logPartitionInfo(partition *estructuras.Partition, partType string) {
	utils.LogSuccess("FDISK", fmt.Sprintf("  → Nombre: %s", partition.GetName()))
	utils.LogSuccess("FDISK", fmt.Sprintf("  → Tipo: %s", partType))
	utils.LogSuccess("FDISK", fmt.Sprintf("  → Tamaño: %d bytes", partition.PartSize))
	utils.LogSuccess("FDISK", fmt.Sprintf("  → Inicio: %d", partition.PartStart))
	utils.LogSuccess("FDISK", fmt.Sprintf("  → Fin: %d", partition.GetEndPosition()))
	utils.LogSuccess("FDISK", fmt.Sprintf("  → Ajuste: %s", partition.GetFitString()))
}
