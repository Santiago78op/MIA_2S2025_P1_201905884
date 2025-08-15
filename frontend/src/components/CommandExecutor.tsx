import React, { useState } from 'react';
import { useCommandExecution } from '../hooks/useApi';
import { commandParser, validateAndFormatCommand } from '../services/commandParser';
import './CommandExecutor.css';

const CommandExecutor: React.FC = () => {
  const [command, setCommand] = useState('');
  const [validationError, setValidationError] = useState<string | null>(null);
  const [showHelp, setShowHelp] = useState(false);
  const [selectedHelpCommand, setSelectedHelpCommand] = useState<string>('');
  const [activeTab, setActiveTab] = useState<'command' | 'script'>('command');
  const [scriptContent, setScriptContent] = useState('');
  const [isExecutingScript, setIsExecutingScript] = useState(false);
  const [scriptResults, setScriptResults] = useState<Array<{
    command: string;
    success: boolean;
    message: string;
    data?: any;
  }>>([]);
  
  const { executeCommand, isExecuting, lastResult, error, clearError } = useCommandExecution();

  // Manejar envío de comando individual
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!command.trim()) return;

    // Validar comando antes de enviarlo
    const validation = validateAndFormatCommand(command);
    
    if (!validation.isValid) {
      const errorMessage = validation.errors.join(', ');
      setValidationError(errorMessage);
      return;
    }

    try {
      await executeCommand(validation.formattedCommand);
      setCommand(''); // Limpiar comando después de ejecutar
      setValidationError(null);
    } catch (err: any) {
      // El error ya se maneja en useCommandExecution
      console.error('Error ejecutando comando:', err);
    }
  };

  // Manejar cambios en el input
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newCommand = e.target.value;
    setCommand(newCommand);
    
    // Limpiar errores anteriores
    if (error) clearError();
    if (validationError) setValidationError(null);
  };

  // Insertar comando predefinido
  const insertPredefinedCommand = (cmd: string) => {
    setCommand(cmd);
    setValidationError(null);
    if (error) clearError();
  };

  // Manejar carga de script
  const handleFileLoad = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    if (!file.name.endsWith('.smia')) {
      alert('Por favor selecciona un archivo con extensión .smia');
      return;
    }

    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      setScriptContent(content);
      setActiveTab('script');
    };
    reader.readAsText(file);
  };

  // Parsear script
  const parseScript = (script: string): string[] => {
    return script
      .split('\n')
      .map(line => line.trim())
      .filter(line => line && !line.startsWith('#'));
  };

  // Ejecutar script
  const executeScript = async () => {
    if (!scriptContent.trim()) return;

    setIsExecutingScript(true);
    setScriptResults([]);

    const commands = parseScript(scriptContent);
    const results: Array<{
      command: string;
      success: boolean;
      message: string;
      data?: any;
    }> = [];

    for (const command of commands) {
      try {
        const validation = validateAndFormatCommand(command);
        
        if (!validation.isValid) {
          const errorMessage = `Error de validación: ${validation.errors.join(', ')}`;
          results.push({
            command: command,
            success: false,
            message: errorMessage
          });
          continue;
        }

        const response = await executeCommand(validation.formattedCommand);
        results.push({
          command: command,
          success: true,
          message: response.message,
          data: response.data
        });

        await new Promise(resolve => setTimeout(resolve, 500));
        
      } catch (error: any) {
        const errorMessage = error.message || 'Error desconocido';
        results.push({
          command: command,
          success: false,
          message: errorMessage
        });
      }

      setScriptResults([...results]);
    }

    setIsExecutingScript(false);
    
    // Resumen de ejecución del script
    const successCount = results.filter(r => r.success).length;
    const errorCount = results.filter(r => !r.success).length;
    
    console.log(`Script ejecutado: ${successCount} exitosos, ${errorCount} errores`);
  };

  // Descargar script de ejemplo
  const downloadSampleScript = () => {
    const sampleScript = `# Script de ejemplo para sistema de archivos EXT2
# Autor: MIA 2S2025

# Crear disco de 3GB
mkdisk -size=3000 -unit=M -path="/home/disk1.mia"

# Crear partición primaria de 300MB
fdisk -size=300 -unit=M -path="/home/disk1.mia" -name="Particion1"

# Montar la partición
mount -path="/home/disk1.mia" -name="Particion1"

# Formatear la partición
mkfs -id="A1" -type=full

# Iniciar sesión como root
login -user="root" -pass="123" -id="A1"

# Crear un grupo
mkgrp -name="usuarios"

# Crear un usuario
mkusr -user="juan" -pass="123" -grp="usuarios"

# Crear una carpeta
mkdir -path="/home/documentos" -p

# Crear un archivo
mkfile -path="/home/documentos/archivo.txt" -size=100

# Generar reportes
rep -id="A1" -path="/home/mbr_report.jpg" -name=mbr
rep -id="A1" -path="/home/disk_report.jpg" -name=disk

# Cerrar sesión
logout
`;

    const blob = new Blob([sampleScript], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'ejemplo.smia';
    a.click();
    URL.revokeObjectURL(url);
  };

  // Mostrar ayuda
  const handleShowHelp = (commandName?: string) => {
    setSelectedHelpCommand(commandName || '');
    setShowHelp(true);
  };

  // Renderizar ayuda
  const renderHelp = () => {
    if (selectedHelpCommand) {
      return (
        <div className="command-help">
          <pre>{commandParser.getCommandHelp(selectedHelpCommand)}</pre>
        </div>
      );
    } else {
      return (
        <div className="command-help">
          <h4>Comandos disponibles:</h4>
          <div className="command-list">
            {commandParser.listCommands().map(cmd => (
              <button
                key={cmd}
                onClick={() => setSelectedHelpCommand(cmd)}
                className="help-command-button"
              >
                {cmd}
              </button>
            ))}
          </div>
        </div>
      );
    }
  };

  // Comandos predefinidos organizados por categorías
  const predefinedCommands = [
    { 
      category: 'Administración de Discos',
      commands: [
        { label: 'Crear disco de 3GB', command: 'mkdisk -size=3000 -unit=M -path="/home/disk1.mia"' },
        { label: 'Crear disco de 1GB', command: 'mkdisk -size=1 -unit=M -path="/home/smalldisk.mia"' },
        { label: 'Eliminar disco', command: 'rmdisk -path="/home/disk1.mia"' },
      ]
    },
    {
      category: 'Particiones',
      commands: [
        { label: 'Crear partición primaria', command: 'fdisk -size=300 -unit=M -path="/home/disk1.mia" -name="Particion1"' },
        { label: 'Crear partición extendida', command: 'fdisk -size=500 -unit=M -path="/home/disk1.mia" -name="Extendida1" -type=E' },
        { label: 'Crear partición lógica', command: 'fdisk -size=100 -unit=M -path="/home/disk1.mia" -name="Logica1" -type=L' },
        { label: 'Montar partición', command: 'mount -path="/home/disk1.mia" -name="Particion1"' },
      ]
    },
    {
      category: 'Sistema de Archivos',
      commands: [
        { label: 'Formatear partición', command: 'mkfs -id="A1" -type=full' },
        { label: 'Iniciar sesión', command: 'login -user="root" -pass="123" -id="A1"' },
        { label: 'Cerrar sesión', command: 'logout' },
        { label: 'Ver particiones montadas', command: 'mounted' },
      ]
    },
    {
      category: 'Usuarios y Grupos',
      commands: [
        { label: 'Crear grupo', command: 'mkgrp -name="usuarios"' },
        { label: 'Crear usuario', command: 'mkusr -user="juan" -pass="123" -grp="usuarios"' },
        { label: 'Cambiar grupo de usuario', command: 'chgrp -user="juan" -grp="admin"' },
      ]
    },
    {
      category: 'Archivos y Carpetas',
      commands: [
        { label: 'Crear archivo', command: 'mkfile -path="/home/archivo.txt" -size=100' },
        { label: 'Crear carpeta', command: 'mkdir -path="/home/nueva_carpeta" -p' },
        { label: 'Ver contenido de archivo', command: 'cat -file1="/home/archivo.txt"' },
      ]
    },
    {
      category: 'Reportes',
      commands: [
        { label: 'Reporte MBR', command: 'rep -id="A1" -path="/home/mbr_report.jpg" -name=mbr' },
        { label: 'Reporte de disco', command: 'rep -id="A1" -path="/home/disk_report.jpg" -name=disk' },
        { label: 'Reporte de árbol', command: 'rep -id="A1" -path="/home/tree_report.jpg" -name=tree' },
      ]
    }
  ];

  return (
    <div className="command-executor">
      <div className="executor-header">
        <h3>Ejecutar Comandos</h3>
        <div className="header-actions">
          <div className="tab-buttons">
            <button 
              onClick={() => setActiveTab('command')}
              className={`tab-button ${activeTab === 'command' ? 'active' : ''}`}
            >
              💻 Comandos
            </button>
            <button 
              onClick={() => setActiveTab('script')}
              className={`tab-button ${activeTab === 'script' ? 'active' : ''}`}
            >
              📜 Scripts
            </button>
          </div>
          <button 
            onClick={() => handleShowHelp()} 
            className="help-button"
            title="Ver ayuda de comandos"
          >
            ❓ Ayuda
          </button>
        </div>
      </div>

      {/* Contenido de las pestañas */}
      {activeTab === 'command' ? (
        <>
          <form onSubmit={handleSubmit} className="command-form">
            <div className="input-group">
              <input
                type="text"
                value={command}
                onChange={handleInputChange}
                placeholder="Ingrese un comando (ej: mkdisk -size=1000 -unit=M -path=/home/disk1.mia)"
                className={`command-input ${validationError || error ? 'error' : ''}`}
                disabled={isExecuting}
              />
              <button
                type="submit"
                disabled={isExecuting || !command.trim() || !!validationError}
                className="execute-button"
              >
                {isExecuting ? '⏳' : '▶️'} 
                {isExecuting ? 'Ejecutando...' : 'Ejecutar'}
              </button>
            </div>
            
            {validationError && (
              <div className="validation-error">
                <span>⚠️ {validationError}</span>
              </div>
            )}
          </form>

          {/* Comandos predefinidos */}
          <div className="predefined-commands">
            <h4>Comandos de ejemplo:</h4>
            {predefinedCommands.map((category, categoryIndex) => (
              <div key={categoryIndex} className="command-category">
                <h5>{category.category}</h5>
                <div className="command-buttons">
                  {category.commands.map((cmd, index) => (
                    <button
                      key={index}
                      onClick={() => insertPredefinedCommand(cmd.command)}
                      className="predefined-button"
                      disabled={isExecuting}
                      title={cmd.command}
                    >
                      {cmd.label}
                    </button>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </>
      ) : (
        // Pestaña de Scripts
        <div className="script-executor">
          <div className="script-header">
            <h4>📜 Ejecutor de Scripts (.smia)</h4>
            <div className="script-actions">
              <button 
                onClick={downloadSampleScript}
                className="download-sample-button"
                title="Descargar script de ejemplo"
              >
                📥 Ejemplo
              </button>
              <label className="load-script-button">
                📁 Cargar Script
                <input
                  type="file"
                  accept=".smia"
                  onChange={handleFileLoad}
                  style={{ display: 'none' }}
                />
              </label>
              <button 
                onClick={executeScript}
                disabled={!scriptContent.trim() || isExecutingScript}
                className="execute-script-button"
              >
                {isExecutingScript ? '⏳ Ejecutando...' : '▶️ Ejecutar Script'}
              </button>
            </div>
          </div>

          {scriptContent && (
            <div className="script-content">
              <h5>Contenido del Script:</h5>
              <textarea
                value={scriptContent}
                onChange={(e) => setScriptContent(e.target.value)}
                className="script-textarea"
                placeholder="Contenido del script aparecerá aquí..."
                rows={10}
              />
              <div className="script-info">
                <small>
                  Líneas totales: {scriptContent.split('\n').length} | 
                  Comandos: {parseScript(scriptContent).length} | 
                  Comentarios: {scriptContent.split('\n').filter(line => line.trim().startsWith('#')).length}
                </small>
              </div>
            </div>
          )}

          {scriptResults.length > 0 && (
            <div className="script-results">
              <h5>Resultados de la Ejecución:</h5>
              <div className="results-list">
                {scriptResults.map((result, index) => (
                  <div key={index} className={`result-item ${result.success ? 'success' : 'error'}`}>
                    <div className="result-command">
                      <code>{result.command}</code>
                    </div>
                    <div className="result-message">
                      <span className={`result-icon ${result.success ? 'success' : 'error'}`}>
                        {result.success ? '✅' : '❌'}
                      </span>
                      {result.message}
                    </div>
                    {result.data && (
                      <details className="result-details">
                        <summary>Ver datos</summary>
                        <pre>{JSON.stringify(result.data, null, 2)}</pre>
                      </details>
                    )}
                  </div>
                ))}
              </div>
              
              <div className="execution-summary">
                <div className="summary-stats">
                  <span className="stat success">
                    ✅ Exitosos: {scriptResults.filter(r => r.success).length}
                  </span>
                  <span className="stat error">
                    ❌ Errores: {scriptResults.filter(r => !r.success).length}
                  </span>
                  <span className="stat total">
                    📊 Total: {scriptResults.length}
                  </span>
                </div>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Sección de ayuda */}
      {showHelp && (
        <div className="help-section">
          <div className="help-header">
            <h4>
              {selectedHelpCommand ? `Ayuda: ${selectedHelpCommand}` : 'Ayuda de Comandos'}
            </h4>
            <div className="help-actions">
              {selectedHelpCommand && (
                <button 
                  onClick={() => setSelectedHelpCommand('')}
                  className="back-button"
                >
                  ← Volver
                </button>
              )}
              <button 
                onClick={() => setShowHelp(false)}
                className="close-help-button"
              >
                ✕ Cerrar
              </button>
            </div>
          </div>
          {renderHelp()}
        </div>
      )}

      {/* Resultados de comandos individuales */}
      {lastResult && (
        <div className="command-result success">
          <h4>✅ Resultado:</h4>
          <div className="result-content">
            <p><strong>Mensaje:</strong> {lastResult.message}</p>
            {lastResult.data && (
              <div className="result-data">
                <strong>Datos:</strong>
                <pre>{JSON.stringify(lastResult.data, null, 2)}</pre>
              </div>
            )}
            <p><strong>Estado:</strong> <span className="status-badge success">{lastResult.status}</span></p>
          </div>
        </div>
      )}

      {/* Errores */}
      {error && (
        <div className="command-result error">
          <h4>❌ Error:</h4>
          <p>{error}</p>
          <button onClick={clearError} className="clear-error-button">
            Limpiar error
          </button>
        </div>
      )}
    </div>
  );
};

export default CommandExecutor;