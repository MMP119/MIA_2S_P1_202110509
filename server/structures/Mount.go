package structures

import (
	"fmt"
	//global "server/global"
	util "server/util"
)

type MOUNT struct {
	Path string 
	Name string 
}

func CommandMount(mount *MOUNT) (string, error) {
	
	var mbr MBR

	msg, err := mbr.DeserializeMBR(mount.Path)
	if err != nil {
		return msg, fmt.Errorf("error leyendo el MBR del disco: %s", err)
	}

	// buscar la particion con el nombre proporcionado
	partition, indexPartition, msg := mbr.GetPartitionByName(mount.Name, mount.Path)
	if partition == nil {
		return msg, fmt.Errorf("no se encontró la partición con el nombre: %s", mount.Name)
	}

	// se genera un id único para la partición
	id, msg, err := GenerateIdPartition(mount, indexPartition)
	if err != nil {
		return msg, fmt.Errorf("error generando id de partición: %s", err)
	}

	//guardar la particion montada en la lista de montajes globales
	//util.GlobalMounts[id] = mount.Path
	//global.MountedPartitions[id] = mount.Path

	// modificar la particion para indicar que está montada
	partition.MountPartition(indexPartition, id)

	// guardar la particion mod en el mbr
	mbr.Mbr_partitions[indexPartition] = *partition

	// serializar el mbr
	msg, err = mbr.SerializeMBR(mount.Path)
	if err != nil {
		return msg, fmt.Errorf("error escribiendo el MBR en el disco: %s", err)
	}
	
	return "", nil
}



func GenerateIdPartition(mount *MOUNT, indexPartition int) (string, string, error) {
	// Asignar una letra a la partición
	letter, err := util.GetLetter(mount.Path)
	if err != nil {
		fmt.Println("Error obteniendo la letra:", err)
		return "", "Error obteniendo la letra en mount",err
	}

	// Crear id de partición
	idPartition := fmt.Sprintf("%s%d%s", util.Carnet, indexPartition+1, letter)

	return idPartition, "", nil
}