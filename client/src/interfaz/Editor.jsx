// eslint-disable-next-line no-unused-vars
import React, { useEffect, useState } from 'react';
import './Editor.css';
import NavBar from './NavBar';

function Editor() {

    const[code, setCode] = useState(''); // Estado para el código
    const [consoleContent, setConsoleContent] = useState(''); // Estado para el contenido de la consola

    const handleFileLoad = (content) => {
        setCode(content);
    };

    const updateLineNumbers = (textarea, lineNumberElement) => {
        const lines = textarea.value.split('\n').length;
        lineNumberElement.innerHTML = Array.from({ length: lines }, (_, i) => `<span>${i + 1}</span>`).join('');
    };

    const syncScroll = (textarea, lineNumberElement) => {
        lineNumberElement.scrollTop = textarea.scrollTop;
    };
    
    //--------------------------------------------------------------------------------
    useEffect(() => {
    const codeInput = document.getElementById('codeInput');
    const consoleOutput = document.getElementById('consoleOutput');
    const lineNumbers = document.getElementById('lineNumbers');
    const consoleLineNumbers = document.getElementById('consoleLineNumbers');

    codeInput.addEventListener('input', () => updateLineNumbers(codeInput, lineNumbers)); // Actualizar los números de línea
    codeInput.addEventListener('scroll', () => syncScroll(codeInput, lineNumbers)); // Sincronizar el scroll
    consoleOutput.addEventListener('input', () => updateLineNumbers(consoleOutput, consoleLineNumbers)); // Actualizar los números de línea
    consoleOutput.addEventListener('scroll', () => syncScroll(consoleOutput, consoleLineNumbers)); // Sincronizar el scroll

    // Inicializar los números de línea
    updateLineNumbers(codeInput, lineNumbers); // Números de línea del editor de código
    updateLineNumbers(consoleOutput, consoleLineNumbers); // Números de línea de la consola

    // Limpiar los event listeners al desmontar el componente
    return () => {
        codeInput.removeEventListener('input', () => updateLineNumbers(codeInput, lineNumbers)); 
        codeInput.removeEventListener('scroll', () => syncScroll(codeInput, lineNumbers));
        consoleOutput.removeEventListener('input', () => updateLineNumbers(consoleOutput, consoleLineNumbers));
        consoleOutput.removeEventListener('scroll', () => syncScroll(consoleOutput, consoleLineNumbers));
    };
    }, []);

    //--------------------------------------------------------------------------------
    // Actualizar los números de línea AL cargar un archivo

    useEffect(() => {
        const codeInput = document.getElementById('codeInput');
        const lineNumbers = document.getElementById('lineNumbers');

        //actualiozar los números de línea al cargar un archivo
        if(codeInput){
            updateLineNumbers(codeInput, lineNumbers);
        }

    }, [code]);


    //--------------------------------------------------------------------------------

    return (
    <div className="editor">
        <NavBar onFileLoad={handleFileLoad}/> {/* pasa la función handleFileLoad como prop */}
        <div className="main">
        <div className="entrada">
            <div className="line-numbers" id="lineNumbers"></div>
            <textarea id="codeInput" placeholder="Write your command here..." value={code} onChange={(e) => setCode(e.target.value)}></textarea>
        </div>
        <div className="console">
            <div className="line-numbers" id="consoleLineNumbers"></div>
            <textarea id="consoleOutput" readOnly placeholder="Console output..." value={consoleContent} onChange={(e) => setConsoleContent(e.target.value)}></textarea>
        </div>
        </div>
    </div>
    );
}

export default Editor;