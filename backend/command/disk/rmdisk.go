package disk

/*
 * Este comando elimina un archivo que representa a un disco duro virtual.
 * El archivo debe existir para poder ser eliminado.
 */

import (
	utils "backend/Utils"
	"fmt"
	"os"
	"strings"
)

/*
| PARÁMETRO | CATEGORÍA  | DESCRIPCIÓN                                                                                                    |
|-----------|------------|----------------------------------------------------------------------------------------------------------------|
| -path     | Obligatorio | Ruta del archivo que representa el disco duro a eliminar. Si el archivo no existe, debe mostrar un error. |
*/

// RmDisk elimina un archivo que representa un disco duro virtual
func RmDisk(path string) error {
	utils.LogInfo("RmDisk", fmt.Sprintf("Iniciando eliminación de disco: path=%s", path))

	// Validar que el path sea obligatorio
	if path == "" {
		utils.LogError("RmDisk", "El parámetro -path es obligatorio")
		return fmt.Errorf("el parámetro -path es obligatorio")
	}

	// Limpiar el path
	path = strings.TrimSpace(path)

	// Verificar si el archivo existe
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		utils.LogError("RmDisk", fmt.Sprintf("El archivo no existe en la ruta: %s", path))
		return fmt.Errorf("el archivo no existe en la ruta: %s", path)
	}
	if err != nil {
		utils.LogError("RmDisk", fmt.Sprintf("Error al acceder al archivo: %v", err))
		return fmt.Errorf("error al acceder al archivo: %v", err)
	}

	// Verificar que sea un archivo regular (no un directorio)
	if fileInfo.IsDir() {
		utils.LogError("RmDisk", fmt.Sprintf("La ruta especificada es un directorio, no un archivo: %s", path))
		return fmt.Errorf("la ruta especificada es un directorio, no un archivo: %s", path)
	}

	// Obtener información del archivo antes de eliminarlo
	fileSize := fileInfo.Size()
	utils.LogInfo("RmDisk", fmt.Sprintf("Archivo encontrado: %s (%.2f MB)", path, float64(fileSize)/(1024*1024)))

	// Verificar que tenga la extensión .mia (advertencia, no error)
	if !strings.HasSuffix(strings.ToLower(path), ".mia") {
		utils.LogWarning("RmDisk", "El archivo no tiene la extensión .mia típica de discos virtuales")
	}

	// Opcional: Leer y verificar el MBR antes de eliminar (para confirmar que es un disco válido)
	// Esto es útil para evitar eliminar archivos que no son discos por error
	/*
		_, err = estructuras.ReadMBR(path)
		if err != nil {
			utils.LogWarning("RmDisk", "El archivo no parece tener un MBR válido, pero se procederá con la eliminación")
		}
	*/

	// Confirmar la eliminación (en un entorno real, aquí se podría pedir confirmación al usuario)
	utils.LogWarning("RmDisk", fmt.Sprintf("Se eliminará permanentemente el archivo: %s", path))

	// Eliminar el archivo
	err = os.Remove(path)
	if err != nil {
		utils.LogError("RmDisk", fmt.Sprintf("Error al eliminar el archivo: %v", err))
		return fmt.Errorf("error al eliminar el archivo: %v", err)
	}

	// Verificar que el archivo fue eliminado correctamente
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		utils.LogError("RmDisk", "El archivo no fue eliminado correctamente")
		return fmt.Errorf("error: el archivo no fue eliminado correctamente")
	}

	utils.LogSuccess("RmDisk", "Disco eliminado exitosamente:")
	utils.LogSuccess("RmDisk", fmt.Sprintf("  → Ruta: %s", path))
	utils.LogSuccess("RmDisk", fmt.Sprintf("  → Tamaño eliminado: %.2f MB", float64(fileSize)/(1024*1024)))

	return nil
}
