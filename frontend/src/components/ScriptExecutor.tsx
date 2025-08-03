// components/ScriptExecutor.tsx
import React, { useState, useRef } from 'react';
import { useCommandExecution } from '../hooks/useApi';
import { validateAndFormatCommand } from '../services/commandParser';

interface ScriptExecutorProps {
  onScriptLoad: (script: string) => void;
}

const ScriptExecutor: React.FC<ScriptExecutorProps> = ({ onScriptLoad }) => {
  const [scriptContent, setScriptContent] = useState('');
  const [isExecutingScript, setIsExecutingScript] = useState(false);
  const [scriptResults, setScriptResults] = useState<Array<{
    command: string;
    success: boolean;
    message: string;
    data?: any;
  }>>([]);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { executeCommand } = useCommandExecution();

  const handleFileLoad = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    // Verificar extensión del archivo
    if (!file.name.endsWith('.smia')) {
      alert('Por favor selecciona un archivo con extensión .smia');
      return;
    }

    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      setScriptContent(content);
      onScriptLoad(content);
    };
    reader.readAsText(file);
  };

  const parseScript = (script: string): string[] => {
    return script
      .split('\n')
      .map(line => line.trim())
      .filter(line => line && !line.startsWith('#')); // Filtrar líneas vacías y comentarios
  };

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
        // Validar comando antes de ejecutar
        const validation = validateAndFormatCommand(command);
        
        if (!validation.isValid) {
          results.push({
            command: command,
            success: false,
            message: `Error de validación: ${validation.errors.join(', ')}`
          });
          continue;
        }

        // Ejecutar comando
        const response = await executeCommand(validation.formattedCommand);
        results.push({
          command: command,
          success: true,
          message: response.message,
          data: response.data
        });

        // Pequeña pausa entre comandos para evitar sobrecarga
        await new Promise(resolve => setTimeout(resolve, 500));
        
      } catch (error) {
        results.push({
          command: command,
          success: false,
          message: error instanceof Error ? error.message : 'Error desconocido'
        });
      }

      // Actualizar resultados en tiempo real
      setScriptResults([...results]);
    }

    setIsExecutingScript(false);
  };

  const downloadSampleScript = () => {
    const sampleScript = `# Script de ejemplo para sistema de archivos EXT2
# Autor: MIA 2S2025
# Descripción: Script básico para crear y configurar un sistema de archivos

# Crear disco de 3GB
mkdisk -size=3000 -unit=M -path="/home/disk1.mia"

# Crear partición primaria de 300MB
fdisk -size=300 -unit=M -path="/home/disk1.mia" -name="Particion1"

# Crear partición extendida de 500MB
fdisk -size=500 -unit=M -path="/home/disk1.mia" -name="Extendida1" -type=E

# Crear partición lógica de 100MB
fdisk -size=100 -unit=M -path="/home/disk1.mia" -name="Logica1" -type=L

# Montar la primera partición
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
rep -id="A1" -path="/home/tree_report.jpg" -name=tree

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

  return (
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
          <button 
            onClick={() => fileInputRef.current?.click()}
            className="load-script-button"
            disabled={isExecutingScript}
          >
            📁 Cargar Script
          </button>
          <button 
            onClick={executeScript}
            disabled={!scriptContent.trim() || isExecutingScript}
            className="execute-script-button"
          >
            {isExecutingScript ? '⏳ Ejecutando...' : '▶️ Ejecutar Script'}
          </button>
        </div>
      </div>

      <input
        ref={fileInputRef}
        type="file"
        accept=".smia"
        onChange={handleFileLoad}
        style={{ display: 'none' }}
      />

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
  );
};

export default ScriptExecutor;