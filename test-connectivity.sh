#!/bin/bash
# 🧪 Script para probar la conectividad backend-frontend

echo "🧪 Probando conectividad Backend ↔ Frontend"
echo "=============================================="

# Verificar que el backend esté ejecutándose
echo "🔍 Verificando backend en http://localhost:8080..."

# Probar endpoint de salud
echo "📡 Probando endpoint /api/health..."
health_response=$(curl -s -w "%{http_code}" http://localhost:8080/api/health)
http_code="${health_response: -3}"

if [ "$http_code" = "200" ]; then
    echo "✅ Backend funcionando correctamente"
    echo "📄 Respuesta: ${health_response%???}"
else
    echo "❌ Backend no responde (HTTP $http_code)"
    echo "🔧 Asegúrate de que el backend esté ejecutándose:"
    echo "   cd backend && go run main.go"
    exit 1
fi

echo ""

# Probar endpoint de sistemas de archivos
echo "📡 Probando endpoint /api/filesystems..."
filesystems_response=$(curl -s -w "%{http_code}" http://localhost:8080/api/filesystems)
http_code="${filesystems_response: -3}"

if [ "$http_code" = "200" ]; then
    echo "✅ Endpoint de sistemas de archivos funcionando"
    echo "📄 Respuesta: ${filesystems_response%???}"
else
    echo "❌ Endpoint de sistemas de archivos no responde (HTTP $http_code)"
fi

echo ""

# Probar endpoint de ejecución de comandos
echo "📡 Probando endpoint /api/execute..."
command_response=$(curl -s -w "%{http_code}" -X POST \
  -H "Content-Type: application/json" \
  -d '{"command":"test command"}' \
  http://localhost:8080/api/execute)
http_code="${command_response: -3}"

if [ "$http_code" = "200" ]; then
    echo "✅ Endpoint de ejecución de comandos funcionando"
    echo "📄 Respuesta: ${command_response%???}"
else
    echo "❌ Endpoint de ejecución de comandos no responde (HTTP $http_code)"
fi

echo ""

# Verificar CORS
echo "🔐 Verificando configuración CORS..."
cors_response=$(curl -s -I -X OPTIONS \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: GET" \
  http://localhost:8080/api/health)

if echo "$cors_response" | grep -q "Access-Control-Allow-Origin"; then
    echo "✅ CORS configurado correctamente"
else
    echo "⚠️  CORS podría no estar configurado correctamente"
fi

echo ""
echo "🌐 URLs para acceder:"
echo "   📋 Backend API: http://localhost:8080/api/health"
echo "   🌍 Frontend: http://localhost:3000"
echo ""
echo "🚀 Para iniciar los servicios:"
echo "   Backend: cd backend && go run main.go"
echo "   Frontend: cd frontend && npm start"
