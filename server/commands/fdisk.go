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
	cmd := &structures.FDISK{} // Crea una nueva instancia de FDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando fdisk
	re := regexp.MustCompile(`(?i)-size=\d+|(?i)-unit=[kKmM]|(?i)-fit=[bBfF]{2}|(?i)-path="[^"]+"|(?i)-path=[^\s]+|(?i)-type=[pPeElL]|(?i)-name="[^"]+"|(?i)-name=[^\s]+`)
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
			if value != "K" && value != "M" {
				return nil, "", errors.New("la unidad debe ser K o M")
			}
			cmd.Unit = strings.ToUpper(value)
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
		case "-type":
			// Verifica que el tipo sea "P", "E" o "L"
			value = strings.ToUpper(value)
			if value != "P" && value != "E" && value != "L" {
				return nil, "", errors.New("el tipo debe ser P, E o L")
			}
			cmd.TypE = value
		case "-name":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return nil, "", errors.New("el nombre no puede estar vacío")
			}
			cmd.Name = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size, -path y -name hayan sido proporcionados
	if cmd.Size == 0 {
		return nil, "", errors.New("faltan parámetros requeridos: -size")
	}
	if cmd.Path == "" {
		return nil, "", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.Name == "" {
		return nil, "", errors.New("faltan parámetros requeridos: -name")
	}

	// Si no se proporcionó la unidad, se establece por defecto a "M"
	if cmd.Unit == "" {
		cmd.Unit = "M"
	}

	// Si no se proporcionó el ajuste, se establece por defecto a "FF"
	if cmd.Fit == "" {
		cmd.Fit = "WF"
	}

	// Si no se proporcionó el tipo, se establece por defecto a "P"
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
	// Imprimir las particiones
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