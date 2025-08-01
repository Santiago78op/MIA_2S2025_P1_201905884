# 🗺️ Roadmap de Desarrollo - Sistema de Archivos EXT2

## 📅 **Cronograma General**

**🎯 Fecha de Entrega:** 7 de septiembre de 2025, 23:59 horas
**⏰ Tiempo Restante:** ~5 semanas

---

## ✅ **COMPLETADO** 

### **Fase 0: Infraestructura Base** ✅ **100% COMPLETADO**
- [x] Configuración del proyecto Go
- [x] Frontend React + TypeScript
- [x] Integración Frontend-Backend
- [x] Sistema de logging con WebSockets/SSE
- [x] API REST completa
- [x] Estructuras MBR y Partition
- [x] Comando `mkdisk` funcional
- [x] Documentación inicial

**🕐 Completado:** Semana 1 (1-7 agosto 2025)

---

## 🔄 **EN DESARROLLO**

### **Fase 1: Gestión Completa de Discos** 🔄 **25% COMPLETADO**

#### **Comandos por Implementar:**
- [ ] `rmdisk` - Eliminar discos virtuales
- [ ] `fdisk` - Gestión de particiones (crear, eliminar, mostrar)
- [ ] `mount` - Montar particiones en el sistema
- [ ] `unmount` - Desmontar particiones

#### **Estructuras por Implementar:**
- [ ] EBR (Extended Boot Record) para particiones lógicas
- [ ] Sistema de montaje de particiones
- [ ] Validaciones de integridad de disco

**🎯 Meta:** 31 de agosto de 2025
**⏱️ Tiempo estimado:** 1.5 semanas

---

## 📋 **PENDIENTE**

### **Fase 2: Sistema de Archivos EXT2** 🔄 **0% COMPLETADO**

#### **Comandos:**
- [ ] `mkfs` - Formatear partición con EXT2
- [ ] `login` - Autenticación de usuarios
- [ ] `logout` - Cerrar sesión
- [ ] `loss` - Simulador de pérdida de datos

#### **Estructuras EXT2:**
- [ ] SuperBloque
- [ ] Tabla de Inodos
- [ ] Bloques de datos
- [ ] Journaling
- [ ] Bitmap de Inodos
- [ ] Bitmap de Bloques

**🎯 Meta:** 7 de septiembre de 2025
**⏱️ Tiempo estimado:** 1.5 semanas

---

### **Fase 3: Gestión de Usuarios y Grupos** 🔄 **0% COMPLETADO**

#### **Comandos:**
- [ ] `mkgrp` - Crear grupos
- [ ] `rmgrp` - Eliminar grupos
- [ ] `mkusr` - Crear usuarios
- [ ] `rmusr` - Eliminar usuarios
- [ ] `chgrp` - Cambiar grupo de archivos/directorios

#### **Estructuras:**
- [ ] Archivo users.txt
- [ ] Sistema de permisos
- [ ] Validaciones de autenticación

**🎯 Meta:** 5 de septiembre de 2025
**⏱️ Tiempo estimado:** 1 semana

---

### **Fase 4: Gestión de Archivos y Directorios** 🔄 **0% COMPLETADO**

#### **Comandos:**
- [ ] `mkdir` - Crear directorios
- [ ] `mkfile` - Crear archivos
- [ ] `cat` - Mostrar contenido de archivos
- [ ] `remove` - Eliminar archivos/directorios
- [ ] `copy` - Copiar archivos
- [ ] `move` - Mover archivos
- [ ] `rename` - Renombrar archivos
- [ ] `find` - Buscar archivos
- [ ] `chown` - Cambiar propietario
- [ ] `chmod` - Cambiar permisos

**🎯 Meta:** 6 de septiembre de 2025
**⏱️ Tiempo estimado:** 1.5 semanas

---

### **Fase 5: Reportes y Visualización** 🔄 **0% COMPLETADO**

#### **Comandos de Reporte:**
- [ ] `rep` - Generar reportes

