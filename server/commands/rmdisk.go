package commands

import (
	"errors" 
	"fmt"    
	"regexp" 
	util "server/util" 
	"strings" 
)

// RMDISK estructura que representa el comando rmdisk con su parámetro
type RMDISK struct {
	path string // Ruta del archivo del disco
}


// CommandRmdisk parsea el comando rmdisk y devuelve una instancia de RMDISK
func ParserRmdisk(tokens []string) (*RMDISK, string, error) {
	cmd := &RMDISK{} // Crea una nueva instancia de RMDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar el parámetro del comando rmdisk
	re := regexp.MustCompile(`(?i)-path="[^"]+"|(?i)-path=[^\s]+`)
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

		// Switch para manejar el parámetro -path
		switch key {
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, "", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -path haya sido proporcionado
	if cmd.path == "" {
		return nil, "", errors.New("faltan parámetros requeridos: -path")
	}

	// SE DEBE PEDIR LA CONFIRMACION DE ELIMINAR EL DISCO, SI NO SE CONFIRMA, NO SE ELIMINA (ESTO DESDE EL FRONTEND), POR AHORA SE ELIMINA DIRECTAMENTE
	successMsg, err := util.DeleteBinaryFile(cmd.path) // Elimina el archivo binario del disco
	if err != nil {
		return nil, "", err // Devuelve un error si no se pudo eliminar el disco
	}

	return cmd, successMsg ,nil // Devuelve el comando RMDISK creado
}
