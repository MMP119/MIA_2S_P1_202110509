// eslint-disable-next-line no-unused-vars
import React, { useState } from 'react';
import NavBar from './interfaz/NavBar';
import Editor from './interfaz/Editor';
import './App.css';

function App(){

  const [code, setCode] = useState(''); // Estado para el código
  const [consoleOutput, setConsoleOutput] = useState(''); // Estado para la salida de la consola

  const handleFileLoad = (content) => { // Función para manejar la carga de un archivo
    setCode(content);
  };

  const handleClear = () => { // Función para manejar el click en limpiar
    setCode('');
    setConsoleOutput('');
  };

  const handleRun = () => { // Función para manejar el click en run
    setConsoleOutput(code);
  }

  return (
    <div className="App">
      <NavBar onFileLoad={handleFileLoad} onClear={handleClear} onRun={handleRun}/>
      <h2>Command Input:</h2>
      <div className='main'>
        <Editor code={code} />
        <h2>Console:</h2>
        <textarea id="consoleOutput" className='console-output' readOnly value={consoleOutput}></textarea>
      </div>
    </div>
  );
}

export default App;