#### **Tipos de Reporte:**
- [ ] `mbr` - Estructura del MBR
- [ ] `disk` - Estado del disco
- [ ] `inode` - Información de inodos
- [ ] `journaling` - Estado del journal
- [ ] `block` - Contenido de bloques
- [ ] `bm_inode` - Bitmap de inodos
- [ ] `bm_block` - Bitmap de bloques
- [ ] `tree` - Árbol de directorios
- [ ] `sb` - SuperBloque
- [ ] `file` - Contenido de archivo
- [ ] `ls` - Listado de archivos

#### **Integración con Graphviz:**
- [ ] Generación de gráficos DOT
- [ ] Visualización en frontend
- [ ] Exportación de reportes

**🎯 Meta:** 7 de septiembre de 2025
**⏱️ Tiempo estimado:** 1 semana

---

## 📊 **Progreso General**

### **Estado por Componente:**

| Componente | Estado | Progreso |
|------------|--------|----------|
| 🏗️ **Infraestructura** | ✅ Completado | 100% |
| 💾 **Gestión de Discos** | 🔄 En desarrollo | 25% |
| 📁 **Sistema EXT2** | ⏳ Pendiente | 0% |
| 👥 **Usuarios/Grupos** | ⏳ Pendiente | 0% |
| 📄 **Archivos/Dirs** | ⏳ Pendiente | 0% |
| 📊 **Reportes** | ⏳ Pendiente | 0% |

### **Progreso Total:** 🔄 **20% COMPLETADO**

---

## 🎯 **Prioridades Semanales**

### **Semana 2 (8-14 agosto):** Completar Gestión de Discos
- [ ] Implementar `rmdisk`
- [ ] Implementar `fdisk` básico
- [ ] Implementar `mount`/`unmount`
- [ ] Estructuras EBR

### **Semana 3 (15-21 agosto):** Sistema EXT2 Core
- [ ] SuperBloque y estructuras básicas
- [ ] `mkfs` básico
- [ ] Sistema de autenticación

### **Semana 4 (22-28 agosto):** Archivos y Usuarios
- [ ] Gestión de usuarios/grupos
- [ ] Comandos básicos de archivos
- [ ] Sistema de permisos

### **Semana 5 (29 agosto - 5 septiembre):** Comandos Avanzados
- [ ] Comandos avanzados de archivos
- [ ] Validaciones y testing
- [ ] Optimizaciones

### **Semana 6 (6-7 septiembre):** Reportes y Entrega
- [ ] Sistema de reportes
- [ ] Documentación final
- [ ] Testing y entrega

---

## 🚨 **Riesgos y Mitigaciones**

### **Riesgos Identificados:**
1. **⏰ Tiempo insuficiente** para reportes complejos
   - **Mitigación:** Implementar reportes básicos primero
   
2. **🔧 Complejidad de EXT2** mayor a la esperada
   - **Mitigación:** Simplificar implementación inicial
   
3. **🐛 Bugs en integración** frontend-backend
   - **Mitigación:** Testing continuo desde ahora

### **Plan de Contingencia:**
- **Prioridad 1:** Comandos básicos funcionando
- **Prioridad 2:** Sistema EXT2 básico
- **Prioridad 3:** Reportes simplificados
- **Prioridad 4:** Features avanzadas

---

## 📈 **Métricas de Éxito**

### **Requisitos Mínimos (70%):**
- [x] ✅ Aplicación web funcional
- [x] ✅ Creación de particiones 
- [ ] 🔄 Ejecución completa de script
- [ ] ⏳ Reportes básicos
- [x] ✅ Documentación

### **Requisitos Deseables (85%):**
- [ ] ⏳ Todos los comandos implementados
- [ ] ⏳ Sistema EXT2 completo
- [ ] ⏳ Reportes con Graphviz

### **Requisitos Excepcionales (95%+):**
- [ ] ⏳ Interface web avanzada
- [ ] ⏳ Optimizaciones de rendimiento
- [ ] ⏳ Features adicionales
