package commands

import (
	"errors"
	"fmt"
	"regexp"
	structures "server/structures"
	"strconv"
	"strings"
)



func ParserMkdisk(tokens []string) (*structures.MKDISK, string,error) {

	cmd := &structures.MKDISK{} // Crea una nueva instancia de MKDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")

	// Expresión regular para encontrar los parámetros del comando mkdisk
	re := regexp.MustCompile(`(?i)-size=\d+|(?i)-unit=[kKmM]|(?i)-fit=[bBfFwW]{2}|(?i)-path="[^"]+"|(?i)-path=[^\s]+`)

	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, "", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {

			case "-size":
				// Convierte el valor del tamaño a un entero
				size, err := strconv.Atoi(value)
				if err != nil || size <= 0 {
					return nil, "", errors.New("el tamaño debe ser un número entero positivo")
				}
				cmd.Size = size

			case "-unit":
				// Verifica que la unidad sea "K" o "M"
				value = strings.ToUpper(value)
				if value != "K" && value != "M" {
					return nil, "", errors.New("la unidad debe ser K o M")
				}
				cmd.Unit = value

			case "-fit":
				// Verifica que el ajuste sea "BF", "FF" o "WF"
				value = strings.ToUpper(value)
				if value != "BF" && value != "FF" && value != "WF" {
					return nil, "", errors.New("el ajuste debe ser BF, FF o WF")
				}
				cmd.Fit = value

			case "-path":
				// Verifica que el path no esté vacío
				if value == "" {
					return nil, "", errors.New("el path no puede estar vacío")
				}
				cmd.Path = value

			default:
				// Si el parámetro no es reconocido, devuelve un error
				return nil, "", fmt.Errorf("parámetro desconocido: %s", key)
			}
	}

	// Verifica que los parámetros -size y -path hayan sido proporcionados
	if cmd.Size == 0 {
		return nil, "", errors.New("faltan parámetros requeridos: -size")
	}
	if cmd.Path == "" {
		return nil, "", errors.New("faltan parámetros requeridos: -path")
	}

	// Si no se proporcionó la unidad, se establece por defecto a "M"
	if cmd.Unit == "" {
		cmd.Unit = "M"
	}

	// Si no se proporcionó el ajuste, se establece por defecto a "FF"
	if cmd.Fit == "" {
		cmd.Fit = "FF"
	}

	// Crear el disco con los parámetros proporcionados
	err := structures.CommandMkdisk(cmd)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return cmd, "Disco Creado Exitosamente", nil // Devuelve el comando MKDISK creado
}

