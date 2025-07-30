#!/bin/bash

# 🐧 Script de instalación automática para Pop!_OS
# Ejecutar con: bash install-popos.sh

set -e  # Salir si hay errores

echo "🚀 Instalando proyecto MIA en Pop!_OS..."

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para imprimir mensajes coloridos
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

# Verificar si el comando existe
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Actualizar sistema
print_status "Actualizando sistema..."
sudo apt update && sudo apt upgrade -y
print_success "Sistema actualizado"

# Instalar Node.js y npm
if ! command_exists node; then
    print_status "Instalando Node.js y npm..."
    
    # Usar NodeSource para versión reciente
    curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
    sudo apt-get install -y nodejs
    
    print_success "Node.js $(node --version) instalado"
    print_success "npm $(npm --version) instalado"
else
    print_success "Node.js $(node --version) ya está instalado"
fi

# Instalar Go
if ! command_exists go; then
    print_status "Instalando Go..."
    
    # Descargar Go
    GO_VERSION="1.21.0"
    wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    
    # Instalar
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    
    # Agregar al PATH
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    fi
    
    export PATH=$PATH:/usr/local/go/bin
    rm go${GO_VERSION}.linux-amd64.tar.gz
    
    print_success "Go $(go version) instalado"
else
    print_success "Go ya está instalado"
fi

# Instalar Git
if ! command_exists git; then
    print_status "Instalando Git..."
    sudo apt install git -y
    print_success "Git instalado"
else
    print_success "Git ya está instalado"
fi

# Instalar Docker (opcional)
read -p "¿Instalar Docker? (y/n): " install_docker
if [[ $install_docker =~ ^[Yy]$ ]]; then
    if ! command_exists docker; then
        print_status "Instalando Docker..."
        sudo apt install docker.io docker-compose -y
        sudo usermod -aG docker $USER
        print_success "Docker instalado"
        print_warning "Debes cerrar y reabrir la terminal para usar Docker sin sudo"
    else
        print_success "Docker ya está instalado"
    fi
fi

# Configurar proyecto
if [ -f "package.json" ]; then
    print_status "Configurando proyecto..."
    
    # Frontend
    if [ -d "frontend" ]; then
        print_status "Instalando dependencias del frontend..."
        cd frontend
        npm install
        print_success "Frontend configurado"
        cd ..
    fi
    
    # Backend
    if [ -d "backend" ]; then
        print_status "Configurando backend..."
        cd backend
        
        if [ ! -f "go.mod" ]; then
            go mod init backend
        fi
        
        go mod tidy
        print_success "Backend configurado"
        cd ..
    fi
    
    print_success "Proyecto configurado exitosamente"
else
    print_warning "No se encontró package.json. Asegúrate de estar en el directorio del proyecto."
fi

# Crear scripts de ejecución
print_status "Creando scripts de ejecución..."

# Script para ejecutar frontend
cat > run-frontend.sh << 'EOF'
#!/bin/bash
echo "🌐 Iniciando frontend React..."
cd frontend
npm start
EOF

# Script para ejecutar backend
cat > run-backend.sh << 'EOF'
#!/bin/bash
echo "⚙️ Iniciando backend Go..."
cd backend
go run main.go
EOF

# Script para ejecutar ambos
cat > run-all.sh << 'EOF'
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
cd backend && go run main.go &
BACKEND_PID=$!

# Esperar un poco para que el backend inicie
sleep 3

# Iniciar frontend en segundo plano
echo "🌐 Iniciando frontend..."
cd ../frontend && npm start &
FRONTEND_PID=$!

echo "✅ Servicios iniciados:"
echo "   - Backend: http://localhost:8080"
echo "   - Frontend: http://localhost:3000"
echo ""
echo "Presiona Ctrl+C para detener todos los servicios"

# Esperar a que terminen los procesos
wait $FRONTEND_PID $BACKEND_PID
EOF

# Hacer scripts ejecutables
chmod +x run-frontend.sh run-backend.sh run-all.sh

print_success "Scripts creados:"
print_status "  - ./run-frontend.sh - Solo frontend"
print_status "  - ./run-backend.sh - Solo backend"  
print_status "  - ./run-all.sh - Ambos servicios"

echo ""
echo "🎉 ¡Instalación completada!"
echo ""
echo "📋 Próximos pasos:"
echo "   1. Para ejecutar todo: ./run-all.sh"
echo "   2. Frontend solo: ./run-frontend.sh"
echo "   3. Backend solo: ./run-backend.sh"
echo ""
echo "🌐 URLs:"
echo "   - Frontend: http://localhost:3000"
echo "   - Backend: http://localhost:8080"
echo ""
echo "📚 Documentación: README_POPOS.md"
