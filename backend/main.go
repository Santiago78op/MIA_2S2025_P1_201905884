package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
}

type FileSystemInfo struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Size       int64  `json:"size"`
	MountPoint string `json:"mountPoint"`
}

// Handlers
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "OK",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Message: "Backend funcionando correctamente",
		Data:    response,
		Status:  "success",
	})
}

func getFileSystemsHandler(w http.ResponseWriter, r *http.Request) {
	// Datos de ejemplo (luego se conectará con la lógica real)
	fileSystems := []FileSystemInfo{
		{
			Name:       "disk1.dk",
			Type:       "EXT2",
			Size:       1024000,
			MountPoint: "/mnt/disk1",
		},
		{
			Name:       "disk2.dk",
			Type:       "EXT3",
			Size:       2048000,
			MountPoint: "/mnt/disk2",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Message: "Sistemas de archivos obtenidos",
		Data:    fileSystems,
		Status:  "success",
	})
}

func createPartitionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var requestData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Aquí iría la lógica para crear partición
	log.Printf("Creando partición: %+v", requestData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Message: "Partición creada exitosamente",
		Data:    requestData,
		Status:  "success",
	})
}

func executeCommandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Command string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Aquí iría la lógica para ejecutar comandos
	log.Printf("Ejecutando comando: %s", requestData.Command)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Message: "Comando ejecutado",
		Data: map[string]string{
			"command": requestData.Command,
			"result":  "Comando ejecutado exitosamente",
		},
		Status: "success",
	})
}

func main() {
	// Crear router
	router := mux.NewRouter()

	// Rutas de API
	api := router.PathPrefix("/api").Subrouter()
	
	// Endpoints
	api.HandleFunc("/health", healthHandler).Methods("GET")
	api.HandleFunc("/filesystems", getFileSystemsHandler).Methods("GET")
	api.HandleFunc("/partition", createPartitionHandler).Methods("POST")
	api.HandleFunc("/execute", executeCommandHandler).Methods("POST")

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Aplicar middleware CORS
	handler := c.Handler(router)

	// Configurar servidor
	port := ":8080"
	fmt.Printf("🚀 Servidor iniciado en http://localhost%s\n", port)
	fmt.Println("📡 Endpoints disponibles:")
	fmt.Println("   GET  /api/health       - Estado del servidor")
	fmt.Println("   GET  /api/filesystems  - Listar sistemas de archivos")
	fmt.Println("   POST /api/partition     - Crear partición")
	fmt.Println("   POST /api/execute       - Ejecutar comando")

	// Iniciar servidor
	log.Fatal(http.ListenAndServe(port, handler))
}
