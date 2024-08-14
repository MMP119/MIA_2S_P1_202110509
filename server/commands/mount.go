package commands

import (
	"errors" 
	"fmt"    
	"regexp" 
	"strings" 
)


type MOUNT struct {
	path string 
	name string 
}

func ParserMount(tokens []string) (*MOUNT, string, error) {
	cmd := &MOUNT{}

	args := strings.Join(tokens, " ")

	re := regexp.MustCompile(`(?i)-path="[^"]+"|(?i)-path=[^\s]+|(?i)-name="[^"]+"|(?i)-name=[^\s]+`)

	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {

		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, "ERROR: formato de parámetro inválido", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]


		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
			case "-path":
				if value == "" {
					return nil,"ERROR: el path no puede estar vacío", errors.New("el path no puede estar vacío")
				}
				cmd.path = value
			case "-name":
				if value == "" {
					return nil, "ERROR: el nombre no puede estar vacío", errors.New("el nombre no puede estar vacío")
				}
				cmd.name = value
			default:
				return nil, "ERROR: parámetro desconocido", fmt.Errorf("parámetro desconocido: %s", key)
			}
	}

	// Verifica que los parámetros -path y -name hayan sido proporcionados
	if cmd.path == "" {
		return nil, "ERROR: faltan parámetros requeridos: -path", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return nil, "ERROR: faltan parámetros requeridos: -name", errors.New("faltan parámetros requeridos: -name")
	}

	/*
		PRÓXIMAMENTE
	*/
	return cmd, "", nil // Devuelve el comando MOUNT creado
}