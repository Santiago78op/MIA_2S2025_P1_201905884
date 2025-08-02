package main

import (
	utils "backend/Utils"
	"backend/command"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

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
	// Por ahora devolvemos datos de ejemplo
	// TODO: Implementar lógica real para obtener sistemas montados
	fileSystems := []FileSystemInfo{
		{
			Name:       "disk1.mia",
			Type:       "EXT2",
			Size:       1048576, // 1MB en bytes
			MountPoint: "/mnt/341A",
			Status:     "mounted",
		},
		{
			Name:       "disk2.mia",
			Type:       "EXT2",
			Size:       2097152, // 2MB en bytes
			MountPoint: "/mnt/341B",
			Status:     "unmounted",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Message: "Sistemas de archivos obtenidos",
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
	api.HandleFunc("/validate", validateCommandHandler).Methods("POST")
	api.HandleFunc("/execute", executeCommandHandler).Methods("POST")

	// Endpoints para logs y comunicación en tiempo real
	api.HandleFunc("/logs", utils.GetLogsHandler).Methods("GET")
	api.HandleFunc("/logs/stream", utils.SSEHandler).Methods("GET")
	router.HandleFunc("/ws", wsHandler).Methods("GET")

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Aplicar middleware CORS
	handler := c.Handler(router)

	// Configurar servidor
	port := ":8080"

	// Mensaje de inicio
	fmt.Println("🚀 ====================================")
	fmt.Println("🚀    ExtreamFS Backend Server")
	fmt.Println("🚀 ====================================")
	fmt.Printf("🚀 Servidor iniciado en http://localhost%s\n", port)
	fmt.Println("🚀 Proyecto: Simulador de Sistema de Archivos EXT2")
	fmt.Println("🚀 Versión: 1.0.0")
	fmt.Println("🚀")
	fmt.Println("📡 Endpoints disponibles:")
	fmt.Println("   GET  /api/health         - Estado del servidor")
	fmt.Println("   GET  /api/filesystems    - Listar sistemas de archivos")
	fmt.Println("   GET  /api/commands       - Comandos soportados")
	fmt.Println("   POST /api/validate       - Validar sintaxis de comandos")
	fmt.Println("   POST /api/execute        - Ejecutar comandos")
	fmt.Println("   GET  /api/logs           - Obtener logs (polling)")
	fmt.Println("   GET  /api/logs/stream    - Stream de logs (SSE)")
	fmt.Println("   GET  /ws                 - WebSocket para logs en tiempo real")
	fmt.Println("🚀")
	fmt.Println("📝 Comandos implementados:")

	commands := commandParser.GetSupportedCommands()
	for i, cmd := range commands {
		if i < 6 { // Mostrar solo los primeros comandos implementados
			status := "✅"
			if cmd == "mkdisk" || cmd == "rmdisk" {
				status = "✅ IMPLEMENTADO"
			} else {
				status = "🚧 EN DESARROLLO"
			}
			fmt.Printf("   %-12s %s\n", cmd, status)
		}
	}
	fmt.Printf("   %-12s %s\n", "...", fmt.Sprintf("y %d comandos más", len(commands)-6))

	fmt.Println("🚀")
	fmt.Println("🔧 Para probar el sistema:")
	fmt.Println("   curl http://localhost:8080/api/health")
	fmt.Println("🚀 ====================================")

	// Log inicial
	utils.LogSuccess("SYSTEM", "ExtreamFS Backend iniciado correctamente")
	utils.LogInfo("SYSTEM", fmt.Sprintf("Servidor escuchando en puerto %s", port))

	// Iniciar servidor
	log.Fatal(http.ListenAndServe(port, handler))
}
