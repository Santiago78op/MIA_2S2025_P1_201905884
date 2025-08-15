import React from 'react';
import ConnectionStatus from './components/ConnectionStatus';
import FileSystemList from './components/FileSystemList';
import CommandExecutor from './components/CommandExecutor';
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>🗃️ Sistema de Archivos EXT2</h1>
        <p>Simulador de sistema de archivos - MIA 2S2025</p>
      </header>
      
      <main className="App-main">
        <ConnectionStatus />
        <CommandExecutor />
        <FileSystemList />
      </main>
      
      <footer className="App-footer">
        <p>Proyecto MIA - Sistema de Archivos EXT2</p>
        <small>Frontend: React + TypeScript | Backend: Go</small>
      </footer>
    </div>
  );
}

export default App;
