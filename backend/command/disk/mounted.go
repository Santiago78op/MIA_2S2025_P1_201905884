package disk

/*
 * MOUNTED - Este comando mostrará todas las particiones montadas en memoria.
 */

import (
	utils "backend/Utils"
	"fmt"
	"sort"
	"strings"
)

/*
| COMANDO   | DESCRIPCIÓN                                                                                                                                                                                                                                                                                                                                                      |
|-----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| mounted   | Muestra todas las particiones montadas actualmente en el sistema. No requiere parámetros.                                                                                                                                                                                                                                                                        |
*/

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

// ShowMountedPartitions muestra todas las particiones montadas (comando MOUNTED)
func ShowMountedPartitions() string {
	utils.LogInfo("MOUNTED", "Mostrando particiones montadas en el sistema")

	partitions := GetMountedPartitions()

	if len(partitions) == 0 {
		utils.LogInfo("MOUNTED", "No hay particiones montadas en el sistema")
		return "No hay particiones montadas en el sistema"
	}

	var result strings.Builder
	result.WriteString("=== PARTICIONES MONTADAS ===\n")
	result.WriteString(fmt.Sprintf("Total: %d particiones\n\n", len(partitions)))

	// Estadísticas generales
	stats := GetMountSystemStats()
	result.WriteString("Estadísticas del sistema:\n")
	result.WriteString(fmt.Sprintf("  • Discos únicos: %d\n", stats["unique_disks"]))
	result.WriteString(fmt.Sprintf("  • Próxima letra: %s\n", stats["next_letter"]))
	result.WriteString(fmt.Sprintf("  • Sufijo carnet: %s\n\n", stats["carnet_suffix"]))

	// Agrupar por disco para mejor presentación
	diskGroups := groupPartitionsByDisk(partitions)

	for diskPath, diskPartitions := range diskGroups {
		result.WriteString(fmt.Sprintf("📁 DISCO: %s\n", diskPath))
		result.WriteString(strings.Repeat("-", 50) + "\n")

		for _, partition := range diskPartitions {
			result.WriteString(fmt.Sprintf("  🔹 ID: %s\n", partition.ID))
			result.WriteString(fmt.Sprintf("     Nombre: %s\n", partition.Name))
			result.WriteString(fmt.Sprintf("     Tipo: %s\n", partition.Type))
			result.WriteString(fmt.Sprintf("     Tamaño: %s\n", formatSize(partition.Size)))
			result.WriteString(fmt.Sprintf("     Correlativo: %d\n", partition.Correlative))
			result.WriteString(fmt.Sprintf("     Montado: %s\n", formatTimestamp(partition.MountTime)))

			if partition.PartitionIndex >= 0 {
				result.WriteString(fmt.Sprintf("     Índice MBR: %d\n", partition.PartitionIndex))
			}

			if partition.EBRPosition > 0 {
				result.WriteString(fmt.Sprintf("     Posición EBR: %d\n", partition.EBRPosition))
			}

			result.WriteString("\n")
		}
		result.WriteString("\n")
	}

	utils.LogSuccess("MOUNTED", fmt.Sprintf("Se mostraron %d particiones montadas", len(partitions)))
	utils.LogInfo("MOUNTED", result.String())
	return result.String()
}

// ShowMountedPartitionsTable muestra las particiones en formato tabla
func ShowMountedPartitionsTable() string {
	utils.LogInfo("MOUNTED", "Generando tabla de particiones montadas")

	partitions := GetMountedPartitions()

	if len(partitions) == 0 {
		return "No hay particiones montadas en el sistema"
	}

	var result strings.Builder

	// Encabezado de la tabla
	result.WriteString("┌──────────┬────────────────┬──────────────┬──────────────┬────────────────┬─────────────┐\n")
	result.WriteString("│    ID    │     NOMBRE     │     TIPO     │   TAMAÑO     │     DISCO      │ CORRELATIVO │\n")
	result.WriteString("├──────────┼────────────────┼──────────────┼──────────────┼────────────────┼─────────────┤\n")

	// Filas de datos
	for _, partition := range partitions {
		id := truncateString(partition.ID, 8)
		name := truncateString(partition.Name, 14)
		pType := truncateString(partition.Type, 12)
		size := truncateString(formatSize(partition.Size), 12)
		disk := truncateString(getFileName(partition.Path), 14)
		correlative := fmt.Sprintf("%d", partition.Correlative)

		result.WriteString(fmt.Sprintf("│ %-8s │ %-14s │ %-12s │ %-12s │ %-14s │ %-11s │\n",
			id, name, pType, size, disk, correlative))
	}

	result.WriteString("└──────────┴────────────────┴──────────────┴──────────────┴────────────────┴─────────────┘\n")
	result.WriteString(fmt.Sprintf("\nTotal: %d particiones montadas\n", len(partitions)))

	return result.String()
}

