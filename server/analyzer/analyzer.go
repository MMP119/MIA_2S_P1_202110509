package analyzer

import (
	commands "server/commands" 
	"errors"                    
	"fmt"                       
	"os"                        
	"os/exec"                   
	"strings"                  
)

// Analyzer analiza el comando de entrada y ejecuta la acción correspondiente
func Analyzer(input string) (interface{}, error) {
	// Divide la entrada en tokens usando espacios en blanco como delimitadores
	tokens := strings.Fields(input)

	// Si no se proporcionó ningún comando, devuelve un error
	if len(tokens) == 0 {
		return nil, errors.New("no se proporcionó ningún comando")
	}

	// pasar el comando a minúsculas
	tokens[0] = strings.ToLower(tokens[0])

	// Switch para manejar diferentes comandos
	switch tokens[0] {

		case "mkdisk":
			return commands.ParserMkdisk(tokens[1:])

		case "rmdisk":
			return commands.ParserRmdisk(tokens[1:])

		case "fdisk":
			return commands.ParserFdisk(tokens[1:])

		case "mount":
			return commands.ParserMount(tokens[1:])

		case "clear":
			// Crea un comando para limpiar la terminal
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout // Redirige la salida del comando a la salida estándar
			err := cmd.Run()       // Ejecuta el comando
			if err != nil {
				// Si hay un error al ejecutar el comando, devuelve un error
				return nil, errors.New("no se pudo limpiar la terminal")
			}
			return nil, nil // Devuelve nil si el comando se ejecutó correctamente
		default:
			// Si el comando no es reconocido, devuelve un error
			return nil, fmt.Errorf("comando desconocido: %s", tokens[0])
		}
}