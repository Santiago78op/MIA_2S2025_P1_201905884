package utils

/*
 * Logger es una utilidad del sistema para mensajes del frontend,
 * permite registrar información relevante durante la ejecución de la aplicación.
 * Probee diferentes niveles de log para categorizar la importancia de los mensajes.
 * INFO, WARNING, ERROR, SUCCESS.
 */

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Constantes para los niveles de log
const (
	INFO    = "INFO"
	WARNING = "WARNING"
	ERROR   = "ERROR"
	SUCCESS = "SUCCESS"
)

// LogMessage es una estructura que maneja los mensajes del sistema
type LogMessage struct {
	Type    string `json:"type"`    // Tipo de mensaje: INFO, WARNING, ERROR, SUCCESS
	Command string `json:"command"` // Comando asociado al mensaje
	Message string `json:"message"` // Mensaje a registrar
	Time    string `json:"time"`    // Hora del mensaje
}

// Var para manejar conexiones WebSocket
var (
	wsConnections    []*websocket.Conn // Conexiones WebSocket activas
	connectionsMutex sync.RWMutex      // Mutex para manejar concurrencia en conexiones
)

// AddWSConnection agrega una nueva conexión WebSocket a la lista de conexiones activas
func AddWSConnection(conn *websocket.Conn) {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	// Agregar una nueva conexión a wsConnections
	wsConnections = append(wsConnections, conn)
}

// RemoveWSConnection remueve una conexión WebSocket de la lista de conexiones activas
func RemoveWSConnection(conn *websocket.Conn) {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	// Buscar y eliminar la conexión de wsConnections
	for i, c := range wsConnections {
		if c == conn {
			wsConnections = append(wsConnections[:i], wsConnections[i+1:]...)
			break
		}
	}
}

// NewLogger crea una nueva instancia de Logger
func NewLogger(tipo, comando, mensaje string) *LogMessage {
	return &LogMessage{
		Type:    tipo,
		Command: comando,
		Message: mensaje,
		Time:    fmt.Sprintf("%d", time.Now().Unix()), // Usar timestamp como hora
	}
}

// Log imprime el mensaje formateado según el tipo
func (l *LogMessage) Log() {
	// Formatear el mensaje
	formattedMessage := fmt.Sprintf("[%s] [%s] %s: %s", l.Time, l.Type, l.Command, l.Message)

	// Imprimir en consola
	fmt.Println(formattedMessage)

	// Debug: Mostrar cuántas conexiones WebSocket hay activas
	connectionsMutex.RLock()
	connectionCount := len(wsConnections)
	connectionsMutex.RUnlock()

	fmt.Printf("DEBUG: Enviando mensaje a %d conexiones WebSocket\n", connectionCount)

	// Enviar a todas las conexiones WebSocket activas como JSON
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	var activeConnections []*websocket.Conn
	for _, conn := range wsConnections {
		if err := conn.WriteJSON(*l); err != nil {
			fmt.Printf("Error al enviar mensaje a WebSocket: %v\n", err)
			// No agregar conexiones con error a la lista activa
		} else {
			activeConnections = append(activeConnections, conn)
		}
	}

	// Actualizar la lista solo con conexiones activas
	wsConnections = activeConnections

	// Registrar en un archivo de log (solo para errores)
	if l.Type == ERROR {
		logToFile(formattedMessage)
	}
}

// logToFile guarda el mensaje en un archivo de log de errores
func logToFile(message string) {
	f, err := openLogFile()
	if err != nil {
		fmt.Printf("No se pudo abrir el archivo de log: %v\n", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(message + "\n")
	if err != nil {
		fmt.Printf("No se pudo escribir en el archivo de log: %v\n", err)
	}
}

// openLogFile abre (o crea) el archivo de log de errores
func openLogFile() (*os.File, error) {
	// Crear el directorio si no existe - usar ruta absoluta
	logDir := "./logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, err
	}
	logPath := logDir + "/error.log"
	return os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

// LogInfo registra un mensaje de información
func LogInfo(comando, mensaje string) {
	logger := NewLogger(INFO, comando, mensaje)
	logger.Log()
}

// LogWarning registra un mensaje de advertencia
func LogWarning(comando, mensaje string) {
	logger := NewLogger(WARNING, comando, mensaje)
	logger.Log()
}

// LogError registra un mensaje de error
func LogError(comando, mensaje string) {
	logger := NewLogger(ERROR, comando, mensaje)
	logger.Log()
	logToFile(fmt.Sprintf("ERROR: %s - %s", comando, mensaje))
}

// LogSuccess registra un mensaje de éxito
func LogSuccess(comando, mensaje string) {
	logger := NewLogger(SUCCESS, comando, mensaje)
	logger.Log()
}

// SSEHandler maneja las conexiones Server-Sent Events para logs en tiempo real
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	// Configurar headers para SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// Crear un canal para esta conexión
	messageChan := make(chan LogMessage, 10)

	// Enviar mensajes al cliente
	for {
		select {
		case msg := <-messageChan:
			data, _ := json.Marshal(msg)
			fmt.Fprintf(w, "data: %s\n\n", data)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}

// GetLogsHandler maneja las peticiones HTTP para obtener logs (polling)
func GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	// Leer los logs del archivo (implementación básica)
	logFile := "./logs/error.log"

	// Verificar si el archivo existe
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		// Si no existe el archivo, devolver logs vacíos
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"logs":   []string{},
			"status": "success",
		})
		return
	}

	// Leer el contenido del archivo
	content, err := os.ReadFile(logFile)
	if err != nil {
		http.Error(w, "Error leyendo logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"logs":   string(content),
		"status": "success",
	})
}
