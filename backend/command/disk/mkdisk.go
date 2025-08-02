package disk

/*
 * Este comando creará un archivo binario que simulará un disco, estos archivos
 * binarios tendrán la extensión .mia y su contenido al inicio será 0 binarios.
 * Deberá ocupar físicamente el tamaño indicado por los parámetros
 */

import (
	utils "backend/Utils"
	action "backend/action"
	estructuras "backend/struct"
	"fmt"
	"strings"
)

/*
| PARÁMETRO | CATEGORÍA  | DESCRIPCIÓN                                                                                                                                                                                                                                                                                                                                                      |
|-----------|------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| -size     | Obligatorio | Recibe un número que indica el tamaño del disco a crear. Debe ser positivo y mayor que cero, de lo contrario se mostrará un error.                                                                                                                                                                                                                              |
| -fit      | Opcional    | Indica el ajuste para crear particiones dentro del disco. Valores posibles: <br>BF: Mejor ajuste (Best Fit)<br>FF: Primer ajuste (First Fit)<br>WF: Peor ajuste (Worst Fit)<br>Si no se especifica, se usa FF. Si se usa otro valor, se muestra un mensaje de error.                                                     |
| -unit     | Opcional    | Recibe una letra que indica las unidades para el parámetro size. Valores posibles:<br>K: Kilobytes (1024 bytes)<br>M: Megabytes (1024 * 1024 bytes)<br>Si no se especifica, se usa Megabytes. Si se usa otro valor, se muestra un mensaje de error.                                                                      |
| -path     | Obligatorio | Ruta donde se creará el archivo que representa el disco duro. Si las carpetas de la ruta no existen, deben crearse.                                                                                                                                                                                                      |
*/

// MkDisk crea un archivo binario que simula un disco duro
func MkDisk(size int64, fit string, unit string, path string) error {
	utils.LogInfo("MkDisk", fmt.Sprintf("Iniciando creación de disco: size=%d, fit=%s, unit=%s, path=%s", size, fit, unit, path))

	// Validar tamaño
	if size <= 0 {
		utils.LogError("MkDisk", "El tamaño del disco debe ser un número positivo mayor que cero")
		return fmt.Errorf("el tamaño del disco debe ser un número positivo mayor que cero")
	}

	// Validar y normalizar fit
	fit = strings.ToUpper(strings.TrimSpace(fit))
	if fit == "" {
		fit = "FF" // Default to First Fit
		utils.LogInfo("MkDisk", "Usando ajuste por defecto: First Fit (FF)")
	}

	validFits := map[string]bool{"BF": true, "FF": true, "WF": true}
	if !validFits[fit] {
		utils.LogError("MkDisk", fmt.Sprintf("Ajuste no válido '%s', use BF, FF o WF", fit))
		return fmt.Errorf("ajuste no válido '%s', use BF, FF o WF", fit)
	}

	// Validar y procesar unidad
	unit = strings.ToUpper(strings.TrimSpace(unit))
	unitMultiplier := int64(1024 * 1024) // Por defecto Megabytes
	unitName := "Megabytes"

	if unit != "" {
		switch unit {
		case "K":
			unitMultiplier = 1024
			unitName = "Kilobytes"
		case "M":
			unitMultiplier = 1024 * 1024
			unitName = "Megabytes"
		default:
			utils.LogError("MkDisk", fmt.Sprintf("Unidad no válida '%s', use K para Kilobytes o M para Megabytes", unit))
			return fmt.Errorf("unidad no válida '%s', use K para Kilobytes o M para Megabytes", unit)
		}
	} else {
		utils.LogInfo("MkDisk", "Usando unidad por defecto: Megabytes")
	}

	// Calcular el tamaño en bytes
	sizeInBytes := size * unitMultiplier
	utils.LogInfo("MkDisk", fmt.Sprintf("Tamaño calculado: %d %s = %d bytes", size, unitName, sizeInBytes))

	// Validar que el path sea obligatorio
	if path == "" {
		utils.LogError("MkDisk", "El parámetro -path es obligatorio")
		return fmt.Errorf("el parámetro -path es obligatorio")
	}

	// Verificar que el path tenga la extensión .mia
	if !strings.HasSuffix(strings.ToLower(path), ".mia") {
		utils.LogWarning("MkDisk", "Se recomienda usar la extensión .mia para archivos de disco")
	}

	// Crear el Disco físico
	utils.LogInfo("MkDisk", "Creando archivo físico del disco...")
	err := action.NewDisk(path, sizeInBytes)
	if err != nil {
		utils.LogError("MkDisk", fmt.Sprintf("Error al crear el archivo del disco: %v", err))
		return fmt.Errorf("error al crear el archivo del disco: %v", err)
	}

	// Crear y escribir el MBR (Master Boot Record) al disco
	utils.LogInfo("MkDisk", "Escribiendo MBR (Master Boot Record)...")
	err = estructuras.WriteMBR(path, sizeInBytes, fit)
	if err != nil {
		utils.LogError("MkDisk", fmt.Sprintf("Error al escribir el MBR: %v", err))
		return fmt.Errorf("error al escribir el MBR: %v", err)
	}

	// Verificar que el MBR se escribió correctamente
	utils.LogInfo("MkDisk", "Verificando integridad del MBR...")
	mbr, err := estructuras.ReadMBR(path)
	if err != nil {
		utils.LogError("MkDisk", fmt.Sprintf("Error al verificar el MBR: %v", err))
		return fmt.Errorf("error al verificar el MBR: %v", err)
	}

	// Validar los datos del MBR
	if mbr.MbrTamanio != sizeInBytes {
		utils.LogError("MkDisk", "El tamaño en el MBR no coincide con el tamaño especificado")
		return fmt.Errorf("error de integridad: el tamaño en el MBR no coincide")
	}

	if !mbr.ValidarFit() {
		utils.LogError("MkDisk", "El tipo de ajuste en el MBR no es válido")
		return fmt.Errorf("error de integridad: tipo de ajuste inválido en el MBR")
	}

	utils.LogSuccess("MkDisk", "Disco creado exitosamente:")
	utils.LogSuccess("MkDisk", fmt.Sprintf("  → Ruta: %s", path))
	utils.LogSuccess("MkDisk", fmt.Sprintf("  → Tamaño: %d %s (%d bytes)", size, unitName, sizeInBytes))
	utils.LogSuccess("MkDisk", fmt.Sprintf("  → Ajuste: %s", fit))
	utils.LogSuccess("MkDisk", fmt.Sprintf("  → Signature: %d", mbr.MbrDiskSignature))
	utils.LogSuccess("MkDisk", fmt.Sprintf("  → Fecha creación: %d", mbr.MbrFechaCreacion))

	return nil
}
