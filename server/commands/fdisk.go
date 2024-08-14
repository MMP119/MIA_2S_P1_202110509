package commands

import (
	structures "server/structures" 
	util "server/util"
	"errors"  
	"fmt"     
	"regexp"  
	"strconv" 
	"strings" 
)



// CommandFdisk parsea el comando fdisk y devuelve una instancia de FDISK
func ParserFdisk(tokens []string) (*structures.FDISK, string, error) {
	cmd := &structures.FDISK{} 

	args := strings.Join(tokens, " ")

	re := regexp.MustCompile(`(?i)-size=\d+|(?i)-unit=[kKmM]|(?i)-fit=[bBfF]{2}|(?i)-path="[^"]+"|(?i)-path=[^\s]+|(?i)-type=[pPeElL]|(?i)-name="[^"]+"|(?i)-name=[^\s]+`)

	matches := re.FindAllString(args, -1)


	for _, match := range matches {

		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, "ERROR: formato de parámetro inválido", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
		case "-size":

			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return nil, "ERROR: el tamaño debe ser un número entero positivo", errors.New("el tamaño debe ser un número entero positivo")
			}
			cmd.Size = size
		case "-unit":

			if value != "K" && value != "M" {
				return nil, "ERROR: La unidad debe ser K o M", errors.New("la unidad debe ser K o M")
			}
			cmd.Unit = strings.ToUpper(value)
		case "-fit":

			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return nil, "ERROR: el ajuste debe ser BF, FF o WF", errors.New("el ajuste debe ser BF, FF o WF")
			}
			cmd.Fit = value
		case "-path":

			if value == "" {
				return nil, "ERROR: el path no puede estar vacío", errors.New("el path no puede estar vacío")
			}
			cmd.Path = value
		case "-type":

			value = strings.ToUpper(value)
			if value != "P" && value != "E" && value != "L" {
				return nil, "ERROR: el tipo debe ser P, E o L", errors.New("el tipo debe ser P, E o L")
			}
			cmd.TypE = value
		case "-name":

			if value == "" {
				return nil, "ERROR: el nombre no puede estar vacío", errors.New("el nombre no puede estar vacío")
			}
			cmd.Name = value
		default:

			return nil, "ERROR: parámetro desconocido", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size, -path y -name hayan sido proporcionados
	if cmd.Size == 0 {
		return nil, "ERROR: faltan parámetros requeridos: -size", errors.New("faltan parámetros requeridos: -size")
	}
	if cmd.Path == "" {
		return nil, "ERROR: faltan parámetros requeridos: -path", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.Name == "" {
		return nil, "ERROR: Faltan parámetros requeridos: -name", errors.New("faltan parámetros requeridos: -name")
	}

	if cmd.Unit == "" {
		cmd.Unit = "M"
	}

	if cmd.Fit == "" {
		cmd.Fit = "WF"
	}

	if cmd.TypE == "" {
		cmd.TypE = "P"
	}

	/* ---------- Ejemplo para ver MBR y Particiones ----------*/
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(cmd.Path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return nil, "", err
	}

	// Imprimir la estructura
	mbr.Print()
	fmt.Println("-----------------------------")

	mbr.PrintPartitions()

	return cmd, "", nil // Devuelve el comando FDISK creado
}

func commandFdisk(fdisk *structures.FDISK) error {
	// Convertir el tamaño a bytes
	_, err := util.ConvertToBytes(fdisk.Size, fdisk.Unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}

	/*
		PRÓXIMAMENTE.........................................
	*/
	return nil
}