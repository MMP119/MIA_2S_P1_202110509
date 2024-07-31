import './NavBar.css';
// eslint-disable-next-line no-unused-vars
import React, {useRef} from 'react';

// eslint-disable-next-line react/prop-types
function NavBar({onFileLoad}){

    const fileInputRef = useRef(null); // constante para manejar el archivo

    const handleOpenClick = () => { // constante para manejar el click en abrir
        fileInputRef.current.click();
    };

    const handleFileChange = (event) => { // constante para manejar el cambio de archivo
        const file = event.target.files[0];
        if(file){
            const reader = new FileReader();
            reader.onload = (e) => {
                const content = e.target.result;
                onFileLoad(content);
                fileInputRef.current.value = '';
            };
            reader.readAsText(file);
        }
    };

    const handleClearClick = () => { // constante para manejar el click en limpiar

        document.getElementById('codeInput').value = '';
        document.getElementById('consoleOutput').value = '';

        // Limpiar los números de línea
        const lineNumbers = document.getElementById('lineNumbers');
        const consoleLineNumbers = document.getElementById('consoleLineNumbers');
        lineNumbers.innerHTML = '1';
        consoleLineNumbers.innerHTML = '1';

    }

    const handleRunClick = () => { // constante para manejar el click en run

        //poner el código en la consola
        const codeInput = document.getElementById('codeInput');
        const consoleOutput = document.getElementById('consoleOutput');
        consoleOutput.value = codeInput.value;

    }

    return (
        <div className="header">
            <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200" />
            <button id="openButton" onClick={handleOpenClick}>Open <span className="material-symbols-outlined">upload</span></button>
            <input
                type = "file"
                ref={fileInputRef}
                style={{display: 'none'}}
                onChange={handleFileChange}
            />
            <button id="clearButton" onClick={handleClearClick}>Clear <span className="material-symbols-outlined">mop</span></button>
            <button id="runButton" onClick={handleRunClick}>Run <span className="material-symbols-outlined">play_arrow</span> </button>
        </div>
    );
}

export default NavBar;