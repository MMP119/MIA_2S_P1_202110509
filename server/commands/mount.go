package commands

import (
	"errors" 
	"fmt"    
	"regexp" 
	"strings" 
	structures "server/structures"
)




func ParserMount(tokens []string) (*structures.MOUNT, string, error) {
	cmd := &structures.MOUNT{}

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
				cmd.Path = value
			case "-name":
				if value == "" {
					return nil, "ERROR: el nombre no puede estar vacío", errors.New("el nombre no puede estar vacío")
				}
				cmd.Name = value
			default:
				return nil, "ERROR: parámetro desconocido", fmt.Errorf("parámetro desconocido: %s", key)
			}
	}

	if cmd.Path == "" {
		return nil, "ERROR: faltan parámetros requeridos: -path", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.Name == "" {
		return nil, "ERROR: faltan parámetros requeridos: -name", errors.New("faltan parámetros requeridos: -name")
	}

	// se monta la partición 

	msg, err := structures.CommandMount(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, msg, err
	}

	return cmd, "", nil // Devuelve el comando MOUNT creado
}