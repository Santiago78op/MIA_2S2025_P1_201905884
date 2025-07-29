# Ejecutar el proyecto React en Linux

## Prerrequisitos
Asegúrate de tener instalado en Linux:
```bash
# Verificar Node.js
node --version

# Verificar npm
npm --version

# Si no están instalados, instalar Node.js (incluye npm):
# Ubuntu/Debian:
sudo apt update
sudo apt install nodejs npm

# CentOS/RHEL/Fedora:
sudo dnf install nodejs npm
# o
sudo yum install nodejs npm

# Arch Linux:
sudo pacman -S nodejs npm
```

## Opción 1: Si transfieres el proyecto completo
```bash
# Navegar al directorio del proyecto
cd /ruta/a/tu/proyecto/frontend

# Instalar dependencias (si node_modules no se transfirió)
npm install

# Iniciar el servidor de desarrollo
npm start

# O ejecutar en segundo plano
nohup npm start &

# Para producción, construir la aplicación
npm run build

# Servir archivos estáticos (usando serve)
npm install -g serve
serve -s build -l 3000
```

## Opción 2: Clonar solo el código fuente
Si solo transfieres el código fuente (sin node_modules):

```bash
# Copiar archivos: package.json, package-lock.json, src/, public/, tsconfig.json

# Instalar dependencias
npm install

# Iniciar desarrollo
npm start
```

## Opción 3: Usando Docker (recomendado para producción)
```bash
# Crear Dockerfile en la carpeta frontend
# Ejecutar contenedor
docker build -t react-app .
docker run -p 3000:3000 react-app
```

## Comandos útiles en Linux:
```bash
# Ver procesos de Node
ps aux | grep node

# Matar proceso en puerto 3000
sudo lsof -t -i tcp:3000 | xargs kill -9

# Ejecutar en puerto específico
PORT=8080 npm start

# Variables de entorno
REACT_APP_API_URL=http://localhost:8000 npm start
```

## Diferencias principales:
- En Linux usas `/` en lugar de `\` para rutas
- Comandos de terminal son bash/zsh en lugar de PowerShell
- Permisos de archivos pueden requerir `sudo` en algunos casos
- Variables de entorno se setean con `export VARIABLE=valor`
