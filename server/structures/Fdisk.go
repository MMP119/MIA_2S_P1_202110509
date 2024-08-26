package structures

import (
	"encoding/binary"
	"fmt"
	util "server/util"
)

// FDISK estructura que representa el comando fdisk con sus parámetros
type FDISK struct {
	Size int    // Tamaño de la partición
	Unit string // Unidad de medida del tamaño (B, K o M); por defecto K
	Path string // Ruta del archivo del disco
	TypE  string // Tipo de partición (P, E, L)
	Fit  string // Tipo de ajuste (BF, FF, WF); por defecto WF
	Name string // Nombre de la partición
}

func CommandFdisk(fdisk *FDISK) (string, error) {
	
	// Convertir el tamaño a bytes
	sizeBytes, err := util.ConvertToBytes(fdisk.Size, fdisk.Unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return "Error converting size en Fdisk", err
	}

	var msg string

	// para crear la particion primaria
	if(fdisk.TypE == "P"){
		msg, err = CreatePrimaryPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creating primary partition:", err)
			return msg, err
		}
	}else if(fdisk.TypE == "E"){

		msg, err = CreatePrimaryPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creating extended partition:", err)
			return msg, err
		}


	}else if(fdisk.TypE == "L"){

		msg, err = CreateLogicalPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creating logical partition:", err)
			return msg, err
		}

	}
	return "",nil
}


func CreatePrimaryPartition(fdisk *FDISK, sizeBytes int)(string, error){
	
	var mbr MBR
	
	msg, err := mbr.DeserializeMBR(fdisk.Path)
	if err != nil {
		return msg, fmt.Errorf("error leyendo el MBR del disco: %s", err)
	}

	// Contar el número de particiones primarias y extendidas
	primaryCount := 0
	extendedExists := 0
	for _, partition := range mbr.Mbr_partitions {
		if partition.Part_status[0] != '2' {
			if partition.Part_type[0] == 'P' {
				primaryCount++
			} else if partition.Part_type[0] == 'E' {
				extendedExists++
			}
		}
	}

	// Verificar que no se exceda el límite de 4 particiones
	if primaryCount >4 {
		return "ERROR: No se pueden crear más particiones primarias", fmt.Errorf("límite de particiones primarias alcanzado")
	}

	if extendedExists >=2 {
		return "ERROR: No se pueden crear más particiones extendidas, ya existe una en el disco", fmt.Errorf("ya existe una partición extendida")
	}

	// se obtiene la primera particion libre
	particionDisponible, inicioParticion, indexParticion, msg:= mbr.GetFirstPartitionAvailable()
	if particionDisponible == nil {
		return msg, fmt.Errorf("no hay particiones disponibles")
	}

	// crear la particion con los parámetros proporcionados 
	particionDisponible.CreatePartition(inicioParticion, sizeBytes, fdisk.TypE, fdisk.Fit, fdisk.Name)

	// montar la particion
	mbr.Mbr_partitions[indexParticion] = *particionDisponible //asignar la particion al MBR

	// Serialiazar el MBR modificado
	msg, err = mbr.SerializeMBR(fdisk.Path)
	if err != nil {
		return msg, fmt.Errorf("error escribiendo el MBR al disco: %s", err)
	}	
	return "",nil
}


func CreateLogicalPartition(fdisk *FDISK, sizeBytes int) (string, error) {
    var mbr MBR
    // Deserializar el MBR del disco
    msg, err := mbr.DeserializeMBR(fdisk.Path)
    if err != nil {
        return msg, fmt.Errorf("error leyendo el MBR del disco: %s", err)
    }

    // Verificar que exista una partición extendida
    var extendedPartition *PARTITION
    for _, partition := range mbr.Mbr_partitions {
        if partition.Part_type[0] == 'E' {
            extendedPartition = &partition
            break
        }
    }
    if extendedPartition == nil {
        return "ERROR: no existe una partición extendida", fmt.Errorf("no existe una partición extendida")
    }

    // Intentar deserializar el primer EBR en la partición extendida
    var ebr EBR
    ebrPosition := extendedPartition.Part_start
    _, err = ebr.DeserializeEBR(fdisk.Path, ebrPosition)

	if err != nil { // Si no existe un EBR, crear el primero
		fmt.Println("Entra acá, creando primer EBR")
		msg, err = CreateEBR(fdisk.Path, int32(sizeBytes), fdisk, ebrPosition)
		if err != nil {
			return msg, fmt.Errorf("error creando el primer EBR: %s", err)
		}

		// Crear la partición lógica justo después del primer EBR
		logicalPartitionStart := ebrPosition + int32(binary.Size(ebr)) // EBR ocupa un espacio
		var logicalPartition PARTITION
		logicalPartition.CreatePartition(int(logicalPartitionStart), sizeBytes, "L", fdisk.Fit, fdisk.Name)

		// Serializar el primer EBR
		_, err = ebr.SerializeEBR(fdisk.Path, ebrPosition)
		if err != nil {
			return "Error serializando el primer EBR", fmt.Errorf("error escribiendo el EBR al disco: %s", err)
		}

	} else { // Existe al menos un EBR, buscar el último y crear el nuevo
		fmt.Println("Creando nuevo EBR, no es el primero")

		// Recorrer hasta el último EBR
		for ebr.Part_next != -1 {
			ebrPosition = ebr.Part_next
			_, err = ebr.DeserializeEBR(fdisk.Path, ebrPosition)
			if err != nil {
				return "Error deserializando EBR existente", fmt.Errorf("error deserializando: %s", err)
			}
		}

		// Crear un nuevo EBR después del último EBR y su partición lógica
		newEBRPosition := ebrPosition + int32(binary.Size(ebr)) + ebr.Part_size
		msg, err = CreateEBR(fdisk.Path, int32(sizeBytes), fdisk, newEBRPosition)
		if err != nil {
			return msg, fmt.Errorf("error creando el nuevo EBR: %s", err)
		}

		// Actualizar el campo Part_next del EBR anterior
		ebr.Part_next = newEBRPosition
		_, err = ebr.SerializeEBR(fdisk.Path, ebrPosition)
		if err != nil {
			return "Error actualizando EBR anterior", fmt.Errorf("error escribiendo el EBR anterior al disco: %s", err)
		}
	}

    return "EBR creado exitosamente", nil
}
