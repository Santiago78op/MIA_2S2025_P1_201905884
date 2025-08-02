#!/bin/bash
echo "🚀 Iniciando proyecto completo..."

# Función para limpiar procesos al salir
cleanup() {
    echo "🛑 Deteniendo servicios..."
    kill $FRONTEND_PID $BACKEND_PID 2>/dev/null
    exit
}

trap cleanup SIGINT SIGTERM

# Iniciar backend en segundo plano
echo "⚙️ Iniciando backend..."
(cd backend && go run main.go) &
BACKEND_PID=$!

# Esperar un poco para que el backend inicie
sleep 3

# Iniciar frontend en segundo plano
echo "🌐 Iniciando frontend..."
(cd frontend && npm start) &
FRONTEND_PID=$!

echo "✅ Servicios iniciados:"
echo "   - Backend: http://localhost:8080"
echo "   - Frontend: http://localhost:3000"
echo ""
echo "Presiona Ctrl+C para detener todos los servicios"

# Esperar a que terminen los procesos
wait $FRONTEND_PID $BACKEND_PID
