package main

import (
	utils "backend/Utils"
	"backend/command"
	diskCommands "backend/command/disk"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

// ConfigureCarnet configura el sufijo del carnet para generar IDs de particiones
func ConfigureCarnet() {
	// Intentar obtener el carnet de argumentos de línea de comandos
	carnet := flag.String("carnet", "", "Número de carnet completo (ej: 202401234)")
	flag.Parse()

	var suffix string

	if *carnet != "" {
		// Extraer últimos dos dígitos del carnet
		if len(*carnet) >= 2 {
			suffix = (*carnet)[len(*carnet)-2:]
		}
	} else {
		// Intentar obtener de variable de entorno
		carnetEnv := os.Getenv("STUDENT_CARNET")
		if carnetEnv != "" && len(carnetEnv) >= 2 {
			suffix = carnetEnv[len(carnetEnv)-2:]
		}
	}

	// Validar que sean dígitos
	if suffix != "" {
		matched, _ := regexp.MatchString(`^\d{2}$`, suffix)
		if matched {
			diskCommands.SetCarnetSuffix(suffix)
			fmt.Printf("✅ Carnet configurado: IDs de partición usarán el sufijo '%s'\n", suffix)
		} else {
			fmt.Printf("⚠️  Carnet inválido, usando sufijo por defecto '34'\n")
		}
	} else {
		fmt.Printf("ℹ️  No se especificó carnet, usando sufijo por defecto '34'\n")
		fmt.Printf("   Para configurar: -carnet=202401234 o export STUDENT_CARNET=202401234\n")
	}
}

// Estructuras de respuesta
type ApiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}

type FileSystemInfo struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Size       int64  `json:"size"`
	MountPoint string `json:"mountPoint"`
	Status     string `json:"status"`
	Path       string `json:"path"`
}

type ExecuteRequest struct {
	Command string `json:"command"`
	Script  string `json:"script,omitempty"`
}

type ExecuteResponse struct {
	Results []command.CommandResult `json:"results"`
	Summary ExecuteSummary          `json:"summary"`
}

type ExecuteSummary struct {
	TotalCommands      int    `json:"total_commands"`
	SuccessfulCommands int    `json:"successful_commands"`
	FailedCommands     int    `json:"failed_commands"`
	ExecutionTime      string `json:"execution_time"`
}

// Variables globales
var (
	startTime     = time.Now()
	commandParser = command.NewCommandParser()
)

// Handlers
func healthHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime).Round(time.Second)

	response := HealthResponse{
		Status:    "OK",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    uptime.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Message: "Backend del sistema de archivos ExtreamFS funcionando correctamente",
		Data:    response,
		Status:  "success",
	})
}

