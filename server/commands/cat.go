package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
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
	
	//se debe leer el archivo y retornar el contenido del archivo, si no existe el archivo se debe retornar un error
	file, err := os.Open(cmd.Filen)
	if err != nil {
		return "Error: no se pudo abrir el archivo", fmt.Errorf("error: no se pudo abrir el archivo: '%s'", err)
	}

	defer file.Close()

	//leer el contenido del archivo
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "Error: no se pudo leer el archivo", fmt.Errorf("error: no se pudo leer el archivo: '%s'", err)
	}

	return string(content), nil

}