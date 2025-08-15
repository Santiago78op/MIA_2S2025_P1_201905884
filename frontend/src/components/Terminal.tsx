import React, { useState, useRef, useEffect, KeyboardEvent } from 'react';
import { useCommandExecution } from '../hooks/useApi';
import { validateAndFormatCommand } from '../services/commandParser';
import './Terminal.css';

interface TerminalProps {
  onCommandExecuted?: (command: string, result: any) => void;
  onError?: (error: string) => void;
}

interface CommandHistoryItem {
  command: string;
  timestamp: string;
  success: boolean;
  output: string;
}

const Terminal: React.FC<TerminalProps> = ({ onCommandExecuted, onError }) => {
  const [currentCommand, setCurrentCommand] = useState('');
  const [commandHistory, setCommandHistory] = useState<CommandHistoryItem[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);
  const [isVisible, setIsVisible] = useState(true);
  const [autoComplete, setAutoComplete] = useState<string[]>([]);
  const [showAutoComplete, setShowAutoComplete] = useState(false);
  
  const inputRef = useRef<HTMLInputElement>(null);
  const terminalRef = useRef<HTMLDivElement>(null);
  
  const { executeCommand, isExecuting } = useCommandExecution();

  // Comandos disponibles para autocompletado
  const availableCommands = [
    'mkdisk', 'rmdisk', 'fdisk', 'mount', 'unmount', 'mounted',
    'mkfs', 'login', 'logout', 'mkgrp', 'rmgrp', 'mkusr', 'rmusr',
    'mkdir', 'mkfile', 'cat', 'rep', 'clear', 'help', 'exit'
  ];

  // Parámetros comunes para autocompletado
  const commonParams = [
    '-size=', '-unit=', '-path=', '-name=', '-type=', '-fit=',
    '-user=', '-pass=', '-id=', '-grp=', '-file1=', '-p'
  ];

  // Enfocar terminal al montar
  useEffect(() => {
    if (inputRef.current && isVisible) {
      inputRef.current.focus();
    }
  }, [isVisible]);

  // Scroll automático al agregar nuevos comandos
  useEffect(() => {
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  }, [commandHistory]);

  // Manejar envío de comando
  const handleSubmitCommand = async () => {
    if (!currentCommand.trim() || isExecuting) return;

    const command = currentCommand.trim();
    const timestamp = new Date().toLocaleTimeString();

    // Comandos especiales de terminal
    if (command === 'clear') {
      setCommandHistory([]);
      setCurrentCommand('');
      return;
    }

    if (command === 'help') {
      const helpOutput = `Comandos disponibles:
  mkdisk    - Crear disco virtual
  rmdisk    - Eliminar disco
  fdisk     - Gestionar particiones
  mount     - Montar partición
  unmount   - Desmontar partición
  mounted   - Mostrar particiones montadas
  mkfs      - Formatear partición
  login     - Iniciar sesión
  logout    - Cerrar sesión
  mkgrp     - Crear grupo
  rmgrp     - Eliminar grupo
  mkusr     - Crear usuario
  rmusr     - Eliminar usuario
  mkdir     - Crear directorio
  mkfile    - Crear archivo
  cat       - Mostrar contenido
  rep       - Generar reportes
  clear     - Limpiar terminal
  help      - Mostrar esta ayuda
  exit      - Minimizar terminal

Ejemplo: mkdisk -size=1000 -unit=M -path="/home/disk1.mia"`;

      const historyItem: CommandHistoryItem = {
        command,
        timestamp,
        success: true,
        output: helpOutput
      };

      setCommandHistory(prev => [...prev, historyItem]);
      setCurrentCommand('');
      return;
    }

    if (command === 'exit') {
      setIsVisible(false);
      setCurrentCommand('');
      return;
    }

    // Agregar comando al historial antes de ejecutar
    const historyItem: CommandHistoryItem = {
      command,
      timestamp,
      success: false,
      output: 'Ejecutando...'
    };

    setCommandHistory(prev => [...prev, historyItem]);
    setCurrentCommand('');
    setHistoryIndex(-1);

    try {
      // Validar comando
      const validation = validateAndFormatCommand(command);
      
      if (!validation.isValid) {
        const errorMessage = `Error de validación: ${validation.errors.join(', ')}`;
        updateLastHistoryItem(false, errorMessage);
        onError?.(errorMessage);
        return;
      }

      // Ejecutar comando
      const result = await executeCommand(validation.formattedCommand);
      
      const output = typeof result.message === 'string' 
        ? result.message 
        : JSON.stringify(result, null, 2);

      updateLastHistoryItem(true, output);
      onCommandExecuted?.(command, result);

    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Error desconocido';
      updateLastHistoryItem(false, errorMessage);
      onError?.(errorMessage);
    }
  };

  // Actualizar último elemento del historial
  const updateLastHistoryItem = (success: boolean, output: string) => {
    setCommandHistory(prev => {
      const newHistory = [...prev];
      if (newHistory.length > 0) {
        newHistory[newHistory.length - 1] = {
          ...newHistory[newHistory.length - 1],
          success,
          output
        };
      }
      return newHistory;
    });
  };

  // Manejar teclas especiales
  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleSubmitCommand();
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      navigateHistory('up');
    } else if (e.key === 'ArrowDown') {
      e.preventDefault();
      navigateHistory('down');
    } else if (e.key === 'Tab') {
      e.preventDefault();
      handleAutoComplete();
    } else if (e.key === 'Escape') {
      setShowAutoComplete(false);
      setAutoComplete([]);
    }
  };

  // Navegar en el historial de comandos
  const navigateHistory = (direction: 'up' | 'down') => {
    const history = commandHistory.map(item => item.command);
    
    if (history.length === 0) return;

    let newIndex = historyIndex;
    
    if (direction === 'up') {
      newIndex = historyIndex === -1 ? history.length - 1 : Math.max(0, historyIndex - 1);
    } else {
      newIndex = historyIndex === -1 ? -1 : Math.min(history.length - 1, historyIndex + 1);
      if (newIndex === history.length - 1 && historyIndex === history.length - 1) {
        newIndex = -1;
      }
    }

    setHistoryIndex(newIndex);
    setCurrentCommand(newIndex === -1 ? '' : history[newIndex]);
  };

  // Autocompletado
  const handleAutoComplete = () => {
    const words = currentCommand.split(' ');
    const currentWord = words[words.length - 1];

    if (words.length === 1) {
      // Autocompletar comando
      const matches = availableCommands.filter(cmd => 
        cmd.startsWith(currentWord.toLowerCase())
      );
      
      if (matches.length === 1) {
        setCurrentCommand(matches[0] + ' ');
      } else if (matches.length > 1) {
        setAutoComplete(matches);
        setShowAutoComplete(true);
      }
    } else {
      // Autocompletar parámetros
      const matches = commonParams.filter(param => 
        param.startsWith(currentWord)
      );
      
      if (matches.length === 1) {
        words[words.length - 1] = matches[0];
        setCurrentCommand(words.join(' '));
      } else if (matches.length > 1) {
        setAutoComplete(matches);
        setShowAutoComplete(true);
      }
    }
  };

  // Seleccionar sugerencia de autocompletado
  const selectAutoComplete = (suggestion: string) => {
    const words = currentCommand.split(' ');
    words[words.length - 1] = suggestion + (suggestion.endsWith('=') ? '' : ' ');
    setCurrentCommand(words.join(' '));
    setShowAutoComplete(false);
    setAutoComplete([]);
    inputRef.current?.focus();
  };

  // Obtener prompt del sistema
  const getPrompt = () => {
    return `mia@ext2-simulator:~$ `;
  };

  if (!isVisible) {
    return (
      <div className="terminal-minimized" onClick={() => setIsVisible(true)}>
        <span>📟 Terminal (click para abrir)</span>
      </div>
    );
  }

  return (
    <div className="terminal-container">
      <div className="terminal-header">
        <div className="terminal-title">
          <span className="terminal-icon">📟</span>
          <span>Terminal MIA</span>
        </div>
        <div className="terminal-controls">
          <button 
            className="terminal-control minimize"
            onClick={() => setIsVisible(false)}
            title="Minimizar"
          >
            —
          </button>
          <button 
            className="terminal-control clear"
            onClick={() => setCommandHistory([])}
            title="Limpiar terminal"
          >
            🗑️
          </button>
        </div>
      </div>

      <div className="terminal-body" ref={terminalRef}>
        <div className="terminal-welcome">
          <div>Bienvenido al Terminal del Sistema de Archivos EXT2</div>
          <div>Escribe 'help' para ver comandos disponibles</div>
          <div>Usa Tab para autocompletado, ↑↓ para historial</div>
          <div className="terminal-separator">{'─'.repeat(60)}</div>
        </div>

        {/* Historial de comandos */}
        {commandHistory.map((item, index) => (
          <div key={index} className="terminal-entry">
            <div className="terminal-command">
              <span className="terminal-prompt">{getPrompt()}</span>
              <span className="command-text">{item.command}</span>
              <span className="command-timestamp">[{item.timestamp}]</span>
            </div>
            <div className={`terminal-output ${item.success ? 'success' : 'error'}`}>
              <pre>{item.output}</pre>
            </div>
          </div>
        ))}

        {/* Línea de comando actual */}
        <div className="terminal-input-line">
          <span className="terminal-prompt">{getPrompt()}</span>
          <div className="terminal-input-container">
            <input
              ref={inputRef}
              type="text"
              value={currentCommand}
              onChange={(e) => setCurrentCommand(e.target.value)}
              onKeyDown={handleKeyDown}
              className="terminal-input"
              placeholder="Ingrese un comando..."
              disabled={isExecuting}
              autoComplete="off"
              spellCheck="false"
            />
          </div>
        </div>

        {/* Autocompletado */}
        {showAutoComplete && autoComplete.length > 0 && (
          <div className="terminal-autocomplete">
            {autoComplete.map((suggestion, index) => (
              <div
                key={index}
                className="autocomplete-item"
                onClick={() => selectAutoComplete(suggestion)}
              >
                {suggestion}
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="terminal-status">
        <span className={`status-indicator ${isExecuting ? 'busy' : 'ready'}`}>
          {isExecuting ? '🔄 Ejecutando...' : '✅ Listo'}
        </span>
        <span className="terminal-help">
          Tab: Autocompletar | ↑↓: Historial | Ctrl+L: Limpiar
        </span>
      </div>
    </div>
  );
};

export default Terminal;