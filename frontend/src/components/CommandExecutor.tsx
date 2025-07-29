import React, { useState } from 'react';
import { useCommandExecution } from '../hooks/useApi';
import './CommandExecutor.css';

const CommandExecutor: React.FC = () => {
  const [command, setCommand] = useState('');
  const { executeCommand, isExecuting, lastResult, error, clearError } = useCommandExecution();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!command.trim()) return;

    try {
      await executeCommand(command.trim());
      setCommand(''); // Limpiar comando después de ejecutar
    } catch (err) {
      // Error ya manejado en el hook
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setCommand(e.target.value);
    if (error) clearError();
  };

  const predefinedCommands = [
    { label: 'mkdisk -size=3000 -unit=M -path=/home/disk1.dk', command: 'mkdisk -size=3000 -unit=M -path=/home/disk1.dk' },
    { label: 'fdisk -size=300 -unit=M -path=/home/disk1.dk -name=Particion1', command: 'fdisk -size=300 -unit=M -path=/home/disk1.dk -name=Particion1' },
    { label: 'mount -path=/home/disk1.dk -name=Particion1', command: 'mount -path=/home/disk1.dk -name=Particion1' },
    { label: 'rep -id=A1 -path=/home/reporte.jpg -name=mbr', command: 'rep -id=A1 -path=/home/reporte.jpg -name=mbr' },
  ];

  return (
    <div className="command-executor">
      <h3>Ejecutar Comandos</h3>
      
      <form onSubmit={handleSubmit} className="command-form">
        <div className="input-group">
          <input
            type="text"
            value={command}
            onChange={handleInputChange}
            placeholder="Ingrese un comando (ej: mkdisk -size=1000 -unit=M -path=/home/disk1.dk)"
            className={`command-input ${error ? 'error' : ''}`}
            disabled={isExecuting}
          />
          <button
            type="submit"
            disabled={isExecuting || !command.trim()}
            className="execute-button"
          >
            {isExecuting ? '⏳' : '▶️'} 
            {isExecuting ? 'Ejecutando...' : 'Ejecutar'}
          </button>
        </div>
      </form>

      {/* Comandos predefinidos */}
      <div className="predefined-commands">
        <h4>Comandos de ejemplo:</h4>
        <div className="command-buttons">
          {predefinedCommands.map((cmd, index) => (
            <button
              key={index}
              onClick={() => setCommand(cmd.command)}
              className="predefined-button"
              disabled={isExecuting}
            >
              {cmd.label}
            </button>
          ))}
        </div>
      </div>

      {/* Resultado */}
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

      {/* Error */}
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
