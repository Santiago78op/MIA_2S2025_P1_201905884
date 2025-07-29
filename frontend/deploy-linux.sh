#!/bin/bash

# Script de despliegue para Linux
echo "🚀 Desplegando aplicación React en Linux..."

# Función para verificar si un comando existe
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Verificar Node.js
if ! command_exists node; then
    echo "❌ Node.js no está instalado"
    echo "Instala Node.js desde: https://nodejs.org/"
    exit 1
fi

# Verificar npm
if ! command_exists npm; then
    echo "❌ npm no está instalado"
    exit 1
fi

echo "✅ Node.js $(node --version) detectado"
echo "✅ npm $(npm --version) detectado"

# Instalar dependencias
echo "📦 Instalando dependencias..."
npm install

if [ $? -ne 0 ]; then
    echo "❌ Error instalando dependencias"
    exit 1
fi

# Opción de construcción o desarrollo
echo "Selecciona el modo de ejecución:"
echo "1) Desarrollo (npm start)"
echo "2) Construcción para producción (npm run build)"
echo "3) Servir aplicación construida"

read -p "Ingresa tu opción (1-3): " option

case $option in
    1)
        echo "🔧 Iniciando servidor de desarrollo..."
        npm start
        ;;
    2)
        echo "🏗️ Construyendo aplicación para producción..."
        npm run build
        
        if [ $? -eq 0 ]; then
            echo "✅ Aplicación construida exitosamente en ./build/"
        else
            echo "❌ Error en la construcción"
            exit 1
        fi
        ;;
    3)
        if [ ! -d "build" ]; then
            echo "🏗️ No existe carpeta build, construyendo..."
            npm run build
        fi
        
        if command_exists serve; then
            echo "🌐 Sirviendo aplicación en http://localhost:3000"
            serve -s build -l 3000
        else
            echo "📦 Instalando serve globalmente..."
            npm install -g serve
            echo "🌐 Sirviendo aplicación en http://localhost:3000"
            serve -s build -l 3000
        fi
        ;;
    *)
        echo "❌ Opción inválida"
        exit 1
        ;;
esac
