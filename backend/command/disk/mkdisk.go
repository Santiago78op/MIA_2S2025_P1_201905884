package disk

/*
 * Este comando creará un archivo binario que simulará un disco, estos archivos
 * binarios tendrán la extensión .mia y su contenido al inicio será 0 binarios.
 * Deberá ocupar físicamente el tamaño indicado por los parámetros
 */

import (
	utils "backend/Utils"
	"fmt"
	"os"
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
	// Validar tamaño
	if size <= 0 {
		// Error en log
		utils.LogError("MkDisk", "El tamaño del disco debe ser un número positivo mayor que cero")
	}

	// Validar fit
	validFits := map[string]bool{"BF": true, "FF": true, "WF": true}
	if fit != "" && !validFits[fit] {
		utils.LogError("MkDisk", "Ajuste no válido, use BF, FF o WF")
		return fmt.Errorf("ajuste no válido, use BF, FF o WF")
	}
	if fit == "" {
		fit = "FF" // Default to First Fit
	}

	// Validar unidad
	unitMultiplier := int64(1024 * 1024) // Por defecto Megabytes
	if unit != "" {
		switch unit {
		case "K":
			unitMultiplier = 1024
		case "M":
			unitMultiplier = 1024 * 1024
		default:
			utils.LogError("MkDisk", "Unidad no válida, use K para Kilobytes o M para Megabytes")
			return fmt.Errorf("unidad no válida, use K para Kilobytes o M para Megabytes")
		}
	}

	// Calcular el tamaño en bytes
	sizeInBytes := size * unitMultiplier

	// Validar que el path sea obligatorio
	if path == "" {
		utils.LogError("MkDisk", "El parámetro -path es obligatorio")
		return fmt.Errorf("el parámetro -path es obligatorio")
	}

	// Crear el archivo del disco
	file, err := os.Create(path)
	if err != nil {
		utils.LogError("MkDisk", fmt.Sprintf("Error al crear el archivo: %v", err))
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	// Escribir ceros en el archivo hasta alcanzar el tamaño especificado
	if _, err := file.Write(make([]byte, sizeInBytes)); err != nil {
		utils.LogError("MkDisk", fmt.Sprintf("Error al escribir en el archivo: %v", err))
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	// Crear el Disco con los parámetros especificados

	// Logica del MBR

	utils.LogInfo("MkDisk", fmt.Sprintf("Disco creado con éxito en %s de tamaño %d bytes con ajuste %s", path, sizeInBytes, fit))
	return nil
}
