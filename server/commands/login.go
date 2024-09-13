package commands

import (
	"errors"
	"fmt"
	"regexp"
	global "server/global"
	structures "server/structures"
	"strings"
)

type LOGIN struct {
	User string
	Pass string
	Id 	 string
}


func ParseLogin(tokens []string)(*LOGIN, string, error){
	cmd := &LOGIN{}

	args := strings.Join(tokens, " ")

	re:= regexp.MustCompile(`(?i)-user=[^\s]+|(?i)-pass=[^\s]+|(?i)-id=[^\s]+`)

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

		switch key{
			case "-user":
				if value == "" {
					return nil, "ERROR: el usuario no puede estar vacío", errors.New("el usuario no puede estar vacío")
				}
				cmd.User = value

			case "-pass":
				if value == "" {
					return nil, "ERROR: la contraseña no puede estar vacía", errors.New("la contraseña no puede estar vacía")
				}
				cmd.Pass = value

			case "-id":
				if value == "" {
					return nil, "ERROR: el id no puede estar vacío", errors.New("el id no puede estar vacío")
				}
				cmd.Id = value
			
			default:
				return nil, "ERROR: parámetro desconocido", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	if cmd.User == "" || cmd.Pass == "" || cmd.Id == "" {
		return nil, "ERROR: faltan parámetros requeridos: -user, -pass, -id", errors.New("faltan parámetros requeridos: -user, -pass, -id")
	}

	msg, err:= CommandLogin(cmd)
	if err != nil {
		return nil, msg, err
	}

	return cmd, "Comando LOGIN: realizado correctamente: "+msg, nil

}


func CommandLogin(login *LOGIN) (string, error){

	// ir al archivo users.txt y buscar el usuario y la contraseña

	//obtener la particion con el id en donde se realizará el login
	partitionSuperblock, partition, partitionPath, err := global.GetMountedPartitionSuperblock(login.Id)
	if err != nil {
		return "Error al obtener la partición montada en el comando login", fmt.Errorf("error al obtener la partición montada: %v", err)
	}


	inode := &structures.Inode{}

	err = inode.Deserialize(partitionPath, int64(partitionSuperblock.S_inode_start+(0*partitionSuperblock.S_inode_size)))
	if err != nil {
		return "error al obtener el inodo raiz",fmt.Errorf("error al obtener el inodo raiz: %v", err)
	}


	//verificar que el primer i nodo esté en cero
	if (inode.I_block[0] == 0){
		//moverme al bloque 0
		folderBlock := &structures.FolderBlock{}

		err = folderBlock.Deserialize(partitionPath, int64(partitionSuperblock.S_block_start+(inode.I_block[0]*partitionSuperblock.S_block_size)))
		if err != nil {
			return "Error al obtener el bloque 0",fmt.Errorf("error al obtener el bloque 0: %v", err)
		}

		//recorrer los contenidos del bloque 0
		for _, contenido := range folderBlock.B_content{
			name := strings.Trim(string(contenido.B_name[:]), "\x00") // Elimina caracteres nulos
			apuntador := contenido.B_inodo
			if (name == "users.txt"){

				//moverme al inodo que apunta el contenido
				err = inode.Deserialize(partitionPath, int64(partitionSuperblock.S_inode_start+(apuntador*partitionSuperblock.S_inode_size)))
				if err != nil {
					return "Error al obtener el inodo del archivo users.txt",fmt.Errorf("error al obtener el inodo del archivo users.txt: %v", err)
				}

				//verificar que el primer i nodo esté en 1
				if (inode.I_block[0] == 1){
					//moverme al bloque 1
					fileBlock := &structures.FileBlock{}

					err = fileBlock.Deserialize(partitionPath, int64(partitionSuperblock.S_block_start+(inode.I_block[0]*partitionSuperblock.S_block_size)))
					if err != nil {
						return "error al obtener el bloque 1 del archivo users.txt",fmt.Errorf("error al obtener el bloque 1 del archivo users.txt: %v", err)
					}


					/*
						el file block tiene 64 bytes en donde se guarda lo siguiente
						1,G,root\n1,U,root,root,123\n

						1,G,root
						1,U,root,root,123

						donde:
						GUI, TIPO, GRUPO
						UID, TIPO, GRUPO, USUARIO, CONTRASEÑA
					*/

					// obtener el usuario y la contraseña
					users := strings.Split(strings.TrimSpace(string(fileBlock.B_content[:])), "\n")
					for _, user := range users {
						if user != "" {
							// separar los valores
							values := strings.Split(user, ",")
							
							// Verificar si es una línea de usuario y no de grupo
							if len(values) >= 5 {
								//fmt.Println(values) // para depuración

								// verificar si el usuario y la contraseña son correctos
								if values[3] == login.User && values[4] == login.Pass {
									
									//verificar si ya hay una sesión activa
									if global.IsSessionActive() {
										//fmt.Println("YA HAY UNA SESION ACTIVA")
										return "YA HAY UNA SESION ACTIVA, DEBE HACER LOGOUT!", nil

									} else {
										//fmt.Println("USUARIO Y CONTRASEÑA CORRECTOS, SESION ACTIVA")	
										global.ActivateSession()
										return "USUARIO Y CONTRASEÑA CORRECTOS, SESION ACTIVA", nil

									}
								} else {
									//fmt.Println("USUARIO Y CONTRASEÑA INCORRECTOS")
									return "USUARIO Y CONTRASEÑA INCORRECTOS", nil
								}
							}
						}
					}

				}

			}
		}

		
	}


	err = partitionSuperblock.Serialize(partitionPath, int64(partition.Part_start))
	if err != nil {
		return "error al serializar el superbloque de la partición",fmt.Errorf("error al serializar el superbloque de la partición: %v", err)
	}

	return "",nil
}