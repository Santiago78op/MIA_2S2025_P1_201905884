# 🐧 Instalación en Pop!_OS

## Prerrequisitos

### 1. Actualizar el sistema
```bash
sudo apt update && sudo apt upgrade -y
```

### 2. Instalar Node.js y npm
```bash
# Opción 1: Desde repositorios de Ubuntu (más simple)
sudo apt install nodejs npm -y

# Opción 2: Usar NodeSource para versión más reciente (recomendado)
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Verificar instalación
node --version
npm --version
```

### 3. Instalar Go (para el backend)
```bash
# Descargar Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# Extraer e instalar
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Agregar al PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verificar
go version
```

### 4. Instalar Git (si no está instalado)
```bash
sudo apt install git -y
```

## 🚀 Instalación del Proyecto

### 1. Clonar el repositorio
```bash
git clone https://github.com/TU-USUARIO/TU-REPOSITORIO.git
cd TU-REPOSITORIO
```

### 2. Configurar Frontend (React)
```bash
cd frontend

# Instalar dependencias
npm install

# Iniciar servidor de desarrollo
npm start
# La aplicación estará disponible en http://localhost:3000
```

### 3. Configurar Backend (Go)
```bash
# En otra terminal
cd backend

# Inicializar módulo Go (si es necesario)
go mod init backend

# Instalar dependencias
go mod tidy

# Ejecutar servidor
go run main.go
# El servidor estará disponible en http://localhost:8080
```

## 🐳 Alternativa con Docker

Si prefieres usar Docker:

```bash
# Instalar Docker
sudo apt install docker.io docker-compose -y
sudo usermod -aG docker $USER
# Cerrar y reabrir terminal

# Ejecutar proyecto completo
docker-compose up
```

## 🔧 Comandos útiles en Pop!_OS

### Gestión de procesos
```bash
# Ver procesos de Node.js
ps aux | grep node

# Matar proceso en puerto específico
sudo kill -9 $(sudo lsof -t -i:3000)

# Ver puertos ocupados
sudo netstat -tlnp | grep :3000
```

### Variables de entorno
```bash
# Establecer variables temporales
export REACT_APP_API_URL=http://localhost:8080
export PORT=3001

# Para que persistan, agregar a ~/.bashrc
echo 'export REACT_APP_API_URL=http://localhost:8080' >> ~/.bashrc
source ~/.bashrc
```

### Ejecución en segundo plano
```bash
# Frontend en segundo plano
cd frontend
nohup npm start > frontend.log 2>&1 &

# Backend en segundo plano  
cd backend
nohup go run main.go > backend.log 2>&1 &

# Ver logs
tail -f frontend.log
tail -f backend.log
```

## 🌐 Construcción para Producción

```bash
# Frontend - Construir aplicación optimizada
cd frontend
npm run build

# Servir con servidor estático
npm install -g serve
serve -s build -l 3000

# Backend - Construir binario
cd backend
go build -o app main.go
./app
```

## 🔍 Solución de Problemas

### Problema: Permisos de npm
```bash
# Configurar directorio global de npm
mkdir ~/.npm-global
npm config set prefix '~/.npm-global'
echo 'export PATH=~/.npm-global/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

### Problema: Puerto ocupado
```bash
# Cambiar puerto de React
PORT=3001 npm start

# O crear archivo .env en frontend/
echo "PORT=3001" > frontend/.env
```

### Problema: CORS en desarrollo
El backend debe incluir headers CORS para permitir requests desde el frontend.

## 📝 Notas específicas para Pop!_OS

- Pop!_OS viene con **Pop!_Shop** donde puedes instalar Node.js gráficamente
- Usa **Super** + **T** para abrir terminal rápidamente
- **Ctrl + Alt + T** también abre terminal
- Pop!_OS incluye **Flatpak** por defecto
- Considera usar **VS Code** desde Pop!_Shop para desarrollo

## 🎯 Flujo de trabajo recomendado

1. **Desarrollo**: Usar `npm start` y `go run main.go`
2. **Testing**: Usar los comandos de prueba específicos
3. **Producción**: Construir con `npm run build` y `go build`
4. **Despliegue**: Usar Docker o binarios compilados
