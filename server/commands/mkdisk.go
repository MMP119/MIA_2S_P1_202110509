package commands

import (
	//para acceder a las funciones de structures
	structures "server/structures"
	"errors" //para manejar errores
	"fmt"
	"strconv" //para convertir cadenas a otros tipos de datos
	"strings"
)

//estructura que representa un comando mkdisk

type MKDISK struct {

	size int	//tamaño del disco
	fit string	//tipo de ajuste del disco (BF, FF, WF)
	unit string	//unidad de medida del tamaño del disco (K, M)
	path string	//ruta donde se creará el disco
}

//función que se encarga de crear un disco, analiza los parámetros del comando mkdisk

func ParseMkdisk(tokens []string)(*MKDISK, error){ //retorna un puntero a MKDISK y un error
	cmd := &MKDISK{} //se crea un nuevo comando mkdisk (una instancia de MKDISK)
	for _, token := range tokens {

		//dividimos cada token en clave y valor usando el signo igual como separador
		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("error: '%s' no es un parámetro válido", token)
		}

		//se obtiene la clave y el valor de cada token
		key, value := strings.ToLower(parts[0]), parts[1]

		//switch para analizar cada parámetro del comando mkdisk
		switch key{

			case "-size":
				//se convierte el valor a entero
				size, err := strconv.Atoi(value)
				if err != nil  || size <= 0 {
					return nil, errors.New("error: el tamaño debe ser un número entero positivo")
				}
				cmd.size = size

			
			case "-fit":
				//se convierte el valor a mayúsculas
				fit := strings.ToUpper(value)
				if fit != "BF" && fit != "FF" && fit != "WF" {
					return nil, errors.New("error: el tipo de ajuste debe ser BF, FF o WF")
				}
				cmd.fit = fit
	
				
			case "-unit":
				//se convierte el valor a mayúsculas
				unit := strings.ToUpper(value)
				if unit != "K" && unit != "M" {
					return nil, errors.New("error: la unidad de medida debe ser K o M")
				}
				cmd.unit = unit
			

			case "-path":
				//se establece la ruta donde se creará el disco, si no existen las carpetas, se crean
				cmd.path = value
			
			default:
				return nil, fmt.Errorf("error: parámetro '%s' no reconocido", key)

		}

	}

	//verificamos que el parametro -size se haya ingresado
	if cmd.size == 0 {
		return nil, errors.New("error: el parámetro -size es obligatorio")
	}

	//verificamos que el parametro -fit se haya ingresado, si no se establece, por defecto es FF
	if cmd.fit == "" {
		cmd.fit = "FF"
	}

	// SI no se establece la unidad, por defecto es M
	if cmd.unit == "" {
		cmd.unit = "M"
	}

	if cmd.path == "" {
		return nil, errors.New("error: el parámetro -path es obligatorio")
	}

	//llamamos a la funcion CreateBinaryFile para crear el disco
	err := structures.CreateBinaryFile(cmd.size, cmd.fit ,cmd.unit, cmd.path)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}


	return cmd, nil //retorna el comando mkdisk y nil (sin errores)
}
