# 🚀 Guía GitHub → Pop!_OS

## 📤 Subir a GitHub (desde Windows)

### 1. Crear repositorio en GitHub
1. Ve a [github.com](https://github.com) y crea un nuevo repositorio
2. **Nombre**: `mia-proyecto-filesystem` (o el que prefieras)
3. **Descripción**: "Simulador de Sistema de Archivos EXT2 - MIA 2S2025"
4. ✅ **Público** o **Privado** (según tu preferencia)
5. ❌ **NO** marcar "Add a README file" (ya tenemos uno)
6. ❌ **NO** agregar .gitignore (ya tenemos uno)

### 2. Conectar repositorio local con GitHub
```bash
# Configurar Git (si es primera vez)
git config --global user.name "Tu Nombre"
git config --global user.email "tu.email@gmail.com"

# Agregar origen remoto (reemplaza TU-USUARIO y TU-REPOSITORIO)
git remote add origin https://github.com/TU-USUARIO/TU-REPOSITORIO.git

# Subir código
git branch -M main
git push -u origin main
```

### 3. Verificar archivos subidos
Tu repositorio debería contener:
```
📁 backend/
   └── main.go
📁 frontend/
   ├── 📁 src/
   ├── 📁 public/
   ├── package.json
   ├── Dockerfile
   └── ...
📁 docs/
📄 README.md
📄 README_POPOS.md
📄 install-popos.sh
📄 docker-compose.yml
📄 .gitignore
```

---

## 📥 Clonar en Pop!_OS

### 1. Preparar Pop!_OS
```bash
# Actualizar sistema
sudo apt update && sudo apt upgrade -y

# Instalar Git (si no está)
sudo apt install git -y
```

### 2. Clonar repositorio
```bash
# Clonar (reemplaza TU-USUARIO y TU-REPOSITORIO)
git clone https://github.com/TU-USUARIO/TU-REPOSITORIO.git

# Entrar al directorio
cd TU-REPOSITORIO

# Ver contenido
ls -la
```

### 3. Instalación automática
```bash
# Método 1: Script automatizado (RECOMENDADO)
chmod +x install-popos.sh
./install-popos.sh

# El script instala:
# ✅ Node.js y npm
# ✅ Go
# ✅ Docker (opcional)
# ✅ Dependencias del proyecto
# ✅ Scripts de ejecución
```

### 4. Instalación manual (alternativa)
```bash
# Instalar Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Instalar Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Configurar proyecto
cd frontend && npm install && cd ..
cd backend && go mod init backend && go mod tidy && cd ..
```

---

## 🚀 Ejecutar proyecto en Pop!_OS

### Opción 1: Scripts automatizados
```bash
# Todo junto
./run-all.sh

# Solo frontend
./run-frontend.sh

# Solo backend  
./run-backend.sh
```

### Opción 2: Comandos manuales
```bash
# Terminal 1: Backend
cd backend
go run main.go

# Terminal 2: Frontend  
cd frontend
npm start
```

### Opción 3: Docker
```bash
# Todo con Docker
docker-compose up

# Solo frontend
docker-compose up frontend
```

---

## 🌐 URLs de acceso

- **Frontend React**: http://localhost:3000
- **Backend Go**: http://localhost:8080

---

## 🔧 Comandos útiles en Pop!_OS

```bash
# Ver procesos activos
ps aux | grep -E "(node|go)"

# Liberar puertos
sudo kill -9 $(sudo lsof -t -i:3000)
sudo kill -9 $(sudo lsof -t -i:8080)

# Variables de entorno
export REACT_APP_API_URL=http://localhost:8080
export PORT=3001

# Logs en tiempo real
tail -f frontend.log
tail -f backend.log

# Estado del repositorio
git status
git pull origin main
```

---

## 📱 Acceso desde otros dispositivos

En Pop!_OS, para acceder desde otros dispositivos en la red:

```bash
# Obtener IP local
ip addr show | grep "inet 192.168"

# Ejecutar con IP específica
REACT_APP_HOST=0.0.0.0 npm start

# Acceso desde otros dispositivos:
# http://192.168.X.X:3000
```

---

## 🔄 Flujo de trabajo Git

```bash
# Hacer cambios y subirlos
git add .
git commit -m "Descripción de cambios"
git push origin main

# Bajar cambios desde GitHub
git pull origin main

# Ver historial
git log --oneline
```

---

## 🆘 Solución de problemas

### Puerto ocupado
```bash
# Cambiar puerto React
PORT=3001 npm start

# O crear .env
echo "PORT=3001" > frontend/.env
```

### Permisos npm
```bash
mkdir ~/.npm-global
npm config set prefix '~/.npm-global'
echo 'export PATH=~/.npm-global/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

### Problemas de CORS
Verificar que el backend incluya headers CORS apropiados.

---

## ✅ Lista de verificación

**Antes de subir a GitHub:**
- [ ] Git configurado con tu nombre y email
- [ ] .gitignore incluye node_modules/ y archivos temporales
- [ ] README.md actualizado
- [ ] Scripts de instalación probados

**Después de clonar en Pop!_OS:**
- [ ] Script install-popos.sh ejecutado sin errores
- [ ] Frontend inicia en http://localhost:3000
- [ ] Backend inicia en http://localhost:8080
- [ ] Sin errores en consola
- [ ] Variables de entorno configuradas
