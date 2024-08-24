package global

import (
	structures "server/structures"
	"errors"
)


const Carnet string = "09" //202110509

var (
	MountedPartitions map[string]string = make(map[string]string)
)

// GetMountedPartition obtiene la partición montada con el id especificado
func GetMountedPartition(id string) (*structures.PARTITION, string, error) {
	// Obtener el path de la partición montada
	path := MountedPartitions[id]
	if path == "" {
		return nil, "", errors.New("la partición no está montada")
	}

	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	msg,err := mbr.DeserializeMBR(path)
	if err != nil {
		return nil,msg, err
	}

	// Buscar la partición con el id especificado
	partition, err:= mbr.GetPartitionByID(id)
	if partition == nil {
		return nil, "", err
	}

	return partition, path, nil
}