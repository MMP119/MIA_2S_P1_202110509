package commands

// para eliminar un archivo que representa un disco, solo tiene el parámetro -path

import (
	"errors"
	"fmt"
	structures "server/structures"
	"strings"
)

// RmDisk es la estructura para el comando rmdisk
type RmDisk struct {
	path string
}


// ParseRmDisk analiza los parámetros del comando rmdisk
func ParseRmDisk(tokens []string) (*RmDisk, error) {
	
	cmd := &RmDisk{} // se crea un nuevo comando rmdisk (una instancia de RmDisk)

	// Reunir todos los tokens en una sola cadena para manejar las comillas
	input := strings.Join(tokens, " ")

	// Usar un analizador que respete las comillas
	args := parseArgs(input)

	for _, arg := range args {
		// Dividimos cada argumento en clave y valor usando el signo igual como separador
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("error: '%s' no es un parámetro válido", arg)
		}

		// Se obtiene la clave y el valor de cada argumento
		key, value := strings.ToLower(parts[0]), parts[1]

		// Switch para analizar cada parámetro del comando rmdisk
		switch key {

		case "-path":
			cmd.path = value

		default:
			return nil, fmt.Errorf("error: '%s' no es un parámetro válido para el comando rmdisk", key)
		}
	}

	// Verificar que se haya ingresado el parámetro -path
	if cmd.path == "" {
		return nil, errors.New("error: falta el parámetro obligatorio -path")
	}

	// SE DEBE PEDIR LA CONFIRMACION DE ELIMINAR EL DISCO, SI NO SE CONFIRMA, NO SE ELIMINA (ESTO DESDE EL FRONTEND), POR AHORA SE ELIMINA DIRECTAMENTE
	err := structures.DeleteBinaryFile(cmd.path)
	if err != nil {
		return nil, fmt.Errorf("error: no se pudo eliminar el disco: '%s'", err)
	}

	return cmd, nil
}
