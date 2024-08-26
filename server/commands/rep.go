package commands

import (
	"errors"
	"fmt"
	"regexp"
	global "server/global"
	structures "server/structures"
	"strings"

)

type REP struct {
	//name
	Path string 
	//id
	//path_file_ls
}

func ParseRep(tokens []string)(*REP, string, error){
	
	cmd := &REP{}

	args := strings.Join(tokens, " ")

	re := regexp.MustCompile(`-size=\d+|-unit=[kKmM]|-fit=[bBfF]{2}|-path="[^"]+"|-path=[^\s]+|-type=[pPeElL]|-name="[^"]+"|-name=[^\s]+`)

	matches := re.FindAllString(args, -1)

	
	for _, math := range matches{
		kv := strings.SplitN(math, "=", 2)
		if len(kv) != 2 {
			return nil, "ERROR: formato de parámetro inválido", fmt.Errorf("formato de parámetro inválido: %s", math)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key{
			case "-path":
				if value == "" {
					return nil, "ERROR: el path es obligatorio", errors.New("el path es obligatorio")
				}
				cmd.Path = value
			
			default: 
				return nil, "ERROR: parámetro no reconocido", fmt.Errorf("parámetro no reconocido: %s", key)
		}
	}

	if cmd.Path == "" {
		return nil, "ERROR: el path es obligatorio", errors.New("el path es obligatorio")
	}

	msg, err := CommandRep(cmd)
	if err != nil {
		return nil, msg, err
	}

	
	return cmd, "", nil
}

func CommandRep(cmd *REP) (string, error) {
	// Crear una nueva estructura MBR
	mbr := &structures.MBR{}

	// Deserializar la estructura MBR desde el archivo binario
	msg, err := mbr.DeserializeMBR(cmd.Path)
	if err != nil {
		return msg, err
	}

	// Imprimir la información del MBR
	fmt.Println("\nMBR\n----------------")
	mbr.Print()

	// Imprimir la información de cada partición
	fmt.Println("\nParticiones\n----------------")
	mbr.PrintPartitions()

	// Imprimir partidas montadas
	fmt.Println("\nParticiones montadas\n----------------")

	for id, path := range global.MountedPartitions {
		fmt.Printf("ID: %s, PATH: %s\n", id, path)
	}

	// Imprimir el SuperBloque de cada partición montada
	index := 0
	// Iterar sobre cada partición montada
	for id, path := range global.MountedPartitions {
		// Crear una nueva estructura SuperBloque
		sb := &structures.SuperBlock{}
		// Deserializar la estructura SuperBloque desde el archivo binario
		err := sb.Deserialize(path, int64(mbr.Mbr_partitions[index].Part_start))
		if err != nil {
			fmt.Printf("Error al leer el SuperBloque de la partición %s: %s\n", id, err)
			continue
		}
		fmt.Printf("\nPartición %s\n----------------", id)

		// Imprimir la información del SuperBloque
		fmt.Println("\nSuperBloque:")
		sb.Print()

		// Imprimir los inodos
		sb.PrintInodes(path)

		// Imprimir los bloques
		sb.PrintBlocks(path)

		index++
	}

	return "", nil
}