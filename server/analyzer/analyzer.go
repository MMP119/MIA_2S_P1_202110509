package analyzer

import (
	commands "server/commands"
	"fmt"
	"errors"	//para manejar errores
	"os"
	"os/exec"
	"strings"
)


// Para analizar cualquier comando de entrada 
func Analyzer(input string) (interface{}, error){ //retorna un interface{} y un error

	// Se divide el input en palabras
	words := strings.Fields(input)

	if len(words) == 0 {
		return nil, errors.New("no se ingresó ningún comando")
	}

	// Se obtiene el comando y se maneja con un switch

	switch words[0] {

		case "mkdisk":
			//aquí se llama a la funcion que se encarga de crear un disco
			return commands.ParseMkdisk(words[1:])

		case "rmdisk":
			return commands.ParseRmDisk(words[1:])

		case "fdisk":
			return nil, errors.New("comando fdisk no implementado")

		case "mount":
			return nil, errors.New("comando mount no implementado")
		
		case "mkfs":
			return nil, errors.New("comando mkfs no implementado")

		case "cat":
			return nil, errors.New("comando cat no implementado")

		case "login":
			return nil, errors.New("comando login no implementado")

		case "logout":
			return nil, errors.New("comando logout no implementado")

		case "mkgrp":
			return nil, errors.New("comando mkgrp no implementado")
		
		case "rmgrp":
			return nil, errors.New("comando rmgrp no implementado")
		
		case "mkusr":
			return nil, errors.New("comando mkusr no implementado")
		
		case "rmusr":
			return nil, errors.New("comando rmusr no implementado")
		
		case "chgrp":
			return nil, errors.New("comando chgrp no implementado")
		
		case "mkfile":
			return nil, errors.New("comando mkfile no implementado")
		
		case "mkdir":
			return nil, errors.New("comando mkdir no implementado")
		
		case "rep": //reportes
			return nil, errors.New("comando rep no implementado")


		case "clear":
			//para limpiar la consola
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout //redirige la salida estandar a la consola
			err := cmd.Run()		//ejecuta el comando

			if err != nil{
				//si hay un error, se retorna el error
				return nil, errors.New("error al limpiar la consola")
			}
			return nil, nil //si no hay errores, retorna nil

		default:
			return nil, fmt.Errorf("comando no reconocido: %s", words[0])

	}

}