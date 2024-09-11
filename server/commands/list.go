package commands

import (
	"errors"
	"fmt"
	"regexp"
	global "server/global"
	"strings"
)

func ParseList(tokens []string) (string, string, error) {

	args := strings.Join(tokens, " ")

	re := regexp.MustCompile(`(?i)-path="[^"]+"|(?i)-path=[^\s]+|(?i)-name="[^"]+"|(?i)-name=[^\s]+`)

	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {

		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", "ERROR: formato de parámetro inválido", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]


		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
			case "-path":
				if value == "" {
					return "","ERROR: el path no puede estar vacío", errors.New("el path no puede estar vacío")
				}

				partitions := []string{}
				for id, disk:= range global.MountedPartitions {
					if(disk == value){
						partitions = append(partitions, id)	
					}		
				}

				// Verificar si hay particiones montadas
				if len(partitions) == 0 {
					return "", "\n No hay particiones montadas\n", nil
				}

				// Crear mensaje de salida
				msg := "\n Particiones montadas en el disco:\n"
				for _, partition := range partitions {
					msg += fmt.Sprintf("- %s\n", partition)
				}

				return "", msg, nil


			}
	}

	return "", "ERROR: faltan parámetros requeridos en el list: -path", errors.New("faltan parámetros requeridos en el list: -path")

}