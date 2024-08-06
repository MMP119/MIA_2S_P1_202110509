package commands

import (
	// para acceder a las funciones de structures
	structures "server/structures"
	"errors" // para manejar errores
	"fmt"
	"strconv" // para convertir cadenas a otros tipos de datos
	"strings"
)

// estructura que representa un comando mkdisk
type MKDISK struct {
	size int    // tamaño del disco
	fit  string // tipo de ajuste del disco (BF, FF, WF)
	unit string // unidad de medida del tamaño del disco (K, M)
	path string // ruta donde se creará el disco
}

// función que se encarga de crear un disco, analiza los parámetros del comando mkdisk
func ParseMkdisk(tokens []string) (*MKDISK, error) { // retorna un puntero a MKDISK y un error
	cmd := &MKDISK{} // se crea un nuevo comando mkdisk (una instancia de MKDISK)

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

		// Switch para analizar cada parámetro del comando mkdisk
		switch key {

		case "-size":
			// Se convierte el valor a entero
			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return nil, errors.New("error: el tamaño debe ser un número entero positivo")
			}
			cmd.size = size

		case "-fit":
			// Se convierte el valor a mayúsculas
			fit := strings.ToUpper(value)
			if fit != "BF" && fit != "FF" && fit != "WF" {
				return nil, errors.New("error: el tipo de ajuste debe ser BF, FF o WF")
			}
			cmd.fit = fit

		case "-unit":
			// Se convierte el valor a mayúsculas
			unit := strings.ToUpper(value)
			if unit != "K" && unit != "M" {
				return nil, errors.New("error: la unidad de medida debe ser K o M")
			}
			cmd.unit = unit

		case "-path":
			// Se establece la ruta donde se creará el disco
			cmd.path = strings.Trim(value, "\"") // Elimina las comillas dobles alrededor de la ruta si existen

		default:
			return nil, fmt.Errorf("error: parámetro '%s' no reconocido", key)
		}
	}

	// Verificamos que el parámetro -size se haya ingresado
	if cmd.size == 0 {
		return nil, errors.New("error: el parámetro -size es obligatorio")
	}

	// Verificamos que el parámetro -fit se haya ingresado, si no se establece, por defecto es FF
	if cmd.fit == "" {
		cmd.fit = "FF"
	}

	// Si no se establece la unidad, por defecto es M
	if cmd.unit == "" {
		cmd.unit = "M"
	}

	if cmd.path == "" {
		return nil, errors.New("error: el parámetro -path es obligatorio")
	}

	// Llamamos a la función CreateBinaryFile para crear el disco
	err := structures.CreateBinaryFile(cmd.size, cmd.fit, cmd.unit, cmd.path)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	return cmd, nil // Retorna el comando mkdisk y nil (sin errores)
}

// parseArgs analiza una cadena de entrada en una lista de argumentos, respetando las comillas
func parseArgs(input string) []string {
	var args []string
	var currentArg strings.Builder
	inQuotes := false

	for _, char := range input {
		switch char {
		case ' ':
			if inQuotes {
				currentArg.WriteRune(char)
			} else if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
		case '"':
			inQuotes = !inQuotes
		default:
			currentArg.WriteRune(char)
		}
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}