// ShowMountedPartitionsJSON retorna las particiones en formato JSON para API
func ShowMountedPartitionsJSON() map[string]interface{} {
	partitions := GetMountedPartitions()
	stats := GetMountSystemStats()

	return map[string]interface{}{
		"partitions": partitions,
		"total":      len(partitions),
		"stats":      stats,
	}
}

// groupPartitionsByDisk agrupa las particiones por disco
func groupPartitionsByDisk(partitions []*MountedPartition) map[string][]*MountedPartition {
	groups := make(map[string][]*MountedPartition)

	for _, partition := range partitions {
		groups[partition.Path] = append(groups[partition.Path], partition)
	}

	// Ordenar particiones dentro de cada grupo
	for diskPath := range groups {
		sort.Slice(groups[diskPath], func(i, j int) bool {
			return groups[diskPath][i].ID < groups[diskPath][j].ID
		})
	}

	return groups
}

// formatSize formatea el tamaño en bytes a una representación legible
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

// formatTimestamp convierte un timestamp a formato legible
func formatTimestamp(timestamp string) string {
	// En una implementación real, convertir el timestamp Unix a fecha legible
	return fmt.Sprintf("Unix: %s", timestamp)
}

// truncateString trunca una cadena a la longitud especificada
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// getFileName extrae el nombre del archivo de una ruta
func getFileName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return path
}

// ShowMountedPartitionsDetailed muestra información detallada de cada partición
func ShowMountedPartitionsDetailed() string {
	utils.LogInfo("MOUNTED", "Generando reporte detallado de particiones montadas")

	partitions := GetMountedPartitions()

	if len(partitions) == 0 {
		return "No hay particiones montadas en el sistema"
	}

	var result strings.Builder
	result.WriteString("=== REPORTE DETALLADO DE PARTICIONES MONTADAS ===\n\n")

	for i, partition := range partitions {
		result.WriteString(fmt.Sprintf("PARTICIÓN #%d\n", i+1))
		result.WriteString(strings.Repeat("=", 40) + "\n")
		result.WriteString(fmt.Sprintf("ID de Montaje........: %s\n", partition.ID))
		result.WriteString(fmt.Sprintf("Nombre...............: %s\n", partition.Name))
		result.WriteString(fmt.Sprintf("Ruta del Disco.......: %s\n", partition.Path))
		result.WriteString(fmt.Sprintf("Tipo de Partición....: %s\n", partition.Type))
		result.WriteString(fmt.Sprintf("Tamaño...............: %s (%d bytes)\n", formatSize(partition.Size), partition.Size))
		result.WriteString(fmt.Sprintf("Correlativo..........: %d\n", partition.Correlative))
		result.WriteString(fmt.Sprintf("Firma del Disco......: %d\n", partition.DiskSignature))
		result.WriteString(fmt.Sprintf("Fecha de Montaje.....: %s\n", formatTimestamp(partition.MountTime)))

		if partition.PartitionIndex >= 0 {
			result.WriteString(fmt.Sprintf("Índice en MBR........: %d\n", partition.PartitionIndex))
		} else {
			result.WriteString("Índice en MBR........: N/A (Partición Lógica)\n")
		}

		if partition.EBRPosition > 0 {
			result.WriteString(fmt.Sprintf("Posición EBR.........: %d\n", partition.EBRPosition))
		} else {
			result.WriteString("Posición EBR.........: N/A (Partición Primaria)\n")
		}

		result.WriteString("\n")
	}

	// Estadísticas del sistema
	stats := GetMountSystemStats()
	result.WriteString("=== ESTADÍSTICAS DEL SISTEMA ===\n")
	result.WriteString(fmt.Sprintf("Total de particiones montadas: %d\n", len(partitions)))
	result.WriteString(fmt.Sprintf("Discos únicos................ : %d\n", stats["unique_disks"]))
	result.WriteString(fmt.Sprintf("Próxima letra disponible.... : %s\n", stats["next_letter"]))
	result.WriteString(fmt.Sprintf("Sufijo del carnet........... : %s\n", stats["carnet_suffix"]))

	if diskStats, ok := stats["partitions_per_disk"].(map[string]int); ok {
		result.WriteString("\nParticiones por disco:\n")
		for diskSig, count := range diskStats {
			result.WriteString(fmt.Sprintf("  Disco %s: %d particiones\n", diskSig, count))
		}
	}

	return result.String()
}