// getFileSystemsHandler obtiene información de los sistemas de archivos montados
func getFileSystemsHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetro de ruta de la query string
	searchPath := r.URL.Query().Get("path")
	if searchPath == "" {
		searchPath = "./Discos" // Ruta por defecto
	}

	// Limpiar y expandir la ruta
	searchPath = strings.TrimSpace(searchPath)

	// Convertir rutas relativas a absolutas si es necesario
	if strings.HasPrefix(searchPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			searchPath = filepath.Join(homeDir, searchPath[2:])
		}
	}

	// Limpiar la ruta
	searchPath = filepath.Clean(searchPath)

	var fileSystems []FileSystemInfo

	// Verificar si el directorio existe
	if _, err := os.Stat(searchPath); os.IsNotExist(err) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{
			Message: fmt.Sprintf("Directorio no encontrado: %s", searchPath),
			Data:    []FileSystemInfo{},
			Status:  "warning",
		})
		return
	}

	// Buscar archivos .mia del directorio especificado
	pattern := filepath.Join(searchPath, "*.mia")
	files, err := filepath.Glob(pattern)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ApiResponse{
			Message: fmt.Sprintf("Error buscando archivos .mia: %v", err),
			Data:    []FileSystemInfo{},
			Status:  "error",
		})
		return
	}

	// También buscar archivos .dsk
	dskPattern := filepath.Join(searchPath, "*.dsk")
	dskFiles, err := filepath.Glob(dskPattern)
	if err == nil {
		files = append(files, dskFiles...)
	}

	// Búsqueda manual adicional si Glob falla
	if len(files) == 0 {
		entries, err := os.ReadDir(searchPath)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					name := entry.Name()
					if strings.HasSuffix(strings.ToLower(name), ".mia") || strings.HasSuffix(strings.ToLower(name), ".dsk") {
						fullPath := filepath.Join(searchPath, name)
						files = append(files, fullPath)
					}
				}
			}
		}
	}

	for _, file := range files {
		if stat, err := os.Stat(file); err == nil {
			fileName := filepath.Base(file)
			fileExt := strings.ToLower(filepath.Ext(fileName))

			// Determinar tipo por extensión
			diskType := "EXT2"
			if fileExt == ".dsk" {
				diskType = "DSK"
			}

			// Determinar si está montado
			status := "unmounted"
			mountPoint := ""

			fileSystems = append(fileSystems, FileSystemInfo{
				Name:       fileName,
				Type:       diskType,
				Size:       stat.Size(),
				MountPoint: mountPoint,
				Status:     status,
				Path:       file,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Message: fmt.Sprintf("Encontrados %d sistemas de archivos", len(fileSystems)),
		Data:    fileSystems,
		Status:  "success",
	})
}

// validateCommandHandler valida la sintaxis de un comando sin ejecutarlo
func validateCommandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var requestData ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	var errors []string
	var validCommands []string

	if requestData.Command != "" {
		// Validar comando único
		if err := commandParser.ValidateCommand(requestData.Command); err != nil {
			errors = append(errors, err.Error())
		} else {
			validCommands = append(validCommands, requestData.Command)
		}
	}

	if requestData.Script != "" {
		// Validar script línea por línea
		lines := strings.Split(requestData.Script, "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			if err := commandParser.ValidateCommand(line); err != nil {
				errors = append(errors, fmt.Sprintf("Línea %d: %s", i+1, err.Error()))
			} else {
				validCommands = append(validCommands, line)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if len(errors) > 0 {
		json.NewEncoder(w).Encode(ApiResponse{
			Message: "Se encontraron errores de validación",
			Data: map[string]interface{}{
				"errors":         errors,
				"valid_commands": validCommands,
			},
			Status: "error",
		})
	} else {
		json.NewEncoder(w).Encode(ApiResponse{
			Message: "Todos los comandos son válidos",
			Data: map[string]interface{}{
				"valid_commands": validCommands,
			},
			Status: "success",
		})
	}
}

// executeCommandHandler ejecuta uno o más comandos
func executeCommandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var requestData ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	startExecution := time.Now()
	var results []command.CommandResult

	if requestData.Command != "" && requestData.Script != "" {
		// No permitir ambos al mismo tiempo
		http.Error(w, "Debe especificar solo 'command' o 'script', no ambos", http.StatusBadRequest)
		return
	}

	if requestData.Command != "" {
		// Ejecutar comando único
		utils.LogInfo("API", fmt.Sprintf("Ejecutando comando único: %s", requestData.Command))
		result := commandParser.ParseAndExecute(requestData.Command)
		results = append(results, *result)
	} else if requestData.Script != "" {
		// Ejecutar script
		utils.LogInfo("API", "Ejecutando script con múltiples comandos")
		scriptResults := commandParser.ExecuteScript(requestData.Script)
		for _, result := range scriptResults {
			results = append(results, *result)
		}
	} else {
		http.Error(w, "Debe especificar 'command' o 'script'", http.StatusBadRequest)
		return
	}

	executionTime := time.Since(startExecution)

	// Calcular estadísticas
	successful := 0
	failed := 0
	for _, result := range results {
		if result.Success {
			successful++
		} else {
			failed++
		}
	}

	summary := ExecuteSummary{
		TotalCommands:      len(results),
		SuccessfulCommands: successful,
		FailedCommands:     failed,
		ExecutionTime:      executionTime.String(),
	}

	response := ExecuteResponse{
		Results: results,
		Summary: summary,
	}

	// Log del resumen
	utils.LogInfo("API", fmt.Sprintf("Ejecución completada: %d comandos (%d exitosos, %d fallidos) en %v",
		summary.TotalCommands, summary.SuccessfulCommands, summary.FailedCommands, executionTime))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Message: fmt.Sprintf("Ejecutados %d comandos", len(results)),
		Data:    response,
		Status:  "success",
	})
}

