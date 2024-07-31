// eslint-disable-next-line no-unused-vars
import React, { useEffect } from 'react';
import './Editor.css';

function Editor() {
    
    useEffect(() => {
    const codeInput = document.getElementById('codeInput');
    const consoleOutput = document.getElementById('consoleOutput');
    const lineNumbers = document.getElementById('lineNumbers');
    const consoleLineNumbers = document.getElementById('consoleLineNumbers');

    const updateLineNumbers = (textarea, lineNumberElement) => {
        const lines = textarea.value.split('\n').length;
        lineNumberElement.innerHTML = Array.from({ length: lines }, (_, i) => `<span>${i + 1}</span>`).join('');
    };

    const syncScroll = (textarea, lineNumberElement) => {
        lineNumberElement.scrollTop = textarea.scrollTop;
    };

    codeInput.addEventListener('input', () => updateLineNumbers(codeInput, lineNumbers));
    codeInput.addEventListener('scroll', () => syncScroll(codeInput, lineNumbers));
    consoleOutput.addEventListener('input', () => updateLineNumbers(consoleOutput, consoleLineNumbers));
    consoleOutput.addEventListener('scroll', () => syncScroll(consoleOutput, consoleLineNumbers));

    // Inicializar los números de línea
    updateLineNumbers(codeInput, lineNumbers);
    updateLineNumbers(consoleOutput, consoleLineNumbers);

    // Limpiar los event listeners al desmontar el componente
    return () => {
        codeInput.removeEventListener('input', () => updateLineNumbers(codeInput, lineNumbers));
        codeInput.removeEventListener('scroll', () => syncScroll(codeInput, lineNumbers));
        consoleOutput.removeEventListener('input', () => updateLineNumbers(consoleOutput, consoleLineNumbers));
        consoleOutput.removeEventListener('scroll', () => syncScroll(consoleOutput, consoleLineNumbers));
    };
    }, []);

    return (
    <div className="editor">
        <div className="main">
        <div className="entrada">
            <div className="line-numbers" id="lineNumbers"></div>
            <textarea id="codeInput" placeholder="Write your command here..."></textarea>
        </div>
        <div className="console">
            <div className="line-numbers" id="consoleLineNumbers"></div>
            <textarea id="consoleOutput" placeholder="Console output..."></textarea>
        </div>
        </div>
    </div>
    );
}

export default Editor;