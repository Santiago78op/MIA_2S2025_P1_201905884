package command

import (
	utils "backend/Utils"
	diskCommands "backend/command/disk"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CommandResult representa el resultado de ejecutar un comando
type CommandResult struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// CommandParser maneja el parsing y ejecución de comandos
type CommandParser struct {
	// Se puede agregar estado del sistema aquí si es necesario
}

// NewCommandParser crea una nueva instancia del parser de comandos
func NewCommandParser() *CommandParser {
	return &CommandParser{}
}

// ParseAndExecute parsea y ejecuta un comando dado como string
func (cp *CommandParser) ParseAndExecute(commandLine string) *CommandResult {
	// Limpiar el comando
	commandLine = strings.TrimSpace(commandLine)

	// Ignorar líneas vacías y comentarios
	if commandLine == "" || strings.HasPrefix(commandLine, "#") {
		return &CommandResult{
			Success: true,
			Message: "Línea ignorada (vacía o comentario)",
		}
	}

	utils.LogInfo("Parser", fmt.Sprintf("Procesando comando: %s", commandLine))

	// Parsear el comando y sus parámetros
	parts, err := cp.parseCommandLine(commandLine)
	if err != nil {
		utils.LogError("Parser", fmt.Sprintf("Error al parsear comando: %v", err))
		return &CommandResult{
			Success: false,
			Error:   fmt.Sprintf("Error de sintaxis: %v", err),
		}
	}

	if len(parts) == 0 {
		return &CommandResult{
			Success: false,
			Error:   "Comando vacío",
		}
	}

	// Obtener el comando principal
	command := strings.ToLower(parts[0])

	// Parsear parámetros
	params, err := cp.parseParameters(parts[1:])
	if err != nil {
		utils.LogError("Parser", fmt.Sprintf("Error al parsear parámetros: %v", err))
		return &CommandResult{
			Success: false,
			Error:   fmt.Sprintf("Error en parámetros: %v", err),
		}
	}

	// Ejecutar el comando correspondiente
	return cp.executeCommand(command, params)
}

// parseCommandLine divide la línea de comando en partes, respetando comillas
func (cp *CommandParser) parseCommandLine(commandLine string) ([]string, error) {
	// Regex para dividir respetando comillas
	re := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)'`)
	matches := re.FindAllStringSubmatch(commandLine, -1)

	var parts []string
	for _, match := range matches {
		if match[1] != "" {
			// Contenido entre comillas dobles
			parts = append(parts, match[1])
		} else if match[2] != "" {
			// Contenido entre comillas simples
			parts = append(parts, match[2])
		} else {
			// Contenido sin comillas
			parts = append(parts, match[0])
		}
	}

	return parts, nil
}

// parseParameters parsea los parámetros del comando (formato -param=value)
func (cp *CommandParser) parseParameters(paramParts []string) (map[string]string, error) {
	params := make(map[string]string)

	for _, part := range paramParts {
		if !strings.HasPrefix(part, "-") {
			return nil, fmt.Errorf("parámetro inválido: %s (debe comenzar con -)", part)
		}

		// Remover el guión inicial
		part = part[1:]

		// Dividir en nombre=valor
		equalIndex := strings.Index(part, "=")
		if equalIndex == -1 {
			// Parámetro sin valor (flag)
			params[part] = "true"
		} else {
			paramName := part[:equalIndex]
			paramValue := part[equalIndex+1:]
			params[paramName] = paramValue
		}
	}

	return params, nil
}

// executeCommand ejecuta el comando especificado con los parámetros dados
func (cp *CommandParser) executeCommand(command string, params map[string]string) *CommandResult {
	switch command {
	case "mkdisk":
		return cp.executeMkDisk(params)
	case "rmdisk":
		return cp.executeRmDisk(params)
	case "fdisk":
		return cp.executeFdisk(params)
	case "mount":
		return cp.executeMount(params)
	case "unmount":
		return cp.executeUnmount(params)
	case "mkfs":
		return cp.executeMkfs(params)
	default:
		utils.LogError("Parser", fmt.Sprintf("Comando no reconocido: %s", command))
		return &CommandResult{
			Success: false,
			Error:   fmt.Sprintf("Comando no reconocido: %s", command),
		}
	}
}

// executeMkDisk ejecuta el comando mkdisk
func (cp *CommandParser) executeMkDisk(params map[string]string) *CommandResult {
	// Validar parámetros obligatorios
	sizeStr, hasSizeParam := params["size"]
	path, hasPath := params["path"]

	if !hasSizeParam {
		return &CommandResult{
			Success: false,
			Error:   "El parámetro -size es obligatorio",
		}
	}

	if !hasPath {
		return &CommandResult{
			Success: false,
			Error:   "El parámetro -path es obligatorio",
		}
	}

	// Convertir size a número
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return &CommandResult{
			Success: false,
			Error:   fmt.Sprintf("El valor de -size debe ser un número válido: %v", err),
		}
	}

	// Parámetros opcionales
	fit := params["fit"]
	unit := params["unit"]

	// Ejecutar el comando
	err = diskCommands.MkDisk(size, fit, unit, path)
	if err != nil {
		return &CommandResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	return &CommandResult{
		Success: true,
		Message: fmt.Sprintf("Disco creado exitosamente en %s", path),
		Data: map[string]interface{}{
			"path": path,
			"size": size,
			"fit":  fit,
			"unit": unit,
		},
	}
}