// getSupportedCommandsHandler retorna la lista de comandos soportados
func getSupportedCommandsHandler(w http.ResponseWriter, r *http.Request) {
	commands := commandParser.GetSupportedCommands()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Message: "Comandos soportados",
		Data: map[string]interface{}{
			"commands": commands,
			"total":    len(commands),
		},
		Status: "success",
	})
}

// wsHandler maneja las conexiones WebSocket para logs en tiempo real
func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	utils.AddWSConnection(conn)
	defer utils.RemoveWSConnection(conn)

	// Enviar mensaje de bienvenida
	welcome := utils.LogMessage{
		Type:    utils.INFO,
		Command: "SYSTEM",
		Message: "Conexión WebSocket establecida para logs en tiempo real",
		Time:    fmt.Sprintf("%d", time.Now().Unix()),
	}
	conn.WriteJSON(welcome)

	// Mantener la conexión viva y manejar mensajes del cliente
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			break
		}

		// Procesar mensajes del cliente si es necesario
		if msgType, ok := msg["type"].(string); ok && msgType == "ping" {
			pong := map[string]interface{}{
				"type":      "pong",
				"timestamp": time.Now().Unix(),
			}
			conn.WriteJSON(pong)
		}
	}
}

// debugPathHandler para diagnosticar problemas de rutas
func debugPathHandler(w http.ResponseWriter, r *http.Request) {
	searchPath := r.URL.Query().Get("path")
	if searchPath == "" {
		searchPath = "./Discos"
	}

	debug := map[string]interface{}{
		"requested_path": searchPath,
		"cleaned_path":   filepath.Clean(searchPath),
	}

	// Información del directorio actual
	if wd, err := os.Getwd(); err == nil {
		debug["working_directory"] = wd
	}

	// Verificar si existe
	if stat, err := os.Stat(searchPath); err == nil {
		debug["path_exists"] = true
		debug["is_directory"] = stat.IsDir()
		debug["size"] = stat.Size()
		debug["permissions"] = stat.Mode().String()
	} else {
		debug["path_exists"] = false
		debug["error"] = err.Error()
	}

	// Listar contenido si es directorio
	if entries, err := os.ReadDir(searchPath); err == nil {
		var files []string
		for _, entry := range entries {
			files = append(files, entry.Name())
		}
		debug["directory_contents"] = files
	}

	// Probar patrones de búsqueda
	patterns := []string{"*.mia", "*.dsk", "*"}
	for _, pattern := range patterns {
		fullPattern := filepath.Join(searchPath, pattern)
		if matches, err := filepath.Glob(fullPattern); err == nil {
			debug[fmt.Sprintf("pattern_%s", pattern)] = matches
		} else {
			debug[fmt.Sprintf("pattern_%s_error", pattern)] = err.Error()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Message: "Información de debug de ruta",
		Data:    debug,
		Status:  "success",
	})
}

func main() {
	// Crear router
	router := mux.NewRouter()

	// Rutas de API
	api := router.PathPrefix("/api").Subrouter()

	// Endpoints básicos
	api.HandleFunc("/health", healthHandler).Methods("GET")
	api.HandleFunc("/filesystems", getFileSystemsHandler).Methods("GET")
	api.HandleFunc("/commands", getSupportedCommandsHandler).Methods("GET")

	// Endpoints para comandos
	api.HandleFunc("/execute", executeCommandHandler).Methods("POST")
	api.HandleFunc("/validate", validateCommandHandler).Methods("POST")

	// Endpoints para logs y comunicación en tiempo real
	api.HandleFunc("/ws", wsHandler)

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Aplicar middleware CORS
	handler := c.Handler(router)

	// Configurar servidor
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Mensaje de inicio
	fmt.Println("🚀 Servidor iniciado en http://localhost:8080")
	fmt.Println("📡 API disponible en http://localhost:8080/api")
	fmt.Println("🔗 WebSocket disponible en ws://localhost:8080/api/ws")

	// Log inicial
	utils.LogInfo("SERVER", "Backend del sistema de archivos ExtreamFS iniciado correctamente")

	// Iniciar servidor
	log.Fatal(server.ListenAndServe())
}
