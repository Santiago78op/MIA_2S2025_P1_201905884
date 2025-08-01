package action

import (
	utils "backend/Utils"
	"fmt"
	"os"
	"path/filepath"
)

/*
 * AdminDisk es una acción que permite administrar discos duros virtuales.
 */
func NewDisk(path string, sizeInBytes int64) error {

	// Verificar si el archivo ya existe
	if _, err := os.Stat(path); err == nil {
		utils.LogError("NewDisk", fmt.Sprintf("El archivo ya existe en la ruta: %s", path))
		return fmt.Errorf("el archivo ya existe en la ruta: %s", path)
	}

	// Crear directorios padre si no existen
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		utils.LogError("NewDisk", fmt.Sprintf("Error al crear directorios: %v", err))
		return fmt.Errorf("error al crear directorios: %v", err)
	}

	// Crear el archivo del disco
	file, err := os.Create(path)
	if err != nil {
		utils.LogError("NewDisk", fmt.Sprintf("Error al crear el archivo: %v", err))
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	// Método más eficiente para archivos grandes: usar Seek para establecer el tamaño
	if sizeInBytes > 0 {
		// Posicionarse al final del tamaño deseado
		_, err := file.Seek(sizeInBytes-1, 0)
		if err != nil {
			utils.LogError("NewDisk", fmt.Sprintf("Error al posicionar en el archivo: %v", err))
			return fmt.Errorf("error al posicionar en el archivo: %v", err)
		}

		// Escribir un byte al final para establecer el tamaño del archivo
		_, err = file.Write([]byte{0})
		if err != nil {
			utils.LogError("NewDisk", fmt.Sprintf("Error al escribir en el archivo: %v", err))
			return fmt.Errorf("error al escribir en el archivo: %v", err)
		}
	}

	utils.LogInfo("NewDisk", fmt.Sprintf("Disco creado con éxito en %s de tamaño %d bytes", path, sizeInBytes))
	return nil
}
