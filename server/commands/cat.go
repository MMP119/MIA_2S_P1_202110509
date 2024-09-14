package commands

import (
	"errors"
	"fmt"
	"regexp"
	global "server/global"
	"strings"
)

type CAT struct{
	Filen string
}


func ParseCat(tokens []string) (*CAT, string, error) {
	cmd := &CAT{}

	args := strings.Join(tokens, " ")

	re := regexp.MustCompile(`(?i)-file[1-9][0-9]*="[^"]+"|(?i)-file[1-9][0-9]*=[^\s]+`)

	matches := re.FindAllString(args, -1)

	var allContent strings.Builder // Acumulador para el contenido de todos los archivos

	for _, match := range matches {

		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, "ERROR: formato de parámetro inválido", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		if strings.HasPrefix(key, "-file") {
			if value == "" {
				return nil, "ERROR: la ruta del archivo no puede estar vacía", errors.New("el nombre del archivo no puede estar vacío")
			}
			cmd.Filen = value
		} else {
			return nil, "ERROR: parámetro no reconocido", fmt.Errorf("parámetro no reconocido: %s", key)
		}

		if cmd.Filen == "" {
			return nil, "ERROR: el nombre del archivo es obligatorio", errors.New("el nombre del archivo es obligatorio")
		}

		// Leer el contenido del archivo
		content, err := CommandCAT(cmd)
		if err != nil {
			return nil, content, err
		}

		// Concatenar el contenido al acumulador
		allContent.WriteString(content)
	}

	// Retornar el contenido concatenado
	return cmd, "Comando CAT: realizado correctamente\n" + allContent.String(), nil
}

func CommandCAT(cmd *CAT) (string, error){
	
	// leer un archivo que esté en la ruta especificada dentro del bloque
	// inodo -> bloque -> contenido

	//la ruta del archivo es cmd.Filen, donde está el inodo -> bloque -> contenido

	//obtener el id de la particion donde se está logueado
	idPartition := global.GetIDSession()

	//obtenemos primero el superbloque para obtener el inodo raíz y luego el inodo del archivo
	partitionSuperblock, partition, partitionPath, err := global.GetMountedPartitionSuperblock(idPartition)
	if err != nil {
		return "Error al obtener la partición montada en el comando login", fmt.Errorf("error al obtener la partición montada: %v", err)
	}


	





}