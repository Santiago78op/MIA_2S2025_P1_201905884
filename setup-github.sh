#!/bin/bash
# 🔗 Script para conectar con GitHub desde Windows
# Ejecutar en PowerShell o Git Bash

echo "🚀 Configurando repositorio para GitHub..."

# Verificar que estamos en un repositorio Git
if [ ! -d ".git" ]; then
    echo "❌ Error: No estás en un repositorio Git"
    echo "Ejecuta: git init"
    exit 1
fi

# Solicitar datos del usuario
echo "📝 Configuración de Git:"
read -p "Tu nombre: " git_name
read -p "Tu email: " git_email
read -p "Usuario de GitHub: " github_user
read -p "Nombre del repositorio: " repo_name

# Configurar Git
git config --global user.name "$git_name"
git config --global user.email "$git_email"

echo "✅ Git configurado"

# Agregar origen remoto
git remote add origin "https://github.com/$github_user/$repo_name.git"

echo "✅ Origen remoto agregado"

# Cambiar branch a main
git branch -M main

echo "📤 Subiendo código a GitHub..."

# Push inicial
git push -u origin main

if [ $? -eq 0 ]; then
    echo "🎉 ¡Código subido exitosamente!"
    echo ""
    echo "🌐 Tu repositorio: https://github.com/$github_user/$repo_name"
    echo ""
    echo "📋 Para clonar en Pop!_OS:"
    echo "   git clone https://github.com/$github_user/$repo_name.git"
    echo "   cd $repo_name"
    echo "   ./install-popos.sh"
else
    echo "❌ Error al subir código"
    echo "Verifica que el repositorio exista en GitHub"
fi
