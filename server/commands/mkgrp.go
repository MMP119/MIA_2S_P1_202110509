package commands

import (
	"errors"
	"fmt"
	"regexp"
	global "server/global"
	structures "server/structures"
	"strconv"
	"strings"
)

type MKGRP struct {
	Name string
}


//este comando crea un grupo para los usuarios de la particion
func ParseMkgrp(tokens []string)(*MKGRP, string, error){

	cmd := &MKGRP{}

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mkfile
	re := regexp.MustCompile(`-name="[^"]+"|-name=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	for _, match := range matches{
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, "ERROR: formato de parámetro inválido", fmt.Errorf("formato de parámetro inválido: %s", match)
		}

		key, value := strings.ToLower(kv[0]), kv[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
			case "-name":
				if cmd.Name != "" {
					return nil, "ERROR: nombre no puede estar vacío", errors.New("nombre no puede estar vacío")
				}
				cmd.Name = value
			default:
				return nil, "ERROR: parámetro inválido", fmt.Errorf("parámetro inválido: %s", key)
		}
	}

	if cmd.Name == "" {
		return nil, "ERROR: nombre no puede estar vacío", errors.New("nombre no puede estar vacío")
	}

	msg, err := CommandMkgrp(cmd)
	if err != nil {
		return nil, msg, err
	}

	return cmd, "", nil

}


func CommandMkgrp(cmd *MKGRP) (string, error) {
	idPartition := global.GetIDSession()
	fmt.Println("ID de partición:", idPartition)

	partitionSuperblock, _, partitionPath, err := global.GetMountedPartitionSuperblock(idPartition)
	if err != nil {
		return "Error al obtener la partición montada en el comando login", fmt.Errorf("error al obtener la partición montada: %v", err)
	}
	fmt.Println("Ruta de partición:", partitionPath)

	inode := &structures.Inode{}
	err = inode.Deserialize(partitionPath, int64(partitionSuperblock.S_inode_start+(0*partitionSuperblock.S_inode_size)))
	if err != nil {
		return "error al obtener el inodo raiz", fmt.Errorf("error al obtener el inodo raiz: %v", err)
	}
	fmt.Println("Inodo raíz obtenido exitosamente")

	if inode.I_block[0] == 0 {
		fmt.Println("El primer inodo está en cero, moviéndose al bloque 0")
		folderBlock := &structures.FolderBlock{}
		err = folderBlock.Deserialize(partitionPath, int64(partitionSuperblock.S_block_start+(inode.I_block[0]*partitionSuperblock.S_block_size)))
		if err != nil {
			return "Error al obtener el bloque 0", fmt.Errorf("error al obtener el bloque 0: %v", err)
		}

		var usersFileInodeIndex int64
		foundUsersFile := false
		for _, contenido := range folderBlock.B_content {
			name := strings.Trim(string(contenido.B_name[:]), "\x00")
			if name == "users.txt" {
				usersFileInodeIndex = int64(contenido.B_inodo)
				foundUsersFile = true
				break
			}
		}

		if !foundUsersFile {
			return "El archivo users.txt no se encontró en el bloque 0", nil
		}

		fmt.Println("Archivo users.txt encontrado con inodo:", usersFileInodeIndex)

		err = inode.Deserialize(partitionPath, int64(partitionSuperblock.S_inode_start)+(usersFileInodeIndex*int64(partitionSuperblock.S_inode_size)))
		if err != nil {
			return "Error al obtener el inodo del archivo users.txt", fmt.Errorf("error al obtener el inodo del archivo users.txt: %v", err)
		}

		if inode.I_block[0] == 1 {
			fmt.Println("El primer inodo está en 1, moviéndose al bloque 1")
			fileBlock := &structures.FileBlock{}
			err = fileBlock.Deserialize(partitionPath, int64(partitionSuperblock.S_block_start+(inode.I_block[0]*partitionSuperblock.S_block_size)))
			if err != nil {
				return "error al obtener el bloque 1 del archivo users.txt", fmt.Errorf("error al obtener el bloque 1 del archivo users.txt: %v", err)
			}

			users := strings.Split(strings.TrimSpace(string(fileBlock.B_content[:])), "\n")
			fmt.Println("Contenido actual de users.txt:", users)

			newGroupName := cmd.Name
			for _, user := range users {
				if user != "" {
					values := strings.Split(user, ",")
					if len(values) > 1 && values[1] == "G" && values[2] == newGroupName {
						return "El grupo ya existe", nil
					}
				}
			}

			lastNumber := 0
			for _, user := range users {
				if user != "" {
					values := strings.Split(user, ",")
					num, err := strconv.Atoi(values[0])
					if err == nil && num > lastNumber {
						lastNumber = num
					}
				}
			}

			fmt.Println("Último número encontrado:", lastNumber)

			newGroupEntry := fmt.Sprintf("%d,G,%s\n", lastNumber+1, newGroupName)
			users = append(users, newGroupEntry)
			copy(fileBlock.B_content[:], []byte(strings.Join(users, "\n")))

			fmt.Println("Contenido a escribir en users.txt:", string(fileBlock.B_content[:]))

			err = fileBlock.Serialize(partitionPath, int64(partitionSuperblock.S_block_start+(inode.I_block[0]*partitionSuperblock.S_block_size)))
			if err != nil {
				return "Error al actualizar el bloque 1 del archivo users.txt", fmt.Errorf("error al actualizar el bloque 1 del archivo users.txt: %v", err)
			}

			err = fileBlock.Deserialize(partitionPath, int64(partitionSuperblock.S_block_start+(inode.I_block[0]*partitionSuperblock.S_block_size)))
			if err != nil {
				return "Error al volver a obtener el bloque 1 del archivo users.txt", fmt.Errorf("error al volver a obtener el bloque 1 del archivo users.txt: %v", err)
			}
			updatedUsers := strings.Split(strings.TrimSpace(string(fileBlock.B_content[:])), "\n")
			fmt.Println("Contenido actualizado de users.txt:", updatedUsers)

			fmt.Println("Grupo creado exitosamente")
			return "Grupo creado exitosamente", nil
		}
	}

	return "Error en la estructura del sistema de archivos", nil
}