// executeRmDisk ejecuta el comando rmdisk
func (cp *CommandParser) executeRmDisk(params map[string]string) *CommandResult {
	// Validar parámetros obligatorios
	path, hasPath := params["path"]

	if !hasPath {
		return &CommandResult{
			Success: false,
			Error:   "El parámetro -path es obligatorio",
		}
	}

	// Ejecutar el comando
	err := diskCommands.RmDisk(path)
	if err != nil {
		return &CommandResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	return &CommandResult{
		Success: true,
		Message: fmt.Sprintf("Disco eliminado exitosamente: %s", path),
		Data: map[string]interface{}{
			"path": path,
		},
	}
}

// executeFdisk ejecuta el comando fdisk (placeholder por ahora)
func (cp *CommandParser) executeFdisk(params map[string]string) *CommandResult {
	utils.LogWarning("Parser", "Comando FDISK aún no implementado")
	return &CommandResult{
		Success: false,
		Error:   "Comando FDISK no implementado aún",
	}
}

// executeMount ejecuta el comando mount (placeholder por ahora)
func (cp *CommandParser) executeMount(params map[string]string) *CommandResult {
	utils.LogWarning("Parser", "Comando MOUNT aún no implementado")
	return &CommandResult{
		Success: false,
		Error:   "Comando MOUNT no implementado aún",
	}
}

// executeUnmount ejecuta el comando unmount (placeholder por ahora)
func (cp *CommandParser) executeUnmount(params map[string]string) *CommandResult {
	utils.LogWarning("Parser", "Comando UNMOUNT aún no implementado")
	return &CommandResult{
		Success: false,
		Error:   "Comando UNMOUNT no implementado aún",
	}
}

// executeMkfs ejecuta el comando mkfs (placeholder por ahora)
func (cp *CommandParser) executeMkfs(params map[string]string) *CommandResult {
	utils.LogWarning("Parser", "Comando MKFS aún no implementado")
	return &CommandResult{
		Success: false,
		Error:   "Comando MKFS no implementado aún",
	}
}

// ExecuteScript ejecuta un script con múltiples comandos
func (cp *CommandParser) ExecuteScript(script string) []*CommandResult {
	lines := strings.Split(script, "\n")
	var results []*CommandResult

	utils.LogInfo("Parser", fmt.Sprintf("Ejecutando script con %d líneas", len(lines)))

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Agregar información de línea para debugging
		utils.LogInfo("Parser", fmt.Sprintf("Línea %d: %s", i+1, line))

		result := cp.ParseAndExecute(line)

		// Agregar información de línea al resultado
		if result.Data == nil {
			result.Data = make(map[string]interface{})
		}
		if dataMap, ok := result.Data.(map[string]interface{}); ok {
			dataMap["line_number"] = i + 1
		}

		results = append(results, result)

		// Si hay un error crítico, se puede decidir si continuar o no
		if !result.Success && result.Error != "" {
			utils.LogError("Parser", fmt.Sprintf("Error en línea %d: %s", i+1, result.Error))
			// Por ahora continuamos con las siguientes líneas
		}
	}

	utils.LogInfo("Parser", fmt.Sprintf("Script ejecutado. Resultados: %d comandos procesados", len(results)))
	return results
}

// ValidateCommand valida la sintaxis de un comando sin ejecutarlo
func (cp *CommandParser) ValidateCommand(commandLine string) error {
	commandLine = strings.TrimSpace(commandLine)

	// Ignorar líneas vacías y comentarios
	if commandLine == "" || strings.HasPrefix(commandLine, "#") {
		return nil
	}

	// Parsear el comando
	parts, err := cp.parseCommandLine(commandLine)
	if err != nil {
		return fmt.Errorf("error de sintaxis: %v", err)
	}

	if len(parts) == 0 {
		return fmt.Errorf("comando vacío")
	}

	// Validar que el comando existe
	command := strings.ToLower(parts[0])
	validCommands := []string{"mkdisk", "rmdisk", "fdisk", "mount", "unmount", "mkfs", "login", "logout"}

	found := false
	for _, validCmd := range validCommands {
		if command == validCmd {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("comando no reconocido: %s", command)
	}

	// Validar parámetros
	_, err = cp.parseParameters(parts[1:])
	if err != nil {
		return fmt.Errorf("error en parámetros: %v", err)
	}

	return nil
}

// GetSupportedCommands retorna la lista de comandos soportados
func (cp *CommandParser) GetSupportedCommands() []string {
	return []string{
		"mkdisk",  // Crear disco
		"rmdisk",  // Eliminar disco
		"fdisk",   // Administrar particiones
		"mount",   // Montar partición
		"unmount", // Desmontar partición
		"mkfs",    // Formatear partición
		"login",   // Iniciar sesión
		"logout",  // Cerrar sesión
		"mkgrp",   // Crear grupo
		"rmgrp",   // Eliminar grupo
		"mkusr",   // Crear usuario
		"rmusr",   // Eliminar usuario
		"chgrp",   // Cambiar grupo
		"mkfile",  // Crear archivo
		"mkdir",   // Crear directorio
		"cat",     // Mostrar contenido
		"rep",     // Generar reportes
	}
}
